# Release History

## 2.0.0 (2022-09-01)
### Breaking Changes

- Function `*LinkedServerClient.Delete` has been removed
- Function `*Client.Update` has been removed
- Struct `ClientUpdateOptions` has been removed
- Struct `LinkedServerClientDeleteOptions` has been removed

### Features Added

- New function `*Client.BeginUpdate(context.Context, string, string, UpdateParameters, *ClientBeginUpdateOptions) (*runtime.Poller[ClientUpdateResponse], error)`
- New function `*LinkedServerClient.BeginDelete(context.Context, string, string, string, *LinkedServerClientBeginDeleteOptions) (*runtime.Poller[LinkedServerClientDeleteResponse], error)`
- New struct `ClientBeginUpdateOptions`
- New struct `LinkedServerClientBeginDeleteOptions`
- New field `Authnotrequired` in struct `CommonPropertiesRedisConfiguration`
- New field `AofBackupEnabled` in struct `CommonPropertiesRedisConfiguration`
- New field `PreferredDataArchiveAuthMethod` in struct `ImportRDBParameters`
- New field `PreferredDataArchiveAuthMethod` in struct `ExportRDBParameters`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redis/armredis` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).