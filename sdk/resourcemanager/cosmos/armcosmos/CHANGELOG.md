# Release History

## 2.7.0 (2024-01-26)
### Features Added

- New value `OperationTypeRecreate` added to enum type `OperationType`
- New struct `ComputedProperty`
- New struct `ResourceRestoreParameters`
- New struct `RestoreParametersBase`
- New field `CustomerManagedKeyStatus`, `EnableBurstCapacity` in struct `DatabaseAccountCreateUpdateProperties`
- New field `CustomerManagedKeyStatus`, `EnableBurstCapacity` in struct `DatabaseAccountGetProperties`
- New field `CustomerManagedKeyStatus`, `EnableBurstCapacity` in struct `DatabaseAccountUpdateProperties`
- New field `CreateMode`, `RestoreParameters` in struct `GremlinDatabaseGetPropertiesResource`
- New field `CreateMode`, `RestoreParameters` in struct `GremlinDatabaseResource`
- New field `CreateMode`, `RestoreParameters` in struct `GremlinGraphGetPropertiesResource`
- New field `CreateMode`, `RestoreParameters` in struct `GremlinGraphResource`
- New field `CreateMode`, `RestoreParameters` in struct `MongoDBCollectionGetPropertiesResource`
- New field `CreateMode`, `RestoreParameters` in struct `MongoDBCollectionResource`
- New field `CreateMode`, `RestoreParameters` in struct `MongoDBDatabaseGetPropertiesResource`
- New field `CreateMode`, `RestoreParameters` in struct `MongoDBDatabaseResource`
- New field `CanUndelete`, `CanUndeleteReason` in struct `RestorableGremlinDatabasePropertiesResource`
- New field `CanUndelete`, `CanUndeleteReason` in struct `RestorableGremlinGraphPropertiesResource`
- New field `CanUndelete`, `CanUndeleteReason` in struct `RestorableMongodbCollectionPropertiesResource`
- New field `CanUndelete`, `CanUndeleteReason` in struct `RestorableMongodbDatabasePropertiesResource`
- New field `CanUndelete`, `CanUndeleteReason` in struct `RestorableSQLContainerPropertiesResource`
- New field `ComputedProperties`, `CreateMode`, `RestoreParameters` in struct `RestorableSQLContainerPropertiesResourceContainer`
- New field `CanUndelete`, `CanUndeleteReason` in struct `RestorableSQLDatabasePropertiesResource`
- New field `CreateMode`, `RestoreParameters` in struct `RestorableSQLDatabasePropertiesResourceDatabase`
- New field `CanUndelete`, `CanUndeleteReason` in struct `RestorableTablePropertiesResource`
- New field `ComputedProperties`, `CreateMode`, `RestoreParameters` in struct `SQLContainerGetPropertiesResource`
- New field `ComputedProperties`, `CreateMode`, `RestoreParameters` in struct `SQLContainerResource`
- New field `CreateMode`, `RestoreParameters` in struct `SQLDatabaseGetPropertiesResource`
- New field `CreateMode`, `RestoreParameters` in struct `SQLDatabaseResource`
- New field `CreateMode`, `RestoreParameters` in struct `TableGetPropertiesResource`
- New field `CreateMode`, `RestoreParameters` in struct `TableResource`


## 3.0.0-beta.2 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.6.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 3.0.0-beta.1 (2023-06-23)
### Breaking Changes

- `PublicNetworkAccessSecuredByPerimeter` from enum `PublicNetworkAccess` has been removed
- Field `InstantMaximumThroughput`, `SoftAllowedMaximumThroughput` of struct `ThroughputSettingsGetPropertiesResource` has been removed
- Field `InstantMaximumThroughput`, `SoftAllowedMaximumThroughput` of struct `ThroughputSettingsResource` has been removed

### Features Added

