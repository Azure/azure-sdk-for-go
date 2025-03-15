//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

const (
	moduleName    = "github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	moduleVersion = "v0.2.3"
)

const (
	// AcrAudience is the audience for Azure Container Registry.
	// acr audience token provides permissions with a smaller scope than ARM audience token,
	// and is preferable as the default audience for Azure Container Registry in all clouds.
	AcrAudience = "https://containerregistry.azure.net"
)
