// credentailfetcher fetches workload credentials through platform plugins.
package credentialfetcher

import (
	"fmt"

	"istio.io/istio/pkg/security"
	"istio.io/istio/security/pkg/credentialfetcher/plugin"
)

func NewCredFetcher(credtype, trustdomain, jwtPath string) (security.CredFetcher, error) {
	switch credtype {
	case security.GCE:
		return plugin.CreateGCEPlugin(trustdomain, jwtPath), nil
	case security.Mock: // for test only
		return plugin.CreateMockPlugin(), nil
	default:
		return nil, fmt.Errorf("invalid credential fetcher type %s", credtype)
	}
}
