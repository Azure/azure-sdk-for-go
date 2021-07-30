// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const defaultScope = "/.default"

const (
	// AzureChina is the Azure Resource Manager China cloud endpoint.
	AzureChina = "https://management.chinacloudapi.cn/"
	// AzureGermany is the Azure Resource Manager Germany cloud endpoint.
	AzureGermany = "https://management.microsoftazure.de/"
	// AzureGovernment is the Azure Resource Manager US government cloud endpoint.
	AzureGovernment = "https://management.usgovcloudapi.net/"
	// AzurePublicCloud is the Azure Resource Manager public cloud endpoint.
	AzurePublicCloud = "https://management.azure.com/"
)

// ConnectionOptions contains configuration settings for the connection's pipeline.
// All zero-value fields will be initialized with their default values.
type ConnectionOptions struct {
	// AuxiliaryTenants contains a list of additional tenants to be used to authenticate
	// across multiple tenants.
	AuxiliaryTenants []string

	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient azcore.Transport

	// Retry configures the built-in retry policy behavior.
	Retry azcore.RetryOptions

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions

	// Logging configures the built-in logging policy behavior.
	Logging azcore.LogOptions

	// DisableRPRegistration disables the auto-RP registration policy.
	// The default value is false.
	DisableRPRegistration bool

	// PerCallPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request.
	PerCallPolicies []azcore.Policy

	// PerRetryPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request, and for each retry request.
	PerRetryPolicies []azcore.Policy
}

// Connection is a connection to an Azure Resource Manager endpoint.
// It contains the base ARM endpoint and a pipeline for making requests.
type Connection struct {
	u string
	p azcore.Pipeline
}

// NewDefaultConnection creates an instance of the Connection type using the AzurePublicCloud.
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func NewDefaultConnection(cred azcore.TokenCredential, options *ConnectionOptions) *Connection {
	return NewConnection(AzurePublicCloud, cred, options)
}

// NewConnection creates an instance of the Connection type with the specified endpoint.
// Use this when connecting to clouds other than the Azure public cloud (stack/sovereign clouds).
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func NewConnection(endpoint string, cred azcore.TokenCredential, options *ConnectionOptions) *Connection {
	if options == nil {
		options = &ConnectionOptions{}
	} else {
		// create a copy so we don't modify the original
		cp := *options
		options = &cp
	}
	if options.Telemetry.Value == "" {
		options.Telemetry.Value = UserAgent
	} else {
		options.Telemetry.Value += " " + UserAgent
	}
	policies := []azcore.Policy{
		azcore.NewTelemetryPolicy(&options.Telemetry),
	}
	if !options.DisableRPRegistration {
		regRPOpts := RegistrationOptions{
			HTTPClient: options.HTTPClient,
			Logging:    options.Logging,
			Retry:      options.Retry,
			Telemetry:  options.Telemetry,
		}
		policies = append(policies, NewRPRegistrationPolicy(endpoint, cred, &regRPOpts))
	}
	policies = append(policies, options.PerCallPolicies...)
	policies = append(policies, azcore.NewRetryPolicy(&options.Retry))
	policies = append(policies, options.PerRetryPolicies...)
	policies = append(policies,
		cred.NewAuthenticationPolicy(
			azcore.AuthenticationOptions{
				TokenRequest: azcore.TokenRequestOptions{
					Scopes: []string{endpointToScope(endpoint)},
				},
				AuxiliaryTenants: options.AuxiliaryTenants,
			},
		),
		azcore.NewLogPolicy(&options.Logging))
	p := azcore.NewPipeline(options.HTTPClient, policies...)
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
