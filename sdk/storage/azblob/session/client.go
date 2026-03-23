// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package session

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

// Client provides session functionality for the underlying service client.
type Client struct {
	client *service.Client
}

// NewClient creates a session client for the provided service client.
//   - client - an instance of a service client
func NewClient(serviceClient *service.Client) (*Client, error) {
	return &Client{
		client: serviceClient,
	}, nil
}

func (c *Client) generatedContainer(containerName string) *generated.ContainerClient {
	cClient := c.client.NewContainerClient(containerName)
	return base.InnerClient((*base.Client[generated.ContainerClient])(cClient))
}

// TODO: CreateContainerSession or other name ideas?
func (c *Client) ContainerCreateSession(ctx context.Context, containerName string) (ContainerCreateSessionResponse, error) {
	resp, err := c.generatedContainer(containerName).CreateSession(ctx, ContainerCreateSessionConfiguration{AuthenticationType: to.Ptr(generated.AuthenticationTypeHMAC)}, nil)
	return resp, err
}
