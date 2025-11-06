# Release History

## 2.0.0 (2025-11-06)
### Breaking Changes

- Struct `ARMBaseModel` has been removed
- Struct `MoveRequest` has been removed

### Features Added

- New field `KubernetesWorkloadProfile` in struct `DeviceProperties`
- New field `SystemData` in struct `Job`
- New field `IPRange` in struct `LoadBalancerConfig`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databoxedge/armdataboxedge` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).