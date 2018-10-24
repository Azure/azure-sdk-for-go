package automation

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

// StatisticsClient is the automation Client
type StatisticsClient struct {
    BaseClient
}
// NewStatisticsClient creates an instance of the StatisticsClient client.
func NewStatisticsClient(subscriptionID string, countType1 CountType) StatisticsClient {
    return NewStatisticsClientWithBaseURI(DefaultBaseURI, subscriptionID, countType1)
}

// NewStatisticsClientWithBaseURI creates an instance of the StatisticsClient client.
    func NewStatisticsClientWithBaseURI(baseURI string, subscriptionID string, countType1 CountType) StatisticsClient {
        return StatisticsClient{ NewWithBaseURI(baseURI, subscriptionID, countType1)}
    }

// ListByAutomationAccount retrieve the statistics for the account.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
        // filter - the filter to apply on the operation.
func (client StatisticsClient) ListByAutomationAccount(ctx context.Context, resourceGroupName string, automationAccountName string, filter string) (result StatisticsListResult, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/StatisticsClient.ListByAutomationAccount")
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
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("automation.StatisticsClient", "ListByAutomationAccount", err.Error())
            }

                req, err := client.ListByAutomationAccountPreparer(ctx, resourceGroupName, automationAccountName, filter)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.StatisticsClient", "ListByAutomationAccount", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByAutomationAccountSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.StatisticsClient", "ListByAutomationAccount", resp, "Failure sending request")
            return
            }

            result, err = client.ListByAutomationAccountResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.StatisticsClient", "ListByAutomationAccount", resp, "Failure responding to request")
            }

    return
    }

    // ListByAutomationAccountPreparer prepares the ListByAutomationAccount request.
    func (client StatisticsClient) ListByAutomationAccountPreparer(ctx context.Context, resourceGroupName string, automationAccountName string, filter string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "automationAccountName": autorest.Encode("path",automationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-10-31"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if len(filter) > 0 {
            queryParameters["$filter"] = autorest.Encode("query",filter)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/statistics",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByAutomationAccountSender sends the ListByAutomationAccount request. The method will close the
    // http.Response Body if it receives an error.
    func (client StatisticsClient) ListByAutomationAccountSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByAutomationAccountResponder handles the response to the ListByAutomationAccount request. The method always
// closes the http.Response Body.
func (client StatisticsClient) ListByAutomationAccountResponder(resp *http.Response) (result StatisticsListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

