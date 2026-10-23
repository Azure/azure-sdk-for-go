# Release History

## 0.4.0 (2026-09-15)
### Breaking Changes

- Type of `ConfigTemplateVersionProperties.Configurations` has been changed from `*string` to `map[string]any`
- Type of `SchemaVersionProperties.Value` has been changed from `*string` to `map[string]any`
- Type of `SolutionTemplateVersionProperties.Configurations` has been changed from `*string` to `map[string]any`
- Type of `SolutionVersionProperties.Configuration` has been changed from `*string` to `map[string]any`
- Type of `SolutionVersionProperties.TargetLevelConfiguration` has been changed from `*string` to `map[string]any`

### Features Added

- New value `JobTypePublish`, `JobTypeUninstall` added to enum type `JobType`
- New value `StateNotApplicable` added to enum type `State`
- New enum type `CMStages` with values `CMStagesConfiguration`, `CMStagesDeployment`, `CMStagesExternalValidation`, `CMStagesPublish`, `CMStagesStaging`, `CMStagesUninstallation`, `CMStagesUnstaging`
- New enum type `ConfigTemplateConfigurationState` with values `ConfigTemplateConfigurationStateConfigurationCompleted`, `ConfigTemplateConfigurationStateConfigurationPending`
- New enum type `ConfigurationState` with values `ConfigurationStateConfigurationCompleted`, `ConfigurationStateConfigurationPending`
- New enum type `InternalState` with values `InternalStatePendingValidation`, `InternalStateValidated`, `InternalStateValidatedWithSchema`, `InternalStateValidatedWithoutSchema`
- New enum type `StateCategory` with values `StateCategoryCompleted`, `StateCategoryFailed`, `StateCategoryInProgress`, `StateCategoryNone`, `StateCategoryPending`
- New function `*ClientFactory.NewConfigTemplateMetadatasClient() *ConfigTemplateMetadatasClient`
- New function `*ClientFactory.NewConfigTemplateSchemasClient() *ConfigTemplateSchemasClient`
- New function `*ClientFactory.NewHierarchyConfigurationMetadataVersionsClient() *HierarchyConfigurationMetadataVersionsClient`
- New function `*ClientFactory.NewHierarchyConfigurationMetadatasClient() *HierarchyConfigurationMetadatasClient`
- New function `*ClientFactory.NewSolutionDeploymentsClient() *SolutionDeploymentsClient`
- New function `*ClientFactory.NewSolutionMetadataVersionsClient() *SolutionMetadataVersionsClient`
- New function `*ClientFactory.NewSolutionMetadatasClient() *SolutionMetadatasClient`
- New function `*ClientFactory.NewSolutionSchemasClient() *SolutionSchemasClient`
- New function `NewConfigTemplateMetadatasClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConfigTemplateMetadatasClient, error)`
- New function `*ConfigTemplateMetadatasClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, configTemplateName string, configTemplateMetadataName string, resource ConfigTemplateMetadata, options *ConfigTemplateMetadatasClientBeginCreateOrUpdateOptions) (*runtime.Poller[ConfigTemplateMetadatasClientCreateOrUpdateResponse], error)`
- New function `*ConfigTemplateMetadatasClient.BeginDelete(ctx context.Context, resourceGroupName string, configTemplateName string, configTemplateMetadataName string, options *ConfigTemplateMetadatasClientBeginDeleteOptions) (*runtime.Poller[ConfigTemplateMetadatasClientDeleteResponse], error)`
- New function `*ConfigTemplateMetadatasClient.Get(ctx context.Context, resourceGroupName string, configTemplateName string, configTemplateMetadataName string, options *ConfigTemplateMetadatasClientGetOptions) (ConfigTemplateMetadatasClientGetResponse, error)`
- New function `*ConfigTemplateMetadatasClient.NewListByConfigTemplatePager(resourceGroupName string, configTemplateName string, options *ConfigTemplateMetadatasClientListByConfigTemplateOptions) *runtime.Pager[ConfigTemplateMetadatasClientListByConfigTemplateResponse]`
- New function `*ConfigTemplateMetadatasClient.BeginUpdate(ctx context.Context, resourceGroupName string, configTemplateName string, configTemplateMetadataName string, properties ConfigTemplateMetadataUpdate, options *ConfigTemplateMetadatasClientBeginUpdateOptions) (*runtime.Poller[ConfigTemplateMetadatasClientUpdateResponse], error)`
- New function `NewConfigTemplateSchemasClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConfigTemplateSchemasClient, error)`
- New function `*ConfigTemplateSchemasClient.Get(ctx context.Context, resourceGroupName string, configTemplateName string, configTemplateVersionName string, configTemplateSchemaName string, options *ConfigTemplateSchemasClientGetOptions) (ConfigTemplateSchemasClientGetResponse, error)`
- New function `*ConfigTemplateSchemasClient.NewListByConfigTemplateVersionPager(resourceGroupName string, configTemplateName string, configTemplateVersionName string, options *ConfigTemplateSchemasClientListByConfigTemplateVersionOptions) *runtime.Pager[ConfigTemplateSchemasClientListByConfigTemplateVersionResponse]`
- New function `*ConfigTemplateVersionsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, configTemplateName string, configTemplateVersionName string, resource ConfigTemplateVersion, options *ConfigTemplateVersionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ConfigTemplateVersionsClientCreateOrUpdateResponse], error)`
- New function `*ConfigTemplateVersionsClient.BeginDelete(ctx context.Context, resourceGroupName string, configTemplateName string, configTemplateVersionName string, options *ConfigTemplateVersionsClientBeginDeleteOptions) (*runtime.Poller[ConfigTemplateVersionsClientDeleteResponse], error)`
- New function `*ConfigTemplateVersionsClient.Update(ctx context.Context, resourceGroupName string, configTemplateName string, configTemplateVersionName string, properties ConfigTemplateVersion, options *ConfigTemplateVersionsClientUpdateOptions) (ConfigTemplateVersionsClientUpdateResponse, error)`
- New function `*ConfigTemplatesClient.BeginLinkToHierarchies(ctx context.Context, resourceGroupName string, configTemplateName string, body HierarchySelector, options *ConfigTemplatesClientBeginLinkToHierarchiesOptions) (*runtime.Poller[ConfigTemplatesClientLinkToHierarchiesResponse], error)`
- New function `*ConfigTemplatesClient.BeginUnLinkFromHierarchies(ctx context.Context, resourceGroupName string, configTemplateName string, body HierarchySelector, options *ConfigTemplatesClientBeginUnLinkFromHierarchiesOptions) (*runtime.Poller[ConfigTemplatesClientUnLinkFromHierarchiesResponse], error)`
- New function `NewHierarchyConfigurationMetadataVersionsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*HierarchyConfigurationMetadataVersionsClient, error)`
- New function `*HierarchyConfigurationMetadataVersionsClient.Get(ctx context.Context, resourceURI string, hierarchyConfigurationMetadataName string, hierarchyConfigurationMetadataVersionName string, options *HierarchyConfigurationMetadataVersionsClientGetOptions) (HierarchyConfigurationMetadataVersionsClientGetResponse, error)`
- New function `*HierarchyConfigurationMetadataVersionsClient.NewListByParentPager(resourceURI string, hierarchyConfigurationMetadataName string, options *HierarchyConfigurationMetadataVersionsClientListByParentOptions) *runtime.Pager[HierarchyConfigurationMetadataVersionsClientListByParentResponse]`
- New function `NewHierarchyConfigurationMetadatasClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*HierarchyConfigurationMetadatasClient, error)`
- New function `*HierarchyConfigurationMetadatasClient.Get(ctx context.Context, resourceURI string, hierarchyConfigurationMetadataName string, options *HierarchyConfigurationMetadatasClientGetOptions) (HierarchyConfigurationMetadatasClientGetResponse, error)`
- New function `*HierarchyConfigurationMetadatasClient.NewListByParentPager(resourceURI string, options *HierarchyConfigurationMetadatasClientListByParentOptions) *runtime.Pager[HierarchyConfigurationMetadatasClientListByParentResponse]`
- New function `*PublishJobParameter.GetJobParameterBase() *JobParameterBase`
- New function `*PublishJobStepStatistics.GetJobStepStatisticsBase() *JobStepStatisticsBase`
- New function `*SchemaReferencesClient.BeginCreateOrUpdate(ctx context.Context, resourceURI string, schemaReferenceName string, resource SchemaReference, options *SchemaReferencesClientBeginCreateOrUpdateOptions) (*runtime.Poller[SchemaReferencesClientCreateOrUpdateResponse], error)`
- New function `*SchemaReferencesClient.BeginDelete(ctx context.Context, resourceURI string, schemaReferenceName string, options *SchemaReferencesClientBeginDeleteOptions) (*runtime.Poller[SchemaReferencesClientDeleteResponse], error)`
- New function `*SchemaReferencesClient.Update(ctx context.Context, resourceURI string, schemaReferenceName string, properties SchemaReference, options *SchemaReferencesClientUpdateOptions) (SchemaReferencesClientUpdateResponse, error)`
- New function `NewSolutionDeploymentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SolutionDeploymentsClient, error)`
- New function `*SolutionDeploymentsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, solutionDeploymentName string, resource SolutionDeployment, options *SolutionDeploymentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[SolutionDeploymentsClientCreateOrUpdateResponse], error)`
- New function `*SolutionDeploymentsClient.BeginDelete(ctx context.Context, resourceGroupName string, solutionDeploymentName string, options *SolutionDeploymentsClientBeginDeleteOptions) (*runtime.Poller[SolutionDeploymentsClientDeleteResponse], error)`
- New function `*SolutionDeploymentsClient.Get(ctx context.Context, resourceGroupName string, solutionDeploymentName string, options *SolutionDeploymentsClientGetOptions) (SolutionDeploymentsClientGetResponse, error)`
- New function `*SolutionDeploymentsClient.NewListByResourceGroupPager(resourceGroupName string, options *SolutionDeploymentsClientListByResourceGroupOptions) *runtime.Pager[SolutionDeploymentsClientListByResourceGroupResponse]`
- New function `*SolutionDeploymentsClient.NewListBySubscriptionPager(options *SolutionDeploymentsClientListBySubscriptionOptions) *runtime.Pager[SolutionDeploymentsClientListBySubscriptionResponse]`
- New function `*SolutionDeploymentsClient.Update(ctx context.Context, resourceGroupName string, solutionDeploymentName string, properties SolutionDeploymentUpdate, options *SolutionDeploymentsClientUpdateOptions) (SolutionDeploymentsClientUpdateResponse, error)`
- New function `NewSolutionMetadataVersionsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*SolutionMetadataVersionsClient, error)`
- New function `*SolutionMetadataVersionsClient.Get(ctx context.Context, resourceURI string, solutionMetadataName string, solutionMetadataVersionName string, options *SolutionMetadataVersionsClientGetOptions) (SolutionMetadataVersionsClientGetResponse, error)`
- New function `*SolutionMetadataVersionsClient.NewListByParentPager(resourceURI string, solutionMetadataName string, options *SolutionMetadataVersionsClientListByParentOptions) *runtime.Pager[SolutionMetadataVersionsClientListByParentResponse]`
- New function `NewSolutionMetadatasClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*SolutionMetadatasClient, error)`
- New function `*SolutionMetadatasClient.Get(ctx context.Context, resourceURI string, solutionMetadataName string, options *SolutionMetadatasClientGetOptions) (SolutionMetadatasClientGetResponse, error)`
- New function `*SolutionMetadatasClient.NewListByParentPager(resourceURI string, options *SolutionMetadatasClientListByParentOptions) *runtime.Pager[SolutionMetadatasClientListByParentResponse]`
- New function `NewSolutionSchemasClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SolutionSchemasClient, error)`
- New function `*SolutionSchemasClient.Get(ctx context.Context, resourceGroupName string, solutionTemplateName string, solutionTemplateVersionName string, solutionSchemaName string, options *SolutionSchemasClientGetOptions) (SolutionSchemasClientGetResponse, error)`
- New function `*SolutionSchemasClient.NewListBySolutionTemplateVersionPager(resourceGroupName string, solutionTemplateName string, solutionTemplateVersionName string, options *SolutionSchemasClientListBySolutionTemplateVersionOptions) *runtime.Pager[SolutionSchemasClientListBySolutionTemplateVersionResponse]`
- New function `*SolutionTemplateVersionsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, solutionTemplateName string, solutionTemplateVersionName string, resource SolutionTemplateVersion, options *SolutionTemplateVersionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[SolutionTemplateVersionsClientCreateOrUpdateResponse], error)`
- New function `*SolutionTemplateVersionsClient.BeginDelete(ctx context.Context, resourceGroupName string, solutionTemplateName string, solutionTemplateVersionName string, options *SolutionTemplateVersionsClientBeginDeleteOptions) (*runtime.Poller[SolutionTemplateVersionsClientDeleteResponse], error)`
- New function `*SolutionTemplateVersionsClient.Update(ctx context.Context, resourceGroupName string, solutionTemplateName string, solutionTemplateVersionName string, properties SolutionTemplateVersion, options *SolutionTemplateVersionsClientUpdateOptions) (SolutionTemplateVersionsClientUpdateResponse, error)`
- New function `*SolutionTemplateVersionsClient.BeginBulkReviewSolution(ctx context.Context, resourceGroupName string, solutionTemplateName string, solutionTemplateVersionName string, body BulkReviewSolutionParameter, options *SolutionTemplateVersionsClientBeginBulkReviewSolutionOptions) (*runtime.Poller[SolutionTemplateVersionsClientBulkReviewSolutionResponse], error)`
- New function `*TargetsClient.BeginUnstageSolutionVersion(ctx context.Context, resourceGroupName string, targetName string, body SolutionVersionParameter, options *TargetsClientBeginUnstageSolutionVersionOptions) (*runtime.Poller[TargetsClientUnstageSolutionVersionResponse], error)`
- New function `*UninstallJobParameter.GetJobParameterBase() *JobParameterBase`
- New function `*UninstallJobStepStatistics.GetJobStepStatisticsBase() *JobStepStatisticsBase`
- New struct `AdditionalData`
- New struct `BulkReviewSolutionParameter`
- New struct `BulkReviewTargetDetails`
- New struct `ConfigTemplateMetadata`
- New struct `ConfigTemplateMetadataListResult`
- New struct `ConfigTemplateMetadataProperties`
- New struct `ConfigTemplateMetadataUpdate`
- New struct `ConfigTemplateMetadataUpdateProperties`
- New struct `ConfigTemplateSchema`
- New struct `ConfigTemplateSchemaListResult`
- New struct `ConfigTemplateSchemaProperties`
- New struct `HierarchyConfigurationMetadata`
- New struct `HierarchyConfigurationMetadataListResult`
- New struct `HierarchyConfigurationMetadataProperties`
- New struct `HierarchyConfigurationMetadataVersion`
- New struct `HierarchyConfigurationMetadataVersionListResult`
- New struct `HierarchyConfigurationMetadataVersionProperties`
- New struct `HierarchyMetadata`
- New struct `HierarchySelector`
- New struct `PublishJobParameter`
- New struct `PublishJobStepStatistics`
- New struct `SolutionDeployment`
- New struct `SolutionDeploymentListResult`
- New struct `SolutionDeploymentProperties`
- New struct `SolutionDeploymentUpdate`
- New struct `SolutionDeploymentUpdateProperties`
- New struct `SolutionMetadata`
- New struct `SolutionMetadataListResult`
- New struct `SolutionMetadataProperties`
- New struct `SolutionMetadataVersion`
- New struct `SolutionMetadataVersionListResult`
- New struct `SolutionMetadataVersionProperties`
- New struct `SolutionSchema`
- New struct `SolutionSchemaListResult`
- New struct `SolutionSchemaProperties`
- New struct `SolutionTemplateMetadata`
- New struct `SolutionTemplateMetadataUpdate`
- New struct `StageMap`
- New struct `TargetMetadata`
- New struct `UninstallJobParameter`
- New struct `UninstallJobStepStatistics`
- New field `SolutionConfiguration` in struct `BulkPublishSolutionParameter`
- New field `SolutionConfiguration`, `SolutionDependencies`, `SolutionVersionID` in struct `BulkPublishTargetDetails`
- New field `UniqueIdentifier` in struct `ConfigTemplateProperties`
- New field `UniqueIdentifier` in struct `ContextProperties`
- New field `ProvisioningState` in struct `DiagnosticUpdateProperties`
- New field `DisplayName` in struct `DynamicSchemaProperties`
- New field `AdditionalData` in struct `JobProperties`
- New field `CurrentVersion`, `ProvisioningState` in struct `SchemaUpdateProperties`
- New field `DisplayName` in struct `SolutionProperties`
- New field `UniqueIdentifier` in struct `SolutionTemplateProperties`
- New field `InternalState` in struct `SolutionTemplateVersionProperties`
- New field `AvailableSolutionTemplateVersions`, `DisplayName`, `ProvisioningState`, `SolutionTemplateID` in struct `SolutionUpdateProperties`
- New field `CurrentStage`, `LatestActionTriggeredBy`, `Stages` in struct `SolutionVersionProperties`


