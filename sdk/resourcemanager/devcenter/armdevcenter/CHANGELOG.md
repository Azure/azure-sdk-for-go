# Release History

## 1.2.0-beta.1 (2023-11-30)
### Features Added

- New enum type `CatalogConnectionState` with values `CatalogConnectionStateConnected`, `CatalogConnectionStateDisconnected`
- New enum type `CatalogResourceValidationStatus` with values `CatalogResourceValidationStatusFailed`, `CatalogResourceValidationStatusPending`, `CatalogResourceValidationStatusSucceeded`, `CatalogResourceValidationStatusUnknown`
- New enum type `CatalogSyncType` with values `CatalogSyncTypeManual`, `CatalogSyncTypeScheduled`
- New enum type `CustomizationTaskInputType` with values `CustomizationTaskInputTypeBoolean`, `CustomizationTaskInputTypeNumber`, `CustomizationTaskInputTypeString`
- New enum type `IdentityType` with values `IdentityTypeDelegatedResourceIdentity`, `IdentityTypeSystemAssignedIdentity`, `IdentityTypeUserAssignedIdentity`
- New enum type `ParameterType` with values `ParameterTypeArray`, `ParameterTypeBoolean`, `ParameterTypeInteger`, `ParameterTypeNumber`, `ParameterTypeObject`, `ParameterTypeString`
- New enum type `SingleSignOnStatus` with values `SingleSignOnStatusDisabled`, `SingleSignOnStatusEnabled`
- New enum type `VirtualNetworkType` with values `VirtualNetworkTypeManaged`, `VirtualNetworkTypeUnmanaged`
- New function `NewCatalogDevBoxDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CatalogDevBoxDefinitionsClient, error)`
- New function `*CatalogDevBoxDefinitionsClient.Get(context.Context, string, string, string, string, *CatalogDevBoxDefinitionsClientGetOptions) (CatalogDevBoxDefinitionsClientGetResponse, error)`
- New function `*CatalogDevBoxDefinitionsClient.GetErrorDetails(context.Context, string, string, string, string, *CatalogDevBoxDefinitionsClientGetErrorDetailsOptions) (CatalogDevBoxDefinitionsClientGetErrorDetailsResponse, error)`
- New function `*CatalogDevBoxDefinitionsClient.NewListByCatalogPager(string, string, string, *CatalogDevBoxDefinitionsClientListByCatalogOptions) *runtime.Pager[CatalogDevBoxDefinitionsClientListByCatalogResponse]`
- New function `*CatalogsClient.BeginConnect(context.Context, string, string, string, *CatalogsClientBeginConnectOptions) (*runtime.Poller[CatalogsClientConnectResponse], error)`
- New function `*CatalogsClient.GetSyncErrorDetails(context.Context, string, string, string, *CatalogsClientGetSyncErrorDetailsOptions) (CatalogsClientGetSyncErrorDetailsResponse, error)`
- New function `*ClientFactory.NewCatalogDevBoxDefinitionsClient() *CatalogDevBoxDefinitionsClient`
- New function `*ClientFactory.NewCustomizationTasksClient() *CustomizationTasksClient`
- New function `*ClientFactory.NewEnvironmentDefinitionsClient() *EnvironmentDefinitionsClient`
- New function `NewCustomizationTasksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CustomizationTasksClient, error)`
- New function `*CustomizationTasksClient.Get(context.Context, string, string, string, string, *CustomizationTasksClientGetOptions) (CustomizationTasksClientGetResponse, error)`
- New function `*CustomizationTasksClient.GetErrorDetails(context.Context, string, string, string, string, *CustomizationTasksClientGetErrorDetailsOptions) (CustomizationTasksClientGetErrorDetailsResponse, error)`
- New function `*CustomizationTasksClient.NewListByCatalogPager(string, string, string, *CustomizationTasksClientListByCatalogOptions) *runtime.Pager[CustomizationTasksClientListByCatalogResponse]`
- New function `NewEnvironmentDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EnvironmentDefinitionsClient, error)`
- New function `*EnvironmentDefinitionsClient.Get(context.Context, string, string, string, string, *EnvironmentDefinitionsClientGetOptions) (EnvironmentDefinitionsClientGetResponse, error)`
- New function `*EnvironmentDefinitionsClient.GetErrorDetails(context.Context, string, string, string, string, *EnvironmentDefinitionsClientGetErrorDetailsOptions) (EnvironmentDefinitionsClientGetErrorDetailsResponse, error)`
- New function `*EnvironmentDefinitionsClient.NewListByCatalogPager(string, string, string, *EnvironmentDefinitionsClientListByCatalogOptions) *runtime.Pager[EnvironmentDefinitionsClientListByCatalogResponse]`
- New struct `CatalogConflictError`
- New struct `CatalogErrorDetails`
- New struct `CatalogResourceValidationErrorDetails`
- New struct `CatalogSyncError`
- New struct `CustomerManagedKeyEncryption`
- New struct `CustomerManagedKeyEncryptionKeyIdentity`
- New struct `CustomizationTask`
- New struct `CustomizationTaskInput`
- New struct `CustomizationTaskListResult`
- New struct `CustomizationTaskProperties`
- New struct `Encryption`
- New struct `EnvironmentDefinition`
- New struct `EnvironmentDefinitionListResult`
- New struct `EnvironmentDefinitionParameter`
- New struct `EnvironmentDefinitionProperties`
- New struct `EnvironmentTypeUpdateProperties`
- New struct `SyncErrorDetails`
- New struct `SyncStats`
- New struct `UpdateProperties`
- New field `DisplayName` in struct `AllowedEnvironmentTypeProperties`
- New field `ConnectionState`, `LastConnectionTime`, `LastSyncStats`, `SyncType` in struct `CatalogProperties`
- New field `SyncType` in struct `CatalogUpdateProperties`
- New field `ValidationStatus` in struct `DevBoxDefinitionProperties`
- New field `DisplayName` in struct `EnvironmentTypeProperties`
- New field `Properties` in struct `EnvironmentTypeUpdate`
- New field `DevBoxCount`, `DisplayName`, `ManagedVirtualNetworkRegions`, `SingleSignOnStatus`, `VirtualNetworkType` in struct `PoolProperties`
- New field `DisplayName`, `ManagedVirtualNetworkRegions`, `SingleSignOnStatus`, `VirtualNetworkType` in struct `PoolUpdateProperties`
- New field `DisplayName`, `EnvironmentCount` in struct `ProjectEnvironmentTypeProperties`
- New field `DisplayName` in struct `ProjectProperties`
- New field `DisplayName` in struct `ProjectUpdateProperties`
- New field `DisplayName`, `Encryption` in struct `Properties`
- New field `Properties` in struct `Update`
- New field `ID` in struct `Usage`


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0-beta.1 (2023-10-27)
### Features Added

