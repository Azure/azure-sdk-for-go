package automation

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
	"net/http"
)

// ConnectionTypeClient is the automation Client
type ConnectionTypeClient struct {
	BaseClient
}

// NewConnectionTypeClient creates an instance of the ConnectionTypeClient client.
func NewConnectionTypeClient(subscriptionID string, resourceGroupName string) ConnectionTypeClient {
	return NewConnectionTypeClientWithBaseURI(DefaultBaseURI, subscriptionID, resourceGroupName)
}

// NewConnectionTypeClientWithBaseURI creates an instance of the ConnectionTypeClient client.
func NewConnectionTypeClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string) ConnectionTypeClient {
	return ConnectionTypeClient{NewWithBaseURI(baseURI, subscriptionID, resourceGroupName)}
}

// CreateOrUpdate create a connectiontype.
//
// automationAccountName is the automation account name. connectionTypeName is the parameters supplied to the
// create or update connectiontype operation. parameters is the parameters supplied to the create or update
// connectiontype operation.
func (client ConnectionTypeClient) CreateOrUpdate(ctx context.Context, automationAccountName string, connectionTypeName string, parameters ConnectionTypeCreateOrUpdateParameters) (result ConnectionType, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.ResourceGroupName,
			Constraints: []validation.Constraint{{Target: "client.ResourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "client.ResourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "client.ResourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil}}},
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.Name", Name: validation.Null, Rule: true, Chain: nil},
				{Target: "parameters.ConnectionTypeCreateOrUpdateProperties", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "parameters.ConnectionTypeCreateOrUpdateProperties.FieldDefinitions", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewError("automation.ConnectionTypeClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, automationAccountName, connectionTypeName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client ConnectionTypeClient) CreateOrUpdatePreparer(ctx context.Context, automationAccountName string, connectionTypeName string, parameters ConnectionTypeCreateOrUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"automationAccountName": autorest.Encode("path", automationAccountName),
		"connectionTypeName":    autorest.Encode("path", connectionTypeName),
		"resourceGroupName":     autorest.Encode("path", client.ResourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-31"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/connectionTypes/{connectionTypeName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client ConnectionTypeClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client ConnectionTypeClient) CreateOrUpdateResponder(resp *http.Response) (result ConnectionType, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusConflict),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete delete the connectiontype.
//
// automationAccountName is the automation account name. connectionTypeName is the name of connectiontype.
func (client ConnectionTypeClient) Delete(ctx context.Context, automationAccountName string, connectionTypeName string) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.ResourceGroupName,
			Constraints: []validation.Constraint{{Target: "client.ResourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "client.ResourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "client.ResourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("automation.ConnectionTypeClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, automationAccountName, connectionTypeName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client ConnectionTypeClient) DeletePreparer(ctx context.Context, automationAccountName string, connectionTypeName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"automationAccountName": autorest.Encode("path", automationAccountName),
		"connectionTypeName":    autorest.Encode("path", connectionTypeName),
		"resourceGroupName":     autorest.Encode("path", client.ResourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-31"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/connectionTypes/{connectionTypeName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client ConnectionTypeClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ConnectionTypeClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get retrieve the connectiontype identified by connectiontype name.
//
// automationAccountName is the automation account name. connectionTypeName is the name of connectiontype.
func (client ConnectionTypeClient) Get(ctx context.Context, automationAccountName string, connectionTypeName string) (result ConnectionType, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.ResourceGroupName,
			Constraints: []validation.Constraint{{Target: "client.ResourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "client.ResourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "client.ResourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("automation.ConnectionTypeClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, automationAccountName, connectionTypeName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ConnectionTypeClient) GetPreparer(ctx context.Context, automationAccountName string, connectionTypeName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"automationAccountName": autorest.Encode("path", automationAccountName),
		"connectionTypeName":    autorest.Encode("path", connectionTypeName),
		"resourceGroupName":     autorest.Encode("path", client.ResourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-31"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/connectionTypes/{connectionTypeName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ConnectionTypeClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ConnectionTypeClient) GetResponder(resp *http.Response) (result ConnectionType, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByAutomationAccount retrieve a list of connectiontypes.
//
// automationAccountName is the automation account name.
func (client ConnectionTypeClient) ListByAutomationAccount(ctx context.Context, automationAccountName string) (result ConnectionTypeListResultPage, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.ResourceGroupName,
			Constraints: []validation.Constraint{{Target: "client.ResourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "client.ResourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "client.ResourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("automation.ConnectionTypeClient", "ListByAutomationAccount", err.Error())
	}

	result.fn = client.listByAutomationAccountNextResults
	req, err := client.ListByAutomationAccountPreparer(ctx, automationAccountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "ListByAutomationAccount", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByAutomationAccountSender(req)
	if err != nil {
		result.ctlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "ListByAutomationAccount", resp, "Failure sending request")
		return
	}

	result.ctlr, err = client.ListByAutomationAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "ListByAutomationAccount", resp, "Failure responding to request")
	}

	return
}

// ListByAutomationAccountPreparer prepares the ListByAutomationAccount request.
func (client ConnectionTypeClient) ListByAutomationAccountPreparer(ctx context.Context, automationAccountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"automationAccountName": autorest.Encode("path", automationAccountName),
		"resourceGroupName":     autorest.Encode("path", client.ResourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-31"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/connectionTypes", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByAutomationAccountSender sends the ListByAutomationAccount request. The method will close the
// http.Response Body if it receives an error.
func (client ConnectionTypeClient) ListByAutomationAccountSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListByAutomationAccountResponder handles the response to the ListByAutomationAccount request. The method always
// closes the http.Response Body.
func (client ConnectionTypeClient) ListByAutomationAccountResponder(resp *http.Response) (result ConnectionTypeListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByAutomationAccountNextResults retrieves the next set of results, if any.
func (client ConnectionTypeClient) listByAutomationAccountNextResults(lastResults ConnectionTypeListResult) (result ConnectionTypeListResult, err error) {
	req, err := lastResults.connectionTypeListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "listByAutomationAccountNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByAutomationAccountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "listByAutomationAccountNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByAutomationAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ConnectionTypeClient", "listByAutomationAccountNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByAutomationAccountComplete enumerates all values, automatically crossing page boundaries as required.
func (client ConnectionTypeClient) ListByAutomationAccountComplete(ctx context.Context, automationAccountName string) (result ConnectionTypeListResultIterator, err error) {
	result.page, err = client.ListByAutomationAccount(ctx, automationAccountName)
	return
}
