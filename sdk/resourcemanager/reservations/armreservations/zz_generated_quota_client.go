//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armreservations

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

// QuotaClient contains the methods for the Quota group.
// Don't use this type directly, use NewQuotaClient() instead.
type QuotaClient struct {
	host string
	pl   runtime.Pipeline
}

// NewQuotaClient creates a new instance of QuotaClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewQuotaClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*QuotaClient, error) {
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
	client := &QuotaClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update the quota (service limits) of a resource to the requested value. Steps:
// 1. Make the Get request to get the quota information for specific resource.
// 2. To increase the quota, update the limit field in the response from Get request to new value.
// 3. Submit the JSON to the quota request API to update the quota. The Create quota request may be constructed as follows.
// The PUT operation can be used to update the quota.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-25
// subscriptionID - Azure subscription ID.
// providerID - Azure resource provider ID.
// location - Azure region.
// resourceName - The resource name for a resource provider, such as SKU name for Microsoft.Compute, Sku or TotalLowPriorityCores
// for Microsoft.MachineLearningServices
// createQuotaRequest - Quota requests payload.
// options - QuotaClientBeginCreateOrUpdateOptions contains the optional parameters for the QuotaClient.BeginCreateOrUpdate
// method.
func (client *QuotaClient) BeginCreateOrUpdate(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginCreateOrUpdateOptions) (*runtime.Poller[QuotaClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, subscriptionID, providerID, location, resourceName, createQuotaRequest, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[QuotaClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaOriginalURI,
		})
	} else {
		return runtime.NewPollerFromResumeToken[QuotaClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - Create or update the quota (service limits) of a resource to the requested value. Steps:
// 1. Make the Get request to get the quota information for specific resource.
// 2. To increase the quota, update the limit field in the response from Get request to new value.
// 3. Submit the JSON to the quota request API to update the quota. The Create quota request may be constructed as follows.
// The PUT operation can be used to update the quota.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-25
func (client *QuotaClient) createOrUpdate(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, subscriptionID, providerID, location, resourceName, createQuotaRequest, options)
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
func (client *QuotaClient) createOrUpdateCreateRequest(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/resourceProviders/{providerId}/locations/{location}/serviceLimits/{resourceName}"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	if providerID == "" {
		return nil, errors.New("parameter providerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{providerId}", url.PathEscape(providerID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-25")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, createQuotaRequest)
}

// Get - Get the current quota (service limit) and usage of a resource. You can use the response from the GET operation to
// submit quota update request.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-25
// subscriptionID - Azure subscription ID.
// providerID - Azure resource provider ID.
// location - Azure region.
// resourceName - The resource name for a resource provider, such as SKU name for Microsoft.Compute, Sku or TotalLowPriorityCores
// for Microsoft.MachineLearningServices
// options - QuotaClientGetOptions contains the optional parameters for the QuotaClient.Get method.
func (client *QuotaClient) Get(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, options *QuotaClientGetOptions) (QuotaClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, subscriptionID, providerID, location, resourceName, options)
	if err != nil {
		return QuotaClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return QuotaClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return QuotaClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *QuotaClient) getCreateRequest(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, options *QuotaClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/resourceProviders/{providerId}/locations/{location}/serviceLimits/{resourceName}"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	if providerID == "" {
		return nil, errors.New("parameter providerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{providerId}", url.PathEscape(providerID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-25")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *QuotaClient) getHandleResponse(resp *http.Response) (QuotaClientGetResponse, error) {
	result := QuotaClientGetResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.CurrentQuotaLimitBase); err != nil {
		return QuotaClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of current quotas (service limits) and usage for all resources. The response from the list quota
// operation can be leveraged to request quota updates.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-25
// subscriptionID - Azure subscription ID.
// providerID - Azure resource provider ID.
// location - Azure region.
// options - QuotaClientListOptions contains the optional parameters for the QuotaClient.List method.
func (client *QuotaClient) NewListPager(subscriptionID string, providerID string, location string, options *QuotaClientListOptions) *runtime.Pager[QuotaClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[QuotaClientListResponse]{
		More: func(page QuotaClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *QuotaClientListResponse) (QuotaClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, subscriptionID, providerID, location, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return QuotaClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return QuotaClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return QuotaClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *QuotaClient) listCreateRequest(ctx context.Context, subscriptionID string, providerID string, location string, options *QuotaClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/resourceProviders/{providerId}/locations/{location}/serviceLimits"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	if providerID == "" {
		return nil, errors.New("parameter providerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{providerId}", url.PathEscape(providerID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-25")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *QuotaClient) listHandleResponse(resp *http.Response) (QuotaClientListResponse, error) {
	result := QuotaClientListResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.QuotaLimits); err != nil {
		return QuotaClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update the quota (service limits) of this resource to the requested value. • To get the quota information
// for specific resource, send a GET request. • To increase the quota, update the limit field
// from the GET response to a new value. • To update the quota value, submit the JSON response to the quota request API to
// update the quota. • To update the quota. use the PATCH operation.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-25
// subscriptionID - Azure subscription ID.
// providerID - Azure resource provider ID.
// location - Azure region.
// resourceName - The resource name for a resource provider, such as SKU name for Microsoft.Compute, Sku or TotalLowPriorityCores
// for Microsoft.MachineLearningServices
// createQuotaRequest - Payload for the quota request.
// options - QuotaClientBeginUpdateOptions contains the optional parameters for the QuotaClient.BeginUpdate method.
func (client *QuotaClient) BeginUpdate(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginUpdateOptions) (*runtime.Poller[QuotaClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, subscriptionID, providerID, location, resourceName, createQuotaRequest, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[QuotaClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaOriginalURI,
		})
	} else {
		return runtime.NewPollerFromResumeToken[QuotaClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - Update the quota (service limits) of this resource to the requested value. • To get the quota information for
// specific resource, send a GET request. • To increase the quota, update the limit field
// from the GET response to a new value. • To update the quota value, submit the JSON response to the quota request API to
// update the quota. • To update the quota. use the PATCH operation.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-25
func (client *QuotaClient) update(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, subscriptionID, providerID, location, resourceName, createQuotaRequest, options)
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

// updateCreateRequest creates the Update request.
func (client *QuotaClient) updateCreateRequest(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/resourceProviders/{providerId}/locations/{location}/serviceLimits/{resourceName}"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	if providerID == "" {
		return nil, errors.New("parameter providerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{providerId}", url.PathEscape(providerID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-25")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, createQuotaRequest)
}
