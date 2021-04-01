package visualstudio

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

// ExtensionsClient is the use these APIs to manage Visual Studio Team Services resources through the Azure Resource
// Manager. All task operations conform to the HTTP/1.1 protocol specification and each operation returns an
// x-ms-request-id header that can be used to obtain information about the request. You must make sure that requests
// made to these resources are secure. For more information, see https://docs.microsoft.com/en-us/rest/api/index.
type ExtensionsClient struct {
	BaseClient
}

// NewExtensionsClient creates an instance of the ExtensionsClient client.
func NewExtensionsClient(subscriptionID string) ExtensionsClient {
	return NewExtensionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewExtensionsClientWithBaseURI creates an instance of the ExtensionsClient client using a custom endpoint.  Use this
// when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewExtensionsClientWithBaseURI(baseURI string, subscriptionID string) ExtensionsClient {
	return ExtensionsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Create registers the extension with a Visual Studio Team Services account.
// Parameters:
// resourceGroupName - name of the resource group within the Azure subscription.
// body - an object containing additional information related to the extension request.
// accountResourceName - the name of the Visual Studio Team Services account resource.
// extensionResourceName - the name of the extension.
func (client ExtensionsClient) Create(ctx context.Context, resourceGroupName string, body ExtensionResourceRequest, accountResourceName string, extensionResourceName string) (result ExtensionResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ExtensionsClient.Create")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CreatePreparer(ctx, resourceGroupName, body, accountResourceName, extensionResourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Create", resp, "Failure responding to request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client ExtensionsClient) CreatePreparer(ctx context.Context, resourceGroupName string, body ExtensionResourceRequest, accountResourceName string, extensionResourceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountResourceName":   autorest.Encode("path", accountResourceName),
		"extensionResourceName": autorest.Encode("path", extensionResourceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2014-04-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/microsoft.visualstudio/account/{accountResourceName}/extension/{extensionResourceName}", pathParameters),
		autorest.WithJSON(body),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client ExtensionsClient) CreateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client ExtensionsClient) CreateResponder(resp *http.Response) (result ExtensionResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete removes an extension resource registration for a Visual Studio Team Services account.
// Parameters:
// resourceGroupName - name of the resource group within the Azure subscription.
// accountResourceName - the name of the Visual Studio Team Services account resource.
// extensionResourceName - the name of the extension.
func (client ExtensionsClient) Delete(ctx context.Context, resourceGroupName string, accountResourceName string, extensionResourceName string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ExtensionsClient.Delete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, accountResourceName, extensionResourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Delete", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client ExtensionsClient) DeletePreparer(ctx context.Context, resourceGroupName string, accountResourceName string, extensionResourceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountResourceName":   autorest.Encode("path", accountResourceName),
		"extensionResourceName": autorest.Encode("path", extensionResourceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2014-04-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/microsoft.visualstudio/account/{accountResourceName}/extension/{extensionResourceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client ExtensionsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ExtensionsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the details of an extension associated with a Visual Studio Team Services account resource.
// Parameters:
// resourceGroupName - name of the resource group within the Azure subscription.
// accountResourceName - the name of the Visual Studio Team Services account resource.
// extensionResourceName - the name of the extension.
func (client ExtensionsClient) Get(ctx context.Context, resourceGroupName string, accountResourceName string, extensionResourceName string) (result ExtensionResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ExtensionsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, accountResourceName, extensionResourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client ExtensionsClient) GetPreparer(ctx context.Context, resourceGroupName string, accountResourceName string, extensionResourceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountResourceName":   autorest.Encode("path", accountResourceName),
		"extensionResourceName": autorest.Encode("path", extensionResourceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2014-04-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/microsoft.visualstudio/account/{accountResourceName}/extension/{extensionResourceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ExtensionsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ExtensionsClient) GetResponder(resp *http.Response) (result ExtensionResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNotFound),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByAccount gets the details of the extension resources created within the resource group.
// Parameters:
// resourceGroupName - name of the resource group within the Azure subscription.
// accountResourceName - the name of the Visual Studio Team Services account resource.
func (client ExtensionsClient) ListByAccount(ctx context.Context, resourceGroupName string, accountResourceName string) (result ExtensionResourceListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ExtensionsClient.ListByAccount")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListByAccountPreparer(ctx, resourceGroupName, accountResourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "ListByAccount", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByAccountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "ListByAccount", resp, "Failure sending request")
		return
	}

	result, err = client.ListByAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "ListByAccount", resp, "Failure responding to request")
		return
	}

	return
}

// ListByAccountPreparer prepares the ListByAccount request.
func (client ExtensionsClient) ListByAccountPreparer(ctx context.Context, resourceGroupName string, accountResourceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountResourceName": autorest.Encode("path", accountResourceName),
		"resourceGroupName":   autorest.Encode("path", resourceGroupName),
		"subscriptionId":      autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2014-04-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/microsoft.visualstudio/account/{accountResourceName}/extension", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByAccountSender sends the ListByAccount request. The method will close the
// http.Response Body if it receives an error.
func (client ExtensionsClient) ListByAccountSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByAccountResponder handles the response to the ListByAccount request. The method always
// closes the http.Response Body.
func (client ExtensionsClient) ListByAccountResponder(resp *http.Response) (result ExtensionResourceListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Update updates an existing extension registration for the Visual Studio Team Services account.
// Parameters:
// resourceGroupName - name of the resource group within the Azure subscription.
// body - an object containing additional information related to the extension request.
// accountResourceName - the name of the Visual Studio Team Services account resource.
// extensionResourceName - the name of the extension.
func (client ExtensionsClient) Update(ctx context.Context, resourceGroupName string, body ExtensionResourceRequest, accountResourceName string, extensionResourceName string) (result ExtensionResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ExtensionsClient.Update")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdatePreparer(ctx, resourceGroupName, body, accountResourceName, extensionResourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "visualstudio.ExtensionsClient", "Update", resp, "Failure responding to request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client ExtensionsClient) UpdatePreparer(ctx context.Context, resourceGroupName string, body ExtensionResourceRequest, accountResourceName string, extensionResourceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountResourceName":   autorest.Encode("path", accountResourceName),
		"extensionResourceName": autorest.Encode("path", extensionResourceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2014-04-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/microsoft.visualstudio/account/{accountResourceName}/extension/{extensionResourceName}", pathParameters),
		autorest.WithJSON(body),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client ExtensionsClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client ExtensionsClient) UpdateResponder(resp *http.Response) (result ExtensionResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
