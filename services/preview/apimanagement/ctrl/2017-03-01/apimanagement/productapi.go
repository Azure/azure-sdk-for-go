package apimanagement

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

// ProductAPIClient is the client for the ProductAPI methods of the Apimanagement service.
type ProductAPIClient struct {
    BaseClient
}
// NewProductAPIClient creates an instance of the ProductAPIClient client.
func NewProductAPIClient() ProductAPIClient {
    return ProductAPIClient{ New()}
}

// CreateOrUpdate adds an API to the specified product.
    // Parameters:
        // apimBaseURL - the management endpoint of the API Management service, for example
        // https://myapimservice.management.azure-api.net.
        // productID - product identifier. Must be unique in the current API Management service instance.
        // apiid - API identifier. Must be unique in the current API Management service instance.
func (client ProductAPIClient) CreateOrUpdate(ctx context.Context, apimBaseURL string, productID string, apiid string) (result APIContract, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ProductAPIClient.CreateOrUpdate")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: productID,
             Constraints: []validation.Constraint{	{Target: "productID", Name: validation.MaxLength, Rule: 256, Chain: nil },
            	{Target: "productID", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "productID", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil }}},
            { TargetValue: apiid,
             Constraints: []validation.Constraint{	{Target: "apiid", Name: validation.MaxLength, Rule: 256, Chain: nil },
            	{Target: "apiid", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "apiid", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.ProductAPIClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, apimBaseURL, productID, apiid)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "CreateOrUpdate", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "CreateOrUpdate", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client ProductAPIClient) CreateOrUpdatePreparer(ctx context.Context, apimBaseURL string, productID string, apiid string) (*http.Request, error) {
            urlParameters := map[string]interface{} {
            "apimBaseUrl": apimBaseURL,
            }

            pathParameters := map[string]interface{} {
            "apiId": autorest.Encode("path",apiid),
            "productId": autorest.Encode("path",productID),
            }

                        const APIVersion = "2017-03-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsPut(),
    autorest.WithCustomBaseURL("{apimBaseUrl}", urlParameters),
    autorest.WithPathParameters("/products/{productId}/apis/{apiId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client ProductAPIClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client ProductAPIClient) CreateOrUpdateResponder(resp *http.Response) (result APIContract, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete deletes the specified API from the specified product.
    // Parameters:
        // apimBaseURL - the management endpoint of the API Management service, for example
        // https://myapimservice.management.azure-api.net.
        // productID - product identifier. Must be unique in the current API Management service instance.
        // apiid - API identifier. Must be unique in the current API Management service instance.
func (client ProductAPIClient) Delete(ctx context.Context, apimBaseURL string, productID string, apiid string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ProductAPIClient.Delete")
        defer func() {
            sc := -1
            if result.Response != nil {
                sc = result.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: productID,
             Constraints: []validation.Constraint{	{Target: "productID", Name: validation.MaxLength, Rule: 256, Chain: nil },
            	{Target: "productID", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "productID", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil }}},
            { TargetValue: apiid,
             Constraints: []validation.Constraint{	{Target: "apiid", Name: validation.MaxLength, Rule: 256, Chain: nil },
            	{Target: "apiid", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "apiid", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.ProductAPIClient", "Delete", err.Error())
            }

                req, err := client.DeletePreparer(ctx, apimBaseURL, productID, apiid)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client ProductAPIClient) DeletePreparer(ctx context.Context, apimBaseURL string, productID string, apiid string) (*http.Request, error) {
            urlParameters := map[string]interface{} {
            "apimBaseUrl": apimBaseURL,
            }

            pathParameters := map[string]interface{} {
            "apiId": autorest.Encode("path",apiid),
            "productId": autorest.Encode("path",productID),
            }

                        const APIVersion = "2017-03-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithCustomBaseURL("{apimBaseUrl}", urlParameters),
    autorest.WithPathParameters("/products/{productId}/apis/{apiId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client ProductAPIClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ProductAPIClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// ListByProduct lists a collection of the APIs associated with a product.
    // Parameters:
        // apimBaseURL - the management endpoint of the API Management service, for example
        // https://myapimservice.management.azure-api.net.
        // productID - product identifier. Must be unique in the current API Management service instance.
        // filter - | Field       | Supported operators    | Supported functions                         |
        // |-------------|------------------------|---------------------------------------------|
        // | id          | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
        // | name        | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
        // | description | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
        // | serviceUrl  | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
        // | path        | ge, le, eq, ne, gt, lt | substringof, contains, startswith, endswith |
        // top - number of records to return.
        // skip - number of records to skip.
func (client ProductAPIClient) ListByProduct(ctx context.Context, apimBaseURL string, productID string, filter string, top *int32, skip *int32) (result APICollectionPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ProductAPIClient.ListByProduct")
        defer func() {
            sc := -1
            if result.ac.Response.Response != nil {
                sc = result.ac.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: productID,
             Constraints: []validation.Constraint{	{Target: "productID", Name: validation.MaxLength, Rule: 256, Chain: nil },
            	{Target: "productID", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "productID", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil }}},
            { TargetValue: top,
             Constraints: []validation.Constraint{	{Target: "top", Name: validation.Null, Rule: false ,
            Chain: []validation.Constraint{	{Target: "top", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil },
            }}}},
            { TargetValue: skip,
             Constraints: []validation.Constraint{	{Target: "skip", Name: validation.Null, Rule: false ,
            Chain: []validation.Constraint{	{Target: "skip", Name: validation.InclusiveMinimum, Rule: 0, Chain: nil },
            }}}}}); err != nil {
            return result, validation.NewError("apimanagement.ProductAPIClient", "ListByProduct", err.Error())
            }

                        result.fn = client.listByProductNextResults
    req, err := client.ListByProductPreparer(ctx, apimBaseURL, productID, filter, top, skip)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "ListByProduct", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByProductSender(req)
            if err != nil {
            result.ac.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "ListByProduct", resp, "Failure sending request")
            return
            }

            result.ac, err = client.ListByProductResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "ListByProduct", resp, "Failure responding to request")
            }

    return
    }

    // ListByProductPreparer prepares the ListByProduct request.
    func (client ProductAPIClient) ListByProductPreparer(ctx context.Context, apimBaseURL string, productID string, filter string, top *int32, skip *int32) (*http.Request, error) {
            urlParameters := map[string]interface{} {
            "apimBaseUrl": apimBaseURL,
            }

            pathParameters := map[string]interface{} {
            "productId": autorest.Encode("path",productID),
            }

                        const APIVersion = "2017-03-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if len(filter) > 0 {
            queryParameters["$filter"] = autorest.Encode("query",filter)
            }
            if top != nil {
            queryParameters["$top"] = autorest.Encode("query",*top)
            }
            if skip != nil {
            queryParameters["$skip"] = autorest.Encode("query",*skip)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithCustomBaseURL("{apimBaseUrl}", urlParameters),
    autorest.WithPathParameters("/products/{productId}/apis",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByProductSender sends the ListByProduct request. The method will close the
    // http.Response Body if it receives an error.
    func (client ProductAPIClient) ListByProductSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// ListByProductResponder handles the response to the ListByProduct request. The method always
// closes the http.Response Body.
func (client ProductAPIClient) ListByProductResponder(resp *http.Response) (result APICollection, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listByProductNextResults retrieves the next set of results, if any.
            func (client ProductAPIClient) listByProductNextResults(ctx context.Context, lastResults APICollection) (result APICollection, err error) {
            req, err := lastResults.aPICollectionPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "listByProductNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListByProductSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "listByProductNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListByProductResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductAPIClient", "listByProductNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListByProductComplete enumerates all values, automatically crossing page boundaries as required.
    func (client ProductAPIClient) ListByProductComplete(ctx context.Context, apimBaseURL string, productID string, filter string, top *int32, skip *int32) (result APICollectionIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/ProductAPIClient.ListByProduct")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.ListByProduct(ctx, apimBaseURL, productID, filter, top, skip)
                return
        }

