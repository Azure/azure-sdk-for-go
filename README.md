# Azure SDK for Go

[![godoc](https://godoc.org/github.com/Azure/azure-sdk-for-go?status.svg)](https://godoc.org/github.com/Azure/azure-sdk-for-go)
[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/Azure.azure-sdk-for-go?branchName=main)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=640&branchName=main)

This repository is for active development of the Azure SDK for Go. For consumers of the SDK you can follow the links below to visit the documentation you are interested in
* [Overview of Azure SDK for Go](https://docs.microsoft.com/azure/developer/go/)
* [SDK Reference](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go)
* [Code Samples for Azure Go SDK](https://github.com/azure-samples/azure-sdk-for-go-samples)
* [Azure REST API Docs](https://docs.microsoft.com/rest/api/)
* [General Azure Docs](https://docs.microsoft.com/azure)

## Getting Started

To get started with a library, see the README.md file located in the library's project folder.  You can find these library folders grouped by service in the `/sdk` directory.

> NOTE: Go **1.18** or later is required.

## Packages available

Each service might have a number of libraries available from each of the following categories:
* [Client - New Releases](#client-new-releases)
* [Client - Previous Versions](#client-previous-versions)
* [Management - New Releases](#management-new-releases)
* [Management - Previous Versions](#management-previous-versions)

### Client: New Releases

We have a new wave of packages that are being announced as **stable** and several that are currently released in **beta**. These libraries allow you to use, consume, and interact with existing resources, for example, uploading a blob. These libraries share a number of core functionalities including retries, logging, transport protocols, authentication protocols, etc. that can be found in the [azcore](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azcore) library. You can learn more about these libraries by reading about the [Azure SDK Go guidelines](https://azure.github.io/azure-sdk/golang_introduction.html).

You can find the most up-to-date list of new packages on our [latest page](https://azure.github.io/azure-sdk/releases/latest/index.html#go). These new libraries can be identified by locating them under the [`sdk`](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk) directory in the repository.

> NOTE: If you need to ensure your code is ready for production use one of the stable, non-beta libraries.

### Client: Previous Versions

The last stable versions of packages that have been provided for usage with Azure are production-ready. These libraries might not implement the [Azure Go SDK guidelines](https://azure.github.io/azure-sdk/golang_introduction.html) or have the same feature set as the New releases, however they do offer a wider coverage of services.

Previous Go SDK packages are located under [/services folder](https://github.com/Azure/azure-sdk-for-go/tree/master/services), and you can see the full list [on this page](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/services). 

### Management: New Releases
A new set of management libraries that follow the [Azure SDK Design Guidelines for Go](https://azure.github.io/azure-sdk/golang_introduction.html) are available at `sdk/resourcemanagement`. These new libraries provide a number of core capabilities that are shared amongst all Azure SDKs, including the intuitive Azure Identity library, an HTTP Pipeline with custom policies, error-handling, distributed tracing, and much more.

To get started, please follow the [quickstart guide here](https://aka.ms/azsdk/go/mgmt). To see the benefits of migrating and how to migrate to the new libraries, please visit the [migration guide](https://aka.ms/azsdk/go/mgmt/migration).

You can find the [most up to date list of all of the new packages on our page](https://azure.github.io/azure-sdk/releases/latest/mgmt/go.html)

> NOTE: If you need to ensure your code is ready for production use one of the stable, non-beta libraries. Also, if you are experiencing authentication issues with the management libraries after upgrading certain packages, it's possible that you upgraded to the new versions of SDK without changing the authentication code. Please refer to the migration guide for proper instructions.

* [Quickstart tutorial for new releases](https://aka.ms/azsdk/go/mgmt). Documentation is also available at each readme file of the individual module (Example: [Readme for Compute Module](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/compute/armcompute))

### Management: Previous Versions
For a complete list of management libraries which enable you to provision and manage Azure resources, please [check here](https://azure.github.io/azure-sdk/releases/latest/all/go.html). They might not have the same feature set as the new releases but they do offer wider coverage of services.

Previous packages are located under [/services folder](https://github.com/Azure/azure-sdk-for-go/tree/master/services), and you can see the full list [on this page](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/services).

* [Quickstart tutorial for previous versions](https://aka.ms/azsdk/go/mgmt/previous)

## Other Azure Go Packages

Azure provides several other packages for using services from Go, listed below.  These packages do NOT follow the New Release guidelines.

| Service              | Import Path/Repo                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------- |
| Storage - Blobs      | [github.com/Azure/azure-storage-blob-go](https://github.com/Azure/azure-storage-blob-go)           |
| Storage - Files      | [github.com/Azure/azure-storage-file-go](https://github.com/Azure/azure-storage-file-go)           |
| Storage - Queues     | [github.com/Azure/azure-storage-queue-go](https://github.com/Azure/azure-storage-queue-go)         |
| Event Hubs           | [github.com/Azure/azure-event-hubs-go](https://github.com/Azure/azure-event-hubs-go)               |
| Application Insights | [github.com/Microsoft/ApplicationInsights-go](https://github.com/Microsoft/ApplicationInsights-go) |

## Samples

More code samples for using the management library for Go SDK can be found in the following locations
- [Go SDK Code Samples Repo(New version)](https://aka.ms/azsdk/go/mgmt/samples)
- [Go SDK Code Samples Repo](https://github.com/azure-samples/azure-sdk-for-go-samples)
- Example files under each package. For example, examples for Network packages can be [found here](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/resourcemanager/network/armnetwork/ze_generated_example_loadbalancernetworkinterfaces_client_test.go)

## Reporting security issues and security bugs

Security issues and bugs should be reported privately, via email, to the Microsoft Security Response Center (MSRC) <secure@microsoft.com>. You should receive a response within 24 hours. If for some reason you do not, please follow up via email to ensure we received your original message. Further information, including the MSRC PGP key, can be found in the [Security TechCenter](https://www.microsoft.com/msrc/faqs-report-an-issue).

## Need help?

* File an issue via [Github Issues](https://github.com/Azure/azure-sdk-for-go/issues)
* Check [previous questions](https://stackoverflow.com/questions/tagged/azure+go) or ask new ones on StackOverflow using `azure` and `go` tags.

## Community

* Chat with us in the **[#Azure SDK
channel](https://gophers.slack.com/messages/CA7HK8EEP)** on the [Gophers
Slack](https://gophers.slack.com/). Sign up
[here](https://invite.slack.golangbridge.org) first if necessary.

## Contribute

See [CONTRIBUTING.md](https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md).

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft trademarks or logos is subject to and must follow [Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general). Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship. Any use of third-party trademarks or logos are subject to those third-party's policies.
