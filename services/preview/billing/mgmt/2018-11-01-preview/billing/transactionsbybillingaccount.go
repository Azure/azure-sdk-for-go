package billing

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

// TransactionsByBillingAccountClient is the billing client provides access to billing resources for Azure
// subscriptions.
type TransactionsByBillingAccountClient struct {
	BaseClient
}

// NewTransactionsByBillingAccountClient creates an instance of the TransactionsByBillingAccountClient client.
func NewTransactionsByBillingAccountClient(subscriptionID string) TransactionsByBillingAccountClient {
	return NewTransactionsByBillingAccountClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewTransactionsByBillingAccountClientWithBaseURI creates an instance of the TransactionsByBillingAccountClient
// client.
func NewTransactionsByBillingAccountClientWithBaseURI(baseURI string, subscriptionID string) TransactionsByBillingAccountClient {
	return TransactionsByBillingAccountClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// List lists the transactions by billingAccountId for given start and end date.
// Parameters:
// billingAccountID - billing Account Id.
// startDate - start date
// endDate - end date
// filter - may be used to filter by transaction kind. The filter supports 'eq', 'lt', 'gt', 'le', 'ge', and
// 'and'. It does not currently support 'ne', 'or', or 'not'. Tag filter is a key value pair string where key
// and value is separated by a colon (:).
func (client TransactionsByBillingAccountClient) List(ctx context.Context, billingAccountID string, startDate string, endDate string, filter string) (result TransactionsListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TransactionsByBillingAccountClient.List")
		defer func() {
			sc := -1
			if result.tlr.Response.Response != nil {
				sc = result.tlr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, billingAccountID, startDate, endDate, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.TransactionsByBillingAccountClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.tlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "billing.TransactionsByBillingAccountClient", "List", resp, "Failure sending request")
		return
	}

	result.tlr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.TransactionsByBillingAccountClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client TransactionsByBillingAccountClient) ListPreparer(ctx context.Context, billingAccountID string, startDate string, endDate string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"billingAccountId": autorest.Encode("path", billingAccountID),
	}

	const APIVersion = "2018-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
		"endDate":     autorest.Encode("query", endDate),
		"startDate":   autorest.Encode("query", startDate),
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/transactions", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client TransactionsByBillingAccountClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client TransactionsByBillingAccountClient) ListResponder(resp *http.Response) (result TransactionsListResult, err error) {
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
func (client TransactionsByBillingAccountClient) listNextResults(ctx context.Context, lastResults TransactionsListResult) (result TransactionsListResult, err error) {
	req, err := lastResults.transactionsListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "billing.TransactionsByBillingAccountClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "billing.TransactionsByBillingAccountClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.TransactionsByBillingAccountClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client TransactionsByBillingAccountClient) ListComplete(ctx context.Context, billingAccountID string, startDate string, endDate string, filter string) (result TransactionsListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TransactionsByBillingAccountClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, billingAccountID, startDate, endDate, filter)
	return
}
