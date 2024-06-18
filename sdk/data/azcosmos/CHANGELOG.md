# Release History

## 1.0.4 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

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
