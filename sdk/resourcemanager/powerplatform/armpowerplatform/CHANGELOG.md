# Release History

## 0.4.0 (2026-03-24)
### Breaking Changes

- Type of `PropertiesNetworkInjection.VirtualNetworks` has been changed from `*VirtualNetworkPropertiesList` to `[]*VirtualNetworkProperties`
- Operation `*PrivateEndpointConnectionsClient.BeginCreateOrUpdate` has been changed to non-LRO, use `*PrivateEndpointConnectionsClient.CreateOrUpdate` instead.
- Operation `*PrivateEndpointConnectionsClient.BeginDelete` has been changed to non-LRO, use `*PrivateEndpointConnectionsClient.Delete` instead.
- Struct `ErrorAdditionalInfo` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorResponse` has been removed
- Struct `PatchTrackedResource` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `TrackedResource` has been removed
- Struct `VirtualNetworkPropertiesList` has been removed

### Features Added

- New value `EnterprisePolicyKindIdentity` added to enum type `EnterprisePolicyKind`
- New enum type `HealthStatus` with values `HealthStatusHealthy`, `HealthStatusUndetermined`, `HealthStatusUnhealthy`, `HealthStatusWarning`
- New field `SystemID` in struct `AccountProperties`
- New field `NextLink` in struct `PrivateEndpointConnectionListResult`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `NextLink` in struct `PrivateLinkResourceListResult`
- New field `HealthStatus`, `SystemID` in struct `Properties`


## 0.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.2.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.1.0 (2022-06-10)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/powerplatform/armpowerplatform` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).