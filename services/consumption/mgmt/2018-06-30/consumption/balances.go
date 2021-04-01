package consumption

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
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

// BalancesClient is the consumption management client provides access to consumption resources for Azure Enterprise
// Subscriptions.
type BalancesClient struct {
	BaseClient
}

// NewBalancesClient creates an instance of the BalancesClient client.
func NewBalancesClient(subscriptionID string) BalancesClient {
	return NewBalancesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewBalancesClientWithBaseURI creates an instance of the BalancesClient client using a custom endpoint.  Use this
// when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewBalancesClientWithBaseURI(baseURI string, subscriptionID string) BalancesClient {
	return BalancesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// GetByBillingAccount gets the balances for a scope by billingAccountId. Balances are available via this API only for
// May 1, 2014 or later.
// Parameters:
// billingAccountID - billingAccount ID
func (client BalancesClient) GetByBillingAccount(ctx context.Context, billingAccountID string) (result Balance, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BalancesClient.GetByBillingAccount")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetByBillingAccountPreparer(ctx, billingAccountID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.BalancesClient", "GetByBillingAccount", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetByBillingAccountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "consumption.BalancesClient", "GetByBillingAccount", resp, "Failure sending request")
		return
	}

	result, err = client.GetByBillingAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.BalancesClient", "GetByBillingAccount", resp, "Failure responding to request")
		return
	}

	return
}

// GetByBillingAccountPreparer prepares the GetByBillingAccount request.
func (client BalancesClient) GetByBillingAccountPreparer(ctx context.Context, billingAccountID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"billingAccountId": autorest.Encode("path", billingAccountID),
	}

	const APIVersion = "2018-06-30"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/providers/Microsoft.Consumption/balances", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetByBillingAccountSender sends the GetByBillingAccount request. The method will close the
// http.Response Body if it receives an error.
func (client BalancesClient) GetByBillingAccountSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetByBillingAccountResponder handles the response to the GetByBillingAccount request. The method always
// closes the http.Response Body.
func (client BalancesClient) GetByBillingAccountResponder(resp *http.Response) (result Balance, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetForBillingPeriodByBillingAccount gets the balances for a scope by billing period and billingAccountId. Balances
// are available via this API only for May 1, 2014 or later.
// Parameters:
// billingAccountID - billingAccount ID
// billingPeriodName - billing Period Name.
func (client BalancesClient) GetForBillingPeriodByBillingAccount(ctx context.Context, billingAccountID string, billingPeriodName string) (result Balance, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BalancesClient.GetForBillingPeriodByBillingAccount")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetForBillingPeriodByBillingAccountPreparer(ctx, billingAccountID, billingPeriodName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.BalancesClient", "GetForBillingPeriodByBillingAccount", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetForBillingPeriodByBillingAccountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "consumption.BalancesClient", "GetForBillingPeriodByBillingAccount", resp, "Failure sending request")
		return
	}

	result, err = client.GetForBillingPeriodByBillingAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumption.BalancesClient", "GetForBillingPeriodByBillingAccount", resp, "Failure responding to request")
		return
	}

	return
}

// GetForBillingPeriodByBillingAccountPreparer prepares the GetForBillingPeriodByBillingAccount request.
func (client BalancesClient) GetForBillingPeriodByBillingAccountPreparer(ctx context.Context, billingAccountID string, billingPeriodName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"billingAccountId":  autorest.Encode("path", billingAccountID),
		"billingPeriodName": autorest.Encode("path", billingPeriodName),
	}

	const APIVersion = "2018-06-30"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/providers/Microsoft.Billing/billingPeriods/{billingPeriodName}/providers/Microsoft.Consumption/balances", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetForBillingPeriodByBillingAccountSender sends the GetForBillingPeriodByBillingAccount request. The method will close the
// http.Response Body if it receives an error.
func (client BalancesClient) GetForBillingPeriodByBillingAccountSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetForBillingPeriodByBillingAccountResponder handles the response to the GetForBillingPeriodByBillingAccount request. The method always
// closes the http.Response Body.
func (client BalancesClient) GetForBillingPeriodByBillingAccountResponder(resp *http.Response) (result Balance, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
