//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatashare

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
	"strings"
)

// SharesClient contains the methods for the Shares group.
// Don't use this type directly, use NewSharesClient() instead.
type SharesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewSharesClient creates a new instance of SharesClient with the specified values.
// subscriptionID - The subscription identifier
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewSharesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SharesClient, error) {
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
	client := &SharesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Create - Create a share
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share.
// share - The share payload
// options - SharesClientCreateOptions contains the optional parameters for the SharesClient.Create method.
func (client *SharesClient) Create(ctx context.Context, resourceGroupName string, accountName string, shareName string, share Share, options *SharesClientCreateOptions) (SharesClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, accountName, shareName, share, options)
	if err != nil {
		return SharesClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SharesClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return SharesClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *SharesClient) createCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, share Share, options *SharesClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}"
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
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, share)
}

// createHandleResponse handles the Create response.
func (client *SharesClient) createHandleResponse(resp *http.Response) (SharesClientCreateResponse, error) {
	result := SharesClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Share); err != nil {
		return SharesClientCreateResponse{}, err
	}
	return result, nil
}

// BeginDelete - Delete a share
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share.
// options - SharesClientBeginDeleteOptions contains the optional parameters for the SharesClient.BeginDelete method.
func (client *SharesClient) BeginDelete(ctx context.Context, resourceGroupName string, accountName string, shareName string, options *SharesClientBeginDeleteOptions) (*armruntime.Poller[SharesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, accountName, shareName, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[SharesClientDeleteResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[SharesClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Delete a share
// If the operation fails it returns an *azcore.ResponseError type.
func (client *SharesClient) deleteOperation(ctx context.Context, resourceGroupName string, accountName string, shareName string, options *SharesClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, shareName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SharesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, options *SharesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}"
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
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Get a share
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share to retrieve.
// options - SharesClientGetOptions contains the optional parameters for the SharesClient.Get method.
func (client *SharesClient) Get(ctx context.Context, resourceGroupName string, accountName string, shareName string, options *SharesClientGetOptions) (SharesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, shareName, options)
	if err != nil {
		return SharesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SharesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SharesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SharesClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, options *SharesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}"
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
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SharesClient) getHandleResponse(resp *http.Response) (SharesClientGetResponse, error) {
	result := SharesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Share); err != nil {
		return SharesClientGetResponse{}, err
	}
	return result, nil
}

// ListByAccount - List shares in an account
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// options - SharesClientListByAccountOptions contains the optional parameters for the SharesClient.ListByAccount method.
func (client *SharesClient) ListByAccount(resourceGroupName string, accountName string, options *SharesClientListByAccountOptions) *runtime.Pager[SharesClientListByAccountResponse] {
	return runtime.NewPager(runtime.PageProcessor[SharesClientListByAccountResponse]{
		More: func(page SharesClientListByAccountResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SharesClientListByAccountResponse) (SharesClientListByAccountResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByAccountCreateRequest(ctx, resourceGroupName, accountName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SharesClientListByAccountResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return SharesClientListByAccountResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SharesClientListByAccountResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByAccountHandleResponse(resp)
		},
	})
}

// listByAccountCreateRequest creates the ListByAccount request.
func (client *SharesClient) listByAccountCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *SharesClientListByAccountOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares"
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
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByAccountHandleResponse handles the ListByAccount response.
func (client *SharesClient) listByAccountHandleResponse(resp *http.Response) (SharesClientListByAccountResponse, error) {
	result := SharesClientListByAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ShareList); err != nil {
		return SharesClientListByAccountResponse{}, err
	}
	return result, nil
}

// ListSynchronizationDetails - List synchronization details
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share.
// shareSynchronization - Share Synchronization payload.
// options - SharesClientListSynchronizationDetailsOptions contains the optional parameters for the SharesClient.ListSynchronizationDetails
// method.
func (client *SharesClient) ListSynchronizationDetails(resourceGroupName string, accountName string, shareName string, shareSynchronization ShareSynchronization, options *SharesClientListSynchronizationDetailsOptions) *runtime.Pager[SharesClientListSynchronizationDetailsResponse] {
	return runtime.NewPager(runtime.PageProcessor[SharesClientListSynchronizationDetailsResponse]{
		More: func(page SharesClientListSynchronizationDetailsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SharesClientListSynchronizationDetailsResponse) (SharesClientListSynchronizationDetailsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listSynchronizationDetailsCreateRequest(ctx, resourceGroupName, accountName, shareName, shareSynchronization, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SharesClientListSynchronizationDetailsResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return SharesClientListSynchronizationDetailsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SharesClientListSynchronizationDetailsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listSynchronizationDetailsHandleResponse(resp)
		},
	})
}

// listSynchronizationDetailsCreateRequest creates the ListSynchronizationDetails request.
func (client *SharesClient) listSynchronizationDetailsCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, shareSynchronization ShareSynchronization, options *SharesClientListSynchronizationDetailsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/listSynchronizationDetails"
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
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, shareSynchronization)
}

// listSynchronizationDetailsHandleResponse handles the ListSynchronizationDetails response.
func (client *SharesClient) listSynchronizationDetailsHandleResponse(resp *http.Response) (SharesClientListSynchronizationDetailsResponse, error) {
	result := SharesClientListSynchronizationDetailsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SynchronizationDetailsList); err != nil {
		return SharesClientListSynchronizationDetailsResponse{}, err
	}
	return result, nil
}

// ListSynchronizations - List synchronizations of a share
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share.
// options - SharesClientListSynchronizationsOptions contains the optional parameters for the SharesClient.ListSynchronizations
// method.
func (client *SharesClient) ListSynchronizations(resourceGroupName string, accountName string, shareName string, options *SharesClientListSynchronizationsOptions) *runtime.Pager[SharesClientListSynchronizationsResponse] {
	return runtime.NewPager(runtime.PageProcessor[SharesClientListSynchronizationsResponse]{
		More: func(page SharesClientListSynchronizationsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SharesClientListSynchronizationsResponse) (SharesClientListSynchronizationsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listSynchronizationsCreateRequest(ctx, resourceGroupName, accountName, shareName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SharesClientListSynchronizationsResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return SharesClientListSynchronizationsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SharesClientListSynchronizationsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listSynchronizationsHandleResponse(resp)
		},
	})
}

// listSynchronizationsCreateRequest creates the ListSynchronizations request.
func (client *SharesClient) listSynchronizationsCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, options *SharesClientListSynchronizationsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/listSynchronizations"
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
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listSynchronizationsHandleResponse handles the ListSynchronizations response.
func (client *SharesClient) listSynchronizationsHandleResponse(resp *http.Response) (SharesClientListSynchronizationsResponse, error) {
	result := SharesClientListSynchronizationsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ShareSynchronizationList); err != nil {
		return SharesClientListSynchronizationsResponse{}, err
	}
	return result, nil
}
