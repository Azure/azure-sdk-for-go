//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/template/aztemplate/internal"
)

type ClientOptions struct{}

type Client struct {
	client *internal.TemplateClient
}

func NewClient(cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {

	//if options == nil {
	//options = &ClientOptions{}
	//}
	tc := internal.NewTemplateClient()

	return &Client{client: tc}, nil

}

func (c *Client) SomeServiceAction() {
	c.client.PrintInfo()
}
