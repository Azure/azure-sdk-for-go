# Release History

## 1.1.1 (2025-10-31)
### Other Changes


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.0.0 (2023-08-25)
### Breaking Changes

- `ProviderKubeVirt`, `ProviderOpenStack` from enum `Provider` has been removed

### Features Added

- New value `StatusETCDSnapshotFailed`, `StatusValidatingETCDHealth`, `StatusValidatingImageDownload`, `StatusValidatingImageUpload`, `StatusValidatingSFSConnectivity` added to enum type `Status`
- New field `ArtifactType` in struct `AppliancesClientListKeysOptions`


## 0.4.0 (2023-04-28)
### Breaking Changes

- Type alias `SSHKeyType` has been removed
- Function `*AppliancesClient.ListClusterCustomerUserCredential` has been removed
- Struct `ApplianceListClusterCustomerUserCredentialResults` has been removed

### Features Added

- New value `StatusImageDeprovisioning`, `StatusImageDownloaded`, `StatusImageDownloading`, `StatusImagePending`, `StatusImageProvisioned`, `StatusImageProvisioning`, `StatusImageUnknown`, `StatusUpgradingKVAIO`, `StatusWaitingForKVAIO` added to enum type `Status`
- New function `*AppliancesClient.GetTelemetryConfig(context.Context, *AppliancesClientGetTelemetryConfigOptions) (AppliancesClientGetTelemetryConfigResponse, error)`
- New function `*AppliancesClient.ListKeys(context.Context, string, string, *AppliancesClientListKeysOptions) (AppliancesClientListKeysResponse, error)`
- New struct `ApplianceGetTelemetryConfigResult`
- New struct `ApplianceListKeysResults`
- New struct `ArtifactProfile`
- New field `Certificate` in struct `SSHKey`
- New field `CreationTimeStamp` in struct `SSHKey`
- New field `ExpirationTimeStamp` in struct `SSHKey`


## 0.3.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.3.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.2.0 (2022-06-28)
### Features Added

- New const `StatusOffline`
- New const `StatusPreUpgrade`
- New const `StatusPostUpgrade`
- New const `StatusUpdatingCAPI`
- New const `StatusUpdatingCloudOperator`
- New const `StatusUpgradeFailed`
- New const `AccessProfileTypeClusterCustomerUser`
- New const `StatusUpdatingCluster`
- New const `StatusUpgradePrerequisitesCompleted`
- New const `StatusConnecting`
- New const `ProviderOpenStack`
- New const `StatusUpgradeComplete`
- New const `ProviderKubeVirt`
- New const `StatusUpgradeClusterExtensionFailedToDelete`
- New const `StatusPreparingForUpgrade`
- New const `SSHKeyTypeSSHCustomerUser`
- New const `StatusWaitingForCloudOperator`
- New const `StatusNone`
- New function `*AppliancesClient.GetUpgradeGraph(context.Context, string, string, string, *AppliancesClientGetUpgradeGraphOptions) (AppliancesClientGetUpgradeGraphResponse, error)`
- New function `PossibleSSHKeyTypeValues() []SSHKeyType`
- New function `*AppliancesClient.ListClusterCustomerUserCredential(context.Context, string, string, *AppliancesClientListClusterCustomerUserCredentialOptions) (AppliancesClientListClusterCustomerUserCredentialResponse, error)`
- New struct `ApplianceListClusterCustomerUserCredentialResults`
- New struct `AppliancesClientGetUpgradeGraphOptions`
- New struct `AppliancesClientGetUpgradeGraphResponse`
- New struct `AppliancesClientListClusterCustomerUserCredentialOptions`
- New struct `AppliancesClientListClusterCustomerUserCredentialResponse`
- New struct `SSHKey`
- New struct `SupportedVersion`
- New struct `SupportedVersionCatalogVersion`
- New struct `SupportedVersionCatalogVersionData`
- New struct `SupportedVersionMetadata`
- New struct `UpgradeGraph`
- New struct `UpgradeGraphProperties`


## 0.1.0 (2022-06-10)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourceconnector/armresourceconnector` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).