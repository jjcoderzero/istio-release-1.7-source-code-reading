package util

import (
	"regexp"

	"istio.io/istio/pkg/config/constants"
)

const (
	DefaultKubernetesDomain   = "svc." + constants.DefaultKubernetesDomain
	ExportToNamespaceLocal    = "."
	ExportToAllNamespaces     = "*"
	IstioProxyName            = "istio-proxy"
	MeshGateway               = "mesh"
	Wildcard                  = "*"
	MeshConfigName            = "istio"
	InjectionLabelName        = "istio-injection"
	InjectionLabelEnableValue = "enabled"
)

var (
	fqdnPattern = regexp.MustCompile(`^(.+)\.(.+)\.svc\.cluster\.local$`)
)
