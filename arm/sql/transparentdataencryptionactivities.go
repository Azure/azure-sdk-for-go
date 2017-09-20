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
// Code generated by Microsoft (R) AutoRest Code Generator 2.2.21.0
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// TransparentDataEncryptionActivitiesClient is the the Azure SQL Database management API provides a RESTful set of web
// services that interact with Azure SQL Database services to manage your databases. The API enables you to create,
// retrieve, update, and delete databases.
type TransparentDataEncryptionActivitiesClient struct {
	ManagementClient
}

// NewTransparentDataEncryptionActivitiesClient creates an instance of the TransparentDataEncryptionActivitiesClient
// client.
func NewTransparentDataEncryptionActivitiesClient(subscriptionID string) TransparentDataEncryptionActivitiesClient {
	return NewTransparentDataEncryptionActivitiesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewTransparentDataEncryptionActivitiesClientWithBaseURI creates an instance of the
// TransparentDataEncryptionActivitiesClient client.
func NewTransparentDataEncryptionActivitiesClientWithBaseURI(baseURI string, subscriptionID string) TransparentDataEncryptionActivitiesClient {
	return TransparentDataEncryptionActivitiesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ListByConfiguration returns a database's transparent data encryption operation result.
//
// resourceGroupName is the name of the resource group that contains the resource. You can obtain this value from the
// Azure Resource Manager API or the portal. serverName is the name of the server. databaseName is the name of the
// database for which the transparent data encryption applies. transparentDataEncryptionName is the name of the
// transparent data encryption configuration.
func (client TransparentDataEncryptionActivitiesClient) ListByConfiguration(resourceGroupName string, serverName string, databaseName string, transparentDataEncryptionName string) (result TransparentDataEncryptionActivityListResult, err error) {
	req, err := client.ListByConfigurationPreparer(resourceGroupName, serverName, databaseName, transparentDataEncryptionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sql.TransparentDataEncryptionActivitiesClient", "ListByConfiguration", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByConfigurationSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "sql.TransparentDataEncryptionActivitiesClient", "ListByConfiguration", resp, "Failure sending request")
		return
	}

	result, err = client.ListByConfigurationResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sql.TransparentDataEncryptionActivitiesClient", "ListByConfiguration", resp, "Failure responding to request")
	}

	return
}

// ListByConfigurationPreparer prepares the ListByConfiguration request.
func (client TransparentDataEncryptionActivitiesClient) ListByConfigurationPreparer(resourceGroupName string, serverName string, databaseName string, transparentDataEncryptionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"databaseName":                  autorest.Encode("path", databaseName),
		"resourceGroupName":             autorest.Encode("path", resourceGroupName),
		"serverName":                    autorest.Encode("path", serverName),
		"subscriptionId":                autorest.Encode("path", client.SubscriptionID),
		"transparentDataEncryptionName": autorest.Encode("path", transparentDataEncryptionName),
	}

	const APIVersion = "2014-04-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/transparentDataEncryption/{transparentDataEncryptionName}/operationResults", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByConfigurationSender sends the ListByConfiguration request. The method will close the
// http.Response Body if it receives an error.
func (client TransparentDataEncryptionActivitiesClient) ListByConfigurationSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListByConfigurationResponder handles the response to the ListByConfiguration request. The method always
// closes the http.Response Body.
func (client TransparentDataEncryptionActivitiesClient) ListByConfigurationResponder(resp *http.Response) (result TransparentDataEncryptionActivityListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
