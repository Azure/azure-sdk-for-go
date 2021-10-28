//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

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

// VirtualApplianceSitesClient contains the methods for the VirtualApplianceSites group.
// Don't use this type directly, use NewVirtualApplianceSitesClient() instead.
type VirtualApplianceSitesClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewVirtualApplianceSitesClient creates a new instance of VirtualApplianceSitesClient with the specified values.
func NewVirtualApplianceSitesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *VirtualApplianceSitesClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &VirtualApplianceSitesClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// BeginCreateOrUpdate - Creates or updates the specified Network Virtual Appliance Site.
// If the operation fails it returns the *CloudError error type.
func (client *VirtualApplianceSitesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, siteName string, parameters VirtualApplianceSite, options *VirtualApplianceSitesBeginCreateOrUpdateOptions) (VirtualApplianceSitesCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, networkVirtualApplianceName, siteName, parameters, options)
	if err != nil {
		return VirtualApplianceSitesCreateOrUpdatePollerResponse{}, err
	}
	result := VirtualApplianceSitesCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("VirtualApplianceSitesClient.CreateOrUpdate", "azure-async-operation", resp, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return VirtualApplianceSitesCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &VirtualApplianceSitesCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Creates or updates the specified Network Virtual Appliance Site.
// If the operation fails it returns the *CloudError error type.
func (client *VirtualApplianceSitesClient) createOrUpdate(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, siteName string, parameters VirtualApplianceSite, options *VirtualApplianceSitesBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, networkVirtualApplianceName, siteName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *VirtualApplianceSitesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, siteName string, parameters VirtualApplianceSite, options *VirtualApplianceSitesBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkVirtualAppliances/{networkVirtualApplianceName}/virtualApplianceSites/{siteName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkVirtualApplianceName == "" {
		return nil, errors.New("parameter networkVirtualApplianceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkVirtualApplianceName}", url.PathEscape(networkVirtualApplianceName))
	if siteName == "" {
		return nil, errors.New("parameter siteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{siteName}", url.PathEscape(siteName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *VirtualApplianceSitesClient) createOrUpdateHandleError(resp *http.Response) error {
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

// BeginDelete - Deletes the specified site from a Virtual Appliance.
// If the operation fails it returns the *CloudError error type.
func (client *VirtualApplianceSitesClient) BeginDelete(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, siteName string, options *VirtualApplianceSitesBeginDeleteOptions) (VirtualApplianceSitesDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, networkVirtualApplianceName, siteName, options)
	if err != nil {
		return VirtualApplianceSitesDeletePollerResponse{}, err
	}
	result := VirtualApplianceSitesDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("VirtualApplianceSitesClient.Delete", "location", resp, client.pl, client.deleteHandleError)
	if err != nil {
		return VirtualApplianceSitesDeletePollerResponse{}, err
	}
	result.Poller = &VirtualApplianceSitesDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Deletes the specified site from a Virtual Appliance.
// If the operation fails it returns the *CloudError error type.
func (client *VirtualApplianceSitesClient) deleteOperation(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, siteName string, options *VirtualApplianceSitesBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, networkVirtualApplianceName, siteName, options)
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
func (client *VirtualApplianceSitesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, siteName string, options *VirtualApplianceSitesBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkVirtualAppliances/{networkVirtualApplianceName}/virtualApplianceSites/{siteName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkVirtualApplianceName == "" {
		return nil, errors.New("parameter networkVirtualApplianceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkVirtualApplianceName}", url.PathEscape(networkVirtualApplianceName))
	if siteName == "" {
		return nil, errors.New("parameter siteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{siteName}", url.PathEscape(siteName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *VirtualApplianceSitesClient) deleteHandleError(resp *http.Response) error {
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

// Get - Gets the specified Virtual Appliance Site.
// If the operation fails it returns the *CloudError error type.
func (client *VirtualApplianceSitesClient) Get(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, siteName string, options *VirtualApplianceSitesGetOptions) (VirtualApplianceSitesGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, networkVirtualApplianceName, siteName, options)
	if err != nil {
		return VirtualApplianceSitesGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualApplianceSitesGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VirtualApplianceSitesGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VirtualApplianceSitesClient) getCreateRequest(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, siteName string, options *VirtualApplianceSitesGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkVirtualAppliances/{networkVirtualApplianceName}/virtualApplianceSites/{siteName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkVirtualApplianceName == "" {
		return nil, errors.New("parameter networkVirtualApplianceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkVirtualApplianceName}", url.PathEscape(networkVirtualApplianceName))
	if siteName == "" {
		return nil, errors.New("parameter siteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{siteName}", url.PathEscape(siteName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VirtualApplianceSitesClient) getHandleResponse(resp *http.Response) (VirtualApplianceSitesGetResponse, error) {
	result := VirtualApplianceSitesGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualApplianceSite); err != nil {
		return VirtualApplianceSitesGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *VirtualApplianceSitesClient) getHandleError(resp *http.Response) error {
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

// List - Lists all Network Virtual Appliance Sites in a Network Virtual Appliance resource.
// If the operation fails it returns the *CloudError error type.
func (client *VirtualApplianceSitesClient) List(resourceGroupName string, networkVirtualApplianceName string, options *VirtualApplianceSitesListOptions) *VirtualApplianceSitesListPager {
	return &VirtualApplianceSitesListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, networkVirtualApplianceName, options)
		},
		advancer: func(ctx context.Context, resp VirtualApplianceSitesListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.NetworkVirtualApplianceSiteListResult.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *VirtualApplianceSitesClient) listCreateRequest(ctx context.Context, resourceGroupName string, networkVirtualApplianceName string, options *VirtualApplianceSitesListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkVirtualAppliances/{networkVirtualApplianceName}/virtualApplianceSites"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkVirtualApplianceName == "" {
		return nil, errors.New("parameter networkVirtualApplianceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkVirtualApplianceName}", url.PathEscape(networkVirtualApplianceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *VirtualApplianceSitesClient) listHandleResponse(resp *http.Response) (VirtualApplianceSitesListResponse, error) {
	result := VirtualApplianceSitesListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkVirtualApplianceSiteListResult); err != nil {
		return VirtualApplianceSitesListResponse{}, err
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *VirtualApplianceSitesClient) listHandleError(resp *http.Response) error {
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
