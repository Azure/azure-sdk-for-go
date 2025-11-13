# Release History

## 0.4.1 (2025-11-13)
### Other Changes
Please note, this package has been deprecated. The service backing this library is retired on January 31, 2025. For more details on the Azure Service Azure HDInsight on AKS retirement, please visit: <https://azure.microsoft.com/en-us/updates?id=azure-hdinsight-on-aks-retirement/>.

## 0.4.0 (2024-08-22)
### Breaking Changes

- Field `ClusterIdentity` of struct `KafkaProfile` has been removed

### Features Added

- New value `ClusterAvailableUpgradeTypePatchVersionUpgrade` added to enum type `ClusterAvailableUpgradeType`
- New value `ClusterUpgradeTypePatchVersionUpgrade` added to enum type `ClusterUpgradeType`
- New enum type `Category` with values `CategoryCustom`, `CategoryPredefined`
- New enum type `ClusterPoolUpgradeHistoryType` with values `ClusterPoolUpgradeHistoryTypeAKSPatchUpgrade`, `ClusterPoolUpgradeHistoryTypeNodeOsUpgrade`
- New enum type `ClusterPoolUpgradeHistoryUpgradeResultType` with values `ClusterPoolUpgradeHistoryUpgradeResultTypeFailed`, `ClusterPoolUpgradeHistoryUpgradeResultTypeSucceed`
- New enum type `ClusterUpgradeHistorySeverityType` with values `ClusterUpgradeHistorySeverityTypeCritical`, `ClusterUpgradeHistorySeverityTypeHigh`, `ClusterUpgradeHistorySeverityTypeLow`, `ClusterUpgradeHistorySeverityTypeMedium`
- New enum type `ClusterUpgradeHistoryType` with values `ClusterUpgradeHistoryTypeAKSPatchUpgrade`, `ClusterUpgradeHistoryTypeHotfixUpgrade`, `ClusterUpgradeHistoryTypeHotfixUpgradeRollback`, `ClusterUpgradeHistoryTypePatchVersionUpgrade`, `ClusterUpgradeHistoryTypePatchVersionUpgradeRollback`
- New enum type `ClusterUpgradeHistoryUpgradeResultType` with values `ClusterUpgradeHistoryUpgradeResultTypeFailed`, `ClusterUpgradeHistoryUpgradeResultTypeSucceed`
- New enum type `LibraryManagementAction` with values `LibraryManagementActionInstall`, `LibraryManagementActionUninstall`
- New enum type `ManagedIdentityType` with values `ManagedIdentityTypeCluster`, `ManagedIdentityTypeInternal`, `ManagedIdentityTypeUser`
- New enum type `Status` with values `StatusINSTALLED`, `StatusINSTALLFAILED`, `StatusINSTALLING`, `StatusUNINSTALLFAILED`, `StatusUNINSTALLING`
- New enum type `Type` with values `TypeMaven`, `TypePypi`
- New function `*ClientFactory.NewClusterLibrariesClient() *ClusterLibrariesClient`
- New function `*ClientFactory.NewClusterPoolUpgradeHistoriesClient() *ClusterPoolUpgradeHistoriesClient`
- New function `*ClientFactory.NewClusterUpgradeHistoriesClient() *ClusterUpgradeHistoriesClient`
- New function `*ClusterAksPatchUpgradeHistoryProperties.GetClusterUpgradeHistoryProperties() *ClusterUpgradeHistoryProperties`
- New function `*ClusterAvailableInPlaceUpgradeProperties.GetClusterAvailableInPlaceUpgradeProperties() *ClusterAvailableInPlaceUpgradeProperties`
- New function `*ClusterAvailableInPlaceUpgradeProperties.GetClusterAvailableUpgradeProperties() *ClusterAvailableUpgradeProperties`
- New function `*ClusterAvailableUpgradeHotfixUpgradeProperties.GetClusterAvailableInPlaceUpgradeProperties() *ClusterAvailableInPlaceUpgradeProperties`
- New function `*ClusterAvailableUpgradePatchVersionUpgradeProperties.GetClusterAvailableInPlaceUpgradeProperties() *ClusterAvailableInPlaceUpgradeProperties`
- New function `*ClusterAvailableUpgradePatchVersionUpgradeProperties.GetClusterAvailableUpgradeProperties() *ClusterAvailableUpgradeProperties`
- New function `*ClusterHotfixUpgradeHistoryProperties.GetClusterInPlaceUpgradeHistoryProperties() *ClusterInPlaceUpgradeHistoryProperties`
- New function `*ClusterHotfixUpgradeHistoryProperties.GetClusterUpgradeHistoryProperties() *ClusterUpgradeHistoryProperties`
- New function `*ClusterHotfixUpgradeProperties.GetClusterInPlaceUpgradeProperties() *ClusterInPlaceUpgradeProperties`
- New function `*ClusterHotfixUpgradeRollbackHistoryProperties.GetClusterInPlaceUpgradeHistoryProperties() *ClusterInPlaceUpgradeHistoryProperties`
- New function `*ClusterHotfixUpgradeRollbackHistoryProperties.GetClusterUpgradeHistoryProperties() *ClusterUpgradeHistoryProperties`
- New function `*ClusterInPlaceUpgradeHistoryProperties.GetClusterInPlaceUpgradeHistoryProperties() *ClusterInPlaceUpgradeHistoryProperties`
- New function `*ClusterInPlaceUpgradeHistoryProperties.GetClusterUpgradeHistoryProperties() *ClusterUpgradeHistoryProperties`
- New function `*ClusterInPlaceUpgradeProperties.GetClusterInPlaceUpgradeProperties() *ClusterInPlaceUpgradeProperties`
- New function `*ClusterInPlaceUpgradeProperties.GetClusterUpgradeProperties() *ClusterUpgradeProperties`
- New function `NewClusterLibrariesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ClusterLibrariesClient, error)`
- New function `*ClusterLibrariesClient.NewListPager(string, string, string, Category, *ClusterLibrariesClientListOptions) *runtime.Pager[ClusterLibrariesClientListResponse]`
- New function `*ClusterLibrariesClient.BeginManageLibraries(context.Context, string, string, string, ClusterLibraryManagementOperation, *ClusterLibrariesClientBeginManageLibrariesOptions) (*runtime.Poller[ClusterLibrariesClientManageLibrariesResponse], error)`
- New function `*ClusterLibraryProperties.GetClusterLibraryProperties() *ClusterLibraryProperties`
- New function `*ClusterPatchVersionUpgradeHistoryProperties.GetClusterInPlaceUpgradeHistoryProperties() *ClusterInPlaceUpgradeHistoryProperties`
- New function `*ClusterPatchVersionUpgradeHistoryProperties.GetClusterUpgradeHistoryProperties() *ClusterUpgradeHistoryProperties`
- New function `*ClusterPatchVersionUpgradeProperties.GetClusterInPlaceUpgradeProperties() *ClusterInPlaceUpgradeProperties`
- New function `*ClusterPatchVersionUpgradeProperties.GetClusterUpgradeProperties() *ClusterUpgradeProperties`
- New function `*ClusterPatchVersionUpgradeRollbackHistoryProperties.GetClusterInPlaceUpgradeHistoryProperties() *ClusterInPlaceUpgradeHistoryProperties`
- New function `*ClusterPatchVersionUpgradeRollbackHistoryProperties.GetClusterUpgradeHistoryProperties() *ClusterUpgradeHistoryProperties`
- New function `*ClusterPoolAksPatchUpgradeHistoryProperties.GetClusterPoolUpgradeHistoryProperties() *ClusterPoolUpgradeHistoryProperties`
- New function `*ClusterPoolNodeOsUpgradeHistoryProperties.GetClusterPoolUpgradeHistoryProperties() *ClusterPoolUpgradeHistoryProperties`
- New function `NewClusterPoolUpgradeHistoriesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ClusterPoolUpgradeHistoriesClient, error)`
- New function `*ClusterPoolUpgradeHistoriesClient.NewListPager(string, string, *ClusterPoolUpgradeHistoriesClientListOptions) *runtime.Pager[ClusterPoolUpgradeHistoriesClientListResponse]`
- New function `*ClusterPoolUpgradeHistoryProperties.GetClusterPoolUpgradeHistoryProperties() *ClusterPoolUpgradeHistoryProperties`
- New function `NewClusterUpgradeHistoriesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ClusterUpgradeHistoriesClient, error)`
- New function `*ClusterUpgradeHistoriesClient.NewListPager(string, string, string, *ClusterUpgradeHistoriesClientListOptions) *runtime.Pager[ClusterUpgradeHistoriesClientListResponse]`
- New function `*ClusterUpgradeHistoryProperties.GetClusterUpgradeHistoryProperties() *ClusterUpgradeHistoryProperties`
- New function `*ClustersClient.BeginUpgradeManualRollback(context.Context, string, string, string, ClusterUpgradeRollback, *ClustersClientBeginUpgradeManualRollbackOptions) (*runtime.Poller[ClustersClientUpgradeManualRollbackResponse], error)`
- New function `*MavenLibraryProperties.GetClusterLibraryProperties() *ClusterLibraryProperties`
- New function `*PyPiLibraryProperties.GetClusterLibraryProperties() *ClusterLibraryProperties`
- New struct `ClusterAksPatchUpgradeHistoryProperties`
- New struct `ClusterAvailableUpgradePatchVersionUpgradeProperties`
- New struct `ClusterHotfixUpgradeHistoryProperties`
- New struct `ClusterHotfixUpgradeRollbackHistoryProperties`
- New struct `ClusterLibrary`
- New struct `ClusterLibraryList`
- New struct `ClusterLibraryManagementOperation`
- New struct `ClusterLibraryManagementOperationProperties`
- New struct `ClusterPatchVersionUpgradeHistoryProperties`
- New struct `ClusterPatchVersionUpgradeProperties`
- New struct `ClusterPatchVersionUpgradeRollbackHistoryProperties`
- New struct `ClusterPoolAksPatchUpgradeHistoryProperties`
- New struct `ClusterPoolNodeOsUpgradeHistoryProperties`
- New struct `ClusterPoolUpgradeHistory`
- New struct `ClusterPoolUpgradeHistoryListResult`
- New struct `ClusterUpgradeHistory`
- New struct `ClusterUpgradeHistoryListResult`
- New struct `ClusterUpgradeRollback`
- New struct `ClusterUpgradeRollbackProperties`
- New struct `IPTag`
- New struct `ManagedIdentityProfile`
- New struct `ManagedIdentitySpec`
- New struct `MavenLibraryProperties`
- New struct `PyPiLibraryProperties`
- New field `PublicIPTag` in struct `ClusterPoolResourcePropertiesClusterPoolProfile`
- New field `AvailabilityZones` in struct `ClusterPoolResourcePropertiesComputeProfile`
- New field `ManagedIdentityProfile` in struct `ClusterProfile`
- New field `AvailabilityZones` in struct `ComputeProfile`
- New field `VMSize` in struct `SSHProfile`
- New field `SecretsProfile`, `TrinoProfile` in struct `UpdatableClusterProfile`


