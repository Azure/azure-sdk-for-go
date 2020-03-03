package msi

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

// SystemAssignedIdentitiesClient is the the Managed Service Identity Client.
type SystemAssignedIdentitiesClient struct {
	BaseClient
}

// NewSystemAssignedIdentitiesClient creates an instance of the SystemAssignedIdentitiesClient client.
func NewSystemAssignedIdentitiesClient(subscriptionID string) SystemAssignedIdentitiesClient {
	return NewSystemAssignedIdentitiesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewSystemAssignedIdentitiesClientWithBaseURI creates an instance of the SystemAssignedIdentitiesClient client using
// a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign
// clouds, Azure stack).
func NewSystemAssignedIdentitiesClientWithBaseURI(baseURI string, subscriptionID string) SystemAssignedIdentitiesClient {
	return SystemAssignedIdentitiesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// GetByScope gets the systemAssignedIdentity available under the specified RP scope.
// Parameters:
// scope - the resource provider scope of the resource. Parent resource being extended by Managed Identities.
func (client SystemAssignedIdentitiesClient) GetByScope(ctx context.Context, scope string) (result SystemAssignedIdentity, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SystemAssignedIdentitiesClient.GetByScope")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetByScopePreparer(ctx, scope)
	if err != nil {
		err = autorest.NewErrorWithError(err, "msi.SystemAssignedIdentitiesClient", "GetByScope", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetByScopeSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "msi.SystemAssignedIdentitiesClient", "GetByScope", resp, "Failure sending request")
		return
	}

	result, err = client.GetByScopeResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "msi.SystemAssignedIdentitiesClient", "GetByScope", resp, "Failure responding to request")
	}

	return
}

// GetByScopePreparer prepares the GetByScope request.
func (client SystemAssignedIdentitiesClient) GetByScopePreparer(ctx context.Context, scope string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"scope": scope,
	}

	const APIVersion = "2018-11-30"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{scope}/providers/Microsoft.ManagedIdentity/identities/default", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetByScopeSender sends the GetByScope request. The method will close the
// http.Response Body if it receives an error.
func (client SystemAssignedIdentitiesClient) GetByScopeSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetByScopeResponder handles the response to the GetByScope request. The method always
// closes the http.Response Body.
func (client SystemAssignedIdentitiesClient) GetByScopeResponder(resp *http.Response) (result SystemAssignedIdentity, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
