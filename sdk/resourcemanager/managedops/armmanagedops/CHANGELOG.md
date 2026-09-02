# Release History

## 0.3.0 (2026-09-02)
### Breaking Changes

- Type of `SKU.Name` has been changed from `*string` to `*SKUName`
- Type of `SKU.Tier` has been changed from `*string` to `*SKUTier`

### Features Added

- New value `ProvisioningStateAccepted` added to enum type `ProvisioningState`
- New enum type `SKUName` with values `SKUNameManagedOps`
- New enum type `SKUTier` with values `SKUTierEssential`
- New struct `ErrorDetails`
- New field `ErrorDetails` in struct `AzureMonitorInformation`
- New field `ErrorDetails` in struct `ChangeTrackingInformation`
- New field `ErrorDetails` in struct `DefenderCspmInformation`
- New field `ErrorDetails` in struct `DefenderForServersInformation`
- New field `ErrorDetails` in struct `GuestConfigurationInformation`
- New field `ErrorDetails` in struct `UpdateManagerInformation`


## 0.2.0 (2026-03-06)
### Breaking Changes

- Type of `AzureMonitorInformation.EnablementStatus` has been changed from `*ChangeTrackingInformationEnablementStatus` to `*EnablementState`
- Type of `ChangeTrackingInformation.EnablementStatus` has been changed from `*ChangeTrackingInformationEnablementStatus` to `*EnablementState`
- Type of `DefenderCspmInformation.EnablementStatus` has been changed from `*ChangeTrackingInformationEnablementStatus` to `*EnablementState`
- Type of `DefenderForServersInformation.EnablementStatus` has been changed from `*ChangeTrackingInformationEnablementStatus` to `*EnablementState`
- Type of `DesiredConfiguration.DefenderCspm` has been changed from `*DesiredConfigurationDefenderForServers` to `*DesiredEnablementState`
- Type of `DesiredConfiguration.DefenderForServers` has been changed from `*DesiredConfigurationDefenderForServers` to `*DesiredEnablementState`
- Type of `DesiredConfigurationUpdate.DefenderCspm` has been changed from `*DesiredConfigurationDefenderForServers` to `*DesiredEnablementState`
- Type of `DesiredConfigurationUpdate.DefenderForServers` has been changed from `*DesiredConfigurationDefenderForServers` to `*DesiredEnablementState`
- Type of `GuestConfigurationInformation.EnablementStatus` has been changed from `*ChangeTrackingInformationEnablementStatus` to `*EnablementState`
- Type of `UpdateManagerInformation.EnablementStatus` has been changed from `*ChangeTrackingInformationEnablementStatus` to `*EnablementState`
- Enum `ChangeTrackingInformationEnablementStatus` has been removed
- Enum `DesiredConfigurationDefenderForServers` has been removed

### Features Added

- New enum type `DesiredEnablementState` with values `DesiredEnablementStateDisable`, `DesiredEnablementStateEnable`
- New enum type `EnablementState` with values `EnablementStateDisabled`, `EnablementStateEnabled`, `EnablementStateFailed`, `EnablementStateInProgress`


## 0.1.0 (2026-02-13)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managedops/armmanagedops` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).