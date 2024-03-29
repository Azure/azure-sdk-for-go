//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsynapse

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// KustoPoolDatabasePrincipalAssignmentsClient contains the methods for the KustoPoolDatabasePrincipalAssignments group.
// Don't use this type directly, use NewKustoPoolDatabasePrincipalAssignmentsClient() instead.
type KustoPoolDatabasePrincipalAssignmentsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewKustoPoolDatabasePrincipalAssignmentsClient creates a new instance of KustoPoolDatabasePrincipalAssignmentsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewKustoPoolDatabasePrincipalAssignmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*KustoPoolDatabasePrincipalAssignmentsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &KustoPoolDatabasePrincipalAssignmentsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CheckNameAvailability - Checks that the database principal assignment is valid and is not already in use.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - databaseName - The name of the database in the Kusto pool.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - principalAssignmentName - The name of the resource.
//   - options - KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityOptions contains the optional parameters for
//     the KustoPoolDatabasePrincipalAssignmentsClient.CheckNameAvailability method.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) CheckNameAvailability(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, resourceGroupName string, principalAssignmentName DatabasePrincipalAssignmentCheckNameRequest, options *KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityOptions) (KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse, error) {
	var err error
	const operationName = "KustoPoolDatabasePrincipalAssignmentsClient.CheckNameAvailability"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.checkNameAvailabilityCreateRequest(ctx, workspaceName, kustoPoolName, databaseName, resourceGroupName, principalAssignmentName, options)
	if err != nil {
		return KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.checkNameAvailabilityHandleResponse(httpResp)
	return resp, err
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) checkNameAvailabilityCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, resourceGroupName string, principalAssignmentName DatabasePrincipalAssignmentCheckNameRequest, options *KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/databases/{databaseName}/checkPrincipalAssignmentNameAvailability"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, principalAssignmentName); err != nil {
		return nil, err
	}
	return req, nil
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) checkNameAvailabilityHandleResponse(resp *http.Response) (KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse, error) {
	result := KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckNameResult); err != nil {
		return KustoPoolDatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	return result, nil
}

