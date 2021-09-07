package apimanagement

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

// QuotaByPeriodKeysClient is the client for the QuotaByPeriodKeys methods of the Apimanagement service.
type QuotaByPeriodKeysClient struct {
	BaseClient
}

// NewQuotaByPeriodKeysClient creates an instance of the QuotaByPeriodKeysClient client.
func NewQuotaByPeriodKeysClient() QuotaByPeriodKeysClient {
	return QuotaByPeriodKeysClient{New()}
}

// Get gets the value of the quota counter associated with the counter-key in the policy for the specific period in
// service instance.
// Parameters:
// apimBaseURL - the management endpoint of the API Management service, for example
// https://myapimservice.management.azure-api.net.
// quotaCounterKey - quota counter key identifier.This is the result of expression defined in counter-key
// attribute of the quota-by-key policy.For Example, if you specify counter-key="boo" in the policy, then it’s
// accessible by "boo" counter key. But if it’s defined as counter-key="@("b"+"a")" then it will be accessible
// by "ba" key
// quotaPeriodKey - quota period key identifier.
func (client QuotaByPeriodKeysClient) Get(ctx context.Context, apimBaseURL string, quotaCounterKey string, quotaPeriodKey string) (result QuotaCounterContract, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/QuotaByPeriodKeysClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, apimBaseURL, quotaCounterKey, quotaPeriodKey)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client QuotaByPeriodKeysClient) GetPreparer(ctx context.Context, apimBaseURL string, quotaCounterKey string, quotaPeriodKey string) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":            autorest.Encode("path"),
		"apimBaseUrl": apimBaseURL,
	}

	pathParameters := map[string]interface{}{
		"quotaCounterKey": autorest.Encode("path", quotaCounterKey),
		"quotaPeriodKey":  autorest.Encode("path", quotaPeriodKey),
	}

	const APIVersion = "2017-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{apimBaseUrl}", urlParameters),
		autorest.WithPathParameters("/quotas/{quotaCounterKey}/{quotaPeriodKey}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client QuotaByPeriodKeysClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client QuotaByPeriodKeysClient) GetResponder(resp *http.Response) (result QuotaCounterContract, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Update updates an existing quota counter value in the specified service instance.
// Parameters:
// apimBaseURL - the management endpoint of the API Management service, for example
// https://myapimservice.management.azure-api.net.
// quotaCounterKey - quota counter key identifier.This is the result of expression defined in counter-key
// attribute of the quota-by-key policy.For Example, if you specify counter-key="boo" in the policy, then it’s
// accessible by "boo" counter key. But if it’s defined as counter-key="@("b"+"a")" then it will be accessible
// by "ba" key
// quotaPeriodKey - quota period key identifier.
// parameters - the value of the Quota counter to be applied on the specified period.
func (client QuotaByPeriodKeysClient) Update(ctx context.Context, apimBaseURL string, quotaCounterKey string, quotaPeriodKey string, parameters QuotaCounterValueContractProperties) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/QuotaByPeriodKeysClient.Update")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdatePreparer(ctx, apimBaseURL, quotaCounterKey, quotaPeriodKey, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.QuotaByPeriodKeysClient", "Update", resp, "Failure responding to request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client QuotaByPeriodKeysClient) UpdatePreparer(ctx context.Context, apimBaseURL string, quotaCounterKey string, quotaPeriodKey string, parameters QuotaCounterValueContractProperties) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"":            autorest.Encode("path"),
		"apimBaseUrl": apimBaseURL,
	}

	pathParameters := map[string]interface{}{
		"quotaCounterKey": autorest.Encode("path", quotaCounterKey),
		"quotaPeriodKey":  autorest.Encode("path", quotaPeriodKey),
	}

	const APIVersion = "2017-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithCustomBaseURL("{apimBaseUrl}", urlParameters),
		autorest.WithPathParameters("/quotas/{quotaCounterKey}/{quotaPeriodKey}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client QuotaByPeriodKeysClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client QuotaByPeriodKeysClient) UpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}
