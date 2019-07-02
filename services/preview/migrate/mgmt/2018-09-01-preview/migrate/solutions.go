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

// SolutionsClient is the migrate your workloads to Azure.
type SolutionsClient struct {
	BaseClient
}

// NewSolutionsClient creates an instance of the SolutionsClient client.
func NewSolutionsClient(subscriptionID string, acceptLanguage string) SolutionsClient {
	return NewSolutionsClientWithBaseURI(DefaultBaseURI, subscriptionID, acceptLanguage)
}

// NewSolutionsClientWithBaseURI creates an instance of the SolutionsClient client.
func NewSolutionsClientWithBaseURI(baseURI string, subscriptionID string, acceptLanguage string) SolutionsClient {
	return SolutionsClient{NewWithBaseURI(baseURI, subscriptionID, acceptLanguage)}
}

// CleanupSolutionData sends the cleanup solution data request.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// solutionName - unique name of a migration solution within a migrate project.
func (client SolutionsClient) CleanupSolutionData(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SolutionsClient.CleanupSolutionData")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CleanupSolutionDataPreparer(ctx, resourceGroupName, migrateProjectName, solutionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "CleanupSolutionData", nil, "Failure preparing request")
		return
	}

	resp, err := client.CleanupSolutionDataSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "CleanupSolutionData", resp, "Failure sending request")
		return
	}

	result, err = client.CleanupSolutionDataResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "CleanupSolutionData", resp, "Failure responding to request")
	}

	return
}

// CleanupSolutionDataPreparer prepares the CleanupSolutionData request.
func (client SolutionsClient) CleanupSolutionDataPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"solutionName":       autorest.Encode("path", solutionName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/migrateProjects/{migrateProjectName}/solutions/{solutionName}/cleanupData", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CleanupSolutionDataSender sends the CleanupSolutionData request. The method will close the
// http.Response Body if it receives an error.
func (client SolutionsClient) CleanupSolutionDataSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CleanupSolutionDataResponder handles the response to the CleanupSolutionData request. The method always
// closes the http.Response Body.
func (client SolutionsClient) CleanupSolutionDataResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// DeleteSolution delete the solution. Deleting non-existent project is a no-operation.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// solutionName - unique name of a migration solution within a migrate project.
func (client SolutionsClient) DeleteSolution(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SolutionsClient.DeleteSolution")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeleteSolutionPreparer(ctx, resourceGroupName, migrateProjectName, solutionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "DeleteSolution", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSolutionSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "DeleteSolution", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteSolutionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "DeleteSolution", resp, "Failure responding to request")
	}

	return
}

// DeleteSolutionPreparer prepares the DeleteSolution request.
func (client SolutionsClient) DeleteSolutionPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"solutionName":       autorest.Encode("path", solutionName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/migrateProjects/{migrateProjectName}/solutions/{solutionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if len(client.AcceptLanguage) > 0 {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithHeader("Accept-Language", autorest.String(client.AcceptLanguage)))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSolutionSender sends the DeleteSolution request. The method will close the
// http.Response Body if it receives an error.
func (client SolutionsClient) DeleteSolutionSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteSolutionResponder handles the response to the DeleteSolution request. The method always
// closes the http.Response Body.
func (client SolutionsClient) DeleteSolutionResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// EnumerateSolutions sends the enumerate solutions request.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
func (client SolutionsClient) EnumerateSolutions(ctx context.Context, resourceGroupName string, migrateProjectName string) (result SolutionsCollection, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SolutionsClient.EnumerateSolutions")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.EnumerateSolutionsPreparer(ctx, resourceGroupName, migrateProjectName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "EnumerateSolutions", nil, "Failure preparing request")
		return
	}

	resp, err := client.EnumerateSolutionsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "EnumerateSolutions", resp, "Failure sending request")
		return
	}

	result, err = client.EnumerateSolutionsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "EnumerateSolutions", resp, "Failure responding to request")
	}

	return
}

// EnumerateSolutionsPreparer prepares the EnumerateSolutions request.
func (client SolutionsClient) EnumerateSolutionsPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/migrateProjects/{migrateProjectName}/solutions", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// EnumerateSolutionsSender sends the EnumerateSolutions request. The method will close the
// http.Response Body if it receives an error.
func (client SolutionsClient) EnumerateSolutionsSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// EnumerateSolutionsResponder handles the response to the EnumerateSolutions request. The method always
// closes the http.Response Body.
func (client SolutionsClient) EnumerateSolutionsResponder(resp *http.Response) (result SolutionsCollection, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetConfig sends the get config request.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// solutionName - unique name of a migration solution within a migrate project.
func (client SolutionsClient) GetConfig(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string) (result SolutionConfig, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SolutionsClient.GetConfig")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetConfigPreparer(ctx, resourceGroupName, migrateProjectName, solutionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "GetConfig", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetConfigSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "GetConfig", resp, "Failure sending request")
		return
	}

	result, err = client.GetConfigResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "GetConfig", resp, "Failure responding to request")
	}

	return
}

