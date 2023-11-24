# Release History

## 1.0.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.

### Other Changes

- Release stable version.


## 0.3.0 (2023-05-26)
### Breaking Changes

- Type of `AssociationProperties.AssociationType` has been changed from `*string` to `*AssociationType`
- Type of `AssociationUpdateProperties.AssociationType` has been changed from `*string` to `*AssociationType`
- Type of `AssociationUpdateProperties.Subnet` has been changed from `*AssociationSubnet` to `*AssociationSubnetUpdate`
- Enum `FrontendIPAddressVersion` has been removed
- Struct `FrontendPropertiesIPAddress` has been removed
- Struct `FrontendUpdateProperties` has been removed
- Field `IPAddressVersion`, `Mode`, `PublicIPAddress` of struct `FrontendProperties` has been removed
- Field `Properties` of struct `FrontendUpdate` has been removed
- Field `Properties` of struct `TrafficControllerUpdate` has been removed

### Features Added

- New enum type `AssociationType` with values `AssociationTypeSubnets`
- New struct `AssociationSubnetUpdate`
- New field `Fqdn` in struct `FrontendProperties`


## 0.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.2.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.1.0 (2023-01-11)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicenetworking/armservicenetworking` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).