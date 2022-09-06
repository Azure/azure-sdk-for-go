//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdevtestlabs

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// VirtualMachineSchedulesClient contains the methods for the VirtualMachineSchedules group.
// Don't use this type directly, use NewVirtualMachineSchedulesClient() instead.
type VirtualMachineSchedulesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewVirtualMachineSchedulesClient creates a new instance of VirtualMachineSchedulesClient with the specified values.
// subscriptionID - The subscription ID.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewVirtualMachineSchedulesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VirtualMachineSchedulesClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &VirtualMachineSchedulesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Create or replace an existing schedule.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-15
// resourceGroupName - The name of the resource group.
// labName - The name of the lab.
// virtualMachineName - The name of the virtual machine.
// name - The name of the schedule.
// schedule - A schedule.
// options - VirtualMachineSchedulesClientCreateOrUpdateOptions contains the optional parameters for the VirtualMachineSchedulesClient.CreateOrUpdate
// method.
func (client *VirtualMachineSchedulesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, schedule Schedule, options *VirtualMachineSchedulesClientCreateOrUpdateOptions) (VirtualMachineSchedulesClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, labName, virtualMachineName, name, schedule, options)
	if err != nil {
		return VirtualMachineSchedulesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineSchedulesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return VirtualMachineSchedulesClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *VirtualMachineSchedulesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, schedule Schedule, options *VirtualMachineSchedulesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/virtualmachines/{virtualMachineName}/schedules/{name}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if labName == "" {
		return nil, errors.New("parameter labName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{labName}", url.PathEscape(labName))
	if virtualMachineName == "" {
		return nil, errors.New("parameter virtualMachineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineName}", url.PathEscape(virtualMachineName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, schedule)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *VirtualMachineSchedulesClient) createOrUpdateHandleResponse(resp *http.Response) (VirtualMachineSchedulesClientCreateOrUpdateResponse, error) {
	result := VirtualMachineSchedulesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Schedule); err != nil {
		return VirtualMachineSchedulesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete schedule.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-15
// resourceGroupName - The name of the resource group.
// labName - The name of the lab.
// virtualMachineName - The name of the virtual machine.
// name - The name of the schedule.
// options - VirtualMachineSchedulesClientDeleteOptions contains the optional parameters for the VirtualMachineSchedulesClient.Delete
// method.
func (client *VirtualMachineSchedulesClient) Delete(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, options *VirtualMachineSchedulesClientDeleteOptions) (VirtualMachineSchedulesClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, labName, virtualMachineName, name, options)
	if err != nil {
		return VirtualMachineSchedulesClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineSchedulesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return VirtualMachineSchedulesClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return VirtualMachineSchedulesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *VirtualMachineSchedulesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, options *VirtualMachineSchedulesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/virtualmachines/{virtualMachineName}/schedules/{name}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if labName == "" {
		return nil, errors.New("parameter labName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{labName}", url.PathEscape(labName))
	if virtualMachineName == "" {
		return nil, errors.New("parameter virtualMachineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineName}", url.PathEscape(virtualMachineName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginExecute - Execute a schedule. This operation can take a while to complete.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-15
// resourceGroupName - The name of the resource group.
// labName - The name of the lab.
// virtualMachineName - The name of the virtual machine.
// name - The name of the schedule.
// options - VirtualMachineSchedulesClientBeginExecuteOptions contains the optional parameters for the VirtualMachineSchedulesClient.BeginExecute
// method.
func (client *VirtualMachineSchedulesClient) BeginExecute(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, options *VirtualMachineSchedulesClientBeginExecuteOptions) (*runtime.Poller[VirtualMachineSchedulesClientExecuteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.execute(ctx, resourceGroupName, labName, virtualMachineName, name, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[VirtualMachineSchedulesClientExecuteResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[VirtualMachineSchedulesClientExecuteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Execute - Execute a schedule. This operation can take a while to complete.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-15
func (client *VirtualMachineSchedulesClient) execute(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, options *VirtualMachineSchedulesClientBeginExecuteOptions) (*http.Response, error) {
	req, err := client.executeCreateRequest(ctx, resourceGroupName, labName, virtualMachineName, name, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// executeCreateRequest creates the Execute request.
func (client *VirtualMachineSchedulesClient) executeCreateRequest(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, options *VirtualMachineSchedulesClientBeginExecuteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/virtualmachines/{virtualMachineName}/schedules/{name}/execute"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if labName == "" {
		return nil, errors.New("parameter labName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{labName}", url.PathEscape(labName))
	if virtualMachineName == "" {
		return nil, errors.New("parameter virtualMachineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineName}", url.PathEscape(virtualMachineName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get schedule.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-15
// resourceGroupName - The name of the resource group.
// labName - The name of the lab.
// virtualMachineName - The name of the virtual machine.
// name - The name of the schedule.
// options - VirtualMachineSchedulesClientGetOptions contains the optional parameters for the VirtualMachineSchedulesClient.Get
// method.
func (client *VirtualMachineSchedulesClient) Get(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, options *VirtualMachineSchedulesClientGetOptions) (VirtualMachineSchedulesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, labName, virtualMachineName, name, options)
	if err != nil {
		return VirtualMachineSchedulesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineSchedulesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VirtualMachineSchedulesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VirtualMachineSchedulesClient) getCreateRequest(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, options *VirtualMachineSchedulesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/virtualmachines/{virtualMachineName}/schedules/{name}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if labName == "" {
		return nil, errors.New("parameter labName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{labName}", url.PathEscape(labName))
	if virtualMachineName == "" {
		return nil, errors.New("parameter virtualMachineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineName}", url.PathEscape(virtualMachineName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	reqQP.Set("api-version", "2018-09-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VirtualMachineSchedulesClient) getHandleResponse(resp *http.Response) (VirtualMachineSchedulesClientGetResponse, error) {
	result := VirtualMachineSchedulesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Schedule); err != nil {
		return VirtualMachineSchedulesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List schedules in a given virtual machine.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-15
// resourceGroupName - The name of the resource group.
// labName - The name of the lab.
// virtualMachineName - The name of the virtual machine.
// options - VirtualMachineSchedulesClientListOptions contains the optional parameters for the VirtualMachineSchedulesClient.List
// method.
func (client *VirtualMachineSchedulesClient) NewListPager(resourceGroupName string, labName string, virtualMachineName string, options *VirtualMachineSchedulesClientListOptions) *runtime.Pager[VirtualMachineSchedulesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[VirtualMachineSchedulesClientListResponse]{
		More: func(page VirtualMachineSchedulesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VirtualMachineSchedulesClientListResponse) (VirtualMachineSchedulesClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, labName, virtualMachineName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VirtualMachineSchedulesClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return VirtualMachineSchedulesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VirtualMachineSchedulesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *VirtualMachineSchedulesClient) listCreateRequest(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, options *VirtualMachineSchedulesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/virtualmachines/{virtualMachineName}/schedules"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if labName == "" {
		return nil, errors.New("parameter labName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{labName}", url.PathEscape(labName))
	if virtualMachineName == "" {
		return nil, errors.New("parameter virtualMachineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineName}", url.PathEscape(virtualMachineName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	reqQP.Set("api-version", "2018-09-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *VirtualMachineSchedulesClient) listHandleResponse(resp *http.Response) (VirtualMachineSchedulesClientListResponse, error) {
	result := VirtualMachineSchedulesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ScheduleList); err != nil {
		return VirtualMachineSchedulesClientListResponse{}, err
	}
	return result, nil
}

// Update - Allows modifying tags of schedules. All other properties will be ignored.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-15
// resourceGroupName - The name of the resource group.
// labName - The name of the lab.
// virtualMachineName - The name of the virtual machine.
// name - The name of the schedule.
// schedule - A schedule.
// options - VirtualMachineSchedulesClientUpdateOptions contains the optional parameters for the VirtualMachineSchedulesClient.Update
// method.
func (client *VirtualMachineSchedulesClient) Update(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, schedule ScheduleFragment, options *VirtualMachineSchedulesClientUpdateOptions) (VirtualMachineSchedulesClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, labName, virtualMachineName, name, schedule, options)
	if err != nil {
		return VirtualMachineSchedulesClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VirtualMachineSchedulesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VirtualMachineSchedulesClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *VirtualMachineSchedulesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, labName string, virtualMachineName string, name string, schedule ScheduleFragment, options *VirtualMachineSchedulesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/virtualmachines/{virtualMachineName}/schedules/{name}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if labName == "" {
		return nil, errors.New("parameter labName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{labName}", url.PathEscape(labName))
	if virtualMachineName == "" {
		return nil, errors.New("parameter virtualMachineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{virtualMachineName}", url.PathEscape(virtualMachineName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, schedule)
}

// updateHandleResponse handles the Update response.
func (client *VirtualMachineSchedulesClient) updateHandleResponse(resp *http.Response) (VirtualMachineSchedulesClientUpdateResponse, error) {
	result := VirtualMachineSchedulesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Schedule); err != nil {
		return VirtualMachineSchedulesClientUpdateResponse{}, err
	}
	return result, nil
}
