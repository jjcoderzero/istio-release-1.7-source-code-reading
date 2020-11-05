package memory_test

import (
	"testing"

	"istio.io/istio/pilot/pkg/config/memory"
	"istio.io/istio/pilot/test/mock"
	"istio.io/istio/pkg/config/schema/collections"
)

func TestStoreInvariant(t *testing.T) {
	store := memory.Make(collections.Mocks)
	mock.CheckMapInvariant(store, t, "some-namespace", 10)
}

func TestIstioConfig(t *testing.T) {
	store := memory.Make(collections.Pilot)
	mock.CheckIstioConfigTypes(store, "some-namespace", t)
}
