# Azure Event Grid Namespaces Client Module for Go

[Azure Event Grid](https://learn.microsoft.com/azure/event-grid/overview) is a highly scalable, fully managed Pub Sub message distribution service that offers flexible message consumption patterns. For more information about Event Grid see: [link](https://learn.microsoft.com/azure/event-grid/overview).

This client module allows you to publish events and receive events using the [Pull delivery](https://learn.microsoft.com/azure/event-grid/pull-delivery-overview) API.

> NOTE: This client does not work with Event Grid Basic. Use the [azeventgrid.Client][godoc_egbasic_client] in the `azeventgrid` package instead.

Key links:
- [Source code][source]
- [API Reference Documentation][godoc]
- [Product documentation](https://azure.microsoft.com/services/event-grid/)
- [Samples][godoc_examples]

## Getting started

### Install the package

Install the Azure Event Grid Namespaces client module for Go with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces
```

### Prerequisites

- Go, version 1.18 or higher
- An [Azure subscription](https://azure.microsoft.com/free/)
- An [Event Grid namespace][ms_namespace]. You can create an Event Grid namespace using the [Azure Portal][ms_create_namespace].
- An [Event Grid namespace topic][ms_topic]. You can create an Event Grid namespace topic using the [Azure Portal][ms_create_topic].

### Authenticate the client

Event Grid Namespaces clients (SenderClient and ReceiverClient) authenticate using a TokenCredential, provided by the [`azidentity` module](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity), or using a shared key credential. 

Using a TokenCredential:
* **SenderClient**: [NewSenderClient example][godoc_example_newsender].
* **ReceiverClient**: [NewReceiverClient example][godoc_example_newreceiver].

Using a shared key:
* **SenderClient**: [NewSenderClientWithSharedKeyCredential example][godoc_example_newsender_sharedkey].
* **ReceiverClient**: [NewReceiverClientWithSharedKeyCredential example][godoc_example_newreceiver_sharedkey].

# Key concepts

An Event Grid namespace is a container for multiple types of resources, including [**namespace topics**][ms_topic]:
- A [**namespace topic**][ms_topic] contains CloudEvents that you publish, via [SenderClient.SendEvents][godoc_send].
- A [**topic subscription**][ms_subscription], associated with a single topic, can be used to receive events via [ReceiverClient.ReceiveEvents][godoc_receive].

Namespaces also offer access using MQTT, although that is not covered in this package.

# Examples

Examples for various scenarios can be found on [pkg.go.dev][godoc_examples] or in the example*_test.go files in our GitHub repo for [aznamespaces][source].

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
[source]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/messaging/eventgrid/aznamespaces
[godoc]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces
[godoc_send]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces#SenderClient.SendEvents
[godoc_receive]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces#ReceiverClient.ReceiveEvents
[godoc_examples]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces#pkg-examples
[godoc_example_newsender_sharedkey]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces#example-NewSenderClientWithSharedKeyCredential
[godoc_example_newreceiver_sharedkey]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces#example-NewReceiverClientWithSharedKeyCredential
[godoc_example_newsender]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces#example-NewSenderClient
[godoc_example_newreceiver]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces#example-NewReceiverClient
[godoc_egbasic_client]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces#Client
[ms_namespace]: https://learn.microsoft.com/azure/event-grid/concepts-pull-delivery#namespaces
[ms_topic]: https://learn.microsoft.com/azure/event-grid/concepts-pull-delivery#namespace-topics
[ms_subscription]: https://learn.microsoft.com/azure/event-grid/concepts-pull-delivery#event-subscriptions
[ms_create_namespace]: https://learn.microsoft.com/azure/event-grid/create-view-manage-namespaces
[ms_create_topic]: https://learn.microsoft.com/azure/event-grid/create-view-manage-namespace-topics
