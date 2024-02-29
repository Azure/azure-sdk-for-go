//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces/internal"
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

// RejectCloudEvents - Reject batch of Cloud Events. The server responds with an HTTP 200 status code if the request is successfully
// accepted. The response body will include the set of successfully rejected lockTokens,
// along with other failed lockTokens with their corresponding error information.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - topicName - Topic Name.
//   - eventSubscriptionName - Event Subscription Name.
//   - lockTokens - slice of lock tokens.
//   - options - RejectCloudEventsOptions contains the optional parameters for the Client.RejectCloudEvents method.
func (client *Client) RejectCloudEvents(ctx context.Context, topicName string, eventSubscriptionName string, lockTokens []string, options *RejectCloudEventsOptions) (RejectCloudEventsResponse, error) {
	return client.internalRejectCloudEvents(ctx, topicName, eventSubscriptionName, rejectOptions{LockTokens: lockTokens}, options)
}

// AcknowledgeCloudEvents - Acknowledge batch of Cloud Events. The server responds with an HTTP 200 status code if the request
// is successfully accepted. The response body will include the set of successfully acknowledged
// lockTokens, along with other failed lockTokens with their corresponding error information. Successfully acknowledged events
// will no longer be available to any consumer.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - topicName - Topic Name.
//   - eventSubscriptionName - Event Subscription Name.
//   - lockTokens - slice of lock tokens.
//   - options - AcknowledgeCloudEventsOptions contains the optional parameters for the Client.AcknowledgeCloudEvents method.
func (client *Client) AcknowledgeCloudEvents(ctx context.Context, topicName string, eventSubscriptionName string, lockTokens []string, options *AcknowledgeCloudEventsOptions) (AcknowledgeCloudEventsResponse, error) {
	return client.internalAcknowledgeCloudEvents(ctx, topicName, eventSubscriptionName, acknowledgeOptions{LockTokens: lockTokens}, options)
}

// ReleaseCloudEvents - Release batch of Cloud Events. The server responds with an HTTP 200 status code if the request is
// successfully accepted. The response body will include the set of successfully released lockTokens,
// along with other failed lockTokens with their corresponding error information.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - topicName - Topic Name.
//   - eventSubscriptionName - Event Subscription Name.
//   - lockTokens - slice of lock tokens.
//   - options - ReleaseCloudEventsOptions contains the optional parameters for the Client.ReleaseCloudEvents method.
func (client *Client) ReleaseCloudEvents(ctx context.Context, topicName string, eventSubscriptionName string, lockTokens []string, options *ReleaseCloudEventsOptions) (ReleaseCloudEventsResponse, error) {
	return client.internalReleaseCloudEvents(ctx, topicName, eventSubscriptionName, releaseOptions{LockTokens: lockTokens}, options)
}

// RenewCloudEventLocks - Renew lock for batch of Cloud Events. The server responds with an HTTP 200 status code if the request
// is successfully accepted. The response body will include the set of successfully renewed
// lockTokens, along with other failed lockTokens with their corresponding error information.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - topicName - Topic Name.
//   - eventSubscriptionName - Event Subscription Name.
//   - lockTokens - slice of lock tokens.
//   - options - RenewCloudEventLocksOptions contains the optional parameters for the Client.RenewCloudEventLocks method.
func (client *Client) RenewCloudEventLocks(ctx context.Context, topicName string, eventSubscriptionName string, lockTokens []string, options *RenewCloudEventLocksOptions) (RenewCloudEventLocksResponse, error) {
	return client.internalRenewCloudEventLocks(ctx, topicName, eventSubscriptionName, renewLockOptions{LockTokens: lockTokens}, options)
}
