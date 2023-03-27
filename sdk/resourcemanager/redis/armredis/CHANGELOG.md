# Release History

## 2.2.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 2.1.0 (2022-10-14)

### Features Added

- New field `PrimaryHostName` in struct `LinkedServerCreateProperties`
- New field `GeoReplicatedPrimaryHostName` in struct `LinkedServerCreateProperties`
- New field `GeoReplicatedPrimaryHostName` in struct `LinkedServerProperties`
- New field `PrimaryHostName` in struct `LinkedServerProperties`


## 2.0.0 (2022-09-01)
### Breaking Changes

- Operation `LinkedServerClient.Delete` has been changed to LRO, use `LinkedServerClient.BeginDelete` instead
- Operation `*Client.Update` has been changed to LRO, use `Client.BeginUpdate` instead

### Features Added

- New field `Authnotrequired` in struct `CommonPropertiesRedisConfiguration`
- New field `AofBackupEnabled` in struct `CommonPropertiesRedisConfiguration`
- New field `PreferredDataArchiveAuthMethod` in struct `ImportRDBParameters`
- New field `PreferredDataArchiveAuthMethod` in struct `ExportRDBParameters`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redis/armredis` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
