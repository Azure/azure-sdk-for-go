// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
)

// Client allows you to administer resources in a Service Bus Namespace.
// For example, you can create queues, enabling capabilities like partitioning, duplicate detection, etc..
// NOTE: For sending and receiving messages you'll need to use the `Client` type instead.
type Client struct {
	em *atom.EntityManager
}

type AdminClientOptions struct {
	// for future expansion
}

// NewClientFromConnectionString creates an Client authenticating using a connection string.
func NewClientFromConnectionString(connectionString string, options *AdminClientOptions) (*Client, error) {
	em, err := atom.NewEntityManagerWithConnectionString(connectionString, internal.Version)

	if err != nil {
		return nil, err
	}

	return &Client{em: em}, nil
}

// NewClient creates an Client authenticating using a TokenCredential.
func NewClient(fullyQualifiedNamespace string, tokenCredential azcore.TokenCredential, options *AdminClientOptions) (*Client, error) {
	em, err := atom.NewEntityManager(fullyQualifiedNamespace, tokenCredential, internal.Version)

	if err != nil {
		return nil, err
	}

	return &Client{em: em}, nil
}

func (ac *Client) GetNamespaceProperties() {
	panic("not yet done")
}
