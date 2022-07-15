# Release History

## 1.1.0-beta.1 (2022-07-15)
### Features Added

- New const `ResetAllToDefaultTrue`
- New const `ResetAllToDefaultFalse`
- New function `*CheckNameAvailabilityWithoutLocationClient.Execute(context.Context, NameAvailabilityRequest, *CheckNameAvailabilityWithoutLocationClientExecuteOptions) (CheckNameAvailabilityWithoutLocationClientExecuteResponse, error)`
- New function `NewLogFilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LogFilesClient, error)`
- New function `*LogFilesClient.NewListByServerPager(string, string, *LogFilesClientListByServerOptions) *runtime.Pager[LogFilesClientListByServerResponse]`
- New function `PossibleResetAllToDefaultValues() []ResetAllToDefault`
- New function `*BackupsClient.Put(context.Context, string, string, string, *BackupsClientPutOptions) (BackupsClientPutResponse, error)`
- New function `NewCheckNameAvailabilityWithoutLocationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CheckNameAvailabilityWithoutLocationClient, error)`
- New struct `BackupsClientPutOptions`
- New struct `BackupsClientPutResponse`
- New struct `CheckNameAvailabilityWithoutLocationClient`
- New struct `CheckNameAvailabilityWithoutLocationClientExecuteOptions`
- New struct `CheckNameAvailabilityWithoutLocationClientExecuteResponse`
- New struct `LogFile`
- New struct `LogFileListResult`
- New struct `LogFileProperties`
- New struct `LogFilesClient`
- New struct `LogFilesClientListByServerOptions`
- New struct `LogFilesClientListByServerResponse`
- New field `PageSize` in struct `ConfigurationsClientListByServerOptions`
- New field `Tags` in struct `ConfigurationsClientListByServerOptions`
- New field `Keyword` in struct `ConfigurationsClientListByServerOptions`
- New field `Page` in struct `ConfigurationsClientListByServerOptions`
- New field `ResetAllToDefault` in struct `ConfigurationListForBatchUpdate`
- New field `Version` in struct `ServerPropertiesForUpdate`
- New field `Location` in struct `VirtualNetworkSubnetUsageResult`
- New field `SubscriptionID` in struct `VirtualNetworkSubnetUsageResult`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysqlflexibleservers` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).