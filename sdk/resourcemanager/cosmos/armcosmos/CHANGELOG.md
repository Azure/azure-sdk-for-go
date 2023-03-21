# Release History

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