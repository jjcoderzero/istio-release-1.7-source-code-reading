package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"istio.io/api/annotation"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
	"go.uber.org/zap"

	"istio.io/pkg/log"
)

var (
	nsSetupBinDir          = "/opt/cni/bin"
	injectAnnotationKey    = annotation.SidecarInject.Name
	sidecarStatusKey       = annotation.SidecarStatus.Name
	interceptRuleMgrType   = defInterceptRuleMgrType
	loggingOptions         = log.DefaultOptions()
	podRetrievalMaxRetries = 30
	podRetrievalInterval   = 1 * time.Second
)

const ISTIOINIT = "istio-init"

// Kubernetes一个K8s特定的结构来保存配置
type Kubernetes struct {
	K8sAPIRoot           string   `json:"k8s_api_root"`
	Kubeconfig           string   `json:"kubeconfig"`
	InterceptRuleMgrType string   `json:"intercept_type"`
	NodeName             string   `json:"node_name"`
	ExcludeNamespaces    []string `json:"exclude_namespaces"`
	CNIBinDir            string   `json:"cni_bin_dir"`
}

// PluginConf是您所期望的配置json。这是通过stdin传入的内容。你的插件可能希望通过运行时参数来公开它的功能
type PluginConf struct {
	types.NetConf // 您可能不希望嵌套此类型
	RuntimeConfig *struct {
		// SampleConfig map[string]interface{} `json:"sample"`
	} `json:"runtimeConfig"`

	RawPrevResult *map[string]interface{} `json:"prevResult"`
	PrevResult    *current.Result         `json:"-"`

	// Add plugin-specific flags here
	LogLevel   string     `json:"log_level"`
	Kubernetes Kubernetes `json:"kubernetes"`
}

// K8sArgs是Kubernetes使用的有效CNI_ARGS，字段名需要与kubelet args中的精确键匹配才能进行数据编组
type K8sArgs struct {
	types.CommonArgs
	IP                         net.IP
	K8S_POD_NAME               types.UnmarshallableString // nolint: golint, stylecheck
	K8S_POD_NAMESPACE          types.UnmarshallableString // nolint: golint, stylecheck
	K8S_POD_INFRA_CONTAINER_ID types.UnmarshallableString // nolint: golint, stylecheck
}

// parseConfig从stdin中解析提供的配置(和prevResult)
func parseConfig(stdin []byte) (*PluginConf, error) {
	conf := PluginConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}

	// 解析之前的结果。如果你的插件没有链接，删除它.
	if conf.RawPrevResult != nil {
		resultBytes, err := json.Marshal(conf.RawPrevResult)
		if err != nil {
			return nil, fmt.Errorf("could not serialize prevResult: %v", err)
		}
		res, err := version.NewResult(conf.CNIVersion, resultBytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse prevResult: %v", err)
		}
		conf.RawPrevResult = nil
		conf.PrevResult, err = current.NewResultFromResult(res)
		if err != nil {
			return nil, fmt.Errorf("could not convert result to current version: %v", err)
		}
	}
	// End previous result parsing

	return &conf, nil
}

