package consumption

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
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// ReservationsSummariesClient is the consumption management client provides access to consumption resources for Azure
// Enterprise Subscriptions.
type ReservationsSummariesClient struct {
	BaseClient
}

// NewReservationsSummariesClient creates an instance of the ReservationsSummariesClient client.
func NewReservationsSummariesClient(subscriptionID string) ReservationsSummariesClient {
	return NewReservationsSummariesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewReservationsSummariesClientWithBaseURI creates an instance of the ReservationsSummariesClient client using a
// custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds,
// Azure stack).
func NewReservationsSummariesClientWithBaseURI(baseURI string, subscriptionID string) ReservationsSummariesClient {
	return ReservationsSummariesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ListByReservationOrder lists the reservations summaries for daily or monthly grain.
// Parameters:
// reservationOrderID - order Id of the reservation
// grain - can be daily or monthly
// filter - required only for daily grain. The properties/UsageDate for start date and end date. The filter
// supports 'le' and  'ge'
func (client ReservationsSummariesClient) ListByReservationOrder(ctx context.Context, reservationOrderID string, grain Datagrain, filter string) (result ReservationSummariesListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ReservationsSummariesClient.ListByReservationOrder")
		defer func() {
			sc := -1
			if result.rslr.Response.Response != nil {
				sc = result.rslr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listByReservationOrderNextResults
	req, err := client.ListByReservationOrderPreparer(ctx, reservationOrderID, grain, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "ListByReservationOrder", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByReservationOrderSender(req)
	if err != nil {
		result.rslr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "ListByReservationOrder", resp, "Failure sending request")
		return
	}

	result.rslr, err = client.ListByReservationOrderResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "ListByReservationOrder", resp, "Failure responding to request")
	}

	return
}

// ListByReservationOrderPreparer prepares the ListByReservationOrder request.
func (client ReservationsSummariesClient) ListByReservationOrderPreparer(ctx context.Context, reservationOrderID string, grain Datagrain, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"reservationOrderId": autorest.Encode("path", reservationOrderID),
	}

	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
		"grain":       autorest.Encode("query", grain),
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Capacity/reservationorders/{reservationOrderId}/providers/Microsoft.Consumption/reservationSummaries", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByReservationOrderSender sends the ListByReservationOrder request. The method will close the
// http.Response Body if it receives an error.
func (client ReservationsSummariesClient) ListByReservationOrderSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListByReservationOrderResponder handles the response to the ListByReservationOrder request. The method always
// closes the http.Response Body.
func (client ReservationsSummariesClient) ListByReservationOrderResponder(resp *http.Response) (result ReservationSummariesListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByReservationOrderNextResults retrieves the next set of results, if any.
func (client ReservationsSummariesClient) listByReservationOrderNextResults(ctx context.Context, lastResults ReservationSummariesListResult) (result ReservationSummariesListResult, err error) {
	req, err := lastResults.reservationSummariesListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "listByReservationOrderNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByReservationOrderSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "listByReservationOrderNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByReservationOrderResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "listByReservationOrderNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByReservationOrderComplete enumerates all values, automatically crossing page boundaries as required.
func (client ReservationsSummariesClient) ListByReservationOrderComplete(ctx context.Context, reservationOrderID string, grain Datagrain, filter string) (result ReservationSummariesListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ReservationsSummariesClient.ListByReservationOrder")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByReservationOrder(ctx, reservationOrderID, grain, filter)
	return
}

// ListByReservationOrderAndReservation lists the reservations summaries for daily or monthly grain.
// Parameters:
// reservationOrderID - order Id of the reservation
// reservationID - id of the reservation
// grain - can be daily or monthly
// filter - required only for daily grain. The properties/UsageDate for start date and end date. The filter
// supports 'le' and  'ge'
func (client ReservationsSummariesClient) ListByReservationOrderAndReservation(ctx context.Context, reservationOrderID string, reservationID string, grain Datagrain, filter string) (result ReservationSummariesListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ReservationsSummariesClient.ListByReservationOrderAndReservation")
		defer func() {
			sc := -1
			if result.rslr.Response.Response != nil {
				sc = result.rslr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listByReservationOrderAndReservationNextResults
	req, err := client.ListByReservationOrderAndReservationPreparer(ctx, reservationOrderID, reservationID, grain, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "ListByReservationOrderAndReservation", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByReservationOrderAndReservationSender(req)
	if err != nil {
		result.rslr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "ListByReservationOrderAndReservation", resp, "Failure sending request")
		return
	}

	result.rslr, err = client.ListByReservationOrderAndReservationResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "ListByReservationOrderAndReservation", resp, "Failure responding to request")
	}

	return
}

// ListByReservationOrderAndReservationPreparer prepares the ListByReservationOrderAndReservation request.
func (client ReservationsSummariesClient) ListByReservationOrderAndReservationPreparer(ctx context.Context, reservationOrderID string, reservationID string, grain Datagrain, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"reservationId":      autorest.Encode("path", reservationID),
		"reservationOrderId": autorest.Encode("path", reservationOrderID),
	}

	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
		"grain":       autorest.Encode("query", grain),
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Capacity/reservationorders/{reservationOrderId}/reservations/{reservationId}/providers/Microsoft.Consumption/reservationSummaries", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByReservationOrderAndReservationSender sends the ListByReservationOrderAndReservation request. The method will close the
// http.Response Body if it receives an error.
func (client ReservationsSummariesClient) ListByReservationOrderAndReservationSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListByReservationOrderAndReservationResponder handles the response to the ListByReservationOrderAndReservation request. The method always
// closes the http.Response Body.
func (client ReservationsSummariesClient) ListByReservationOrderAndReservationResponder(resp *http.Response) (result ReservationSummariesListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByReservationOrderAndReservationNextResults retrieves the next set of results, if any.
func (client ReservationsSummariesClient) listByReservationOrderAndReservationNextResults(ctx context.Context, lastResults ReservationSummariesListResult) (result ReservationSummariesListResult, err error) {
	req, err := lastResults.reservationSummariesListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "listByReservationOrderAndReservationNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByReservationOrderAndReservationSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "listByReservationOrderAndReservationNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByReservationOrderAndReservationResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.ReservationsSummariesClient", "listByReservationOrderAndReservationNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByReservationOrderAndReservationComplete enumerates all values, automatically crossing page boundaries as required.
func (client ReservationsSummariesClient) ListByReservationOrderAndReservationComplete(ctx context.Context, reservationOrderID string, reservationID string, grain Datagrain, filter string) (result ReservationSummariesListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ReservationsSummariesClient.ListByReservationOrderAndReservation")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByReservationOrderAndReservation(ctx, reservationOrderID, reservationID, grain, filter)
	return
}
