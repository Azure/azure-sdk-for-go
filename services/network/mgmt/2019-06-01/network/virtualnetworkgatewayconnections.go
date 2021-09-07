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
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// VirtualNetworkGatewayConnectionsClient is the network Client
type VirtualNetworkGatewayConnectionsClient struct {
	BaseClient
}

// NewVirtualNetworkGatewayConnectionsClient creates an instance of the VirtualNetworkGatewayConnectionsClient client.
func NewVirtualNetworkGatewayConnectionsClient(subscriptionID string) VirtualNetworkGatewayConnectionsClient {
	return NewVirtualNetworkGatewayConnectionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewVirtualNetworkGatewayConnectionsClientWithBaseURI creates an instance of the
// VirtualNetworkGatewayConnectionsClient client using a custom endpoint.  Use this when interacting with an Azure
// cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewVirtualNetworkGatewayConnectionsClientWithBaseURI(baseURI string, subscriptionID string) VirtualNetworkGatewayConnectionsClient {
	return VirtualNetworkGatewayConnectionsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates a virtual network gateway connection in the specified resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
// virtualNetworkGatewayConnectionName - the name of the virtual network gateway connection.
// parameters - parameters supplied to the create or update virtual network gateway connection operation.
func (client VirtualNetworkGatewayConnectionsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string, parameters VirtualNetworkGatewayConnection) (result VirtualNetworkGatewayConnectionsCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.VirtualNetworkGatewayConnectionPropertiesFormat", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.VirtualNetworkGatewayConnectionPropertiesFormat.VirtualNetworkGateway1", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "parameters.VirtualNetworkGatewayConnectionPropertiesFormat.VirtualNetworkGateway1.VirtualNetworkGatewayPropertiesFormat", Name: validation.Null, Rule: true, Chain: nil}}},
					{Target: "parameters.VirtualNetworkGatewayConnectionPropertiesFormat.VirtualNetworkGateway2", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.VirtualNetworkGatewayConnectionPropertiesFormat.VirtualNetworkGateway2.VirtualNetworkGatewayPropertiesFormat", Name: validation.Null, Rule: true, Chain: nil}}},
					{Target: "parameters.VirtualNetworkGatewayConnectionPropertiesFormat.LocalNetworkGateway2", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.VirtualNetworkGatewayConnectionPropertiesFormat.LocalNetworkGateway2.LocalNetworkGatewayPropertiesFormat", Name: validation.Null, Rule: true, Chain: nil}}},
				}}}}}); err != nil {
		return result, validation.NewError("network.VirtualNetworkGatewayConnectionsClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, virtualNetworkGatewayConnectionName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client VirtualNetworkGatewayConnectionsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string, parameters VirtualNetworkGatewayConnection) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":                   autorest.Encode("path", resourceGroupName),
		"subscriptionId":                      autorest.Encode("path", client.SubscriptionID),
		"virtualNetworkGatewayConnectionName": autorest.Encode("path", virtualNetworkGatewayConnectionName),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/connections/{virtualNetworkGatewayConnectionName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualNetworkGatewayConnectionsClient) CreateOrUpdateSender(req *http.Request) (future VirtualNetworkGatewayConnectionsCreateOrUpdateFuture, err error) {
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
func (client VirtualNetworkGatewayConnectionsClient) CreateOrUpdateResponder(resp *http.Response) (result VirtualNetworkGatewayConnection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes the specified virtual network Gateway connection.
// Parameters:
// resourceGroupName - the name of the resource group.
// virtualNetworkGatewayConnectionName - the name of the virtual network gateway connection.
func (client VirtualNetworkGatewayConnectionsClient) Delete(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string) (result VirtualNetworkGatewayConnectionsDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, virtualNetworkGatewayConnectionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client VirtualNetworkGatewayConnectionsClient) DeletePreparer(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":                   autorest.Encode("path", resourceGroupName),
		"subscriptionId":                      autorest.Encode("path", client.SubscriptionID),
		"virtualNetworkGatewayConnectionName": autorest.Encode("path", virtualNetworkGatewayConnectionName),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/connections/{virtualNetworkGatewayConnectionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualNetworkGatewayConnectionsClient) DeleteSender(req *http.Request) (future VirtualNetworkGatewayConnectionsDeleteFuture, err error) {
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
func (client VirtualNetworkGatewayConnectionsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the specified virtual network gateway connection by resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
// virtualNetworkGatewayConnectionName - the name of the virtual network gateway connection.
func (client VirtualNetworkGatewayConnectionsClient) Get(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string) (result VirtualNetworkGatewayConnection, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, virtualNetworkGatewayConnectionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client VirtualNetworkGatewayConnectionsClient) GetPreparer(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":                   autorest.Encode("path", resourceGroupName),
		"subscriptionId":                      autorest.Encode("path", client.SubscriptionID),
		"virtualNetworkGatewayConnectionName": autorest.Encode("path", virtualNetworkGatewayConnectionName),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/connections/{virtualNetworkGatewayConnectionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualNetworkGatewayConnectionsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client VirtualNetworkGatewayConnectionsClient) GetResponder(resp *http.Response) (result VirtualNetworkGatewayConnection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetSharedKey the Get VirtualNetworkGatewayConnectionSharedKey operation retrieves information about the specified
// virtual network gateway connection shared key through Network resource provider.
// Parameters:
// resourceGroupName - the name of the resource group.
// virtualNetworkGatewayConnectionName - the virtual network gateway connection shared key name.
func (client VirtualNetworkGatewayConnectionsClient) GetSharedKey(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string) (result ConnectionSharedKey, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.GetSharedKey")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetSharedKeyPreparer(ctx, resourceGroupName, virtualNetworkGatewayConnectionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "GetSharedKey", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSharedKeySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "GetSharedKey", resp, "Failure sending request")
		return
	}

	result, err = client.GetSharedKeyResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "GetSharedKey", resp, "Failure responding to request")
		return
	}

	return
}

// GetSharedKeyPreparer prepares the GetSharedKey request.
func (client VirtualNetworkGatewayConnectionsClient) GetSharedKeyPreparer(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":                   autorest.Encode("path", resourceGroupName),
		"subscriptionId":                      autorest.Encode("path", client.SubscriptionID),
		"virtualNetworkGatewayConnectionName": autorest.Encode("path", virtualNetworkGatewayConnectionName),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/connections/{virtualNetworkGatewayConnectionName}/sharedkey", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSharedKeySender sends the GetSharedKey request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualNetworkGatewayConnectionsClient) GetSharedKeySender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetSharedKeyResponder handles the response to the GetSharedKey request. The method always
// closes the http.Response Body.
func (client VirtualNetworkGatewayConnectionsClient) GetSharedKeyResponder(resp *http.Response) (result ConnectionSharedKey, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List the List VirtualNetworkGatewayConnections operation retrieves all the virtual network gateways connections
// created.
// Parameters:
// resourceGroupName - the name of the resource group.
func (client VirtualNetworkGatewayConnectionsClient) List(ctx context.Context, resourceGroupName string) (result VirtualNetworkGatewayConnectionListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.List")
		defer func() {
			sc := -1
			if result.vngclr.Response.Response != nil {
				sc = result.vngclr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.vngclr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "List", resp, "Failure sending request")
		return
	}

	result.vngclr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "List", resp, "Failure responding to request")
		return
	}
	if result.vngclr.hasNextLink() && result.vngclr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client VirtualNetworkGatewayConnectionsClient) ListPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/connections", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualNetworkGatewayConnectionsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client VirtualNetworkGatewayConnectionsClient) ListResponder(resp *http.Response) (result VirtualNetworkGatewayConnectionListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client VirtualNetworkGatewayConnectionsClient) listNextResults(ctx context.Context, lastResults VirtualNetworkGatewayConnectionListResult) (result VirtualNetworkGatewayConnectionListResult, err error) {
	req, err := lastResults.virtualNetworkGatewayConnectionListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client VirtualNetworkGatewayConnectionsClient) ListComplete(ctx context.Context, resourceGroupName string) (result VirtualNetworkGatewayConnectionListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.List")
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

// ResetSharedKey the VirtualNetworkGatewayConnectionResetSharedKey operation resets the virtual network gateway
// connection shared key for passed virtual network gateway connection in the specified resource group through Network
// resource provider.
// Parameters:
// resourceGroupName - the name of the resource group.
// virtualNetworkGatewayConnectionName - the virtual network gateway connection reset shared key Name.
// parameters - parameters supplied to the begin reset virtual network gateway connection shared key operation
// through network resource provider.
func (client VirtualNetworkGatewayConnectionsClient) ResetSharedKey(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string, parameters ConnectionResetSharedKey) (result VirtualNetworkGatewayConnectionsResetSharedKeyFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.ResetSharedKey")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.KeyLength", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.KeyLength", Name: validation.InclusiveMaximum, Rule: int64(128), Chain: nil},
					{Target: "parameters.KeyLength", Name: validation.InclusiveMinimum, Rule: int64(1), Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("network.VirtualNetworkGatewayConnectionsClient", "ResetSharedKey", err.Error())
	}

	req, err := client.ResetSharedKeyPreparer(ctx, resourceGroupName, virtualNetworkGatewayConnectionName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "ResetSharedKey", nil, "Failure preparing request")
		return
	}

	result, err = client.ResetSharedKeySender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "ResetSharedKey", result.Response(), "Failure sending request")
		return
	}

	return
}

// ResetSharedKeyPreparer prepares the ResetSharedKey request.
func (client VirtualNetworkGatewayConnectionsClient) ResetSharedKeyPreparer(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string, parameters ConnectionResetSharedKey) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":                   autorest.Encode("path", resourceGroupName),
		"subscriptionId":                      autorest.Encode("path", client.SubscriptionID),
		"virtualNetworkGatewayConnectionName": autorest.Encode("path", virtualNetworkGatewayConnectionName),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/connections/{virtualNetworkGatewayConnectionName}/sharedkey/reset", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ResetSharedKeySender sends the ResetSharedKey request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualNetworkGatewayConnectionsClient) ResetSharedKeySender(req *http.Request) (future VirtualNetworkGatewayConnectionsResetSharedKeyFuture, err error) {
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

// ResetSharedKeyResponder handles the response to the ResetSharedKey request. The method always
// closes the http.Response Body.
func (client VirtualNetworkGatewayConnectionsClient) ResetSharedKeyResponder(resp *http.Response) (result ConnectionResetSharedKey, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// SetSharedKey the Put VirtualNetworkGatewayConnectionSharedKey operation sets the virtual network gateway connection
// shared key for passed virtual network gateway connection in the specified resource group through Network resource
// provider.
// Parameters:
// resourceGroupName - the name of the resource group.
// virtualNetworkGatewayConnectionName - the virtual network gateway connection name.
// parameters - parameters supplied to the Begin Set Virtual Network Gateway connection Shared key operation
// throughNetwork resource provider.
func (client VirtualNetworkGatewayConnectionsClient) SetSharedKey(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string, parameters ConnectionSharedKey) (result VirtualNetworkGatewayConnectionsSetSharedKeyFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.SetSharedKey")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Value", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("network.VirtualNetworkGatewayConnectionsClient", "SetSharedKey", err.Error())
	}

	req, err := client.SetSharedKeyPreparer(ctx, resourceGroupName, virtualNetworkGatewayConnectionName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "SetSharedKey", nil, "Failure preparing request")
		return
	}

	result, err = client.SetSharedKeySender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "SetSharedKey", result.Response(), "Failure sending request")
		return
	}

	return
}

