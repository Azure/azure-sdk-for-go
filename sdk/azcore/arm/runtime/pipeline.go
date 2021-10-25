//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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
	perCallPolicies := []azpolicy.Policy{}
	if !options.DisableRPRegistration {
		regRPOpts := armpolicy.RegistrationOptions{ClientOptions: options.ClientOptions}
		perCallPolicies = append(perCallPolicies, NewRPRegistrationPolicy(string(ep), cred, &regRPOpts))
	}
	perRetryPolicies := []azpolicy.Policy{
		NewBearerTokenPolicy(cred, &armpolicy.BearerTokenOptions{
			Scopes:           []string{shared.EndpointToScope(string(ep))},
			AuxiliaryTenants: options.AuxiliaryTenants,
		}),
	}
	return azruntime.NewPipeline(module, version, perCallPolicies, perRetryPolicies, &options.ClientOptions)
}
