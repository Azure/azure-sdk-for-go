//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armrecoveryservices

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// VaultExtendedInfoClient contains the methods for the VaultExtendedInfo group.
// Don't use this type directly, use NewVaultExtendedInfoClient() instead.
type VaultExtendedInfoClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewVaultExtendedInfoClient creates a new instance of VaultExtendedInfoClient with the specified values.
func NewVaultExtendedInfoClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *VaultExtendedInfoClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &VaultExtendedInfoClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// CreateOrUpdate - Create vault extended info.
// If the operation fails it returns the *CloudError error type.
func (client *VaultExtendedInfoClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, resourceExtendedInfoDetails VaultExtendedInfoResource, options *VaultExtendedInfoCreateOrUpdateOptions) (VaultExtendedInfoCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, vaultName, resourceExtendedInfoDetails, options)
	if err != nil {
		return VaultExtendedInfoCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VaultExtendedInfoCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VaultExtendedInfoCreateOrUpdateResponse{}, client.createOrUpdateHandleError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *VaultExtendedInfoClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, resourceExtendedInfoDetails VaultExtendedInfoResource, options *VaultExtendedInfoCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/extendedInformation/vaultExtendedInfo"
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
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, resourceExtendedInfoDetails)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *VaultExtendedInfoClient) createOrUpdateHandleResponse(resp *http.Response) (VaultExtendedInfoCreateOrUpdateResponse, error) {
	result := VaultExtendedInfoCreateOrUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.VaultExtendedInfoResource); err != nil {
		return VaultExtendedInfoCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *VaultExtendedInfoClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Get the vault extended info.
// If the operation fails it returns the *CloudError error type.
func (client *VaultExtendedInfoClient) Get(ctx context.Context, resourceGroupName string, vaultName string, options *VaultExtendedInfoGetOptions) (VaultExtendedInfoGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, vaultName, options)
	if err != nil {
		return VaultExtendedInfoGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VaultExtendedInfoGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VaultExtendedInfoGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VaultExtendedInfoClient) getCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, options *VaultExtendedInfoGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/extendedInformation/vaultExtendedInfo"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VaultExtendedInfoClient) getHandleResponse(resp *http.Response) (VaultExtendedInfoGetResponse, error) {
	result := VaultExtendedInfoGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.VaultExtendedInfoResource); err != nil {
		return VaultExtendedInfoGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *VaultExtendedInfoClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Update - Update vault extended info.
// If the operation fails it returns the *CloudError error type.
func (client *VaultExtendedInfoClient) Update(ctx context.Context, resourceGroupName string, vaultName string, resourceExtendedInfoDetails VaultExtendedInfoResource, options *VaultExtendedInfoUpdateOptions) (VaultExtendedInfoUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, vaultName, resourceExtendedInfoDetails, options)
	if err != nil {
		return VaultExtendedInfoUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VaultExtendedInfoUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VaultExtendedInfoUpdateResponse{}, client.updateHandleError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *VaultExtendedInfoClient) updateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, resourceExtendedInfoDetails VaultExtendedInfoResource, options *VaultExtendedInfoUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/extendedInformation/vaultExtendedInfo"
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
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, resourceExtendedInfoDetails)
}

// updateHandleResponse handles the Update response.
func (client *VaultExtendedInfoClient) updateHandleResponse(resp *http.Response) (VaultExtendedInfoUpdateResponse, error) {
	result := VaultExtendedInfoUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.VaultExtendedInfoResource); err != nil {
		return VaultExtendedInfoUpdateResponse{}, err
	}
	return result, nil
}

// updateHandleError handles the Update error response.
func (client *VaultExtendedInfoClient) updateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
