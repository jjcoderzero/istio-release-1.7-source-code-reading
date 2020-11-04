package service

import (
	"istio.io/istio/galley/pkg/config/analysis"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/util"
	"istio.io/istio/galley/pkg/config/analysis/msg"
	configKube "istio.io/istio/pkg/config/kube"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"

	v1 "k8s.io/api/core/v1"
)

// PortNameAnalyzer checks the port name of the service
type PortNameAnalyzer struct{}

var _ analysis.Analyzer = &PortNameAnalyzer{}

// Metadata implements Analyzer
func (s *PortNameAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "service.PortNameAnalyzer",
		Description: "Checks the port names associated with each service",
		Inputs: collection.Names{
			collections.K8SCoreV1Services.Name(),
		},
	}
}

// Analyze implements Analyzer
func (s *PortNameAnalyzer) Analyze(c analysis.Context) {
	c.ForEach(collections.K8SCoreV1Services.Name(), func(r *resource.Instance) bool {
		svcNs := r.Metadata.FullName.Namespace

		// Skip system namespaces entirely
		if util.IsSystemNamespace(svcNs) {
			return true
		}

		// Skip port name check for istio control plane
		if util.IsIstioControlPlane(r) {
			return true
		}

		s.analyzeService(r, c)
		return true
	})
}

func (s *PortNameAnalyzer) analyzeService(r *resource.Instance, c analysis.Context) {
	svc := r.Message.(*v1.ServiceSpec)
	for _, port := range svc.Ports {
		if instance := configKube.ConvertProtocol(port.Port, port.Name, port.Protocol, port.AppProtocol); instance.IsUnsupported() {
			c.Report(collections.K8SCoreV1Services.Name(), msg.NewPortNameIsNotUnderNamingConvention(
				r, port.Name, int(port.Port), port.TargetPort.String()))
		}
	}
}
