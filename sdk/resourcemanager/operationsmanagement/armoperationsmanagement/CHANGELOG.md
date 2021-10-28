# Release History

## 0.2.0 (2021-10-28)
### Breaking Changes

- Function `NewManagementConfigurationsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewSolutionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewManagementAssociationsClient` parameter(s) have been changed from `(*arm.Connection, string, string, string, string)` to `(string, string, string, string, azcore.TokenCredential, *arm.ClientOptions)`

### New Content


Total 4 breaking change(s), 0 additive change(s).


## 0.1.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.1.0 (2021-10-20)

- Initial preview release.
