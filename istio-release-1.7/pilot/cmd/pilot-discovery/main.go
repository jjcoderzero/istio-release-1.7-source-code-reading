package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"google.golang.org/grpc/grpclog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/pkg/collateral"
	"istio.io/pkg/ctrlz"
	"istio.io/pkg/log"
	"istio.io/pkg/version"

	"istio.io/istio/pilot/pkg/bootstrap"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/serviceregistry"
	"istio.io/istio/pkg/cmd"
	"istio.io/istio/pkg/config/constants"
	"istio.io/istio/pkg/spiffe"
)

const (
	defaultMCPMaxMessageSize        = 1024 * 1024 * 4  // 默认gRPC最大消息大小
	defaultMCPInitialConnWindowSize = 1024 * 1024      // 默认gRPC初始化窗口大小
	defaultMCPInitialWindowSize     = 1024 * 1024      // 默认gRPC连接窗口大小
)

var (
	serverArgs *bootstrap.PilotArgs

	loggingOptions = log.DefaultOptions()

	rootCmd = &cobra.Command{
		Use:          "pilot-discovery",
		Short:        "Istio Pilot.",
		Long:         "Istio Pilot provides fleet-wide traffic management capabilities in the Istio Service Mesh.",
		SilenceUsage: true,
	}

	discoveryCmd = &cobra.Command{
		Use:   "discovery",
		Short: "Start Istio proxy discovery service.",
		Args:  cobra.ExactArgs(0),
		RunE: func(c *cobra.Command, args []string) error {
			cmd.PrintFlags(c.Flags())
			if err := log.Configure(loggingOptions); err != nil {
				return err
			}
			grpclog.SetLoggerV2(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, ioutil.Discard))

			spiffe.SetTrustDomain(spiffe.DetermineTrustDomain(serverArgs.RegistryOptions.KubeOptions.TrustDomain, hasKubeRegistry()))

			// 创建一个接收空结构的stop channel用来停止所有servers
			stop := make(chan struct{})

			// 创建服务发现的 Server
			discoveryServer, err := bootstrap.NewServer(serverArgs)
			if err != nil {
				return fmt.Errorf("failed to create discovery service: %v", err)
			}

			// 运行 Server 中注册的所有服务
			if err := discoveryServer.Start(stop); err != nil {
				return fmt.Errorf("failed to start discovery service: %v", err)
			}
			// 等待 SIGINT 和 SIGTERM 信号并关闭 stop channel
			cmd.WaitSignal(stop)
			// Wait until we shut down. In theory this could block forever; in practice we will get
			// forcibly shut down after 30s in Kubernetes.
			discoveryServer.WaitUntilCompletion()
			return nil
		},
	}
)

// 当我们运行在k8s上，默认的trust domain 是 'cluster.local', 否则是空的字符串
func hasKubeRegistry() bool {
	for _, r := range serverArgs.RegistryOptions.Registries {
		if serviceregistry.ProviderID(r) == serviceregistry.Kubernetes {
			return true
		}
	}
	return false
}

