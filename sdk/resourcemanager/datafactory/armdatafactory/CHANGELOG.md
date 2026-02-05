# Release History

## 11.0.0 (2026-02-05)
### Breaking Changes

- Type of `JiraObjectDataset.TypeProperties` has been changed from `*GenericDatasetTypeProperties` to `*JiraTableDatasetTypeProperties`

### Features Added

- Type of `ExpressionV2.Value` has been changed from `*string` to `any`
- New enum type `AmazonRdsForOracleAuthenticationType` with values `AmazonRdsForOracleAuthenticationTypeBasic`
- New enum type `HDInsightClusterAuthenticationType` with values `HDInsightClusterAuthenticationTypeBasicAuth`, `HDInsightClusterAuthenticationTypeSystemAssignedManagedIdentity`, `HDInsightClusterAuthenticationTypeUserAssignedManagedIdentity`
- New enum type `HDInsightOndemandClusterResourceGroupAuthenticationType` with values `HDInsightOndemandClusterResourceGroupAuthenticationTypeServicePrincipalKey`, `HDInsightOndemandClusterResourceGroupAuthenticationTypeSystemAssignedManagedIdentity`, `HDInsightOndemandClusterResourceGroupAuthenticationTypeUserAssignedManagedIdentity`
- New enum type `ImpalaThriftTransportProtocol` with values `ImpalaThriftTransportProtocolBinary`, `ImpalaThriftTransportProtocolHTTP`
- New enum type `InteractiveCapabilityStatus` with values `InteractiveCapabilityStatusDisabled`, `InteractiveCapabilityStatusDisabling`, `InteractiveCapabilityStatusEnabled`, `InteractiveCapabilityStatusEnabling`
- New enum type `LakehouseAuthenticationType` with values `LakehouseAuthenticationTypeServicePrincipal`, `LakehouseAuthenticationTypeSystemAssignedManagedIdentity`, `LakehouseAuthenticationTypeUserAssignedManagedIdentity`
- New enum type `NetezzaSecurityLevelType` with values `NetezzaSecurityLevelTypeOnlyUnSecured`, `NetezzaSecurityLevelTypePreferredUnSecured`
- New enum type `WarehouseAuthenticationType` with values `WarehouseAuthenticationTypeServicePrincipal`, `WarehouseAuthenticationTypeSystemAssignedManagedIdentity`, `WarehouseAuthenticationTypeUserAssignedManagedIdentity`
- New function `*ClientFactory.NewIntegrationRuntimeClient() *IntegrationRuntimeClient`
- New function `*DatabricksJobActivity.GetActivity() *Activity`
- New function `*DatabricksJobActivity.GetExecutionActivity() *ExecutionActivity`
- New function `NewIntegrationRuntimeClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*IntegrationRuntimeClient, error)`
- New function `*IntegrationRuntimeClient.BeginDisableInteractiveQuery(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, options *IntegrationRuntimeClientBeginDisableInteractiveQueryOptions) (*runtime.Poller[IntegrationRuntimeClientDisableInteractiveQueryResponse], error)`
- New function `*IntegrationRuntimeClient.BeginEnableInteractiveQuery(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, enableInteractiveQueryRequest EnableInteractiveQueryRequest, options *IntegrationRuntimeClientBeginEnableInteractiveQueryOptions) (*runtime.Poller[IntegrationRuntimeClientEnableInteractiveQueryResponse], error)`
- New struct `DatabricksJobActivity`
- New struct `DatabricksJobActivityTypeProperties`
- New struct `EnableInteractiveQueryRequest`
- New struct `InteractiveQueryProperties`
- New struct `JiraTableDatasetTypeProperties`
- New field `AuthenticationType`, `CryptoChecksumClient`, `CryptoChecksumTypesClient`, `EnableBulkLoad`, `EncryptionClient`, `EncryptionTypesClient`, `FetchSize`, `FetchTswtzAsTimestamp`, `InitialLobFetchSize`, `InitializationString`, `Server`, `StatementCacheSize`, `SupportV1DataTypes`, `Username` in struct `AmazonRdsForLinkedServiceTypeProperties`
- New field `NumberPrecision`, `NumberScale` in struct `AmazonRdsForOracleSource`
- New field `DataSecurityMode` in struct `AzureDatabricksLinkedServiceTypeProperties`
- New field `ClusterAuthType`, `Credential` in struct `HDInsightLinkedServiceTypeProperties`
- New field `ClusterResourceGroupAuthType` in struct `HDInsightOnDemandLinkedServiceTypeProperties`
- New field `EnableServerCertificateValidation` in struct `HiveLinkedServiceTypeProperties`
- New field `EnableServerCertificateValidation`, `ThriftTransportProtocol` in struct `ImpalaLinkedServiceTypeProperties`
- New field `AuthenticationType`, `Credential` in struct `LakeHouseLinkedServiceTypeProperties`
- New field `TreatDecimalAsString` in struct `LookupActivityTypeProperties`
- New field `InteractiveQuery` in struct `ManagedIntegrationRuntimeTypeProperties`
- New field `Database`, `Port`, `SecurityLevel`, `Server`, `UID` in struct `NetezzaLinkedServiceTypeProperties`
- New field `NumberPrecision`, `NumberScale` in struct `OracleSource`
- New field `RefreshToken` in struct `QuickBooksLinkedServiceTypeProperties`
- New field `PartitionOption` in struct `SalesforceV2Source`
- New field `TreatDecimalAsString` in struct `ScriptActivityTypeProperties`
- New field `UseUTCTimestamps` in struct `SnowflakeLinkedV2ServiceTypeProperties`
- New field `EnableServerCertificateValidation` in struct `SparkLinkedServiceTypeProperties`
- New field `AuthenticationType`, `Credential` in struct `WarehouseLinkedServiceTypeProperties`


## 10.0.0 (2025-04-24)
### Breaking Changes

- Type of `ServiceNowV2ObjectDataset.TypeProperties` has been changed from `*GenericDatasetTypeProperties` to `*ServiceNowV2DatasetTypeProperties`

### Features Added

