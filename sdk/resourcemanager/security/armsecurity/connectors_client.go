//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

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

// ConnectorsClient contains the methods for the SecurityConnectors group.
// Don't use this type directly, use NewConnectorsClient() instead.
type ConnectorsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewConnectorsClient creates a new instance of ConnectorsClient with the specified values.
//   - subscriptionID - Azure subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewConnectorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConnectorsClient, error) {
	cl, err := arm.NewClient(moduleName+".ConnectorsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ConnectorsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates a security connector. If a security connector is already created and a subsequent request
// is issued for the same security connector id, then it will be updated.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-03-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - securityConnectorName - The security connector name.
//   - securityConnector - The security connector resource
//   - options - ConnectorsClientCreateOrUpdateOptions contains the optional parameters for the ConnectorsClient.CreateOrUpdate
//     method.
func (client *ConnectorsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, securityConnectorName string, securityConnector Connector, options *ConnectorsClientCreateOrUpdateOptions) (ConnectorsClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, securityConnectorName, securityConnector, options)
	if err != nil {
		return ConnectorsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ConnectorsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return ConnectorsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ConnectorsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, securityConnectorName string, securityConnector Connector, options *ConnectorsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/securityConnectors/{securityConnectorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if securityConnectorName == "" {
		return nil, errors.New("parameter securityConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{securityConnectorName}", url.PathEscape(securityConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, securityConnector); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ConnectorsClient) createOrUpdateHandleResponse(resp *http.Response) (ConnectorsClientCreateOrUpdateResponse, error) {
	result := ConnectorsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Connector); err != nil {
		return ConnectorsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a security connector.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-03-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - securityConnectorName - The security connector name.
//   - options - ConnectorsClientDeleteOptions contains the optional parameters for the ConnectorsClient.Delete method.
func (client *ConnectorsClient) Delete(ctx context.Context, resourceGroupName string, securityConnectorName string, options *ConnectorsClientDeleteOptions) (ConnectorsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, securityConnectorName, options)
	if err != nil {
		return ConnectorsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ConnectorsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return ConnectorsClientDeleteResponse{}, err
	}
	return ConnectorsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ConnectorsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, securityConnectorName string, options *ConnectorsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/securityConnectors/{securityConnectorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if securityConnectorName == "" {
		return nil, errors.New("parameter securityConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{securityConnectorName}", url.PathEscape(securityConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Retrieves details of a specific security connector
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-03-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - securityConnectorName - The security connector name.
//   - options - ConnectorsClientGetOptions contains the optional parameters for the ConnectorsClient.Get method.
func (client *ConnectorsClient) Get(ctx context.Context, resourceGroupName string, securityConnectorName string, options *ConnectorsClientGetOptions) (ConnectorsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, securityConnectorName, options)
	if err != nil {
		return ConnectorsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ConnectorsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ConnectorsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ConnectorsClient) getCreateRequest(ctx context.Context, resourceGroupName string, securityConnectorName string, options *ConnectorsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/securityConnectors/{securityConnectorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if securityConnectorName == "" {
		return nil, errors.New("parameter securityConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{securityConnectorName}", url.PathEscape(securityConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ConnectorsClient) getHandleResponse(resp *http.Response) (ConnectorsClientGetResponse, error) {
	result := ConnectorsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Connector); err != nil {
		return ConnectorsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all the security connectors in the specified subscription. Use the 'nextLink' property in the response
// to get the next page of security connectors for the specified subscription.
//
// Generated from API version 2023-03-01-preview
//   - options - ConnectorsClientListOptions contains the optional parameters for the ConnectorsClient.NewListPager method.
func (client *ConnectorsClient) NewListPager(options *ConnectorsClientListOptions) (*runtime.Pager[ConnectorsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ConnectorsClientListResponse]{
		More: func(page ConnectorsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ConnectorsClientListResponse) (ConnectorsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ConnectorsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ConnectorsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ConnectorsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ConnectorsClient) listCreateRequest(ctx context.Context, options *ConnectorsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/securityConnectors"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ConnectorsClient) listHandleResponse(resp *http.Response) (ConnectorsClientListResponse, error) {
	result := ConnectorsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConnectorsList); err != nil {
		return ConnectorsClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Lists all the security connectors in the specified resource group. Use the 'nextLink' property
// in the response to get the next page of security connectors for the specified resource group.
//
// Generated from API version 2023-03-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - options - ConnectorsClientListByResourceGroupOptions contains the optional parameters for the ConnectorsClient.NewListByResourceGroupPager
//     method.
func (client *ConnectorsClient) NewListByResourceGroupPager(resourceGroupName string, options *ConnectorsClientListByResourceGroupOptions) (*runtime.Pager[ConnectorsClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ConnectorsClientListByResourceGroupResponse]{
		More: func(page ConnectorsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ConnectorsClientListByResourceGroupResponse) (ConnectorsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ConnectorsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ConnectorsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ConnectorsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ConnectorsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ConnectorsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/securityConnectors"
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
	reqQP.Set("api-version", "2023-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ConnectorsClient) listByResourceGroupHandleResponse(resp *http.Response) (ConnectorsClientListByResourceGroupResponse, error) {
	result := ConnectorsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConnectorsList); err != nil {
		return ConnectorsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// Update - Updates a security connector
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-03-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - securityConnectorName - The security connector name.
//   - securityConnector - The security connector resource
//   - options - ConnectorsClientUpdateOptions contains the optional parameters for the ConnectorsClient.Update method.
func (client *ConnectorsClient) Update(ctx context.Context, resourceGroupName string, securityConnectorName string, securityConnector Connector, options *ConnectorsClientUpdateOptions) (ConnectorsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, securityConnectorName, securityConnector, options)
	if err != nil {
		return ConnectorsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ConnectorsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ConnectorsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *ConnectorsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, securityConnectorName string, securityConnector Connector, options *ConnectorsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/securityConnectors/{securityConnectorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if securityConnectorName == "" {
		return nil, errors.New("parameter securityConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{securityConnectorName}", url.PathEscape(securityConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, securityConnector); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *ConnectorsClient) updateHandleResponse(resp *http.Response) (ConnectorsClientUpdateResponse, error) {
	result := ConnectorsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Connector); err != nil {
		return ConnectorsClientUpdateResponse{}, err
	}
	return result, nil
}