// GetConfigPreparer prepares the GetConfig request.
func (client SolutionsClient) GetConfigPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"solutionName":       autorest.Encode("path", solutionName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/migrateProjects/{migrateProjectName}/solutions/{solutionName}/getConfig", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetConfigSender sends the GetConfig request. The method will close the
// http.Response Body if it receives an error.
func (client SolutionsClient) GetConfigSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetConfigResponder handles the response to the GetConfig request. The method always
// closes the http.Response Body.
func (client SolutionsClient) GetConfigResponder(resp *http.Response) (result SolutionConfig, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetSolution sends the get solution request.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// solutionName - unique name of a migration solution within a migrate project.
func (client SolutionsClient) GetSolution(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string) (result Solution, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SolutionsClient.GetSolution")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetSolutionPreparer(ctx, resourceGroupName, migrateProjectName, solutionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "GetSolution", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSolutionSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "GetSolution", resp, "Failure sending request")
		return
	}

	result, err = client.GetSolutionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "GetSolution", resp, "Failure responding to request")
	}

	return
}

// GetSolutionPreparer prepares the GetSolution request.
func (client SolutionsClient) GetSolutionPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"solutionName":       autorest.Encode("path", solutionName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/migrateProjects/{migrateProjectName}/solutions/{solutionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSolutionSender sends the GetSolution request. The method will close the
// http.Response Body if it receives an error.
func (client SolutionsClient) GetSolutionSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetSolutionResponder handles the response to the GetSolution request. The method always
// closes the http.Response Body.
func (client SolutionsClient) GetSolutionResponder(resp *http.Response) (result Solution, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// PatchSolution update a solution with specified name. Supports partial updates, for example only tags can be
// provided.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// solutionName - unique name of a migration solution within a migrate project.
// solutionInput - the input for the solution.
func (client SolutionsClient) PatchSolution(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string, solutionInput Solution) (result Solution, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SolutionsClient.PatchSolution")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.PatchSolutionPreparer(ctx, resourceGroupName, migrateProjectName, solutionName, solutionInput)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "PatchSolution", nil, "Failure preparing request")
		return
	}

	resp, err := client.PatchSolutionSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "PatchSolution", resp, "Failure sending request")
		return
	}

	result, err = client.PatchSolutionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "PatchSolution", resp, "Failure responding to request")
	}

	return
}

// PatchSolutionPreparer prepares the PatchSolution request.
func (client SolutionsClient) PatchSolutionPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string, solutionInput Solution) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"solutionName":       autorest.Encode("path", solutionName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	solutionInput.ID = nil
	solutionInput.Name = nil
	solutionInput.Type = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/migrateProjects/{migrateProjectName}/solutions/{solutionName}", pathParameters),
		autorest.WithJSON(solutionInput),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PatchSolutionSender sends the PatchSolution request. The method will close the
// http.Response Body if it receives an error.
func (client SolutionsClient) PatchSolutionSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// PatchSolutionResponder handles the response to the PatchSolution request. The method always
// closes the http.Response Body.
func (client SolutionsClient) PatchSolutionResponder(resp *http.Response) (result Solution, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// PutSolution sends the put solution request.
// Parameters:
// resourceGroupName - name of the Azure Resource Group that migrate project is part of.
// migrateProjectName - name of the Azure Migrate project.
// solutionName - unique name of a migration solution within a migrate project.
// solutionInput - the input for the solution.
func (client SolutionsClient) PutSolution(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string, solutionInput Solution) (result Solution, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SolutionsClient.PutSolution")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.PutSolutionPreparer(ctx, resourceGroupName, migrateProjectName, solutionName, solutionInput)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "PutSolution", nil, "Failure preparing request")
		return
	}

	resp, err := client.PutSolutionSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "PutSolution", resp, "Failure sending request")
		return
	}

	result, err = client.PutSolutionResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrate.SolutionsClient", "PutSolution", resp, "Failure responding to request")
	}

	return
}

// PutSolutionPreparer prepares the PutSolution request.
func (client SolutionsClient) PutSolutionPreparer(ctx context.Context, resourceGroupName string, migrateProjectName string, solutionName string, solutionInput Solution) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"migrateProjectName": autorest.Encode("path", migrateProjectName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"solutionName":       autorest.Encode("path", solutionName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	solutionInput.ID = nil
	solutionInput.Name = nil
	solutionInput.Type = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/migrateProjects/{migrateProjectName}/solutions/{solutionName}", pathParameters),
		autorest.WithJSON(solutionInput),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PutSolutionSender sends the PutSolution request. The method will close the
// http.Response Body if it receives an error.
func (client SolutionsClient) PutSolutionSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// PutSolutionResponder handles the response to the PutSolution request. The method always
// closes the http.Response Body.
func (client SolutionsClient) PutSolutionResponder(resp *http.Response) (result Solution, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