- New enum type `AzurePostgreSQLWriteMethodEnum` with values `AzurePostgreSQLWriteMethodEnumBulkInsert`, `AzurePostgreSQLWriteMethodEnumCopyCommand`, `AzurePostgreSQLWriteMethodEnumUpsert`
- New enum type `GreenplumAuthenticationType` with values `GreenplumAuthenticationTypeBasic`
- New enum type `OracleAuthenticationType` with values `OracleAuthenticationTypeBasic`
- New enum type `ValueType` with values `ValueTypeActual`, `ValueTypeDisplay`
- New function `*TeradataImportCommand.GetImportSettings() *ImportSettings`
- New function `*TeradataSink.GetCopySink() *CopySink`
- New struct `AzurePostgreSQLSinkUpsertSettings`
- New struct `ServiceNowV2DatasetTypeProperties`
- New struct `TeradataImportCommand`
- New struct `TeradataSink`
- New field `AzureCloudType`, `Credential`, `ServicePrincipalCredentialType`, `ServicePrincipalEmbeddedCert`, `ServicePrincipalEmbeddedCertPassword`, `ServicePrincipalID`, `ServicePrincipalKey`, `Tenant` in struct `AzurePostgreSQLLinkedServiceTypeProperties`
- New field `UpsertSettings`, `WriteMethod` in struct `AzurePostgreSQLSink`
- New field `BypassBusinessLogicExecution`, `BypassPowerAutomateFlows` in struct `CommonDataServiceForAppsSink`
- New field `BypassBusinessLogicExecution`, `BypassPowerAutomateFlows` in struct `DynamicsCrmSink`
- New field `BypassBusinessLogicExecution`, `BypassPowerAutomateFlows` in struct `DynamicsSink`
- New field `AuthenticationType`, `CommandTimeout`, `ConnectionTimeout`, `Database`, `Host`, `Port`, `SSLMode`, `Username` in struct `GreenplumLinkedServiceTypeProperties`
- New field `ServicePrincipalCredentialType`, `ServicePrincipalEmbeddedCert`, `ServicePrincipalEmbeddedCertPassword` in struct `Office365LinkedServiceTypeProperties`
- New field `AuthenticationType`, `CryptoChecksumClient`, `CryptoChecksumTypesClient`, `EnableBulkLoad`, `EncryptionClient`, `EncryptionTypesClient`, `FetchSize`, `FetchTswtzAsTimestamp`, `InitialLobFetchSize`, `InitializationString`, `Server`, `StatementCacheSize`, `SupportV1DataTypes`, `Username` in struct `OracleLinkedServiceTypeProperties`
- New field `EnableServerCertificateValidation` in struct `PrestoLinkedServiceTypeProperties`
- New field `ReturnMultistatementResult` in struct `ScriptActivityTypeProperties`
- New field `Role`, `Schema` in struct `SnowflakeLinkedV2ServiceTypeProperties`
- New field `CharacterSet`, `HTTPSPortNumber`, `MaxRespSize`, `PortNumber`, `SSLMode`, `UseDataEncryption` in struct `TeradataLinkedServiceTypeProperties`


## 9.1.0 (2024-12-26)
### Features Added

- New function `*IcebergDataset.GetDataset() *Dataset`
- New function `*IcebergSink.GetCopySink() *CopySink`
- New function `*IcebergWriteSettings.GetFormatWriteSettings() *FormatWriteSettings`
- New struct `IcebergDataset`
- New struct `IcebergDatasetTypeProperties`
- New struct `IcebergSink`
- New struct `IcebergWriteSettings`
- New field `CommandTimeout`, `Database`, `Encoding`, `Port`, `ReadBufferSize`, `SSLMode`, `Server`, `Timeout`, `Timezone`, `TrustServerCertificate`, `Username` in struct `AzurePostgreSQLLinkedServiceTypeProperties`
- New field `SSLMode`, `UseSystemTrustStore` in struct `MariaDBLinkedServiceTypeProperties`
- New field `AllowZeroDateTime`, `ConnectionTimeout`, `ConvertZeroDateTime`, `GUIDFormat`, `SSLCert`, `SSLKey`, `TreatTinyAsBoolean` in struct `MySQLLinkedServiceTypeProperties`
- New field `AuthenticationType` in struct `PostgreSQLV2LinkedServiceTypeProperties`
- New field `PageSize` in struct `SalesforceV2Source`
- New field `PageSize` in struct `ServiceNowV2Source`
- New field `Host` in struct `SnowflakeLinkedV2ServiceTypeProperties`


## 9.0.0 (2024-08-23)
### Breaking Changes

- Type of `AzureTableStorageLinkedService.TypeProperties` has been changed from `*AzureStorageLinkedServiceTypeProperties` to `*AzureTableStorageLinkedServiceTypeProperties`

### Features Added

