//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsphere

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

// CatalogsClient contains the methods for the Catalogs group.
// Don't use this type directly, use NewCatalogsClient() instead.
type CatalogsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewCatalogsClient creates a new instance of CatalogsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCatalogsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CatalogsClient, error) {
	cl, err := arm.NewClient(moduleName+".CatalogsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CatalogsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CountDevices - Counts devices in catalog.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - options - CatalogsClientCountDevicesOptions contains the optional parameters for the CatalogsClient.CountDevices method.
func (client *CatalogsClient) CountDevices(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientCountDevicesOptions) (CatalogsClientCountDevicesResponse, error) {
	var err error
	req, err := client.countDevicesCreateRequest(ctx, resourceGroupName, catalogName, options)
	if err != nil {
		return CatalogsClientCountDevicesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CatalogsClientCountDevicesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CatalogsClientCountDevicesResponse{}, err
	}
	resp, err := client.countDevicesHandleResponse(httpResp)
	return resp, err
}

// countDevicesCreateRequest creates the CountDevices request.
func (client *CatalogsClient) countDevicesCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientCountDevicesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}/countDevices"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// countDevicesHandleResponse handles the CountDevices response.
func (client *CatalogsClient) countDevicesHandleResponse(resp *http.Response) (CatalogsClientCountDevicesResponse, error) {
	result := CatalogsClientCountDevicesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CountDeviceResponse); err != nil {
		return CatalogsClientCountDevicesResponse{}, err
	}
	return result, nil
}

