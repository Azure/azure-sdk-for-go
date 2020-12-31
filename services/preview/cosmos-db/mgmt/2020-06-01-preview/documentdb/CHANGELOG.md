Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Constants

1. DatabaseAccountKind.GlobalDocumentDB
1. DatabaseAccountKind.Parse
1. TriggerOperation.All
1. TriggerOperation.Update

### Removed Funcs

1. *CassandraResourcesCreateUpdateCassandraKeyspaceFuture.Result(CassandraResourcesClient) (CassandraKeyspaceGetResults, error)
1. *CassandraResourcesCreateUpdateCassandraTableFuture.Result(CassandraResourcesClient) (CassandraTableGetResults, error)
1. *CassandraResourcesDeleteCassandraKeyspaceFuture.Result(CassandraResourcesClient) (autorest.Response, error)
1. *CassandraResourcesDeleteCassandraTableFuture.Result(CassandraResourcesClient) (autorest.Response, error)
1. *CassandraResourcesUpdateCassandraKeyspaceThroughputFuture.Result(CassandraResourcesClient) (ThroughputSettingsGetResults, error)
1. *CassandraResourcesUpdateCassandraTableThroughputFuture.Result(CassandraResourcesClient) (ThroughputSettingsGetResults, error)
1. *DatabaseAccountsCreateOrUpdateFuture.Result(DatabaseAccountsClient) (DatabaseAccountGetResults, error)
1. *DatabaseAccountsDeleteFuture.Result(DatabaseAccountsClient) (autorest.Response, error)
1. *DatabaseAccountsFailoverPriorityChangeFuture.Result(DatabaseAccountsClient) (autorest.Response, error)
1. *DatabaseAccountsOfflineRegionFuture.Result(DatabaseAccountsClient) (autorest.Response, error)
1. *DatabaseAccountsOnlineRegionFuture.Result(DatabaseAccountsClient) (autorest.Response, error)
1. *DatabaseAccountsRegenerateKeyFuture.Result(DatabaseAccountsClient) (autorest.Response, error)
1. *DatabaseAccountsUpdateFuture.Result(DatabaseAccountsClient) (DatabaseAccountGetResults, error)
1. *GremlinResourcesCreateUpdateGremlinDatabaseFuture.Result(GremlinResourcesClient) (GremlinDatabaseGetResults, error)
1. *GremlinResourcesCreateUpdateGremlinGraphFuture.Result(GremlinResourcesClient) (GremlinGraphGetResults, error)
1. *GremlinResourcesDeleteGremlinDatabaseFuture.Result(GremlinResourcesClient) (autorest.Response, error)
1. *GremlinResourcesDeleteGremlinGraphFuture.Result(GremlinResourcesClient) (autorest.Response, error)
1. *GremlinResourcesUpdateGremlinDatabaseThroughputFuture.Result(GremlinResourcesClient) (ThroughputSettingsGetResults, error)
1. *GremlinResourcesUpdateGremlinGraphThroughputFuture.Result(GremlinResourcesClient) (ThroughputSettingsGetResults, error)
1. *MongoDBResourcesCreateUpdateMongoDBCollectionFuture.Result(MongoDBResourcesClient) (MongoDBCollectionGetResults, error)
1. *MongoDBResourcesCreateUpdateMongoDBDatabaseFuture.Result(MongoDBResourcesClient) (MongoDBDatabaseGetResults, error)
1. *MongoDBResourcesDeleteMongoDBCollectionFuture.Result(MongoDBResourcesClient) (autorest.Response, error)
1. *MongoDBResourcesDeleteMongoDBDatabaseFuture.Result(MongoDBResourcesClient) (autorest.Response, error)
1. *MongoDBResourcesUpdateMongoDBCollectionThroughputFuture.Result(MongoDBResourcesClient) (ThroughputSettingsGetResults, error)
1. *MongoDBResourcesUpdateMongoDBDatabaseThroughputFuture.Result(MongoDBResourcesClient) (ThroughputSettingsGetResults, error)
1. *NotebookWorkspacesCreateOrUpdateFuture.Result(NotebookWorkspacesClient) (NotebookWorkspace, error)
1. *NotebookWorkspacesDeleteFuture.Result(NotebookWorkspacesClient) (autorest.Response, error)
1. *NotebookWorkspacesRegenerateAuthTokenFuture.Result(NotebookWorkspacesClient) (autorest.Response, error)
1. *NotebookWorkspacesStartFuture.Result(NotebookWorkspacesClient) (autorest.Response, error)
1. *PrivateEndpointConnectionsCreateOrUpdateFuture.Result(PrivateEndpointConnectionsClient) (PrivateEndpointConnection, error)
1. *PrivateEndpointConnectionsDeleteFuture.Result(PrivateEndpointConnectionsClient) (autorest.Response, error)
1. *SQLResourcesCreateUpdateSQLContainerFuture.Result(SQLResourcesClient) (SQLContainerGetResults, error)
1. *SQLResourcesCreateUpdateSQLDatabaseFuture.Result(SQLResourcesClient) (SQLDatabaseGetResults, error)
1. *SQLResourcesCreateUpdateSQLRoleAssignmentFuture.Result(SQLResourcesClient) (SQLRoleAssignmentGetResults, error)
1. *SQLResourcesCreateUpdateSQLRoleDefinitionFuture.Result(SQLResourcesClient) (SQLRoleDefinitionGetResults, error)
1. *SQLResourcesCreateUpdateSQLStoredProcedureFuture.Result(SQLResourcesClient) (SQLStoredProcedureGetResults, error)
1. *SQLResourcesCreateUpdateSQLTriggerFuture.Result(SQLResourcesClient) (SQLTriggerGetResults, error)
1. *SQLResourcesCreateUpdateSQLUserDefinedFunctionFuture.Result(SQLResourcesClient) (SQLUserDefinedFunctionGetResults, error)
1. *SQLResourcesDeleteSQLContainerFuture.Result(SQLResourcesClient) (autorest.Response, error)
1. *SQLResourcesDeleteSQLDatabaseFuture.Result(SQLResourcesClient) (autorest.Response, error)
1. *SQLResourcesDeleteSQLRoleAssignmentFuture.Result(SQLResourcesClient) (autorest.Response, error)
1. *SQLResourcesDeleteSQLRoleDefinitionFuture.Result(SQLResourcesClient) (autorest.Response, error)
1. *SQLResourcesDeleteSQLStoredProcedureFuture.Result(SQLResourcesClient) (autorest.Response, error)
1. *SQLResourcesDeleteSQLTriggerFuture.Result(SQLResourcesClient) (autorest.Response, error)
1. *SQLResourcesDeleteSQLUserDefinedFunctionFuture.Result(SQLResourcesClient) (autorest.Response, error)
1. *SQLResourcesUpdateSQLContainerThroughputFuture.Result(SQLResourcesClient) (ThroughputSettingsGetResults, error)
1. *SQLResourcesUpdateSQLDatabaseThroughputFuture.Result(SQLResourcesClient) (ThroughputSettingsGetResults, error)
1. *TableResourcesCreateUpdateTableFuture.Result(TableResourcesClient) (TableGetResults, error)
1. *TableResourcesDeleteTableFuture.Result(TableResourcesClient) (autorest.Response, error)
1. *TableResourcesUpdateTableThroughputFuture.Result(TableResourcesClient) (ThroughputSettingsGetResults, error)

