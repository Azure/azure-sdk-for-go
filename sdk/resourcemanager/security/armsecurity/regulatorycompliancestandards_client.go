//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

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

// RegulatoryComplianceStandardsClient contains the methods for the RegulatoryComplianceStandards group.
// Don't use this type directly, use NewRegulatoryComplianceStandardsClient() instead.
type RegulatoryComplianceStandardsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewRegulatoryComplianceStandardsClient creates a new instance of RegulatoryComplianceStandardsClient with the specified values.
//   - subscriptionID - Azure subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewRegulatoryComplianceStandardsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RegulatoryComplianceStandardsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &RegulatoryComplianceStandardsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Supported regulatory compliance details state for selected standard
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-01-01-preview
//   - regulatoryComplianceStandardName - Name of the regulatory compliance standard object
//   - options - RegulatoryComplianceStandardsClientGetOptions contains the optional parameters for the RegulatoryComplianceStandardsClient.Get
//     method.
func (client *RegulatoryComplianceStandardsClient) Get(ctx context.Context, regulatoryComplianceStandardName string, options *RegulatoryComplianceStandardsClientGetOptions) (RegulatoryComplianceStandardsClientGetResponse, error) {
	var err error
	const operationName = "RegulatoryComplianceStandardsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, regulatoryComplianceStandardName, options)
	if err != nil {
		return RegulatoryComplianceStandardsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RegulatoryComplianceStandardsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RegulatoryComplianceStandardsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *RegulatoryComplianceStandardsClient) getCreateRequest(ctx context.Context, regulatoryComplianceStandardName string, options *RegulatoryComplianceStandardsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/regulatoryComplianceStandards/{regulatoryComplianceStandardName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if regulatoryComplianceStandardName == "" {
		return nil, errors.New("parameter regulatoryComplianceStandardName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{regulatoryComplianceStandardName}", url.PathEscape(regulatoryComplianceStandardName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RegulatoryComplianceStandardsClient) getHandleResponse(resp *http.Response) (RegulatoryComplianceStandardsClientGetResponse, error) {
	result := RegulatoryComplianceStandardsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RegulatoryComplianceStandard); err != nil {
		return RegulatoryComplianceStandardsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Supported regulatory compliance standards details and state
//
// Generated from API version 2019-01-01-preview
//   - options - RegulatoryComplianceStandardsClientListOptions contains the optional parameters for the RegulatoryComplianceStandardsClient.NewListPager
//     method.
func (client *RegulatoryComplianceStandardsClient) NewListPager(options *RegulatoryComplianceStandardsClientListOptions) *runtime.Pager[RegulatoryComplianceStandardsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[RegulatoryComplianceStandardsClientListResponse]{
		More: func(page RegulatoryComplianceStandardsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RegulatoryComplianceStandardsClientListResponse) (RegulatoryComplianceStandardsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "RegulatoryComplianceStandardsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return RegulatoryComplianceStandardsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *RegulatoryComplianceStandardsClient) listCreateRequest(ctx context.Context, options *RegulatoryComplianceStandardsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/regulatoryComplianceStandards"
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
	reqQP.Set("api-version", "2019-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *RegulatoryComplianceStandardsClient) listHandleResponse(resp *http.Response) (RegulatoryComplianceStandardsClientListResponse, error) {
	result := RegulatoryComplianceStandardsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RegulatoryComplianceStandardList); err != nil {
		return RegulatoryComplianceStandardsClientListResponse{}, err
	}
	return result, nil
}
