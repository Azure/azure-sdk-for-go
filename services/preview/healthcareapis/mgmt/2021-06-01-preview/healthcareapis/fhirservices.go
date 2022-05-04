package healthcareapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// FhirServicesClient is the azure Healthcare APIs Client
type FhirServicesClient struct {
	BaseClient
}

// NewFhirServicesClient creates an instance of the FhirServicesClient client.
func NewFhirServicesClient(subscriptionID string) FhirServicesClient {
	return NewFhirServicesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewFhirServicesClientWithBaseURI creates an instance of the FhirServicesClient client using a custom endpoint.  Use
// this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewFhirServicesClientWithBaseURI(baseURI string, subscriptionID string) FhirServicesClient {
	return FhirServicesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates a FHIR Service resource with the specified parameters.
// Parameters:
// resourceGroupName - the name of the resource group that contains the service instance.
// workspaceName - the name of workspace resource.
// fhirServiceName - the name of FHIR Service resource.
// fhirservice - the parameters for creating or updating a Fhir Service resource.
func (client FhirServicesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, fhirServiceName string, fhirservice FhirService) (result FhirServicesCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FhirServicesClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: fhirServiceName,
			Constraints: []validation.Constraint{{Target: "fhirServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "fhirServiceName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: fhirservice,
			Constraints: []validation.Constraint{{Target: "fhirservice.FhirServiceProperties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "fhirservice.FhirServiceProperties.CorsConfiguration", Name: validation.Null, Rule: false,
					Chain: []validation.Constraint{{Target: "fhirservice.FhirServiceProperties.CorsConfiguration.MaxAge", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "fhirservice.FhirServiceProperties.CorsConfiguration.MaxAge", Name: validation.InclusiveMaximum, Rule: int64(99999), Chain: nil},
							{Target: "fhirservice.FhirServiceProperties.CorsConfiguration.MaxAge", Name: validation.InclusiveMinimum, Rule: int64(0), Chain: nil},
						}},
					}},
				}}}}}); err != nil {
		return result, validation.NewError("healthcareapis.FhirServicesClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, workspaceName, fhirServiceName, fhirservice)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client FhirServicesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, workspaceName string, fhirServiceName string, fhirservice FhirService) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fhirServiceName":   autorest.Encode("path", fhirServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/fhirservices/{fhirServiceName}", pathParameters),
		autorest.WithJSON(fhirservice),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client FhirServicesClient) CreateOrUpdateSender(req *http.Request) (future FhirServicesCreateOrUpdateFuture, err error) {
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
func (client FhirServicesClient) CreateOrUpdateResponder(resp *http.Response) (result FhirService, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a FHIR Service.
// Parameters:
// resourceGroupName - the name of the resource group that contains the service instance.
// fhirServiceName - the name of FHIR Service resource.
// workspaceName - the name of workspace resource.
func (client FhirServicesClient) Delete(ctx context.Context, resourceGroupName string, fhirServiceName string, workspaceName string) (result FhirServicesDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FhirServicesClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: fhirServiceName,
			Constraints: []validation.Constraint{{Target: "fhirServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "fhirServiceName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewError("healthcareapis.FhirServicesClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, resourceGroupName, fhirServiceName, workspaceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client FhirServicesClient) DeletePreparer(ctx context.Context, resourceGroupName string, fhirServiceName string, workspaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fhirServiceName":   autorest.Encode("path", fhirServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/fhirservices/{fhirServiceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client FhirServicesClient) DeleteSender(req *http.Request) (future FhirServicesDeleteFuture, err error) {
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
func (client FhirServicesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the properties of the specified FHIR Service.
// Parameters:
// resourceGroupName - the name of the resource group that contains the service instance.
// workspaceName - the name of workspace resource.
// fhirServiceName - the name of FHIR Service resource.
func (client FhirServicesClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, fhirServiceName string) (result FhirService, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FhirServicesClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: fhirServiceName,
			Constraints: []validation.Constraint{{Target: "fhirServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "fhirServiceName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewError("healthcareapis.FhirServicesClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, workspaceName, fhirServiceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client FhirServicesClient) GetPreparer(ctx context.Context, resourceGroupName string, workspaceName string, fhirServiceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fhirServiceName":   autorest.Encode("path", fhirServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/fhirservices/{fhirServiceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client FhirServicesClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client FhirServicesClient) GetResponder(resp *http.Response) (result FhirService, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByWorkspace lists all FHIR Services for the given workspace
// Parameters:
// resourceGroupName - the name of the resource group that contains the service instance.
// workspaceName - the name of workspace resource.
func (client FhirServicesClient) ListByWorkspace(ctx context.Context, resourceGroupName string, workspaceName string) (result FhirServiceCollectionPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FhirServicesClient.ListByWorkspace")
		defer func() {
			sc := -1
			if result.fsc.Response.Response != nil {
				sc = result.fsc.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewError("healthcareapis.FhirServicesClient", "ListByWorkspace", err.Error())
	}

	result.fn = client.listByWorkspaceNextResults
	req, err := client.ListByWorkspacePreparer(ctx, resourceGroupName, workspaceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "ListByWorkspace", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByWorkspaceSender(req)
	if err != nil {
		result.fsc.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "ListByWorkspace", resp, "Failure sending request")
		return
	}

	result.fsc, err = client.ListByWorkspaceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "ListByWorkspace", resp, "Failure responding to request")
		return
	}
	if result.fsc.hasNextLink() && result.fsc.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListByWorkspacePreparer prepares the ListByWorkspace request.
func (client FhirServicesClient) ListByWorkspacePreparer(ctx context.Context, resourceGroupName string, workspaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/fhirservices", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByWorkspaceSender sends the ListByWorkspace request. The method will close the
// http.Response Body if it receives an error.
func (client FhirServicesClient) ListByWorkspaceSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByWorkspaceResponder handles the response to the ListByWorkspace request. The method always
// closes the http.Response Body.
func (client FhirServicesClient) ListByWorkspaceResponder(resp *http.Response) (result FhirServiceCollection, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByWorkspaceNextResults retrieves the next set of results, if any.
func (client FhirServicesClient) listByWorkspaceNextResults(ctx context.Context, lastResults FhirServiceCollection) (result FhirServiceCollection, err error) {
	req, err := lastResults.fhirServiceCollectionPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "listByWorkspaceNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByWorkspaceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "listByWorkspaceNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByWorkspaceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "listByWorkspaceNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByWorkspaceComplete enumerates all values, automatically crossing page boundaries as required.
func (client FhirServicesClient) ListByWorkspaceComplete(ctx context.Context, resourceGroupName string, workspaceName string) (result FhirServiceCollectionIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FhirServicesClient.ListByWorkspace")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByWorkspace(ctx, resourceGroupName, workspaceName)
	return
}

// Update patch FHIR Service details.
// Parameters:
// resourceGroupName - the name of the resource group that contains the service instance.
// fhirServiceName - the name of FHIR Service resource.
// workspaceName - the name of workspace resource.
// fhirservicePatchResource - the parameters for updating a Fhir Service.
func (client FhirServicesClient) Update(ctx context.Context, resourceGroupName string, fhirServiceName string, workspaceName string, fhirservicePatchResource FhirServicePatchResource) (result FhirServicesUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/FhirServicesClient.Update")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: fhirServiceName,
			Constraints: []validation.Constraint{{Target: "fhirServiceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "fhirServiceName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 3, Chain: nil}}}}); err != nil {
		return result, validation.NewError("healthcareapis.FhirServicesClient", "Update", err.Error())
	}

	req, err := client.UpdatePreparer(ctx, resourceGroupName, fhirServiceName, workspaceName, fhirservicePatchResource)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthcareapis.FhirServicesClient", "Update", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client FhirServicesClient) UpdatePreparer(ctx context.Context, resourceGroupName string, fhirServiceName string, workspaceName string, fhirservicePatchResource FhirServicePatchResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fhirServiceName":   autorest.Encode("path", fhirServiceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workspaceName":     autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2021-06-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/fhirservices/{fhirServiceName}", pathParameters),
		autorest.WithJSON(fhirservicePatchResource),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client FhirServicesClient) UpdateSender(req *http.Request) (future FhirServicesUpdateFuture, err error) {
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

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client FhirServicesClient) UpdateResponder(resp *http.Response) (result FhirService, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
