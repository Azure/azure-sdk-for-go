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

// PolicySnippetsClient is the client for the PolicySnippets methods of the Apimanagement service.
type PolicySnippetsClient struct {
	BaseClient
}

// NewPolicySnippetsClient creates an instance of the PolicySnippetsClient client.
func NewPolicySnippetsClient() PolicySnippetsClient {
	return PolicySnippetsClient{New()}
}

// List lists all policy snippets.
// Parameters:
// apimBaseURL - the management endpoint of the API Management service, for example
// https://myapimservice.management.azure-api.net.
// scope - policy scope.
func (client PolicySnippetsClient) List(ctx context.Context, apimBaseURL string, scope PolicyScopeContract) (result PolicySnippetsCollection, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PolicySnippetsClient.List")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListPreparer(ctx, apimBaseURL, scope)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.PolicySnippetsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "apimanagement.PolicySnippetsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.PolicySnippetsClient", "List", resp, "Failure responding to request")
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client PolicySnippetsClient) ListPreparer(ctx context.Context, apimBaseURL string, scope PolicyScopeContract) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"apimBaseUrl": apimBaseURL,
	}

	const APIVersion = "2017-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(string(scope)) > 0 {
		queryParameters["scope"] = autorest.Encode("query", scope)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{apimBaseUrl}", urlParameters),
		autorest.WithPath("/policySnippets"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client PolicySnippetsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client PolicySnippetsClient) ListResponder(resp *http.Response) (result PolicySnippetsCollection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
