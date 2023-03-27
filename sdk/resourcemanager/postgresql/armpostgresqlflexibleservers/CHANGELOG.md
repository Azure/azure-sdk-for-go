# Release History

## 2.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module

## 2.0.0 (2023-01-27)
### Breaking Changes

- Function `*CheckNameAvailabilityClient.Execute` parameter(s) have been changed from `(context.Context, NameAvailabilityRequest, *CheckNameAvailabilityClientExecuteOptions)` to `(context.Context, CheckNameAvailabilityRequest, *CheckNameAvailabilityClientExecuteOptions)`
- Function `*ConfigurationsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, Configuration, *ConfigurationsClientBeginUpdateOptions)` to `(context.Context, string, string, string, ConfigurationForUpdate, *ConfigurationsClientBeginUpdateOptions)`
- Type of `NameAvailability.Reason` has been changed from `*Reason` to `*CheckNameAvailabilityReason`
- Type alias `Reason` has been removed
- Struct `NameAvailabilityRequest` has been removed
- Field `Location` of struct `ServerForUpdate` has been removed

### Features Added

- New value `CreateModeGeoRestore`, `CreateModeReplica` added to type alias `CreateMode`
- New value `HighAvailabilityModeSameZone` added to type alias `HighAvailabilityMode`
- New type alias `ActiveDirectoryAuthEnum` with values `ActiveDirectoryAuthEnumDisabled`, `ActiveDirectoryAuthEnumEnabled`
- New type alias `ArmServerKeyType` with values `ArmServerKeyTypeAzureKeyVault`, `ArmServerKeyTypeSystemAssigned`
- New type alias `CheckNameAvailabilityReason` with values `CheckNameAvailabilityReasonAlreadyExists`, `CheckNameAvailabilityReasonInvalid`
- New type alias `IdentityType` with values `IdentityTypeNone`, `IdentityTypeSystemAssigned`, `IdentityTypeUserAssigned`
- New type alias `Origin` with values `OriginFull`
- New type alias `PasswordAuthEnum` with values `PasswordAuthEnumDisabled`, `PasswordAuthEnumEnabled`
- New type alias `PrincipalType` with values `PrincipalTypeGroup`, `PrincipalTypeServicePrincipal`, `PrincipalTypeUnknown`, `PrincipalTypeUser`
- New type alias `ReplicationRole` with values `ReplicationRoleAsyncReplica`, `ReplicationRoleGeoAsyncReplica`, `ReplicationRoleGeoSyncReplica`, `ReplicationRoleNone`, `ReplicationRolePrimary`, `ReplicationRoleSecondary`, `ReplicationRoleSyncReplica`, `ReplicationRoleWalReplica`
- New function `NewAdministratorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AdministratorsClient, error)`
- New function `*AdministratorsClient.BeginCreate(context.Context, string, string, string, ActiveDirectoryAdministratorAdd, *AdministratorsClientBeginCreateOptions) (*runtime.Poller[AdministratorsClientCreateResponse], error)`
- New function `*AdministratorsClient.BeginDelete(context.Context, string, string, string, *AdministratorsClientBeginDeleteOptions) (*runtime.Poller[AdministratorsClientDeleteResponse], error)`
- New function `*AdministratorsClient.Get(context.Context, string, string, string, *AdministratorsClientGetOptions) (AdministratorsClientGetResponse, error)`
- New function `*AdministratorsClient.NewListByServerPager(string, string, *AdministratorsClientListByServerOptions) *runtime.Pager[AdministratorsClientListByServerResponse]`
- New function `NewBackupsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsClient, error)`
- New function `*BackupsClient.Get(context.Context, string, string, string, *BackupsClientGetOptions) (BackupsClientGetResponse, error)`
- New function `*BackupsClient.NewListByServerPager(string, string, *BackupsClientListByServerOptions) *runtime.Pager[BackupsClientListByServerResponse]`
- New function `NewCheckNameAvailabilityWithLocationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CheckNameAvailabilityWithLocationClient, error)`
- New function `*CheckNameAvailabilityWithLocationClient.Execute(context.Context, string, CheckNameAvailabilityRequest, *CheckNameAvailabilityWithLocationClientExecuteOptions) (CheckNameAvailabilityWithLocationClientExecuteResponse, error)`
- New function `NewReplicasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ReplicasClient, error)`
- New function `*ReplicasClient.NewListByServerPager(string, string, *ReplicasClientListByServerOptions) *runtime.Pager[ReplicasClientListByServerResponse]`
- New struct `ActiveDirectoryAdministrator`
- New struct `ActiveDirectoryAdministratorAdd`
- New struct `AdministratorListResult`
- New struct `AdministratorProperties`
- New struct `AdministratorPropertiesForAdd`
- New struct `AdministratorsClient`
- New struct `AdministratorsClientCreateResponse`
- New struct `AdministratorsClientDeleteResponse`
- New struct `AdministratorsClientListByServerResponse`
- New struct `AuthConfig`
- New struct `BackupsClient`
- New struct `BackupsClientListByServerResponse`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityWithLocationClient`
- New struct `ConfigurationForUpdate`
- New struct `DataEncryption`
- New struct `FastProvisioningEditionCapability`
- New struct `ReplicasClient`
- New struct `ReplicasClientListByServerResponse`
- New struct `ServerBackup`
- New struct `ServerBackupListResult`
- New struct `ServerBackupProperties`
- New struct `StorageTierCapability`
- New struct `UserAssignedIdentity`
- New struct `UserIdentity`
- New field `FastProvisioningSupported` in struct `CapabilityProperties`
- New field `SupportedFastProvisioningEditions` in struct `CapabilityProperties`
- New field `Identity` in struct `Server`
- New field `Identity` in struct `ServerForUpdate`
- New field `AuthConfig` in struct `ServerProperties`
- New field `DataEncryption` in struct `ServerProperties`
- New field `ReplicaCapacity` in struct `ServerProperties`
- New field `ReplicationRole` in struct `ServerProperties`
- New field `AuthConfig` in struct `ServerPropertiesForUpdate`
- New field `DataEncryption` in struct `ServerPropertiesForUpdate`
- New field `ReplicationRole` in struct `ServerPropertiesForUpdate`
- New field `Version` in struct `ServerPropertiesForUpdate`
- New field `SupportedVersionsToUpgrade` in struct `ServerVersionCapability`
- New field `SupportedUpgradableTierList` in struct `StorageMBCapability`


## 1.1.0 (2022-07-21)
### Features Added

- New const `ServerVersionFourteen`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).