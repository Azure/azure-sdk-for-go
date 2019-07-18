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
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// UpdateAccessClient is the ARM applications
type UpdateAccessClient struct {
	BaseClient
}

// NewUpdateAccessClient creates an instance of the UpdateAccessClient client.
func NewUpdateAccessClient(subscriptionID string) UpdateAccessClient {
	return NewUpdateAccessClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewUpdateAccessClientWithBaseURI creates an instance of the UpdateAccessClient client.
func NewUpdateAccessClientWithBaseURI(baseURI string, subscriptionID string) UpdateAccessClient {
	return UpdateAccessClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Post update Access on application.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// applicationName - the name of the managed application.
// parameters - parameters supplied to the update managed application access.
func (client UpdateAccessClient) Post(ctx context.Context, resourceGroupName string, applicationName string, parameters JitUpdateAccessDefinition) (result UpdateAccessPostFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/UpdateAccessClient.Post")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\p{L}\._\(\)\w]+$`, Chain: nil}}},
		{TargetValue: applicationName,
			Constraints: []validation.Constraint{{Target: "applicationName", Name: validation.MaxLength, Rule: 64, Chain: nil},
				{Target: "applicationName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewError("managedapplications.UpdateAccessClient", "Post", err.Error())
	}

	req, err := client.PostPreparer(ctx, resourceGroupName, applicationName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapplications.UpdateAccessClient", "Post", nil, "Failure preparing request")
		return
	}

	result, err = client.PostSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapplications.UpdateAccessClient", "Post", result.Response(), "Failure sending request")
		return
	}

	return
}

// PostPreparer prepares the Post request.
func (client UpdateAccessClient) PostPreparer(ctx context.Context, resourceGroupName string, applicationName string, parameters JitUpdateAccessDefinition) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationName":   autorest.Encode("path", applicationName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-03-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Solutions/applications/{applicationName}/updateAccess", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PostSender sends the Post request. The method will close the
// http.Response Body if it receives an error.
func (client UpdateAccessClient) PostSender(req *http.Request) (future UpdateAccessPostFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// PostResponder handles the response to the Post request. The method always
// closes the http.Response Body.
func (client UpdateAccessClient) PostResponder(resp *http.Response) (result JitUpdateAccessDefinition, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
