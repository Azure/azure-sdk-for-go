package cdn

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
    "github.com/Azure/go-autorest/autorest/validation"
)

// OriginsClient is the use these APIs to manage Azure CDN resources through the Azure Resource Manager. You must make
// sure that requests made to these resources are secure.
type OriginsClient struct {
    BaseClient
}
// NewOriginsClient creates an instance of the OriginsClient client.
func NewOriginsClient(subscriptionID string) OriginsClient {
    return NewOriginsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewOriginsClientWithBaseURI creates an instance of the OriginsClient client.
    func NewOriginsClientWithBaseURI(baseURI string, subscriptionID string) OriginsClient {
        return OriginsClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// Get gets an existing origin within an endpoint.
    // Parameters:
        // resourceGroupName - name of the Resource group within the Azure subscription.
        // profileName - name of the CDN profile which is unique within the resource group.
        // endpointName - name of the endpoint under the profile which is unique globally.
        // originName - name of the origin which is unique within the endpoint.
func (client OriginsClient) Get(ctx context.Context, resourceGroupName string, profileName string, endpointName string, originName string) (result Origin, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/OriginsClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("cdn.OriginsClient", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, profileName, endpointName, originName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client OriginsClient) GetPreparer(ctx context.Context, resourceGroupName string, profileName string, endpointName string, originName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "endpointName": autorest.Encode("path",endpointName),
            "originName": autorest.Encode("path",originName),
            "profileName": autorest.Encode("path",profileName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-12"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/endpoints/{endpointName}/origins/{originName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client OriginsClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client OriginsClient) GetResponder(resp *http.Response) (result Origin, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// ListByEndpoint lists all of the existing origins within an endpoint.
    // Parameters:
        // resourceGroupName - name of the Resource group within the Azure subscription.
        // profileName - name of the CDN profile which is unique within the resource group.
        // endpointName - name of the endpoint under the profile which is unique globally.
func (client OriginsClient) ListByEndpoint(ctx context.Context, resourceGroupName string, profileName string, endpointName string) (result OriginListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/OriginsClient.ListByEndpoint")
        defer func() {
            sc := -1
            if result.olr.Response.Response != nil {
                sc = result.olr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("cdn.OriginsClient", "ListByEndpoint", err.Error())
            }

                        result.fn = client.listByEndpointNextResults
    req, err := client.ListByEndpointPreparer(ctx, resourceGroupName, profileName, endpointName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "ListByEndpoint", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByEndpointSender(req)
            if err != nil {
            result.olr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "ListByEndpoint", resp, "Failure sending request")
            return
            }

            result.olr, err = client.ListByEndpointResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "ListByEndpoint", resp, "Failure responding to request")
            }

    return
    }

    // ListByEndpointPreparer prepares the ListByEndpoint request.
    func (client OriginsClient) ListByEndpointPreparer(ctx context.Context, resourceGroupName string, profileName string, endpointName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "endpointName": autorest.Encode("path",endpointName),
            "profileName": autorest.Encode("path",profileName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-12"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/endpoints/{endpointName}/origins",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByEndpointSender sends the ListByEndpoint request. The method will close the
    // http.Response Body if it receives an error.
    func (client OriginsClient) ListByEndpointSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByEndpointResponder handles the response to the ListByEndpoint request. The method always
// closes the http.Response Body.
func (client OriginsClient) ListByEndpointResponder(resp *http.Response) (result OriginListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listByEndpointNextResults retrieves the next set of results, if any.
            func (client OriginsClient) listByEndpointNextResults(ctx context.Context, lastResults OriginListResult) (result OriginListResult, err error) {
            req, err := lastResults.originListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "cdn.OriginsClient", "listByEndpointNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListByEndpointSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "cdn.OriginsClient", "listByEndpointNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListByEndpointResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "listByEndpointNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListByEndpointComplete enumerates all values, automatically crossing page boundaries as required.
    func (client OriginsClient) ListByEndpointComplete(ctx context.Context, resourceGroupName string, profileName string, endpointName string) (result OriginListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/OriginsClient.ListByEndpoint")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.ListByEndpoint(ctx, resourceGroupName, profileName, endpointName)
                return
        }

// Update updates an existing origin within an endpoint.
    // Parameters:
        // resourceGroupName - name of the Resource group within the Azure subscription.
        // profileName - name of the CDN profile which is unique within the resource group.
        // endpointName - name of the endpoint under the profile which is unique globally.
        // originName - name of the origin which is unique within the endpoint.
        // originUpdateProperties - origin properties
func (client OriginsClient) Update(ctx context.Context, resourceGroupName string, profileName string, endpointName string, originName string, originUpdateProperties OriginUpdateParameters) (result OriginsUpdateFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/OriginsClient.Update")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("cdn.OriginsClient", "Update", err.Error())
            }

                req, err := client.UpdatePreparer(ctx, resourceGroupName, profileName, endpointName, originName, originUpdateProperties)
    if err != nil {
    err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "Update", nil , "Failure preparing request")
    return
    }

            result, err = client.UpdateSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "Update", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // UpdatePreparer prepares the Update request.
    func (client OriginsClient) UpdatePreparer(ctx context.Context, resourceGroupName string, profileName string, endpointName string, originName string, originUpdateProperties OriginUpdateParameters) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "endpointName": autorest.Encode("path",endpointName),
            "originName": autorest.Encode("path",originName),
            "profileName": autorest.Encode("path",profileName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-12"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPatch(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/endpoints/{endpointName}/origins/{originName}",pathParameters),
    autorest.WithJSON(originUpdateProperties),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // UpdateSender sends the Update request. The method will close the
    // http.Response Body if it receives an error.
    func (client OriginsClient) UpdateSender(req *http.Request) (future OriginsUpdateFuture, err error) {
            var resp *http.Response
            resp, err = autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            if err != nil {
            return
            }
            future.Future, err = azure.NewFutureFromResponse(resp)
            return
            }

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client OriginsClient) UpdateResponder(resp *http.Response) (result Origin, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

