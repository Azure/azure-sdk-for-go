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

// CrossRegionRestoreClient is the open API 2.0 Specs for Azure RecoveryServices Backup service
type CrossRegionRestoreClient struct {
	BaseClient
}

// NewCrossRegionRestoreClient creates an instance of the CrossRegionRestoreClient client.
func NewCrossRegionRestoreClient(subscriptionID string) CrossRegionRestoreClient {
	return NewCrossRegionRestoreClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewCrossRegionRestoreClientWithBaseURI creates an instance of the CrossRegionRestoreClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewCrossRegionRestoreClientWithBaseURI(baseURI string, subscriptionID string) CrossRegionRestoreClient {
	return CrossRegionRestoreClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Trigger sends the trigger request.
// Parameters:
// azureRegion - azure region to hit Api
// parameters - resource cross region restore request
func (client CrossRegionRestoreClient) Trigger(ctx context.Context, azureRegion string, parameters CrossRegionRestoreRequestResource) (result CrossRegionRestoreTriggerFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CrossRegionRestoreClient.Trigger")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.TriggerPreparer(ctx, azureRegion, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.CrossRegionRestoreClient", "Trigger", nil, "Failure preparing request")
		return
	}

	result, err = client.TriggerSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.CrossRegionRestoreClient", "Trigger", result.Response(), "Failure sending request")
		return
	}

	return
}

// TriggerPreparer prepares the Trigger request.
func (client CrossRegionRestoreClient) TriggerPreparer(ctx context.Context, azureRegion string, parameters CrossRegionRestoreRequestResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"azureRegion":    autorest.Encode("path", azureRegion),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-12-20"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/providers/Microsoft.RecoveryServices/locations/{azureRegion}/backupCrossRegionRestore", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// TriggerSender sends the Trigger request. The method will close the
// http.Response Body if it receives an error.
func (client CrossRegionRestoreClient) TriggerSender(req *http.Request) (future CrossRegionRestoreTriggerFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// TriggerResponder handles the response to the Trigger request. The method always
// closes the http.Response Body.
func (client CrossRegionRestoreClient) TriggerResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
