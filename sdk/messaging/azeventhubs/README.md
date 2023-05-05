# Azure Event Hubs Client Module for Go

[Azure Event Hubs][eventhubs_docs] is a big data streaming platform and event ingestion service from Microsoft. For more information about Event Hubs, see: [link][eventhubs_about].

Use the client module `github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs` in your application to:

- Send events to an event hub.
- Consume events from an event hub.

Key links:

[Source code][azeventhubs_repo] | [Package (pkg.go.dev)][azeventhubs_pkg_go_docs] | [REST API documentation][eventhubs_rest_docs] | [Product documentation][eventhubs_docs] | [Samples][azeventhubs_samples]

## Getting started

### Prerequisites

- Go, version 1.18 or higher - [Install Go][go_install]
- Azure subscription - [Create a free account][azure_sub]
- An [Event Hub namespace][eventhubs_namespace]
- An Event Hub. You can create an event hub in your Event Hubs Namespace using the [Azure Portal][eventhubs_create_portal], or the [Azure CLI][eventhubs_create_cli].

### Install the package

Install the Azure Event Hubs client module for Go with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs
```

### Authenticate the client

Event Hub clients are created using an Event Hub a credential from the [Azure Identity package][azure_identity_pkg], like [DefaultAzureCredential][default_azure_credential].

Alternatively, you can create a client using a connection string.

#### Using a service principal
 - ConsumerClient: [link][azeventhubs_serviceprinciple_consumerclient]
 - ProducerClient: [link][azeventhubs_serviceprinciple_producerclient]

#### Using a connection string
 - ConsumerClient: [link][azeventhubs_connectionstring_consumerclient]
 - ProducerClient: [link][azeventhubs_connectionstring_producerclient]

# Key concepts

An Event Hub [**namespace**][eventhubs_namespace] can have multiple event hubs. Each event hub, in turn, contains [**partitions**][eventhubs_partition] which store events.

Events are published to an event hub using an [event publisher][eventhubs_eventpublisher]. In this package, the event publisher is the [ProducerClient][azeventhubs_producerclient].

Events can be consumed from an event hub using an [event consumer][eventhubs_eventconsumer]. In this package there are two types for consuming events: 
- The basic event consumer is the, in the [ConsumerClient][azeventhubs_consumerclient]. This consumer is useful if you already know which partitions you want to receive from.
- A distributed event consumer, which uses Azure Blobs for checkpointing and coordination. This is implemented in the [Processor][azeventhubs_processor]. This is useful when you want to have the partition assignment be dynamically chosen, and balanced with other Processor instances.

For more information about Event Hubs features and terminology can be found here: [link][eventhubs_features_terminology]

# Examples

Examples for various scenarios can be found on [pkg.go.dev][azeventhubs_samples] or in the example*_test.go files in our GitHub repo for [azeventhubs][azeventhubs_repo_main].

# Troubleshooting

### Logging

This module uses the classification-based logging implementation in `azcore`. To enable console logging for all SDK modules, set the environment variable `AZURE_SDK_GO_LOGGING` to `all`. 

Use the `azcore/log` package to control log event output or to enable logs for `azeventhubs` only. For example:

```go
import (
  "fmt"
  azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
)

// print log output to stdout
azlog.SetListener(func(event azlog.Event, s string) {
    fmt.Printf("[%s] %s\n", event, s)
})

