//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid

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

// NewClientWithSharedKeyCredential creates a [Client] using a shared key.
func NewClientWithSharedKeyCredential(endpoint string, keyCred *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	azc, err := azcore.NewClient(internal.ModuleName+".Client", internal.ModuleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewKeyCredentialPolicy(keyCred, "Authorization", &runtime.KeyCredentialPolicyOptions{
				Prefix: "SharedAccessKey ",
			}),
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

// PublishCloudEvents - Publish Batch Cloud Event to namespace topic. In case of success, the server responds with an HTTP
// 200 status code with an empty JSON object in response. Otherwise, the server can return various error
// codes. For example, 401: which indicates authorization failure, 403: which indicates quota exceeded or message is too large,
// 410: which indicates that specific topic is not found, 400: for bad
// request, and 500: for internal server error.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01-preview
//   - endpoint - The host name of the namespace, e.g. namespaceName1.westus-1.eventgrid.azure.net
//   - topicName - Topic Name.
//   - events - Array of Cloud Events being published.
//   - options - ClientPublishCloudEventsOptions contains the optional parameters for the Client.PublishCloudEvents method.
func (client *Client) PublishCloudEvents(ctx context.Context, topicName string, events []messaging.CloudEvent, options *PublishCloudEventsOptions) (PublishCloudEventsResponse, error) {
	ctx = runtime.WithHTTPHeader(ctx, http.Header{
		"Content-type": []string{"application/cloudevents-batch+json; charset=utf-8"},
	})

	return client.internalPublishCloudEvents(ctx, topicName, events, options)
}
