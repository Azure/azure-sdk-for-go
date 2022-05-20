//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/template/aztemplate/internal"
)

// ClientOptions contains optional parameters for NewClient
type ClientOptions struct {
	azcore.ClientOptions
}

// Client is the client to interact with.
// Don't use this type directly, use NewClient() instead.
type Client struct {
	client *internal.TemplateClient
}

// NewClient returns a pointer to a Client
func NewClient(cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(cred, []string{"service_scope"}, nil),
		},
	}, &options.ClientOptions)

	return &Client{client: internal.NewTemplateClient(pl)}, nil

}

// SomeServiceActionOptions contains the optional values for the Client.SomeServiceAction method.
type SomeServiceActionOptions struct {
	// placeholder for future options
}

// SomeServiceAction does some service action
func (c *Client) SomeServiceAction(ctx context.Context, options *SomeServiceActionOptions) {
	c.client.SomeAPI(ctx, nil)
}
