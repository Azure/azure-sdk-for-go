# Azure Monitor Metrics client module for Go

Source code | Package (pkg.go.dev) | [REST API documentation][monitor_rest_docs] | [Product documentation][monitor_docs] | Samples

## Getting started

### Prerequisites

* Go version 1.18 or higher - [Install Go](https://go.dev/doc/install)
* Azure subscription - [Create a free account][azure_sub]
* The resource URI of an Azure resource (Storage Account, Key Vault, CosmosDB, etc.) that you plan to monitor

### Install the packages

Install the `azmetrics` and `azidentity` modules with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

The [azidentity][azure_identity] module is used for client authentication.

### Authentication

An authenticated client object is required to execute a query. The examples demonstrate using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate; however, the client accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credential types.

The clients default to the Azure public cloud. For other cloud configurations, see the [cloud][cloud_documentation] package documentation.

#### Create a client

Example client

## Key concepts

### Metrics data structure

Each set of metric values is a time series with the following characteristics:

- The time the value was collected
- The resource associated with the value
- A namespace that acts like a category for the metric
- A metric name
- The value itself
- Some metrics may have multiple dimensions as described in [multi-dimensional metrics][multi-metrics]. Custom metrics can have up to 10 dimensions.

## Examples

Get started with our examples.

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a [Contributor License Agreement (CLA)][cla] declaring that you have the right to, and actually do, grant us the rights to use your contribution.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate
the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to
do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][coc]. For more information, see
the [Code of Conduct FAQ][coc_faq] or contact [opencode@microsoft.com][coc_contact] with any additional questions or
comments.

<!-- LINKS -->
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[cloud_documentation]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud
[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[monitor_docs]: https://learn.microsoft.com/azure/azure-monitor/
[monitor_rest_docs]: https://learn.microsoft.com/rest/api/monitor/
[multi-metrics]: https://learn.microsoft.com/azure/azure-monitor/essentials/data-platform-metrics#multi-dimensional-metrics
[cla]: https://cla.microsoft.com
[coc]: https://opensource.microsoft.com/codeofconduct/
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc_contact]: mailto:opencode@microsoft.com
