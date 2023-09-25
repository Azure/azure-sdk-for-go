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

// FqdnListGlobalRulestackClient contains the methods for the FqdnListGlobalRulestack group.
// Don't use this type directly, use NewFqdnListGlobalRulestackClient() instead.
type FqdnListGlobalRulestackClient struct {
	internal *arm.Client
}

// NewFqdnListGlobalRulestackClient creates a new instance of FqdnListGlobalRulestackClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewFqdnListGlobalRulestackClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*FqdnListGlobalRulestackClient, error) {
	cl, err := arm.NewClient(moduleName+".FqdnListGlobalRulestackClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &FqdnListGlobalRulestackClient{
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a FqdnListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - globalRulestackName - GlobalRulestack resource name
//   - name - fqdn list name
//   - resource - Resource create parameters.
//   - options - FqdnListGlobalRulestackClientBeginCreateOrUpdateOptions contains the optional parameters for the FqdnListGlobalRulestackClient.BeginCreateOrUpdate
//     method.
func (client *FqdnListGlobalRulestackClient) BeginCreateOrUpdate(ctx context.Context, globalRulestackName string, name string, resource FqdnListGlobalRulestackResource, options *FqdnListGlobalRulestackClientBeginCreateOrUpdateOptions) (*runtime.Poller[FqdnListGlobalRulestackClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, globalRulestackName, name, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[FqdnListGlobalRulestackClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[FqdnListGlobalRulestackClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create a FqdnListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
func (client *FqdnListGlobalRulestackClient) createOrUpdate(ctx context.Context, globalRulestackName string, name string, resource FqdnListGlobalRulestackResource, options *FqdnListGlobalRulestackClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, globalRulestackName, name, resource, options)
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
func (client *FqdnListGlobalRulestackClient) createOrUpdateCreateRequest(ctx context.Context, globalRulestackName string, name string, resource FqdnListGlobalRulestackResource, options *FqdnListGlobalRulestackClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/fqdnlists/{name}"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
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

// BeginDelete - Delete a FqdnListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - globalRulestackName - GlobalRulestack resource name
//   - name - fqdn list name
//   - options - FqdnListGlobalRulestackClientBeginDeleteOptions contains the optional parameters for the FqdnListGlobalRulestackClient.BeginDelete
//     method.
func (client *FqdnListGlobalRulestackClient) BeginDelete(ctx context.Context, globalRulestackName string, name string, options *FqdnListGlobalRulestackClientBeginDeleteOptions) (*runtime.Poller[FqdnListGlobalRulestackClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, globalRulestackName, name, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[FqdnListGlobalRulestackClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[FqdnListGlobalRulestackClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Delete a FqdnListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
func (client *FqdnListGlobalRulestackClient) deleteOperation(ctx context.Context, globalRulestackName string, name string, options *FqdnListGlobalRulestackClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, globalRulestackName, name, options)
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
func (client *FqdnListGlobalRulestackClient) deleteCreateRequest(ctx context.Context, globalRulestackName string, name string, options *FqdnListGlobalRulestackClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/fqdnlists/{name}"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
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

// Get - Get a FqdnListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - globalRulestackName - GlobalRulestack resource name
//   - name - fqdn list name
//   - options - FqdnListGlobalRulestackClientGetOptions contains the optional parameters for the FqdnListGlobalRulestackClient.Get
//     method.
func (client *FqdnListGlobalRulestackClient) Get(ctx context.Context, globalRulestackName string, name string, options *FqdnListGlobalRulestackClientGetOptions) (FqdnListGlobalRulestackClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, globalRulestackName, name, options)
	if err != nil {
		return FqdnListGlobalRulestackClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FqdnListGlobalRulestackClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FqdnListGlobalRulestackClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *FqdnListGlobalRulestackClient) getCreateRequest(ctx context.Context, globalRulestackName string, name string, options *FqdnListGlobalRulestackClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/fqdnlists/{name}"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
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
func (client *FqdnListGlobalRulestackClient) getHandleResponse(resp *http.Response) (FqdnListGlobalRulestackClientGetResponse, error) {
	result := FqdnListGlobalRulestackClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FqdnListGlobalRulestackResource); err != nil {
		return FqdnListGlobalRulestackClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List FqdnListGlobalRulestackResource resources by Tenant
//
// Generated from API version 2022-08-29
//   - globalRulestackName - GlobalRulestack resource name
//   - options - FqdnListGlobalRulestackClientListOptions contains the optional parameters for the FqdnListGlobalRulestackClient.NewListPager
//     method.
func (client *FqdnListGlobalRulestackClient) NewListPager(globalRulestackName string, options *FqdnListGlobalRulestackClientListOptions) (*runtime.Pager[FqdnListGlobalRulestackClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[FqdnListGlobalRulestackClientListResponse]{
		More: func(page FqdnListGlobalRulestackClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *FqdnListGlobalRulestackClientListResponse) (FqdnListGlobalRulestackClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, globalRulestackName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return FqdnListGlobalRulestackClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return FqdnListGlobalRulestackClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return FqdnListGlobalRulestackClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *FqdnListGlobalRulestackClient) listCreateRequest(ctx context.Context, globalRulestackName string, options *FqdnListGlobalRulestackClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/fqdnlists"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
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

// listHandleResponse handles the List response.
func (client *FqdnListGlobalRulestackClient) listHandleResponse(resp *http.Response) (FqdnListGlobalRulestackClientListResponse, error) {
	result := FqdnListGlobalRulestackClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FqdnListGlobalRulestackResourceListResult); err != nil {
		return FqdnListGlobalRulestackClientListResponse{}, err
	}
	return result, nil
}