## Struct Changes

### Removed Struct Fields

1. CassandraResourcesCreateUpdateCassandraKeyspaceFuture.azure.Future
1. CassandraResourcesCreateUpdateCassandraTableFuture.azure.Future
1. CassandraResourcesDeleteCassandraKeyspaceFuture.azure.Future
1. CassandraResourcesDeleteCassandraTableFuture.azure.Future
1. CassandraResourcesUpdateCassandraKeyspaceThroughputFuture.azure.Future
1. CassandraResourcesUpdateCassandraTableThroughputFuture.azure.Future
1. DatabaseAccountsCreateOrUpdateFuture.azure.Future
1. DatabaseAccountsDeleteFuture.azure.Future
1. DatabaseAccountsFailoverPriorityChangeFuture.azure.Future
1. DatabaseAccountsOfflineRegionFuture.azure.Future
1. DatabaseAccountsOnlineRegionFuture.azure.Future
1. DatabaseAccountsRegenerateKeyFuture.azure.Future
1. DatabaseAccountsUpdateFuture.azure.Future
1. GremlinResourcesCreateUpdateGremlinDatabaseFuture.azure.Future
1. GremlinResourcesCreateUpdateGremlinGraphFuture.azure.Future
1. GremlinResourcesDeleteGremlinDatabaseFuture.azure.Future
1. GremlinResourcesDeleteGremlinGraphFuture.azure.Future
1. GremlinResourcesUpdateGremlinDatabaseThroughputFuture.azure.Future
1. GremlinResourcesUpdateGremlinGraphThroughputFuture.azure.Future
1. MongoDBResourcesCreateUpdateMongoDBCollectionFuture.azure.Future
1. MongoDBResourcesCreateUpdateMongoDBDatabaseFuture.azure.Future
1. MongoDBResourcesDeleteMongoDBCollectionFuture.azure.Future
1. MongoDBResourcesDeleteMongoDBDatabaseFuture.azure.Future
1. MongoDBResourcesUpdateMongoDBCollectionThroughputFuture.azure.Future
1. MongoDBResourcesUpdateMongoDBDatabaseThroughputFuture.azure.Future
1. NotebookWorkspacesCreateOrUpdateFuture.azure.Future
1. NotebookWorkspacesDeleteFuture.azure.Future
1. NotebookWorkspacesRegenerateAuthTokenFuture.azure.Future
1. NotebookWorkspacesStartFuture.azure.Future
1. PrivateEndpointConnectionsCreateOrUpdateFuture.azure.Future
1. PrivateEndpointConnectionsDeleteFuture.azure.Future
1. SQLResourcesCreateUpdateSQLContainerFuture.azure.Future
1. SQLResourcesCreateUpdateSQLDatabaseFuture.azure.Future
1. SQLResourcesCreateUpdateSQLRoleAssignmentFuture.azure.Future
1. SQLResourcesCreateUpdateSQLRoleDefinitionFuture.azure.Future
1. SQLResourcesCreateUpdateSQLStoredProcedureFuture.azure.Future
1. SQLResourcesCreateUpdateSQLTriggerFuture.azure.Future
1. SQLResourcesCreateUpdateSQLUserDefinedFunctionFuture.azure.Future
1. SQLResourcesDeleteSQLContainerFuture.azure.Future
1. SQLResourcesDeleteSQLDatabaseFuture.azure.Future
1. SQLResourcesDeleteSQLRoleAssignmentFuture.azure.Future
1. SQLResourcesDeleteSQLRoleDefinitionFuture.azure.Future
1. SQLResourcesDeleteSQLStoredProcedureFuture.azure.Future
1. SQLResourcesDeleteSQLTriggerFuture.azure.Future
1. SQLResourcesDeleteSQLUserDefinedFunctionFuture.azure.Future
1. SQLResourcesUpdateSQLContainerThroughputFuture.azure.Future
1. SQLResourcesUpdateSQLDatabaseThroughputFuture.azure.Future
1. TableResourcesCreateUpdateTableFuture.azure.Future
1. TableResourcesDeleteTableFuture.azure.Future
1. TableResourcesUpdateTableThroughputFuture.azure.Future

