# Release History

## 1.2.0 (2023-02-24)
### Features Added

- New type alias `ImmutabilityState` with values `ImmutabilityStateDisabled`, `ImmutabilityStateLocked`, `ImmutabilityStateUnlocked`
- New type alias `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New type alias `VaultSubResourceType` with values `VaultSubResourceTypeAzureBackup`, `VaultSubResourceTypeAzureBackupSecondary`, `VaultSubResourceTypeAzureSiteRecovery`
- New function `*Client.Capabilities(context.Context, string, ResourceCapabilities, *ClientCapabilitiesOptions) (ClientCapabilitiesResponse, error)`
- New struct `CapabilitiesProperties`
- New struct `CapabilitiesResponse`
- New struct `CapabilitiesResponseProperties`
- New struct `DNSZone`
- New struct `DNSZoneResponse`
- New struct `ImmutabilitySettings`
- New struct `ResourceCapabilities`
- New struct `ResourceCapabilitiesBase`
- New struct `SecuritySettings`
- New field `GroupIDs` in struct `PrivateEndpointConnection`
- New field `PublicNetworkAccess` in struct `VaultProperties`
- New field `SecuritySettings` in struct `VaultProperties`


## 1.1.0 (2022-07-22)
### Features Added

- New const `StandardTierStorageRedundancyGeoRedundant`
- New const `StandardTierStorageRedundancyZoneRedundant`
- New const `CrossRegionRestoreEnabled`
- New const `CrossRegionRestoreDisabled`
- New const `StandardTierStorageRedundancyLocallyRedundant`
- New function `PossibleCrossRegionRestoreValues() []CrossRegionRestore`
- New function `PossibleStandardTierStorageRedundancyValues() []StandardTierStorageRedundancy`
- New struct `VaultPropertiesRedundancySettings`
- New field `RedundancySettings` in struct `VaultProperties`
- New field `AADAudience` in struct `ResourceCertificateAndAADDetails`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).