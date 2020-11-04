package util

import (
	"istio.io/istio/pkg/config/resource"
)

// IsSystemNamespace returns true for system namespaces
func IsSystemNamespace(ns resource.Namespace) bool {
	return ns == "kube-system" || ns == "kube-public"
}

// IsIstioControlPlane returns true for resources that are part of the Istio control plane
func IsIstioControlPlane(r *resource.Instance) bool {
	if _, ok := r.Metadata.Labels["istio"]; ok {
		return true
	}
	if r.Metadata.Labels["release"] == "istio" {
		return true
	}
	return false
}
