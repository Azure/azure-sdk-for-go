# Release History

## 1.1.0 (2023-09-22)
### Features Added

- New enum type `CredentialHealthStatus` with values `CredentialHealthStatusHealthy`, `CredentialHealthStatusUnhealthy`
- New enum type `CredentialName` with values `CredentialNameCredential1`
- New function `NewCacheRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CacheRulesClient, error)`
- New function `*CacheRulesClient.BeginCreate(context.Context, string, string, string, CacheRule, *CacheRulesClientBeginCreateOptions) (*runtime.Poller[CacheRulesClientCreateResponse], error)`
- New function `*CacheRulesClient.BeginDelete(context.Context, string, string, string, *CacheRulesClientBeginDeleteOptions) (*runtime.Poller[CacheRulesClientDeleteResponse], error)`
- New function `*CacheRulesClient.Get(context.Context, string, string, string, *CacheRulesClientGetOptions) (CacheRulesClientGetResponse, error)`
- New function `*CacheRulesClient.NewListPager(string, string, *CacheRulesClientListOptions) *runtime.Pager[CacheRulesClientListResponse]`
- New function `*CacheRulesClient.BeginUpdate(context.Context, string, string, string, CacheRuleUpdateParameters, *CacheRulesClientBeginUpdateOptions) (*runtime.Poller[CacheRulesClientUpdateResponse], error)`
- New function `*ClientFactory.NewCacheRulesClient() *CacheRulesClient`
- New function `*ClientFactory.NewCredentialSetsClient() *CredentialSetsClient`
- New function `NewCredentialSetsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CredentialSetsClient, error)`
- New function `*CredentialSetsClient.BeginCreate(context.Context, string, string, string, CredentialSet, *CredentialSetsClientBeginCreateOptions) (*runtime.Poller[CredentialSetsClientCreateResponse], error)`
- New function `*CredentialSetsClient.BeginDelete(context.Context, string, string, string, *CredentialSetsClientBeginDeleteOptions) (*runtime.Poller[CredentialSetsClientDeleteResponse], error)`
- New function `*CredentialSetsClient.Get(context.Context, string, string, string, *CredentialSetsClientGetOptions) (CredentialSetsClientGetResponse, error)`
- New function `*CredentialSetsClient.NewListPager(string, string, *CredentialSetsClientListOptions) *runtime.Pager[CredentialSetsClientListResponse]`
- New function `*CredentialSetsClient.BeginUpdate(context.Context, string, string, string, CredentialSetUpdateParameters, *CredentialSetsClientBeginUpdateOptions) (*runtime.Poller[CredentialSetsClientUpdateResponse], error)`
- New struct `AuthCredential`
- New struct `CacheRule`
- New struct `CacheRuleProperties`
- New struct `CacheRuleUpdateParameters`
- New struct `CacheRuleUpdateProperties`
- New struct `CacheRulesListResult`
- New struct `CredentialHealth`
- New struct `CredentialSet`
- New struct `CredentialSetListResult`
- New struct `CredentialSetProperties`
- New struct `CredentialSetUpdateParameters`
- New struct `CredentialSetUpdateProperties`


## 1.0.0 (2023-03-24)
### Breaking Changes

