# Release History

## 0.3.4 (Unreleased)

### Features Added
* Added `NullPartitionKey` variable to create and query documents with null partition key in CosmosDB

### Breaking Changes

### Bugs Fixed

### Other Changes

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
