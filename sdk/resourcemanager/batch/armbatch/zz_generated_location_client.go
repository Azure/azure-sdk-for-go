//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armbatch

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

// LocationClient contains the methods for the Location group.
// Don't use this type directly, use NewLocationClient() instead.
type LocationClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewLocationClient creates a new instance of LocationClient with the specified values.
// subscriptionID - The Azure subscription ID. This is a GUID-formatted string (e.g. 00000000-0000-0000-0000-000000000000)
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewLocationClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LocationClient, error) {
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
	client := &LocationClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CheckNameAvailability - Checks whether the Batch account name is available in the specified region.
// If the operation fails it returns an *azcore.ResponseError type.
// locationName - The desired region for the name check.
// parameters - Properties needed to check the availability of a name.
// options - LocationClientCheckNameAvailabilityOptions contains the optional parameters for the LocationClient.CheckNameAvailability
// method.
func (client *LocationClient) CheckNameAvailability(ctx context.Context, locationName string, parameters CheckNameAvailabilityParameters, options *LocationClientCheckNameAvailabilityOptions) (LocationClientCheckNameAvailabilityResponse, error) {
	req, err := client.checkNameAvailabilityCreateRequest(ctx, locationName, parameters, options)
	if err != nil {
		return LocationClientCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LocationClientCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LocationClientCheckNameAvailabilityResponse{}, runtime.NewResponseError(resp)
	}
	return client.checkNameAvailabilityHandleResponse(resp)
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *LocationClient) checkNameAvailabilityCreateRequest(ctx context.Context, locationName string, parameters CheckNameAvailabilityParameters, options *LocationClientCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Batch/locations/{locationName}/checkNameAvailability"
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *LocationClient) checkNameAvailabilityHandleResponse(resp *http.Response) (LocationClientCheckNameAvailabilityResponse, error) {
	result := LocationClientCheckNameAvailabilityResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckNameAvailabilityResult); err != nil {
		return LocationClientCheckNameAvailabilityResponse{}, err
	}
	return result, nil
}

// GetQuotas - Gets the Batch service quotas for the specified subscription at the given location.
// If the operation fails it returns an *azcore.ResponseError type.
// locationName - The region for which to retrieve Batch service quotas.
// options - LocationClientGetQuotasOptions contains the optional parameters for the LocationClient.GetQuotas method.
func (client *LocationClient) GetQuotas(ctx context.Context, locationName string, options *LocationClientGetQuotasOptions) (LocationClientGetQuotasResponse, error) {
	req, err := client.getQuotasCreateRequest(ctx, locationName, options)
	if err != nil {
		return LocationClientGetQuotasResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LocationClientGetQuotasResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LocationClientGetQuotasResponse{}, runtime.NewResponseError(resp)
	}
	return client.getQuotasHandleResponse(resp)
}

// getQuotasCreateRequest creates the GetQuotas request.
func (client *LocationClient) getQuotasCreateRequest(ctx context.Context, locationName string, options *LocationClientGetQuotasOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Batch/locations/{locationName}/quotas"
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getQuotasHandleResponse handles the GetQuotas response.
func (client *LocationClient) getQuotasHandleResponse(resp *http.Response) (LocationClientGetQuotasResponse, error) {
	result := LocationClientGetQuotasResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocationQuota); err != nil {
		return LocationClientGetQuotasResponse{}, err
	}
	return result, nil
}

