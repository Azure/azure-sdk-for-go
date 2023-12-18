# Release History

## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 2.0.0 (2023-03-24)
### Breaking Changes

- Type alias `APIVersionParameter` has been removed
- Function `*AvailableGroundStationsClient.Get` has been removed
- Field `Etag` of struct `Contact` has been removed
- Field `Etag` of struct `ContactProfile` has been removed
- Field `Etag` of struct `Spacecraft` has been removed

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New struct `ContactProfileThirdPartyConfiguration`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New field `ThirdPartyConfigurations` in struct `ContactProfileProperties`
- New field `ThirdPartyConfigurations` in struct `ContactProfilesProperties`
- New field `NextLink` in struct `OperationResult`
- New field `Value` in struct `OperationResult`


## 1.0.0 (2022-05-19)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/orbital/armorbital` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).