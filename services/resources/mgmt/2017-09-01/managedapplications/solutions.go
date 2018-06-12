package managedapplications

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

// SolutionsClient is the ARM applications
type SolutionsClient struct {
	BaseClient
}

// NewSolutionsClient creates an instance of the SolutionsClient client.
func NewSolutionsClient(subscriptionID string) SolutionsClient {
	return NewSolutionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewSolutionsClientWithBaseURI creates an instance of the SolutionsClient client.
func NewSolutionsClientWithBaseURI(baseURI string, subscriptionID string) SolutionsClient {
	return SolutionsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Operations gets all the preview solution operations.
func (client SolutionsClient) Operations(ctx context.Context) (result OperationsList, err error) {
	req, err := client.OperationsPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapplications.SolutionsClient", "Operations", nil, "Failure preparing request")
		return
	}

	resp, err := client.OperationsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "managedapplications.SolutionsClient", "Operations", resp, "Failure sending request")
		return
	}

	result, err = client.OperationsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapplications.SolutionsClient", "Operations", resp, "Failure responding to request")
	}

	return
}

// OperationsPreparer prepares the Operations request.
func (client SolutionsClient) OperationsPreparer(ctx context.Context) (*http.Request, error) {
	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/providers/Microsoft.Solutions/operations"))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// OperationsSender sends the Operations request. The method will close the
// http.Response Body if it receives an error.
func (client SolutionsClient) OperationsSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// OperationsResponder handles the response to the Operations request. The method always
// closes the http.Response Body.
func (client SolutionsClient) OperationsResponder(resp *http.Response) (result OperationsList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
