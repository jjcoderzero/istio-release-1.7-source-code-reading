package crdclient

import (
	"istio.io/pkg/monitoring"

	"istio.io/istio/pilot/pkg/model"
)

var (
	typeTag  = monitoring.MustCreateLabel("type")
	eventTag = monitoring.MustCreateLabel("event")
	nameTag  = monitoring.MustCreateLabel("name")

	k8sEvents = monitoring.NewSum(
		"pilot_k8s_cfg_events",
		"Events from k8s config.",
		monitoring.WithLabels(typeTag, eventTag),
	)

	k8sErrors = monitoring.NewGauge(
		"pilot_k8s_object_errors",
		"Errors converting k8s CRDs",
		monitoring.WithLabels(nameTag),
	)

	k8sTotalErrors = monitoring.NewSum(
		"pilot_total_k8s_object_errors",
		"Total Errors converting k8s CRDs",
	)
)

func init() {
	monitoring.MustRegister(k8sEvents)
	monitoring.MustRegister(k8sErrors)
	monitoring.MustRegister(k8sTotalErrors)
}

func incrementEvent(kind, event string) {
	k8sEvents.With(typeTag.Value(kind), eventTag.Value(event)).Increment()
}

func handleValidationFailure(obj *model.Config, err error) {
	key := obj.Namespace + "/" + obj.Name
	scope.Debugf("CRD validation failed: %s (%v): %v", key, obj.GroupVersionKind, err)
	k8sErrors.With(nameTag.Value(key)).Record(1)

	k8sTotalErrors.Increment()
}
