//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azingest

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

// Cloud Service Names for Monitor Ingestion, used to identify the respective cloud.ServiceConfiguration
const (
	ServiceNameIngestion cloud.ServiceName = "azingest"
)

func init() {
	cloud.AzureChina.Services[ServiceNameIngestion] = cloud.ServiceConfiguration{
		Audience: "https://monitor.azure.cn/",
	}
	cloud.AzureGovernment.Services[ServiceNameIngestion] = cloud.ServiceConfiguration{
		Audience: "https://monitor.azure.us/",
	}
	cloud.AzurePublic.Services[ServiceNameIngestion] = cloud.ServiceConfiguration{
		Audience: "https://monitor.azure.com/",
	}
}
