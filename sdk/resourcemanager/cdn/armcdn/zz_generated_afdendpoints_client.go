//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcdn

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

// AFDEndpointsClient contains the methods for the AFDEndpoints group.
// Don't use this type directly, use NewAFDEndpointsClient() instead.
type AFDEndpointsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewAFDEndpointsClient creates a new instance of AFDEndpointsClient with the specified values.
// subscriptionID - Azure Subscription ID.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewAFDEndpointsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AFDEndpointsClient, error) {
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
	client := &AFDEndpointsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreate - Creates a new AzureFrontDoor endpoint with the specified endpoint name under the specified subscription,
// resource group and profile.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// endpointName - Name of the endpoint under the profile which is unique globally.
// endpoint - Endpoint properties
// options - AFDEndpointsClientBeginCreateOptions contains the optional parameters for the AFDEndpointsClient.BeginCreate
// method.
func (client *AFDEndpointsClient) BeginCreate(ctx context.Context, resourceGroupName string, profileName string, endpointName string, endpoint AFDEndpoint, options *AFDEndpointsClientBeginCreateOptions) (*armruntime.Poller[AFDEndpointsClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, profileName, endpointName, endpoint, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[AFDEndpointsClientCreateResponse]{
			FinalStateVia: armruntime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[AFDEndpointsClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - Creates a new AzureFrontDoor endpoint with the specified endpoint name under the specified subscription, resource
// group and profile.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *AFDEndpointsClient) create(ctx context.Context, resourceGroupName string, profileName string, endpointName string, endpoint AFDEndpoint, options *AFDEndpointsClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, profileName, endpointName, endpoint, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *AFDEndpointsClient) createCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointName string, endpoint AFDEndpoint, options *AFDEndpointsClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints/{endpointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, endpoint)
}

// BeginDelete - Deletes an existing AzureFrontDoor endpoint with the specified endpoint name under the specified subscription,
// resource group and profile.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// endpointName - Name of the endpoint under the profile which is unique globally.
// options - AFDEndpointsClientBeginDeleteOptions contains the optional parameters for the AFDEndpointsClient.BeginDelete
// method.
func (client *AFDEndpointsClient) BeginDelete(ctx context.Context, resourceGroupName string, profileName string, endpointName string, options *AFDEndpointsClientBeginDeleteOptions) (*armruntime.Poller[AFDEndpointsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, profileName, endpointName, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[AFDEndpointsClientDeleteResponse]{
			FinalStateVia: armruntime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[AFDEndpointsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Deletes an existing AzureFrontDoor endpoint with the specified endpoint name under the specified subscription,
// resource group and profile.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *AFDEndpointsClient) deleteOperation(ctx context.Context, resourceGroupName string, profileName string, endpointName string, options *AFDEndpointsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, profileName, endpointName, options)
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
func (client *AFDEndpointsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointName string, options *AFDEndpointsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints/{endpointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets an existing AzureFrontDoor endpoint with the specified endpoint name under the specified subscription, resource
// group and profile.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// endpointName - Name of the endpoint under the profile which is unique globally.
// options - AFDEndpointsClientGetOptions contains the optional parameters for the AFDEndpointsClient.Get method.
func (client *AFDEndpointsClient) Get(ctx context.Context, resourceGroupName string, profileName string, endpointName string, options *AFDEndpointsClientGetOptions) (AFDEndpointsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, profileName, endpointName, options)
	if err != nil {
		return AFDEndpointsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AFDEndpointsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AFDEndpointsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AFDEndpointsClient) getCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointName string, options *AFDEndpointsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints/{endpointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AFDEndpointsClient) getHandleResponse(resp *http.Response) (AFDEndpointsClientGetResponse, error) {
	result := AFDEndpointsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AFDEndpoint); err != nil {
		return AFDEndpointsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByProfilePager - Lists existing AzureFrontDoor endpoints.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// options - AFDEndpointsClientListByProfileOptions contains the optional parameters for the AFDEndpointsClient.ListByProfile
// method.
func (client *AFDEndpointsClient) NewListByProfilePager(resourceGroupName string, profileName string, options *AFDEndpointsClientListByProfileOptions) *runtime.Pager[AFDEndpointsClientListByProfileResponse] {
	return runtime.NewPager(runtime.PageProcessor[AFDEndpointsClientListByProfileResponse]{
		More: func(page AFDEndpointsClientListByProfileResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AFDEndpointsClientListByProfileResponse) (AFDEndpointsClientListByProfileResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByProfileCreateRequest(ctx, resourceGroupName, profileName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AFDEndpointsClientListByProfileResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return AFDEndpointsClientListByProfileResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AFDEndpointsClientListByProfileResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByProfileHandleResponse(resp)
		},
	})
}

// listByProfileCreateRequest creates the ListByProfile request.
func (client *AFDEndpointsClient) listByProfileCreateRequest(ctx context.Context, resourceGroupName string, profileName string, options *AFDEndpointsClientListByProfileOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByProfileHandleResponse handles the ListByProfile response.
func (client *AFDEndpointsClient) listByProfileHandleResponse(resp *http.Response) (AFDEndpointsClientListByProfileResponse, error) {
	result := AFDEndpointsClientListByProfileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AFDEndpointListResult); err != nil {
		return AFDEndpointsClientListByProfileResponse{}, err
	}
	return result, nil
}

// NewListResourceUsagePager - Checks the quota and actual usage of endpoints under the given CDN profile.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// endpointName - Name of the endpoint under the profile which is unique globally.
// options - AFDEndpointsClientListResourceUsageOptions contains the optional parameters for the AFDEndpointsClient.ListResourceUsage
// method.
func (client *AFDEndpointsClient) NewListResourceUsagePager(resourceGroupName string, profileName string, endpointName string, options *AFDEndpointsClientListResourceUsageOptions) *runtime.Pager[AFDEndpointsClientListResourceUsageResponse] {
	return runtime.NewPager(runtime.PageProcessor[AFDEndpointsClientListResourceUsageResponse]{
		More: func(page AFDEndpointsClientListResourceUsageResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AFDEndpointsClientListResourceUsageResponse) (AFDEndpointsClientListResourceUsageResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listResourceUsageCreateRequest(ctx, resourceGroupName, profileName, endpointName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AFDEndpointsClientListResourceUsageResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return AFDEndpointsClientListResourceUsageResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AFDEndpointsClientListResourceUsageResponse{}, runtime.NewResponseError(resp)
			}
			return client.listResourceUsageHandleResponse(resp)
		},
	})
}

// listResourceUsageCreateRequest creates the ListResourceUsage request.
func (client *AFDEndpointsClient) listResourceUsageCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointName string, options *AFDEndpointsClientListResourceUsageOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints/{endpointName}/usages"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listResourceUsageHandleResponse handles the ListResourceUsage response.
func (client *AFDEndpointsClient) listResourceUsageHandleResponse(resp *http.Response) (AFDEndpointsClientListResourceUsageResponse, error) {
	result := AFDEndpointsClientListResourceUsageResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UsagesListResult); err != nil {
		return AFDEndpointsClientListResourceUsageResponse{}, err
	}
	return result, nil
}

// BeginPurgeContent - Removes a content from AzureFrontDoor.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// endpointName - Name of the endpoint under the profile which is unique globally.
// contents - The list of paths to the content and the list of linked domains to be purged. Path can be a full URL, e.g. '/pictures/city.png'
// which removes a single file, or a directory with a wildcard, e.g.
// '/pictures/*' which removes all folders and files in the directory.
// options - AFDEndpointsClientBeginPurgeContentOptions contains the optional parameters for the AFDEndpointsClient.BeginPurgeContent
// method.
func (client *AFDEndpointsClient) BeginPurgeContent(ctx context.Context, resourceGroupName string, profileName string, endpointName string, contents AfdPurgeParameters, options *AFDEndpointsClientBeginPurgeContentOptions) (*armruntime.Poller[AFDEndpointsClientPurgeContentResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.purgeContent(ctx, resourceGroupName, profileName, endpointName, contents, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[AFDEndpointsClientPurgeContentResponse]{
			FinalStateVia: armruntime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[AFDEndpointsClientPurgeContentResponse](options.ResumeToken, client.pl, nil)
	}
}

// PurgeContent - Removes a content from AzureFrontDoor.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *AFDEndpointsClient) purgeContent(ctx context.Context, resourceGroupName string, profileName string, endpointName string, contents AfdPurgeParameters, options *AFDEndpointsClientBeginPurgeContentOptions) (*http.Response, error) {
	req, err := client.purgeContentCreateRequest(ctx, resourceGroupName, profileName, endpointName, contents, options)
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

// purgeContentCreateRequest creates the PurgeContent request.
func (client *AFDEndpointsClient) purgeContentCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointName string, contents AfdPurgeParameters, options *AFDEndpointsClientBeginPurgeContentOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints/{endpointName}/purge"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, contents)
}

// BeginUpdate - Updates an existing AzureFrontDoor endpoint with the specified endpoint name under the specified subscription,
// resource group and profile. Only tags can be updated after creating an endpoint. To
// update origins, use the Update Origin operation. To update origin groups, use the Update Origin group operation. To update
// domains, use the Update Custom Domain operation.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// endpointName - Name of the endpoint under the profile which is unique globally.
// endpointUpdateProperties - Endpoint update properties
// options - AFDEndpointsClientBeginUpdateOptions contains the optional parameters for the AFDEndpointsClient.BeginUpdate
// method.
func (client *AFDEndpointsClient) BeginUpdate(ctx context.Context, resourceGroupName string, profileName string, endpointName string, endpointUpdateProperties AFDEndpointUpdateParameters, options *AFDEndpointsClientBeginUpdateOptions) (*armruntime.Poller[AFDEndpointsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, profileName, endpointName, endpointUpdateProperties, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[AFDEndpointsClientUpdateResponse]{
			FinalStateVia: armruntime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[AFDEndpointsClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - Updates an existing AzureFrontDoor endpoint with the specified endpoint name under the specified subscription,
// resource group and profile. Only tags can be updated after creating an endpoint. To
// update origins, use the Update Origin operation. To update origin groups, use the Update Origin group operation. To update
// domains, use the Update Custom Domain operation.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *AFDEndpointsClient) update(ctx context.Context, resourceGroupName string, profileName string, endpointName string, endpointUpdateProperties AFDEndpointUpdateParameters, options *AFDEndpointsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, profileName, endpointName, endpointUpdateProperties, options)
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
func (client *AFDEndpointsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointName string, endpointUpdateProperties AFDEndpointUpdateParameters, options *AFDEndpointsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints/{endpointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, endpointUpdateProperties)
}

// ValidateCustomDomain - Validates the custom domain mapping to ensure it maps to the correct CDN endpoint in DNS.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - Name of the Resource group within the Azure subscription.
// profileName - Name of the Azure Front Door Standard or Azure Front Door Premium profile which is unique within the resource
// group.
// endpointName - Name of the endpoint under the profile which is unique globally.
// customDomainProperties - Custom domain to be validated.
// options - AFDEndpointsClientValidateCustomDomainOptions contains the optional parameters for the AFDEndpointsClient.ValidateCustomDomain
// method.
func (client *AFDEndpointsClient) ValidateCustomDomain(ctx context.Context, resourceGroupName string, profileName string, endpointName string, customDomainProperties ValidateCustomDomainInput, options *AFDEndpointsClientValidateCustomDomainOptions) (AFDEndpointsClientValidateCustomDomainResponse, error) {
	req, err := client.validateCustomDomainCreateRequest(ctx, resourceGroupName, profileName, endpointName, customDomainProperties, options)
	if err != nil {
		return AFDEndpointsClientValidateCustomDomainResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AFDEndpointsClientValidateCustomDomainResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AFDEndpointsClientValidateCustomDomainResponse{}, runtime.NewResponseError(resp)
	}
	return client.validateCustomDomainHandleResponse(resp)
}

// validateCustomDomainCreateRequest creates the ValidateCustomDomain request.
func (client *AFDEndpointsClient) validateCustomDomainCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointName string, customDomainProperties ValidateCustomDomainInput, options *AFDEndpointsClientValidateCustomDomainOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints/{endpointName}/validateCustomDomain"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, customDomainProperties)
}

// validateCustomDomainHandleResponse handles the ValidateCustomDomain response.
func (client *AFDEndpointsClient) validateCustomDomainHandleResponse(resp *http.Response) (AFDEndpointsClientValidateCustomDomainResponse, error) {
	result := AFDEndpointsClientValidateCustomDomainResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ValidateCustomDomainOutput); err != nil {
		return AFDEndpointsClientValidateCustomDomainResponse{}, err
	}
	return result, nil
}