- New value `SQLServerAuthenticationTypeUserAssignedManagedIdentity` added to enum type `SQLServerAuthenticationType`
- New struct `AzureTableStorageLinkedServiceTypeProperties`
- New struct `ContinuationSettingsReference`
- New field `Version` in struct `AmazonMWSLinkedService`
- New field `Version` in struct `AmazonRdsForOracleLinkedService`
- New field `Version` in struct `AmazonRdsForSQLServerLinkedService`
- New field `Version` in struct `AmazonRedshiftLinkedService`
- New field `Version` in struct `AmazonS3CompatibleLinkedService`
- New field `Version` in struct `AmazonS3LinkedService`
- New field `Version` in struct `AppFiguresLinkedService`
- New field `Version` in struct `AsanaLinkedService`
- New field `Version` in struct `AzureBatchLinkedService`
- New field `Version` in struct `AzureBlobFSLinkedService`
- New field `Version` in struct `AzureBlobStorageLinkedService`
- New field `Version` in struct `AzureDataExplorerLinkedService`
- New field `Version` in struct `AzureDataLakeAnalyticsLinkedService`
- New field `Version` in struct `AzureDataLakeStoreLinkedService`
- New field `Version` in struct `AzureDatabricksDeltaLakeLinkedService`
- New field `Version` in struct `AzureDatabricksLinkedService`
- New field `Version` in struct `AzureFileStorageLinkedService`
- New field `Credential`, `ServiceEndpoint` in struct `AzureFileStorageLinkedServiceTypeProperties`
- New field `Version` in struct `AzureFunctionLinkedService`
- New field `Version` in struct `AzureKeyVaultLinkedService`
- New field `Version` in struct `AzureMLLinkedService`
- New field `Version` in struct `AzureMLServiceLinkedService`
- New field `Version` in struct `AzureMariaDBLinkedService`
- New field `Version` in struct `AzureMySQLLinkedService`
- New field `Version` in struct `AzurePostgreSQLLinkedService`
- New field `Version` in struct `AzureSQLDWLinkedService`
- New field `Version` in struct `AzureSQLDatabaseLinkedService`
- New field `Version` in struct `AzureSQLMILinkedService`
- New field `Version` in struct `AzureSearchLinkedService`
- New field `Version` in struct `AzureStorageLinkedService`
- New field `Version` in struct `AzureSynapseArtifactsLinkedService`
- New field `Version` in struct `AzureTableStorageLinkedService`
- New field `Version` in struct `CassandraLinkedService`
- New field `Version` in struct `CommonDataServiceForAppsLinkedService`
- New field `Domain` in struct `CommonDataServiceForAppsLinkedServiceTypeProperties`
- New field `Version` in struct `ConcurLinkedService`
- New field `Version` in struct `CosmosDbLinkedService`
- New field `Version` in struct `CosmosDbMongoDbAPILinkedService`
- New field `Version` in struct `CouchbaseLinkedService`
- New field `Version` in struct `CustomDataSourceLinkedService`
- New field `Version` in struct `DataworldLinkedService`
- New field `Version` in struct `Db2LinkedService`
- New field `Version` in struct `DrillLinkedService`
- New field `Version` in struct `DynamicsAXLinkedService`
- New field `Version` in struct `DynamicsCrmLinkedService`
- New field `Domain` in struct `DynamicsCrmLinkedServiceTypeProperties`
- New field `Version` in struct `DynamicsLinkedService`
- New field `Domain` in struct `DynamicsLinkedServiceTypeProperties`
- New field `Version` in struct `EloquaLinkedService`
- New field `ContinuationSettings` in struct `ExecuteDataFlowActivityTypeProperties`
- New field `ContinuationSettings` in struct `ExecutePowerQueryActivityTypeProperties`
- New field `Version` in struct `FileServerLinkedService`
- New field `Version` in struct `FtpServerLinkedService`
- New field `Version` in struct `GoogleAdWordsLinkedService`
- New field `Version` in struct `GoogleBigQueryLinkedService`
- New field `Version` in struct `GoogleBigQueryV2LinkedService`
- New field `Version` in struct `GoogleCloudStorageLinkedService`
- New field `Version` in struct `GoogleSheetsLinkedService`
- New field `Version` in struct `GreenplumLinkedService`
- New field `Version` in struct `HBaseLinkedService`
- New field `Version` in struct `HDInsightLinkedService`
- New field `Version` in struct `HDInsightOnDemandLinkedService`
- New field `Version` in struct `HTTPLinkedService`
- New field `Version` in struct `HdfsLinkedService`
- New field `Version` in struct `HiveLinkedService`
- New field `Version` in struct `HubspotLinkedService`
- New field `Version` in struct `ImpalaLinkedService`
- New field `Version` in struct `InformixLinkedService`
- New field `Version` in struct `JiraLinkedService`
- New field `Version` in struct `LakeHouseLinkedService`
- New field `Version` in struct `LinkedService`
- New field `Version` in struct `MagentoLinkedService`
- New field `Version` in struct `MariaDBLinkedService`
- New field `Version` in struct `MarketoLinkedService`
- New field `Version` in struct `MicrosoftAccessLinkedService`
- New field `Version` in struct `MongoDbAtlasLinkedService`
- New field `Version` in struct `MongoDbLinkedService`
- New field `Version` in struct `MongoDbV2LinkedService`
- New field `Version` in struct `MySQLLinkedService`
- New field `Version` in struct `NetezzaLinkedService`
- New field `Version` in struct `ODataLinkedService`
- New field `Version` in struct `OdbcLinkedService`
- New field `Version` in struct `Office365LinkedService`
- New field `Version` in struct `OracleCloudStorageLinkedService`
- New field `Version` in struct `OracleLinkedService`
- New field `Version` in struct `OracleServiceCloudLinkedService`
- New field `Version` in struct `PaypalLinkedService`
- New field `Version` in struct `PhoenixLinkedService`
- New field `Version` in struct `PostgreSQLLinkedService`
- New field `Version` in struct `PostgreSQLV2LinkedService`
- New field `Version` in struct `PrestoLinkedService`
- New field `Version` in struct `QuickBooksLinkedService`
- New field `Version` in struct `QuickbaseLinkedService`
- New field `Version` in struct `ResponsysLinkedService`
- New field `Version` in struct `RestServiceLinkedService`
- New field `ServicePrincipalCredentialType`, `ServicePrincipalEmbeddedCert`, `ServicePrincipalEmbeddedCertPassword` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Version` in struct `SQLServerLinkedService`
- New field `Credential` in struct `SQLServerLinkedServiceTypeProperties`
- New field `Version` in struct `SalesforceLinkedService`
- New field `Version` in struct `SalesforceMarketingCloudLinkedService`
- New field `Version` in struct `SalesforceServiceCloudLinkedService`
- New field `Version` in struct `SalesforceServiceCloudV2LinkedService`
- New field `Version` in struct `SalesforceV2LinkedService`
- New field `Version` in struct `SapBWLinkedService`
- New field `Version` in struct `SapCloudForCustomerLinkedService`
- New field `Version` in struct `SapEccLinkedService`
- New field `Version` in struct `SapHanaLinkedService`
- New field `Version` in struct `SapOdpLinkedService`
- New field `Version` in struct `SapOpenHubLinkedService`
- New field `Version` in struct `SapTableLinkedService`
- New field `Version` in struct `ServiceNowLinkedService`
- New field `Version` in struct `ServiceNowV2LinkedService`
- New field `Version` in struct `SftpServerLinkedService`
- New field `Version` in struct `SharePointOnlineListLinkedService`
- New field `ServicePrincipalCredentialType`, `ServicePrincipalEmbeddedCert`, `ServicePrincipalEmbeddedCertPassword` in struct `SharePointOnlineListLinkedServiceTypeProperties`
- New field `Version` in struct `ShopifyLinkedService`
- New field `Version` in struct `SmartsheetLinkedService`
- New field `StorageIntegration` in struct `SnowflakeExportCopyCommand`
- New field `StorageIntegration` in struct `SnowflakeImportCopyCommand`
- New field `Version` in struct `SnowflakeLinkedService`
- New field `Version` in struct `SnowflakeV2LinkedService`
- New field `Version` in struct `SparkLinkedService`
- New field `Version` in struct `SquareLinkedService`
- New field `Version` in struct `SybaseLinkedService`
- New field `Version` in struct `TeamDeskLinkedService`
- New field `Version` in struct `TeradataLinkedService`
- New field `Version` in struct `TwilioLinkedService`
- New field `Version` in struct `VerticaLinkedService`
- New field `Database`, `Port`, `Server`, `UID` in struct `VerticaLinkedServiceTypeProperties`
- New field `Version` in struct `WarehouseLinkedService`
- New field `Version` in struct `WebLinkedService`
- New field `Version` in struct `XeroLinkedService`
- New field `Version` in struct `ZendeskLinkedService`
- New field `Version` in struct `ZohoLinkedService`


## 8.0.0 (2024-06-05)
### Breaking Changes

- Enum `ScriptType` has been removed
- Field `Operator` of struct `ExpressionV2` has been removed

### Features Added

- Type of `ScriptActivityScriptBlock.Type` has been changed from `*ScriptType` to `any`
- New value `ExpressionV2TypeNAry` added to enum type `ExpressionV2Type`
- New enum type `AmazonRdsForSQLAuthenticationType` with values `AmazonRdsForSQLAuthenticationTypeSQL`, `AmazonRdsForSQLAuthenticationTypeWindows`
- New enum type `AzureSQLDWAuthenticationType` with values `AzureSQLDWAuthenticationTypeSQL`, `AzureSQLDWAuthenticationTypeServicePrincipal`, `AzureSQLDWAuthenticationTypeSystemAssignedManagedIdentity`, `AzureSQLDWAuthenticationTypeUserAssignedManagedIdentity`
- New enum type `AzureSQLDatabaseAuthenticationType` with values `AzureSQLDatabaseAuthenticationTypeSQL`, `AzureSQLDatabaseAuthenticationTypeServicePrincipal`, `AzureSQLDatabaseAuthenticationTypeSystemAssignedManagedIdentity`, `AzureSQLDatabaseAuthenticationTypeUserAssignedManagedIdentity`
- New enum type `AzureSQLMIAuthenticationType` with values `AzureSQLMIAuthenticationTypeSQL`, `AzureSQLMIAuthenticationTypeServicePrincipal`, `AzureSQLMIAuthenticationTypeSystemAssignedManagedIdentity`, `AzureSQLMIAuthenticationTypeUserAssignedManagedIdentity`
- New enum type `SQLServerAuthenticationType` with values `SQLServerAuthenticationTypeSQL`, `SQLServerAuthenticationTypeWindows`
- New struct `ManagedIdentityTypeProperties`
- New field `ApplicationIntent`, `AuthenticationType`, `CommandTimeout`, `ConnectRetryCount`, `ConnectRetryInterval`, `ConnectTimeout`, `Database`, `Encrypt`, `FailoverPartner`, `HostNameInCertificate`, `IntegratedSecurity`, `LoadBalanceTimeout`, `MaxPoolSize`, `MinPoolSize`, `MultiSubnetFailover`, `MultipleActiveResultSets`, `PacketSize`, `Pooling`, `Server`, `TrustServerCertificate` in struct `AmazonRdsForSQLServerLinkedServiceTypeProperties`
- New field `ApplicationIntent`, `AuthenticationType`, `CommandTimeout`, `ConnectRetryCount`, `ConnectRetryInterval`, `ConnectTimeout`, `Database`, `Encrypt`, `FailoverPartner`, `HostNameInCertificate`, `IntegratedSecurity`, `LoadBalanceTimeout`, `MaxPoolSize`, `MinPoolSize`, `MultiSubnetFailover`, `MultipleActiveResultSets`, `PacketSize`, `Pooling`, `Server`, `ServicePrincipalCredential`, `ServicePrincipalCredentialType`, `TrustServerCertificate`, `UserName` in struct `AzureSQLDWLinkedServiceTypeProperties`
- New field `ApplicationIntent`, `AuthenticationType`, `CommandTimeout`, `ConnectRetryCount`, `ConnectRetryInterval`, `ConnectTimeout`, `Database`, `Encrypt`, `FailoverPartner`, `HostNameInCertificate`, `IntegratedSecurity`, `LoadBalanceTimeout`, `MaxPoolSize`, `MinPoolSize`, `MultiSubnetFailover`, `MultipleActiveResultSets`, `PacketSize`, `Pooling`, `Server`, `ServicePrincipalCredential`, `ServicePrincipalCredentialType`, `TrustServerCertificate`, `UserName` in struct `AzureSQLDatabaseLinkedServiceTypeProperties`
- New field `ApplicationIntent`, `AuthenticationType`, `CommandTimeout`, `ConnectRetryCount`, `ConnectRetryInterval`, `ConnectTimeout`, `Database`, `Encrypt`, `FailoverPartner`, `HostNameInCertificate`, `IntegratedSecurity`, `LoadBalanceTimeout`, `MaxPoolSize`, `MinPoolSize`, `MultiSubnetFailover`, `MultipleActiveResultSets`, `PacketSize`, `Pooling`, `Server`, `ServicePrincipalCredential`, `ServicePrincipalCredentialType`, `TrustServerCertificate`, `UserName` in struct `AzureSQLMILinkedServiceTypeProperties`
- New field `Credential` in struct `DynamicsCrmLinkedServiceTypeProperties`
- New field `Operators` in struct `ExpressionV2`
- New field `Schema` in struct `LakeHouseTableDatasetTypeProperties`
- New field `TypeProperties` in struct `ManagedIdentityCredential`
- New field `ApplicationIntent`, `AuthenticationType`, `CommandTimeout`, `ConnectRetryCount`, `ConnectRetryInterval`, `ConnectTimeout`, `Database`, `Encrypt`, `FailoverPartner`, `HostNameInCertificate`, `IntegratedSecurity`, `LoadBalanceTimeout`, `MaxPoolSize`, `MinPoolSize`, `MultiSubnetFailover`, `MultipleActiveResultSets`, `PacketSize`, `Pooling`, `Server`, `TrustServerCertificate` in struct `SQLServerLinkedServiceTypeProperties`
- New field `Query` in struct `SalesforceServiceCloudV2Source`
- New field `Query` in struct `SalesforceV2Source`


## 7.0.0 (2024-04-04)
### Breaking Changes

- Function `*CredentialOperationsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, ManagedIdentityCredentialResource, *CredentialOperationsClientCreateOrUpdateOptions)` to `(context.Context, string, string, string, CredentialResource, *CredentialOperationsClientCreateOrUpdateOptions)`
- Type of `AzureFunctionActivityTypeProperties.Headers` has been changed from `map[string]*string` to `map[string]any`
- Type of `CredentialListResponse.Value` has been changed from `[]*ManagedIdentityCredentialResource` to `[]*CredentialResource`
- Type of `WebActivityTypeProperties.Headers` has been changed from `map[string]*string` to `map[string]any`
- Type of `WebHookActivityTypeProperties.Headers` has been changed from `map[string]*string` to `map[string]any`
- Struct `ManagedIdentityCredentialResource` has been removed
- Struct `ManagedIdentityTypeProperties` has been removed
- Field `ManagedIdentityCredentialResource` of struct `CredentialOperationsClientCreateOrUpdateResponse` has been removed
- Field `ManagedIdentityCredentialResource` of struct `CredentialOperationsClientGetResponse` has been removed
- Field `TypeProperties` of struct `ManagedIdentityCredential` has been removed

### Features Added

- New struct `CredentialResource`
- New anonymous field `CredentialResource` in struct `CredentialOperationsClientCreateOrUpdateResponse`
- New anonymous field `CredentialResource` in struct `CredentialOperationsClientGetResponse`


## 6.1.0 (2024-03-22)
### Features Added

- New enum type `ExpressionV2Type` with values `ExpressionV2TypeBinary`, `ExpressionV2TypeConstant`, `ExpressionV2TypeField`, `ExpressionV2TypeUnary`
- New enum type `GoogleBigQueryV2AuthenticationType` with values `GoogleBigQueryV2AuthenticationTypeServiceAuthentication`, `GoogleBigQueryV2AuthenticationTypeUserAuthentication`
- New enum type `ServiceNowV2AuthenticationType` with values `ServiceNowV2AuthenticationTypeBasic`, `ServiceNowV2AuthenticationTypeOAuth2`
- New function `*GoogleBigQueryV2LinkedService.GetLinkedService() *LinkedService`
- New function `*GoogleBigQueryV2ObjectDataset.GetDataset() *Dataset`
- New function `*GoogleBigQueryV2Source.GetCopySource() *CopySource`
- New function `*GoogleBigQueryV2Source.GetTabularSource() *TabularSource`
- New function `*PostgreSQLV2LinkedService.GetLinkedService() *LinkedService`
- New function `*PostgreSQLV2Source.GetCopySource() *CopySource`
- New function `*PostgreSQLV2Source.GetTabularSource() *TabularSource`
- New function `*PostgreSQLV2TableDataset.GetDataset() *Dataset`
- New function `*ServiceNowV2LinkedService.GetLinkedService() *LinkedService`
- New function `*ServiceNowV2ObjectDataset.GetDataset() *Dataset`
- New function `*ServiceNowV2Source.GetCopySource() *CopySource`
- New function `*ServiceNowV2Source.GetTabularSource() *TabularSource`
- New struct `ExpressionV2`
- New struct `GoogleBigQueryV2DatasetTypeProperties`
- New struct `GoogleBigQueryV2LinkedService`
- New struct `GoogleBigQueryV2LinkedServiceTypeProperties`
- New struct `GoogleBigQueryV2ObjectDataset`
- New struct `GoogleBigQueryV2Source`
- New struct `PostgreSQLV2LinkedService`
- New struct `PostgreSQLV2LinkedServiceTypeProperties`
- New struct `PostgreSQLV2Source`
- New struct `PostgreSQLV2TableDataset`
- New struct `PostgreSQLV2TableDatasetTypeProperties`
- New struct `ServiceNowV2LinkedService`
- New struct `ServiceNowV2LinkedServiceTypeProperties`
- New struct `ServiceNowV2ObjectDataset`
- New struct `ServiceNowV2Source`


## 6.0.0 (2024-02-23)
### Breaking Changes

- Type of `AzureFunctionActivityTypeProperties.Headers` has been changed from `any` to `map[string]*string`
- Type of `WebActivityTypeProperties.Headers` has been changed from `any` to `map[string]*string`
- Type of `WebHookActivityTypeProperties.Headers` has been changed from `any` to `map[string]*string`
- Field `ReadBehavior` of struct `SalesforceServiceCloudV2Source` has been removed
- Field `ReadBehavior` of struct `SalesforceV2Source` has been removed

### Features Added

- New enum type `SnowflakeAuthenticationType` with values `SnowflakeAuthenticationTypeAADServicePrincipal`, `SnowflakeAuthenticationTypeBasic`, `SnowflakeAuthenticationTypeKeyPair`
- New function `*SnowflakeV2Dataset.GetDataset() *Dataset`
- New function `*SnowflakeV2LinkedService.GetLinkedService() *LinkedService`
- New function `*SnowflakeV2Sink.GetCopySink() *CopySink`
- New function `*SnowflakeV2Source.GetCopySource() *CopySource`
- New struct `SnowflakeLinkedV2ServiceTypeProperties`
- New struct `SnowflakeV2Dataset`
- New struct `SnowflakeV2LinkedService`
- New struct `SnowflakeV2Sink`
- New struct `SnowflakeV2Source`
- New field `AuthenticationType` in struct `SalesforceServiceCloudV2LinkedServiceTypeProperties`
- New field `IncludeDeletedObjects` in struct `SalesforceServiceCloudV2Source`
- New field `AuthenticationType` in struct `SalesforceV2LinkedServiceTypeProperties`
- New field `IncludeDeletedObjects` in struct `SalesforceV2Source`


## 5.0.0 (2024-01-26)
### Breaking Changes

- Field `Pwd` of struct `MariaDBLinkedServiceTypeProperties` has been removed

### Features Added

- New enum type `SalesforceV2SinkWriteBehavior` with values `SalesforceV2SinkWriteBehaviorInsert`, `SalesforceV2SinkWriteBehaviorUpsert`
- New function `*SalesforceServiceCloudV2LinkedService.GetLinkedService() *LinkedService`
- New function `*SalesforceServiceCloudV2ObjectDataset.GetDataset() *Dataset`
- New function `*SalesforceServiceCloudV2Sink.GetCopySink() *CopySink`
- New function `*SalesforceServiceCloudV2Source.GetCopySource() *CopySource`
- New function `*SalesforceV2LinkedService.GetLinkedService() *LinkedService`
- New function `*SalesforceV2ObjectDataset.GetDataset() *Dataset`
- New function `*SalesforceV2Sink.GetCopySink() *CopySink`
- New function `*SalesforceV2Source.GetCopySource() *CopySource`
- New function `*SalesforceV2Source.GetTabularSource() *TabularSource`
- New function `*WarehouseLinkedService.GetLinkedService() *LinkedService`
- New function `*WarehouseSink.GetCopySink() *CopySink`
- New function `*WarehouseSource.GetCopySource() *CopySource`
- New function `*WarehouseSource.GetTabularSource() *TabularSource`
- New function `*WarehouseTableDataset.GetDataset() *Dataset`
- New struct `SalesforceServiceCloudV2LinkedService`
- New struct `SalesforceServiceCloudV2LinkedServiceTypeProperties`
- New struct `SalesforceServiceCloudV2ObjectDataset`
- New struct `SalesforceServiceCloudV2ObjectDatasetTypeProperties`
- New struct `SalesforceServiceCloudV2Sink`
- New struct `SalesforceServiceCloudV2Source`
- New struct `SalesforceV2LinkedService`
- New struct `SalesforceV2LinkedServiceTypeProperties`
- New struct `SalesforceV2ObjectDataset`
- New struct `SalesforceV2ObjectDatasetTypeProperties`
- New struct `SalesforceV2Sink`
- New struct `SalesforceV2Source`
- New struct `WarehouseLinkedService`
- New struct `WarehouseLinkedServiceTypeProperties`
- New struct `WarehouseSink`
- New struct `WarehouseSource`
- New struct `WarehouseTableDataset`
- New struct `WarehouseTableDatasetTypeProperties`
- New field `Metadata` in struct `AzureBlobFSWriteSettings`
- New field `Metadata` in struct `AzureBlobStorageWriteSettings`
- New field `Metadata` in struct `AzureDataLakeStoreWriteSettings`
- New field `Metadata` in struct `AzureFileStorageWriteSettings`
- New field `Metadata` in struct `FileServerWriteSettings`
- New field `Metadata` in struct `LakeHouseWriteSettings`
- New field `Database`, `DriverVersion`, `Password`, `Port`, `Server`, `Username` in struct `MariaDBLinkedServiceTypeProperties`
- New field `Database`, `DriverVersion`, `Port`, `SSLMode`, `Server`, `UseSystemTrustStore`, `Username` in struct `MySQLLinkedServiceTypeProperties`
- New field `Metadata` in struct `SftpWriteSettings`
- New field `Metadata` in struct `StoreWriteSettings`
- New field `HTTPRequestTimeout`, `TurnOffAsync` in struct `WebActivityTypeProperties`


## 4.0.0 (2023-12-22)
### Breaking Changes

- Type of `AmazonMWSLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AmazonRdsForLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AmazonRdsForSQLServerLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AmazonRedshiftLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AmazonS3CompatibleLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AmazonS3LinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AsanaLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureBatchLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureBlobFSLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureDataLakeAnalyticsLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureDataLakeStoreLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureDatabricksDetltaLakeLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureDatabricksLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureFileStorageLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureFunctionLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureMLLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureMLServiceLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureMariaDBLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureMySQLLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzurePostgreSQLLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureSQLDWLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureSQLDatabaseLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureSQLMILinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `AzureSearchLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `CassandraLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `CommonDataServiceForAppsLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `ConcurLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `CosmosDbLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `CouchbaseLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `DataworldLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `Db2LinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `DrillLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `DynamicsAXLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `DynamicsCrmLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `DynamicsLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `EloquaLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `FileServerLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `FtpServerLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `GoogleAdWordsLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `GoogleBigQueryLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `GoogleCloudStorageLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `GoogleSheetsLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `GreenplumLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `HBaseLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `HDInsightHiveActivityTypeProperties.Variables` has been changed from `[]any` to `map[string]any`
- Type of `HDInsightLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `HDInsightOnDemandLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `HTTPLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `HdfsLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `HiveLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `HubspotLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `ImpalaLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `InformixLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `JiraLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `MagentoLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `MariaDBLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `MarketoLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `MicrosoftAccessLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `MongoDbLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `MySQLLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `NetezzaLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `ODataLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `OdbcLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `Office365LinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `OracleCloudStorageLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `OracleLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `OracleServiceCloudLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `PaypalLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `PhoenixLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `PostgreSQLLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `PrestoLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `QuickBooksLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `QuickbaseLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `ResponsysLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `RestResourceDatasetTypeProperties.AdditionalHeaders` has been changed from `any` to `map[string]any`
- Type of `RestResourceDatasetTypeProperties.PaginationRules` has been changed from `any` to `map[string]any`
- Type of `RestServiceLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SQLServerLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SalesforceLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SalesforceMarketingCloudLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SalesforceServiceCloudLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SapBWLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SapCloudForCustomerLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SapHanaLinkedServiceProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SapOdpLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SapOpenHubLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SapTableLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `ServiceNowLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SftpServerLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SharePointOnlineListLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `ShopifyLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SmartsheetLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SnowflakeLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SparkLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SquareLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `SybaseLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `TeamDeskLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `TeradataLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `VerticaLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `XeroLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `ZendeskLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Type of `ZohoLinkedServiceTypeProperties.EncryptedCredential` has been changed from `any` to `*string`
- Enum `CosmosDbServicePrincipalCredentialType` has been removed
- Enum `SalesforceSourceReadBehavior` has been removed
- Field `EnablePartitionDiscovery`, `PartitionRootPath` of struct `HTTPReadSettings` has been removed

### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- Type of `AmazonS3CompatibleReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `AmazonS3ReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `AzureBlobFSReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `AzureBlobStorageLinkedServiceTypeProperties.AccountKind` has been changed from `*string` to `any`
- Type of `AzureBlobStorageLinkedServiceTypeProperties.ServiceEndpoint` has been changed from `*string` to `any`
- Type of `AzureBlobStorageReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `AzureDataLakeStoreReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `AzureFileStorageReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `CosmosDbLinkedServiceTypeProperties.ServicePrincipalCredentialType` has been changed from `*CosmosDbServicePrincipalCredentialType` to `any`
- Type of `FileServerReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `FtpReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `FtpReadSettings.UseBinaryTransfer` has been changed from `*bool` to `any`
- Type of `GoogleCloudStorageReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `HdfsReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `OracleCloudStorageReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `SalesforceServiceCloudSource.ReadBehavior` has been changed from `*SalesforceSourceReadBehavior` to `any`
- Type of `SalesforceSource.ReadBehavior` has been changed from `*SalesforceSourceReadBehavior` to `any`
- Type of `SapEccLinkedServiceTypeProperties.URL` has been changed from `*string` to `any`
- Type of `SapEccLinkedServiceTypeProperties.Username` has been changed from `*string` to `any`
- Type of `SftpReadSettings.EnablePartitionDiscovery` has been changed from `*bool` to `any`
- Type of `SynapseNotebookActivityTypeProperties.NumExecutors` has been changed from `*int32` to `any`
- New enum type `ActivityOnInactiveMarkAs` with values `ActivityOnInactiveMarkAsFailed`, `ActivityOnInactiveMarkAsSkipped`, `ActivityOnInactiveMarkAsSucceeded`
- New enum type `ActivityState` with values `ActivityStateActive`, `ActivityStateInactive`
- New enum type `ConnectionType` with values `ConnectionTypeLinkedservicetype`
- New enum type `FrequencyType` with values `FrequencyTypeHour`, `FrequencyTypeMinute`, `FrequencyTypeSecond`
- New enum type `MappingType` with values `MappingTypeAggregate`, `MappingTypeDerived`, `MappingTypeDirect`
- New function `NewChangeDataCaptureClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ChangeDataCaptureClient, error)`
- New function `*ChangeDataCaptureClient.CreateOrUpdate(context.Context, string, string, string, ChangeDataCaptureResource, *ChangeDataCaptureClientCreateOrUpdateOptions) (ChangeDataCaptureClientCreateOrUpdateResponse, error)`
- New function `*ChangeDataCaptureClient.Delete(context.Context, string, string, string, *ChangeDataCaptureClientDeleteOptions) (ChangeDataCaptureClientDeleteResponse, error)`
- New function `*ChangeDataCaptureClient.Get(context.Context, string, string, string, *ChangeDataCaptureClientGetOptions) (ChangeDataCaptureClientGetResponse, error)`
- New function `*ChangeDataCaptureClient.NewListByFactoryPager(string, string, *ChangeDataCaptureClientListByFactoryOptions) *runtime.Pager[ChangeDataCaptureClientListByFactoryResponse]`
- New function `*ChangeDataCaptureClient.Start(context.Context, string, string, string, *ChangeDataCaptureClientStartOptions) (ChangeDataCaptureClientStartResponse, error)`
- New function `*ChangeDataCaptureClient.Status(context.Context, string, string, string, *ChangeDataCaptureClientStatusOptions) (ChangeDataCaptureClientStatusResponse, error)`
- New function `*ChangeDataCaptureClient.Stop(context.Context, string, string, string, *ChangeDataCaptureClientStopOptions) (ChangeDataCaptureClientStopResponse, error)`
- New function `*ClientFactory.NewChangeDataCaptureClient() *ChangeDataCaptureClient`
- New function `*LakeHouseLinkedService.GetLinkedService() *LinkedService`
- New function `*LakeHouseLocation.GetDatasetLocation() *DatasetLocation`
- New function `*LakeHouseReadSettings.GetStoreReadSettings() *StoreReadSettings`
- New function `*LakeHouseTableDataset.GetDataset() *Dataset`
- New function `*LakeHouseTableSink.GetCopySink() *CopySink`
- New function `*LakeHouseTableSource.GetCopySource() *CopySource`
- New function `*LakeHouseWriteSettings.GetStoreWriteSettings() *StoreWriteSettings`
- New function `*ParquetReadSettings.GetFormatReadSettings() *FormatReadSettings`
- New struct `ChangeDataCapture`
- New struct `ChangeDataCaptureFolder`
- New struct `ChangeDataCaptureListResponse`
- New struct `ChangeDataCaptureResource`
- New struct `DataMapperMapping`
- New struct `IntegrationRuntimeDataFlowPropertiesCustomPropertiesItem`
- New struct `LakeHouseLinkedService`
- New struct `LakeHouseLinkedServiceTypeProperties`
- New struct `LakeHouseLocation`
- New struct `LakeHouseReadSettings`
- New struct `LakeHouseTableDataset`
- New struct `LakeHouseTableDatasetTypeProperties`
- New struct `LakeHouseTableSink`
- New struct `LakeHouseTableSource`
- New struct `LakeHouseWriteSettings`
- New struct `MapperAttributeMapping`
- New struct `MapperAttributeMappings`
- New struct `MapperAttributeReference`
- New struct `MapperConnection`
- New struct `MapperConnectionReference`
- New struct `MapperDslConnectorProperties`
- New struct `MapperPolicy`
- New struct `MapperPolicyRecurrence`
- New struct `MapperSourceConnectionsInfo`
- New struct `MapperTable`
- New struct `MapperTableProperties`
- New struct `MapperTableSchema`
- New struct `MapperTargetConnectionsInfo`
- New struct `ParquetReadSettings`
- New struct `SecureInputOutputPolicy`
- New field `OnInactiveMarkAs`, `State` in struct `Activity`
- New field `IsolationLevel` in struct `AmazonRdsForSQLServerSource`
- New field `OnInactiveMarkAs`, `State` in struct `AppendVariableActivity`
- New field `OnInactiveMarkAs`, `State` in struct `AzureDataExplorerCommandActivity`
- New field `OnInactiveMarkAs`, `State` in struct `AzureFunctionActivity`
- New field `OnInactiveMarkAs`, `State` in struct `AzureMLBatchExecutionActivity`
- New field `OnInactiveMarkAs`, `State` in struct `AzureMLExecutePipelineActivity`
- New field `Authentication` in struct `AzureMLServiceLinkedServiceTypeProperties`
- New field `OnInactiveMarkAs`, `State` in struct `AzureMLUpdateResourceActivity`
- New field `IsolationLevel` in struct `AzureSQLSource`
- New field `OnInactiveMarkAs`, `State` in struct `ControlActivity`
- New field `OnInactiveMarkAs`, `State` in struct `CopyActivity`
- New field `OnInactiveMarkAs`, `State` in struct `CustomActivity`
- New field `OnInactiveMarkAs`, `State` in struct `DataLakeAnalyticsUSQLActivity`
- New field `OnInactiveMarkAs`, `State` in struct `DatabricksNotebookActivity`
- New field `OnInactiveMarkAs`, `State` in struct `DatabricksSparkJarActivity`
- New field `OnInactiveMarkAs`, `State` in struct `DatabricksSparkPythonActivity`
- New field `OnInactiveMarkAs`, `State` in struct `DeleteActivity`
- New field `OnInactiveMarkAs`, `State` in struct `ExecuteDataFlowActivity`
- New field `OnInactiveMarkAs`, `State` in struct `ExecutePipelineActivity`
- New field `OnInactiveMarkAs`, `State` in struct `ExecuteSSISPackageActivity`
- New field `OnInactiveMarkAs`, `State` in struct `ExecuteWranglingDataflowActivity`
- New field `OnInactiveMarkAs`, `State` in struct `ExecutionActivity`
- New field `OnInactiveMarkAs`, `State` in struct `FailActivity`
- New field `OnInactiveMarkAs`, `State` in struct `FilterActivity`
- New field `OnInactiveMarkAs`, `State` in struct `ForEachActivity`
- New field `OnInactiveMarkAs`, `State` in struct `GetMetadataActivity`
- New field `GoogleAdsAPIVersion`, `LoginCustomerID`, `PrivateKey`, `SupportLegacyDataTypes` in struct `GoogleAdWordsLinkedServiceTypeProperties`
- New field `OnInactiveMarkAs`, `State` in struct `HDInsightHiveActivity`
- New field `OnInactiveMarkAs`, `State` in struct `HDInsightMapReduceActivity`
- New field `OnInactiveMarkAs`, `State` in struct `HDInsightPigActivity`
- New field `OnInactiveMarkAs`, `State` in struct `HDInsightSparkActivity`
- New field `OnInactiveMarkAs`, `State` in struct `HDInsightStreamingActivity`
- New field `AdditionalColumns` in struct `HTTPReadSettings`
- New field `OnInactiveMarkAs`, `State` in struct `IfConditionActivity`
- New field `CustomProperties` in struct `IntegrationRuntimeDataFlowProperties`
- New field `OnInactiveMarkAs`, `State` in struct `LookupActivity`
- New field `DriverVersion` in struct `MongoDbAtlasLinkedServiceTypeProperties`
- New field `FormatSettings` in struct `ParquetSource`
- New field `NumberOfExternalNodes`, `NumberOfPipelineNodes` in struct `PipelineExternalComputeScaleProperties`
- New field `IsolationLevel` in struct `SQLDWSource`
- New field `IsolationLevel` in struct `SQLMISource`
- New field `IsolationLevel` in struct `SQLServerSource`
- New field `OnInactiveMarkAs`, `State` in struct `SQLServerStoredProcedureActivity`
- New field `OnInactiveMarkAs`, `State` in struct `ScriptActivity`
- New field `SelfContainedInteractiveAuthoringEnabled` in struct `SelfHostedIntegrationRuntimeStatusTypeProperties`
- New field `SelfContainedInteractiveAuthoringEnabled` in struct `SelfHostedIntegrationRuntimeTypeProperties`
- New field `OnInactiveMarkAs`, `Policy`, `State` in struct `SetVariableActivity`
- New field `SetSystemVariable` in struct `SetVariableActivityTypeProperties`
- New field `OnInactiveMarkAs`, `State` in struct `SwitchActivity`
- New field `OnInactiveMarkAs`, `State` in struct `SynapseNotebookActivity`
- New field `ConfigurationType`, `SparkConfig`, `TargetSparkConfiguration` in struct `SynapseNotebookActivityTypeProperties`
- New field `OnInactiveMarkAs`, `State` in struct `SynapseSparkJobDefinitionActivity`
- New field `OnInactiveMarkAs`, `State` in struct `UntilActivity`
- New field `OnInactiveMarkAs`, `State` in struct `ValidationActivity`
- New field `OnInactiveMarkAs`, `State` in struct `WaitActivity`
- New field `OnInactiveMarkAs`, `State` in struct `WebActivity`
- New field `OnInactiveMarkAs`, `Policy`, `State` in struct `WebHookActivity`


