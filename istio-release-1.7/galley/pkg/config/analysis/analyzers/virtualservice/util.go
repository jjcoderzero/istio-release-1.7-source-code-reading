package virtualservice

import (
	"istio.io/api/networking/v1alpha3"
)

func getRouteDestinations(vs *v1alpha3.VirtualService) []*v1alpha3.Destination {
	destinations := make([]*v1alpha3.Destination, 0)

	for _, r := range vs.GetTcp() {
		for _, rd := range r.GetRoute() {
			destinations = append(destinations, rd.GetDestination())
		}
	}
	for _, r := range vs.GetTls() {
		for _, rd := range r.GetRoute() {
			destinations = append(destinations, rd.GetDestination())
		}
	}
	for _, r := range vs.GetHttp() {
		for _, rd := range r.GetRoute() {
			destinations = append(destinations, rd.GetDestination())
		}
	}

	return destinations
}

func getHTTPMirrorDestinations(vs *v1alpha3.VirtualService) []*v1alpha3.Destination {
	var destinations []*v1alpha3.Destination

	for _, r := range vs.GetHttp() {
		if m := r.GetMirror(); m != nil {
			destinations = append(destinations, m)
		}
	}

	return destinations
}
