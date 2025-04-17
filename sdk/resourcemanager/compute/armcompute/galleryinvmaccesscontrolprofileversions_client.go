// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcompute

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

// GalleryInVMAccessControlProfileVersionsClient contains the methods for the GalleryInVMAccessControlProfileVersions group.
// Don't use this type directly, use NewGalleryInVMAccessControlProfileVersionsClient() instead.
type GalleryInVMAccessControlProfileVersionsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewGalleryInVMAccessControlProfileVersionsClient creates a new instance of GalleryInVMAccessControlProfileVersionsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewGalleryInVMAccessControlProfileVersionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*GalleryInVMAccessControlProfileVersionsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &GalleryInVMAccessControlProfileVersionsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update a gallery inVMAccessControlProfile version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-03
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - galleryName - The name of the Shared Image Gallery.
//   - inVMAccessControlProfileName - The name of the gallery inVMAccessControlProfile to be retrieved.
//   - inVMAccessControlProfileVersionName - The name of the gallery inVMAccessControlProfile version to be retrieved.
//   - galleryInVMAccessControlProfileVersion - Parameters supplied to the create or update gallery inVMAccessControlProfile version
//     operation.
//   - options - GalleryInVMAccessControlProfileVersionsClientBeginCreateOrUpdateOptions contains the optional parameters for
//     the GalleryInVMAccessControlProfileVersionsClient.BeginCreateOrUpdate method.
func (client *GalleryInVMAccessControlProfileVersionsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, galleryInVMAccessControlProfileVersion GalleryInVMAccessControlProfileVersion, options *GalleryInVMAccessControlProfileVersionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GalleryInVMAccessControlProfileVersionsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, galleryName, inVMAccessControlProfileName, inVMAccessControlProfileVersionName, galleryInVMAccessControlProfileVersion, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[GalleryInVMAccessControlProfileVersionsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[GalleryInVMAccessControlProfileVersionsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create or update a gallery inVMAccessControlProfile version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-03
func (client *GalleryInVMAccessControlProfileVersionsClient) createOrUpdate(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, galleryInVMAccessControlProfileVersion GalleryInVMAccessControlProfileVersion, options *GalleryInVMAccessControlProfileVersionsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "GalleryInVMAccessControlProfileVersionsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, galleryName, inVMAccessControlProfileName, inVMAccessControlProfileVersionName, galleryInVMAccessControlProfileVersion, options)
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
func (client *GalleryInVMAccessControlProfileVersionsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, galleryInVMAccessControlProfileVersion GalleryInVMAccessControlProfileVersion, _ *GalleryInVMAccessControlProfileVersionsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/inVMAccessControlProfiles/{inVMAccessControlProfileName}/versions/{inVMAccessControlProfileVersionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if galleryName == "" {
		return nil, errors.New("parameter galleryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	if inVMAccessControlProfileName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileName}", url.PathEscape(inVMAccessControlProfileName))
	if inVMAccessControlProfileVersionName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileVersionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileVersionName}", url.PathEscape(inVMAccessControlProfileVersionName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-03")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, galleryInVMAccessControlProfileVersion); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete a gallery inVMAccessControlProfile version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-03
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - galleryName - The name of the Shared Image Gallery.
//   - inVMAccessControlProfileName - The name of the gallery inVMAccessControlProfile to be retrieved.
//   - inVMAccessControlProfileVersionName - The name of the gallery inVMAccessControlProfile version to be retrieved.
//   - options - GalleryInVMAccessControlProfileVersionsClientBeginDeleteOptions contains the optional parameters for the GalleryInVMAccessControlProfileVersionsClient.BeginDelete
//     method.
func (client *GalleryInVMAccessControlProfileVersionsClient) BeginDelete(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, options *GalleryInVMAccessControlProfileVersionsClientBeginDeleteOptions) (*runtime.Poller[GalleryInVMAccessControlProfileVersionsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, galleryName, inVMAccessControlProfileName, inVMAccessControlProfileVersionName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[GalleryInVMAccessControlProfileVersionsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[GalleryInVMAccessControlProfileVersionsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a gallery inVMAccessControlProfile version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-03
func (client *GalleryInVMAccessControlProfileVersionsClient) deleteOperation(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, options *GalleryInVMAccessControlProfileVersionsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "GalleryInVMAccessControlProfileVersionsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, galleryName, inVMAccessControlProfileName, inVMAccessControlProfileVersionName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *GalleryInVMAccessControlProfileVersionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, _ *GalleryInVMAccessControlProfileVersionsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/inVMAccessControlProfiles/{inVMAccessControlProfileName}/versions/{inVMAccessControlProfileVersionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if galleryName == "" {
		return nil, errors.New("parameter galleryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	if inVMAccessControlProfileName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileName}", url.PathEscape(inVMAccessControlProfileName))
	if inVMAccessControlProfileVersionName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileVersionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileVersionName}", url.PathEscape(inVMAccessControlProfileVersionName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-03")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Retrieves information about a gallery inVMAccessControlProfile version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-03
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - galleryName - The name of the Shared Image Gallery.
//   - inVMAccessControlProfileName - The name of the gallery inVMAccessControlProfile to be retrieved.
//   - inVMAccessControlProfileVersionName - The name of the gallery inVMAccessControlProfile version to be retrieved.
//   - options - GalleryInVMAccessControlProfileVersionsClientGetOptions contains the optional parameters for the GalleryInVMAccessControlProfileVersionsClient.Get
//     method.
func (client *GalleryInVMAccessControlProfileVersionsClient) Get(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, options *GalleryInVMAccessControlProfileVersionsClientGetOptions) (GalleryInVMAccessControlProfileVersionsClientGetResponse, error) {
	var err error
	const operationName = "GalleryInVMAccessControlProfileVersionsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, galleryName, inVMAccessControlProfileName, inVMAccessControlProfileVersionName, options)
	if err != nil {
		return GalleryInVMAccessControlProfileVersionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GalleryInVMAccessControlProfileVersionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return GalleryInVMAccessControlProfileVersionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *GalleryInVMAccessControlProfileVersionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, _ *GalleryInVMAccessControlProfileVersionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/inVMAccessControlProfiles/{inVMAccessControlProfileName}/versions/{inVMAccessControlProfileVersionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if galleryName == "" {
		return nil, errors.New("parameter galleryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	if inVMAccessControlProfileName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileName}", url.PathEscape(inVMAccessControlProfileName))
	if inVMAccessControlProfileVersionName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileVersionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileVersionName}", url.PathEscape(inVMAccessControlProfileVersionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-03")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *GalleryInVMAccessControlProfileVersionsClient) getHandleResponse(resp *http.Response) (GalleryInVMAccessControlProfileVersionsClientGetResponse, error) {
	result := GalleryInVMAccessControlProfileVersionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GalleryInVMAccessControlProfileVersion); err != nil {
		return GalleryInVMAccessControlProfileVersionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByGalleryInVMAccessControlProfilePager - List gallery inVMAccessControlProfile versions in a gallery inVMAccessControlProfile
//
// Generated from API version 2024-03-03
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - galleryName - The name of the Shared Image Gallery.
//   - inVMAccessControlProfileName - The name of the gallery inVMAccessControlProfile to be retrieved.
//   - options - GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileOptions contains the optional
//     parameters for the GalleryInVMAccessControlProfileVersionsClient.NewListByGalleryInVMAccessControlProfilePager method.
func (client *GalleryInVMAccessControlProfileVersionsClient) NewListByGalleryInVMAccessControlProfilePager(resourceGroupName string, galleryName string, inVMAccessControlProfileName string, options *GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileOptions) *runtime.Pager[GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse] {
	return runtime.NewPager(runtime.PagingHandler[GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse]{
		More: func(page GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse) (GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "GalleryInVMAccessControlProfileVersionsClient.NewListByGalleryInVMAccessControlProfilePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByGalleryInVMAccessControlProfileCreateRequest(ctx, resourceGroupName, galleryName, inVMAccessControlProfileName, options)
			}, nil)
			if err != nil {
				return GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse{}, err
			}
			return client.listByGalleryInVMAccessControlProfileHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByGalleryInVMAccessControlProfileCreateRequest creates the ListByGalleryInVMAccessControlProfile request.
func (client *GalleryInVMAccessControlProfileVersionsClient) listByGalleryInVMAccessControlProfileCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, _ *GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/inVMAccessControlProfiles/{inVMAccessControlProfileName}/versions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if galleryName == "" {
		return nil, errors.New("parameter galleryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	if inVMAccessControlProfileName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileName}", url.PathEscape(inVMAccessControlProfileName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-03")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByGalleryInVMAccessControlProfileHandleResponse handles the ListByGalleryInVMAccessControlProfile response.
func (client *GalleryInVMAccessControlProfileVersionsClient) listByGalleryInVMAccessControlProfileHandleResponse(resp *http.Response) (GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse, error) {
	result := GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GalleryInVMAccessControlProfileVersionList); err != nil {
		return GalleryInVMAccessControlProfileVersionsClientListByGalleryInVMAccessControlProfileResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update a gallery inVMAccessControlProfile version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-03
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - galleryName - The name of the Shared Image Gallery.
//   - inVMAccessControlProfileName - The name of the gallery inVMAccessControlProfile to be retrieved.
//   - inVMAccessControlProfileVersionName - The name of the gallery inVMAccessControlProfile version to be retrieved.
//   - galleryInVMAccessControlProfileVersion - Parameters supplied to the update gallery inVMAccessControlProfile version operation.
//   - options - GalleryInVMAccessControlProfileVersionsClientBeginUpdateOptions contains the optional parameters for the GalleryInVMAccessControlProfileVersionsClient.BeginUpdate
//     method.
func (client *GalleryInVMAccessControlProfileVersionsClient) BeginUpdate(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, galleryInVMAccessControlProfileVersion GalleryInVMAccessControlProfileVersionUpdate, options *GalleryInVMAccessControlProfileVersionsClientBeginUpdateOptions) (*runtime.Poller[GalleryInVMAccessControlProfileVersionsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, galleryName, inVMAccessControlProfileName, inVMAccessControlProfileVersionName, galleryInVMAccessControlProfileVersion, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[GalleryInVMAccessControlProfileVersionsClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[GalleryInVMAccessControlProfileVersionsClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Update a gallery inVMAccessControlProfile version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-03
func (client *GalleryInVMAccessControlProfileVersionsClient) update(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, galleryInVMAccessControlProfileVersion GalleryInVMAccessControlProfileVersionUpdate, options *GalleryInVMAccessControlProfileVersionsClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "GalleryInVMAccessControlProfileVersionsClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, galleryName, inVMAccessControlProfileName, inVMAccessControlProfileVersionName, galleryInVMAccessControlProfileVersion, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *GalleryInVMAccessControlProfileVersionsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, inVMAccessControlProfileVersionName string, galleryInVMAccessControlProfileVersion GalleryInVMAccessControlProfileVersionUpdate, _ *GalleryInVMAccessControlProfileVersionsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/inVMAccessControlProfiles/{inVMAccessControlProfileName}/versions/{inVMAccessControlProfileVersionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if galleryName == "" {
		return nil, errors.New("parameter galleryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	if inVMAccessControlProfileName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileName}", url.PathEscape(inVMAccessControlProfileName))
	if inVMAccessControlProfileVersionName == "" {
		return nil, errors.New("parameter inVMAccessControlProfileVersionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inVMAccessControlProfileVersionName}", url.PathEscape(inVMAccessControlProfileVersionName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-03")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, galleryInVMAccessControlProfileVersion); err != nil {
		return nil, err
	}
	return req, nil
}
