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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// VaultCertificatesClient contains the methods for the VaultCertificates group.
// Don't use this type directly, use NewVaultCertificatesClient() instead.
type VaultCertificatesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewVaultCertificatesClient creates a new instance of VaultCertificatesClient with the specified values.
// subscriptionID - The subscription Id.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewVaultCertificatesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *VaultCertificatesClient {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := options.Endpoint
	if len(ep) == 0 {
		ep = arm.AzurePublicCloud
	}
	client := &VaultCertificatesClient{
		subscriptionID: subscriptionID,
		host:           string(ep),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options),
	}
	return client
}

// Create - Uploads a certificate for a resource.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// vaultName - The name of the recovery services vault.
// certificateName - Certificate friendly name.
// certificateRequest - Input parameters for uploading the vault certificate.
// options - VaultCertificatesClientCreateOptions contains the optional parameters for the VaultCertificatesClient.Create
// method.
func (client *VaultCertificatesClient) Create(ctx context.Context, resourceGroupName string, vaultName string, certificateName string, certificateRequest CertificateRequest, options *VaultCertificatesClientCreateOptions) (VaultCertificatesClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, vaultName, certificateName, certificateRequest, options)
	if err != nil {
		return VaultCertificatesClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VaultCertificatesClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VaultCertificatesClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *VaultCertificatesClient) createCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, certificateName string, certificateRequest CertificateRequest, options *VaultCertificatesClientCreateOptions) (*policy.Request, error) {
	urlPath := "/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/certificates/{certificateName}"
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
	if certificateName == "" {
		return nil, errors.New("parameter certificateName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{certificateName}", url.PathEscape(certificateName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, certificateRequest)
}

// createHandleResponse handles the Create response.
func (client *VaultCertificatesClient) createHandleResponse(resp *http.Response) (VaultCertificatesClientCreateResponse, error) {
	result := VaultCertificatesClientCreateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.VaultCertificateResponse); err != nil {
		return VaultCertificatesClientCreateResponse{}, err
	}
	return result, nil
}