// ListSupportedCloudServiceSKUs - Gets the list of Batch supported Cloud Service VM sizes available at the given location.
// If the operation fails it returns an *azcore.ResponseError type.
// locationName - The region for which to retrieve Batch service supported SKUs.
// options - LocationClientListSupportedCloudServiceSKUsOptions contains the optional parameters for the LocationClient.ListSupportedCloudServiceSKUs
// method.
func (client *LocationClient) ListSupportedCloudServiceSKUs(locationName string, options *LocationClientListSupportedCloudServiceSKUsOptions) *runtime.Pager[LocationClientListSupportedCloudServiceSKUsResponse] {
	return runtime.NewPager(runtime.PageProcessor[LocationClientListSupportedCloudServiceSKUsResponse]{
		More: func(page LocationClientListSupportedCloudServiceSKUsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *LocationClientListSupportedCloudServiceSKUsResponse) (LocationClientListSupportedCloudServiceSKUsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listSupportedCloudServiceSKUsCreateRequest(ctx, locationName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return LocationClientListSupportedCloudServiceSKUsResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return LocationClientListSupportedCloudServiceSKUsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LocationClientListSupportedCloudServiceSKUsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listSupportedCloudServiceSKUsHandleResponse(resp)
		},
	})
}

// listSupportedCloudServiceSKUsCreateRequest creates the ListSupportedCloudServiceSKUs request.
func (client *LocationClient) listSupportedCloudServiceSKUsCreateRequest(ctx context.Context, locationName string, options *LocationClientListSupportedCloudServiceSKUsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Batch/locations/{locationName}/cloudServiceSkus"
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Maxresults != nil {
		reqQP.Set("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	reqQP.Set("api-version", "2022-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listSupportedCloudServiceSKUsHandleResponse handles the ListSupportedCloudServiceSKUs response.
func (client *LocationClient) listSupportedCloudServiceSKUsHandleResponse(resp *http.Response) (LocationClientListSupportedCloudServiceSKUsResponse, error) {
	result := LocationClientListSupportedCloudServiceSKUsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SupportedSKUsResult); err != nil {
		return LocationClientListSupportedCloudServiceSKUsResponse{}, err
	}
	return result, nil
}

// ListSupportedVirtualMachineSKUs - Gets the list of Batch supported Virtual Machine VM sizes available at the given location.
// If the operation fails it returns an *azcore.ResponseError type.
// locationName - The region for which to retrieve Batch service supported SKUs.
// options - LocationClientListSupportedVirtualMachineSKUsOptions contains the optional parameters for the LocationClient.ListSupportedVirtualMachineSKUs
// method.
func (client *LocationClient) ListSupportedVirtualMachineSKUs(locationName string, options *LocationClientListSupportedVirtualMachineSKUsOptions) *runtime.Pager[LocationClientListSupportedVirtualMachineSKUsResponse] {
	return runtime.NewPager(runtime.PageProcessor[LocationClientListSupportedVirtualMachineSKUsResponse]{
		More: func(page LocationClientListSupportedVirtualMachineSKUsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *LocationClientListSupportedVirtualMachineSKUsResponse) (LocationClientListSupportedVirtualMachineSKUsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listSupportedVirtualMachineSKUsCreateRequest(ctx, locationName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return LocationClientListSupportedVirtualMachineSKUsResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return LocationClientListSupportedVirtualMachineSKUsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LocationClientListSupportedVirtualMachineSKUsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listSupportedVirtualMachineSKUsHandleResponse(resp)
		},
	})
}

// listSupportedVirtualMachineSKUsCreateRequest creates the ListSupportedVirtualMachineSKUs request.
func (client *LocationClient) listSupportedVirtualMachineSKUsCreateRequest(ctx context.Context, locationName string, options *LocationClientListSupportedVirtualMachineSKUsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Batch/locations/{locationName}/virtualMachineSkus"
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Maxresults != nil {
		reqQP.Set("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	reqQP.Set("api-version", "2022-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listSupportedVirtualMachineSKUsHandleResponse handles the ListSupportedVirtualMachineSKUs response.
func (client *LocationClient) listSupportedVirtualMachineSKUsHandleResponse(resp *http.Response) (LocationClientListSupportedVirtualMachineSKUsResponse, error) {
	result := LocationClientListSupportedVirtualMachineSKUsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SupportedSKUsResult); err != nil {
		return LocationClientListSupportedVirtualMachineSKUsResponse{}, err
	}
	return result, nil
}
