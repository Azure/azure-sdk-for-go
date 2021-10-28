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

// AzureFirewallsClient contains the methods for the AzureFirewalls group.
// Don't use this type directly, use NewAzureFirewallsClient() instead.
type AzureFirewallsClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewAzureFirewallsClient creates a new instance of AzureFirewallsClient with the specified values.
func NewAzureFirewallsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *AzureFirewallsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &AzureFirewallsClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// BeginCreateOrUpdate - Creates or updates the specified Azure Firewall.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, azureFirewallName string, parameters AzureFirewall, options *AzureFirewallsBeginCreateOrUpdateOptions) (AzureFirewallsCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, azureFirewallName, parameters, options)
	if err != nil {
		return AzureFirewallsCreateOrUpdatePollerResponse{}, err
	}
	result := AzureFirewallsCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("AzureFirewallsClient.CreateOrUpdate", "azure-async-operation", resp, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return AzureFirewallsCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &AzureFirewallsCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Creates or updates the specified Azure Firewall.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) createOrUpdate(ctx context.Context, resourceGroupName string, azureFirewallName string, parameters AzureFirewall, options *AzureFirewallsBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, azureFirewallName, parameters, options)
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
func (client *AzureFirewallsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, azureFirewallName string, parameters AzureFirewall, options *AzureFirewallsBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/azureFirewalls/{azureFirewallName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if azureFirewallName == "" {
		return nil, errors.New("parameter azureFirewallName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureFirewallName}", url.PathEscape(azureFirewallName))
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
func (client *AzureFirewallsClient) createOrUpdateHandleError(resp *http.Response) error {
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

// BeginDelete - Deletes the specified Azure Firewall.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) BeginDelete(ctx context.Context, resourceGroupName string, azureFirewallName string, options *AzureFirewallsBeginDeleteOptions) (AzureFirewallsDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, azureFirewallName, options)
	if err != nil {
		return AzureFirewallsDeletePollerResponse{}, err
	}
	result := AzureFirewallsDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("AzureFirewallsClient.Delete", "location", resp, client.pl, client.deleteHandleError)
	if err != nil {
		return AzureFirewallsDeletePollerResponse{}, err
	}
	result.Poller = &AzureFirewallsDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Deletes the specified Azure Firewall.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) deleteOperation(ctx context.Context, resourceGroupName string, azureFirewallName string, options *AzureFirewallsBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, azureFirewallName, options)
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
func (client *AzureFirewallsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, azureFirewallName string, options *AzureFirewallsBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/azureFirewalls/{azureFirewallName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if azureFirewallName == "" {
		return nil, errors.New("parameter azureFirewallName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureFirewallName}", url.PathEscape(azureFirewallName))
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
func (client *AzureFirewallsClient) deleteHandleError(resp *http.Response) error {
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

// Get - Gets the specified Azure Firewall.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) Get(ctx context.Context, resourceGroupName string, azureFirewallName string, options *AzureFirewallsGetOptions) (AzureFirewallsGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, azureFirewallName, options)
	if err != nil {
		return AzureFirewallsGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AzureFirewallsGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AzureFirewallsGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AzureFirewallsClient) getCreateRequest(ctx context.Context, resourceGroupName string, azureFirewallName string, options *AzureFirewallsGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/azureFirewalls/{azureFirewallName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if azureFirewallName == "" {
		return nil, errors.New("parameter azureFirewallName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureFirewallName}", url.PathEscape(azureFirewallName))
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
func (client *AzureFirewallsClient) getHandleResponse(resp *http.Response) (AzureFirewallsGetResponse, error) {
	result := AzureFirewallsGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureFirewall); err != nil {
		return AzureFirewallsGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *AzureFirewallsClient) getHandleError(resp *http.Response) error {
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

// List - Lists all Azure Firewalls in a resource group.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) List(resourceGroupName string, options *AzureFirewallsListOptions) *AzureFirewallsListPager {
	return &AzureFirewallsListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, options)
		},
		advancer: func(ctx context.Context, resp AzureFirewallsListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.AzureFirewallListResult.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *AzureFirewallsClient) listCreateRequest(ctx context.Context, resourceGroupName string, options *AzureFirewallsListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/azureFirewalls"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
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
func (client *AzureFirewallsClient) listHandleResponse(resp *http.Response) (AzureFirewallsListResponse, error) {
	result := AzureFirewallsListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureFirewallListResult); err != nil {
		return AzureFirewallsListResponse{}, err
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *AzureFirewallsClient) listHandleError(resp *http.Response) error {
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

// ListAll - Gets all the Azure Firewalls in a subscription.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) ListAll(options *AzureFirewallsListAllOptions) *AzureFirewallsListAllPager {
	return &AzureFirewallsListAllPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listAllCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp AzureFirewallsListAllResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.AzureFirewallListResult.NextLink)
		},
	}
}

// listAllCreateRequest creates the ListAll request.
func (client *AzureFirewallsClient) listAllCreateRequest(ctx context.Context, options *AzureFirewallsListAllOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/azureFirewalls"
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

// listAllHandleResponse handles the ListAll response.
func (client *AzureFirewallsClient) listAllHandleResponse(resp *http.Response) (AzureFirewallsListAllResponse, error) {
	result := AzureFirewallsListAllResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureFirewallListResult); err != nil {
		return AzureFirewallsListAllResponse{}, err
	}
	return result, nil
}

// listAllHandleError handles the ListAll error response.
func (client *AzureFirewallsClient) listAllHandleError(resp *http.Response) error {
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

// BeginUpdateTags - Updates tags of an Azure Firewall resource.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) BeginUpdateTags(ctx context.Context, resourceGroupName string, azureFirewallName string, parameters TagsObject, options *AzureFirewallsBeginUpdateTagsOptions) (AzureFirewallsUpdateTagsPollerResponse, error) {
	resp, err := client.updateTags(ctx, resourceGroupName, azureFirewallName, parameters, options)
	if err != nil {
		return AzureFirewallsUpdateTagsPollerResponse{}, err
	}
	result := AzureFirewallsUpdateTagsPollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("AzureFirewallsClient.UpdateTags", "azure-async-operation", resp, client.pl, client.updateTagsHandleError)
	if err != nil {
		return AzureFirewallsUpdateTagsPollerResponse{}, err
	}
	result.Poller = &AzureFirewallsUpdateTagsPoller{
		pt: pt,
	}
	return result, nil
}

// UpdateTags - Updates tags of an Azure Firewall resource.
// If the operation fails it returns the *CloudError error type.
func (client *AzureFirewallsClient) updateTags(ctx context.Context, resourceGroupName string, azureFirewallName string, parameters TagsObject, options *AzureFirewallsBeginUpdateTagsOptions) (*http.Response, error) {
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, azureFirewallName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.updateTagsHandleError(resp)
	}
	return resp, nil
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *AzureFirewallsClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, azureFirewallName string, parameters TagsObject, options *AzureFirewallsBeginUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/azureFirewalls/{azureFirewallName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if azureFirewallName == "" {
		return nil, errors.New("parameter azureFirewallName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureFirewallName}", url.PathEscape(azureFirewallName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateTagsHandleError handles the UpdateTags error response.
func (client *AzureFirewallsClient) updateTagsHandleError(resp *http.Response) error {
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
