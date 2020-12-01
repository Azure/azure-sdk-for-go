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

// LotsClient is the consumption management client provides access to consumption resources for Azure Enterprise
// Subscriptions.
type LotsClient struct {
	BaseClient
}

// NewLotsClient creates an instance of the LotsClient client.
func NewLotsClient(subscriptionID string) LotsClient {
	return NewLotsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewLotsClientWithBaseURI creates an instance of the LotsClient client using a custom endpoint.  Use this when
// interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewLotsClientWithBaseURI(baseURI string, subscriptionID string) LotsClient {
	return LotsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// List lists the lots by billingAccountId and billingProfileId.
// Parameters:
// billingAccountID - billingAccount ID
// billingProfileID - azure Billing Profile ID.
func (client LotsClient) List(ctx context.Context, billingAccountID string, billingProfileID string) (result LotsPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/LotsClient.List")
		defer func() {
			sc := -1
			if result.l.Response.Response != nil {
				sc = result.l.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, billingAccountID, billingProfileID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.LotsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.l.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "consumption.LotsClient", "List", resp, "Failure sending request")
		return
	}

	result.l, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.LotsClient", "List", resp, "Failure responding to request")
	}
	if result.l.hasNextLink() && result.l.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListPreparer prepares the List request.
func (client LotsClient) ListPreparer(ctx context.Context, billingAccountID string, billingProfileID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"billingAccountId": autorest.Encode("path", billingAccountID),
		"billingProfileId": autorest.Encode("path", billingProfileID),
	}

	const APIVersion = "2019-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}/providers/Microsoft.Consumption/lots", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client LotsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client LotsClient) ListResponder(resp *http.Response) (result Lots, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client LotsClient) listNextResults(ctx context.Context, lastResults Lots) (result Lots, err error) {
	req, err := lastResults.lotsPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "consumption.LotsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "consumption.LotsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.LotsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client LotsClient) ListComplete(ctx context.Context, billingAccountID string, billingProfileID string) (result LotsIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/LotsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, billingAccountID, billingProfileID)
	return
}
