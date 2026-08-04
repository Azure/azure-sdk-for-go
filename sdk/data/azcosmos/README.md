# Azure Cosmos DB SDK for Go

## Introduction

This client library enables client applications to connect to Azure Cosmos DB via the NoSQL API. Azure Cosmos DB is a globally distributed, multi-model database service.

## Status: v2 is under construction

This is the v2 major version of the module and it is **not usable yet**. The v2 surface is being
assembled incrementally so that it can be reviewed as it lands. This release covers the error and
response model, partition keys, client construction, and reading and creating single items; the
operations are not wired to the driver yet and report that they are not implemented.

v2 replaces the v1 pure-Go implementation with a binding to the shared Rust Cosmos driver, so that
routing, retries, session handling, failover behavior and query fan-out are consistent across the
Cosmos DB SDKs. The decision record is ADR 0001, added by
[PR 27238](https://github.com/Azure/azure-sdk-for-go/pull/27238).

Because v2 binds to a native driver, it requires cgo (`CGO_ENABLED=1`). v1's WebAssembly support
does not carry over.

Until v2 is ready, use v1:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos
```

Usage documentation and samples return to this README as the v2 surface lands. Migration guidance
from v1 is in [`migrationguide.md`](migrationguide.md).

## Getting Started

### Prerequisites

* Go 1.25 or higher
* An Azure subscription or free Azure Cosmos DB trial account
* A C toolchain, because v2 requires cgo

Note: If you don't have an Azure subscription, create a free account before you begin.
You can Try Azure Cosmos DB for free without an Azure subscription, free of charge and commitments, or create an Azure Cosmos DB free tier account, with the first 400 RU/s and 5 GB of storage for free. You can also use the Azure Cosmos DB Emulator with a URI of https://localhost:8081. For the key to use with the emulator, see [how to develop with the emulator](https://learn.microsoft.com/azure/cosmos-db/how-to-develop-emulator).

### Create an Azure Cosmos DB account

You can create an Azure Cosmos DB account using:

* [Azure Portal](https://portal.azure.com).
* [Azure CLI](https://learn.microsoft.com/cli/azure).
* [Azure ARM](https://learn.microsoft.com/azure/cosmos-db/quick-create-template).

## Next steps

- [Resource Model of Azure Cosmos DB Service](https://learn.microsoft.com/azure/cosmos-db/sql-api-resources)
- [Azure Cosmos DB Resource URI](https://learn.microsoft.com/rest/api/documentdb/documentdb-resource-uri-syntax-for-rest)
- [Partitioning](https://learn.microsoft.com/azure/cosmos-db/partition-data)
- [Using emulator](https://github.com/Azure/azure-documentdb-dotnet/blob/master/docs/documentdb-nosql-local-emulator.md)

## License

This project is licensed under MIT.

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Cosmos` label.

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a Contributor License
Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For
details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate
the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to
do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
