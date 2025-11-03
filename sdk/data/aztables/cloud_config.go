// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

// ServiceName is the [cloud.ServiceName] for Azure Tables, used to identify the respective [cloud.ServiceConfiguration].
//
// NOTE: ServiceConfiguration omits the Endpoint as that's explicitly passed to client constructors.
const ServiceName cloud.ServiceName = "data/aztables"

func init() {
	cloud.AzureChina.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://storage.azure.com",
		CosmosAudience: "https://cosmos.azure.cn",
	}
	cloud.AzureGovernment.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://storage.azure.com",
		CosmosAudience: "https://cosmos.azure.us",
	}
	cloud.AzurePublic.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://storage.azure.com",
		CosmosAudience: "https://cosmos.azure.com",
	}
}
