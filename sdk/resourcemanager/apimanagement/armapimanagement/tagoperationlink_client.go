// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// TagOperationLinkClient contains the methods for the TagOperationLink group.
// Don't use this type directly, use NewTagOperationLinkClient() instead.
type TagOperationLinkClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewTagOperationLinkClient creates a new instance of TagOperationLinkClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewTagOperationLinkClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TagOperationLinkClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &TagOperationLinkClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Adds an operation to the specified tag via link.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - tagID - Tag identifier. Must be unique in the current API Management service instance.
//   - operationLinkID - Tag-operation link identifier. Must be unique in the current API Management service instance.
//   - parameters - Create or update parameters.
//   - options - TagOperationLinkClientCreateOrUpdateOptions contains the optional parameters for the TagOperationLinkClient.CreateOrUpdate
//     method.
func (client *TagOperationLinkClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, tagID string, operationLinkID string, parameters TagOperationLinkContract, options *TagOperationLinkClientCreateOrUpdateOptions) (TagOperationLinkClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "TagOperationLinkClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, tagID, operationLinkID, parameters, options)
	if err != nil {
		return TagOperationLinkClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TagOperationLinkClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return TagOperationLinkClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *TagOperationLinkClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, tagID string, operationLinkID string, parameters TagOperationLinkContract, _ *TagOperationLinkClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/tags/{tagId}/operationLinks/{operationLinkId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if tagID == "" {
		return nil, errors.New("parameter tagID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{tagId}", url.PathEscape(tagID))
	if operationLinkID == "" {
		return nil, errors.New("parameter operationLinkID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationLinkId}", url.PathEscape(operationLinkID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *TagOperationLinkClient) createOrUpdateHandleResponse(resp *http.Response) (TagOperationLinkClientCreateOrUpdateResponse, error) {
	result := TagOperationLinkClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TagOperationLinkContract); err != nil {
		return TagOperationLinkClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the specified operation from the specified tag.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - tagID - Tag identifier. Must be unique in the current API Management service instance.
//   - operationLinkID - Tag-operation link identifier. Must be unique in the current API Management service instance.
//   - options - TagOperationLinkClientDeleteOptions contains the optional parameters for the TagOperationLinkClient.Delete method.
func (client *TagOperationLinkClient) Delete(ctx context.Context, resourceGroupName string, serviceName string, tagID string, operationLinkID string, options *TagOperationLinkClientDeleteOptions) (TagOperationLinkClientDeleteResponse, error) {
	var err error
	const operationName = "TagOperationLinkClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, tagID, operationLinkID, options)
	if err != nil {
		return TagOperationLinkClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TagOperationLinkClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return TagOperationLinkClientDeleteResponse{}, err
	}
	return TagOperationLinkClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *TagOperationLinkClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, tagID string, operationLinkID string, _ *TagOperationLinkClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/tags/{tagId}/operationLinks/{operationLinkId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if tagID == "" {
		return nil, errors.New("parameter tagID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{tagId}", url.PathEscape(tagID))
	if operationLinkID == "" {
		return nil, errors.New("parameter operationLinkID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationLinkId}", url.PathEscape(operationLinkID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the operation link for the tag.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - tagID - Tag identifier. Must be unique in the current API Management service instance.
//   - operationLinkID - Tag-operation link identifier. Must be unique in the current API Management service instance.
//   - options - TagOperationLinkClientGetOptions contains the optional parameters for the TagOperationLinkClient.Get method.
func (client *TagOperationLinkClient) Get(ctx context.Context, resourceGroupName string, serviceName string, tagID string, operationLinkID string, options *TagOperationLinkClientGetOptions) (TagOperationLinkClientGetResponse, error) {
	var err error
	const operationName = "TagOperationLinkClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, tagID, operationLinkID, options)
	if err != nil {
		return TagOperationLinkClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TagOperationLinkClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TagOperationLinkClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *TagOperationLinkClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, tagID string, operationLinkID string, _ *TagOperationLinkClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/tags/{tagId}/operationLinks/{operationLinkId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if tagID == "" {
		return nil, errors.New("parameter tagID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{tagId}", url.PathEscape(tagID))
	if operationLinkID == "" {
		return nil, errors.New("parameter operationLinkID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationLinkId}", url.PathEscape(operationLinkID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *TagOperationLinkClient) getHandleResponse(resp *http.Response) (TagOperationLinkClientGetResponse, error) {
	result := TagOperationLinkClientGetResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.TagOperationLinkContract); err != nil {
		return TagOperationLinkClientGetResponse{}, err
	}
	return result, nil
}

// NewListByProductPager - Lists a collection of the operation links associated with a tag.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - tagID - Tag identifier. Must be unique in the current API Management service instance.
//   - options - TagOperationLinkClientListByProductOptions contains the optional parameters for the TagOperationLinkClient.NewListByProductPager
//     method.
func (client *TagOperationLinkClient) NewListByProductPager(resourceGroupName string, serviceName string, tagID string, options *TagOperationLinkClientListByProductOptions) *runtime.Pager[TagOperationLinkClientListByProductResponse] {
	return runtime.NewPager(runtime.PagingHandler[TagOperationLinkClientListByProductResponse]{
		More: func(page TagOperationLinkClientListByProductResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TagOperationLinkClientListByProductResponse) (TagOperationLinkClientListByProductResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "TagOperationLinkClient.NewListByProductPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByProductCreateRequest(ctx, resourceGroupName, serviceName, tagID, options)
			}, nil)
			if err != nil {
				return TagOperationLinkClientListByProductResponse{}, err
			}
			return client.listByProductHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByProductCreateRequest creates the ListByProduct request.
func (client *TagOperationLinkClient) listByProductCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, tagID string, options *TagOperationLinkClientListByProductOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/tags/{tagId}/operationLinks"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if tagID == "" {
		return nil, errors.New("parameter tagID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{tagId}", url.PathEscape(tagID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByProductHandleResponse handles the ListByProduct response.
func (client *TagOperationLinkClient) listByProductHandleResponse(resp *http.Response) (TagOperationLinkClientListByProductResponse, error) {
	result := TagOperationLinkClientListByProductResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TagOperationLinkCollection); err != nil {
		return TagOperationLinkClientListByProductResponse{}, err
	}
	return result, nil
}
