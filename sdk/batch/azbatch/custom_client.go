// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azbatch

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ClientOptions contains optional settings for Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient constructs a Client
func NewClient(endpoint string, credential azcore.TokenCredential, opts *ClientOptions) (*Client, error) {
	return newClient(
		endpoint,
		runtime.NewBearerTokenPolicy(
			credential,
			[]string{"https://batch.core.windows.net//.default"},
			&policy.BearerTokenOptions{InsecureAllowCredentialWithHTTP: opts.InsecureAllowCredentialWithHTTP},
		),
		opts,
	)
}

func newClient(endpoint string, authPolicy policy.Policy, opts *ClientOptions) (*Client, error) {
	if opts == nil {
		opts = &ClientOptions{}
	}
	c, err := azcore.NewClient(moduleName, version, runtime.PipelineOptions{
		APIVersion: runtime.APIVersionOptions{
			Location: runtime.APIVersionLocationQueryParam,
			Name:     "api-version",
		},
		PerRetry: []policy.Policy{authPolicy},
		Tracing: runtime.TracingOptions{
			// TODO
			Namespace: "Azure.Compute.Batch",
		},
	}, &opts.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{endpoint: endpoint, internal: c}, nil
}
