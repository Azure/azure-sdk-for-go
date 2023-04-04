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

// WorkbooksClient contains the methods for the Workbooks group.
// Don't use this type directly, use NewWorkbooksClient() instead.
type WorkbooksClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewWorkbooksClient creates a new instance of WorkbooksClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWorkbooksClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkbooksClient, error) {
	cl, err := arm.NewClient(moduleName+".WorkbooksClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WorkbooksClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Create a new workbook.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the resource.
//   - workbookProperties - Properties that need to be specified to create a new workbook.
//   - options - WorkbooksClientCreateOrUpdateOptions contains the optional parameters for the WorkbooksClient.CreateOrUpdate
//     method.
func (client *WorkbooksClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties Workbook, options *WorkbooksClientCreateOrUpdateOptions) (WorkbooksClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, resourceName, workbookProperties, options)
	if err != nil {
		return WorkbooksClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkbooksClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return WorkbooksClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *WorkbooksClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties Workbook, options *WorkbooksClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks/{resourceName}"
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
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, workbookProperties)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *WorkbooksClient) createOrUpdateHandleResponse(resp *http.Response) (WorkbooksClientCreateOrUpdateResponse, error) {
	result := WorkbooksClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Workbook); err != nil {
		return WorkbooksClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a workbook.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the resource.
//   - options - WorkbooksClientDeleteOptions contains the optional parameters for the WorkbooksClient.Delete method.
func (client *WorkbooksClient) Delete(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbooksClientDeleteOptions) (WorkbooksClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return WorkbooksClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkbooksClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return WorkbooksClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return WorkbooksClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *WorkbooksClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbooksClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks/{resourceName}"
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
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a single workbook by its resourceName.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the resource.
//   - options - WorkbooksClientGetOptions contains the optional parameters for the WorkbooksClient.Get method.
func (client *WorkbooksClient) Get(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbooksClientGetOptions) (WorkbooksClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return WorkbooksClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkbooksClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return WorkbooksClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *WorkbooksClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbooksClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks/{resourceName}"
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
	reqQP.Set("api-version", "2022-04-01")
	if options != nil && options.CanFetchContent != nil {
		reqQP.Set("canFetchContent", strconv.FormatBool(*options.CanFetchContent))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *WorkbooksClient) getHandleResponse(resp *http.Response) (WorkbooksClientGetResponse, error) {
	result := WorkbooksClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Workbook); err != nil {
		return WorkbooksClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Get all Workbooks defined within a specified resource group and category.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - category - Category of workbook to return.
//   - options - WorkbooksClientListByResourceGroupOptions contains the optional parameters for the WorkbooksClient.NewListByResourceGroupPager
//     method.
func (client *WorkbooksClient) NewListByResourceGroupPager(resourceGroupName string, category CategoryType, options *WorkbooksClientListByResourceGroupOptions) *runtime.Pager[WorkbooksClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[WorkbooksClientListByResourceGroupResponse]{
		More: func(page WorkbooksClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WorkbooksClientListByResourceGroupResponse) (WorkbooksClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, category, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return WorkbooksClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WorkbooksClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WorkbooksClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *WorkbooksClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, category CategoryType, options *WorkbooksClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks"
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
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *WorkbooksClient) listByResourceGroupHandleResponse(resp *http.Response) (WorkbooksClientListByResourceGroupResponse, error) {
	result := WorkbooksClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkbooksListResult); err != nil {
		return WorkbooksClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Get all Workbooks defined within a specified subscription and category.
//
// Generated from API version 2022-04-01
//   - category - Category of workbook to return.
//   - options - WorkbooksClientListBySubscriptionOptions contains the optional parameters for the WorkbooksClient.NewListBySubscriptionPager
//     method.
func (client *WorkbooksClient) NewListBySubscriptionPager(category CategoryType, options *WorkbooksClientListBySubscriptionOptions) *runtime.Pager[WorkbooksClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[WorkbooksClientListBySubscriptionResponse]{
		More: func(page WorkbooksClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WorkbooksClientListBySubscriptionResponse) (WorkbooksClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, category, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return WorkbooksClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WorkbooksClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WorkbooksClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *WorkbooksClient) listBySubscriptionCreateRequest(ctx context.Context, category CategoryType, options *WorkbooksClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Insights/workbooks"
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
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *WorkbooksClient) listBySubscriptionHandleResponse(resp *http.Response) (WorkbooksClientListBySubscriptionResponse, error) {
	result := WorkbooksClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkbooksListResult); err != nil {
		return WorkbooksClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// RevisionGet - Get a single workbook revision defined by its revisionId.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the resource.
//   - revisionID - The id of the workbook's revision.
//   - options - WorkbooksClientRevisionGetOptions contains the optional parameters for the WorkbooksClient.RevisionGet method.
func (client *WorkbooksClient) RevisionGet(ctx context.Context, resourceGroupName string, resourceName string, revisionID string, options *WorkbooksClientRevisionGetOptions) (WorkbooksClientRevisionGetResponse, error) {
	req, err := client.revisionGetCreateRequest(ctx, resourceGroupName, resourceName, revisionID, options)
	if err != nil {
		return WorkbooksClientRevisionGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkbooksClientRevisionGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return WorkbooksClientRevisionGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.revisionGetHandleResponse(resp)
}

// revisionGetCreateRequest creates the RevisionGet request.
func (client *WorkbooksClient) revisionGetCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, revisionID string, options *WorkbooksClientRevisionGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks/{resourceName}/revisions/{revisionId}"
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
	if revisionID == "" {
		return nil, errors.New("parameter revisionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{revisionId}", url.PathEscape(revisionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// revisionGetHandleResponse handles the RevisionGet response.
func (client *WorkbooksClient) revisionGetHandleResponse(resp *http.Response) (WorkbooksClientRevisionGetResponse, error) {
	result := WorkbooksClientRevisionGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Workbook); err != nil {
		return WorkbooksClientRevisionGetResponse{}, err
	}
	return result, nil
}

// NewRevisionsListPager - Get the revisions for the workbook defined by its resourceName.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the resource.
//   - options - WorkbooksClientRevisionsListOptions contains the optional parameters for the WorkbooksClient.NewRevisionsListPager
//     method.
func (client *WorkbooksClient) NewRevisionsListPager(resourceGroupName string, resourceName string, options *WorkbooksClientRevisionsListOptions) *runtime.Pager[WorkbooksClientRevisionsListResponse] {
	return runtime.NewPager(runtime.PagingHandler[WorkbooksClientRevisionsListResponse]{
		More: func(page WorkbooksClientRevisionsListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WorkbooksClientRevisionsListResponse) (WorkbooksClientRevisionsListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.revisionsListCreateRequest(ctx, resourceGroupName, resourceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return WorkbooksClientRevisionsListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WorkbooksClientRevisionsListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WorkbooksClientRevisionsListResponse{}, runtime.NewResponseError(resp)
			}
			return client.revisionsListHandleResponse(resp)
		},
	})
}

// revisionsListCreateRequest creates the RevisionsList request.
func (client *WorkbooksClient) revisionsListCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbooksClientRevisionsListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks/{resourceName}/revisions"
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
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// revisionsListHandleResponse handles the RevisionsList response.
func (client *WorkbooksClient) revisionsListHandleResponse(resp *http.Response) (WorkbooksClientRevisionsListResponse, error) {
	result := WorkbooksClientRevisionsListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkbooksListResult); err != nil {
		return WorkbooksClientRevisionsListResponse{}, err
	}
	return result, nil
}

// Update - Updates a workbook that has already been added.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the resource.
//   - options - WorkbooksClientUpdateOptions contains the optional parameters for the WorkbooksClient.Update method.
func (client *WorkbooksClient) Update(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbooksClientUpdateOptions) (WorkbooksClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return WorkbooksClientUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkbooksClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated) {
		return WorkbooksClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *WorkbooksClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *WorkbooksClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks/{resourceName}"
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
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.WorkbookUpdateParameters != nil {
		return req, runtime.MarshalAsJSON(req, *options.WorkbookUpdateParameters)
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *WorkbooksClient) updateHandleResponse(resp *http.Response) (WorkbooksClientUpdateResponse, error) {
	result := WorkbooksClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Workbook); err != nil {
		return WorkbooksClientUpdateResponse{}, err
	}
	return result, nil
}
