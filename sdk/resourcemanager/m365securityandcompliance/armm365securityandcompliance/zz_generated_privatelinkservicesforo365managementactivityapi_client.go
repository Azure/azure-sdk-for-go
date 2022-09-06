//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armm365securityandcompliance

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

// PrivateLinkServicesForO365ManagementActivityAPIClient contains the methods for the PrivateLinkServicesForO365ManagementActivityAPI group.
// Don't use this type directly, use NewPrivateLinkServicesForO365ManagementActivityAPIClient() instead.
type PrivateLinkServicesForO365ManagementActivityAPIClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewPrivateLinkServicesForO365ManagementActivityAPIClient creates a new instance of PrivateLinkServicesForO365ManagementActivityAPIClient with the specified values.
// subscriptionID - The subscription identifier.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewPrivateLinkServicesForO365ManagementActivityAPIClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PrivateLinkServicesForO365ManagementActivityAPIClient, error) {
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
	client := &PrivateLinkServicesForO365ManagementActivityAPIClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update the metadata of a privateLinkServicesForO365ManagementActivityAPI instance.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
// resourceGroupName - The name of the resource group that contains the service instance.
// resourceName - The name of the service instance.
// privateLinkServicesForO365ManagementActivityAPIDescription - The service instance metadata.
// options - PrivateLinkServicesForO365ManagementActivityAPIClientBeginCreateOrUpdateOptions contains the optional parameters
// for the PrivateLinkServicesForO365ManagementActivityAPIClient.BeginCreateOrUpdate method.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, privateLinkServicesForO365ManagementActivityAPIDescription PrivateLinkServicesForO365ManagementActivityAPIDescription, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginCreateOrUpdateOptions) (*runtime.Poller[PrivateLinkServicesForO365ManagementActivityAPIClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, resourceName, privateLinkServicesForO365ManagementActivityAPIDescription, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[PrivateLinkServicesForO365ManagementActivityAPIClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[PrivateLinkServicesForO365ManagementActivityAPIClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - Create or update the metadata of a privateLinkServicesForO365ManagementActivityAPI instance.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) createOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, privateLinkServicesForO365ManagementActivityAPIDescription PrivateLinkServicesForO365ManagementActivityAPIDescription, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, resourceName, privateLinkServicesForO365ManagementActivityAPIDescription, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, privateLinkServicesForO365ManagementActivityAPIDescription PrivateLinkServicesForO365ManagementActivityAPIDescription, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.M365SecurityAndCompliance/privateLinkServicesForO365ManagementActivityAPI/{resourceName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-25-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, privateLinkServicesForO365ManagementActivityAPIDescription)
}

