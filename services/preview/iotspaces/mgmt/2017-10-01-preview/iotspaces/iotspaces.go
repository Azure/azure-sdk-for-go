package iotspaces

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
    "github.com/satori/go.uuid"
)

// Client is the use this API to manage the IoTSpaces service instances in your Azure subscription.
type Client struct {
    BaseClient
}
// NewClient creates an instance of the Client client.
func NewClient(subscriptionID uuid.UUID) Client {
    return NewClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewClientWithBaseURI creates an instance of the Client client.
    func NewClientWithBaseURI(baseURI string, subscriptionID uuid.UUID) Client {
        return Client{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CheckNameAvailability check if an IoTSpaces instance name is available.
    // Parameters:
        // operationInputs - set the name parameter in the OperationInputs structure to the name of the IoTSpaces
        // instance to check.
func (client Client) CheckNameAvailability(ctx context.Context, operationInputs OperationInputs) (result NameAvailabilityInfo, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/Client.CheckNameAvailability")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: operationInputs,
             Constraints: []validation.Constraint{	{Target: "operationInputs.Name", Name: validation.Null, Rule: true, Chain: nil }}}}); err != nil {
            return result, validation.NewError("iotspaces.Client", "CheckNameAvailability", err.Error())
            }

                req, err := client.CheckNameAvailabilityPreparer(ctx, operationInputs)
    if err != nil {
    err = autorest.NewErrorWithError(err, "iotspaces.Client", "CheckNameAvailability", nil , "Failure preparing request")
    return
    }

            resp, err := client.CheckNameAvailabilitySender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "CheckNameAvailability", resp, "Failure sending request")
            return
            }

            result, err = client.CheckNameAvailabilityResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "CheckNameAvailability", resp, "Failure responding to request")
            }

    return
    }

    // CheckNameAvailabilityPreparer prepares the CheckNameAvailability request.
    func (client Client) CheckNameAvailabilityPreparer(ctx context.Context, operationInputs OperationInputs) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPost(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.IoTSpaces/checkNameAvailability",pathParameters),
    autorest.WithJSON(operationInputs),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CheckNameAvailabilitySender sends the CheckNameAvailability request. The method will close the
    // http.Response Body if it receives an error.
    func (client Client) CheckNameAvailabilitySender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CheckNameAvailabilityResponder handles the response to the CheckNameAvailability request. The method always
// closes the http.Response Body.
func (client Client) CheckNameAvailabilityResponder(resp *http.Response) (result NameAvailabilityInfo, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// CreateOrUpdate create or update the metadata of an IoTSpaces instance. The usual pattern to modify a property is to
// retrieve the IoTSpaces instance metadata and security metadata, and then combine them with the modified values in a
// new body to update the IoTSpaces instance.
    // Parameters:
        // resourceGroupName - the name of the resource group that contains the IoTSpaces instance.
        // resourceName - the name of the IoTSpaces instance.
        // iotSpaceDescription - the IoTSpaces instance metadata and security metadata.
func (client Client) CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, iotSpaceDescription Description) (result CreateOrUpdateFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/Client.CreateOrUpdate")
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
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil }}},
            { TargetValue: resourceName,
             Constraints: []validation.Constraint{	{Target: "resourceName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceName", Name: validation.MinLength, Rule: 1, Chain: nil }}},
            { TargetValue: iotSpaceDescription,
             Constraints: []validation.Constraint{	{Target: "iotSpaceDescription.Sku", Name: validation.Null, Rule: true, Chain: nil }}}}); err != nil {
            return result, validation.NewError("iotspaces.Client", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, resourceName, iotSpaceDescription)
    if err != nil {
    err = autorest.NewErrorWithError(err, "iotspaces.Client", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            result, err = client.CreateOrUpdateSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "CreateOrUpdate", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client Client) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, resourceName string, iotSpaceDescription Description) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "resourceName": autorest.Encode("path",resourceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.IoTSpaces/Graph/{resourceName}",pathParameters),
    autorest.WithJSON(iotSpaceDescription),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client Client) CreateOrUpdateSender(req *http.Request) (future CreateOrUpdateFuture, err error) {
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
func (client Client) CreateOrUpdateResponder(resp *http.Response) (result Description, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete delete an IoTSpaces instance.
    // Parameters:
        // resourceGroupName - the name of the resource group that contains the IoTSpaces instance.
        // resourceName - the name of the IoTSpaces instance.
func (client Client) Delete(ctx context.Context, resourceGroupName string, resourceName string) (result DeleteFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/Client.Delete")
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
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil }}},
            { TargetValue: resourceName,
             Constraints: []validation.Constraint{	{Target: "resourceName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceName", Name: validation.MinLength, Rule: 1, Chain: nil }}}}); err != nil {
            return result, validation.NewError("iotspaces.Client", "Delete", err.Error())
            }

                req, err := client.DeletePreparer(ctx, resourceGroupName, resourceName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "iotspaces.Client", "Delete", nil , "Failure preparing request")
    return
    }

            result, err = client.DeleteSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "Delete", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client Client) DeletePreparer(ctx context.Context, resourceGroupName string, resourceName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "resourceName": autorest.Encode("path",resourceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.IoTSpaces/Graph/{resourceName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client Client) DeleteSender(req *http.Request) (future DeleteFuture, err error) {
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
func (client Client) DeleteResponder(resp *http.Response) (result Description, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted,http.StatusNoContent),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Get get the metadata of a IoTSpaces instance.
    // Parameters:
        // resourceGroupName - the name of the resource group that contains the IoTSpaces instance.
        // resourceName - the name of the IoTSpaces instance.
func (client Client) Get(ctx context.Context, resourceGroupName string, resourceName string) (result Description, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/Client.Get")
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
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil }}},
            { TargetValue: resourceName,
             Constraints: []validation.Constraint{	{Target: "resourceName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceName", Name: validation.MinLength, Rule: 1, Chain: nil }}}}); err != nil {
            return result, validation.NewError("iotspaces.Client", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, resourceName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "iotspaces.Client", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client Client) GetPreparer(ctx context.Context, resourceGroupName string, resourceName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "resourceName": autorest.Encode("path",resourceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.IoTSpaces/Graph/{resourceName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client Client) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client Client) GetResponder(resp *http.Response) (result Description, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List get all the IoTSpaces instances in a subscription.
func (client Client) List(ctx context.Context) (result DescriptionListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/Client.List")
        defer func() {
            sc := -1
            if result.dlr.Response.Response != nil {
                sc = result.dlr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx)
    if err != nil {
    err = autorest.NewErrorWithError(err, "iotspaces.Client", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.dlr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "List", resp, "Failure sending request")
            return
            }

            result.dlr, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client Client) ListPreparer(ctx context.Context) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.IoTSpaces/Graph",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client Client) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client Client) ListResponder(resp *http.Response) (result DescriptionListResult, err error) {
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
            func (client Client) listNextResults(ctx context.Context, lastResults DescriptionListResult) (result DescriptionListResult, err error) {
            req, err := lastResults.descriptionListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "iotspaces.Client", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "iotspaces.Client", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client Client) ListComplete(ctx context.Context) (result DescriptionListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/Client.List")
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

// ListByResourceGroup get all the IoTSpaces instances in a resource group.
    // Parameters:
        // resourceGroupName - the name of the resource group that contains the IoTSpaces instance.
func (client Client) ListByResourceGroup(ctx context.Context, resourceGroupName string) (result DescriptionListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/Client.ListByResourceGroup")
        defer func() {
            sc := -1
            if result.dlr.Response.Response != nil {
                sc = result.dlr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil }}}}); err != nil {
            return result, validation.NewError("iotspaces.Client", "ListByResourceGroup", err.Error())
            }

                        result.fn = client.listByResourceGroupNextResults
    req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "iotspaces.Client", "ListByResourceGroup", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByResourceGroupSender(req)
            if err != nil {
            result.dlr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "ListByResourceGroup", resp, "Failure sending request")
            return
            }

            result.dlr, err = client.ListByResourceGroupResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "ListByResourceGroup", resp, "Failure responding to request")
            }

    return
    }

    // ListByResourceGroupPreparer prepares the ListByResourceGroup request.
    func (client Client) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.IoTSpaces/Graph",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
    // http.Response Body if it receives an error.
    func (client Client) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client Client) ListByResourceGroupResponder(resp *http.Response) (result DescriptionListResult, err error) {
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
            func (client Client) listByResourceGroupNextResults(ctx context.Context, lastResults DescriptionListResult) (result DescriptionListResult, err error) {
            req, err := lastResults.descriptionListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "iotspaces.Client", "listByResourceGroupNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListByResourceGroupSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "iotspaces.Client", "listByResourceGroupNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListByResourceGroupResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListByResourceGroupComplete enumerates all values, automatically crossing page boundaries as required.
    func (client Client) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result DescriptionListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/Client.ListByResourceGroup")
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

// Update update the metadata of a IoTSpaces instance.
    // Parameters:
        // resourceGroupName - the name of the resource group that contains the IoTSpaces instance.
        // resourceName - the name of the IoTSpaces instance.
        // iotSpacePatchDescription - the IoTSpaces instance metadata and security metadata.
func (client Client) Update(ctx context.Context, resourceGroupName string, resourceName string, iotSpacePatchDescription PatchDescription) (result UpdateFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/Client.Update")
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
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil }}},
            { TargetValue: resourceName,
             Constraints: []validation.Constraint{	{Target: "resourceName", Name: validation.MaxLength, Rule: 64, Chain: nil },
            	{Target: "resourceName", Name: validation.MinLength, Rule: 1, Chain: nil }}}}); err != nil {
            return result, validation.NewError("iotspaces.Client", "Update", err.Error())
            }

                req, err := client.UpdatePreparer(ctx, resourceGroupName, resourceName, iotSpacePatchDescription)
    if err != nil {
    err = autorest.NewErrorWithError(err, "iotspaces.Client", "Update", nil , "Failure preparing request")
    return
    }

            result, err = client.UpdateSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.Client", "Update", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // UpdatePreparer prepares the Update request.
    func (client Client) UpdatePreparer(ctx context.Context, resourceGroupName string, resourceName string, iotSpacePatchDescription PatchDescription) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "resourceName": autorest.Encode("path",resourceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2017-10-01-preview"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPatch(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.IoTSpaces/Graph/{resourceName}",pathParameters),
    autorest.WithJSON(iotSpacePatchDescription),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // UpdateSender sends the Update request. The method will close the
    // http.Response Body if it receives an error.
    func (client Client) UpdateSender(req *http.Request) (future UpdateFuture, err error) {
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
func (client Client) UpdateResponder(resp *http.Response) (result Description, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

