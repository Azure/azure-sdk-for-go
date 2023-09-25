//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armorbital

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

// ContactProfilesClient contains the methods for the ContactProfiles group.
// Don't use this type directly, use NewContactProfilesClient() instead.
type ContactProfilesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewContactProfilesClient creates a new instance of ContactProfilesClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewContactProfilesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ContactProfilesClient, error) {
	cl, err := arm.NewClient(moduleName+".ContactProfilesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ContactProfilesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a contact profile.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - contactProfileName - Contact Profile name.
//   - parameters - The parameters to provide for the created Contact Profile.
//   - options - ContactProfilesClientBeginCreateOrUpdateOptions contains the optional parameters for the ContactProfilesClient.BeginCreateOrUpdate
//     method.
func (client *ContactProfilesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, contactProfileName string, parameters ContactProfile, options *ContactProfilesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ContactProfilesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, contactProfileName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContactProfilesClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ContactProfilesClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a contact profile.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
func (client *ContactProfilesClient) createOrUpdate(ctx context.Context, resourceGroupName string, contactProfileName string, parameters ContactProfile, options *ContactProfilesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, contactProfileName, parameters, options)
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
func (client *ContactProfilesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, contactProfileName string, parameters ContactProfile, options *ContactProfilesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Orbital/contactProfiles/{contactProfileName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if contactProfileName == "" {
		return nil, errors.New("parameter contactProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{contactProfileName}", url.PathEscape(contactProfileName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes a specified contact profile resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - contactProfileName - Contact Profile name.
//   - options - ContactProfilesClientBeginDeleteOptions contains the optional parameters for the ContactProfilesClient.BeginDelete
//     method.
func (client *ContactProfilesClient) BeginDelete(ctx context.Context, resourceGroupName string, contactProfileName string, options *ContactProfilesClientBeginDeleteOptions) (*runtime.Poller[ContactProfilesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, contactProfileName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContactProfilesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ContactProfilesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes a specified contact profile resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
func (client *ContactProfilesClient) deleteOperation(ctx context.Context, resourceGroupName string, contactProfileName string, options *ContactProfilesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, contactProfileName, options)
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
func (client *ContactProfilesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, contactProfileName string, options *ContactProfilesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Orbital/contactProfiles/{contactProfileName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if contactProfileName == "" {
		return nil, errors.New("parameter contactProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{contactProfileName}", url.PathEscape(contactProfileName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the specified contact Profile in a specified resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - contactProfileName - Contact Profile name.
//   - options - ContactProfilesClientGetOptions contains the optional parameters for the ContactProfilesClient.Get method.
func (client *ContactProfilesClient) Get(ctx context.Context, resourceGroupName string, contactProfileName string, options *ContactProfilesClientGetOptions) (ContactProfilesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, contactProfileName, options)
	if err != nil {
		return ContactProfilesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ContactProfilesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ContactProfilesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ContactProfilesClient) getCreateRequest(ctx context.Context, resourceGroupName string, contactProfileName string, options *ContactProfilesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Orbital/contactProfiles/{contactProfileName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if contactProfileName == "" {
		return nil, errors.New("parameter contactProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{contactProfileName}", url.PathEscape(contactProfileName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ContactProfilesClient) getHandleResponse(resp *http.Response) (ContactProfilesClientGetResponse, error) {
	result := ContactProfilesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContactProfile); err != nil {
		return ContactProfilesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Returns list of contact profiles by Resource Group.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - ContactProfilesClientListOptions contains the optional parameters for the ContactProfilesClient.NewListPager
//     method.
func (client *ContactProfilesClient) NewListPager(resourceGroupName string, options *ContactProfilesClientListOptions) (*runtime.Pager[ContactProfilesClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ContactProfilesClientListResponse]{
		More: func(page ContactProfilesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ContactProfilesClientListResponse) (ContactProfilesClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ContactProfilesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ContactProfilesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ContactProfilesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ContactProfilesClient) listCreateRequest(ctx context.Context, resourceGroupName string, options *ContactProfilesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Orbital/contactProfiles"
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
	reqQP.Set("api-version", "2022-11-01")
	if options != nil && options.Skiptoken != nil {
		reqQP.Set("$skiptoken", *options.Skiptoken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ContactProfilesClient) listHandleResponse(resp *http.Response) (ContactProfilesClientListResponse, error) {
	result := ContactProfilesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContactProfileListResult); err != nil {
		return ContactProfilesClientListResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Returns list of contact profiles by Subscription.
//
// Generated from API version 2022-11-01
//   - options - ContactProfilesClientListBySubscriptionOptions contains the optional parameters for the ContactProfilesClient.NewListBySubscriptionPager
//     method.
func (client *ContactProfilesClient) NewListBySubscriptionPager(options *ContactProfilesClientListBySubscriptionOptions) (*runtime.Pager[ContactProfilesClientListBySubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ContactProfilesClientListBySubscriptionResponse]{
		More: func(page ContactProfilesClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ContactProfilesClientListBySubscriptionResponse) (ContactProfilesClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ContactProfilesClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ContactProfilesClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ContactProfilesClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *ContactProfilesClient) listBySubscriptionCreateRequest(ctx context.Context, options *ContactProfilesClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Orbital/contactProfiles"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	if options != nil && options.Skiptoken != nil {
		reqQP.Set("$skiptoken", *options.Skiptoken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *ContactProfilesClient) listBySubscriptionHandleResponse(resp *http.Response) (ContactProfilesClientListBySubscriptionResponse, error) {
	result := ContactProfilesClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContactProfileListResult); err != nil {
		return ContactProfilesClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginUpdateTags - Updates the specified contact profile tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - contactProfileName - Contact Profile name.
//   - parameters - Parameters supplied to update contact profile tags.
//   - options - ContactProfilesClientBeginUpdateTagsOptions contains the optional parameters for the ContactProfilesClient.BeginUpdateTags
//     method.
func (client *ContactProfilesClient) BeginUpdateTags(ctx context.Context, resourceGroupName string, contactProfileName string, parameters TagsObject, options *ContactProfilesClientBeginUpdateTagsOptions) (*runtime.Poller[ContactProfilesClientUpdateTagsResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.updateTags(ctx, resourceGroupName, contactProfileName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContactProfilesClientUpdateTagsResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ContactProfilesClientUpdateTagsResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// UpdateTags - Updates the specified contact profile tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
func (client *ContactProfilesClient) updateTags(ctx context.Context, resourceGroupName string, contactProfileName string, parameters TagsObject, options *ContactProfilesClientBeginUpdateTagsOptions) (*http.Response, error) {
	var err error
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, contactProfileName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *ContactProfilesClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, contactProfileName string, parameters TagsObject, options *ContactProfilesClientBeginUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Orbital/contactProfiles/{contactProfileName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if contactProfileName == "" {
		return nil, errors.New("parameter contactProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{contactProfileName}", url.PathEscape(contactProfileName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

