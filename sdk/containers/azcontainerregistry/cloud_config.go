//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcontainerregistry

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

const (
	// ServiceName is the cloud service name for Azure Container Registry
	ServiceName cloud.ServiceName = "azcontainerregistry"
)

func init() {
	cloud.AzureChina.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://management.chinacloudapi.cn/",
	}
	cloud.AzureGovernment.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://management.usgovcloudapi.net/",
	}
	cloud.AzurePublic.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://management.azure.com",
	}
}
