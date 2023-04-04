//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armsecuritydevops

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

// GitHubOwnerClient contains the methods for the GitHubOwner group.
// Don't use this type directly, use NewGitHubOwnerClient() instead.
type GitHubOwnerClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewGitHubOwnerClient creates a new instance of GitHubOwnerClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewGitHubOwnerClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*GitHubOwnerClient, error) {
	cl, err := arm.NewClient(moduleName+".GitHubOwnerClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &GitHubOwnerClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update a monitored GitHub owner.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - gitHubConnectorName - Name of the GitHub Connector.
//   - gitHubOwnerName - Name of the GitHub Owner.
//   - gitHubOwner - Github owner.
//   - options - GitHubOwnerClientBeginCreateOrUpdateOptions contains the optional parameters for the GitHubOwnerClient.BeginCreateOrUpdate
//     method.
func (client *GitHubOwnerClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, gitHubConnectorName string, gitHubOwnerName string, gitHubOwner GitHubOwner, options *GitHubOwnerClientBeginCreateOrUpdateOptions) (*runtime.Poller[GitHubOwnerClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, gitHubConnectorName, gitHubOwnerName, gitHubOwner, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[GitHubOwnerClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[GitHubOwnerClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create or update a monitored GitHub owner.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
func (client *GitHubOwnerClient) createOrUpdate(ctx context.Context, resourceGroupName string, gitHubConnectorName string, gitHubOwnerName string, gitHubOwner GitHubOwner, options *GitHubOwnerClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, gitHubConnectorName, gitHubOwnerName, gitHubOwner, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *GitHubOwnerClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, gitHubConnectorName string, gitHubOwnerName string, gitHubOwner GitHubOwner, options *GitHubOwnerClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/gitHubConnectors/{gitHubConnectorName}/owners/{gitHubOwnerName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if gitHubConnectorName == "" {
		return nil, errors.New("parameter gitHubConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gitHubConnectorName}", url.PathEscape(gitHubConnectorName))
	if gitHubOwnerName == "" {
		return nil, errors.New("parameter gitHubOwnerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gitHubOwnerName}", url.PathEscape(gitHubOwnerName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, gitHubOwner)
}

// Get - Returns a monitored GitHub repository.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - gitHubConnectorName - Name of the GitHub Connector.
//   - gitHubOwnerName - Name of the GitHub Owner.
//   - options - GitHubOwnerClientGetOptions contains the optional parameters for the GitHubOwnerClient.Get method.
func (client *GitHubOwnerClient) Get(ctx context.Context, resourceGroupName string, gitHubConnectorName string, gitHubOwnerName string, options *GitHubOwnerClientGetOptions) (GitHubOwnerClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, gitHubConnectorName, gitHubOwnerName, options)
	if err != nil {
		return GitHubOwnerClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GitHubOwnerClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return GitHubOwnerClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *GitHubOwnerClient) getCreateRequest(ctx context.Context, resourceGroupName string, gitHubConnectorName string, gitHubOwnerName string, options *GitHubOwnerClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/gitHubConnectors/{gitHubConnectorName}/owners/{gitHubOwnerName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if gitHubConnectorName == "" {
		return nil, errors.New("parameter gitHubConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gitHubConnectorName}", url.PathEscape(gitHubConnectorName))
	if gitHubOwnerName == "" {
		return nil, errors.New("parameter gitHubOwnerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gitHubOwnerName}", url.PathEscape(gitHubOwnerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *GitHubOwnerClient) getHandleResponse(resp *http.Response) (GitHubOwnerClientGetResponse, error) {
	result := GitHubOwnerClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GitHubOwner); err != nil {
		return GitHubOwnerClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Returns a list of monitored GitHub owners.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - gitHubConnectorName - Name of the GitHub Connector.
//   - options - GitHubOwnerClientListOptions contains the optional parameters for the GitHubOwnerClient.NewListPager method.
func (client *GitHubOwnerClient) NewListPager(resourceGroupName string, gitHubConnectorName string, options *GitHubOwnerClientListOptions) *runtime.Pager[GitHubOwnerClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[GitHubOwnerClientListResponse]{
		More: func(page GitHubOwnerClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *GitHubOwnerClientListResponse) (GitHubOwnerClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, gitHubConnectorName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return GitHubOwnerClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return GitHubOwnerClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return GitHubOwnerClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *GitHubOwnerClient) listCreateRequest(ctx context.Context, resourceGroupName string, gitHubConnectorName string, options *GitHubOwnerClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/gitHubConnectors/{gitHubConnectorName}/owners"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if gitHubConnectorName == "" {
		return nil, errors.New("parameter gitHubConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gitHubConnectorName}", url.PathEscape(gitHubConnectorName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *GitHubOwnerClient) listHandleResponse(resp *http.Response) (GitHubOwnerClientListResponse, error) {
	result := GitHubOwnerClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GitHubOwnerListResponse); err != nil {
		return GitHubOwnerClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Patch a monitored GitHub repository.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - gitHubConnectorName - Name of the GitHub Connector.
//   - gitHubOwnerName - Name of the GitHub Owner.
//   - gitHubOwner - Github owner.
//   - options - GitHubOwnerClientBeginUpdateOptions contains the optional parameters for the GitHubOwnerClient.BeginUpdate method.
func (client *GitHubOwnerClient) BeginUpdate(ctx context.Context, resourceGroupName string, gitHubConnectorName string, gitHubOwnerName string, gitHubOwner GitHubOwner, options *GitHubOwnerClientBeginUpdateOptions) (*runtime.Poller[GitHubOwnerClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, gitHubConnectorName, gitHubOwnerName, gitHubOwner, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[GitHubOwnerClientUpdateResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[GitHubOwnerClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Patch a monitored GitHub repository.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
func (client *GitHubOwnerClient) update(ctx context.Context, resourceGroupName string, gitHubConnectorName string, gitHubOwnerName string, gitHubOwner GitHubOwner, options *GitHubOwnerClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, gitHubConnectorName, gitHubOwnerName, gitHubOwner, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *GitHubOwnerClient) updateCreateRequest(ctx context.Context, resourceGroupName string, gitHubConnectorName string, gitHubOwnerName string, gitHubOwner GitHubOwner, options *GitHubOwnerClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SecurityDevOps/gitHubConnectors/{gitHubConnectorName}/owners/{gitHubOwnerName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if gitHubConnectorName == "" {
		return nil, errors.New("parameter gitHubConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gitHubConnectorName}", url.PathEscape(gitHubConnectorName))
	if gitHubOwnerName == "" {
		return nil, errors.New("parameter gitHubOwnerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gitHubOwnerName}", url.PathEscape(gitHubOwnerName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, gitHubOwner)
}
