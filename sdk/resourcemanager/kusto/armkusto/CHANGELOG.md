# Release History

## 2.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.1.0 (2023-10-27)
### Features Added

- New value `LanguageExtensionImageNamePython3108DL`, `LanguageExtensionImageNamePythonCustomImage` added to enum type `LanguageExtensionImageName`
- New enum type `Language` with values `LanguagePython`
- New enum type `VnetState` with values `VnetStateDisabled`, `VnetStateEnabled`
- New function `*ClientFactory.NewSandboxCustomImagesClient() *SandboxCustomImagesClient`
- New function `NewSandboxCustomImagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SandboxCustomImagesClient, error)`
- New function `*SandboxCustomImagesClient.CheckNameAvailability(context.Context, string, string, SandboxCustomImagesCheckNameRequest, *SandboxCustomImagesClientCheckNameAvailabilityOptions) (SandboxCustomImagesClientCheckNameAvailabilityResponse, error)`
- New function `*SandboxCustomImagesClient.BeginCreateOrUpdate(context.Context, string, string, string, SandboxCustomImage, *SandboxCustomImagesClientBeginCreateOrUpdateOptions) (*runtime.Poller[SandboxCustomImagesClientCreateOrUpdateResponse], error)`
- New function `*SandboxCustomImagesClient.BeginDelete(context.Context, string, string, string, *SandboxCustomImagesClientBeginDeleteOptions) (*runtime.Poller[SandboxCustomImagesClientDeleteResponse], error)`
- New function `*SandboxCustomImagesClient.Get(context.Context, string, string, string, *SandboxCustomImagesClientGetOptions) (SandboxCustomImagesClientGetResponse, error)`
- New function `*SandboxCustomImagesClient.NewListByClusterPager(string, string, *SandboxCustomImagesClientListByClusterOptions) *runtime.Pager[SandboxCustomImagesClientListByClusterResponse]`
- New function `*SandboxCustomImagesClient.BeginUpdate(context.Context, string, string, string, SandboxCustomImage, *SandboxCustomImagesClientBeginUpdateOptions) (*runtime.Poller[SandboxCustomImagesClientUpdateResponse], error)`
- New struct `SandboxCustomImage`
- New struct `SandboxCustomImageProperties`
- New struct `SandboxCustomImagesCheckNameRequest`
- New struct `SandboxCustomImagesListResult`
- New field `Zones` in struct `ClusterUpdate`
- New field `IPAddress` in struct `EndpointDetail`
- New field `LanguageExtensionCustomImageName` in struct `LanguageExtension`
- New field `State` in struct `VirtualNetworkConfiguration`


## 2.0.0 (2023-07-28)
### Breaking Changes

- `LanguageExtensionImageNamePython3912`, `LanguageExtensionImageNamePython3912IncludeDeepLearning` from enum `LanguageExtensionImageName` has been removed

### Features Added

- New value `StateMigrated` added to enum type `State`
- New enum type `MigrationClusterRole` with values `MigrationClusterRoleDestination`, `MigrationClusterRoleSource`
- New function `*ClientFactory.NewDatabaseClient() *DatabaseClient`
- New function `*ClustersClient.BeginMigrate(context.Context, string, string, ClusterMigrateRequest, *ClustersClientBeginMigrateOptions) (*runtime.Poller[ClustersClientMigrateResponse], error)`
- New function `NewDatabaseClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DatabaseClient, error)`
- New function `*DatabaseClient.InviteFollower(context.Context, string, string, string, DatabaseInviteFollowerRequest, *DatabaseClientInviteFollowerOptions) (DatabaseClientInviteFollowerResponse, error)`
- New struct `ClusterMigrateRequest`
- New struct `DatabaseInviteFollowerRequest`
- New struct `DatabaseInviteFollowerResult`
- New struct `MigrationClusterProperties`
- New struct `SuspensionDetails`
- New field `MigrationCluster` in struct `ClusterProperties`
- New field `NextLink` in struct `DatabaseListResult`
- New field `Skiptoken`, `Top` in struct `DatabasesClientListByClusterOptions`
- New field `AzureAsyncOperation` in struct `OperationsResultsLocationClientGetResponse`
- New field `SuspensionDetails` in struct `ReadOnlyFollowingDatabaseProperties`
- New field `KeyVaultProperties`, `SuspensionDetails` in struct `ReadWriteDatabaseProperties`


## 1.3.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.3.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.2.0 (2023-02-24)
### Features Added

- New value `AzureSKUNameStandardL32AsV3`, `AzureSKUNameStandardL32SV3` added to type alias `AzureSKUName`
- New value `DataConnectionKindCosmosDb` added to type alias `DataConnectionKind`
- New value `ProvisioningStateCanceled` added to type alias `ProvisioningState`
- New type alias `LanguageExtensionImageName` with values `LanguageExtensionImageNamePython3108`, `LanguageExtensionImageNamePython365`, `LanguageExtensionImageNamePython3912`, `LanguageExtensionImageNamePython3912IncludeDeepLearning`, `LanguageExtensionImageNameR`
- New function `*CosmosDbDataConnection.GetDataConnection() *DataConnection`
- New function `NewSKUsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SKUsClient, error)`
- New function `*SKUsClient.NewListPager(string, *SKUsClientListOptions) *runtime.Pager[SKUsClientListResponse]`
- New struct `CosmosDbDataConnection`
- New struct `CosmosDbDataConnectionProperties`
- New struct `ResourceSKUCapabilities`
- New struct `ResourceSKUZoneDetails`
- New struct `SKUsClient`
- New struct `SKUsClientListResponse`
- New field `LanguageExtensionImageName` in struct `LanguageExtension`
- New field `ZoneDetails` in struct `SKULocationInfoItem`
- New field `FunctionsToExclude` in struct `TableLevelSharingProperties`
- New field `FunctionsToInclude` in struct `TableLevelSharingProperties`


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
