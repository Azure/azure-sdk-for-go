package sql

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
    "github.com/Azure/go-autorest/autorest"
    "github.com/Azure/go-autorest/autorest/azure"
    "net/http"
    "context"
    "github.com/Azure/go-autorest/tracing"
)

// ElasticPoolDatabaseActivitiesClient is the the Azure SQL Database management API provides a RESTful set of web
// services that interact with Azure SQL Database services to manage your databases. The API enables you to create,
// retrieve, update, and delete databases.
type ElasticPoolDatabaseActivitiesClient struct {
    BaseClient
}
// NewElasticPoolDatabaseActivitiesClient creates an instance of the ElasticPoolDatabaseActivitiesClient client.
func NewElasticPoolDatabaseActivitiesClient(subscriptionID string) ElasticPoolDatabaseActivitiesClient {
    return NewElasticPoolDatabaseActivitiesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewElasticPoolDatabaseActivitiesClientWithBaseURI creates an instance of the ElasticPoolDatabaseActivitiesClient
// client.
    func NewElasticPoolDatabaseActivitiesClientWithBaseURI(baseURI string, subscriptionID string) ElasticPoolDatabaseActivitiesClient {
        return ElasticPoolDatabaseActivitiesClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// ListByElasticPool returns activity on databases inside of an elastic pool.
    // Parameters:
        // resourceGroupName - the name of the resource group that contains the resource. You can obtain this value
        // from the Azure Resource Manager API or the portal.
        // serverName - the name of the server.
        // elasticPoolName - the name of the elastic pool.
func (client ElasticPoolDatabaseActivitiesClient) ListByElasticPool(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string) (result ElasticPoolDatabaseActivityListResult, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ElasticPoolDatabaseActivitiesClient.ListByElasticPool")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.ListByElasticPoolPreparer(ctx, resourceGroupName, serverName, elasticPoolName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "sql.ElasticPoolDatabaseActivitiesClient", "ListByElasticPool", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByElasticPoolSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "sql.ElasticPoolDatabaseActivitiesClient", "ListByElasticPool", resp, "Failure sending request")
            return
            }

            result, err = client.ListByElasticPoolResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "sql.ElasticPoolDatabaseActivitiesClient", "ListByElasticPool", resp, "Failure responding to request")
            }

    return
    }

    // ListByElasticPoolPreparer prepares the ListByElasticPool request.
    func (client ElasticPoolDatabaseActivitiesClient) ListByElasticPoolPreparer(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "elasticPoolName": autorest.Encode("path",elasticPoolName),
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "serverName": autorest.Encode("path",serverName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2014-04-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}/elasticPoolDatabaseActivity",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByElasticPoolSender sends the ListByElasticPool request. The method will close the
    // http.Response Body if it receives an error.
    func (client ElasticPoolDatabaseActivitiesClient) ListByElasticPoolSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByElasticPoolResponder handles the response to the ListByElasticPool request. The method always
// closes the http.Response Body.
func (client ElasticPoolDatabaseActivitiesClient) ListByElasticPoolResponder(resp *http.Response) (result ElasticPoolDatabaseActivityListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

