# Azure Cosmos DB SDK for Go

<!-- cSpell:ignore azsdk failback itemdb libc PPAF PPCB -->

## Introduction

This client library enables client applications to connect to Azure Cosmos DB via the NoSQL API. Azure Cosmos DB is a globally distributed, multi-model database service.

## Status: v2 is under construction

This is the v2 major version of the module and it is **not usable yet**. The v2 surface is being
assembled incrementally so that it can be reviewed as it lands. This release covers the error and
response model, partition keys, client construction, and creating, reading, replacing, upserting,
deleting, and patching single items, database and container lifecycle operations, transactional
batches with complete hierarchical partition keys, and paged queries.

v2 replaces the v1 pure-Go implementation with a binding to the shared Rust Cosmos driver, so that
routing, retries, session handling, failover behavior and query fan-out are consistent across the
Cosmos DB SDKs.

### Building with the driver

This module requires **Go 1.27 or later** for concrete generic methods. Unlike the
repository's usual two-release support policy, this provisional API does not support Go 1.26.
Its `ci.yml` pins analysis compiler and linter versions. Repository-wide PR analysis
honors those pins only while checking this module; other packages keep their existing tools.

The driver binding is selected automatically when cgo is enabled on glibc `linux/amd64` or
`darwin/arm64`. No build tag or linker environment variable is required:

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

The v0.2.0 driver owns throttle and session retries; Go submits once and exposes the driver's
attempt count and final status on responses and errors. When attached by the driver,
`Response.Diagnostics` and `Error.Diagnostics` contain a Go-owned snapshot of total charge,
elapsed time, regions contacted, retained per-attempt status/substatus and latency, and the
driver's JSON rendering at `ClientOptions.DiagnosticsVerbosity` (a compact per-region summary by
default; the full per-attempt rendering only when set to `DiagnosticsVerbosityDetailed`). The
native driver can compact old attempts; `AttemptCount`
remains the total even when `len(Diagnostics.Attempts)` is smaller. Diagnostics are not available
when cancellation returns before a native completion.

The caller's session token is forwarded unchanged, and the returned token can be supplied on a
later request. The v0.2 native ABI does not support on-demand operation cancellation: cancelling
a Go context returns promptly and releases its Go-side resources, but the native request can
continue until its own timeout or retry budget ends. Avoid assuming a cancelled write did not
reach the service.
If an item or batch operation's native completion is already available when cancellation is observed,
the SDK returns that outcome, including its response body and ETag, rather than a cancellation error.
When no end-to-end timeout is configured or inherited for an item or batch operation, the remaining context
deadline supplies its native timeout, subject to its one-second minimum. Explicit and inherited timeouts
remain unchanged; an earlier context deadline still limits the Go wait.
Transactional batches use the same cancellation behavior as point operations. A cancellation
or execution error is not proof that writes did not reach the service.

### Client initialization

`NewClient(ctx, ...)` and `NewClientWithKey(ctx, ...)` initialize eagerly: they probe the account's
HTTP capabilities, fetch account properties, seed routing state, and create the account transport.
Failure returns a nil client. The context bounds the wait; cancelled native work drains in the
background before its handles and runtime attachment are released. Builds without native support
return a driver-unavailable error rather than a client.

Container metadata is not fetched during client construction because the client does not know which
containers the application will use. Both `Client.NewContainer(ctx, ...)` and
`DatabaseClient.NewContainer(ctx, ...)` resolve metadata and topology through the native cache
before returning. Cache misses contact the service and can consume RUs; neither factory creates
a resource. Temporary native references are released, and operations continue resolving through
the driver so container recreation is not pinned to a stale reference. `NewDatabase` remains local.

One limit applies to the driver-backed build today: v1's WebAssembly support does not carry over.

## Advanced availability

The native driver owns regional retry and per-partition routing. Per-partition automatic failover
(PPAF) and the per-partition circuit breaker (PPCB) are distinct paths: a request is eligible for
one or the other, never both.

### Per-partition automatic failover (PPAF)

PPAF applies to writes on single-master accounts. It is enabled automatically when the account's
server properties advertise the capability; Go has no PPAF enable/disable option.

