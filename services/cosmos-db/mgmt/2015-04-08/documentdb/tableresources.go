package documentdb

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

// TableResourcesClient is the azure Cosmos DB Database Service Resource Provider REST API
type TableResourcesClient struct {
	BaseClient
}

// NewTableResourcesClient creates an instance of the TableResourcesClient client.
func NewTableResourcesClient(subscriptionID string) TableResourcesClient {
	return NewTableResourcesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewTableResourcesClientWithBaseURI creates an instance of the TableResourcesClient client.
func NewTableResourcesClientWithBaseURI(baseURI string, subscriptionID string) TableResourcesClient {
	return TableResourcesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateUpdateTable create or update an Azure Cosmos DB Table
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
// tableName - cosmos DB table name.
// createUpdateTableParameters - the parameters to provide for the current Table.
func (client TableResourcesClient) CreateUpdateTable(ctx context.Context, resourceGroupName string, accountName string, tableName string, createUpdateTableParameters TableCreateUpdateParameters) (result TableResourcesCreateUpdateTableFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TableResourcesClient.CreateUpdateTable")
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
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "accountName", Name: validation.Pattern, Rule: `^[a-z0-9]+(-[a-z0-9]+)*`, Chain: nil}}},
		{TargetValue: createUpdateTableParameters,
			Constraints: []validation.Constraint{{Target: "createUpdateTableParameters.TableCreateUpdateProperties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "createUpdateTableParameters.TableCreateUpdateProperties.Resource", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "createUpdateTableParameters.TableCreateUpdateProperties.Resource.ID", Name: validation.Null, Rule: true, Chain: nil}}},
					{Target: "createUpdateTableParameters.TableCreateUpdateProperties.Options", Name: validation.Null, Rule: true, Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("documentdb.TableResourcesClient", "CreateUpdateTable", err.Error())
	}

	req, err := client.CreateUpdateTablePreparer(ctx, resourceGroupName, accountName, tableName, createUpdateTableParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "CreateUpdateTable", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateUpdateTableSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "CreateUpdateTable", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateUpdateTablePreparer prepares the CreateUpdateTable request.
func (client TableResourcesClient) CreateUpdateTablePreparer(ctx context.Context, resourceGroupName string, accountName string, tableName string, createUpdateTableParameters TableCreateUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tableName":         autorest.Encode("path", tableName),
	}

	const APIVersion = "2015-04-08"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/tables/{tableName}", pathParameters),
		autorest.WithJSON(createUpdateTableParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateUpdateTableSender sends the CreateUpdateTable request. The method will close the
// http.Response Body if it receives an error.
func (client TableResourcesClient) CreateUpdateTableSender(req *http.Request) (future TableResourcesCreateUpdateTableFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// CreateUpdateTableResponder handles the response to the CreateUpdateTable request. The method always
// closes the http.Response Body.
func (client TableResourcesClient) CreateUpdateTableResponder(resp *http.Response) (result Table, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// DeleteTable deletes an existing Azure Cosmos DB Table.
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
// tableName - cosmos DB table name.
func (client TableResourcesClient) DeleteTable(ctx context.Context, resourceGroupName string, accountName string, tableName string) (result TableResourcesDeleteTableFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TableResourcesClient.DeleteTable")
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
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "accountName", Name: validation.Pattern, Rule: `^[a-z0-9]+(-[a-z0-9]+)*`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("documentdb.TableResourcesClient", "DeleteTable", err.Error())
	}

	req, err := client.DeleteTablePreparer(ctx, resourceGroupName, accountName, tableName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "DeleteTable", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteTableSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "DeleteTable", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeleteTablePreparer prepares the DeleteTable request.
func (client TableResourcesClient) DeleteTablePreparer(ctx context.Context, resourceGroupName string, accountName string, tableName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tableName":         autorest.Encode("path", tableName),
	}

	const APIVersion = "2015-04-08"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/tables/{tableName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteTableSender sends the DeleteTable request. The method will close the
// http.Response Body if it receives an error.
func (client TableResourcesClient) DeleteTableSender(req *http.Request) (future TableResourcesDeleteTableFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// DeleteTableResponder handles the response to the DeleteTable request. The method always
// closes the http.Response Body.
func (client TableResourcesClient) DeleteTableResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// GetTable gets the Tables under an existing Azure Cosmos DB database account with the provided name.
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
// tableName - cosmos DB table name.
func (client TableResourcesClient) GetTable(ctx context.Context, resourceGroupName string, accountName string, tableName string) (result Table, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TableResourcesClient.GetTable")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "accountName", Name: validation.Pattern, Rule: `^[a-z0-9]+(-[a-z0-9]+)*`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("documentdb.TableResourcesClient", "GetTable", err.Error())
	}

	req, err := client.GetTablePreparer(ctx, resourceGroupName, accountName, tableName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "GetTable", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetTableSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "GetTable", resp, "Failure sending request")
		return
	}

	result, err = client.GetTableResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "GetTable", resp, "Failure responding to request")
	}

	return
}

// GetTablePreparer prepares the GetTable request.
func (client TableResourcesClient) GetTablePreparer(ctx context.Context, resourceGroupName string, accountName string, tableName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tableName":         autorest.Encode("path", tableName),
	}

	const APIVersion = "2015-04-08"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/tables/{tableName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetTableSender sends the GetTable request. The method will close the
// http.Response Body if it receives an error.
func (client TableResourcesClient) GetTableSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// GetTableResponder handles the response to the GetTable request. The method always
// closes the http.Response Body.
func (client TableResourcesClient) GetTableResponder(resp *http.Response) (result Table, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetTableThroughput gets the RUs per second of the Table under an existing Azure Cosmos DB database account with the
// provided name.
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
// tableName - cosmos DB table name.
func (client TableResourcesClient) GetTableThroughput(ctx context.Context, resourceGroupName string, accountName string, tableName string) (result Throughput, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TableResourcesClient.GetTableThroughput")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "accountName", Name: validation.Pattern, Rule: `^[a-z0-9]+(-[a-z0-9]+)*`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("documentdb.TableResourcesClient", "GetTableThroughput", err.Error())
	}

	req, err := client.GetTableThroughputPreparer(ctx, resourceGroupName, accountName, tableName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "GetTableThroughput", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetTableThroughputSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "GetTableThroughput", resp, "Failure sending request")
		return
	}

	result, err = client.GetTableThroughputResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "GetTableThroughput", resp, "Failure responding to request")
	}

	return
}

// GetTableThroughputPreparer prepares the GetTableThroughput request.
func (client TableResourcesClient) GetTableThroughputPreparer(ctx context.Context, resourceGroupName string, accountName string, tableName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tableName":         autorest.Encode("path", tableName),
	}

	const APIVersion = "2015-04-08"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/tables/{tableName}/settings/throughput", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetTableThroughputSender sends the GetTableThroughput request. The method will close the
// http.Response Body if it receives an error.
func (client TableResourcesClient) GetTableThroughputSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// GetTableThroughputResponder handles the response to the GetTableThroughput request. The method always
// closes the http.Response Body.
func (client TableResourcesClient) GetTableThroughputResponder(resp *http.Response) (result Throughput, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListTables lists the Tables under an existing Azure Cosmos DB database account.
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
func (client TableResourcesClient) ListTables(ctx context.Context, resourceGroupName string, accountName string) (result TableListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TableResourcesClient.ListTables")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "accountName", Name: validation.Pattern, Rule: `^[a-z0-9]+(-[a-z0-9]+)*`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("documentdb.TableResourcesClient", "ListTables", err.Error())
	}

	req, err := client.ListTablesPreparer(ctx, resourceGroupName, accountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "ListTables", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListTablesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "ListTables", resp, "Failure sending request")
		return
	}

	result, err = client.ListTablesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "ListTables", resp, "Failure responding to request")
	}

	return
}

// ListTablesPreparer prepares the ListTables request.
func (client TableResourcesClient) ListTablesPreparer(ctx context.Context, resourceGroupName string, accountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-04-08"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/tables", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListTablesSender sends the ListTables request. The method will close the
// http.Response Body if it receives an error.
func (client TableResourcesClient) ListTablesSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// ListTablesResponder handles the response to the ListTables request. The method always
// closes the http.Response Body.
func (client TableResourcesClient) ListTablesResponder(resp *http.Response) (result TableListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// UpdateTableThroughput update RUs per second of an Azure Cosmos DB Table
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
// tableName - cosmos DB table name.
// updateThroughputParameters - the parameters to provide for the RUs per second of the current Table.
func (client TableResourcesClient) UpdateTableThroughput(ctx context.Context, resourceGroupName string, accountName string, tableName string, updateThroughputParameters ThroughputUpdateParameters) (result TableResourcesUpdateTableThroughputFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TableResourcesClient.UpdateTableThroughput")
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
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil},
				{Target: "accountName", Name: validation.Pattern, Rule: `^[a-z0-9]+(-[a-z0-9]+)*`, Chain: nil}}},
		{TargetValue: updateThroughputParameters,
			Constraints: []validation.Constraint{{Target: "updateThroughputParameters.ThroughputUpdateProperties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "updateThroughputParameters.ThroughputUpdateProperties.Resource", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "updateThroughputParameters.ThroughputUpdateProperties.Resource.Throughput", Name: validation.Null, Rule: true, Chain: nil}}},
				}}}}}); err != nil {
		return result, validation.NewError("documentdb.TableResourcesClient", "UpdateTableThroughput", err.Error())
	}

	req, err := client.UpdateTableThroughputPreparer(ctx, resourceGroupName, accountName, tableName, updateThroughputParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "UpdateTableThroughput", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateTableThroughputSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.TableResourcesClient", "UpdateTableThroughput", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdateTableThroughputPreparer prepares the UpdateTableThroughput request.
func (client TableResourcesClient) UpdateTableThroughputPreparer(ctx context.Context, resourceGroupName string, accountName string, tableName string, updateThroughputParameters ThroughputUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tableName":         autorest.Encode("path", tableName),
	}

	const APIVersion = "2015-04-08"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/tables/{tableName}/settings/throughput", pathParameters),
		autorest.WithJSON(updateThroughputParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateTableThroughputSender sends the UpdateTableThroughput request. The method will close the
// http.Response Body if it receives an error.
func (client TableResourcesClient) UpdateTableThroughputSender(req *http.Request) (future TableResourcesUpdateTableThroughputFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// UpdateTableThroughputResponder handles the response to the UpdateTableThroughput request. The method always
// closes the http.Response Body.
func (client TableResourcesClient) UpdateTableThroughputResponder(resp *http.Response) (result Throughput, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
