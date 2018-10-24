package compute

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

// VirtualMachineExtensionsClient is the compute Client
type VirtualMachineExtensionsClient struct {
    BaseClient
}
// NewVirtualMachineExtensionsClient creates an instance of the VirtualMachineExtensionsClient client.
func NewVirtualMachineExtensionsClient(subscriptionID string) VirtualMachineExtensionsClient {
    return NewVirtualMachineExtensionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewVirtualMachineExtensionsClientWithBaseURI creates an instance of the VirtualMachineExtensionsClient client.
    func NewVirtualMachineExtensionsClientWithBaseURI(baseURI string, subscriptionID string) VirtualMachineExtensionsClient {
        return VirtualMachineExtensionsClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// CreateOrUpdate the operation to create or update the extension.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // VMName - the name of the virtual machine where the extension should be created or updated.
        // VMExtensionName - the name of the virtual machine extension.
        // extensionParameters - parameters supplied to the Create Virtual Machine Extension operation.
func (client VirtualMachineExtensionsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, VMName string, VMExtensionName string, extensionParameters VirtualMachineExtension) (result VirtualMachineExtensionsCreateOrUpdateFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualMachineExtensionsClient.CreateOrUpdate")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, VMName, VMExtensionName, extensionParameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "CreateOrUpdate", nil , "Failure preparing request")
    return
    }

            result, err = client.CreateOrUpdateSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "CreateOrUpdate", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // CreateOrUpdatePreparer prepares the CreateOrUpdate request.
    func (client VirtualMachineExtensionsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, VMName string, VMExtensionName string, extensionParameters VirtualMachineExtension) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "vmExtensionName": autorest.Encode("path",VMExtensionName),
            "vmName": autorest.Encode("path",VMName),
            }

                        const APIVersion = "2018-04-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPut(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{vmName}/extensions/{vmExtensionName}",pathParameters),
    autorest.WithJSON(extensionParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualMachineExtensionsClient) CreateOrUpdateSender(req *http.Request) (future VirtualMachineExtensionsCreateOrUpdateFuture, err error) {
            var resp *http.Response
            resp, err = autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            if err != nil {
            return
            }
            future.Future, err = azure.NewFutureFromResponse(resp)
            return
            }

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client VirtualMachineExtensionsClient) CreateOrUpdateResponder(resp *http.Response) (result VirtualMachineExtension, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusCreated),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Delete the operation to delete the extension.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // VMName - the name of the virtual machine where the extension should be deleted.
        // VMExtensionName - the name of the virtual machine extension.
func (client VirtualMachineExtensionsClient) Delete(ctx context.Context, resourceGroupName string, VMName string, VMExtensionName string) (result VirtualMachineExtensionsDeleteFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualMachineExtensionsClient.Delete")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.DeletePreparer(ctx, resourceGroupName, VMName, VMExtensionName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "Delete", nil , "Failure preparing request")
    return
    }

            result, err = client.DeleteSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "Delete", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // DeletePreparer prepares the Delete request.
    func (client VirtualMachineExtensionsClient) DeletePreparer(ctx context.Context, resourceGroupName string, VMName string, VMExtensionName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "vmExtensionName": autorest.Encode("path",VMExtensionName),
            "vmName": autorest.Encode("path",VMName),
            }

                        const APIVersion = "2018-04-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsDelete(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{vmName}/extensions/{vmExtensionName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // DeleteSender sends the Delete request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualMachineExtensionsClient) DeleteSender(req *http.Request) (future VirtualMachineExtensionsDeleteFuture, err error) {
            var resp *http.Response
            resp, err = autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            if err != nil {
            return
            }
            future.Future, err = azure.NewFutureFromResponse(resp)
            return
            }

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client VirtualMachineExtensionsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK,http.StatusAccepted,http.StatusNoContent),
    autorest.ByClosing())
    result.Response = resp
        return
    }

// Get the operation to get the extension.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // VMName - the name of the virtual machine containing the extension.
        // VMExtensionName - the name of the virtual machine extension.
        // expand - the expand expression to apply on the operation.
