//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"fmt"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// NewPipeline creates a pipeline from connection options.
// The telemetry policy, when enabled, will use the specified module and version info.
// This method panics when the ClientOptions.Cloud field is set with a Configuration object
// that's missing Azure Resource Manager settings. A future version will return an error instead.
func NewPipeline(module, version string, cred shared.TokenCredential, plOpts azruntime.PipelineOptions, options *arm.ClientOptions) pipeline.Pipeline {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	conf, err := getConfiguration(&options.ClientOptions)
	if err != nil {
		// TODO: return the error
		panic(err)
	}
	scope := conf.Audiences[0] + "/.default"
	ep := conf.Endpoint
	if options.Endpoint != "" {
		scope = shared.EndpointToScope(string(options.Endpoint))
		ep = string(options.Endpoint)
	}
	authPolicy := NewBearerTokenPolicy(cred, &armpolicy.BearerTokenOptions{
		Scopes:           []string{scope},
		AuxiliaryTenants: options.AuxiliaryTenants,
	})
	perRetry := make([]pipeline.Policy, 0, len(plOpts.PerRetry)+1)
	copy(perRetry, plOpts.PerRetry)
	plOpts.PerRetry = append(perRetry, authPolicy)
	if !options.DisableRPRegistration {
		regRPOpts := armpolicy.RegistrationOptions{ClientOptions: options.ClientOptions}
		regPolicy := NewRPRegistrationPolicy(ep, cred, &regRPOpts)
		perCall := make([]pipeline.Policy, 0, len(plOpts.PerCall)+1)
		copy(perCall, plOpts.PerCall)
		plOpts.PerCall = append(perCall, regPolicy)
	}
	return azruntime.NewPipeline(module, version, plOpts, &options.ClientOptions)
}

func getConfiguration(o *azpolicy.ClientOptions) (cloud.ServiceConfiguration, error) {
	c := cloud.WellKnownClouds[cloud.AzurePublicCloud]
	if !reflect.ValueOf(o.Cloud).IsZero() {
		c = o.Cloud
	}
	if conf, ok := c.Services[cloud.ResourceManager]; ok && conf.Endpoint != "" && len(conf.Audiences) > 0 {
		return conf, nil
	} else {
		return conf, fmt.Errorf(`missing Azure Resource Manager configuration for cloud "%s"`, c.Name)
	}
}
