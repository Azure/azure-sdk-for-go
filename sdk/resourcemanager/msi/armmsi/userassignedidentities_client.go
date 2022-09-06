//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armmsi

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

// UserAssignedIdentitiesClient contains the methods for the UserAssignedIdentities group.
// Don't use this type directly, use NewUserAssignedIdentitiesClient() instead.
type UserAssignedIdentitiesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewUserAssignedIdentitiesClient creates a new instance of UserAssignedIdentitiesClient with the specified values.
// subscriptionID - The Id of the Subscription to which the identity belongs.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewUserAssignedIdentitiesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*UserAssignedIdentitiesClient, error) {
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
	client := &UserAssignedIdentitiesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update an identity in the specified subscription and resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-31-preview
// resourceGroupName - The name of the Resource Group to which the identity belongs.
// resourceName - The name of the identity resource.
// parameters - Parameters to create or update the identity
// options - UserAssignedIdentitiesClientCreateOrUpdateOptions contains the optional parameters for the UserAssignedIdentitiesClient.CreateOrUpdate
// method.
func (client *UserAssignedIdentitiesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, parameters Identity, options *UserAssignedIdentitiesClientCreateOrUpdateOptions) (UserAssignedIdentitiesClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, resourceName, parameters, options)
	if err != nil {
		return UserAssignedIdentitiesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return UserAssignedIdentitiesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return UserAssignedIdentitiesClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *UserAssignedIdentitiesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, parameters Identity, options *UserAssignedIdentitiesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}"
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
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *UserAssignedIdentitiesClient) createOrUpdateHandleResponse(resp *http.Response) (UserAssignedIdentitiesClientCreateOrUpdateResponse, error) {
	result := UserAssignedIdentitiesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Identity); err != nil {
		return UserAssignedIdentitiesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the identity.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-31-preview
// resourceGroupName - The name of the Resource Group to which the identity belongs.
// resourceName - The name of the identity resource.
// options - UserAssignedIdentitiesClientDeleteOptions contains the optional parameters for the UserAssignedIdentitiesClient.Delete
// method.
func (client *UserAssignedIdentitiesClient) Delete(ctx context.Context, resourceGroupName string, resourceName string, options *UserAssignedIdentitiesClientDeleteOptions) (UserAssignedIdentitiesClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return UserAssignedIdentitiesClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return UserAssignedIdentitiesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return UserAssignedIdentitiesClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return UserAssignedIdentitiesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *UserAssignedIdentitiesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *UserAssignedIdentitiesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}"
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
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the identity.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-31-preview
// resourceGroupName - The name of the Resource Group to which the identity belongs.
// resourceName - The name of the identity resource.
// options - UserAssignedIdentitiesClientGetOptions contains the optional parameters for the UserAssignedIdentitiesClient.Get
// method.
func (client *UserAssignedIdentitiesClient) Get(ctx context.Context, resourceGroupName string, resourceName string, options *UserAssignedIdentitiesClientGetOptions) (UserAssignedIdentitiesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return UserAssignedIdentitiesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return UserAssignedIdentitiesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return UserAssignedIdentitiesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *UserAssignedIdentitiesClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *UserAssignedIdentitiesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}"
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
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *UserAssignedIdentitiesClient) getHandleResponse(resp *http.Response) (UserAssignedIdentitiesClientGetResponse, error) {
	result := UserAssignedIdentitiesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Identity); err != nil {
		return UserAssignedIdentitiesClientGetResponse{}, err
	}
	return result, nil
}

// NewListAssociatedResourcesPager - Lists the associated resources for this identity.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-31-preview
// resourceGroupName - The name of the Resource Group to which the identity belongs.
// resourceName - The name of the identity resource.
// options - UserAssignedIdentitiesClientListAssociatedResourcesOptions contains the optional parameters for the UserAssignedIdentitiesClient.ListAssociatedResources
// method.
func (client *UserAssignedIdentitiesClient) NewListAssociatedResourcesPager(resourceGroupName string, resourceName string, options *UserAssignedIdentitiesClientListAssociatedResourcesOptions) *runtime.Pager[UserAssignedIdentitiesClientListAssociatedResourcesResponse] {
	return runtime.NewPager(runtime.PagingHandler[UserAssignedIdentitiesClientListAssociatedResourcesResponse]{
		More: func(page UserAssignedIdentitiesClientListAssociatedResourcesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *UserAssignedIdentitiesClientListAssociatedResourcesResponse) (UserAssignedIdentitiesClientListAssociatedResourcesResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listAssociatedResourcesCreateRequest(ctx, resourceGroupName, resourceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return UserAssignedIdentitiesClientListAssociatedResourcesResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return UserAssignedIdentitiesClientListAssociatedResourcesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return UserAssignedIdentitiesClientListAssociatedResourcesResponse{}, runtime.NewResponseError(resp)
			}
			return client.listAssociatedResourcesHandleResponse(resp)
		},
	})
}

