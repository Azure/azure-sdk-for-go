// Package resources implements the Azure ARM Resources service API version 2018-05-01.
//
// Provides operations for working with resources and resource groups.
package resources

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
	"net/http"
)

const (
	// DefaultBaseURI is the default URI used for the service Resources
	DefaultBaseURI = "https://management.azure.com"
)

// BaseClient is the base client for Resources.
type BaseClient struct {
	autorest.Client
	BaseURI        string
	SubscriptionID string
}

// New creates an instance of the BaseClient client.
func New(subscriptionID string) BaseClient {
	return NewWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewWithBaseURI creates an instance of the BaseClient client.
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return BaseClient{
		Client:         autorest.NewClientWithUserAgent(UserAgent()),
		BaseURI:        baseURI,
		SubscriptionID: subscriptionID,
	}
}

// ListOperations lists all of the available Microsoft.Resources REST API operations.
func (client BaseClient) ListOperations(ctx context.Context) (result OperationListResultPage, err error) {
	result.fn = client.listOperationsNextResults
	req, err := client.ListOperationsPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.BaseClient", "ListOperations", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListOperationsSender(req)
	if err != nil {
		result.olr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.BaseClient", "ListOperations", resp, "Failure sending request")
		return
	}

	result.olr, err = client.ListOperationsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.BaseClient", "ListOperations", resp, "Failure responding to request")
	}

	return
}

// ListOperationsPreparer prepares the ListOperations request.
func (client BaseClient) ListOperationsPreparer(ctx context.Context) (*http.Request, error) {
	const APIVersion = "2018-05-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/providers/Microsoft.Resources/operations"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListOperationsSender sends the ListOperations request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) ListOperationsSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListOperationsResponder handles the response to the ListOperations request. The method always
// closes the http.Response Body.
func (client BaseClient) ListOperationsResponder(resp *http.Response) (result OperationListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listOperationsNextResults retrieves the next set of results, if any.
func (client BaseClient) listOperationsNextResults(lastResults OperationListResult) (result OperationListResult, err error) {
	req, err := lastResults.operationListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "resources.BaseClient", "listOperationsNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListOperationsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "resources.BaseClient", "listOperationsNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListOperationsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.BaseClient", "listOperationsNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListOperationsComplete enumerates all values, automatically crossing page boundaries as required.
func (client BaseClient) ListOperationsComplete(ctx context.Context) (result OperationListResultIterator, err error) {
	result.page, err = client.ListOperations(ctx)
	return
}
