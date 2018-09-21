package deploymentmanager

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
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// ServiceTopologiesClient is the REST APIs for orchestrating deployments using the Azure Deployment Manager (ADM).
type ServiceTopologiesClient struct {
	BaseClient
}

// NewServiceTopologiesClient creates an instance of the ServiceTopologiesClient client.
func NewServiceTopologiesClient(subscriptionID string) ServiceTopologiesClient {
	return NewServiceTopologiesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewServiceTopologiesClientWithBaseURI creates an instance of the ServiceTopologiesClient client.
func NewServiceTopologiesClientWithBaseURI(baseURI string, subscriptionID string) ServiceTopologiesClient {
	return ServiceTopologiesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Create this is an asynchronous operation and can be polled to completion using the operation resource returned by
// this operation.
// Parameters:
// serviceTopologyInfo - source topology object defines the resource.
// resourceGroupName - the name of the resource group. The name is case insensitive.
// serviceTopologyName - the name of the service topology .
func (client ServiceTopologiesClient) Create(ctx context.Context, serviceTopologyInfo ServiceTopologyResource, resourceGroupName string, serviceTopologyName string) (result ServiceTopologyResource, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: serviceTopologyInfo,
			Constraints: []validation.Constraint{{Target: "serviceTopologyInfo.ServiceTopologyResourceProperties", Name: validation.Null, Rule: true, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("deploymentmanager.ServiceTopologiesClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, serviceTopologyInfo, resourceGroupName, serviceTopologyName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Create", resp, "Failure responding to request")
	}

	return
}

// CreatePreparer prepares the Create request.
func (client ServiceTopologiesClient) CreatePreparer(ctx context.Context, serviceTopologyInfo ServiceTopologyResource, resourceGroupName string, serviceTopologyName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":   autorest.Encode("path", resourceGroupName),
		"serviceTopologyName": autorest.Encode("path", serviceTopologyName),
		"subscriptionId":      autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DeploymentManager/serviceTopologies/{serviceTopologyName}", pathParameters),
		autorest.WithJSON(serviceTopologyInfo),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceTopologiesClient) CreateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client ServiceTopologiesClient) CreateResponder(resp *http.Response) (result ServiceTopologyResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete sends the delete request.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// serviceTopologyName - the name of the service topology .
func (client ServiceTopologiesClient) Delete(ctx context.Context, resourceGroupName string, serviceTopologyName string) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("deploymentmanager.ServiceTopologiesClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, resourceGroupName, serviceTopologyName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client ServiceTopologiesClient) DeletePreparer(ctx context.Context, resourceGroupName string, serviceTopologyName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":   autorest.Encode("path", resourceGroupName),
		"serviceTopologyName": autorest.Encode("path", serviceTopologyName),
		"subscriptionId":      autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DeploymentManager/serviceTopologies/{serviceTopologyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceTopologiesClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ServiceTopologiesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get sends the get request.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// serviceTopologyName - the name of the service topology .
func (client ServiceTopologiesClient) Get(ctx context.Context, resourceGroupName string, serviceTopologyName string) (result ServiceTopologyResource, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("deploymentmanager.ServiceTopologiesClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, serviceTopologyName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentmanager.ServiceTopologiesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ServiceTopologiesClient) GetPreparer(ctx context.Context, resourceGroupName string, serviceTopologyName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":   autorest.Encode("path", resourceGroupName),
		"serviceTopologyName": autorest.Encode("path", serviceTopologyName),
		"subscriptionId":      autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DeploymentManager/serviceTopologies/{serviceTopologyName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceTopologiesClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ServiceTopologiesClient) GetResponder(resp *http.Response) (result ServiceTopologyResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
