# Azure Cosmos DB SDK for Go

<!-- cSpell:ignore azsdk itemdb libc -->

## Introduction

This client library enables client applications to connect to Azure Cosmos DB via the NoSQL API. Azure Cosmos DB is a globally distributed, multi-model database service.

## Status: v2 is under construction

This is the v2 major version of the module and it is **not usable yet**. The v2 surface is being
assembled incrementally so that it can be reviewed as it lands. This release covers the error and
response model, partition keys, client construction, and creating, reading, replacing, upserting,
deleting, and patching single items, plus paged queries within a complete logical partition.

v2 replaces the v1 pure-Go implementation with a binding to the shared Rust Cosmos driver, so that
routing, retries, session handling, failover behavior and query fan-out are consistent across the
Cosmos DB SDKs.

### Building with the driver

The driver binding is selected automatically when cgo is enabled on glibc `linux/amd64` or
`darwin/arm64`. The binding uses the published **native ABI 0.2.0** platform modules.
No local replacements, build tag, or linker environment variable is required:

```sh
go build ./...
```

The native archives are distributed in target-specific modules from
[azure-cosmos-driver](https://github.com/Azure/azure-cosmos-driver), so Go preserves them in module
zips and vendored builds and links the matching package automatically. No separate driver build or
runtime sidecar is required; the resulting executable is self-contained.

`CGO_ENABLED=0` and unsupported platforms select `driver_stub.go`. That diagnostic build keeps the
API compilable, but operations report that the driver is unavailable.

Alpine and other musl-based Linux distributions are not supported yet. The Linux archive
is built for glibc, and the build reports that limitation explicitly rather than linking Rust code
compiled for a different libc ABI.

`azurecosmosdriver.h` is copied byte-for-byte from the pinned driver distribution modules because
this package's cgo files need the ABI declarations in their own include path. The header version
and linked library version must both match the pin in `driver.go` before any struct-sensitive ABI
call during construction. Mismatched versions can otherwise cause incompatible struct layouts.

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

One limit applies to the driver-backed build today: v1's WebAssembly support does not carry over.

### Operation defaults and shared runtimes

`OperationOptions` is accepted by each item request, query page fetch, `ClientOptions.Operation`, and
`RuntimeOptions.Operation`. The Rust driver resolves declared environment overrides above request,
client, runtime, ordinary supported environment settings, and driver defaults, in that order.
Database and container handles do not introduce extra configuration layers.

Nil pointers inherit; explicit false and zero values override. Nil excluded-region slices and
custom-header maps inherit; non-nil empty values clear the inherited list/map. Nonempty values
replace rather than extend them. Throughput and throttling members inherit independently.
Binary encoding and availability strategies replace their entire groups.

Leaving `ClientOptions.Runtime` nil uses a lazily initialized process-wide runtime, as in Rust.
`NewRuntime` creates an explicit runtime with construction-time defaults. Multiple clients for the
same account may share it, including after another client closes. Native caches retain account
references and their credentials: a cached container may use the identity that first resolved it,
not the identity supplied to a later client. Use different explicit runtimes when credentials
must be isolated. The global runtime is retained for the process lifetime.
Configure client content-response defaults through `ClientOptions.Operation.EnableContentResponseOnWrite`.

Runtime defaults cannot be updated after construction. Requests capture a native snapshot before
lazy initialization. `EndToEndTimeout` is optional: nil inherits; explicit zero or subsecond values
clamp to one second at every scope. A caller's context deadline can stop waiting earlier without
cancelling native work. Throttling retry budgets apply per transport invocation, not to
the entire logical operation. The hedging master switch and its environment override can override
the chosen availability strategy. Environment settings are captured at runtime construction.

`Client.Close` drains and closes only that client. It does not cancel credentials retained by a
shared native cache. `Runtime.Close` releases the owner's reference and rejects new attachments;
existing clients keep working until they close. Native resources are freed after the last owner
and attached client release them. Both methods are idempotent and concurrency-safe.
Do not call Close from inside a credential callback.

### Context cancellation

Native ABI 0.2.0 does not support cancelling submitted operations. Cancelling a Go context
stops the caller's wait, but the native request continues and a write may still commit.
An already delivered completion takes precedence over cancellation. Otherwise, a cancelled
wait returns `CodeOperationCancelled` wrapping the context error, with the outcome unknown.
Response metadata from a later completion cannot be returned to that caller.

The binding drains abandoned waits in the background and retains native handles and credentials
until completion. `Client.Close` waits for that drain and can take longer
than the caller's context deadline. End-to-end timeout snapshots bound native item execution,
but a shorter Go deadline does not cancel native initialization or metadata resolution.
For retryable client-side patches, supply `TrackingID` before submission to retain the identity
after cancellation; a native-generated tracking ID is only available after completion.

### Binary response compatibility change

**Raw item response bytes now use the driver's binary JSON default**, including for existing
callers that leave options unset. Rust-compatible `ClientOptions.BinaryEncoding` resolves once from
an explicit group, then `AZURE_COSMOS_BINARY_ENCODING_ENABLED`, then enabled by default.
Create/read/replace/upsert resolve request encoding over that SDK client default rather than
runtime/client `Operation.BinaryEncoding`. PATCH and delete retain native layered resolution.
Before using `encoding/json`, select text responses on the request or dedicated client group:

```go
options := azcosmos.OperationOptions{
    BinaryEncoding: &azcosmos.BinaryEncodingOptions{
        RequestTextResponse: true,
    },
}
```

This keeps binary wire encoding while asking the driver to return text JSON. Alternatively,
`&azcosmos.BinaryEncodingOptions{Enabled: to.Ptr(false)}` disables binary wire encoding
(`to` is `github.com/Azure/azure-sdk-for-go/sdk/azcore/to`).
The zero-valued group enables binary, as Rust's default does. Text conversion preserves JSON
values, not necessarily byte-for-byte formatting. Go never decodes application item schemas.

### SDK identity

Service requests include `azsdk-go-azcosmos/<version>` in the User-Agent header alongside the native
Cosmos driver identity and feature flags. `RuntimeOptions.ApplicationID` is the optional application
suffix; SDK identity does not consume its length allowance. Per-client suffix overrides require
native support not present in v0.2.0, so they are not exposed.

### Additional configuration

`RuntimeOptions.CPURefreshInterval` configures the existing native CPU/memory sampler (one to sixty
seconds). `ClientOptions.FaultInjectionRules` forwards copied, typed test rules through the existing
versioned native driver-options API. Rules can select item or metadata operations, transport and
region, predefined failures or custom HTTP responses, delays, probability, hit limits, and lifetime.

`Response.Diagnostics` and `Error.Diagnostics` contain immutable native snapshots when available.
`ClientOptions.DiagnosticsHandler` observes completed Go item calls after the client lifetime guard
is released. It receives absent native diagnostics when a call stops waiting before completion.
JSON rendering errors are reported through `Diagnostics.Err()` without changing operation outcomes.
Snapshot JSON uses native default verbosity and a driver-defined schema.

Connection-pool configuration, partition-failover tuning, backup endpoints, runtime diagnostic
retention/verbosity configuration, and per-client identity overrides remain deferred: the published
ABI cannot carry those settings. Submillisecond precision is also limited by its millisecond fields.
Constructors stay network-free; use `Initialize(ctx)` for eager account initialization.

### Patching items

> [!IMPORTANT]
> `PatchItem`, `PatchOperations`, and the related PATCH options are provisional while PATCH support
> in the shared driver is in preview. They may change or be removed before azcosmos/v2 reaches a
> stable release.

`PatchOperations` owns a JSON snapshot of each value when it is appended and supports `add`, `set`,
`replace`, `remove`, `incr`, and `move`. Paths use RFC 6901 JSON Pointer syntax. There is no
Go-side ten-operation limit: the driver automatically chooses a server PATCH or a client-side
read-modify-write execution strategy by default. `PatchItemOptions.Strategy` can explicitly select
automatic, client-side, or server-side execution. Explicit server-side execution does not fall
back to read-modify-write when a request exceeds the service limit.

Client-side execution of a patch that is not intrinsically retry-safe permanently adds the
`_azsdkPatchTracking` property to the item. The driver uses it to deduplicate retries within one
`PatchItem` call. `PatchItemOptions.TrackingID` accepts a UUID (simple, hyphenated, braced, or URN) to reuse across
application retries after an ambiguous outcome. The effective ID is returned on `ItemResponse`
or `Error`, including cancellation, when supplied by the driver. Duplicate suppression is bounded
by tracking capacity and retention; it is not a permanent exactly-once guarantee.

`MaxAttempts` (1..255), `TrackingCapacity` (1..65535), and `TrackingRetention` configure client-side
patching and are inert for server-side execution. Retention floors to whole seconds with a minimum
of one second even for explicit zero and a maximum of uint32 seconds; larger values clamp.
Nil uses the driver default. Capacity pressure can evict tracking entries sooner.
PATCH exposes only If-Match. Other item APIs expose mutually exclusive If-Match and If-None-Match
preconditions, whose service support depends on the operation. Patch-specific `Strategy` overrides
the shared `Operation.PatchStrategy`.

### Querying items

`ContainerClient.NewQueryItemsPager` accepts a `Query`, an explicit `FeedScope`, and optional
`QueryOptions`. Use `NewQuery(sql).WithParameter(name, value)` to capture JSON parameter values,
and `NewFeedScopeForPartitionKey(pk)` to target a complete logical partition. For hierarchical
partition keys, supply every component. Query values are immutable; `WithParameter` returns a
new query and an error if the value cannot be serialized.

The API follows the Rust SDK's separation of query, scope, and feed options, with Go's standard
`More`/`NextPage` pager. `QueryOptions.Feed.PageSizeHint` is a positive page-size hint; zero leaves
sizing to the driver, and negative hints are rejected. Results are raw JSON values in
`QueryItemsResponse.Items`, including scalar `SELECT VALUE` results. An empty page does not
necessarily end the query: use `More`, not the number of items.

Queries default to text wire encoding, overriding inherited binary-encoding preferences so the
pager can split the JSON feed envelope. To enable binary wire encoding for a query, set both
`QueryOptions.Operation.BinaryEncoding.Enabled` and `RequestTextResponse` to true.
Explicit raw binary responses are rejected. Other operation options inherit normally; each
`NextPage` captures a fresh runtime snapshot and one budget for initialization, metadata, and query.

Save `QueryItemsResponse.ContinuationToken` and pass it in `QueryOptions.Feed` to resume with
the same query and scope. This is an opaque **driver planner token**, not the service's
`x-ms-continuation` header. Do not parse or modify it. An empty token indicates completion.
Each fetch resolves its container reference through the native driver's cache, allowing queries
started after container recreation to use the replacement. Tokens from the deleted container
cannot be used to resume against the replacement.
With the currently pinned azcore pager, a failed `NextPage` does not advance the continuation;
callers must handle the error rather than blindly continuing a `More` loop.

Pager construction performs no network I/O and snapshots its inputs. Each pager's first fetch
(including a resumed pager) performs an additional container-metadata read to verify that the
key is complete. Its request charge is included in the first fetch's response or error, and
it shares the page's context and end-to-end timeout. This requires permission to read container
metadata. Pagers are not safe for concurrent use; independent pagers can share a client.

Cross-partition and hierarchical-prefix queries are intentionally not exposed yet. The pinned
native ABI returns only the first item of pre-split driver pages, so those pipelines cannot be
safely exposed as general query support without a native fix. Complete logical-partition
queries use the driver's direct request pipeline and return complete JSON feed envelopes.

### Running the end-to-end tests

The tests in `emulator_test.go` run real operations against a service. They need a driver-backed
build and the `EMULATOR` environment variable, and they skip otherwise.

They run against the driver's own in-memory emulator, which is the same one the driver's Rust tests
use, so the binding is exercised against what the driver is developed against. It is a plain
process rather than a container, it creates the test database and container from its config, and it
reports its endpoints as JSON on stdout, so the endpoint is read rather than assumed:

Build the emulator from the pinned Rust source using the
[emulator build instructions](internal/testdata/README.md). No emulator executable is checked
into this module. CI builds it with a pinned Rust toolchain and locked Cargo dependencies.

```sh
/path/to/emulator-target/release/azure_data_cosmos_emulator \
  --config path/to/azcosmos/internal/testdata/emulator-config.json
# {"event":"ready","accountEndpoint":"http://127.0.0.1:49151/", ...}

EMULATOR=1 AZCOSMOS_ENDPOINT=http://127.0.0.1:49151/ go test -run TestEmulator ./...
```

The container the tests use is declared in `internal/testdata/emulator-config.json`; its ids
default to `itemdb` and `items` and can be overridden with `AZCOSMOS_DATABASE` and
`AZCOSMOS_CONTAINER`. Query scope tests also use the `query-hierarchical` container declared in
the same configuration, under the selected database.

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
