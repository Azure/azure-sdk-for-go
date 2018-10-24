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
)

// WorkflowVersionTriggersClient is the REST API for Azure Logic Apps.
type WorkflowVersionTriggersClient struct {
    BaseClient
}
// NewWorkflowVersionTriggersClient creates an instance of the WorkflowVersionTriggersClient client.
func NewWorkflowVersionTriggersClient(subscriptionID string) WorkflowVersionTriggersClient {
    return NewWorkflowVersionTriggersClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewWorkflowVersionTriggersClientWithBaseURI creates an instance of the WorkflowVersionTriggersClient client.
    func NewWorkflowVersionTriggersClientWithBaseURI(baseURI string, subscriptionID string) WorkflowVersionTriggersClient {
        return WorkflowVersionTriggersClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// ListCallbackURL get the callback url for a trigger of a workflow version.
    // Parameters:
        // resourceGroupName - the resource group name.
        // workflowName - the workflow name.
        // versionID - the workflow versionId.
        // triggerName - the workflow trigger name.
        // parameters - the callback URL parameters.
func (client WorkflowVersionTriggersClient) ListCallbackURL(ctx context.Context, resourceGroupName string, workflowName string, versionID string, triggerName string, parameters *GetCallbackURLParameters) (result WorkflowTriggerCallbackURL, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/WorkflowVersionTriggersClient.ListCallbackURL")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.ListCallbackURLPreparer(ctx, resourceGroupName, workflowName, versionID, triggerName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "logic.WorkflowVersionTriggersClient", "ListCallbackURL", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListCallbackURLSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "logic.WorkflowVersionTriggersClient", "ListCallbackURL", resp, "Failure sending request")
            return
            }

            result, err = client.ListCallbackURLResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "logic.WorkflowVersionTriggersClient", "ListCallbackURL", resp, "Failure responding to request")
            }

    return
    }

    // ListCallbackURLPreparer prepares the ListCallbackURL request.
    func (client WorkflowVersionTriggersClient) ListCallbackURLPreparer(ctx context.Context, resourceGroupName string, workflowName string, versionID string, triggerName string, parameters *GetCallbackURLParameters) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "triggerName": autorest.Encode("path",triggerName),
            "versionId": autorest.Encode("path",versionID),
            "workflowName": autorest.Encode("path",workflowName),
            }

                        const APIVersion = "2018-07-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPost(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/versions/{versionId}/triggers/{triggerName}/listCallbackUrl",pathParameters),
    autorest.WithQueryParameters(queryParameters))
            if parameters != nil {
            preparer = autorest.DecoratePreparer(preparer,
            autorest.WithJSON(parameters))
            }
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListCallbackURLSender sends the ListCallbackURL request. The method will close the
    // http.Response Body if it receives an error.
    func (client WorkflowVersionTriggersClient) ListCallbackURLSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListCallbackURLResponder handles the response to the ListCallbackURL request. The method always
// closes the http.Response Body.
func (client WorkflowVersionTriggersClient) ListCallbackURLResponder(resp *http.Response) (result WorkflowTriggerCallbackURL, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

