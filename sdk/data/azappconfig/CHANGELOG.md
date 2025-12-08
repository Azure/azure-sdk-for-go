# Release History

## 2.1.0 (Unreleased)

### Features Added

* Added internal pipeline policy to normalize query parameters for deterministic request URLs
  * Query parameter names are converted to lowercase
  * Parameters are sorted in case-insensitive alphabetical order

### Breaking Changes

### Bugs Fixed

### Other Changes

## 2.0.0 (2025-10-15)

### Features Added

* Support assigning Tags and filtering by Tags
  * Added `Tags` field to `AddSettingOptions` and `SetSettingOptions` structs, enabling users to assign key-value metadata tags when creating or updating configuration settings.
  * Added `TagsFilter` field to `SettingSelector` struct, allowing retrieval of settings and revisions filtered by tags.

### Breaking Changes

* Changed `Tags` field type in `Setting` from `map[string]string` to `map[string]*string` to support null tag values and maintain compatibility with the Azure App Configuration service backend.

## 1.2.0 (2025-05-06)

### Other Changes
* Updated dependencies.

## 1.2.0-beta.1 (2024-06-11)

### Features Added
* Support ETag-per-page
  * Added field `MatchConditions` to `ListSettingsOptions` which allows specifying request conditions when iterating over pages of settings.
  * Added field `ETag` to `ListSettingsPageResponse` which contains the ETag for a page of configuration settings.

### Other Changes
* Updated dependencies.

## 1.1.0 (2024-01-17)

### Features Added
* Added support for [`Snapshots`](https://learn.microsoft.com/azure/azure-app-configuration/concept-snapshots).

### Other Changes
* Updated to latest version of `azcore`.
* Enabled spans for distributed tracing.

## 1.0.0 (2023-10-11)

### Bugs Fixed
* Check for a `Sync-Token` value before updating the cache.

### Other Changes
* Cleaned up docs and added examples.

## 0.6.0 (2023-09-20)

### Features Added
* Handle setting content type in `AddSetting` and `SetSetting` ([#19797](https://github.com/Azure/azure-sdk-for-go/issues/19797))
* Added type `SyncToken` for better type safety when handling Sync-Token header values.

### Breaking Changes
* Response types `ListRevisionsPage` and `ListSettingsPage` now have the suffix `Response` in their names.
* Method `UpdateSyncToken` on type `Client` has been replaced with `SetSyncToken`.
* Response types' `SyncToken` field type has changed from `*string` to `SyncToken`.

### Bugs Fixed
* Fixed an issue that could cause HTTP requests to fail with `http.StatusUnauthorized` in some cases.
* The pipeline policy for setting the `Sync-Token` header in HTTP requests now properly formats the value.
* The caching mechanism for `Sync-Token` values is now goroutine safe.

### Other Changes
* `NewClientFromConnectionString()` will return a more descriptive error message when parsing the connection string fails.

## 0.5.0 (2022-11-08)

### Breaking Changes
* Changed type of `OnlyIfChanged` and `OnlyIfUnchanged` option fields from `bool` to `*azcore.ETag`.

### Bugs Fixed
* `OnlyIfChanged` and `OnlyIfUnchanged` option fields have no effect
  ([#19297](https://github.com/Azure/azure-sdk-for-go/issues/19297))

## 0.4.3 (2022-10-31)

### Bugs Fixed
* Fixed missing host URL when iterating over pages.

### Other Changes
* Regenerated internal code with latest Autorest Go code generator.

## 0.4.2 (2022-10-20)

### Bugs Fixed
* Fixed a bug in `syncTokenPolicy` that could cause a panic in some conditions.

## 0.4.1 (2022-09-22)

### Features Added
* Added `NewListSettingsPager`.

## 0.4.0 (2022-05-18)

### Breaking Changes
* Moved to new location

## 0.3.1 (2022-05-12)

### Other Changes
* Update to latest `azcore`

## 0.3.0 (2022-05-10)

### Breaking Changes
* Changed argument semantics of `AddSetting`, `DeleteSetting`, `GetSetting`, `SetSetting`, and `SetReadOnly`.

## 0.2.0 (2022-04-20)

### Breaking Changes
* Upgraded to latest `azcore` which requires Go 1.18 or later.
* Renamed method `ListRevisions` to `NewListRevisionsPager` and removed `ListRevisionsPager` type.

### Other Changes
* Regenerated internal code with latest code generator.

### Bugs Fixed
* Fixed authentication in Germany West Central using connection string (#17424).

## 0.1.0 (2022-03-09)

### Features Added
* Initial release
