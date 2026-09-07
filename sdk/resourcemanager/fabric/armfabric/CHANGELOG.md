# Release History

## 1.1.0-beta.1 (2026-08-10)
### Features Added

- New enum type `CapacityOverageState` with values `CapacityOverageStateDisabled`, `CapacityOverageStateEnabled`
- New function `*CapacitiesClient.NewListUsagesPager(location string, options *CapacitiesClientListUsagesOptions) *runtime.Pager[CapacitiesClientListUsagesResponse]`
- New struct `CapacityOverageProperties`
- New struct `PagedQuota`
- New struct `Quota`
- New struct `QuotaName`
- New field `Overage` in struct `CapacityProperties`
- New field `Overage` in struct `CapacityUpdateProperties`


## 1.0.0 (2024-10-22)
### Other Changes

Release stable version.

## 0.1.0 (2024-09-26)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fabric/armfabric` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).