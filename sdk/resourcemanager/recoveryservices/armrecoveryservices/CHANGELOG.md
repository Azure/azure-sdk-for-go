# Release History

## 0.3.0 (2021-12-28)
### Breaking Changes

### Features Added

- New const `BackupStorageVersionUnassigned`
- New const `BackupStorageVersionV2`
- New const `BackupStorageVersionV1`
- New function `PossibleBackupStorageVersionValues() []BackupStorageVersion`
- New function `BackupStorageVersion.ToPtr() *BackupStorageVersion`
- New field `BackupStorageVersion` in struct `VaultProperties`


## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-10-20)

- Initial preview release.
