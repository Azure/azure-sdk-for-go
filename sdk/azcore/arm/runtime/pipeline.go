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
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// NewPipeline creates a pipeline from connection options.
// The telemetry policy, when enabled, will use the specified module and version info.
func NewPipeline(module, version string, cred azcore.TokenCredential, plOpts azruntime.PipelineOptions, options *arm.ClientOptions) pipeline.Pipeline {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := options.Host
	if len(ep) == 0 {
		ep = arm.AzurePublicCloud
	}
	authPolicy := NewBearerTokenPolicy(cred, &armpolicy.BearerTokenOptions{
		Scopes:           []string{shared.EndpointToScope(string(ep))},
		AuxiliaryTenants: options.AuxiliaryTenants,
	})
	o := plOpts
	o.PerRetry = append(o.PerRetry, authPolicy)
	if !options.DisableRPRegistration {
		regRPOpts := armpolicy.RegistrationOptions{ClientOptions: options.ClientOptions}
		o.PerCall = append(o.PerCall, NewRPRegistrationPolicy(string(ep), cred, &regRPOpts))
	}
	return azruntime.NewPipeline(module, version, o, &options.ClientOptions)
}
