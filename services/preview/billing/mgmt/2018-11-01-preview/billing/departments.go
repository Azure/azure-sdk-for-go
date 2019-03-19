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
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// DepartmentsClient is the billing client provides access to billing resources for Azure subscriptions.
type DepartmentsClient struct {
	BaseClient
}

// NewDepartmentsClient creates an instance of the DepartmentsClient client.
func NewDepartmentsClient(subscriptionID string) DepartmentsClient {
	return NewDepartmentsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDepartmentsClientWithBaseURI creates an instance of the DepartmentsClient client.
func NewDepartmentsClientWithBaseURI(baseURI string, subscriptionID string) DepartmentsClient {
	return DepartmentsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get get the department by id.
// Parameters:
// billingAccountID - billing Account Id.
// departmentID - department Id.
// expand - may be used to expand the enrollmentAccounts.
// filter - the filter supports 'eq', 'lt', 'gt', 'le', 'ge', and 'and'. It does not currently support 'ne',
// 'or', or 'not'. Tag filter is a key value pair string where key and value is separated by a colon (:).
func (client DepartmentsClient) Get(ctx context.Context, billingAccountID string, departmentID string, expand string, filter string) (result Department, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/DepartmentsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, billingAccountID, departmentID, expand, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.DepartmentsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "billing.DepartmentsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "billing.DepartmentsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client DepartmentsClient) GetPreparer(ctx context.Context, billingAccountID string, departmentID string, expand string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"billingAccountId": autorest.Encode("path", billingAccountID),
		"departmentId":     autorest.Encode("path", departmentID),
	}

	const APIVersion = "2018-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(expand) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/departments/{departmentId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client DepartmentsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DepartmentsClient) GetResponder(resp *http.Response) (result Department, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