## 0.3.0 (2025-08-28)
### Breaking Changes

- Function `*SolutionTemplatesClient.Update` parameter(s) have been changed from `(context.Context, string, string, SolutionTemplate, *SolutionTemplatesClientUpdateOptions)` to `(context.Context, string, string, SolutionTemplateUpdate, *SolutionTemplatesClientUpdateOptions)`

### Features Added

- New struct `SolutionTemplateUpdate`
- New struct `SolutionTemplateUpdateProperties`


## 0.2.0 (2025-08-27)
### Breaking Changes

- Function `*ConfigTemplatesClient.Update` parameter(s) have been changed from `(context.Context, string, string, ConfigTemplate, *ConfigTemplatesClientUpdateOptions)` to `(context.Context, string, string, ConfigTemplateUpdate, *ConfigTemplatesClientUpdateOptions)`
- Function `*ContextsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, Context, *ContextsClientBeginUpdateOptions)` to `(context.Context, string, string, ContextUpdate, *ContextsClientBeginUpdateOptions)`
- Function `*DiagnosticsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, Diagnostic, *DiagnosticsClientBeginUpdateOptions)` to `(context.Context, string, string, DiagnosticUpdate, *DiagnosticsClientBeginUpdateOptions)`
- Function `*SchemasClient.Update` parameter(s) have been changed from `(context.Context, string, string, Schema, *SchemasClientUpdateOptions)` to `(context.Context, string, string, SchemaUpdate, *SchemasClientUpdateOptions)`
- Function `*SolutionsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, Solution, *SolutionsClientBeginUpdateOptions)` to `(context.Context, string, string, string, SolutionUpdate, *SolutionsClientBeginUpdateOptions)`
- Function `*TargetsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, Target, *TargetsClientBeginUpdateOptions)` to `(context.Context, string, string, TargetUpdate, *TargetsClientBeginUpdateOptions)`

### Features Added

- New struct `ConfigTemplateUpdate`
- New struct `ConfigTemplateUpdateProperties`
- New struct `ContextUpdate`
- New struct `ContextUpdateProperties`
- New struct `DiagnosticUpdate`
- New struct `DiagnosticUpdateProperties`
- New struct `SchemaUpdate`
- New struct `SchemaUpdateProperties`
- New struct `SolutionUpdate`
- New struct `SolutionUpdateProperties`
- New struct `TargetUpdate`
- New struct `TargetUpdateProperties`


## 0.1.0 (2025-08-13)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/workloadorchestration/armworkloadorchestration` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).