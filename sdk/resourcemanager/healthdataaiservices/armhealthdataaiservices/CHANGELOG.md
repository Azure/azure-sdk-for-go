# Release History

## 1.1.0 (2026-01-27)
### Features Added

- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierStandard`
- New struct `SKU`
- New struct `SKUUpdate`
- New field `SKU` in struct `DeidService`
- New field `SKU` in struct `DeidUpdate`


## 1.0.0 (2024-11-20)
### Breaking Changes

- `ManagedServiceIdentityTypeSystemAndUserAssigned` from enum `ManagedServiceIdentityType` has been removed

### Features Added

- New value `ManagedServiceIdentityTypeSystemAssignedUserAssigned` added to enum type `ManagedServiceIdentityType`


## 0.1.0 (2024-08-15)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/healthdataaiservices/armhealthdataaiservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).