// BeginCreateOrUpdate - Creates a Kusto pool database principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - databaseName - The name of the database in the Kusto pool.
//   - principalAssignmentName - The name of the Kusto principalAssignment.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - parameters - The Kusto principalAssignments parameters supplied for the operation.
//   - options - KustoPoolDatabasePrincipalAssignmentsClientBeginCreateOrUpdateOptions contains the optional parameters for the
//     KustoPoolDatabasePrincipalAssignmentsClient.BeginCreateOrUpdate method.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) BeginCreateOrUpdate(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, principalAssignmentName string, resourceGroupName string, parameters DatabasePrincipalAssignment, options *KustoPoolDatabasePrincipalAssignmentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[KustoPoolDatabasePrincipalAssignmentsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, workspaceName, kustoPoolName, databaseName, principalAssignmentName, resourceGroupName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[KustoPoolDatabasePrincipalAssignmentsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[KustoPoolDatabasePrincipalAssignmentsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Creates a Kusto pool database principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
func (client *KustoPoolDatabasePrincipalAssignmentsClient) createOrUpdate(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, principalAssignmentName string, resourceGroupName string, parameters DatabasePrincipalAssignment, options *KustoPoolDatabasePrincipalAssignmentsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "KustoPoolDatabasePrincipalAssignmentsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, workspaceName, kustoPoolName, databaseName, principalAssignmentName, resourceGroupName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) createOrUpdateCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, principalAssignmentName string, resourceGroupName string, parameters DatabasePrincipalAssignment, options *KustoPoolDatabasePrincipalAssignmentsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/databases/{databaseName}/principalAssignments/{principalAssignmentName}"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes a Kusto pool principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - databaseName - The name of the database in the Kusto pool.
//   - principalAssignmentName - The name of the Kusto principalAssignment.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - KustoPoolDatabasePrincipalAssignmentsClientBeginDeleteOptions contains the optional parameters for the KustoPoolDatabasePrincipalAssignmentsClient.BeginDelete
//     method.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) BeginDelete(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, principalAssignmentName string, resourceGroupName string, options *KustoPoolDatabasePrincipalAssignmentsClientBeginDeleteOptions) (*runtime.Poller[KustoPoolDatabasePrincipalAssignmentsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, workspaceName, kustoPoolName, databaseName, principalAssignmentName, resourceGroupName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[KustoPoolDatabasePrincipalAssignmentsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[KustoPoolDatabasePrincipalAssignmentsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes a Kusto pool principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
func (client *KustoPoolDatabasePrincipalAssignmentsClient) deleteOperation(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, principalAssignmentName string, resourceGroupName string, options *KustoPoolDatabasePrincipalAssignmentsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "KustoPoolDatabasePrincipalAssignmentsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, workspaceName, kustoPoolName, databaseName, principalAssignmentName, resourceGroupName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) deleteCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, principalAssignmentName string, resourceGroupName string, options *KustoPoolDatabasePrincipalAssignmentsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/databases/{databaseName}/principalAssignments/{principalAssignmentName}"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets a Kusto pool database principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - databaseName - The name of the database in the Kusto pool.
//   - principalAssignmentName - The name of the Kusto principalAssignment.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - KustoPoolDatabasePrincipalAssignmentsClientGetOptions contains the optional parameters for the KustoPoolDatabasePrincipalAssignmentsClient.Get
//     method.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) Get(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, principalAssignmentName string, resourceGroupName string, options *KustoPoolDatabasePrincipalAssignmentsClientGetOptions) (KustoPoolDatabasePrincipalAssignmentsClientGetResponse, error) {
	var err error
	const operationName = "KustoPoolDatabasePrincipalAssignmentsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, workspaceName, kustoPoolName, databaseName, principalAssignmentName, resourceGroupName, options)
	if err != nil {
		return KustoPoolDatabasePrincipalAssignmentsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return KustoPoolDatabasePrincipalAssignmentsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return KustoPoolDatabasePrincipalAssignmentsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) getCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, principalAssignmentName string, resourceGroupName string, options *KustoPoolDatabasePrincipalAssignmentsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/databases/{databaseName}/principalAssignments/{principalAssignmentName}"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) getHandleResponse(resp *http.Response) (KustoPoolDatabasePrincipalAssignmentsClientGetResponse, error) {
	result := KustoPoolDatabasePrincipalAssignmentsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DatabasePrincipalAssignment); err != nil {
		return KustoPoolDatabasePrincipalAssignmentsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all Kusto pool database principalAssignments.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - databaseName - The name of the database in the Kusto pool.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - KustoPoolDatabasePrincipalAssignmentsClientListOptions contains the optional parameters for the KustoPoolDatabasePrincipalAssignmentsClient.NewListPager
//     method.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) NewListPager(workspaceName string, kustoPoolName string, databaseName string, resourceGroupName string, options *KustoPoolDatabasePrincipalAssignmentsClientListOptions) *runtime.Pager[KustoPoolDatabasePrincipalAssignmentsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[KustoPoolDatabasePrincipalAssignmentsClientListResponse]{
		More: func(page KustoPoolDatabasePrincipalAssignmentsClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *KustoPoolDatabasePrincipalAssignmentsClientListResponse) (KustoPoolDatabasePrincipalAssignmentsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "KustoPoolDatabasePrincipalAssignmentsClient.NewListPager")
			req, err := client.listCreateRequest(ctx, workspaceName, kustoPoolName, databaseName, resourceGroupName, options)
			if err != nil {
				return KustoPoolDatabasePrincipalAssignmentsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return KustoPoolDatabasePrincipalAssignmentsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return KustoPoolDatabasePrincipalAssignmentsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) listCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, databaseName string, resourceGroupName string, options *KustoPoolDatabasePrincipalAssignmentsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/databases/{databaseName}/principalAssignments"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *KustoPoolDatabasePrincipalAssignmentsClient) listHandleResponse(resp *http.Response) (KustoPoolDatabasePrincipalAssignmentsClientListResponse, error) {
	result := KustoPoolDatabasePrincipalAssignmentsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DatabasePrincipalAssignmentListResult); err != nil {
		return KustoPoolDatabasePrincipalAssignmentsClientListResponse{}, err
	}
	return result, nil
}
