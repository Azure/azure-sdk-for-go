# Release History

## 4.0.0-beta.1 (2025-11-24)
### Breaking Changes

- Function `*AccessPolicyAssignmentClient.BeginCreateUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, AccessPolicyAssignment, *AccessPolicyAssignmentClientBeginCreateUpdateOptions)` to `(context.Context, string, string, string, string, AccessPolicyAssignment, *BeginCreateUpdateOptions)`
- Function `*AccessPolicyAssignmentClient.BeginCreateUpdate` return value(s) have been changed from `(*runtime.Poller[AccessPolicyAssignmentClientCreateUpdateResponse], error)` to `(*runtime.Poller[CreateUpdateResponse], error)`
- Function `*AccessPolicyAssignmentClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, string, *AccessPolicyAssignmentClientBeginDeleteOptions)` to `(context.Context, string, string, string, string, *BeginDeleteOptions)`
- Function `*AccessPolicyAssignmentClient.BeginDelete` return value(s) have been changed from `(*runtime.Poller[AccessPolicyAssignmentClientDeleteResponse], error)` to `(*runtime.Poller[DeleteResponse], error)`
- Function `*AccessPolicyAssignmentClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, *AccessPolicyAssignmentClientGetOptions)` to `(context.Context, string, string, string, string, *GetOptions)`
- Function `*AccessPolicyAssignmentClient.Get` return value(s) have been changed from `(AccessPolicyAssignmentClientGetResponse, error)` to `(GetResponse, error)`
- Function `*AccessPolicyAssignmentClient.NewListPager` parameter(s) have been changed from `(string, string, string, *AccessPolicyAssignmentClientListOptions)` to `(string, string, string, *ListOptions)`
- Function `*AccessPolicyAssignmentClient.NewListPager` return value(s) have been changed from `(*runtime.Pager[AccessPolicyAssignmentClientListResponse])` to `(*runtime.Pager[ListResponse])`
- Function `*Client.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, Cluster, *ClientBeginCreateOptions)` to `(context.Context, string, string, Cluster, *BeginCreateOptions)`
- Function `*Client.BeginCreate` return value(s) have been changed from `(*runtime.Poller[ClientCreateResponse], error)` to `(*runtime.Poller[CreateResponse], error)`
- Function `*Client.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *ClientBeginDeleteOptions)` to `(context.Context, string, string, *BeginDeleteOptions)`
- Function `*Client.BeginDelete` return value(s) have been changed from `(*runtime.Poller[ClientDeleteResponse], error)` to `(*runtime.Poller[DeleteResponse], error)`
- Function `*Client.Get` parameter(s) have been changed from `(context.Context, string, string, *ClientGetOptions)` to `(context.Context, string, string, *GetOptions)`
- Function `*Client.Get` return value(s) have been changed from `(ClientGetResponse, error)` to `(GetResponse, error)`
- Function `*Client.NewListByResourceGroupPager` parameter(s) have been changed from `(string, *ClientListByResourceGroupOptions)` to `(string, *ListByResourceGroupOptions)`
- Function `*Client.NewListByResourceGroupPager` return value(s) have been changed from `(*runtime.Pager[ClientListByResourceGroupResponse])` to `(*runtime.Pager[ListByResourceGroupResponse])`
- Function `*Client.NewListPager` parameter(s) have been changed from `(*ClientListOptions)` to `(*ListOptions)`
- Function `*Client.NewListPager` return value(s) have been changed from `(*runtime.Pager[ClientListResponse])` to `(*runtime.Pager[ListResponse])`
- Function `*Client.ListSKUsForScaling` parameter(s) have been changed from `(context.Context, string, string, *ClientListSKUsForScalingOptions)` to `(context.Context, string, string, *ListSKUsForScalingOptions)`
- Function `*Client.ListSKUsForScaling` return value(s) have been changed from `(ClientListSKUsForScalingResponse, error)` to `(ListSKUsForScalingResponse, error)`
- Function `*Client.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, ClusterUpdate, *ClientBeginUpdateOptions)` to `(context.Context, string, string, ClusterUpdate, *BeginUpdateOptions)`
- Function `*Client.BeginUpdate` return value(s) have been changed from `(*runtime.Poller[ClientUpdateResponse], error)` to `(*runtime.Poller[UpdateResponse], error)`
- Function `*DatabasesClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, Database, *DatabasesClientBeginCreateOptions)` to `(context.Context, string, string, string, Database, *BeginCreateOptions)`
- Function `*DatabasesClient.BeginCreate` return value(s) have been changed from `(*runtime.Poller[DatabasesClientCreateResponse], error)` to `(*runtime.Poller[CreateResponse], error)`
- Function `*DatabasesClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, *DatabasesClientBeginDeleteOptions)` to `(context.Context, string, string, string, *BeginDeleteOptions)`
- Function `*DatabasesClient.BeginDelete` return value(s) have been changed from `(*runtime.Poller[DatabasesClientDeleteResponse], error)` to `(*runtime.Poller[DeleteResponse], error)`
- Function `*DatabasesClient.BeginExport` parameter(s) have been changed from `(context.Context, string, string, string, ExportClusterParameters, *DatabasesClientBeginExportOptions)` to `(context.Context, string, string, string, ExportClusterParameters, *BeginExportOptions)`
- Function `*DatabasesClient.BeginExport` return value(s) have been changed from `(*runtime.Poller[DatabasesClientExportResponse], error)` to `(*runtime.Poller[ExportResponse], error)`
- Function `*DatabasesClient.BeginFlush` parameter(s) have been changed from `(context.Context, string, string, string, *DatabasesClientBeginFlushOptions)` to `(context.Context, string, string, string, *BeginFlushOptions)`
- Function `*DatabasesClient.BeginFlush` return value(s) have been changed from `(*runtime.Poller[DatabasesClientFlushResponse], error)` to `(*runtime.Poller[FlushResponse], error)`
- Function `*DatabasesClient.BeginForceLinkToReplicationGroup` parameter(s) have been changed from `(context.Context, string, string, string, ForceLinkParameters, *DatabasesClientBeginForceLinkToReplicationGroupOptions)` to `(context.Context, string, string, string, ForceLinkParameters, *BeginForceLinkToReplicationGroupOptions)`
- Function `*DatabasesClient.BeginForceLinkToReplicationGroup` return value(s) have been changed from `(*runtime.Poller[DatabasesClientForceLinkToReplicationGroupResponse], error)` to `(*runtime.Poller[ForceLinkToReplicationGroupResponse], error)`
- Function `*DatabasesClient.BeginForceUnlink` parameter(s) have been changed from `(context.Context, string, string, string, ForceUnlinkParameters, *DatabasesClientBeginForceUnlinkOptions)` to `(context.Context, string, string, string, ForceUnlinkParameters, *BeginForceUnlinkOptions)`
- Function `*DatabasesClient.BeginForceUnlink` return value(s) have been changed from `(*runtime.Poller[DatabasesClientForceUnlinkResponse], error)` to `(*runtime.Poller[ForceUnlinkResponse], error)`
- Function `*DatabasesClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *DatabasesClientGetOptions)` to `(context.Context, string, string, string, *GetOptions)`
- Function `*DatabasesClient.Get` return value(s) have been changed from `(DatabasesClientGetResponse, error)` to `(GetResponse, error)`
- Function `*DatabasesClient.BeginImport` parameter(s) have been changed from `(context.Context, string, string, string, ImportClusterParameters, *DatabasesClientBeginImportOptions)` to `(context.Context, string, string, string, ImportClusterParameters, *BeginImportOptions)`
- Function `*DatabasesClient.BeginImport` return value(s) have been changed from `(*runtime.Poller[DatabasesClientImportResponse], error)` to `(*runtime.Poller[ImportResponse], error)`
- Function `*DatabasesClient.NewListByClusterPager` parameter(s) have been changed from `(string, string, *DatabasesClientListByClusterOptions)` to `(string, string, *ListByClusterOptions)`
- Function `*DatabasesClient.NewListByClusterPager` return value(s) have been changed from `(*runtime.Pager[DatabasesClientListByClusterResponse])` to `(*runtime.Pager[ListByClusterResponse])`
- Function `*DatabasesClient.ListKeys` parameter(s) have been changed from `(context.Context, string, string, string, *DatabasesClientListKeysOptions)` to `(context.Context, string, string, string, *ListKeysOptions)`
- Function `*DatabasesClient.ListKeys` return value(s) have been changed from `(DatabasesClientListKeysResponse, error)` to `(ListKeysResponse, error)`
- Function `*DatabasesClient.BeginRegenerateKey` parameter(s) have been changed from `(context.Context, string, string, string, RegenerateKeyParameters, *DatabasesClientBeginRegenerateKeyOptions)` to `(context.Context, string, string, string, RegenerateKeyParameters, *BeginRegenerateKeyOptions)`
- Function `*DatabasesClient.BeginRegenerateKey` return value(s) have been changed from `(*runtime.Poller[DatabasesClientRegenerateKeyResponse], error)` to `(*runtime.Poller[RegenerateKeyResponse], error)`
- Function `*DatabasesClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, DatabaseUpdate, *DatabasesClientBeginUpdateOptions)` to `(context.Context, string, string, string, DatabaseUpdate, *BeginUpdateOptions)`
- Function `*DatabasesClient.BeginUpdate` return value(s) have been changed from `(*runtime.Poller[DatabasesClientUpdateResponse], error)` to `(*runtime.Poller[UpdateResponse], error)`
- Function `*DatabasesClient.BeginUpgradeDBRedisVersion` parameter(s) have been changed from `(context.Context, string, string, string, *DatabasesClientBeginUpgradeDBRedisVersionOptions)` to `(context.Context, string, string, string, *BeginUpgradeDBRedisVersionOptions)`
- Function `*DatabasesClient.BeginUpgradeDBRedisVersion` return value(s) have been changed from `(*runtime.Poller[DatabasesClientUpgradeDBRedisVersionResponse], error)` to `(*runtime.Poller[UpgradeDBRedisVersionResponse], error)`
- Function `*OperationsClient.NewListPager` parameter(s) have been changed from `(*OperationsClientListOptions)` to `(*ListOptions)`
- Function `*OperationsClient.NewListPager` return value(s) have been changed from `(*runtime.Pager[OperationsClientListResponse])` to `(*runtime.Pager[ListResponse])`
- Function `*OperationsStatusClient.Get` parameter(s) have been changed from `(context.Context, string, string, *OperationsStatusClientGetOptions)` to `(context.Context, string, string, *GetOptions)`
- Function `*OperationsStatusClient.Get` return value(s) have been changed from `(OperationsStatusClientGetResponse, error)` to `(GetResponse, error)`
- Function `*PrivateEndpointConnectionsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions)` to `(context.Context, string, string, string, *BeginDeleteOptions)`
- Function `*PrivateEndpointConnectionsClient.BeginDelete` return value(s) have been changed from `(*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)` to `(*runtime.Poller[DeleteResponse], error)`
- Function `*PrivateEndpointConnectionsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions)` to `(context.Context, string, string, string, *GetOptions)`
- Function `*PrivateEndpointConnectionsClient.Get` return value(s) have been changed from `(PrivateEndpointConnectionsClientGetResponse, error)` to `(GetResponse, error)`
- Function `*PrivateEndpointConnectionsClient.NewListPager` parameter(s) have been changed from `(string, string, *PrivateEndpointConnectionsClientListOptions)` to `(string, string, *ListOptions)`
- Function `*PrivateEndpointConnectionsClient.NewListPager` return value(s) have been changed from `(*runtime.Pager[PrivateEndpointConnectionsClientListResponse])` to `(*runtime.Pager[ListResponse])`
- Function `*PrivateEndpointConnectionsClient.BeginPut` parameter(s) have been changed from `(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientBeginPutOptions)` to `(context.Context, string, string, string, PrivateEndpointConnection, *BeginPutOptions)`
- Function `*PrivateEndpointConnectionsClient.BeginPut` return value(s) have been changed from `(*runtime.Poller[PrivateEndpointConnectionsClientPutResponse], error)` to `(*runtime.Poller[PutResponse], error)`
- Function `*PrivateLinkResourcesClient.NewListByClusterPager` parameter(s) have been changed from `(string, string, *PrivateLinkResourcesClientListByClusterOptions)` to `(string, string, *ListByClusterOptions)`
- Function `*PrivateLinkResourcesClient.NewListByClusterPager` return value(s) have been changed from `(*runtime.Pager[PrivateLinkResourcesClientListByClusterResponse])` to `(*runtime.Pager[ListByClusterResponse])`
- Type of `ClusterCreateProperties.Encryption` has been changed from `*ClusterCommonPropertiesEncryption` to `*ClusterPropertiesEncryption`
- Type of `ClusterUpdateProperties.Encryption` has been changed from `*ClusterCommonPropertiesEncryption` to `*ClusterPropertiesEncryption`
- Type of `DatabaseCreateProperties.GeoReplication` has been changed from `*DatabaseCommonPropertiesGeoReplication` to `*DatabasePropertiesGeoReplication`
- Type of `DatabaseUpdateProperties.GeoReplication` has been changed from `*DatabaseCommonPropertiesGeoReplication` to `*DatabasePropertiesGeoReplication`
- Struct `AccessPolicyAssignmentClientBeginCreateUpdateOptions` has been removed
- Struct `AccessPolicyAssignmentClientBeginDeleteOptions` has been removed
- Struct `AccessPolicyAssignmentClientCreateUpdateResponse` has been removed
- Struct `AccessPolicyAssignmentClientDeleteResponse` has been removed
- Struct `AccessPolicyAssignmentClientGetOptions` has been removed
- Struct `AccessPolicyAssignmentClientGetResponse` has been removed
- Struct `AccessPolicyAssignmentClientListOptions` has been removed
- Struct `AccessPolicyAssignmentClientListResponse` has been removed
- Struct `ClientBeginCreateOptions` has been removed
- Struct `ClientBeginDeleteOptions` has been removed
- Struct `ClientBeginUpdateOptions` has been removed
- Struct `ClientCreateResponse` has been removed
- Struct `ClientDeleteResponse` has been removed
- Struct `ClientGetOptions` has been removed
- Struct `ClientGetResponse` has been removed
- Struct `ClientListByResourceGroupOptions` has been removed
- Struct `ClientListByResourceGroupResponse` has been removed
- Struct `ClientListOptions` has been removed
- Struct `ClientListResponse` has been removed
- Struct `ClientListSKUsForScalingOptions` has been removed
- Struct `ClientListSKUsForScalingResponse` has been removed
- Struct `ClientUpdateResponse` has been removed
- Struct `ClusterCommonProperties` has been removed
- Struct `ClusterCommonPropertiesEncryption` has been removed
- Struct `ClusterCommonPropertiesEncryptionCustomerManagedKeyEncryption` has been removed
- Struct `ClusterCommonPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity` has been removed
- Struct `DatabaseCommonProperties` has been removed
- Struct `DatabaseCommonPropertiesGeoReplication` has been removed
- Struct `DatabasesClientBeginCreateOptions` has been removed
- Struct `DatabasesClientBeginDeleteOptions` has been removed
- Struct `DatabasesClientBeginExportOptions` has been removed
- Struct `DatabasesClientBeginFlushOptions` has been removed
- Struct `DatabasesClientBeginForceLinkToReplicationGroupOptions` has been removed
- Struct `DatabasesClientBeginForceUnlinkOptions` has been removed
- Struct `DatabasesClientBeginImportOptions` has been removed
- Struct `DatabasesClientBeginRegenerateKeyOptions` has been removed
- Struct `DatabasesClientBeginUpdateOptions` has been removed
- Struct `DatabasesClientBeginUpgradeDBRedisVersionOptions` has been removed
- Struct `DatabasesClientCreateResponse` has been removed
- Struct `DatabasesClientDeleteResponse` has been removed
- Struct `DatabasesClientExportResponse` has been removed
- Struct `DatabasesClientFlushResponse` has been removed
- Struct `DatabasesClientForceLinkToReplicationGroupResponse` has been removed
- Struct `DatabasesClientForceUnlinkResponse` has been removed
- Struct `DatabasesClientGetOptions` has been removed
- Struct `DatabasesClientGetResponse` has been removed
- Struct `DatabasesClientImportResponse` has been removed
- Struct `DatabasesClientListByClusterOptions` has been removed
- Struct `DatabasesClientListByClusterResponse` has been removed
- Struct `DatabasesClientListKeysOptions` has been removed
- Struct `DatabasesClientListKeysResponse` has been removed
- Struct `DatabasesClientRegenerateKeyResponse` has been removed
- Struct `DatabasesClientUpdateResponse` has been removed
- Struct `DatabasesClientUpgradeDBRedisVersionResponse` has been removed
- Struct `ErrorDetailAutoGenerated` has been removed
- Struct `ErrorResponseAutoGenerated` has been removed
- Struct `OperationsClientListOptions` has been removed
- Struct `OperationsClientListResponse` has been removed
- Struct `OperationsStatusClientGetOptions` has been removed
- Struct `OperationsStatusClientGetResponse` has been removed
- Struct `PrivateEndpointConnectionsClientBeginDeleteOptions` has been removed
- Struct `PrivateEndpointConnectionsClientBeginPutOptions` has been removed
- Struct `PrivateEndpointConnectionsClientDeleteResponse` has been removed
- Struct `PrivateEndpointConnectionsClientGetOptions` has been removed
- Struct `PrivateEndpointConnectionsClientGetResponse` has been removed
- Struct `PrivateEndpointConnectionsClientListOptions` has been removed
- Struct `PrivateEndpointConnectionsClientListResponse` has been removed
- Struct `PrivateEndpointConnectionsClientPutResponse` has been removed
- Struct `PrivateLinkResourcesClientListByClusterOptions` has been removed
- Struct `PrivateLinkResourcesClientListByClusterResponse` has been removed
- Struct `ProxyResource` has been removed
- Struct `ProxyResourceAutoGenerated` has been removed
- Struct `Resource` has been removed
- Struct `ResourceAutoGenerated` has been removed
- Struct `TrackedResource` has been removed

