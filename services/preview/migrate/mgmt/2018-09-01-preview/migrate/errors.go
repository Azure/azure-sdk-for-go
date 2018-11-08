package migrate

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
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// ErrorsClient is the migrate your workloads to Azure.
type ErrorsClient struct {
	BaseClient
}

// NewErrorsClient creates an instance of the ErrorsClient client.
func NewErrorsClient(subscriptionID string, acceptLanguage string) ErrorsClient {
	return NewErrorsClientWithBaseURI(DefaultBaseURI, subscriptionID, acceptLanguage)
}

// NewErrorsClientWithBaseURI creates an instance of the ErrorsClient client.
func NewErrorsClientWithBaseURI(baseURI string, subscriptionID string, acceptLanguage string) ErrorsClient {
	return ErrorsClient{NewWithBaseURI(baseURI, subscriptionID, acceptLanguage)}
}

// DeleteError delete the migrate error. Deleting non-existent migrate error is a no-operation.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// errorName - unique name of an error within a migrate project.
func (client ErrorsClient) DeleteError(ctx context.Context, resourceGroupName string, migrateProjectName string, errorName string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ErrorsClient.DeleteError")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeleteErrorPreparer(ctx, resourceGroupName, migrateProjectName, errorName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "DeleteError", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteErrorSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "DeleteError", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteErrorResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "DeleteError", resp, "Failure responding to request")
	}

	return
}

// DeleteErrorPreparer prepares the DeleteError request.
func (client ErrorsClient) DeleteErrorPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, errorName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"errorName":          autorest.Encode("path", errorName),
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/MigrateProjects/{migrateProjectName}/MigrateErrors/{errorName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteErrorSender sends the DeleteError request. The method will close the
// http.Response Body if it receives an error.
func (client ErrorsClient) DeleteErrorSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteErrorResponder handles the response to the DeleteError request. The method always
// closes the http.Response Body.
func (client ErrorsClient) DeleteErrorResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// EnumerateErrors sends the enumerate errors request.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// continuationToken - the continuation token.
func (client ErrorsClient) EnumerateErrors(ctx context.Context, resourceGroupName string, migrateProjectName string, continuationToken string) (result ErrorCollection, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ErrorsClient.EnumerateErrors")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.EnumerateErrorsPreparer(ctx, resourceGroupName, migrateProjectName, continuationToken)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "EnumerateErrors", nil, "Failure preparing request")
		return
	}

	resp, err := client.EnumerateErrorsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "EnumerateErrors", resp, "Failure sending request")
		return
	}

	result, err = client.EnumerateErrorsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "EnumerateErrors", resp, "Failure responding to request")
	}

	return
}

// EnumerateErrorsPreparer prepares the EnumerateErrors request.
func (client ErrorsClient) EnumerateErrorsPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, continuationToken string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(continuationToken) > 0 {
		queryParameters["continuationToken"] = autorest.Encode("query", continuationToken)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/MigrateProjects/{migrateProjectName}/MigrateErrors", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(client.AcceptLanguage) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("Accept-Language", autorest.String(client.AcceptLanguage)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// EnumerateErrorsSender sends the EnumerateErrors request. The method will close the
// http.Response Body if it receives an error.
func (client ErrorsClient) EnumerateErrorsSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// EnumerateErrorsResponder handles the response to the EnumerateErrors request. The method always
// closes the http.Response Body.
func (client ErrorsClient) EnumerateErrorsResponder(resp *http.Response) (result ErrorCollection, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetError sends the get error request.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// errorName - unique name of an error within a migrate project.
func (client ErrorsClient) GetError(ctx context.Context, resourceGroupName string, migrateProjectName string, errorName string) (result Error, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ErrorsClient.GetError")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetErrorPreparer(ctx, resourceGroupName, migrateProjectName, errorName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "GetError", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetErrorSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "GetError", resp, "Failure sending request")
		return
	}

	result, err = client.GetErrorResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.ErrorsClient", "GetError", resp, "Failure responding to request")
	}

	return
}

// GetErrorPreparer prepares the GetError request.
func (client ErrorsClient) GetErrorPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, errorName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"errorName":          autorest.Encode("path", errorName),
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/MigrateProjects/{migrateProjectName}/MigrateErrors/{errorName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetErrorSender sends the GetError request. The method will close the
// http.Response Body if it receives an error.
func (client ErrorsClient) GetErrorSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetErrorResponder handles the response to the GetError request. The method always
// closes the http.Response Body.
func (client ErrorsClient) GetErrorResponder(resp *http.Response) (result Error, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