// BeginCreateOrUpdate - Create a Catalog
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - resource - Resource create parameters.
//   - options - CatalogsClientBeginCreateOrUpdateOptions contains the optional parameters for the CatalogsClient.BeginCreateOrUpdate
//     method.
func (client *CatalogsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, catalogName string, resource Catalog, options *CatalogsClientBeginCreateOrUpdateOptions) (*runtime.Poller[CatalogsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, catalogName, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CatalogsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CatalogsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create a Catalog
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
func (client *CatalogsClient) createOrUpdate(ctx context.Context, resourceGroupName string, catalogName string, resource Catalog, options *CatalogsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, catalogName, resource, options)
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
func (client *CatalogsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, resource Catalog, options *CatalogsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Delete a Catalog
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - options - CatalogsClientBeginDeleteOptions contains the optional parameters for the CatalogsClient.BeginDelete method.
func (client *CatalogsClient) BeginDelete(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientBeginDeleteOptions) (*runtime.Poller[CatalogsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, catalogName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CatalogsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CatalogsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Delete a Catalog
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
func (client *CatalogsClient) deleteOperation(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, catalogName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *CatalogsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a Catalog
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - options - CatalogsClientGetOptions contains the optional parameters for the CatalogsClient.Get method.
func (client *CatalogsClient) Get(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientGetOptions) (CatalogsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, catalogName, options)
	if err != nil {
		return CatalogsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CatalogsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CatalogsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *CatalogsClient) getCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CatalogsClient) getHandleResponse(resp *http.Response) (CatalogsClientGetResponse, error) {
	result := CatalogsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Catalog); err != nil {
		return CatalogsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - List Catalog resources by resource group
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - CatalogsClientListByResourceGroupOptions contains the optional parameters for the CatalogsClient.NewListByResourceGroupPager
//     method.
func (client *CatalogsClient) NewListByResourceGroupPager(resourceGroupName string, options *CatalogsClientListByResourceGroupOptions) (*runtime.Pager[CatalogsClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[CatalogsClientListByResourceGroupResponse]{
		More: func(page CatalogsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CatalogsClientListByResourceGroupResponse) (CatalogsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CatalogsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CatalogsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CatalogsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *CatalogsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *CatalogsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs"
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
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *CatalogsClient) listByResourceGroupHandleResponse(resp *http.Response) (CatalogsClientListByResourceGroupResponse, error) {
	result := CatalogsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CatalogListResult); err != nil {
		return CatalogsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - List Catalog resources by subscription ID
//
// Generated from API version 2022-09-01-preview
//   - options - CatalogsClientListBySubscriptionOptions contains the optional parameters for the CatalogsClient.NewListBySubscriptionPager
//     method.
func (client *CatalogsClient) NewListBySubscriptionPager(options *CatalogsClientListBySubscriptionOptions) (*runtime.Pager[CatalogsClientListBySubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[CatalogsClientListBySubscriptionResponse]{
		More: func(page CatalogsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CatalogsClientListBySubscriptionResponse) (CatalogsClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CatalogsClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CatalogsClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CatalogsClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *CatalogsClient) listBySubscriptionCreateRequest(ctx context.Context, options *CatalogsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.AzureSphere/catalogs"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *CatalogsClient) listBySubscriptionHandleResponse(resp *http.Response) (CatalogsClientListBySubscriptionResponse, error) {
	result := CatalogsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CatalogListResult); err != nil {
		return CatalogsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// NewListDeploymentsPager - Lists deployments for catalog.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - options - CatalogsClientListDeploymentsOptions contains the optional parameters for the CatalogsClient.NewListDeploymentsPager
//     method.
func (client *CatalogsClient) NewListDeploymentsPager(resourceGroupName string, catalogName string, options *CatalogsClientListDeploymentsOptions) (*runtime.Pager[CatalogsClientListDeploymentsResponse]) {
	return runtime.NewPager(runtime.PagingHandler[CatalogsClientListDeploymentsResponse]{
		More: func(page CatalogsClientListDeploymentsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CatalogsClientListDeploymentsResponse) (CatalogsClientListDeploymentsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listDeploymentsCreateRequest(ctx, resourceGroupName, catalogName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CatalogsClientListDeploymentsResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CatalogsClientListDeploymentsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CatalogsClientListDeploymentsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listDeploymentsHandleResponse(resp)
		},
	})
}

// listDeploymentsCreateRequest creates the ListDeployments request.
func (client *CatalogsClient) listDeploymentsCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientListDeploymentsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}/listDeployments"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Maxpagesize != nil {
		reqQP.Set("$maxpagesize", strconv.FormatInt(int64(*options.Maxpagesize), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listDeploymentsHandleResponse handles the ListDeployments response.
func (client *CatalogsClient) listDeploymentsHandleResponse(resp *http.Response) (CatalogsClientListDeploymentsResponse, error) {
	result := CatalogsClientListDeploymentsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeploymentListResult); err != nil {
		return CatalogsClientListDeploymentsResponse{}, err
	}
	return result, nil
}

// NewListDeviceGroupsPager - List the device groups for the catalog.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - listDeviceGroupsRequest - List device groups for catalog.
//   - options - CatalogsClientListDeviceGroupsOptions contains the optional parameters for the CatalogsClient.NewListDeviceGroupsPager
//     method.
func (client *CatalogsClient) NewListDeviceGroupsPager(resourceGroupName string, catalogName string, listDeviceGroupsRequest ListDeviceGroupsRequest, options *CatalogsClientListDeviceGroupsOptions) (*runtime.Pager[CatalogsClientListDeviceGroupsResponse]) {
	return runtime.NewPager(runtime.PagingHandler[CatalogsClientListDeviceGroupsResponse]{
		More: func(page CatalogsClientListDeviceGroupsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CatalogsClientListDeviceGroupsResponse) (CatalogsClientListDeviceGroupsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listDeviceGroupsCreateRequest(ctx, resourceGroupName, catalogName, listDeviceGroupsRequest, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CatalogsClientListDeviceGroupsResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CatalogsClientListDeviceGroupsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CatalogsClientListDeviceGroupsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listDeviceGroupsHandleResponse(resp)
		},
	})
}

// listDeviceGroupsCreateRequest creates the ListDeviceGroups request.
func (client *CatalogsClient) listDeviceGroupsCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, listDeviceGroupsRequest ListDeviceGroupsRequest, options *CatalogsClientListDeviceGroupsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}/listDeviceGroups"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Maxpagesize != nil {
		reqQP.Set("$maxpagesize", strconv.FormatInt(int64(*options.Maxpagesize), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, listDeviceGroupsRequest); err != nil {
	return nil, err
}
	return req, nil
}

// listDeviceGroupsHandleResponse handles the ListDeviceGroups response.
func (client *CatalogsClient) listDeviceGroupsHandleResponse(resp *http.Response) (CatalogsClientListDeviceGroupsResponse, error) {
	result := CatalogsClientListDeviceGroupsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeviceGroupListResult); err != nil {
		return CatalogsClientListDeviceGroupsResponse{}, err
	}
	return result, nil
}

// NewListDeviceInsightsPager - Lists device insights for catalog.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - options - CatalogsClientListDeviceInsightsOptions contains the optional parameters for the CatalogsClient.NewListDeviceInsightsPager
//     method.
func (client *CatalogsClient) NewListDeviceInsightsPager(resourceGroupName string, catalogName string, options *CatalogsClientListDeviceInsightsOptions) (*runtime.Pager[CatalogsClientListDeviceInsightsResponse]) {
	return runtime.NewPager(runtime.PagingHandler[CatalogsClientListDeviceInsightsResponse]{
		More: func(page CatalogsClientListDeviceInsightsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CatalogsClientListDeviceInsightsResponse) (CatalogsClientListDeviceInsightsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listDeviceInsightsCreateRequest(ctx, resourceGroupName, catalogName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CatalogsClientListDeviceInsightsResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CatalogsClientListDeviceInsightsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CatalogsClientListDeviceInsightsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listDeviceInsightsHandleResponse(resp)
		},
	})
}

// listDeviceInsightsCreateRequest creates the ListDeviceInsights request.
func (client *CatalogsClient) listDeviceInsightsCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientListDeviceInsightsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}/listDeviceInsights"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Maxpagesize != nil {
		reqQP.Set("$maxpagesize", strconv.FormatInt(int64(*options.Maxpagesize), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listDeviceInsightsHandleResponse handles the ListDeviceInsights response.
func (client *CatalogsClient) listDeviceInsightsHandleResponse(resp *http.Response) (CatalogsClientListDeviceInsightsResponse, error) {
	result := CatalogsClientListDeviceInsightsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PagedDeviceInsight); err != nil {
		return CatalogsClientListDeviceInsightsResponse{}, err
	}
	return result, nil
}

// NewListDevicesPager - Lists devices for catalog.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - options - CatalogsClientListDevicesOptions contains the optional parameters for the CatalogsClient.NewListDevicesPager
//     method.
func (client *CatalogsClient) NewListDevicesPager(resourceGroupName string, catalogName string, options *CatalogsClientListDevicesOptions) (*runtime.Pager[CatalogsClientListDevicesResponse]) {
	return runtime.NewPager(runtime.PagingHandler[CatalogsClientListDevicesResponse]{
		More: func(page CatalogsClientListDevicesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CatalogsClientListDevicesResponse) (CatalogsClientListDevicesResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listDevicesCreateRequest(ctx, resourceGroupName, catalogName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CatalogsClientListDevicesResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CatalogsClientListDevicesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CatalogsClientListDevicesResponse{}, runtime.NewResponseError(resp)
			}
			return client.listDevicesHandleResponse(resp)
		},
	})
}

// listDevicesCreateRequest creates the ListDevices request.
func (client *CatalogsClient) listDevicesCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, options *CatalogsClientListDevicesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}/listDevices"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Maxpagesize != nil {
		reqQP.Set("$maxpagesize", strconv.FormatInt(int64(*options.Maxpagesize), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listDevicesHandleResponse handles the ListDevices response.
func (client *CatalogsClient) listDevicesHandleResponse(resp *http.Response) (CatalogsClientListDevicesResponse, error) {
	result := CatalogsClientListDevicesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeviceListResult); err != nil {
		return CatalogsClientListDevicesResponse{}, err
	}
	return result, nil
}

// Update - Update a Catalog
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - catalogName - Name of catalog
//   - properties - The resource properties to be updated.
//   - options - CatalogsClientUpdateOptions contains the optional parameters for the CatalogsClient.Update method.
func (client *CatalogsClient) Update(ctx context.Context, resourceGroupName string, catalogName string, properties CatalogUpdate, options *CatalogsClientUpdateOptions) (CatalogsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, catalogName, properties, options)
	if err != nil {
		return CatalogsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CatalogsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CatalogsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *CatalogsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, catalogName string, properties CatalogUpdate, options *CatalogsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureSphere/catalogs/{catalogName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if catalogName == "" {
		return nil, errors.New("parameter catalogName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{catalogName}", url.PathEscape(catalogName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, properties); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *CatalogsClient) updateHandleResponse(resp *http.Response) (CatalogsClientUpdateResponse, error) {
	result := CatalogsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Catalog); err != nil {
		return CatalogsClientUpdateResponse{}, err
	}
	return result, nil
}

