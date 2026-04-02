# Release History

## 3.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Type of `CustomerManagedKeyEncryption.KeyEncryptionKeyIdentity` has been changed from `*CustomerManagedKeyEncryptionKeyIdentity` to `*KeyEncryptionKeyIdentity`
- Type of `OperationStatus.PercentComplete` has been changed from `*float32` to `*float64`
- Type of `OperationStatus.Properties` has been changed from `any` to `map[string]any`
- Type of `OperationStatusResult.PercentComplete` has been changed from `*float32` to `*float64`
- Enum `IdentityType` has been removed
- Struct `CustomerManagedKeyEncryptionKeyIdentity` has been removed
- Field `ResourceID` of struct `OperationStatus` has been removed
- Field `ResourceID` of struct `OperationStatusResult` has been removed
- Field `Top` of struct `SchedulesClientGetOptions` has been removed

### Features Added

- New value `CatalogItemTypeImageDefinition` added to enum type `CatalogItemType`
- New value `DomainJoinTypeNone` added to enum type `DomainJoinType`
- New value `HealthCheckStatusInformational` added to enum type `HealthCheckStatus`
- New enum type `AutoImageBuildStatus` with values `AutoImageBuildStatusDisabled`, `AutoImageBuildStatusEnabled`
- New enum type `AutoStartEnableStatus` with values `AutoStartEnableStatusDisabled`, `AutoStartEnableStatusEnabled`
- New enum type `AzureAiServicesMode` with values `AzureAiServicesModeAutoDeploy`, `AzureAiServicesModeDisabled`
- New enum type `CmkIdentityType` with values `CmkIdentityTypeSystemAssigned`, `CmkIdentityTypeUserAssigned`
- New enum type `CustomizationTaskInputType` with values `CustomizationTaskInputTypeBoolean`, `CustomizationTaskInputTypeNumber`, `CustomizationTaskInputTypeString`
- New enum type `DayOfWeek` with values `DayOfWeekFriday`, `DayOfWeekMonday`, `DayOfWeekSaturday`, `DayOfWeekSunday`, `DayOfWeekThursday`, `DayOfWeekTuesday`, `DayOfWeekWednesday`
- New enum type `DevBoxDeleteMode` with values `DevBoxDeleteModeAuto`, `DevBoxDeleteModeManual`
- New enum type `DevBoxTunnelEnableStatus` with values `DevBoxTunnelEnableStatusDisabled`, `DevBoxTunnelEnableStatusEnabled`
- New enum type `DevboxDisksEncryptionEnableStatus` with values `DevboxDisksEncryptionEnableStatusDisabled`, `DevboxDisksEncryptionEnableStatusEnabled`
- New enum type `ImageDefinitionBuildStatus` with values `ImageDefinitionBuildStatusCancelled`, `ImageDefinitionBuildStatusFailed`, `ImageDefinitionBuildStatusRunning`, `ImageDefinitionBuildStatusSucceeded`, `ImageDefinitionBuildStatusTimedOut`, `ImageDefinitionBuildStatusValidationFailed`
- New enum type `InstallAzureMonitorAgentEnableStatus` with values `InstallAzureMonitorAgentEnableStatusDisabled`, `InstallAzureMonitorAgentEnableStatusEnabled`
- New enum type `KeepAwakeEnableStatus` with values `KeepAwakeEnableStatusDisabled`, `KeepAwakeEnableStatusEnabled`
- New enum type `KeyEncryptionKeyIdentityType` with values `KeyEncryptionKeyIdentityTypeDelegatedResourceIdentity`, `KeyEncryptionKeyIdentityTypeSystemAssignedIdentity`, `KeyEncryptionKeyIdentityTypeUserAssignedIdentity`
- New enum type `MicrosoftHostedNetworkEnableStatus` with values `MicrosoftHostedNetworkEnableStatusDisabled`, `MicrosoftHostedNetworkEnableStatusEnabled`
- New enum type `PolicyAction` with values `PolicyActionAllow`, `PolicyActionDeny`
- New enum type `PoolDevBoxDefinitionType` with values `PoolDevBoxDefinitionTypeReference`, `PoolDevBoxDefinitionTypeValue`
- New enum type `ProjectCustomizationIdentityType` with values `ProjectCustomizationIdentityTypeSystemAssignedIdentity`, `ProjectCustomizationIdentityTypeUserAssignedIdentity`
- New enum type `ResourceType` with values `ResourceTypeAttachedNetworks`, `ResourceTypeImages`, `ResourceTypeSKUs`
- New enum type `ServerlessGpuSessionsMode` with values `ServerlessGpuSessionsModeAutoDeploy`, `ServerlessGpuSessionsModeDisabled`
- New enum type `StopOnNoConnectEnableStatus` with values `StopOnNoConnectEnableStatusDisabled`, `StopOnNoConnectEnableStatusEnabled`
- New enum type `UserCustomizationsEnableStatus` with values `UserCustomizationsEnableStatusDisabled`, `UserCustomizationsEnableStatusEnabled`
- New enum type `WorkspaceStorageMode` with values `WorkspaceStorageModeAutoDeploy`, `WorkspaceStorageModeDisabled`
- New function `NewCatalogImageDefinitionBuildClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CatalogImageDefinitionBuildClient, error)`
- New function `*CatalogImageDefinitionBuildClient.BeginCancel(ctx context.Context, resourceGroupName string, devCenterName string, catalogName string, imageDefinitionName string, buildName string, options *CatalogImageDefinitionBuildClientBeginCancelOptions) (*runtime.Poller[CatalogImageDefinitionBuildClientCancelResponse], error)`
- New function `*CatalogImageDefinitionBuildClient.Get(ctx context.Context, resourceGroupName string, devCenterName string, catalogName string, imageDefinitionName string, buildName string, options *CatalogImageDefinitionBuildClientGetOptions) (CatalogImageDefinitionBuildClientGetResponse, error)`
- New function `*CatalogImageDefinitionBuildClient.GetBuildDetails(ctx context.Context, resourceGroupName string, devCenterName string, catalogName string, imageDefinitionName string, buildName string, options *CatalogImageDefinitionBuildClientGetBuildDetailsOptions) (CatalogImageDefinitionBuildClientGetBuildDetailsResponse, error)`
- New function `NewCatalogImageDefinitionBuildsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CatalogImageDefinitionBuildsClient, error)`
- New function `*CatalogImageDefinitionBuildsClient.NewListByImageDefinitionPager(resourceGroupName string, devCenterName string, catalogName string, imageDefinitionName string, options *CatalogImageDefinitionBuildsClientListByImageDefinitionOptions) *runtime.Pager[CatalogImageDefinitionBuildsClientListByImageDefinitionResponse]`
- New function `NewCatalogImageDefinitionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CatalogImageDefinitionsClient, error)`
- New function `*CatalogImageDefinitionsClient.GetByDevCenterCatalog(ctx context.Context, resourceGroupName string, devCenterName string, catalogName string, imageDefinitionName string, options *CatalogImageDefinitionsClientGetByDevCenterCatalogOptions) (CatalogImageDefinitionsClientGetByDevCenterCatalogResponse, error)`
- New function `*CatalogImageDefinitionsClient.GetErrorDetails(ctx context.Context, resourceGroupName string, devCenterName string, catalogName string, imageDefinitionName string, options *CatalogImageDefinitionsClientGetErrorDetailsOptions) (CatalogImageDefinitionsClientGetErrorDetailsResponse, error)`
- New function `*CatalogImageDefinitionsClient.NewListByDevCenterCatalogPager(resourceGroupName string, devCenterName string, catalogName string, options *CatalogImageDefinitionsClientListByDevCenterCatalogOptions) *runtime.Pager[CatalogImageDefinitionsClientListByDevCenterCatalogResponse]`
- New function `*CatalogImageDefinitionsClient.BeginBuildImage(ctx context.Context, resourceGroupName string, devCenterName string, catalogName string, imageDefinitionName string, options *CatalogImageDefinitionsClientBeginBuildImageOptions) (*runtime.Poller[CatalogImageDefinitionsClientBuildImageResponse], error)`
- New function `*ClientFactory.NewCatalogImageDefinitionBuildClient() *CatalogImageDefinitionBuildClient`
- New function `*ClientFactory.NewCatalogImageDefinitionBuildsClient() *CatalogImageDefinitionBuildsClient`
- New function `*ClientFactory.NewCatalogImageDefinitionsClient() *CatalogImageDefinitionsClient`
- New function `*ClientFactory.NewCustomizationTasksClient() *CustomizationTasksClient`
- New function `*ClientFactory.NewEncryptionSetsClient() *EncryptionSetsClient`
- New function `*ClientFactory.NewProjectCatalogImageDefinitionBuildClient() *ProjectCatalogImageDefinitionBuildClient`
- New function `*ClientFactory.NewProjectCatalogImageDefinitionBuildsClient() *ProjectCatalogImageDefinitionBuildsClient`
- New function `*ClientFactory.NewProjectCatalogImageDefinitionsClient() *ProjectCatalogImageDefinitionsClient`
- New function `*ClientFactory.NewProjectPoliciesClient() *ProjectPoliciesClient`
- New function `NewCustomizationTasksClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CustomizationTasksClient, error)`
- New function `*CustomizationTasksClient.Get(ctx context.Context, resourceGroupName string, devCenterName string, catalogName string, taskName string, options *CustomizationTasksClientGetOptions) (CustomizationTasksClientGetResponse, error)`
- New function `*CustomizationTasksClient.GetErrorDetails(ctx context.Context, resourceGroupName string, devCenterName string, catalogName string, taskName string, options *CustomizationTasksClientGetErrorDetailsOptions) (CustomizationTasksClientGetErrorDetailsResponse, error)`
- New function `*CustomizationTasksClient.NewListByCatalogPager(resourceGroupName string, devCenterName string, catalogName string, options *CustomizationTasksClientListByCatalogOptions) *runtime.Pager[CustomizationTasksClientListByCatalogResponse]`
- New function `NewEncryptionSetsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*EncryptionSetsClient, error)`
- New function `*EncryptionSetsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, devCenterName string, encryptionSetName string, body EncryptionSet, options *EncryptionSetsClientBeginCreateOrUpdateOptions) (*runtime.Poller[EncryptionSetsClientCreateOrUpdateResponse], error)`
- New function `*EncryptionSetsClient.BeginDelete(ctx context.Context, resourceGroupName string, devCenterName string, encryptionSetName string, options *EncryptionSetsClientBeginDeleteOptions) (*runtime.Poller[EncryptionSetsClientDeleteResponse], error)`
- New function `*EncryptionSetsClient.Get(ctx context.Context, resourceGroupName string, devCenterName string, encryptionSetName string, options *EncryptionSetsClientGetOptions) (EncryptionSetsClientGetResponse, error)`
- New function `*EncryptionSetsClient.NewListPager(resourceGroupName string, devCenterName string, options *EncryptionSetsClientListOptions) *runtime.Pager[EncryptionSetsClientListResponse]`
- New function `*EncryptionSetsClient.BeginUpdate(ctx context.Context, resourceGroupName string, devCenterName string, encryptionSetName string, body EncryptionSetUpdate, options *EncryptionSetsClientBeginUpdateOptions) (*runtime.Poller[EncryptionSetsClientUpdateResponse], error)`
- New function `*ImageVersionsClient.GetByProject(ctx context.Context, resourceGroupName string, projectName string, imageName string, versionName string, options *ImageVersionsClientGetByProjectOptions) (ImageVersionsClientGetByProjectResponse, error)`
- New function `*ImageVersionsClient.NewListByProjectPager(resourceGroupName string, projectName string, imageName string, options *ImageVersionsClientListByProjectOptions) *runtime.Pager[ImageVersionsClientListByProjectResponse]`
- New function `*ImagesClient.GetByProject(ctx context.Context, resourceGroupName string, projectName string, imageName string, options *ImagesClientGetByProjectOptions) (ImagesClientGetByProjectResponse, error)`
- New function `*ImagesClient.NewListByProjectPager(resourceGroupName string, projectName string, options *ImagesClientListByProjectOptions) *runtime.Pager[ImagesClientListByProjectResponse]`
- New function `NewProjectCatalogImageDefinitionBuildClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ProjectCatalogImageDefinitionBuildClient, error)`
- New function `*ProjectCatalogImageDefinitionBuildClient.BeginCancel(ctx context.Context, resourceGroupName string, projectName string, catalogName string, imageDefinitionName string, buildName string, options *ProjectCatalogImageDefinitionBuildClientBeginCancelOptions) (*runtime.Poller[ProjectCatalogImageDefinitionBuildClientCancelResponse], error)`
- New function `*ProjectCatalogImageDefinitionBuildClient.Get(ctx context.Context, resourceGroupName string, projectName string, catalogName string, imageDefinitionName string, buildName string, options *ProjectCatalogImageDefinitionBuildClientGetOptions) (ProjectCatalogImageDefinitionBuildClientGetResponse, error)`
- New function `*ProjectCatalogImageDefinitionBuildClient.GetBuildDetails(ctx context.Context, resourceGroupName string, projectName string, catalogName string, imageDefinitionName string, buildName string, options *ProjectCatalogImageDefinitionBuildClientGetBuildDetailsOptions) (ProjectCatalogImageDefinitionBuildClientGetBuildDetailsResponse, error)`
- New function `NewProjectCatalogImageDefinitionBuildsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ProjectCatalogImageDefinitionBuildsClient, error)`
- New function `*ProjectCatalogImageDefinitionBuildsClient.NewListByImageDefinitionPager(resourceGroupName string, projectName string, catalogName string, imageDefinitionName string, options *ProjectCatalogImageDefinitionBuildsClientListByImageDefinitionOptions) *runtime.Pager[ProjectCatalogImageDefinitionBuildsClientListByImageDefinitionResponse]`
- New function `NewProjectCatalogImageDefinitionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ProjectCatalogImageDefinitionsClient, error)`
- New function `*ProjectCatalogImageDefinitionsClient.GetByProjectCatalog(ctx context.Context, resourceGroupName string, projectName string, catalogName string, imageDefinitionName string, options *ProjectCatalogImageDefinitionsClientGetByProjectCatalogOptions) (ProjectCatalogImageDefinitionsClientGetByProjectCatalogResponse, error)`
- New function `*ProjectCatalogImageDefinitionsClient.GetErrorDetails(ctx context.Context, resourceGroupName string, projectName string, catalogName string, imageDefinitionName string, options *ProjectCatalogImageDefinitionsClientGetErrorDetailsOptions) (ProjectCatalogImageDefinitionsClientGetErrorDetailsResponse, error)`
- New function `*ProjectCatalogImageDefinitionsClient.NewListByProjectCatalogPager(resourceGroupName string, projectName string, catalogName string, options *ProjectCatalogImageDefinitionsClientListByProjectCatalogOptions) *runtime.Pager[ProjectCatalogImageDefinitionsClientListByProjectCatalogResponse]`
- New function `*ProjectCatalogImageDefinitionsClient.BeginBuildImage(ctx context.Context, resourceGroupName string, projectName string, catalogName string, imageDefinitionName string, options *ProjectCatalogImageDefinitionsClientBeginBuildImageOptions) (*runtime.Poller[ProjectCatalogImageDefinitionsClientBuildImageResponse], error)`
- New function `NewProjectPoliciesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ProjectPoliciesClient, error)`
- New function `*ProjectPoliciesClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, devCenterName string, projectPolicyName string, body ProjectPolicy, options *ProjectPoliciesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ProjectPoliciesClientCreateOrUpdateResponse], error)`
- New function `*ProjectPoliciesClient.BeginDelete(ctx context.Context, resourceGroupName string, devCenterName string, projectPolicyName string, options *ProjectPoliciesClientBeginDeleteOptions) (*runtime.Poller[ProjectPoliciesClientDeleteResponse], error)`
- New function `*ProjectPoliciesClient.Get(ctx context.Context, resourceGroupName string, devCenterName string, projectPolicyName string, options *ProjectPoliciesClientGetOptions) (ProjectPoliciesClientGetResponse, error)`
- New function `*ProjectPoliciesClient.NewListByDevCenterPager(resourceGroupName string, devCenterName string, options *ProjectPoliciesClientListByDevCenterOptions) *runtime.Pager[ProjectPoliciesClientListByDevCenterResponse]`
- New function `*ProjectPoliciesClient.BeginUpdate(ctx context.Context, resourceGroupName string, devCenterName string, projectPolicyName string, body ProjectPolicyUpdate, options *ProjectPoliciesClientBeginUpdateOptions) (*runtime.Poller[ProjectPoliciesClientUpdateResponse], error)`
- New function `*ProjectsClient.GetInheritedSettings(ctx context.Context, resourceGroupName string, projectName string, options *ProjectsClientGetInheritedSettingsOptions) (ProjectsClientGetInheritedSettingsResponse, error)`
- New function `*SKUsClient.NewListByProjectPager(resourceGroupName string, projectName string, options *SKUsClientListByProjectOptions) *runtime.Pager[SKUsClientListByProjectResponse]`
- New struct `ActiveHoursConfiguration`
- New struct `AzureAiServicesSettings`
- New struct `CustomizationTask`
- New struct `CustomizationTaskInput`
- New struct `CustomizationTaskInstance`
- New struct `CustomizationTaskListResult`
- New struct `CustomizationTaskProperties`
- New struct `DefinitionParametersItem`
- New struct `DevBoxAutoDeleteSettings`
- New struct `DevBoxProvisioningSettings`
- New struct `EncryptionSet`
- New struct `EncryptionSetListResult`
- New struct `EncryptionSetProperties`
- New struct `EncryptionSetUpdate`
- New struct `EncryptionSetUpdateProperties`
- New struct `ImageCreationErrorDetails`
- New struct `ImageDefinition`
- New struct `ImageDefinitionBuild`
- New struct `ImageDefinitionBuildDetails`
- New struct `ImageDefinitionBuildListResult`
- New struct `ImageDefinitionBuildProperties`
- New struct `ImageDefinitionBuildTask`
- New struct `ImageDefinitionBuildTaskGroup`
- New struct `ImageDefinitionBuildTaskParametersItem`
- New struct `ImageDefinitionListResult`
- New struct `ImageDefinitionProperties`
- New struct `ImageDefinitionReference`
- New struct `InheritedSettingsForProject`
- New struct `KeyEncryptionKeyIdentity`
- New struct `LatestImageBuild`
- New struct `NetworkSettings`
- New struct `PoolDevBoxDefinition`
- New struct `ProjectCustomizationManagedIdentity`
- New struct `ProjectCustomizationSettings`
- New struct `ProjectNetworkSettings`
- New struct `ProjectPolicy`
- New struct `ProjectPolicyListResult`
- New struct `ProjectPolicyProperties`
- New struct `ProjectPolicyUpdate`
- New struct `ProjectPolicyUpdateProperties`
- New struct `ResourcePolicy`
- New struct `ServerlessGpuSessionsSettings`
- New struct `StopOnNoConnectConfiguration`
- New struct `WorkspaceStorageSettings`
- New field `ActiveHoursConfiguration`, `DevBoxDefinition`, `DevBoxDefinitionType`, `DevBoxTunnelEnableStatus`, `StopOnNoConnect` in struct `PoolProperties`
- New field `ActiveHoursConfiguration`, `DevBoxDefinition`, `DevBoxDefinitionType`, `DevBoxTunnelEnableStatus`, `StopOnNoConnect` in struct `PoolUpdateProperties`
- New field `AzureAiServicesSettings`, `CustomizationSettings`, `DevBoxAutoDeleteSettings`, `ServerlessGpuSessionsSettings`, `WorkspaceStorageSettings` in struct `ProjectProperties`
- New field `AzureAiServicesSettings`, `CustomizationSettings`, `DevBoxAutoDeleteSettings`, `ServerlessGpuSessionsSettings`, `WorkspaceStorageSettings` in struct `ProjectUpdateProperties`
- New field `DevBoxProvisioningSettings`, `NetworkSettings` in struct `Properties`
- New field `DevBoxProvisioningSettings`, `NetworkSettings` in struct `UpdateProperties`