### Features Added

- New enum type `MaintenanceDayOfWeek` with values `MaintenanceDayOfWeekFriday`, `MaintenanceDayOfWeekMonday`, `MaintenanceDayOfWeekSaturday`, `MaintenanceDayOfWeekSunday`, `MaintenanceDayOfWeekThursday`, `MaintenanceDayOfWeekTuesday`, `MaintenanceDayOfWeekWednesday`
- New enum type `MaintenanceWindowType` with values `MaintenanceWindowTypeWeekly`
- New enum type `MigrationProvisioningState` with values `MigrationProvisioningStateAccepted`, `MigrationProvisioningStateCancelled`, `MigrationProvisioningStateCancelling`, `MigrationProvisioningStateFailed`, `MigrationProvisioningStateInProgress`, `MigrationProvisioningStateReadyForDNSSwitch`, `MigrationProvisioningStateSucceeded`
- New enum type `SourceType` with values `SourceTypeAzureCacheForRedis`
- New function `*AzureCacheForRedisMigrationProperties.GetMigrationProperties() *MigrationProperties`
- New function `*ClientFactory.NewMigrationClient() *MigrationClient`
- New function `NewMigrationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MigrationClient, error)`
- New function `*MigrationClient.BeginCancel(context.Context, string, string, *BeginCancelOptions) (*runtime.Poller[CancelResponse], error)`
- New function `*MigrationClient.Get(context.Context, string, string, *GetOptions) (GetResponse, error)`
- New function `*MigrationClient.NewListPager(string, string, *ListOptions) *runtime.Pager[ListResponse]`
- New function `*MigrationClient.BeginStart(context.Context, string, string, Migration, *BeginStartOptions) (*runtime.Poller[StartResponse], error)`
- New function `*MigrationProperties.GetMigrationProperties() *MigrationProperties`
- New struct `AzureCacheForRedisMigrationProperties`
- New struct `BeginCreateOptions`
- New struct `BeginCreateUpdateOptions`
- New struct `BeginDeleteOptions`
- New struct `BeginExportOptions`
- New struct `BeginFlushOptions`
- New struct `BeginForceLinkToReplicationGroupOptions`
- New struct `BeginForceUnlinkOptions`
- New struct `BeginImportOptions`
- New struct `BeginPutOptions`
- New struct `BeginRegenerateKeyOptions`
- New struct `BeginUpdateOptions`
- New struct `BeginUpgradeDBRedisVersionOptions`
- New struct `ClusterPropertiesEncryption`
- New struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryption`
- New struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity`
- New struct `CreateResponse`
- New struct `CreateUpdateResponse`
- New struct `DatabasePropertiesGeoReplication`
- New struct `DeleteResponse`
- New struct `ExportResponse`
- New struct `FlushResponse`
- New struct `ForceLinkToReplicationGroupResponse`
- New struct `ForceUnlinkResponse`
- New struct `ImportResponse`
- New struct `ListByClusterOptions`
- New struct `ListByClusterResponse`
- New struct `ListByResourceGroupOptions`
- New struct `ListByResourceGroupResponse`
- New struct `ListKeysOptions`
- New struct `ListKeysResponse`
- New struct `ListSKUsForScalingOptions`
- New struct `ListSKUsForScalingResponse`
- New struct `MaintenanceConfiguration`
- New struct `MaintenanceWindow`
- New struct `MaintenanceWindowSchedule`
- New struct `Migration`
- New struct `MigrationList`
- New struct `PutResponse`
- New struct `RegenerateKeyResponse`
- New struct `UpdateResponse`
- New struct `UpgradeDBRedisVersionResponse`
- New field `SystemData` in struct `AccessPolicyAssignment`
- New field `SystemData` in struct `Cluster`
- New field `MaintenanceConfiguration` in struct `ClusterCreateProperties`
- New field `MaintenanceConfiguration` in struct `ClusterUpdateProperties`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `NextLink` in struct `PrivateEndpointConnectionListResult`
- New field `GroupIDs` in struct `PrivateEndpointConnectionProperties`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `NextLink` in struct `PrivateLinkResourceListResult`


