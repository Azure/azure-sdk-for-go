//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

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

// ExpressRouteGatewaysClient contains the methods for the ExpressRouteGateways group.
// Don't use this type directly, use NewExpressRouteGatewaysClient() instead.
type ExpressRouteGatewaysClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewExpressRouteGatewaysClient creates a new instance of ExpressRouteGatewaysClient with the specified values.
//   - subscriptionID - The subscription credentials which uniquely identify the Microsoft Azure subscription. The subscription
//     ID forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewExpressRouteGatewaysClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ExpressRouteGatewaysClient, error) {
	cl, err := arm.NewClient(moduleName+".ExpressRouteGatewaysClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ExpressRouteGatewaysClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a ExpressRoute gateway in a specified resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group.
//   - expressRouteGatewayName - The name of the ExpressRoute gateway.
//   - putExpressRouteGatewayParameters - Parameters required in an ExpressRoute gateway PUT operation.
//   - options - ExpressRouteGatewaysClientBeginCreateOrUpdateOptions contains the optional parameters for the ExpressRouteGatewaysClient.BeginCreateOrUpdate
//     method.
func (client *ExpressRouteGatewaysClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, putExpressRouteGatewayParameters ExpressRouteGateway, options *ExpressRouteGatewaysClientBeginCreateOrUpdateOptions) (*runtime.Poller[ExpressRouteGatewaysClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, expressRouteGatewayName, putExpressRouteGatewayParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ExpressRouteGatewaysClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ExpressRouteGatewaysClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a ExpressRoute gateway in a specified resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
func (client *ExpressRouteGatewaysClient) createOrUpdate(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, putExpressRouteGatewayParameters ExpressRouteGateway, options *ExpressRouteGatewaysClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, expressRouteGatewayName, putExpressRouteGatewayParameters, options)
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
func (client *ExpressRouteGatewaysClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, putExpressRouteGatewayParameters ExpressRouteGateway, options *ExpressRouteGatewaysClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways/{expressRouteGatewayName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if expressRouteGatewayName == "" {
		return nil, errors.New("parameter expressRouteGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{expressRouteGatewayName}", url.PathEscape(expressRouteGatewayName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, putExpressRouteGatewayParameters); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes the specified ExpressRoute gateway in a resource group. An ExpressRoute gateway resource can only
// be deleted when there are no connection subresources.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group.
//   - expressRouteGatewayName - The name of the ExpressRoute gateway.
//   - options - ExpressRouteGatewaysClientBeginDeleteOptions contains the optional parameters for the ExpressRouteGatewaysClient.BeginDelete
//     method.
func (client *ExpressRouteGatewaysClient) BeginDelete(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, options *ExpressRouteGatewaysClientBeginDeleteOptions) (*runtime.Poller[ExpressRouteGatewaysClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, expressRouteGatewayName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ExpressRouteGatewaysClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ExpressRouteGatewaysClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the specified ExpressRoute gateway in a resource group. An ExpressRoute gateway resource can only be deleted
// when there are no connection subresources.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
func (client *ExpressRouteGatewaysClient) deleteOperation(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, options *ExpressRouteGatewaysClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, expressRouteGatewayName, options)
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
func (client *ExpressRouteGatewaysClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, options *ExpressRouteGatewaysClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways/{expressRouteGatewayName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if expressRouteGatewayName == "" {
		return nil, errors.New("parameter expressRouteGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{expressRouteGatewayName}", url.PathEscape(expressRouteGatewayName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Fetches the details of a ExpressRoute gateway in a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group.
//   - expressRouteGatewayName - The name of the ExpressRoute gateway.
//   - options - ExpressRouteGatewaysClientGetOptions contains the optional parameters for the ExpressRouteGatewaysClient.Get
//     method.
func (client *ExpressRouteGatewaysClient) Get(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, options *ExpressRouteGatewaysClientGetOptions) (ExpressRouteGatewaysClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, expressRouteGatewayName, options)
	if err != nil {
		return ExpressRouteGatewaysClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExpressRouteGatewaysClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExpressRouteGatewaysClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ExpressRouteGatewaysClient) getCreateRequest(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, options *ExpressRouteGatewaysClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways/{expressRouteGatewayName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if expressRouteGatewayName == "" {
		return nil, errors.New("parameter expressRouteGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{expressRouteGatewayName}", url.PathEscape(expressRouteGatewayName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ExpressRouteGatewaysClient) getHandleResponse(resp *http.Response) (ExpressRouteGatewaysClientGetResponse, error) {
	result := ExpressRouteGatewaysClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExpressRouteGateway); err != nil {
		return ExpressRouteGatewaysClientGetResponse{}, err
	}
	return result, nil
}

// ListByResourceGroup - Lists ExpressRoute gateways in a given resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group.
//   - options - ExpressRouteGatewaysClientListByResourceGroupOptions contains the optional parameters for the ExpressRouteGatewaysClient.ListByResourceGroup
//     method.
func (client *ExpressRouteGatewaysClient) ListByResourceGroup(ctx context.Context, resourceGroupName string, options *ExpressRouteGatewaysClientListByResourceGroupOptions) (ExpressRouteGatewaysClientListByResourceGroupResponse, error) {
	var err error
	req, err := client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
	if err != nil {
		return ExpressRouteGatewaysClientListByResourceGroupResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExpressRouteGatewaysClientListByResourceGroupResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExpressRouteGatewaysClientListByResourceGroupResponse{}, err
	}
	resp, err := client.listByResourceGroupHandleResponse(httpResp)
	return resp, err
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ExpressRouteGatewaysClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ExpressRouteGatewaysClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ExpressRouteGatewaysClient) listByResourceGroupHandleResponse(resp *http.Response) (ExpressRouteGatewaysClientListByResourceGroupResponse, error) {
	result := ExpressRouteGatewaysClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExpressRouteGatewayList); err != nil {
		return ExpressRouteGatewaysClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// ListBySubscription - Lists ExpressRoute gateways under a given subscription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - options - ExpressRouteGatewaysClientListBySubscriptionOptions contains the optional parameters for the ExpressRouteGatewaysClient.ListBySubscription
//     method.
func (client *ExpressRouteGatewaysClient) ListBySubscription(ctx context.Context, options *ExpressRouteGatewaysClientListBySubscriptionOptions) (ExpressRouteGatewaysClientListBySubscriptionResponse, error) {
	var err error
	req, err := client.listBySubscriptionCreateRequest(ctx, options)
	if err != nil {
		return ExpressRouteGatewaysClientListBySubscriptionResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExpressRouteGatewaysClientListBySubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExpressRouteGatewaysClientListBySubscriptionResponse{}, err
	}
	resp, err := client.listBySubscriptionHandleResponse(httpResp)
	return resp, err
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *ExpressRouteGatewaysClient) listBySubscriptionCreateRequest(ctx context.Context, options *ExpressRouteGatewaysClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/expressRouteGateways"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *ExpressRouteGatewaysClient) listBySubscriptionHandleResponse(resp *http.Response) (ExpressRouteGatewaysClientListBySubscriptionResponse, error) {
	result := ExpressRouteGatewaysClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExpressRouteGatewayList); err != nil {
		return ExpressRouteGatewaysClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginUpdateTags - Updates express route gateway tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The resource group name of the ExpressRouteGateway.
//   - expressRouteGatewayName - The name of the gateway.
//   - expressRouteGatewayParameters - Parameters supplied to update a virtual wan express route gateway tags.
//   - options - ExpressRouteGatewaysClientBeginUpdateTagsOptions contains the optional parameters for the ExpressRouteGatewaysClient.BeginUpdateTags
//     method.
func (client *ExpressRouteGatewaysClient) BeginUpdateTags(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, expressRouteGatewayParameters TagsObject, options *ExpressRouteGatewaysClientBeginUpdateTagsOptions) (*runtime.Poller[ExpressRouteGatewaysClientUpdateTagsResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.updateTags(ctx, resourceGroupName, expressRouteGatewayName, expressRouteGatewayParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ExpressRouteGatewaysClientUpdateTagsResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ExpressRouteGatewaysClientUpdateTagsResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// UpdateTags - Updates express route gateway tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
func (client *ExpressRouteGatewaysClient) updateTags(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, expressRouteGatewayParameters TagsObject, options *ExpressRouteGatewaysClientBeginUpdateTagsOptions) (*http.Response, error) {
	var err error
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, expressRouteGatewayName, expressRouteGatewayParameters, options)
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

// updateTagsCreateRequest creates the UpdateTags request.
func (client *ExpressRouteGatewaysClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, expressRouteGatewayName string, expressRouteGatewayParameters TagsObject, options *ExpressRouteGatewaysClientBeginUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/expressRouteGateways/{expressRouteGatewayName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if expressRouteGatewayName == "" {
		return nil, errors.New("parameter expressRouteGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{expressRouteGatewayName}", url.PathEscape(expressRouteGatewayName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, expressRouteGatewayParameters); err != nil {
	return nil, err
}
	return req, nil
}

