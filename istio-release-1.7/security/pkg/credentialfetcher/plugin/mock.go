// Test only: this is the mock plugin of credentialfetcher.
package plugin

import (
	"istio.io/pkg/log"

	"istio.io/istio/pkg/security"
)

var (
	mockcredLog = log.RegisterScope("mockcred", "Mock credential fetcher for istio agent", 0)
)

// The plugin object.
type MockPlugin struct {
}

// CreateMockPlugin creates a mock credential fetcher plugin. Return the pointer to the created plugin.
func CreateMockPlugin() *MockPlugin {
	p := &MockPlugin{}
	return p
}

// GetPlatformCredential returns a constant token string.
func (p *MockPlugin) GetPlatformCredential() (string, error) {
	mockcredLog.Debugf("mock plugin returns a constant token.")
	return "test_token", nil
}

// GetType returns credential fetcher type.
func (p *MockPlugin) GetType() string {
	return security.Mock
}

// GetIdentityProvider returns the name of the identity provider that can authenticate the workload credential.
func (p *MockPlugin) GetIdentityProvider() string {
	return "fakeIDP"
}
