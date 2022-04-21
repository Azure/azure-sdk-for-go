package synapse

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
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

// SQLPoolConnectionPoliciesClient is the azure Synapse Analytics Management Client
type SQLPoolConnectionPoliciesClient struct {
	BaseClient
}

// NewSQLPoolConnectionPoliciesClient creates an instance of the SQLPoolConnectionPoliciesClient client.
func NewSQLPoolConnectionPoliciesClient(subscriptionID string) SQLPoolConnectionPoliciesClient {
	return NewSQLPoolConnectionPoliciesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewSQLPoolConnectionPoliciesClientWithBaseURI creates an instance of the SQLPoolConnectionPoliciesClient client
// using a custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign
// clouds, Azure stack).
func NewSQLPoolConnectionPoliciesClientWithBaseURI(baseURI string, subscriptionID string) SQLPoolConnectionPoliciesClient {
	return SQLPoolConnectionPoliciesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get get a Sql pool's connection policy, which is used with table auditing.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - the name of the workspace
// SQLPoolName - SQL pool name
func (client SQLPoolConnectionPoliciesClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, SQLPoolName string) (result SQLPoolConnectionPolicy, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SQLPoolConnectionPoliciesClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("synapse.SQLPoolConnectionPoliciesClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, workspaceName, SQLPoolName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolConnectionPoliciesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolConnectionPoliciesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolConnectionPoliciesClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client SQLPoolConnectionPoliciesClient) GetPreparer(ctx context.Context, resourceGroupName string, workspaceName string, SQLPoolName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"connectionPolicyName": autorest.Encode("path", "default"),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"sqlPoolName":          autorest.Encode("path", SQLPoolName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
		"workspaceName":        autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2020-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/connectionPolicies/{connectionPolicyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client SQLPoolConnectionPoliciesClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client SQLPoolConnectionPoliciesClient) GetResponder(resp *http.Response) (result SQLPoolConnectionPolicy, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
