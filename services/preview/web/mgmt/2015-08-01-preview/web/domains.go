package web

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

// DomainsClient is the webSite Management Client
type DomainsClient struct {
    BaseClient
}
// NewDomainsClient creates an instance of the DomainsClient client.
func NewDomainsClient(subscriptionID string) DomainsClient {
    return NewDomainsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDomainsClientWithBaseURI creates an instance of the DomainsClient client.
    func NewDomainsClientWithBaseURI(baseURI string, subscriptionID string) DomainsClient {
        return DomainsClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdateDomain sends the create or update domain request.
    // Parameters:
        // resourceGroupName - &gt;Name of the resource group
        // domainName - name of the domain
        // domain - domain registration information
func (client DomainsClient) CreateOrUpdateDomain(ctx context.Context, resourceGroupName string, domainName string, domain Domain) (result Domain, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DomainsClient.CreateOrUpdateDomain")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.CreateOrUpdateDomainPreparer(ctx, resourceGroupName, domainName, domain)
    if err != nil {
    err = autorest.NewErrorWithError(err, "web.DomainsClient", "CreateOrUpdateDomain", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateOrUpdateDomainSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "CreateOrUpdateDomain", resp, "Failure sending request")
            return
            }

            result, err = client.CreateOrUpdateDomainResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "CreateOrUpdateDomain", resp, "Failure responding to request")
            }

    return
    }

    // CreateOrUpdateDomainPreparer prepares the CreateOrUpdateDomain request.
    func (client DomainsClient) CreateOrUpdateDomainPreparer(ctx context.Context, resourceGroupName string, domainName string, domain Domain) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "domainName": autorest.Encode("path",domainName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-08-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DomainRegistration/domains/{domainName}",pathParameters),
    autorest.WithJSON(domain),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateDomainSender sends the CreateOrUpdateDomain request. The method will close the
    // http.Response Body if it receives an error.
    func (client DomainsClient) CreateOrUpdateDomainSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// CreateOrUpdateDomainResponder handles the response to the CreateOrUpdateDomain request. The method always
// closes the http.Response Body.
func (client DomainsClient) CreateOrUpdateDomainResponder(resp *http.Response) (result Domain, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// DeleteDomain sends the delete domain request.
    // Parameters:
        // resourceGroupName - name of the resource group
        // domainName - name of the domain
        // forceHardDeleteDomain - if true then the domain will be deleted immediately instead of after 24 hours
func (client DomainsClient) DeleteDomain(ctx context.Context, resourceGroupName string, domainName string, forceHardDeleteDomain *bool) (result SetObject, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DomainsClient.DeleteDomain")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.DeleteDomainPreparer(ctx, resourceGroupName, domainName, forceHardDeleteDomain)
    if err != nil {
    err = autorest.NewErrorWithError(err, "web.DomainsClient", "DeleteDomain", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteDomainSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "DeleteDomain", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteDomainResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "DeleteDomain", resp, "Failure responding to request")
            }

    return
    }

    // DeleteDomainPreparer prepares the DeleteDomain request.
    func (client DomainsClient) DeleteDomainPreparer(ctx context.Context, resourceGroupName string, domainName string, forceHardDeleteDomain *bool) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "domainName": autorest.Encode("path",domainName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-08-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if forceHardDeleteDomain != nil {
            queryParameters["forceHardDeleteDomain"] = autorest.Encode("query",*forceHardDeleteDomain)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DomainRegistration/domains/{domainName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteDomainSender sends the DeleteDomain request. The method will close the
    // http.Response Body if it receives an error.
    func (client DomainsClient) DeleteDomainSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// DeleteDomainResponder handles the response to the DeleteDomain request. The method always
// closes the http.Response Body.
func (client DomainsClient) DeleteDomainResponder(resp *http.Response) (result SetObject, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusNoContent),
    autorest.ByUnmarshallingJSON(&result.Value),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// GetDomain sends the get domain request.
    // Parameters:
        // resourceGroupName - name of the resource group
        // domainName - name of the domain
func (client DomainsClient) GetDomain(ctx context.Context, resourceGroupName string, domainName string) (result Domain, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DomainsClient.GetDomain")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetDomainPreparer(ctx, resourceGroupName, domainName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomain", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetDomainSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomain", resp, "Failure sending request")
            return
            }

            result, err = client.GetDomainResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomain", resp, "Failure responding to request")
            }

    return
    }

    // GetDomainPreparer prepares the GetDomain request.
    func (client DomainsClient) GetDomainPreparer(ctx context.Context, resourceGroupName string, domainName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "domainName": autorest.Encode("path",domainName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-08-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DomainRegistration/domains/{domainName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetDomainSender sends the GetDomain request. The method will close the
    // http.Response Body if it receives an error.
    func (client DomainsClient) GetDomainSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetDomainResponder handles the response to the GetDomain request. The method always
// closes the http.Response Body.
func (client DomainsClient) GetDomainResponder(resp *http.Response) (result Domain, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// GetDomainOperation sends the get domain operation request.
    // Parameters:
        // resourceGroupName - name of the resource group
        // domainName - name of the domain
        // operationID - domain purchase operation Id
func (client DomainsClient) GetDomainOperation(ctx context.Context, resourceGroupName string, domainName string, operationID string) (result Domain, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DomainsClient.GetDomainOperation")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetDomainOperationPreparer(ctx, resourceGroupName, domainName, operationID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomainOperation", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetDomainOperationSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomainOperation", resp, "Failure sending request")
            return
            }

            result, err = client.GetDomainOperationResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomainOperation", resp, "Failure responding to request")
            }

    return
    }

    // GetDomainOperationPreparer prepares the GetDomainOperation request.
    func (client DomainsClient) GetDomainOperationPreparer(ctx context.Context, resourceGroupName string, domainName string, operationID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "domainName": autorest.Encode("path",domainName),
            "operationId": autorest.Encode("path",operationID),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-08-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DomainRegistration/domains/{domainName}/operationresults/{operationId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetDomainOperationSender sends the GetDomainOperation request. The method will close the
    // http.Response Body if it receives an error.
    func (client DomainsClient) GetDomainOperationSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetDomainOperationResponder handles the response to the GetDomainOperation request. The method always
// closes the http.Response Body.
func (client DomainsClient) GetDomainOperationResponder(resp *http.Response) (result Domain, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted,http.StatusInternalServerError),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// GetDomains sends the get domains request.
    // Parameters:
        // resourceGroupName - name of the resource group
func (client DomainsClient) GetDomains(ctx context.Context, resourceGroupName string) (result DomainCollectionPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DomainsClient.GetDomains")
        defer func() {
            sc := -1
            if result.dc.Response.Response != nil {
                sc = result.dc.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.getDomainsNextResults
    req, err := client.GetDomainsPreparer(ctx, resourceGroupName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomains", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetDomainsSender(req)
            if err != nil {
            result.dc.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomains", resp, "Failure sending request")
            return
            }

            result.dc, err = client.GetDomainsResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "GetDomains", resp, "Failure responding to request")
            }

    return
    }

    // GetDomainsPreparer prepares the GetDomains request.
    func (client DomainsClient) GetDomainsPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-08-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DomainRegistration/domains",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetDomainsSender sends the GetDomains request. The method will close the
    // http.Response Body if it receives an error.
    func (client DomainsClient) GetDomainsSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetDomainsResponder handles the response to the GetDomains request. The method always
// closes the http.Response Body.
func (client DomainsClient) GetDomainsResponder(resp *http.Response) (result DomainCollection, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // getDomainsNextResults retrieves the next set of results, if any.
            func (client DomainsClient) getDomainsNextResults(ctx context.Context, lastResults DomainCollection) (result DomainCollection, err error) {
            req, err := lastResults.domainCollectionPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "web.DomainsClient", "getDomainsNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.GetDomainsSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "web.DomainsClient", "getDomainsNextResults", resp, "Failure sending next results request")
            }
            result, err = client.GetDomainsResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "getDomainsNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // GetDomainsComplete enumerates all values, automatically crossing page boundaries as required.
    func (client DomainsClient) GetDomainsComplete(ctx context.Context, resourceGroupName string) (result DomainCollectionIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/DomainsClient.GetDomains")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.GetDomains(ctx, resourceGroupName)
                return
        }

// UpdateDomain sends the update domain request.
    // Parameters:
        // resourceGroupName - &gt;Name of the resource group
        // domainName - name of the domain
        // domain - domain registration information
func (client DomainsClient) UpdateDomain(ctx context.Context, resourceGroupName string, domainName string, domain Domain) (result Domain, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DomainsClient.UpdateDomain")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.UpdateDomainPreparer(ctx, resourceGroupName, domainName, domain)
    if err != nil {
    err = autorest.NewErrorWithError(err, "web.DomainsClient", "UpdateDomain", nil , "Failure preparing request")
    return
    }

            resp, err := client.UpdateDomainSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "UpdateDomain", resp, "Failure sending request")
            return
            }

            result, err = client.UpdateDomainResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "web.DomainsClient", "UpdateDomain", resp, "Failure responding to request")
            }

    return
    }

    // UpdateDomainPreparer prepares the UpdateDomain request.
    func (client DomainsClient) UpdateDomainPreparer(ctx context.Context, resourceGroupName string, domainName string, domain Domain) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "domainName": autorest.Encode("path",domainName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-08-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPatch(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DomainRegistration/domains/{domainName}",pathParameters),
    autorest.WithJSON(domain),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // UpdateDomainSender sends the UpdateDomain request. The method will close the
    // http.Response Body if it receives an error.
    func (client DomainsClient) UpdateDomainSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// UpdateDomainResponder handles the response to the UpdateDomain request. The method always
// closes the http.Response Body.
func (client DomainsClient) UpdateDomainResponder(resp *http.Response) (result Domain, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

