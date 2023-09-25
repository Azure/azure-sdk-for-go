//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpolicyinsights

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

// PolicyTrackedResourcesClient contains the methods for the PolicyTrackedResources group.
// Don't use this type directly, use NewPolicyTrackedResourcesClient() instead.
type PolicyTrackedResourcesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewPolicyTrackedResourcesClient creates a new instance of PolicyTrackedResourcesClient with the specified values.
//   - subscriptionID - Microsoft Azure subscription ID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPolicyTrackedResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PolicyTrackedResourcesClient, error) {
	cl, err := arm.NewClient(moduleName+".PolicyTrackedResourcesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PolicyTrackedResourcesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// NewListQueryResultsForManagementGroupPager - Queries policy tracked resources under the management group.
//
// Generated from API version 2018-07-01-preview
//   - managementGroupName - Management group name.
//   - policyTrackedResourcesResource - The name of the virtual resource under PolicyTrackedResources resource type; only "default"
//     is allowed.
//   - QueryOptions - QueryOptions contains a group of parameters for the PolicyTrackedResourcesClient.ListQueryResultsForManagementGroup
//     method.
//   - options - PolicyTrackedResourcesClientListQueryResultsForManagementGroupOptions contains the optional parameters for the
//     PolicyTrackedResourcesClient.NewListQueryResultsForManagementGroupPager method.
func (client *PolicyTrackedResourcesClient) NewListQueryResultsForManagementGroupPager(managementGroupName string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForManagementGroupOptions) (*runtime.Pager[PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse]{
		More: func(page PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse) (PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listQueryResultsForManagementGroupCreateRequest(ctx, managementGroupName, policyTrackedResourcesResource, queryOptions, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listQueryResultsForManagementGroupHandleResponse(resp)
		},
	})
}

// listQueryResultsForManagementGroupCreateRequest creates the ListQueryResultsForManagementGroup request.
func (client *PolicyTrackedResourcesClient) listQueryResultsForManagementGroupCreateRequest(ctx context.Context, managementGroupName string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForManagementGroupOptions) (*policy.Request, error) {
	urlPath := "/providers/{managementGroupsNamespace}/managementGroups/{managementGroupName}/providers/Microsoft.PolicyInsights/policyTrackedResources/{policyTrackedResourcesResource}/queryResults"
	urlPath = strings.ReplaceAll(urlPath, "{managementGroupsNamespace}", url.PathEscape("Microsoft.Management"))
	if managementGroupName == "" {
		return nil, errors.New("parameter managementGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managementGroupName}", url.PathEscape(managementGroupName))
	if policyTrackedResourcesResource == "" {
		return nil, errors.New("parameter policyTrackedResourcesResource cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyTrackedResourcesResource}", url.PathEscape(string(policyTrackedResourcesResource)))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if queryOptions != nil && queryOptions.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*queryOptions.Top), 10))
	}
	if queryOptions != nil && queryOptions.Filter != nil {
		reqQP.Set("$filter", *queryOptions.Filter)
	}
	reqQP.Set("api-version", "2018-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listQueryResultsForManagementGroupHandleResponse handles the ListQueryResultsForManagementGroup response.
func (client *PolicyTrackedResourcesClient) listQueryResultsForManagementGroupHandleResponse(resp *http.Response) (PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse, error) {
	result := PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyTrackedResourcesQueryResults); err != nil {
		return PolicyTrackedResourcesClientListQueryResultsForManagementGroupResponse{}, err
	}
	return result, nil
}

// NewListQueryResultsForResourcePager - Queries policy tracked resources under the resource.
//
// Generated from API version 2018-07-01-preview
//   - resourceID - Resource ID.
//   - policyTrackedResourcesResource - The name of the virtual resource under PolicyTrackedResources resource type; only "default"
//     is allowed.
//   - QueryOptions - QueryOptions contains a group of parameters for the PolicyTrackedResourcesClient.ListQueryResultsForManagementGroup
//     method.
//   - options - PolicyTrackedResourcesClientListQueryResultsForResourceOptions contains the optional parameters for the PolicyTrackedResourcesClient.NewListQueryResultsForResourcePager
//     method.
func (client *PolicyTrackedResourcesClient) NewListQueryResultsForResourcePager(resourceID string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForResourceOptions) (*runtime.Pager[PolicyTrackedResourcesClientListQueryResultsForResourceResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PolicyTrackedResourcesClientListQueryResultsForResourceResponse]{
		More: func(page PolicyTrackedResourcesClientListQueryResultsForResourceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PolicyTrackedResourcesClientListQueryResultsForResourceResponse) (PolicyTrackedResourcesClientListQueryResultsForResourceResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listQueryResultsForResourceCreateRequest(ctx, resourceID, policyTrackedResourcesResource, queryOptions, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PolicyTrackedResourcesClientListQueryResultsForResourceResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PolicyTrackedResourcesClientListQueryResultsForResourceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PolicyTrackedResourcesClientListQueryResultsForResourceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listQueryResultsForResourceHandleResponse(resp)
		},
	})
}

