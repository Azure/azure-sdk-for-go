# Azure Monitor Ingestion client module for Go

The Azure Monitor Ingestion client module is used to send custom logs to [Azure Monitor][azure_monitor_overview] using the [Logs Ingestion API][ingestion_overview].

This library allows you to send data from virtually any source to supported built-in tables or to custom tables that you create in Log Analytics workspaces. You can even extend the schema of built-in tables with custom columns.

[Source code][module_source] | Package (pkg.go.dev) | [Product documentation][azure_monitor_overview] | [Samples][ingest_samples]

## Getting started

### Prerequisites

* A supported Go version (the Azure SDK supports the two most recent Go releases)
* An [Azure subscription][azure_subscription]
* An [Azure Log Analytics workspace][azure_monitor_create_using_portal]
* A [Data Collection Endpoint][data_collection_endpoint]
* A [Data Collection Rule][data_collection_rule]

### Install the package

Install the `azingest` and `azidentity` modules with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

The [azidentity][azure_identity] module is used for Azure Active Directory authentication while creating the client.

### Authentication

An authenticated client object is required to upload logs. The examples demonstrate using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate; however, the client accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credential types.

#### Create a client

Example [client][example_client]

## Key concepts

### Data Collection Endpoint

Data Collection Endpoints (DCEs) allow you to uniquely configure ingestion settings for Azure Monitor. [This article][data_collection_endpoint] provides an overview of data collection endpoints including their contents and structure and how you can create and work with them.

### Data Collection Rule

Data collection rules (DCR) define data collected by Azure Monitor and specify how and where that data should be sent or stored. The REST API call must specify a DCR to use. A single DCE can support multiple DCRs, so you can specify a different DCR for different sources and target tables.

The DCR must understand the structure of the input data and the structure of the target table. If the two don't match, it can use a transformation to convert the source data to match the target table. You may also use the transform to filter source data and perform any other calculations or conversions.

For more information, see [Data collection rules in Azure Monitor][data_collection_rule], and see [this article][data_collection_rule_structure] for details about a DCR's structure. For information on how to retrieve a DCR ID, see [this tutorial][data_collection_rule_tutorial].

### Log Analytics workspace tables

Custom logs can send data to any custom table that you create and to certain built-in tables in your Log Analytics workspace. The target table must exist before you can send data to it. The following built-in tables are currently supported:

- [CommonSecurityLog](https://learn.microsoft.com/azure/azure-monitor/reference/tables/commonsecuritylog)
- [SecurityEvents](https://learn.microsoft.com/azure/azure-monitor/reference/tables/securityevent)
- [Syslog](https://learn.microsoft.com/azure/azure-monitor/reference/tables/syslog)
- [WindowsEvents](https://learn.microsoft.com/azure/azure-monitor/reference/tables/windowsevent)

### Logs retrieval

The logs that were uploaded using this module can be queried using the [azquery][azure_monitor_query] module (Azure Monitor Query).

## Examples

Get started with our [examples][ingest_samples].

## Troubleshooting

### Error Handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Monitor Query.

### Logging

This module uses the logging implementation in `azcore`. To turn on logging for all Azure SDK modules, set `AZURE_SDK_GO_LOGGING` to `all`. By default, the logger writes to stderr. Use the `azcore/log` package to control log output. For example, logging only HTTP request and response events, and printing them to stdout:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"

// Print log events to stdout
azlog.SetListener(func(cls azlog.Event, msg string) {
	fmt.Println(msg)
})

// Includes only requests and responses in credential logs
azlog.SetEvents(azlog.EventRequest, azlog.EventResponse)
```

## Next steps

To learn more about Azure Monitor, see the [Azure Monitor service documentation][azure_monitor_overview].

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit [cla.microsoft.com][cla].

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repositories using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][code_of_conduct]. For more information see the [Code of Conduct FAQ][coc_faq] or contact [opencode@microsoft.com][coc_contact] with any additional questions or comments.

<!-- LINKS -->
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_monitor_create_using_portal]: https://learn.microsoft.com/azure/azure-monitor/logs/quick-create-workspace
[azure_monitor_overview]: https://learn.microsoft.com/azure/azure-monitor/
[azure_monitor_query]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery
[azure_subscription]: https://azure.microsoft.com/free/
[data_collection_endpoint]: https://learn.microsoft.com/azure/azure-monitor/essentials/data-collection-endpoint-overview
[data_collection_rule]: https://learn.microsoft.com/azure/azure-monitor/essentials/data-collection-rule-overview
[data_collection_rule_structure]: https://learn.microsoft.com/azure/azure-monitor/essentials/data-collection-rule-structure
[data_collection_rule_tutorial]: https://learn.microsoft.com/azure/azure-monitor/logs/tutorial-logs-ingestion-portal#collect-information-from-the-dcr
[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[example_client]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest#example-NewClient
[ingest_samples]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest#pkg-examples
[ingestion_overview]: https://learn.microsoft.com/azure/azure-monitor/logs/logs-ingestion-api-overview
[module_source]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/monitor/azingest

[cla]: https://cla.microsoft.com
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc_contact]: mailto:opencode@microsoft.com