# Release History

## 0.7.0 (2026-05-20)
### Breaking Changes

- Type of `GitHubWorkflowProfile.DeploymentProperties` has been changed from `*DeploymentProperties` to `*Deployment`
- `QuickStartTemplateTypeALL` from enum `QuickStartTemplateType` has been removed
- Struct `DeploymentProperties` has been removed
- Field `NumberOfStores` of struct `ScaleProperty` has been removed
- Field `ScaleProperties` of struct `ScaleTemplateRequest` has been removed

### Features Added

- New enum type `ParameterKind` with values `ParameterKindAzureContainerRegistry`, `ParameterKindAzureKeyvaultURI`, `ParameterKindAzureManagedCluster`, `ParameterKindAzureResourceGroup`, `ParameterKindAzureServiceConnection`, `ParameterKindClusterResourceType`, `ParameterKindContainerImageName`, `ParameterKindContainerImageVersion`, `ParameterKindDirPath`, `ParameterKindDockerFileName`, `ParameterKindEnvVarMap`, `ParameterKindFilePath`, `ParameterKindFlag`, `ParameterKindHelmChartOverrides`, `ParameterKindImagePullPolicy`, `ParameterKindIngressHostName`, `ParameterKindKubernetesNamespace`, `ParameterKindKubernetesProbeDelay`, `ParameterKindKubernetesProbeHTTPPath`, `ParameterKindKubernetesProbePeriod`, `ParameterKindKubernetesProbeThreshold`, `ParameterKindKubernetesProbeTimeout`, `ParameterKindKubernetesProbeType`, `ParameterKindKubernetesResourceLimit`, `ParameterKindKubernetesResourceName`, `ParameterKindKubernetesResourceRequest`, `ParameterKindLabel`, `ParameterKindPort`, `ParameterKindReplicaCount`, `ParameterKindRepositoryBranch`, `ParameterKindResourceLimit`, `ParameterKindScalingResourceType`, `ParameterKindScalingResourceUtilization`, `ParameterKindWorkflowAuthType`, `ParameterKindWorkflowName`
- New enum type `ParameterType` with values `ParameterTypeBool`, `ParameterTypeFloat`, `ParameterTypeInt`, `ParameterTypeObject`, `ParameterTypeString`
- New enum type `RepositoryProviderType` with values `RepositoryProviderTypeAdo`, `RepositoryProviderTypeGithub`
- New enum type `TemplateType` with values `TemplateTypeDeployment`, `TemplateTypeDockerfile`, `TemplateTypeManifest`, `TemplateTypeWorkflow`
- New function `NewADOOAuthClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ADOOAuthClient, error)`
- New function `*ADOOAuthClient.Get(ctx context.Context, location string, options *ADOOAuthClientGetOptions) (ADOOAuthClientGetResponse, error)`
- New function `*ADOOAuthClient.NewListPager(location string, options *ADOOAuthClientListOptions) *runtime.Pager[ADOOAuthClientListResponse]`
- New function `*ClientFactory.NewADOOAuthClient() *ADOOAuthClient`
- New function `*ClientFactory.NewTemplateClient() *TemplateClient`
- New function `*ClientFactory.NewVersionedTemplateClient() *VersionedTemplateClient`
- New function `*DeveloperHubServiceClient.GetADOOAuthInfo(ctx context.Context, location string, options *DeveloperHubServiceClientGetADOOAuthInfoOptions) (DeveloperHubServiceClientGetADOOAuthInfoResponse, error)`
- New function `NewTemplateClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TemplateClient, error)`
- New function `*TemplateClient.Get(ctx context.Context, templateName string, options *TemplateClientGetOptions) (TemplateClientGetResponse, error)`
- New function `*TemplateClient.NewListPager(options *TemplateClientListOptions) *runtime.Pager[TemplateClientListResponse]`
- New function `NewVersionedTemplateClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VersionedTemplateClient, error)`
- New function `*VersionedTemplateClient.Generate(ctx context.Context, templateName string, templateVersion string, parameters map[string]*string, options *VersionedTemplateClientGenerateOptions) (VersionedTemplateClientGenerateResponse, error)`
- New function `*VersionedTemplateClient.Get(ctx context.Context, templateName string, templateVersion string, options *VersionedTemplateClientGetOptions) (VersionedTemplateClientGetResponse, error)`
- New function `*VersionedTemplateClient.NewListPager(templateName string, options *VersionedTemplateClientListOptions) *runtime.Pager[VersionedTemplateClientListResponse]`
- New struct `ADOOAuth`
- New struct `ADOOAuthCallRequest`
- New struct `ADOOAuthInfoResponse`
- New struct `ADOOAuthListResponse`
- New struct `ADOOAuthResponse`
- New struct `ADOProviderProfile`
- New struct `ADORepository`
- New struct `AzurePipelineProfile`
- New struct `Build`
- New struct `Deployment`
- New struct `GenerateVersionedTemplateResponse`
- New struct `GitHubProviderProfile`
- New struct `GitHubRepository`
- New struct `OidcCredentials`
- New struct `Parameter`
- New struct `ParameterDefault`
- New struct `PullRequest`
- New struct `Template`
- New struct `TemplateListResult`
- New struct `TemplateProperties`
- New struct `TemplateReference`
- New struct `TemplateWorkflowProfile`
- New struct `VersionedTemplate`
- New struct `VersionedTemplateListResult`
- New struct `VersionedTemplateProperties`
- New field `NumberOfStore` in struct `ScaleProperty`
- New field `ScaleRequirement` in struct `ScaleTemplateRequest`
- New field `AzurePipelineProfile`, `TemplateWorkflowProfile` in struct `WorkflowProperties`


