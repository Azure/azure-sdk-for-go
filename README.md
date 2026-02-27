# Azure SDK for Go

[![godoc](https://godoc.org/github.com/Azure/azure-sdk-for-go?status.svg)](https://godoc.org/github.com/Azure/azure-sdk-for-go)

This repository is for active development of the Azure SDK for Go. For consumers of the SDK you can follow the links below to visit the documentation you are interested in
* [Overview of Azure SDK for Go](https://docs.microsoft.com/azure/developer/go/)
* [SDK Reference](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk)
* [Code Samples for Azure Go SDK](https://github.com/azure-samples/azure-sdk-for-go-samples)
* [Azure REST API Docs](https://docs.microsoft.com/rest/api/)
* [General Azure Docs](https://docs.microsoft.com/azure)
* [Share your feedback to our Azure SDK](https://www.surveymonkey.com/r/FWPGFGG)

## Getting Started

To get started with a module, see the README.md file located in the module's project folder.  You can find these module folders grouped by service in the `/sdk` directory.

<a id="go-version-support"></a>
> [!IMPORTANT]
> Our libraries are compatible with the two most recent major Go releases, the same [policy](https://go.dev/doc/devel/release#policy) the Go programming language follows.

> [!IMPORTANT]
> Projects are highly encouraged to use the latest version of Go. This ensures your product has all the latest security fixes and is included in [Go's support lifecycle](https://go.dev/doc/devel/release).

> [!WARNING]
> The [root azure-sdk-for-go Go module](https://godoc.org/github.com/Azure/azure-sdk-for-go) which contains subpaths of `/services/**/mgmt/**` (also known as track 1) is [deprecated and no longer recieving support](https://azure.github.io/azure-sdk/releases/deprecated/go.html). See [the migration guide](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/ARM/MIGRATION_GUIDE.md) to learn how to migrate to the current version.

## Packages available

Each service can have both 'client' and 'management' modules. 'Client' modules are used to consume the service, whereas 'management' modules are used to configure and manage the service.

### Client modules

Our client modules follow the [Azure Go SDK guidelines](https://azure.github.io/azure-sdk/golang_introduction.html). These modules allow you to use, consume, and interact with existing resources, for example, uploading a blob. They also share a number of core functionalities including retries, logging, transport protocols, authentication protocols, etc. that can be found in the [azcore](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azcore) module.

You can find the most up-to-date list of new modules on our [latest page](https://azure.github.io/azure-sdk/releases/latest/index.html#go).

> [!NOTE]
> If you need to ensure your code is ready for production use one of the stable, non-beta modules.

### Management modules
Similar to our client modules, the management modules follow the [Azure Go SDK guidelines](https://azure.github.io/azure-sdk/golang_introduction.html). All management modules are available at `sdk/resourcemanager`. These modules provide a number of core capabilities that are shared amongst all Azure SDKs, including the intuitive Azure Identity module, an HTTP Pipeline with custom policies, error-handling, distributed tracing, and much more.

To get started, please follow the [quickstart guide here](https://aka.ms/azsdk/go/mgmt). To see the benefits of migrating and how to migrate to the new modules, please visit the [migration guide](https://aka.ms/azsdk/go/mgmt/migration).

You can find the [most up to date list of all of the new packages on our page](https://azure.github.io/azure-sdk/releases/latest/mgmt/go.html)

> [!NOTE]
> If you need to ensure your code is ready for production use one of the stable, non-beta modules. Also, if you are experiencing authentication issues with the management modules after upgrading certain packages, it's possible that you upgraded to the new versions of SDK without changing the authentication code. Please refer to the migration guide for proper instructions.

* [Quickstart tutorial for new releases](https://aka.ms/azsdk/go/mgmt). Documentation is also available at each readme file of the individual module (Example: [Readme for Compute Module](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/compute/armcompute))

## Samples

More code samples for using the management modules for Go SDK can be found in the following locations
- [Go SDK Code Samples Repo](https://aka.ms/azsdk/go/mgmt/samples)
- Example files under each package. For example, examples for Network packages can be [found here](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/resourcemanager/network/armnetwork/loadbalancernetworkinterfaces_client_example_test.go)

## Historical releases

Note that the latest modules from Microsoft are grouped by service in the `/sdk` directory. If you're using packages with prefix `github.com/Azure/azure-sdk-for-go/services`, `github.com/Azure/azure-sdk-for-go/storage`, `github.com/Azure/azure-sdk-for-go/profiles`, please consider migrating to the latest modules. You can find a mapping table from these historical releases to their equivalent [here](https://azure.github.io/azure-sdk/releases/deprecated/index.html#go). 

## Reporting security issues and security bugs

Security issues and bugs should be reported privately, via email, to the Microsoft Security Response Center (MSRC) <secure@microsoft.com>. You should receive a response within 24 hours. If for some reason you do not, please follow up via email to ensure we received your original message. Further information, including the MSRC PGP key, can be found in the [Security TechCenter](https://www.microsoft.com/msrc/faqs-report-an-issue).

## Need help?

* File an issue via [Github Issues](https://github.com/Azure/azure-sdk-for-go/issues)
* Check [previous questions](https://stackoverflow.com/questions/tagged/azure+go) or ask new ones on StackOverflow using `azure` and `go` tags.

## Data Collection
The software may collect information about you and your use of the software and send it to Microsoft. Microsoft may use this information to provide services and improve our products and services. You may turn off the telemetry as described below. You can learn more about data collection and use in the help documentation and Microsoft’s [privacy statement](https://go.microsoft.com/fwlink/?LinkID=824704). For more information on the data collected by the Azure SDK, please visit the [Telemetry Guidelines](https://azure.github.io/azure-sdk/general_azurecore.html#telemetry-policy) page.

### Telemetry Configuration
Telemetry collection is on by default.

To opt out, you can disable telemetry at client and credential construction. Set `Disabled` to true in `ClientOptions.Telemetry`. This will disable telemetry for all methods in the client. Do this for every new client and credential created.

The example below uses the `azblob` module. In your code, you can replace `azblob` with the package you are using.

```go
package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func main() {
	// set http client options
	clientOpts := policy.ClientOptions{
		Telemetry: policy.TelemetryOptions{
			Disabled: true,
		},
	}
	// set identity client options
	credOpts := azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: clientOpts,
	}
	// set service client options
	azblobOpts := azblob.ClientOptions{
		ClientOptions: clientOpts,
	}

	// authenticate with Microsoft Entra ID
	cred, err := azidentity.NewManagedIdentityCredential(&credOpts)
	// TODO: handle error

	// create a client for the specified storage account
	client, err := azblob.NewClient(account, cred, &azblobOpts)
	// TODO: handle error
  	// TODO: do something with the client
}
```
> [!NOTE]
> Please note that `AzureDeveloperCLICredential` and `AzureCLICredential` do not include `ClientOptions.Telemetry`. Therefore, it is unnecessary to set options in these credentials.


## Community

* Chat with us in the **[#Azure SDK
channel](https://gophers.slack.com/messages/CA7HK8EEP)** on the [Gophers
Slack](https://gophers.slack.com/). Sign up
[here](https://invite.slack.golangbridge.org) first if necessary.

## Contribute

See [CONTRIBUTING.md](https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md).

For AI agents and automated tools, see [AGENTS.md](https://github.com/Azure/azure-sdk-for-go/blob/main/AGENTS.md) for guidance on repository workflows, automation boundaries, and best practices.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/). For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft trademarks or logos is subject to and must follow [Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general). Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship. Any use of third-party trademarks or logos are subject to those third-party's policies.
