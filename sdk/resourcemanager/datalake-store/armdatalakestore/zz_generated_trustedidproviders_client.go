//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatalakestore

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

// TrustedIDProvidersClient contains the methods for the TrustedIDProviders group.
// Don't use this type directly, use NewTrustedIDProvidersClient() instead.
type TrustedIDProvidersClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewTrustedIDProvidersClient creates a new instance of TrustedIDProvidersClient with the specified values.
// subscriptionID - Gets subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID
// forms part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewTrustedIDProvidersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TrustedIDProvidersClient, error) {
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
	client := &TrustedIDProvidersClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates the specified trusted identity provider. During update, the trusted identity provider
// with the specified name will be replaced with this new provider
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the Azure resource group.
// accountName - The name of the Data Lake Store account.
// trustedIDProviderName - The name of the trusted identity provider. This is used for differentiation of providers in the
// account.
// parameters - Parameters supplied to create or replace the trusted identity provider.
// options - TrustedIDProvidersClientCreateOrUpdateOptions contains the optional parameters for the TrustedIDProvidersClient.CreateOrUpdate
// method.
func (client *TrustedIDProvidersClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, trustedIDProviderName string, parameters CreateOrUpdateTrustedIDProviderParameters, options *TrustedIDProvidersClientCreateOrUpdateOptions) (TrustedIDProvidersClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, accountName, trustedIDProviderName, parameters, options)
	if err != nil {
		return TrustedIDProvidersClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return TrustedIDProvidersClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return TrustedIDProvidersClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *TrustedIDProvidersClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, trustedIDProviderName string, parameters CreateOrUpdateTrustedIDProviderParameters, options *TrustedIDProvidersClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/trustedIdProviders/{trustedIdProviderName}"
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
	if trustedIDProviderName == "" {
		return nil, errors.New("parameter trustedIDProviderName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trustedIdProviderName}", url.PathEscape(trustedIDProviderName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *TrustedIDProvidersClient) createOrUpdateHandleResponse(resp *http.Response) (TrustedIDProvidersClientCreateOrUpdateResponse, error) {
	result := TrustedIDProvidersClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TrustedIDProvider); err != nil {
		return TrustedIDProvidersClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the specified trusted identity provider from the specified Data Lake Store account
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the Azure resource group.
// accountName - The name of the Data Lake Store account.
// trustedIDProviderName - The name of the trusted identity provider to delete.
// options - TrustedIDProvidersClientDeleteOptions contains the optional parameters for the TrustedIDProvidersClient.Delete
// method.
func (client *TrustedIDProvidersClient) Delete(ctx context.Context, resourceGroupName string, accountName string, trustedIDProviderName string, options *TrustedIDProvidersClientDeleteOptions) (TrustedIDProvidersClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, trustedIDProviderName, options)
	if err != nil {
		return TrustedIDProvidersClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return TrustedIDProvidersClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return TrustedIDProvidersClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return TrustedIDProvidersClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *TrustedIDProvidersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, trustedIDProviderName string, options *TrustedIDProvidersClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/trustedIdProviders/{trustedIdProviderName}"
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
	if trustedIDProviderName == "" {
		return nil, errors.New("parameter trustedIDProviderName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trustedIdProviderName}", url.PathEscape(trustedIDProviderName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Gets the specified Data Lake Store trusted identity provider.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the Azure resource group.
// accountName - The name of the Data Lake Store account.
// trustedIDProviderName - The name of the trusted identity provider to retrieve.
// options - TrustedIDProvidersClientGetOptions contains the optional parameters for the TrustedIDProvidersClient.Get method.
func (client *TrustedIDProvidersClient) Get(ctx context.Context, resourceGroupName string, accountName string, trustedIDProviderName string, options *TrustedIDProvidersClientGetOptions) (TrustedIDProvidersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, trustedIDProviderName, options)
	if err != nil {
		return TrustedIDProvidersClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return TrustedIDProvidersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return TrustedIDProvidersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *TrustedIDProvidersClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, trustedIDProviderName string, options *TrustedIDProvidersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/trustedIdProviders/{trustedIdProviderName}"
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
	if trustedIDProviderName == "" {
		return nil, errors.New("parameter trustedIDProviderName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trustedIdProviderName}", url.PathEscape(trustedIDProviderName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *TrustedIDProvidersClient) getHandleResponse(resp *http.Response) (TrustedIDProvidersClientGetResponse, error) {
	result := TrustedIDProvidersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TrustedIDProvider); err != nil {
		return TrustedIDProvidersClientGetResponse{}, err
	}
	return result, nil
}

// ListByAccount - Lists the Data Lake Store trusted identity providers within the specified Data Lake Store account.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the Azure resource group.
// accountName - The name of the Data Lake Store account.
// options - TrustedIDProvidersClientListByAccountOptions contains the optional parameters for the TrustedIDProvidersClient.ListByAccount
// method.
func (client *TrustedIDProvidersClient) ListByAccount(resourceGroupName string, accountName string, options *TrustedIDProvidersClientListByAccountOptions) *runtime.Pager[TrustedIDProvidersClientListByAccountResponse] {
	return runtime.NewPager(runtime.PageProcessor[TrustedIDProvidersClientListByAccountResponse]{
		More: func(page TrustedIDProvidersClientListByAccountResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TrustedIDProvidersClientListByAccountResponse) (TrustedIDProvidersClientListByAccountResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByAccountCreateRequest(ctx, resourceGroupName, accountName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return TrustedIDProvidersClientListByAccountResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return TrustedIDProvidersClientListByAccountResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return TrustedIDProvidersClientListByAccountResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByAccountHandleResponse(resp)
		},
	})
}

// listByAccountCreateRequest creates the ListByAccount request.
func (client *TrustedIDProvidersClient) listByAccountCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *TrustedIDProvidersClientListByAccountOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/trustedIdProviders"
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
	reqQP.Set("api-version", "2016-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByAccountHandleResponse handles the ListByAccount response.
func (client *TrustedIDProvidersClient) listByAccountHandleResponse(resp *http.Response) (TrustedIDProvidersClientListByAccountResponse, error) {
	result := TrustedIDProvidersClientListByAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TrustedIDProviderListResult); err != nil {
		return TrustedIDProvidersClientListByAccountResponse{}, err
	}
	return result, nil
}

// Update - Updates the specified trusted identity provider.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the Azure resource group.
// accountName - The name of the Data Lake Store account.
// trustedIDProviderName - The name of the trusted identity provider. This is used for differentiation of providers in the
// account.
// options - TrustedIDProvidersClientUpdateOptions contains the optional parameters for the TrustedIDProvidersClient.Update
// method.
func (client *TrustedIDProvidersClient) Update(ctx context.Context, resourceGroupName string, accountName string, trustedIDProviderName string, options *TrustedIDProvidersClientUpdateOptions) (TrustedIDProvidersClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, accountName, trustedIDProviderName, options)
	if err != nil {
		return TrustedIDProvidersClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return TrustedIDProvidersClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return TrustedIDProvidersClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *TrustedIDProvidersClient) updateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, trustedIDProviderName string, options *TrustedIDProvidersClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/trustedIdProviders/{trustedIdProviderName}"
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
	if trustedIDProviderName == "" {
		return nil, errors.New("parameter trustedIDProviderName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trustedIdProviderName}", url.PathEscape(trustedIDProviderName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	if options != nil && options.Parameters != nil {
		return req, runtime.MarshalAsJSON(req, *options.Parameters)
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *TrustedIDProvidersClient) updateHandleResponse(resp *http.Response) (TrustedIDProvidersClientUpdateResponse, error) {
	result := TrustedIDProvidersClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TrustedIDProvider); err != nil {
		return TrustedIDProvidersClientUpdateResponse{}, err
	}
	return result, nil
}
