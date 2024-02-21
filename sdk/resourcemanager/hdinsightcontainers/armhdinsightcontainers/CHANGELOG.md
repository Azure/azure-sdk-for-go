# Release History

## 0.3.0 (2024-02-23)
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