PPAF uses service-directed, success-time discovery. A failed attempt can move the current write
through other regions, but the driver does not persist a per-partition routing override from the
failure alone. It records the override only after a later region definitively processes the write,
shown by an HTTP 2xx, 409 Conflict, or 412 Precondition Failed response. The confirming response
causes the driver to apply the deferred failed-region effects and record an override past those
endpoints. Subsequent writes to that physical partition use the recorded override. PPAF overrides
are not governed by PPCB's timer-based probe and failback sweep.

### Per-partition circuit breaker (PPCB)

PPCB applies to reads on any account and to writes on multi-master accounts. The native driver owns
it and enables it by default. Health and failure counts are tracked independently for each physical
partition and region, so a problem with one partition does not move unrelated partitions away from
their preferred region.

The circuit trips after the configured read or write failure threshold is reached. Its counters
reset when the gap since the last recorded failure exceeds the counter-reset window; failures spaced
closer than that can accumulate toward the threshold over a longer period. After the configured
unavailability duration and a bounded background sweep, one request probes the original region. A
successful probe restores preferred routing; a failed probe keeps that partition routed away until
a later probe is eligible.

### Read hedging

Hedging is enabled by default for eligible reads. The built-in threshold is the lesser of one second
and half the request timeout, falling back to one second when no request timeout is available. No
option is needed for this default behavior. Use `OperationOptions.HedgingEnabled` to override the
master switch, or `OperationOptions.AvailabilityStrategy` with `HedgingAvailability` to change the
threshold. Set `ClientOptions.Routing` with `PreferredRegions` so the driver has an explicit
preferred order. For example, this client uses a 500 ms threshold:

```go
func newHedgedClient(
	ctx context.Context,
	endpoint string,
	credential azcore.TokenCredential,
) (*azcosmos.Client, error) {
	return azcosmos.NewClient(ctx, endpoint, credential, &azcosmos.ClientOptions{
		Routing: azcosmos.PreferredRegions(
			azcosmos.RegionWestUS,
			azcosmos.RegionEastUS,
		),
		Operation: azcosmos.OperationOptions{
			AvailabilityStrategy: azcosmos.HedgingAvailability(500 * time.Millisecond),
		},
	})
}
```

This example requires imports for `context`, `time`,
`github.com/Azure/azure-sdk-for-go/sdk/azcore`, and
`github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2`.

### PPCB configuration

Go currently supports environment-driven PPCB configuration. It does not expose the Rust driver's
`PartitionFailoverOptions` programmatically. The native driver captures these values when it is
constructed by either account constructor. Changing the process environment
afterward does not reconfigure an existing driver or client.

The pinned native driver v0.2.0 uses these defaults:

| Environment variable | Default | Purpose |
| --- | --- | --- |
| `AZURE_COSMOS_PPCB_ENABLED` | `true` | Enables the driver's PPCB setting. |
| `AZURE_COSMOS_PPCB_ENABLED_OVERRIDE` | unset | Authoritative incident-response kill switch. When set, it overrides both `AZURE_COSMOS_PPCB_ENABLED` and account-property enablement. |
| `AZURE_COSMOS_PPCB_READ_FAILURE_THRESHOLD` | `10` | Read failures before a partition/region circuit trips. |
| `AZURE_COSMOS_PPCB_WRITE_FAILURE_THRESHOLD` | `5` | Multi-master write failures before a partition/region circuit trips. |
| `AZURE_COSMOS_PPCB_COUNTER_RESET_WINDOW_MS` | `300000` | Window after which accumulated failure counters reset. |
| `AZURE_COSMOS_PPCB_PARTITION_UNAVAILABILITY_DURATION_MS` | `5000` | Minimum time a tripped entry remains unavailable before probing is eligible. |
| `AZURE_COSMOS_PPCB_FAILBACK_SWEEP_INTERVAL_MS` | `300000` | Interval between background failback sweeps. |
| `AZURE_COSMOS_PPCB_CONSECUTIVE_HEDGE_WIN_THRESHOLD` | `5` | Consecutive alternate-region hedge wins before PPCB trips the primary partition/region. |
| `AZURE_COSMOS_PARTITION_TOPOLOGY_CACHE_MODE` | `Eager` | Loads physical partition topology during container resolution; `Lazy` defers it until first use. |

