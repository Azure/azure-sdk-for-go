package storage

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

// AccountsClient is the the Storage Management Client.
type AccountsClient struct {
    BaseClient
}
// NewAccountsClient creates an instance of the AccountsClient client.
func NewAccountsClient(subscriptionID string) AccountsClient {
    return NewAccountsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewAccountsClientWithBaseURI creates an instance of the AccountsClient client.
    func NewAccountsClientWithBaseURI(baseURI string, subscriptionID string) AccountsClient {
        return AccountsClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CheckNameAvailability checks that account name is valid and is not in use.
    // Parameters:
        // accountName - the name of the storage account within the specified resource group. Storage account names
        // must be between 3 and 24 characters in length and use numbers and lower-case letters only.
func (client AccountsClient) CheckNameAvailability(ctx context.Context, accountName AccountCheckNameAvailabilityParameters) (result CheckNameAvailabilityResult, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.CheckNameAvailability")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.CheckNameAvailabilityPreparer(ctx, accountName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "CheckNameAvailability", nil , "Failure preparing request")
    return
    }

            resp, err := client.CheckNameAvailabilitySender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "CheckNameAvailability", resp, "Failure sending request")
            return
            }

            result, err = client.CheckNameAvailabilityResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "CheckNameAvailability", resp, "Failure responding to request")
            }

    return
    }

    // CheckNameAvailabilityPreparer prepares the CheckNameAvailability request.
    func (client AccountsClient) CheckNameAvailabilityPreparer(ctx context.Context, accountName AccountCheckNameAvailabilityParameters) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPost(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Storage/checkNameAvailability",pathParameters),
    autorest.WithJSON(accountName),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CheckNameAvailabilitySender sends the CheckNameAvailability request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) CheckNameAvailabilitySender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CheckNameAvailabilityResponder handles the response to the CheckNameAvailability request. The method always
