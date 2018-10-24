package reservations

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

// OrderClient is the this API describe Azure Reservation
type OrderClient struct {
    BaseClient
}
// NewOrderClient creates an instance of the OrderClient client.
func NewOrderClient() OrderClient {
    return NewOrderClientWithBaseURI(DefaultBaseURI, )
}

// NewOrderClientWithBaseURI creates an instance of the OrderClient client.
    func NewOrderClientWithBaseURI(baseURI string, ) OrderClient {
        return OrderClient{ NewWithBaseURI(baseURI, )}
    }

// Get get the details of the `ReservationOrder`.
    // Parameters:
        // reservationOrderID - order Id of the reservation
func (client OrderClient) Get(ctx context.Context, reservationOrderID string) (result OrderResponse, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/OrderClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetPreparer(ctx, reservationOrderID)
    if err != nil {
    err = autorest.NewErrorWithError(err, "reservations.OrderClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "reservations.OrderClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "reservations.OrderClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client OrderClient) GetPreparer(ctx context.Context, reservationOrderID string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "reservationOrderId": autorest.Encode("path",reservationOrderID),
            }

                        const APIVersion = "2018-06-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/providers/Microsoft.Capacity/reservationOrders/{reservationOrderId}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client OrderClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client OrderClient) GetResponder(resp *http.Response) (result OrderResponse, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List list of all the `ReservationOrder`s that the user has access to in the current tenant.
func (client OrderClient) List(ctx context.Context) (result OrderListPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/OrderClient.List")
        defer func() {
            sc := -1
            if result.ol.Response.Response != nil {
                sc = result.ol.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx)
    if err != nil {
    err = autorest.NewErrorWithError(err, "reservations.OrderClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.ol.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "reservations.OrderClient", "List", resp, "Failure sending request")
            return
            }

            result.ol, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "reservations.OrderClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client OrderClient) ListPreparer(ctx context.Context) (*http.Request, error) {
                    const APIVersion = "2018-06-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPath("/providers/Microsoft.Capacity/reservationOrders"),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client OrderClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client OrderClient) ListResponder(resp *http.Response) (result OrderList, err error) {
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
            func (client OrderClient) listNextResults(ctx context.Context, lastResults OrderList) (result OrderList, err error) {
            req, err := lastResults.orderListPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "reservations.OrderClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "reservations.OrderClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "reservations.OrderClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client OrderClient) ListComplete(ctx context.Context) (result OrderListIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/OrderClient.List")
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

