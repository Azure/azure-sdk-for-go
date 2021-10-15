//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// NewDefaultPipeline creates a pipeline from connection options, with any additional policies as specified.
// module, version: used by the telemetry policy, when enabled
// perCallPolicies: additional policies to invoke once per request
// perRetryPolicies: additional policies to invoke once per request and once per retry of that request
func NewDefaultPipeline(module, version string, perCallPolicies []policy.Policy, perRetryPolicies []policy.Policy, options *policy.ClientOptions) pipeline.Pipeline {
	if options == nil {
		options = &policy.ClientOptions{}
	}
	policies := []policy.Policy{}
	if !options.Telemetry.Disabled {
		policies = append(policies, NewTelemetryPolicy(module, version, &options.Telemetry))
	}
	policies = append(policies, options.PerCallPolicies...)
	policies = append(policies, perCallPolicies...)
	policies = append(policies, NewRetryPolicy(&options.Retry))
	policies = append(policies, options.PerRetryPolicies...)
	policies = append(policies, perRetryPolicies...)
	policies = append(policies, NewLogPolicy(&options.Logging))
	policies = append(policies, pipeline.PolicyFunc(httpHeaderPolicy), pipeline.PolicyFunc(bodyDownloadPolicy))
	transport := options.Transport
	if transport == nil {
		transport = defaultHTTPClient
	}
	return pipeline.NewPipeline(transport, policies...)
}
