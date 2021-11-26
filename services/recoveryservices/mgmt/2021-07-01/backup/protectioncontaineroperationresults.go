package backup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
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

// ProtectionContainerOperationResultsClient is the open API 2.0 Specs for Azure RecoveryServices Backup service
type ProtectionContainerOperationResultsClient struct {
	BaseClient
}

// NewProtectionContainerOperationResultsClient creates an instance of the ProtectionContainerOperationResultsClient
// client.
func NewProtectionContainerOperationResultsClient(subscriptionID string) ProtectionContainerOperationResultsClient {
	return NewProtectionContainerOperationResultsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewProtectionContainerOperationResultsClientWithBaseURI creates an instance of the
// ProtectionContainerOperationResultsClient client using a custom endpoint.  Use this when interacting with an Azure
// cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewProtectionContainerOperationResultsClientWithBaseURI(baseURI string, subscriptionID string) ProtectionContainerOperationResultsClient {
	return ProtectionContainerOperationResultsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get fetches the result of any operation on the container.
// Parameters:
// vaultName - the name of the recovery services vault.
// resourceGroupName - the name of the resource group where the recovery services vault is present.
// fabricName - fabric name associated with the container.
// containerName - container name whose information should be fetched.
// operationID - operation ID which represents the operation whose result needs to be fetched.
func (client ProtectionContainerOperationResultsClient) Get(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string, operationID string) (result ProtectionContainerResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProtectionContainerOperationResultsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, vaultName, resourceGroupName, fabricName, containerName, operationID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainerOperationResultsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainerOperationResultsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainerOperationResultsClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client ProtectionContainerOperationResultsClient) GetPreparer(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string, operationID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerName":     autorest.Encode("path", containerName),
		"fabricName":        autorest.Encode("path", fabricName),
		"operationId":       autorest.Encode("path", operationID),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2021-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupFabrics/{fabricName}/protectionContainers/{containerName}/operationResults/{operationId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ProtectionContainerOperationResultsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ProtectionContainerOperationResultsClient) GetResponder(resp *http.Response) (result ProtectionContainerResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