// listQueryResultsForResourceCreateRequest creates the ListQueryResultsForResource request.
func (client *PolicyTrackedResourcesClient) listQueryResultsForResourceCreateRequest(ctx context.Context, resourceID string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForResourceOptions) (*policy.Request, error) {
	urlPath := "/{resourceId}/providers/Microsoft.PolicyInsights/policyTrackedResources/{policyTrackedResourcesResource}/queryResults"
	urlPath = strings.ReplaceAll(urlPath, "{resourceId}", resourceID)
	if policyTrackedResourcesResource == "" {
		return nil, errors.New("parameter policyTrackedResourcesResource cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyTrackedResourcesResource}", url.PathEscape(string(policyTrackedResourcesResource)))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if queryOptions != nil && queryOptions.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*queryOptions.Top), 10))
	}
	if queryOptions != nil && queryOptions.Filter != nil {
		reqQP.Set("$filter", *queryOptions.Filter)
	}
	reqQP.Set("api-version", "2018-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listQueryResultsForResourceHandleResponse handles the ListQueryResultsForResource response.
func (client *PolicyTrackedResourcesClient) listQueryResultsForResourceHandleResponse(resp *http.Response) (PolicyTrackedResourcesClientListQueryResultsForResourceResponse, error) {
	result := PolicyTrackedResourcesClientListQueryResultsForResourceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyTrackedResourcesQueryResults); err != nil {
		return PolicyTrackedResourcesClientListQueryResultsForResourceResponse{}, err
	}
	return result, nil
}

// NewListQueryResultsForResourceGroupPager - Queries policy tracked resources under the resource group.
//
// Generated from API version 2018-07-01-preview
//   - resourceGroupName - Resource group name.
//   - policyTrackedResourcesResource - The name of the virtual resource under PolicyTrackedResources resource type; only "default"
//     is allowed.
//   - QueryOptions - QueryOptions contains a group of parameters for the PolicyTrackedResourcesClient.ListQueryResultsForManagementGroup
//     method.
//   - options - PolicyTrackedResourcesClientListQueryResultsForResourceGroupOptions contains the optional parameters for the
//     PolicyTrackedResourcesClient.NewListQueryResultsForResourceGroupPager method.
func (client *PolicyTrackedResourcesClient) NewListQueryResultsForResourceGroupPager(resourceGroupName string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForResourceGroupOptions) (*runtime.Pager[PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse]{
		More: func(page PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse) (PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listQueryResultsForResourceGroupCreateRequest(ctx, resourceGroupName, policyTrackedResourcesResource, queryOptions, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listQueryResultsForResourceGroupHandleResponse(resp)
		},
	})
}

// listQueryResultsForResourceGroupCreateRequest creates the ListQueryResultsForResourceGroup request.
func (client *PolicyTrackedResourcesClient) listQueryResultsForResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PolicyInsights/policyTrackedResources/{policyTrackedResourcesResource}/queryResults"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if policyTrackedResourcesResource == "" {
		return nil, errors.New("parameter policyTrackedResourcesResource cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyTrackedResourcesResource}", url.PathEscape(string(policyTrackedResourcesResource)))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if queryOptions != nil && queryOptions.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*queryOptions.Top), 10))
	}
	if queryOptions != nil && queryOptions.Filter != nil {
		reqQP.Set("$filter", *queryOptions.Filter)
	}
	reqQP.Set("api-version", "2018-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listQueryResultsForResourceGroupHandleResponse handles the ListQueryResultsForResourceGroup response.
func (client *PolicyTrackedResourcesClient) listQueryResultsForResourceGroupHandleResponse(resp *http.Response) (PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse, error) {
	result := PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyTrackedResourcesQueryResults); err != nil {
		return PolicyTrackedResourcesClientListQueryResultsForResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListQueryResultsForSubscriptionPager - Queries policy tracked resources under the subscription.
//
// Generated from API version 2018-07-01-preview
//   - policyTrackedResourcesResource - The name of the virtual resource under PolicyTrackedResources resource type; only "default"
//     is allowed.
//   - QueryOptions - QueryOptions contains a group of parameters for the PolicyTrackedResourcesClient.ListQueryResultsForManagementGroup
//     method.
//   - options - PolicyTrackedResourcesClientListQueryResultsForSubscriptionOptions contains the optional parameters for the PolicyTrackedResourcesClient.NewListQueryResultsForSubscriptionPager
//     method.
func (client *PolicyTrackedResourcesClient) NewListQueryResultsForSubscriptionPager(policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForSubscriptionOptions) (*runtime.Pager[PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse]{
		More: func(page PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse) (PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listQueryResultsForSubscriptionCreateRequest(ctx, policyTrackedResourcesResource, queryOptions, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listQueryResultsForSubscriptionHandleResponse(resp)
		},
	})
}

// listQueryResultsForSubscriptionCreateRequest creates the ListQueryResultsForSubscription request.
func (client *PolicyTrackedResourcesClient) listQueryResultsForSubscriptionCreateRequest(ctx context.Context, policyTrackedResourcesResource PolicyTrackedResourcesResourceType, queryOptions *QueryOptions, options *PolicyTrackedResourcesClientListQueryResultsForSubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.PolicyInsights/policyTrackedResources/{policyTrackedResourcesResource}/queryResults"
	if policyTrackedResourcesResource == "" {
		return nil, errors.New("parameter policyTrackedResourcesResource cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyTrackedResourcesResource}", url.PathEscape(string(policyTrackedResourcesResource)))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if queryOptions != nil && queryOptions.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*queryOptions.Top), 10))
	}
	if queryOptions != nil && queryOptions.Filter != nil {
		reqQP.Set("$filter", *queryOptions.Filter)
	}
	reqQP.Set("api-version", "2018-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listQueryResultsForSubscriptionHandleResponse handles the ListQueryResultsForSubscription response.
func (client *PolicyTrackedResourcesClient) listQueryResultsForSubscriptionHandleResponse(resp *http.Response) (PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse, error) {
	result := PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyTrackedResourcesQueryResults); err != nil {
		return PolicyTrackedResourcesClientListQueryResultsForSubscriptionResponse{}, err
	}
	return result, nil
}

