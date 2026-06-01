# Release History

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