- Type alias `ActivationStatus` has been removed
- Type alias `AuditLogStatus` has been removed
- Type alias `AzureADAuthenticationAsArmPolicyStatus` has been removed
- Type alias `CertificateType` has been removed
- Type alias `ConnectedRegistryMode` has been removed
- Type alias `ConnectionState` has been removed
- Type alias `CredentialHealthStatus` has been removed
- Type alias `CredentialName` has been removed
- Type alias `LogLevel` has been removed
- Type alias `PipelineOptions` has been removed
- Type alias `PipelineRunSourceType` has been removed
- Type alias `PipelineRunTargetType` has been removed
- Type alias `PipelineSourceType` has been removed
- Type alias `TLSStatus` has been removed
- Function `NewCacheRulesClient` has been removed
- Function `*CacheRulesClient.BeginCreate` has been removed
- Function `*CacheRulesClient.BeginDelete` has been removed
- Function `*CacheRulesClient.Get` has been removed
- Function `*CacheRulesClient.NewListPager` has been removed
- Function `*CacheRulesClient.BeginUpdate` has been removed
- Function `NewConnectedRegistriesClient` has been removed
- Function `*ConnectedRegistriesClient.BeginCreate` has been removed
- Function `*ConnectedRegistriesClient.BeginDeactivate` has been removed
- Function `*ConnectedRegistriesClient.BeginDelete` has been removed
- Function `*ConnectedRegistriesClient.Get` has been removed
- Function `*ConnectedRegistriesClient.NewListPager` has been removed
- Function `*ConnectedRegistriesClient.BeginUpdate` has been removed
- Function `NewCredentialSetsClient` has been removed
- Function `*CredentialSetsClient.BeginCreate` has been removed
- Function `*CredentialSetsClient.BeginDelete` has been removed
- Function `*CredentialSetsClient.Get` has been removed
- Function `*CredentialSetsClient.NewListPager` has been removed
- Function `*CredentialSetsClient.BeginUpdate` has been removed
- Function `NewExportPipelinesClient` has been removed
- Function `*ExportPipelinesClient.BeginCreate` has been removed
- Function `*ExportPipelinesClient.BeginDelete` has been removed
- Function `*ExportPipelinesClient.Get` has been removed
- Function `*ExportPipelinesClient.NewListPager` has been removed
- Function `NewImportPipelinesClient` has been removed
- Function `*ImportPipelinesClient.BeginCreate` has been removed
- Function `*ImportPipelinesClient.BeginDelete` has been removed
- Function `*ImportPipelinesClient.Get` has been removed
- Function `*ImportPipelinesClient.NewListPager` has been removed
- Function `NewPipelineRunsClient` has been removed
- Function `*PipelineRunsClient.BeginCreate` has been removed
- Function `*PipelineRunsClient.BeginDelete` has been removed
- Function `*PipelineRunsClient.Get` has been removed
- Function `*PipelineRunsClient.NewListPager` has been removed
- Struct `ActivationProperties` has been removed
- Struct `AuthCredential` has been removed
- Struct `AzureADAuthenticationAsArmPolicy` has been removed
- Struct `CacheRule` has been removed
- Struct `CacheRuleProperties` has been removed
- Struct `CacheRuleUpdateParameters` has been removed
- Struct `CacheRuleUpdateProperties` has been removed
- Struct `CacheRulesClient` has been removed
- Struct `CacheRulesListResult` has been removed
- Struct `ConnectedRegistriesClient` has been removed
- Struct `ConnectedRegistry` has been removed
- Struct `ConnectedRegistryListResult` has been removed
- Struct `ConnectedRegistryProperties` has been removed
- Struct `ConnectedRegistryUpdateParameters` has been removed
- Struct `ConnectedRegistryUpdateProperties` has been removed
- Struct `CredentialHealth` has been removed
- Struct `CredentialSet` has been removed
- Struct `CredentialSetListResult` has been removed
- Struct `CredentialSetProperties` has been removed
- Struct `CredentialSetUpdateParameters` has been removed
- Struct `CredentialSetUpdateProperties` has been removed
- Struct `CredentialSetsClient` has been removed
- Struct `ExportPipeline` has been removed
- Struct `ExportPipelineListResult` has been removed
- Struct `ExportPipelineProperties` has been removed
- Struct `ExportPipelineTargetProperties` has been removed
- Struct `ExportPipelinesClient` has been removed
- Struct `ImportPipeline` has been removed
- Struct `ImportPipelineListResult` has been removed
- Struct `ImportPipelineProperties` has been removed
- Struct `ImportPipelineSourceProperties` has been removed
- Struct `ImportPipelinesClient` has been removed
- Struct `LoggingProperties` has been removed
- Struct `LoginServerProperties` has been removed
- Struct `ParentProperties` has been removed
- Struct `PipelineRun` has been removed
- Struct `PipelineRunListResult` has been removed
- Struct `PipelineRunProperties` has been removed
- Struct `PipelineRunRequest` has been removed
- Struct `PipelineRunResponse` has been removed
- Struct `PipelineRunSourceProperties` has been removed
- Struct `PipelineRunTargetProperties` has been removed
- Struct `PipelineRunsClient` has been removed
- Struct `PipelineSourceTriggerDescriptor` has been removed
- Struct `PipelineSourceTriggerProperties` has been removed
- Struct `PipelineTriggerDescriptor` has been removed
- Struct `PipelineTriggerProperties` has been removed
- Struct `ProgressProperties` has been removed
- Struct `SoftDeletePolicy` has been removed
- Struct `StatusDetailProperties` has been removed
- Struct `SyncProperties` has been removed
- Struct `SyncUpdateProperties` has been removed
- Struct `TLSCertificateProperties` has been removed
- Struct `TLSProperties` has been removed
- Field `AzureADAuthenticationAsArmPolicy` of struct `Policies` has been removed
- Field `SoftDeletePolicy` of struct `Policies` has been removed
- Field `AnonymousPullEnabled` of struct `RegistryProperties` has been removed
- Field `AnonymousPullEnabled` of struct `RegistryPropertiesUpdateParameters` has been removed

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.7.0 (2023-02-24)
### Features Added

