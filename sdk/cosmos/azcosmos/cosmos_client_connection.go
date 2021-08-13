// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// cosmosClientConnection maintains a Pipeline for the client.
// The Pipeline is build based on the CosmosClientOptions.
type cosmosClientConnection struct {
	endpoint string
	Pipeline azcore.Pipeline
}

// newConnection creates an instance of the connection type with the specified endpoint.
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func newCosmosClientConnection(endpoint string, cred azcore.Credential, options *CosmosClientOptions) *cosmosClientConnection {
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
	return &cosmosClientConnection{endpoint: endpoint, Pipeline: azcore.NewPipeline(options.HTTPClient, policies...)}
}

func (cc *cosmosClientConnection) getPath(parentPath string, pathSegment string, id string) string {
	var completePath strings.Builder
	parentPathLength := len(parentPath)
	completePath.Grow(parentPathLength + 2 + len(pathSegment) + len(id))
	if parentPathLength > 0 {
		completePath.WriteString(parentPath)
		completePath.WriteString("/")
	}
	completePath.WriteString(pathSegment)
	completePath.WriteString("/")
	completePath.WriteString(url.QueryEscape(id))
	return completePath.String()
}