- New enum type `CatalogConnectionState` with values `CatalogConnectionStateConnected`, `CatalogConnectionStateDisconnected`
- New enum type `CatalogResourceValidationStatus` with values `CatalogResourceValidationStatusFailed`, `CatalogResourceValidationStatusPending`, `CatalogResourceValidationStatusSucceeded`, `CatalogResourceValidationStatusUnknown`
- New enum type `CatalogSyncType` with values `CatalogSyncTypeManual`, `CatalogSyncTypeScheduled`
- New enum type `CustomizationTaskInputType` with values `CustomizationTaskInputTypeBoolean`, `CustomizationTaskInputTypeNumber`, `CustomizationTaskInputTypeString`
- New enum type `IdentityType` with values `IdentityTypeDelegatedResourceIdentity`, `IdentityTypeSystemAssignedIdentity`, `IdentityTypeUserAssignedIdentity`
- New enum type `ParameterType` with values `ParameterTypeArray`, `ParameterTypeBoolean`, `ParameterTypeInteger`, `ParameterTypeNumber`, `ParameterTypeObject`, `ParameterTypeString`
- New enum type `SingleSignOnStatus` with values `SingleSignOnStatusDisabled`, `SingleSignOnStatusEnabled`
- New enum type `VirtualNetworkType` with values `VirtualNetworkTypeManaged`, `VirtualNetworkTypeUnmanaged`
- New function `NewCatalogDevBoxDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CatalogDevBoxDefinitionsClient, error)`
- New function `*CatalogDevBoxDefinitionsClient.Get(context.Context, string, string, string, string, *CatalogDevBoxDefinitionsClientGetOptions) (CatalogDevBoxDefinitionsClientGetResponse, error)`
- New function `*CatalogDevBoxDefinitionsClient.GetErrorDetails(context.Context, string, string, string, string, *CatalogDevBoxDefinitionsClientGetErrorDetailsOptions) (CatalogDevBoxDefinitionsClientGetErrorDetailsResponse, error)`
- New function `*CatalogDevBoxDefinitionsClient.NewListByCatalogPager(string, string, string, *CatalogDevBoxDefinitionsClientListByCatalogOptions) *runtime.Pager[CatalogDevBoxDefinitionsClientListByCatalogResponse]`
- New function `*CatalogsClient.BeginConnect(context.Context, string, string, string, *CatalogsClientBeginConnectOptions) (*runtime.Poller[CatalogsClientConnectResponse], error)`
- New function `*CatalogsClient.GetSyncErrorDetails(context.Context, string, string, string, *CatalogsClientGetSyncErrorDetailsOptions) (CatalogsClientGetSyncErrorDetailsResponse, error)`
- New function `*ClientFactory.NewCatalogDevBoxDefinitionsClient() *CatalogDevBoxDefinitionsClient`
- New function `*ClientFactory.NewCustomizationTasksClient() *CustomizationTasksClient`
- New function `*ClientFactory.NewEnvironmentDefinitionsClient() *EnvironmentDefinitionsClient`
- New function `NewCustomizationTasksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CustomizationTasksClient, error)`
- New function `*CustomizationTasksClient.Get(context.Context, string, string, string, string, *CustomizationTasksClientGetOptions) (CustomizationTasksClientGetResponse, error)`
- New function `*CustomizationTasksClient.GetErrorDetails(context.Context, string, string, string, string, *CustomizationTasksClientGetErrorDetailsOptions) (CustomizationTasksClientGetErrorDetailsResponse, error)`
- New function `*CustomizationTasksClient.NewListByCatalogPager(string, string, string, *CustomizationTasksClientListByCatalogOptions) *runtime.Pager[CustomizationTasksClientListByCatalogResponse]`
- New function `NewEnvironmentDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EnvironmentDefinitionsClient, error)`
- New function `*EnvironmentDefinitionsClient.Get(context.Context, string, string, string, string, *EnvironmentDefinitionsClientGetOptions) (EnvironmentDefinitionsClientGetResponse, error)`
- New function `*EnvironmentDefinitionsClient.GetErrorDetails(context.Context, string, string, string, string, *EnvironmentDefinitionsClientGetErrorDetailsOptions) (EnvironmentDefinitionsClientGetErrorDetailsResponse, error)`
- New function `*EnvironmentDefinitionsClient.NewListByCatalogPager(string, string, string, *EnvironmentDefinitionsClientListByCatalogOptions) *runtime.Pager[EnvironmentDefinitionsClientListByCatalogResponse]`
- New struct `CatalogConflictError`
- New struct `CatalogErrorDetails`
- New struct `CatalogResourceValidationErrorDetails`
- New struct `CatalogSyncError`
- New struct `CustomerManagedKeyEncryption`
- New struct `CustomerManagedKeyEncryptionKeyIdentity`
- New struct `CustomizationTask`
- New struct `CustomizationTaskInput`
- New struct `CustomizationTaskListResult`
- New struct `CustomizationTaskProperties`
- New struct `Encryption`
- New struct `EnvironmentDefinition`
- New struct `EnvironmentDefinitionListResult`
- New struct `EnvironmentDefinitionParameter`
- New struct `EnvironmentDefinitionProperties`
- New struct `EnvironmentTypeUpdateProperties`
- New struct `SyncErrorDetails`
- New struct `SyncStats`
- New struct `UpdateProperties`
- New field `DisplayName` in struct `AllowedEnvironmentTypeProperties`
- New field `ConnectionState`, `LastConnectionTime`, `LastSyncStats`, `SyncType` in struct `CatalogProperties`
- New field `SyncType` in struct `CatalogUpdateProperties`
- New field `ValidationStatus` in struct `DevBoxDefinitionProperties`
- New field `DisplayName` in struct `EnvironmentTypeProperties`
- New field `Properties` in struct `EnvironmentTypeUpdate`
- New field `DevBoxCount`, `DisplayName`, `ManagedVirtualNetworkRegions`, `SingleSignOnStatus`, `VirtualNetworkType` in struct `PoolProperties`
- New field `DisplayName`, `ManagedVirtualNetworkRegions`, `SingleSignOnStatus`, `VirtualNetworkType` in struct `PoolUpdateProperties`
- New field `DisplayName`, `EnvironmentCount` in struct `ProjectEnvironmentTypeProperties`
- New field `DisplayName` in struct `ProjectProperties`
- New field `DisplayName` in struct `ProjectUpdateProperties`
- New field `DisplayName`, `Encryption` in struct `Properties`
- New field `Properties` in struct `Update`
- New field `ID` in struct `Usage`