// cmdAdd is called for ADD requests
func cmdAdd(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		log.Errorf("istio-cni cmdAdd parsing config %v", err)
		return err
	}

	var loggedPrevResult interface{}
	if conf.PrevResult == nil {
		loggedPrevResult = "none"
	} else {
		loggedPrevResult = conf.PrevResult
	}

	log.Info("CmdAdd config parsed",
		zap.String("version", conf.CNIVersion),
		zap.Reflect("prevResult", loggedPrevResult))

	// Determine if running under k8s by checking the CNI args
	k8sArgs := K8sArgs{}
	if err := types.LoadArgs(args.Args, &k8sArgs); err != nil {
		return err
	}
	log.Infof("Getting identifiers with arguments: %s", args.Args)
	log.Infof("Loaded k8s arguments: %v", k8sArgs)
	if conf.Kubernetes.CNIBinDir != "" {
		nsSetupBinDir = conf.Kubernetes.CNIBinDir
	}
	if conf.Kubernetes.InterceptRuleMgrType != "" {
		interceptRuleMgrType = conf.Kubernetes.InterceptRuleMgrType
	}

	log.Info("",
		zap.String("ContainerID", args.ContainerID),
		zap.String("Pod", string(k8sArgs.K8S_POD_NAME)),
		zap.String("Namespace", string(k8sArgs.K8S_POD_NAMESPACE)),
		zap.String("InterceptType", interceptRuleMgrType))

	// Check if the workload is running under Kubernetes.
	if string(k8sArgs.K8S_POD_NAMESPACE) != "" && string(k8sArgs.K8S_POD_NAME) != "" {
		excludePod := false
		for _, excludeNs := range conf.Kubernetes.ExcludeNamespaces {
			if string(k8sArgs.K8S_POD_NAMESPACE) == excludeNs {
				excludePod = true
				break
			}
		}
		if !excludePod {
			client, err := newKubeClient(*conf)
			if err != nil {
				return err
			}
			log.Debug("Created Kubernetes client", zap.Reflect("client", client))
			var containers []string
			var initContainersMap map[string]struct{}
			var annotations map[string]string
			var k8sErr error
			for attempt := 1; attempt <= podRetrievalMaxRetries; attempt++ {
				containers, initContainersMap, _, annotations, k8sErr = getKubePodInfo(client, string(k8sArgs.K8S_POD_NAME), string(k8sArgs.K8S_POD_NAMESPACE))
				if k8sErr == nil {
					break
				}
				log.Warn("Waiting for pod metadata", zap.Error(k8sErr), zap.Int("attempt", attempt))
				time.Sleep(podRetrievalInterval)
			}
			if k8sErr != nil {
				log.Error("Failed to get pod data", zap.Error(k8sErr))
				return k8sErr
			}

			// Check if istio-init container is present; in that case exclude pod
			if _, present := initContainersMap[ISTIOINIT]; present {
				log.Info("Pod excluded due to being already injected with istio-init container",
					zap.String("pod", string(k8sArgs.K8S_POD_NAME)),
					zap.String("namespace", string(k8sArgs.K8S_POD_NAMESPACE)))
				excludePod = true
			}

			log.Infof("Found containers %v", containers)
			if len(containers) > 1 {
				log.Info("Checking annotations prior to redirect for Istio proxy",
					zap.String("ContainerID", args.ContainerID),
					zap.String("netns", args.Netns),
					zap.String("pod", string(k8sArgs.K8S_POD_NAME)),
					zap.String("Namespace", string(k8sArgs.K8S_POD_NAMESPACE)),
					zap.Reflect("annotations", annotations))
				if val, ok := annotations[injectAnnotationKey]; ok {
					log.Infof("Pod %s contains inject annotation: %s", string(k8sArgs.K8S_POD_NAME), val)
					if injectEnabled, err := strconv.ParseBool(val); err == nil {
						if !injectEnabled {
							log.Infof("Pod excluded due to inject-disabled annotation")
							excludePod = true
						}
					}
				}
				if _, ok := annotations[sidecarStatusKey]; !ok {
					log.Infof("Pod %s excluded due to not containing sidecar annotation", string(k8sArgs.K8S_POD_NAME))
					excludePod = true
				}
				if !excludePod {
					log.Infof("setting up redirect")
					if redirect, redirErr := NewRedirect(annotations); redirErr != nil {
						log.Errorf("Pod redirect failed due to bad params: %v", redirErr)
					} else {
						log.Infof("Redirect local ports: %v", redirect.includePorts)
						// Get the constructor for the configured type of InterceptRuleMgr
						interceptMgrCtor := GetInterceptRuleMgrCtor(interceptRuleMgrType)
						if interceptMgrCtor == nil {
							log.Errorf("Pod redirect failed due to unavailable InterceptRuleMgr of type %s",
								interceptRuleMgrType)
						} else {
							rulesMgr := interceptMgrCtor()
							if err := rulesMgr.Program(args.Netns, redirect); err != nil {
								return err
							}
						}
					}
				}
			}
		} else {
			log.Infof("Pod excluded")
		}
	} else {
		log.Infof("No Kubernetes Data")
	}

	var result *current.Result
	if conf.PrevResult == nil {
		result = &current.Result{
			CNIVersion: current.ImplementedSpecVersion,
		}
	} else {
		// Pass through the result for the next plugin
		result = conf.PrevResult
	}
	return types.PrintResult(result, conf.CNIVersion)
}

func cmdGet(args *skel.CmdArgs) error {
	log.Info("cmdGet not implemented")
	// TODO: implement
	return fmt.Errorf("not implemented")
}

// cmdDel is called for DELETE requests
func cmdDel(args *skel.CmdArgs) error {
	log.Info("istio-cni cmdDel parsing config")
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}
	_ = conf

	// Do your delete here

	return nil
}

func main() {
	loggingOptions.OutputPaths = []string{"stderr"}
	loggingOptions.JSONEncoding = true
	if err := log.Configure(loggingOptions); err != nil {
		os.Exit(1)
	}
	// TODO: implement plugin version
	skel.PluginMain(cmdAdd, cmdGet, cmdDel, version.All, "istio-cni")
}
