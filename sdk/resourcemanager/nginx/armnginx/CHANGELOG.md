# Release History

## 3.1.0-beta.3 (2025-10-09)
### Breaking Changes

- Function `*APIKeysClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, *APIKeysClientCreateOrUpdateOptions)` to `(context.Context, string, string, string, DeploymentAPIKeyRequest, *APIKeysClientCreateOrUpdateOptions)`
- Function `*CertificatesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, *CertificatesClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, string, Certificate, *CertificatesClientBeginCreateOrUpdateOptions)`
- Function `*ConfigurationsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, *ConfigurationsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, string, ConfigurationRequest, *ConfigurationsClientBeginCreateOrUpdateOptions)`
- Function `*DeploymentsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, *DeploymentsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, Deployment, *DeploymentsClientBeginCreateOrUpdateOptions)`
- Function `*DeploymentsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, *DeploymentsClientBeginUpdateOptions)` to `(context.Context, string, string, DeploymentUpdateParameters, *DeploymentsClientBeginUpdateOptions)`
- Type of `OperationListResult.Value` has been changed from `[]*OperationResult` to `[]*Operation`
- Struct `OperationResult` has been removed
- Field `Body` of struct `APIKeysClientCreateOrUpdateOptions` has been removed
- Field `Body` of struct `CertificatesClientBeginCreateOrUpdateOptions` has been removed
- Field `Body` of struct `ConfigurationsClientBeginCreateOrUpdateOptions` has been removed
- Field `Body` of struct `DeploymentsClientBeginCreateOrUpdateOptions` has been removed
- Field `Body` of struct `DeploymentsClientBeginUpdateOptions` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `NginxDeploymentWafPolicyApplyingStatusCode` with values `NginxDeploymentWafPolicyApplyingStatusCodeApplying`, `NginxDeploymentWafPolicyApplyingStatusCodeFailed`, `NginxDeploymentWafPolicyApplyingStatusCodeNotApplied`, `NginxDeploymentWafPolicyApplyingStatusCodeRemoving`, `NginxDeploymentWafPolicyApplyingStatusCodeSucceeded`
- New enum type `NginxDeploymentWafPolicyCompilingStatusCode` with values `NginxDeploymentWafPolicyCompilingStatusCodeFailed`, `NginxDeploymentWafPolicyCompilingStatusCodeInProgress`, `NginxDeploymentWafPolicyCompilingStatusCodeNotStarted`, `NginxDeploymentWafPolicyCompilingStatusCodeSucceeded`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New function `*ClientFactory.NewDefaultWafPolicyClient() *DefaultWafPolicyClient`
- New function `*ClientFactory.NewWafPolicyClient() *WafPolicyClient`
- New function `NewDefaultWafPolicyClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DefaultWafPolicyClient, error)`
- New function `*DefaultWafPolicyClient.List(context.Context, string, string, *DefaultWafPolicyClientListOptions) (DefaultWafPolicyClientListResponse, error)`
- New function `NewWafPolicyClient(string, azcore.TokenCredential, *arm.ClientOptions) (*WafPolicyClient, error)`
- New function `*WafPolicyClient.BeginCreate(context.Context, string, string, string, DeploymentWafPolicy, *WafPolicyClientBeginCreateOptions) (*runtime.Poller[WafPolicyClientCreateResponse], error)`
- New function `*WafPolicyClient.BeginDelete(context.Context, string, string, string, *WafPolicyClientBeginDeleteOptions) (*runtime.Poller[WafPolicyClientDeleteResponse], error)`
- New function `*WafPolicyClient.Get(context.Context, string, string, string, *WafPolicyClientGetOptions) (WafPolicyClientGetResponse, error)`
- New function `*WafPolicyClient.NewListPager(string, string, *WafPolicyClientListOptions) *runtime.Pager[WafPolicyClientListResponse]`
- New struct `DeploymentDefaultWafPolicyListResponse`
- New struct `DeploymentDefaultWafPolicyProperties`
- New struct `DeploymentWafPolicy`
- New struct `DeploymentWafPolicyApplyingStatus`
- New struct `DeploymentWafPolicyCompilingStatus`
- New struct `DeploymentWafPolicyListResponse`
- New struct `DeploymentWafPolicyMetadata`
- New struct `DeploymentWafPolicyMetadataProperties`
- New struct `DeploymentWafPolicyProperties`
- New struct `Operation`
- New field `SystemData` in struct `DeploymentAPIKeyRequest`
- New field `SystemData` in struct `DeploymentAPIKeyResponse`
- New field `WafRelease` in struct `WebApplicationFirewallStatus`


## 3.1.0-beta.2 (2025-02-27)
### Breaking Changes

- Type of `AnalysisCreateConfig.ProtectedFiles` has been changed from `[]*ConfigurationFile` to `[]*ConfigurationProtectedFileRequest`
- Type of `ConfigurationListResponse.Value` has been changed from `[]*Configuration` to `[]*ConfigurationResponse`
- Type of `ConfigurationsClientBeginCreateOrUpdateOptions.Body` has been changed from `*Configuration` to `*ConfigurationRequest`
- Struct `Configuration` has been removed
- Struct `ConfigurationProperties` has been removed
- Field `Configuration` of struct `ConfigurationsClientCreateOrUpdateResponse` has been removed
- Field `Configuration` of struct `ConfigurationsClientGetResponse` has been removed
- Field `ManagedResourceGroup` of struct `DeploymentProperties` has been removed

