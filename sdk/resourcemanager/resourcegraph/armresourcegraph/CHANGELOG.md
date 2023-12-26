# Release History

## 0.9.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.8.2 (2023-10-09)

### Other Changes

- Updated to latest `azcore` beta.

## 0.8.1 (2023-07-19)

### Bug Fixes

- Fixed a potential panic in faked paged and long-running operations.

## 0.8.0 (2023-06-13)

### Features Added

- Support for test fakes and OpenTelemetry trace spans.

## 0.7.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.7.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.6.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).