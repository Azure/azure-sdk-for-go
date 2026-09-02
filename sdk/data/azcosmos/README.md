# Azure Cosmos DB SDK for Go

## Introduction

This client library enables client applications to connect to Azure Cosmos DB via the NoSQL API. Azure Cosmos DB is a globally distributed, multi-model database service.

## Status: v2 is under construction

This is the v2 major version of the module and it is **not usable yet**. The v2 surface is being
assembled incrementally so that it can be reviewed as it lands. This release covers the error and
response model, partition keys, client construction, and reading and creating single items.

v2 replaces the v1 pure-Go implementation with a binding to the shared Rust Cosmos driver, so that
routing, retries, session handling, failover behavior and query fan-out are consistent across the
Cosmos DB SDKs. The decision record is
[`docs/adr/0001-go-v2-uses-ffi.md`](docs/adr/0001-go-v2-uses-ffi.md).

### Building with the driver

The driver binding is selected automatically when cgo is enabled on glibc `linux/amd64` or
`darwin/arm64`. No build tag or linker environment variable is required:

```sh
go build ./...
```

The native archives are committed as `.syso` files in target-specific internal packages, so Go
preserves them in module zips and vendored builds and links the matching package automatically. No
separate driver build or runtime sidecar is required; the resulting executable is self-contained.

`CGO_ENABLED=0` and unsupported platforms select `driver_stub.go`. That diagnostic build keeps the
API compilable, but operations report that the driver is unavailable.

Alpine and other musl-based Linux distributions are not supported yet. The bundled Linux archive
is built for glibc, and the build reports that limitation explicitly rather than linking Rust code
compiled for a different libc ABI.

**Those committed binaries are temporary and are meant to be deleted.** Checking a build artifact
into the repository is not the plan of record: the distribution design puts each target's library
in its own Go module in [azure-cosmos-driver](https://github.com/Azure/azure-cosmos-driver),
selected by `GOOS`/`GOARCH`. That repository exists and already carries a darwin/arm64 module, but
not yet one for `linux/amd64`, which is the platform CI runs on — so the copies here stand in until
it does.

The trigger to remove them, and the steps, are recorded in
[`internal/native/lib/README.md`](internal/native/lib/README.md). Only `linux/amd64` and
`darwin/arm64` are present, because those are the platforms actually built and tested; another
platform needs its own build from
[azure_data_cosmos_driver_native](https://github.com/Azure/azure-sdk-for-rust/tree/main/sdk/cosmos/azure_data_cosmos_driver_native)
dropped into the matching directory, which is the cost these files impose.

`azurecosmosdriver.h` is the header the driver generates, vendored here and pinned to the version
in `driver.go`. That version is checked against the linked archive when a client is constructed,
because a header and a library from different versions do not fail to compile — they fail as moved
struct offsets somewhere far from the cause.

### Client initialization

`NewClient` and `NewClientWithKey` perform no network I/O. Call `Client.Initialize` with a context
to initialize eagerly: it probes the account's HTTP capabilities, fetches and caches account
properties, seeds the routing state, and creates the account transport. An unreachable or
unauthorized account is then reported before the first operation, and the context bounds the work.

Calling `Initialize` is optional. The first operation performs the same initialization lazily when
it has not already run. The diagnostic build has no driver to initialize and reports that through
`Initialize`.

Container metadata is not fetched during client construction because the client does not know which
containers the application will use. The first operation on a container resolves and caches that
container's metadata.

Two limits apply to the driver-backed build today, both upstream gaps rather than choices this
module makes:

* Only account keys work. The C ABI has no constructor for a token credential, so [`NewClient`]
  reports that it is unsupported and [`NewClientWithKey`] is the working path.
* v1's WebAssembly support does not carry over.

### Running the end-to-end tests

The tests in `emulator_test.go` run real operations against a service. They need a driver-backed
build and the `EMULATOR` environment variable, and they skip otherwise.

They run against the driver's own in-memory emulator, which is the same one the driver's Rust tests
use, so the binding is exercised against what the driver is developed against. It is a plain
process rather than a container, it creates the test database and container from its config, and it
reports its endpoints as JSON on stdout, so the endpoint is read rather than assumed:

```sh
internal/native/lib/linux_amd64/azure_data_cosmos_emulator \
  --config path/to/azcosmos/internal/testdata/emulator-config.json
# {"event":"ready","accountEndpoint":"http://127.0.0.1:49151/", ...}

EMULATOR=1 AZCOSMOS_ENDPOINT=http://127.0.0.1:49151/ go test -run TestEmulator ./...
```

The container the tests use is declared in `internal/testdata/emulator-config.json`; its ids
default to `itemdb` and `items` and can be overridden with `AZCOSMOS_DATABASE` and
`AZCOSMOS_CONTAINER`.

## Getting Started

### Prerequisites

* An Azure subscription or free Azure Cosmos DB trial account
* A C toolchain on a supported driver platform

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
