<!-- cspell:ignore azcosmos azcore azidentity cgo PPAF PPCB HPK RNTBD upsert -->
# Cosmos DB Go v1 to v2 public API surface

> **Status:** API direction established; diagnostics and `RawResponse`
> recommendations remain to be ratified.  
> **Audience:** Cosmos SDK engineers, leads, and Go Central SDK reviewers.  
> **Scope:** Public Go API shape only. Delivery sequencing and native-binary
> distribution are covered by their respective planning documents.

## 1. Purpose and framing

Go v2 uses the Rust Cosmos driver through FFI, but it should still look and feel
like an Azure SDK for Go. The Rust SDK is the reference for capabilities and
behavior; it is not a template to copy literally into Go.

The API surface is organized into three categories:

| Category                               | What belongs here                                                              | Default approach                                                         |
| -------------------------------------- | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| **Brand-new Go v2 APIs**               | Rust-driver capabilities that Go v1 does not expose                            | Design an idiomatic Go surface around the driver capability              |
| **Go v1 APIs with real shape changes** | Existing APIs that cannot remain unchanged because of FFI or expanded behavior | Preserve familiar Go patterns while making the changed contract explicit |
| **Small naming or cosmetic changes**   | Existing behavior where only names or grouping might change                    | Keep the v1 shape unless the change has clear customer value             |

### Sources

The comparison is grounded in:

- the exported Go v1 API under
  `azure-sdk-for-go/sdk/data/azcosmos`;
- the public Rust SDK under
  `azure-sdk-for-rust/sdk/cosmos/azure_data_cosmos`;
- the C ABI in `azure_data_cosmos_driver_native`;
- the driver specifications for availability, routing, Gateway 2.0, change
  feed, patch, diagnostics, and distributed transactions.

## 2. Go conventions that carry forward

These choices apply across the full API and keep Go v2 familiar to existing
customers.

### Client construction

Keep the three Go v1 constructors:

```go
NewClient(endpoint, tokenCredential, options)
NewClientWithKey(endpoint, keyCredential, options)
NewClientFromConnectionString(connectionString, options)
```

The Rust SDK uses a unified credential abstraction internally. Go does not need
to expose that abstraction. Key conversion, token-refresh callbacks, and FFI
marshalling remain behind the existing constructors.

### Paging and continuation

Keep `runtime.Pager[T]` and the established continuation-token model:

```go
pager := container.NewQueryItemsPager(query, partitionKey, options)
for pager.More() {
    page, err := pager.NextPage(ctx)
    // page.ContinuationToken can be persisted and used to resume later.
}
```

Continuation tokens remain:

- on page responses, so callers can save them; and
- in query or feed options, so callers can resume.

The Rust iterator's live `to_continuation_token()` method should not be copied
into Go. An opaque token value is simpler, already resumable across process
boundaries, and consistent with the Azure SDK for Go.

### Options

Keep per-operation `*XxxOptions` structs with `nil` meaning defaults. Do not use
Rust-style builders as the primary API.

A builder or functional-options convenience layer can be added later as an
additive feature after the options structs stabilize.

## 3. Brand-new Go v2 APIs

These capabilities do not have a Go v1 surface to preserve.

| Capability                       | Proposed Go v2 surface                                                             | Default and availability                         |
| -------------------------------- | ---------------------------------------------------------------------------------- | ------------------------------------------------ |
| Cross-region request hedging     | `AvailabilityStrategy` on client options with a per-operation override             | Disabled or driver-default at the zero value; GA |
| Per-partition automatic failover | Client-level enablement or tuning only when the behavior is not fully automatic    | Match the driver default; GA                     |
| Per-partition circuit breaker    | Driver-aligned tuning and environment-variable support                             | Enabled by default; GA                           |
| Gateway 2.0                      | Automatic selection with a narrow disable or transport escape hatch only if needed | Automatic for eligible accounts; GA              |
| Throughput-control groups        | Client-level group registration plus a per-operation group selector                | Opt-in; GA                                       |
| Local query planning             | No public API                                                                      | Internal and enabled by default; GA              |
| User-Agent feature token         | No new customer-facing type; retain the existing user-agent suffix option          | Internal feature advertising; GA                 |
| Distributed transactions         | Preview-only API surface                                                           | Excluded from the initial GA API                 |

### Availability strategy

Cross-region request hedging reduces tail latency. The Go surface should allow a
client-wide default and a per-operation override:

```go
type AvailabilityStrategy struct {
    // Hedging mode and thresholds.
    // The zero value follows the driver default.
}
```

### Per-partition resilience

Per-partition automatic failover and the per-partition circuit breaker should
remain primarily driver-owned behaviors. Public options should expose only
customer-relevant enablement or tuning.

The circuit breaker is enabled by default in the driver. Go v2 must preserve
that default and honor the established environment variables so operational
behavior remains consistent across SDKs.

### Gateway 2.0

