package insights

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

// DiagnosticSettingsCategoryGroupClient is the monitor Management Client
type DiagnosticSettingsCategoryGroupClient struct {
	BaseClient
}

// NewDiagnosticSettingsCategoryGroupClient creates an instance of the DiagnosticSettingsCategoryGroupClient client.
func NewDiagnosticSettingsCategoryGroupClient(subscriptionID string) DiagnosticSettingsCategoryGroupClient {
	return NewDiagnosticSettingsCategoryGroupClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDiagnosticSettingsCategoryGroupClientWithBaseURI creates an instance of the DiagnosticSettingsCategoryGroupClient
// client.
func NewDiagnosticSettingsCategoryGroupClientWithBaseURI(baseURI string, subscriptionID string) DiagnosticSettingsCategoryGroupClient {
	return DiagnosticSettingsCategoryGroupClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets the diagnostic settings category for the specified resource.
//
// resourceURI is the identifier of the resource. name is the name of the diagnostic setting.
func (client DiagnosticSettingsCategoryGroupClient) Get(ctx context.Context, resourceURI string, name string) (result DiagnosticSettingsCategoryResource, err error) {
	req, err := client.GetPreparer(ctx, resourceURI, name)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.DiagnosticSettingsCategoryGroupClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.DiagnosticSettingsCategoryGroupClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.DiagnosticSettingsCategoryGroupClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client DiagnosticSettingsCategoryGroupClient) GetPreparer(ctx context.Context, resourceURI string, name string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"name":        autorest.Encode("path", name),
		"resourceUri": resourceURI,
	}

	const APIVersion = "2017-05-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{resourceUri}/providers/microsoft.insights/diagnosticSettingsCategories/{name}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client DiagnosticSettingsCategoryGroupClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DiagnosticSettingsCategoryGroupClient) GetResponder(resp *http.Response) (result DiagnosticSettingsCategoryResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List lists the diagnostic settings categories for the specified resource.
//
// resourceURI is the identifier of the resource.
func (client DiagnosticSettingsCategoryGroupClient) List(ctx context.Context, resourceURI string) (result DiagnosticSettingsCategoryResourceCollection, err error) {
	req, err := client.ListPreparer(ctx, resourceURI)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.DiagnosticSettingsCategoryGroupClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.DiagnosticSettingsCategoryGroupClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.DiagnosticSettingsCategoryGroupClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client DiagnosticSettingsCategoryGroupClient) ListPreparer(ctx context.Context, resourceURI string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceUri": resourceURI,
	}

	const APIVersion = "2017-05-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{resourceUri}/providers/microsoft.insights/diagnosticSettingsCategories", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client DiagnosticSettingsCategoryGroupClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client DiagnosticSettingsCategoryGroupClient) ListResponder(resp *http.Response) (result DiagnosticSettingsCategoryResourceCollection, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
