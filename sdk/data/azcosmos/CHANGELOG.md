# Release History

## 1.5.0-beta.3 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 1.5.0-beta.2 (2025-11-11)

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

* Added limited support for cross-partition queries that can be served by the gateway. See [PR 23926](https://github.com/Azure/azure-sdk-for-go/pull/23926) and <https://learn.microsoft.com/rest/api/cosmos-db/querying-cosmosdb-resources-using-the-rest-api#queries-that-cannot-be-served-by-gateway> for more details.

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
* Added cross-region availability and failover mechanics supporting [Azure Cosmos DB SDK multiregional environment behavior](https://learn.microsoft.com/azure/cosmos-db/nosql/troubleshoot-sdk-availability)
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