## 1.0.0 (2023-05-26)
### Breaking Changes

- Type of `ProjectEnvironmentTypeProperties.Status` has been changed from `*EnableStatus` to `*EnvironmentTypeEnableStatus`
- Type of `ProjectEnvironmentTypeUpdateProperties.Status` has been changed from `*EnableStatus` to `*EnvironmentTypeEnableStatus`
- Type of `ScheduleProperties.State` has been changed from `*EnableStatus` to `*ScheduleEnableStatus`
- Type of `ScheduleUpdateProperties.State` has been changed from `*EnableStatus` to `*ScheduleEnableStatus`
- Enum `EnableStatus` has been removed
- Field `Offer`, `Publisher`, `SKU` of struct `ImageReference` has been removed

### Features Added

- New enum type `EnvironmentTypeEnableStatus` with values `EnvironmentTypeEnableStatusDisabled`, `EnvironmentTypeEnableStatusEnabled`
- New enum type `HealthStatus` with values `HealthStatusHealthy`, `HealthStatusPending`, `HealthStatusUnhealthy`, `HealthStatusUnknown`, `HealthStatusWarning`
- New enum type `ScheduleEnableStatus` with values `ScheduleEnableStatusDisabled`, `ScheduleEnableStatusEnabled`
- New enum type `StopOnDisconnectEnableStatus` with values `StopOnDisconnectEnableStatusDisabled`, `StopOnDisconnectEnableStatusEnabled`
- New function `*PoolsClient.BeginRunHealthChecks(context.Context, string, string, string, *PoolsClientBeginRunHealthChecksOptions) (*runtime.Poller[PoolsClientRunHealthChecksResponse], error)`
- New function `*NetworkConnectionsClient.NewListOutboundNetworkDependenciesEndpointsPager(string, string, *NetworkConnectionsClientListOutboundNetworkDependenciesEndpointsOptions) *runtime.Pager[NetworkConnectionsClientListOutboundNetworkDependenciesEndpointsResponse]`
- New struct `EndpointDependency`
- New struct `EndpointDetail`
- New struct `HealthStatusDetail`
- New struct `OutboundEnvironmentEndpoint`
- New struct `OutboundEnvironmentEndpointCollection`
- New struct `StopOnDisconnectConfiguration`
- New field `HibernateSupport` in struct `ImageProperties`
- New field `HealthStatus`, `HealthStatusDetails`, `StopOnDisconnect` in struct `PoolProperties`
- New field `StopOnDisconnect` in struct `PoolUpdateProperties`
- New field `MaxDevBoxesPerUser` in struct `ProjectProperties`
- New field `MaxDevBoxesPerUser` in struct `ProjectUpdateProperties`


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