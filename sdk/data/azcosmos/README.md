# Azure Cosmos DB SDK for Go

## Introduction

This client library enables client applications to connect to Azure Cosmos DB via the NoSQL API. Azure Cosmos DB is a globally distributed, multi-model database service.

## Status: v2 is under construction

This is the v2 major version of the module and it is **not usable yet**. The v2 surface is being
assembled incrementally so that it can be reviewed as it lands. This release covers the error and
response model, partition keys, client construction, and reading and creating single items; the
item operations are not wired to the driver yet and report that they are not implemented.

v2 replaces the v1 pure-Go implementation with a binding to the shared Rust Cosmos driver, so that
routing, retries, session handling, failover behavior and query fan-out are consistent across the
Cosmos DB SDKs. The decision record is
[`docs/adr/0001-go-v2-uses-ffi.md`](docs/adr/0001-go-v2-uses-ffi.md).

### Building against the driver

The driver binding is behind the `azcosmos_driver` build tag and is **off by default**, so the
default build needs no C toolchain and works with `CGO_ENABLED=0`. In that build the client can be
constructed and the API compiled against, and operations report that they are not implemented.

Selecting the binding needs cgo and the `azcosmos_driver` tag. The native library is committed
under `internal/native/lib/<goos>_<goarch>`, so no separate build step is required for the
platforms that are present:

```sh
CGO_ENABLED=1 CGO_LDFLAGS_ALLOW='^-Wl,-rpath,(@(executable_path|loader_path)|\$ORIGIN)$' \
  go build -tags azcosmos_driver ./...
```

`CGO_LDFLAGS_ALLOW` is required because the library is linked dynamically with executable-relative
rpaths, which cgo rejects by default. That matches how the driver modules in
[azure-cosmos-driver](https://github.com/Azure/azure-cosmos-driver) link: a library shipped beside a
binary is found, and nothing else is searched.

Nothing else is searched on purpose. The build directory is not an rpath, because it is an absolute
path on whoever did the build — a module cache path for anyone consuming this — and baking it into
a shipped binary makes the loader search a directory that someone else may be able to create on the
host it runs on. Running the tests from the tree therefore needs that directory passed explicitly,
which the commands below do; it reaches the test binary and never a consumer's.

**Those committed binaries are temporary and are meant to be deleted.** Checking a build artifact
into the repository is not the plan of record: the distribution design puts each target's library
in its own Go module in [azure-cosmos-driver](https://github.com/Azure/azure-cosmos-driver),
selected by `GOOS`/`GOARCH`. That repository exists and already carries a darwin/arm64 module, but
not yet one for `linux/amd64`, which is the platform CI runs on — so the copies here stand in until
it does.

The trigger to remove them, and the steps, are recorded next to the files in
[`internal/native/lib/README.md`](internal/native/lib/README.md). Only `linux/amd64` and
`darwin/arm64` are present, because those are the platforms actually built and tested; another
platform needs its own build from
[azure_data_cosmos_driver_native](https://github.com/Azure/azure-sdk-for-rust/tree/main/sdk/cosmos/azure_data_cosmos_driver_native)
dropped into the matching directory, which is the cost these files impose.

`internal/native/azurecosmosdriver.h` is the header the driver generates, vendored here and pinned
to the version in `driver.go`. That version is checked against the linked library when a client
first reaches the driver, because a header and a library from different versions do not fail to
compile — they fail as moved struct offsets somewhere far from the cause.

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

# LIBDIR is darwin_arm64 or linux_amd64. It is passed through -ldflags rather than set as an
# environment variable, because macOS strips DYLD_* from anything a protected process spawns.
LIBDIR=$PWD/internal/native/lib/$(go env GOOS)_$(go env GOARCH)

EMULATOR=1 AZCOSMOS_ENDPOINT=http://127.0.0.1:49151/ CGO_ENABLED=1 \
  CGO_LDFLAGS_ALLOW='^-Wl,-rpath,(@(executable_path|loader_path)|\$ORIGIN)$' \
  go test -tags azcosmos_driver -ldflags "-extldflags=-Wl,-rpath,$LIBDIR" \
  -run TestEmulator ./...
```

The container the tests use is declared in `internal/testdata/emulator-config.json`; its ids
default to `itemdb` and `items` and can be overridden with `AZCOSMOS_DATABASE` and
`AZCOSMOS_CONTAINER`.

## Getting Started

### Prerequisites

* An Azure subscription or free Azure Cosmos DB trial account
* A C toolchain, but only to build with the `azcosmos_driver` tag

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
