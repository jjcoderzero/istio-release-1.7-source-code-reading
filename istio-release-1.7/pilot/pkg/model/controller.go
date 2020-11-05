package model

// Controller defines an event controller loop.  Proxy agent registers itself
// with the controller loop and receives notifications on changes to the
// service topology or changes to the configuration artifacts.
//
// The controller guarantees the following consistency requirement: registry
// view in the controller is as AT LEAST as fresh as the moment notification
// arrives, but MAY BE more fresh (e.g. "delete" cancels an "add" event).  For
// example, an event for a service creation will see a service registry without
// the service if the event is immediately followed by the service deletion
// event.
//
// Handlers execute on the single worker queue in the order they are appended.
// Handlers receive the notification event and the associated object.  Note
// that all handlers must be appended before starting the controller.
type Controller interface {
	// AppendServiceHandler notifies about changes to the service catalog.
	AppendServiceHandler(f func(*Service, Event)) error

	// AppendInstanceHandler notifies about changes to the service instances
	// for a service.
	AppendInstanceHandler(f func(*ServiceInstance, Event)) error

	// AppendWorkloadHandler notifies about changes to workloads. This differs from InstanceHandler,
	// which deals with service instances (the result of a merge of Service and Workload)
	AppendWorkloadHandler(f func(*WorkloadInstance, Event)) error

	// Run until a signal is received
	Run(stop <-chan struct{})

	// HasSynced returns true after initial cache synchronization is complete
	HasSynced() bool
}

// Event represents a registry update event
type Event int

const (
	// EventAdd is sent when an object is added
	EventAdd Event = iota

	// EventUpdate is sent when an object is modified
	// Captures the modified object
	EventUpdate

	// EventDelete is sent when an object is deleted
	// Captures the object at the last known state
	EventDelete
)

func (event Event) String() string {
	out := "unknown"
	switch event {
	case EventAdd:
		out = "add"
	case EventUpdate:
		out = "update"
	case EventDelete:
		out = "delete"
	}
	return out
}