## 2.0.0 (2024-04-26)
### Breaking Changes

- Field `Tags` of struct `CatalogUpdate` has been removed
- Field `Location`, `Tags` of struct `ScheduleUpdate` has been removed

### Features Added

- New enum type `CatalogConnectionState` with values `CatalogConnectionStateConnected`, `CatalogConnectionStateDisconnected`
- New enum type `CatalogItemSyncEnableStatus` with values `CatalogItemSyncEnableStatusDisabled`, `CatalogItemSyncEnableStatusEnabled`
- New enum type `CatalogItemType` with values `CatalogItemTypeEnvironmentDefinition`
- New enum type `CatalogResourceValidationStatus` with values `CatalogResourceValidationStatusFailed`, `CatalogResourceValidationStatusPending`, `CatalogResourceValidationStatusSucceeded`, `CatalogResourceValidationStatusUnknown`
- New enum type `CatalogSyncType` with values `CatalogSyncTypeManual`, `CatalogSyncTypeScheduled`
- New enum type `IdentityType` with values `IdentityTypeDelegatedResourceIdentity`, `IdentityTypeSystemAssignedIdentity`, `IdentityTypeUserAssignedIdentity`
- New enum type `ParameterType` with values `ParameterTypeArray`, `ParameterTypeBoolean`, `ParameterTypeInteger`, `ParameterTypeNumber`, `ParameterTypeObject`, `ParameterTypeString`
- New enum type `SingleSignOnStatus` with values `SingleSignOnStatusDisabled`, `SingleSignOnStatusEnabled`
- New enum type `VirtualNetworkType` with values `VirtualNetworkTypeManaged`, `VirtualNetworkTypeUnmanaged`
- New function `*CatalogsClient.BeginConnect(context.Context, string, string, string, *CatalogsClientBeginConnectOptions) (*runtime.Poller[CatalogsClientConnectResponse], error)`
- New function `*CatalogsClient.GetSyncErrorDetails(context.Context, string, string, string, *CatalogsClientGetSyncErrorDetailsOptions) (CatalogsClientGetSyncErrorDetailsResponse, error)`
- New function `NewCheckScopedNameAvailabilityClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CheckScopedNameAvailabilityClient, error)`
- New function `*CheckScopedNameAvailabilityClient.Execute(context.Context, CheckScopedNameAvailabilityRequest, *CheckScopedNameAvailabilityClientExecuteOptions) (CheckScopedNameAvailabilityClientExecuteResponse, error)`
- New function `*ClientFactory.NewCheckScopedNameAvailabilityClient() *CheckScopedNameAvailabilityClient`
- New function `*ClientFactory.NewEnvironmentDefinitionsClient() *EnvironmentDefinitionsClient`
- New function `*ClientFactory.NewProjectCatalogEnvironmentDefinitionsClient() *ProjectCatalogEnvironmentDefinitionsClient`
- New function `*ClientFactory.NewProjectCatalogsClient() *ProjectCatalogsClient`
- New function `NewEnvironmentDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EnvironmentDefinitionsClient, error)`
- New function `*EnvironmentDefinitionsClient.Get(context.Context, string, string, string, string, *EnvironmentDefinitionsClientGetOptions) (EnvironmentDefinitionsClientGetResponse, error)`
- New function `*EnvironmentDefinitionsClient.GetByProjectCatalog(context.Context, string, string, string, string, *EnvironmentDefinitionsClientGetByProjectCatalogOptions) (EnvironmentDefinitionsClientGetByProjectCatalogResponse, error)`
- New function `*EnvironmentDefinitionsClient.GetErrorDetails(context.Context, string, string, string, string, *EnvironmentDefinitionsClientGetErrorDetailsOptions) (EnvironmentDefinitionsClientGetErrorDetailsResponse, error)`
- New function `*EnvironmentDefinitionsClient.NewListByCatalogPager(string, string, string, *EnvironmentDefinitionsClientListByCatalogOptions) *runtime.Pager[EnvironmentDefinitionsClientListByCatalogResponse]`
- New function `*EnvironmentDefinitionsClient.NewListByProjectCatalogPager(string, string, string, *EnvironmentDefinitionsClientListByProjectCatalogOptions) *runtime.Pager[EnvironmentDefinitionsClientListByProjectCatalogResponse]`
- New function `NewProjectCatalogEnvironmentDefinitionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectCatalogEnvironmentDefinitionsClient, error)`
- New function `*ProjectCatalogEnvironmentDefinitionsClient.GetErrorDetails(context.Context, string, string, string, string, *ProjectCatalogEnvironmentDefinitionsClientGetErrorDetailsOptions) (ProjectCatalogEnvironmentDefinitionsClientGetErrorDetailsResponse, error)`
- New function `NewProjectCatalogsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectCatalogsClient, error)`
- New function `*ProjectCatalogsClient.BeginConnect(context.Context, string, string, string, *ProjectCatalogsClientBeginConnectOptions) (*runtime.Poller[ProjectCatalogsClientConnectResponse], error)`
- New function `*ProjectCatalogsClient.BeginCreateOrUpdate(context.Context, string, string, string, Catalog, *ProjectCatalogsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ProjectCatalogsClientCreateOrUpdateResponse], error)`
- New function `*ProjectCatalogsClient.BeginDelete(context.Context, string, string, string, *ProjectCatalogsClientBeginDeleteOptions) (*runtime.Poller[ProjectCatalogsClientDeleteResponse], error)`
- New function `*ProjectCatalogsClient.Get(context.Context, string, string, string, *ProjectCatalogsClientGetOptions) (ProjectCatalogsClientGetResponse, error)`
- New function `*ProjectCatalogsClient.GetSyncErrorDetails(context.Context, string, string, string, *ProjectCatalogsClientGetSyncErrorDetailsOptions) (ProjectCatalogsClientGetSyncErrorDetailsResponse, error)`
- New function `*ProjectCatalogsClient.NewListPager(string, string, *ProjectCatalogsClientListOptions) *runtime.Pager[ProjectCatalogsClientListResponse]`
- New function `*ProjectCatalogsClient.BeginPatch(context.Context, string, string, string, CatalogUpdate, *ProjectCatalogsClientBeginPatchOptions) (*runtime.Poller[ProjectCatalogsClientPatchResponse], error)`
- New function `*ProjectCatalogsClient.BeginSync(context.Context, string, string, string, *ProjectCatalogsClientBeginSyncOptions) (*runtime.Poller[ProjectCatalogsClientSyncResponse], error)`
- New struct `CatalogConflictError`
- New struct `CatalogErrorDetails`
- New struct `CatalogResourceValidationErrorDetails`
- New struct `CatalogSyncError`
- New struct `CheckScopedNameAvailabilityRequest`
- New struct `CustomerManagedKeyEncryption`
- New struct `CustomerManagedKeyEncryptionKeyIdentity`
- New struct `Encryption`
- New struct `EnvironmentDefinition`
- New struct `EnvironmentDefinitionListResult`
- New struct `EnvironmentDefinitionParameter`
- New struct `EnvironmentDefinitionProperties`
- New struct `EnvironmentTypeUpdateProperties`
- New struct `ProjectCatalogSettings`
- New struct `ProjectCatalogSettingsInfo`
- New struct `SyncErrorDetails`
- New struct `SyncStats`
- New struct `UpdateProperties`
- New field `DisplayName` in struct `AllowedEnvironmentTypeProperties`
- New field `ConnectionState`, `LastConnectionTime`, `LastSyncStats`, `SyncType`, `Tags` in struct `CatalogProperties`
- New field `SyncType`, `Tags` in struct `CatalogUpdateProperties`
- New field `ValidationStatus` in struct `DevBoxDefinitionProperties`
- New field `DisplayName` in struct `EnvironmentTypeProperties`
- New field `Properties` in struct `EnvironmentTypeUpdate`
- New field `ResourceID` in struct `OperationStatusResult`
- New field `Location` in struct `OperationStatusesClientGetResponse`
- New field `DevBoxCount`, `DisplayName`, `ManagedVirtualNetworkRegions`, `SingleSignOnStatus`, `VirtualNetworkType` in struct `PoolProperties`
- New field `DisplayName`, `ManagedVirtualNetworkRegions`, `SingleSignOnStatus`, `VirtualNetworkType` in struct `PoolUpdateProperties`
- New field `Identity` in struct `Project`
- New field `DisplayName`, `EnvironmentCount` in struct `ProjectEnvironmentTypeProperties`
- New field `DisplayName` in struct `ProjectEnvironmentTypeUpdateProperties`
- New field `CatalogSettings`, `DisplayName` in struct `ProjectProperties`
- New field `Identity` in struct `ProjectUpdate`
- New field `CatalogSettings`, `DisplayName` in struct `ProjectUpdateProperties`
- New field `DisplayName`, `Encryption`, `ProjectCatalogSettings` in struct `Properties`
- New field `Location`, `Tags` in struct `ScheduleProperties`
- New field `Location`, `Tags` in struct `ScheduleUpdateProperties`
- New field `Properties` in struct `Update`
- New field `ID` in struct `Usage`


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