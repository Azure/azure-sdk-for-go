# Azure SDK for Go

[![godoc](https://godoc.org/github.com/Azure/azure-sdk-for-go?status.svg)](https://godoc.org/github.com/Azure/azure-sdk-for-go)
[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/Azure.azure-sdk-for-go?branchName=main)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=640&branchName=main)

This repository is for active development of the Azure SDK for Go. For consumers of the SDK you can follow the links below to visit the documentation you are interested in
* [Overview of Azure SDK for Go](https://docs.microsoft.com/azure/developer/go/)
* [SDK Reference](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go)
* [Code Samples for Azure Go SDK](https://github.com/azure-samples/azure-sdk-for-go-samples)
* [Azure REST API Docs](https://docs.microsoft.com/rest/api/)
* [General Azure Docs](https://docs.microsoft.com/azure)


## New Releases

A new set of management libraries for Go that follow the [Azure SDK Design Guidelines for Go](https://azure.github.io/azure-sdk/golang_introduction.html) are now available. These new libraries provide a number of core capabilities that are shared amongst all Azure SDKs, including the intuitive Azure Identity library, an HTTP Pipeline with custom policies, error-handling, distributed tracing, and much more. Those new libraries can be identified by locating them under the `/sdk` directory in the repo, and you can find the corresponding references [here](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk).

To get started, please follow the [quickstart guide here](https://aka.ms/azsdk/go/mgmt). To see the benefits of migrating to the new libraries, please visit this [migration guide](https://aka.ms/azsdk/go/mgmt/migration) that shows how to transition from older versions of libraries.

You can find the most up to date list of all of the new packages [on this page](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk) as well as our [official releases page](https://azure.github.io/azure-sdk/releases/latest/go.html)

## Previous Versions

Previous Go SDK packages are located under [/services folder](https://github.com/Azure/azure-sdk-for-go/tree/master/services), and you can see the full list [on this page](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/services). They might not have the same feature set as the new releases but they do offer wider coverage of services.


## Getting Started

For instructions and documentation on how to use our Azure SDK for Go, we have provided quick-start tutorials for both new and previous releases. 

* [Quickstart tutorial for new releases](https://aka.ms/azsdk/go/mgmt). Documentation is also available at each readme file of the individual module (Example: [Readme for Compute Module](https://github.com/Azure/azure-sdk-for-go/tree/master/sdk/compute/armcompute))
* [Quickstart tutorial for previous versions](https://aka.ms/azsdk/go/mgmt/previous)

## Other Azure Go Packages

Azure provides several other packages for using services from Go, listed below.
If a package you need isn't available please open an issue and let us know.

| Service              | Import Path/Repo                                                                                   |
| -------------------- | -------------------------------------------------------------------------------------------------- |
| Storage - Blobs      | [github.com/Azure/azure-storage-blob-go](https://github.com/Azure/azure-storage-blob-go)           |
| Storage - Files      | [github.com/Azure/azure-storage-file-go](https://github.com/Azure/azure-storage-file-go)           |
| Storage - Queues     | [github.com/Azure/azure-storage-queue-go](https://github.com/Azure/azure-storage-queue-go)         |
| Service Bus          | [github.com/Azure/azure-service-bus-go](https://github.com/Azure/azure-service-bus-go)             |
| Event Hubs           | [github.com/Azure/azure-event-hubs-go](https://github.com/Azure/azure-event-hubs-go)               |
| Application Insights | [github.com/Microsoft/ApplicationInsights-go](https://github.com/Microsoft/ApplicationInsights-go) |

Packages that are still in public preview can be found under the `services/preview` directory. Please be aware that since these packages are in preview they are subject to change, including breaking changes, outside of a major semver bump.

## Samples

More code samples for using the management library for Go SDK can be found in the following locations
- [Go SDK Code Samples Repo](https://github.com/azure-samples/azure-sdk-for-go-samples)
- Example files under each package. For example, examples for Network packages can be [found here](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/network/armnetwork/example_networkinterfaces_test.go)

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