- New value `CreateModePointInTimeRestore` added to enum type `CreateMode`
- New value `OperationTypeRecreate` added to enum type `OperationType`
- New enum type `CheckNameAvailabilityReason` with values `CheckNameAvailabilityReasonAlreadyExists`, `CheckNameAvailabilityReasonInvalid`
- New enum type `DataTransferComponent` with values `DataTransferComponentAzureBlobStorage`, `DataTransferComponentCosmosDBCassandra`, `DataTransferComponentCosmosDBMongo`, `DataTransferComponentCosmosDBSQL`
- New enum type `EnableFullTextQuery` with values `EnableFullTextQueryFalse`, `EnableFullTextQueryNone`, `EnableFullTextQueryTrue`
- New enum type `MongoClusterStatus` with values `MongoClusterStatusDropping`, `MongoClusterStatusProvisioning`, `MongoClusterStatusReady`, `MongoClusterStatusStarting`, `MongoClusterStatusStopped`, `MongoClusterStatusStopping`, `MongoClusterStatusUpdating`
- New enum type `NodeKind` with values `NodeKindShard`
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateDropping`, `ProvisioningStateFailed`, `ProvisioningStateInProgress`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New enum type `ThroughputPolicyType` with values `ThroughputPolicyTypeCustom`, `ThroughputPolicyTypeEqual`, `ThroughputPolicyTypeNone`
- New function `*AzureBlobDataTransferDataSourceSink.GetDataTransferDataSourceSink() *DataTransferDataSourceSink`
- New function `*CassandraClustersClient.GetBackup(context.Context, string, string, string, *CassandraClustersClientGetBackupOptions) (CassandraClustersClientGetBackupResponse, error)`
- New function `*CassandraClustersClient.NewListBackupsPager(string, string, *CassandraClustersClientListBackupsOptions) *runtime.Pager[CassandraClustersClientListBackupsResponse]`
- New function `*CassandraDataTransferDataSourceSink.GetDataTransferDataSourceSink() *DataTransferDataSourceSink`
- New function `*CassandraResourcesClient.BeginCreateUpdateCassandraView(context.Context, string, string, string, string, CassandraViewCreateUpdateParameters, *CassandraResourcesClientBeginCreateUpdateCassandraViewOptions) (*runtime.Poller[CassandraResourcesClientCreateUpdateCassandraViewResponse], error)`
- New function `*CassandraResourcesClient.BeginDeleteCassandraView(context.Context, string, string, string, string, *CassandraResourcesClientBeginDeleteCassandraViewOptions) (*runtime.Poller[CassandraResourcesClientDeleteCassandraViewResponse], error)`
- New function `*CassandraResourcesClient.GetCassandraView(context.Context, string, string, string, string, *CassandraResourcesClientGetCassandraViewOptions) (CassandraResourcesClientGetCassandraViewResponse, error)`
- New function `*CassandraResourcesClient.GetCassandraViewThroughput(context.Context, string, string, string, string, *CassandraResourcesClientGetCassandraViewThroughputOptions) (CassandraResourcesClientGetCassandraViewThroughputResponse, error)`
- New function `*CassandraResourcesClient.NewListCassandraViewsPager(string, string, string, *CassandraResourcesClientListCassandraViewsOptions) *runtime.Pager[CassandraResourcesClientListCassandraViewsResponse]`
- New function `*CassandraResourcesClient.BeginMigrateCassandraViewToAutoscale(context.Context, string, string, string, string, *CassandraResourcesClientBeginMigrateCassandraViewToAutoscaleOptions) (*runtime.Poller[CassandraResourcesClientMigrateCassandraViewToAutoscaleResponse], error)`
- New function `*CassandraResourcesClient.BeginMigrateCassandraViewToManualThroughput(context.Context, string, string, string, string, *CassandraResourcesClientBeginMigrateCassandraViewToManualThroughputOptions) (*runtime.Poller[CassandraResourcesClientMigrateCassandraViewToManualThroughputResponse], error)`
- New function `*CassandraResourcesClient.BeginUpdateCassandraViewThroughput(context.Context, string, string, string, string, ThroughputSettingsUpdateParameters, *CassandraResourcesClientBeginUpdateCassandraViewThroughputOptions) (*runtime.Poller[CassandraResourcesClientUpdateCassandraViewThroughputResponse], error)`
- New function `*ClientFactory.NewDataTransferJobsClient() *DataTransferJobsClient`
- New function `*ClientFactory.NewGraphResourcesClient() *GraphResourcesClient`
- New function `*ClientFactory.NewMongoClustersClient() *MongoClustersClient`
- New function `*DataTransferDataSourceSink.GetDataTransferDataSourceSink() *DataTransferDataSourceSink`
- New function `NewDataTransferJobsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DataTransferJobsClient, error)`
- New function `*DataTransferJobsClient.Cancel(context.Context, string, string, string, *DataTransferJobsClientCancelOptions) (DataTransferJobsClientCancelResponse, error)`
- New function `*DataTransferJobsClient.Create(context.Context, string, string, string, CreateJobRequest, *DataTransferJobsClientCreateOptions) (DataTransferJobsClientCreateResponse, error)`
- New function `*DataTransferJobsClient.Get(context.Context, string, string, string, *DataTransferJobsClientGetOptions) (DataTransferJobsClientGetResponse, error)`
- New function `*DataTransferJobsClient.NewListByDatabaseAccountPager(string, string, *DataTransferJobsClientListByDatabaseAccountOptions) *runtime.Pager[DataTransferJobsClientListByDatabaseAccountResponse]`
- New function `*DataTransferJobsClient.Pause(context.Context, string, string, string, *DataTransferJobsClientPauseOptions) (DataTransferJobsClientPauseResponse, error)`
- New function `*DataTransferJobsClient.Resume(context.Context, string, string, string, *DataTransferJobsClientResumeOptions) (DataTransferJobsClientResumeResponse, error)`
- New function `NewGraphResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GraphResourcesClient, error)`
- New function `*GraphResourcesClient.BeginCreateUpdateGraph(context.Context, string, string, string, GraphResourceCreateUpdateParameters, *GraphResourcesClientBeginCreateUpdateGraphOptions) (*runtime.Poller[GraphResourcesClientCreateUpdateGraphResponse], error)`
- New function `*GraphResourcesClient.BeginDeleteGraphResource(context.Context, string, string, string, *GraphResourcesClientBeginDeleteGraphResourceOptions) (*runtime.Poller[GraphResourcesClientDeleteGraphResourceResponse], error)`
- New function `*GraphResourcesClient.GetGraph(context.Context, string, string, string, *GraphResourcesClientGetGraphOptions) (GraphResourcesClientGetGraphResponse, error)`
- New function `*GraphResourcesClient.NewListGraphsPager(string, string, *GraphResourcesClientListGraphsOptions) *runtime.Pager[GraphResourcesClientListGraphsResponse]`
- New function `NewMongoClustersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MongoClustersClient, error)`
- New function `*MongoClustersClient.CheckNameAvailability(context.Context, string, CheckNameAvailabilityRequest, *MongoClustersClientCheckNameAvailabilityOptions) (MongoClustersClientCheckNameAvailabilityResponse, error)`
- New function `*MongoClustersClient.BeginCreateOrUpdate(context.Context, string, string, MongoCluster, *MongoClustersClientBeginCreateOrUpdateOptions) (*runtime.Poller[MongoClustersClientCreateOrUpdateResponse], error)`
- New function `*MongoClustersClient.BeginCreateOrUpdateFirewallRule(context.Context, string, string, string, FirewallRule, *MongoClustersClientBeginCreateOrUpdateFirewallRuleOptions) (*runtime.Poller[MongoClustersClientCreateOrUpdateFirewallRuleResponse], error)`
- New function `*MongoClustersClient.BeginDelete(context.Context, string, string, *MongoClustersClientBeginDeleteOptions) (*runtime.Poller[MongoClustersClientDeleteResponse], error)`
- New function `*MongoClustersClient.BeginDeleteFirewallRule(context.Context, string, string, string, *MongoClustersClientBeginDeleteFirewallRuleOptions) (*runtime.Poller[MongoClustersClientDeleteFirewallRuleResponse], error)`
- New function `*MongoClustersClient.Get(context.Context, string, string, *MongoClustersClientGetOptions) (MongoClustersClientGetResponse, error)`
- New function `*MongoClustersClient.GetFirewallRule(context.Context, string, string, string, *MongoClustersClientGetFirewallRuleOptions) (MongoClustersClientGetFirewallRuleResponse, error)`
- New function `*MongoClustersClient.NewListByResourceGroupPager(string, *MongoClustersClientListByResourceGroupOptions) *runtime.Pager[MongoClustersClientListByResourceGroupResponse]`
- New function `*MongoClustersClient.ListConnectionStrings(context.Context, string, string, *MongoClustersClientListConnectionStringsOptions) (MongoClustersClientListConnectionStringsResponse, error)`
- New function `*MongoClustersClient.NewListFirewallRulesPager(string, string, *MongoClustersClientListFirewallRulesOptions) *runtime.Pager[MongoClustersClientListFirewallRulesResponse]`
- New function `*MongoClustersClient.NewListPager(*MongoClustersClientListOptions) *runtime.Pager[MongoClustersClientListResponse]`
- New function `*MongoClustersClient.BeginUpdate(context.Context, string, string, MongoClusterUpdate, *MongoClustersClientBeginUpdateOptions) (*runtime.Poller[MongoClustersClientUpdateResponse], error)`
- New function `*MongoDBResourcesClient.BeginListMongoDBCollectionPartitionMerge(context.Context, string, string, string, string, MergeParameters, *MongoDBResourcesClientBeginListMongoDBCollectionPartitionMergeOptions) (*runtime.Poller[MongoDBResourcesClientListMongoDBCollectionPartitionMergeResponse], error)`
- New function `*MongoDBResourcesClient.BeginMongoDBContainerRedistributeThroughput(context.Context, string, string, string, string, RedistributeThroughputParameters, *MongoDBResourcesClientBeginMongoDBContainerRedistributeThroughputOptions) (*runtime.Poller[MongoDBResourcesClientMongoDBContainerRedistributeThroughputResponse], error)`
- New function `*MongoDBResourcesClient.BeginMongoDBContainerRetrieveThroughputDistribution(context.Context, string, string, string, string, RetrieveThroughputParameters, *MongoDBResourcesClientBeginMongoDBContainerRetrieveThroughputDistributionOptions) (*runtime.Poller[MongoDBResourcesClientMongoDBContainerRetrieveThroughputDistributionResponse], error)`
- New function `*MongoDBResourcesClient.BeginMongoDBDatabasePartitionMerge(context.Context, string, string, string, MergeParameters, *MongoDBResourcesClientBeginMongoDBDatabasePartitionMergeOptions) (*runtime.Poller[MongoDBResourcesClientMongoDBDatabasePartitionMergeResponse], error)`
- New function `*MongoDBResourcesClient.BeginMongoDBDatabaseRedistributeThroughput(context.Context, string, string, string, RedistributeThroughputParameters, *MongoDBResourcesClientBeginMongoDBDatabaseRedistributeThroughputOptions) (*runtime.Poller[MongoDBResourcesClientMongoDBDatabaseRedistributeThroughputResponse], error)`
- New function `*MongoDBResourcesClient.BeginMongoDBDatabaseRetrieveThroughputDistribution(context.Context, string, string, string, RetrieveThroughputParameters, *MongoDBResourcesClientBeginMongoDBDatabaseRetrieveThroughputDistributionOptions) (*runtime.Poller[MongoDBResourcesClientMongoDBDatabaseRetrieveThroughputDistributionResponse], error)`
- New function `*MongoDataTransferDataSourceSink.GetDataTransferDataSourceSink() *DataTransferDataSourceSink`
- New function `*SQLDataTransferDataSourceSink.GetDataTransferDataSourceSink() *DataTransferDataSourceSink`
- New function `*SQLResourcesClient.BeginListSQLContainerPartitionMerge(context.Context, string, string, string, string, MergeParameters, *SQLResourcesClientBeginListSQLContainerPartitionMergeOptions) (*runtime.Poller[SQLResourcesClientListSQLContainerPartitionMergeResponse], error)`
- New function `*SQLResourcesClient.BeginSQLContainerRedistributeThroughput(context.Context, string, string, string, string, RedistributeThroughputParameters, *SQLResourcesClientBeginSQLContainerRedistributeThroughputOptions) (*runtime.Poller[SQLResourcesClientSQLContainerRedistributeThroughputResponse], error)`
- New function `*SQLResourcesClient.BeginSQLContainerRetrieveThroughputDistribution(context.Context, string, string, string, string, RetrieveThroughputParameters, *SQLResourcesClientBeginSQLContainerRetrieveThroughputDistributionOptions) (*runtime.Poller[SQLResourcesClientSQLContainerRetrieveThroughputDistributionResponse], error)`
- New function `*SQLResourcesClient.BeginSQLDatabasePartitionMerge(context.Context, string, string, string, MergeParameters, *SQLResourcesClientBeginSQLDatabasePartitionMergeOptions) (*runtime.Poller[SQLResourcesClientSQLDatabasePartitionMergeResponse], error)`
- New function `*SQLResourcesClient.BeginSQLDatabaseRedistributeThroughput(context.Context, string, string, string, RedistributeThroughputParameters, *SQLResourcesClientBeginSQLDatabaseRedistributeThroughputOptions) (*runtime.Poller[SQLResourcesClientSQLDatabaseRedistributeThroughputResponse], error)`
- New function `*SQLResourcesClient.BeginSQLDatabaseRetrieveThroughputDistribution(context.Context, string, string, string, RetrieveThroughputParameters, *SQLResourcesClientBeginSQLDatabaseRetrieveThroughputDistributionOptions) (*runtime.Poller[SQLResourcesClientSQLDatabaseRetrieveThroughputDistributionResponse], error)`
- New struct `AzureBlobDataTransferDataSourceSink`
- New struct `BackupResource`
- New struct `BackupResourceProperties`
- New struct `CassandraDataTransferDataSourceSink`
- New struct `CassandraViewCreateUpdateParameters`
- New struct `CassandraViewCreateUpdateProperties`
- New struct `CassandraViewGetProperties`
- New struct `CassandraViewGetPropertiesOptions`
- New struct `CassandraViewGetPropertiesResource`
- New struct `CassandraViewGetResults`
- New struct `CassandraViewListResult`
- New struct `CassandraViewResource`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityResponse`
- New struct `ConnectionString`
- New struct `CreateJobRequest`
- New struct `DataTransferJobFeedResults`
- New struct `DataTransferJobGetResults`
- New struct `DataTransferJobProperties`
- New struct `DiagnosticLogSettings`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponseAutoGenerated`
- New struct `FirewallRule`
- New struct `FirewallRuleListResult`
- New struct `FirewallRuleProperties`
- New struct `GraphResource`
- New struct `GraphResourceCreateUpdateParameters`
- New struct `GraphResourceCreateUpdateProperties`
- New struct `GraphResourceGetProperties`
- New struct `GraphResourceGetPropertiesOptions`
- New struct `GraphResourceGetPropertiesResource`
- New struct `GraphResourceGetResults`
- New struct `GraphResourcesListResult`
- New struct `ListBackups`
- New struct `ListConnectionStringsResult`
- New struct `MaterializedViewDefinition`
- New struct `MergeParameters`
- New struct `MongoCluster`
- New struct `MongoClusterListResult`
- New struct `MongoClusterProperties`
- New struct `MongoClusterRestoreParameters`
- New struct `MongoClusterUpdate`
- New struct `MongoDataTransferDataSourceSink`
- New struct `NodeGroupProperties`
- New struct `NodeGroupSpec`
- New struct `PhysicalPartitionID`
- New struct `PhysicalPartitionStorageInfo`
- New struct `PhysicalPartitionStorageInfoCollection`
- New struct `PhysicalPartitionThroughputInfoProperties`
- New struct `PhysicalPartitionThroughputInfoResource`
- New struct `PhysicalPartitionThroughputInfoResult`
- New struct `PhysicalPartitionThroughputInfoResultProperties`
- New struct `PhysicalPartitionThroughputInfoResultPropertiesResource`
- New struct `ProxyResourceAutoGenerated`
- New struct `RedistributeThroughputParameters`
- New struct `RedistributeThroughputProperties`
- New struct `RedistributeThroughputPropertiesResource`
- New struct `ResourceAutoGenerated`
- New struct `ResourceRestoreParameters`
- New struct `RestoreParametersBase`
- New struct `RetrieveThroughputParameters`
- New struct `RetrieveThroughputProperties`
- New struct `RetrieveThroughputPropertiesResource`
- New struct `SQLDataTransferDataSourceSink`
- New struct `TrackedResource`
- New field `Identity` in struct `ARMResourceProperties`
- New field `Identity` in struct `CassandraKeyspaceCreateUpdateParameters`
- New field `Identity` in struct `CassandraKeyspaceGetResults`
- New field `Identity` in struct `CassandraTableCreateUpdateParameters`
- New field `Identity` in struct `CassandraTableGetResults`
- New field `DiagnosticLogSettings`, `EnableBurstCapacity`, `EnableMaterializedViews` in struct `DatabaseAccountCreateUpdateProperties`
- New field `DiagnosticLogSettings`, `EnableBurstCapacity`, `EnableMaterializedViews` in struct `DatabaseAccountGetProperties`
- New field `DiagnosticLogSettings`, `EnableBurstCapacity`, `EnableMaterializedViews` in struct `DatabaseAccountUpdateProperties`
- New field `Identity` in struct `GremlinDatabaseCreateUpdateParameters`
- New field `CreateMode`, `RestoreParameters` in struct `GremlinDatabaseGetPropertiesResource`
- New field `Identity` in struct `GremlinDatabaseGetResults`
- New field `CreateMode`, `RestoreParameters` in struct `GremlinDatabaseResource`
- New field `Identity` in struct `GremlinGraphCreateUpdateParameters`
- New field `CreateMode`, `RestoreParameters` in struct `GremlinGraphGetPropertiesResource`
- New field `Identity` in struct `GremlinGraphGetResults`
- New field `CreateMode`, `RestoreParameters` in struct `GremlinGraphResource`
- New field `Identity` in struct `MongoDBCollectionCreateUpdateParameters`
- New field `CreateMode`, `RestoreParameters` in struct `MongoDBCollectionGetPropertiesResource`
- New field `Identity` in struct `MongoDBCollectionGetResults`
- New field `CreateMode`, `RestoreParameters` in struct `MongoDBCollectionResource`
- New field `Identity` in struct `MongoDBDatabaseCreateUpdateParameters`
- New field `CreateMode`, `RestoreParameters` in struct `MongoDBDatabaseGetPropertiesResource`
- New field `Identity` in struct `MongoDBDatabaseGetResults`
- New field `CreateMode`, `RestoreParameters` in struct `MongoDBDatabaseResource`
- New field `CreateMode`, `MaterializedViewDefinition`, `RestoreParameters` in struct `RestorableSQLContainerPropertiesResourceContainer`
- New field `CreateMode`, `RestoreParameters` in struct `RestorableSQLDatabasePropertiesResourceDatabase`
- New field `SourceBackupLocation` in struct `RestoreParameters`
- New field `Identity` in struct `SQLContainerCreateUpdateParameters`
- New field `CreateMode`, `MaterializedViewDefinition`, `RestoreParameters` in struct `SQLContainerGetPropertiesResource`
- New field `Identity` in struct `SQLContainerGetResults`
- New field `CreateMode`, `MaterializedViewDefinition`, `RestoreParameters` in struct `SQLContainerResource`
- New field `Identity` in struct `SQLDatabaseCreateUpdateParameters`
- New field `CreateMode`, `RestoreParameters` in struct `SQLDatabaseGetPropertiesResource`
- New field `Identity` in struct `SQLDatabaseGetResults`
- New field `CreateMode`, `RestoreParameters` in struct `SQLDatabaseResource`
- New field `Identity` in struct `SQLStoredProcedureCreateUpdateParameters`
- New field `Identity` in struct `SQLStoredProcedureGetResults`
- New field `Identity` in struct `SQLTriggerCreateUpdateParameters`
- New field `Identity` in struct `SQLTriggerGetResults`
- New field `Identity` in struct `SQLUserDefinedFunctionCreateUpdateParameters`
- New field `Identity` in struct `SQLUserDefinedFunctionGetResults`
- New field `Identity` in struct `TableCreateUpdateParameters`
- New field `CreateMode`, `RestoreParameters` in struct `TableGetPropertiesResource`
- New field `Identity` in struct `TableGetResults`
- New field `CreateMode`, `RestoreParameters` in struct `TableResource`
- New field `Identity` in struct `ThroughputSettingsGetResults`
- New field `Identity` in struct `ThroughputSettingsUpdateParameters`


## 2.5.0 (2023-05-26)
### Features Added

- New value `PublicNetworkAccessSecuredByPerimeter` added to enum type `PublicNetworkAccess`
- New enum type `ContinuousTier` with values `ContinuousTierContinuous30Days`, `ContinuousTierContinuous7Days`
- New struct `ContinuousModeProperties`
- New field `ContinuousModeProperties` in struct `ContinuousModeBackupPolicy`
- New field `OldestRestorableTime` in struct `RestorableDatabaseAccountProperties`
- New field `InstantMaximumThroughput`, `SoftAllowedMaximumThroughput` in struct `ThroughputSettingsGetPropertiesResource`
- New field `InstantMaximumThroughput`, `SoftAllowedMaximumThroughput` in struct `ThroughputSettingsResource`


## 2.4.0 (2023-04-28)
### Features Added

- New value `AuthenticationMethodLdap` added to enum type `AuthenticationMethod`
- New enum type `Kind` with values `KindPrimary`, `KindPrimaryReadonly`, `KindSecondary`, `KindSecondaryReadonly`
- New enum type `Status` with values `StatusDeleting`, `StatusInitializing`, `StatusInternallyReady`, `StatusOnline`, `StatusUninitialized`
- New enum type `Type` with values `TypeCassandra`, `TypeCassandraConnectorMetadata`, `TypeGremlin`, `TypeGremlinV2`, `TypeMongoDB`, `TypeSQL`, `TypeSQLDedicatedGateway`, `TypeTable`, `TypeUndefined`
- New struct `AuthenticationMethodLdapProperties`
- New struct `CassandraError`
- New field `Errors` in struct `CassandraClusterPublicStatus`
- New field `ProvisionError` in struct `ClusterResourceProperties`
- New field `CassandraProcessStatus` in struct `ComponentsM9L909SchemasCassandraclusterpublicstatusPropertiesDatacentersItemsPropertiesNodesItems`
- New field `AuthenticationMethodLdapProperties`, `Deallocated`, `ProvisionError` in struct `DataCenterResourceProperties`
- New field `KeyKind`, `Type` in struct `DatabaseAccountConnectionString`
- New field `IsSubscriptionRegionAccessAllowedForAz`, `IsSubscriptionRegionAccessAllowedForRegular`, `Status` in struct `LocationProperties`


## 2.3.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module

## 2.2.0 (2023-02-24)
### Features Added

- New type alias `MinimalTLSVersion` with values `MinimalTLSVersionTLS`, `MinimalTLSVersionTls11`, `MinimalTLSVersionTls12`
- New function `*GremlinResourcesClient.BeginRetrieveContinuousBackupInformation(context.Context, string, string, string, string, ContinuousBackupRestoreLocation, *GremlinResourcesClientBeginRetrieveContinuousBackupInformationOptions) (*runtime.Poller[GremlinResourcesClientRetrieveContinuousBackupInformationResponse], error)`
- New function `NewRestorableGremlinDatabasesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RestorableGremlinDatabasesClient, error)`
- New function `*RestorableGremlinDatabasesClient.NewListPager(string, string, *RestorableGremlinDatabasesClientListOptions) *runtime.Pager[RestorableGremlinDatabasesClientListResponse]`
- New function `NewRestorableGremlinGraphsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RestorableGremlinGraphsClient, error)`
- New function `*RestorableGremlinGraphsClient.NewListPager(string, string, *RestorableGremlinGraphsClientListOptions) *runtime.Pager[RestorableGremlinGraphsClientListResponse]`
- New function `NewRestorableGremlinResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RestorableGremlinResourcesClient, error)`
- New function `*RestorableGremlinResourcesClient.NewListPager(string, string, *RestorableGremlinResourcesClientListOptions) *runtime.Pager[RestorableGremlinResourcesClientListResponse]`
- New function `NewRestorableTableResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RestorableTableResourcesClient, error)`
- New function `*RestorableTableResourcesClient.NewListPager(string, string, *RestorableTableResourcesClientListOptions) *runtime.Pager[RestorableTableResourcesClientListResponse]`
- New function `NewRestorableTablesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RestorableTablesClient, error)`
- New function `*RestorableTablesClient.NewListPager(string, string, *RestorableTablesClientListOptions) *runtime.Pager[RestorableTablesClientListResponse]`
- New function `*SQLResourcesClient.BeginCreateUpdateClientEncryptionKey(context.Context, string, string, string, string, ClientEncryptionKeyCreateUpdateParameters, *SQLResourcesClientBeginCreateUpdateClientEncryptionKeyOptions) (*runtime.Poller[SQLResourcesClientCreateUpdateClientEncryptionKeyResponse], error)`
- New function `*SQLResourcesClient.GetClientEncryptionKey(context.Context, string, string, string, string, *SQLResourcesClientGetClientEncryptionKeyOptions) (SQLResourcesClientGetClientEncryptionKeyResponse, error)`
- New function `*SQLResourcesClient.NewListClientEncryptionKeysPager(string, string, string, *SQLResourcesClientListClientEncryptionKeysOptions) *runtime.Pager[SQLResourcesClientListClientEncryptionKeysResponse]`
- New function `*TableResourcesClient.BeginRetrieveContinuousBackupInformation(context.Context, string, string, string, ContinuousBackupRestoreLocation, *TableResourcesClientBeginRetrieveContinuousBackupInformationOptions) (*runtime.Poller[TableResourcesClientRetrieveContinuousBackupInformationResponse], error)`
- New struct `ClientEncryptionIncludedPath`
- New struct `ClientEncryptionKeyCreateUpdateParameters`
- New struct `ClientEncryptionKeyCreateUpdateProperties`
- New struct `ClientEncryptionKeyGetProperties`
- New struct `ClientEncryptionKeyGetPropertiesResource`
- New struct `ClientEncryptionKeyGetResults`
- New struct `ClientEncryptionKeyResource`
- New struct `ClientEncryptionKeysListResult`
- New struct `ClientEncryptionPolicy`
- New struct `GremlinDatabaseRestoreResource`
- New struct `KeyWrapMetadata`
- New struct `RestorableGremlinDatabaseGetResult`
- New struct `RestorableGremlinDatabaseProperties`
- New struct `RestorableGremlinDatabasePropertiesResource`
- New struct `RestorableGremlinDatabasesClient`
- New struct `RestorableGremlinDatabasesListResult`
- New struct `RestorableGremlinGraphGetResult`
- New struct `RestorableGremlinGraphProperties`
- New struct `RestorableGremlinGraphPropertiesResource`
- New struct `RestorableGremlinGraphsClient`
- New struct `RestorableGremlinGraphsListResult`
- New struct `RestorableGremlinResourcesClient`
- New struct `RestorableGremlinResourcesGetResult`
- New struct `RestorableGremlinResourcesListResult`
- New struct `RestorableTableGetResult`
- New struct `RestorableTableProperties`
- New struct `RestorableTablePropertiesResource`
- New struct `RestorableTableResourcesClient`
- New struct `RestorableTableResourcesGetResult`
- New struct `RestorableTableResourcesListResult`
- New struct `RestorableTablesClient`
- New struct `RestorableTablesListResult`
- New field `MinimalTLSVersion` in struct `DatabaseAccountCreateUpdateProperties`
- New field `MinimalTLSVersion` in struct `DatabaseAccountGetProperties`
- New field `MinimalTLSVersion` in struct `DatabaseAccountUpdateProperties`
- New field `EndTime` in struct `RestorableMongodbCollectionsClientListOptions`
- New field `StartTime` in struct `RestorableMongodbCollectionsClientListOptions`
- New field `ClientEncryptionPolicy` in struct `RestorableSQLContainerPropertiesResourceContainer`
- New field `GremlinDatabasesToRestore` in struct `RestoreParameters`
- New field `TablesToRestore` in struct `RestoreParameters`
- New field `ClientEncryptionPolicy` in struct `SQLContainerGetPropertiesResource`
- New field `ClientEncryptionPolicy` in struct `SQLContainerResource`


## 2.1.0 (2022-09-06)
### Features Added

- New const `MongoRoleDefinitionTypeBuiltInRole`
- New const `MongoRoleDefinitionTypeCustomRole`
- New type alias `MongoRoleDefinitionType`
- New function `*MongoDBResourcesClient.BeginCreateUpdateMongoRoleDefinition(context.Context, string, string, string, MongoRoleDefinitionCreateUpdateParameters, *MongoDBResourcesClientBeginCreateUpdateMongoRoleDefinitionOptions) (*runtime.Poller[MongoDBResourcesClientCreateUpdateMongoRoleDefinitionResponse], error)`
- New function `*MongoDBResourcesClient.BeginCreateUpdateMongoUserDefinition(context.Context, string, string, string, MongoUserDefinitionCreateUpdateParameters, *MongoDBResourcesClientBeginCreateUpdateMongoUserDefinitionOptions) (*runtime.Poller[MongoDBResourcesClientCreateUpdateMongoUserDefinitionResponse], error)`
- New function `*MongoDBResourcesClient.GetMongoRoleDefinition(context.Context, string, string, string, *MongoDBResourcesClientGetMongoRoleDefinitionOptions) (MongoDBResourcesClientGetMongoRoleDefinitionResponse, error)`
- New function `*MongoDBResourcesClient.NewListMongoUserDefinitionsPager(string, string, *MongoDBResourcesClientListMongoUserDefinitionsOptions) *runtime.Pager[MongoDBResourcesClientListMongoUserDefinitionsResponse]`
- New function `*MongoDBResourcesClient.BeginDeleteMongoRoleDefinition(context.Context, string, string, string, *MongoDBResourcesClientBeginDeleteMongoRoleDefinitionOptions) (*runtime.Poller[MongoDBResourcesClientDeleteMongoRoleDefinitionResponse], error)`
- New function `*MongoDBResourcesClient.NewListMongoRoleDefinitionsPager(string, string, *MongoDBResourcesClientListMongoRoleDefinitionsOptions) *runtime.Pager[MongoDBResourcesClientListMongoRoleDefinitionsResponse]`
- New function `PossibleMongoRoleDefinitionTypeValues() []MongoRoleDefinitionType`
- New function `*MongoDBResourcesClient.GetMongoUserDefinition(context.Context, string, string, string, *MongoDBResourcesClientGetMongoUserDefinitionOptions) (MongoDBResourcesClientGetMongoUserDefinitionResponse, error)`
- New function `*MongoDBResourcesClient.BeginDeleteMongoUserDefinition(context.Context, string, string, string, *MongoDBResourcesClientBeginDeleteMongoUserDefinitionOptions) (*runtime.Poller[MongoDBResourcesClientDeleteMongoUserDefinitionResponse], error)`
- New struct `AccountKeyMetadata`
- New struct `DatabaseAccountKeysMetadata`
- New struct `MongoDBResourcesClientBeginCreateUpdateMongoRoleDefinitionOptions`
- New struct `MongoDBResourcesClientBeginCreateUpdateMongoUserDefinitionOptions`
- New struct `MongoDBResourcesClientBeginDeleteMongoRoleDefinitionOptions`
- New struct `MongoDBResourcesClientBeginDeleteMongoUserDefinitionOptions`
- New struct `MongoDBResourcesClientCreateUpdateMongoRoleDefinitionResponse`
- New struct `MongoDBResourcesClientCreateUpdateMongoUserDefinitionResponse`
- New struct `MongoDBResourcesClientDeleteMongoRoleDefinitionResponse`
- New struct `MongoDBResourcesClientDeleteMongoUserDefinitionResponse`
- New struct `MongoDBResourcesClientGetMongoRoleDefinitionOptions`
- New struct `MongoDBResourcesClientGetMongoRoleDefinitionResponse`
- New struct `MongoDBResourcesClientGetMongoUserDefinitionOptions`
- New struct `MongoDBResourcesClientGetMongoUserDefinitionResponse`
- New struct `MongoDBResourcesClientListMongoRoleDefinitionsOptions`
- New struct `MongoDBResourcesClientListMongoRoleDefinitionsResponse`
- New struct `MongoDBResourcesClientListMongoUserDefinitionsOptions`
- New struct `MongoDBResourcesClientListMongoUserDefinitionsResponse`
- New struct `MongoRoleDefinitionCreateUpdateParameters`
- New struct `MongoRoleDefinitionGetResults`
- New struct `MongoRoleDefinitionListResult`
- New struct `MongoRoleDefinitionResource`
- New struct `MongoUserDefinitionCreateUpdateParameters`
- New struct `MongoUserDefinitionGetResults`
- New struct `MongoUserDefinitionListResult`
- New struct `MongoUserDefinitionResource`
- New struct `Privilege`
- New struct `PrivilegeResource`
- New struct `Role`
- New field `EnablePartitionMerge` in struct `DatabaseAccountGetProperties`
- New field `KeysMetadata` in struct `DatabaseAccountGetProperties`
- New field `KeysMetadata` in struct `DatabaseAccountUpdateProperties`
- New field `EnablePartitionMerge` in struct `DatabaseAccountUpdateProperties`
- New field `KeysMetadata` in struct `DatabaseAccountCreateUpdateProperties`
- New field `EnablePartitionMerge` in struct `DatabaseAccountCreateUpdateProperties`


## 2.0.0 (2022-07-18)
### Breaking Changes

- Type of `RestorableMongodbResourcesListResult.Value` has been changed from `[]*DatabaseRestoreResource` to `[]*RestorableMongodbResourcesGetResult`
- Type of `RestorableSQLResourcesListResult.Value` has been changed from `[]*DatabaseRestoreResource` to `[]*RestorableSQLResourcesGetResult`

### Features Added

- New const `ServiceStatusRunning`
- New const `ServiceSizeCosmosD16S`
- New const `ServiceTypeSQLDedicatedGateway`
- New const `ServiceTypeGraphAPICompute`
- New const `ServiceStatusDeleting`
- New const `ServiceSizeCosmosD8S`
- New const `ServiceStatusStopped`
- New const `ServiceStatusCreating`
- New const `ServiceStatusUpdating`
- New const `ServiceSizeCosmosD4S`
- New const `ServiceTypeMaterializedViewsBuilder`
- New const `ServiceStatusError`
- New const `ServiceTypeDataTransfer`
- New function `*SQLDedicatedGatewayServiceResourceProperties.GetServiceResourceProperties() *ServiceResourceProperties`
- New function `PossibleServiceTypeValues() []ServiceType`
- New function `PossibleServiceStatusValues() []ServiceStatus`
- New function `*ServiceClient.BeginCreate(context.Context, string, string, string, ServiceResourceCreateUpdateParameters, *ServiceClientBeginCreateOptions) (*runtime.Poller[ServiceClientCreateResponse], error)`
- New function `*ServiceClient.NewListPager(string, string, *ServiceClientListOptions) *runtime.Pager[ServiceClientListResponse]`
- New function `*ServiceClient.BeginDelete(context.Context, string, string, string, *ServiceClientBeginDeleteOptions) (*runtime.Poller[ServiceClientDeleteResponse], error)`
- New function `*GraphAPIComputeServiceResourceProperties.GetServiceResourceProperties() *ServiceResourceProperties`
- New function `*ServiceClient.Get(context.Context, string, string, string, *ServiceClientGetOptions) (ServiceClientGetResponse, error)`
- New function `*MaterializedViewsBuilderServiceResourceProperties.GetServiceResourceProperties() *ServiceResourceProperties`
- New function `*DataTransferServiceResourceProperties.GetServiceResourceProperties() *ServiceResourceProperties`
- New function `*ServiceResourceProperties.GetServiceResourceProperties() *ServiceResourceProperties`
- New function `PossibleServiceSizeValues() []ServiceSize`
- New function `NewServiceClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServiceClient, error)`
- New struct `DataTransferRegionalServiceResource`
- New struct `DataTransferServiceResource`
- New struct `DataTransferServiceResourceProperties`
- New struct `GraphAPIComputeRegionalServiceResource`
- New struct `GraphAPIComputeServiceResource`
- New struct `GraphAPIComputeServiceResourceProperties`
- New struct `MaterializedViewsBuilderRegionalServiceResource`
- New struct `MaterializedViewsBuilderServiceResource`
- New struct `MaterializedViewsBuilderServiceResourceProperties`
- New struct `RegionalServiceResource`
- New struct `RestorableMongodbResourcesGetResult`
- New struct `RestorableSQLResourcesGetResult`
- New struct `SQLDedicatedGatewayRegionalServiceResource`
- New struct `SQLDedicatedGatewayServiceResource`
- New struct `SQLDedicatedGatewayServiceResourceProperties`
- New struct `ServiceClient`
- New struct `ServiceClientBeginCreateOptions`
- New struct `ServiceClientBeginDeleteOptions`
- New struct `ServiceClientCreateResponse`
- New struct `ServiceClientDeleteResponse`
- New struct `ServiceClientGetOptions`
- New struct `ServiceClientGetResponse`
- New struct `ServiceClientListOptions`
- New struct `ServiceClientListResponse`
- New struct `ServiceResource`
- New struct `ServiceResourceCreateUpdateParameters`
- New struct `ServiceResourceCreateUpdateProperties`
- New struct `ServiceResourceListResult`
- New struct `ServiceResourceProperties`
- New field `AnalyticalStorageTTL` in struct `GremlinGraphGetPropertiesResource`
- New field `AnalyticalStorageTTL` in struct `GremlinGraphResource`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).