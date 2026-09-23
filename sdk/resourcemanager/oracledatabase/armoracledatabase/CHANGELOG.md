# Release History

## 3.0.0 (2026-09-17)
### Breaking Changes

- Field `Hostname`, `ScanDNSName` of struct `CloudVMClusterProperties` has been removed
- Field `Hostname`, `ScanDNSName` of struct `ExadbVMClusterProperties` has been removed

### Features Added

- New value `BaseDbSystemShapesVMBaseDBX86` added to enum type `BaseDbSystemShapes`
- New value `SystemShapesExadataX11MV` added to enum type `SystemShapes`
- New value `WorkloadTypeLH` added to enum type `WorkloadType`
- New enum type `BackupDestinationType` with values `BackupDestinationTypeAzure`, `BackupDestinationTypeOci`
- New enum type `CategoryType` with values `CategoryTypeDataReplication`, `CategoryTypeDataTransforms`, `CategoryTypeStreamAnalytics`
- New enum type `ConnectionLifecycleState` with values `ConnectionLifecycleStateActive`, `ConnectionLifecycleStateCreating`, `ConnectionLifecycleStateDeleted`, `ConnectionLifecycleStateDeleting`, `ConnectionLifecycleStateFailed`, `ConnectionLifecycleStateUpdating`
- New enum type `ConnectionType` with values `ConnectionTypeAmazonKinesis`, `ConnectionTypeAmazonRedshift`, `ConnectionTypeAmazonS3`, `ConnectionTypeAzureDataLakeStorage`, `ConnectionTypeAzureSynapseAnalytics`, `ConnectionTypeDatabricks`, `ConnectionTypeDb2Connection`, `ConnectionTypeElasticsearch`, `ConnectionTypeGeneric`, `ConnectionTypeGoldenGate`, `ConnectionTypeGoogleBigQuery`, `ConnectionTypeGoogleCloudStorage`, `ConnectionTypeGooglePubSub`, `ConnectionTypeHdfs`, `ConnectionTypeIceberg`, `ConnectionTypeJavaMessageService`, `ConnectionTypeKafka`, `ConnectionTypeKafkaSchemaRegistry`, `ConnectionTypeMicrosoftFabric`, `ConnectionTypeMicrosoftSQLServer`, `ConnectionTypeMongoDbConnection`, `ConnectionTypeMySQL`, `ConnectionTypeOciObjectStorage`, `ConnectionTypeOracle`, `ConnectionTypeOracleNoSQL`, `ConnectionTypePostgreSQL`, `ConnectionTypeRedis`, `ConnectionTypeSnowflake`
- New enum type `CredentialType` with values `CredentialTypeGoldenGate`, `CredentialTypeIam`
- New enum type `DeploymentLifecycleState` with values `DeploymentLifecycleStateActive`, `DeploymentLifecycleStateCanceled`, `DeploymentLifecycleStateCanceling`, `DeploymentLifecycleStateCreating`, `DeploymentLifecycleStateDeleted`, `DeploymentLifecycleStateDeleting`, `DeploymentLifecycleStateFailed`, `DeploymentLifecycleStateInActive`, `DeploymentLifecycleStateInProgress`, `DeploymentLifecycleStateNeedsAttention`, `DeploymentLifecycleStateSucceeded`, `DeploymentLifecycleStateUpdating`, `DeploymentLifecycleStateWaiting`
- New enum type `DeploymentType` with values `DeploymentTypeBigData`, `DeploymentTypeDataTransforms`, `DeploymentTypeDatabaseDB2I`, `DeploymentTypeDatabaseDB2ZOS`, `DeploymentTypeDatabaseMicrosoftSQLServer`, `DeploymentTypeDatabaseMySQL`, `DeploymentTypeDatabaseOracle`, `DeploymentTypeDatabasePostGreSQL`, `DeploymentTypeGgsa`, `DeploymentTypeOgg`
- New enum type `FrequencyType` with values `FrequencyTypeDaily`, `FrequencyTypeMonthly`, `FrequencyTypeWeekly`
- New enum type `GiMinorVersionSortOrder` with values `GiMinorVersionSortOrderAsc`, `GiMinorVersionSortOrderDesc`
- New enum type `GoldenGateConnectionAssignmentLifecycleState` with values `GoldenGateConnectionAssignmentLifecycleStateActive`, `GoldenGateConnectionAssignmentLifecycleStateCreating`, `GoldenGateConnectionAssignmentLifecycleStateDeleted`, `GoldenGateConnectionAssignmentLifecycleStateDeleting`, `GoldenGateConnectionAssignmentLifecycleStateFailed`, `GoldenGateConnectionAssignmentLifecycleStateUpdating`
- New enum type `KafkaConnectionTechnologyType` with values `KafkaConnectionTechnologyTypeApacheKafka`, `KafkaConnectionTechnologyTypeAzureEventHubs`, `KafkaConnectionTechnologyTypeConfluentKafka`, `KafkaConnectionTechnologyTypeOciStreaming`
- New enum type `MicrosoftFabricConnectionTechnologyType` with values `MicrosoftFabricConnectionTechnologyTypeMicrosoftFabricLakehouse`, `MicrosoftFabricConnectionTechnologyTypeMicrosoftFabricMirror`
- New enum type `OracleConnectionTechnologyType` with values `OracleConnectionTechnologyTypeAmazonRdsOracle`, `OracleConnectionTechnologyTypeOciAutonomousDatabase`, `OracleConnectionTechnologyTypeOracleAutonomousDatabaseAtAws`, `OracleConnectionTechnologyTypeOracleAutonomousDatabaseAtAzure`, `OracleConnectionTechnologyTypeOracleAutonomousDatabaseAtGoogleCloud`, `OracleConnectionTechnologyTypeOracleDatabase`, `OracleConnectionTechnologyTypeOracleExadata`, `OracleConnectionTechnologyTypeOracleExadataDatabaseAtAws`, `OracleConnectionTechnologyTypeOracleExadataDatabaseAtAzure`, `OracleConnectionTechnologyTypeOracleExadataDatabaseAtGoogleCloud`
- New enum type `ProximityPlacementGroupEntityType` with values `ProximityPlacementGroupEntityTypeCloudExadataInfrastructure`, `ProximityPlacementGroupEntityTypeOtherProducts`
- New enum type `RoutingMethod` with values `RoutingMethodDedicatedEndpoint`, `RoutingMethodSharedDeploymentEndpoint`, `RoutingMethodSharedServiceEndpoint`
- New enum type `SessionMode` with values `SessionModeDirect`, `SessionModeRedirect`
- New enum type `SetupType` with values `SetupTypeDevelopmentOrTesting`, `SetupTypeProduction`
- New function `*ClientFactory.NewDatabaseEditionsClient() *DatabaseEditionsClient`
- New function `*ClientFactory.NewDatabaseSystemShapeResourcesClient() *DatabaseSystemShapeResourcesClient`
- New function `*ClientFactory.NewGoldenGateConnectionsClient() *GoldenGateConnectionsClient`
- New function `*ClientFactory.NewGoldenGateDeploymentsClient() *GoldenGateDeploymentsClient`
- New function `*ConnectionBaseProperties.GetConnectionBaseProperties() *ConnectionBaseProperties`
- New function `NewDatabaseEditionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DatabaseEditionsClient, error)`
- New function `*DatabaseEditionsClient.Get(ctx context.Context, location string, databaseeditionname string, options *DatabaseEditionsClientGetOptions) (DatabaseEditionsClientGetResponse, error)`
- New function `*DatabaseEditionsClient.NewListByLocationPager(location string, options *DatabaseEditionsClientListByLocationOptions) *runtime.Pager[DatabaseEditionsClientListByLocationResponse]`
- New function `NewDatabaseSystemShapeResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DatabaseSystemShapeResourcesClient, error)`
- New function `*DatabaseSystemShapeResourcesClient.Get(ctx context.Context, location string, databasesystemshapename string, options *DatabaseSystemShapeResourcesClientGetOptions) (DatabaseSystemShapeResourcesClientGetResponse, error)`
- New function `*DatabaseSystemShapeResourcesClient.NewListByLocationPager(location string, options *DatabaseSystemShapeResourcesClientListByLocationOptions) *runtime.Pager[DatabaseSystemShapeResourcesClientListByLocationResponse]`
- New function `NewGoldenGateConnectionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*GoldenGateConnectionsClient, error)`
- New function `*GoldenGateConnectionsClient.BeginAssignDeployment(ctx context.Context, resourceGroupName string, goldenGateConnectionName string, body AssignUnassignDeployment, options *GoldenGateConnectionsClientBeginAssignDeploymentOptions) (*runtime.Poller[GoldenGateConnectionsClientAssignDeploymentResponse], error)`
- New function `*GoldenGateConnectionsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, goldenGateConnectionName string, resource GoldenGateConnection, options *GoldenGateConnectionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GoldenGateConnectionsClientCreateOrUpdateResponse], error)`
- New function `*GoldenGateConnectionsClient.BeginDelete(ctx context.Context, resourceGroupName string, goldenGateConnectionName string, options *GoldenGateConnectionsClientBeginDeleteOptions) (*runtime.Poller[GoldenGateConnectionsClientDeleteResponse], error)`
- New function `*GoldenGateConnectionsClient.Get(ctx context.Context, resourceGroupName string, goldenGateConnectionName string, options *GoldenGateConnectionsClientGetOptions) (GoldenGateConnectionsClientGetResponse, error)`
- New function `*GoldenGateConnectionsClient.GetAssignedDeployment(ctx context.Context, resourceGroupName string, goldenGateConnectionName string, assignmentID string, options *GoldenGateConnectionsClientGetAssignedDeploymentOptions) (GoldenGateConnectionsClientGetAssignedDeploymentResponse, error)`
- New function `*GoldenGateConnectionsClient.NewListAssignedDeploymentsByParentPager(resourceGroupName string, goldenGateConnectionName string, options *GoldenGateConnectionsClientListAssignedDeploymentsByParentOptions) *runtime.Pager[GoldenGateConnectionsClientListAssignedDeploymentsByParentResponse]`
- New function `*GoldenGateConnectionsClient.NewListByResourceGroupPager(resourceGroupName string, options *GoldenGateConnectionsClientListByResourceGroupOptions) *runtime.Pager[GoldenGateConnectionsClientListByResourceGroupResponse]`
- New function `*GoldenGateConnectionsClient.NewListBySubscriptionPager(options *GoldenGateConnectionsClientListBySubscriptionOptions) *runtime.Pager[GoldenGateConnectionsClientListBySubscriptionResponse]`
- New function `*GoldenGateConnectionsClient.BeginUnassignDeployment(ctx context.Context, resourceGroupName string, goldenGateConnectionName string, body AssignUnassignDeployment, options *GoldenGateConnectionsClientBeginUnassignDeploymentOptions) (*runtime.Poller[GoldenGateConnectionsClientUnassignDeploymentResponse], error)`
- New function `*GoldenGateConnectionsClient.BeginUpdate(ctx context.Context, resourceGroupName string, goldenGateConnectionName string, properties GoldenGateConnectionUpdate, options *GoldenGateConnectionsClientBeginUpdateOptions) (*runtime.Poller[GoldenGateConnectionsClientUpdateResponse], error)`
- New function `NewGoldenGateDeploymentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*GoldenGateDeploymentsClient, error)`
- New function `*GoldenGateDeploymentsClient.BeginAssignConnection(ctx context.Context, resourceGroupName string, goldenGateDeploymentName string, body AssignUnassignConnection, options *GoldenGateDeploymentsClientBeginAssignConnectionOptions) (*runtime.Poller[GoldenGateDeploymentsClientAssignConnectionResponse], error)`
- New function `*GoldenGateDeploymentsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, goldenGateDeploymentName string, resource GoldenGateDeployment, options *GoldenGateDeploymentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GoldenGateDeploymentsClientCreateOrUpdateResponse], error)`
- New function `*GoldenGateDeploymentsClient.BeginDelete(ctx context.Context, resourceGroupName string, goldenGateDeploymentName string, options *GoldenGateDeploymentsClientBeginDeleteOptions) (*runtime.Poller[GoldenGateDeploymentsClientDeleteResponse], error)`
- New function `*GoldenGateDeploymentsClient.Get(ctx context.Context, resourceGroupName string, goldenGateDeploymentName string, options *GoldenGateDeploymentsClientGetOptions) (GoldenGateDeploymentsClientGetResponse, error)`
- New function `*GoldenGateDeploymentsClient.GetAssignedConnection(ctx context.Context, resourceGroupName string, goldenGateDeploymentName string, assignmentID string, options *GoldenGateDeploymentsClientGetAssignedConnectionOptions) (GoldenGateDeploymentsClientGetAssignedConnectionResponse, error)`
- New function `*GoldenGateDeploymentsClient.NewListAssignedConnectionsByParentPager(resourceGroupName string, goldenGateDeploymentName string, options *GoldenGateDeploymentsClientListAssignedConnectionsByParentOptions) *runtime.Pager[GoldenGateDeploymentsClientListAssignedConnectionsByParentResponse]`
- New function `*GoldenGateDeploymentsClient.NewListByResourceGroupPager(resourceGroupName string, options *GoldenGateDeploymentsClientListByResourceGroupOptions) *runtime.Pager[GoldenGateDeploymentsClientListByResourceGroupResponse]`
- New function `*GoldenGateDeploymentsClient.NewListBySubscriptionPager(options *GoldenGateDeploymentsClientListBySubscriptionOptions) *runtime.Pager[GoldenGateDeploymentsClientListBySubscriptionResponse]`
- New function `*GoldenGateDeploymentsClient.BeginUnassignConnection(ctx context.Context, resourceGroupName string, goldenGateDeploymentName string, body AssignUnassignConnection, options *GoldenGateDeploymentsClientBeginUnassignConnectionOptions) (*runtime.Poller[GoldenGateDeploymentsClientUnassignConnectionResponse], error)`
- New function `*GoldenGateDeploymentsClient.BeginUpdate(ctx context.Context, resourceGroupName string, goldenGateDeploymentName string, properties GoldenGateDeploymentUpdate, options *GoldenGateDeploymentsClientBeginUpdateOptions) (*runtime.Poller[GoldenGateDeploymentsClientUpdateResponse], error)`
- New function `*KafkaConnectionDetails.GetConnectionBaseProperties() *ConnectionBaseProperties`
- New function `*MicrosoftFabricConnectionDetails.GetConnectionBaseProperties() *ConnectionBaseProperties`
- New function `*OracleConnectionDetails.GetConnectionBaseProperties() *ConnectionBaseProperties`
- New struct `AssignUnassignConnection`
- New struct `AssignUnassignDeployment`
- New struct `AssignedConnection`
- New struct `AssignedConnectionListResult`
- New struct `AssignedDeployment`
- New struct `AssignedDeploymentListResult`
- New struct `BackupScheduleType`
- New struct `DatabaseEdition`
- New struct `DatabaseEditionListResult`
- New struct `DatabaseEditionProperties`
- New struct `DatabaseSystemShape`
- New struct `DatabaseSystemShapeListResult`
- New struct `DatabaseSystemShapeProperties`
- New struct `DeploymentConnectionAssignmentProperties`
- New struct `DeploymentProperties`
- New struct `GoldenGateConnection`
- New struct `GoldenGateConnectionListResult`
- New struct `GoldenGateConnectionUpdate`
- New struct `GoldenGateConnectionUpdateProperties`
- New struct `GoldenGateDeployment`
- New struct `GoldenGateDeploymentListResult`
- New struct `GoldenGateDeploymentUpdate`
- New struct `GoldenGateDeploymentUpdateProperties`
- New struct `GroupToRolesMappingDetails`
- New struct `KafkaBootstrapServer`
- New struct `KafkaConnectionDetails`
- New struct `MaintenanceConfigurationType`
- New struct `MaintenanceWindowType`
- New struct `MicrosoftFabricConnectionDetails`
- New struct `OggDeploymentDetails`
- New struct `OracleConnectionDetails`
- New struct `ProximityPlacementGroup`
- New field `BackupDestination` in struct `AutonomousDatabaseBackupProperties`
- New field `BackupDestination`, `IsScheduleAzUpdateToEarliest`, `NetworkAnchorID`, `ResourceAnchorID`, `TimeScheduledAzUpdate`, `Zone` in struct `AutonomousDatabaseBaseProperties`
- New field `BackupDestination`, `IsScheduleAzUpdateToEarliest`, `NetworkAnchorID`, `ResourceAnchorID`, `TimeScheduledAzUpdate`, `Zone` in struct `AutonomousDatabaseCloneProperties`
- New field `BackupDestination`, `IsScheduleAzUpdateToEarliest`, `NetworkAnchorID`, `ResourceAnchorID`, `TimeScheduledAzUpdate`, `Zone` in struct `AutonomousDatabaseCrossRegionDisasterRecoveryProperties`
- New field `BackupDestination`, `IsScheduleAzUpdateToEarliest`, `NetworkAnchorID`, `ResourceAnchorID`, `TimeScheduledAzUpdate`, `Zone` in struct `AutonomousDatabaseFromBackupTimestampProperties`
- New field `BackupDestination`, `IsScheduleAzUpdateToEarliest`, `NetworkAnchorID`, `ResourceAnchorID`, `TimeScheduledAzUpdate`, `Zone` in struct `AutonomousDatabaseProperties`
- New field `ProximityPlacementGroup`, `ResourceAnchorID` in struct `CloudExadataInfrastructureProperties`
- New field `HostnameV2`, `IsAcceleratedNetworkEnabled`, `NetworkAnchorID`, `ProximityPlacementGroup`, `RecoStoragePercentage`, `ResourceAnchorID`, `ScanDNSNameV2`, `SparseStoragePercentage` in struct `CloudVMClusterProperties`
- New field `IsAcceleratedNetworkEnabled` in struct `CloudVMClusterUpdateProperties`
- New field `CharacterSet`, `DataCollectionOptions`, `NcharacterSet` in struct `DbSystemBaseProperties`
- New field `CharacterSet`, `DataCollectionOptions`, `NcharacterSet` in struct `DbSystemProperties`
- New field `HostnameV2`, `ScanDNSNameV2` in struct `ExadbVMClusterProperties`
- New field `AutoscaleLimitInGbs`, `IsAutoscaleEnabled` in struct `ExascaleDbStorageVaultProperties`
- New field `IsGiVersionForProvisioning`, `Shape`, `SortOrder` in struct `GiMinorVersionsClientListByParentOptions`
- New field `ProximityPlacementGroup` in struct `NetworkAnchorProperties`


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
