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
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// AddressClient is the billing client provides access to billing resources for Azure subscriptions.
type AddressClient struct {
	BaseClient
}

// NewAddressClient creates an instance of the AddressClient client.
func NewAddressClient(subscriptionID string) AddressClient {
	return NewAddressClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewAddressClientWithBaseURI creates an instance of the AddressClient client using a custom endpoint.  Use this when
// interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewAddressClientWithBaseURI(baseURI string, subscriptionID string) AddressClient {
	return AddressClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Validate validates an address. Use the operation to validate an address before using it as soldTo or a billTo
// address.
func (client AddressClient) Validate(ctx context.Context, address AddressDetails) (result ValidateAddressResponse, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/AddressClient.Validate")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: address,
			Constraints: []validation.Constraint{{Target: "address.AddressLine1", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "address.Country", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("billing.AddressClient", "Validate", err.Error())
	}

	req, err := client.ValidatePreparer(ctx, address)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.AddressClient", "Validate", nil, "Failure preparing request")
		return
	}

	resp, err := client.ValidateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "billing.AddressClient", "Validate", resp, "Failure sending request")
		return
	}

	result, err = client.ValidateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.AddressClient", "Validate", resp, "Failure responding to request")
	}

	return
}

// ValidatePreparer prepares the Validate request.
func (client AddressClient) ValidatePreparer(ctx context.Context, address AddressDetails) (*http.Request, error) {
	const APIVersion = "2020-05-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/providers/Microsoft.Billing/validateAddress"),
		autorest.WithJSON(address),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ValidateSender sends the Validate request. The method will close the
// http.Response Body if it receives an error.
func (client AddressClient) ValidateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ValidateResponder handles the response to the Validate request. The method always
// closes the http.Response Body.
func (client AddressClient) ValidateResponder(resp *http.Response) (result ValidateAddressResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
