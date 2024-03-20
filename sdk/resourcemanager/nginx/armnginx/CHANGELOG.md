# Release History

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
