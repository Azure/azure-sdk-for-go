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

// VariableClient is the automation Client
type VariableClient struct {
    BaseClient
}
// NewVariableClient creates an instance of the VariableClient client.
func NewVariableClient(subscriptionID string, countType1 CountType) VariableClient {
    return NewVariableClientWithBaseURI(DefaultBaseURI, subscriptionID, countType1)
}

// NewVariableClientWithBaseURI creates an instance of the VariableClient client.
    func NewVariableClientWithBaseURI(baseURI string, subscriptionID string, countType1 CountType) VariableClient {
        return VariableClient{ NewWithBaseURI(baseURI, subscriptionID, countType1)}
    }

// CreateOrUpdate create a variable.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
        // variableName - the variable name.
        // parameters - the parameters supplied to the create or update variable operation.
func (client VariableClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, automationAccountName string, variableName string, parameters VariableCreateOrUpdateParameters) (result Variable, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VariableClient.CreateOrUpdate")
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
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil }}},
            { TargetValue: parameters,
             Constraints: []validation.Constraint{	{Target: "parameters.Name", Name: validation.Null, Rule: true, Chain: nil },
            	{Target: "parameters.VariableCreateOrUpdateProperties", Name: validation.Null, Rule: true, Chain: nil }}}}); err != nil {
            return result, validation.NewError("automation.VariableClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, automationAccountName, variableName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.VariableClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "CreateOrUpdate", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "CreateOrUpdate", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client VariableClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, automationAccountName string, variableName string, parameters VariableCreateOrUpdateParameters) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "automationAccountName": autorest.Encode("path",automationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "variableName": autorest.Encode("path",variableName),
            }

                        const APIVersion = "2015-10-31"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/variables/{variableName}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client VariableClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client VariableClient) CreateOrUpdateResponder(resp *http.Response) (result Variable, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete delete the variable.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
        // variableName - the name of variable.
func (client VariableClient) Delete(ctx context.Context, resourceGroupName string, automationAccountName string, variableName string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VariableClient.Delete")
        defer func() {
            sc := -1
            if result.Response != nil {
                sc = result.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("automation.VariableClient", "Delete", err.Error())
            }

                req, err := client.DeletePreparer(ctx, resourceGroupName, automationAccountName, variableName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.VariableClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client VariableClient) DeletePreparer(ctx context.Context, resourceGroupName string, automationAccountName string, variableName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "automationAccountName": autorest.Encode("path",automationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "variableName": autorest.Encode("path",variableName),
            }

                        const APIVersion = "2015-10-31"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/variables/{variableName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client VariableClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client VariableClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get retrieve the variable identified by variable name.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
        // variableName - the name of variable.
func (client VariableClient) Get(ctx context.Context, resourceGroupName string, automationAccountName string, variableName string) (result Variable, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VariableClient.Get")
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
            return result, validation.NewError("automation.VariableClient", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, automationAccountName, variableName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.VariableClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client VariableClient) GetPreparer(ctx context.Context, resourceGroupName string, automationAccountName string, variableName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "automationAccountName": autorest.Encode("path",automationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "variableName": autorest.Encode("path",variableName),
            }

                        const APIVersion = "2015-10-31"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/variables/{variableName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client VariableClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client VariableClient) GetResponder(resp *http.Response) (result Variable, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// ListByAutomationAccount retrieve a list of variables.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
func (client VariableClient) ListByAutomationAccount(ctx context.Context, resourceGroupName string, automationAccountName string) (result VariableListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VariableClient.ListByAutomationAccount")
        defer func() {
            sc := -1
            if result.vlr.Response.Response != nil {
                sc = result.vlr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("automation.VariableClient", "ListByAutomationAccount", err.Error())
            }

                        result.fn = client.listByAutomationAccountNextResults
    req, err := client.ListByAutomationAccountPreparer(ctx, resourceGroupName, automationAccountName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.VariableClient", "ListByAutomationAccount", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByAutomationAccountSender(req)
            if err != nil {
            result.vlr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "ListByAutomationAccount", resp, "Failure sending request")
            return
            }

            result.vlr, err = client.ListByAutomationAccountResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "ListByAutomationAccount", resp, "Failure responding to request")
            }

    return
    }

    // ListByAutomationAccountPreparer prepares the ListByAutomationAccount request.
    func (client VariableClient) ListByAutomationAccountPreparer(ctx context.Context, resourceGroupName string, automationAccountName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "automationAccountName": autorest.Encode("path",automationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-10-31"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/variables",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByAutomationAccountSender sends the ListByAutomationAccount request. The method will close the
    // http.Response Body if it receives an error.
    func (client VariableClient) ListByAutomationAccountSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByAutomationAccountResponder handles the response to the ListByAutomationAccount request. The method always
// closes the http.Response Body.
func (client VariableClient) ListByAutomationAccountResponder(resp *http.Response) (result VariableListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listByAutomationAccountNextResults retrieves the next set of results, if any.
            func (client VariableClient) listByAutomationAccountNextResults(ctx context.Context, lastResults VariableListResult) (result VariableListResult, err error) {
            req, err := lastResults.variableListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "automation.VariableClient", "listByAutomationAccountNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListByAutomationAccountSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "automation.VariableClient", "listByAutomationAccountNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListByAutomationAccountResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "listByAutomationAccountNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListByAutomationAccountComplete enumerates all values, automatically crossing page boundaries as required.
    func (client VariableClient) ListByAutomationAccountComplete(ctx context.Context, resourceGroupName string, automationAccountName string) (result VariableListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/VariableClient.ListByAutomationAccount")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.ListByAutomationAccount(ctx, resourceGroupName, automationAccountName)
                return
        }

// Update update a variable.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
        // variableName - the variable name.
        // parameters - the parameters supplied to the update variable operation.
func (client VariableClient) Update(ctx context.Context, resourceGroupName string, automationAccountName string, variableName string, parameters VariableUpdateParameters) (result Variable, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VariableClient.Update")
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
            return result, validation.NewError("automation.VariableClient", "Update", err.Error())
            }

                req, err := client.UpdatePreparer(ctx, resourceGroupName, automationAccountName, variableName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.VariableClient", "Update", nil , "Failure preparing request")
    return
    }

            resp, err := client.UpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "Update", resp, "Failure sending request")
            return
            }

            result, err = client.UpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.VariableClient", "Update", resp, "Failure responding to request")
            }

    return
    }

    // UpdatePreparer prepares the Update request.
    func (client VariableClient) UpdatePreparer(ctx context.Context, resourceGroupName string, automationAccountName string, variableName string, parameters VariableUpdateParameters) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "automationAccountName": autorest.Encode("path",automationAccountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "variableName": autorest.Encode("path",variableName),
            }

                        const APIVersion = "2015-10-31"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPatch(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/variables/{variableName}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // UpdateSender sends the Update request. The method will close the
    // http.Response Body if it receives an error.
    func (client VariableClient) UpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client VariableClient) UpdateResponder(resp *http.Response) (result Variable, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

