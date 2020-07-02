package avs

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

// PrivateCloudsClient is the azure VMware Solution API
type PrivateCloudsClient struct {
	BaseClient
}

// NewPrivateCloudsClient creates an instance of the PrivateCloudsClient client.
func NewPrivateCloudsClient(subscriptionID string) PrivateCloudsClient {
	return NewPrivateCloudsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewPrivateCloudsClientWithBaseURI creates an instance of the PrivateCloudsClient client using a custom endpoint.
// Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewPrivateCloudsClientWithBaseURI(baseURI string, subscriptionID string) PrivateCloudsClient {
	return PrivateCloudsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate sends the create or update request.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// privateCloudName - name of the private cloud
// privateCloud - the private cloud
func (client PrivateCloudsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, privateCloudName string, privateCloud PrivateCloud) (result PrivateCloudsCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: privateCloud,
			Constraints: []validation.Constraint{{Target: "privateCloud.Sku", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "privateCloud.Sku.Name", Name: validation.Null, Rule: true, Chain: nil}}},
				{Target: "privateCloud.PrivateCloudProperties", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "privateCloud.PrivateCloudProperties.NetworkBlock", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("avs.PrivateCloudsClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, privateCloudName, privateCloud)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client PrivateCloudsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, privateCloudName string, privateCloud PrivateCloud) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"privateCloudName":  autorest.Encode("path", privateCloudName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-03-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}", pathParameters),
		autorest.WithJSON(privateCloud),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateCloudsClient) CreateOrUpdateSender(req *http.Request) (future PrivateCloudsCreateOrUpdateFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client PrivateCloudsClient) CreateOrUpdateResponder(resp *http.Response) (result PrivateCloud, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete sends the delete request.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// privateCloudName - name of the private cloud
func (client PrivateCloudsClient) Delete(ctx context.Context, resourceGroupName string, privateCloudName string) (result PrivateCloudsDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.Delete")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("avs.PrivateCloudsClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, resourceGroupName, privateCloudName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client PrivateCloudsClient) DeletePreparer(ctx context.Context, resourceGroupName string, privateCloudName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"privateCloudName":  autorest.Encode("path", privateCloudName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-03-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateCloudsClient) DeleteSender(req *http.Request) (future PrivateCloudsDeleteFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client PrivateCloudsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get sends the get request.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// privateCloudName - name of the private cloud
func (client PrivateCloudsClient) Get(ctx context.Context, resourceGroupName string, privateCloudName string) (result PrivateCloud, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("avs.PrivateCloudsClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, privateCloudName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client PrivateCloudsClient) GetPreparer(ctx context.Context, resourceGroupName string, privateCloudName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"privateCloudName":  autorest.Encode("path", privateCloudName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-03-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateCloudsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client PrivateCloudsClient) GetResponder(resp *http.Response) (result PrivateCloud, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List sends the list request.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
func (client PrivateCloudsClient) List(ctx context.Context, resourceGroupName string) (result PrivateCloudListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.List")
		defer func() {
			sc := -1
			if result.pcl.Response.Response != nil {
				sc = result.pcl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("avs.PrivateCloudsClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.pcl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "List", resp, "Failure sending request")
		return
	}

	result.pcl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client PrivateCloudsClient) ListPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-03-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateCloudsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client PrivateCloudsClient) ListResponder(resp *http.Response) (result PrivateCloudList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client PrivateCloudsClient) listNextResults(ctx context.Context, lastResults PrivateCloudList) (result PrivateCloudList, err error) {
	req, err := lastResults.privateCloudListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client PrivateCloudsClient) ListComplete(ctx context.Context, resourceGroupName string) (result PrivateCloudListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, resourceGroupName)
	return
}

// ListAdminCredentials sends the list admin credentials request.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// privateCloudName - name of the private cloud
func (client PrivateCloudsClient) ListAdminCredentials(ctx context.Context, resourceGroupName string, privateCloudName string) (result AdminCredentials, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.ListAdminCredentials")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("avs.PrivateCloudsClient", "ListAdminCredentials", err.Error())
	}

	req, err := client.ListAdminCredentialsPreparer(ctx, resourceGroupName, privateCloudName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "ListAdminCredentials", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListAdminCredentialsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "ListAdminCredentials", resp, "Failure sending request")
		return
	}

	result, err = client.ListAdminCredentialsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "ListAdminCredentials", resp, "Failure responding to request")
	}

	return
}

// ListAdminCredentialsPreparer prepares the ListAdminCredentials request.
func (client PrivateCloudsClient) ListAdminCredentialsPreparer(ctx context.Context, resourceGroupName string, privateCloudName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"privateCloudName":  autorest.Encode("path", privateCloudName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-03-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}/listAdminCredentials", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListAdminCredentialsSender sends the ListAdminCredentials request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateCloudsClient) ListAdminCredentialsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListAdminCredentialsResponder handles the response to the ListAdminCredentials request. The method always
// closes the http.Response Body.
func (client PrivateCloudsClient) ListAdminCredentialsResponder(resp *http.Response) (result AdminCredentials, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListInSubscription sends the list in subscription request.
func (client PrivateCloudsClient) ListInSubscription(ctx context.Context) (result PrivateCloudListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.ListInSubscription")
		defer func() {
			sc := -1
			if result.pcl.Response.Response != nil {
				sc = result.pcl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("avs.PrivateCloudsClient", "ListInSubscription", err.Error())
	}

	result.fn = client.listInSubscriptionNextResults
	req, err := client.ListInSubscriptionPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "ListInSubscription", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListInSubscriptionSender(req)
	if err != nil {
		result.pcl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "ListInSubscription", resp, "Failure sending request")
		return
	}

	result.pcl, err = client.ListInSubscriptionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "ListInSubscription", resp, "Failure responding to request")
	}

	return
}

// ListInSubscriptionPreparer prepares the ListInSubscription request.
func (client PrivateCloudsClient) ListInSubscriptionPreparer(ctx context.Context) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-03-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.AVS/privateClouds", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListInSubscriptionSender sends the ListInSubscription request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateCloudsClient) ListInSubscriptionSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListInSubscriptionResponder handles the response to the ListInSubscription request. The method always
// closes the http.Response Body.
func (client PrivateCloudsClient) ListInSubscriptionResponder(resp *http.Response) (result PrivateCloudList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listInSubscriptionNextResults retrieves the next set of results, if any.
func (client PrivateCloudsClient) listInSubscriptionNextResults(ctx context.Context, lastResults PrivateCloudList) (result PrivateCloudList, err error) {
	req, err := lastResults.privateCloudListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "listInSubscriptionNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListInSubscriptionSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "listInSubscriptionNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListInSubscriptionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "listInSubscriptionNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListInSubscriptionComplete enumerates all values, automatically crossing page boundaries as required.
func (client PrivateCloudsClient) ListInSubscriptionComplete(ctx context.Context) (result PrivateCloudListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.ListInSubscription")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListInSubscription(ctx)
	return
}

// Update sends the update request.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// privateCloudName - name of the private cloud
// privateCloudUpdate - the private cloud properties to be updated
func (client PrivateCloudsClient) Update(ctx context.Context, resourceGroupName string, privateCloudName string, privateCloudUpdate PrivateCloudUpdate) (result PrivateCloudsUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateCloudsClient.Update")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("avs.PrivateCloudsClient", "Update", err.Error())
	}

	req, err := client.UpdatePreparer(ctx, resourceGroupName, privateCloudName, privateCloudUpdate)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "avs.PrivateCloudsClient", "Update", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client PrivateCloudsClient) UpdatePreparer(ctx context.Context, resourceGroupName string, privateCloudName string, privateCloudUpdate PrivateCloudUpdate) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"privateCloudName":  autorest.Encode("path", privateCloudName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-03-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}", pathParameters),
		autorest.WithJSON(privateCloudUpdate),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateCloudsClient) UpdateSender(req *http.Request) (future PrivateCloudsUpdateFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client PrivateCloudsClient) UpdateResponder(resp *http.Response) (result PrivateCloud, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