## 3.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 3.2.0 (2023-03-24)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New enum type `AzureStorageAuthenticationType` with values `AzureStorageAuthenticationTypeAccountKey`, `AzureStorageAuthenticationTypeAnonymous`, `AzureStorageAuthenticationTypeMsi`, `AzureStorageAuthenticationTypeSasURI`, `AzureStorageAuthenticationTypeServicePrincipal`
- New struct `CopyComputeScaleProperties`
- New struct `PipelineExternalComputeScaleProperties`
- New field `SasToken` in struct `AzureBlobFSLinkedServiceTypeProperties`
- New field `SasURI` in struct `AzureBlobFSLinkedServiceTypeProperties`
- New field `AuthenticationType` in struct `AzureBlobStorageLinkedServiceTypeProperties`
- New field `ContainerURI` in struct `AzureBlobStorageLinkedServiceTypeProperties`
- New field `CopyComputeScaleProperties` in struct `IntegrationRuntimeComputeProperties`
- New field `PipelineExternalComputeScaleProperties` in struct `IntegrationRuntimeComputeProperties`


## 3.1.0 (2023-02-24)
### Features Added

- Type of `SynapseSparkJobActivityTypeProperties.NumExecutors` has been changed from `*int32` to `any`
- New type alias `ConfigurationType` with values `ConfigurationTypeArtifact`, `ConfigurationTypeCustomized`, `ConfigurationTypeDefault`
- New type alias `SparkConfigurationReferenceType` with values `SparkConfigurationReferenceTypeSparkConfigurationReference`
- New function `*Credential.GetCredential() *Credential`
- New function `NewCredentialOperationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CredentialOperationsClient, error)`
- New function `*CredentialOperationsClient.CreateOrUpdate(context.Context, string, string, string, ManagedIdentityCredentialResource, *CredentialOperationsClientCreateOrUpdateOptions) (CredentialOperationsClientCreateOrUpdateResponse, error)`
- New function `*CredentialOperationsClient.Delete(context.Context, string, string, string, *CredentialOperationsClientDeleteOptions) (CredentialOperationsClientDeleteResponse, error)`
- New function `*CredentialOperationsClient.Get(context.Context, string, string, string, *CredentialOperationsClientGetOptions) (CredentialOperationsClientGetResponse, error)`
- New function `*CredentialOperationsClient.NewListByFactoryPager(string, string, *CredentialOperationsClientListByFactoryOptions) *runtime.Pager[CredentialOperationsClientListByFactoryResponse]`
- New function `*ManagedIdentityCredential.GetCredential() *Credential`
- New function `*ServicePrincipalCredential.GetCredential() *Credential`
- New struct `CredentialListResponse`
- New struct `CredentialOperationsClient`
- New struct `CredentialOperationsClientListByFactoryResponse`
- New struct `ManagedIdentityCredential`
- New struct `ManagedIdentityCredentialResource`
- New struct `ManagedIdentityTypeProperties`
- New struct `ServicePrincipalCredential`
- New struct `ServicePrincipalCredentialTypeProperties`
- New struct `SparkConfigurationParametrizationReference`
- New field `ConfigurationType` in struct `SynapseSparkJobActivityTypeProperties`
- New field `ScanFolder` in struct `SynapseSparkJobActivityTypeProperties`
- New field `SparkConfig` in struct `SynapseSparkJobActivityTypeProperties`
- New field `TargetSparkConfiguration` in struct `SynapseSparkJobActivityTypeProperties`


