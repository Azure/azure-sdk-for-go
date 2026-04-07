# Release History

## 1.0.0 (2026-04-07)
### Features Added

- New value `AKSIdentityTypeWorkload` added to enum type `AKSIdentityType`
- New enum type `AutoUpgradeMode` with values `AutoUpgradeModeCompatible`, `AutoUpgradeModeNone`, `AutoUpgradeModePatch`
- New struct `AccessDetail`
- New struct `AdditionalDetails`
- New struct `ManagementDetails`
- New field `ManagedBy` in struct `Extension`
- New field `AdditionalDetails`, `AutoUpgradeMode`, `ExtensionState`, `ManagementDetails` in struct `ExtensionProperties`
- New field `ClientID`, `ObjectID`, `ResourceID` in struct `ExtensionPropertiesAksAssignedIdentity`
- New field `AutoUpgradeMode` in struct `PatchExtensionProperties`


## 0.1.0 (2025-05-13)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kubernetesconfiguration/armextensions` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
