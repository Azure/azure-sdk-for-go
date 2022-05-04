//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

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

// DeviceSecurityGroupsClient contains the methods for the DeviceSecurityGroups group.
// Don't use this type directly, use NewDeviceSecurityGroupsClient() instead.
type DeviceSecurityGroupsClient struct {
	host string
	pl   runtime.Pipeline
}

// NewDeviceSecurityGroupsClient creates a new instance of DeviceSecurityGroupsClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewDeviceSecurityGroupsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*DeviceSecurityGroupsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &DeviceSecurityGroupsClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// CreateOrUpdate - Use this method to creates or updates the device security group on a specified IoT Hub resource.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceID - The identifier of the resource.
// deviceSecurityGroupName - The name of the device security group. Note that the name of the device security group is case
// insensitive.
// deviceSecurityGroup - Security group object.
// options - DeviceSecurityGroupsClientCreateOrUpdateOptions contains the optional parameters for the DeviceSecurityGroupsClient.CreateOrUpdate
// method.
func (client *DeviceSecurityGroupsClient) CreateOrUpdate(ctx context.Context, resourceID string, deviceSecurityGroupName string, deviceSecurityGroup DeviceSecurityGroup, options *DeviceSecurityGroupsClientCreateOrUpdateOptions) (DeviceSecurityGroupsClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceID, deviceSecurityGroupName, deviceSecurityGroup, options)
	if err != nil {
		return DeviceSecurityGroupsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DeviceSecurityGroupsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return DeviceSecurityGroupsClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *DeviceSecurityGroupsClient) createOrUpdateCreateRequest(ctx context.Context, resourceID string, deviceSecurityGroupName string, deviceSecurityGroup DeviceSecurityGroup, options *DeviceSecurityGroupsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{resourceId}/providers/Microsoft.Security/deviceSecurityGroups/{deviceSecurityGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
	if deviceSecurityGroupName == "" {
		return nil, errors.New("parameter deviceSecurityGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{deviceSecurityGroupName}", url.PathEscape(deviceSecurityGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, deviceSecurityGroup)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *DeviceSecurityGroupsClient) createOrUpdateHandleResponse(resp *http.Response) (DeviceSecurityGroupsClientCreateOrUpdateResponse, error) {
	result := DeviceSecurityGroupsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeviceSecurityGroup); err != nil {
		return DeviceSecurityGroupsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - User this method to deletes the device security group.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceID - The identifier of the resource.
// deviceSecurityGroupName - The name of the device security group. Note that the name of the device security group is case
// insensitive.
// options - DeviceSecurityGroupsClientDeleteOptions contains the optional parameters for the DeviceSecurityGroupsClient.Delete
// method.
func (client *DeviceSecurityGroupsClient) Delete(ctx context.Context, resourceID string, deviceSecurityGroupName string, options *DeviceSecurityGroupsClientDeleteOptions) (DeviceSecurityGroupsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceID, deviceSecurityGroupName, options)
	if err != nil {
		return DeviceSecurityGroupsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DeviceSecurityGroupsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return DeviceSecurityGroupsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return DeviceSecurityGroupsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DeviceSecurityGroupsClient) deleteCreateRequest(ctx context.Context, resourceID string, deviceSecurityGroupName string, options *DeviceSecurityGroupsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/{resourceId}/providers/Microsoft.Security/deviceSecurityGroups/{deviceSecurityGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
	if deviceSecurityGroupName == "" {
		return nil, errors.New("parameter deviceSecurityGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{deviceSecurityGroupName}", url.PathEscape(deviceSecurityGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Use this method to get the device security group for the specified IoT Hub resource.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceID - The identifier of the resource.
// deviceSecurityGroupName - The name of the device security group. Note that the name of the device security group is case
// insensitive.
// options - DeviceSecurityGroupsClientGetOptions contains the optional parameters for the DeviceSecurityGroupsClient.Get
// method.
func (client *DeviceSecurityGroupsClient) Get(ctx context.Context, resourceID string, deviceSecurityGroupName string, options *DeviceSecurityGroupsClientGetOptions) (DeviceSecurityGroupsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceID, deviceSecurityGroupName, options)
	if err != nil {
		return DeviceSecurityGroupsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DeviceSecurityGroupsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DeviceSecurityGroupsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DeviceSecurityGroupsClient) getCreateRequest(ctx context.Context, resourceID string, deviceSecurityGroupName string, options *DeviceSecurityGroupsClientGetOptions) (*policy.Request, error) {
	urlPath := "/{resourceId}/providers/Microsoft.Security/deviceSecurityGroups/{deviceSecurityGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
	if deviceSecurityGroupName == "" {
		return nil, errors.New("parameter deviceSecurityGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{deviceSecurityGroupName}", url.PathEscape(deviceSecurityGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DeviceSecurityGroupsClient) getHandleResponse(resp *http.Response) (DeviceSecurityGroupsClientGetResponse, error) {
	result := DeviceSecurityGroupsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeviceSecurityGroup); err != nil {
		return DeviceSecurityGroupsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Use this method get the list of device security groups for the specified IoT Hub resource.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceID - The identifier of the resource.
// options - DeviceSecurityGroupsClientListOptions contains the optional parameters for the DeviceSecurityGroupsClient.List
// method.
func (client *DeviceSecurityGroupsClient) NewListPager(resourceID string, options *DeviceSecurityGroupsClientListOptions) *runtime.Pager[DeviceSecurityGroupsClientListResponse] {
	return runtime.NewPager(runtime.PageProcessor[DeviceSecurityGroupsClientListResponse]{
		More: func(page DeviceSecurityGroupsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DeviceSecurityGroupsClientListResponse) (DeviceSecurityGroupsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceID, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return DeviceSecurityGroupsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return DeviceSecurityGroupsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DeviceSecurityGroupsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *DeviceSecurityGroupsClient) listCreateRequest(ctx context.Context, resourceID string, options *DeviceSecurityGroupsClientListOptions) (*policy.Request, error) {
	urlPath := "/{resourceId}/providers/Microsoft.Security/deviceSecurityGroups"
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *DeviceSecurityGroupsClient) listHandleResponse(resp *http.Response) (DeviceSecurityGroupsClientListResponse, error) {
	result := DeviceSecurityGroupsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeviceSecurityGroupList); err != nil {
		return DeviceSecurityGroupsClientListResponse{}, err
	}
	return result, nil
}
