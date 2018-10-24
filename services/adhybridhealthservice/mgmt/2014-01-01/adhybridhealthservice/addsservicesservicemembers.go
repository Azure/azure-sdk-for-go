package adhybridhealthservice

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

// AddsServicesServiceMembersClient is the REST APIs for Azure Active Drectory Connect Health
type AddsServicesServiceMembersClient struct {
    BaseClient
}
// NewAddsServicesServiceMembersClient creates an instance of the AddsServicesServiceMembersClient client.
func NewAddsServicesServiceMembersClient() AddsServicesServiceMembersClient {
    return NewAddsServicesServiceMembersClientWithBaseURI(DefaultBaseURI, )
}

// NewAddsServicesServiceMembersClientWithBaseURI creates an instance of the AddsServicesServiceMembersClient client.
    func NewAddsServicesServiceMembersClientWithBaseURI(baseURI string, ) AddsServicesServiceMembersClient {
        return AddsServicesServiceMembersClient{ NewWithBaseURI(baseURI, )}
    }

// Add onboards  a server, for a given Active Directory Domain Controller service, to Azure Active Directory Connect
// Health Service.
    // Parameters:
        // serviceName - the name of the service under which the server is to be onboarded.
        // serviceMember - the server object.
func (client AddsServicesServiceMembersClient) Add(ctx context.Context, serviceName string, serviceMember ServiceMember) (result ServiceMember, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AddsServicesServiceMembersClient.Add")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.AddPreparer(ctx, serviceName, serviceMember)
    if err != nil {
    err = autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "Add", nil , "Failure preparing request")
    return
    }

            resp, err := client.AddSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "Add", resp, "Failure sending request")
            return
            }

            result, err = client.AddResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "Add", resp, "Failure responding to request")
            }

    return
    }

    // AddPreparer prepares the Add request.
    func (client AddsServicesServiceMembersClient) AddPreparer(ctx context.Context, serviceName string, serviceMember ServiceMember) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "serviceName": autorest.Encode("path",serviceName),
            }

                        const APIVersion = "2014-01-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPost(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/providers/Microsoft.ADHybridHealthService/addsservices/{serviceName}/servicemembers",pathParameters),
    autorest.WithJSON(serviceMember),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // AddSender sends the Add request. The method will close the
    // http.Response Body if it receives an error.
    func (client AddsServicesServiceMembersClient) AddSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// AddResponder handles the response to the Add request. The method always
// closes the http.Response Body.
func (client AddsServicesServiceMembersClient) AddResponder(resp *http.Response) (result ServiceMember, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List gets the details of the servers, for a given Active Directory Domain Controller service, that are onboarded to
// Azure Active Directory Connect Health Service.
    // Parameters:
        // serviceName - the name of the service.
        // filter - the server property filter to apply.
        // dimensionType - the server specific dimension.
        // dimensionSignature - the value of the dimension.
func (client AddsServicesServiceMembersClient) List(ctx context.Context, serviceName string, filter string, dimensionType string, dimensionSignature string) (result ServiceMembersPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/AddsServicesServiceMembersClient.List")
        defer func() {
            sc := -1
            if result.sm.Response.Response != nil {
                sc = result.sm.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx, serviceName, filter, dimensionType, dimensionSignature)
    if err != nil {
    err = autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.sm.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "List", resp, "Failure sending request")
            return
            }

            result.sm, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client AddsServicesServiceMembersClient) ListPreparer(ctx context.Context, serviceName string, filter string, dimensionType string, dimensionSignature string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "serviceName": autorest.Encode("path",serviceName),
            }

                        const APIVersion = "2014-01-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if len(filter) > 0 {
            queryParameters["$filter"] = autorest.Encode("query",filter)
            }
            if len(dimensionType) > 0 {
            queryParameters["dimensionType"] = autorest.Encode("query",dimensionType)
            }
            if len(dimensionSignature) > 0 {
            queryParameters["dimensionSignature"] = autorest.Encode("query",dimensionSignature)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/providers/Microsoft.ADHybridHealthService/addsservices/{serviceName}/servicemembers",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client AddsServicesServiceMembersClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client AddsServicesServiceMembersClient) ListResponder(resp *http.Response) (result ServiceMembers, err error) {
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
            func (client AddsServicesServiceMembersClient) listNextResults(ctx context.Context, lastResults ServiceMembers) (result ServiceMembers, err error) {
            req, err := lastResults.serviceMembersPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "adhybridhealthservice.AddsServicesServiceMembersClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client AddsServicesServiceMembersClient) ListComplete(ctx context.Context, serviceName string, filter string, dimensionType string, dimensionSignature string) (result ServiceMembersIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/AddsServicesServiceMembersClient.List")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.List(ctx, serviceName, filter, dimensionType, dimensionSignature)
                return
        }

