package appconfiguration

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
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// PrivateLinkResourcesClient is the client for the PrivateLinkResources methods of the Appconfiguration service.
type PrivateLinkResourcesClient struct {
	BaseClient
}

// NewPrivateLinkResourcesClient creates an instance of the PrivateLinkResourcesClient client.
func NewPrivateLinkResourcesClient(subscriptionID string) PrivateLinkResourcesClient {
	return NewPrivateLinkResourcesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewPrivateLinkResourcesClientWithBaseURI creates an instance of the PrivateLinkResourcesClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewPrivateLinkResourcesClientWithBaseURI(baseURI string, subscriptionID string) PrivateLinkResourcesClient {
	return PrivateLinkResourcesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets a private link resource that need to be created for a configuration store.
// Parameters:
// resourceGroupName - the name of the resource group to which the container registry belongs.
// configStoreName - the name of the configuration store.
// groupName - the name of the private link resource group.
func (client PrivateLinkResourcesClient) Get(ctx context.Context, resourceGroupName string, configStoreName string, groupName string) (result PrivateLinkResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateLinkResourcesClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: configStoreName,
			Constraints: []validation.Constraint{{Target: "configStoreName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "configStoreName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "configStoreName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9_-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("appconfiguration.PrivateLinkResourcesClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, configStoreName, groupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client PrivateLinkResourcesClient) GetPreparer(ctx context.Context, resourceGroupName string, configStoreName string, groupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"configStoreName":   autorest.Encode("path", configStoreName),
		"groupName":         autorest.Encode("path", groupName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppConfiguration/configurationStores/{configStoreName}/privateLinkResources/{groupName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateLinkResourcesClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client PrivateLinkResourcesClient) GetResponder(resp *http.Response) (result PrivateLinkResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByConfigurationStore gets the private link resources that need to be created for a configuration store.
// Parameters:
// resourceGroupName - the name of the resource group to which the container registry belongs.
// configStoreName - the name of the configuration store.
func (client PrivateLinkResourcesClient) ListByConfigurationStore(ctx context.Context, resourceGroupName string, configStoreName string) (result PrivateLinkResourceListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateLinkResourcesClient.ListByConfigurationStore")
		defer func() {
			sc := -1
			if result.plrlr.Response.Response != nil {
				sc = result.plrlr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: configStoreName,
			Constraints: []validation.Constraint{{Target: "configStoreName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "configStoreName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "configStoreName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9_-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("appconfiguration.PrivateLinkResourcesClient", "ListByConfigurationStore", err.Error())
	}

	result.fn = client.listByConfigurationStoreNextResults
	req, err := client.ListByConfigurationStorePreparer(ctx, resourceGroupName, configStoreName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "ListByConfigurationStore", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByConfigurationStoreSender(req)
	if err != nil {
		result.plrlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "ListByConfigurationStore", resp, "Failure sending request")
		return
	}

	result.plrlr, err = client.ListByConfigurationStoreResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "ListByConfigurationStore", resp, "Failure responding to request")
	}
	if result.plrlr.hasNextLink() && result.plrlr.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListByConfigurationStorePreparer prepares the ListByConfigurationStore request.
func (client PrivateLinkResourcesClient) ListByConfigurationStorePreparer(ctx context.Context, resourceGroupName string, configStoreName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"configStoreName":   autorest.Encode("path", configStoreName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppConfiguration/configurationStores/{configStoreName}/privateLinkResources", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByConfigurationStoreSender sends the ListByConfigurationStore request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateLinkResourcesClient) ListByConfigurationStoreSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByConfigurationStoreResponder handles the response to the ListByConfigurationStore request. The method always
// closes the http.Response Body.
func (client PrivateLinkResourcesClient) ListByConfigurationStoreResponder(resp *http.Response) (result PrivateLinkResourceListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByConfigurationStoreNextResults retrieves the next set of results, if any.
func (client PrivateLinkResourcesClient) listByConfigurationStoreNextResults(ctx context.Context, lastResults PrivateLinkResourceListResult) (result PrivateLinkResourceListResult, err error) {
	req, err := lastResults.privateLinkResourceListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "listByConfigurationStoreNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByConfigurationStoreSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "listByConfigurationStoreNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByConfigurationStoreResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appconfiguration.PrivateLinkResourcesClient", "listByConfigurationStoreNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByConfigurationStoreComplete enumerates all values, automatically crossing page boundaries as required.
func (client PrivateLinkResourcesClient) ListByConfigurationStoreComplete(ctx context.Context, resourceGroupName string, configStoreName string) (result PrivateLinkResourceListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateLinkResourcesClient.ListByConfigurationStore")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByConfigurationStore(ctx, resourceGroupName, configStoreName)
	return
}
