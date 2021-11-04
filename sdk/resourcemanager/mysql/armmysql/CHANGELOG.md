# Release History

## 0.2.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-10-20)

- Initial preview release.
