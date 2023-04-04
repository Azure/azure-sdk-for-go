//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcdn

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

// PoliciesClient contains the methods for the Policies group.
// Don't use this type directly, use NewPoliciesClient() instead.
type PoliciesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewPoliciesClient creates a new instance of PoliciesClient with the specified values.
//   - subscriptionID - Azure Subscription ID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPoliciesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PoliciesClient, error) {
	cl, err := arm.NewClient(moduleName+".PoliciesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PoliciesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update policy with specified rule set name within a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - Name of the Resource group within the Azure subscription.
//   - policyName - The name of the CdnWebApplicationFirewallPolicy.
//   - cdnWebApplicationFirewallPolicy - Policy to be created.
//   - options - PoliciesClientBeginCreateOrUpdateOptions contains the optional parameters for the PoliciesClient.BeginCreateOrUpdate
//     method.
func (client *PoliciesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, policyName string, cdnWebApplicationFirewallPolicy WebApplicationFirewallPolicy, options *PoliciesClientBeginCreateOrUpdateOptions) (*runtime.Poller[PoliciesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, policyName, cdnWebApplicationFirewallPolicy, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[PoliciesClientCreateOrUpdateResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[PoliciesClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create or update policy with specified rule set name within a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
func (client *PoliciesClient) createOrUpdate(ctx context.Context, resourceGroupName string, policyName string, cdnWebApplicationFirewallPolicy WebApplicationFirewallPolicy, options *PoliciesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, policyName, cdnWebApplicationFirewallPolicy, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *PoliciesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, policyName string, cdnWebApplicationFirewallPolicy WebApplicationFirewallPolicy, options *PoliciesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies/{policyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if policyName == "" {
		return nil, errors.New("parameter policyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyName}", url.PathEscape(policyName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, cdnWebApplicationFirewallPolicy)
}

// Delete - Deletes Policy
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - Name of the Resource group within the Azure subscription.
//   - policyName - The name of the CdnWebApplicationFirewallPolicy.
//   - options - PoliciesClientDeleteOptions contains the optional parameters for the PoliciesClient.Delete method.
func (client *PoliciesClient) Delete(ctx context.Context, resourceGroupName string, policyName string, options *PoliciesClientDeleteOptions) (PoliciesClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, policyName, options)
	if err != nil {
		return PoliciesClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PoliciesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return PoliciesClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return PoliciesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *PoliciesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, policyName string, options *PoliciesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies/{policyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if policyName == "" {
		return nil, errors.New("parameter policyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyName}", url.PathEscape(policyName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Retrieve protection policy with specified name within a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - Name of the Resource group within the Azure subscription.
//   - policyName - The name of the CdnWebApplicationFirewallPolicy.
//   - options - PoliciesClientGetOptions contains the optional parameters for the PoliciesClient.Get method.
func (client *PoliciesClient) Get(ctx context.Context, resourceGroupName string, policyName string, options *PoliciesClientGetOptions) (PoliciesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, policyName, options)
	if err != nil {
		return PoliciesClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PoliciesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PoliciesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *PoliciesClient) getCreateRequest(ctx context.Context, resourceGroupName string, policyName string, options *PoliciesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies/{policyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if policyName == "" {
		return nil, errors.New("parameter policyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyName}", url.PathEscape(policyName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PoliciesClient) getHandleResponse(resp *http.Response) (PoliciesClientGetResponse, error) {
	result := PoliciesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WebApplicationFirewallPolicy); err != nil {
		return PoliciesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all of the protection policies within a resource group.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - Name of the Resource group within the Azure subscription.
//   - options - PoliciesClientListOptions contains the optional parameters for the PoliciesClient.NewListPager method.
func (client *PoliciesClient) NewListPager(resourceGroupName string, options *PoliciesClientListOptions) *runtime.Pager[PoliciesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[PoliciesClientListResponse]{
		More: func(page PoliciesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PoliciesClientListResponse) (PoliciesClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PoliciesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PoliciesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PoliciesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *PoliciesClient) listCreateRequest(ctx context.Context, resourceGroupName string, options *PoliciesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies"
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
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *PoliciesClient) listHandleResponse(resp *http.Response) (PoliciesClientListResponse, error) {
	result := PoliciesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WebApplicationFirewallPolicyList); err != nil {
		return PoliciesClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update an existing CdnWebApplicationFirewallPolicy with the specified policy name under the specified subscription
// and resource group
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - Name of the Resource group within the Azure subscription.
//   - policyName - The name of the CdnWebApplicationFirewallPolicy.
//   - cdnWebApplicationFirewallPolicyPatchParameters - CdnWebApplicationFirewallPolicy parameters to be patched.
//   - options - PoliciesClientBeginUpdateOptions contains the optional parameters for the PoliciesClient.BeginUpdate method.
func (client *PoliciesClient) BeginUpdate(ctx context.Context, resourceGroupName string, policyName string, cdnWebApplicationFirewallPolicyPatchParameters WebApplicationFirewallPolicyPatchParameters, options *PoliciesClientBeginUpdateOptions) (*runtime.Poller[PoliciesClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, policyName, cdnWebApplicationFirewallPolicyPatchParameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[PoliciesClientUpdateResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[PoliciesClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Update an existing CdnWebApplicationFirewallPolicy with the specified policy name under the specified subscription
// and resource group
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
func (client *PoliciesClient) update(ctx context.Context, resourceGroupName string, policyName string, cdnWebApplicationFirewallPolicyPatchParameters WebApplicationFirewallPolicyPatchParameters, options *PoliciesClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, policyName, cdnWebApplicationFirewallPolicyPatchParameters, options)
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
func (client *PoliciesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, policyName string, cdnWebApplicationFirewallPolicyPatchParameters WebApplicationFirewallPolicyPatchParameters, options *PoliciesClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies/{policyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if policyName == "" {
		return nil, errors.New("parameter policyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyName}", url.PathEscape(policyName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, cdnWebApplicationFirewallPolicyPatchParameters)
}
