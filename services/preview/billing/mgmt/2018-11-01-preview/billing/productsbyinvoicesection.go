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

// ProductsByInvoiceSectionClient is the billing client provides access to billing resources for Azure subscriptions.
type ProductsByInvoiceSectionClient struct {
	BaseClient
}

// NewProductsByInvoiceSectionClient creates an instance of the ProductsByInvoiceSectionClient client.
func NewProductsByInvoiceSectionClient(subscriptionID string) ProductsByInvoiceSectionClient {
	return NewProductsByInvoiceSectionClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewProductsByInvoiceSectionClientWithBaseURI creates an instance of the ProductsByInvoiceSectionClient client.
func NewProductsByInvoiceSectionClientWithBaseURI(baseURI string, subscriptionID string) ProductsByInvoiceSectionClient {
	return ProductsByInvoiceSectionClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// List lists products by invoiceSectionId.
// Parameters:
// billingAccountID - billing Account Id.
// invoiceSectionID - invoiceSection Id.
// filter - may be used to filter by product type. The filter supports 'eq', 'lt', 'gt', 'le', 'ge', and 'and'.
// It does not currently support 'ne', 'or', or 'not'. Tag filter is a key value pair string where key and
// value is separated by a colon (:).
func (client ProductsByInvoiceSectionClient) List(ctx context.Context, billingAccountID string, invoiceSectionID string, filter string) (result ProductsListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProductsByInvoiceSectionClient.List")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListPreparer(ctx, billingAccountID, invoiceSectionID, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.ProductsByInvoiceSectionClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "billing.ProductsByInvoiceSectionClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.ProductsByInvoiceSectionClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client ProductsByInvoiceSectionClient) ListPreparer(ctx context.Context, billingAccountID string, invoiceSectionID string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"billingAccountId": autorest.Encode("path", billingAccountID),
		"invoiceSectionId": autorest.Encode("path", invoiceSectionID),
	}

	const APIVersion = "2018-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/invoiceSections/{invoiceSectionId}/products", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ProductsByInvoiceSectionClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ProductsByInvoiceSectionClient) ListResponder(resp *http.Response) (result ProductsListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
