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
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// AlertsClient is the REST APIs for Azure Active Directory Connect Health
type AlertsClient struct {
	BaseClient
}

// NewAlertsClient creates an instance of the AlertsClient client.
func NewAlertsClient() AlertsClient {
	return NewAlertsClientWithBaseURI(DefaultBaseURI)
}

// NewAlertsClientWithBaseURI creates an instance of the AlertsClient client using a custom endpoint.  Use this when
// interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewAlertsClientWithBaseURI(baseURI string) AlertsClient {
	return AlertsClient{NewWithBaseURI(baseURI)}
}

// ListAddsAlerts gets the alerts for a given Active Directory Domain Service.
// Parameters:
// serviceName - the name of the service.
// filter - the alert property filter to apply.
// state - the alert state to query for.
// from - the start date to query for.
// toParameter - the end date till when to query for.
func (client AlertsClient) ListAddsAlerts(ctx context.Context, serviceName string, filter string, state string, from *date.Time, toParameter *date.Time) (result AlertsPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/AlertsClient.ListAddsAlerts")
		defer func() {
			sc := -1
			if result.a.Response.Response != nil {
				sc = result.a.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listAddsAlertsNextResults
	req, err := client.ListAddsAlertsPreparer(ctx, serviceName, filter, state, from, toParameter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "adhybridhealthservice.AlertsClient", "ListAddsAlerts", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListAddsAlertsSender(req)
	if err != nil {
		result.a.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "adhybridhealthservice.AlertsClient", "ListAddsAlerts", resp, "Failure sending request")
		return
	}

	result.a, err = client.ListAddsAlertsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "adhybridhealthservice.AlertsClient", "ListAddsAlerts", resp, "Failure responding to request")
	}
	if result.a.hasNextLink() && result.a.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListAddsAlertsPreparer prepares the ListAddsAlerts request.
func (client AlertsClient) ListAddsAlertsPreparer(ctx context.Context, serviceName string, filter string, state string, from *date.Time, toParameter *date.Time) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"serviceName": autorest.Encode("path", serviceName),
	}

	const APIVersion = "2014-01-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if len(state) > 0 {
		queryParameters["state"] = autorest.Encode("query", state)
	}
	if from != nil {
		queryParameters["from"] = autorest.Encode("query", *from)
	}
	if toParameter != nil {
		queryParameters["to"] = autorest.Encode("query", *toParameter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.ADHybridHealthService/addsservices/{serviceName}/alerts", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListAddsAlertsSender sends the ListAddsAlerts request. The method will close the
// http.Response Body if it receives an error.
func (client AlertsClient) ListAddsAlertsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListAddsAlertsResponder handles the response to the ListAddsAlerts request. The method always
// closes the http.Response Body.
func (client AlertsClient) ListAddsAlertsResponder(resp *http.Response) (result Alerts, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listAddsAlertsNextResults retrieves the next set of results, if any.
func (client AlertsClient) listAddsAlertsNextResults(ctx context.Context, lastResults Alerts) (result Alerts, err error) {
	req, err := lastResults.alertsPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "adhybridhealthservice.AlertsClient", "listAddsAlertsNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListAddsAlertsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "adhybridhealthservice.AlertsClient", "listAddsAlertsNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListAddsAlertsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "adhybridhealthservice.AlertsClient", "listAddsAlertsNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListAddsAlertsComplete enumerates all values, automatically crossing page boundaries as required.
func (client AlertsClient) ListAddsAlertsComplete(ctx context.Context, serviceName string, filter string, state string, from *date.Time, toParameter *date.Time) (result AlertsIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/AlertsClient.ListAddsAlerts")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListAddsAlerts(ctx, serviceName, filter, state, from, toParameter)
	return
}
