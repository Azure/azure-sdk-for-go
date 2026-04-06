# Release History

## 3.0.0 (2025-11-12)
### Breaking Changes

- Function `*ProjectCapabilityHostsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, CapabilityHost, *ProjectCapabilityHostsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, string, string, ProjectCapabilityHost, *ProjectCapabilityHostsClientBeginCreateOrUpdateOptions)`
- Field `CapabilityHost` of struct `ProjectCapabilityHostsClientCreateOrUpdateResponse` has been removed
- Field `CapabilityHost` of struct `ProjectCapabilityHostsClientGetResponse` has been removed

### Features Added

- New value `ConnectionCategoryAzureStorageAccount` added to enum type `ConnectionCategory`
- New value `ModelLifecycleStatusLegacy` added to enum type `ModelLifecycleStatus`
- New enum type `DeprecationStatus` with values `DeprecationStatusPlanned`, `DeprecationStatusTentative`
- New enum type `TierUpgradePolicy` with values `TierUpgradePolicyNoAutoUpgrade`, `TierUpgradePolicyOnceUpgradeIsAvailable`
- New enum type `UpgradeAvailabilityStatus` with values `UpgradeAvailabilityStatusAvailable`, `UpgradeAvailabilityStatusNotAvailable`
- New function `*AccountCapabilityHostsClient.NewListPager(string, string, *AccountCapabilityHostsClientListOptions) *runtime.Pager[AccountCapabilityHostsClientListResponse]`
- New function `*ClientFactory.NewQuotaTiersClient() *QuotaTiersClient`
- New function `*ClientFactory.NewRaiTopicsClient() *RaiTopicsClient`
- New function `*ProjectCapabilityHostsClient.NewListPager(string, string, string, *ProjectCapabilityHostsClientListOptions) *runtime.Pager[ProjectCapabilityHostsClientListResponse]`
- New function `NewQuotaTiersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*QuotaTiersClient, error)`
- New function `*QuotaTiersClient.CreateOrUpdate(context.Context, string, QuotaTier, *QuotaTiersClientCreateOrUpdateOptions) (QuotaTiersClientCreateOrUpdateResponse, error)`
- New function `*QuotaTiersClient.Get(context.Context, string, *QuotaTiersClientGetOptions) (QuotaTiersClientGetResponse, error)`
- New function `*QuotaTiersClient.NewListBySubscriptionPager(*QuotaTiersClientListBySubscriptionOptions) *runtime.Pager[QuotaTiersClientListBySubscriptionResponse]`
- New function `*QuotaTiersClient.Update(context.Context, string, QuotaTier, *QuotaTiersClientUpdateOptions) (QuotaTiersClientUpdateResponse, error)`
- New function `NewRaiTopicsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiTopicsClient, error)`
- New function `*RaiTopicsClient.CreateOrUpdate(context.Context, string, string, string, RaiTopic, *RaiTopicsClientCreateOrUpdateOptions) (RaiTopicsClientCreateOrUpdateResponse, error)`
- New function `*RaiTopicsClient.BeginDelete(context.Context, string, string, string, *RaiTopicsClientBeginDeleteOptions) (*runtime.Poller[RaiTopicsClientDeleteResponse], error)`
- New function `*RaiTopicsClient.Get(context.Context, string, string, string, *RaiTopicsClientGetOptions) (RaiTopicsClientGetResponse, error)`
- New function `*RaiTopicsClient.NewListPager(string, string, *RaiTopicsClientListOptions) *runtime.Pager[RaiTopicsClientListResponse]`
- New struct `CapabilityHostResourceArmPaginatedResult`
- New struct `CustomTopicConfig`
- New struct `ProjectCapabilityHost`
- New struct `ProjectCapabilityHostProperties`
- New struct `ProjectCapabilityHostResourceArmPaginatedResult`
- New struct `QuotaTier`
- New struct `QuotaTierListResult`
- New struct `QuotaTierProperties`
- New struct `QuotaTierUpgradeEligibilityInfo`
- New struct `RaiTopic`
- New struct `RaiTopicConfig`
- New struct `RaiTopicProperties`
- New struct `RaiTopicResult`
- New struct `ReplacementConfig`
- New field `ModelCatalogAssetID`, `ReplacementConfig` in struct `AccountModel`
- New field `StoredCompletionsDisabled` in struct `AccountProperties`
- New field `SystemData` in struct `AzureEntityResource`
- New field `SystemData` in struct `CapabilityHost`
- New field `SystemData` in struct `ConnectionPropertiesV2BasicResource`
- New field `SystemData` in struct `ModelCapacityListResultValueItem`
- New field `DeprecationStatus` in struct `ModelDeprecationInfo`
- New field `SystemData` in struct `NetworkSecurityPerimeterConfiguration`
- New field `SystemData` in struct `PrivateLinkResource`
- New anonymous field `ProjectCapabilityHost` in struct `ProjectCapabilityHostsClientCreateOrUpdateResponse`
- New anonymous field `ProjectCapabilityHost` in struct `ProjectCapabilityHostsClientGetResponse`
- New field `SystemData` in struct `ProxyResource`
- New field `SystemData` in struct `RaiContentFilter`
- New field `CustomTopics` in struct `RaiPolicyProperties`
- New field `SystemData` in struct `Resource`


## 2.0.0 (2025-09-15)
### Breaking Changes

- Type of `AccountProperties.NetworkInjections` has been changed from `*NetworkInjections` to `[]*NetworkInjection`
- Struct `NetworkInjections` has been removed

### Features Added

- New struct `NetworkInjection`


## 1.8.0 (2025-07-25)
### Features Added

