//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armrecoveryservicesbackup

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

// BackupResourceStorageConfigsNonCRRClient contains the methods for the BackupResourceStorageConfigsNonCRR group.
// Don't use this type directly, use NewBackupResourceStorageConfigsNonCRRClient() instead.
type BackupResourceStorageConfigsNonCRRClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewBackupResourceStorageConfigsNonCRRClient creates a new instance of BackupResourceStorageConfigsNonCRRClient with the specified values.
func NewBackupResourceStorageConfigsNonCRRClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *BackupResourceStorageConfigsNonCRRClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &BackupResourceStorageConfigsNonCRRClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// Get - Fetches resource storage config.
// If the operation fails it returns the *NewErrorResponse error type.
func (client *BackupResourceStorageConfigsNonCRRClient) Get(ctx context.Context, vaultName string, resourceGroupName string, options *BackupResourceStorageConfigsNonCRRGetOptions) (BackupResourceStorageConfigsNonCRRGetResponse, error) {
	req, err := client.getCreateRequest(ctx, vaultName, resourceGroupName, options)
	if err != nil {
		return BackupResourceStorageConfigsNonCRRGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackupResourceStorageConfigsNonCRRGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return BackupResourceStorageConfigsNonCRRGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *BackupResourceStorageConfigsNonCRRClient) getCreateRequest(ctx context.Context, vaultName string, resourceGroupName string, options *BackupResourceStorageConfigsNonCRRGetOptions) (*policy.Request, error) {
	urlPath := "/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupstorageconfig/vaultstorageconfig"
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
func (client *BackupResourceStorageConfigsNonCRRClient) getHandleResponse(resp *http.Response) (BackupResourceStorageConfigsNonCRRGetResponse, error) {
	result := BackupResourceStorageConfigsNonCRRGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.BackupResourceConfigResource); err != nil {
		return BackupResourceStorageConfigsNonCRRGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *BackupResourceStorageConfigsNonCRRClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := NewErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Patch - Updates vault storage model type.
// If the operation fails it returns the *NewErrorResponse error type.
func (client *BackupResourceStorageConfigsNonCRRClient) Patch(ctx context.Context, vaultName string, resourceGroupName string, parameters BackupResourceConfigResource, options *BackupResourceStorageConfigsNonCRRPatchOptions) (BackupResourceStorageConfigsNonCRRPatchResponse, error) {
	req, err := client.patchCreateRequest(ctx, vaultName, resourceGroupName, parameters, options)
	if err != nil {
		return BackupResourceStorageConfigsNonCRRPatchResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackupResourceStorageConfigsNonCRRPatchResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusNoContent) {
		return BackupResourceStorageConfigsNonCRRPatchResponse{}, client.patchHandleError(resp)
	}
	return BackupResourceStorageConfigsNonCRRPatchResponse{RawResponse: resp}, nil
}

// patchCreateRequest creates the Patch request.
func (client *BackupResourceStorageConfigsNonCRRClient) patchCreateRequest(ctx context.Context, vaultName string, resourceGroupName string, parameters BackupResourceConfigResource, options *BackupResourceStorageConfigsNonCRRPatchOptions) (*policy.Request, error) {
	urlPath := "/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupstorageconfig/vaultstorageconfig"
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
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// patchHandleError handles the Patch error response.
func (client *BackupResourceStorageConfigsNonCRRClient) patchHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := NewErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Update - Updates vault storage model type.
// If the operation fails it returns the *NewErrorResponse error type.
func (client *BackupResourceStorageConfigsNonCRRClient) Update(ctx context.Context, vaultName string, resourceGroupName string, parameters BackupResourceConfigResource, options *BackupResourceStorageConfigsNonCRRUpdateOptions) (BackupResourceStorageConfigsNonCRRUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, vaultName, resourceGroupName, parameters, options)
	if err != nil {
		return BackupResourceStorageConfigsNonCRRUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackupResourceStorageConfigsNonCRRUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return BackupResourceStorageConfigsNonCRRUpdateResponse{}, client.updateHandleError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *BackupResourceStorageConfigsNonCRRClient) updateCreateRequest(ctx context.Context, vaultName string, resourceGroupName string, parameters BackupResourceConfigResource, options *BackupResourceStorageConfigsNonCRRUpdateOptions) (*policy.Request, error) {
	urlPath := "/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupstorageconfig/vaultstorageconfig"
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
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleResponse handles the Update response.
func (client *BackupResourceStorageConfigsNonCRRClient) updateHandleResponse(resp *http.Response) (BackupResourceStorageConfigsNonCRRUpdateResponse, error) {
	result := BackupResourceStorageConfigsNonCRRUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.BackupResourceConfigResource); err != nil {
		return BackupResourceStorageConfigsNonCRRUpdateResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// updateHandleError handles the Update error response.
func (client *BackupResourceStorageConfigsNonCRRClient) updateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := NewErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
