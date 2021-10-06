package network

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
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

// ServiceEndpointPoliciesClient is the network Client
type ServiceEndpointPoliciesClient struct {
	BaseClient
}

// NewServiceEndpointPoliciesClient creates an instance of the ServiceEndpointPoliciesClient client.
func NewServiceEndpointPoliciesClient(subscriptionID string) ServiceEndpointPoliciesClient {
	return NewServiceEndpointPoliciesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewServiceEndpointPoliciesClientWithBaseURI creates an instance of the ServiceEndpointPoliciesClient client using a
// custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds,
// Azure stack).
func NewServiceEndpointPoliciesClientWithBaseURI(baseURI string, subscriptionID string) ServiceEndpointPoliciesClient {
	return ServiceEndpointPoliciesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates a service Endpoint Policies.
// Parameters:
// resourceGroupName - the name of the resource group.
// serviceEndpointPolicyName - the name of the service endpoint policy.
// parameters - parameters supplied to the create or update service endpoint policy operation.
func (client ServiceEndpointPoliciesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, serviceEndpointPolicyName string, parameters ServiceEndpointPolicy) (result ServiceEndpointPoliciesCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceEndpointPoliciesClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, serviceEndpointPolicyName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client ServiceEndpointPoliciesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, serviceEndpointPolicyName string, parameters ServiceEndpointPolicy) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":         autorest.Encode("path", resourceGroupName),
		"serviceEndpointPolicyName": autorest.Encode("path", serviceEndpointPolicyName),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	parameters.Etag = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/serviceEndpointPolicies/{serviceEndpointPolicyName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceEndpointPoliciesClient) CreateOrUpdateSender(req *http.Request) (future ServiceEndpointPoliciesCreateOrUpdateFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client ServiceEndpointPoliciesClient) CreateOrUpdateResponder(resp *http.Response) (result ServiceEndpointPolicy, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes the specified service endpoint policy.
// Parameters:
// resourceGroupName - the name of the resource group.
// serviceEndpointPolicyName - the name of the service endpoint policy.
func (client ServiceEndpointPoliciesClient) Delete(ctx context.Context, resourceGroupName string, serviceEndpointPolicyName string) (result ServiceEndpointPoliciesDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceEndpointPoliciesClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, serviceEndpointPolicyName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client ServiceEndpointPoliciesClient) DeletePreparer(ctx context.Context, resourceGroupName string, serviceEndpointPolicyName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":         autorest.Encode("path", resourceGroupName),
		"serviceEndpointPolicyName": autorest.Encode("path", serviceEndpointPolicyName),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/serviceEndpointPolicies/{serviceEndpointPolicyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceEndpointPoliciesClient) DeleteSender(req *http.Request) (future ServiceEndpointPoliciesDeleteFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ServiceEndpointPoliciesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the specified service Endpoint Policies in a specified resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
// serviceEndpointPolicyName - the name of the service endpoint policy.
// expand - expands referenced resources.
func (client ServiceEndpointPoliciesClient) Get(ctx context.Context, resourceGroupName string, serviceEndpointPolicyName string, expand string) (result ServiceEndpointPolicy, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceEndpointPoliciesClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, serviceEndpointPolicyName, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client ServiceEndpointPoliciesClient) GetPreparer(ctx context.Context, resourceGroupName string, serviceEndpointPolicyName string, expand string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":         autorest.Encode("path", resourceGroupName),
		"serviceEndpointPolicyName": autorest.Encode("path", serviceEndpointPolicyName),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(expand) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/serviceEndpointPolicies/{serviceEndpointPolicyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceEndpointPoliciesClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ServiceEndpointPoliciesClient) GetResponder(resp *http.Response) (result ServiceEndpointPolicy, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets all the service endpoint policies in a subscription.
func (client ServiceEndpointPoliciesClient) List(ctx context.Context) (result ServiceEndpointPolicyListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceEndpointPoliciesClient.List")
		defer func() {
			sc := -1
			if result.seplr.Response.Response != nil {
				sc = result.seplr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.seplr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "List", resp, "Failure sending request")
		return
	}

	result.seplr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "List", resp, "Failure responding to request")
		return
	}
	if result.seplr.hasNextLink() && result.seplr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client ServiceEndpointPoliciesClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Network/ServiceEndpointPolicies", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceEndpointPoliciesClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ServiceEndpointPoliciesClient) ListResponder(resp *http.Response) (result ServiceEndpointPolicyListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client ServiceEndpointPoliciesClient) listNextResults(ctx context.Context, lastResults ServiceEndpointPolicyListResult) (result ServiceEndpointPolicyListResult, err error) {
	req, err := lastResults.serviceEndpointPolicyListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client ServiceEndpointPoliciesClient) ListComplete(ctx context.Context) (result ServiceEndpointPolicyListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceEndpointPoliciesClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx)
	return
}

// ListByResourceGroup gets all service endpoint Policies in a resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
func (client ServiceEndpointPoliciesClient) ListByResourceGroup(ctx context.Context, resourceGroupName string) (result ServiceEndpointPolicyListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceEndpointPoliciesClient.ListByResourceGroup")
		defer func() {
			sc := -1
			if result.seplr.Response.Response != nil {
				sc = result.seplr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listByResourceGroupNextResults
	req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "ListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.seplr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "ListByResourceGroup", resp, "Failure sending request")
		return
	}

	result.seplr, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "ListByResourceGroup", resp, "Failure responding to request")
		return
	}
	if result.seplr.hasNextLink() && result.seplr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListByResourceGroupPreparer prepares the ListByResourceGroup request.
func (client ServiceEndpointPoliciesClient) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/serviceEndpointPolicies", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceEndpointPoliciesClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client ServiceEndpointPoliciesClient) ListByResourceGroupResponder(resp *http.Response) (result ServiceEndpointPolicyListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByResourceGroupNextResults retrieves the next set of results, if any.
func (client ServiceEndpointPoliciesClient) listByResourceGroupNextResults(ctx context.Context, lastResults ServiceEndpointPolicyListResult) (result ServiceEndpointPolicyListResult, err error) {
	req, err := lastResults.serviceEndpointPolicyListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "listByResourceGroupNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "listByResourceGroupNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByResourceGroupComplete enumerates all values, automatically crossing page boundaries as required.
func (client ServiceEndpointPoliciesClient) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result ServiceEndpointPolicyListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceEndpointPoliciesClient.ListByResourceGroup")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByResourceGroup(ctx, resourceGroupName)
	return
}

// UpdateTags updates tags of a service endpoint policy.
// Parameters:
// resourceGroupName - the name of the resource group.
// serviceEndpointPolicyName - the name of the service endpoint policy.
// parameters - parameters supplied to update service endpoint policy tags.
func (client ServiceEndpointPoliciesClient) UpdateTags(ctx context.Context, resourceGroupName string, serviceEndpointPolicyName string, parameters TagsObject) (result ServiceEndpointPolicy, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceEndpointPoliciesClient.UpdateTags")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdateTagsPreparer(ctx, resourceGroupName, serviceEndpointPolicyName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "UpdateTags", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateTagsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "UpdateTags", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateTagsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.ServiceEndpointPoliciesClient", "UpdateTags", resp, "Failure responding to request")
		return
	}

	return
}

// UpdateTagsPreparer prepares the UpdateTags request.
func (client ServiceEndpointPoliciesClient) UpdateTagsPreparer(ctx context.Context, resourceGroupName string, serviceEndpointPolicyName string, parameters TagsObject) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":         autorest.Encode("path", resourceGroupName),
		"serviceEndpointPolicyName": autorest.Encode("path", serviceEndpointPolicyName),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/serviceEndpointPolicies/{serviceEndpointPolicyName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateTagsSender sends the UpdateTags request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceEndpointPoliciesClient) UpdateTagsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// UpdateTagsResponder handles the response to the UpdateTags request. The method always
// closes the http.Response Body.
func (client ServiceEndpointPoliciesClient) UpdateTagsResponder(resp *http.Response) (result ServiceEndpointPolicy, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
