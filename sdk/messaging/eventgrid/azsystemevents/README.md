# Azure Event Grid System Events Module for Go

Azure Event Grid system events are published by Azure services to system topics. The models in this package map to events sent by various Azure services.

Key links:
- [Source code][source]
- [API Reference Documentation][godoc]
- [Product documentation][product_docs]
- [Samples][godoc_examples]

## Getting started

### Install the package

Install the Azure Event Grid system events module for Go with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents
```

### Prerequisites

- Go, version 1.18 or higher

# Key concepts

Event subscriptions can be used to forward events from an [Event Grid system topic][system_topics] to a data source, like [Azure Storage Queues][event_handler_storage_queues]. The payload will be formatted as an array of events, using the event envelope (Cloud Event Schema or Event Grid Schema) configured within the subscription.

To consume events, use the client package for that service. For example, if the Event Grid subscription uses an an Azure Storage Queue, we would use the [azqeueue](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue) package to consume it.

# Examples

Examples for deserializing system events can be found on [pkg.go.dev][godoc_examples] or in the example*_test.go files in our GitHub repo for [azsystemevents][source].

# Troubleshooting

### Logging

This module uses the classification-based logging implementation in `azcore`. To enable console logging for all SDK modules, set the environment variable `AZURE_SDK_GO_LOGGING` to `all`. 

Use the `azcore/log` package to control log event output.

```go
import (
  "fmt"
  azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
)

// print log output to stdout
azlog.SetListener(func(event azlog.Event, s string) {
    fmt.Printf("[%s] %s\n", event, s)
})
```

# Next steps

More sample code should go here, along with links out to the appropriate example tests.

## Contributing
For details on contributing to this repository, see the [contributing guide][azure_sdk_for_go_contributing].

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

### Additional Helpful Links for Contributors  
Many people all over the world have helped make this project better.  You'll want to check out:

* [What are some good first issues for new contributors to the repo?](https://github.com/azure/azure-sdk-for-go/issues?q=is%3Aopen+is%3Aissue+label%3A%22up+for+grabs%22)
* [How to build and test your change][azure_sdk_for_go_contributing_developer_guide]
* [How you can make a change happen!][azure_sdk_for_go_contributing_pull_requests]
* Frequently Asked Questions (FAQ) and Conceptual Topics in the detailed [Azure SDK for Go wiki](https://github.com/azure/azure-sdk-for-go/wiki).

<!-- ### Community-->
### Reporting security issues and security bugs

Security issues and bugs should be reported privately, via email, to the Microsoft Security Response Center (MSRC) <secure@microsoft.com>. You should receive a response within 24 hours. If for some reason you do not, please follow up via email to ensure we received your original message. Further information, including the MSRC PGP key, can be found in the [Security TechCenter](https://www.microsoft.com/msrc/faqs-report-an-issue).

### License

Azure SDK for Go is licensed under the [MIT](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/template/aztemplate/LICENSE.txt) license.

<!-- LINKS -->
[azure_sdk_for_go_contributing]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md
[azure_sdk_for_go_contributing_developer_guide]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#developer-guide
[azure_sdk_for_go_contributing_pull_requests]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#pull-requests
[azure_cli]: https://docs.microsoft.com/cli/azure
[azure_portal]: https://portal.azure.com
[azure_sub]: https://azure.microsoft.com/free/
[event_handler_storage_queues]: https://learn.microsoft.com/azure/event-grid/handler-storage-queues
[event_handlers]: https://learn.microsoft.com/azure/event-grid/overview#event-handlers
[product_docs]: https://learn.microsoft.com/azure/event-grid/overview
[system_topics]: https://learn.microsoft.com/azure/event-grid/system-topics
 [source]: https://aka.ms/azsdk/go/systemevents/src
 [godoc_examples]: https://aka.ms/azsdk/go/systemevents/pkg#pkg-examples
 [godoc]: https://aka.ms/azsdk/go/systemevents/pkg
