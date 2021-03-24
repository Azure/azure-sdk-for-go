package dataprotection

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

// FindRestorableTimeRangesClient is the open API 2.0 Specs for Azure Data Protection service
type FindRestorableTimeRangesClient struct {
	BaseClient
}

// NewFindRestorableTimeRangesClient creates an instance of the FindRestorableTimeRangesClient client.
func NewFindRestorableTimeRangesClient(subscriptionID string) FindRestorableTimeRangesClient {
	return NewFindRestorableTimeRangesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewFindRestorableTimeRangesClientWithBaseURI creates an instance of the FindRestorableTimeRangesClient client using
// a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign
// clouds, Azure stack).
func NewFindRestorableTimeRangesClientWithBaseURI(baseURI string, subscriptionID string) FindRestorableTimeRangesClient {
	return FindRestorableTimeRangesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Post sends the post request.
// Parameters:
// vaultName - the name of the backup vault.
// resourceGroupName - the name of the resource group where the backup vault is present.
// parameters - request body for operation
func (client FindRestorableTimeRangesClient) Post(ctx context.Context, vaultName string, resourceGroupName string, backupInstances string, parameters AzureBackupFindRestorableTimeRangesRequest) (result AzureBackupFindRestorableTimeRangesResponseResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FindRestorableTimeRangesClient.Post")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.StartTime", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "parameters.EndTime", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("dataprotection.FindRestorableTimeRangesClient", "Post", err.Error())
	}

	req, err := client.PostPreparer(ctx, vaultName, resourceGroupName, backupInstances, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataprotection.FindRestorableTimeRangesClient", "Post", nil, "Failure preparing request")
		return
	}

	resp, err := client.PostSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "dataprotection.FindRestorableTimeRangesClient", "Post", resp, "Failure sending request")
		return
	}

	result, err = client.PostResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataprotection.FindRestorableTimeRangesClient", "Post", resp, "Failure responding to request")
		return
	}

	return
}

// PostPreparer prepares the Post request.
func (client FindRestorableTimeRangesClient) PostPreparer(ctx context.Context, vaultName string, resourceGroupName string, backupInstances string, parameters AzureBackupFindRestorableTimeRangesRequest) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"backupInstances":   autorest.Encode("path", backupInstances),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2021-01-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{vaultName}/backupInstances/{backupInstances}/findRestorableTimeRanges", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PostSender sends the Post request. The method will close the
// http.Response Body if it receives an error.
func (client FindRestorableTimeRangesClient) PostSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// PostResponder handles the response to the Post request. The method always
// closes the http.Response Body.
func (client FindRestorableTimeRangesClient) PostResponder(resp *http.Response) (result AzureBackupFindRestorableTimeRangesResponseResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