## Signature Changes

### Const Types

1. Create changed type from TriggerOperation to OperationType
1. Delete changed type from TriggerOperation to OperationType
1. MongoDB changed type from DatabaseAccountKind to APIType
1. Replace changed type from TriggerOperation to OperationType

### New Constants

1. APIType.Cassandra
1. APIType.Gremlin
1. APIType.GremlinV2
1. APIType.SQL
1. APIType.Table
1. DatabaseAccountKind.DatabaseAccountKindGlobalDocumentDB
1. DatabaseAccountKind.DatabaseAccountKindMongoDB
1. DatabaseAccountKind.DatabaseAccountKindParse
1. OperationType.SystemOperation
1. TriggerOperation.TriggerOperationAll
1. TriggerOperation.TriggerOperationCreate
1. TriggerOperation.TriggerOperationDelete
1. TriggerOperation.TriggerOperationReplace
1. TriggerOperation.TriggerOperationUpdate

### New Funcs

1. *RestorableMongodbCollectionGetResult.UnmarshalJSON([]byte) error
1. *RestorableMongodbDatabaseGetResult.UnmarshalJSON([]byte) error
1. *RestorableSQLContainerGetResult.UnmarshalJSON([]byte) error
1. *RestorableSQLDatabaseGetResult.UnmarshalJSON([]byte) error
1. NewRestorableMongodbCollectionsClient(string) RestorableMongodbCollectionsClient
1. NewRestorableMongodbCollectionsClientWithBaseURI(string, string) RestorableMongodbCollectionsClient
1. NewRestorableMongodbDatabasesClient(string) RestorableMongodbDatabasesClient
1. NewRestorableMongodbDatabasesClientWithBaseURI(string, string) RestorableMongodbDatabasesClient
1. NewRestorableMongodbResourcesClient(string) RestorableMongodbResourcesClient
1. NewRestorableMongodbResourcesClientWithBaseURI(string, string) RestorableMongodbResourcesClient
1. NewRestorableSQLContainersClient(string) RestorableSQLContainersClient
1. NewRestorableSQLContainersClientWithBaseURI(string, string) RestorableSQLContainersClient
1. NewRestorableSQLDatabasesClient(string) RestorableSQLDatabasesClient
1. NewRestorableSQLDatabasesClientWithBaseURI(string, string) RestorableSQLDatabasesClient
1. NewRestorableSQLResourcesClient(string) RestorableSQLResourcesClient
1. NewRestorableSQLResourcesClientWithBaseURI(string, string) RestorableSQLResourcesClient
1. PossibleAPITypeValues() []APIType
1. PossibleOperationTypeValues() []OperationType
1. RestorableDatabaseAccountProperties.MarshalJSON() ([]byte, error)
1. RestorableMongodbCollectionGetResult.MarshalJSON() ([]byte, error)
1. RestorableMongodbCollectionsClient.List(context.Context, string, string, string) (RestorableMongodbCollectionsListResult, error)
1. RestorableMongodbCollectionsClient.ListPreparer(context.Context, string, string, string) (*http.Request, error)
1. RestorableMongodbCollectionsClient.ListResponder(*http.Response) (RestorableMongodbCollectionsListResult, error)
1. RestorableMongodbCollectionsClient.ListSender(*http.Request) (*http.Response, error)
1. RestorableMongodbDatabaseGetResult.MarshalJSON() ([]byte, error)
1. RestorableMongodbDatabasesClient.List(context.Context, string, string) (RestorableMongodbDatabasesListResult, error)
1. RestorableMongodbDatabasesClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. RestorableMongodbDatabasesClient.ListResponder(*http.Response) (RestorableMongodbDatabasesListResult, error)
1. RestorableMongodbDatabasesClient.ListSender(*http.Request) (*http.Response, error)
1. RestorableMongodbResourcesClient.List(context.Context, string, string, string, string) (RestorableMongodbResourcesListResult, error)
1. RestorableMongodbResourcesClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. RestorableMongodbResourcesClient.ListResponder(*http.Response) (RestorableMongodbResourcesListResult, error)
1. RestorableMongodbResourcesClient.ListSender(*http.Request) (*http.Response, error)
1. RestorableSQLContainerGetResult.MarshalJSON() ([]byte, error)
1. RestorableSQLContainerPropertiesResource.MarshalJSON() ([]byte, error)
1. RestorableSQLContainerPropertiesResourceContainer.MarshalJSON() ([]byte, error)
1. RestorableSQLContainersClient.List(context.Context, string, string, string) (RestorableSQLContainersListResult, error)
1. RestorableSQLContainersClient.ListPreparer(context.Context, string, string, string) (*http.Request, error)
1. RestorableSQLContainersClient.ListResponder(*http.Response) (RestorableSQLContainersListResult, error)
1. RestorableSQLContainersClient.ListSender(*http.Request) (*http.Response, error)
1. RestorableSQLDatabaseGetResult.MarshalJSON() ([]byte, error)
1. RestorableSQLDatabasePropertiesResource.MarshalJSON() ([]byte, error)
1. RestorableSQLDatabasePropertiesResourceDatabase.MarshalJSON() ([]byte, error)
1. RestorableSQLDatabasesClient.List(context.Context, string, string) (RestorableSQLDatabasesListResult, error)
1. RestorableSQLDatabasesClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. RestorableSQLDatabasesClient.ListResponder(*http.Response) (RestorableSQLDatabasesListResult, error)
1. RestorableSQLDatabasesClient.ListSender(*http.Request) (*http.Response, error)
1. RestorableSQLResourcesClient.List(context.Context, string, string, string, string) (RestorableSQLResourcesListResult, error)
1. RestorableSQLResourcesClient.ListPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. RestorableSQLResourcesClient.ListResponder(*http.Response) (RestorableSQLResourcesListResult, error)
1. RestorableSQLResourcesClient.ListSender(*http.Request) (*http.Response, error)