func (client VirtualMachineExtensionsClient) Get(ctx context.Context, resourceGroupName string, VMName string, VMExtensionName string, expand string) (result VirtualMachineExtension, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualMachineExtensionsClient.Get")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.GetPreparer(ctx, resourceGroupName, VMName, VMExtensionName, expand)
    if err != nil {
    err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "Get", nil , "Failure preparing request")
    return
    }

            resp, err := client.GetSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "Get", resp, "Failure sending request")
            return
            }

            result, err = client.GetResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "Get", resp, "Failure responding to request")
            }

    return
    }

    // GetPreparer prepares the Get request.
    func (client VirtualMachineExtensionsClient) GetPreparer(ctx context.Context, resourceGroupName string, VMName string, VMExtensionName string, expand string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "vmExtensionName": autorest.Encode("path",VMExtensionName),
            "vmName": autorest.Encode("path",VMName),
            }

                        const APIVersion = "2018-04-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if len(expand) > 0 {
            queryParameters["$expand"] = autorest.Encode("query",expand)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{vmName}/extensions/{vmExtensionName}",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // GetSender sends the Get request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualMachineExtensionsClient) GetSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client VirtualMachineExtensionsClient) GetResponder(resp *http.Response) (result VirtualMachineExtension, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// List the operation to get all extensions of a Virtual Machine.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // VMName - the name of the virtual machine containing the extension.
        // expand - the expand expression to apply on the operation.
func (client VirtualMachineExtensionsClient) List(ctx context.Context, resourceGroupName string, VMName string, expand string) (result VirtualMachineExtensionsListResult, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualMachineExtensionsClient.List")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.ListPreparer(ctx, resourceGroupName, VMName, expand)
    if err != nil {
    err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "List", resp, "Failure sending request")
            return
            }

            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client VirtualMachineExtensionsClient) ListPreparer(ctx context.Context, resourceGroupName string, VMName string, expand string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "vmName": autorest.Encode("path",VMName),
            }

                        const APIVersion = "2018-04-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }
            if len(expand) > 0 {
            queryParameters["$expand"] = autorest.Encode("query",expand)
            }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{vmName}/extensions",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualMachineExtensionsClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client VirtualMachineExtensionsClient) ListResponder(resp *http.Response) (result VirtualMachineExtensionsListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

// Update the operation to update the extension.
    // Parameters:
        // resourceGroupName - the name of the resource group.
        // VMName - the name of the virtual machine where the extension should be updated.
        // VMExtensionName - the name of the virtual machine extension.
        // extensionParameters - parameters supplied to the Update Virtual Machine Extension operation.
func (client VirtualMachineExtensionsClient) Update(ctx context.Context, resourceGroupName string, VMName string, VMExtensionName string, extensionParameters VirtualMachineExtensionUpdate) (result VirtualMachineExtensionsUpdateFuture, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/VirtualMachineExtensionsClient.Update")
        defer func() {
            sc := -1
            if result.Response() != nil {
                sc = result.Response().StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.UpdatePreparer(ctx, resourceGroupName, VMName, VMExtensionName, extensionParameters)
    if err != nil {
    err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "Update", nil , "Failure preparing request")
    return
    }

            result, err = client.UpdateSender(req)
            if err != nil {
            err = autorest.NewErrorWithError(err, "compute.VirtualMachineExtensionsClient", "Update", result.Response(), "Failure sending request")
            return
            }

    return
    }

    // UpdatePreparer prepares the Update request.
    func (client VirtualMachineExtensionsClient) UpdatePreparer(ctx context.Context, resourceGroupName string, VMName string, VMExtensionName string, extensionParameters VirtualMachineExtensionUpdate) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "vmExtensionName": autorest.Encode("path",VMExtensionName),
            "vmName": autorest.Encode("path",VMName),
            }

                        const APIVersion = "2018-04-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsContentType("application/json; charset=utf-8"),
    autorest.AsPatch(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{vmName}/extensions/{vmExtensionName}",pathParameters),
    autorest.WithJSON(extensionParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // UpdateSender sends the Update request. The method will close the
    // http.Response Body if it receives an error.
    func (client VirtualMachineExtensionsClient) UpdateSender(req *http.Request) (future VirtualMachineExtensionsUpdateFuture, err error) {
            var resp *http.Response
            resp, err = autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            if err != nil {
            return
            }
            future.Future, err = azure.NewFutureFromResponse(resp)
            return
            }

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client VirtualMachineExtensionsClient) UpdateResponder(resp *http.Response) (result VirtualMachineExtension, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

