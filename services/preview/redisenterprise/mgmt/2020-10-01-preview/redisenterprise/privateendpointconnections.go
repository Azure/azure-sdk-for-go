package redisenterprise

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

// PrivateEndpointConnectionsClient is the REST API for managing Redis Enterprise resources in Azure.
type PrivateEndpointConnectionsClient struct {
	BaseClient
}

// NewPrivateEndpointConnectionsClient creates an instance of the PrivateEndpointConnectionsClient client.
func NewPrivateEndpointConnectionsClient(subscriptionID string) PrivateEndpointConnectionsClient {
	return NewPrivateEndpointConnectionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewPrivateEndpointConnectionsClientWithBaseURI creates an instance of the PrivateEndpointConnectionsClient client
// using a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign
// clouds, Azure stack).
func NewPrivateEndpointConnectionsClientWithBaseURI(baseURI string, subscriptionID string) PrivateEndpointConnectionsClient {
	return PrivateEndpointConnectionsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Delete deletes the specified private endpoint connection associated with the RedisEnterprise cluster.
// Parameters:
// resourceGroupName - the name of the resource group.
// clusterName - the name of the RedisEnterprise cluster.
// privateEndpointConnectionName - the name of the private endpoint connection associated with the Azure
// resource
func (client PrivateEndpointConnectionsClient) Delete(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateEndpointConnectionsClient.Delete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, clusterName, privateEndpointConnectionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "Delete", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client PrivateEndpointConnectionsClient) DeletePreparer(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clusterName":                   autorest.Encode("path", clusterName),
		"privateEndpointConnectionName": autorest.Encode("path", privateEndpointConnectionName),
		"resourceGroupName":             autorest.Encode("path", resourceGroupName),
		"subscriptionId":                autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-10-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redisEnterprise/{clusterName}/privateEndpointConnections/{privateEndpointConnectionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateEndpointConnectionsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client PrivateEndpointConnectionsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// GetMethod gets the specified private endpoint connection associated with the RedisEnterprise cluster.
// Parameters:
// resourceGroupName - the name of the resource group.
// clusterName - the name of the RedisEnterprise cluster.
// privateEndpointConnectionName - the name of the private endpoint connection associated with the Azure
// resource
func (client PrivateEndpointConnectionsClient) GetMethod(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string) (result PrivateEndpointConnection, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateEndpointConnectionsClient.GetMethod")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetMethodPreparer(ctx, resourceGroupName, clusterName, privateEndpointConnectionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "GetMethod", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetMethodSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "GetMethod", resp, "Failure sending request")
		return
	}

	result, err = client.GetMethodResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "GetMethod", resp, "Failure responding to request")
		return
	}

	return
}

// GetMethodPreparer prepares the GetMethod request.
func (client PrivateEndpointConnectionsClient) GetMethodPreparer(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clusterName":                   autorest.Encode("path", clusterName),
		"privateEndpointConnectionName": autorest.Encode("path", privateEndpointConnectionName),
		"resourceGroupName":             autorest.Encode("path", resourceGroupName),
		"subscriptionId":                autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-10-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redisEnterprise/{clusterName}/privateEndpointConnections/{privateEndpointConnectionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetMethodSender sends the GetMethod request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateEndpointConnectionsClient) GetMethodSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetMethodResponder handles the response to the GetMethod request. The method always
// closes the http.Response Body.
func (client PrivateEndpointConnectionsClient) GetMethodResponder(resp *http.Response) (result PrivateEndpointConnection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List lists all the private endpoint connections associated with the RedisEnterprise cluster.
// Parameters:
// resourceGroupName - the name of the resource group.
// clusterName - the name of the RedisEnterprise cluster.
func (client PrivateEndpointConnectionsClient) List(ctx context.Context, resourceGroupName string, clusterName string) (result PrivateEndpointConnectionListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateEndpointConnectionsClient.List")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListPreparer(ctx, resourceGroupName, clusterName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "List", resp, "Failure responding to request")
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client PrivateEndpointConnectionsClient) ListPreparer(ctx context.Context, resourceGroupName string, clusterName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clusterName":       autorest.Encode("path", clusterName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-10-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redisEnterprise/{clusterName}/privateEndpointConnections", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateEndpointConnectionsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client PrivateEndpointConnectionsClient) ListResponder(resp *http.Response) (result PrivateEndpointConnectionListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Put updates the state of the specified private endpoint connection associated with the RedisEnterprise cluster.
// Parameters:
// resourceGroupName - the name of the resource group.
// clusterName - the name of the RedisEnterprise cluster.
// privateEndpointConnectionName - the name of the private endpoint connection associated with the Azure
// resource
// properties - the private endpoint connection properties.
func (client PrivateEndpointConnectionsClient) Put(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string, properties PrivateEndpointConnection) (result PrivateEndpointConnectionsPutFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateEndpointConnectionsClient.Put")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: properties,
			Constraints: []validation.Constraint{{Target: "properties.PrivateEndpointConnectionProperties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "properties.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("redisenterprise.PrivateEndpointConnectionsClient", "Put", err.Error())
	}

	req, err := client.PutPreparer(ctx, resourceGroupName, clusterName, privateEndpointConnectionName, properties)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "Put", nil, "Failure preparing request")
		return
	}

	result, err = client.PutSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.PrivateEndpointConnectionsClient", "Put", result.Response(), "Failure sending request")
		return
	}

	return
}

// PutPreparer prepares the Put request.
func (client PrivateEndpointConnectionsClient) PutPreparer(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string, properties PrivateEndpointConnection) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"clusterName":                   autorest.Encode("path", clusterName),
		"privateEndpointConnectionName": autorest.Encode("path", privateEndpointConnectionName),
		"resourceGroupName":             autorest.Encode("path", resourceGroupName),
		"subscriptionId":                autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-10-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cache/redisEnterprise/{clusterName}/privateEndpointConnections/{privateEndpointConnectionName}", pathParameters),
		autorest.WithJSON(properties),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PutSender sends the Put request. The method will close the
// http.Response Body if it receives an error.
func (client PrivateEndpointConnectionsClient) PutSender(req *http.Request) (future PrivateEndpointConnectionsPutFuture, err error) {
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

// PutResponder handles the response to the Put request. The method always
// closes the http.Response Body.
func (client PrivateEndpointConnectionsClient) PutResponder(resp *http.Response) (result PrivateEndpointConnection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
