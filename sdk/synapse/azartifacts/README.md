# Azure Synapse Artifacts Module for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/Azure/azure-sdk-for-go/sdk/synapse/azartifacts)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/synapse/azartifacts)

The `azartifacts` module provides operations for working with Azure synapse.

[Source code](https://github.com/Azure/azure-sdk-for-go/tree/master/sdk/synapse/azartifacts)

# Getting started

## Prerequisites

- an [Azure subscription](https://azure.microsoft.com/free/)
- Go 1.13 or above

## Install the package

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

Install the Azure Synapse Artifacts module:

```sh
go get github.com/Azure/azure-sdk-for-go/sdk/synapse/azartifacts
```

## Authorization

When creating a client, you will need to provide a credential for authenticating with Azure Synapse.  The `azidentity` module provides facilities for various ways of authenticating with Azure including client/secret, certificate, managed identity, and more.

```go
cred, err := azidentity.NewDefaultAzureCredential(nil)
```

For more information on authentication, please see the documentation for `azidentity` at [pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity).

## Connecting to Azure Synapse

Once you have a credential, create a connection to the desired Synapse endpoint URL.

```go
con := azartifacts.NewConnection("<endpoint>", cred, nil)
```

## Clients

Azure Synapse artifacts modules consist of one or more clients.  A client groups a set of related APIs, providing access to its functionality for the associated endpoint.  Create one or more clients to access the APIs you require using your `azartifacts.Connection`.

```go
client := azartifacts.NewNotebookClient(con)
```

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Synapse` label.

# Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.microsoft.com](https://cla.microsoft.com).

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any
additional questions or comments.
