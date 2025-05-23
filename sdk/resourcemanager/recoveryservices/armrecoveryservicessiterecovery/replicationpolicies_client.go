// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armrecoveryservicessiterecovery

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

// ReplicationPoliciesClient contains the methods for the ReplicationPolicies group.
// Don't use this type directly, use NewReplicationPoliciesClient() instead.
type ReplicationPoliciesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewReplicationPoliciesClient creates a new instance of ReplicationPoliciesClient with the specified values.
//   - subscriptionID - The subscription Id.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewReplicationPoliciesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ReplicationPoliciesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ReplicationPoliciesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreate - The operation to create a replication policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group where the recovery services vault is present.
//   - resourceName - The name of the recovery services vault.
//   - policyName - Replication policy name.
//   - input - Create policy input.
//   - options - ReplicationPoliciesClientBeginCreateOptions contains the optional parameters for the ReplicationPoliciesClient.BeginCreate
//     method.
func (client *ReplicationPoliciesClient) BeginCreate(ctx context.Context, resourceGroupName string, resourceName string, policyName string, input CreatePolicyInput, options *ReplicationPoliciesClientBeginCreateOptions) (*runtime.Poller[ReplicationPoliciesClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, resourceName, policyName, input, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReplicationPoliciesClientCreateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ReplicationPoliciesClientCreateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Create - The operation to create a replication policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
func (client *ReplicationPoliciesClient) create(ctx context.Context, resourceGroupName string, resourceName string, policyName string, input CreatePolicyInput, options *ReplicationPoliciesClientBeginCreateOptions) (*http.Response, error) {
	var err error
	const operationName = "ReplicationPoliciesClient.BeginCreate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createCreateRequest(ctx, resourceGroupName, resourceName, policyName, input, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createCreateRequest creates the Create request.
func (client *ReplicationPoliciesClient) createCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, policyName string, input CreatePolicyInput, _ *ReplicationPoliciesClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationPolicies/{policyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if policyName == "" {
		return nil, errors.New("parameter policyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyName}", url.PathEscape(policyName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, input); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - The operation to delete a replication policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group where the recovery services vault is present.
//   - resourceName - The name of the recovery services vault.
//   - policyName - Replication policy name.
//   - options - ReplicationPoliciesClientBeginDeleteOptions contains the optional parameters for the ReplicationPoliciesClient.BeginDelete
//     method.
func (client *ReplicationPoliciesClient) BeginDelete(ctx context.Context, resourceGroupName string, resourceName string, policyName string, options *ReplicationPoliciesClientBeginDeleteOptions) (*runtime.Poller[ReplicationPoliciesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, resourceName, policyName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReplicationPoliciesClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ReplicationPoliciesClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - The operation to delete a replication policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
func (client *ReplicationPoliciesClient) deleteOperation(ctx context.Context, resourceGroupName string, resourceName string, policyName string, options *ReplicationPoliciesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ReplicationPoliciesClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, policyName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ReplicationPoliciesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, policyName string, _ *ReplicationPoliciesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationPolicies/{policyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if policyName == "" {
		return nil, errors.New("parameter policyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyName}", url.PathEscape(policyName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Gets the details of a replication policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group where the recovery services vault is present.
//   - resourceName - The name of the recovery services vault.
//   - policyName - Replication policy name.
//   - options - ReplicationPoliciesClientGetOptions contains the optional parameters for the ReplicationPoliciesClient.Get method.
func (client *ReplicationPoliciesClient) Get(ctx context.Context, resourceGroupName string, resourceName string, policyName string, options *ReplicationPoliciesClientGetOptions) (ReplicationPoliciesClientGetResponse, error) {
	var err error
	const operationName = "ReplicationPoliciesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, policyName, options)
	if err != nil {
		return ReplicationPoliciesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ReplicationPoliciesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ReplicationPoliciesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ReplicationPoliciesClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, policyName string, _ *ReplicationPoliciesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationPolicies/{policyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if policyName == "" {
		return nil, errors.New("parameter policyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyName}", url.PathEscape(policyName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ReplicationPoliciesClient) getHandleResponse(resp *http.Response) (ReplicationPoliciesClientGetResponse, error) {
	result := ReplicationPoliciesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Policy); err != nil {
		return ReplicationPoliciesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists the replication policies for a vault.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group where the recovery services vault is present.
//   - resourceName - The name of the recovery services vault.
//   - options - ReplicationPoliciesClientListOptions contains the optional parameters for the ReplicationPoliciesClient.NewListPager
//     method.
func (client *ReplicationPoliciesClient) NewListPager(resourceGroupName string, resourceName string, options *ReplicationPoliciesClientListOptions) *runtime.Pager[ReplicationPoliciesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReplicationPoliciesClientListResponse]{
		More: func(page ReplicationPoliciesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReplicationPoliciesClientListResponse) (ReplicationPoliciesClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ReplicationPoliciesClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, resourceName, options)
			}, nil)
			if err != nil {
				return ReplicationPoliciesClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *ReplicationPoliciesClient) listCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, _ *ReplicationPoliciesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationPolicies"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ReplicationPoliciesClient) listHandleResponse(resp *http.Response) (ReplicationPoliciesClientListResponse, error) {
	result := ReplicationPoliciesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyCollection); err != nil {
		return ReplicationPoliciesClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - The operation to update a replication policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group where the recovery services vault is present.
//   - resourceName - The name of the recovery services vault.
//   - policyName - Policy Id.
//   - input - Update Policy Input.
//   - options - ReplicationPoliciesClientBeginUpdateOptions contains the optional parameters for the ReplicationPoliciesClient.BeginUpdate
//     method.
func (client *ReplicationPoliciesClient) BeginUpdate(ctx context.Context, resourceGroupName string, resourceName string, policyName string, input UpdatePolicyInput, options *ReplicationPoliciesClientBeginUpdateOptions) (*runtime.Poller[ReplicationPoliciesClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, resourceName, policyName, input, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReplicationPoliciesClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ReplicationPoliciesClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - The operation to update a replication policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
func (client *ReplicationPoliciesClient) update(ctx context.Context, resourceGroupName string, resourceName string, policyName string, input UpdatePolicyInput, options *ReplicationPoliciesClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ReplicationPoliciesClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, policyName, input, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *ReplicationPoliciesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, policyName string, input UpdatePolicyInput, _ *ReplicationPoliciesClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationPolicies/{policyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if policyName == "" {
		return nil, errors.New("parameter policyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyName}", url.PathEscape(policyName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, input); err != nil {
		return nil, err
	}
	return req, nil
}