- New type alias `CredentialHealthStatus` with values `CredentialHealthStatusHealthy`, `CredentialHealthStatusUnhealthy`
- New type alias `CredentialName` with values `CredentialNameCredential1`
- New function `NewCacheRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CacheRulesClient, error)`
- New function `*CacheRulesClient.BeginCreate(context.Context, string, string, string, CacheRule, *CacheRulesClientBeginCreateOptions) (*runtime.Poller[CacheRulesClientCreateResponse], error)`
- New function `*CacheRulesClient.BeginDelete(context.Context, string, string, string, *CacheRulesClientBeginDeleteOptions) (*runtime.Poller[CacheRulesClientDeleteResponse], error)`
- New function `*CacheRulesClient.Get(context.Context, string, string, string, *CacheRulesClientGetOptions) (CacheRulesClientGetResponse, error)`
- New function `*CacheRulesClient.NewListPager(string, string, *CacheRulesClientListOptions) *runtime.Pager[CacheRulesClientListResponse]`
- New function `*CacheRulesClient.BeginUpdate(context.Context, string, string, string, CacheRuleUpdateParameters, *CacheRulesClientBeginUpdateOptions) (*runtime.Poller[CacheRulesClientUpdateResponse], error)`
- New function `NewCredentialSetsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CredentialSetsClient, error)`
- New function `*CredentialSetsClient.BeginCreate(context.Context, string, string, string, CredentialSet, *CredentialSetsClientBeginCreateOptions) (*runtime.Poller[CredentialSetsClientCreateResponse], error)`
- New function `*CredentialSetsClient.BeginDelete(context.Context, string, string, string, *CredentialSetsClientBeginDeleteOptions) (*runtime.Poller[CredentialSetsClientDeleteResponse], error)`
- New function `*CredentialSetsClient.Get(context.Context, string, string, string, *CredentialSetsClientGetOptions) (CredentialSetsClientGetResponse, error)`
- New function `*CredentialSetsClient.NewListPager(string, string, *CredentialSetsClientListOptions) *runtime.Pager[CredentialSetsClientListResponse]`
- New function `*CredentialSetsClient.BeginUpdate(context.Context, string, string, string, CredentialSetUpdateParameters, *CredentialSetsClientBeginUpdateOptions) (*runtime.Poller[CredentialSetsClientUpdateResponse], error)`
- New struct `AuthCredential`
- New struct `CacheRule`
- New struct `CacheRuleProperties`
- New struct `CacheRuleUpdateParameters`
- New struct `CacheRuleUpdateProperties`
- New struct `CacheRulesClient`
- New struct `CacheRulesClientCreateResponse`
- New struct `CacheRulesClientDeleteResponse`
- New struct `CacheRulesClientListResponse`
- New struct `CacheRulesClientUpdateResponse`
- New struct `CacheRulesListResult`
- New struct `CredentialHealth`
- New struct `CredentialSet`
- New struct `CredentialSetListResult`
- New struct `CredentialSetProperties`
- New struct `CredentialSetUpdateParameters`
- New struct `CredentialSetUpdateProperties`
- New struct `CredentialSetsClient`
- New struct `CredentialSetsClientCreateResponse`
- New struct `CredentialSetsClientDeleteResponse`
- New struct `CredentialSetsClientListResponse`
- New struct `CredentialSetsClientUpdateResponse`


## 0.6.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).