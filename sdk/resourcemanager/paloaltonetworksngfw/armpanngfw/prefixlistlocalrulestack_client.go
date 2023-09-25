//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpanngfw

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

// PrefixListLocalRulestackClient contains the methods for the PrefixListLocalRulestack group.
// Don't use this type directly, use NewPrefixListLocalRulestackClient() instead.
type PrefixListLocalRulestackClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewPrefixListLocalRulestackClient creates a new instance of PrefixListLocalRulestackClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPrefixListLocalRulestackClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PrefixListLocalRulestackClient, error) {
	cl, err := arm.NewClient(moduleName+".PrefixListLocalRulestackClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PrefixListLocalRulestackClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a PrefixListResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - localRulestackName - LocalRulestack resource name
//   - name - Local Rule priority
//   - resource - Resource create parameters.
//   - options - PrefixListLocalRulestackClientBeginCreateOrUpdateOptions contains the optional parameters for the PrefixListLocalRulestackClient.BeginCreateOrUpdate
//     method.
func (client *PrefixListLocalRulestackClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, localRulestackName string, name string, resource PrefixListResource, options *PrefixListLocalRulestackClientBeginCreateOrUpdateOptions) (*runtime.Poller[PrefixListLocalRulestackClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, localRulestackName, name, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PrefixListLocalRulestackClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[PrefixListLocalRulestackClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create a PrefixListResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
func (client *PrefixListLocalRulestackClient) createOrUpdate(ctx context.Context, resourceGroupName string, localRulestackName string, name string, resource PrefixListResource, options *PrefixListLocalRulestackClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, localRulestackName, name, resource, options)
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
func (client *PrefixListLocalRulestackClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, localRulestackName string, name string, resource PrefixListResource, options *PrefixListLocalRulestackClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/{localRulestackName}/prefixlists/{name}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if localRulestackName == "" {
		return nil, errors.New("parameter localRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{localRulestackName}", url.PathEscape(localRulestackName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Delete a PrefixListResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - localRulestackName - LocalRulestack resource name
//   - name - Local Rule priority
//   - options - PrefixListLocalRulestackClientBeginDeleteOptions contains the optional parameters for the PrefixListLocalRulestackClient.BeginDelete
//     method.
func (client *PrefixListLocalRulestackClient) BeginDelete(ctx context.Context, resourceGroupName string, localRulestackName string, name string, options *PrefixListLocalRulestackClientBeginDeleteOptions) (*runtime.Poller[PrefixListLocalRulestackClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, localRulestackName, name, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PrefixListLocalRulestackClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[PrefixListLocalRulestackClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Delete a PrefixListResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
func (client *PrefixListLocalRulestackClient) deleteOperation(ctx context.Context, resourceGroupName string, localRulestackName string, name string, options *PrefixListLocalRulestackClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, localRulestackName, name, options)
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
func (client *PrefixListLocalRulestackClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, localRulestackName string, name string, options *PrefixListLocalRulestackClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/{localRulestackName}/prefixlists/{name}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if localRulestackName == "" {
		return nil, errors.New("parameter localRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{localRulestackName}", url.PathEscape(localRulestackName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a PrefixListResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - localRulestackName - LocalRulestack resource name
//   - name - Local Rule priority
//   - options - PrefixListLocalRulestackClientGetOptions contains the optional parameters for the PrefixListLocalRulestackClient.Get
//     method.
func (client *PrefixListLocalRulestackClient) Get(ctx context.Context, resourceGroupName string, localRulestackName string, name string, options *PrefixListLocalRulestackClientGetOptions) (PrefixListLocalRulestackClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, localRulestackName, name, options)
	if err != nil {
		return PrefixListLocalRulestackClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PrefixListLocalRulestackClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PrefixListLocalRulestackClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PrefixListLocalRulestackClient) getCreateRequest(ctx context.Context, resourceGroupName string, localRulestackName string, name string, options *PrefixListLocalRulestackClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/{localRulestackName}/prefixlists/{name}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if localRulestackName == "" {
		return nil, errors.New("parameter localRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{localRulestackName}", url.PathEscape(localRulestackName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PrefixListLocalRulestackClient) getHandleResponse(resp *http.Response) (PrefixListLocalRulestackClientGetResponse, error) {
	result := PrefixListLocalRulestackClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrefixListResource); err != nil {
		return PrefixListLocalRulestackClientGetResponse{}, err
	}
	return result, nil
}

// NewListByLocalRulestacksPager - List PrefixListResource resources by LocalRulestacks
//
// Generated from API version 2022-08-29
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - localRulestackName - LocalRulestack resource name
//   - options - PrefixListLocalRulestackClientListByLocalRulestacksOptions contains the optional parameters for the PrefixListLocalRulestackClient.NewListByLocalRulestacksPager
//     method.
func (client *PrefixListLocalRulestackClient) NewListByLocalRulestacksPager(resourceGroupName string, localRulestackName string, options *PrefixListLocalRulestackClientListByLocalRulestacksOptions) (*runtime.Pager[PrefixListLocalRulestackClientListByLocalRulestacksResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PrefixListLocalRulestackClientListByLocalRulestacksResponse]{
		More: func(page PrefixListLocalRulestackClientListByLocalRulestacksResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PrefixListLocalRulestackClientListByLocalRulestacksResponse) (PrefixListLocalRulestackClientListByLocalRulestacksResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByLocalRulestacksCreateRequest(ctx, resourceGroupName, localRulestackName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PrefixListLocalRulestackClientListByLocalRulestacksResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PrefixListLocalRulestackClientListByLocalRulestacksResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PrefixListLocalRulestackClientListByLocalRulestacksResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByLocalRulestacksHandleResponse(resp)
		},
	})
}

// listByLocalRulestacksCreateRequest creates the ListByLocalRulestacks request.
func (client *PrefixListLocalRulestackClient) listByLocalRulestacksCreateRequest(ctx context.Context, resourceGroupName string, localRulestackName string, options *PrefixListLocalRulestackClientListByLocalRulestacksOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/{localRulestackName}/prefixlists"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if localRulestackName == "" {
		return nil, errors.New("parameter localRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{localRulestackName}", url.PathEscape(localRulestackName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByLocalRulestacksHandleResponse handles the ListByLocalRulestacks response.
func (client *PrefixListLocalRulestackClient) listByLocalRulestacksHandleResponse(resp *http.Response) (PrefixListLocalRulestackClientListByLocalRulestacksResponse, error) {
	result := PrefixListLocalRulestackClientListByLocalRulestacksResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrefixListResourceListResult); err != nil {
		return PrefixListLocalRulestackClientListByLocalRulestacksResponse{}, err
	}
	return result, nil
}

