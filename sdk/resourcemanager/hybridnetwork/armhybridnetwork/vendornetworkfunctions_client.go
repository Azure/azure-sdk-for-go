//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridnetwork

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

// VendorNetworkFunctionsClient contains the methods for the VendorNetworkFunctions group.
// Don't use this type directly, use NewVendorNetworkFunctionsClient() instead.
type VendorNetworkFunctionsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewVendorNetworkFunctionsClient creates a new instance of VendorNetworkFunctionsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewVendorNetworkFunctionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VendorNetworkFunctionsClient, error) {
	cl, err := arm.NewClient(moduleName+".VendorNetworkFunctionsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &VendorNetworkFunctionsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a vendor network function. This operation can take up to 6 hours to complete.
// This is expected service behavior.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-05-01
//   - locationName - The Azure region where the network function resource was created by the customer.
//   - vendorName - The name of the vendor.
//   - serviceKey - The GUID for the vendor network function.
//   - parameters - Parameters supplied to the create or update vendor network function operation.
//   - options - VendorNetworkFunctionsClientBeginCreateOrUpdateOptions contains the optional parameters for the VendorNetworkFunctionsClient.BeginCreateOrUpdate
//     method.
func (client *VendorNetworkFunctionsClient) BeginCreateOrUpdate(ctx context.Context, locationName string, vendorName string, serviceKey string, parameters VendorNetworkFunction, options *VendorNetworkFunctionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[VendorNetworkFunctionsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, locationName, vendorName, serviceKey, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[VendorNetworkFunctionsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[VendorNetworkFunctionsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a vendor network function. This operation can take up to 6 hours to complete. This
// is expected service behavior.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-05-01
func (client *VendorNetworkFunctionsClient) createOrUpdate(ctx context.Context, locationName string, vendorName string, serviceKey string, parameters VendorNetworkFunction, options *VendorNetworkFunctionsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, locationName, vendorName, serviceKey, parameters, options)
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
func (client *VendorNetworkFunctionsClient) createOrUpdateCreateRequest(ctx context.Context, locationName string, vendorName string, serviceKey string, parameters VendorNetworkFunction, options *VendorNetworkFunctionsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/locations/{locationName}/vendors/{vendorName}/networkFunctions/{serviceKey}"
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if vendorName == "" {
		return nil, errors.New("parameter vendorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vendorName}", url.PathEscape(vendorName))
	if serviceKey == "" {
		return nil, errors.New("parameter serviceKey cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceKey}", url.PathEscape(serviceKey))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// Get - Gets information about the specified vendor network function.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-05-01
//   - locationName - The Azure region where the network function resource was created by the customer.
//   - vendorName - The name of the vendor.
//   - serviceKey - The GUID for the vendor network function.
//   - options - VendorNetworkFunctionsClientGetOptions contains the optional parameters for the VendorNetworkFunctionsClient.Get
//     method.
func (client *VendorNetworkFunctionsClient) Get(ctx context.Context, locationName string, vendorName string, serviceKey string, options *VendorNetworkFunctionsClientGetOptions) (VendorNetworkFunctionsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, locationName, vendorName, serviceKey, options)
	if err != nil {
		return VendorNetworkFunctionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return VendorNetworkFunctionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return VendorNetworkFunctionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *VendorNetworkFunctionsClient) getCreateRequest(ctx context.Context, locationName string, vendorName string, serviceKey string, options *VendorNetworkFunctionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/locations/{locationName}/vendors/{vendorName}/networkFunctions/{serviceKey}"
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if vendorName == "" {
		return nil, errors.New("parameter vendorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vendorName}", url.PathEscape(vendorName))
	if serviceKey == "" {
		return nil, errors.New("parameter serviceKey cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceKey}", url.PathEscape(serviceKey))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VendorNetworkFunctionsClient) getHandleResponse(resp *http.Response) (VendorNetworkFunctionsClientGetResponse, error) {
	result := VendorNetworkFunctionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VendorNetworkFunction); err != nil {
		return VendorNetworkFunctionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all the vendor network function sub resources in an Azure region, filtered by skuType, skuName, vendorProvisioningState.
//
// Generated from API version 2021-05-01
//   - locationName - The Azure region where the network function resource was created by the customer.
//   - vendorName - The name of the vendor.
//   - options - VendorNetworkFunctionsClientListOptions contains the optional parameters for the VendorNetworkFunctionsClient.NewListPager
//     method.
func (client *VendorNetworkFunctionsClient) NewListPager(locationName string, vendorName string, options *VendorNetworkFunctionsClientListOptions) (*runtime.Pager[VendorNetworkFunctionsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[VendorNetworkFunctionsClientListResponse]{
		More: func(page VendorNetworkFunctionsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VendorNetworkFunctionsClientListResponse) (VendorNetworkFunctionsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, locationName, vendorName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VendorNetworkFunctionsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return VendorNetworkFunctionsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VendorNetworkFunctionsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *VendorNetworkFunctionsClient) listCreateRequest(ctx context.Context, locationName string, vendorName string, options *VendorNetworkFunctionsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/locations/{locationName}/vendors/{vendorName}/networkFunctions"
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if vendorName == "" {
		return nil, errors.New("parameter vendorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vendorName}", url.PathEscape(vendorName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *VendorNetworkFunctionsClient) listHandleResponse(resp *http.Response) (VendorNetworkFunctionsClientListResponse, error) {
	result := VendorNetworkFunctionsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VendorNetworkFunctionListResult); err != nil {
		return VendorNetworkFunctionsClientListResponse{}, err
	}
	return result, nil
}