// listAssociatedResourcesCreateRequest creates the ListAssociatedResources request.
func (client *UserAssignedIdentitiesClient) listAssociatedResourcesCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *UserAssignedIdentitiesClientListAssociatedResourcesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}/listAssociatedResources"
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
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Skiptoken != nil {
		reqQP.Set("$skiptoken", *options.Skiptoken)
	}
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listAssociatedResourcesHandleResponse handles the ListAssociatedResources response.
func (client *UserAssignedIdentitiesClient) listAssociatedResourcesHandleResponse(resp *http.Response) (UserAssignedIdentitiesClientListAssociatedResourcesResponse, error) {
	result := UserAssignedIdentitiesClientListAssociatedResourcesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AssociatedResourcesListResult); err != nil {
		return UserAssignedIdentitiesClientListAssociatedResourcesResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Lists all the userAssignedIdentities available under the specified ResourceGroup.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-31-preview
// resourceGroupName - The name of the Resource Group to which the identity belongs.
// options - UserAssignedIdentitiesClientListByResourceGroupOptions contains the optional parameters for the UserAssignedIdentitiesClient.ListByResourceGroup
// method.
func (client *UserAssignedIdentitiesClient) NewListByResourceGroupPager(resourceGroupName string, options *UserAssignedIdentitiesClientListByResourceGroupOptions) *runtime.Pager[UserAssignedIdentitiesClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[UserAssignedIdentitiesClientListByResourceGroupResponse]{
		More: func(page UserAssignedIdentitiesClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *UserAssignedIdentitiesClientListByResourceGroupResponse) (UserAssignedIdentitiesClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return UserAssignedIdentitiesClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return UserAssignedIdentitiesClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return UserAssignedIdentitiesClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *UserAssignedIdentitiesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *UserAssignedIdentitiesClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities"
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
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *UserAssignedIdentitiesClient) listByResourceGroupHandleResponse(resp *http.Response) (UserAssignedIdentitiesClientListByResourceGroupResponse, error) {
	result := UserAssignedIdentitiesClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UserAssignedIdentitiesListResult); err != nil {
		return UserAssignedIdentitiesClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Lists all the userAssignedIdentities available under the specified subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-31-preview
// options - UserAssignedIdentitiesClientListBySubscriptionOptions contains the optional parameters for the UserAssignedIdentitiesClient.ListBySubscription
// method.
func (client *UserAssignedIdentitiesClient) NewListBySubscriptionPager(options *UserAssignedIdentitiesClientListBySubscriptionOptions) *runtime.Pager[UserAssignedIdentitiesClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[UserAssignedIdentitiesClientListBySubscriptionResponse]{
		More: func(page UserAssignedIdentitiesClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *UserAssignedIdentitiesClientListBySubscriptionResponse) (UserAssignedIdentitiesClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return UserAssignedIdentitiesClientListBySubscriptionResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return UserAssignedIdentitiesClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return UserAssignedIdentitiesClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *UserAssignedIdentitiesClient) listBySubscriptionCreateRequest(ctx context.Context, options *UserAssignedIdentitiesClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ManagedIdentity/userAssignedIdentities"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *UserAssignedIdentitiesClient) listBySubscriptionHandleResponse(resp *http.Response) (UserAssignedIdentitiesClientListBySubscriptionResponse, error) {
	result := UserAssignedIdentitiesClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UserAssignedIdentitiesListResult); err != nil {
		return UserAssignedIdentitiesClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Update an identity in the specified subscription and resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-31-preview
// resourceGroupName - The name of the Resource Group to which the identity belongs.
// resourceName - The name of the identity resource.
// parameters - Parameters to update the identity
// options - UserAssignedIdentitiesClientUpdateOptions contains the optional parameters for the UserAssignedIdentitiesClient.Update
// method.
func (client *UserAssignedIdentitiesClient) Update(ctx context.Context, resourceGroupName string, resourceName string, parameters IdentityUpdate, options *UserAssignedIdentitiesClientUpdateOptions) (UserAssignedIdentitiesClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, parameters, options)
	if err != nil {
		return UserAssignedIdentitiesClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return UserAssignedIdentitiesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return UserAssignedIdentitiesClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *UserAssignedIdentitiesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, parameters IdentityUpdate, options *UserAssignedIdentitiesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}"
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
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleResponse handles the Update response.
func (client *UserAssignedIdentitiesClient) updateHandleResponse(resp *http.Response) (UserAssignedIdentitiesClientUpdateResponse, error) {
	result := UserAssignedIdentitiesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Identity); err != nil {
		return UserAssignedIdentitiesClientUpdateResponse{}, err
	}
	return result, nil
}