- New value `ProvisioningStateCanceled` added to enum type `ProvisioningState`
- New enum type `CapabilityHostKind` with values `CapabilityHostKindAgents`
- New enum type `CapabilityHostProvisioningState` with values `CapabilityHostProvisioningStateCanceled`, `CapabilityHostProvisioningStateCreating`, `CapabilityHostProvisioningStateDeleting`, `CapabilityHostProvisioningStateFailed`, `CapabilityHostProvisioningStateSucceeded`, `CapabilityHostProvisioningStateUpdating`
- New enum type `ConnectionAuthType` with values `ConnectionAuthTypeAAD`, `ConnectionAuthTypeAPIKey`, `ConnectionAuthTypeAccessKey`, `ConnectionAuthTypeAccountKey`, `ConnectionAuthTypeCustomKeys`, `ConnectionAuthTypeManagedIdentity`, `ConnectionAuthTypeNone`, `ConnectionAuthTypeOAuth2`, `ConnectionAuthTypePAT`, `ConnectionAuthTypeSAS`, `ConnectionAuthTypeServicePrincipal`, `ConnectionAuthTypeUsernamePassword`
- New enum type `ConnectionCategory` with values `ConnectionCategoryADLSGen2`, `ConnectionCategoryAIServices`, `ConnectionCategoryAPIKey`, `ConnectionCategoryAmazonMws`, `ConnectionCategoryAmazonRdsForOracle`, `ConnectionCategoryAmazonRdsForSQLServer`, `ConnectionCategoryAmazonRedshift`, `ConnectionCategoryAmazonS3Compatible`, `ConnectionCategoryAzureBlob`, `ConnectionCategoryAzureDataExplorer`, `ConnectionCategoryAzureDatabricksDeltaLake`, `ConnectionCategoryAzureMariaDb`, `ConnectionCategoryAzureMySQLDb`, `ConnectionCategoryAzureOneLake`, `ConnectionCategoryAzureOpenAI`, `ConnectionCategoryAzurePostgresDb`, `ConnectionCategoryAzureSQLDb`, `ConnectionCategoryAzureSQLMi`, `ConnectionCategoryAzureSynapseAnalytics`, `ConnectionCategoryAzureTableStorage`, `ConnectionCategoryBingLLMSearch`, `ConnectionCategoryCassandra`, `ConnectionCategoryCognitiveSearch`, `ConnectionCategoryCognitiveService`, `ConnectionCategoryConcur`, `ConnectionCategoryContainerRegistry`, `ConnectionCategoryCosmosDb`, `ConnectionCategoryCosmosDbMongoDbAPI`, `ConnectionCategoryCouchbase`, `ConnectionCategoryCustomKeys`, `ConnectionCategoryDb2`, `ConnectionCategoryDrill`, `ConnectionCategoryDynamics`, `ConnectionCategoryDynamicsAx`, `ConnectionCategoryDynamicsCrm`, `ConnectionCategoryElasticsearch`, `ConnectionCategoryEloqua`, `ConnectionCategoryFileServer`, `ConnectionCategoryFtpServer`, `ConnectionCategoryGenericContainerRegistry`, `ConnectionCategoryGenericHTTP`, `ConnectionCategoryGenericRest`, `ConnectionCategoryGit`, `ConnectionCategoryGoogleAdWords`, `ConnectionCategoryGoogleBigQuery`, `ConnectionCategoryGoogleCloudStorage`, `ConnectionCategoryGreenplum`, `ConnectionCategoryHbase`, `ConnectionCategoryHdfs`, `ConnectionCategoryHive`, `ConnectionCategoryHubspot`, `ConnectionCategoryImpala`, `ConnectionCategoryInformix`, `ConnectionCategoryJira`, `ConnectionCategoryMagento`, `ConnectionCategoryManagedOnlineEndpoint`, `ConnectionCategoryMariaDb`, `ConnectionCategoryMarketo`, `ConnectionCategoryMicrosoftAccess`, `ConnectionCategoryMongoDbAtlas`, `ConnectionCategoryMongoDbV2`, `ConnectionCategoryMySQL`, `ConnectionCategoryNetezza`, `ConnectionCategoryODataRest`, `ConnectionCategoryOdbc`, `ConnectionCategoryOffice365`, `ConnectionCategoryOpenAI`, `ConnectionCategoryOracle`, `ConnectionCategoryOracleCloudStorage`, `ConnectionCategoryOracleServiceCloud`, `ConnectionCategoryPayPal`, `ConnectionCategoryPhoenix`, `ConnectionCategoryPinecone`, `ConnectionCategoryPostgreSQL`, `ConnectionCategoryPresto`, `ConnectionCategoryPythonFeed`, `ConnectionCategoryQuickBooks`, `ConnectionCategoryRedis`, `ConnectionCategoryResponsys`, `ConnectionCategoryS3`, `ConnectionCategorySQLServer`, `ConnectionCategorySalesforce`, `ConnectionCategorySalesforceMarketingCloud`, `ConnectionCategorySalesforceServiceCloud`, `ConnectionCategorySapBw`, `ConnectionCategorySapCloudForCustomer`, `ConnectionCategorySapEcc`, `ConnectionCategorySapHana`, `ConnectionCategorySapOpenHub`, `ConnectionCategorySapTable`, `ConnectionCategorySerp`, `ConnectionCategoryServerless`, `ConnectionCategoryServiceNow`, `ConnectionCategorySftp`, `ConnectionCategorySharePointOnlineList`, `ConnectionCategoryShopify`, `ConnectionCategorySnowflake`, `ConnectionCategorySpark`, `ConnectionCategorySquare`, `ConnectionCategorySybase`, `ConnectionCategoryTeradata`, `ConnectionCategoryVertica`, `ConnectionCategoryWebTable`, `ConnectionCategoryXero`, `ConnectionCategoryZoho`
- New enum type `ConnectionGroup` with values `ConnectionGroupAzure`, `ConnectionGroupAzureAI`, `ConnectionGroupDatabase`, `ConnectionGroupFile`, `ConnectionGroupGenericProtocol`, `ConnectionGroupNoSQL`, `ConnectionGroupServicesAndApps`
- New enum type `ManagedPERequirement` with values `ManagedPERequirementNotApplicable`, `ManagedPERequirementNotRequired`, `ManagedPERequirementRequired`
- New enum type `ManagedPEStatus` with values `ManagedPEStatusActive`, `ManagedPEStatusInactive`, `ManagedPEStatusNotApplicable`
- New enum type `ScenarioType` with values `ScenarioTypeAgent`, `ScenarioTypeNone`
- New function `*AADAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*APIKeyAuthConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*AccessKeyAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `NewAccountCapabilityHostsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccountCapabilityHostsClient, error)`
- New function `*AccountCapabilityHostsClient.BeginCreateOrUpdate(context.Context, string, string, string, CapabilityHost, *AccountCapabilityHostsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AccountCapabilityHostsClientCreateOrUpdateResponse], error)`
- New function `*AccountCapabilityHostsClient.BeginDelete(context.Context, string, string, string, *AccountCapabilityHostsClientBeginDeleteOptions) (*runtime.Poller[AccountCapabilityHostsClientDeleteResponse], error)`
- New function `*AccountCapabilityHostsClient.Get(context.Context, string, string, string, *AccountCapabilityHostsClientGetOptions) (AccountCapabilityHostsClientGetResponse, error)`
- New function `NewAccountConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccountConnectionsClient, error)`
- New function `*AccountConnectionsClient.Create(context.Context, string, string, string, *AccountConnectionsClientCreateOptions) (AccountConnectionsClientCreateResponse, error)`
- New function `*AccountConnectionsClient.Delete(context.Context, string, string, string, *AccountConnectionsClientDeleteOptions) (AccountConnectionsClientDeleteResponse, error)`
- New function `*AccountConnectionsClient.Get(context.Context, string, string, string, *AccountConnectionsClientGetOptions) (AccountConnectionsClientGetResponse, error)`
- New function `*AccountConnectionsClient.NewListPager(string, string, *AccountConnectionsClientListOptions) *runtime.Pager[AccountConnectionsClientListResponse]`
- New function `*AccountConnectionsClient.Update(context.Context, string, string, string, *AccountConnectionsClientUpdateOptions) (AccountConnectionsClientUpdateResponse, error)`
- New function `*AccountKeyAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*ClientFactory.NewAccountCapabilityHostsClient() *AccountCapabilityHostsClient`
- New function `*ClientFactory.NewAccountConnectionsClient() *AccountConnectionsClient`
- New function `*ClientFactory.NewProjectCapabilityHostsClient() *ProjectCapabilityHostsClient`
- New function `*ClientFactory.NewProjectConnectionsClient() *ProjectConnectionsClient`
- New function `*ClientFactory.NewProjectsClient() *ProjectsClient`
- New function `*ConnectionPropertiesV2.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*CustomKeysConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*ManagedIdentityAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*OAuth2AuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*PATAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `NewProjectCapabilityHostsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectCapabilityHostsClient, error)`
- New function `*ProjectCapabilityHostsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, CapabilityHost, *ProjectCapabilityHostsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ProjectCapabilityHostsClientCreateOrUpdateResponse], error)`
- New function `*ProjectCapabilityHostsClient.BeginDelete(context.Context, string, string, string, string, *ProjectCapabilityHostsClientBeginDeleteOptions) (*runtime.Poller[ProjectCapabilityHostsClientDeleteResponse], error)`
- New function `*ProjectCapabilityHostsClient.Get(context.Context, string, string, string, string, *ProjectCapabilityHostsClientGetOptions) (ProjectCapabilityHostsClientGetResponse, error)`
- New function `NewProjectConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectConnectionsClient, error)`
- New function `*ProjectConnectionsClient.Create(context.Context, string, string, string, string, *ProjectConnectionsClientCreateOptions) (ProjectConnectionsClientCreateResponse, error)`
- New function `*ProjectConnectionsClient.Delete(context.Context, string, string, string, string, *ProjectConnectionsClientDeleteOptions) (ProjectConnectionsClientDeleteResponse, error)`
- New function `*ProjectConnectionsClient.Get(context.Context, string, string, string, string, *ProjectConnectionsClientGetOptions) (ProjectConnectionsClientGetResponse, error)`
- New function `*ProjectConnectionsClient.NewListPager(string, string, string, *ProjectConnectionsClientListOptions) *runtime.Pager[ProjectConnectionsClientListResponse]`
- New function `*ProjectConnectionsClient.Update(context.Context, string, string, string, string, *ProjectConnectionsClientUpdateOptions) (ProjectConnectionsClientUpdateResponse, error)`
- New function `NewProjectsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectsClient, error)`
- New function `*ProjectsClient.BeginCreate(context.Context, string, string, string, Project, *ProjectsClientBeginCreateOptions) (*runtime.Poller[ProjectsClientCreateResponse], error)`
- New function `*ProjectsClient.BeginDelete(context.Context, string, string, string, *ProjectsClientBeginDeleteOptions) (*runtime.Poller[ProjectsClientDeleteResponse], error)`
- New function `*ProjectsClient.Get(context.Context, string, string, string, *ProjectsClientGetOptions) (ProjectsClientGetResponse, error)`
- New function `*ProjectsClient.NewListPager(string, string, *ProjectsClientListOptions) *runtime.Pager[ProjectsClientListResponse]`
- New function `*ProjectsClient.BeginUpdate(context.Context, string, string, string, Project, *ProjectsClientBeginUpdateOptions) (*runtime.Poller[ProjectsClientUpdateResponse], error)`
- New function `*SASAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*ServicePrincipalAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*UsernamePasswordAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*NoneAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New struct `AADAuthTypeConnectionProperties`
- New struct `APIKeyAuthConnectionProperties`
- New struct `AccessKeyAuthTypeConnectionProperties`
- New struct `AccountKeyAuthTypeConnectionProperties`
- New struct `CapabilityHost`
- New struct `CapabilityHostProperties`
- New struct `ConnectionAPIKey`
- New struct `ConnectionAccessKey`
- New struct `ConnectionAccountKey`
- New struct `ConnectionManagedIdentity`
- New struct `ConnectionOAuth2`
- New struct `ConnectionPersonalAccessToken`
- New struct `ConnectionPropertiesV2BasicResource`
- New struct `ConnectionPropertiesV2BasicResourceArmPaginatedResult`
- New struct `ConnectionServicePrincipal`
- New struct `ConnectionSharedAccessSignature`
- New struct `ConnectionUpdateContent`
- New struct `ConnectionUsernamePassword`
- New struct `CustomKeys`
- New struct `CustomKeysConnectionProperties`
- New struct `ManagedIdentityAuthTypeConnectionProperties`
- New struct `NetworkInjections`
- New struct `NoneAuthTypeConnectionProperties`
- New struct `OAuth2AuthTypeConnectionProperties`
- New struct `PATAuthTypeConnectionProperties`
- New struct `Project`
- New struct `ProjectListResult`
- New struct `ProjectProperties`
- New struct `ResourceBase`
- New struct `SASAuthTypeConnectionProperties`
- New struct `ServicePrincipalAuthTypeConnectionProperties`
- New struct `UsernamePasswordAuthTypeConnectionProperties`
- New field `AllowProjectManagement`, `AssociatedProjects`, `DefaultProject`, `NetworkInjections` in struct `AccountProperties`
- New field `SpilloverDeploymentName` in struct `DeploymentProperties`


## 1.8.0-beta.1 (2025-05-12)
### Features Added

- New value `ProvisioningStateCanceled` added to enum type `ProvisioningState`
- New enum type `CapabilityHostKind` with values `CapabilityHostKindAgents`
- New enum type `CapabilityHostProvisioningState` with values `CapabilityHostProvisioningStateCanceled`, `CapabilityHostProvisioningStateCreating`, `CapabilityHostProvisioningStateDeleting`, `CapabilityHostProvisioningStateFailed`, `CapabilityHostProvisioningStateSucceeded`, `CapabilityHostProvisioningStateUpdating`
- New enum type `ConnectionAuthType` with values `ConnectionAuthTypeAAD`, `ConnectionAuthTypeAPIKey`, `ConnectionAuthTypeAccessKey`, `ConnectionAuthTypeAccountKey`, `ConnectionAuthTypeCustomKeys`, `ConnectionAuthTypeManagedIdentity`, `ConnectionAuthTypeNone`, `ConnectionAuthTypeOAuth2`, `ConnectionAuthTypePAT`, `ConnectionAuthTypeSAS`, `ConnectionAuthTypeServicePrincipal`, `ConnectionAuthTypeUsernamePassword`
- New enum type `ConnectionCategory` with values `ConnectionCategoryADLSGen2`, `ConnectionCategoryAIServices`, `ConnectionCategoryAPIKey`, `ConnectionCategoryAmazonMws`, `ConnectionCategoryAmazonRdsForOracle`, `ConnectionCategoryAmazonRdsForSQLServer`, `ConnectionCategoryAmazonRedshift`, `ConnectionCategoryAmazonS3Compatible`, `ConnectionCategoryAzureBlob`, `ConnectionCategoryAzureDataExplorer`, `ConnectionCategoryAzureDatabricksDeltaLake`, `ConnectionCategoryAzureMariaDb`, `ConnectionCategoryAzureMySQLDb`, `ConnectionCategoryAzureOneLake`, `ConnectionCategoryAzureOpenAI`, `ConnectionCategoryAzurePostgresDb`, `ConnectionCategoryAzureSQLDb`, `ConnectionCategoryAzureSQLMi`, `ConnectionCategoryAzureSynapseAnalytics`, `ConnectionCategoryAzureTableStorage`, `ConnectionCategoryBingLLMSearch`, `ConnectionCategoryCassandra`, `ConnectionCategoryCognitiveSearch`, `ConnectionCategoryCognitiveService`, `ConnectionCategoryConcur`, `ConnectionCategoryContainerRegistry`, `ConnectionCategoryCosmosDb`, `ConnectionCategoryCosmosDbMongoDbAPI`, `ConnectionCategoryCouchbase`, `ConnectionCategoryCustomKeys`, `ConnectionCategoryDb2`, `ConnectionCategoryDrill`, `ConnectionCategoryDynamics`, `ConnectionCategoryDynamicsAx`, `ConnectionCategoryDynamicsCrm`, `ConnectionCategoryElasticsearch`, `ConnectionCategoryEloqua`, `ConnectionCategoryFileServer`, `ConnectionCategoryFtpServer`, `ConnectionCategoryGenericContainerRegistry`, `ConnectionCategoryGenericHTTP`, `ConnectionCategoryGenericRest`, `ConnectionCategoryGit`, `ConnectionCategoryGoogleAdWords`, `ConnectionCategoryGoogleBigQuery`, `ConnectionCategoryGoogleCloudStorage`, `ConnectionCategoryGreenplum`, `ConnectionCategoryHbase`, `ConnectionCategoryHdfs`, `ConnectionCategoryHive`, `ConnectionCategoryHubspot`, `ConnectionCategoryImpala`, `ConnectionCategoryInformix`, `ConnectionCategoryJira`, `ConnectionCategoryMagento`, `ConnectionCategoryManagedOnlineEndpoint`, `ConnectionCategoryMariaDb`, `ConnectionCategoryMarketo`, `ConnectionCategoryMicrosoftAccess`, `ConnectionCategoryMongoDbAtlas`, `ConnectionCategoryMongoDbV2`, `ConnectionCategoryMySQL`, `ConnectionCategoryNetezza`, `ConnectionCategoryODataRest`, `ConnectionCategoryOdbc`, `ConnectionCategoryOffice365`, `ConnectionCategoryOpenAI`, `ConnectionCategoryOracle`, `ConnectionCategoryOracleCloudStorage`, `ConnectionCategoryOracleServiceCloud`, `ConnectionCategoryPayPal`, `ConnectionCategoryPhoenix`, `ConnectionCategoryPinecone`, `ConnectionCategoryPostgreSQL`, `ConnectionCategoryPresto`, `ConnectionCategoryPythonFeed`, `ConnectionCategoryQuickBooks`, `ConnectionCategoryRedis`, `ConnectionCategoryResponsys`, `ConnectionCategoryS3`, `ConnectionCategorySQLServer`, `ConnectionCategorySalesforce`, `ConnectionCategorySalesforceMarketingCloud`, `ConnectionCategorySalesforceServiceCloud`, `ConnectionCategorySapBw`, `ConnectionCategorySapCloudForCustomer`, `ConnectionCategorySapEcc`, `ConnectionCategorySapHana`, `ConnectionCategorySapOpenHub`, `ConnectionCategorySapTable`, `ConnectionCategorySerp`, `ConnectionCategoryServerless`, `ConnectionCategoryServiceNow`, `ConnectionCategorySftp`, `ConnectionCategorySharePointOnlineList`, `ConnectionCategoryShopify`, `ConnectionCategorySnowflake`, `ConnectionCategorySpark`, `ConnectionCategorySquare`, `ConnectionCategorySybase`, `ConnectionCategoryTeradata`, `ConnectionCategoryVertica`, `ConnectionCategoryWebTable`, `ConnectionCategoryXero`, `ConnectionCategoryZoho`
- New enum type `ConnectionGroup` with values `ConnectionGroupAzure`, `ConnectionGroupAzureAI`, `ConnectionGroupDatabase`, `ConnectionGroupFile`, `ConnectionGroupGenericProtocol`, `ConnectionGroupNoSQL`, `ConnectionGroupServicesAndApps`
- New enum type `ManagedPERequirement` with values `ManagedPERequirementNotApplicable`, `ManagedPERequirementNotRequired`, `ManagedPERequirementRequired`
- New enum type `ManagedPEStatus` with values `ManagedPEStatusActive`, `ManagedPEStatusInactive`, `ManagedPEStatusNotApplicable`
- New enum type `ScenarioType` with values `ScenarioTypeAgent`, `ScenarioTypeNone`
- New function `*AADAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*APIKeyAuthConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*AccessKeyAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `NewAccountCapabilityHostsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccountCapabilityHostsClient, error)`
- New function `*AccountCapabilityHostsClient.BeginCreateOrUpdate(context.Context, string, string, string, CapabilityHost, *AccountCapabilityHostsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AccountCapabilityHostsClientCreateOrUpdateResponse], error)`
- New function `*AccountCapabilityHostsClient.BeginDelete(context.Context, string, string, string, *AccountCapabilityHostsClientBeginDeleteOptions) (*runtime.Poller[AccountCapabilityHostsClientDeleteResponse], error)`
- New function `*AccountCapabilityHostsClient.Get(context.Context, string, string, string, *AccountCapabilityHostsClientGetOptions) (AccountCapabilityHostsClientGetResponse, error)`
- New function `NewAccountConnectionClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccountConnectionClient, error)`
- New function `*AccountConnectionClient.Create(context.Context, string, string, string, *AccountConnectionClientCreateOptions) (AccountConnectionClientCreateResponse, error)`
- New function `*AccountConnectionClient.Delete(context.Context, string, string, string, *AccountConnectionClientDeleteOptions) (AccountConnectionClientDeleteResponse, error)`
- New function `*AccountConnectionClient.Get(context.Context, string, string, string, *AccountConnectionClientGetOptions) (AccountConnectionClientGetResponse, error)`
- New function `*AccountConnectionClient.NewListPager(string, string, *AccountConnectionClientListOptions) *runtime.Pager[AccountConnectionClientListResponse]`
- New function `*AccountConnectionClient.Update(context.Context, string, string, string, *AccountConnectionClientUpdateOptions) (AccountConnectionClientUpdateResponse, error)`
- New function `*AccountKeyAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*ClientFactory.NewAccountCapabilityHostsClient() *AccountCapabilityHostsClient`
- New function `*ClientFactory.NewAccountConnectionClient() *AccountConnectionClient`
- New function `*ClientFactory.NewProjectCapabilityHostsClient() *ProjectCapabilityHostsClient`
- New function `*ClientFactory.NewProjectConnectionClient() *ProjectConnectionClient`
- New function `*ClientFactory.NewProjectsClient() *ProjectsClient`
- New function `*ConnectionPropertiesV2.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*CustomKeysConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*ManagedIdentityAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*OAuth2AuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*PATAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `NewProjectCapabilityHostsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectCapabilityHostsClient, error)`
- New function `*ProjectCapabilityHostsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, CapabilityHost, *ProjectCapabilityHostsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ProjectCapabilityHostsClientCreateOrUpdateResponse], error)`
- New function `*ProjectCapabilityHostsClient.BeginDelete(context.Context, string, string, string, string, *ProjectCapabilityHostsClientBeginDeleteOptions) (*runtime.Poller[ProjectCapabilityHostsClientDeleteResponse], error)`
- New function `*ProjectCapabilityHostsClient.Get(context.Context, string, string, string, string, *ProjectCapabilityHostsClientGetOptions) (ProjectCapabilityHostsClientGetResponse, error)`
- New function `NewProjectConnectionClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectConnectionClient, error)`
- New function `*ProjectConnectionClient.Create(context.Context, string, string, string, string, *ProjectConnectionClientCreateOptions) (ProjectConnectionClientCreateResponse, error)`
- New function `*ProjectConnectionClient.Delete(context.Context, string, string, string, string, *ProjectConnectionClientDeleteOptions) (ProjectConnectionClientDeleteResponse, error)`
- New function `*ProjectConnectionClient.Get(context.Context, string, string, string, string, *ProjectConnectionClientGetOptions) (ProjectConnectionClientGetResponse, error)`
- New function `*ProjectConnectionClient.NewListPager(string, string, string, *ProjectConnectionClientListOptions) *runtime.Pager[ProjectConnectionClientListResponse]`
- New function `*ProjectConnectionClient.Update(context.Context, string, string, string, string, *ProjectConnectionClientUpdateOptions) (ProjectConnectionClientUpdateResponse, error)`
- New function `NewProjectsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ProjectsClient, error)`
- New function `*ProjectsClient.BeginCreate(context.Context, string, string, string, Project, *ProjectsClientBeginCreateOptions) (*runtime.Poller[ProjectsClientCreateResponse], error)`
- New function `*ProjectsClient.BeginDelete(context.Context, string, string, string, *ProjectsClientBeginDeleteOptions) (*runtime.Poller[ProjectsClientDeleteResponse], error)`
- New function `*ProjectsClient.Get(context.Context, string, string, string, *ProjectsClientGetOptions) (ProjectsClientGetResponse, error)`
- New function `*ProjectsClient.NewListPager(string, string, *ProjectsClientListOptions) *runtime.Pager[ProjectsClientListResponse]`
- New function `*ProjectsClient.BeginUpdate(context.Context, string, string, string, Project, *ProjectsClientBeginUpdateOptions) (*runtime.Poller[ProjectsClientUpdateResponse], error)`
- New function `*SASAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*ServicePrincipalAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*UsernamePasswordAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New function `*NoneAuthTypeConnectionProperties.GetConnectionPropertiesV2() *ConnectionPropertiesV2`
- New struct `AADAuthTypeConnectionProperties`
- New struct `APIKeyAuthConnectionProperties`
- New struct `AccessKeyAuthTypeConnectionProperties`
- New struct `AccountKeyAuthTypeConnectionProperties`
- New struct `CapabilityHost`
- New struct `CapabilityHostProperties`
- New struct `ConnectionAPIKey`
- New struct `ConnectionAccessKey`
- New struct `ConnectionAccountKey`
- New struct `ConnectionManagedIdentity`
- New struct `ConnectionOAuth2`
- New struct `ConnectionPersonalAccessToken`
- New struct `ConnectionPropertiesV2BasicResource`
- New struct `ConnectionPropertiesV2BasicResourceArmPaginatedResult`
- New struct `ConnectionServicePrincipal`
- New struct `ConnectionSharedAccessSignature`
- New struct `ConnectionUpdateContent`
- New struct `ConnectionUsernamePassword`
- New struct `CustomKeys`
- New struct `CustomKeysConnectionProperties`
- New struct `ManagedIdentityAuthTypeConnectionProperties`
- New struct `NetworkInjections`
- New struct `NoneAuthTypeConnectionProperties`
- New struct `OAuth2AuthTypeConnectionProperties`
- New struct `PATAuthTypeConnectionProperties`
- New struct `Project`
- New struct `ProjectListResult`
- New struct `ProjectProperties`
- New struct `ResourceBase`
- New struct `SASAuthTypeConnectionProperties`
- New struct `ServicePrincipalAuthTypeConnectionProperties`
- New struct `UsernamePasswordAuthTypeConnectionProperties`
- New field `AllowProjectManagement`, `AssociatedProjects`, `DefaultProject`, `NetworkInjections` in struct `AccountProperties`
- New field `SpilloverDeploymentName` in struct `DeploymentProperties`


