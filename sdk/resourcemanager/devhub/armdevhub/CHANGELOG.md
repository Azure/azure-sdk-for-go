# Release History

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