# Release History

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
