# Release History

## 1.1.0 (2025-07-11)
### Features Added

- New enum type `InstanceFeatureMode` with values `InstanceFeatureModeDisabled`, `InstanceFeatureModePreview`, `InstanceFeatureModeStable`
- New struct `InstanceFeature`
- New field `Features` in struct `InstanceProperties`


## 1.0.0 (2024-12-12)
### Breaking Changes

- `ManagedServiceIdentityTypeSystemAndUserAssigned` from enum `ManagedServiceIdentityType` has been removed

### Features Added

- New value `ManagedServiceIdentityTypeSystemAssignedUserAssigned` added to enum type `ManagedServiceIdentityType`


## 0.1.0 (2024-10-24)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotoperations/armiotoperations` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).