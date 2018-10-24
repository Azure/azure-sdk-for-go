package logic

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

// IntegrationAccountSessionsClient is the REST API for Azure Logic Apps.
type IntegrationAccountSessionsClient struct {
    BaseClient
}
// NewIntegrationAccountSessionsClient creates an instance of the IntegrationAccountSessionsClient client.
func NewIntegrationAccountSessionsClient(subscriptionID string) IntegrationAccountSessionsClient {
    return NewIntegrationAccountSessionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewIntegrationAccountSessionsClientWithBaseURI creates an instance of the IntegrationAccountSessionsClient client.
    func NewIntegrationAccountSessionsClientWithBaseURI(baseURI string, subscriptionID string) IntegrationAccountSessionsClient {
        return IntegrationAccountSessionsClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate creates or updates an integration account session.
    // Parameters:
        // resourceGroupName - the resource group name.
        // integrationAccountName - the integration account name.
        // sessionName - the integration account session name.
        // session - the integration account session.
func (client IntegrationAccountSessionsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, integrationAccountName string, sessionName string, session IntegrationAccountSession) (result IntegrationAccountSession, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/IntegrationAccountSessionsClient.CreateOrUpdate")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: session,
             Constraints: []validation.Constraint{	{Target: "session.IntegrationAccountSessionProperties", Name: validation.Null, Rule: true, Chain: nil }}}}); err != nil {
            return result, validation.NewError("logic.IntegrationAccountSessionsClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, integrationAccountName, sessionName, session)
    if err != nil {
    err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "CreateOrUpdate", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "CreateOrUpdate", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client IntegrationAccountSessionsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, integrationAccountName string, sessionName string, session IntegrationAccountSession) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "integrationAccountName": autorest.Encode("path",integrationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "sessionName": autorest.Encode("path",sessionName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-07-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/sessions/{sessionName}",pathParameters),
    autorest.WithJSON(session),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client IntegrationAccountSessionsClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client IntegrationAccountSessionsClient) CreateOrUpdateResponder(resp *http.Response) (result IntegrationAccountSession, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete deletes an integration account session.
    // Parameters:
        // resourceGroupName - the resource group name.
        // integrationAccountName - the integration account name.
        // sessionName - the integration account session name.
func (client IntegrationAccountSessionsClient) Delete(ctx context.Context, resourceGroupName string, integrationAccountName string, sessionName string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/IntegrationAccountSessionsClient.Delete")
        defer func() {
            sc := -1
            if result.Response != nil {
                sc = result.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.DeletePreparer(ctx, resourceGroupName, integrationAccountName, sessionName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client IntegrationAccountSessionsClient) DeletePreparer(ctx context.Context, resourceGroupName string, integrationAccountName string, sessionName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "integrationAccountName": autorest.Encode("path",integrationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "sessionName": autorest.Encode("path",sessionName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-07-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/sessions/{sessionName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client IntegrationAccountSessionsClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client IntegrationAccountSessionsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get gets an integration account session.
    // Parameters:
        // resourceGroupName - the resource group name.
        // integrationAccountName - the integration account name.
        // sessionName - the integration account session name.
func (client IntegrationAccountSessionsClient) Get(ctx context.Context, resourceGroupName string, integrationAccountName string, sessionName string) (result IntegrationAccountSession, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/IntegrationAccountSessionsClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetPreparer(ctx, resourceGroupName, integrationAccountName, sessionName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client IntegrationAccountSessionsClient) GetPreparer(ctx context.Context, resourceGroupName string, integrationAccountName string, sessionName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "integrationAccountName": autorest.Encode("path",integrationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "sessionName": autorest.Encode("path",sessionName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-07-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/sessions/{sessionName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client IntegrationAccountSessionsClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client IntegrationAccountSessionsClient) GetResponder(resp *http.Response) (result IntegrationAccountSession, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List gets a list of integration account sessions.
    // Parameters:
        // resourceGroupName - the resource group name.
        // integrationAccountName - the integration account name.
        // top - the number of items to be included in the result.
        // filter - the filter to apply on the operation. Options for filters include: ChangedTime.
func (client IntegrationAccountSessionsClient) List(ctx context.Context, resourceGroupName string, integrationAccountName string, top *int32, filter string) (result IntegrationAccountSessionListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/IntegrationAccountSessionsClient.List")
        defer func() {
            sc := -1
            if result.iaslr.Response.Response != nil {
                sc = result.iaslr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx, resourceGroupName, integrationAccountName, top, filter)
    if err != nil {
    err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.iaslr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "List", resp, "Failure sending request")
            return
            }

            result.iaslr, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client IntegrationAccountSessionsClient) ListPreparer(ctx context.Context, resourceGroupName string, integrationAccountName string, top *int32, filter string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "integrationAccountName": autorest.Encode("path",integrationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-07-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if top != nil {
            queryParameters["$top"] = autorest.Encode("query",*top)
            }
            if len(filter) > 0 {
            queryParameters["$filter"] = autorest.Encode("query",filter)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/sessions",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client IntegrationAccountSessionsClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client IntegrationAccountSessionsClient) ListResponder(resp *http.Response) (result IntegrationAccountSessionListResult, err error) {
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
            func (client IntegrationAccountSessionsClient) listNextResults(ctx context.Context, lastResults IntegrationAccountSessionListResult) (result IntegrationAccountSessionListResult, err error) {
            req, err := lastResults.integrationAccountSessionListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "logic.IntegrationAccountSessionsClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client IntegrationAccountSessionsClient) ListComplete(ctx context.Context, resourceGroupName string, integrationAccountName string, top *int32, filter string) (result IntegrationAccountSessionListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/IntegrationAccountSessionsClient.List")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.List(ctx, resourceGroupName, integrationAccountName, top, filter)
                return
        }