## Struct Changes

### New Structs

1. RestorableLocationResource
1. RestorableMongodbCollectionGetResult
1. RestorableMongodbCollectionProperties
1. RestorableMongodbCollectionPropertiesResource
1. RestorableMongodbCollectionsClient
1. RestorableMongodbCollectionsListResult
1. RestorableMongodbDatabaseGetResult
1. RestorableMongodbDatabaseProperties
1. RestorableMongodbDatabasePropertiesResource
1. RestorableMongodbDatabasesClient
1. RestorableMongodbDatabasesListResult
1. RestorableMongodbResourcesClient
1. RestorableMongodbResourcesListResult
1. RestorableSQLContainerGetResult
1. RestorableSQLContainerProperties
1. RestorableSQLContainerPropertiesResource
1. RestorableSQLContainerPropertiesResourceContainer
1. RestorableSQLContainersClient
1. RestorableSQLContainersListResult
1. RestorableSQLDatabaseGetResult
1. RestorableSQLDatabaseProperties
1. RestorableSQLDatabasePropertiesResource
1. RestorableSQLDatabasePropertiesResourceDatabase
1. RestorableSQLDatabasesClient
1. RestorableSQLDatabasesListResult
1. RestorableSQLResourcesClient
1. RestorableSQLResourcesListResult

