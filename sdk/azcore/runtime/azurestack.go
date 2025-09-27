//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package runtime

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	// AzureStackStorageAPIVersionEnvVar is the environment variable name for configuring the API version
	AzureStackStorageAPIVersionEnvVar = "AZURESTACK_STORAGE_API_VERSION"

	// DefaultAzureStackAPIVersion is the default API version to use for AzureStack
	DefaultAzureStackAPIVersion = "2019-07-07"
)

// GetAzureStackAPIVersion returns the API version to use for AzureStack if configured
// via AZURESTACK_STORAGE_API_VERSION environment variable.
func GetAzureStackAPIVersion() string {
	version := os.Getenv(AzureStackStorageAPIVersionEnvVar)
	if version == "" {
		return ""
	}
	// If environment variable is set but empty, use the default version
	if version == "*" {
		return DefaultAzureStackAPIVersion
	}
	return version
}

// ApplyAzureStackAPIVersion applies the AzureStack API version to ClientOptions if configured
func ApplyAzureStackAPIVersion(options *policy.ClientOptions) {
	if options == nil {
		return
	}

	// Only set API version if AzureStack environment variable is configured
	if apiVersion := GetAzureStackAPIVersion(); apiVersion != "" {
		options.APIVersion = apiVersion
	}
}