## 1.7.0 (2024-12-27)
### Features Added

- New value `ModelLifecycleStatusDeprecated`, `ModelLifecycleStatusDeprecating`, `ModelLifecycleStatusStable` added to enum type `ModelLifecycleStatus`
- New enum type `ByPassSelection` with values `ByPassSelectionAzureServices`, `ByPassSelectionNone`
- New enum type `ContentLevel` with values `ContentLevelHigh`, `ContentLevelLow`, `ContentLevelMedium`
- New enum type `DefenderForAISettingState` with values `DefenderForAISettingStateDisabled`, `DefenderForAISettingStateEnabled`
- New enum type `EncryptionScopeProvisioningState` with values `EncryptionScopeProvisioningStateAccepted`, `EncryptionScopeProvisioningStateCanceled`, `EncryptionScopeProvisioningStateCreating`, `EncryptionScopeProvisioningStateDeleting`, `EncryptionScopeProvisioningStateFailed`, `EncryptionScopeProvisioningStateMoving`, `EncryptionScopeProvisioningStateSucceeded`
- New enum type `EncryptionScopeState` with values `EncryptionScopeStateDisabled`, `EncryptionScopeStateEnabled`
- New enum type `NspAccessRuleDirection` with values `NspAccessRuleDirectionInbound`, `NspAccessRuleDirectionOutbound`
- New enum type `RaiPolicyContentSource` with values `RaiPolicyContentSourceCompletion`, `RaiPolicyContentSourcePrompt`
- New enum type `RaiPolicyMode` with values `RaiPolicyModeAsynchronousFilter`, `RaiPolicyModeBlocking`, `RaiPolicyModeDefault`, `RaiPolicyModeDeferred`
- New enum type `RaiPolicyType` with values `RaiPolicyTypeSystemManaged`, `RaiPolicyTypeUserManaged`
- New function `*ClientFactory.NewDefenderForAISettingsClient() *DefenderForAISettingsClient`
- New function `*ClientFactory.NewEncryptionScopesClient() *EncryptionScopesClient`
- New function `*ClientFactory.NewLocationBasedModelCapacitiesClient() *LocationBasedModelCapacitiesClient`
- New function `*ClientFactory.NewModelCapacitiesClient() *ModelCapacitiesClient`
- New function `*ClientFactory.NewNetworkSecurityPerimeterConfigurationsClient() *NetworkSecurityPerimeterConfigurationsClient`
- New function `*ClientFactory.NewRaiBlocklistItemsClient() *RaiBlocklistItemsClient`
- New function `*ClientFactory.NewRaiBlocklistsClient() *RaiBlocklistsClient`
- New function `*ClientFactory.NewRaiContentFiltersClient() *RaiContentFiltersClient`
- New function `*ClientFactory.NewRaiPoliciesClient() *RaiPoliciesClient`
- New function `NewDefenderForAISettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DefenderForAISettingsClient, error)`
- New function `*DefenderForAISettingsClient.CreateOrUpdate(context.Context, string, string, string, DefenderForAISetting, *DefenderForAISettingsClientCreateOrUpdateOptions) (DefenderForAISettingsClientCreateOrUpdateResponse, error)`
- New function `*DefenderForAISettingsClient.Get(context.Context, string, string, string, *DefenderForAISettingsClientGetOptions) (DefenderForAISettingsClientGetResponse, error)`
- New function `*DefenderForAISettingsClient.NewListPager(string, string, *DefenderForAISettingsClientListOptions) *runtime.Pager[DefenderForAISettingsClientListResponse]`
- New function `*DefenderForAISettingsClient.Update(context.Context, string, string, string, DefenderForAISetting, *DefenderForAISettingsClientUpdateOptions) (DefenderForAISettingsClientUpdateResponse, error)`
- New function `*DeploymentsClient.NewListSKUsPager(string, string, string, *DeploymentsClientListSKUsOptions) *runtime.Pager[DeploymentsClientListSKUsResponse]`
- New function `*DeploymentsClient.BeginUpdate(context.Context, string, string, string, PatchResourceTagsAndSKU, *DeploymentsClientBeginUpdateOptions) (*runtime.Poller[DeploymentsClientUpdateResponse], error)`
- New function `NewEncryptionScopesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*EncryptionScopesClient, error)`
- New function `*EncryptionScopesClient.CreateOrUpdate(context.Context, string, string, string, EncryptionScope, *EncryptionScopesClientCreateOrUpdateOptions) (EncryptionScopesClientCreateOrUpdateResponse, error)`
- New function `*EncryptionScopesClient.BeginDelete(context.Context, string, string, string, *EncryptionScopesClientBeginDeleteOptions) (*runtime.Poller[EncryptionScopesClientDeleteResponse], error)`
- New function `*EncryptionScopesClient.Get(context.Context, string, string, string, *EncryptionScopesClientGetOptions) (EncryptionScopesClientGetResponse, error)`
- New function `*EncryptionScopesClient.NewListPager(string, string, *EncryptionScopesClientListOptions) *runtime.Pager[EncryptionScopesClientListResponse]`
- New function `NewLocationBasedModelCapacitiesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LocationBasedModelCapacitiesClient, error)`
- New function `*LocationBasedModelCapacitiesClient.NewListPager(string, string, string, string, *LocationBasedModelCapacitiesClientListOptions) *runtime.Pager[LocationBasedModelCapacitiesClientListResponse]`
- New function `*ManagementClient.CalculateModelCapacity(context.Context, CalculateModelCapacityParameter, *ManagementClientCalculateModelCapacityOptions) (ManagementClientCalculateModelCapacityResponse, error)`
- New function `NewModelCapacitiesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ModelCapacitiesClient, error)`
- New function `*ModelCapacitiesClient.NewListPager(string, string, string, *ModelCapacitiesClientListOptions) *runtime.Pager[ModelCapacitiesClientListResponse]`
- New function `NewRaiBlocklistItemsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiBlocklistItemsClient, error)`
- New function `*RaiBlocklistItemsClient.BatchAdd(context.Context, string, string, string, []*RaiBlocklistItemBulkRequest, *RaiBlocklistItemsClientBatchAddOptions) (RaiBlocklistItemsClientBatchAddResponse, error)`
- New function `*RaiBlocklistItemsClient.BatchDelete(context.Context, string, string, string, any, *RaiBlocklistItemsClientBatchDeleteOptions) (RaiBlocklistItemsClientBatchDeleteResponse, error)`
- New function `*RaiBlocklistItemsClient.CreateOrUpdate(context.Context, string, string, string, string, RaiBlocklistItem, *RaiBlocklistItemsClientCreateOrUpdateOptions) (RaiBlocklistItemsClientCreateOrUpdateResponse, error)`
- New function `*RaiBlocklistItemsClient.BeginDelete(context.Context, string, string, string, string, *RaiBlocklistItemsClientBeginDeleteOptions) (*runtime.Poller[RaiBlocklistItemsClientDeleteResponse], error)`
- New function `*RaiBlocklistItemsClient.Get(context.Context, string, string, string, string, *RaiBlocklistItemsClientGetOptions) (RaiBlocklistItemsClientGetResponse, error)`
- New function `*RaiBlocklistItemsClient.NewListPager(string, string, string, *RaiBlocklistItemsClientListOptions) *runtime.Pager[RaiBlocklistItemsClientListResponse]`
- New function `NewRaiBlocklistsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiBlocklistsClient, error)`
- New function `*RaiBlocklistsClient.CreateOrUpdate(context.Context, string, string, string, RaiBlocklist, *RaiBlocklistsClientCreateOrUpdateOptions) (RaiBlocklistsClientCreateOrUpdateResponse, error)`
- New function `*RaiBlocklistsClient.BeginDelete(context.Context, string, string, string, *RaiBlocklistsClientBeginDeleteOptions) (*runtime.Poller[RaiBlocklistsClientDeleteResponse], error)`
- New function `*RaiBlocklistsClient.Get(context.Context, string, string, string, *RaiBlocklistsClientGetOptions) (RaiBlocklistsClientGetResponse, error)`
- New function `*RaiBlocklistsClient.NewListPager(string, string, *RaiBlocklistsClientListOptions) *runtime.Pager[RaiBlocklistsClientListResponse]`
- New function `NewRaiContentFiltersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiContentFiltersClient, error)`
- New function `*RaiContentFiltersClient.Get(context.Context, string, string, *RaiContentFiltersClientGetOptions) (RaiContentFiltersClientGetResponse, error)`
- New function `*RaiContentFiltersClient.NewListPager(string, *RaiContentFiltersClientListOptions) *runtime.Pager[RaiContentFiltersClientListResponse]`
- New function `NewRaiPoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RaiPoliciesClient, error)`
- New function `*RaiPoliciesClient.CreateOrUpdate(context.Context, string, string, string, RaiPolicy, *RaiPoliciesClientCreateOrUpdateOptions) (RaiPoliciesClientCreateOrUpdateResponse, error)`
- New function `*RaiPoliciesClient.BeginDelete(context.Context, string, string, string, *RaiPoliciesClientBeginDeleteOptions) (*runtime.Poller[RaiPoliciesClientDeleteResponse], error)`
- New function `*RaiPoliciesClient.Get(context.Context, string, string, string, *RaiPoliciesClientGetOptions) (RaiPoliciesClientGetResponse, error)`
- New function `*RaiPoliciesClient.NewListPager(string, string, *RaiPoliciesClientListOptions) *runtime.Pager[RaiPoliciesClientListResponse]`
- New function `NewNetworkSecurityPerimeterConfigurationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.Get(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientGetOptions) (NetworkSecurityPerimeterConfigurationsClientGetResponse, error)`
- New function `*NetworkSecurityPerimeterConfigurationsClient.NewListPager(string, string, *NetworkSecurityPerimeterConfigurationsClientListOptions) *runtime.Pager[NetworkSecurityPerimeterConfigurationsClientListResponse]`
- New function `*NetworkSecurityPerimeterConfigurationsClient.BeginReconcile(context.Context, string, string, string, *NetworkSecurityPerimeterConfigurationsClientBeginReconcileOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientReconcileResponse], error)`
- New struct `BillingMeterInfo`
- New struct `CalculateModelCapacityParameter`
- New struct `CalculateModelCapacityResult`
- New struct `CalculateModelCapacityResultEstimatedCapacity`
- New struct `CustomBlocklistConfig`
- New struct `DefenderForAISetting`
- New struct `DefenderForAISettingProperties`
- New struct `DefenderForAISettingResult`
- New struct `DeploymentCapacitySettings`
- New struct `DeploymentSKUListResult`
- New struct `EncryptionScope`
- New struct `EncryptionScopeListResult`
- New struct `EncryptionScopeProperties`
- New struct `ModelCapacityCalculatorWorkload`
- New struct `ModelCapacityCalculatorWorkloadRequestParam`
- New struct `ModelCapacityListResult`
- New struct `ModelCapacityListResultValueItem`
- New struct `ModelSKUCapacityProperties`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterAccessRule`
- New struct `NetworkSecurityPerimeterAccessRuleProperties`
- New struct `NetworkSecurityPerimeterAccessRulePropertiesSubscriptionsItem`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationAssociationInfo`
- New struct `NetworkSecurityPerimeterConfigurationList`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityPerimeterProfileInfo`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `RaiBlockListItemsResult`
- New struct `RaiBlockListResult`
- New struct `RaiBlocklist`
- New struct `RaiBlocklistConfig`
- New struct `RaiBlocklistItem`
- New struct `RaiBlocklistItemBulkRequest`
- New struct `RaiBlocklistItemProperties`
- New struct `RaiBlocklistProperties`
- New struct `RaiContentFilter`
- New struct `RaiContentFilterListResult`
- New struct `RaiContentFilterProperties`
- New struct `RaiMonitorConfig`
- New struct `RaiPolicy`
- New struct `RaiPolicyContentFilter`
- New struct `RaiPolicyListResult`
- New struct `RaiPolicyProperties`
- New struct `SKUResource`
- New struct `UserOwnedAmlWorkspace`
- New field `Publisher`, `SourceAccount` in struct `AccountModel`
- New field `AmlWorkspace`, `RaiMonitorConfig` in struct `AccountProperties`
- New field `AllowedValues` in struct `CapacityConfig`
- New field `Tags` in struct `CommitmentPlanAccountAssociation`
- New field `Tags` in struct `Deployment`
- New field `Publisher`, `SourceAccount` in struct `DeploymentModel`
- New field `CapacitySettings`, `CurrentCapacity`, `DynamicThrottlingEnabled`, `ParentDeploymentName` in struct `DeploymentProperties`
- New field `Description` in struct `Model`
- New field `Cost` in struct `ModelSKU`
- New field `Bypass` in struct `NetworkRuleSet`


## 1.6.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.5.0 (2023-07-28)
### Features Added

- New value `DeploymentProvisioningStateCanceled`, `DeploymentProvisioningStateDisabled` added to enum type `DeploymentProvisioningState`
- New value `HostingModelProvisionedWeb` added to enum type `HostingModel`
- New enum type `AbusePenaltyAction` with values `AbusePenaltyActionBlock`, `AbusePenaltyActionThrottle`
- New enum type `DeploymentModelVersionUpgradeOption` with values `DeploymentModelVersionUpgradeOptionNoAutoUpgrade`, `DeploymentModelVersionUpgradeOptionOnceCurrentVersionExpired`, `DeploymentModelVersionUpgradeOptionOnceNewDefaultVersionAvailable`
- New function `*ClientFactory.NewModelsClient() *ModelsClient`
- New function `*ClientFactory.NewUsagesClient() *UsagesClient`
- New function `NewModelsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ModelsClient, error)`
- New function `*ModelsClient.NewListPager(string, *ModelsClientListOptions) *runtime.Pager[ModelsClientListResponse]`
- New function `NewUsagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UsagesClient, error)`
- New function `*UsagesClient.NewListPager(string, *UsagesClientListOptions) *runtime.Pager[UsagesClientListResponse]`
- New struct `AbusePenalty`
- New struct `CapacityConfig`
- New struct `Model`
- New struct `ModelListResult`
- New struct `ModelSKU`
- New field `IsDefaultVersion`, `SKUs`, `Source` in struct `AccountModel`
- New field `AbusePenalty` in struct `AccountProperties`
- New field `ProvisioningIssues` in struct `CommitmentPlanProperties`
- New field `SKU` in struct `Deployment`
- New field `Source` in struct `DeploymentModel`
- New field `RateLimits`, `VersionUpgradeOption` in struct `DeploymentProperties`
- New field `NextLink` in struct `UsageListResult`


## 1.4.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.4.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.3.0 (2023-02-24)
### Features Added

- New type alias `CommitmentPlanProvisioningState` with values `CommitmentPlanProvisioningStateAccepted`, `CommitmentPlanProvisioningStateCanceled`, `CommitmentPlanProvisioningStateCreating`, `CommitmentPlanProvisioningStateDeleting`, `CommitmentPlanProvisioningStateFailed`, `CommitmentPlanProvisioningStateMoving`, `CommitmentPlanProvisioningStateSucceeded`
- New type alias `ModelLifecycleStatus` with values `ModelLifecycleStatusGenerallyAvailable`, `ModelLifecycleStatusPreview`
- New type alias `RoutingMethods` with values `RoutingMethodsPerformance`, `RoutingMethodsPriority`, `RoutingMethodsWeighted`
- New function `*CommitmentPlansClient.BeginCreateOrUpdateAssociation(context.Context, string, string, string, CommitmentPlanAccountAssociation, *CommitmentPlansClientBeginCreateOrUpdateAssociationOptions) (*runtime.Poller[CommitmentPlansClientCreateOrUpdateAssociationResponse], error)`
- New function `*CommitmentPlansClient.BeginCreateOrUpdatePlan(context.Context, string, string, CommitmentPlan, *CommitmentPlansClientBeginCreateOrUpdatePlanOptions) (*runtime.Poller[CommitmentPlansClientCreateOrUpdatePlanResponse], error)`
- New function `*CommitmentPlansClient.BeginDeleteAssociation(context.Context, string, string, string, *CommitmentPlansClientBeginDeleteAssociationOptions) (*runtime.Poller[CommitmentPlansClientDeleteAssociationResponse], error)`
- New function `*CommitmentPlansClient.BeginDeletePlan(context.Context, string, string, *CommitmentPlansClientBeginDeletePlanOptions) (*runtime.Poller[CommitmentPlansClientDeletePlanResponse], error)`
- New function `*CommitmentPlansClient.GetAssociation(context.Context, string, string, string, *CommitmentPlansClientGetAssociationOptions) (CommitmentPlansClientGetAssociationResponse, error)`
- New function `*CommitmentPlansClient.GetPlan(context.Context, string, string, *CommitmentPlansClientGetPlanOptions) (CommitmentPlansClientGetPlanResponse, error)`
- New function `*CommitmentPlansClient.NewListAssociationsPager(string, string, *CommitmentPlansClientListAssociationsOptions) *runtime.Pager[CommitmentPlansClientListAssociationsResponse]`
- New function `*CommitmentPlansClient.NewListPlansByResourceGroupPager(string, *CommitmentPlansClientListPlansByResourceGroupOptions) *runtime.Pager[CommitmentPlansClientListPlansByResourceGroupResponse]`
- New function `*CommitmentPlansClient.NewListPlansBySubscriptionPager(*CommitmentPlansClientListPlansBySubscriptionOptions) *runtime.Pager[CommitmentPlansClientListPlansBySubscriptionResponse]`
- New function `*CommitmentPlansClient.BeginUpdatePlan(context.Context, string, string, PatchResourceTagsAndSKU, *CommitmentPlansClientBeginUpdatePlanOptions) (*runtime.Poller[CommitmentPlansClientUpdatePlanResponse], error)`
- New struct `CommitmentPlanAccountAssociation`
- New struct `CommitmentPlanAccountAssociationListResult`
- New struct `CommitmentPlanAccountAssociationProperties`
- New struct `CommitmentPlanAssociation`
- New struct `CommitmentPlansClientCreateOrUpdateAssociationResponse`
- New struct `CommitmentPlansClientCreateOrUpdatePlanResponse`
- New struct `CommitmentPlansClientDeleteAssociationResponse`
- New struct `CommitmentPlansClientDeletePlanResponse`
- New struct `CommitmentPlansClientListAssociationsResponse`
- New struct `CommitmentPlansClientListPlansByResourceGroupResponse`
- New struct `CommitmentPlansClientListPlansBySubscriptionResponse`
- New struct `CommitmentPlansClientUpdatePlanResponse`
- New struct `MultiRegionSettings`
- New struct `PatchResourceTags`
- New struct `PatchResourceTagsAndSKU`
- New struct `RegionSetting`
- New field `FinetuneCapabilities` in struct `AccountModel`
- New field `LifecycleStatus` in struct `AccountModel`
- New field `CommitmentPlanAssociations` in struct `AccountProperties`
- New field `Locations` in struct `AccountProperties`
- New field `Kind` in struct `CommitmentPlan`
- New field `Location` in struct `CommitmentPlan`
- New field `SKU` in struct `CommitmentPlan`
- New field `Tags` in struct `CommitmentPlan`
- New field `CommitmentPlanGUID` in struct `CommitmentPlanProperties`
- New field `ProvisioningState` in struct `CommitmentPlanProperties`


## 1.2.0 (2022-10-20)
### Features Added

- New field `CallRateLimit` in struct `DeploymentProperties`
- New field `Capabilities` in struct `DeploymentProperties`
- New field `RaiPolicyName` in struct `DeploymentProperties`
- New field `CallRateLimit` in struct `AccountModel`
- New field `CallRateLimit` in struct `DeploymentModel`


## 1.1.0 (2022-06-09)
### Features Added

- New const `DeploymentScaleTypeStandard`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).