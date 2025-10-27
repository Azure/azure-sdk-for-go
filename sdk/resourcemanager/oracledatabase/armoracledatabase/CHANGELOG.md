# Release History

## 2.0.0 (2025-09-23)
### Breaking Changes

- Field `ScheduledOperations` of struct `AutonomousDatabaseBaseProperties` has been removed
- Field `ScheduledOperations` of struct `AutonomousDatabaseCloneProperties` has been removed
- Field `ScheduledOperations` of struct `AutonomousDatabaseCrossRegionDisasterRecoveryProperties` has been removed
- Field `ScheduledOperations` of struct `AutonomousDatabaseFromBackupTimestampProperties` has been removed
- Field `ScheduledOperations` of struct `AutonomousDatabaseProperties` has been removed
- Field `ScheduledOperations` of struct `AutonomousDatabaseUpdateProperties` has been removed

### Features Added

- New enum type `AutonomousDatabaseLifecycleActionEnum` with values `AutonomousDatabaseLifecycleActionEnumRestart`, `AutonomousDatabaseLifecycleActionEnumStart`, `AutonomousDatabaseLifecycleActionEnumStop`
- New enum type `BaseDbSystemShapes` with values `BaseDbSystemShapesVMStandardX86`
- New enum type `DbSystemDatabaseEditionType` with values `DbSystemDatabaseEditionTypeEnterpriseEdition`, `DbSystemDatabaseEditionTypeEnterpriseEditionDeveloper`, `DbSystemDatabaseEditionTypeEnterpriseEditionExtreme`, `DbSystemDatabaseEditionTypeEnterpriseEditionHighPerformance`, `DbSystemDatabaseEditionTypeStandardEdition`
- New enum type `DbSystemLifecycleState` with values `DbSystemLifecycleStateAvailable`, `DbSystemLifecycleStateFailed`, `DbSystemLifecycleStateMaintenanceInProgress`, `DbSystemLifecycleStateMigrated`, `DbSystemLifecycleStateNeedsAttention`, `DbSystemLifecycleStateProvisioning`, `DbSystemLifecycleStateTerminated`, `DbSystemLifecycleStateTerminating`, `DbSystemLifecycleStateUpdating`, `DbSystemLifecycleStateUpgrading`
- New enum type `DbSystemSourceType` with values `DbSystemSourceTypeNone`
- New enum type `DiskRedundancyType` with values `DiskRedundancyTypeHigh`, `DiskRedundancyTypeNormal`
- New enum type `ExadataVMClusterStorageManagementType` with values `ExadataVMClusterStorageManagementTypeASM`, `ExadataVMClusterStorageManagementTypeExascale`
- New enum type `ShapeAttribute` with values `ShapeAttributeBLOCKSTORAGE`, `ShapeAttributeSMARTSTORAGE`
- New enum type `ShapeFamilyType` with values `ShapeFamilyTypeExadata`, `ShapeFamilyTypeExadbXs`, `ShapeFamilyTypeSingleNode`, `ShapeFamilyTypeVirtualMachine`
- New enum type `StorageManagementType` with values `StorageManagementTypeLVM`
- New enum type `StorageVolumePerformanceMode` with values `StorageVolumePerformanceModeBalanced`, `StorageVolumePerformanceModeHighPerformance`
- New function `*AutonomousDatabasesClient.BeginAction(context.Context, string, string, AutonomousDatabaseLifecycleAction, *AutonomousDatabasesClientBeginActionOptions) (*runtime.Poller[AutonomousDatabasesClientActionResponse], error)`
- New function `*ClientFactory.NewDbSystemsClient() *DbSystemsClient`
- New function `*ClientFactory.NewDbVersionsClient() *DbVersionsClient`
- New function `*ClientFactory.NewNetworkAnchorsClient() *NetworkAnchorsClient`
- New function `*ClientFactory.NewResourceAnchorsClient() *ResourceAnchorsClient`
- New function `*CloudExadataInfrastructuresClient.BeginConfigureExascale(context.Context, string, string, ConfigureExascaleCloudExadataInfrastructureDetails, *CloudExadataInfrastructuresClientBeginConfigureExascaleOptions) (*runtime.Poller[CloudExadataInfrastructuresClientConfigureExascaleResponse], error)`
- New function `*DbSystemBaseProperties.GetDbSystemBaseProperties() *DbSystemBaseProperties`
- New function `*DbSystemProperties.GetDbSystemBaseProperties() *DbSystemBaseProperties`
- New function `NewDbSystemsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DbSystemsClient, error)`
- New function `*DbSystemsClient.BeginCreateOrUpdate(context.Context, string, string, DbSystem, *DbSystemsClientBeginCreateOrUpdateOptions) (*runtime.Poller[DbSystemsClientCreateOrUpdateResponse], error)`
- New function `*DbSystemsClient.BeginDelete(context.Context, string, string, *DbSystemsClientBeginDeleteOptions) (*runtime.Poller[DbSystemsClientDeleteResponse], error)`
- New function `*DbSystemsClient.Get(context.Context, string, string, *DbSystemsClientGetOptions) (DbSystemsClientGetResponse, error)`
- New function `*DbSystemsClient.NewListByResourceGroupPager(string, *DbSystemsClientListByResourceGroupOptions) *runtime.Pager[DbSystemsClientListByResourceGroupResponse]`
- New function `*DbSystemsClient.NewListBySubscriptionPager(*DbSystemsClientListBySubscriptionOptions) *runtime.Pager[DbSystemsClientListBySubscriptionResponse]`
- New function `*DbSystemsClient.BeginUpdate(context.Context, string, string, DbSystemUpdate, *DbSystemsClientBeginUpdateOptions) (*runtime.Poller[DbSystemsClientUpdateResponse], error)`
- New function `NewDbVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DbVersionsClient, error)`
- New function `*DbVersionsClient.Get(context.Context, string, string, *DbVersionsClientGetOptions) (DbVersionsClientGetResponse, error)`
- New function `*DbVersionsClient.NewListByLocationPager(string, *DbVersionsClientListByLocationOptions) *runtime.Pager[DbVersionsClientListByLocationResponse]`
- New function `NewResourceAnchorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ResourceAnchorsClient, error)`
- New function `*ResourceAnchorsClient.BeginCreateOrUpdate(context.Context, string, string, ResourceAnchor, *ResourceAnchorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ResourceAnchorsClientCreateOrUpdateResponse], error)`
- New function `*ResourceAnchorsClient.BeginDelete(context.Context, string, string, *ResourceAnchorsClientBeginDeleteOptions) (*runtime.Poller[ResourceAnchorsClientDeleteResponse], error)`
- New function `*ResourceAnchorsClient.Get(context.Context, string, string, *ResourceAnchorsClientGetOptions) (ResourceAnchorsClientGetResponse, error)`
- New function `*ResourceAnchorsClient.NewListByResourceGroupPager(string, *ResourceAnchorsClientListByResourceGroupOptions) *runtime.Pager[ResourceAnchorsClientListByResourceGroupResponse]`
- New function `*ResourceAnchorsClient.NewListBySubscriptionPager(*ResourceAnchorsClientListBySubscriptionOptions) *runtime.Pager[ResourceAnchorsClientListBySubscriptionResponse]`
- New function `*ResourceAnchorsClient.BeginUpdate(context.Context, string, string, ResourceAnchorUpdate, *ResourceAnchorsClientBeginUpdateOptions) (*runtime.Poller[ResourceAnchorsClientUpdateResponse], error)`
- New function `NewNetworkAnchorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkAnchorsClient, error)`
- New function `*NetworkAnchorsClient.BeginCreateOrUpdate(context.Context, string, string, NetworkAnchor, *NetworkAnchorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[NetworkAnchorsClientCreateOrUpdateResponse], error)`
- New function `*NetworkAnchorsClient.BeginDelete(context.Context, string, string, *NetworkAnchorsClientBeginDeleteOptions) (*runtime.Poller[NetworkAnchorsClientDeleteResponse], error)`
- New function `*NetworkAnchorsClient.Get(context.Context, string, string, *NetworkAnchorsClientGetOptions) (NetworkAnchorsClientGetResponse, error)`
- New function `*NetworkAnchorsClient.NewListByResourceGroupPager(string, *NetworkAnchorsClientListByResourceGroupOptions) *runtime.Pager[NetworkAnchorsClientListByResourceGroupResponse]`
- New function `*NetworkAnchorsClient.NewListBySubscriptionPager(*NetworkAnchorsClientListBySubscriptionOptions) *runtime.Pager[NetworkAnchorsClientListBySubscriptionResponse]`
- New function `*NetworkAnchorsClient.BeginUpdate(context.Context, string, string, NetworkAnchorUpdate, *NetworkAnchorsClientBeginUpdateOptions) (*runtime.Poller[NetworkAnchorsClientUpdateResponse], error)`
- New struct `AutonomousDatabaseLifecycleAction`
- New struct `ConfigureExascaleCloudExadataInfrastructureDetails`
- New struct `DNSForwardingRule`
- New struct `DbSystem`
- New struct `DbSystemListResult`
- New struct `DbSystemOptions`
- New struct `DbSystemProperties`
- New struct `DbSystemUpdate`
- New struct `DbSystemUpdateProperties`
- New struct `DbVersion`
- New struct `DbVersionListResult`
- New struct `DbVersionProperties`
- New struct `ExascaleConfigDetails`
- New struct `NetworkAnchor`
- New struct `NetworkAnchorListResult`
- New struct `NetworkAnchorProperties`
- New struct `NetworkAnchorUpdate`
- New struct `NetworkAnchorUpdateProperties`
- New struct `ResourceAnchor`
- New struct `ResourceAnchorListResult`
- New struct `ResourceAnchorProperties`
- New struct `ResourceAnchorUpdate`
- New field `ScheduledOperationsList` in struct `AutonomousDatabaseBaseProperties`
- New field `ScheduledOperationsList` in struct `AutonomousDatabaseCloneProperties`
- New field `ScheduledOperationsList` in struct `AutonomousDatabaseCrossRegionDisasterRecoveryProperties`
- New field `ScheduledOperationsList` in struct `AutonomousDatabaseFromBackupTimestampProperties`
- New field `ScheduledOperationsList` in struct `AutonomousDatabaseProperties`
- New field `ScheduledOperationsList` in struct `AutonomousDatabaseUpdateProperties`
- New field `ExascaleConfig` in struct `CloudExadataInfrastructureProperties`
- New field `ExascaleDbStorageVaultID`, `StorageManagementType` in struct `CloudVMClusterProperties`
- New field `ShapeAttributes` in struct `DbSystemShapeProperties`
- New field `ShapeAttribute` in struct `DbSystemShapesClientListByLocationOptions`
- New field `ShapeAttribute` in struct `ExadbVMClusterProperties`
- New field `AttachedShapeAttributes`, `ExadataInfrastructureID` in struct `ExascaleDbStorageVaultProperties`
- New field `ShapeAttribute` in struct `GiVersionsClientListByLocationOptions`