Gateway 2.0 should be selected automatically for eligible accounts. A public
option is justified only as a narrow escape hatch; customers should not need to
understand transport internals for normal use.

### Throughput-control groups

Go v1 already exposes request priority. Go v2 adds throughput-group
registration and a per-operation group selector:

```go
type ThroughputControlGroup struct {
    // Group name and target throughput configuration.
}

// Register a group on the client, then select it in operation options.
```

### Local query planning

Do not expose a public toggle. Local planning is an internal driver
optimization that avoids a gateway round trip. The Rust public SDK does not
expose this setting, and Go should not add one without a concrete customer need.

### Distributed transactions

Do not include distributed transactions in the initial generally available
surface. The driver marks the feature as work in progress and gates it behind a
disabled-by-default preview feature. Go should mirror that through a preview
mechanism until the driver and service behavior stabilize.

## 4. Go v1 APIs with real shape changes

| Area                                      | Go v1 shape                                                  | Go v2 direction                                                                                        | Migration impact                            |
| ----------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------------- |
| Error handling                            | `azcore.ResponseError`; Cosmos details read from raw headers | Add typed `CosmosError`, preserve `ResponseError` compatibility, and separate native-boundary failures | Header-parsing callers move to typed fields |
| Change feed                               | Timestamp-only `StartFrom *time.Time`                        | Add typed start positions and `ChangeFeedMode` inside `ChangeFeedOptions`                              | Timestamp callers update one field          |
| Response metadata                         | Typed fields plus `RawResponse *http.Response`               | Keep typed metadata; recommend removing `RawResponse`                                                  | Direct `RawResponse` users migrate          |
| Diagnostics                               | Existing accessors and `DiagnosticsFromError`                | Recommend extending the existing type with verbosity and summary data                                  | Additive if accepted                        |
| Authentication                            | Go credential constructors                                   | Keep the public entry points; marshal token refresh through FFI                                        | No call-site change                         |
| Item, read-many, patch, and batch results | Metadata sourced from HTTP responses                         | Preserve Go result shapes; populate fields from structured driver results                              | Mostly unchanged; verify field parity       |
| Hierarchical partition keys               | Typed constructors and builder                               | Keep both; support driver-backed multi-hash and prefix behavior                                        | Existing callers remain valid               |

### Error handling

#### What the native boundary returns

The C ABI returns a structured error rather than an `http.Response`. It
provides:

- status and substatus;
- activity ID;
- session token;
- ETag;
- retry-after duration;
- message and backtrace;
- whether the error came from the service;
- a native error code that distinguishes service errors from argument,
  authentication, FFI, and fatal failures.

Request charge is response metadata and is not part of the error.

#### Go v2 shape

Introduce a typed `CosmosError` with accessors for the structured fields:

```go
var cosmosError *cosmos.CosmosError
if errors.As(err, &cosmosError) {
    // Use typed status, substatus, activity ID, retry-after, and diagnostics.
}
```

Preserve the established Go error pattern:

```go
var responseError *azcore.ResponseError
if errors.As(err, &responseError) {
    // Existing status-code handling continues to work.
}
```

`CosmosError` should implement `Unwrap` or `As` so existing
`azcore.ResponseError` checks continue to succeed. Errors that originate in
argument validation, authentication conversion, FFI plumbing, or another
native-boundary failure should use a separate error kind. This makes "Cosmos
returned 404" distinguishable from "the native call failed."

#### Alternatives considered

| Approach                                                               | Customer effect                                              | Assessment                                                                   |
| ---------------------------------------------------------------------- | ------------------------------------------------------------ | ---------------------------------------------------------------------------- |
| Return only `azcore.ResponseError` and synthesize an HTTP response     | Maximum source compatibility                                 | Recreates headers from typed data and presents a response that never existed |
| Add typed `CosmosError` while preserving `ResponseError` compatibility | Existing checks work; new typed fields remove header parsing | **Selected direction**                                                       |
| Replace `ResponseError` entirely                                       | Cleanest new type                                            | Unnecessary break at every error-handling call site                          |

### Change feed

Go v1's `StartFrom *time.Time` can express only a timestamp. The driver supports
explicit positions for now, beginning, timestamp, and continuation.

Go v2 should add:

```go
type ChangeFeedStartFrom struct {
    // Represents now, beginning, timestamp, or continuation.
}

type ChangeFeedMode int
```

Both remain inside `ChangeFeedOptions`; start position should not become a
required positional argument. The pager and continuation-string model stay
unchanged.

Initial mode support is `LatestVersion`, with room to add
`AllVersionsAndDeletes` when the capability is ready.

Existing timestamp callers migrate from:

```go
ChangeFeedOptions{StartFrom: &startTime}
```

to the timestamp form of `ChangeFeedStartFrom`.

### Response metadata

The following fields map directly from the driver and remain typed:

- request charge;
- activity ID;
- ETag;
- session token;
- diagnostics.

The exception is `Response.RawResponse *http.Response`. Go v2 does not execute
the service request through Go's HTTP stack, and the FFI contract does not
return an HTTP response.