## 3.0.0 (2025-10-17)
### Breaking Changes

- Function `*DatabasesClient.BeginFlush` parameter(s) have been changed from `(context.Context, string, string, string, FlushParameters, *DatabasesClientBeginFlushOptions)` to `(context.Context, string, string, string, *DatabasesClientBeginFlushOptions)`
- Type of `Cluster.Properties` has been changed from `*ClusterProperties` to `*ClusterCreateProperties`
- Type of `ClusterUpdate.Properties` has been changed from `*ClusterProperties` to `*ClusterUpdateProperties`
- Type of `Database.Properties` has been changed from `*DatabaseProperties` to `*DatabaseCreateProperties`
- Type of `DatabaseUpdate.Properties` has been changed from `*DatabaseProperties` to `*DatabaseUpdateProperties`
- Struct `ClusterProperties` has been removed
- Struct `ClusterPropertiesEncryption` has been removed
- Struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryption` has been removed
- Struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity` has been removed
- Struct `DatabaseProperties` has been removed
- Struct `DatabasePropertiesGeoReplication` has been removed

### Features Added

- New value `ClusteringPolicyNoCluster` added to enum type `ClusteringPolicy`
- New value `ResourceStateMoving` added to enum type `ResourceState`
- New value `SKUNameBalancedB0`, `SKUNameBalancedB1`, `SKUNameBalancedB10`, `SKUNameBalancedB100`, `SKUNameBalancedB1000`, `SKUNameBalancedB150`, `SKUNameBalancedB20`, `SKUNameBalancedB250`, `SKUNameBalancedB3`, `SKUNameBalancedB350`, `SKUNameBalancedB5`, `SKUNameBalancedB50`, `SKUNameBalancedB500`, `SKUNameBalancedB700`, `SKUNameComputeOptimizedX10`, `SKUNameComputeOptimizedX100`, `SKUNameComputeOptimizedX150`, `SKUNameComputeOptimizedX20`, `SKUNameComputeOptimizedX250`, `SKUNameComputeOptimizedX3`, `SKUNameComputeOptimizedX350`, `SKUNameComputeOptimizedX5`, `SKUNameComputeOptimizedX50`, `SKUNameComputeOptimizedX500`, `SKUNameComputeOptimizedX700`, `SKUNameEnterpriseE1`, `SKUNameEnterpriseE200`, `SKUNameEnterpriseE400`, `SKUNameEnterpriseE5`, `SKUNameFlashOptimizedA1000`, `SKUNameFlashOptimizedA1500`, `SKUNameFlashOptimizedA2000`, `SKUNameFlashOptimizedA250`, `SKUNameFlashOptimizedA4500`, `SKUNameFlashOptimizedA500`, `SKUNameFlashOptimizedA700`, `SKUNameMemoryOptimizedM10`, `SKUNameMemoryOptimizedM100`, `SKUNameMemoryOptimizedM1000`, `SKUNameMemoryOptimizedM150`, `SKUNameMemoryOptimizedM1500`, `SKUNameMemoryOptimizedM20`, `SKUNameMemoryOptimizedM2000`, `SKUNameMemoryOptimizedM250`, `SKUNameMemoryOptimizedM350`, `SKUNameMemoryOptimizedM50`, `SKUNameMemoryOptimizedM500`, `SKUNameMemoryOptimizedM700` added to enum type `SKUName`
- New enum type `AccessKeysAuthentication` with values `AccessKeysAuthenticationDisabled`, `AccessKeysAuthenticationEnabled`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DeferUpgradeSetting` with values `DeferUpgradeSettingDeferred`, `DeferUpgradeSettingNotDeferred`
- New enum type `HighAvailability` with values `HighAvailabilityDisabled`, `HighAvailabilityEnabled`
- New enum type `Kind` with values `KindV1`, `KindV2`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New enum type `RedundancyMode` with values `RedundancyModeLR`, `RedundancyModeNone`, `RedundancyModeZR`
- New function `NewAccessPolicyAssignmentClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessPolicyAssignmentClient, error)`
- New function `*AccessPolicyAssignmentClient.BeginCreateUpdate(context.Context, string, string, string, string, AccessPolicyAssignment, *AccessPolicyAssignmentClientBeginCreateUpdateOptions) (*runtime.Poller[AccessPolicyAssignmentClientCreateUpdateResponse], error)`
- New function `*AccessPolicyAssignmentClient.BeginDelete(context.Context, string, string, string, string, *AccessPolicyAssignmentClientBeginDeleteOptions) (*runtime.Poller[AccessPolicyAssignmentClientDeleteResponse], error)`
- New function `*AccessPolicyAssignmentClient.Get(context.Context, string, string, string, string, *AccessPolicyAssignmentClientGetOptions) (AccessPolicyAssignmentClientGetResponse, error)`
- New function `*AccessPolicyAssignmentClient.NewListPager(string, string, string, *AccessPolicyAssignmentClientListOptions) *runtime.Pager[AccessPolicyAssignmentClientListResponse]`
- New function `*Client.ListSKUsForScaling(context.Context, string, string, *ClientListSKUsForScalingOptions) (ClientListSKUsForScalingResponse, error)`
- New function `*ClientFactory.NewAccessPolicyAssignmentClient() *AccessPolicyAssignmentClient`
- New function `*DatabasesClient.BeginForceLinkToReplicationGroup(context.Context, string, string, string, ForceLinkParameters, *DatabasesClientBeginForceLinkToReplicationGroupOptions) (*runtime.Poller[DatabasesClientForceLinkToReplicationGroupResponse], error)`
- New function `*DatabasesClient.BeginUpgradeDBRedisVersion(context.Context, string, string, string, *DatabasesClientBeginUpgradeDBRedisVersionOptions) (*runtime.Poller[DatabasesClientUpgradeDBRedisVersionResponse], error)`
- New struct `AccessPolicyAssignment`
- New struct `AccessPolicyAssignmentList`
- New struct `AccessPolicyAssignmentProperties`
- New struct `AccessPolicyAssignmentPropertiesUser`
- New struct `ClusterCommonProperties`
- New struct `ClusterCommonPropertiesEncryption`
- New struct `ClusterCommonPropertiesEncryptionCustomerManagedKeyEncryption`
- New struct `ClusterCommonPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity`
- New struct `ClusterCreateProperties`
- New struct `ClusterUpdateProperties`
- New struct `DatabaseCommonProperties`
- New struct `DatabaseCommonPropertiesGeoReplication`
- New struct `DatabaseCreateProperties`
- New struct `DatabaseUpdateProperties`
- New struct `ErrorDetailAutoGenerated`
- New struct `ErrorResponseAutoGenerated`
- New struct `ForceLinkParameters`
- New struct `ForceLinkParametersGeoReplication`
- New struct `ProxyResourceAutoGenerated`
- New struct `ResourceAutoGenerated`
- New struct `SKUDetails`
- New struct `SKUDetailsList`
- New struct `SystemData`
- New field `Kind` in struct `Cluster`
- New field `SystemData` in struct `Database`
- New field `Parameters` in struct `DatabasesClientBeginFlushOptions`
- New field `SystemData` in struct `ProxyResource`


## 2.1.0-beta.3 (2025-04-23)
### Breaking Changes

- Function `*DatabasesClient.BeginFlush` parameter(s) have been changed from `(context.Context, string, string, string, FlushParameters, *DatabasesClientBeginFlushOptions)` to `(context.Context, string, string, string, *DatabasesClientBeginFlushOptions)`
- Field `GroupNickname`, `LinkedDatabases` of struct `ForceLinkParameters` has been removed

### Features Added

- New value `ClusteringPolicyNoCluster` added to enum type `ClusteringPolicy`
- New value `ResourceStateMoving` added to enum type `ResourceState`
- New enum type `Kind` with values `KindV1`, `KindV2`
- New function `*Client.ListSKUsForScaling(context.Context, string, string, *ClientListSKUsForScalingOptions) (ClientListSKUsForScalingResponse, error)`
- New struct `ErrorDetailAutoGenerated`
- New struct `ErrorResponseAutoGenerated`
- New struct `ForceLinkParametersGeoReplication`
- New struct `SKUDetails`
- New struct `SKUDetailsList`
- New field `Kind` in struct `Cluster`
- New field `Parameters` in struct `DatabasesClientBeginFlushOptions`
- New field `GeoReplication` in struct `ForceLinkParameters`


## 2.1.0-beta.2 (2024-09-27)
### Features Added

- New value `SKUNameBalancedB0`, `SKUNameBalancedB1`, `SKUNameBalancedB10`, `SKUNameBalancedB100`, `SKUNameBalancedB1000`, `SKUNameBalancedB150`, `SKUNameBalancedB20`, `SKUNameBalancedB250`, `SKUNameBalancedB3`, `SKUNameBalancedB350`, `SKUNameBalancedB5`, `SKUNameBalancedB50`, `SKUNameBalancedB500`, `SKUNameBalancedB700`, `SKUNameComputeOptimizedX10`, `SKUNameComputeOptimizedX100`, `SKUNameComputeOptimizedX150`, `SKUNameComputeOptimizedX20`, `SKUNameComputeOptimizedX250`, `SKUNameComputeOptimizedX3`, `SKUNameComputeOptimizedX350`, `SKUNameComputeOptimizedX5`, `SKUNameComputeOptimizedX50`, `SKUNameComputeOptimizedX500`, `SKUNameComputeOptimizedX700`, `SKUNameEnterpriseE1`, `SKUNameEnterpriseE200`, `SKUNameEnterpriseE400`, `SKUNameFlashOptimizedA1000`, `SKUNameFlashOptimizedA1500`, `SKUNameFlashOptimizedA2000`, `SKUNameFlashOptimizedA250`, `SKUNameFlashOptimizedA4500`, `SKUNameFlashOptimizedA500`, `SKUNameFlashOptimizedA700`, `SKUNameMemoryOptimizedM10`, `SKUNameMemoryOptimizedM100`, `SKUNameMemoryOptimizedM1000`, `SKUNameMemoryOptimizedM150`, `SKUNameMemoryOptimizedM1500`, `SKUNameMemoryOptimizedM20`, `SKUNameMemoryOptimizedM2000`, `SKUNameMemoryOptimizedM250`, `SKUNameMemoryOptimizedM350`, `SKUNameMemoryOptimizedM50`, `SKUNameMemoryOptimizedM500`, `SKUNameMemoryOptimizedM700` added to enum type `SKUName`
- New enum type `AccessKeysAuthentication` with values `AccessKeysAuthenticationDisabled`, `AccessKeysAuthenticationEnabled`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `HighAvailability` with values `HighAvailabilityDisabled`, `HighAvailabilityEnabled`
- New enum type `RedundancyMode` with values `RedundancyModeLR`, `RedundancyModeNone`, `RedundancyModeZR`
- New function `NewAccessPolicyAssignmentClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessPolicyAssignmentClient, error)`
- New function `*AccessPolicyAssignmentClient.BeginCreateUpdate(context.Context, string, string, string, string, AccessPolicyAssignment, *AccessPolicyAssignmentClientBeginCreateUpdateOptions) (*runtime.Poller[AccessPolicyAssignmentClientCreateUpdateResponse], error)`
- New function `*AccessPolicyAssignmentClient.BeginDelete(context.Context, string, string, string, string, *AccessPolicyAssignmentClientBeginDeleteOptions) (*runtime.Poller[AccessPolicyAssignmentClientDeleteResponse], error)`
- New function `*AccessPolicyAssignmentClient.Get(context.Context, string, string, string, string, *AccessPolicyAssignmentClientGetOptions) (AccessPolicyAssignmentClientGetResponse, error)`
- New function `*AccessPolicyAssignmentClient.NewListPager(string, string, string, *AccessPolicyAssignmentClientListOptions) *runtime.Pager[AccessPolicyAssignmentClientListResponse]`
- New function `*ClientFactory.NewAccessPolicyAssignmentClient() *AccessPolicyAssignmentClient`
- New struct `AccessPolicyAssignment`
- New struct `AccessPolicyAssignmentList`
- New struct `AccessPolicyAssignmentProperties`
- New struct `AccessPolicyAssignmentPropertiesUser`
- New struct `ProxyResourceAutoGenerated`
- New struct `ResourceAutoGenerated`
- New struct `SystemData`
- New field `HighAvailability`, `RedundancyMode` in struct `ClusterProperties`
- New field `SystemData` in struct `Database`
- New field `AccessKeysAuthentication` in struct `DatabaseProperties`
- New field `SystemData` in struct `ProxyResource`


## 2.1.0-beta.1 (2024-05-24)
### Features Added

- New value `SKUNameEnterpriseE5` added to enum type `SKUName`
- New enum type `DeferUpgradeSetting` with values `DeferUpgradeSettingDeferred`, `DeferUpgradeSettingNotDeferred`
- New function `*DatabasesClient.BeginForceLinkToReplicationGroup(context.Context, string, string, string, ForceLinkParameters, *DatabasesClientBeginForceLinkToReplicationGroupOptions) (*runtime.Poller[DatabasesClientForceLinkToReplicationGroupResponse], error)`
- New function `*DatabasesClient.BeginUpgradeDBRedisVersion(context.Context, string, string, string, *DatabasesClientBeginUpgradeDBRedisVersionOptions) (*runtime.Poller[DatabasesClientUpgradeDBRedisVersionResponse], error)`
- New struct `ForceLinkParameters`
- New field `DeferUpgrade`, `RedisVersion` in struct `DatabaseProperties`


## 2.0.0 (2024-02-23)
### Breaking Changes

- Operation `*PrivateEndpointConnectionsClient.Delete` has been changed to LRO, use `*PrivateEndpointConnectionsClient.BeginDelete` instead.

### Features Added

- New value `ResourceStateScaling`, `ResourceStateScalingFailed` added to enum type `ResourceState`
- New enum type `CmkIdentityType` with values `CmkIdentityTypeSystemAssignedIdentity`, `CmkIdentityTypeUserAssignedIdentity`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New function `*DatabasesClient.BeginFlush(context.Context, string, string, string, FlushParameters, *DatabasesClientBeginFlushOptions) (*runtime.Poller[DatabasesClientFlushResponse], error)`
- New struct `ClusterPropertiesEncryption`
- New struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryption`
- New struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity`
- New struct `FlushParameters`
- New struct `ManagedServiceIdentity`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Cluster`
- New field `Encryption` in struct `ClusterProperties`
- New field `Identity` in struct `ClusterUpdate`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).