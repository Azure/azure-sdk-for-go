package compute

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

// VirtualMachineScaleSetRollingUpgradesClient is the compute Client
type VirtualMachineScaleSetRollingUpgradesClient struct {
	BaseClient
}

// NewVirtualMachineScaleSetRollingUpgradesClient creates an instance of the
// VirtualMachineScaleSetRollingUpgradesClient client.
func NewVirtualMachineScaleSetRollingUpgradesClient(subscriptionID string) VirtualMachineScaleSetRollingUpgradesClient {
	return NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI creates an instance of the
// VirtualMachineScaleSetRollingUpgradesClient client using a custom endpoint.  Use this when interacting with an Azure
// cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI(baseURI string, subscriptionID string) VirtualMachineScaleSetRollingUpgradesClient {
	return VirtualMachineScaleSetRollingUpgradesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Cancel cancels the current virtual machine scale set rolling upgrade.
// Parameters:
// resourceGroupName - the name of the resource group.
// VMScaleSetName - the name of the VM scale set.
func (client VirtualMachineScaleSetRollingUpgradesClient) Cancel(ctx context.Context, resourceGroupName string, VMScaleSetName string) (result VirtualMachineScaleSetRollingUpgradesCancelFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualMachineScaleSetRollingUpgradesClient.Cancel")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CancelPreparer(ctx, resourceGroupName, VMScaleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "Cancel", nil, "Failure preparing request")
		return
	}

	result, err = client.CancelSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "Cancel", result.Response(), "Failure sending request")
		return
	}

	return
}

// CancelPreparer prepares the Cancel request.
func (client VirtualMachineScaleSetRollingUpgradesClient) CancelPreparer(ctx context.Context, resourceGroupName string, VMScaleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vmScaleSetName":    autorest.Encode("path", VMScaleSetName),
	}

	const APIVersion = "2017-03-30"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/rollingUpgrades/cancel", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CancelSender sends the Cancel request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineScaleSetRollingUpgradesClient) CancelSender(req *http.Request) (future VirtualMachineScaleSetRollingUpgradesCancelFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// CancelResponder handles the response to the Cancel request. The method always
// closes the http.Response Body.
func (client VirtualMachineScaleSetRollingUpgradesClient) CancelResponder(resp *http.Response) (result OperationStatusResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetLatest gets the status of the latest virtual machine scale set rolling upgrade.
// Parameters:
// resourceGroupName - the name of the resource group.
// VMScaleSetName - the name of the VM scale set.
func (client VirtualMachineScaleSetRollingUpgradesClient) GetLatest(ctx context.Context, resourceGroupName string, VMScaleSetName string) (result RollingUpgradeStatusInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualMachineScaleSetRollingUpgradesClient.GetLatest")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetLatestPreparer(ctx, resourceGroupName, VMScaleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "GetLatest", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetLatestSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "GetLatest", resp, "Failure sending request")
		return
	}

	result, err = client.GetLatestResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "GetLatest", resp, "Failure responding to request")
		return
	}

	return
}

// GetLatestPreparer prepares the GetLatest request.
func (client VirtualMachineScaleSetRollingUpgradesClient) GetLatestPreparer(ctx context.Context, resourceGroupName string, VMScaleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vmScaleSetName":    autorest.Encode("path", VMScaleSetName),
	}

	const APIVersion = "2017-03-30"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/rollingUpgrades/latest", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetLatestSender sends the GetLatest request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineScaleSetRollingUpgradesClient) GetLatestSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetLatestResponder handles the response to the GetLatest request. The method always
// closes the http.Response Body.
func (client VirtualMachineScaleSetRollingUpgradesClient) GetLatestResponder(resp *http.Response) (result RollingUpgradeStatusInfo, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// StartOSUpgrade starts a rolling upgrade to move all virtual machine scale set instances to the latest available
// Platform Image OS version. Instances which are already running the latest available OS version are not affected.
// Parameters:
// resourceGroupName - the name of the resource group.
// VMScaleSetName - the name of the VM scale set.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartOSUpgrade(ctx context.Context, resourceGroupName string, VMScaleSetName string) (result VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualMachineScaleSetRollingUpgradesClient.StartOSUpgrade")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.StartOSUpgradePreparer(ctx, resourceGroupName, VMScaleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "StartOSUpgrade", nil, "Failure preparing request")
		return
	}

	result, err = client.StartOSUpgradeSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "StartOSUpgrade", result.Response(), "Failure sending request")
		return
	}

	return
}

// StartOSUpgradePreparer prepares the StartOSUpgrade request.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartOSUpgradePreparer(ctx context.Context, resourceGroupName string, VMScaleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vmScaleSetName":    autorest.Encode("path", VMScaleSetName),
	}

	const APIVersion = "2017-03-30"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/osRollingUpgrade", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// StartOSUpgradeSender sends the StartOSUpgrade request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartOSUpgradeSender(req *http.Request) (future VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// StartOSUpgradeResponder handles the response to the StartOSUpgrade request. The method always
// closes the http.Response Body.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartOSUpgradeResponder(resp *http.Response) (result OperationStatusResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
