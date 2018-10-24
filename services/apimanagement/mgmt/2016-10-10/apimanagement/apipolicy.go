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
    "io"
)

// APIPolicyClient is the apiManagement Client
type APIPolicyClient struct {
    BaseClient
}
// NewAPIPolicyClient creates an instance of the APIPolicyClient client.
func NewAPIPolicyClient(subscriptionID string) APIPolicyClient {
    return NewAPIPolicyClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewAPIPolicyClientWithBaseURI creates an instance of the APIPolicyClient client.
    func NewAPIPolicyClientWithBaseURI(baseURI string, subscriptionID string) APIPolicyClient {
        return APIPolicyClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate creates or updates policy configuration for the API.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // serviceName - the name of the API Management service.
        // apiid - API identifier. Must be unique in the current API Management service instance.
        // parameters - the policy contents to apply.
        // ifMatch - the entity state (Etag) version of the Api Policy to update. A value of "*" can be used for
        // If-Match to unconditionally apply the operation.
func (client APIPolicyClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, apiid string, parameters io.ReadCloser, ifMatch string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/APIPolicyClient.CreateOrUpdate")
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
            { TargetValue: apiid,
             Constraints: []validation.Constraint{	{Target: "apiid", Name: validation.MaxLength, Rule: 256, Chain: nil },
            	{Target: "apiid", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "apiid", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.APIPolicyClient", "CreateOrUpdate", err.Error())
            }

                req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, serviceName, apiid, parameters, ifMatch)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "CreateOrUpdate", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "CreateOrUpdate", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client APIPolicyClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, serviceName string, apiid string, parameters io.ReadCloser, ifMatch string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "apiId": autorest.Encode("path",apiid),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serviceName": autorest.Encode("path",serviceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-10-10"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/vnd.ms-azure-apim.policy+xml"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/apis/{apiId}/policy",pathParameters),
    autorest.WithFile(parameters),
    autorest.WithQueryParameters(queryParameters),
    autorest.WithHeader("If-Match", autorest.String(ifMatch)))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client APIPolicyClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client APIPolicyClient) CreateOrUpdateResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Delete deletes the policy configuration at the Api.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // serviceName - the name of the API Management service.
        // apiid - API identifier. Must be unique in the current API Management service instance.
        // ifMatch - the entity state (Etag) version of the Api policy to update. A value of "*" can be used for
        // If-Match to unconditionally apply the operation.
func (client APIPolicyClient) Delete(ctx context.Context, resourceGroupName string, serviceName string, apiid string, ifMatch string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/APIPolicyClient.Delete")
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
            { TargetValue: apiid,
             Constraints: []validation.Constraint{	{Target: "apiid", Name: validation.MaxLength, Rule: 256, Chain: nil },
            	{Target: "apiid", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "apiid", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.APIPolicyClient", "Delete", err.Error())
            }

                req, err := client.DeletePreparer(ctx, resourceGroupName, serviceName, apiid, ifMatch)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client APIPolicyClient) DeletePreparer(ctx context.Context, resourceGroupName string, serviceName string, apiid string, ifMatch string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "apiId": autorest.Encode("path",apiid),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serviceName": autorest.Encode("path",serviceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-10-10"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/apis/{apiId}/policy",pathParameters),
    autorest.WithQueryParameters(queryParameters),
    autorest.WithHeader("If-Match", autorest.String(ifMatch)))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client APIPolicyClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client APIPolicyClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get get the policy configuration at the API level.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // serviceName - the name of the API Management service.
        // apiid - API identifier. Must be unique in the current API Management service instance.
func (client APIPolicyClient) Get(ctx context.Context, resourceGroupName string, serviceName string, apiid string) (result ReadCloser, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/APIPolicyClient.Get")
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
            { TargetValue: apiid,
             Constraints: []validation.Constraint{	{Target: "apiid", Name: validation.MaxLength, Rule: 256, Chain: nil },
            	{Target: "apiid", Name: validation.MinLength, Rule: 1, Chain: nil },
            	{Target: "apiid", Name: validation.Pattern, Rule: `^[^*#&+:<>?]+$`, Chain: nil }}}}); err != nil {
            return result, validation.NewError("apimanagement.APIPolicyClient", "Get", err.Error())
            }

                req, err := client.GetPreparer(ctx, resourceGroupName, serviceName, apiid)
    if err != nil {
    err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "apimanagement.APIPolicyClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client APIPolicyClient) GetPreparer(ctx context.Context, resourceGroupName string, serviceName string, apiid string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "apiId": autorest.Encode("path",apiid),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serviceName": autorest.Encode("path",serviceName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-10-10"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/apis/{apiId}/policy",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client APIPolicyClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client APIPolicyClient) GetResponder(resp *http.Response) (result ReadCloser, err error) {
    result.Value = &resp.Body
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK))
    result.Response = autorest.Response{Response: resp}
        return
    }

