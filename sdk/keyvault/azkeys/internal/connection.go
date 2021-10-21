//go:build go1.16
// +build go1.16

// Code generated by Microsoft (R) AutoRest Code Generator (autorest: 3.6.2, generator: @autorest/go@4.0.0-preview.29)
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

var scopes = []string{"https://vault.azure.net/.default"}
// ConnectionOptions contains configuration settings for the connection's pipeline.
// All zero-value fields will be initialized with their default values.
type ConnectionOptions struct {
	// Transport sets the transport for making HTTP requests.
	Transport policy.Transporter
	// Retry configures the built-in retry policy behavior.
	Retry policy.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry policy.TelemetryOptions
	// Logging configures the built-in logging policy behavior.
	Logging policy.LogOptions
	// PerCallPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request.
	PerCallPolicies []policy.Policy
	// PerRetryPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request, and for each retry request.
	PerRetryPolicies []policy.Policy
}

// connection - The key vault client performs cryptographic key operations and vault operations against the Key Vault service.
type connection struct {
	p runtime.Pipeline
}

// NewConnection creates an instance of the connection type with the specified endpoint.
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func NewConnection(cred azcore.Credential, options *ConnectionOptions) *connection {
	if options == nil {
		options = &ConnectionOptions{}
	}
	policies := []policy.Policy{}
	if !options.Telemetry.Disabled {
		policies = append(policies, runtime.NewTelemetryPolicy(module, version, &options.Telemetry))
	}
	policies = append(policies, options.PerCallPolicies...)
	policies = append(policies, runtime.NewRetryPolicy(&options.Retry))
	policies = append(policies, options.PerRetryPolicies...)
		policies = append(policies, cred.NewAuthenticationPolicy(runtime.AuthenticationOptions{TokenRequest: policy.TokenRequestOptions{Scopes: scopes}}))
	policies = append(policies, runtime.NewLogPolicy(&options.Logging))
	client := &connection{
		p: runtime.NewPipeline(options.Transport, policies...),
	}
	return client
}

// Pipeline returns the connection's pipeline.
func (c *connection) Pipeline() (runtime.Pipeline) {
	return c.p
}

