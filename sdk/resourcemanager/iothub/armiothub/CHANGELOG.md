# Release History

## 0.2.0 (2021-10-29)
### Breaking Changes

- Function `NewResourceProviderCommonClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateLinkResourcesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIotHubResourceClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewCertificatesClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewOperationsClient` parameter(s) have been changed from `(*arm.Connection)` to `(azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewIotHubClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewPrivateEndpointConnectionsClient` parameter(s) have been changed from `(*arm.Connection, string)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`

### New Content


Total 7 breaking change(s), 0 additive change(s).


## 0.1.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.1.0 (2021-10-18)

- Initial preview release.
