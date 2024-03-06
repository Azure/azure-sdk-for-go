# Release History

## 1.3.0-beta.1 (2023-11-30)
### Features Added

- New field `EnableSecureChannel` in struct `StorageAccount`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0-beta.1 (2023-07-28)
### Features Added

- New field `EnableSecureChannel` in struct `StorageAccount`
- Added feature to support selecting use secure channel during creation. The paramter would force to true if the cluster created based on a stroage account that secure transfer enabled, no matter it use 'blob' or 'dfs' type.

## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hdinsight/armhdinsight` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