## 1.1.0 (2025-06-04)

### Features Added

- New value `DataBaseTypeCloneFromBackupTimestamp`, `DataBaseTypeCrossRegionDisasterRecovery` added to enum type `DataBaseType`
- New enum type `AddSubscriptionOperationState` with values `AddSubscriptionOperationStateFailed`, `AddSubscriptionOperationStateSucceeded`, `AddSubscriptionOperationStateUpdating`
- New enum type `ExadbVMClusterLifecycleState` with values `ExadbVMClusterLifecycleStateAvailable`, `ExadbVMClusterLifecycleStateFailed`, `ExadbVMClusterLifecycleStateMaintenanceInProgress`, `ExadbVMClusterLifecycleStateProvisioning`, `ExadbVMClusterLifecycleStateTerminated`, `ExadbVMClusterLifecycleStateTerminating`, `ExadbVMClusterLifecycleStateUpdating`
- New enum type `ExascaleDbStorageVaultLifecycleState` with values `ExascaleDbStorageVaultLifecycleStateAvailable`, `ExascaleDbStorageVaultLifecycleStateFailed`, `ExascaleDbStorageVaultLifecycleStateProvisioning`, `ExascaleDbStorageVaultLifecycleStateTerminated`, `ExascaleDbStorageVaultLifecycleStateTerminating`, `ExascaleDbStorageVaultLifecycleStateUpdating`
- New enum type `GridImageType` with values `GridImageTypeCustomImage`, `GridImageTypeReleaseUpdate`
- New enum type `HardwareType` with values `HardwareTypeCELL`, `HardwareTypeCOMPUTE`
- New enum type `ShapeFamily` with values `ShapeFamilyExadata`, `ShapeFamilyExadbXs`
- New enum type `SystemShapes` with values `SystemShapesExaDbXs`, `SystemShapesExadataX11M`, `SystemShapesExadataX9M`
- New function `*AutonomousDatabaseCrossRegionDisasterRecoveryProperties.GetAutonomousDatabaseBaseProperties() *AutonomousDatabaseBaseProperties`
- New function `*AutonomousDatabaseFromBackupTimestampProperties.GetAutonomousDatabaseBaseProperties() *AutonomousDatabaseBaseProperties`
- New function `*AutonomousDatabasesClient.BeginChangeDisasterRecoveryConfiguration(context.Context, string, string, DisasterRecoveryConfigurationDetails, *AutonomousDatabasesClientBeginChangeDisasterRecoveryConfigurationOptions) (*runtime.Poller[AutonomousDatabasesClientChangeDisasterRecoveryConfigurationResponse], error)`
- New function `*ClientFactory.NewExadbVMClustersClient() *ExadbVMClustersClient`
- New function `*ClientFactory.NewExascaleDbNodesClient() *ExascaleDbNodesClient`
- New function `*ClientFactory.NewExascaleDbStorageVaultsClient() *ExascaleDbStorageVaultsClient`
- New function `*ClientFactory.NewFlexComponentsClient() *FlexComponentsClient`
- New function `*ClientFactory.NewGiMinorVersionsClient() *GiMinorVersionsClient`
- New function `NewExadbVMClustersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExadbVMClustersClient, error)`
- New function `*ExadbVMClustersClient.BeginCreateOrUpdate(context.Context, string, string, ExadbVMCluster, *ExadbVMClustersClientBeginCreateOrUpdateOptions) (*runtime.Poller[ExadbVMClustersClientCreateOrUpdateResponse], error)`
- New function `*ExadbVMClustersClient.BeginDelete(context.Context, string, string, *ExadbVMClustersClientBeginDeleteOptions) (*runtime.Poller[ExadbVMClustersClientDeleteResponse], error)`
- New function `*ExadbVMClustersClient.Get(context.Context, string, string, *ExadbVMClustersClientGetOptions) (ExadbVMClustersClientGetResponse, error)`
- New function `*ExadbVMClustersClient.NewListByResourceGroupPager(string, *ExadbVMClustersClientListByResourceGroupOptions) *runtime.Pager[ExadbVMClustersClientListByResourceGroupResponse]`
- New function `*ExadbVMClustersClient.NewListBySubscriptionPager(*ExadbVMClustersClientListBySubscriptionOptions) *runtime.Pager[ExadbVMClustersClientListBySubscriptionResponse]`
- New function `*ExadbVMClustersClient.BeginRemoveVMs(context.Context, string, string, RemoveVirtualMachineFromExadbVMClusterDetails, *ExadbVMClustersClientBeginRemoveVMsOptions) (*runtime.Poller[ExadbVMClustersClientRemoveVMsResponse], error)`
- New function `*ExadbVMClustersClient.BeginUpdate(context.Context, string, string, ExadbVMClusterUpdate, *ExadbVMClustersClientBeginUpdateOptions) (*runtime.Poller[ExadbVMClustersClientUpdateResponse], error)`
- New function `NewExascaleDbNodesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExascaleDbNodesClient, error)`
- New function `*ExascaleDbNodesClient.BeginAction(context.Context, string, string, string, DbNodeAction, *ExascaleDbNodesClientBeginActionOptions) (*runtime.Poller[ExascaleDbNodesClientActionResponse], error)`
- New function `*ExascaleDbNodesClient.Get(context.Context, string, string, string, *ExascaleDbNodesClientGetOptions) (ExascaleDbNodesClientGetResponse, error)`
- New function `*ExascaleDbNodesClient.NewListByParentPager(string, string, *ExascaleDbNodesClientListByParentOptions) *runtime.Pager[ExascaleDbNodesClientListByParentResponse]`
- New function `NewExascaleDbStorageVaultsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExascaleDbStorageVaultsClient, error)`
- New function `*ExascaleDbStorageVaultsClient.BeginCreate(context.Context, string, string, ExascaleDbStorageVault, *ExascaleDbStorageVaultsClientBeginCreateOptions) (*runtime.Poller[ExascaleDbStorageVaultsClientCreateResponse], error)`
- New function `*ExascaleDbStorageVaultsClient.BeginDelete(context.Context, string, string, *ExascaleDbStorageVaultsClientBeginDeleteOptions) (*runtime.Poller[ExascaleDbStorageVaultsClientDeleteResponse], error)`
- New function `*ExascaleDbStorageVaultsClient.Get(context.Context, string, string, *ExascaleDbStorageVaultsClientGetOptions) (ExascaleDbStorageVaultsClientGetResponse, error)`
- New function `*ExascaleDbStorageVaultsClient.NewListByResourceGroupPager(string, *ExascaleDbStorageVaultsClientListByResourceGroupOptions) *runtime.Pager[ExascaleDbStorageVaultsClientListByResourceGroupResponse]`
- New function `*ExascaleDbStorageVaultsClient.NewListBySubscriptionPager(*ExascaleDbStorageVaultsClientListBySubscriptionOptions) *runtime.Pager[ExascaleDbStorageVaultsClientListBySubscriptionResponse]`
- New function `*ExascaleDbStorageVaultsClient.BeginUpdate(context.Context, string, string, ExascaleDbStorageVaultTagsUpdate, *ExascaleDbStorageVaultsClientBeginUpdateOptions) (*runtime.Poller[ExascaleDbStorageVaultsClientUpdateResponse], error)`
- New function `NewFlexComponentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FlexComponentsClient, error)`
- New function `*FlexComponentsClient.Get(context.Context, string, string, *FlexComponentsClientGetOptions) (FlexComponentsClientGetResponse, error)`
- New function `*FlexComponentsClient.NewListByParentPager(string, *FlexComponentsClientListByParentOptions) *runtime.Pager[FlexComponentsClientListByParentResponse]`
- New function `NewGiMinorVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GiMinorVersionsClient, error)`
- New function `*GiMinorVersionsClient.Get(context.Context, string, string, string, *GiMinorVersionsClientGetOptions) (GiMinorVersionsClientGetResponse, error)`
- New function `*GiMinorVersionsClient.NewListByParentPager(string, string, *GiMinorVersionsClientListByParentOptions) *runtime.Pager[GiMinorVersionsClientListByParentResponse]`
- New function `*OracleSubscriptionsClient.BeginAddAzureSubscriptions(context.Context, AzureSubscriptions, *OracleSubscriptionsClientBeginAddAzureSubscriptionsOptions) (*runtime.Poller[OracleSubscriptionsClientAddAzureSubscriptionsResponse], error)`
- New struct `AutonomousDatabaseCrossRegionDisasterRecoveryProperties`
- New struct `AutonomousDatabaseFromBackupTimestampProperties`
- New struct `AzureSubscriptions`
- New struct `DbActionResponse`
- New struct `DbNodeDetails`
- New struct `DefinedFileSystemConfiguration`
- New struct `DisasterRecoveryConfigurationDetails`
- New struct `ExadbVMCluster`
- New struct `ExadbVMClusterListResult`
- New struct `ExadbVMClusterProperties`
- New struct `ExadbVMClusterStorageDetails`
- New struct `ExadbVMClusterUpdate`
- New struct `ExadbVMClusterUpdateProperties`
- New struct `ExascaleDbNode`
- New struct `ExascaleDbNodeListResult`
- New struct `ExascaleDbNodeProperties`
- New struct `ExascaleDbStorageDetails`
- New struct `ExascaleDbStorageInputDetails`
- New struct `ExascaleDbStorageVault`
- New struct `ExascaleDbStorageVaultListResult`
- New struct `ExascaleDbStorageVaultProperties`
- New struct `ExascaleDbStorageVaultTagsUpdate`
- New struct `FileSystemConfigurationDetails`
- New struct `FlexComponent`
- New struct `FlexComponentListResult`
- New struct `FlexComponentProperties`
- New struct `GiMinorVersion`
- New struct `GiMinorVersionListResult`
- New struct `GiMinorVersionProperties`
- New struct `RemoveVirtualMachineFromExadbVMClusterDetails`
- New field `RemoteDisasterRecoveryConfiguration`, `TimeDisasterRecoveryRoleChanged` in struct `AutonomousDatabaseCloneProperties`
- New field `RemoteDisasterRecoveryConfiguration`, `TimeDisasterRecoveryRoleChanged` in struct `AutonomousDatabaseProperties`
- New field `ComputeModel`, `DatabaseServerType`, `DefinedFileSystemConfiguration`, `StorageServerType` in struct `CloudExadataInfrastructureProperties`
- New field `ComputeModel`, `FileSystemConfigurationDetails` in struct `CloudVMClusterProperties`
- New field `FileSystemConfigurationDetails` in struct `CloudVMClusterUpdateProperties`
- New field `ComputeModel` in struct `DbServerProperties`
- New field `AreServerTypesSupported`, `ComputeModel`, `DisplayName`, `ShapeName` in struct `DbSystemShapeProperties`
- New field `Zone` in struct `DbSystemShapesClientListByLocationOptions`
- New field `Shape`, `Zone` in struct `GiVersionsClientListByLocationOptions`
- New field `AddSubscriptionOperationState`, `AzureSubscriptionIDs`, `LastOperationStatusDetail` in struct `OracleSubscriptionProperties`
- New field `PeerDbLocation`, `PeerDbOcid` in struct `PeerDbDetails`


