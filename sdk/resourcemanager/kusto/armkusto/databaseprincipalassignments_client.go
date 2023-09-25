//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkusto

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

// DatabasePrincipalAssignmentsClient contains the methods for the DatabasePrincipalAssignments group.
// Don't use this type directly, use NewDatabasePrincipalAssignmentsClient() instead.
type DatabasePrincipalAssignmentsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewDatabasePrincipalAssignmentsClient creates a new instance of DatabasePrincipalAssignmentsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewDatabasePrincipalAssignmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DatabasePrincipalAssignmentsClient, error) {
	cl, err := arm.NewClient(moduleName+".DatabasePrincipalAssignmentsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &DatabasePrincipalAssignmentsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CheckNameAvailability - Checks that the database principal assignment is valid and is not already in use.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-02
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the Kusto cluster.
//   - databaseName - The name of the database in the Kusto cluster.
//   - principalAssignmentName - The name of the resource.
//   - options - DatabasePrincipalAssignmentsClientCheckNameAvailabilityOptions contains the optional parameters for the DatabasePrincipalAssignmentsClient.CheckNameAvailability
//     method.
func (client *DatabasePrincipalAssignmentsClient) CheckNameAvailability(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName DatabasePrincipalAssignmentCheckNameRequest, options *DatabasePrincipalAssignmentsClientCheckNameAvailabilityOptions) (DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse, error) {
	var err error
	req, err := client.checkNameAvailabilityCreateRequest(ctx, resourceGroupName, clusterName, databaseName, principalAssignmentName, options)
	if err != nil {
		return DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.checkNameAvailabilityHandleResponse(httpResp)
	return resp, err
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *DatabasePrincipalAssignmentsClient) checkNameAvailabilityCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName DatabasePrincipalAssignmentCheckNameRequest, options *DatabasePrincipalAssignmentsClientCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/checkPrincipalAssignmentNameAvailability"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-02")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, principalAssignmentName); err != nil {
	return nil, err
}
	return req, nil
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *DatabasePrincipalAssignmentsClient) checkNameAvailabilityHandleResponse(resp *http.Response) (DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse, error) {
	result := DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckNameResult); err != nil {
		return DatabasePrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	return result, nil
}

// BeginCreateOrUpdate - Creates a Kusto cluster database principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-02
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the Kusto cluster.
//   - databaseName - The name of the database in the Kusto cluster.
//   - principalAssignmentName - The name of the Kusto principalAssignment.
//   - parameters - The Kusto principalAssignments parameters supplied for the operation.
//   - options - DatabasePrincipalAssignmentsClientBeginCreateOrUpdateOptions contains the optional parameters for the DatabasePrincipalAssignmentsClient.BeginCreateOrUpdate
//     method.
func (client *DatabasePrincipalAssignmentsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string, parameters DatabasePrincipalAssignment, options *DatabasePrincipalAssignmentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[DatabasePrincipalAssignmentsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, clusterName, databaseName, principalAssignmentName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[DatabasePrincipalAssignmentsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[DatabasePrincipalAssignmentsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates a Kusto cluster database principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-02
func (client *DatabasePrincipalAssignmentsClient) createOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string, parameters DatabasePrincipalAssignment, options *DatabasePrincipalAssignmentsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, clusterName, databaseName, principalAssignmentName, parameters, options)
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
func (client *DatabasePrincipalAssignmentsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string, parameters DatabasePrincipalAssignment, options *DatabasePrincipalAssignmentsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/principalAssignments/{principalAssignmentName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-02")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes a Kusto principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-02
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the Kusto cluster.
//   - databaseName - The name of the database in the Kusto cluster.
//   - principalAssignmentName - The name of the Kusto principalAssignment.
//   - options - DatabasePrincipalAssignmentsClientBeginDeleteOptions contains the optional parameters for the DatabasePrincipalAssignmentsClient.BeginDelete
//     method.
func (client *DatabasePrincipalAssignmentsClient) BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string, options *DatabasePrincipalAssignmentsClientBeginDeleteOptions) (*runtime.Poller[DatabasePrincipalAssignmentsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, clusterName, databaseName, principalAssignmentName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[DatabasePrincipalAssignmentsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[DatabasePrincipalAssignmentsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes a Kusto principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-02
func (client *DatabasePrincipalAssignmentsClient) deleteOperation(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string, options *DatabasePrincipalAssignmentsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, clusterName, databaseName, principalAssignmentName, options)
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
func (client *DatabasePrincipalAssignmentsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string, options *DatabasePrincipalAssignmentsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/principalAssignments/{principalAssignmentName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-02")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets a Kusto cluster database principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-02
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the Kusto cluster.
//   - databaseName - The name of the database in the Kusto cluster.
//   - principalAssignmentName - The name of the Kusto principalAssignment.
//   - options - DatabasePrincipalAssignmentsClientGetOptions contains the optional parameters for the DatabasePrincipalAssignmentsClient.Get
//     method.
func (client *DatabasePrincipalAssignmentsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string, options *DatabasePrincipalAssignmentsClientGetOptions) (DatabasePrincipalAssignmentsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, clusterName, databaseName, principalAssignmentName, options)
	if err != nil {
		return DatabasePrincipalAssignmentsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DatabasePrincipalAssignmentsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DatabasePrincipalAssignmentsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *DatabasePrincipalAssignmentsClient) getCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string, options *DatabasePrincipalAssignmentsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/principalAssignments/{principalAssignmentName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-02")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DatabasePrincipalAssignmentsClient) getHandleResponse(resp *http.Response) (DatabasePrincipalAssignmentsClientGetResponse, error) {
	result := DatabasePrincipalAssignmentsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DatabasePrincipalAssignment); err != nil {
		return DatabasePrincipalAssignmentsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all Kusto cluster database principalAssignments.
//
// Generated from API version 2023-05-02
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the Kusto cluster.
//   - databaseName - The name of the database in the Kusto cluster.
//   - options - DatabasePrincipalAssignmentsClientListOptions contains the optional parameters for the DatabasePrincipalAssignmentsClient.NewListPager
//     method.
func (client *DatabasePrincipalAssignmentsClient) NewListPager(resourceGroupName string, clusterName string, databaseName string, options *DatabasePrincipalAssignmentsClientListOptions) (*runtime.Pager[DatabasePrincipalAssignmentsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[DatabasePrincipalAssignmentsClientListResponse]{
		More: func(page DatabasePrincipalAssignmentsClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *DatabasePrincipalAssignmentsClientListResponse) (DatabasePrincipalAssignmentsClientListResponse, error) {
			req, err := client.listCreateRequest(ctx, resourceGroupName, clusterName, databaseName, options)
			if err != nil {
				return DatabasePrincipalAssignmentsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return DatabasePrincipalAssignmentsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DatabasePrincipalAssignmentsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *DatabasePrincipalAssignmentsClient) listCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, databaseName string, options *DatabasePrincipalAssignmentsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/databases/{databaseName}/principalAssignments"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-02")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *DatabasePrincipalAssignmentsClient) listHandleResponse(resp *http.Response) (DatabasePrincipalAssignmentsClientListResponse, error) {
	result := DatabasePrincipalAssignmentsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DatabasePrincipalAssignmentListResult); err != nil {
		return DatabasePrincipalAssignmentsClientListResponse{}, err
	}
	return result, nil
}

