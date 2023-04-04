//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armalertsmanagement

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

// SmartGroupsClient contains the methods for the SmartGroups group.
// Don't use this type directly, use NewSmartGroupsClient() instead.
type SmartGroupsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewSmartGroupsClient creates a new instance of SmartGroupsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSmartGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SmartGroupsClient, error) {
	cl, err := arm.NewClient(moduleName+".SmartGroupsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SmartGroupsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// ChangeState - Change the state of a Smart Group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-05-preview
//   - smartGroupID - Smart group unique id.
//   - newState - New state of the alert.
//   - options - SmartGroupsClientChangeStateOptions contains the optional parameters for the SmartGroupsClient.ChangeState method.
func (client *SmartGroupsClient) ChangeState(ctx context.Context, smartGroupID string, newState AlertState, options *SmartGroupsClientChangeStateOptions) (SmartGroupsClientChangeStateResponse, error) {
	req, err := client.changeStateCreateRequest(ctx, smartGroupID, newState, options)
	if err != nil {
		return SmartGroupsClientChangeStateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SmartGroupsClientChangeStateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SmartGroupsClientChangeStateResponse{}, runtime.NewResponseError(resp)
	}
	return client.changeStateHandleResponse(resp)
}

// changeStateCreateRequest creates the ChangeState request.
func (client *SmartGroupsClient) changeStateCreateRequest(ctx context.Context, smartGroupID string, newState AlertState, options *SmartGroupsClientChangeStateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.AlertsManagement/smartGroups/{smartGroupId}/changeState"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if smartGroupID == "" {
		return nil, errors.New("parameter smartGroupID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{smartGroupId}", url.PathEscape(smartGroupID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-05-preview")
	reqQP.Set("newState", string(newState))
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// changeStateHandleResponse handles the ChangeState response.
func (client *SmartGroupsClient) changeStateHandleResponse(resp *http.Response) (SmartGroupsClientChangeStateResponse, error) {
	result := SmartGroupsClientChangeStateResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.SmartGroup); err != nil {
		return SmartGroupsClientChangeStateResponse{}, err
	}
	return result, nil
}

// NewGetAllPager - List all the Smart Groups within a specified subscription.
//
// Generated from API version 2019-05-05-preview
//   - options - SmartGroupsClientGetAllOptions contains the optional parameters for the SmartGroupsClient.NewGetAllPager method.
func (client *SmartGroupsClient) NewGetAllPager(options *SmartGroupsClientGetAllOptions) *runtime.Pager[SmartGroupsClientGetAllResponse] {
	return runtime.NewPager(runtime.PagingHandler[SmartGroupsClientGetAllResponse]{
		More: func(page SmartGroupsClientGetAllResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SmartGroupsClientGetAllResponse) (SmartGroupsClientGetAllResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.getAllCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SmartGroupsClientGetAllResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return SmartGroupsClientGetAllResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SmartGroupsClientGetAllResponse{}, runtime.NewResponseError(resp)
			}
			return client.getAllHandleResponse(resp)
		},
	})
}

// getAllCreateRequest creates the GetAll request.
func (client *SmartGroupsClient) getAllCreateRequest(ctx context.Context, options *SmartGroupsClientGetAllOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.AlertsManagement/smartGroups"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.TargetResource != nil {
		reqQP.Set("targetResource", *options.TargetResource)
	}
	if options != nil && options.TargetResourceGroup != nil {
		reqQP.Set("targetResourceGroup", *options.TargetResourceGroup)
	}
	if options != nil && options.TargetResourceType != nil {
		reqQP.Set("targetResourceType", *options.TargetResourceType)
	}
	if options != nil && options.MonitorService != nil {
		reqQP.Set("monitorService", string(*options.MonitorService))
	}
	if options != nil && options.MonitorCondition != nil {
		reqQP.Set("monitorCondition", string(*options.MonitorCondition))
	}
	if options != nil && options.Severity != nil {
		reqQP.Set("severity", string(*options.Severity))
	}
	if options != nil && options.SmartGroupState != nil {
		reqQP.Set("smartGroupState", string(*options.SmartGroupState))
	}
	if options != nil && options.TimeRange != nil {
		reqQP.Set("timeRange", string(*options.TimeRange))
	}
	if options != nil && options.PageCount != nil {
		reqQP.Set("pageCount", strconv.FormatInt(*options.PageCount, 10))
	}
	if options != nil && options.SortBy != nil {
		reqQP.Set("sortBy", string(*options.SortBy))
	}
	if options != nil && options.SortOrder != nil {
		reqQP.Set("sortOrder", string(*options.SortOrder))
	}
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getAllHandleResponse handles the GetAll response.
func (client *SmartGroupsClient) getAllHandleResponse(resp *http.Response) (SmartGroupsClientGetAllResponse, error) {
	result := SmartGroupsClientGetAllResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SmartGroupsList); err != nil {
		return SmartGroupsClientGetAllResponse{}, err
	}
	return result, nil
}

// GetByID - Get information related to a specific Smart Group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-05-preview
//   - smartGroupID - Smart group unique id.
//   - options - SmartGroupsClientGetByIDOptions contains the optional parameters for the SmartGroupsClient.GetByID method.
func (client *SmartGroupsClient) GetByID(ctx context.Context, smartGroupID string, options *SmartGroupsClientGetByIDOptions) (SmartGroupsClientGetByIDResponse, error) {
	req, err := client.getByIDCreateRequest(ctx, smartGroupID, options)
	if err != nil {
		return SmartGroupsClientGetByIDResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SmartGroupsClientGetByIDResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SmartGroupsClientGetByIDResponse{}, runtime.NewResponseError(resp)
	}
	return client.getByIDHandleResponse(resp)
}

// getByIDCreateRequest creates the GetByID request.
func (client *SmartGroupsClient) getByIDCreateRequest(ctx context.Context, smartGroupID string, options *SmartGroupsClientGetByIDOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.AlertsManagement/smartGroups/{smartGroupId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if smartGroupID == "" {
		return nil, errors.New("parameter smartGroupID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{smartGroupId}", url.PathEscape(smartGroupID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getByIDHandleResponse handles the GetByID response.
func (client *SmartGroupsClient) getByIDHandleResponse(resp *http.Response) (SmartGroupsClientGetByIDResponse, error) {
	result := SmartGroupsClientGetByIDResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.SmartGroup); err != nil {
		return SmartGroupsClientGetByIDResponse{}, err
	}
	return result, nil
}

// GetHistory - Get the history a smart group, which captures any Smart Group state changes (New/Acknowledged/Closed) .
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-05-preview
//   - smartGroupID - Smart group unique id.
//   - options - SmartGroupsClientGetHistoryOptions contains the optional parameters for the SmartGroupsClient.GetHistory method.
func (client *SmartGroupsClient) GetHistory(ctx context.Context, smartGroupID string, options *SmartGroupsClientGetHistoryOptions) (SmartGroupsClientGetHistoryResponse, error) {
	req, err := client.getHistoryCreateRequest(ctx, smartGroupID, options)
	if err != nil {
		return SmartGroupsClientGetHistoryResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SmartGroupsClientGetHistoryResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SmartGroupsClientGetHistoryResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHistoryHandleResponse(resp)
}

// getHistoryCreateRequest creates the GetHistory request.
func (client *SmartGroupsClient) getHistoryCreateRequest(ctx context.Context, smartGroupID string, options *SmartGroupsClientGetHistoryOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.AlertsManagement/smartGroups/{smartGroupId}/history"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if smartGroupID == "" {
		return nil, errors.New("parameter smartGroupID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{smartGroupId}", url.PathEscape(smartGroupID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHistoryHandleResponse handles the GetHistory response.
func (client *SmartGroupsClient) getHistoryHandleResponse(resp *http.Response) (SmartGroupsClientGetHistoryResponse, error) {
	result := SmartGroupsClientGetHistoryResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SmartGroupModification); err != nil {
		return SmartGroupsClientGetHistoryResponse{}, err
	}
	return result, nil
}