// SetSharedKeyPreparer prepares the SetSharedKey request.
func (client VirtualNetworkGatewayConnectionsClient) SetSharedKeyPreparer(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string, parameters ConnectionSharedKey) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":                   autorest.Encode("path", resourceGroupName),
		"subscriptionId":                      autorest.Encode("path", client.SubscriptionID),
		"virtualNetworkGatewayConnectionName": autorest.Encode("path", virtualNetworkGatewayConnectionName),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/connections/{virtualNetworkGatewayConnectionName}/sharedkey", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetSharedKeySender sends the SetSharedKey request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualNetworkGatewayConnectionsClient) SetSharedKeySender(req *http.Request) (future VirtualNetworkGatewayConnectionsSetSharedKeyFuture, err error) {
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

// SetSharedKeyResponder handles the response to the SetSharedKey request. The method always
// closes the http.Response Body.
func (client VirtualNetworkGatewayConnectionsClient) SetSharedKeyResponder(resp *http.Response) (result ConnectionSharedKey, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// UpdateTags updates a virtual network gateway connection tags.
// Parameters:
// resourceGroupName - the name of the resource group.
// virtualNetworkGatewayConnectionName - the name of the virtual network gateway connection.
// parameters - parameters supplied to update virtual network gateway connection tags.
func (client VirtualNetworkGatewayConnectionsClient) UpdateTags(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string, parameters TagsObject) (result VirtualNetworkGatewayConnectionsUpdateTagsFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualNetworkGatewayConnectionsClient.UpdateTags")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdateTagsPreparer(ctx, resourceGroupName, virtualNetworkGatewayConnectionName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "UpdateTags", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateTagsSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewayConnectionsClient", "UpdateTags", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdateTagsPreparer prepares the UpdateTags request.
func (client VirtualNetworkGatewayConnectionsClient) UpdateTagsPreparer(ctx context.Context, resourceGroupName string, virtualNetworkGatewayConnectionName string, parameters TagsObject) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":                   autorest.Encode("path", resourceGroupName),
		"subscriptionId":                      autorest.Encode("path", client.SubscriptionID),
		"virtualNetworkGatewayConnectionName": autorest.Encode("path", virtualNetworkGatewayConnectionName),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/connections/{virtualNetworkGatewayConnectionName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateTagsSender sends the UpdateTags request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualNetworkGatewayConnectionsClient) UpdateTagsSender(req *http.Request) (future VirtualNetworkGatewayConnectionsUpdateTagsFuture, err error) {
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

// UpdateTagsResponder handles the response to the UpdateTags request. The method always
// closes the http.Response Body.
func (client VirtualNetworkGatewayConnectionsClient) UpdateTagsResponder(resp *http.Response) (result VirtualNetworkGatewayConnection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
