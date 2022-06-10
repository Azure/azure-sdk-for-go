//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmigrate

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
	"strings"
)

// GroupsClient contains the methods for the Groups group.
// Don't use this type directly, use NewGroupsClient() instead.
type GroupsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewGroupsClient creates a new instance of GroupsClient with the specified values.
// subscriptionID - Azure Subscription Id in which project was created.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*GroupsClient, error) {
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
	client := &GroupsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Create - Create a new group by sending a json object of type 'group' as given in Models section as part of the Request
// Body. The group name in a project is unique.
// This operation is Idempotent.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// options - GroupsClientCreateOptions contains the optional parameters for the GroupsClient.Create method.
func (client *GroupsClient) Create(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *GroupsClientCreateOptions) (GroupsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, projectName, groupName, options)
	if err != nil {
		return GroupsClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return GroupsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return GroupsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *GroupsClient) createCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *GroupsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if projectName == "" {
		return nil, errors.New("parameter projectName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{projectName}", url.PathEscape(projectName))
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Group != nil {
		return req, runtime.MarshalAsJSON(req, *options.Group)
	}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *GroupsClient) createHandleResponse(resp *http.Response) (GroupsClientCreateResponse, error) {
	result := GroupsClientCreateResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Group); err != nil {
		return GroupsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Delete the group from the project. The machines remain in the project. Deleting a non-existent group results in
// a no-operation.
// A group is an aggregation mechanism for machines in a project. Therefore, deleting group does not delete machines in it.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// options - GroupsClientDeleteOptions contains the optional parameters for the GroupsClient.Delete method.
func (client *GroupsClient) Delete(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *GroupsClientDeleteOptions) (GroupsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, projectName, groupName, options)
	if err != nil {
		return GroupsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return GroupsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return GroupsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return client.deleteHandleResponse(resp)
}

// deleteCreateRequest creates the Delete request.
func (client *GroupsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *GroupsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if projectName == "" {
		return nil, errors.New("parameter projectName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{projectName}", url.PathEscape(projectName))
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *GroupsClient) deleteHandleResponse(resp *http.Response) (GroupsClientDeleteResponse, error) {
	result := GroupsClientDeleteResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	return result, nil
}

// Get - Get information related to a specific group in the project. Returns a json object of type 'group' as specified in
// the models section.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// options - GroupsClientGetOptions contains the optional parameters for the GroupsClient.Get method.
func (client *GroupsClient) Get(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *GroupsClientGetOptions) (GroupsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, projectName, groupName, options)
	if err != nil {
		return GroupsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return GroupsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return GroupsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *GroupsClient) getCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *GroupsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if projectName == "" {
		return nil, errors.New("parameter projectName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{projectName}", url.PathEscape(projectName))
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *GroupsClient) getHandleResponse(resp *http.Response) (GroupsClientGetResponse, error) {
	result := GroupsClientGetResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Group); err != nil {
		return GroupsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByProjectPager - Get all groups created in the project. Returns a json array of objects of type 'group' as specified
// in the Models section.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// options - GroupsClientListByProjectOptions contains the optional parameters for the GroupsClient.ListByProject method.
func (client *GroupsClient) NewListByProjectPager(resourceGroupName string, projectName string, options *GroupsClientListByProjectOptions) *runtime.Pager[GroupsClientListByProjectResponse] {
	return runtime.NewPager(runtime.PagingHandler[GroupsClientListByProjectResponse]{
		More: func(page GroupsClientListByProjectResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *GroupsClientListByProjectResponse) (GroupsClientListByProjectResponse, error) {
			req, err := client.listByProjectCreateRequest(ctx, resourceGroupName, projectName, options)
			if err != nil {
				return GroupsClientListByProjectResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return GroupsClientListByProjectResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return GroupsClientListByProjectResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByProjectHandleResponse(resp)
		},
	})
}

// listByProjectCreateRequest creates the ListByProject request.
func (client *GroupsClient) listByProjectCreateRequest(ctx context.Context, resourceGroupName string, projectName string, options *GroupsClientListByProjectOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if projectName == "" {
		return nil, errors.New("parameter projectName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{projectName}", url.PathEscape(projectName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByProjectHandleResponse handles the ListByProject response.
func (client *GroupsClient) listByProjectHandleResponse(resp *http.Response) (GroupsClientListByProjectResponse, error) {
	result := GroupsClientListByProjectResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.GroupResultList); err != nil {
		return GroupsClientListByProjectResponse{}, err
	}
	return result, nil
}

// UpdateMachines - Update machines in group by adding or removing machines.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2019-10-01
// resourceGroupName - Name of the Azure Resource Group that project is part of.
// projectName - Name of the Azure Migrate project.
// groupName - Unique name of a group within a project.
// options - GroupsClientUpdateMachinesOptions contains the optional parameters for the GroupsClient.UpdateMachines method.
func (client *GroupsClient) UpdateMachines(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *GroupsClientUpdateMachinesOptions) (GroupsClientUpdateMachinesResponse, error) {
	req, err := client.updateMachinesCreateRequest(ctx, resourceGroupName, projectName, groupName, options)
	if err != nil {
		return GroupsClientUpdateMachinesResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return GroupsClientUpdateMachinesResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return GroupsClientUpdateMachinesResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateMachinesHandleResponse(resp)
}

// updateMachinesCreateRequest creates the UpdateMachines request.
func (client *GroupsClient) updateMachinesCreateRequest(ctx context.Context, resourceGroupName string, projectName string, groupName string, options *GroupsClientUpdateMachinesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/assessmentProjects/{projectName}/groups/{groupName}/updateMachines"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if projectName == "" {
		return nil, errors.New("parameter projectName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{projectName}", url.PathEscape(projectName))
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.GroupUpdateProperties != nil {
		return req, runtime.MarshalAsJSON(req, *options.GroupUpdateProperties)
	}
	return req, nil
}

// updateMachinesHandleResponse handles the UpdateMachines response.
func (client *GroupsClient) updateMachinesHandleResponse(resp *http.Response) (GroupsClientUpdateMachinesResponse, error) {
	result := GroupsClientUpdateMachinesResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Group); err != nil {
		return GroupsClientUpdateMachinesResponse{}, err
	}
	return result, nil
}