Without `AZURE_COSMOS_PPCB_ENABLED_OVERRIDE`, PPCB is enabled when either the driver setting or the
account property enables it. The override is the authoritative operational safety valve in either
direction and should be set before driver construction; changing it for already-created drivers
requires restarting the process.

### Observing regional behavior

Successful responses and Cosmos errors can include `Diagnostics`. Use
`Diagnostics.RegionsContacted` for the distinct regions in first-contact order and
`Diagnostics.Attempts` for retained per-attempt region, endpoint, status/substatus, and latency.
These fields show the request's observable regional path, but diagnostics do not expose the
driver's full internal circuit-breaker state.

### Managing databases and containers

Use `Client.CreateDatabase`, `DatabaseClient.CreateContainer`, and the resource clients' `Read` and
`Delete` methods to manage databases and containers. Creates and deletes require a client created
with `NewClientWithKey`. Metadata reads also support Microsoft Entra ID through `NewClient` when
the identity has the `Microsoft.DocumentDB/databaseAccounts/readMetadata` data action.

Container creates require at least one partition key path. An omitted partition key kind is inferred
as `Hash` for one path and `MultiHash` for multiple paths. An omitted partition key version creates
the container with version 2. The definition is immutable after creation.

Database and container creates always request the response body so the returned properties include
the service-assigned resource ID, self-link, ETag, and last-modified time. For these operations,
`OperationOptions.EnableContentResponseOnWrite` is ignored. Deleting a database also deletes every
container and item in it; deleting a container deletes every item in it.

Provisioned-throughput creation and throughput read/replace APIs are not included yet. Configure
throughput separately until the follow-up throughput surface lands.

### SDK identity

Service requests include `azsdk-go-azcosmos/<version>` in the User-Agent header alongside the native
Cosmos driver identity and feature flags. Configure an application suffix through
`RuntimeOptions.ApplicationID`, then supply that runtime in `ClientOptions.Runtime`.
SDK identity does not consume the suffix's length allowance.

### Application models and serialization

`CreateItem`, `ReplaceItem`, and `UpsertItem` accept models and infer their type parameter:

```go
type Product struct {
    ID       string `json:"id"`
    Category string `json:"category"`
    Quantity int    `json:"quantity"`
}

item := Product{ID: "item-1", Category: "catalog", Quantity: 3}
response, err := container.CreateItem(ctx, azcosmos.NewPartitionKeyString("catalog"), item.ID, item, nil)
```

IDs and partition keys remain explicit; the SDK never infers or changes model fields.
`ReadItem`, content-returning writes, and `PatchItem` return `ItemResponse`. Decode it
with `response.IntoModel[Product]()`, or inspect `response.Value` directly.
Model conversion is local and leaves response metadata and bytes intact, including after
`Client.Close`. Delete and body-disabled writes have no model body; calling `IntoModel`
on an absent body returns an error. Check for a body before decoding an optional write response.

