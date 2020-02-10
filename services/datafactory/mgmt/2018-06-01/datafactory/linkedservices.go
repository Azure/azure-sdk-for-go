package datafactory

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

// LinkedServicesClient is the the Azure Data Factory V2 management API provides a RESTful set of web services that
// interact with Azure Data Factory V2 services.
type LinkedServicesClient struct {
    BaseClient
}
// NewLinkedServicesClient creates an instance of the LinkedServicesClient client.
func NewLinkedServicesClient(subscriptionID string) LinkedServicesClient {
    return NewLinkedServicesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewLinkedServicesClientWithBaseURI creates an instance of the LinkedServicesClient client using a custom endpoint.
// Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
    func NewLinkedServicesClientWithBaseURI(baseURI string, subscriptionID string) LinkedServicesClient {
        return LinkedServicesClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate creates or updates a linked service.
    // Parameters:
        // resourceGroupName - the resource group name.
        // factoryName - the factory name.
        // linkedServiceName - the linked service name.
        // linkedService - linked service resource definition.
        // ifMatch - eTag of the linkedService entity.  Should only be specified for update, for which it should match
        // existing entity or can be * for unconditional update.
func (client LinkedServicesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string, linkedService LinkedServiceResource, ifMatch string) (result LinkedServiceResource, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/LinkedServicesClient.CreateOrUpdate")
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
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}},
            { TargetValue: factoryName,
             Constraints: []validation.Constraint{	{Target: "factoryName", Name: validation.MaxLength, Rule: 63, Chain: nil },
            	{Target: "factoryName", Name: validation.MinLength, Rule: 3, Chain: nil },
            	{Target: "factoryName", Name: validation.Pattern, Rule: `^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`, Chain: nil }}},
            { TargetValue: linkedServiceName,
             Constraints: []validation.Constraint{	{Target: "linkedServiceName", Name: validation.MaxLength, Rule: 260, Chain: nil },
            	{Target: "linkedServiceName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "linkedServiceName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil }}},
            { TargetValue: linkedService,
             Constraints: []validation.Constraint{	{Target: "linkedService.Properties", Name: validation.Null, Rule: true ,
            Chain: []validation.Constraint{	{Target: "linkedService.Properties.ConnectVia", Name: validation.Null, Rule: false ,
            Chain: []validation.Constraint{	{Target: "linkedService.Properties.ConnectVia.Type", Name: validation.Null, Rule: true, Chain: nil },
            	{Target: "linkedService.Properties.ConnectVia.ReferenceName", Name: validation.Null, Rule: true, Chain: nil },
            }},
            }}}}}); err != nil {
            return result, validation.NewError("datafactory.LinkedServicesClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, factoryName, linkedServiceName, linkedService, ifMatch)
    if err != nil {
    err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "CreateOrUpdate", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "CreateOrUpdate", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client LinkedServicesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string, linkedService LinkedServiceResource, ifMatch string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "factoryName": autorest.Encode("path",factoryName),
            "linkedServiceName": autorest.Encode("path",linkedServiceName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-06-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

        preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/linkedservices/{linkedServiceName}",pathParameters),
    autorest.WithJSON(linkedService),
    autorest.WithQueryParameters(queryParameters))
            if len(ifMatch) > 0 {
            preparer = autorest.DecoratePreparer(preparer,
            autorest.WithHeader("If-Match",autorest.String(ifMatch)))
            }
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client LinkedServicesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
        sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
            return autorest.SendWithSender(client, req, sd...)
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client LinkedServicesClient) CreateOrUpdateResponder(resp *http.Response) (result LinkedServiceResource, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete deletes a linked service.
    // Parameters:
        // resourceGroupName - the resource group name.
        // factoryName - the factory name.
        // linkedServiceName - the linked service name.
func (client LinkedServicesClient) Delete(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/LinkedServicesClient.Delete")
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
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}},
            { TargetValue: factoryName,
             Constraints: []validation.Constraint{	{Target: "factoryName", Name: validation.MaxLength, Rule: 63, Chain: nil },
            	{Target: "factoryName", Name: validation.MinLength, Rule: 3, Chain: nil },
            	{Target: "factoryName", Name: validation.Pattern, Rule: `^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`, Chain: nil }}},
            { TargetValue: linkedServiceName,
             Constraints: []validation.Constraint{	{Target: "linkedServiceName", Name: validation.MaxLength, Rule: 260, Chain: nil },
            	{Target: "linkedServiceName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "linkedServiceName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("datafactory.LinkedServicesClient", "Delete", err.Error())
            }

                req, err := client.DeletePreparer(ctx, resourceGroupName, factoryName, linkedServiceName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client LinkedServicesClient) DeletePreparer(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "factoryName": autorest.Encode("path",factoryName),
            "linkedServiceName": autorest.Encode("path",linkedServiceName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-06-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

        preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/linkedservices/{linkedServiceName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client LinkedServicesClient) DeleteSender(req *http.Request) (*http.Response, error) {
        sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
            return autorest.SendWithSender(client, req, sd...)
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client LinkedServicesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get gets a linked service.
    // Parameters:
        // resourceGroupName - the resource group name.
        // factoryName - the factory name.
        // linkedServiceName - the linked service name.
        // ifNoneMatch - eTag of the linked service entity. Should only be specified for get. If the ETag matches the
        // existing entity tag, or if * was provided, then no content will be returned.
func (client LinkedServicesClient) Get(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string, ifNoneMatch string) (result LinkedServiceResource, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/LinkedServicesClient.Get")
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
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}},
            { TargetValue: factoryName,
             Constraints: []validation.Constraint{	{Target: "factoryName", Name: validation.MaxLength, Rule: 63, Chain: nil },
            	{Target: "factoryName", Name: validation.MinLength, Rule: 3, Chain: nil },
            	{Target: "factoryName", Name: validation.Pattern, Rule: `^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`, Chain: nil }}},
            { TargetValue: linkedServiceName,
             Constraints: []validation.Constraint{	{Target: "linkedServiceName", Name: validation.MaxLength, Rule: 260, Chain: nil },
            	{Target: "linkedServiceName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "linkedServiceName", Name: validation.Pattern, Rule: `^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("datafactory.LinkedServicesClient", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, factoryName, linkedServiceName, ifNoneMatch)
    if err != nil {
    err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client LinkedServicesClient) GetPreparer(ctx context.Context, resourceGroupName string, factoryName string, linkedServiceName string, ifNoneMatch string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "factoryName": autorest.Encode("path",factoryName),
            "linkedServiceName": autorest.Encode("path",linkedServiceName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-06-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

        preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/linkedservices/{linkedServiceName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
            if len(ifNoneMatch) > 0 {
            preparer = autorest.DecoratePreparer(preparer,
            autorest.WithHeader("If-None-Match",autorest.String(ifNoneMatch)))
            }
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client LinkedServicesClient) GetSender(req *http.Request) (*http.Response, error) {
        sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
            return autorest.SendWithSender(client, req, sd...)
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client LinkedServicesClient) GetResponder(resp *http.Response) (result LinkedServiceResource, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNotModified),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// ListByFactory lists linked services.
    // Parameters:
        // resourceGroupName - the resource group name.
        // factoryName - the factory name.
func (client LinkedServicesClient) ListByFactory(ctx context.Context, resourceGroupName string, factoryName string) (result LinkedServiceListResponsePage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/LinkedServicesClient.ListByFactory")
        defer func() {
            sc := -1
            if result.lslr.Response.Response != nil {
                sc = result.lslr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: resourceGroupName,
             Constraints: []validation.Constraint{	{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil }}},
            { TargetValue: factoryName,
             Constraints: []validation.Constraint{	{Target: "factoryName", Name: validation.MaxLength, Rule: 63, Chain: nil },
            	{Target: "factoryName", Name: validation.MinLength, Rule: 3, Chain: nil },
            	{Target: "factoryName", Name: validation.Pattern, Rule: `^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("datafactory.LinkedServicesClient", "ListByFactory", err.Error())
            }

                        result.fn = client.listByFactoryNextResults
    req, err := client.ListByFactoryPreparer(ctx, resourceGroupName, factoryName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "ListByFactory", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByFactorySender(req)
            if err != nil {
            result.lslr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "ListByFactory", resp, "Failure sending request")
            return
            }

            result.lslr, err = client.ListByFactoryResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "ListByFactory", resp, "Failure responding to request")
            }

    return
    }

    // ListByFactoryPreparer prepares the ListByFactory request.
    func (client LinkedServicesClient) ListByFactoryPreparer(ctx context.Context, resourceGroupName string, factoryName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "factoryName": autorest.Encode("path",factoryName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-06-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

        preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/linkedservices",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByFactorySender sends the ListByFactory request. The method will close the
    // http.Response Body if it receives an error.
    func (client LinkedServicesClient) ListByFactorySender(req *http.Request) (*http.Response, error) {
        sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
            return autorest.SendWithSender(client, req, sd...)
            }

// ListByFactoryResponder handles the response to the ListByFactory request. The method always
// closes the http.Response Body.
func (client LinkedServicesClient) ListByFactoryResponder(resp *http.Response) (result LinkedServiceListResponse, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listByFactoryNextResults retrieves the next set of results, if any.
            func (client LinkedServicesClient) listByFactoryNextResults(ctx context.Context, lastResults LinkedServiceListResponse) (result LinkedServiceListResponse, err error) {
            req, err := lastResults.linkedServiceListResponsePreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "listByFactoryNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListByFactorySender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "listByFactoryNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListByFactoryResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "datafactory.LinkedServicesClient", "listByFactoryNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListByFactoryComplete enumerates all values, automatically crossing page boundaries as required.
    func (client LinkedServicesClient) ListByFactoryComplete(ctx context.Context, resourceGroupName string, factoryName string) (result LinkedServiceListResponseIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/LinkedServicesClient.ListByFactory")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.ListByFactory(ctx, resourceGroupName, factoryName)
                return
        }

