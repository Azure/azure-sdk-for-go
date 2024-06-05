# Release History

## 2.1.0-beta.1 (2024-05-24)
### Features Added

- New value `SKUNameEnterpriseE5` added to enum type `SKUName`
- New enum type `DeferUpgradeSetting` with values `DeferUpgradeSettingDeferred`, `DeferUpgradeSettingNotDeferred`
- New function `*DatabasesClient.BeginForceLinkToReplicationGroup(context.Context, string, string, string, ForceLinkParameters, *DatabasesClientBeginForceLinkToReplicationGroupOptions) (*runtime.Poller[DatabasesClientForceLinkToReplicationGroupResponse], error)`
- New function `*DatabasesClient.BeginUpgradeDBRedisVersion(context.Context, string, string, string, *DatabasesClientBeginUpgradeDBRedisVersionOptions) (*runtime.Poller[DatabasesClientUpgradeDBRedisVersionResponse], error)`
- New struct `ForceLinkParameters`
- New field `DeferUpgrade`, `RedisVersion` in struct `DatabaseProperties`


## 2.0.0 (2024-02-23)
### Breaking Changes

- Operation `*PrivateEndpointConnectionsClient.Delete` has been changed to LRO, use `*PrivateEndpointConnectionsClient.BeginDelete` instead.

### Features Added

- New value `ResourceStateScaling`, `ResourceStateScalingFailed` added to enum type `ResourceState`
- New enum type `CmkIdentityType` with values `CmkIdentityTypeSystemAssignedIdentity`, `CmkIdentityTypeUserAssignedIdentity`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New function `*DatabasesClient.BeginFlush(context.Context, string, string, string, FlushParameters, *DatabasesClientBeginFlushOptions) (*runtime.Poller[DatabasesClientFlushResponse], error)`
- New struct `ClusterPropertiesEncryption`
- New struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryption`
- New struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity`
- New struct `FlushParameters`
- New struct `ManagedServiceIdentity`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Cluster`
- New field `Encryption` in struct `ClusterProperties`
- New field `Identity` in struct `ClusterUpdate`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-04-07)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).