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

// ProtectionContainersClient is the open API 2.0 Specs for Azure RecoveryServices Backup service
type ProtectionContainersClient struct {
	BaseClient
}

// NewProtectionContainersClient creates an instance of the ProtectionContainersClient client.
func NewProtectionContainersClient(subscriptionID string) ProtectionContainersClient {
	return NewProtectionContainersClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewProtectionContainersClientWithBaseURI creates an instance of the ProtectionContainersClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewProtectionContainersClientWithBaseURI(baseURI string, subscriptionID string) ProtectionContainersClient {
	return ProtectionContainersClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets details of the specific container registered to your Recovery Services Vault.
// Parameters:
// vaultName - the name of the recovery services vault.
// resourceGroupName - the name of the resource group where the recovery services vault is present.
// fabricName - name of the fabric where the container belongs.
// containerName - name of the container whose details need to be fetched.
func (client ProtectionContainersClient) Get(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string) (result ProtectionContainerResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProtectionContainersClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, vaultName, resourceGroupName, fabricName, containerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client ProtectionContainersClient) GetPreparer(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerName":     autorest.Encode("path", containerName),
		"fabricName":        autorest.Encode("path", fabricName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2021-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupFabrics/{fabricName}/protectionContainers/{containerName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ProtectionContainersClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ProtectionContainersClient) GetResponder(resp *http.Response) (result ProtectionContainerResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Inquire this is an async operation and the results should be tracked using location header or Azure-async-url.
// Parameters:
// vaultName - the name of the recovery services vault.
// resourceGroupName - the name of the resource group where the recovery services vault is present.
// fabricName - fabric Name associated with the container.
// containerName - name of the container in which inquiry needs to be triggered.
// filter - oData filter options.
func (client ProtectionContainersClient) Inquire(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string, filter string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProtectionContainersClient.Inquire")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.InquirePreparer(ctx, vaultName, resourceGroupName, fabricName, containerName, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Inquire", nil, "Failure preparing request")
		return
	}

	resp, err := client.InquireSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Inquire", resp, "Failure sending request")
		return
	}

	result, err = client.InquireResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Inquire", resp, "Failure responding to request")
		return
	}

	return
}

// InquirePreparer prepares the Inquire request.
func (client ProtectionContainersClient) InquirePreparer(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerName":     autorest.Encode("path", containerName),
		"fabricName":        autorest.Encode("path", fabricName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2021-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupFabrics/{fabricName}/protectionContainers/{containerName}/inquire", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// InquireSender sends the Inquire request. The method will close the
// http.Response Body if it receives an error.
func (client ProtectionContainersClient) InquireSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// InquireResponder handles the response to the Inquire request. The method always
// closes the http.Response Body.
func (client ProtectionContainersClient) InquireResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Refresh discovers all the containers in the subscription that can be backed up to Recovery Services Vault. This is
// an
// asynchronous operation. To know the status of the operation, call GetRefreshOperationResult API.
// Parameters:
// vaultName - the name of the recovery services vault.
// resourceGroupName - the name of the resource group where the recovery services vault is present.
// fabricName - fabric name associated the container.
// filter - oData filter options.
func (client ProtectionContainersClient) Refresh(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, filter string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProtectionContainersClient.Refresh")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.RefreshPreparer(ctx, vaultName, resourceGroupName, fabricName, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Refresh", nil, "Failure preparing request")
		return
	}

	resp, err := client.RefreshSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Refresh", resp, "Failure sending request")
		return
	}

	result, err = client.RefreshResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Refresh", resp, "Failure responding to request")
		return
	}

	return
}

// RefreshPreparer prepares the Refresh request.
func (client ProtectionContainersClient) RefreshPreparer(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fabricName":        autorest.Encode("path", fabricName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2021-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupFabrics/{fabricName}/refreshContainers", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RefreshSender sends the Refresh request. The method will close the
// http.Response Body if it receives an error.
func (client ProtectionContainersClient) RefreshSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// RefreshResponder handles the response to the Refresh request. The method always
// closes the http.Response Body.
func (client ProtectionContainersClient) RefreshResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Register registers the container with Recovery Services vault.
// This is an asynchronous operation. To track the operation status, use location header to call get latest status of
// the operation.
// Parameters:
// vaultName - the name of the recovery services vault.
// resourceGroupName - the name of the resource group where the recovery services vault is present.
// fabricName - fabric name associated with the container.
// containerName - name of the container to be registered.
// parameters - request body for operation
func (client ProtectionContainersClient) Register(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string, parameters ProtectionContainerResource) (result ProtectionContainerResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProtectionContainersClient.Register")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.RegisterPreparer(ctx, vaultName, resourceGroupName, fabricName, containerName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Register", nil, "Failure preparing request")
		return
	}

	resp, err := client.RegisterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Register", resp, "Failure sending request")
		return
	}

	result, err = client.RegisterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Register", resp, "Failure responding to request")
		return
	}

	return
}

// RegisterPreparer prepares the Register request.
func (client ProtectionContainersClient) RegisterPreparer(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string, parameters ProtectionContainerResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerName":     autorest.Encode("path", containerName),
		"fabricName":        autorest.Encode("path", fabricName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2021-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupFabrics/{fabricName}/protectionContainers/{containerName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RegisterSender sends the Register request. The method will close the
// http.Response Body if it receives an error.
func (client ProtectionContainersClient) RegisterSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// RegisterResponder handles the response to the Register request. The method always
// closes the http.Response Body.
func (client ProtectionContainersClient) RegisterResponder(resp *http.Response) (result ProtectionContainerResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Unregister unregisters the given container from your Recovery Services Vault. This is an asynchronous operation. To
// determine
// whether the backend service has finished processing the request, call Get Container Operation Result API.
// Parameters:
// vaultName - the name of the recovery services vault.
// resourceGroupName - the name of the resource group where the recovery services vault is present.
// fabricName - name of the fabric where the container belongs.
// containerName - name of the container which needs to be unregistered from the Recovery Services Vault.
func (client ProtectionContainersClient) Unregister(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProtectionContainersClient.Unregister")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UnregisterPreparer(ctx, vaultName, resourceGroupName, fabricName, containerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Unregister", nil, "Failure preparing request")
		return
	}

	resp, err := client.UnregisterSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Unregister", resp, "Failure sending request")
		return
	}

	result, err = client.UnregisterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ProtectionContainersClient", "Unregister", resp, "Failure responding to request")
		return
	}

	return
}

// UnregisterPreparer prepares the Unregister request.
func (client ProtectionContainersClient) UnregisterPreparer(ctx context.Context, vaultName string, resourceGroupName string, fabricName string, containerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerName":     autorest.Encode("path", containerName),
		"fabricName":        autorest.Encode("path", fabricName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2021-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupFabrics/{fabricName}/protectionContainers/{containerName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UnregisterSender sends the Unregister request. The method will close the
// http.Response Body if it receives an error.
func (client ProtectionContainersClient) UnregisterSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// UnregisterResponder handles the response to the Unregister request. The method always
// closes the http.Response Body.
func (client ProtectionContainersClient) UnregisterResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}
