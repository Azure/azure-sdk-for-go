package appplatform

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

// BuildpackBindingClient is the REST API for Azure Spring Cloud
type BuildpackBindingClient struct {
	BaseClient
}

// NewBuildpackBindingClient creates an instance of the BuildpackBindingClient client.
func NewBuildpackBindingClient(subscriptionID string) BuildpackBindingClient {
	return NewBuildpackBindingClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewBuildpackBindingClientWithBaseURI creates an instance of the BuildpackBindingClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewBuildpackBindingClientWithBaseURI(baseURI string, subscriptionID string) BuildpackBindingClient {
	return BuildpackBindingClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate create or update a buildpack binding.
// Parameters:
// resourceGroupName - the name of the resource group that contains the resource. You can obtain this value
// from the Azure Resource Manager API or the portal.
// serviceName - the name of the Service resource.
// buildServiceName - the name of the build service resource.
// builderName - the name of the builder resource.
// buildpackBindingName - the name of the Buildpack Binding Name
// buildpackBinding - the target buildpack binding for the create or update operation
func (client BuildpackBindingClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string, buildpackBindingName string, buildpackBinding BuildpackBindingResource) (result BuildpackBindingCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BuildpackBindingClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, serviceName, buildServiceName, builderName, buildpackBindingName, buildpackBinding)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client BuildpackBindingClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string, buildpackBindingName string, buildpackBinding BuildpackBindingResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"builderName":          autorest.Encode("path", builderName),
		"buildpackBindingName": autorest.Encode("path", buildpackBindingName),
		"buildServiceName":     autorest.Encode("path", buildServiceName),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"serviceName":          autorest.Encode("path", serviceName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2022-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/buildServices/{buildServiceName}/builders/{builderName}/buildpackBindings/{buildpackBindingName}", pathParameters),
		autorest.WithJSON(buildpackBinding),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client BuildpackBindingClient) CreateOrUpdateSender(req *http.Request) (future BuildpackBindingCreateOrUpdateFuture, err error) {
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

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client BuildpackBindingClient) CreateOrUpdateResponder(resp *http.Response) (result BuildpackBindingResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete operation to delete a Buildpack Binding
// Parameters:
// resourceGroupName - the name of the resource group that contains the resource. You can obtain this value
// from the Azure Resource Manager API or the portal.
// serviceName - the name of the Service resource.
// buildServiceName - the name of the build service resource.
// builderName - the name of the builder resource.
// buildpackBindingName - the name of the Buildpack Binding Name
func (client BuildpackBindingClient) Delete(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string, buildpackBindingName string) (result BuildpackBindingDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BuildpackBindingClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, serviceName, buildServiceName, builderName, buildpackBindingName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client BuildpackBindingClient) DeletePreparer(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string, buildpackBindingName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"builderName":          autorest.Encode("path", builderName),
		"buildpackBindingName": autorest.Encode("path", buildpackBindingName),
		"buildServiceName":     autorest.Encode("path", buildServiceName),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"serviceName":          autorest.Encode("path", serviceName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2022-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/buildServices/{buildServiceName}/builders/{builderName}/buildpackBindings/{buildpackBindingName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client BuildpackBindingClient) DeleteSender(req *http.Request) (future BuildpackBindingDeleteFuture, err error) {
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

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client BuildpackBindingClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get get a buildpack binding by name.
// Parameters:
// resourceGroupName - the name of the resource group that contains the resource. You can obtain this value
// from the Azure Resource Manager API or the portal.
// serviceName - the name of the Service resource.
// buildServiceName - the name of the build service resource.
// builderName - the name of the builder resource.
// buildpackBindingName - the name of the Buildpack Binding Name
func (client BuildpackBindingClient) Get(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string, buildpackBindingName string) (result BuildpackBindingResource, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BuildpackBindingClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, serviceName, buildServiceName, builderName, buildpackBindingName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client BuildpackBindingClient) GetPreparer(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string, buildpackBindingName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"builderName":          autorest.Encode("path", builderName),
		"buildpackBindingName": autorest.Encode("path", buildpackBindingName),
		"buildServiceName":     autorest.Encode("path", buildServiceName),
		"resourceGroupName":    autorest.Encode("path", resourceGroupName),
		"serviceName":          autorest.Encode("path", serviceName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2022-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/buildServices/{buildServiceName}/builders/{builderName}/buildpackBindings/{buildpackBindingName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client BuildpackBindingClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client BuildpackBindingClient) GetResponder(resp *http.Response) (result BuildpackBindingResource, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List handles requests to list all buildpack bindings in a builder.
// Parameters:
// resourceGroupName - the name of the resource group that contains the resource. You can obtain this value
// from the Azure Resource Manager API or the portal.
// serviceName - the name of the Service resource.
// buildServiceName - the name of the build service resource.
// builderName - the name of the builder resource.
func (client BuildpackBindingClient) List(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string) (result BuildpackBindingResourceCollectionPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BuildpackBindingClient.List")
		defer func() {
			sc := -1
			if result.bbrc.Response.Response != nil {
				sc = result.bbrc.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, serviceName, buildServiceName, builderName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.bbrc.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "List", resp, "Failure sending request")
		return
	}

	result.bbrc, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "List", resp, "Failure responding to request")
		return
	}
	if result.bbrc.hasNextLink() && result.bbrc.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client BuildpackBindingClient) ListPreparer(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"builderName":       autorest.Encode("path", builderName),
		"buildServiceName":  autorest.Encode("path", buildServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serviceName":       autorest.Encode("path", serviceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2022-03-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/buildServices/{buildServiceName}/builders/{builderName}/buildpackBindings", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client BuildpackBindingClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client BuildpackBindingClient) ListResponder(resp *http.Response) (result BuildpackBindingResourceCollection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client BuildpackBindingClient) listNextResults(ctx context.Context, lastResults BuildpackBindingResourceCollection) (result BuildpackBindingResourceCollection, err error) {
	req, err := lastResults.buildpackBindingResourceCollectionPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "appplatform.BuildpackBindingClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client BuildpackBindingClient) ListComplete(ctx context.Context, resourceGroupName string, serviceName string, buildServiceName string, builderName string) (result BuildpackBindingResourceCollectionIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BuildpackBindingClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, resourceGroupName, serviceName, buildServiceName, builderName)
	return
}