// pick the set of events to log
azlog.SetEvents(
  azeventhubs.EventConn,
  azeventhubs.EventAuth,
  azeventhubs.EventProducer,
  azeventhubs.EventConsumer,
)
```

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a [Contributor License Agreement (CLA)][cla] declaring that you have the right to, and actually do, grant us the rights to use your contribution.
 
If you'd like to contribute to this library, please read the [contributing guide] [contributing_guide] to learn more about how to build and test the code.
 
When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.
 
This project has adopted the [Microsoft Open Source Code of Conduct][coc]. For more information, see the [Code of Conduct FAQ][coc_faq] or contact [opencode@microsoft.com][coc_contact] with any additional questions or comments.


### Additional Helpful Links for Contributors  

Many people all over the world have helped make this project better.  You'll want to check out:

* [What are some good first issues for new contributors to the repo?][azure_sdk_for_go_contributing_first_issues]
* [How to build and test your change][azure_sdk_for_go_contributing_developer_guide]
* [How you can make a change happen!][azure_sdk_for_go_contributing_pull_requests]
* Frequently Asked Questions (FAQ) and Conceptual Topics in the detailed [Azure SDK for Go wiki][azure_sdk_for_go_wiki].

<!-- ### Community-->
### Reporting security issues and security bugs

Security issues and bugs should be reported privately, via email, to the Microsoft Security Response Center (MSRC) <secure@microsoft.com>. You should receive a response within 24 hours. If for some reason you do not, please follow up via email to ensure we received your original message. Further information, including the MSRC PGP key, can be found in the [Security TechCenter][security_techcenter].

### License

Azure SDK for Go is licensed under the [MIT][azeventhubs_mit_license] license.

<!-- LINKS -->

[azure_identity_pkg]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sdk_for_go_contributing]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md
[azure_sdk_for_go_contributing_developer_guide]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#developer-guide
[azure_sdk_for_go_contributing_first_issues]: https://github.com/azure/azure-sdk-for-go/issues?q=is%3Aopen+is%3Aissue+label%3A%22up+for+grabs%22
[azure_sdk_for_go_contributing_pull_requests]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#pull-requests
[azure_sdk_for_go_wiki]: https://github.com/azure/azure-sdk-for-go/wiki
[azure_sub]: https://azure.microsoft.com/free/
[azeventhubs_connectionstring_consumerclient]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#example-NewConsumerClientFromConnectionString
[azeventhubs_connectionstring_producerclient]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#example-NewProducerClientFromConnectionString
[azeventhubs_consumerclient]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#ConsumerClient
[azeventhubs_mit_license]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/LICENSE.txt
[azeventhubs_repo]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/messaging/azeventhubs
[azeventhubs_repo_main]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs
[azeventhubs_samples]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#pkg-examples
[azeventhubs_pkg_go_docs]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs
[azeventhubs_processor]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#Processor
[azeventhubs_producerclient]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#ProducerClient
[azeventhubs_serviceprinciple_consumerclient]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#example-NewConsumerClient
[azeventhubs_serviceprinciple_producerclient]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#example-NewProducerClient

[contributing_guide]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md
[cla]: https://cla.microsoft.com
[coc]: https://opensource.microsoft.com/codeofconduct/
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc_contact]: mailto:opencode@microsoft.com

[default_azure_credential]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential
[eventhubs_about]: https://docs.microsoft.com/azure/event-hubs/event-hubs-about
[eventhubs_eventconsumer]: https://docs.microsoft.com/azure/event-hubs/event-hubs-features#event-consumers
[eventhubs_create_portal]: https://docs.microsoft.com/azure/event-hubs/event-hubs-create
[eventhubs_create_cli]: https://docs.microsoft.com/azure/event-hubs/event-hubs-quickstart-cli
[eventhubs_docs]: https://azure.microsoft.com/services/event-hubs/
[eventhubs_eventpublisher]: https://docs.microsoft.com/azure/event-hubs/event-hubs-features#event-publishers
[eventhubs_features_terminology]: https://docs.microsoft.com/azure/event-hubs/event-hubs-features
[eventhubs_namespace]: https://docs.microsoft.com/azure/event-hubs/event-hubs-features#namespace
[eventhubs_partition]: https://docs.microsoft.com/azure/event-hubs/event-hubs-features#partitions
[eventhubs_rest_docs]: https://learn.microsoft.com/en-us/rest/api/eventhub/
[go_install]: https://go.dev/doc/install
[security_techcenter]: https://www.microsoft.com/msrc/faqs-report-an-issue

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fmessaging%2Fazeventhubs%2FREADME.png)