#### Recommendation: remove `RawResponse`

Removing the field is the most accurate v2 contract:

- keeping it as always `nil` creates a runtime trap;
- synthesizing it fabricates a response that never existed;
- typed response and error fields already expose the useful metadata.

This is a source-breaking change for callers that directly access
`RawResponse`, but a major-version boundary is the right time to remove a field
that cannot be populated honestly.

If source compatibility is later judged more important than API accuracy, the
fallback is to retain the field and document it as always `nil`. Synthesizing an
HTTP response is not recommended.

Apply the same metadata mapping consistently to item, read-many, patch, and
transactional-batch results.

### Diagnostics

Go v1 exposes:

- `Diagnostics.String()`;
- `Diagnostics.ClientElapsedTime()`;
- `Diagnostics.StartTimeUTC()`;
- `Diagnostics.FailedRequestCount()`;
- `DiagnosticsFromError(err)`.

The driver adds verbosity control and richer summary data.

#### Recommendation: extend the existing type

Keep every v1 accessor and add:

- diagnostics verbosity in client and operation options;
- a summary-JSON or structured accessor;
- support for retrieving diagnostics from `CosmosError`.

This keeps existing debugging code working and makes the new surface additive.
Replacing the type with a driver-shaped structure would break diagnostics
callers for limited customer benefit. A larger redesign can wait for a future
major version if customer usage demonstrates a need.

### Authentication

Keep the public entry point:

```go
NewClient(endpoint, azcore.TokenCredential, options)
```

Token-refresh callbacks and credential marshalling require careful FFI design
and testing, but they are implementation concerns and do not require a
different public Go constructor.

## 5. Small naming or cosmetic changes

Only the module and package names change. The remaining v1 names stay because a
rename would add migration cost without changing behavior.

| Go v1 surface                                                                            | Go v2 direction                     | Outcome    |
| ---------------------------------------------------------------------------------------- | ----------------------------------- | ---------- |
| Module path ending in `/azcosmos`                                                        | Module path ending in `/cosmos/`    | **Rename** |
| Package name `azcosmos`                                                                  | Package name `cosmos`               | **Rename** |
| `ItemOptions`                                                                            | Keep one shared item-options type   | **Keep**   |
| `NewPartitionKeyString`, `NewPartitionKeyNumber`, `NewPartitionKeyBool`, and the builder | Keep typed constructors and builder | **Keep**   |
| Enum-listing helpers such as `ConsistencyLevelValues()`                                  | Keep existing helpers               | **Keep**   |
| `QueryParameter{Name, Value}`                                                            | Keep existing shape                 | **Keep**   |
| `NewManualThroughputProperties` and `NewAutoscaleThroughputProperties`                   | Keep existing names                 | **Keep**   |
| `CreateContainerOptions` and `ReplaceContainerOptions`                                   | Keep existing names                 | **Keep**   |
| `DiagnosticsFromError`                                                                   | Keep existing name                  | **Keep**   |
| `NullPartitionKey` and `NewPartitionKey()`                                               | Keep existing names                 | **Keep**   |

The new module is expected under:

```text
github.com/Azure/azure-sdk-for-go/sdk/data/cosmos/
```

with package name `cosmos`.

## 6. Consolidated migration impact

| Migration type        | What changes                                                                | Customer action                                              |
| --------------------- | --------------------------------------------------------------------------- | ------------------------------------------------------------ |
| Mechanical            | Module and package rename                                                   | Update the import path and package qualifier                 |
| Additive              | Typed `CosmosError`, richer diagnostics, new driver capabilities            | Adopt new fields and options when useful                     |
| Small code change     | Change-feed start position                                                  | Replace the timestamp pointer with the typed timestamp form  |
| Real code change      | Raw-header parsing and direct `RawResponse` access                          | Use typed response or `CosmosError` fields                   |
| Behavioral validation | Retry, session, routing, and resilience behavior moves into the Rust driver | Validate scenarios during implementation and release testing |
| Build and deployment  | Go v2 uses cgo and native driver binaries                                   | Use a supported platform and distributed binary              |

Most v1 applications should need only the import and package rename. Additional
work is limited to callers that parse raw headers, access `RawResponse`, or use
the timestamp-only change-feed field.

## 7. Recommended final API position

Proceed with the direction in this document:

- preserve the Go v1 constructors, pager, continuation-token, and options
  patterns;
- add public API only where the Rust driver introduces a customer-visible
  capability;
- use typed errors and response metadata instead of fabricating HTTP objects;
- keep local query planning internal;
- keep distributed transactions preview-only;
- rename the module and package to `cosmos`, while retaining the rest of the v1
  naming surface.

Before freezing the API, ratify these two recommendations:

1. **Diagnostics:** preserve and extend the v1 type.
2. **Responses:** remove `RawResponse` and direct customers to typed metadata.

No other API-shape choice in this document remains open.
