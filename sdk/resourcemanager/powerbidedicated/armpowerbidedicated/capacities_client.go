//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpowerbidedicated

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

// CapacitiesClient contains the methods for the Capacities group.
// Don't use this type directly, use NewCapacitiesClient() instead.
type CapacitiesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewCapacitiesClient creates a new instance of CapacitiesClient with the specified values.
//   - subscriptionID - A unique identifier for a Microsoft Azure subscription. The subscription ID forms part of the URI for
//     every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCapacitiesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CapacitiesClient, error) {
	cl, err := arm.NewClient(moduleName+".CapacitiesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CapacitiesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CheckNameAvailability - Check the name availability in the target location.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - location - The region name which the operation will lookup into.
//   - capacityParameters - The name of the capacity.
//   - options - CapacitiesClientCheckNameAvailabilityOptions contains the optional parameters for the CapacitiesClient.CheckNameAvailability
//     method.
func (client *CapacitiesClient) CheckNameAvailability(ctx context.Context, location string, capacityParameters CheckCapacityNameAvailabilityParameters, options *CapacitiesClientCheckNameAvailabilityOptions) (CapacitiesClientCheckNameAvailabilityResponse, error) {
	req, err := client.checkNameAvailabilityCreateRequest(ctx, location, capacityParameters, options)
	if err != nil {
		return CapacitiesClientCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CapacitiesClientCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return CapacitiesClientCheckNameAvailabilityResponse{}, runtime.NewResponseError(resp)
	}
	return client.checkNameAvailabilityHandleResponse(resp)
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *CapacitiesClient) checkNameAvailabilityCreateRequest(ctx context.Context, location string, capacityParameters CheckCapacityNameAvailabilityParameters, options *CapacitiesClientCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.PowerBIDedicated/locations/{location}/checkNameAvailability"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, capacityParameters)
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *CapacitiesClient) checkNameAvailabilityHandleResponse(resp *http.Response) (CapacitiesClientCheckNameAvailabilityResponse, error) {
	result := CapacitiesClientCheckNameAvailabilityResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckCapacityNameAvailabilityResult); err != nil {
		return CapacitiesClientCheckNameAvailabilityResponse{}, err
	}
	return result, nil
}

