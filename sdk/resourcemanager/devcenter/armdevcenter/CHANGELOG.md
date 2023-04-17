# Release History

## 0.5.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.5.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.4.0 (2022-11-24)
### Breaking Changes

- Type of `AllowedEnvironmentTypeProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `AttachedNetworkConnectionProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `CatalogProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `DevBoxDefinitionProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `EnvironmentTypeProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `GalleryProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `ImageProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `ImageVersionProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `NetworkProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `PoolProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `ProjectEnvironmentTypeProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `ProjectProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `Properties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `ScheduleProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Operation `*NetworkConnectionsClient.RunHealthChecks` has been changed to LRO, use `*NetworkConnectionsClient.BeginRunHealthChecks` instead.

### Features Added

- New type alias `CheckNameAvailabilityReason`
- New type alias `HibernateSupport`
- New type alias `ProvisioningState`
- New function `NewCheckNameAvailabilityClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CheckNameAvailabilityClient, error)`
- New function `*CheckNameAvailabilityClient.Execute(context.Context, CheckNameAvailabilityRequest, *CheckNameAvailabilityClientExecuteOptions) (CheckNameAvailabilityClientExecuteResponse, error)`
- New struct `CheckNameAvailabilityClient`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityResponse`
- New field `HibernateSupport` in struct `DevBoxDefinitionProperties`
- New field `HibernateSupport` in struct `DevBoxDefinitionUpdateProperties`
- New field `DevCenterURI` in struct `ProjectProperties`
- New field `DevCenterURI` in struct `Properties`


## 0.3.0 (2022-10-27)
### Breaking Changes

- Type of `OperationStatus.Error` has been changed from `*OperationStatusError` to `*ErrorDetail`
- Struct `OperationStatusError` has been removed

### Features Added

- New const `CatalogSyncStateFailed`
- New const `CatalogSyncStateSucceeded`
- New const `CatalogSyncStateInProgress`
- New const `CatalogSyncStateCanceled`
- New type alias `CatalogSyncState`
- New function `PossibleCatalogSyncStateValues() []CatalogSyncState`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `OperationStatusResult`
- New anonymous field `Schedule` in struct `SchedulesClientUpdateResponse`
- New field `Operations` in struct `OperationStatus`
- New field `ResourceID` in struct `OperationStatus`
- New field `SyncState` in struct `CatalogProperties`


## 0.2.0 (2022-09-29)
### Features Added

- New function `NewProjectAllowedEnvironmentTypesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectAllowedEnvironmentTypesClient, error)`
- New function `*ProjectAllowedEnvironmentTypesClient.NewListPager(string, string, *ProjectAllowedEnvironmentTypesClientListOptions) *runtime.Pager[ProjectAllowedEnvironmentTypesClientListResponse]`
- New function `*ProjectAllowedEnvironmentTypesClient.Get(context.Context, string, string, string, *ProjectAllowedEnvironmentTypesClientGetOptions) (ProjectAllowedEnvironmentTypesClientGetResponse, error)`
- New struct `AllowedEnvironmentType`
- New struct `AllowedEnvironmentTypeListResult`
- New struct `AllowedEnvironmentTypeProperties`
- New struct `ProjectAllowedEnvironmentTypesClient`
- New struct `ProjectAllowedEnvironmentTypesClientGetOptions`
- New struct `ProjectAllowedEnvironmentTypesClientGetResponse`
- New struct `ProjectAllowedEnvironmentTypesClientListOptions`
- New struct `ProjectAllowedEnvironmentTypesClientListResponse`


## 0.1.0 (2022-07-25)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devcenter/armdevcenter` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).