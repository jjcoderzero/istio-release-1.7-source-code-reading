package cache

import "istio.io/pkg/monitoring"

const (
	TokenExchange = "token_exchange"
	CSR           = "csr"
)

var (
	RequestType = monitoring.MustCreateLabel("request_type")
)

// Metrics for outgoing requests from citadel agent to external services such as token exchange server or a CA.
// This is different from incoming request metrics (i.e. from Envoy to citadel agent).
var (
	outgoingLatency = monitoring.NewSum(
		"outgoing_latency",
		"The latency of "+
			"outgoing requests (e.g. to a token exchange server, CA, etc.) in milliseconds.",
		monitoring.WithLabels(RequestType), monitoring.WithUnit(monitoring.Milliseconds))

	numOutgoingRequests = monitoring.NewSum(
		"num_outgoing_requests",
		"Number of total outgoing requests (e.g. to a token exchange server, CA, etc.)",
		monitoring.WithLabels(RequestType))

	numOutgoingRetries = monitoring.NewSum(
		"num_outgoing_retries",
		"Number of outgoing retry requests (e.g. to a token exchange server, CA, etc.)",
		monitoring.WithLabels(RequestType))

	numFailedOutgoingRequests = monitoring.NewSum(
		"num_failed_outgoing_requests",
		"Number of failed outgoing requests (e.g. to a token exchange server, CA, etc.)",
		monitoring.WithLabels(RequestType))
)

func init() {
	monitoring.MustRegister(
		outgoingLatency,
		numOutgoingRequests,
		numOutgoingRetries,
		numFailedOutgoingRequests,
	)
}
