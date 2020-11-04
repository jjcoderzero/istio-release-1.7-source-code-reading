package analyzers

import (
	"istio.io/istio/galley/pkg/config/analysis"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/annotations"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/authz"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/deployment"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/deprecation"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/destinationrule"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/gateway"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/injection"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/multicluster"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/schema"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/service"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/sidecar"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/virtualservice"
)

// All returns all analyzers
func All() []analysis.Analyzer {
	analyzers := []analysis.Analyzer{
		// Please keep this list sorted alphabetically by pkg.name for convenience
		&annotations.K8sAnalyzer{},
		&authz.AuthorizationPoliciesAnalyzer{},
		&deployment.ServiceAssociationAnalyzer{},
		&deprecation.FieldAnalyzer{},
		&gateway.IngressGatewayPortAnalyzer{},
		&gateway.SecretAnalyzer{},
		&injection.Analyzer{},
		&injection.ImageAnalyzer{},
		&multicluster.MeshNetworksAnalyzer{},
		&service.PortNameAnalyzer{},
		&sidecar.DefaultSelectorAnalyzer{},
		&sidecar.SelectorAnalyzer{},
		&virtualservice.ConflictingMeshGatewayHostsAnalyzer{},
		&virtualservice.DestinationHostAnalyzer{},
		&virtualservice.DestinationRuleAnalyzer{},
		&virtualservice.GatewayAnalyzer{},
		&virtualservice.RegexAnalyzer{},
		&destinationrule.CaCertificateAnalyzer{},
	}

	analyzers = append(analyzers, schema.AllValidationAnalyzers()...)

	return analyzers
}

// AllCombined returns all analyzers combined as one
func AllCombined() *analysis.CombinedAnalyzer {
	return analysis.Combine("all", All()...)
}
