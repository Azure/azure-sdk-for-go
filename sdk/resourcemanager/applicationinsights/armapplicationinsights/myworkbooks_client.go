//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armapplicationinsights

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

// MyWorkbooksClient contains the methods for the MyWorkbooks group.
// Don't use this type directly, use NewMyWorkbooksClient() instead.
type MyWorkbooksClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewMyWorkbooksClient creates a new instance of MyWorkbooksClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewMyWorkbooksClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MyWorkbooksClient, error) {
	cl, err := arm.NewClient(moduleName+".MyWorkbooksClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &MyWorkbooksClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Create a new private workbook.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-08
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - workbookProperties - Properties that need to be specified to create a new private workbook.
//   - options - MyWorkbooksClientCreateOrUpdateOptions contains the optional parameters for the MyWorkbooksClient.CreateOrUpdate
//     method.
func (client *MyWorkbooksClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties MyWorkbook, options *MyWorkbooksClientCreateOrUpdateOptions) (MyWorkbooksClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, resourceName, workbookProperties, options)
	if err != nil {
		return MyWorkbooksClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return MyWorkbooksClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return MyWorkbooksClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *MyWorkbooksClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties MyWorkbook, options *MyWorkbooksClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/myWorkbooks/{resourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.SourceID != nil {
		reqQP.Set("sourceId", *options.SourceID)
	}
	reqQP.Set("api-version", "2021-03-08")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, workbookProperties)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *MyWorkbooksClient) createOrUpdateHandleResponse(resp *http.Response) (MyWorkbooksClientCreateOrUpdateResponse, error) {
	result := MyWorkbooksClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MyWorkbook); err != nil {
		return MyWorkbooksClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a private workbook.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-08
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - options - MyWorkbooksClientDeleteOptions contains the optional parameters for the MyWorkbooksClient.Delete method.
func (client *MyWorkbooksClient) Delete(ctx context.Context, resourceGroupName string, resourceName string, options *MyWorkbooksClientDeleteOptions) (MyWorkbooksClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return MyWorkbooksClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return MyWorkbooksClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return MyWorkbooksClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return MyWorkbooksClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *MyWorkbooksClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *MyWorkbooksClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/myWorkbooks/{resourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-08")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a single private workbook by its resourceName.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-08
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - options - MyWorkbooksClientGetOptions contains the optional parameters for the MyWorkbooksClient.Get method.
func (client *MyWorkbooksClient) Get(ctx context.Context, resourceGroupName string, resourceName string, options *MyWorkbooksClientGetOptions) (MyWorkbooksClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return MyWorkbooksClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return MyWorkbooksClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return MyWorkbooksClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *MyWorkbooksClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *MyWorkbooksClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/myWorkbooks/{resourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-08")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *MyWorkbooksClient) getHandleResponse(resp *http.Response) (MyWorkbooksClientGetResponse, error) {
	result := MyWorkbooksClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MyWorkbook); err != nil {
		return MyWorkbooksClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Get all private workbooks defined within a specified resource group and category.
//
// Generated from API version 2021-03-08
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - category - Category of workbook to return.
//   - options - MyWorkbooksClientListByResourceGroupOptions contains the optional parameters for the MyWorkbooksClient.NewListByResourceGroupPager
//     method.
func (client *MyWorkbooksClient) NewListByResourceGroupPager(resourceGroupName string, category CategoryType, options *MyWorkbooksClientListByResourceGroupOptions) *runtime.Pager[MyWorkbooksClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[MyWorkbooksClientListByResourceGroupResponse]{
		More: func(page MyWorkbooksClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *MyWorkbooksClientListByResourceGroupResponse) (MyWorkbooksClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, category, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return MyWorkbooksClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return MyWorkbooksClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return MyWorkbooksClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *MyWorkbooksClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, category CategoryType, options *MyWorkbooksClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/myWorkbooks"
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
	reqQP.Set("category", string(category))
	if options != nil && options.Tags != nil {
		reqQP.Set("tags", strings.Join(options.Tags, ","))
	}
	if options != nil && options.SourceID != nil {
		reqQP.Set("sourceId", *options.SourceID)
	}
	if options != nil && options.CanFetchContent != nil {
		reqQP.Set("canFetchContent", strconv.FormatBool(*options.CanFetchContent))
	}
	reqQP.Set("api-version", "2021-03-08")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *MyWorkbooksClient) listByResourceGroupHandleResponse(resp *http.Response) (MyWorkbooksClientListByResourceGroupResponse, error) {
	result := MyWorkbooksClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MyWorkbooksListResult); err != nil {
		return MyWorkbooksClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Get all private workbooks defined within a specified subscription and category.
//
// Generated from API version 2021-03-08
//   - category - Category of workbook to return.
//   - options - MyWorkbooksClientListBySubscriptionOptions contains the optional parameters for the MyWorkbooksClient.NewListBySubscriptionPager
//     method.
func (client *MyWorkbooksClient) NewListBySubscriptionPager(category CategoryType, options *MyWorkbooksClientListBySubscriptionOptions) *runtime.Pager[MyWorkbooksClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[MyWorkbooksClientListBySubscriptionResponse]{
		More: func(page MyWorkbooksClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *MyWorkbooksClientListBySubscriptionResponse) (MyWorkbooksClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, category, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return MyWorkbooksClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return MyWorkbooksClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return MyWorkbooksClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *MyWorkbooksClient) listBySubscriptionCreateRequest(ctx context.Context, category CategoryType, options *MyWorkbooksClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Insights/myWorkbooks"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("category", string(category))
	if options != nil && options.Tags != nil {
		reqQP.Set("tags", strings.Join(options.Tags, ","))
	}
	if options != nil && options.CanFetchContent != nil {
		reqQP.Set("canFetchContent", strconv.FormatBool(*options.CanFetchContent))
	}
	reqQP.Set("api-version", "2021-03-08")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *MyWorkbooksClient) listBySubscriptionHandleResponse(resp *http.Response) (MyWorkbooksClientListBySubscriptionResponse, error) {
	result := MyWorkbooksClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MyWorkbooksListResult); err != nil {
		return MyWorkbooksClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Updates a private workbook that has already been added.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-08
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - workbookProperties - Properties that need to be specified to create a new private workbook.
//   - options - MyWorkbooksClientUpdateOptions contains the optional parameters for the MyWorkbooksClient.Update method.
func (client *MyWorkbooksClient) Update(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties MyWorkbook, options *MyWorkbooksClientUpdateOptions) (MyWorkbooksClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, workbookProperties, options)
	if err != nil {
		return MyWorkbooksClientUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return MyWorkbooksClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated) {
		return MyWorkbooksClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *MyWorkbooksClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties MyWorkbook, options *MyWorkbooksClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/myWorkbooks/{resourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.SourceID != nil {
		reqQP.Set("sourceId", *options.SourceID)
	}
	reqQP.Set("api-version", "2021-03-08")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, workbookProperties)
}

// updateHandleResponse handles the Update response.
func (client *MyWorkbooksClient) updateHandleResponse(resp *http.Response) (MyWorkbooksClientUpdateResponse, error) {
	result := MyWorkbooksClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MyWorkbook); err != nil {
		return MyWorkbooksClientUpdateResponse{}, err
	}
	return result, nil
}
