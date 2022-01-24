//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkusto

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ClusterPrincipalAssignmentsClient contains the methods for the ClusterPrincipalAssignments group.
// Don't use this type directly, use NewClusterPrincipalAssignmentsClient() instead.
type ClusterPrincipalAssignmentsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewClusterPrincipalAssignmentsClient creates a new instance of ClusterPrincipalAssignmentsClient with the specified values.
// subscriptionID - Gets subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID
// forms part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewClusterPrincipalAssignmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *ClusterPrincipalAssignmentsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &ClusterPrincipalAssignmentsClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// CheckNameAvailability - Checks that the principal assignment name is valid and is not already in use.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// principalAssignmentName - The name of the principal assignment.
// options - ClusterPrincipalAssignmentsClientCheckNameAvailabilityOptions contains the optional parameters for the ClusterPrincipalAssignmentsClient.CheckNameAvailability
// method.
func (client *ClusterPrincipalAssignmentsClient) CheckNameAvailability(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName ClusterPrincipalAssignmentCheckNameRequest, options *ClusterPrincipalAssignmentsClientCheckNameAvailabilityOptions) (ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse, error) {
	req, err := client.checkNameAvailabilityCreateRequest(ctx, resourceGroupName, clusterName, principalAssignmentName, options)
	if err != nil {
		return ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse{}, runtime.NewResponseError(resp)
	}
	return client.checkNameAvailabilityHandleResponse(resp)
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *ClusterPrincipalAssignmentsClient) checkNameAvailabilityCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName ClusterPrincipalAssignmentCheckNameRequest, options *ClusterPrincipalAssignmentsClientCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/checkPrincipalAssignmentNameAvailability"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, principalAssignmentName)
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *ClusterPrincipalAssignmentsClient) checkNameAvailabilityHandleResponse(resp *http.Response) (ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse, error) {
	result := ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckNameResult); err != nil {
		return ClusterPrincipalAssignmentsClientCheckNameAvailabilityResponse{}, err
	}
	return result, nil
}

// BeginCreateOrUpdate - Create a Kusto cluster principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// principalAssignmentName - The name of the Kusto principalAssignment.
// parameters - The Kusto cluster principalAssignment's parameters supplied for the operation.
// options - ClusterPrincipalAssignmentsClientBeginCreateOrUpdateOptions contains the optional parameters for the ClusterPrincipalAssignmentsClient.BeginCreateOrUpdate
// method.
func (client *ClusterPrincipalAssignmentsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName string, parameters ClusterPrincipalAssignment, options *ClusterPrincipalAssignmentsClientBeginCreateOrUpdateOptions) (ClusterPrincipalAssignmentsClientCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, clusterName, principalAssignmentName, parameters, options)
	if err != nil {
		return ClusterPrincipalAssignmentsClientCreateOrUpdatePollerResponse{}, err
	}
	result := ClusterPrincipalAssignmentsClientCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("ClusterPrincipalAssignmentsClient.CreateOrUpdate", "", resp, client.pl)
	if err != nil {
		return ClusterPrincipalAssignmentsClientCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &ClusterPrincipalAssignmentsClientCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Create a Kusto cluster principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *ClusterPrincipalAssignmentsClient) createOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName string, parameters ClusterPrincipalAssignment, options *ClusterPrincipalAssignmentsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, clusterName, principalAssignmentName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ClusterPrincipalAssignmentsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName string, parameters ClusterPrincipalAssignment, options *ClusterPrincipalAssignmentsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/principalAssignments/{principalAssignmentName}"
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
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Deletes a Kusto cluster principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// principalAssignmentName - The name of the Kusto principalAssignment.
// options - ClusterPrincipalAssignmentsClientBeginDeleteOptions contains the optional parameters for the ClusterPrincipalAssignmentsClient.BeginDelete
// method.
func (client *ClusterPrincipalAssignmentsClient) BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName string, options *ClusterPrincipalAssignmentsClientBeginDeleteOptions) (ClusterPrincipalAssignmentsClientDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, clusterName, principalAssignmentName, options)
	if err != nil {
		return ClusterPrincipalAssignmentsClientDeletePollerResponse{}, err
	}
	result := ClusterPrincipalAssignmentsClientDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("ClusterPrincipalAssignmentsClient.Delete", "", resp, client.pl)
	if err != nil {
		return ClusterPrincipalAssignmentsClientDeletePollerResponse{}, err
	}
	result.Poller = &ClusterPrincipalAssignmentsClientDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Deletes a Kusto cluster principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *ClusterPrincipalAssignmentsClient) deleteOperation(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName string, options *ClusterPrincipalAssignmentsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, clusterName, principalAssignmentName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ClusterPrincipalAssignmentsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName string, options *ClusterPrincipalAssignmentsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/principalAssignments/{principalAssignmentName}"
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
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets a Kusto cluster principalAssignment.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// principalAssignmentName - The name of the Kusto principalAssignment.
// options - ClusterPrincipalAssignmentsClientGetOptions contains the optional parameters for the ClusterPrincipalAssignmentsClient.Get
// method.
func (client *ClusterPrincipalAssignmentsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName string, options *ClusterPrincipalAssignmentsClientGetOptions) (ClusterPrincipalAssignmentsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, clusterName, principalAssignmentName, options)
	if err != nil {
		return ClusterPrincipalAssignmentsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ClusterPrincipalAssignmentsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ClusterPrincipalAssignmentsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ClusterPrincipalAssignmentsClient) getCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, principalAssignmentName string, options *ClusterPrincipalAssignmentsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/principalAssignments/{principalAssignmentName}"
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
	if principalAssignmentName == "" {
		return nil, errors.New("parameter principalAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{principalAssignmentName}", url.PathEscape(principalAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ClusterPrincipalAssignmentsClient) getHandleResponse(resp *http.Response) (ClusterPrincipalAssignmentsClientGetResponse, error) {
	result := ClusterPrincipalAssignmentsClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.ClusterPrincipalAssignment); err != nil {
		return ClusterPrincipalAssignmentsClientGetResponse{}, err
	}
	return result, nil
}

// List - Lists all Kusto cluster principalAssignments.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group containing the Kusto cluster.
// clusterName - The name of the Kusto cluster.
// options - ClusterPrincipalAssignmentsClientListOptions contains the optional parameters for the ClusterPrincipalAssignmentsClient.List
// method.
func (client *ClusterPrincipalAssignmentsClient) List(ctx context.Context, resourceGroupName string, clusterName string, options *ClusterPrincipalAssignmentsClientListOptions) (ClusterPrincipalAssignmentsClientListResponse, error) {
	req, err := client.listCreateRequest(ctx, resourceGroupName, clusterName, options)
	if err != nil {
		return ClusterPrincipalAssignmentsClientListResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ClusterPrincipalAssignmentsClientListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ClusterPrincipalAssignmentsClientListResponse{}, runtime.NewResponseError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *ClusterPrincipalAssignmentsClient) listCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, options *ClusterPrincipalAssignmentsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/principalAssignments"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ClusterPrincipalAssignmentsClient) listHandleResponse(resp *http.Response) (ClusterPrincipalAssignmentsClientListResponse, error) {
	result := ClusterPrincipalAssignmentsClientListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.ClusterPrincipalAssignmentListResult); err != nil {
		return ClusterPrincipalAssignmentsClientListResponse{}, err
	}
	return result, nil
}