### New Struct Fields

1. CassandraResourcesCreateUpdateCassandraKeyspaceFuture.Result
1. CassandraResourcesCreateUpdateCassandraKeyspaceFuture.azure.FutureAPI
1. CassandraResourcesCreateUpdateCassandraTableFuture.Result
1. CassandraResourcesCreateUpdateCassandraTableFuture.azure.FutureAPI
1. CassandraResourcesDeleteCassandraKeyspaceFuture.Result
1. CassandraResourcesDeleteCassandraKeyspaceFuture.azure.FutureAPI
1. CassandraResourcesDeleteCassandraTableFuture.Result
1. CassandraResourcesDeleteCassandraTableFuture.azure.FutureAPI
1. CassandraResourcesUpdateCassandraKeyspaceThroughputFuture.Result
1. CassandraResourcesUpdateCassandraKeyspaceThroughputFuture.azure.FutureAPI
1. CassandraResourcesUpdateCassandraTableThroughputFuture.Result
1. CassandraResourcesUpdateCassandraTableThroughputFuture.azure.FutureAPI
1. DatabaseAccountsCreateOrUpdateFuture.Result
1. DatabaseAccountsCreateOrUpdateFuture.azure.FutureAPI
1. DatabaseAccountsDeleteFuture.Result
1. DatabaseAccountsDeleteFuture.azure.FutureAPI
1. DatabaseAccountsFailoverPriorityChangeFuture.Result
1. DatabaseAccountsFailoverPriorityChangeFuture.azure.FutureAPI
1. DatabaseAccountsOfflineRegionFuture.Result
1. DatabaseAccountsOfflineRegionFuture.azure.FutureAPI
1. DatabaseAccountsOnlineRegionFuture.Result
1. DatabaseAccountsOnlineRegionFuture.azure.FutureAPI
1. DatabaseAccountsRegenerateKeyFuture.Result
1. DatabaseAccountsRegenerateKeyFuture.azure.FutureAPI
1. DatabaseAccountsUpdateFuture.Result
1. DatabaseAccountsUpdateFuture.azure.FutureAPI
1. GremlinResourcesCreateUpdateGremlinDatabaseFuture.Result
1. GremlinResourcesCreateUpdateGremlinDatabaseFuture.azure.FutureAPI
1. GremlinResourcesCreateUpdateGremlinGraphFuture.Result
1. GremlinResourcesCreateUpdateGremlinGraphFuture.azure.FutureAPI
1. GremlinResourcesDeleteGremlinDatabaseFuture.Result
1. GremlinResourcesDeleteGremlinDatabaseFuture.azure.FutureAPI
1. GremlinResourcesDeleteGremlinGraphFuture.Result
1. GremlinResourcesDeleteGremlinGraphFuture.azure.FutureAPI
1. GremlinResourcesUpdateGremlinDatabaseThroughputFuture.Result
1. GremlinResourcesUpdateGremlinDatabaseThroughputFuture.azure.FutureAPI
1. GremlinResourcesUpdateGremlinGraphThroughputFuture.Result
1. GremlinResourcesUpdateGremlinGraphThroughputFuture.azure.FutureAPI
1. MongoDBResourcesCreateUpdateMongoDBCollectionFuture.Result
1. MongoDBResourcesCreateUpdateMongoDBCollectionFuture.azure.FutureAPI
1. MongoDBResourcesCreateUpdateMongoDBDatabaseFuture.Result
1. MongoDBResourcesCreateUpdateMongoDBDatabaseFuture.azure.FutureAPI
1. MongoDBResourcesDeleteMongoDBCollectionFuture.Result
1. MongoDBResourcesDeleteMongoDBCollectionFuture.azure.FutureAPI
1. MongoDBResourcesDeleteMongoDBDatabaseFuture.Result
1. MongoDBResourcesDeleteMongoDBDatabaseFuture.azure.FutureAPI
1. MongoDBResourcesUpdateMongoDBCollectionThroughputFuture.Result
1. MongoDBResourcesUpdateMongoDBCollectionThroughputFuture.azure.FutureAPI
1. MongoDBResourcesUpdateMongoDBDatabaseThroughputFuture.Result
1. MongoDBResourcesUpdateMongoDBDatabaseThroughputFuture.azure.FutureAPI
1. NotebookWorkspacesCreateOrUpdateFuture.Result
1. NotebookWorkspacesCreateOrUpdateFuture.azure.FutureAPI
1. NotebookWorkspacesDeleteFuture.Result
1. NotebookWorkspacesDeleteFuture.azure.FutureAPI
1. NotebookWorkspacesRegenerateAuthTokenFuture.Result
1. NotebookWorkspacesRegenerateAuthTokenFuture.azure.FutureAPI
1. NotebookWorkspacesStartFuture.Result
1. NotebookWorkspacesStartFuture.azure.FutureAPI
1. PrivateEndpointConnectionsCreateOrUpdateFuture.Result
1. PrivateEndpointConnectionsCreateOrUpdateFuture.azure.FutureAPI
1. PrivateEndpointConnectionsDeleteFuture.Result
1. PrivateEndpointConnectionsDeleteFuture.azure.FutureAPI
1. RestorableDatabaseAccountProperties.APIType
1. RestorableDatabaseAccountProperties.RestorableLocations
1. SQLResourcesCreateUpdateSQLContainerFuture.Result
1. SQLResourcesCreateUpdateSQLContainerFuture.azure.FutureAPI
1. SQLResourcesCreateUpdateSQLDatabaseFuture.Result
1. SQLResourcesCreateUpdateSQLDatabaseFuture.azure.FutureAPI
1. SQLResourcesCreateUpdateSQLRoleAssignmentFuture.Result
1. SQLResourcesCreateUpdateSQLRoleAssignmentFuture.azure.FutureAPI
1. SQLResourcesCreateUpdateSQLRoleDefinitionFuture.Result
1. SQLResourcesCreateUpdateSQLRoleDefinitionFuture.azure.FutureAPI
1. SQLResourcesCreateUpdateSQLStoredProcedureFuture.Result
1. SQLResourcesCreateUpdateSQLStoredProcedureFuture.azure.FutureAPI
1. SQLResourcesCreateUpdateSQLTriggerFuture.Result
1. SQLResourcesCreateUpdateSQLTriggerFuture.azure.FutureAPI
1. SQLResourcesCreateUpdateSQLUserDefinedFunctionFuture.Result
1. SQLResourcesCreateUpdateSQLUserDefinedFunctionFuture.azure.FutureAPI
1. SQLResourcesDeleteSQLContainerFuture.Result
1. SQLResourcesDeleteSQLContainerFuture.azure.FutureAPI
1. SQLResourcesDeleteSQLDatabaseFuture.Result
1. SQLResourcesDeleteSQLDatabaseFuture.azure.FutureAPI
1. SQLResourcesDeleteSQLRoleAssignmentFuture.Result
1. SQLResourcesDeleteSQLRoleAssignmentFuture.azure.FutureAPI
1. SQLResourcesDeleteSQLRoleDefinitionFuture.Result
1. SQLResourcesDeleteSQLRoleDefinitionFuture.azure.FutureAPI
1. SQLResourcesDeleteSQLStoredProcedureFuture.Result
1. SQLResourcesDeleteSQLStoredProcedureFuture.azure.FutureAPI
1. SQLResourcesDeleteSQLTriggerFuture.Result
1. SQLResourcesDeleteSQLTriggerFuture.azure.FutureAPI
1. SQLResourcesDeleteSQLUserDefinedFunctionFuture.Result
1. SQLResourcesDeleteSQLUserDefinedFunctionFuture.azure.FutureAPI
1. SQLResourcesUpdateSQLContainerThroughputFuture.Result
1. SQLResourcesUpdateSQLContainerThroughputFuture.azure.FutureAPI
1. SQLResourcesUpdateSQLDatabaseThroughputFuture.Result
1. SQLResourcesUpdateSQLDatabaseThroughputFuture.azure.FutureAPI
1. TableResourcesCreateUpdateTableFuture.Result
1. TableResourcesCreateUpdateTableFuture.azure.FutureAPI
1. TableResourcesDeleteTableFuture.Result
1. TableResourcesDeleteTableFuture.azure.FutureAPI
1. TableResourcesUpdateTableThroughputFuture.Result
1. TableResourcesUpdateTableThroughputFuture.azure.FutureAPI
