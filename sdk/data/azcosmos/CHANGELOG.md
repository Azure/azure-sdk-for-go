# Release History

## 1.5.0-beta.7 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

* Cross-region failover in the Cosmos retry policy no longer blocks on a synchronous `gem.Update(ctx, true)`. Connection-error failover drops the forced refresh and leaves `retryCount` at 0 (unavailable endpoints are demoted to the tail of the route list, so bumping the index routed back to them), and normalizes `MarkEndpointUnavailable*` URLs to scheme+host so marks are actually visible to the cache. 403/`WriteForbidden` retries refresh fire-and-forget (CAS-gated); 408 retries no longer mark unavailable or refresh. See [PR 26889](https://github.com/Azure/azure-sdk-for-go/pull/26889).
* Connection-error retry policy now retries up to 3 times in-region before one cross-region failover. Writes only fail over for errors that prove the request never reached the service (DNS, dial, TLS, `ECONNREFUSED`); ambiguous transport failures (`ECONNRESET`, `EOF`, transport timeouts) are no longer retried for writes to avoid duplicates. Caller context deadlines short-circuit the policy. See [PR 26858](https://github.com/Azure/azure-sdk-for-go/pull/26858).
* HTTP `408 Request Timeout` is now handled: reads retry once cross-region, writes return immediately to avoid duplicates. See [PR 26858](https://github.com/Azure/azure-sdk-for-go/pull/26858).
* Fixed excessive `GetDatabaseAccount` calls with preferred regions, and stopped data-plane retries from trailing into the default endpoint once topology is populated. See [PR 26815](https://github.com/Azure/azure-sdk-for-go/pull/26815).
* Partition key range cache now serves concurrent callers from a single in-flight refresh per container, keeps the cached routing map readable during refresh, and runs the refresh on `context.Background()` so one caller's cancellation cannot abort the shared fetch. See [PR 26855](https://github.com/Azure/azure-sdk-for-go/pull/26855).
* Partition key range cache change-feed pagination is now resilient to mid-drain 429s (retried indefinitely with capped backoff + jitter) and preserves accumulated pages across refreshes. See [PR 26855](https://github.com/Azure/azure-sdk-for-go/pull/26855).

### Other Changes

* Tightened the default HTTP client: 5s dial timeout (was 30s), 65s wall-clock cap per HTTP attempt (was unbounded), larger idle connection pool (1000/100, was 100/10), and faster HTTP/2 health checks. Caller-supplied `Transport` and shorter `context` deadlines are unaffected. See [PR 26856](https://github.com/Azure/azure-sdk-for-go/pull/26856).

## 1.5.0-beta.6 (2026-05-15)

### Features Added

* Adds `PriorityLevel` and `ThroughputBucket` options at the client and per-request level for item, query, change-feed, batch, and read-many operations. See [PR 26750](https://github.com/Azure/azure-sdk-for-go/pull/26750)
* Added client-level partition key range cache and container properties cache, reducing redundant metadata round-trips for ReadMany and query operations. See [PR 26723](https://github.com/Azure/azure-sdk-for-go/pull/26723)
* Added operation diagnostics on responses and `DiagnosticsFromError` for retrieving diagnostics from failed operations. See [PR 26548](https://github.com/Azure/azure-sdk-for-go/pull/26548)

### Breaking Changes

* Removed `ChangeFeedResponse.PopulateCompositeContinuationToken()`. The method is no longer needed: `GetChangeFeed` now populates `ChangeFeedResponse.ContinuationToken` directly with the multi-range composite token. Callers who built single-range tokens manually can use `GetCompositeContinuationToken()` instead. See [PR 26792](https://github.com/Azure/azure-sdk-for-go/pull/26792).

### Bugs Fixed

* Fixed `GetChangeFeed` to survive partition splits: customer-supplied `FeedRange`s are now overlap-matched against the routing map, `410/Gone` triggers a cache refresh and bounded retry, split parents expand into per-child queue entries (inheriting the parent's ETag), and the continuation token persists multi-range state across calls. Continuation tokens are guarded against cross-container reuse. See [PR 26768](https://github.com/Azure/azure-sdk-for-go/pull/26768).
* Fixed V2 partition key routing: the top 2 bits of the first EPK byte are now masked to stay within the partition key range space [0x00, 0x3F]. Previously, items whose V2 hash started with a byte >= 0x40 could fail routing in ReadMany because the EPK lexicographically exceeded the "FF" range sentinel. See [PR 26723](https://github.com/Azure/azure-sdk-for-go/pull/26723)
* Fixed error handling for partition key range calls which would previously cause panics on any error. See [PR 26723](https://github.com/Azure/azure-sdk-for-go/pull/26723)
* Fixed partition key range cache to use change-feed pagination when fetching ranges, preventing incomplete range sets on containers with many partitions. The incremental refresh path now accumulates all pages before merging, correctly handling cascading splits across multiple change-feed pages. See [PR 26777](https://github.com/Azure/azure-sdk-for-go/pull/26777)

## 1.5.0-beta.5 (2026-03-09)

### Features Added

* Adds support for float 16 datatype for vector embedding policy. See [PR 25707](https://github.com/Azure/azure-sdk-for-go/pull/25707)
* Improved the performance of the built-in ReadMany implementation. See [PR 26007](https://github.com/Azure/azure-sdk-for-go/pull/26007)

### Breaking Changes

* Removed `QueryEngine` field from `ReadManyOptions`. ReadMany now always uses the built-in Go-native implementation.

### Other Changes

* Small performance optimizations to API's using query engine. See [PR 25669](https://github.com/Azure/azure-sdk-for-go/pull/25669)

## 1.5.0-beta.4 (2025-11-24)

### Features Added

* Added client engine support for `ReadManyItems`. See [PR 25458](https://github.com/Azure/azure-sdk-for-go/pull/25458)

## 1.5.0-beta.3 (2025-11-10)

### Features Added

* Adjusted the query engine abstraction to support future enhancements and optimizations. See [PR 25503](https://github.com/Azure/azure-sdk-for-go/pull/25503)

## 1.5.0-beta.2 (2025-11-03)

### Features Added

* Added `ReadManyItems` API to read documents across partitions. See [PR 25522](https://github.com/Azure/azure-sdk-for-go/pull/25522)

## 1.5.0-beta.1 (2025-10-16)

### Features Added

* Added support for BypassIntegratedCache option See [PR 24772](https://github.com/Azure/azure-sdk-for-go/pull/24772)
* Added support for specifying Full-Text Search indexing policies when creating a container. See [PR 24833](https://github.com/Azure/azure-sdk-for-go/pull/24833)
* Added support for specifying Vector Search indexing policies when creating a container. See [PR 24833](https://github.com/Azure/azure-sdk-for-go/pull/24833)
* Added support for reading Feed Ranges from a container. See [PR 24889](https://github.com/Azure/azure-sdk-for-go/pull/24889)
* Added support for reading Change Feed through Feed Ranges from a container. See [PR 24898](https://github.com/Azure/azure-sdk-for-go/pull/24898)
* Additional logging in the query engine integration code. See [PR 25444](https://github.com/Azure/azure-sdk-for-go/pull/25444)

## 1.4.1 (2025-08-27)

### Bugs Fixed

* Fixed bug where the correct header was not being sent for writes on multiple write region accounts. See [PR 25127](https://github.com/Azure/azure-sdk-for-go/pull/25127)

## 1.5.0-beta.0 (2025-06-09)

### Features Added

* Added an initial API for integrating an external client-side Query Engine with the Cosmos DB Go SDK. This API is unstable and not recommended for production use. See [PR 24273](https://github.com/Azure/azure-sdk-for-go/pull/24273) for more details.

## 1.4.0 (2025-04-29)

### Other Changes

* Requests to update region topology (often made automatically as part of other operations) now pass through the same Context as the request that triggered them. This allows for flowing telemetry spans and other Context values through HTTP pipeline policies. However, these requests do NOT use the cancellation signal provided in the original request Context, in order to ensure the region topology is properly updated even if the original request is cancelled. See [PR 24351](https://github.com/Azure/azure-sdk-for-go/issues/24351) for more details.

## 1.3.0 (2025-02-12)

### Features Added

* Added limited support for cross-partition queries that can be served by the gateway. See [PR 23926](https://github.com/Azure/azure-sdk-for-go/pull/23926) and this [querying with Cosmos document](https://learn.microsoft.com/rest/api/cosmos-db/querying-cosmosdb-resources-using-the-rest-api#queries-that-cannot-be-served-by-gateway) for more details.

### Other Changes

* All queries now set the `x-ms-documentdb-query-enablecrosspartition` header. This should not impact single-partition queries, but in the event that it does cause problems for you, this behavior can be disabled by setting the `EnableCrossPartitionQuery` value on `azcosmos.QueryOptions` to `false`.

## 1.2.0 (2024-11-12)

### Features Added

* Added API for creating Hierarchical PartitionKeys. See [PR 23577](https://github.com/Azure/azure-sdk-for-go/pull/23577)
* Set all Telemetry spans to have the Kind of SpanKindClient. See [PR 23618](https://github.com/Azure/azure-sdk-for-go/pull/23618)
* Set request_charge and status_code on all trace spans. See [PR 23652](https://github.com/Azure/azure-sdk-for-go/pull/23652)

### Bugs Fixed

* Pager Telemetry spans are now more consistent with the rest of the spans. See [PR 23658](https://github.com/Azure/azure-sdk-for-go/pull/23658)

## 1.1.0 (2024-09-10)

### Features Added

* Added support for OpenTelemetry trace spans. See [PR 23268](https://github.com/Azure/azure-sdk-for-go/pull/23268)
* Added support for MaxIntegratedCacheStaleness option See [PR 23406](https://github.com/Azure/azure-sdk-for-go/pull/23406)

### Bugs Fixed

* Fixed sending `Prefer` header with `return=minimal` value on metadata operations. See [PR 23335](https://github.com/Azure/azure-sdk-for-go/pull/23335)
* Fixed routing metadata requests to satellite regions when using ClientOptions.PreferredRegions and multiple write region accounts. See [PR 23339](https://github.com/Azure/azure-sdk-for-go/pull/23339)

## 1.0.3 (2024-06-17)

### Bugs Fixed

* Fixed data race on clientRetryPolicy. See [PR 23061](https://github.com/Azure/azure-sdk-for-go/pull/23061)

## 1.0.2 (2024-06-11)

### Bugs Fixed

* Fixed ReplaceThroughput operations on Database and Container. See [PR 22923](https://github.com/Azure/azure-sdk-for-go/pull/22923)

## 1.0.1 (2024-05-02)

### Bugs Fixed

* Reduces minimum required go version to 1.21

## 1.0.0 (2024-04-09)

### Features Added

* Added regional routing support through ClientOptions.PreferredRegions
* Added availability logic and failover mechanics to support cross-regional retries and resiliency enhancements
* Added extended logging for requests, responses, and client configuration

### Breaking Changes

* ItemOptions.SessionToken, QueryOptions.SessionToken, QueryOptions.ContinuationToken, QueryDatabasesOptions.ContinuationToken, QueryContainersOptions.ContinuationToken are now `*string`
* ItemResponse.SessionToken, QueryItemsResponse.ContinuationToken, QueryContainersResponse.ContinuationToken, QueryDatabasesResponse.ContinuationToken are now `*string`

## 0.3.6 (2023-08-18)

### Bugs Fixed

* Fixed PatchItem function to respect EnableContentResponseOnWrite

## 0.3.5 (2023-05-09)

### Features Added

* Added support for accounts with [merge support](https://aka.ms/cosmosdbsdksupportformerge) enabled

### Bugs Fixed

* Fixed unmarshalling error when using projections in value queries

## 0.3.4 (2023-04-11)

### Features Added

* Added `NullPartitionKey` variable to create and query documents with null partition key in CosmosDB

## 0.3.3 (2023-01-10)

### Features Added

* Added `PatchItem` function to patch documents
* Added support for querying databases and containers

## 0.3.2 (2022-08-09)

### Features Added

* Added `NewClientFromConnectionString` function to create client from connection string
* Added support for parametrized queries through `QueryOptions.QueryParameters`

### Bugs Fixed

* Fixed handling of ids with whitespaces and special supported characters

## 0.3.1 (2022-05-12)

### Features Added

* Added Transactional Batch support

### Other Changes

* Update to latest `azcore` and `internal` modules

## 0.3.0 (2022-05-10)

### Features Added

* Added single partition query support.
* Added Azure AD authentication support through `azcosmos.NewClient`

### Breaking Changes

* This module now requires Go 1.18

## 0.2.0 (2022-01-13)

### Features Added

* Failed API calls will now return an `*azcore.ResponseError` type.

### Breaking Changes

* Updated to latest `azcore`. Public surface area is unchanged.  However, the `azcore.HTTPResponse` interface has been removed.

## 0.1.0 (2021-11-09)

* This is the initial preview release of the `azcosmos` library
