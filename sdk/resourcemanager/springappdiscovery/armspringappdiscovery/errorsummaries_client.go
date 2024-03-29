//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armspringappdiscovery

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

// ErrorSummariesClient contains the methods for the ErrorSummaries group.
// Don't use this type directly, use NewErrorSummariesClient() instead.
type ErrorSummariesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewErrorSummariesClient creates a new instance of ErrorSummariesClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewErrorSummariesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ErrorSummariesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ErrorSummariesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Gets the ErrorSummaries resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - siteName - The springbootsites name.
//   - errorSummaryName - The name of error summary
//   - options - ErrorSummariesClientGetOptions contains the optional parameters for the ErrorSummariesClient.Get method.
func (client *ErrorSummariesClient) Get(ctx context.Context, resourceGroupName string, siteName string, errorSummaryName string, options *ErrorSummariesClientGetOptions) (ErrorSummariesClientGetResponse, error) {
	var err error
	const operationName = "ErrorSummariesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, siteName, errorSummaryName, options)
	if err != nil {
		return ErrorSummariesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ErrorSummariesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ErrorSummariesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ErrorSummariesClient) getCreateRequest(ctx context.Context, resourceGroupName string, siteName string, errorSummaryName string, options *ErrorSummariesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OffAzureSpringBoot/springbootsites/{siteName}/errorSummaries/{errorSummaryName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if siteName == "" {
		return nil, errors.New("parameter siteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{siteName}", url.PathEscape(siteName))
	if errorSummaryName == "" {
		return nil, errors.New("parameter errorSummaryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{errorSummaryName}", url.PathEscape(errorSummaryName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ErrorSummariesClient) getHandleResponse(resp *http.Response) (ErrorSummariesClientGetResponse, error) {
	result := ErrorSummariesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ErrorSummary); err != nil {
		return ErrorSummariesClientGetResponse{}, err
	}
	return result, nil
}

// NewListBySitePager - Lists the ErrorSummaries resource in springbootsites.
//
// Generated from API version 2023-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - siteName - The springbootsites name.
//   - options - ErrorSummariesClientListBySiteOptions contains the optional parameters for the ErrorSummariesClient.NewListBySitePager
//     method.
func (client *ErrorSummariesClient) NewListBySitePager(resourceGroupName string, siteName string, options *ErrorSummariesClientListBySiteOptions) *runtime.Pager[ErrorSummariesClientListBySiteResponse] {
	return runtime.NewPager(runtime.PagingHandler[ErrorSummariesClientListBySiteResponse]{
		More: func(page ErrorSummariesClientListBySiteResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ErrorSummariesClientListBySiteResponse) (ErrorSummariesClientListBySiteResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ErrorSummariesClient.NewListBySitePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySiteCreateRequest(ctx, resourceGroupName, siteName, options)
			}, nil)
			if err != nil {
				return ErrorSummariesClientListBySiteResponse{}, err
			}
			return client.listBySiteHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySiteCreateRequest creates the ListBySite request.
func (client *ErrorSummariesClient) listBySiteCreateRequest(ctx context.Context, resourceGroupName string, siteName string, options *ErrorSummariesClientListBySiteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OffAzureSpringBoot/springbootsites/{siteName}/errorSummaries"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if siteName == "" {
		return nil, errors.New("parameter siteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{siteName}", url.PathEscape(siteName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySiteHandleResponse handles the ListBySite response.
func (client *ErrorSummariesClient) listBySiteHandleResponse(resp *http.Response) (ErrorSummariesClientListBySiteResponse, error) {
	result := ErrorSummariesClientListBySiteResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ErrorSummaryList); err != nil {
		return ErrorSummariesClientListBySiteResponse{}, err
	}
	return result, nil
}
