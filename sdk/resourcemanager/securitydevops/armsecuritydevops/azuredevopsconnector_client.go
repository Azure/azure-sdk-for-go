//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecuritydevops

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

// AzureDevOpsConnectorClient contains the methods for the AzureDevOpsConnector group.
// Don't use this type directly, use NewAzureDevOpsConnectorClient() instead.
type AzureDevOpsConnectorClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewAzureDevOpsConnectorClient creates a new instance of AzureDevOpsConnectorClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAzureDevOpsConnectorClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AzureDevOpsConnectorClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AzureDevOpsConnectorClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates an Azure DevOps Connector.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - azureDevOpsConnectorName - Name of the AzureDevOps Connector.
//   - azureDevOpsConnector - Connector resource payload.
//   - options - AzureDevOpsConnectorClientBeginCreateOrUpdateOptions contains the optional parameters for the AzureDevOpsConnectorClient.BeginCreateOrUpdate
//     method.
func (client *AzureDevOpsConnectorClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, azureDevOpsConnector AzureDevOpsConnector, options *AzureDevOpsConnectorClientBeginCreateOrUpdateOptions) (*runtime.Poller[AzureDevOpsConnectorClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, azureDevOpsConnectorName, azureDevOpsConnector, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AzureDevOpsConnectorClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[AzureDevOpsConnectorClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Creates or updates an Azure DevOps Connector.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
func (client *AzureDevOpsConnectorClient) createOrUpdate(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, azureDevOpsConnector AzureDevOpsConnector, options *AzureDevOpsConnectorClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "AzureDevOpsConnectorClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, azureDevOpsConnectorName, azureDevOpsConnector, options)
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
func (client *AzureDevOpsConnectorClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, azureDevOpsConnector AzureDevOpsConnector, options *AzureDevOpsConnectorClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/azureDevOpsConnectors/{azureDevOpsConnectorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if azureDevOpsConnectorName == "" {
		return nil, errors.New("parameter azureDevOpsConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureDevOpsConnectorName}", url.PathEscape(azureDevOpsConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, azureDevOpsConnector); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete monitored AzureDevOps Connector details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - azureDevOpsConnectorName - Name of the AzureDevOps Connector.
//   - options - AzureDevOpsConnectorClientBeginDeleteOptions contains the optional parameters for the AzureDevOpsConnectorClient.BeginDelete
//     method.
func (client *AzureDevOpsConnectorClient) BeginDelete(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, options *AzureDevOpsConnectorClientBeginDeleteOptions) (*runtime.Poller[AzureDevOpsConnectorClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, azureDevOpsConnectorName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AzureDevOpsConnectorClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[AzureDevOpsConnectorClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete monitored AzureDevOps Connector details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
func (client *AzureDevOpsConnectorClient) deleteOperation(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, options *AzureDevOpsConnectorClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "AzureDevOpsConnectorClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, azureDevOpsConnectorName, options)
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
func (client *AzureDevOpsConnectorClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, options *AzureDevOpsConnectorClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/azureDevOpsConnectors/{azureDevOpsConnectorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if azureDevOpsConnectorName == "" {
		return nil, errors.New("parameter azureDevOpsConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureDevOpsConnectorName}", url.PathEscape(azureDevOpsConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Returns a monitored AzureDevOps Connector resource for a given ID.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - azureDevOpsConnectorName - Name of the AzureDevOps Connector.
//   - options - AzureDevOpsConnectorClientGetOptions contains the optional parameters for the AzureDevOpsConnectorClient.Get
//     method.
func (client *AzureDevOpsConnectorClient) Get(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, options *AzureDevOpsConnectorClientGetOptions) (AzureDevOpsConnectorClientGetResponse, error) {
	var err error
	const operationName = "AzureDevOpsConnectorClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, azureDevOpsConnectorName, options)
	if err != nil {
		return AzureDevOpsConnectorClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AzureDevOpsConnectorClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AzureDevOpsConnectorClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AzureDevOpsConnectorClient) getCreateRequest(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, options *AzureDevOpsConnectorClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/azureDevOpsConnectors/{azureDevOpsConnectorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if azureDevOpsConnectorName == "" {
		return nil, errors.New("parameter azureDevOpsConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureDevOpsConnectorName}", url.PathEscape(azureDevOpsConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AzureDevOpsConnectorClient) getHandleResponse(resp *http.Response) (AzureDevOpsConnectorClientGetResponse, error) {
	result := AzureDevOpsConnectorClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureDevOpsConnector); err != nil {
		return AzureDevOpsConnectorClientGetResponse{}, err
	}
	return result, nil
}

//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - AzureDevOpsConnectorClientListByResourceGroupOptions contains the optional parameters for the AzureDevOpsConnectorClient.NewListByResourceGroupPager
//     method.
func (client *AzureDevOpsConnectorClient) NewListByResourceGroupPager(resourceGroupName string, options *AzureDevOpsConnectorClientListByResourceGroupOptions) *runtime.Pager[AzureDevOpsConnectorClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[AzureDevOpsConnectorClientListByResourceGroupResponse]{
		More: func(page AzureDevOpsConnectorClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AzureDevOpsConnectorClientListByResourceGroupResponse) (AzureDevOpsConnectorClientListByResourceGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "AzureDevOpsConnectorClient.NewListByResourceGroupPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			}, nil)
			if err != nil {
				return AzureDevOpsConnectorClientListByResourceGroupResponse{}, err
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *AzureDevOpsConnectorClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *AzureDevOpsConnectorClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/azureDevOpsConnectors"
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
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *AzureDevOpsConnectorClient) listByResourceGroupHandleResponse(resp *http.Response) (AzureDevOpsConnectorClientListByResourceGroupResponse, error) {
	result := AzureDevOpsConnectorClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureDevOpsConnectorListResponse); err != nil {
		return AzureDevOpsConnectorClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Returns a list of monitored AzureDevOps Connectors.
//
// Generated from API version 2022-09-01-preview
//   - options - AzureDevOpsConnectorClientListBySubscriptionOptions contains the optional parameters for the AzureDevOpsConnectorClient.NewListBySubscriptionPager
//     method.
func (client *AzureDevOpsConnectorClient) NewListBySubscriptionPager(options *AzureDevOpsConnectorClientListBySubscriptionOptions) *runtime.Pager[AzureDevOpsConnectorClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[AzureDevOpsConnectorClientListBySubscriptionResponse]{
		More: func(page AzureDevOpsConnectorClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AzureDevOpsConnectorClientListBySubscriptionResponse) (AzureDevOpsConnectorClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "AzureDevOpsConnectorClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return AzureDevOpsConnectorClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *AzureDevOpsConnectorClient) listBySubscriptionCreateRequest(ctx context.Context, options *AzureDevOpsConnectorClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.SecurityDevOps/azureDevOpsConnectors"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *AzureDevOpsConnectorClient) listBySubscriptionHandleResponse(resp *http.Response) (AzureDevOpsConnectorClientListBySubscriptionResponse, error) {
	result := AzureDevOpsConnectorClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureDevOpsConnectorListResponse); err != nil {
		return AzureDevOpsConnectorClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update monitored AzureDevOps Connector details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - azureDevOpsConnectorName - Name of the AzureDevOps Connector.
//   - azureDevOpsConnector - Connector resource payload.
//   - options - AzureDevOpsConnectorClientBeginUpdateOptions contains the optional parameters for the AzureDevOpsConnectorClient.BeginUpdate
//     method.
func (client *AzureDevOpsConnectorClient) BeginUpdate(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, azureDevOpsConnector AzureDevOpsConnector, options *AzureDevOpsConnectorClientBeginUpdateOptions) (*runtime.Poller[AzureDevOpsConnectorClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, azureDevOpsConnectorName, azureDevOpsConnector, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AzureDevOpsConnectorClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[AzureDevOpsConnectorClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Update monitored AzureDevOps Connector details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
func (client *AzureDevOpsConnectorClient) update(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, azureDevOpsConnector AzureDevOpsConnector, options *AzureDevOpsConnectorClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "AzureDevOpsConnectorClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, azureDevOpsConnectorName, azureDevOpsConnector, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *AzureDevOpsConnectorClient) updateCreateRequest(ctx context.Context, resourceGroupName string, azureDevOpsConnectorName string, azureDevOpsConnector AzureDevOpsConnector, options *AzureDevOpsConnectorClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/azureDevOpsConnectors/{azureDevOpsConnectorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if azureDevOpsConnectorName == "" {
		return nil, errors.New("parameter azureDevOpsConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureDevOpsConnectorName}", url.PathEscape(azureDevOpsConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, azureDevOpsConnector); err != nil {
		return nil, err
	}
	return req, nil
}
