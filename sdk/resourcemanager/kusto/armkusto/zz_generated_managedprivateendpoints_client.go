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
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ManagedPrivateEndpointsClient contains the methods for the ManagedPrivateEndpoints group.
// Don't use this type directly, use NewManagedPrivateEndpointsClient() instead.
type ManagedPrivateEndpointsClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewManagedPrivateEndpointsClient creates a new instance of ManagedPrivateEndpointsClient with the specified values.
func NewManagedPrivateEndpointsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *ManagedPrivateEndpointsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &ManagedPrivateEndpointsClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// CheckNameAvailability - Checks that the managed private endpoints resource name is valid and is not already in use.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) CheckNameAvailability(ctx context.Context, resourceGroupName string, clusterName string, resourceName ManagedPrivateEndpointsCheckNameRequest, options *ManagedPrivateEndpointsCheckNameAvailabilityOptions) (ManagedPrivateEndpointsCheckNameAvailabilityResponse, error) {
	req, err := client.checkNameAvailabilityCreateRequest(ctx, resourceGroupName, clusterName, resourceName, options)
	if err != nil {
		return ManagedPrivateEndpointsCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ManagedPrivateEndpointsCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ManagedPrivateEndpointsCheckNameAvailabilityResponse{}, client.checkNameAvailabilityHandleError(resp)
	}
	return client.checkNameAvailabilityHandleResponse(resp)
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *ManagedPrivateEndpointsClient) checkNameAvailabilityCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, resourceName ManagedPrivateEndpointsCheckNameRequest, options *ManagedPrivateEndpointsCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/managedPrivateEndpointsCheckNameAvailability"
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
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, resourceName)
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *ManagedPrivateEndpointsClient) checkNameAvailabilityHandleResponse(resp *http.Response) (ManagedPrivateEndpointsCheckNameAvailabilityResponse, error) {
	result := ManagedPrivateEndpointsCheckNameAvailabilityResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckNameResult); err != nil {
		return ManagedPrivateEndpointsCheckNameAvailabilityResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// checkNameAvailabilityHandleError handles the CheckNameAvailability error response.
func (client *ManagedPrivateEndpointsClient) checkNameAvailabilityHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginCreateOrUpdate - Creates a managed private endpoint.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, parameters ManagedPrivateEndpoint, options *ManagedPrivateEndpointsBeginCreateOrUpdateOptions) (ManagedPrivateEndpointsCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, clusterName, managedPrivateEndpointName, parameters, options)
	if err != nil {
		return ManagedPrivateEndpointsCreateOrUpdatePollerResponse{}, err
	}
	result := ManagedPrivateEndpointsCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("ManagedPrivateEndpointsClient.CreateOrUpdate", "", resp, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return ManagedPrivateEndpointsCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &ManagedPrivateEndpointsCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Creates a managed private endpoint.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) createOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, parameters ManagedPrivateEndpoint, options *ManagedPrivateEndpointsBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, clusterName, managedPrivateEndpointName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ManagedPrivateEndpointsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, parameters ManagedPrivateEndpoint, options *ManagedPrivateEndpointsBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/managedPrivateEndpoints/{managedPrivateEndpointName}"
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
	if managedPrivateEndpointName == "" {
		return nil, errors.New("parameter managedPrivateEndpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedPrivateEndpointName}", url.PathEscape(managedPrivateEndpointName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *ManagedPrivateEndpointsClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginDelete - Deletes a managed private endpoint.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, options *ManagedPrivateEndpointsBeginDeleteOptions) (ManagedPrivateEndpointsDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, clusterName, managedPrivateEndpointName, options)
	if err != nil {
		return ManagedPrivateEndpointsDeletePollerResponse{}, err
	}
	result := ManagedPrivateEndpointsDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("ManagedPrivateEndpointsClient.Delete", "", resp, client.pl, client.deleteHandleError)
	if err != nil {
		return ManagedPrivateEndpointsDeletePollerResponse{}, err
	}
	result.Poller = &ManagedPrivateEndpointsDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Deletes a managed private endpoint.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) deleteOperation(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, options *ManagedPrivateEndpointsBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, clusterName, managedPrivateEndpointName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ManagedPrivateEndpointsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, options *ManagedPrivateEndpointsBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/managedPrivateEndpoints/{managedPrivateEndpointName}"
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
	if managedPrivateEndpointName == "" {
		return nil, errors.New("parameter managedPrivateEndpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedPrivateEndpointName}", url.PathEscape(managedPrivateEndpointName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *ManagedPrivateEndpointsClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Gets a managed private endpoint.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, options *ManagedPrivateEndpointsGetOptions) (ManagedPrivateEndpointsGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, clusterName, managedPrivateEndpointName, options)
	if err != nil {
		return ManagedPrivateEndpointsGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ManagedPrivateEndpointsGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ManagedPrivateEndpointsGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ManagedPrivateEndpointsClient) getCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, options *ManagedPrivateEndpointsGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/managedPrivateEndpoints/{managedPrivateEndpointName}"
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
	if managedPrivateEndpointName == "" {
		return nil, errors.New("parameter managedPrivateEndpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedPrivateEndpointName}", url.PathEscape(managedPrivateEndpointName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
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
func (client *ManagedPrivateEndpointsClient) getHandleResponse(resp *http.Response) (ManagedPrivateEndpointsGetResponse, error) {
	result := ManagedPrivateEndpointsGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedPrivateEndpoint); err != nil {
		return ManagedPrivateEndpointsGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *ManagedPrivateEndpointsClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// List - Returns the list of managed private endpoints.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) List(ctx context.Context, resourceGroupName string, clusterName string, options *ManagedPrivateEndpointsListOptions) (ManagedPrivateEndpointsListResponse, error) {
	req, err := client.listCreateRequest(ctx, resourceGroupName, clusterName, options)
	if err != nil {
		return ManagedPrivateEndpointsListResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ManagedPrivateEndpointsListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ManagedPrivateEndpointsListResponse{}, client.listHandleError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *ManagedPrivateEndpointsClient) listCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, options *ManagedPrivateEndpointsListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/managedPrivateEndpoints"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
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
func (client *ManagedPrivateEndpointsClient) listHandleResponse(resp *http.Response) (ManagedPrivateEndpointsListResponse, error) {
	result := ManagedPrivateEndpointsListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedPrivateEndpointListResult); err != nil {
		return ManagedPrivateEndpointsListResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *ManagedPrivateEndpointsClient) listHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginUpdate - Updates a managed private endpoint.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) BeginUpdate(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, parameters ManagedPrivateEndpoint, options *ManagedPrivateEndpointsBeginUpdateOptions) (ManagedPrivateEndpointsUpdatePollerResponse, error) {
	resp, err := client.update(ctx, resourceGroupName, clusterName, managedPrivateEndpointName, parameters, options)
	if err != nil {
		return ManagedPrivateEndpointsUpdatePollerResponse{}, err
	}
	result := ManagedPrivateEndpointsUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("ManagedPrivateEndpointsClient.Update", "", resp, client.pl, client.updateHandleError)
	if err != nil {
		return ManagedPrivateEndpointsUpdatePollerResponse{}, err
	}
	result.Poller = &ManagedPrivateEndpointsUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// Update - Updates a managed private endpoint.
// If the operation fails it returns the *CloudError error type.
func (client *ManagedPrivateEndpointsClient) update(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, parameters ManagedPrivateEndpoint, options *ManagedPrivateEndpointsBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, clusterName, managedPrivateEndpointName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.updateHandleError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *ManagedPrivateEndpointsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, managedPrivateEndpointName string, parameters ManagedPrivateEndpoint, options *ManagedPrivateEndpointsBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{clusterName}/managedPrivateEndpoints/{managedPrivateEndpointName}"
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
	if managedPrivateEndpointName == "" {
		return nil, errors.New("parameter managedPrivateEndpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedPrivateEndpointName}", url.PathEscape(managedPrivateEndpointName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleError handles the Update error response.
func (client *ManagedPrivateEndpointsClient) updateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
