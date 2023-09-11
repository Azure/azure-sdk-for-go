//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package publisher

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/internal"
)

// ClientOptions contains optional settings for [Client]
type ClientOptions struct {
	azcore.ClientOptions
}

var tokenScopes = []string{"https://eventgrid.azure.net/.default"}

// NewClient creates a [Client] that authenticates using a TokenCredential.
func NewClient(endpoint string, tokenCredential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	azc, err := azcore.NewClient(internal.ModuleName+".Client", internal.ModuleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(tokenCredential, tokenScopes, nil),
		},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azc,
		endpoint: endpoint,
	}, nil
}

// NewClientWithSharedKeyCredential creates a [Client] using a shared key.
func NewClientWithSharedKeyCredential(endpoint string, key string, options *ClientOptions) (*Client, error) {
	const sasKeyHeader = "aeg-sas-key"

	if options == nil {
		options = &ClientOptions{}
	}

	// TODO: I believe we're supposed to allow for dynamically updating the key at any time as well.
	azc, err := azcore.NewClient(internal.ModuleName+".Client", internal.ModuleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			// TODO: Java has a specific policy for this kind of authentication.
			// AzureKeyCredentialPolicy
			&skpolicy{HeaderName: sasKeyHeader, Key: key},
		},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azc,
		endpoint: endpoint,
	}, nil
}

// NewClientWithSharedKeyCredential creates a [Client] using a shared key.
func NewClientWithSAS(endpoint string, sas string, options *ClientOptions) (*Client, error) {
	const sasTokenHeader = "aeg-sas-token"

	if options == nil {
		options = &ClientOptions{}
	}

	// TODO: I believe we're supposed to allow for dynamically updating the key at any time as well.
	azc, err := azcore.NewClient(internal.ModuleName+".PublisherClient", internal.ModuleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			// TODO: Java has a specific policy for this kind of authentication.
			// AzureKeyCredentialPolicy
			&skpolicy{HeaderName: sasTokenHeader, Key: sas},
		},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azc,
		endpoint: endpoint,
	}, nil
}

// TODO: remove in favor of a common policy instead?
type skpolicy struct {
	Key        string
	HeaderName string
}

func (p *skpolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Add(p.HeaderName, p.Key)
	return req.Next()
}

// PublishCloudEvents - Publishes a batch of events to an Azure Event Grid topic.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-01-01
//   - events - An array of events to be published to Event Grid.
//   - options - ClientPublishCloudEventEventsOptions contains the optional parameters for the Client.PublishCloudEvents
//     method.
func (client *Client) PublishCloudEvents(ctx context.Context, events []messaging.CloudEvent, options *PublishCloudEventsOptions) (PublishCloudEventsResponse, error) {
	return client.internalPublishCloudEvents(ctx, events, options)
}
