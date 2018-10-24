package managementpartner

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

// PartnerClient is the this API describe ACE Provisioning ManagementPartner
type PartnerClient struct {
    BaseClient
}
// NewPartnerClient creates an instance of the PartnerClient client.
func NewPartnerClient() PartnerClient {
    return NewPartnerClientWithBaseURI(DefaultBaseURI, )
}

// NewPartnerClientWithBaseURI creates an instance of the PartnerClient client.
    func NewPartnerClientWithBaseURI(baseURI string, ) PartnerClient {
        return PartnerClient{ NewWithBaseURI(baseURI, )}
    }

// Create create a management partner for the objectId and tenantId.
    // Parameters:
        // partnerID - id of the Partner
func (client PartnerClient) Create(ctx context.Context, partnerID string) (result PartnerResponse, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/PartnerClient.Create")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.CreatePreparer(ctx, partnerID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Create", nil , "Failure preparing request")
    return
    }

            resp, err := client.CreateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Create", resp, "Failure sending request")
            return
            }

            result, err = client.CreateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Create", resp, "Failure responding to request")
            }

    return
    }

    // CreatePreparer prepares the Create request.
    func (client PartnerClient) CreatePreparer(ctx context.Context, partnerID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "partnerId": autorest.Encode("path",partnerID),
            }

                        const APIVersion = "2018-02-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/providers/Microsoft.ManagementPartner/partners/{partnerId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateSender sends the Create request. The method will close the
    // http.Response Body if it receives an error.
    func (client PartnerClient) CreateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client PartnerClient) CreateResponder(resp *http.Response) (result PartnerResponse, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete delete the management partner for the objectId and tenantId.
    // Parameters:
        // partnerID - id of the Partner
func (client PartnerClient) Delete(ctx context.Context, partnerID string) (result autorest.Response, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/PartnerClient.Delete")
        defer func() {
            sc := -1
            if result.Response != nil {
                sc = result.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.DeletePreparer(ctx, partnerID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Delete", nil , "Failure preparing request")
    return
    }

            resp, err := client.DeleteSender(req)
            if err != nil {
            result.Response = resp
            err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Delete", resp, "Failure sending request")
            return
            }

            result, err = client.DeleteResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Delete", resp, "Failure responding to request")
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client PartnerClient) DeletePreparer(ctx context.Context, partnerID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "partnerId": autorest.Encode("path",partnerID),
            }

                        const APIVersion = "2018-02-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/providers/Microsoft.ManagementPartner/partners/{partnerId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client PartnerClient) DeleteSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client PartnerClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get get the management partner using the partnerId, objectId and tenantId.
    // Parameters:
        // partnerID - id of the Partner
func (client PartnerClient) Get(ctx context.Context, partnerID string) (result PartnerResponse, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/PartnerClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetPreparer(ctx, partnerID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client PartnerClient) GetPreparer(ctx context.Context, partnerID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "partnerId": autorest.Encode("path",partnerID),
            }

                        const APIVersion = "2018-02-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/providers/Microsoft.ManagementPartner/partners/{partnerId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client PartnerClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client PartnerClient) GetResponder(resp *http.Response) (result PartnerResponse, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Update update the management partner for the objectId and tenantId.
    // Parameters:
        // partnerID - id of the Partner
func (client PartnerClient) Update(ctx context.Context, partnerID string) (result PartnerResponse, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/PartnerClient.Update")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.UpdatePreparer(ctx, partnerID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Update", nil , "Failure preparing request")
    return
    }

            resp, err := client.UpdateSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Update", resp, "Failure sending request")
            return
            }

            result, err = client.UpdateResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "managementpartner.PartnerClient", "Update", resp, "Failure responding to request")
            }

    return
    }

    // UpdatePreparer prepares the Update request.
    func (client PartnerClient) UpdatePreparer(ctx context.Context, partnerID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "partnerId": autorest.Encode("path",partnerID),
            }

                        const APIVersion = "2018-02-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsPatch(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/providers/Microsoft.ManagementPartner/partners/{partnerId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // UpdateSender sends the Update request. The method will close the
    // http.Response Body if it receives an error.
    func (client PartnerClient) UpdateSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client PartnerClient) UpdateResponder(resp *http.Response) (result PartnerResponse, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

