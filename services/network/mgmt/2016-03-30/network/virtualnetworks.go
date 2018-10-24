package network

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
    "github.com/Azure/go-autorest/autorest"
    "github.com/Azure/go-autorest/autorest/azure"
    "net/http"
    "context"
    "github.com/Azure/go-autorest/tracing"
)

// VirtualNetworksClient is the network Client
type VirtualNetworksClient struct {
    BaseClient
}
// NewVirtualNetworksClient creates an instance of the VirtualNetworksClient client.
func NewVirtualNetworksClient(subscriptionID string) VirtualNetworksClient {
    return NewVirtualNetworksClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewVirtualNetworksClientWithBaseURI creates an instance of the VirtualNetworksClient client.
    func NewVirtualNetworksClientWithBaseURI(baseURI string, subscriptionID string) VirtualNetworksClient {
        return VirtualNetworksClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate the Put VirtualNetwork operation creates/updates a virtual network in the specified resource group.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // virtualNetworkName - the name of the virtual network.
        // parameters - parameters supplied to the create/update Virtual Network operation
func (client VirtualNetworksClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, virtualNetworkName string, parameters VirtualNetwork) (result VirtualNetworksCreateOrUpdateFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualNetworksClient.CreateOrUpdate")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, virtualNetworkName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            result, err = client.CreateOrUpdateSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "CreateOrUpdate", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client VirtualNetworksClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, virtualNetworkName string, parameters VirtualNetwork) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "virtualNetworkName": autorest.Encode("path",virtualNetworkName),
            }

                        const APIVersion = "2016-03-30"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{virtualNetworkName}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualNetworksClient) CreateOrUpdateSender(req *http.Request) (future VirtualNetworksCreateOrUpdateFuture, err error) {
            var resp *http.Response
            resp, err = autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            if err != nil {
            return
            }
            future.Future, err = azure.NewFutureFromResponse(resp)
            return
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client VirtualNetworksClient) CreateOrUpdateResponder(resp *http.Response) (result VirtualNetwork, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete the Delete VirtualNetwork operation deletes the specifed virtual network
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // virtualNetworkName - the name of the virtual network.
func (client VirtualNetworksClient) Delete(ctx context.Context, resourceGroupName string, virtualNetworkName string) (result VirtualNetworksDeleteFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualNetworksClient.Delete")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.DeletePreparer(ctx, resourceGroupName, virtualNetworkName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "Delete", nil , "Failure preparing request")
    return
    }

            result, err = client.DeleteSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "Delete", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client VirtualNetworksClient) DeletePreparer(ctx context.Context, resourceGroupName string, virtualNetworkName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "virtualNetworkName": autorest.Encode("path",virtualNetworkName),
            }

                        const APIVersion = "2016-03-30"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{virtualNetworkName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualNetworksClient) DeleteSender(req *http.Request) (future VirtualNetworksDeleteFuture, err error) {
            var resp *http.Response
            resp, err = autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            if err != nil {
            return
            }
            future.Future, err = azure.NewFutureFromResponse(resp)
            return
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client VirtualNetworksClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get the Get VirtualNetwork operation retrieves information about the specified virtual network.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // virtualNetworkName - the name of the virtual network.
        // expand - expand references resources.
func (client VirtualNetworksClient) Get(ctx context.Context, resourceGroupName string, virtualNetworkName string, expand string) (result VirtualNetwork, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualNetworksClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetPreparer(ctx, resourceGroupName, virtualNetworkName, expand)
    if err != nil {
    err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client VirtualNetworksClient) GetPreparer(ctx context.Context, resourceGroupName string, virtualNetworkName string, expand string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "virtualNetworkName": autorest.Encode("path",virtualNetworkName),
            }

                        const APIVersion = "2016-03-30"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if len(expand) > 0 {
            queryParameters["$expand"] = autorest.Encode("query",expand)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{virtualNetworkName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualNetworksClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client VirtualNetworksClient) GetResponder(resp *http.Response) (result VirtualNetwork, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List the list VirtualNetwork returns all Virtual Networks in a resource group
    // Parameters:
        // resourceGroupName - the name of the resource group.
func (client VirtualNetworksClient) List(ctx context.Context, resourceGroupName string) (result VirtualNetworkListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualNetworksClient.List")
        defer func() {
            sc := -1
            if result.vnlr.Response.Response != nil {
                sc = result.vnlr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx, resourceGroupName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.vnlr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "List", resp, "Failure sending request")
            return
            }

            result.vnlr, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client VirtualNetworksClient) ListPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-03-30"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualNetworksClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client VirtualNetworksClient) ListResponder(resp *http.Response) (result VirtualNetworkListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listNextResults retrieves the next set of results, if any.
            func (client VirtualNetworksClient) listNextResults(ctx context.Context, lastResults VirtualNetworkListResult) (result VirtualNetworkListResult, err error) {
            req, err := lastResults.virtualNetworkListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client VirtualNetworksClient) ListComplete(ctx context.Context, resourceGroupName string) (result VirtualNetworkListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/VirtualNetworksClient.List")
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

// ListAll the list VirtualNetwork returns all Virtual Networks in a subscription
func (client VirtualNetworksClient) ListAll(ctx context.Context) (result VirtualNetworkListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualNetworksClient.ListAll")
        defer func() {
            sc := -1
            if result.vnlr.Response.Response != nil {
                sc = result.vnlr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listAllNextResults
    req, err := client.ListAllPreparer(ctx)
    if err != nil {
    err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "ListAll", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListAllSender(req)
            if err != nil {
            result.vnlr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "ListAll", resp, "Failure sending request")
            return
            }

            result.vnlr, err = client.ListAllResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "ListAll", resp, "Failure responding to request")
            }

    return
    }

    // ListAllPreparer prepares the ListAll request.
    func (client VirtualNetworksClient) ListAllPreparer(ctx context.Context) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-03-30"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Network/virtualNetworks",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListAllSender sends the ListAll request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualNetworksClient) ListAllSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListAllResponder handles the response to the ListAll request. The method always
// closes the http.Response Body.
func (client VirtualNetworksClient) ListAllResponder(resp *http.Response) (result VirtualNetworkListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listAllNextResults retrieves the next set of results, if any.
            func (client VirtualNetworksClient) listAllNextResults(ctx context.Context, lastResults VirtualNetworkListResult) (result VirtualNetworkListResult, err error) {
            req, err := lastResults.virtualNetworkListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "listAllNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListAllSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "listAllNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListAllResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.VirtualNetworksClient", "listAllNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListAllComplete enumerates all values, automatically crossing page boundaries as required.
    func (client VirtualNetworksClient) ListAllComplete(ctx context.Context) (result VirtualNetworkListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/VirtualNetworksClient.ListAll")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.ListAll(ctx)
                return
        }

