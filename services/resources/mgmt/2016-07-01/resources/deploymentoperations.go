package resources

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

// DeploymentOperationsClient is the client for the DeploymentOperations methods of the Resources service.
type DeploymentOperationsClient struct {
    BaseClient
}
// NewDeploymentOperationsClient creates an instance of the DeploymentOperationsClient client.
func NewDeploymentOperationsClient(subscriptionID string) DeploymentOperationsClient {
    return NewDeploymentOperationsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDeploymentOperationsClientWithBaseURI creates an instance of the DeploymentOperationsClient client.
    func NewDeploymentOperationsClientWithBaseURI(baseURI string, subscriptionID string) DeploymentOperationsClient {
        return DeploymentOperationsClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// Get get a list of deployments operations.
    // Parameters:
        // resourceGroupName - the name of the resource group. The name is case insensitive.
        // deploymentName - the name of the deployment.
        // operationID - operation Id.
func (client DeploymentOperationsClient) Get(ctx context.Context, resourceGroupName string, deploymentName string, operationID string) (result DeploymentOperation, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DeploymentOperationsClient.Get")
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
            return result, validation.NewError("resources.DeploymentOperationsClient", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, deploymentName, operationID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client DeploymentOperationsClient) GetPreparer(ctx context.Context, resourceGroupName string, deploymentName string, operationID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "deploymentName": autorest.Encode("path",deploymentName),
            "operationId": autorest.Encode("path",operationID),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-07-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/deployments/{deploymentName}/operations/{operationId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client DeploymentOperationsClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DeploymentOperationsClient) GetResponder(resp *http.Response) (result DeploymentOperation, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List gets a list of deployments operations.
    // Parameters:
        // resourceGroupName - the name of the resource group. The name is case insensitive.
        // deploymentName - the name of the deployment.
        // top - query parameters.
func (client DeploymentOperationsClient) List(ctx context.Context, resourceGroupName string, deploymentName string, top *int32) (result DeploymentOperationsListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DeploymentOperationsClient.List")
        defer func() {
            sc := -1
            if result.dolr.Response.Response != nil {
                sc = result.dolr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("resources.DeploymentOperationsClient", "List", err.Error())
            }

                        result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx, resourceGroupName, deploymentName, top)
    if err != nil {
    err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.dolr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "List", resp, "Failure sending request")
            return
            }

            result.dolr, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client DeploymentOperationsClient) ListPreparer(ctx context.Context, resourceGroupName string, deploymentName string, top *int32) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "deploymentName": autorest.Encode("path",deploymentName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-07-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if top != nil {
            queryParameters["$top"] = autorest.Encode("query",*top)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/deployments/{deploymentName}/operations",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client DeploymentOperationsClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client DeploymentOperationsClient) ListResponder(resp *http.Response) (result DeploymentOperationsListResult, err error) {
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
            func (client DeploymentOperationsClient) listNextResults(ctx context.Context, lastResults DeploymentOperationsListResult) (result DeploymentOperationsListResult, err error) {
            req, err := lastResults.deploymentOperationsListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "resources.DeploymentOperationsClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client DeploymentOperationsClient) ListComplete(ctx context.Context, resourceGroupName string, deploymentName string, top *int32) (result DeploymentOperationsListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/DeploymentOperationsClient.List")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.List(ctx, resourceGroupName, deploymentName, top)
                return
        }

