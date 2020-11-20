// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const defaultScope = "/.default"

const (
	// AzureChina is the Azure Resourece Manager China cloud endpoint.
	AzureChina = "https://management.chinacloudapi.cn/"
	// AzureGermany is the Azure Resourece Manager Germany cloud endpoint.
	AzureGermany = "https://management.microsoftazure.de/"
	// AzureGovernment is the Azure Resourece Manager US government cloud endpoint.
	AzureGovernment = "https://management.usgovcloudapi.net/"
	// AzurePublicCloud is the Azure Resourece Manager public cloud endpoint.
	AzurePublicCloud = "https://management.azure.com/"
)

// ConnectionOptions contains configuration settings for the connection's pipeline.
// Call DefaultConnectionOptions() to create an instance populated with default values.
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

// NewDefaultConnection creates an instance of the Connection type using the AzurePublicCloud.
func NewDefaultConnection(cred azcore.TokenCredential, options *ConnectionOptions) *Connection {
	return NewConnection(AzurePublicCloud, cred, options)
}

// NewConnection creates an instance of the Connection type with the specified endpoint.
// Use this when connecting to clouds other than the Azure public cloud (stack/sovereign clouds).
func NewConnection(endpoint string, cred azcore.TokenCredential, options *ConnectionOptions) *Connection {
	if options == nil {
		o := DefaultConnectionOptions()
		options = &o
	}
	if options.Telemetry.Value == "" {
		options.Telemetry.Value = UserAgent
	} else {
		options.Telemetry.Value += " " + UserAgent
	}
	p := azcore.NewPipeline(options.HTTPClient,
		azcore.NewTelemetryPolicy(&options.Telemetry),
		NewRPRegistrationPolicy(endpoint, cred, &options.RegisterRPOptions),
		azcore.NewRetryPolicy(&options.Retry),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{endpointToScope(endpoint)}}}),
		azcore.NewLogPolicy(nil))
	return NewConnectionWithPipeline(endpoint, p)
}

// NewConnectionWithPipeline creates an instance of the Connection type with the specified endpoint and pipeline.
// Use this when a custom pipeline is required.
func NewConnectionWithPipeline(endpoint string, p azcore.Pipeline) *Connection {
	return &Connection{u: endpoint, p: p}
}

// Endpoint returns the connection's ARM endpoint.
func (c *Connection) Endpoint() string {
	return c.u
}

// Pipeline returns the connection's pipeline.
func (c *Connection) Pipeline() azcore.Pipeline {
	return c.p
}

func endpointToScope(endpoint string) string {
	if endpoint[len(endpoint)-1] != '/' {
		endpoint += "/"
	}
	return endpoint + defaultScope
}