### Features Added

- New enum type `ActivationState` with values `ActivationStateDisabled`, `ActivationStateEnabled`
- New enum type `Level` with values `LevelInfo`, `LevelWarning`
- New function `NewAPIKeysClient(string, azcore.TokenCredential, *arm.ClientOptions) (*APIKeysClient, error)`
- New function `*APIKeysClient.CreateOrUpdate(context.Context, string, string, string, *APIKeysClientCreateOrUpdateOptions) (APIKeysClientCreateOrUpdateResponse, error)`
- New function `*APIKeysClient.Delete(context.Context, string, string, string, *APIKeysClientDeleteOptions) (APIKeysClientDeleteResponse, error)`
- New function `*APIKeysClient.Get(context.Context, string, string, string, *APIKeysClientGetOptions) (APIKeysClientGetResponse, error)`
- New function `*APIKeysClient.NewListPager(string, string, *APIKeysClientListOptions) *runtime.Pager[APIKeysClientListResponse]`
- New function `*ClientFactory.NewAPIKeysClient() *APIKeysClient`
- New struct `ConfigurationProtectedFileRequest`
- New struct `ConfigurationProtectedFileResponse`
- New struct `ConfigurationRequest`
- New struct `ConfigurationRequestProperties`
- New struct `ConfigurationResponse`
- New struct `ConfigurationResponseProperties`
- New struct `DeploymentAPIKeyListResponse`
- New struct `DeploymentAPIKeyRequest`
- New struct `DeploymentAPIKeyRequestProperties`
- New struct `DeploymentAPIKeyResponse`
- New struct `DeploymentAPIKeyResponseProperties`
- New struct `DeploymentPropertiesNginxAppProtect`
- New struct `DeploymentUpdatePropertiesNginxAppProtect`
- New struct `DiagnosticItem`
- New struct `WebApplicationFirewallComponentVersions`
- New struct `WebApplicationFirewallPackage`
- New struct `WebApplicationFirewallSettings`
- New struct `WebApplicationFirewallStatus`
- New field `Diagnostics` in struct `AnalysisResultData`
- New anonymous field `ConfigurationResponse` in struct `ConfigurationsClientCreateOrUpdateResponse`
- New anonymous field `ConfigurationResponse` in struct `ConfigurationsClientGetResponse`
- New field `DataplaneAPIEndpoint`, `NginxAppProtect` in struct `DeploymentProperties`
- New field `NetworkProfile`, `NginxAppProtect` in struct `DeploymentUpdateProperties`


## 3.1.0-beta.1 (2024-03-22)
### Features Added

- New function `*ConfigurationsClient.Analysis(context.Context, string, string, string, *ConfigurationsClientAnalysisOptions) (ConfigurationsClientAnalysisResponse, error)`
- New struct `AnalysisCreate`
- New struct `AnalysisCreateConfig`
- New struct `AnalysisDiagnostic`
- New struct `AnalysisResult`
- New struct `AnalysisResultData`
- New struct `AutoUpgradeProfile`
- New struct `CertificateErrorResponseBody`
- New struct `DeploymentScalingPropertiesAutoScaleSettings`
- New struct `ScaleProfile`
- New struct `ScaleProfileCapacity`
- New field `CertificateError`, `KeyVaultSecretCreated`, `KeyVaultSecretVersion`, `SHA1Thumbprint` in struct `CertificateProperties`
- New field `AutoUpgradeProfile` in struct `DeploymentProperties`
- New field `AutoScaleSettings` in struct `DeploymentScalingProperties`
- New field `AutoUpgradeProfile` in struct `DeploymentUpdateProperties`


## 3.0.0 (2023-11-24)
### Breaking Changes

- Field `Tags` of struct `Certificate` has been removed
- Field `Tags` of struct `Configuration` has been removed

### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New struct `DeploymentScalingProperties`
- New struct `DeploymentUserProfile`
- New field `ProtectedFiles` in struct `ConfigurationPackage`
- New field `ScalingProperties`, `UserProfile` in struct `DeploymentProperties`
- New field `ScalingProperties`, `UserProfile` in struct `DeploymentUpdateProperties`


## 2.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 2.0.0 (2022-10-13)
### Breaking Changes

- Function `*CertificatesClient.BeginCreate` has been renamed to `*CertificatesClient.BeginCreateOrUpdate`
- Function `*DeploymentsClient.BeginCreate` has been renamed to `*DeploymentsClient.BeginCreateOrUpdate`
- Struct `CertificatesClientBeginCreateOptions` has been renamed to `CertificatesClientBeginCreateOrUpdateOptions`
- Struct `CertificatesClientCreateResponse` has been renamed to `CertificatesClientCreateOrUpdateResponse`
- Struct `DeploymentsClientBeginCreateOptions` has been renamed to `DeploymentsClientBeginCreateOrUpdateOptions`
- Struct `DeploymentsClientCreateResponse` has been renamed to `DeploymentsClientCreateOrUpdateResponse`

### Features Added

- New field `ProtectedFiles` in struct `ConfigurationProperties`


## 1.0.1 (2022-10-12)
### Other Changes
- Loosen Go version requirement.

## 1.0.0 (2022-08-19)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/nginx/armnginx` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
