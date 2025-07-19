# Release History

## 3.1.0-beta.3 (2025-08-21)
### Breaking Changes

- Function `*APIKeysClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, *APIKeysClientCreateOrUpdateOptions)` to `(context.Context, string, string, string, APIKeyRequest, *APIKeysClientCreateOrUpdateOptions)`
- Function `*CertificatesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, *CertificatesClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, string, Certificate, *CertificatesClientBeginCreateOrUpdateOptions)`
- Function `*DeploymentsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, *DeploymentsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, Deployment, *DeploymentsClientBeginCreateOrUpdateOptions)`
- Type of `AnalysisCreate.Config` has been changed from `*AnalysisCreateConfig` to `*AnalysisConfig`
- Type of `AnalysisDiagnostic.Line` has been changed from `*float32` to `*int64`
- Type of `AnalysisResult.Data` has been changed from `*AnalysisResultData` to `*AnalysisResultContent`
- Type of `CertificateProperties.KeyVaultSecretCreated` has been changed from `*time.Time` to `*string`
- Type of `ConfigurationRequestProperties.ProtectedFiles` has been changed from `[]*ConfigurationProtectedFileRequest` to `[]*ConfigurationProtectedFileContent`
- Type of `ConfigurationResponseProperties.ProtectedFiles` has been changed from `[]*ConfigurationProtectedFileResponse` to `[]*ConfigurationProtectedFileResult`
- Type of `Deployment.Identity` has been changed from `*IdentityProperties` to `*ManagedServiceIdentity`
- Type of `Deployment.SKU` has been changed from `*ResourceSKU` to `*SKU`
- Type of `DeploymentProperties.NginxAppProtect` has been changed from `*DeploymentPropertiesNginxAppProtect` to `*AppProtect`
- Type of `DeploymentScalingProperties.AutoScaleSettings` has been changed from `*DeploymentScalingPropertiesAutoScaleSettings` to `*AutoScaleSettings`
- Type of `DeploymentUpdateProperties.NginxAppProtect` has been changed from `*DeploymentUpdatePropertiesNginxAppProtect` to `*AppProtect`
- Type of `DiagnosticItem.Level` has been changed from `*Level` to `*DiagnosticLevel`
- Type of `DiagnosticItem.Line` has been changed from `*float32` to `*int64`
- Type of `OperationListResult.Value` has been changed from `[]*OperationResult` to `[]*Operation`
- Type of `PrivateIPAddress.PrivateIPAllocationMethod` has been changed from `*NginxPrivateIPAllocationMethod` to `*PrivateIPAllocationMethod`
- Type of `WebApplicationFirewallPackage.RevisionDatetime` has been changed from `*time.Time` to `*string`
- Enum `IdentityType` has been removed
- Enum `Level` has been removed
- Enum `NginxPrivateIPAllocationMethod` has been removed
- Function `*APIKeysClient.NewListPager` has been removed
- Function `*CertificatesClient.NewListPager` has been removed
- Function `*ConfigurationsClient.Analysis` has been removed
- Function `*ConfigurationsClient.NewListPager` has been removed
- Function `*DeploymentsClient.NewListPager` has been removed
- Operation `*APIKeysClient.Delete` has been changed to LRO, use `*APIKeysClient.BeginDelete` instead.
- Operation `*ConfigurationsClient.BeginCreateOrUpdate` has been changed to non-LRO, use `*ConfigurationsClient.CreateOrUpdate` instead.
- Operation `*DeploymentsClient.BeginUpdate` has been changed to non-LRO, use `*DeploymentsClient.Update` instead.
- Struct `AnalysisCreateConfig` has been removed
- Struct `AnalysisResultData` has been removed
- Struct `CertificateListResponse` has been removed
- Struct `ConfigurationListResponse` has been removed
- Struct `ConfigurationProtectedFileRequest` has been removed
- Struct `ConfigurationProtectedFileResponse` has been removed
- Struct `ConfigurationResponse` has been removed
- Struct `DeploymentAPIKeyListResponse` has been removed
- Struct `DeploymentAPIKeyRequest` has been removed
- Struct `DeploymentAPIKeyRequestProperties` has been removed
- Struct `DeploymentAPIKeyResponse` has been removed
- Struct `DeploymentAPIKeyResponseProperties` has been removed
- Struct `DeploymentListResponse` has been removed
- Struct `DeploymentPropertiesNginxAppProtect` has been removed
- Struct `DeploymentScalingPropertiesAutoScaleSettings` has been removed
- Struct `DeploymentUpdateParameters` has been removed
- Struct `DeploymentUpdatePropertiesNginxAppProtect` has been removed
- Struct `IdentityProperties` has been removed
- Struct `OperationResult` has been removed
- Struct `ResourceSKU` has been removed
- Struct `UserIdentityProperties` has been removed
- Field `Body` of struct `APIKeysClientCreateOrUpdateOptions` has been removed
- Field `DeploymentAPIKeyResponse` of struct `APIKeysClientCreateOrUpdateResponse` has been removed
- Field `DeploymentAPIKeyResponse` of struct `APIKeysClientGetResponse` has been removed
- Field `Location` of struct `Certificate` has been removed
- Field `Body` of struct `CertificatesClientBeginCreateOrUpdateOptions` has been removed
- Field `ConfigurationResponse` of struct `ConfigurationsClientGetResponse` has been removed
- Field `Body` of struct `DeploymentsClientBeginCreateOrUpdateOptions` has been removed
- Field `DeploymentListResponse` of struct `DeploymentsClientListByResourceGroupResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `DiagnosticLevel` with values `DiagnosticLevelInfo`, `DiagnosticLevelWarning`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `PrivateIPAllocationMethod` with values `PrivateIPAllocationMethodDynamic`, `PrivateIPAllocationMethodStatic`
- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierPremium`, `SKUTierStandard`
- New function `*APIKeysClient.NewListByDeploymentPager(string, string, *APIKeysClientListByDeploymentOptions) *runtime.Pager[APIKeysClientListByDeploymentResponse]`
- New function `*CertificatesClient.NewListByDeploymentPager(string, string, *CertificatesClientListByDeploymentOptions) *runtime.Pager[CertificatesClientListByDeploymentResponse]`
- New function `*CertificatesClient.Update(context.Context, string, string, string, CertificateUpdate, *CertificatesClientUpdateOptions) (CertificatesClientUpdateResponse, error)`
- New function `*ClientFactory.NewWafPoliciesClient() *WafPoliciesClient`
- New function `*ConfigurationsClient.Analyze(context.Context, string, string, string, AnalysisCreate, *ConfigurationsClientAnalyzeOptions) (ConfigurationsClientAnalyzeResponse, error)`
- New function `*ConfigurationsClient.NewListByDeploymentPager(string, string, *ConfigurationsClientListByDeploymentOptions) *runtime.Pager[ConfigurationsClientListByDeploymentResponse]`
- New function `*ConfigurationsClient.Update(context.Context, string, string, string, ConfigurationUpdate, *ConfigurationsClientUpdateOptions) (ConfigurationsClientUpdateResponse, error)`
- New function `*DeploymentsClient.NewListBySubscriptionPager(*DeploymentsClientListBySubscriptionOptions) *runtime.Pager[DeploymentsClientListBySubscriptionResponse]`
- New function `*DeploymentsClient.ListDefaultWafPolicies(context.Context, string, string, *DeploymentsClientListDefaultWafPoliciesOptions) (DeploymentsClientListDefaultWafPoliciesResponse, error)`
- New function `NewWafPoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*WafPoliciesClient, error)`
- New function `*WafPoliciesClient.BeginCreateOrUpdate(context.Context, string, string, string, WafPolicy, *WafPoliciesClientBeginCreateOrUpdateOptions) (*runtime.Poller[WafPoliciesClientCreateOrUpdateResponse], error)`
- New function `*WafPoliciesClient.BeginDelete(context.Context, string, string, string, *WafPoliciesClientBeginDeleteOptions) (*runtime.Poller[WafPoliciesClientDeleteResponse], error)`
- New function `*WafPoliciesClient.Get(context.Context, string, string, string, *WafPoliciesClientGetOptions) (WafPoliciesClientGetResponse, error)`
- New function `*WafPoliciesClient.NewListByDeploymentPager(string, string, *WafPoliciesClientListByDeploymentOptions) *runtime.Pager[WafPoliciesClientListByDeploymentResponse]`
- New struct `APIKey`
- New struct `APIKeyListResult`
- New struct `APIKeyRequest`
- New struct `APIKeyRequestProperties`
- New struct `APIKeyResponseProperties`
- New struct `AnalysisConfig`
- New struct `AnalysisResultContent`
- New struct `AppProtect`
- New struct `AutoScaleSettings`
- New struct `CertificateListResult`
- New struct `CertificateUpdate`
- New struct `CertificateUpdateProperties`
- New struct `Configuration`
- New struct `ConfigurationListResult`
- New struct `ConfigurationProtectedFileContent`
- New struct `ConfigurationProtectedFileResult`
- New struct `ConfigurationUpdate`
- New struct `ConfigurationUpdateProperties`
- New struct `DeploymentDefaultWafPolicy`
- New struct `DeploymentDefaultWafPolicyListResponse`
- New struct `DeploymentListResult`
- New struct `DeploymentUpdate`
- New struct `DeploymentWafPolicyApplyingStatus`
- New struct `DeploymentWafPolicyCompilingStatus`
- New struct `DeploymentWafPolicyMetadata`
- New struct `DeploymentWafPolicyMetadataListResult`
- New struct `DeploymentWafPolicyMetadataProperties`
- New struct `DeploymentWafPolicyProperties`
- New struct `ManagedServiceIdentity`
- New struct `Operation`
- New struct `SKU`
- New struct `UserAssignedIdentity`
- New struct `WafPolicy`
- New anonymous field `APIKey` in struct `APIKeysClientCreateOrUpdateResponse`
- New field `RetryAfter` in struct `APIKeysClientCreateOrUpdateResponse`
- New anonymous field `APIKey` in struct `APIKeysClientGetResponse`
- New anonymous field `Configuration` in struct `ConfigurationsClientGetResponse`
- New field `DataplaneAPIEndpoint`, `IPAddress`, `NginxVersion` in struct `DeploymentUpdateProperties`
- New anonymous field `DeploymentListResult` in struct `DeploymentsClientListByResourceGroupResponse`
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