## 0.6.0 (2024-09-26)
### Features Added

- New enum type `QuickStartTemplateType` with values `QuickStartTemplateTypeALL`, `QuickStartTemplateTypeHCI`, `QuickStartTemplateTypeHCIAKS`, `QuickStartTemplateTypeHCIARCVM`, `QuickStartTemplateTypeNone`
- New function `*ClientFactory.NewIacProfilesClient() *IacProfilesClient`
- New function `NewIacProfilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*IacProfilesClient, error)`
- New function `*IacProfilesClient.CreateOrUpdate(context.Context, string, string, IacProfile, *IacProfilesClientCreateOrUpdateOptions) (IacProfilesClientCreateOrUpdateResponse, error)`
- New function `*IacProfilesClient.Delete(context.Context, string, string, *IacProfilesClientDeleteOptions) (IacProfilesClientDeleteResponse, error)`
- New function `*IacProfilesClient.Export(context.Context, string, string, ExportTemplateRequest, *IacProfilesClientExportOptions) (IacProfilesClientExportResponse, error)`
- New function `*IacProfilesClient.Get(context.Context, string, string, *IacProfilesClientGetOptions) (IacProfilesClientGetResponse, error)`
- New function `*IacProfilesClient.NewListByResourceGroupPager(string, *IacProfilesClientListByResourceGroupOptions) *runtime.Pager[IacProfilesClientListByResourceGroupResponse]`
- New function `*IacProfilesClient.NewListPager(*IacProfilesClientListOptions) *runtime.Pager[IacProfilesClientListResponse]`
- New function `*IacProfilesClient.Scale(context.Context, string, string, ScaleTemplateRequest, *IacProfilesClientScaleOptions) (IacProfilesClientScaleResponse, error)`
- New function `*IacProfilesClient.Sync(context.Context, string, string, *IacProfilesClientSyncOptions) (IacProfilesClientSyncResponse, error)`
- New function `*IacProfilesClient.UpdateTags(context.Context, string, string, TagsObject, *IacProfilesClientUpdateTagsOptions) (IacProfilesClientUpdateTagsResponse, error)`
- New struct `ExportTemplateRequest`
- New struct `IacGitHubProfile`
- New struct `IacProfile`
- New struct `IacProfileListResult`
- New struct `IacProfileProperties`
- New struct `IacTemplateDetails`
- New struct `IacTemplateProperties`
- New struct `PrLinkResponse`
- New struct `ScaleProperty`
- New struct `ScaleTemplateRequest`
- New struct `StageProperties`
- New struct `TerraformProfile`


## 0.5.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.4.0 (2023-05-26)
### Breaking Changes

- Type of `GitHubWorkflowProfile.AuthStatus` has been changed from `*ManifestType` to `*AuthorizationStatus`

### Features Added

- New enum type `AuthorizationStatus` with values `AuthorizationStatusAuthorized`, `AuthorizationStatusError`, `AuthorizationStatusNotFound`
- New enum type `DockerfileGenerationMode` with values `DockerfileGenerationModeDisabled`, `DockerfileGenerationModeEnabled`
- New enum type `GenerationLanguage` with values `GenerationLanguageClojure`, `GenerationLanguageCsharp`, `GenerationLanguageErlang`, `GenerationLanguageGo`, `GenerationLanguageGomodule`, `GenerationLanguageGradle`, `GenerationLanguageJava`, `GenerationLanguageJavascript`, `GenerationLanguagePhp`, `GenerationLanguagePython`, `GenerationLanguageRuby`, `GenerationLanguageRust`, `GenerationLanguageSwift`
- New enum type `GenerationManifestType` with values `GenerationManifestTypeHelm`, `GenerationManifestTypeKube`
- New enum type `ManifestGenerationMode` with values `ManifestGenerationModeDisabled`, `ManifestGenerationModeEnabled`
- New enum type `WorkflowRunStatus` with values `WorkflowRunStatusCompleted`, `WorkflowRunStatusInprogress`, `WorkflowRunStatusQueued`
- New function `*DeveloperHubServiceClient.GeneratePreviewArtifacts(context.Context, string, ArtifactGenerationProperties, *DeveloperHubServiceClientGeneratePreviewArtifactsOptions) (DeveloperHubServiceClientGeneratePreviewArtifactsResponse, error)`
- New struct `ArtifactGenerationProperties`
- New field `ArtifactGenerationProperties` in struct `WorkflowProperties`
- New field `WorkflowRunStatus` in struct `WorkflowRun`


## 0.3.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.2.0 (2022-10-13)
### Breaking Changes

- Function `NewWorkflowClient` parameter(s) have been changed from `(string, *string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDeveloperHubServiceClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*DeveloperHubServiceClient.GitHubOAuthCallback` parameter(s) have been changed from `(context.Context, string, *DeveloperHubServiceClientGitHubOAuthCallbackOptions)` to `(context.Context, string, string, string, *DeveloperHubServiceClientGitHubOAuthCallbackOptions)`

### Features Added

- New field `ManagedClusterResource` in struct `WorkflowClientListByResourceGroupOptions`


## 0.1.1 (2022-10-12)
### Other Changes
- Loosen Go version requirement.

## 0.1.0 (2022-09-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devhub/armdevhub` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).