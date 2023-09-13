# Release History

## 0.6.0 (Unreleased)

### Features Added
* Handle setting content type in `AddSetting` and `SetSetting` ([#19797](https://github.com/Azure/azure-sdk-for-go/issues/19797))

### Breaking Changes
* Response types `ListRevisionsPage` and `ListSettingsPage` now have the suffix `Response` in their names.

### Bugs Fixed

### Other Changes

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
