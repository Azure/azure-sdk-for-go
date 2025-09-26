# Azure App Configuration Client Module for Go

Azure App Configuration is a managed service that helps developers centralize their application and feature settings simply and securely.
It allows you to create and manage application configuration settings and retrieve their revisions from a specific point in time.

[Source code][azappconfig_src] | [Package (pkg.go.dev)][azappconfig] | [Product documentation][appconfig_docs]

## Getting started

### Install the module

Install `azappconfig` with `go get`:

```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2
```

### Prerequisites

* An [Azure subscription][azure_sub]
* A supported Go version (the Azure SDK supports the two most recent Go releases)
* A Configuration store. If you need to create one, see the App Configuration documentation for instructions on doing so in the [Azure Portal][appconfig_portal] or with the [Azure CLI][appconfig_cli].

### Authentication

Azure App Configuration supports authenticating with Azure Active Directory and connection strings. To authenticate with Azure Active Directory, use the [azappconfig.NewClient][azappconfig_newclient] constructor and to authenticate with a connection string, use the [azappconfig.NewClientFromConnectionString][azappconfig_newclientfromconnectionstring] constructor. For simplicity, the examples demonstrate authenticating with a connection string.

See the [azidentity][azure_identity] documentation for more information about possible Azure Active Directory credential types.

## Key concepts

A [Setting][azappconfig_setting] is the fundamental resource within a Configuration Store. In its simplest form, it is a key and a value. However, there are additional properties such as the modifiable content type and tags fields that allow the value to be interpreted or associated in different ways.

The [Label][label_concept] property of a Setting provides a way to separate Settings into different dimensions. These dimensions are user defined and can take any form. Some common examples of dimensions to use for a label include regions, semantic versions, or environments. Many applications have a required set of configuration keys that have varying values as the application exists across different dimensions.

For example, MaxRequests may be 100 in "NorthAmerica" and 200 in "WestEurope". By creating a Setting named MaxRequests with a label of "NorthAmerica" and another, only with a different value, with a "WestEurope" label, an application can seamlessly retrieve Settings as it runs in these two dimensions.

**Tags** provide additional metadata for configuration settings and enable powerful filtering capabilities. Tags are key-value pairs that can be used to categorize and query settings. For example, you can tag settings with `environment=production` or `version=1.2.3` and then use the `TagsFilter` in `SettingSelector` to retrieve only settings that match specific tag criteria.

## Examples

Examples for various scenarios can be found on [pkg.go.dev][azappconfig_examples] or in the `example*_test.go` files in our GitHub repo for [azappconfig][azappconfig_src].

## Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][code_of_conduct].
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact opencode@microsoft.com with any additional questions or comments.

[azure_identity]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[code_of_conduct]: https://opensource.microsoft.com/codeofconduct/
[appconfig_docs]: https://learn.microsoft.com/azure/azure-app-configuration/
[appconfig_portal]: https://learn.microsoft.com/azure/azure-app-configuration/quickstart-azure-app-configuration-create?tabs=azure-portal
[appconfig_cli]: https://learn.microsoft.com/azure/azure-app-configuration/quickstart-azure-app-configuration-create?tabs=azure-cli
[label_concept]: https://learn.microsoft.com/azure/azure-app-configuration/concept-key-value#label-keys
[azappconfig]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig
[azappconfig_newclient]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig#NewClient
[azappconfig_newclientfromconnectionstring]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig#NewClientFromConnectionString
[azappconfig_examples]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig#pkg-examples
[azappconfig_setting]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig#Setting
[azappconfig_src]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/data/azappconfig


