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

// PrefixListGlobalRulestackClient contains the methods for the PrefixListGlobalRulestack group.
// Don't use this type directly, use NewPrefixListGlobalRulestackClient() instead.
type PrefixListGlobalRulestackClient struct {
	internal *arm.Client
}

// NewPrefixListGlobalRulestackClient creates a new instance of PrefixListGlobalRulestackClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPrefixListGlobalRulestackClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*PrefixListGlobalRulestackClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PrefixListGlobalRulestackClient{
		internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a PrefixListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
//   - globalRulestackName - GlobalRulestack resource name
//   - name - Local Rule priority
//   - resource - Resource create parameters.
//   - options - PrefixListGlobalRulestackClientBeginCreateOrUpdateOptions contains the optional parameters for the PrefixListGlobalRulestackClient.BeginCreateOrUpdate
//     method.
func (client *PrefixListGlobalRulestackClient) BeginCreateOrUpdate(ctx context.Context, globalRulestackName string, name string, resource PrefixListGlobalRulestackResource, options *PrefixListGlobalRulestackClientBeginCreateOrUpdateOptions) (*runtime.Poller[PrefixListGlobalRulestackClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, globalRulestackName, name, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PrefixListGlobalRulestackClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[PrefixListGlobalRulestackClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create a PrefixListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
func (client *PrefixListGlobalRulestackClient) createOrUpdate(ctx context.Context, globalRulestackName string, name string, resource PrefixListGlobalRulestackResource, options *PrefixListGlobalRulestackClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "PrefixListGlobalRulestackClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
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
func (client *PrefixListGlobalRulestackClient) createOrUpdateCreateRequest(ctx context.Context, globalRulestackName string, name string, resource PrefixListGlobalRulestackResource, options *PrefixListGlobalRulestackClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/prefixlists/{name}"
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
	reqQP.Set("api-version", "2023-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete a PrefixListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
//   - globalRulestackName - GlobalRulestack resource name
//   - name - Local Rule priority
//   - options - PrefixListGlobalRulestackClientBeginDeleteOptions contains the optional parameters for the PrefixListGlobalRulestackClient.BeginDelete
//     method.
func (client *PrefixListGlobalRulestackClient) BeginDelete(ctx context.Context, globalRulestackName string, name string, options *PrefixListGlobalRulestackClientBeginDeleteOptions) (*runtime.Poller[PrefixListGlobalRulestackClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, globalRulestackName, name, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PrefixListGlobalRulestackClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[PrefixListGlobalRulestackClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a PrefixListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
func (client *PrefixListGlobalRulestackClient) deleteOperation(ctx context.Context, globalRulestackName string, name string, options *PrefixListGlobalRulestackClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "PrefixListGlobalRulestackClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
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
func (client *PrefixListGlobalRulestackClient) deleteCreateRequest(ctx context.Context, globalRulestackName string, name string, options *PrefixListGlobalRulestackClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/prefixlists/{name}"
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
	reqQP.Set("api-version", "2023-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a PrefixListGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
//   - globalRulestackName - GlobalRulestack resource name
//   - name - Local Rule priority
//   - options - PrefixListGlobalRulestackClientGetOptions contains the optional parameters for the PrefixListGlobalRulestackClient.Get
//     method.
func (client *PrefixListGlobalRulestackClient) Get(ctx context.Context, globalRulestackName string, name string, options *PrefixListGlobalRulestackClientGetOptions) (PrefixListGlobalRulestackClientGetResponse, error) {
	var err error
	const operationName = "PrefixListGlobalRulestackClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, globalRulestackName, name, options)
	if err != nil {
		return PrefixListGlobalRulestackClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PrefixListGlobalRulestackClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PrefixListGlobalRulestackClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PrefixListGlobalRulestackClient) getCreateRequest(ctx context.Context, globalRulestackName string, name string, options *PrefixListGlobalRulestackClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/prefixlists/{name}"
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
	reqQP.Set("api-version", "2023-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PrefixListGlobalRulestackClient) getHandleResponse(resp *http.Response) (PrefixListGlobalRulestackClientGetResponse, error) {
	result := PrefixListGlobalRulestackClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrefixListGlobalRulestackResource); err != nil {
		return PrefixListGlobalRulestackClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List PrefixListGlobalRulestackResource resources by Tenant
//
// Generated from API version 2023-09-01
//   - globalRulestackName - GlobalRulestack resource name
//   - options - PrefixListGlobalRulestackClientListOptions contains the optional parameters for the PrefixListGlobalRulestackClient.NewListPager
//     method.
func (client *PrefixListGlobalRulestackClient) NewListPager(globalRulestackName string, options *PrefixListGlobalRulestackClientListOptions) *runtime.Pager[PrefixListGlobalRulestackClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[PrefixListGlobalRulestackClientListResponse]{
		More: func(page PrefixListGlobalRulestackClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PrefixListGlobalRulestackClientListResponse) (PrefixListGlobalRulestackClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PrefixListGlobalRulestackClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, globalRulestackName, options)
			}, nil)
			if err != nil {
				return PrefixListGlobalRulestackClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *PrefixListGlobalRulestackClient) listCreateRequest(ctx context.Context, globalRulestackName string, options *PrefixListGlobalRulestackClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/prefixlists"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *PrefixListGlobalRulestackClient) listHandleResponse(resp *http.Response) (PrefixListGlobalRulestackClientListResponse, error) {
	result := PrefixListGlobalRulestackClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrefixListGlobalRulestackResourceListResult); err != nil {
		return PrefixListGlobalRulestackClientListResponse{}, err
	}
	return result, nil
}
