//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmediaservices

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

// ClientFactory is a client factory used to create any client in this module.
// Don't use this type directly, use NewClientFactory instead.
type ClientFactory struct {
	subscriptionID string
	credential     azcore.TokenCredential
	options        *arm.ClientOptions
}

// NewClientFactory creates a new instance of ClientFactory with the specified values.
// The parameter values will be propagated to any client created from this factory.
//   - subscriptionID - The unique identifier for a Microsoft Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewClientFactory(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClientFactory, error) {
	_, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	return &ClientFactory{
		subscriptionID: subscriptionID, credential: credential,
		options: options.Clone(),
	}, nil
}

// NewAccountFiltersClient creates a new instance of AccountFiltersClient.
func (c *ClientFactory) NewAccountFiltersClient() *AccountFiltersClient {
	subClient, _ := NewAccountFiltersClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewAssetFiltersClient creates a new instance of AssetFiltersClient.
func (c *ClientFactory) NewAssetFiltersClient() *AssetFiltersClient {
	subClient, _ := NewAssetFiltersClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewAssetTrackOperationResultsClient creates a new instance of AssetTrackOperationResultsClient.
func (c *ClientFactory) NewAssetTrackOperationResultsClient() *AssetTrackOperationResultsClient {
	subClient, _ := NewAssetTrackOperationResultsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewAssetTrackOperationStatusesClient creates a new instance of AssetTrackOperationStatusesClient.
func (c *ClientFactory) NewAssetTrackOperationStatusesClient() *AssetTrackOperationStatusesClient {
	subClient, _ := NewAssetTrackOperationStatusesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewAssetsClient creates a new instance of AssetsClient.
func (c *ClientFactory) NewAssetsClient() *AssetsClient {
	subClient, _ := NewAssetsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewClient creates a new instance of Client.
func (c *ClientFactory) NewClient() *Client {
	subClient, _ := NewClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewContentKeyPoliciesClient creates a new instance of ContentKeyPoliciesClient.
func (c *ClientFactory) NewContentKeyPoliciesClient() *ContentKeyPoliciesClient {
	subClient, _ := NewContentKeyPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewJobsClient creates a new instance of JobsClient.
func (c *ClientFactory) NewJobsClient() *JobsClient {
	subClient, _ := NewJobsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewLiveEventsClient creates a new instance of LiveEventsClient.
func (c *ClientFactory) NewLiveEventsClient() *LiveEventsClient {
	subClient, _ := NewLiveEventsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewLiveOutputsClient creates a new instance of LiveOutputsClient.
func (c *ClientFactory) NewLiveOutputsClient() *LiveOutputsClient {
	subClient, _ := NewLiveOutputsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewLocationsClient creates a new instance of LocationsClient.
func (c *ClientFactory) NewLocationsClient() *LocationsClient {
	subClient, _ := NewLocationsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewOperationResultsClient creates a new instance of OperationResultsClient.
func (c *ClientFactory) NewOperationResultsClient() *OperationResultsClient {
	subClient, _ := NewOperationResultsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewOperationStatusesClient creates a new instance of OperationStatusesClient.
func (c *ClientFactory) NewOperationStatusesClient() *OperationStatusesClient {
	subClient, _ := NewOperationStatusesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewOperationsClient creates a new instance of OperationsClient.
func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	subClient, _ := NewOperationsClient(c.credential, c.options)
	return subClient
}

// NewPrivateEndpointConnectionsClient creates a new instance of PrivateEndpointConnectionsClient.
func (c *ClientFactory) NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient {
	subClient, _ := NewPrivateEndpointConnectionsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewPrivateLinkResourcesClient creates a new instance of PrivateLinkResourcesClient.
func (c *ClientFactory) NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient {
	subClient, _ := NewPrivateLinkResourcesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewStreamingEndpointsClient creates a new instance of StreamingEndpointsClient.
func (c *ClientFactory) NewStreamingEndpointsClient() *StreamingEndpointsClient {
	subClient, _ := NewStreamingEndpointsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewStreamingLocatorsClient creates a new instance of StreamingLocatorsClient.
func (c *ClientFactory) NewStreamingLocatorsClient() *StreamingLocatorsClient {
	subClient, _ := NewStreamingLocatorsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewStreamingPoliciesClient creates a new instance of StreamingPoliciesClient.
func (c *ClientFactory) NewStreamingPoliciesClient() *StreamingPoliciesClient {
	subClient, _ := NewStreamingPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewTracksClient creates a new instance of TracksClient.
func (c *ClientFactory) NewTracksClient() *TracksClient {
	subClient, _ := NewTracksClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

// NewTransformsClient creates a new instance of TransformsClient.
func (c *ClientFactory) NewTransformsClient() *TransformsClient {
	subClient, _ := NewTransformsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}