## 0.3.0 (2024-04-26)
### Breaking Changes

- Type of `ClusterProfile.KafkaProfile` has been changed from `map[string]any` to `*KafkaProfile`
- Field `ID`, `Location`, `Name`, `SystemData`, `Type` of struct `ClusterPatch` has been removed

### Features Added

- New value `ActionLASTSTATEUPDATE`, `ActionRELAUNCH` added to enum type `Action`
- New enum type `ClusterAvailableUpgradeType` with values `ClusterAvailableUpgradeTypeAKSPatchUpgrade`, `ClusterAvailableUpgradeTypeHotfixUpgrade`
- New enum type `ClusterPoolAvailableUpgradeType` with values `ClusterPoolAvailableUpgradeTypeAKSPatchUpgrade`, `ClusterPoolAvailableUpgradeTypeNodeOsUpgrade`
- New enum type `ClusterPoolUpgradeType` with values `ClusterPoolUpgradeTypeAKSPatchUpgrade`, `ClusterPoolUpgradeTypeNodeOsUpgrade`
- New enum type `ClusterUpgradeType` with values `ClusterUpgradeTypeAKSPatchUpgrade`, `ClusterUpgradeTypeHotfixUpgrade`
- New enum type `CurrentClusterAksVersionStatus` with values `CurrentClusterAksVersionStatusDeprecated`, `CurrentClusterAksVersionStatusSupported`
- New enum type `CurrentClusterPoolAksVersionStatus` with values `CurrentClusterPoolAksVersionStatusDeprecated`, `CurrentClusterPoolAksVersionStatusSupported`
- New enum type `DataDiskType` with values `DataDiskTypePremiumSSDLRS`, `DataDiskTypePremiumSSDV2LRS`, `DataDiskTypePremiumSSDZRS`, `DataDiskTypeStandardHDDLRS`, `DataDiskTypeStandardSSDLRS`, `DataDiskTypeStandardSSDZRS`
- New enum type `DbConnectionAuthenticationMode` with values `DbConnectionAuthenticationModeIdentityAuth`, `DbConnectionAuthenticationModeSQLAuth`
- New enum type `DeploymentMode` with values `DeploymentModeApplication`, `DeploymentModeSession`
- New enum type `MetastoreDbConnectionAuthenticationMode` with values `MetastoreDbConnectionAuthenticationModeIdentityAuth`, `MetastoreDbConnectionAuthenticationModeSQLAuth`
- New enum type `OutboundType` with values `OutboundTypeLoadBalancer`, `OutboundTypeUserDefinedRouting`
- New enum type `RangerUsersyncMode` with values `RangerUsersyncModeAutomatic`, `RangerUsersyncModeStatic`
- New enum type `Severity` with values `SeverityCritical`, `SeverityHigh`, `SeverityLow`, `SeverityMedium`
- New enum type `UpgradeMode` with values `UpgradeModeLASTSTATEUPDATE`, `UpgradeModeSTATELESSUPDATE`, `UpgradeModeUPDATE`
- New function `*ClientFactory.NewClusterAvailableUpgradesClient() *ClusterAvailableUpgradesClient`
- New function `*ClientFactory.NewClusterPoolAvailableUpgradesClient() *ClusterPoolAvailableUpgradesClient`
- New function `*ClusterAKSPatchVersionUpgradeProperties.GetClusterUpgradeProperties() *ClusterUpgradeProperties`
- New function `*ClusterAvailableUpgradeAksPatchUpgradeProperties.GetClusterAvailableUpgradeProperties() *ClusterAvailableUpgradeProperties`
- New function `*ClusterAvailableUpgradeHotfixUpgradeProperties.GetClusterAvailableUpgradeProperties() *ClusterAvailableUpgradeProperties`
- New function `*ClusterAvailableUpgradeProperties.GetClusterAvailableUpgradeProperties() *ClusterAvailableUpgradeProperties`
- New function `NewClusterAvailableUpgradesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ClusterAvailableUpgradesClient, error)`
- New function `*ClusterAvailableUpgradesClient.NewListPager(string, string, string, *ClusterAvailableUpgradesClientListOptions) *runtime.Pager[ClusterAvailableUpgradesClientListResponse]`
- New function `*ClusterHotfixUpgradeProperties.GetClusterUpgradeProperties() *ClusterUpgradeProperties`
- New function `*ClusterPoolAKSPatchVersionUpgradeProperties.GetClusterPoolUpgradeProperties() *ClusterPoolUpgradeProperties`
- New function `*ClusterPoolAvailableUpgradeAksPatchUpgradeProperties.GetClusterPoolAvailableUpgradeProperties() *ClusterPoolAvailableUpgradeProperties`
- New function `*ClusterPoolAvailableUpgradeNodeOsUpgradeProperties.GetClusterPoolAvailableUpgradeProperties() *ClusterPoolAvailableUpgradeProperties`
- New function `*ClusterPoolAvailableUpgradeProperties.GetClusterPoolAvailableUpgradeProperties() *ClusterPoolAvailableUpgradeProperties`
- New function `NewClusterPoolAvailableUpgradesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ClusterPoolAvailableUpgradesClient, error)`
- New function `*ClusterPoolAvailableUpgradesClient.NewListPager(string, string, *ClusterPoolAvailableUpgradesClientListOptions) *runtime.Pager[ClusterPoolAvailableUpgradesClientListResponse]`
- New function `*ClusterPoolNodeOsImageUpdateProperties.GetClusterPoolUpgradeProperties() *ClusterPoolUpgradeProperties`
- New function `*ClusterPoolUpgradeProperties.GetClusterPoolUpgradeProperties() *ClusterPoolUpgradeProperties`
- New function `*ClusterPoolsClient.BeginUpgrade(context.Context, string, string, ClusterPoolUpgrade, *ClusterPoolsClientBeginUpgradeOptions) (*runtime.Poller[ClusterPoolsClientUpgradeResponse], error)`
- New function `*ClusterUpgradeProperties.GetClusterUpgradeProperties() *ClusterUpgradeProperties`
- New function `*ClustersClient.BeginUpgrade(context.Context, string, string, string, ClusterUpgrade, *ClustersClientBeginUpgradeOptions) (*runtime.Poller[ClustersClientUpgradeResponse], error)`
- New struct `ClusterAKSPatchVersionUpgradeProperties`
- New struct `ClusterAccessProfile`
- New struct `ClusterAvailableUpgrade`
- New struct `ClusterAvailableUpgradeAksPatchUpgradeProperties`
- New struct `ClusterAvailableUpgradeHotfixUpgradeProperties`
- New struct `ClusterAvailableUpgradeList`
- New struct `ClusterHotfixUpgradeProperties`
- New struct `ClusterPoolAKSPatchVersionUpgradeProperties`
- New struct `ClusterPoolAvailableUpgrade`
- New struct `ClusterPoolAvailableUpgradeAksPatchUpgradeProperties`
- New struct `ClusterPoolAvailableUpgradeList`
- New struct `ClusterPoolAvailableUpgradeNodeOsUpgradeProperties`
- New struct `ClusterPoolNodeOsImageUpdateProperties`
- New struct `ClusterPoolUpgrade`
- New struct `ClusterRangerPluginProfile`
- New struct `ClusterUpgrade`
- New struct `DiskStorageProfile`
- New struct `FlinkJobProfile`
- New struct `KafkaConnectivityEndpoints`
- New struct `KafkaProfile`
- New struct `RangerAdminSpec`
- New struct `RangerAdminSpecDatabase`
- New struct `RangerAuditSpec`
- New struct `RangerProfile`
- New struct `RangerUsersyncSpec`
- New field `Filter` in struct `ClusterJobsClientListOptions`
- New field `APIServerAuthorizedIPRanges`, `EnablePrivateAPIServer`, `OutboundType` in struct `ClusterPoolResourcePropertiesNetworkProfile`
- New field `ClusterAccessProfile`, `RangerPluginProfile`, `RangerProfile` in struct `ClusterProfile`
- New field `PrivateFqdn` in struct `ConnectivityProfileWeb`
- New field `MetastoreDbConnectionAuthenticationMode` in struct `FlinkHiveCatalogOption`
- New field `RunID` in struct `FlinkJobProperties`
- New field `DeploymentMode`, `JobSpec` in struct `FlinkProfile`
- New field `MetastoreDbConnectionAuthenticationMode` in struct `HiveCatalogOption`
- New field `PrivateSSHEndpoint` in struct `SSHConnectivityEndpoint`
- New field `DbConnectionAuthenticationMode` in struct `SparkMetastoreSpec`
- New field `RangerPluginProfile`, `RangerProfile` in struct `UpdatableClusterProfile`


## 0.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.1.0 (2023-08-25)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hdinsightcontainers/armhdinsightcontainers` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).