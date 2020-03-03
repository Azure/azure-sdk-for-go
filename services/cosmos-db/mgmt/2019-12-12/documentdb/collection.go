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

// CollectionClient is the azure Cosmos DB Database Service Resource Provider REST API
type CollectionClient struct {
	BaseClient
}

// NewCollectionClient creates an instance of the CollectionClient client.
func NewCollectionClient(subscriptionID string, subscriptionID1 string) CollectionClient {
	return NewCollectionClientWithBaseURI(DefaultBaseURI, subscriptionID, subscriptionID1)
}

// NewCollectionClientWithBaseURI creates an instance of the CollectionClient client using a custom endpoint.  Use this
// when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewCollectionClientWithBaseURI(baseURI string, subscriptionID string, subscriptionID1 string) CollectionClient {
	return CollectionClient{NewWithBaseURI(baseURI, subscriptionID, subscriptionID1)}
}

// ListMetricDefinitions retrieves metric definitions for the given collection.
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
// databaseRid - cosmos DB database rid.
// collectionRid - cosmos DB collection rid.
func (client CollectionClient) ListMetricDefinitions(ctx context.Context, resourceGroupName string, accountName string, databaseRid string, collectionRid string) (result MetricDefinitionsListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CollectionClient.ListMetricDefinitions")
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
		return result, validation.NewError("documentdb.CollectionClient", "ListMetricDefinitions", err.Error())
	}

	req, err := client.ListMetricDefinitionsPreparer(ctx, resourceGroupName, accountName, databaseRid, collectionRid)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListMetricDefinitions", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListMetricDefinitionsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListMetricDefinitions", resp, "Failure sending request")
		return
	}

	result, err = client.ListMetricDefinitionsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListMetricDefinitions", resp, "Failure responding to request")
	}

	return
}

// ListMetricDefinitionsPreparer prepares the ListMetricDefinitions request.
func (client CollectionClient) ListMetricDefinitionsPreparer(ctx context.Context, resourceGroupName string, accountName string, databaseRid string, collectionRid string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"collectionRid":     autorest.Encode("path", collectionRid),
		"databaseRid":       autorest.Encode("path", databaseRid),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-12-12"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/databases/{databaseRid}/collections/{collectionRid}/metricDefinitions", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListMetricDefinitionsSender sends the ListMetricDefinitions request. The method will close the
// http.Response Body if it receives an error.
func (client CollectionClient) ListMetricDefinitionsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListMetricDefinitionsResponder handles the response to the ListMetricDefinitions request. The method always
// closes the http.Response Body.
func (client CollectionClient) ListMetricDefinitionsResponder(resp *http.Response) (result MetricDefinitionsListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListMetrics retrieves the metrics determined by the given filter for the given database account and collection.
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
// databaseRid - cosmos DB database rid.
// collectionRid - cosmos DB collection rid.
// filter - an OData filter expression that describes a subset of metrics to return. The parameters that can be
// filtered are name.value (name of the metric, can have an or of multiple names), startTime, endTime, and
// timeGrain. The supported operator is eq.
func (client CollectionClient) ListMetrics(ctx context.Context, resourceGroupName string, accountName string, databaseRid string, collectionRid string, filter string) (result MetricListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CollectionClient.ListMetrics")
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
		return result, validation.NewError("documentdb.CollectionClient", "ListMetrics", err.Error())
	}

	req, err := client.ListMetricsPreparer(ctx, resourceGroupName, accountName, databaseRid, collectionRid, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListMetrics", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListMetricsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListMetrics", resp, "Failure sending request")
		return
	}

	result, err = client.ListMetricsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListMetrics", resp, "Failure responding to request")
	}

	return
}

// ListMetricsPreparer prepares the ListMetrics request.
func (client CollectionClient) ListMetricsPreparer(ctx context.Context, resourceGroupName string, accountName string, databaseRid string, collectionRid string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"collectionRid":     autorest.Encode("path", collectionRid),
		"databaseRid":       autorest.Encode("path", databaseRid),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-12-12"
	queryParameters := map[string]interface{}{
		"$filter":     autorest.Encode("query", filter),
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/databases/{databaseRid}/collections/{collectionRid}/metrics", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListMetricsSender sends the ListMetrics request. The method will close the
// http.Response Body if it receives an error.
func (client CollectionClient) ListMetricsSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListMetricsResponder handles the response to the ListMetrics request. The method always
// closes the http.Response Body.
func (client CollectionClient) ListMetricsResponder(resp *http.Response) (result MetricListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListUsages retrieves the usages (most recent storage data) for the given collection.
// Parameters:
// resourceGroupName - name of an Azure resource group.
// accountName - cosmos DB database account name.
// databaseRid - cosmos DB database rid.
// collectionRid - cosmos DB collection rid.
// filter - an OData filter expression that describes a subset of usages to return. The supported parameter is
// name.value (name of the metric, can have an or of multiple names).
func (client CollectionClient) ListUsages(ctx context.Context, resourceGroupName string, accountName string, databaseRid string, collectionRid string, filter string) (result UsagesResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/CollectionClient.ListUsages")
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
		return result, validation.NewError("documentdb.CollectionClient", "ListUsages", err.Error())
	}

	req, err := client.ListUsagesPreparer(ctx, resourceGroupName, accountName, databaseRid, collectionRid, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListUsages", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListUsagesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListUsages", resp, "Failure sending request")
		return
	}

	result, err = client.ListUsagesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "documentdb.CollectionClient", "ListUsages", resp, "Failure responding to request")
	}

	return
}

// ListUsagesPreparer prepares the ListUsages request.
func (client CollectionClient) ListUsagesPreparer(ctx context.Context, resourceGroupName string, accountName string, databaseRid string, collectionRid string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"collectionRid":     autorest.Encode("path", collectionRid),
		"databaseRid":       autorest.Encode("path", databaseRid),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-12-12"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/databases/{databaseRid}/collections/{collectionRid}/usages", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListUsagesSender sends the ListUsages request. The method will close the
// http.Response Body if it receives an error.
func (client CollectionClient) ListUsagesSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListUsagesResponder handles the response to the ListUsages request. The method always
// closes the http.Response Body.
func (client CollectionClient) ListUsagesResponder(resp *http.Response) (result UsagesResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
