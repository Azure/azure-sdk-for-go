# Release History

## 2.0.0 (2024-07-31)
### Breaking Changes

- Type of `FileSystemResourceProperties.StorageSKU` has been changed from `*StorageSKU` to `*string`
- `ProvisioningStateNotSpecified` from enum `ProvisioningState` has been removed
- Enum `StorageSKU` has been removed
- Field `InitialCapacity` of struct `FileSystemResourceProperties` has been removed
- Field `ClusterLoginURL`, `PrivateIPs` of struct `FileSystemResourceUpdateProperties` has been removed

### Features Added

- New field `TermUnit` in struct `MarketplaceDetails`


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.0.0 (2023-05-26)
### Other Changes

- Release stable version.


## 0.2.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.1.0 (2023-02-15)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/liftrqumulo/armqumulo` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
