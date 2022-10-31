//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

const monitorQueryLogs cloud.ServiceName = "azqueryLogs"
const monitorQueryMetrics cloud.ServiceName = "azqueryMetrics"

func init() {
	cloud.AzureChina.Services[monitorQueryLogs] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://api.loganalytics.azure.cn",
	}
	cloud.AzureGovernment.Services[monitorQueryLogs] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://api.loganalytics.us",
	}
	cloud.AzurePublic.Services[monitorQueryLogs] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://api.loganalytics.io",
	}
	cloud.AzureChina.Services[monitorQueryMetrics] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://management.chinacloudapi.cn/",
	}
	cloud.AzureGovernment.Services[monitorQueryMetrics] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://management.usgovcloudapi.net/",
	}
	cloud.AzurePublic.Services[monitorQueryMetrics] = cloud.ServiceConfiguration{
		Audience: "",
		Endpoint: "https://management.azure.com",
	}
}