func init() {

	serverArgs = bootstrap.NewPilotArgs(func(p *bootstrap.PilotArgs) {
		// Set Defaults
		p.CtrlZOptions = ctrlz.DefaultOptions()
		// TODO replace with mesh config?
		p.InjectionOptions = bootstrap.InjectionOptions{
			InjectionDirectory: "./var/lib/istio/inject",
		}
	})

	// Process commandline args.
	discoveryCmd.PersistentFlags().StringSliceVar(&serverArgs.RegistryOptions.Registries, "registries",
		[]string{string(serviceregistry.Kubernetes)},
		fmt.Sprintf("Comma separated list of platform service registries to read from (choose one or more from {%s, %s, %s})",
			serviceregistry.Kubernetes, serviceregistry.Consul, serviceregistry.Mock))
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.RegistryOptions.ClusterRegistriesNamespace, "clusterRegistriesNamespace",
		serverArgs.RegistryOptions.ClusterRegistriesNamespace, "Namespace for ConfigMap which stores clusters configs")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.RegistryOptions.KubeConfig, "kubeconfig", "",
		"Use a Kubernetes configuration file instead of in-cluster configuration")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.MeshConfigFile, "meshConfig", "./etc/istio/config/mesh",
		"File name for Istio mesh configuration. If not specified, a default mesh will be used.")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.NetworksConfigFile, "networksConfig", "/etc/istio/config/meshNetworks",
		"File name for Istio mesh networks configuration. If not specified, a default mesh networks will be used.")
	discoveryCmd.PersistentFlags().StringVarP(&serverArgs.Namespace, "namespace", "n", bootstrap.PodNamespaceVar.Get(),
		"Select a namespace where the controller resides. If not set, uses ${POD_NAMESPACE} environment variable")
	discoveryCmd.PersistentFlags().StringSliceVar(&serverArgs.Plugins, "plugins", bootstrap.DefaultPlugins,
		"comma separated list of networking plugins to enable")

	// MCP client flags
	discoveryCmd.PersistentFlags().IntVar(&serverArgs.MCPOptions.MaxMessageSize, "mcpMaxMsgSize", defaultMCPMaxMessageSize,
		"Max message size received by MCP's gRPC client")
	discoveryCmd.PersistentFlags().IntVar(&serverArgs.MCPOptions.InitialWindowSize, "mcpInitialWindowSize", defaultMCPInitialWindowSize,
		"Initial window size for MCP's gRPC connection")
	discoveryCmd.PersistentFlags().IntVar(&serverArgs.MCPOptions.InitialConnWindowSize, "mcpInitialConnWindowSize", defaultMCPInitialConnWindowSize,
		"Initial connection window size for MCP's gRPC connection")

	// RegistryOptions Controller options
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.RegistryOptions.FileDir, "configDir", "",
		"Directory to watch for updates to config yaml files. If specified, the files will be used as the source of config, rather than a CRD client.")
	discoveryCmd.PersistentFlags().StringVarP(&serverArgs.RegistryOptions.KubeOptions.WatchedNamespaces, "appNamespace", "a", metav1.NamespaceAll,
		"Specify the applications namespace list the controller manages, separated by comma; if not set, controller watches all namespaces")
	discoveryCmd.PersistentFlags().DurationVar(&serverArgs.RegistryOptions.KubeOptions.ResyncPeriod, "resync", 60*time.Second,
		"Controller resync interval")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.RegistryOptions.KubeOptions.DomainSuffix, "domain", constants.DefaultKubernetesDomain,
		"DNS domain suffix")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.RegistryOptions.KubeOptions.ClusterID, "clusterID", features.ClusterName,
		"The ID of the cluster that this Istiod instance resides")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.RegistryOptions.KubeOptions.TrustDomain, "trust-domain", "",
		"The domain serves to identify the system with spiffe")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.RegistryOptions.ConsulServerAddr, "consulserverURL", "",
		"URL for the Consul server")

	// using address, so it can be configured as localhost:.. (possibly UDS in future)
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.ServerOptions.HTTPAddr, "httpAddr", ":8080",
		"Discovery service HTTP address")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.ServerOptions.HTTPSAddr, "httpsAddr", ":15017",
		"Injection and validation service HTTPS address")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.ServerOptions.GRPCAddr, "grpcAddr", ":15010",
		"Discovery service gRPC address")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.ServerOptions.SecureGRPCAddr, "secureGRPCAddr", ":15012",
		"Discovery service secured gRPC address")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.ServerOptions.MonitoringAddr, "monitoringAddr", ":15014",
		"HTTP address to use for pilot's self-monitoring information")
	discoveryCmd.PersistentFlags().BoolVar(&serverArgs.ServerOptions.EnableProfiling, "profile", true,
		"Enable profiling via web interface host:port/debug/pprof")

	// Use TLS certificates if provided.
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.ServerOptions.TLSOptions.CaCertFile, "caCertFile", "",
		"File containing the x509 Server CA Certificate")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.ServerOptions.TLSOptions.CertFile, "tlsCertFile", "",
		"File containing the x509 Server Certificate")
	discoveryCmd.PersistentFlags().StringVar(&serverArgs.ServerOptions.TLSOptions.KeyFile, "tlsKeyFile", "",
		"File containing the x509 private key matching --tlsCertFile")

	// Attach the Istio logging options to the command.
	loggingOptions.AttachCobraFlags(rootCmd)

	// Attach the Istio Ctrlz options to the command.
	serverArgs.CtrlZOptions.AttachCobraFlags(rootCmd)

	// Attach the Istio Keepalive options to the command.
	serverArgs.KeepaliveOptions.AttachCobraFlags(rootCmd)

	cmd.AddFlags(rootCmd)

	rootCmd.AddCommand(discoveryCmd)
	rootCmd.AddCommand(version.CobraCommand())
	rootCmd.AddCommand(collateral.CobraCommand(rootCmd, &doc.GenManHeader{
		Title:   "Istio Pilot Discovery",
		Section: "pilot-discovery CLI",
		Manual:  "Istio Pilot Discovery",
	}))

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errora(err)
		os.Exit(-1)
	}
}
