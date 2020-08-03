package siterecovery

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// ReplicationProtectableItemsClient is the client for the ReplicationProtectableItems methods of the Siterecovery
// service.
type ReplicationProtectableItemsClient struct {
	BaseClient
}

// NewReplicationProtectableItemsClient creates an instance of the ReplicationProtectableItemsClient client.
func NewReplicationProtectableItemsClient(subscriptionID string, resourceGroupName string, resourceName string) ReplicationProtectableItemsClient {
	return NewReplicationProtectableItemsClientWithBaseURI(DefaultBaseURI, subscriptionID, resourceGroupName, resourceName)
}

// NewReplicationProtectableItemsClientWithBaseURI creates an instance of the ReplicationProtectableItemsClient client
// using a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign
// clouds, Azure stack).
func NewReplicationProtectableItemsClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string, resourceName string) ReplicationProtectableItemsClient {
	return ReplicationProtectableItemsClient{NewWithBaseURI(baseURI, subscriptionID, resourceGroupName, resourceName)}
}

// Get the operation to get the details of a protectable item.
// Parameters:
// fabricName - fabric name.
// protectionContainerName - protection container name.
// protectableItemName - protectable item name.
func (client ReplicationProtectableItemsClient) Get(ctx context.Context, fabricName string, protectionContainerName string, protectableItemName string) (result ProtectableItem, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ReplicationProtectableItemsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, fabricName, protectionContainerName, protectableItemName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ReplicationProtectableItemsClient) GetPreparer(ctx context.Context, fabricName string, protectionContainerName string, protectableItemName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fabricName":              autorest.Encode("path", fabricName),
		"protectableItemName":     autorest.Encode("path", protectableItemName),
		"protectionContainerName": autorest.Encode("path", protectionContainerName),
		"resourceGroupName":       autorest.Encode("path", client.ResourceGroupName),
		"resourceName":            autorest.Encode("path", client.ResourceName),
		"subscriptionId":          autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-08-10"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{fabricName}/replicationProtectionContainers/{protectionContainerName}/replicationProtectableItems/{protectableItemName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ReplicationProtectableItemsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ReplicationProtectableItemsClient) GetResponder(resp *http.Response) (result ProtectableItem, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByReplicationProtectionContainers lists the protectable items in a protection container.
// Parameters:
// fabricName - fabric name.
// protectionContainerName - protection container name.
func (client ReplicationProtectableItemsClient) ListByReplicationProtectionContainers(ctx context.Context, fabricName string, protectionContainerName string) (result ProtectableItemCollectionPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ReplicationProtectableItemsClient.ListByReplicationProtectionContainers")
		defer func() {
			sc := -1
			if result.pic.Response.Response != nil {
				sc = result.pic.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listByReplicationProtectionContainersNextResults
	req, err := client.ListByReplicationProtectionContainersPreparer(ctx, fabricName, protectionContainerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByReplicationProtectionContainersSender(req)
	if err != nil {
		result.pic.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", resp, "Failure sending request")
		return
	}

	result.pic, err = client.ListByReplicationProtectionContainersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", resp, "Failure responding to request")
	}
	if result.pic.hasNextLink() && result.pic.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListByReplicationProtectionContainersPreparer prepares the ListByReplicationProtectionContainers request.
func (client ReplicationProtectableItemsClient) ListByReplicationProtectionContainersPreparer(ctx context.Context, fabricName string, protectionContainerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fabricName":              autorest.Encode("path", fabricName),
		"protectionContainerName": autorest.Encode("path", protectionContainerName),
		"resourceGroupName":       autorest.Encode("path", client.ResourceGroupName),
		"resourceName":            autorest.Encode("path", client.ResourceName),
		"subscriptionId":          autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-08-10"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{fabricName}/replicationProtectionContainers/{protectionContainerName}/replicationProtectableItems", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByReplicationProtectionContainersSender sends the ListByReplicationProtectionContainers request. The method will close the
// http.Response Body if it receives an error.
func (client ReplicationProtectableItemsClient) ListByReplicationProtectionContainersSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByReplicationProtectionContainersResponder handles the response to the ListByReplicationProtectionContainers request. The method always
// closes the http.Response Body.
func (client ReplicationProtectableItemsClient) ListByReplicationProtectionContainersResponder(resp *http.Response) (result ProtectableItemCollection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByReplicationProtectionContainersNextResults retrieves the next set of results, if any.
func (client ReplicationProtectableItemsClient) listByReplicationProtectionContainersNextResults(ctx context.Context, lastResults ProtectableItemCollection) (result ProtectableItemCollection, err error) {
	req, err := lastResults.protectableItemCollectionPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "listByReplicationProtectionContainersNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByReplicationProtectionContainersSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "listByReplicationProtectionContainersNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByReplicationProtectionContainersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "siterecovery.ReplicationProtectableItemsClient", "listByReplicationProtectionContainersNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByReplicationProtectionContainersComplete enumerates all values, automatically crossing page boundaries as required.
func (client ReplicationProtectableItemsClient) ListByReplicationProtectionContainersComplete(ctx context.Context, fabricName string, protectionContainerName string) (result ProtectableItemCollectionIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ReplicationProtectableItemsClient.ListByReplicationProtectionContainers")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByReplicationProtectionContainers(ctx, fabricName, protectionContainerName)
	return
}
