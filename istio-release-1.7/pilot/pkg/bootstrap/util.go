package bootstrap

import (
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry"
	"istio.io/pkg/ledger"
)

func hasKubeRegistry(registries []string) bool {
	for _, r := range registries {
		if serviceregistry.ProviderID(r) == serviceregistry.Kubernetes {
			return true
		}
	}
	return false
}

func buildLedger(ca RegistryOptions) ledger.Ledger {
	var result ledger.Ledger
	if ca.DistributionTrackingEnabled {
		result = ledger.Make(ca.DistributionCacheRetention)
	} else {
		result = &model.DisabledLedger{}
	}
	return result
}
