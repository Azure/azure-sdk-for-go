//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// NewPipeline creates a pipeline from connection options.
// The telemetry policy, when enabled, will use the specified module and version info.
func NewPipeline(module, version string, cred azcore.TokenCredential, options *arm.ClientOptions) pipeline.Pipeline {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := options.Host
	if len(ep) == 0 {
		ep = arm.AzurePublicCloud
	}
	policies := []policy.Policy{}
	if !options.Telemetry.Disabled {
		policies = append(policies, azruntime.NewTelemetryPolicy(module, version, &options.Telemetry))
	}
	if !options.DisableRPRegistration {
		regRPOpts := RegistrationOptions{
			HTTPClient: options.Transport,
			Logging:    options.Logging,
			Retry:      options.Retry,
			Telemetry:  options.Telemetry,
		}
		policies = append(policies, NewRPRegistrationPolicy(string(ep), cred, &regRPOpts))
	}
	policies = append(policies, options.PerCallPolicies...)
	policies = append(policies, azruntime.NewRetryPolicy(&options.Retry))
	policies = append(policies, options.PerRetryPolicies...)
	policies = append(policies,
		azruntime.NewBearerTokenPolicy(cred, azruntime.AuthenticationOptions{
			TokenRequest: policy.TokenRequestOptions{
				Scopes: []string{shared.EndpointToScope(string(ep))},
			},
			AuxiliaryTenants: options.AuxiliaryTenants,
		},
		),
		azruntime.NewLogPolicy(&options.Logging))
	return azruntime.NewPipeline(options.Transport, policies...)
}
