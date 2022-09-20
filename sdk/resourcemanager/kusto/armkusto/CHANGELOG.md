# Release History

## 1.1.0 (2022-09-15)
### Features Added

- New const `AzureSKUNameStandardEC8AsV51TBPS`
- New const `AzureSKUNameStandardEC8AdsV5`
- New const `AzureSKUNameStandardL8AsV3`
- New const `AzureSKUNameStandardL8SV3`
- New const `AzureSKUNameStandardE8DV5`
- New const `AzureSKUNameStandardE16DV4`
- New const `AzureSKUNameStandardL16AsV3`
- New const `AzureSKUNameStandardEC16AsV54TBPS`
- New const `AzureSKUNameStandardEC16AsV53TBPS`
- New const `AzureSKUNameStandardE8DV4`
- New const `DatabaseShareOriginDirect`
- New const `AzureSKUNameStandardE2DV4`
- New const `AzureSKUNameStandardE2DV5`
- New const `AzureSKUNameStandardL16SV3`
- New const `DatabaseShareOriginDataShare`
- New const `AzureSKUNameStandardE16DV5`
- New const `CallerRoleAdmin`
- New const `DatabaseShareOriginOther`
- New const `AzureSKUNameStandardEC8AsV52TBPS`
- New const `AzureSKUNameStandardE4DV4`
- New const `AzureSKUNameStandardEC16AdsV5`
- New const `CallerRoleNone`
- New const `AzureSKUNameStandardE4DV5`
- New type alias `CallerRole`
- New type alias `DatabaseShareOrigin`
- New function `PossibleDatabaseShareOriginValues() []DatabaseShareOrigin`
- New function `PossibleCallerRoleValues() []CallerRole`
- New field `CallerRole` in struct `DatabasesClientBeginCreateOrUpdateOptions`
- New field `DatabaseShareOrigin` in struct `FollowerDatabaseDefinition`
- New field `TableLevelSharingProperties` in struct `FollowerDatabaseDefinition`
- New field `DatabaseNameOverride` in struct `AttachedDatabaseConfigurationProperties`
- New field `DatabaseNamePrefix` in struct `AttachedDatabaseConfigurationProperties`
- New field `OriginalDatabaseName` in struct `ReadOnlyFollowingDatabaseProperties`
- New field `TableLevelSharingProperties` in struct `ReadOnlyFollowingDatabaseProperties`
- New field `DatabaseShareOrigin` in struct `ReadOnlyFollowingDatabaseProperties`
- New field `RetrievalStartDate` in struct `IotHubConnectionProperties`
- New field `RetrievalStartDate` in struct `EventHubConnectionProperties`
- New field `CallerRole` in struct `DatabasesClientBeginUpdateOptions`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kusto/armkusto` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
