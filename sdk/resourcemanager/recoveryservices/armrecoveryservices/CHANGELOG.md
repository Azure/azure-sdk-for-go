# Release History

## 2.1.0 (2025-06-12)
### Features Added

- New enum type `IdentityType` with values `IdentityTypeSystemAssigned`, `IdentityTypeUserAssigned`
- New enum type `State` with values `StateDisabled`, `StateEnabled`, `StateInvalid`
- New struct `AssociatedIdentity`
- New struct `SourceScanConfiguration`
- New field `SourceScanConfiguration` in struct `SecuritySettings`


## 2.0.0 (2024-05-24)
### Breaking Changes

- Operation `*VaultsClient.Delete` has been changed to LRO, use `*VaultsClient.BeginDelete` instead.

### Features Added

- New value `StandardTierStorageRedundancyInvalid` added to enum type `StandardTierStorageRedundancy`
- New enum type `BCDRSecurityLevel` with values `BCDRSecurityLevelExcellent`, `BCDRSecurityLevelFair`, `BCDRSecurityLevelGood`, `BCDRSecurityLevelPoor`
- New enum type `EnhancedSecurityState` with values `EnhancedSecurityStateAlwaysON`, `EnhancedSecurityStateDisabled`, `EnhancedSecurityStateEnabled`, `EnhancedSecurityStateInvalid`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New field `AlertsForAllFailoverIssues`, `AlertsForAllReplicationIssues` in struct `AzureMonitorAlertSettings`
- New field `EmailNotificationsForSiteRecovery` in struct `ClassicAlertSettings`
- New field `EnhancedSecurityState` in struct `SoftDeleteSettings`
- New field `BcdrSecurityLevel`, `ResourceGuardOperationRequests` in struct `VaultProperties`
- New field `XMSAuthorizationAuxiliary` in struct `VaultsClientBeginCreateOrUpdateOptions`
- New field `XMSAuthorizationAuxiliary` in struct `VaultsClientBeginUpdateOptions`


## 1.6.0 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.5.0 (2023-09-22)
### Features Added

- New enum type `MultiUserAuthorization` with values `MultiUserAuthorizationDisabled`, `MultiUserAuthorizationEnabled`, `MultiUserAuthorizationInvalid`
- New enum type `SecureScoreLevel` with values `SecureScoreLevelAdequate`, `SecureScoreLevelMaximum`, `SecureScoreLevelMinimum`, `SecureScoreLevelNone`
- New enum type `SoftDeleteState` with values `SoftDeleteStateAlwaysON`, `SoftDeleteStateDisabled`, `SoftDeleteStateEnabled`, `SoftDeleteStateInvalid`
- New struct `SoftDeleteSettings`
- New field `MultiUserAuthorization`, `SoftDeleteSettings` in struct `SecuritySettings`
- New field `SecureScore` in struct `VaultProperties`


## 1.4.0 (2023-06-23)
### Features Added

- New enum type `CrossSubscriptionRestoreState` with values `CrossSubscriptionRestoreStateDisabled`, `CrossSubscriptionRestoreStateEnabled`, `CrossSubscriptionRestoreStatePermanentlyDisabled`
- New struct `CrossSubscriptionRestoreSettings`
- New struct `RestoreSettings`
- New field `RestoreSettings` in struct `VaultProperties`


## 1.3.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.3.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New enum type `ImmutabilityState` with values `ImmutabilityStateDisabled`, `ImmutabilityStateLocked`, `ImmutabilityStateUnlocked`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New enum type `VaultSubResourceType` with values `VaultSubResourceTypeAzureBackup`, `VaultSubResourceTypeAzureBackupSecondary`, `VaultSubResourceTypeAzureSiteRecovery`
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