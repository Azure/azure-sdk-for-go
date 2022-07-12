//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcompute

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

// SharedGalleryImageVersionsClient contains the methods for the SharedGalleryImageVersions group.
// Don't use this type directly, use NewSharedGalleryImageVersionsClient() instead.
type SharedGalleryImageVersionsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewSharedGalleryImageVersionsClient creates a new instance of SharedGalleryImageVersionsClient with the specified values.
// subscriptionID - Subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID forms
// part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewSharedGalleryImageVersionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SharedGalleryImageVersionsClient, error) {
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
	client := &SharedGalleryImageVersionsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Get - Get a shared gallery image version by subscription id or tenant id.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-03
// location - Resource location.
// galleryUniqueName - The unique name of the Shared Gallery.
// galleryImageName - The name of the Shared Gallery Image Definition from which the Image Versions are to be listed.
// galleryImageVersionName - The name of the gallery image version to be created. Needs to follow semantic version name pattern:
// The allowed characters are digit and period. Digits must be within the range of a 32-bit integer.
// Format: ..
// options - SharedGalleryImageVersionsClientGetOptions contains the optional parameters for the SharedGalleryImageVersionsClient.Get
// method.
func (client *SharedGalleryImageVersionsClient) Get(ctx context.Context, location string, galleryUniqueName string, galleryImageName string, galleryImageVersionName string, options *SharedGalleryImageVersionsClientGetOptions) (SharedGalleryImageVersionsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, location, galleryUniqueName, galleryImageName, galleryImageVersionName, options)
	if err != nil {
		return SharedGalleryImageVersionsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SharedGalleryImageVersionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SharedGalleryImageVersionsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SharedGalleryImageVersionsClient) getCreateRequest(ctx context.Context, location string, galleryUniqueName string, galleryImageName string, galleryImageVersionName string, options *SharedGalleryImageVersionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/sharedGalleries/{galleryUniqueName}/images/{galleryImageName}/versions/{galleryImageVersionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if galleryUniqueName == "" {
		return nil, errors.New("parameter galleryUniqueName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryUniqueName}", url.PathEscape(galleryUniqueName))
	if galleryImageName == "" {
		return nil, errors.New("parameter galleryImageName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryImageName}", url.PathEscape(galleryImageName))
	if galleryImageVersionName == "" {
		return nil, errors.New("parameter galleryImageVersionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryImageVersionName}", url.PathEscape(galleryImageVersionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-03")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SharedGalleryImageVersionsClient) getHandleResponse(resp *http.Response) (SharedGalleryImageVersionsClientGetResponse, error) {
	result := SharedGalleryImageVersionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SharedGalleryImageVersion); err != nil {
		return SharedGalleryImageVersionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List shared gallery image versions by subscription id or tenant id.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-03
// location - Resource location.
// galleryUniqueName - The unique name of the Shared Gallery.
// galleryImageName - The name of the Shared Gallery Image Definition from which the Image Versions are to be listed.
// options - SharedGalleryImageVersionsClientListOptions contains the optional parameters for the SharedGalleryImageVersionsClient.List
// method.
func (client *SharedGalleryImageVersionsClient) NewListPager(location string, galleryUniqueName string, galleryImageName string, options *SharedGalleryImageVersionsClientListOptions) *runtime.Pager[SharedGalleryImageVersionsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[SharedGalleryImageVersionsClientListResponse]{
		More: func(page SharedGalleryImageVersionsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SharedGalleryImageVersionsClientListResponse) (SharedGalleryImageVersionsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, location, galleryUniqueName, galleryImageName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SharedGalleryImageVersionsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return SharedGalleryImageVersionsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SharedGalleryImageVersionsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *SharedGalleryImageVersionsClient) listCreateRequest(ctx context.Context, location string, galleryUniqueName string, galleryImageName string, options *SharedGalleryImageVersionsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/locations/{location}/sharedGalleries/{galleryUniqueName}/images/{galleryImageName}/versions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if galleryUniqueName == "" {
		return nil, errors.New("parameter galleryUniqueName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryUniqueName}", url.PathEscape(galleryUniqueName))
	if galleryImageName == "" {
		return nil, errors.New("parameter galleryImageName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryImageName}", url.PathEscape(galleryImageName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-03")
	if options != nil && options.SharedTo != nil {
		reqQP.Set("sharedTo", string(*options.SharedTo))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SharedGalleryImageVersionsClient) listHandleResponse(resp *http.Response) (SharedGalleryImageVersionsClientListResponse, error) {
	result := SharedGalleryImageVersionsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SharedGalleryImageVersionList); err != nil {
		return SharedGalleryImageVersionsClientListResponse{}, err
	}
	return result, nil
}
