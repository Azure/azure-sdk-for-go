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

// SQLPoolTablesClient is the azure Synapse Analytics Management Client
type SQLPoolTablesClient struct {
	BaseClient
}

// NewSQLPoolTablesClient creates an instance of the SQLPoolTablesClient client.
func NewSQLPoolTablesClient(subscriptionID string) SQLPoolTablesClient {
	return NewSQLPoolTablesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewSQLPoolTablesClientWithBaseURI creates an instance of the SQLPoolTablesClient client using a custom endpoint.
// Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewSQLPoolTablesClientWithBaseURI(baseURI string, subscriptionID string) SQLPoolTablesClient {
	return SQLPoolTablesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get get Sql pool table
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - the name of the workspace.
// SQLPoolName - SQL pool name
// schemaName - the name of the schema.
// tableName - the name of the table.
func (client SQLPoolTablesClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, SQLPoolName string, schemaName string, tableName string) (result SQLPoolTable, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SQLPoolTablesClient.Get")
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
		return result, validation.NewError("synapse.SQLPoolTablesClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, workspaceName, SQLPoolName, schemaName, tableName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client SQLPoolTablesClient) GetPreparer(ctx context.Context, resourceGroupName string, workspaceName string, SQLPoolName string, schemaName string, tableName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"schemaName":        autorest.Encode("path", schemaName),
		"sqlPoolName":       autorest.Encode("path", SQLPoolName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tableName":         autorest.Encode("path", tableName),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/schemas/{schemaName}/tables/{tableName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client SQLPoolTablesClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client SQLPoolTablesClient) GetResponder(resp *http.Response) (result SQLPoolTable, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListBySchema gets tables of a given schema in a SQL pool.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// workspaceName - the name of the workspace.
// SQLPoolName - SQL pool name
// schemaName - the name of the schema.
// filter - an OData filter expression that filters elements in the collection.
func (client SQLPoolTablesClient) ListBySchema(ctx context.Context, resourceGroupName string, workspaceName string, SQLPoolName string, schemaName string, filter string) (result SQLPoolTableListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SQLPoolTablesClient.ListBySchema")
		defer func() {
			sc := -1
			if result.sptlr.Response.Response != nil {
				sc = result.sptlr.Response.Response.StatusCode
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
		return result, validation.NewError("synapse.SQLPoolTablesClient", "ListBySchema", err.Error())
	}

	result.fn = client.listBySchemaNextResults
	req, err := client.ListBySchemaPreparer(ctx, resourceGroupName, workspaceName, SQLPoolName, schemaName, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "ListBySchema", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListBySchemaSender(req)
	if err != nil {
		result.sptlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "ListBySchema", resp, "Failure sending request")
		return
	}

	result.sptlr, err = client.ListBySchemaResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "ListBySchema", resp, "Failure responding to request")
		return
	}
	if result.sptlr.hasNextLink() && result.sptlr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListBySchemaPreparer prepares the ListBySchema request.
func (client SQLPoolTablesClient) ListBySchemaPreparer(ctx context.Context, resourceGroupName string, workspaceName string, SQLPoolName string, schemaName string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"schemaName":        autorest.Encode("path", schemaName),
		"sqlPoolName":       autorest.Encode("path", SQLPoolName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/schemas/{schemaName}/tables", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListBySchemaSender sends the ListBySchema request. The method will close the
// http.Response Body if it receives an error.
func (client SQLPoolTablesClient) ListBySchemaSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListBySchemaResponder handles the response to the ListBySchema request. The method always
// closes the http.Response Body.
func (client SQLPoolTablesClient) ListBySchemaResponder(resp *http.Response) (result SQLPoolTableListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listBySchemaNextResults retrieves the next set of results, if any.
func (client SQLPoolTablesClient) listBySchemaNextResults(ctx context.Context, lastResults SQLPoolTableListResult) (result SQLPoolTableListResult, err error) {
	req, err := lastResults.sQLPoolTableListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "listBySchemaNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListBySchemaSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "listBySchemaNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListBySchemaResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synapse.SQLPoolTablesClient", "listBySchemaNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListBySchemaComplete enumerates all values, automatically crossing page boundaries as required.
func (client SQLPoolTablesClient) ListBySchemaComplete(ctx context.Context, resourceGroupName string, workspaceName string, SQLPoolName string, schemaName string, filter string) (result SQLPoolTableListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SQLPoolTablesClient.ListBySchema")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListBySchema(ctx, resourceGroupName, workspaceName, SQLPoolName, schemaName, filter)
	return
}