Serialization uses `encoding/json`, honoring struct tags and model-level
`json.Marshaler` / `json.Unmarshaler` hooks. Clients and responses retain no serializer state.
Model hooks are synchronous and cannot be interrupted by context cancellation. See the compiling
[examples](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2#pkg-examples)
for model-based operations and a model hook that preserves numbers in untyped projections.

Conversion failures return `*Error` with `CodeSerializationFailed` and unwrap their cause.
Response decoding errors retain available RU charge, ETag, session token, status, and
diagnostics while `FromWire` remains false because the failure is local. A write already
completed before `IntoModel` was called; a decoding error never retries or rolls it back.
No partial model or partial query page is returned on error.

#### Migrating provisional v2 callers

There are no separate `Typed` or `Raw` operation variants. Pass models directly instead of
calling `json.Marshal` before each write. For intentionally pre-encoded JSON, pass
`json.RawMessage(payload)`; plain `[]byte` now follows encoding/json's base64 encoding
behavior rather than acting as a passthrough. Select
`NewQueryItemsPager[json.RawMessage]` to obtain encoded query values.

Pass context first to both account constructors and both container factories, and remove calls
to `Client.Initialize`. Account and container factory errors now occur during construction;
`NewDatabase` remains local. `Client.Endpoint`, `DatabaseClient.ID`, and `ContainerClient.ID` remain.
Replace serializer injection with model JSON hooks or `json.RawMessage`.

Move `ClientOptions.BinaryEncoding` to `ClientOptions.Operation.BinaryEncoding`, replacing
`Enabled: &value` with `Disabled: !value`. A nil outer pointer inherits; a non-nil zero group
explicitly enables binary encoding. Precedence is request, client, runtime, environment, then
enabled. The Go binding reads `AZURE_COSMOS_BINARY_ENCODING_ENABLED` at runtime creation because
native 0.2.0 does not read it: `1`, `true`, `yes`, and `on` (case-insensitive, trimmed) enable it;
other set values disable it. The implicit runtime captures this once on first use. This does not
add binary support to PATCH, delete, or transactional batches. Responses remain text JSON.

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
`PatchItem` call. Supplying a stable tracking ID across separate calls is not exposed yet.

### Transactional batches

Build a `TransactionalBatch` with `NewTransactionalBatch(partitionKey)` and append `CreateItem`,
`ReadItem`, `UpsertItem`, `ReplaceItem`, or `DeleteItem` operations. Execute it with
`ContainerClient.ExecuteTransactionalBatch(ctx, batch, options)`. A nil options pointer selects
defaults. Transactional PATCH is not exposed.

The key must be complete: include every component of a hierarchical partition key. Unlike query
scopes, batches cannot use a prefix or span logical partitions. Execution reads the container's
partition key definition before native submission and rejects an incomplete key. This metadata
validation can consume RUs; its charge is included in the returned response or error.
An error can contain setup metadata before the batch itself is submitted; these fields do not
establish that any writes were admitted.
The driver/service checks application-defined item IDs and partition key fields; Go does not
parse these fields or implement a separate transaction transport.

Append methods copy item JSON and conditional options. Extending a copied builder does not modify
the original. Any append error makes that builder unusable, preventing an ignored error from
silently submitting only part of a transaction. A batch must contain 1-100 operations; the encoded
JSON operation envelope sent to the native driver must be at most 2 MiB. Reusing a batch executes
a new transaction. A finished builder can be executed concurrently, but must not be mutated
concurrently.

Read, upsert, replace, and delete accept If-Match; read and upsert also accept If-None-Match.
The two conditions are mutually exclusive. Go forwards these options using the same operation
envelope as the Rust SDK; it does not independently enforce them or rewrite the returned statuses.

**Live-account conditional behavior with native driver v0.2.0:**

| Batch option | Observed result |
| --- | --- |
| Read `IfMatchETag` | A matching ETag returned HTTP 304; a different ETag returned HTTP 200. Accompanying writes committed in both cases. |
| Read `IfNoneMatchETag` | Matching and different ETags both returned HTTP 200. A matching point-read control returned HTTP 304. |
| Upsert `IfNoneMatchETag` | Existing items were replaced, including when the condition was `"*"`. Accompanying writes committed. |

Do not use these batch-read conditions to guard other writes, or upsert If-None-Match to require
creation. Use `TransactionalBatch.CreateItem` for creation-only behavior; an existing item then
causes a conflict and rolls back the batch. These live observations do not establish intended
service guarantees or a Rust-specific defect. Emulator conditional-read behavior differs and
does not establish the live-service contract.

`TransactionalBatchOptions.Operation.EnableContentResponseOnWrite` is request-wide: nil inherits
runtime/client/native defaults, false suppresses write bodies, and true requests them. Native
driver v0.2.0 defaults to false. **Live-account content behavior:** explicit or inherited false
can also suppress HTTP 200 read bodies. Set this option explicitly to true when read bodies
are required; this also requests write bodies. Go preserves absent content rather than fetching
it with a separate request. Emulator read-content behavior differs and does not establish the
live-service contract. `SessionToken` and other shared operation settings are forwarded to the driver.

A received response can contain non-2xx operation statuses without an execution error.
Like the Rust SDK, Go does not promote an individual operation's failure to a batch-level error.
Neither `err == nil` nor a 2xx outer status means every operation succeeded.
`OperationResults` is ordered exactly like the builder and preserves statuses, substatuses,
charges, raw JSON, ETags, and retry delays. Batch-level headers and the raw result envelope are
also retained. HTTP 304 reports a not-modified read; HTTP 424 reports a dependency failure.
HTTP 207 alone does not establish rollback, and mixed results such as `[201, 304]` are preserved
without inferring whether writes committed. Malformed result envelopes return an explicit
decoding error with the native body and metadata.

Context cancellation stops waiting for a batch, not native execution. An admitted batch may still
commit after the method returns a cancellation error. Allow a context and native end-to-end timeout
budget appropriate for the complete batch, including metadata setup and retries; short budgets
can leave the caller without a result for submitted writes. A longer budget reduces avoidable
timeouts but does not guarantee a known commit outcome. Do not blindly replay a cancelled batch:
each execution is a new transaction.

An already available native completion is returned unchanged; otherwise the method returns the
zero response and a cancellation error while native work drains in the background. `Client.Close`
waits for outstanding native work before freeing resources. Responses and errors own their
payloads and remain usable after the client closes.

### Querying items

`ContainerClient.NewQueryItemsPager[T]` accepts a `Query`, an explicit `FeedScope`, and optional
`QueryOptions`. Use `NewQuery(sql).WithParameter(name, value)` to capture JSON parameter values,
and `NewFeedScopeForPartitionKey(pk)` to target a logical partition or hierarchical key prefix.
Use `NewFeedScopeForFullContainer()` for cross-partition queries. Query values are immutable; `WithParameter` returns a
new query and an error if the value cannot be serialized.

The API follows the Rust SDK's separation of query, scope, and feed options. Its Cosmos-owned
`QueryItemsPager[T]` retains the native query plan and provides `More`, `NextPage`, and `Close`.
Use `defer pager.Close()` to release native resources on early exit; exhaustion also releases
them, and `Client.Close` reclaims any remaining cursors. This unreleased-v2 API intentionally
differs from `runtime.Pager[T]` and requires Go SDK API review before stabilization.
`QueryOptions.Feed.PageSizeHint` is a positive per-partition page-size hint; zero leaves
sizing to the driver, and negative hints are rejected. Results are decoded values in
`QueryItemsResponse[T].Items`. Select a document or projection type, or a scalar for
`SELECT VALUE`. Pointer types preserve JSON null as nil with encoding/json;
other null conversions follow encoding/json semantics. An empty page does not
necessarily end the query: use `More`, not the number of items.

Call `pager.ContinuationToken(ctx)` after a page, before exhaustion or Close, and pass the result
in `QueryOptions.Feed.ContinuationToken` to resume with
the same query and scope. This is an opaque **driver planner token**, not the service's
`x-ms-continuation` header. Do not parse or modify it. Snapshotting does not advance the pager.
Some buffered query plans do not support checkpoints but can still be fully iterated.
Checkpoint availability does not determine whether more pages exist. A final empty page may
mark exhaustion; check `More()` again after fetching it.
Each new pager resolves its container reference through the native driver's cache, allowing queries
started after container recreation to use the replacement. Tokens from the deleted container
cannot be used to resume against the replacement.
An execution or model-decoding failure terminates and closes the pager; the native driver
owns execution retries. A failed page delivers no partial results and cannot be checkpointed.
Resume from a checkpoint saved before the failure if needed. An unsupported checkpoint does
not terminate the pager.

Pager construction performs no network I/O and snapshots its inputs. Each pager's first fetch
(including a resumed pager) performs an additional container-metadata read to verify that the
key or prefix fits the container definition. Full-container scopes omit this additional read.
Its request charge is included in the first fetch's response or error when the metadata
completion is available before that fetch returns. If cancellation returns first, charges
reported by a later native completion cannot be surfaced to the caller; a zero reported
charge does not mean the operation consumed no RUs.
The metadata read shares the page's context and end-to-end timeout. This requires permission to read container
metadata. Iteration and checkpoint calls must not be concurrent; independent pagers can share a
client. `Close` safely synchronizes with active calls.
All cursors on a client share a cursor-format completion reactor, separate from the legacy
operation queue. Query callers wait on Go channels rather than each blocking an OS thread.

The pinned native driver v0.2.0 returns complete raw and pre-split query pages.
`FeedOptions.MaxFanOut` limits initial physical-partition fan-out; zero uses the native default
of 100. This is not a runtime limit: resume and later splits do not recheck it. Broad scopes
can consume substantial RUs. `QueryPlanMode` selects local-preferred (default) or gateway-only
planning. Optional `PopulateIndexMetrics` and `PopulateQueryMetrics` flags request the
corresponding response fields; index metrics are already decoded JSON.

The published ABI does not expose arbitrary physical/EPK ranges or Rust's buffered-window
configuration. Buffered query shapes retain the native default window (1000), require finite
TOP/LIMIT where the engine requires it, and may not support continuation snapshots.
Cross-partition aggregates such as COUNT and SUM are rejected by this published engine;
the Go binding preserves that error rather than computing aggregates itself.

Contexts bound query waits; canceling an admitted page or checkpoint terminates that pager
without claiming to cancel native work. Late completions are drained safely, and client
shutdown waits for their cleanup. Native per-page timeouts retain the configured duration
(the driver clamps explicit values below one second, including zero); subsequent page contexts
do not inherit the first page's remaining deadline. With `EndToEndTimeout == nil`, timeout
settings are inherited, and a context deadline limits only the Go wait.
Per-page context deadlines cannot bound native execution with the pinned ABI. Stronger
guarantees require cooperative cancellation
([Azure/azure-sdk-for-rust#5358](https://github.com/Azure/azure-sdk-for-rust/issues/5358))
or per-page deadline support in a future native release.
Item operations also stop waiting when their context ends. Submitted writes may still commit;
a cancellation error is not evidence that a write was rolled back or never reached the service.

#### Rust query-option parity

Against the Rust revision used to publish native v0.2.0, query-option coverage is:

| Rust setting | Go support |
| --- | --- |
| Positive page-size hint, continuation token, maximum fan-out | `FeedOptions.PageSizeHint`, `ContinuationToken`, `MaxFanOut` |
| Query-plan mode, session token, index/query metric flags | Corresponding `QueryOptions` fields |
| Buffered-query window limit | Not exposed by the native ABI; native default of 1000 applies |
| Explicit `MaxItemCountHint::ServerDecides` | Not independently expressible through this ABI; Go zero selects the native unset/default behavior |
| Common per-operation settings | `OperationOptions` exposes consistency, excluded regions, content response on writes, end-to-end timeout, PATCH strategy, throughput bucket/priority, throttle/session/failover retry controls, endpoint-unavailability TTL, session-capture control, hedging/availability strategy, custom headers, and binary encoding |

Common operation settings inherit across runtime, client, and request scopes; their applicability
depends on the operation. Binary wire encoding is enabled by default, with native conversion
returning text JSON to Go. A nil `ExcludedRegions` slice inherits exclusions, while a non-nil
empty slice explicitly clears them. The buffered-window and explicit server-decides gaps
still require a native contract change.

### Running the end-to-end tests

The tests in `emulator_test.go`, `diagnostics_emulator_test.go`, and `retry_native_test.go` run
real operations against a service. They need a driver-backed build and the `EMULATOR` environment
variable, and they skip otherwise.

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

EMULATOR=1 AZCOSMOS_ENDPOINT=http://127.0.0.1:49151/ go test -run 'Test(Emulator|NativeThrottle|NativeStaleSession)' ./...
```

The container the tests use is declared in `internal/testdata/emulator-config.json`; its ids
default to `itemdb` and `items` and can be overridden with `AZCOSMOS_DATABASE` and
`AZCOSMOS_CONTAINER`. Query scope tests also use the `query-hierarchical` container declared in
the same configuration, under the selected database.

Query regression tests cover ordered cross-partition checkpoint/resume on a fresh client,
mid-pagination errors and request-charge accounting, and metric/session option propagation.
The contract tests use a local HTTP proxy to inject known error and metric response headers
while the real native driver executes against the emulator. This validates the binding and
native decoding, not live-service metric generation. Per-cursor test hooks control checkpoint
delivery and cancellation races without relying on network timing.

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
