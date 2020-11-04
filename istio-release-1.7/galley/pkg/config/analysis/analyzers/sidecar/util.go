package sidecar

import "istio.io/istio/pkg/config/resource"

func getNames(entries []*resource.Instance) []string {
	names := make([]string, 0, len(entries))
	for _, rs := range entries {
		names = append(names, string(rs.Metadata.FullName.Name))
	}
	return names
}
