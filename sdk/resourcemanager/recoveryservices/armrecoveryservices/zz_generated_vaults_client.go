//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armrecoveryservices

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

// VaultsClient contains the methods for the Vaults group.
// Don't use this type directly, use NewVaultsClient() instead.
type VaultsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewVaultsClient creates a new instance of VaultsClient with the specified values.
// subscriptionID - The subscription Id.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewVaultsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VaultsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &VaultsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a Recovery Services vault.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// vaultName - The name of the recovery services vault.
// vault - Recovery Services Vault to be created.
// options - VaultsClientBeginCreateOrUpdateOptions contains the optional parameters for the VaultsClient.BeginCreateOrUpdate
// method.
func (client *VaultsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, vault Vault, options *VaultsClientBeginCreateOrUpdateOptions) (*runtime.Poller[VaultsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, vaultName, vault, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[VaultsClientCreateOrUpdateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[VaultsClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - Creates or updates a Recovery Services vault.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
func (client *VaultsClient) createOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, vault Vault, options *VaultsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, vaultName, vault, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *VaultsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, vault Vault, options *VaultsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, vault)
}

// Delete - Deletes a vault.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// vaultName - The name of the recovery services vault.
// options - VaultsClientDeleteOptions contains the optional parameters for the VaultsClient.Delete method.
func (client *VaultsClient) Delete(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsClientDeleteOptions) (VaultsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, vaultName, options)
	if err != nil {
		return VaultsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VaultsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VaultsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return VaultsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *VaultsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get the Vault details.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// vaultName - The name of the recovery services vault.
// options - VaultsClientGetOptions contains the optional parameters for the VaultsClient.Get method.
func (client *VaultsClient) Get(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsClientGetOptions) (VaultsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, vaultName, options)
	if err != nil {
		return VaultsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VaultsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VaultsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VaultsClient) getCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, options *VaultsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VaultsClient) getHandleResponse(resp *http.Response) (VaultsClientGetResponse, error) {
	result := VaultsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Vault); err != nil {
		return VaultsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Retrieve a list of Vaults.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// options - VaultsClientListByResourceGroupOptions contains the optional parameters for the VaultsClient.ListByResourceGroup
// method.
func (client *VaultsClient) NewListByResourceGroupPager(resourceGroupName string, options *VaultsClientListByResourceGroupOptions) *runtime.Pager[VaultsClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[VaultsClientListByResourceGroupResponse]{
		More: func(page VaultsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VaultsClientListByResourceGroupResponse) (VaultsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VaultsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return VaultsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VaultsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *VaultsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *VaultsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *VaultsClient) listByResourceGroupHandleResponse(resp *http.Response) (VaultsClientListByResourceGroupResponse, error) {
	result := VaultsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VaultList); err != nil {
		return VaultsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionIDPager - Fetches all the resources of the specified type in the subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// options - VaultsClientListBySubscriptionIDOptions contains the optional parameters for the VaultsClient.ListBySubscriptionID
// method.
func (client *VaultsClient) NewListBySubscriptionIDPager(options *VaultsClientListBySubscriptionIDOptions) *runtime.Pager[VaultsClientListBySubscriptionIDResponse] {
	return runtime.NewPager(runtime.PagingHandler[VaultsClientListBySubscriptionIDResponse]{
		More: func(page VaultsClientListBySubscriptionIDResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VaultsClientListBySubscriptionIDResponse) (VaultsClientListBySubscriptionIDResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionIDCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VaultsClientListBySubscriptionIDResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return VaultsClientListBySubscriptionIDResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VaultsClientListBySubscriptionIDResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionIDHandleResponse(resp)
		},
	})
}

// listBySubscriptionIDCreateRequest creates the ListBySubscriptionID request.
func (client *VaultsClient) listBySubscriptionIDCreateRequest(ctx context.Context, options *VaultsClientListBySubscriptionIDOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.RecoveryServices/vaults"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionIDHandleResponse handles the ListBySubscriptionID response.
func (client *VaultsClient) listBySubscriptionIDHandleResponse(resp *http.Response) (VaultsClientListBySubscriptionIDResponse, error) {
	result := VaultsClientListBySubscriptionIDResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VaultList); err != nil {
		return VaultsClientListBySubscriptionIDResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Updates the vault.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// vaultName - The name of the recovery services vault.
// vault - Recovery Services Vault to be created.
// options - VaultsClientBeginUpdateOptions contains the optional parameters for the VaultsClient.BeginUpdate method.
func (client *VaultsClient) BeginUpdate(ctx context.Context, resourceGroupName string, vaultName string, vault PatchVault, options *VaultsClientBeginUpdateOptions) (*runtime.Poller[VaultsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, vaultName, vault, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[VaultsClientUpdateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[VaultsClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - Updates the vault.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
func (client *VaultsClient) update(ctx context.Context, resourceGroupName string, vaultName string, vault PatchVault, options *VaultsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, vaultName, vault, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *VaultsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, vault PatchVault, options *VaultsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, vault)
}
