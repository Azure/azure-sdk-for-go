//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// Endpoint is the base URL for Azure Resource Manager.
type Endpoint string

const (
	// AzureChina is the Azure Resource Manager China cloud endpoint.
	AzureChina Endpoint = "https://management.chinacloudapi.cn/"
	// AzureGermany is the Azure Resource Manager Germany cloud endpoint.
	AzureGermany Endpoint = "https://management.microsoftazure.de/"
	// AzureGovernment is the Azure Resource Manager US government cloud endpoint.
	AzureGovernment Endpoint = "https://management.usgovcloudapi.net/"
	// AzurePublicCloud is the Azure Resource Manager public cloud endpoint.
	AzurePublicCloud Endpoint = "https://management.azure.com/"
)

// ConnectionOptions contains configuration settings for the connection's pipeline.
// All zero-value fields will be initialized with their default values.
type ConnectionOptions struct {
	// AuxiliaryTenants contains a list of additional tenants to be used to authenticate
	// across multiple tenants.
	AuxiliaryTenants []string

	// HTTPClient sets the transport for making HTTP requests.
	HTTPClient policy.Transporter

	// Retry configures the built-in retry policy behavior.
	Retry policy.RetryOptions

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry policy.TelemetryOptions

	// Logging configures the built-in logging policy behavior.
	Logging policy.LogOptions

	// DisableRPRegistration disables the auto-RP registration policy.
	// The default value is false.
	DisableRPRegistration bool

	// PerCallPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request.
	PerCallPolicies []policy.Policy

	// PerRetryPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request, and for each retry request.
	PerRetryPolicies []policy.Policy
}

// Connection is a connection to an Azure Resource Manager endpoint.
// It contains the base ARM endpoint and a pipeline for making requests.
type Connection struct {
	ep   Endpoint
	cred azcore.TokenCredential
	opt  ConnectionOptions
}

// NewDefaultConnection creates an instance of the Connection type using the AzurePublicCloud.
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func NewDefaultConnection(cred azcore.TokenCredential, options *ConnectionOptions) *Connection {
	return NewConnection(AzurePublicCloud, cred, options)
}

// NewConnection creates an instance of the Connection type with the specified endpoint.
// Use this when connecting to clouds other than the Azure public cloud (stack/sovereign clouds).
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func NewConnection(endpoint Endpoint, cred azcore.TokenCredential, options *ConnectionOptions) *Connection {
	if options == nil {
		options = &ConnectionOptions{}
	}
	return &Connection{ep: endpoint, cred: cred, opt: *options}
}

// Endpoint returns the connection's ARM endpoint.
func (con *Connection) Endpoint() Endpoint {
	return con.ep
}

// NewPipeline creates a pipeline from the connection's options.
// The telemetry policy, when enabled, will use the specified module and version info.
func (con *Connection) NewPipeline(module, version string) pipeline.Pipeline {
	policies := []policy.Policy{}
	if !con.opt.Telemetry.Disabled {
		policies = append(policies, azruntime.NewTelemetryPolicy(module, version, &con.opt.Telemetry))
	}
	if !con.opt.DisableRPRegistration {
		regRPOpts := armruntime.RegistrationOptions{
			HTTPClient: con.opt.HTTPClient,
			Logging:    con.opt.Logging,
			Retry:      con.opt.Retry,
			Telemetry:  con.opt.Telemetry,
		}
		policies = append(policies, armruntime.NewRPRegistrationPolicy(string(con.ep), con.cred, &regRPOpts))
	}
	policies = append(policies, con.opt.PerCallPolicies...)
	policies = append(policies, azruntime.NewRetryPolicy(&con.opt.Retry))
	policies = append(policies, con.opt.PerRetryPolicies...)
	policies = append(policies,
		con.cred.NewAuthenticationPolicy(
			azruntime.AuthenticationOptions{
				TokenRequest: policy.TokenRequestOptions{
					Scopes: []string{shared.EndpointToScope(string(con.ep))},
				},
				AuxiliaryTenants: con.opt.AuxiliaryTenants,
			},
		),
		azruntime.NewLogPolicy(&con.opt.Logging))
	return azruntime.NewPipeline(con.opt.HTTPClient, policies...)
}
