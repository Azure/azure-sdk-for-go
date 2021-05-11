package machinelearningservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
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

// OnlineEndpointsClient is the these APIs allow end users to operate on Azure Machine Learning Workspace resources.
type OnlineEndpointsClient struct {
	BaseClient
}

// NewOnlineEndpointsClient creates an instance of the OnlineEndpointsClient client.
func NewOnlineEndpointsClient(subscriptionID string) OnlineEndpointsClient {
	return NewOnlineEndpointsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewOnlineEndpointsClientWithBaseURI creates an instance of the OnlineEndpointsClient client using a custom endpoint.
// Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewOnlineEndpointsClientWithBaseURI(baseURI string, subscriptionID string) OnlineEndpointsClient {
	return OnlineEndpointsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate sends the create or update request.
// Parameters:
// endpointName - online Endpoint name.
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - name of Azure Machine Learning workspace.
// body - online Endpoint entity to apply during operation.
func (client OnlineEndpointsClient) CreateOrUpdate(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string, body OnlineEndpointTrackedResource) (result OnlineEndpointsCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: endpointName,
			Constraints: []validation.Constraint{{Target: "endpointName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9][a-zA-Z0-9\-_]{0,254}$`, Chain: nil}}},
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: body,
			Constraints: []validation.Constraint{{Target: "body.Properties", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("machinelearningservices.OnlineEndpointsClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, endpointName, resourceGroupName, workspaceName, body)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "CreateOrUpdate", nil, "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client OnlineEndpointsClient) CreateOrUpdatePreparer(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string, body OnlineEndpointTrackedResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"endpointName":      autorest.Encode("path", endpointName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	body.SystemData = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{workspaceName}/onlineEndpoints/{endpointName}", pathParameters),
		autorest.WithJSON(body),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client OnlineEndpointsClient) CreateOrUpdateSender(req *http.Request) (future OnlineEndpointsCreateOrUpdateFuture, err error) {
	var resp *http.Response
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
func (client OnlineEndpointsClient) CreateOrUpdateResponder(resp *http.Response) (result OnlineEndpointTrackedResource, err error) {
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
// endpointName - online Endpoint name.
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - name of Azure Machine Learning workspace.
func (client OnlineEndpointsClient) Delete(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string) (result OnlineEndpointsDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
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
		return result, validation.NewError("machinelearningservices.OnlineEndpointsClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, endpointName, resourceGroupName, workspaceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "Delete", nil, "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client OnlineEndpointsClient) DeletePreparer(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"endpointName":      autorest.Encode("path", endpointName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{workspaceName}/onlineEndpoints/{endpointName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client OnlineEndpointsClient) DeleteSender(req *http.Request) (future OnlineEndpointsDeleteFuture, err error) {
	var resp *http.Response
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
func (client OnlineEndpointsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get sends the get request.
// Parameters:
// endpointName - online Endpoint name.
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - name of Azure Machine Learning workspace.
func (client OnlineEndpointsClient) Get(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string) (result OnlineEndpointTrackedResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.Get")
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
		return result, validation.NewError("machinelearningservices.OnlineEndpointsClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, endpointName, resourceGroupName, workspaceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client OnlineEndpointsClient) GetPreparer(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"endpointName":      autorest.Encode("path", endpointName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{workspaceName}/onlineEndpoints/{endpointName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client OnlineEndpointsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client OnlineEndpointsClient) GetResponder(resp *http.Response) (result OnlineEndpointTrackedResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetToken sends the get token request.
// Parameters:
// endpointName - online Endpoint name.
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - name of Azure Machine Learning workspace.
func (client OnlineEndpointsClient) GetToken(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string) (result EndpointAuthToken, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.GetToken")
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
		return result, validation.NewError("machinelearningservices.OnlineEndpointsClient", "GetToken", err.Error())
	}

	req, err := client.GetTokenPreparer(ctx, endpointName, resourceGroupName, workspaceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "GetToken", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetTokenSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "GetToken", resp, "Failure sending request")
		return
	}

	result, err = client.GetTokenResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "GetToken", resp, "Failure responding to request")
		return
	}

	return
}

// GetTokenPreparer prepares the GetToken request.
func (client OnlineEndpointsClient) GetTokenPreparer(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"endpointName":      autorest.Encode("path", endpointName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{workspaceName}/onlineEndpoints/{endpointName}/token", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetTokenSender sends the GetToken request. The method will close the
// http.Response Body if it receives an error.
func (client OnlineEndpointsClient) GetTokenSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetTokenResponder handles the response to the GetToken request. The method always
// closes the http.Response Body.
func (client OnlineEndpointsClient) GetTokenResponder(resp *http.Response) (result EndpointAuthToken, err error) {
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
// workspaceName - name of Azure Machine Learning workspace.
// name - name of the endpoint.
// count - number of endpoints to be retrieved in a page of results.
// computeType - endpointComputeType to be filtered by.
// skip - continuation token for pagination.
// tags - a set of tags with which to filter the returned models. It is a comma separated string of tags key or
// tags key=value. Example: tagKey1,tagKey2,tagKey3=value3 .
// properties - a set of properties with which to filter the returned models. It is a comma separated string of
// properties key and/or properties key=value Example: propKey1,propKey2,propKey3=value3 .
// orderBy - the option to order the response.
func (client OnlineEndpointsClient) List(ctx context.Context, resourceGroupName string, workspaceName string, name string, count *int32, computeType EndpointComputeType, skip string, tags string, properties string, orderBy OrderString) (result OnlineEndpointTrackedResourceArmPaginatedResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.List")
		defer func() {
			sc := -1
			if result.oetrapr.Response.Response != nil {
				sc = result.oetrapr.Response.Response.StatusCode
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
		return result, validation.NewError("machinelearningservices.OnlineEndpointsClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, workspaceName, name, count, computeType, skip, tags, properties, orderBy)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.oetrapr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "List", resp, "Failure sending request")
		return
	}

	result.oetrapr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "List", resp, "Failure responding to request")
		return
	}
	if result.oetrapr.hasNextLink() && result.oetrapr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client OnlineEndpointsClient) ListPreparer(ctx context.Context, resourceGroupName string, workspaceName string, name string, count *int32, computeType EndpointComputeType, skip string, tags string, properties string, orderBy OrderString) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(name) > 0 {
		queryParameters["name"] = autorest.Encode("query", name)
	}
	if count != nil {
		queryParameters["count"] = autorest.Encode("query", *count)
	}
	if len(string(computeType)) > 0 {
		queryParameters["computeType"] = autorest.Encode("query", computeType)
	}
	if len(skip) > 0 {
		queryParameters["$skip"] = autorest.Encode("query", skip)
	}
	if len(tags) > 0 {
		queryParameters["tags"] = autorest.Encode("query", tags)
	}
	if len(properties) > 0 {
		queryParameters["properties"] = autorest.Encode("query", properties)
	}
	if len(string(orderBy)) > 0 {
		queryParameters["orderBy"] = autorest.Encode("query", orderBy)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{workspaceName}/onlineEndpoints", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client OnlineEndpointsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client OnlineEndpointsClient) ListResponder(resp *http.Response) (result OnlineEndpointTrackedResourceArmPaginatedResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client OnlineEndpointsClient) listNextResults(ctx context.Context, lastResults OnlineEndpointTrackedResourceArmPaginatedResult) (result OnlineEndpointTrackedResourceArmPaginatedResult, err error) {
	req, err := lastResults.onlineEndpointTrackedResourceArmPaginatedResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client OnlineEndpointsClient) ListComplete(ctx context.Context, resourceGroupName string, workspaceName string, name string, count *int32, computeType EndpointComputeType, skip string, tags string, properties string, orderBy OrderString) (result OnlineEndpointTrackedResourceArmPaginatedResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, resourceGroupName, workspaceName, name, count, computeType, skip, tags, properties, orderBy)
	return
}

// ListKeys sends the list keys request.
// Parameters:
// endpointName - online Endpoint name.
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - name of Azure Machine Learning workspace.
func (client OnlineEndpointsClient) ListKeys(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string) (result EndpointAuthKeys, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.ListKeys")
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
		return result, validation.NewError("machinelearningservices.OnlineEndpointsClient", "ListKeys", err.Error())
	}

	req, err := client.ListKeysPreparer(ctx, endpointName, resourceGroupName, workspaceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "ListKeys", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListKeysSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "ListKeys", resp, "Failure sending request")
		return
	}

	result, err = client.ListKeysResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "ListKeys", resp, "Failure responding to request")
		return
	}

	return
}

// ListKeysPreparer prepares the ListKeys request.
func (client OnlineEndpointsClient) ListKeysPreparer(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"endpointName":      autorest.Encode("path", endpointName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{workspaceName}/onlineEndpoints/{endpointName}/listKeys", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListKeysSender sends the ListKeys request. The method will close the
// http.Response Body if it receives an error.
func (client OnlineEndpointsClient) ListKeysSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListKeysResponder handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (client OnlineEndpointsClient) ListKeysResponder(resp *http.Response) (result EndpointAuthKeys, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// RegenerateKeys sends the regenerate keys request.
// Parameters:
// endpointName - online Endpoint name.
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - name of Azure Machine Learning workspace.
// body - regenerateKeys request .
func (client OnlineEndpointsClient) RegenerateKeys(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string, body RegenerateEndpointKeysRequest) (result OnlineEndpointsRegenerateKeysFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.RegenerateKeys")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
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
		return result, validation.NewError("machinelearningservices.OnlineEndpointsClient", "RegenerateKeys", err.Error())
	}

	req, err := client.RegenerateKeysPreparer(ctx, endpointName, resourceGroupName, workspaceName, body)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "RegenerateKeys", nil, "Failure preparing request")
		return
	}

	result, err = client.RegenerateKeysSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "RegenerateKeys", nil, "Failure sending request")
		return
	}

	return
}

// RegenerateKeysPreparer prepares the RegenerateKeys request.
func (client OnlineEndpointsClient) RegenerateKeysPreparer(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string, body RegenerateEndpointKeysRequest) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"endpointName":      autorest.Encode("path", endpointName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{workspaceName}/onlineEndpoints/{endpointName}/regenerateKeys", pathParameters),
		autorest.WithJSON(body),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RegenerateKeysSender sends the RegenerateKeys request. The method will close the
// http.Response Body if it receives an error.
func (client OnlineEndpointsClient) RegenerateKeysSender(req *http.Request) (future OnlineEndpointsRegenerateKeysFuture, err error) {
	var resp *http.Response
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

// RegenerateKeysResponder handles the response to the RegenerateKeys request. The method always
// closes the http.Response Body.
func (client OnlineEndpointsClient) RegenerateKeysResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Update sends the update request.
// Parameters:
// endpointName - online Endpoint name.
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - name of Azure Machine Learning workspace.
// body - online Endpoint entity to apply during operation.
func (client OnlineEndpointsClient) Update(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string, body PartialOnlineEndpointPartialTrackedResource) (result OnlineEndpointsUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/OnlineEndpointsClient.Update")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
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
		return result, validation.NewError("machinelearningservices.OnlineEndpointsClient", "Update", err.Error())
	}

	req, err := client.UpdatePreparer(ctx, endpointName, resourceGroupName, workspaceName, body)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningservices.OnlineEndpointsClient", "Update", nil, "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client OnlineEndpointsClient) UpdatePreparer(ctx context.Context, endpointName string, resourceGroupName string, workspaceName string, body PartialOnlineEndpointPartialTrackedResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"endpointName":      autorest.Encode("path", endpointName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{workspaceName}/onlineEndpoints/{endpointName}", pathParameters),
		autorest.WithJSON(body),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client OnlineEndpointsClient) UpdateSender(req *http.Request) (future OnlineEndpointsUpdateFuture, err error) {
	var resp *http.Response
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

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client OnlineEndpointsClient) UpdateResponder(resp *http.Response) (result OnlineEndpointTrackedResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
