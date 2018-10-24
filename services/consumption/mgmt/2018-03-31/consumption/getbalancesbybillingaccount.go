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
    "github.com/Azure/go-autorest/autorest"
    "github.com/Azure/go-autorest/autorest/azure"
    "net/http"
    "context"
    "github.com/Azure/go-autorest/tracing"
)

// GetBalancesByBillingAccountClient is the consumption management client provides access to consumption resources for
// Azure Enterprise Subscriptions.
type GetBalancesByBillingAccountClient struct {
    BaseClient
}
// NewGetBalancesByBillingAccountClient creates an instance of the GetBalancesByBillingAccountClient client.
func NewGetBalancesByBillingAccountClient(subscriptionID string) GetBalancesByBillingAccountClient {
    return NewGetBalancesByBillingAccountClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewGetBalancesByBillingAccountClientWithBaseURI creates an instance of the GetBalancesByBillingAccountClient client.
    func NewGetBalancesByBillingAccountClientWithBaseURI(baseURI string, subscriptionID string) GetBalancesByBillingAccountClient {
        return GetBalancesByBillingAccountClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// ByBillingPeriod gets the balances for a scope by billing period and billingAccountId. Balances are available via
// this API only for May 1, 2014 or later.
    // Parameters:
        // billingAccountID - billingAccount ID
        // billingPeriodName - billing Period Name.
func (client GetBalancesByBillingAccountClient) ByBillingPeriod(ctx context.Context, billingAccountID string, billingPeriodName string) (result Balance, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/GetBalancesByBillingAccountClient.ByBillingPeriod")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.ByBillingPeriodPreparer(ctx, billingAccountID, billingPeriodName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "consumption.GetBalancesByBillingAccountClient", "ByBillingPeriod", nil , "Failure preparing request")
    return
    }

            resp, err := client.ByBillingPeriodSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "consumption.GetBalancesByBillingAccountClient", "ByBillingPeriod", resp, "Failure sending request")
            return
            }

            result, err = client.ByBillingPeriodResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "consumption.GetBalancesByBillingAccountClient", "ByBillingPeriod", resp, "Failure responding to request")
            }

    return
    }

    // ByBillingPeriodPreparer prepares the ByBillingPeriod request.
    func (client GetBalancesByBillingAccountClient) ByBillingPeriodPreparer(ctx context.Context, billingAccountID string, billingPeriodName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "billingAccountId": autorest.Encode("path",billingAccountID),
            "billingPeriodName": autorest.Encode("path",billingPeriodName),
            }

                        const APIVersion = "2018-03-31"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/providers/Microsoft.Billing/billingPeriods/{billingPeriodName}/providers/Microsoft.Consumption/balances",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ByBillingPeriodSender sends the ByBillingPeriod request. The method will close the
    // http.Response Body if it receives an error.
    func (client GetBalancesByBillingAccountClient) ByBillingPeriodSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            }

// ByBillingPeriodResponder handles the response to the ByBillingPeriod request. The method always
// closes the http.Response Body.
func (client GetBalancesByBillingAccountClient) ByBillingPeriodResponder(resp *http.Response) (result Balance, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

