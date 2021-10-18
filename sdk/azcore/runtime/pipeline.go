//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// NewPipeline creates a pipeline from connection options, with any additional policies as specified.
// module, version: used by the telemetry policy, when enabled
// perCall: additional policies to invoke once per request
// perRetry: additional policies to invoke once per request and once per retry of that request
func NewPipeline(module, version string, perCall, perRetry []policy.Policy, options *policy.ClientOptions) Pipeline {
	if options == nil {
		options = &policy.ClientOptions{}
	}
	policies := []policy.Policy{}
	if !options.Telemetry.Disabled {
		policies = append(policies, NewTelemetryPolicy(module, version, &options.Telemetry))
	}
	policies = append(policies, options.PerCallPolicies...)
	policies = append(policies, perCall...)
	policies = append(policies, NewRetryPolicy(&options.Retry))
	policies = append(policies, options.PerRetryPolicies...)
	policies = append(policies, perRetry...)
	policies = append(policies, NewLogPolicy(&options.Logging))
	policies = append(policies, pipeline.PolicyFunc(httpHeaderPolicy), pipeline.PolicyFunc(bodyDownloadPolicy))
	transport := options.Transport
	if transport == nil {
		transport = defaultHTTPClient
	}
	return pipeline.NewPipeline(transport, policies...)
}
