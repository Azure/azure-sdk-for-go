# Azure Monitor Query client module for Go

The Azure Monitor Query client library is used to execute read-only queries against [Azure Monitor][azure_monitor_overview]'s two data platforms:

- [Logs](https://docs.microsoft.com/azure/azure-monitor/logs/data-platform-logs) - Collects and organizes log and performance data from monitored resources. Data from different sources such as platform logs from Azure services, log and performance data from virtual machines agents, and usage and performance data from apps can be consolidated into a single [Azure Log Analytics workspace](https://docs.microsoft.com/azure/azure-monitor/logs/data-platform-logs#log-analytics-and-workspaces). The various data types can be analyzed together using the [Kusto Query Language][kusto_query_language].
- [Metrics](https://docs.microsoft.com/azure/azure-monitor/essentials/data-platform-metrics) - Collects numeric data from monitored resources into a time series database. Metrics are numerical values that are collected at regular intervals and describe some aspect of a system at a particular time. Metrics are lightweight and capable of supporting near real-time scenarios, making them particularly useful for alerting and fast detection of issues.

## Getting started

### Install packages

Install `azquery` and `azidentity` with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.

### Prerequisites

* An [Azure subscription][azure_sub]
* A supported Go version (the Azure SDK supports the two most recent Go releases)
* For log queries, a Log Analytics workspace. 
* For metric queries, a Resource URI.

### Authentication

This document demonstrates using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate. This credential type works in both local development and production environments. We recommend using a [managed identity][managed_identity] in production.

[Client][client_docs] accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credential types.

#### Create a logs client

```go
import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewLogsClient(cred, nil)
}
```

#### Create a metrics client

```go
import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client := azkeys.NewMetricsClient(cred, nil)
}
```

### Execute the query

For examples of Logs and Metrics queries, see the [Examples](#examples) section.

## Key concepts

## Examples

- [Logs query](#logs-query)
    - [Query workspace]
    - [Batch query]
- [Metrics query](#metrics-query)
  - [Query resource]