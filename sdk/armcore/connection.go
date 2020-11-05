// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const scope = "https://management.azure.com//.default"

// ConnectionOptions contains configuration settings for the connection's pipeline.
type ConnectionOptions struct {
	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient azcore.Transport
	// Retry configures the built-in retry policy behavior.
	Retry azcore.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions
	// RegisterRPOptions configures the built-in RP registration policy behavior.
	RegisterRPOptions RegistrationOptions
}

// DefaultConnectionOptions creates a ConnectionOptions type initialized with default values.
func DefaultConnectionOptions() ConnectionOptions {
	return ConnectionOptions{
		Retry:             azcore.DefaultRetryOptions(),
		RegisterRPOptions: DefaultRegistrationOptions(),
		Telemetry:         azcore.DefaultTelemetryOptions(),
	}
}

// Connection is a connection to an Azure Resource Manager endpoint.
// It contains the base ARM endpoint and a pipeline for making requests.
type Connection struct {
	u string
	p azcore.Pipeline
}

// DefaultEndpoint is the Azure Resourece Manager public cloud endpoint.
const DefaultEndpoint = "https://management.azure.com"

// NewDefaultConnection creates an instance of the Connection type using the DefaultEndpoint.
func NewDefaultConnection(cred azcore.TokenCredential, options *ConnectionOptions) *Connection {
	return NewConnection(DefaultEndpoint, cred, options)
}

// NewConnection creates an instance of the Connection type with the specified endpoint.
// Use this when connecting to clouds other than the Azure public cloud (stack/sovereign clouds).
func NewConnection(endpoint string, cred azcore.TokenCredential, options *ConnectionOptions) *Connection {
	if options == nil {
		o := DefaultConnectionOptions()
		options = &o
	}
	p := azcore.NewPipeline(options.HTTPClient,
		azcore.NewTelemetryPolicy(&options.Telemetry),
		NewRPRegistrationPolicy(cred, &options.RegisterRPOptions),
		azcore.NewRetryPolicy(&options.Retry),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewLogPolicy(nil))
	return NewConnectionWithPipeline(endpoint, p)
}

// NewConnectionWithPipeline creates an instance of the Connection type with the specified endpoint and pipeline.
// Use this when a custom pipeline is required.
func NewConnectionWithPipeline(endpoint string, p azcore.Pipeline) *Connection {
	return &Connection{u: endpoint, p: p}
}

// Do invokes the Do() method on the connection's pipeline.
func (c *Connection) Do(req *azcore.Request) (*azcore.Response, error) {
	return c.p.Do(req)
}

// Endpoint returns the connection's ARM endpoint.
func (c *Connection) Endpoint() string {
	return c.u
}
