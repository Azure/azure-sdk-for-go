//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

const MonitorQuery cloud.ServiceName = "azquery"

func init() {
	cloud.AzureChina.Services[MonitorQuery] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://api.loganalytics.azure.cn",
	}
	cloud.AzureGovernment.Services[MonitorQuery] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://api.loganalytics.us",
	}
	cloud.AzurePublic.Services[MonitorQuery] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://api.loganalytics.io",
	}
}
