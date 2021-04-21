package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type ClientOptions struct {
	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient azcore.Transport
	// Retry configures the built-in retry policy behavior.
	Retry azcore.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions
}

func (o *ClientOptions) getConnectionOptions() *connectionOptions {
	if o == nil {
		return nil
	}

	return &connectionOptions{
		HTTPClient: o.HTTPClient,
		Retry: o.Retry,
		Telemetry: o.Telemetry,
	}
}