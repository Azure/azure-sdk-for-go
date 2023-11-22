# Release History

## 1.0.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New field `AcceptLanguage` in struct `BitLockerKeysClientListOptions`
- New field `AcceptLanguage` in struct `JobsClientCreateOptions`
- New field `AcceptLanguage` in struct `JobsClientDeleteOptions`
- New field `AcceptLanguage` in struct `JobsClientGetOptions`
- New field `AcceptLanguage` in struct `JobsClientListByResourceGroupOptions`
- New field `AcceptLanguage` in struct `JobsClientListBySubscriptionOptions`
- New field `AcceptLanguage` in struct `JobsClientUpdateOptions`
- New field `AcceptLanguage` in struct `LocationsClientGetOptions`
- New field `AcceptLanguage` in struct `LocationsClientListOptions`
- New field `AcceptLanguage` in struct `OperationsClientListOptions`


## 0.7.0 (Unreleased)

### Breaking Changes

* The `acceptLanguage` parameter has been moved out of client constructors and into each method's options type.

## 0.6.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.6.0 (2023-04-03)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storageimportexport/armstorageimportexport` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).