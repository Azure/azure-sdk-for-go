package mysql

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

// TopQueryStatisticsClient is the the Microsoft Azure management API provides create, read, update, and delete
// functionality for Azure MySQL resources including servers, databases, firewall rules, VNET rules, log files and
// configurations with new business model.
type TopQueryStatisticsClient struct {
	BaseClient
}

// NewTopQueryStatisticsClient creates an instance of the TopQueryStatisticsClient client.
func NewTopQueryStatisticsClient(subscriptionID string) TopQueryStatisticsClient {
	return NewTopQueryStatisticsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewTopQueryStatisticsClientWithBaseURI creates an instance of the TopQueryStatisticsClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewTopQueryStatisticsClientWithBaseURI(baseURI string, subscriptionID string) TopQueryStatisticsClient {
	return TopQueryStatisticsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get retrieve the query statistic for specified identifier.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// serverName - the name of the server.
// queryStatisticID - the Query Statistic identifier.
func (client TopQueryStatisticsClient) Get(ctx context.Context, resourceGroupName string, serverName string, queryStatisticID string) (result QueryStatistic, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TopQueryStatisticsClient.Get")
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
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("mysql.TopQueryStatisticsClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, serverName, queryStatisticID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client TopQueryStatisticsClient) GetPreparer(ctx context.Context, resourceGroupName string, serverName string, queryStatisticID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"queryStatisticId":  autorest.Encode("path", queryStatisticID),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serverName":        autorest.Encode("path", serverName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/servers/{serverName}/topQueryStatistics/{queryStatisticId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client TopQueryStatisticsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client TopQueryStatisticsClient) GetResponder(resp *http.Response) (result QueryStatistic, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByServer retrieve the Query-Store top queries for specified metric and aggregation.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// serverName - the name of the server.
// parameters - the required parameters for retrieving top query statistics.
func (client TopQueryStatisticsClient) ListByServer(ctx context.Context, resourceGroupName string, serverName string, parameters TopQueryStatisticsInput) (result TopQueryStatisticsResultListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TopQueryStatisticsClient.ListByServer")
		defer func() {
			sc := -1
			if result.tqsrl.Response.Response != nil {
				sc = result.tqsrl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.TopQueryStatisticsInputProperties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.TopQueryStatisticsInputProperties.NumberOfTopQueries", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.TopQueryStatisticsInputProperties.AggregationFunction", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.TopQueryStatisticsInputProperties.ObservedMetric", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.TopQueryStatisticsInputProperties.ObservationStartTime", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.TopQueryStatisticsInputProperties.ObservationEndTime", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.TopQueryStatisticsInputProperties.AggregationWindow", Name: validation.Null, Rule: true, Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("mysql.TopQueryStatisticsClient", "ListByServer", err.Error())
	}

	result.fn = client.listByServerNextResults
	req, err := client.ListByServerPreparer(ctx, resourceGroupName, serverName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "ListByServer", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByServerSender(req)
	if err != nil {
		result.tqsrl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "ListByServer", resp, "Failure sending request")
		return
	}

	result.tqsrl, err = client.ListByServerResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "ListByServer", resp, "Failure responding to request")
	}

	return
}

// ListByServerPreparer prepares the ListByServer request.
func (client TopQueryStatisticsClient) ListByServerPreparer(ctx context.Context, resourceGroupName string, serverName string, parameters TopQueryStatisticsInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serverName":        autorest.Encode("path", serverName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/servers/{serverName}/topQueryStatistics", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByServerSender sends the ListByServer request. The method will close the
// http.Response Body if it receives an error.
func (client TopQueryStatisticsClient) ListByServerSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByServerResponder handles the response to the ListByServer request. The method always
// closes the http.Response Body.
func (client TopQueryStatisticsClient) ListByServerResponder(resp *http.Response) (result TopQueryStatisticsResultList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByServerNextResults retrieves the next set of results, if any.
func (client TopQueryStatisticsClient) listByServerNextResults(ctx context.Context, lastResults TopQueryStatisticsResultList) (result TopQueryStatisticsResultList, err error) {
	req, err := lastResults.topQueryStatisticsResultListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "listByServerNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByServerSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "listByServerNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByServerResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mysql.TopQueryStatisticsClient", "listByServerNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByServerComplete enumerates all values, automatically crossing page boundaries as required.
func (client TopQueryStatisticsClient) ListByServerComplete(ctx context.Context, resourceGroupName string, serverName string, parameters TopQueryStatisticsInput) (result TopQueryStatisticsResultListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TopQueryStatisticsClient.ListByServer")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByServer(ctx, resourceGroupName, serverName, parameters)
	return
}
