package storagesync

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
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// WorkflowsClient is the microsoft Storage Sync Service API
type WorkflowsClient struct {
	BaseClient
}

// NewWorkflowsClient creates an instance of the WorkflowsClient client.
func NewWorkflowsClient(subscriptionID string) WorkflowsClient {
	return NewWorkflowsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewWorkflowsClientWithBaseURI creates an instance of the WorkflowsClient client.
func NewWorkflowsClientWithBaseURI(baseURI string, subscriptionID string) WorkflowsClient {
	return WorkflowsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Abort abort the given workflow.
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// storageSyncServiceName - name of Storage Sync Service resource.
// workflowID - workflow Id
func (client WorkflowsClient) Abort(ctx context.Context, resourceGroupName string, storageSyncServiceName string, workflowID string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/WorkflowsClient.Abort")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("storagesync.WorkflowsClient", "Abort", err.Error())
	}

	req, err := client.AbortPreparer(ctx, resourceGroupName, storageSyncServiceName, workflowID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "Abort", nil, "Failure preparing request")
		return
	}

	resp, err := client.AbortSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "Abort", resp, "Failure sending request")
		return
	}

	result, err = client.AbortResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "Abort", resp, "Failure responding to request")
	}

	return
}

// AbortPreparer prepares the Abort request.
func (client WorkflowsClient) AbortPreparer(ctx context.Context, resourceGroupName string, storageSyncServiceName string, workflowID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":      autorest.Encode("path", resourceGroupName),
		"storageSyncServiceName": autorest.Encode("path", storageSyncServiceName),
		"subscriptionId":         autorest.Encode("path", client.SubscriptionID),
		"workflowId":             autorest.Encode("path", workflowID),
	}

	const APIVersion = "2019-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{storageSyncServiceName}/workflows/{workflowId}/abort", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AbortSender sends the Abort request. The method will close the
// http.Response Body if it receives an error.
func (client WorkflowsClient) AbortSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// AbortResponder handles the response to the Abort request. The method always
// closes the http.Response Body.
func (client WorkflowsClient) AbortResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get get Workflows resource
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// storageSyncServiceName - name of Storage Sync Service resource.
// workflowID - workflow Id
func (client WorkflowsClient) Get(ctx context.Context, resourceGroupName string, storageSyncServiceName string, workflowID string) (result Workflow, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/WorkflowsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("storagesync.WorkflowsClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, storageSyncServiceName, workflowID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client WorkflowsClient) GetPreparer(ctx context.Context, resourceGroupName string, storageSyncServiceName string, workflowID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":      autorest.Encode("path", resourceGroupName),
		"storageSyncServiceName": autorest.Encode("path", storageSyncServiceName),
		"subscriptionId":         autorest.Encode("path", client.SubscriptionID),
		"workflowId":             autorest.Encode("path", workflowID),
	}

	const APIVersion = "2019-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{storageSyncServiceName}/workflows/{workflowId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client WorkflowsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client WorkflowsClient) GetResponder(resp *http.Response) (result Workflow, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByStorageSyncService get a Workflow List
// Parameters:
// resourceGroupName - the name of the resource group. The name is case insensitive.
// storageSyncServiceName - name of Storage Sync Service resource.
func (client WorkflowsClient) ListByStorageSyncService(ctx context.Context, resourceGroupName string, storageSyncServiceName string) (result WorkflowArray, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/WorkflowsClient.ListByStorageSyncService")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("storagesync.WorkflowsClient", "ListByStorageSyncService", err.Error())
	}

	req, err := client.ListByStorageSyncServicePreparer(ctx, resourceGroupName, storageSyncServiceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "ListByStorageSyncService", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByStorageSyncServiceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "ListByStorageSyncService", resp, "Failure sending request")
		return
	}

	result, err = client.ListByStorageSyncServiceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagesync.WorkflowsClient", "ListByStorageSyncService", resp, "Failure responding to request")
	}

	return
}

// ListByStorageSyncServicePreparer prepares the ListByStorageSyncService request.
func (client WorkflowsClient) ListByStorageSyncServicePreparer(ctx context.Context, resourceGroupName string, storageSyncServiceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName":      autorest.Encode("path", resourceGroupName),
		"storageSyncServiceName": autorest.Encode("path", storageSyncServiceName),
		"subscriptionId":         autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{storageSyncServiceName}/workflows", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByStorageSyncServiceSender sends the ListByStorageSyncService request. The method will close the
// http.Response Body if it receives an error.
func (client WorkflowsClient) ListByStorageSyncServiceSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListByStorageSyncServiceResponder handles the response to the ListByStorageSyncService request. The method always
// closes the http.Response Body.
func (client WorkflowsClient) ListByStorageSyncServiceResponder(resp *http.Response) (result WorkflowArray, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
