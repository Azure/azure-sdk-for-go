package backup

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

// CrrJobsClient is the open API 2.0 Specs for Azure RecoveryServices Backup service
type CrrJobsClient struct {
	BaseClient
}

// NewCrrJobsClient creates an instance of the CrrJobsClient client.
func NewCrrJobsClient(subscriptionID string) CrrJobsClient {
	return NewCrrJobsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewCrrJobsClientWithBaseURI creates an instance of the CrrJobsClient client using a custom endpoint.  Use this when
// interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewCrrJobsClientWithBaseURI(baseURI string, subscriptionID string) CrrJobsClient {
	return CrrJobsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// List sends the list request.
// Parameters:
// azureRegion - azure region to hit Api
func (client CrrJobsClient) List(ctx context.Context, azureRegion string) (result JobResourceListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CrrJobsClient.List")
		defer func() {
			sc := -1
			if result.jrl.Response.Response != nil {
				sc = result.jrl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, azureRegion)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.CrrJobsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.jrl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "backup.CrrJobsClient", "List", resp, "Failure sending request")
		return
	}

	result.jrl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.CrrJobsClient", "List", resp, "Failure responding to request")
	}
	if result.jrl.hasNextLink() && result.jrl.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListPreparer prepares the List request.
func (client CrrJobsClient) ListPreparer(ctx context.Context, azureRegion string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"azureRegion":    autorest.Encode("path", azureRegion),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-12-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/providers/Microsoft.RecoveryServices/locations/{azureRegion}/backupCrrJobs", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client CrrJobsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client CrrJobsClient) ListResponder(resp *http.Response) (result JobResourceList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client CrrJobsClient) listNextResults(ctx context.Context, lastResults JobResourceList) (result JobResourceList, err error) {
	req, err := lastResults.jobResourceListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "backup.CrrJobsClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "backup.CrrJobsClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.CrrJobsClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client CrrJobsClient) ListComplete(ctx context.Context, azureRegion string) (result JobResourceListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CrrJobsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, azureRegion)
	return
}
