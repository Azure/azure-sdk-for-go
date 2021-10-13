//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"

// Endpoint is the base URL for Azure Resource Manager.
type Endpoint string

const (
	// AzureChina is the Azure Resource Manager China cloud endpoint.
	AzureChina Endpoint = "https://management.chinacloudapi.cn/"
	// AzureGermany is the Azure Resource Manager Germany cloud endpoint.
	AzureGermany Endpoint = "https://management.microsoftazure.de/"
	// AzureGovernment is the Azure Resource Manager US government cloud endpoint.
	AzureGovernment Endpoint = "https://management.usgovcloudapi.net/"
	// AzurePublicCloud is the Azure Resource Manager public cloud endpoint.
	AzurePublicCloud Endpoint = "https://management.azure.com/"
)

// ClientOptions contains configuration settings for a client's pipeline.
type ClientOptions struct {
	policy.ClientOptions

	// AuxiliaryTenants contains a list of additional tenants for cross-tenant requests.
	AuxiliaryTenants []string

	// DisableRPRegistration disables the auto-RP registration policy. Defaults to false.
	DisableRPRegistration bool

	// Host is the base URL for Azure Resource Manager. Defaults to AzurePublicCloud.
	Host Endpoint
}
