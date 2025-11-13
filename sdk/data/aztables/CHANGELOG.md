# Release History

## 1.4.1 (2025-11-13)

### Bugs Fixed

- Fix an issue that the Storage Table token audiences for sovereign clouds are incorrect. (PR#25534)

## 1.4.0 (2025-06-19)

### Features Added

- Added support for sovereign clouds.

### Other Changes

- Update dependencies.

## 1.3.0 (2024-11-05)

### Features Added

- Client/ServiceClient now supports `azcore.TokenCredential` authentication with Azure Cosmos DB for Table.

### Other Changes

- Updated dependencies.

## 1.2.0 (2024-03-11)

### Features Added

- Methods `Client.AddEntity` and `ServiceClient.NewListTablesPager` now include OData metadata in their responses.
- The amount of OData metadata returned has been made configurable for the following methods:
  - `Client.AddEntity`, `Client.GetEntity`, `Client.NewListEntitiesPager`, and `ServiceClient.NewListTablesPager`.
  - Use one of the following constants to specify the amount: `MetadataFormatFull`, `MetadataFormatMinimal`, or `MetadataFormatNone`.

### Bugs Fixed

- Fixed an issue that could cause `Client.NewListEntitiesPager` to skip pages in some cases.
- Fixed an issue that could cause unmarshaling empty time values to fail.

### Other Changes

- Update dependencies.

## 1.1.0 (2023-11-14)

### Features Added

- Enabled spans for distributed tracing.

### Bugs Fixed

- Internal calls in `Client.SubmitTransaction` now honor the caller's context.

### Other Changes

- Updated to latest version of `azcore`.
- Clients now share the underlying `*azcore.Client`.

## 1.0.2 (2023-07-20)

### Bugs Fixed

- Escape single-quote characters in partition and row keys.

### Other Changes

- Update dependencies.

## 1.0.1 (2022-06-16)

### Bugs Fixed

- Accept empty `rowKey` parameter.

## 1.0.0 (2022-05-16)

### Breaking Changes

- For type `EDMEntity` renamed field `Id` to `ID`, `Etag` to `ETag`

## 0.8.1 (2022-05-12)

### Other Changes

- Update to latest `azcore` and `internal` modules

## 0.8.0 (2022-04-20)

### Features Added

- Added `TableErrorCode` to help recover from and understand error responses

### Breaking Changes

- Renamed `InsertEntityResponse/Options` to `UpsertEntityResponse/Options`
- Renamed `PossibleGeoReplicationStatusTypeValues` to `PossibleGeoReplicationStatusValues`
- Renamed the following methods
  - `Client.ListEntities` to `Client.NewListEntitiesPager`
  - `ServiceClient.ListTables` to `ServiceClient.NewListTablesPager`

### Bugs Fixed

- Convert `Start` and `Expiry` times in `AccessPolicy` to UTC format as required by the service.
- Fixed `moduleName` to report the module name as part of telemetry.

### Other Changes

- Fixed bugs in some live tests.

## 0.7.0 (2022-04-05)

### Features Added

- Added the `NextTableName` continuation token option to `ListTablesOptions`
- Added the `TableName` property to `CreateTableResponse`

### Breaking Changes

- This module now requires Go 1.18
- Removed the `ODataID`, `ODataEditLink`, and `ODataType` on `TableProperties`
- Removed `ODataMetadata` on `ListTablesPageResponse`
- Removed `ResponsePreference` on `AddEntityOptions`
- Renamed `ListEntitiesOptions.PartitionKey` to `ListEntitiesOptions.NextPartitionKey`.
- Renamed `ListEntitiesOptionsRowKey` to `ListEntitiesOptions.NextRowKey`
- Renamed `Client.Create` to `Client.CreateTable`
- Renamed `ListEntitiesPageResponse` to `ListEntitiesResponse`
- Removed the `Entity` prefix on `EntityUpdateModeMerge` and `EntityUpdateModeReplace`
- Renamed `Client.InsertEntity` to `Client.UpsertEntity`
- Removed the `Continuation` prefix from `ContinuationNextPartitionKey`, `ContinuationNextRowKey`, and `ContinuationNextTable`
- Removed the `ResponseFormat` type
- Renamed `Client.List` to `Client.ListEntities`
- Renamed `Client.GetTableSASToken` to `Client.GetTableSASURL` and `ServiceClient.GetAccountSASToken` to `ServiceClient.GetAccountSASURL`
- `ServiceClient.GetProperties` returns a `ServiceProperties` struct which can be used on the `ServiceClient.SetProperties`
- Removed the `Type` suffix from `GeoReplicationStatusType`
- `ServiceClient.CreateTable` returns a response struct with the name of the table created, not a `Client`
- `SASSignatureValues.NewSASQueryParameters` is now `SASSignatureValues.Sign` and returns an encoded SAS

## 0.6.0 (2022-03-08)

### Breaking Changes

- Prefixed all `TransactionType` constants with `TransactionType`.
- Prefixed all `EntityUpdateMode` constants with `EntityUpdateMode`.
- Changed the `SharedKeyCredential.ComputeHMACSHA256` method to a private method.
- Changed the `ListTablesPager` and `ListEntitiesPager` to structs.
- Renamed the `ResponseProperties` type to `TableProperties`.
- Removing `ContentType` from the `TransactionResponse` struct.
- Update `ListEntitiesPager` and `ListTablesPager`.
  - The `More` method checks whether there are more pages to retrieve.
  - The `NextPage(context.Context)` method gets the next page and returns a response and an `error`.
- Removed `RawResponse` from all Response structs
- `TransactionResponse` is an empty struct

## 0.5.0 (2022-01-12)

### Other Changes

- Updates `azcore` dependency from `v0.20.0` to `v0.21.0`

## 0.4.0 (2021-11-09)

### Features Added

- Added `NextPagePartitionKey` and `NextPageRowKey` to `ListEntitiesPager` for retrieving continuation tokens.
- Added `PartitionKey` and `RowKey` to `ListEntitiesOptions` for using exposed continuation tokens.

### Bug Fixes

- Fixed a bug on transactional batches where `InsertMerge` and `InsertReplace` failed if the entity did not exist.

## 0.3.0 (2021-11-02)

### Features Added

- Added `NewClientWithNoCredential` and `NewServiceClientWithNoCredential` for authenticating the `Client` and `ServiceClient` with SAS URLs
- Added `NewClientWithSharedKey` and `NewServiceClientWithSharedKey` for authenticating the `Client` and `ServiceClient` with Shared Keys

### Breaking Changes

- `NewClient` and `NewServiceClient` is now used for authenticating the `Client` and `ServiceClient` with credentials from `azidentity` only.
- `ClientOptions` embeds `azcore.ClientOptions` and removes all named fields.

## 0.2.0 (2021-10-05)

### Features Added

- Added the `ClientOptions.PerTryPolicies` for policies that execute once per retry of an operation.

### Breaking Changes

- Changed the `ClientOptions.PerCallOptions` field name to `ClientOptions.PerCallPolicies`
- Changed the `ClientOptions.Transporter` field name to `ClientOptions.Transport`

## 0.1.0 (09-07-2021)

- This is the initial release of the `aztables` library
