package dtl

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

// CustomImagesClient is the the DevTest Labs Client.
type CustomImagesClient struct {
    BaseClient
}
// NewCustomImagesClient creates an instance of the CustomImagesClient client.
func NewCustomImagesClient(subscriptionID string) CustomImagesClient {
    return NewCustomImagesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewCustomImagesClientWithBaseURI creates an instance of the CustomImagesClient client.
    func NewCustomImagesClientWithBaseURI(baseURI string, subscriptionID string) CustomImagesClient {
        return CustomImagesClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate create or replace an existing custom image. This operation can take a while to complete.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // labName - the name of the lab.
        // name - the name of the custom image.
        // customImage - a custom image.
func (client CustomImagesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, labName string, name string, customImage CustomImage) (result CustomImagesCreateOrUpdateFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/CustomImagesClient.CreateOrUpdate")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: customImage,
             Constraints: []validation.Constraint{	{Target: "customImage.CustomImageProperties", Name: validation.Null, Rule: true, Chain: nil }}}}); err != nil {
            return result, validation.NewError("dtl.CustomImagesClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, labName, name, customImage)
    if err != nil {
    err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            result, err = client.CreateOrUpdateSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "CreateOrUpdate", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client CustomImagesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, labName string, name string, customImage CustomImage) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "labName": autorest.Encode("path",labName),
            "name": autorest.Encode("path",name),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-05-15"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/customimages/{name}",pathParameters),
    autorest.WithJSON(customImage),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client CustomImagesClient) CreateOrUpdateSender(req *http.Request) (future CustomImagesCreateOrUpdateFuture, err error) {
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
func (client CustomImagesClient) CreateOrUpdateResponder(resp *http.Response) (result CustomImage, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete delete custom image. This operation can take a while to complete.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // labName - the name of the lab.
        // name - the name of the custom image.
func (client CustomImagesClient) Delete(ctx context.Context, resourceGroupName string, labName string, name string) (result CustomImagesDeleteFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/CustomImagesClient.Delete")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.DeletePreparer(ctx, resourceGroupName, labName, name)
    if err != nil {
    err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "Delete", nil , "Failure preparing request")
    return
    }

            result, err = client.DeleteSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "Delete", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client CustomImagesClient) DeletePreparer(ctx context.Context, resourceGroupName string, labName string, name string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "labName": autorest.Encode("path",labName),
            "name": autorest.Encode("path",name),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-05-15"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/customimages/{name}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client CustomImagesClient) DeleteSender(req *http.Request) (future CustomImagesDeleteFuture, err error) {
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
func (client CustomImagesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get get custom image.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // labName - the name of the lab.
        // name - the name of the custom image.
        // expand - specify the $expand query. Example: 'properties($select=vm)'
func (client CustomImagesClient) Get(ctx context.Context, resourceGroupName string, labName string, name string, expand string) (result CustomImage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/CustomImagesClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetPreparer(ctx, resourceGroupName, labName, name, expand)
    if err != nil {
    err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client CustomImagesClient) GetPreparer(ctx context.Context, resourceGroupName string, labName string, name string, expand string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "labName": autorest.Encode("path",labName),
            "name": autorest.Encode("path",name),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-05-15"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if len(expand) > 0 {
            queryParameters["$expand"] = autorest.Encode("query",expand)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/customimages/{name}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client CustomImagesClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client CustomImagesClient) GetResponder(resp *http.Response) (result CustomImage, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List list custom images in a given lab.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // labName - the name of the lab.
        // expand - specify the $expand query. Example: 'properties($select=vm)'
        // filter - the filter to apply to the operation.
        // top - the maximum number of resources to return from the operation.
        // orderby - the ordering expression for the results, using OData notation.
func (client CustomImagesClient) List(ctx context.Context, resourceGroupName string, labName string, expand string, filter string, top *int32, orderby string) (result ResponseWithContinuationCustomImagePage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/CustomImagesClient.List")
        defer func() {
            sc := -1
            if result.rwcci.Response.Response != nil {
                sc = result.rwcci.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx, resourceGroupName, labName, expand, filter, top, orderby)
    if err != nil {
    err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.rwcci.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "List", resp, "Failure sending request")
            return
            }

            result.rwcci, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client CustomImagesClient) ListPreparer(ctx context.Context, resourceGroupName string, labName string, expand string, filter string, top *int32, orderby string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "labName": autorest.Encode("path",labName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-05-15"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if len(expand) > 0 {
            queryParameters["$expand"] = autorest.Encode("query",expand)
            }
            if len(filter) > 0 {
            queryParameters["$filter"] = autorest.Encode("query",filter)
            }
            if top != nil {
            queryParameters["$top"] = autorest.Encode("query",*top)
            }
            if len(orderby) > 0 {
            queryParameters["$orderby"] = autorest.Encode("query",orderby)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/customimages",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client CustomImagesClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client CustomImagesClient) ListResponder(resp *http.Response) (result ResponseWithContinuationCustomImage, err error) {
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
            func (client CustomImagesClient) listNextResults(ctx context.Context, lastResults ResponseWithContinuationCustomImage) (result ResponseWithContinuationCustomImage, err error) {
            req, err := lastResults.responseWithContinuationCustomImagePreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "dtl.CustomImagesClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client CustomImagesClient) ListComplete(ctx context.Context, resourceGroupName string, labName string, expand string, filter string, top *int32, orderby string) (result ResponseWithContinuationCustomImageIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/CustomImagesClient.List")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.List(ctx, resourceGroupName, labName, expand, filter, top, orderby)
                return
        }

