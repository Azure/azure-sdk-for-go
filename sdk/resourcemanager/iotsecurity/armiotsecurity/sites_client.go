//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armiotsecurity

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strings"
)

// SitesClient contains the methods for the Sites group.
// Don't use this type directly, use NewSitesClient() instead.
type SitesClient struct {
	internal *arm.Client
}

// NewSitesClient creates a new instance of SitesClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSitesClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*SitesClient, error) {
	cl, err := arm.NewClient(moduleName+".SitesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SitesClient{
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update IoT site
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-02-01-preview
//   - scope - Scope of the query (IoT Hub, /providers/Microsoft.Devices/iotHubs/myHub)
//   - siteModel - The IoT sites model
//   - options - SitesClientCreateOrUpdateOptions contains the optional parameters for the SitesClient.CreateOrUpdate method.
func (client *SitesClient) CreateOrUpdate(ctx context.Context, scope string, siteModel SiteModel, options *SitesClientCreateOrUpdateOptions) (SitesClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, scope, siteModel, options)
	if err != nil {
		return SitesClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SitesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return SitesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *SitesClient) createOrUpdateCreateRequest(ctx context.Context, scope string, siteModel SiteModel, options *SitesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.IoTSecurity/sites/default"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-02-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, siteModel); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *SitesClient) createOrUpdateHandleResponse(resp *http.Response) (SitesClientCreateOrUpdateResponse, error) {
	result := SitesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SiteModel); err != nil {
		return SitesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete IoT site
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-02-01-preview
//   - scope - Scope of the query (IoT Hub, /providers/Microsoft.Devices/iotHubs/myHub)
//   - options - SitesClientDeleteOptions contains the optional parameters for the SitesClient.Delete method.
func (client *SitesClient) Delete(ctx context.Context, scope string, options *SitesClientDeleteOptions) (SitesClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, scope, options)
	if err != nil {
		return SitesClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SitesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return SitesClientDeleteResponse{}, err
	}
	return SitesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SitesClient) deleteCreateRequest(ctx context.Context, scope string, options *SitesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.IoTSecurity/sites/default"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-02-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get IoT site
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-02-01-preview
//   - scope - Scope of the query (IoT Hub, /providers/Microsoft.Devices/iotHubs/myHub)
//   - options - SitesClientGetOptions contains the optional parameters for the SitesClient.Get method.
func (client *SitesClient) Get(ctx context.Context, scope string, options *SitesClientGetOptions) (SitesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, scope, options)
	if err != nil {
		return SitesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SitesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return SitesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *SitesClient) getCreateRequest(ctx context.Context, scope string, options *SitesClientGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.IoTSecurity/sites/default"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-02-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SitesClient) getHandleResponse(resp *http.Response) (SitesClientGetResponse, error) {
	result := SitesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SiteModel); err != nil {
		return SitesClientGetResponse{}, err
	}
	return result, nil
}

// List - List IoT sites
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-02-01-preview
//   - scope - Scope of the query (IoT Hub, /providers/Microsoft.Devices/iotHubs/myHub)
//   - options - SitesClientListOptions contains the optional parameters for the SitesClient.List method.
func (client *SitesClient) List(ctx context.Context, scope string, options *SitesClientListOptions) (SitesClientListResponse, error) {
	var err error
	req, err := client.listCreateRequest(ctx, scope, options)
	if err != nil {
		return SitesClientListResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SitesClientListResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return SitesClientListResponse{}, err
	}
	resp, err := client.listHandleResponse(httpResp)
	return resp, err
}

// listCreateRequest creates the List request.
func (client *SitesClient) listCreateRequest(ctx context.Context, scope string, options *SitesClientListOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.IoTSecurity/sites"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-02-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SitesClient) listHandleResponse(resp *http.Response) (SitesClientListResponse, error) {
	result := SitesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SitesList); err != nil {
		return SitesClientListResponse{}, err
	}
	return result, nil
}