// BeginCreate - Provisions the specified Dedicated capacity based on the configuration specified in the request.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - resourceGroupName - The name of the Azure Resource group of which a given PowerBIDedicated capacity is part. This name
//     must be at least 1 character in length, and no more than 90.
//   - dedicatedCapacityName - The name of the Dedicated capacity. It must be a minimum of 3 characters, and a maximum of 63.
//   - capacityParameters - Contains the information used to provision the Dedicated capacity.
//   - options - CapacitiesClientBeginCreateOptions contains the optional parameters for the CapacitiesClient.BeginCreate method.
func (client *CapacitiesClient) BeginCreate(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, capacityParameters DedicatedCapacity, options *CapacitiesClientBeginCreateOptions) (*runtime.Poller[CapacitiesClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, dedicatedCapacityName, capacityParameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[CapacitiesClientCreateResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[CapacitiesClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Provisions the specified Dedicated capacity based on the configuration specified in the request.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
func (client *CapacitiesClient) create(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, capacityParameters DedicatedCapacity, options *CapacitiesClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, dedicatedCapacityName, capacityParameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *CapacitiesClient) createCreateRequest(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, capacityParameters DedicatedCapacity, options *CapacitiesClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerBIDedicated/capacities/{dedicatedCapacityName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dedicatedCapacityName == "" {
		return nil, errors.New("parameter dedicatedCapacityName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dedicatedCapacityName}", url.PathEscape(dedicatedCapacityName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, capacityParameters)
}

// BeginDelete - Deletes the specified Dedicated capacity.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - resourceGroupName - The name of the Azure Resource group of which a given PowerBIDedicated capacity is part. This name
//     must be at least 1 character in length, and no more than 90.
//   - dedicatedCapacityName - The name of the Dedicated capacity. It must be at least 3 characters in length, and no more than
//     63.
//   - options - CapacitiesClientBeginDeleteOptions contains the optional parameters for the CapacitiesClient.BeginDelete method.
func (client *CapacitiesClient) BeginDelete(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginDeleteOptions) (*runtime.Poller[CapacitiesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, dedicatedCapacityName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[CapacitiesClientDeleteResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[CapacitiesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the specified Dedicated capacity.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
func (client *CapacitiesClient) deleteOperation(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, dedicatedCapacityName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *CapacitiesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerBIDedicated/capacities/{dedicatedCapacityName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dedicatedCapacityName == "" {
		return nil, errors.New("parameter dedicatedCapacityName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dedicatedCapacityName}", url.PathEscape(dedicatedCapacityName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// GetDetails - Gets details about the specified dedicated capacity.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - resourceGroupName - The name of the Azure Resource group of which a given PowerBIDedicated capacity is part. This name
//     must be at least 1 character in length, and no more than 90.
//   - dedicatedCapacityName - The name of the dedicated capacity. It must be a minimum of 3 characters, and a maximum of 63.
//   - options - CapacitiesClientGetDetailsOptions contains the optional parameters for the CapacitiesClient.GetDetails method.
func (client *CapacitiesClient) GetDetails(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientGetDetailsOptions) (CapacitiesClientGetDetailsResponse, error) {
	req, err := client.getDetailsCreateRequest(ctx, resourceGroupName, dedicatedCapacityName, options)
	if err != nil {
		return CapacitiesClientGetDetailsResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CapacitiesClientGetDetailsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return CapacitiesClientGetDetailsResponse{}, runtime.NewResponseError(resp)
	}
	return client.getDetailsHandleResponse(resp)
}

// getDetailsCreateRequest creates the GetDetails request.
func (client *CapacitiesClient) getDetailsCreateRequest(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientGetDetailsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerBIDedicated/capacities/{dedicatedCapacityName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dedicatedCapacityName == "" {
		return nil, errors.New("parameter dedicatedCapacityName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dedicatedCapacityName}", url.PathEscape(dedicatedCapacityName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getDetailsHandleResponse handles the GetDetails response.
func (client *CapacitiesClient) getDetailsHandleResponse(resp *http.Response) (CapacitiesClientGetDetailsResponse, error) {
	result := CapacitiesClientGetDetailsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DedicatedCapacity); err != nil {
		return CapacitiesClientGetDetailsResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all the Dedicated capacities for the given subscription.
//
// Generated from API version 2021-01-01
//   - options - CapacitiesClientListOptions contains the optional parameters for the CapacitiesClient.NewListPager method.
func (client *CapacitiesClient) NewListPager(options *CapacitiesClientListOptions) *runtime.Pager[CapacitiesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[CapacitiesClientListResponse]{
		More: func(page CapacitiesClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *CapacitiesClientListResponse) (CapacitiesClientListResponse, error) {
			req, err := client.listCreateRequest(ctx, options)
			if err != nil {
				return CapacitiesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CapacitiesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CapacitiesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *CapacitiesClient) listCreateRequest(ctx context.Context, options *CapacitiesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.PowerBIDedicated/capacities"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *CapacitiesClient) listHandleResponse(resp *http.Response) (CapacitiesClientListResponse, error) {
	result := CapacitiesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DedicatedCapacities); err != nil {
		return CapacitiesClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Gets all the Dedicated capacities for the given resource group.
//
// Generated from API version 2021-01-01
//   - resourceGroupName - The name of the Azure Resource group of which a given PowerBIDedicated capacity is part. This name
//     must be at least 1 character in length, and no more than 90.
//   - options - CapacitiesClientListByResourceGroupOptions contains the optional parameters for the CapacitiesClient.NewListByResourceGroupPager
//     method.
func (client *CapacitiesClient) NewListByResourceGroupPager(resourceGroupName string, options *CapacitiesClientListByResourceGroupOptions) *runtime.Pager[CapacitiesClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[CapacitiesClientListByResourceGroupResponse]{
		More: func(page CapacitiesClientListByResourceGroupResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *CapacitiesClientListByResourceGroupResponse) (CapacitiesClientListByResourceGroupResponse, error) {
			req, err := client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			if err != nil {
				return CapacitiesClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CapacitiesClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CapacitiesClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *CapacitiesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *CapacitiesClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerBIDedicated/capacities"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *CapacitiesClient) listByResourceGroupHandleResponse(resp *http.Response) (CapacitiesClientListByResourceGroupResponse, error) {
	result := CapacitiesClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DedicatedCapacities); err != nil {
		return CapacitiesClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// ListSKUs - Lists eligible SKUs for PowerBI Dedicated resource provider.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - options - CapacitiesClientListSKUsOptions contains the optional parameters for the CapacitiesClient.ListSKUs method.
func (client *CapacitiesClient) ListSKUs(ctx context.Context, options *CapacitiesClientListSKUsOptions) (CapacitiesClientListSKUsResponse, error) {
	req, err := client.listSKUsCreateRequest(ctx, options)
	if err != nil {
		return CapacitiesClientListSKUsResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CapacitiesClientListSKUsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return CapacitiesClientListSKUsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listSKUsHandleResponse(resp)
}

// listSKUsCreateRequest creates the ListSKUs request.
func (client *CapacitiesClient) listSKUsCreateRequest(ctx context.Context, options *CapacitiesClientListSKUsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.PowerBIDedicated/skus"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSKUsHandleResponse handles the ListSKUs response.
func (client *CapacitiesClient) listSKUsHandleResponse(resp *http.Response) (CapacitiesClientListSKUsResponse, error) {
	result := CapacitiesClientListSKUsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SKUEnumerationForNewResourceResult); err != nil {
		return CapacitiesClientListSKUsResponse{}, err
	}
	return result, nil
}

// ListSKUsForCapacity - Lists eligible SKUs for a PowerBI Dedicated resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - resourceGroupName - The name of the Azure Resource group of which a given PowerBIDedicated capacity is part. This name
//     must be at least 1 character in length, and no more than 90.
//   - dedicatedCapacityName - The name of the Dedicated capacity. It must be at least 3 characters in length, and no more than
//     63.
//   - options - CapacitiesClientListSKUsForCapacityOptions contains the optional parameters for the CapacitiesClient.ListSKUsForCapacity
//     method.
func (client *CapacitiesClient) ListSKUsForCapacity(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientListSKUsForCapacityOptions) (CapacitiesClientListSKUsForCapacityResponse, error) {
	req, err := client.listSKUsForCapacityCreateRequest(ctx, resourceGroupName, dedicatedCapacityName, options)
	if err != nil {
		return CapacitiesClientListSKUsForCapacityResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CapacitiesClientListSKUsForCapacityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return CapacitiesClientListSKUsForCapacityResponse{}, runtime.NewResponseError(resp)
	}
	return client.listSKUsForCapacityHandleResponse(resp)
}

// listSKUsForCapacityCreateRequest creates the ListSKUsForCapacity request.
func (client *CapacitiesClient) listSKUsForCapacityCreateRequest(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientListSKUsForCapacityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerBIDedicated/capacities/{dedicatedCapacityName}/skus"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dedicatedCapacityName == "" {
		return nil, errors.New("parameter dedicatedCapacityName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dedicatedCapacityName}", url.PathEscape(dedicatedCapacityName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSKUsForCapacityHandleResponse handles the ListSKUsForCapacity response.
func (client *CapacitiesClient) listSKUsForCapacityHandleResponse(resp *http.Response) (CapacitiesClientListSKUsForCapacityResponse, error) {
	result := CapacitiesClientListSKUsForCapacityResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SKUEnumerationForExistingResourceResult); err != nil {
		return CapacitiesClientListSKUsForCapacityResponse{}, err
	}
	return result, nil
}

// BeginResume - Resumes operation of the specified Dedicated capacity instance.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - resourceGroupName - The name of the Azure Resource group of which a given PowerBIDedicated capacity is part. This name
//     must be at least 1 character in length, and no more than 90.
//   - dedicatedCapacityName - The name of the Dedicated capacity. It must be at least 3 characters in length, and no more than
//     63.
//   - options - CapacitiesClientBeginResumeOptions contains the optional parameters for the CapacitiesClient.BeginResume method.
func (client *CapacitiesClient) BeginResume(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginResumeOptions) (*runtime.Poller[CapacitiesClientResumeResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.resume(ctx, resourceGroupName, dedicatedCapacityName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[CapacitiesClientResumeResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[CapacitiesClientResumeResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Resume - Resumes operation of the specified Dedicated capacity instance.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
func (client *CapacitiesClient) resume(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginResumeOptions) (*http.Response, error) {
	req, err := client.resumeCreateRequest(ctx, resourceGroupName, dedicatedCapacityName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// resumeCreateRequest creates the Resume request.
func (client *CapacitiesClient) resumeCreateRequest(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginResumeOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerBIDedicated/capacities/{dedicatedCapacityName}/resume"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dedicatedCapacityName == "" {
		return nil, errors.New("parameter dedicatedCapacityName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dedicatedCapacityName}", url.PathEscape(dedicatedCapacityName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginSuspend - Suspends operation of the specified dedicated capacity instance.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - resourceGroupName - The name of the Azure Resource group of which a given PowerBIDedicated capacity is part. This name
//     must be at least 1 character in length, and no more than 90.
//   - dedicatedCapacityName - The name of the Dedicated capacity. It must be at least 3 characters in length, and no more than
//     63.
//   - options - CapacitiesClientBeginSuspendOptions contains the optional parameters for the CapacitiesClient.BeginSuspend method.
func (client *CapacitiesClient) BeginSuspend(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginSuspendOptions) (*runtime.Poller[CapacitiesClientSuspendResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.suspend(ctx, resourceGroupName, dedicatedCapacityName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[CapacitiesClientSuspendResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[CapacitiesClientSuspendResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Suspend - Suspends operation of the specified dedicated capacity instance.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
func (client *CapacitiesClient) suspend(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginSuspendOptions) (*http.Response, error) {
	req, err := client.suspendCreateRequest(ctx, resourceGroupName, dedicatedCapacityName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// suspendCreateRequest creates the Suspend request.
func (client *CapacitiesClient) suspendCreateRequest(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, options *CapacitiesClientBeginSuspendOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerBIDedicated/capacities/{dedicatedCapacityName}/suspend"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dedicatedCapacityName == "" {
		return nil, errors.New("parameter dedicatedCapacityName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dedicatedCapacityName}", url.PathEscape(dedicatedCapacityName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginUpdate - Updates the current state of the specified Dedicated capacity.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
//   - resourceGroupName - The name of the Azure Resource group of which a given PowerBIDedicated capacity is part. This name
//     must be at least 1 character in length, and no more than 90.
//   - dedicatedCapacityName - The name of the Dedicated capacity. It must be at least 3 characters in length, and no more than
//     63.
//   - capacityUpdateParameters - Request object that contains the updated information for the capacity.
//   - options - CapacitiesClientBeginUpdateOptions contains the optional parameters for the CapacitiesClient.BeginUpdate method.
func (client *CapacitiesClient) BeginUpdate(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, capacityUpdateParameters DedicatedCapacityUpdateParameters, options *CapacitiesClientBeginUpdateOptions) (*runtime.Poller[CapacitiesClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, dedicatedCapacityName, capacityUpdateParameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[CapacitiesClientUpdateResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[CapacitiesClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Updates the current state of the specified Dedicated capacity.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-01-01
func (client *CapacitiesClient) update(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, capacityUpdateParameters DedicatedCapacityUpdateParameters, options *CapacitiesClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, dedicatedCapacityName, capacityUpdateParameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *CapacitiesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, dedicatedCapacityName string, capacityUpdateParameters DedicatedCapacityUpdateParameters, options *CapacitiesClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerBIDedicated/capacities/{dedicatedCapacityName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dedicatedCapacityName == "" {
		return nil, errors.New("parameter dedicatedCapacityName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dedicatedCapacityName}", url.PathEscape(dedicatedCapacityName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, capacityUpdateParameters)
}
