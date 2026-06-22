# Release History

## 0.8.0 (2026-05-20)
### Breaking Changes

- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseError` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `TrackedResource` has been removed

### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New struct `SystemData`
- New field `NextLink` in struct `OperationList`
- New field `SystemData` in struct `ProviderInstance`
- New field `SystemData` in struct `SapMonitor`


## 0.7.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.6.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hanaonazure/armhanaonazure` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).