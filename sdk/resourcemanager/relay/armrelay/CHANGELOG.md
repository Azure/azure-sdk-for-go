# Release History

## 2.0.0-beta.1 (2026-03-19)
### Breaking Changes

- Type of `Operation.Origin` has been changed from `*string` to `*Origin`
- Struct `ErrorAdditionalInfo` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `ResourceNamespacePatch` has been removed
- Struct `TrackedResource` has been removed
- Field `Properties` of struct `Operation` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New field `PublicNetworkAccess`, `TrustedServiceAccessEnabled` in struct `NetworkRuleSetProperties`
- New field `ActionType` in struct `Operation`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `SystemData` in struct `UpdateParameters`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/relay/armrelay` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).