# Azure SDK for Go

[![godoc](https://godoc.org/github.com/Azure/azure-sdk-for-go?status.svg)](https://godoc.org/github.com/Azure/azure-sdk-for-go)
[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/Azure.azure-sdk-for-go?branchName=master)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=640&branchName=master)

This repository is for active development of the Azure SDK for Go. For consumers of the SDK you can follow the links below to visit the documentation you are interested in
* [Overview of Azure SDK for Go](https://docs.microsoft.com/azure/developer/go/)
* [SDK Reference](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go)
* [Code Samples for Azure Go SDK](https://github.com/azure-samples/azure-sdk-for-go-samples)
* [Azure REST API Docs](https://docs.microsoft.com/rest/api/)
* [General Azure Docs](https://docs.microsoft.com/azure)


## New Releases

A new set of management libraries for Go that follow the [Azure SDK Design Guidelines for Go](https://azure.github.io/azure-sdk/golang_introduction.html) are now available. These new libraries provide a number of core capabilities that are shared amongst all Azure SDKs, including the intuitive Azure Identity library, an HTTP Pipeline with custom policies, error-handling, distributed tracing, and much more.
To get started, please follow the [quickstart guide] here (https://aka.ms/azsdk/go/mgmt). 
To see the benefits of migrating to the new libraries, please visit this migration guideIn addition, a migration guide that shows how to transition from older versions of libraries is located [here](todo-addthislink).
Those new libraries can be identified by locating them under the [/sdk folder](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk) in the repo

You can find the most up to date list of all of the new packages [on this page](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk) as well as our [official releases page](https://azure.github.io/azure-sdk/releases/latest/go.html)

## Previous Versions

Previous Go SDK packages are located under [/services folder](https://github.com/Azure/azure-sdk-for-go/tree/master/services), and you can see the full list [on this page](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/services). They might not have the same feature set as the new releases but they do offer wider coverage of services.


## Getting Started

For instructions and documentation on how to use our Azure SDK for Go, we have provided quickstart tutorials for both new and previous releases. 

* [Quickstart tutorial for new releases](https://github.com/Azure/azure-sdk-for-go/blob/master/documentation/previous-versions-quickstart.md). Documentation is also available at each readme file of the individual module (Example: [Readme for Compute Module](https://github.com/Azure/azure-sdk-for-go/tree/master/sdk/compute/armcompute))
* [Quickstart tutorial for previous versions](https://github.com/Azure/azure-sdk-for-go/blob/master/documentation/previous-versions-quickstart.md)

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

## Contributing

For details on contributing to this repository, see the [contributing guide](https://github.com/Azure/azure-sdk-for-go/blob/master/CONTRIBUTING.md).

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit
https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repositories using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

[samples_repo]: https://github.com/Azure-Samples/azure-sdk-for-go-samples
