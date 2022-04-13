//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmediaservices

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

// StreamingLocatorsClient contains the methods for the StreamingLocators group.
// Don't use this type directly, use NewStreamingLocatorsClient() instead.
type StreamingLocatorsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewStreamingLocatorsClient creates a new instance of StreamingLocatorsClient with the specified values.
// subscriptionID - The unique identifier for a Microsoft Azure subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewStreamingLocatorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*StreamingLocatorsClient, error) {
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
	client := &StreamingLocatorsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Create - Create a Streaming Locator in the Media Services account
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the Azure subscription.
// accountName - The Media Services account name.
// streamingLocatorName - The Streaming Locator name.
// parameters - The request parameters
// options - StreamingLocatorsClientCreateOptions contains the optional parameters for the StreamingLocatorsClient.Create
// method.
func (client *StreamingLocatorsClient) Create(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, parameters StreamingLocator, options *StreamingLocatorsClientCreateOptions) (StreamingLocatorsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, accountName, streamingLocatorName, parameters, options)
	if err != nil {
		return StreamingLocatorsClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return StreamingLocatorsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated) {
		return StreamingLocatorsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *StreamingLocatorsClient) createCreateRequest(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, parameters StreamingLocator, options *StreamingLocatorsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/streamingLocators/{streamingLocatorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if streamingLocatorName == "" {
		return nil, errors.New("parameter streamingLocatorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{streamingLocatorName}", url.PathEscape(streamingLocatorName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createHandleResponse handles the Create response.
func (client *StreamingLocatorsClient) createHandleResponse(resp *http.Response) (StreamingLocatorsClientCreateResponse, error) {
	result := StreamingLocatorsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StreamingLocator); err != nil {
		return StreamingLocatorsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a Streaming Locator in the Media Services account
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the Azure subscription.
// accountName - The Media Services account name.
// streamingLocatorName - The Streaming Locator name.
// options - StreamingLocatorsClientDeleteOptions contains the optional parameters for the StreamingLocatorsClient.Delete
// method.
func (client *StreamingLocatorsClient) Delete(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, options *StreamingLocatorsClientDeleteOptions) (StreamingLocatorsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, streamingLocatorName, options)
	if err != nil {
		return StreamingLocatorsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return StreamingLocatorsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return StreamingLocatorsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return StreamingLocatorsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *StreamingLocatorsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, options *StreamingLocatorsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/streamingLocators/{streamingLocatorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if streamingLocatorName == "" {
		return nil, errors.New("parameter streamingLocatorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{streamingLocatorName}", url.PathEscape(streamingLocatorName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Get the details of a Streaming Locator in the Media Services account
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the Azure subscription.
// accountName - The Media Services account name.
// streamingLocatorName - The Streaming Locator name.
// options - StreamingLocatorsClientGetOptions contains the optional parameters for the StreamingLocatorsClient.Get method.
func (client *StreamingLocatorsClient) Get(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, options *StreamingLocatorsClientGetOptions) (StreamingLocatorsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, streamingLocatorName, options)
	if err != nil {
		return StreamingLocatorsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return StreamingLocatorsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return StreamingLocatorsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *StreamingLocatorsClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, options *StreamingLocatorsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/streamingLocators/{streamingLocatorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if streamingLocatorName == "" {
		return nil, errors.New("parameter streamingLocatorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{streamingLocatorName}", url.PathEscape(streamingLocatorName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *StreamingLocatorsClient) getHandleResponse(resp *http.Response) (StreamingLocatorsClientGetResponse, error) {
	result := StreamingLocatorsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StreamingLocator); err != nil {
		return StreamingLocatorsClientGetResponse{}, err
	}
	return result, nil
}

// List - Lists the Streaming Locators in the account
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the Azure subscription.
// accountName - The Media Services account name.
// options - StreamingLocatorsClientListOptions contains the optional parameters for the StreamingLocatorsClient.List method.
func (client *StreamingLocatorsClient) List(resourceGroupName string, accountName string, options *StreamingLocatorsClientListOptions) *runtime.Pager[StreamingLocatorsClientListResponse] {
	return runtime.NewPager(runtime.PageProcessor[StreamingLocatorsClientListResponse]{
		More: func(page StreamingLocatorsClientListResponse) bool {
			return page.ODataNextLink != nil && len(*page.ODataNextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *StreamingLocatorsClientListResponse) (StreamingLocatorsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, accountName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.ODataNextLink)
			}
			if err != nil {
				return StreamingLocatorsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return StreamingLocatorsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return StreamingLocatorsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *StreamingLocatorsClient) listCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *StreamingLocatorsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/streamingLocators"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *StreamingLocatorsClient) listHandleResponse(resp *http.Response) (StreamingLocatorsClientListResponse, error) {
	result := StreamingLocatorsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StreamingLocatorCollection); err != nil {
		return StreamingLocatorsClientListResponse{}, err
	}
	return result, nil
}

// ListContentKeys - List Content Keys used by this Streaming Locator
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the Azure subscription.
// accountName - The Media Services account name.
// streamingLocatorName - The Streaming Locator name.
// options - StreamingLocatorsClientListContentKeysOptions contains the optional parameters for the StreamingLocatorsClient.ListContentKeys
// method.
func (client *StreamingLocatorsClient) ListContentKeys(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, options *StreamingLocatorsClientListContentKeysOptions) (StreamingLocatorsClientListContentKeysResponse, error) {
	req, err := client.listContentKeysCreateRequest(ctx, resourceGroupName, accountName, streamingLocatorName, options)
	if err != nil {
		return StreamingLocatorsClientListContentKeysResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return StreamingLocatorsClientListContentKeysResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return StreamingLocatorsClientListContentKeysResponse{}, runtime.NewResponseError(resp)
	}
	return client.listContentKeysHandleResponse(resp)
}

// listContentKeysCreateRequest creates the ListContentKeys request.
func (client *StreamingLocatorsClient) listContentKeysCreateRequest(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, options *StreamingLocatorsClientListContentKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/streamingLocators/{streamingLocatorName}/listContentKeys"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if streamingLocatorName == "" {
		return nil, errors.New("parameter streamingLocatorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{streamingLocatorName}", url.PathEscape(streamingLocatorName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listContentKeysHandleResponse handles the ListContentKeys response.
func (client *StreamingLocatorsClient) listContentKeysHandleResponse(resp *http.Response) (StreamingLocatorsClientListContentKeysResponse, error) {
	result := StreamingLocatorsClientListContentKeysResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ListContentKeysResponse); err != nil {
		return StreamingLocatorsClientListContentKeysResponse{}, err
	}
	return result, nil
}

// ListPaths - List Paths supported by this Streaming Locator
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the Azure subscription.
// accountName - The Media Services account name.
// streamingLocatorName - The Streaming Locator name.
// options - StreamingLocatorsClientListPathsOptions contains the optional parameters for the StreamingLocatorsClient.ListPaths
// method.
func (client *StreamingLocatorsClient) ListPaths(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, options *StreamingLocatorsClientListPathsOptions) (StreamingLocatorsClientListPathsResponse, error) {
	req, err := client.listPathsCreateRequest(ctx, resourceGroupName, accountName, streamingLocatorName, options)
	if err != nil {
		return StreamingLocatorsClientListPathsResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return StreamingLocatorsClientListPathsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return StreamingLocatorsClientListPathsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listPathsHandleResponse(resp)
}

// listPathsCreateRequest creates the ListPaths request.
func (client *StreamingLocatorsClient) listPathsCreateRequest(ctx context.Context, resourceGroupName string, accountName string, streamingLocatorName string, options *StreamingLocatorsClientListPathsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/streamingLocators/{streamingLocatorName}/listPaths"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if streamingLocatorName == "" {
		return nil, errors.New("parameter streamingLocatorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{streamingLocatorName}", url.PathEscape(streamingLocatorName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listPathsHandleResponse handles the ListPaths response.
func (client *StreamingLocatorsClient) listPathsHandleResponse(resp *http.Response) (StreamingLocatorsClientListPathsResponse, error) {
	result := StreamingLocatorsClientListPathsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ListPathsResponse); err != nil {
		return StreamingLocatorsClientListPathsResponse{}, err
	}
	return result, nil
}
