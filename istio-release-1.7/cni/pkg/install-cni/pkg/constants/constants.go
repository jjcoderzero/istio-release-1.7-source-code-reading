package constants

// 命令行参数
const (
	MountedCNINetDir     = "mounted-cni-net-dir"
	CNINetDir            = "cni-net-dir"
	CNIConfName          = "cni-conf-name"
	ChainedCNIPlugin     = "chained-cni-plugin"
	Sleep                = "sleep"
	CNINetworkConfigFile = "cni-network-config-file"
	CNINetworkConfig     = "cni-network-config"
	LogLevel             = "log-level"
	KubeconfigFilename   = "kubecfg-file-name"
	KubeconfigMode       = "kubeconfig-mode"
	KubeCAFile           = "kube-ca-file"
	SkipTLSVerify        = "skip-tls-verify"
	SkipCNIBinaries      = "skip-cni-binaries"
	UpdateCNIBinaries    = "update-cni-binaries"
)

// 内部的常量
const (
	CNIBinDir             = "/opt/cni/bin"
	HostCNIBinDir         = "/host/opt/cni/bin"
	SecondaryBinDir       = "/host/secondary-bin-dir"
	ServiceAccountPath    = "/var/run/secrets/kubernetes.io/serviceaccount"
	DefaultKubeconfigMode = 0600
)
