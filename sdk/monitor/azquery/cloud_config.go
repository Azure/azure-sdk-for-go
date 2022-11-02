//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

const (
	MonitorQueryLogs    cloud.ServiceName = "azqueryLogs"
	MonitorQueryMetrics cloud.ServiceName = "azqueryMetrics"
)

func init() {
	cloud.AzureChina.Services[MonitorQueryLogs] = cloud.ServiceConfiguration{
		Audience: "https://api.loganalytics.azure.cn",
		Endpoint: "https://api.loganalytics.azure.cn/v1",
	}
	cloud.AzureGovernment.Services[MonitorQueryLogs] = cloud.ServiceConfiguration{
		Audience: "https://api.loganalytics.us",
		Endpoint: "https://api.loganalytics.us/v1",
	}
	cloud.AzurePublic.Services[MonitorQueryLogs] = cloud.ServiceConfiguration{
		Audience: "https://api.loganalytics.io",
		Endpoint: "https://api.loganalytics.io/v1",
	}
	cloud.AzureChina.Services[MonitorQueryMetrics] = cloud.ServiceConfiguration{
		Audience: "https://management.chinacloudapi.cn/",
		Endpoint: "https://management.chinacloudapi.cn/",
	}
	cloud.AzureGovernment.Services[MonitorQueryMetrics] = cloud.ServiceConfiguration{
		Audience: "https://management.usgovcloudapi.net/",
		Endpoint: "https://management.usgovcloudapi.net/",
	}
	cloud.AzurePublic.Services[MonitorQueryMetrics] = cloud.ServiceConfiguration{
		Audience: "https://management.azure.com",
		Endpoint: "https://management.azure.com",
	}
}
