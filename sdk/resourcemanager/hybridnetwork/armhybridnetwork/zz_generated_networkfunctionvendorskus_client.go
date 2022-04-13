//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridnetwork

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

// NetworkFunctionVendorSKUsClient contains the methods for the NetworkFunctionVendorSKUs group.
// Don't use this type directly, use NewNetworkFunctionVendorSKUsClient() instead.
type NetworkFunctionVendorSKUsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewNetworkFunctionVendorSKUsClient creates a new instance of NetworkFunctionVendorSKUsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewNetworkFunctionVendorSKUsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*NetworkFunctionVendorSKUsClient, error) {
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
	client := &NetworkFunctionVendorSKUsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// ListBySKU - Lists information about network function vendor sku details.
// If the operation fails it returns an *azcore.ResponseError type.
// vendorName - The name of the network function vendor.
// vendorSKUName - The name of the network function sku.
// options - NetworkFunctionVendorSKUsClientListBySKUOptions contains the optional parameters for the NetworkFunctionVendorSKUsClient.ListBySKU
// method.
func (client *NetworkFunctionVendorSKUsClient) ListBySKU(vendorName string, vendorSKUName string, options *NetworkFunctionVendorSKUsClientListBySKUOptions) *runtime.Pager[NetworkFunctionVendorSKUsClientListBySKUResponse] {
	return runtime.NewPager(runtime.PageProcessor[NetworkFunctionVendorSKUsClientListBySKUResponse]{
		More: func(page NetworkFunctionVendorSKUsClientListBySKUResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *NetworkFunctionVendorSKUsClientListBySKUResponse) (NetworkFunctionVendorSKUsClientListBySKUResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySKUCreateRequest(ctx, vendorName, vendorSKUName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return NetworkFunctionVendorSKUsClientListBySKUResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return NetworkFunctionVendorSKUsClientListBySKUResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return NetworkFunctionVendorSKUsClientListBySKUResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySKUHandleResponse(resp)
		},
	})
}

// listBySKUCreateRequest creates the ListBySKU request.
func (client *NetworkFunctionVendorSKUsClient) listBySKUCreateRequest(ctx context.Context, vendorName string, vendorSKUName string, options *NetworkFunctionVendorSKUsClientListBySKUOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/networkFunctionVendors/{vendorName}/vendorSkus/{vendorSkuName}"
	if vendorName == "" {
		return nil, errors.New("parameter vendorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vendorName}", url.PathEscape(vendorName))
	if vendorSKUName == "" {
		return nil, errors.New("parameter vendorSKUName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vendorSkuName}", url.PathEscape(vendorSKUName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listBySKUHandleResponse handles the ListBySKU response.
func (client *NetworkFunctionVendorSKUsClient) listBySKUHandleResponse(resp *http.Response) (NetworkFunctionVendorSKUsClientListBySKUResponse, error) {
	result := NetworkFunctionVendorSKUsClientListBySKUResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFunctionSKUDetails); err != nil {
		return NetworkFunctionVendorSKUsClientListBySKUResponse{}, err
	}
	return result, nil
}

// ListByVendor - Lists all network function vendor sku details in a vendor.
// If the operation fails it returns an *azcore.ResponseError type.
// vendorName - The name of the network function vendor.
// options - NetworkFunctionVendorSKUsClientListByVendorOptions contains the optional parameters for the NetworkFunctionVendorSKUsClient.ListByVendor
// method.
func (client *NetworkFunctionVendorSKUsClient) ListByVendor(vendorName string, options *NetworkFunctionVendorSKUsClientListByVendorOptions) *runtime.Pager[NetworkFunctionVendorSKUsClientListByVendorResponse] {
	return runtime.NewPager(runtime.PageProcessor[NetworkFunctionVendorSKUsClientListByVendorResponse]{
		More: func(page NetworkFunctionVendorSKUsClientListByVendorResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *NetworkFunctionVendorSKUsClientListByVendorResponse) (NetworkFunctionVendorSKUsClientListByVendorResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByVendorCreateRequest(ctx, vendorName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return NetworkFunctionVendorSKUsClientListByVendorResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return NetworkFunctionVendorSKUsClientListByVendorResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return NetworkFunctionVendorSKUsClientListByVendorResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByVendorHandleResponse(resp)
		},
	})
}

// listByVendorCreateRequest creates the ListByVendor request.
func (client *NetworkFunctionVendorSKUsClient) listByVendorCreateRequest(ctx context.Context, vendorName string, options *NetworkFunctionVendorSKUsClientListByVendorOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/networkFunctionVendors/{vendorName}/vendorSkus"
	if vendorName == "" {
		return nil, errors.New("parameter vendorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vendorName}", url.PathEscape(vendorName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByVendorHandleResponse handles the ListByVendor response.
func (client *NetworkFunctionVendorSKUsClient) listByVendorHandleResponse(resp *http.Response) (NetworkFunctionVendorSKUsClientListByVendorResponse, error) {
	result := NetworkFunctionVendorSKUsClientListByVendorResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFunctionSKUListResult); err != nil {
		return NetworkFunctionVendorSKUsClientListByVendorResponse{}, err
	}
	return result, nil
}