// closes the http.Response Body.
func (client AccountsClient) CheckNameAvailabilityResponder(resp *http.Response) (result CheckNameAvailabilityResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Create asynchronously creates a new storage account with the specified parameters. Existing accounts cannot be
// updated with this API and should instead use the Update Storage Account API. If an account is already created and
// subsequent PUT request is issued with exact same set of properties, then HTTP 200 would be returned.
    // Parameters:
        // resourceGroupName - the name of the resource group within the user's subscription.
        // accountName - the name of the storage account within the specified resource group. Storage account names
        // must be between 3 and 24 characters in length and use numbers and lower-case letters only.
        // parameters - the parameters to provide for the created account.
func (client AccountsClient) Create(ctx context.Context, resourceGroupName string, accountName string, parameters AccountCreateParameters) (result AccountsCreateFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.Create")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.CreatePreparer(ctx, resourceGroupName, accountName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Create", nil , "Failure preparing request")
    return
    }

            result, err = client.CreateSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Create", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // CreatePreparer prepares the Create request.
    func (client AccountsClient) CreatePreparer(ctx context.Context, resourceGroupName string, accountName string, parameters AccountCreateParameters) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "accountName": autorest.Encode("path",accountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateSender sends the Create request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) CreateSender(req *http.Request) (future AccountsCreateFuture, err error) {
            var resp *http.Response
            resp, err = autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            if err != nil {
            return
            }
            future.Future, err = azure.NewFutureFromResponse(resp)
            return
            }

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client AccountsClient) CreateResponder(resp *http.Response) (result Account, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete deletes a storage account in Microsoft Azure.
    // Parameters:
        // resourceGroupName - the name of the resource group within the user's subscription.
        // accountName - the name of the storage account within the specified resource group. Storage account names
        // must be between 3 and 24 characters in length and use numbers and lower-case letters only.
func (client AccountsClient) Delete(ctx context.Context, resourceGroupName string, accountName string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.Delete")
        defer func() {
            sc := -1
            if result.Response != nil {
                sc = result.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.DeletePreparer(ctx, resourceGroupName, accountName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client AccountsClient) DeletePreparer(ctx context.Context, resourceGroupName string, accountName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "accountName": autorest.Encode("path",accountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client AccountsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// GetProperties returns the properties for the specified storage account including but not limited to name, account
// type, location, and account status. The ListKeys operation should be used to retrieve storage keys.
    // Parameters:
        // resourceGroupName - the name of the resource group within the user's subscription.
        // accountName - the name of the storage account within the specified resource group. Storage account names
        // must be between 3 and 24 characters in length and use numbers and lower-case letters only.
func (client AccountsClient) GetProperties(ctx context.Context, resourceGroupName string, accountName string) (result Account, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.GetProperties")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetPropertiesPreparer(ctx, resourceGroupName, accountName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "GetProperties", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetPropertiesSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "GetProperties", resp, "Failure sending request")
            return
            }

            result, err = client.GetPropertiesResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "GetProperties", resp, "Failure responding to request")
            }

    return
    }

    // GetPropertiesPreparer prepares the GetProperties request.
    func (client AccountsClient) GetPropertiesPreparer(ctx context.Context, resourceGroupName string, accountName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "accountName": autorest.Encode("path",accountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetPropertiesSender sends the GetProperties request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) GetPropertiesSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetPropertiesResponder handles the response to the GetProperties request. The method always
// closes the http.Response Body.
func (client AccountsClient) GetPropertiesResponder(resp *http.Response) (result Account, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List lists all the storage accounts available under the subscription. Note that storage keys are not returned; use
// the ListKeys operation for this.
func (client AccountsClient) List(ctx context.Context) (result AccountListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.List")
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
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.alr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "List", resp, "Failure sending request")
            return
            }

            result.alr, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client AccountsClient) ListPreparer(ctx context.Context) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Storage/storageAccounts",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client AccountsClient) ListResponder(resp *http.Response) (result AccountListResult, err error) {
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
            func (client AccountsClient) listNextResults(ctx context.Context, lastResults AccountListResult) (result AccountListResult, err error) {
            req, err := lastResults.accountListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client AccountsClient) ListComplete(ctx context.Context) (result AccountListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.List")
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

// ListByResourceGroup lists all the storage accounts available under the given resource group. Note that storage keys
// are not returned; use the ListKeys operation for this.
    // Parameters:
        // resourceGroupName - the name of the resource group within the user's subscription.
func (client AccountsClient) ListByResourceGroup(ctx context.Context, resourceGroupName string) (result AccountListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.ListByResourceGroup")
        defer func() {
            sc := -1
            if result.alr.Response.Response != nil {
                sc = result.alr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listByResourceGroupNextResults
    req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "ListByResourceGroup", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByResourceGroupSender(req)
            if err != nil {
            result.alr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "ListByResourceGroup", resp, "Failure sending request")
            return
            }

            result.alr, err = client.ListByResourceGroupResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "ListByResourceGroup", resp, "Failure responding to request")
            }

    return
    }

    // ListByResourceGroupPreparer prepares the ListByResourceGroup request.
    func (client AccountsClient) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client AccountsClient) ListByResourceGroupResponder(resp *http.Response) (result AccountListResult, err error) {
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
            func (client AccountsClient) listByResourceGroupNextResults(ctx context.Context, lastResults AccountListResult) (result AccountListResult, err error) {
            req, err := lastResults.accountListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "listByResourceGroupNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListByResourceGroupSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "storage.AccountsClient", "listByResourceGroupNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListByResourceGroupResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListByResourceGroupComplete enumerates all values, automatically crossing page boundaries as required.
    func (client AccountsClient) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result AccountListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.ListByResourceGroup")
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

// ListKeys lists the access keys for the specified storage account.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // accountName - the name of the storage account.
func (client AccountsClient) ListKeys(ctx context.Context, resourceGroupName string, accountName string) (result AccountKeys, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.ListKeys")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.ListKeysPreparer(ctx, resourceGroupName, accountName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "ListKeys", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListKeysSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "ListKeys", resp, "Failure sending request")
            return
            }

            result, err = client.ListKeysResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "ListKeys", resp, "Failure responding to request")
            }

    return
    }

    // ListKeysPreparer prepares the ListKeys request.
    func (client AccountsClient) ListKeysPreparer(ctx context.Context, resourceGroupName string, accountName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "accountName": autorest.Encode("path",accountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsPost(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/listKeys",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListKeysSender sends the ListKeys request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) ListKeysSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListKeysResponder handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (client AccountsClient) ListKeysResponder(resp *http.Response) (result AccountKeys, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// RegenerateKey regenerates the access keys for the specified storage account.
    // Parameters:
        // resourceGroupName - the name of the resource group within the user's subscription.
        // accountName - the name of the storage account within the specified resource group. Storage account names
        // must be between 3 and 24 characters in length and use numbers and lower-case letters only.
        // regenerateKey - specifies name of the key which should be regenerated.
func (client AccountsClient) RegenerateKey(ctx context.Context, resourceGroupName string, accountName string, regenerateKey AccountRegenerateKeyParameters) (result AccountKeys, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.RegenerateKey")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.RegenerateKeyPreparer(ctx, resourceGroupName, accountName, regenerateKey)
    if err != nil {
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "RegenerateKey", nil , "Failure preparing request")
    return
    }

            resp, err := client.RegenerateKeySender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "RegenerateKey", resp, "Failure sending request")
            return
            }

            result, err = client.RegenerateKeyResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "RegenerateKey", resp, "Failure responding to request")
            }

    return
    }

    // RegenerateKeyPreparer prepares the RegenerateKey request.
    func (client AccountsClient) RegenerateKeyPreparer(ctx context.Context, resourceGroupName string, accountName string, regenerateKey AccountRegenerateKeyParameters) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "accountName": autorest.Encode("path",accountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPost(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/regenerateKey",pathParameters),
    autorest.WithJSON(regenerateKey),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // RegenerateKeySender sends the RegenerateKey request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) RegenerateKeySender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// RegenerateKeyResponder handles the response to the RegenerateKey request. The method always
// closes the http.Response Body.
func (client AccountsClient) RegenerateKeyResponder(resp *http.Response) (result AccountKeys, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Update updates the account type or tags for a storage account. It can also be used to add a custom domain (note that
// custom domains cannot be added via the Create operation). Only one custom domain is supported per storage account.
// This API can only be used to update one of tags, accountType, or customDomain per call. To update multiple of these
// properties, call the API multiple times with one change per call. This call does not change the storage keys for the
// account. If you want to change storage account keys, use the RegenerateKey operation. The location and name of the
// storage account cannot be changed after creation.
    // Parameters:
        // resourceGroupName - the name of the resource group within the user's subscription.
        // accountName - the name of the storage account within the specified resource group. Storage account names
        // must be between 3 and 24 characters in length and use numbers and lower-case letters only.
        // parameters - the parameters to update on the account. Note that only one property can be changed at a time
        // using this API.
func (client AccountsClient) Update(ctx context.Context, resourceGroupName string, accountName string, parameters AccountUpdateParameters) (result Account, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AccountsClient.Update")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.UpdatePreparer(ctx, resourceGroupName, accountName, parameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Update", nil , "Failure preparing request")
    return
    }

            resp, err := client.UpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Update", resp, "Failure sending request")
            return
            }

            result, err = client.UpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "storage.AccountsClient", "Update", resp, "Failure responding to request")
            }

    return
    }

    // UpdatePreparer prepares the Update request.
    func (client AccountsClient) UpdatePreparer(ctx context.Context, resourceGroupName string, accountName string, parameters AccountUpdateParameters) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "accountName": autorest.Encode("path",accountName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-05-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPatch(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // UpdateSender sends the Update request. The method will close the
    // http.Response Body if it receives an error.
    func (client AccountsClient) UpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client AccountsClient) UpdateResponder(resp *http.Response) (result Account, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

