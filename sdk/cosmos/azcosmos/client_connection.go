// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// clientConnection maintains a Pipeline for the client.
// The Pipeline is build based on the CosmosClientOptions.
type clientConnection struct {
	endpoint string
	pipeline azcore.Pipeline
}

// newConnection creates an instance of the connection type with the specified endpoint.
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func newConnection(endpoint string, cred azcore.Credential, options *CosmosClientOptions) *clientConnection {
	if options == nil {
		options = &CosmosClientOptions{}
	}
	policies := []azcore.Policy{
		azcore.NewTelemetryPolicy(options.enrichTelemetryOptions()),
	}
	policies = append(policies, options.PerCallPolicies...)
	policies = append(policies, azcore.NewRetryPolicy(&options.Retry))
	policies = append(policies, options.PerRetryPolicies...)
	policies = append(policies, options.getSDKInternalPolicies()...)
	policies = append(policies, cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{"none"}}}))
	policies = append(policies, azcore.NewLogPolicy(&options.Logging))
	return &clientConnection{endpoint: endpoint, pipeline: azcore.NewPipeline(options.HTTPClient, policies...)}
}

// Endpoint returns the connection's endpoint.
func (c *clientConnection) Endpoint() string {
	return c.endpoint
}

// Pipeline returns the connection's pipeline.
func (c *clientConnection) Pipeline() azcore.Pipeline {
	return c.pipeline
}