// BeginDelete - Delete a service instance.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
// resourceGroupName - The name of the resource group that contains the service instance.
// resourceName - The name of the service instance.
// options - PrivateLinkServicesForO365ManagementActivityAPIClientBeginDeleteOptions contains the optional parameters for
// the PrivateLinkServicesForO365ManagementActivityAPIClient.BeginDelete method.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) BeginDelete(ctx context.Context, resourceGroupName string, resourceName string, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginDeleteOptions) (*runtime.Poller[PrivateLinkServicesForO365ManagementActivityAPIClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, resourceName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[PrivateLinkServicesForO365ManagementActivityAPIClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[PrivateLinkServicesForO365ManagementActivityAPIClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Delete a service instance.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) deleteOperation(ctx context.Context, resourceGroupName string, resourceName string, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.M365SecurityAndCompliance/privateLinkServicesForO365ManagementActivityAPI/{resourceName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-25-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get the metadata of a privateLinkServicesForO365ManagementActivityAPI resource.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
// resourceGroupName - The name of the resource group that contains the service instance.
// resourceName - The name of the service instance.
// options - PrivateLinkServicesForO365ManagementActivityAPIClientGetOptions contains the optional parameters for the PrivateLinkServicesForO365ManagementActivityAPIClient.Get
// method.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) Get(ctx context.Context, resourceGroupName string, resourceName string, options *PrivateLinkServicesForO365ManagementActivityAPIClientGetOptions) (PrivateLinkServicesForO365ManagementActivityAPIClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return PrivateLinkServicesForO365ManagementActivityAPIClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PrivateLinkServicesForO365ManagementActivityAPIClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PrivateLinkServicesForO365ManagementActivityAPIClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *PrivateLinkServicesForO365ManagementActivityAPIClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.M365SecurityAndCompliance/privateLinkServicesForO365ManagementActivityAPI/{resourceName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-25-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) getHandleResponse(resp *http.Response) (PrivateLinkServicesForO365ManagementActivityAPIClientGetResponse, error) {
	result := PrivateLinkServicesForO365ManagementActivityAPIClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateLinkServicesForO365ManagementActivityAPIDescription); err != nil {
		return PrivateLinkServicesForO365ManagementActivityAPIClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Get all the privateLinkServicesForO365ManagementActivityAPI instances in a subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
// options - PrivateLinkServicesForO365ManagementActivityAPIClientListOptions contains the optional parameters for the PrivateLinkServicesForO365ManagementActivityAPIClient.List
// method.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) NewListPager(options *PrivateLinkServicesForO365ManagementActivityAPIClientListOptions) *runtime.Pager[PrivateLinkServicesForO365ManagementActivityAPIClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[PrivateLinkServicesForO365ManagementActivityAPIClientListResponse]{
		More: func(page PrivateLinkServicesForO365ManagementActivityAPIClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PrivateLinkServicesForO365ManagementActivityAPIClientListResponse) (PrivateLinkServicesForO365ManagementActivityAPIClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PrivateLinkServicesForO365ManagementActivityAPIClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return PrivateLinkServicesForO365ManagementActivityAPIClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PrivateLinkServicesForO365ManagementActivityAPIClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) listCreateRequest(ctx context.Context, options *PrivateLinkServicesForO365ManagementActivityAPIClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.M365SecurityAndCompliance/privateLinkServicesForO365ManagementActivityAPI"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-25-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) listHandleResponse(resp *http.Response) (PrivateLinkServicesForO365ManagementActivityAPIClientListResponse, error) {
	result := PrivateLinkServicesForO365ManagementActivityAPIClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateLinkServicesForO365ManagementActivityAPIDescriptionListResult); err != nil {
		return PrivateLinkServicesForO365ManagementActivityAPIClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Get all the service instances in a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
// resourceGroupName - The name of the resource group that contains the service instance.
// options - PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupOptions contains the optional parameters
// for the PrivateLinkServicesForO365ManagementActivityAPIClient.ListByResourceGroup method.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) NewListByResourceGroupPager(resourceGroupName string, options *PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupOptions) *runtime.Pager[PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse]{
		More: func(page PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse) (PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.M365SecurityAndCompliance/privateLinkServicesForO365ManagementActivityAPI"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-25-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) listByResourceGroupHandleResponse(resp *http.Response) (PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse, error) {
	result := PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateLinkServicesForO365ManagementActivityAPIDescriptionListResult); err != nil {
		return PrivateLinkServicesForO365ManagementActivityAPIClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update the metadata of a privateLinkServicesForO365ManagementActivityAPI instance.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
// resourceGroupName - The name of the resource group that contains the service instance.
// resourceName - The name of the service instance.
// servicePatchDescription - The service instance metadata and security metadata.
// options - PrivateLinkServicesForO365ManagementActivityAPIClientBeginUpdateOptions contains the optional parameters for
// the PrivateLinkServicesForO365ManagementActivityAPIClient.BeginUpdate method.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) BeginUpdate(ctx context.Context, resourceGroupName string, resourceName string, servicePatchDescription ServicesPatchDescription, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginUpdateOptions) (*runtime.Poller[PrivateLinkServicesForO365ManagementActivityAPIClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, resourceName, servicePatchDescription, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[PrivateLinkServicesForO365ManagementActivityAPIClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[PrivateLinkServicesForO365ManagementActivityAPIClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - Update the metadata of a privateLinkServicesForO365ManagementActivityAPI instance.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-03-25-preview
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) update(ctx context.Context, resourceGroupName string, resourceName string, servicePatchDescription ServicesPatchDescription, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, servicePatchDescription, options)
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

// updateCreateRequest creates the Update request.
func (client *PrivateLinkServicesForO365ManagementActivityAPIClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, servicePatchDescription ServicesPatchDescription, options *PrivateLinkServicesForO365ManagementActivityAPIClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.M365SecurityAndCompliance/privateLinkServicesForO365ManagementActivityAPI/{resourceName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-25-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, servicePatchDescription)
}
