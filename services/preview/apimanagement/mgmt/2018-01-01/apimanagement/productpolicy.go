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

// ProductPolicyClient is the apiManagement Client
type ProductPolicyClient struct {
    BaseClient
}
// NewProductPolicyClient creates an instance of the ProductPolicyClient client.
func NewProductPolicyClient(subscriptionID string) ProductPolicyClient {
    return NewProductPolicyClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewProductPolicyClientWithBaseURI creates an instance of the ProductPolicyClient client.
    func NewProductPolicyClientWithBaseURI(baseURI string, subscriptionID string) ProductPolicyClient {
        return ProductPolicyClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate creates or updates policy configuration for the Product.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // serviceName - the name of the API Management service.
        // productID - product identifier. Must be unique in the current API Management service instance.
        // parameters - the policy contents to apply.
        // ifMatch - eTag of the Entity. Not required when creating an entity, but required when updating an entity.
func (client ProductPolicyClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, productID string, parameters PolicyContract, ifMatch string) (result PolicyContract, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ProductPolicyClient.CreateOrUpdate")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: serviceName,
             Constraints: []validation.Constraint{	{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil },
            	{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil }}},
            { TargetValue: productID,
             Constraints: []validation.Constraint{	{Target: "productID", Name: validation.MaxLength, Rule: 80, Chain: nil },
            	{Target: "productID", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "productID", Name: validation.Pattern, Rule: `(^[\w]+$)|(^[\w][\w\-]+[\w]$)`, Chain: nil }}},
            { TargetValue: parameters,
             Constraints: []validation.Constraint{	{Target: "parameters.PolicyContractProperties", Name: validation.Null, Rule: false ,
            Chain: []validation.Constraint{	{Target: "parameters.PolicyContractProperties.PolicyContent", Name: validation.Null, Rule: true, Chain: nil },
            }}}}}); err != nil {
            return result, validation.NewError("apimanagement.ProductPolicyClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, serviceName, productID, parameters, ifMatch)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "CreateOrUpdate", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "CreateOrUpdate", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client ProductPolicyClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, serviceName string, productID string, parameters PolicyContract, ifMatch string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "policyId": autorest.Encode("path", "policy"),
            "productId": autorest.Encode("path",productID),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serviceName": autorest.Encode("path",serviceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-01-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/products/{productId}/policies/{policyId}",pathParameters),
    autorest.WithJSON(parameters),
    autorest.WithQueryParameters(queryParameters))
            if len(ifMatch) > 0 {
            preparer = autorest.DecoratePreparer(preparer,
            autorest.WithHeader("If-Match",autorest.String(ifMatch)))
            }
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client ProductPolicyClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client ProductPolicyClient) CreateOrUpdateResponder(resp *http.Response) (result PolicyContract, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete deletes the policy configuration at the Product.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // serviceName - the name of the API Management service.
        // productID - product identifier. Must be unique in the current API Management service instance.
        // ifMatch - eTag of the Entity. ETag should match the current entity state from the header response of the GET
        // request or it should be * for unconditional update.
func (client ProductPolicyClient) Delete(ctx context.Context, resourceGroupName string, serviceName string, productID string, ifMatch string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ProductPolicyClient.Delete")
        defer func() {
            sc := -1
            if result.Response != nil {
                sc = result.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: serviceName,
             Constraints: []validation.Constraint{	{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil },
            	{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil }}},
            { TargetValue: productID,
             Constraints: []validation.Constraint{	{Target: "productID", Name: validation.MaxLength, Rule: 80, Chain: nil },
            	{Target: "productID", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "productID", Name: validation.Pattern, Rule: `(^[\w]+$)|(^[\w][\w\-]+[\w]$)`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.ProductPolicyClient", "Delete", err.Error())
            }

                req, err := client.DeletePreparer(ctx, resourceGroupName, serviceName, productID, ifMatch)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client ProductPolicyClient) DeletePreparer(ctx context.Context, resourceGroupName string, serviceName string, productID string, ifMatch string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "policyId": autorest.Encode("path", "policy"),
            "productId": autorest.Encode("path",productID),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serviceName": autorest.Encode("path",serviceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-01-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/products/{productId}/policies/{policyId}",pathParameters),
    autorest.WithQueryParameters(queryParameters),
    autorest.WithHeader("If-Match", autorest.String(ifMatch)))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client ProductPolicyClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ProductPolicyClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get get the policy configuration at the Product level.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // serviceName - the name of the API Management service.
        // productID - product identifier. Must be unique in the current API Management service instance.
func (client ProductPolicyClient) Get(ctx context.Context, resourceGroupName string, serviceName string, productID string) (result PolicyContract, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ProductPolicyClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: serviceName,
             Constraints: []validation.Constraint{	{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil },
            	{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil }}},
            { TargetValue: productID,
             Constraints: []validation.Constraint{	{Target: "productID", Name: validation.MaxLength, Rule: 80, Chain: nil },
            	{Target: "productID", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "productID", Name: validation.Pattern, Rule: `(^[\w]+$)|(^[\w][\w\-]+[\w]$)`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.ProductPolicyClient", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, serviceName, productID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client ProductPolicyClient) GetPreparer(ctx context.Context, resourceGroupName string, serviceName string, productID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "policyId": autorest.Encode("path", "policy"),
            "productId": autorest.Encode("path",productID),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serviceName": autorest.Encode("path",serviceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-01-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/products/{productId}/policies/{policyId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client ProductPolicyClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ProductPolicyClient) GetResponder(resp *http.Response) (result PolicyContract, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// GetEntityTag get the ETag of the policy configuration at the Product level.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // serviceName - the name of the API Management service.
        // productID - product identifier. Must be unique in the current API Management service instance.
func (client ProductPolicyClient) GetEntityTag(ctx context.Context, resourceGroupName string, serviceName string, productID string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ProductPolicyClient.GetEntityTag")
        defer func() {
            sc := -1
            if result.Response != nil {
                sc = result.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: serviceName,
             Constraints: []validation.Constraint{	{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil },
            	{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil }}},
            { TargetValue: productID,
             Constraints: []validation.Constraint{	{Target: "productID", Name: validation.MaxLength, Rule: 80, Chain: nil },
            	{Target: "productID", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "productID", Name: validation.Pattern, Rule: `(^[\w]+$)|(^[\w][\w\-]+[\w]$)`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.ProductPolicyClient", "GetEntityTag", err.Error())
            }

                req, err := client.GetEntityTagPreparer(ctx, resourceGroupName, serviceName, productID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "GetEntityTag", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetEntityTagSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "GetEntityTag", resp, "Failure sending request")
            return
            }

            result, err = client.GetEntityTagResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "GetEntityTag", resp, "Failure responding to request")
            }

    return
    }

    // GetEntityTagPreparer prepares the GetEntityTag request.
    func (client ProductPolicyClient) GetEntityTagPreparer(ctx context.Context, resourceGroupName string, serviceName string, productID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "policyId": autorest.Encode("path", "policy"),
            "productId": autorest.Encode("path",productID),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serviceName": autorest.Encode("path",serviceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-01-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsHead(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/products/{productId}/policies/{policyId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetEntityTagSender sends the GetEntityTag request. The method will close the
    // http.Response Body if it receives an error.
    func (client ProductPolicyClient) GetEntityTagSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetEntityTagResponder handles the response to the GetEntityTag request. The method always
// closes the http.Response Body.
func (client ProductPolicyClient) GetEntityTagResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// ListByProduct get the policy configuration at the Product level.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // serviceName - the name of the API Management service.
        // productID - product identifier. Must be unique in the current API Management service instance.
func (client ProductPolicyClient) ListByProduct(ctx context.Context, resourceGroupName string, serviceName string, productID string) (result PolicyCollection, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ProductPolicyClient.ListByProduct")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
            if err := validation.Validate([]validation.Validation{
            { TargetValue: serviceName,
             Constraints: []validation.Constraint{	{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil },
            	{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil }}},
            { TargetValue: productID,
             Constraints: []validation.Constraint{	{Target: "productID", Name: validation.MaxLength, Rule: 80, Chain: nil },
            	{Target: "productID", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "productID", Name: validation.Pattern, Rule: `(^[\w]+$)|(^[\w][\w\-]+[\w]$)`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.ProductPolicyClient", "ListByProduct", err.Error())
            }

                req, err := client.ListByProductPreparer(ctx, resourceGroupName, serviceName, productID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "ListByProduct", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByProductSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "ListByProduct", resp, "Failure sending request")
            return
            }

            result, err = client.ListByProductResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.ProductPolicyClient", "ListByProduct", resp, "Failure responding to request")
            }

    return
    }

    // ListByProductPreparer prepares the ListByProduct request.
    func (client ProductPolicyClient) ListByProductPreparer(ctx context.Context, resourceGroupName string, serviceName string, productID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "productId": autorest.Encode("path",productID),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serviceName": autorest.Encode("path",serviceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2018-01-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/products/{productId}/policies",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByProductSender sends the ListByProduct request. The method will close the
    // http.Response Body if it receives an error.
    func (client ProductPolicyClient) ListByProductSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByProductResponder handles the response to the ListByProduct request. The method always
// closes the http.Response Body.
func (client ProductPolicyClient) ListByProductResponder(resp *http.Response) (result PolicyCollection, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

