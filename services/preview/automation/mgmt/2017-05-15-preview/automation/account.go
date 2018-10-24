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

// AccountClient is the automation Client
type AccountClient struct {
    BaseClient
}
// NewAccountClient creates an instance of the AccountClient client.
func NewAccountClient(subscriptionID string) AccountClient {
    return NewAccountClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewAccountClientWithBaseURI creates an instance of the AccountClient client.
    func NewAccountClientWithBaseURI(baseURI string, subscriptionID string) AccountClient {
        return AccountClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate create or update automation account.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
        // parameters - parameters supplied to the create or update automation account.
func (client AccountClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, automationAccountName string, parameters AccountCreateOrUpdateParameters) (result Account, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountClient.CreateOrUpdate")
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
            return result, validation.NewError("automation.AccountClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, automationAccountName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.AccountClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "CreateOrUpdate", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "CreateOrUpdate", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client AccountClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, automationAccountName string, parameters AccountCreateOrUpdateParameters) (*http.Request, error) {
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
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client AccountClient) CreateOrUpdateResponder(resp *http.Response) (result Account, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete delete an automation account.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
func (client AccountClient) Delete(ctx context.Context, resourceGroupName string, automationAccountName string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountClient.Delete")
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
            return result, validation.NewError("automation.AccountClient", "Delete", err.Error())
            }

                req, err := client.DeletePreparer(ctx, resourceGroupName, automationAccountName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.AccountClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client AccountClient) DeletePreparer(ctx context.Context, resourceGroupName string, automationAccountName string) (*http.Request, error) {
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
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client AccountClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get get information about an Automation Account.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
func (client AccountClient) Get(ctx context.Context, resourceGroupName string, automationAccountName string) (result Account, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountClient.Get")
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
            return result, validation.NewError("automation.AccountClient", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, automationAccountName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.AccountClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client AccountClient) GetPreparer(ctx context.Context, resourceGroupName string, automationAccountName string) (*http.Request, error) {
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
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client AccountClient) GetResponder(resp *http.Response) (result Account, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List retrieve a list of accounts within a given subscription.
func (client AccountClient) List(ctx context.Context) (result AccountListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountClient.List")
        defer func() {
            sc := -1
            if result.alr.Response.Response != nil {
                sc = result.alr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.AccountClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.alr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "List", resp, "Failure sending request")
            return
            }

            result.alr, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client AccountClient) ListPreparer(ctx context.Context) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-10-31"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Automation/automationAccounts",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client AccountClient) ListResponder(resp *http.Response) (result AccountListResult, err error) {
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
            func (client AccountClient) listNextResults(ctx context.Context, lastResults AccountListResult) (result AccountListResult, err error) {
            req, err := lastResults.accountListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "automation.AccountClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "automation.AccountClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client AccountClient) ListComplete(ctx context.Context) (result AccountListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/AccountClient.List")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.List(ctx)
                return
        }

// ListByResourceGroup retrieve a list of accounts within a given resource group.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
func (client AccountClient) ListByResourceGroup(ctx context.Context, resourceGroupName string) (result AccountListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountClient.ListByResourceGroup")
        defer func() {
            sc := -1
            if result.alr.Response.Response != nil {
                sc = result.alr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("automation.AccountClient", "ListByResourceGroup", err.Error())
            }

                        result.fn = client.listByResourceGroupNextResults
    req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.AccountClient", "ListByResourceGroup", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByResourceGroupSender(req)
            if err != nil {
            result.alr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "ListByResourceGroup", resp, "Failure sending request")
            return
            }

            result.alr, err = client.ListByResourceGroupResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "ListByResourceGroup", resp, "Failure responding to request")
            }

    return
    }

    // ListByResourceGroupPreparer prepares the ListByResourceGroup request.
    func (client AccountClient) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
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
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client AccountClient) ListByResourceGroupResponder(resp *http.Response) (result AccountListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listByResourceGroupNextResults retrieves the next set of results, if any.
            func (client AccountClient) listByResourceGroupNextResults(ctx context.Context, lastResults AccountListResult) (result AccountListResult, err error) {
            req, err := lastResults.accountListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "automation.AccountClient", "listByResourceGroupNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListByResourceGroupSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "automation.AccountClient", "listByResourceGroupNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListByResourceGroupResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListByResourceGroupComplete enumerates all values, automatically crossing page boundaries as required.
    func (client AccountClient) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result AccountListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/AccountClient.ListByResourceGroup")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.ListByResourceGroup(ctx, resourceGroupName)
                return
        }

// Update update an automation account.
    // Parameters:
        // resourceGroupName - name of an Azure Resource group.
        // automationAccountName - the name of the automation account.
        // parameters - parameters supplied to the update automation account.
func (client AccountClient) Update(ctx context.Context, resourceGroupName string, automationAccountName string, parameters AccountUpdateParameters) (result Account, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountClient.Update")
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
            return result, validation.NewError("automation.AccountClient", "Update", err.Error())
            }

                req, err := client.UpdatePreparer(ctx, resourceGroupName, automationAccountName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "automation.AccountClient", "Update", nil , "Failure preparing request")
    return
    }

            resp, err := client.UpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "Update", resp, "Failure sending request")
            return
            }

            result, err = client.UpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "automation.AccountClient", "Update", resp, "Failure responding to request")
            }

    return
    }

    // UpdatePreparer prepares the Update request.
    func (client AccountClient) UpdatePreparer(ctx context.Context, resourceGroupName string, automationAccountName string, parameters AccountUpdateParameters) (*http.Request, error) {
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
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPatch(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // UpdateSender sends the Update request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountClient) UpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client AccountClient) UpdateResponder(resp *http.Response) (result Account, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

