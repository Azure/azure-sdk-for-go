# Release History

## 1.0.0 (2022-11-09)
### Breaking Changes

- Function `*AppliancesClient.ListClusterCustomerUserCredential` has been removed
- Struct `ApplianceListClusterCustomerUserCredentialResults` has been removed
- Struct `AppliancesClientListClusterCustomerUserCredentialOptions` has been removed
- Struct `AppliancesClientListClusterCustomerUserCredentialResponse` has been removed

### Features Added

- New const `SSHKeyTypeScopedAccessKey`
- New const `StatusImageDeprovisioning`
- New const `StatusImageProvisioned`
- New const `StatusImageDownloaded`
- New const `StatusImageProvisioning`
- New const `StatusUpgradingKVAIO`
- New const `ArtifactTypeLogsArtifactType`
- New const `SSHKeyTypeLogsKey`
- New const `StatusImagePending`
- New const `StatusWaitingForKVAIO`
- New const `SSHKeyTypeManagementCAKey`
- New const `StatusImageUnknown`
- New const `StatusImageDownloading`
- New type alias `ArtifactType`
- New function `*AppliancesClient.ListKeys(context.Context, string, string, *AppliancesClientListKeysOptions) (AppliancesClientListKeysResponse, error)`
- New function `*AppliancesClient.GetTelemetryConfig(context.Context, *AppliancesClientGetTelemetryConfigOptions) (AppliancesClientGetTelemetryConfigResponse, error)`
- New function `PossibleArtifactTypeValues() []ArtifactType`
- New struct `ApplianceGetTelemetryConfigResult`
- New struct `ApplianceListKeysResults`
- New struct `AppliancesClientGetTelemetryConfigOptions`
- New struct `AppliancesClientGetTelemetryConfigResponse`
- New struct `AppliancesClientListKeysOptions`
- New struct `AppliancesClientListKeysResponse`
- New struct `ArtifactProfile`
- New field `SystemData` in struct `Resource`
- New field `SystemData` in struct `TrackedResource`
- New field `Certificate` in struct `SSHKey`
- New field `CreationTimeStamp` in struct `SSHKey`
- New field `ExpirationTimeStamp` in struct `SSHKey`


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