## 3.0.0 (2022-10-27)
### Breaking Changes

- Type of `SynapseSparkJobReference.ReferenceName` has been changed from `*string` to `interface{}`

### Features Added

- New field `WorkspaceResourceID` in struct `AzureSynapseArtifactsLinkedServiceTypeProperties`
- New field `DisablePublish` in struct `FactoryRepoConfiguration`
- New field `DisablePublish` in struct `FactoryGitHubConfiguration`
- New field `DisablePublish` in struct `FactoryVSTSConfiguration`
- New field `ScriptBlockExecutionTimeout` in struct `ScriptActivityTypeProperties`
- New field `PythonCodeReference` in struct `SynapseSparkJobActivityTypeProperties`
- New field `FilesV2` in struct `SynapseSparkJobActivityTypeProperties`


## 2.0.0 (2022-10-10)
### Breaking Changes

- Type of `SQLMISource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLMISink.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `AzureSQLSink.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLServerSource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLServerSink.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `AzureSQLSource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLSink.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `SQLSource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`
- Type of `AmazonRdsForSQLServerSource.StoredProcedureParameters` has been changed from `map[string]*StoredProcedureParameter` to `interface{}`


## 1.3.0 (2022-09-07)
### Features Added

- New const `NotebookParameterTypeBool`
- New const `NotebookReferenceTypeNotebookReference`
- New const `NotebookParameterTypeString`
- New const `SparkJobReferenceTypeSparkJobDefinitionReference`
- New const `NotebookParameterTypeInt`
- New const `BigDataPoolReferenceTypeBigDataPoolReference`
- New const `NotebookParameterTypeFloat`
- New type alias `NotebookParameterType`
- New type alias `SparkJobReferenceType`
- New type alias `NotebookReferenceType`
- New type alias `BigDataPoolReferenceType`
- New function `*AzureSynapseArtifactsLinkedService.GetLinkedService() *LinkedService`
- New function `PossibleBigDataPoolReferenceTypeValues() []BigDataPoolReferenceType`
- New function `PossibleNotebookParameterTypeValues() []NotebookParameterType`
- New function `*SynapseSparkJobDefinitionActivity.GetExecutionActivity() *ExecutionActivity`
- New function `*GoogleSheetsLinkedService.GetLinkedService() *LinkedService`
- New function `*SynapseNotebookActivity.GetExecutionActivity() *ExecutionActivity`
- New function `PossibleNotebookReferenceTypeValues() []NotebookReferenceType`
- New function `PossibleSparkJobReferenceTypeValues() []SparkJobReferenceType`
- New function `*SynapseNotebookActivity.GetActivity() *Activity`
- New function `*SynapseSparkJobDefinitionActivity.GetActivity() *Activity`
- New struct `AzureSynapseArtifactsLinkedService`
- New struct `AzureSynapseArtifactsLinkedServiceTypeProperties`
- New struct `BigDataPoolParametrizationReference`
- New struct `GoogleSheetsLinkedService`
- New struct `GoogleSheetsLinkedServiceTypeProperties`
- New struct `NotebookParameter`
- New struct `SynapseNotebookActivity`
- New struct `SynapseNotebookActivityTypeProperties`
- New struct `SynapseNotebookReference`
- New struct `SynapseSparkJobActivityTypeProperties`
- New struct `SynapseSparkJobDefinitionActivity`
- New struct `SynapseSparkJobReference`


## 1.2.0 (2022-06-15)
### Features Added

- New field `ClientSecret` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Resource` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Scope` in struct `RestServiceLinkedServiceTypeProperties`
- New field `TokenEndpoint` in struct `RestServiceLinkedServiceTypeProperties`
- New field `ClientID` in struct `RestServiceLinkedServiceTypeProperties`


## 1.1.0 (2022-05-30)
### Features Added

- New function `GlobalParameterResource.MarshalJSON() ([]byte, error)`
- New struct `GlobalParameterListResponse`
- New struct `GlobalParameterResource`
- New struct `GlobalParametersClientCreateOrUpdateOptions`
- New struct `GlobalParametersClientCreateOrUpdateResponse`
- New struct `GlobalParametersClientDeleteOptions`
- New struct `GlobalParametersClientDeleteResponse`
- New struct `GlobalParametersClientGetOptions`
- New struct `GlobalParametersClientGetResponse`
- New struct `GlobalParametersClientListByFactoryOptions`
- New struct `GlobalParametersClientListByFactoryResponse`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).## 1.2.0 (2022-06-15)
### Features Added

- New field `ClientSecret` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Resource` in struct `RestServiceLinkedServiceTypeProperties`
- New field `Scope` in struct `RestServiceLinkedServiceTypeProperties`
- New field `TokenEndpoint` in struct `RestServiceLinkedServiceTypeProperties`
- New field `ClientID` in struct `RestServiceLinkedServiceTypeProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).