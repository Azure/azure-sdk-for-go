// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkeyvault

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SecretsClient contains the methods for the Secrets group.
// Don't use this type directly, use NewSecretsClient() instead.
type SecretsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewSecretsClient creates a new instance of SecretsClient with the specified values.
//   - subscriptionID - Subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSecretsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SecretsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SecretsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update a secret in a key vault in the specified subscription. NOTE: This API is intended for
// internal use in ARM deployments. Users should use the data-plane REST service for interaction
// with vault secrets.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the Resource Group to which the vault belongs.
//   - vaultName - Name of the vault
//   - secretName - Name of the secret. The value you provide may be copied globally for the purpose of running the service. The
//     value provided should not include personally identifiable or sensitive information.
//   - parameters - Parameters to create or update the secret
//   - options - SecretsClientCreateOrUpdateOptions contains the optional parameters for the SecretsClient.CreateOrUpdate method.
func (client *SecretsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, secretName string, parameters SecretCreateOrUpdateParameters, options *SecretsClientCreateOrUpdateOptions) (SecretsClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "SecretsClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, vaultName, secretName, parameters, options)
	if err != nil {
		return SecretsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SecretsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return SecretsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *SecretsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, secretName string, parameters SecretCreateOrUpdateParameters, _ *SecretsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}/secrets/{secretName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if secretName == "" {
		return nil, errors.New("parameter secretName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretName}", url.PathEscape(secretName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *SecretsClient) createOrUpdateHandleResponse(resp *http.Response) (SecretsClientCreateOrUpdateResponse, error) {
	result := SecretsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return SecretsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Get - Gets the specified secret. NOTE: This API is intended for internal use in ARM deployments. Users should use the data-plane
// REST service for interaction with vault secrets.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the Resource Group to which the vault belongs.
//   - vaultName - The name of the vault.
//   - secretName - The name of the secret.
//   - options - SecretsClientGetOptions contains the optional parameters for the SecretsClient.Get method.
func (client *SecretsClient) Get(ctx context.Context, resourceGroupName string, vaultName string, secretName string, options *SecretsClientGetOptions) (SecretsClientGetResponse, error) {
	var err error
	const operationName = "SecretsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, vaultName, secretName, options)
	if err != nil {
		return SecretsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SecretsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return SecretsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *SecretsClient) getCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, secretName string, _ *SecretsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}/secrets/{secretName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if secretName == "" {
		return nil, errors.New("parameter secretName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretName}", url.PathEscape(secretName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SecretsClient) getHandleResponse(resp *http.Response) (SecretsClientGetResponse, error) {
	result := SecretsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return SecretsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - The List operation gets information about the secrets in a vault. NOTE: This API is intended for internal
// use in ARM deployments. Users should use the data-plane REST service for interaction with
// vault secrets.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the Resource Group to which the vault belongs.
//   - vaultName - The name of the vault.
//   - options - SecretsClientListOptions contains the optional parameters for the SecretsClient.NewListPager method.
func (client *SecretsClient) NewListPager(resourceGroupName string, vaultName string, options *SecretsClientListOptions) *runtime.Pager[SecretsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[SecretsClientListResponse]{
		More: func(page SecretsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SecretsClientListResponse) (SecretsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "SecretsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, vaultName, options)
			}, nil)
			if err != nil {
				return SecretsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *SecretsClient) listCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, options *SecretsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}/secrets"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SecretsClient) listHandleResponse(resp *http.Response) (SecretsClientListResponse, error) {
	result := SecretsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretListResult); err != nil {
		return SecretsClientListResponse{}, err
	}
	return result, nil
}

// Update - Update a secret in the specified subscription. NOTE: This API is intended for internal use in ARM deployments.
// Users should use the data-plane REST service for interaction with vault secrets.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the Resource Group to which the vault belongs.
//   - vaultName - Name of the vault
//   - secretName - Name of the secret
//   - parameters - Parameters to patch the secret
//   - options - SecretsClientUpdateOptions contains the optional parameters for the SecretsClient.Update method.
func (client *SecretsClient) Update(ctx context.Context, resourceGroupName string, vaultName string, secretName string, parameters SecretPatchParameters, options *SecretsClientUpdateOptions) (SecretsClientUpdateResponse, error) {
	var err error
	const operationName = "SecretsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, vaultName, secretName, parameters, options)
	if err != nil {
		return SecretsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SecretsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return SecretsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *SecretsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, secretName string, parameters SecretPatchParameters, _ *SecretsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{vaultName}/secrets/{secretName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if secretName == "" {
		return nil, errors.New("parameter secretName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secretName}", url.PathEscape(secretName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *SecretsClient) updateHandleResponse(resp *http.Response) (SecretsClientUpdateResponse, error) {
	result := SecretsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return SecretsClientUpdateResponse{}, err
	}
	return result, nil
}