## 1.0.0 (2024-06-28)
### Other Changes

- Release stable version.


## 0.2.0 (2024-06-26)
### Breaking Changes

- Type of `CloudExadataInfrastructureProperties.DataStorageSizeInTbs` has been changed from `*int32` to `*float64`
- Type of `CloudVMClusterProperties.NsgCidrs` has been changed from `[]*NSGCidr` to `[]*NsgCidr`
- Type of `OracleSubscriptionUpdate.Plan` has been changed from `*ResourcePlanTypeUpdate` to `*PlanUpdate`
- Struct `NSGCidr` has been removed
- Struct `ResourcePlanTypeUpdate` has been removed
- Field `AutonomousDatabaseID`, `DatabaseSizeInTBs`, `SizeInTBs`, `Type` of struct `AutonomousDatabaseBackupProperties` has been removed

### Features Added

- New enum type `RepeatCadenceType` with values `RepeatCadenceTypeMonthly`, `RepeatCadenceTypeOneTime`, `RepeatCadenceTypeWeekly`, `RepeatCadenceTypeYearly`
- New function `*AutonomousDatabasesClient.BeginRestore(context.Context, string, string, RestoreAutonomousDatabaseDetails, *AutonomousDatabasesClientBeginRestoreOptions) (*runtime.Poller[AutonomousDatabasesClientRestoreResponse], error)`
- New function `*AutonomousDatabasesClient.BeginShrink(context.Context, string, string, *AutonomousDatabasesClientBeginShrinkOptions) (*runtime.Poller[AutonomousDatabasesClientShrinkResponse], error)`
- New function `*ClientFactory.NewSystemVersionsClient() *SystemVersionsClient`
- New function `NewSystemVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SystemVersionsClient, error)`
- New function `*SystemVersionsClient.Get(context.Context, string, string, *SystemVersionsClientGetOptions) (SystemVersionsClientGetResponse, error)`
- New function `*SystemVersionsClient.NewListByLocationPager(string, *SystemVersionsClientListByLocationOptions) *runtime.Pager[SystemVersionsClientListByLocationResponse]`
- New struct `LongTermBackUpScheduleDetails`
- New struct `NsgCidr`
- New struct `PlanUpdate`
- New struct `RestoreAutonomousDatabaseDetails`
- New struct `SystemVersion`
- New struct `SystemVersionListResult`
- New struct `SystemVersionProperties`
- New field `AutonomousDatabaseOcid`, `BackupType`, `DatabaseSizeInTbs`, `SizeInTbs`, `TimeStarted` in struct `AutonomousDatabaseBackupProperties`
- New field `LongTermBackupSchedule`, `NextLongTermBackupTimeStamp` in struct `AutonomousDatabaseBaseProperties`
- New field `LongTermBackupSchedule`, `NextLongTermBackupTimeStamp` in struct `AutonomousDatabaseCloneProperties`
- New field `LongTermBackupSchedule`, `NextLongTermBackupTimeStamp` in struct `AutonomousDatabaseProperties`
- New field `LongTermBackupSchedule` in struct `AutonomousDatabaseUpdateProperties`


## 0.1.0 (2024-05-24)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/oracledatabase/armoracledatabase` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
