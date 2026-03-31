// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

// ServiceName is the [cloud.ServiceName] for Azure App Configuration, used to identify the respective [cloud.ServiceConfiguration].
//
// NOTE: ServiceConfiguration omits the Endpoint as that's explicitly passed to client constructors.
const ServiceName cloud.ServiceName = "data/azappconfig"

func init() {
	cloud.AzureChina.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://appconfig.azure.cn",
	}
	cloud.AzureGovernment.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://appconfig.azure.us",
	}
	cloud.AzurePublic.Services[ServiceName] = cloud.ServiceConfiguration{
		Audience: "https://appconfig.azure.com",
	}
}
