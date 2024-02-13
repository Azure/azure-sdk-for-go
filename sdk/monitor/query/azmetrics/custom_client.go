//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azmetrics

// this file contains handwritten additions to the generated code

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a client that accesses Azure Monitor metrics data.
// Client should be used for performing metrics queries on multiple monitored resources in the same region.
// A credential with authorization at the subscription level is required when using this client.
//
// endpoint - The regional endpoint to use, for example https://eastus.metrics.monitor.azure.com.
// The region should match the region of the requested resources. For global resources, the region should be 'global'.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{"https://metrics.monitor.azure.com/.default"}, nil)
	azcoreClient, err := azcore.NewClient(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{endpoint: endpoint, internal: azcoreClient}, nil
}
