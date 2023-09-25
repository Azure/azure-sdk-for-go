//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armrecoveryservicesbackup

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

// BackupResourceEncryptionConfigsClient contains the methods for the BackupResourceEncryptionConfigs group.
// Don't use this type directly, use NewBackupResourceEncryptionConfigsClient() instead.
type BackupResourceEncryptionConfigsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewBackupResourceEncryptionConfigsClient creates a new instance of BackupResourceEncryptionConfigsClient with the specified values.
//   - subscriptionID - The subscription Id.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewBackupResourceEncryptionConfigsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*BackupResourceEncryptionConfigsClient, error) {
	cl, err := arm.NewClient(moduleName+".BackupResourceEncryptionConfigsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &BackupResourceEncryptionConfigsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Get - Fetches Vault Encryption config.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
//   - vaultName - The name of the recovery services vault.
//   - resourceGroupName - The name of the resource group where the recovery services vault is present.
//   - options - BackupResourceEncryptionConfigsClientGetOptions contains the optional parameters for the BackupResourceEncryptionConfigsClient.Get
//     method.
func (client *BackupResourceEncryptionConfigsClient) Get(ctx context.Context, vaultName string, resourceGroupName string, options *BackupResourceEncryptionConfigsClientGetOptions) (BackupResourceEncryptionConfigsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, vaultName, resourceGroupName, options)
	if err != nil {
		return BackupResourceEncryptionConfigsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return BackupResourceEncryptionConfigsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return BackupResourceEncryptionConfigsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *BackupResourceEncryptionConfigsClient) getCreateRequest(ctx context.Context, vaultName string, resourceGroupName string, options *BackupResourceEncryptionConfigsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupEncryptionConfigs/backupResourceEncryptionConfig"
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *BackupResourceEncryptionConfigsClient) getHandleResponse(resp *http.Response) (BackupResourceEncryptionConfigsClientGetResponse, error) {
	result := BackupResourceEncryptionConfigsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.BackupResourceEncryptionConfigExtendedResource); err != nil {
		return BackupResourceEncryptionConfigsClientGetResponse{}, err
	}
	return result, nil
}

// Update - Updates Vault encryption config.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
//   - vaultName - The name of the recovery services vault.
//   - resourceGroupName - The name of the resource group where the recovery services vault is present.
//   - parameters - Vault encryption input config request
//   - options - BackupResourceEncryptionConfigsClientUpdateOptions contains the optional parameters for the BackupResourceEncryptionConfigsClient.Update
//     method.
func (client *BackupResourceEncryptionConfigsClient) Update(ctx context.Context, vaultName string, resourceGroupName string, parameters BackupResourceEncryptionConfigResource, options *BackupResourceEncryptionConfigsClientUpdateOptions) (BackupResourceEncryptionConfigsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, vaultName, resourceGroupName, parameters, options)
	if err != nil {
		return BackupResourceEncryptionConfigsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return BackupResourceEncryptionConfigsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return BackupResourceEncryptionConfigsClientUpdateResponse{}, err
	}
	return BackupResourceEncryptionConfigsClientUpdateResponse{}, nil
}

// updateCreateRequest creates the Update request.
func (client *BackupResourceEncryptionConfigsClient) updateCreateRequest(ctx context.Context, vaultName string, resourceGroupName string, parameters BackupResourceEncryptionConfigResource, options *BackupResourceEncryptionConfigsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupEncryptionConfigs/backupResourceEncryptionConfig"
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

