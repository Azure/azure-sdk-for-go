# Release History

## 3.0.0 (2025-10-15)
### Breaking Changes

- Function `*DatabasesClient.BeginFlush` parameter(s) have been changed from `(context.Context, string, string, string, FlushParameters, *DatabasesClientBeginFlushOptions)` to `(context.Context, string, string, string, *DatabasesClientBeginFlushOptions)`
- Type of `Cluster.Properties` has been changed from `*ClusterProperties` to `*ClusterCreateProperties`
- Type of `ClusterUpdate.Properties` has been changed from `*ClusterProperties` to `*ClusterUpdateProperties`
- Type of `Database.Properties` has been changed from `*DatabaseProperties` to `*DatabaseCreateProperties`
- Type of `DatabaseUpdate.Properties` has been changed from `*DatabaseProperties` to `*DatabaseUpdateProperties`
- Struct `ClusterProperties` has been removed
- Struct `ClusterPropertiesEncryption` has been removed
- Struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryption` has been removed
- Struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity` has been removed
- Struct `DatabaseProperties` has been removed
- Struct `DatabasePropertiesGeoReplication` has been removed

### Features Added

- New value `ClusteringPolicyNoCluster` added to enum type `ClusteringPolicy`
- New value `ResourceStateMoving` added to enum type `ResourceState`
- New value `SKUNameBalancedB0`, `SKUNameBalancedB1`, `SKUNameBalancedB10`, `SKUNameBalancedB100`, `SKUNameBalancedB1000`, `SKUNameBalancedB150`, `SKUNameBalancedB20`, `SKUNameBalancedB250`, `SKUNameBalancedB3`, `SKUNameBalancedB350`, `SKUNameBalancedB5`, `SKUNameBalancedB50`, `SKUNameBalancedB500`, `SKUNameBalancedB700`, `SKUNameComputeOptimizedX10`, `SKUNameComputeOptimizedX100`, `SKUNameComputeOptimizedX150`, `SKUNameComputeOptimizedX20`, `SKUNameComputeOptimizedX250`, `SKUNameComputeOptimizedX3`, `SKUNameComputeOptimizedX350`, `SKUNameComputeOptimizedX5`, `SKUNameComputeOptimizedX50`, `SKUNameComputeOptimizedX500`, `SKUNameComputeOptimizedX700`, `SKUNameEnterpriseE1`, `SKUNameEnterpriseE200`, `SKUNameEnterpriseE400`, `SKUNameEnterpriseE5`, `SKUNameFlashOptimizedA1000`, `SKUNameFlashOptimizedA1500`, `SKUNameFlashOptimizedA2000`, `SKUNameFlashOptimizedA250`, `SKUNameFlashOptimizedA4500`, `SKUNameFlashOptimizedA500`, `SKUNameFlashOptimizedA700`, `SKUNameMemoryOptimizedM10`, `SKUNameMemoryOptimizedM100`, `SKUNameMemoryOptimizedM1000`, `SKUNameMemoryOptimizedM150`, `SKUNameMemoryOptimizedM1500`, `SKUNameMemoryOptimizedM20`, `SKUNameMemoryOptimizedM2000`, `SKUNameMemoryOptimizedM250`, `SKUNameMemoryOptimizedM350`, `SKUNameMemoryOptimizedM50`, `SKUNameMemoryOptimizedM500`, `SKUNameMemoryOptimizedM700` added to enum type `SKUName`
- New enum type `AccessKeysAuthentication` with values `AccessKeysAuthenticationDisabled`, `AccessKeysAuthenticationEnabled`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DeferUpgradeSetting` with values `DeferUpgradeSettingDeferred`, `DeferUpgradeSettingNotDeferred`
- New enum type `HighAvailability` with values `HighAvailabilityDisabled`, `HighAvailabilityEnabled`
- New enum type `Kind` with values `KindV1`, `KindV2`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New enum type `RedundancyMode` with values `RedundancyModeLR`, `RedundancyModeNone`, `RedundancyModeZR`
- New function `NewAccessPolicyAssignmentClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessPolicyAssignmentClient, error)`
- New function `*AccessPolicyAssignmentClient.BeginCreateUpdate(context.Context, string, string, string, string, AccessPolicyAssignment, *AccessPolicyAssignmentClientBeginCreateUpdateOptions) (*runtime.Poller[AccessPolicyAssignmentClientCreateUpdateResponse], error)`
- New function `*AccessPolicyAssignmentClient.BeginDelete(context.Context, string, string, string, string, *AccessPolicyAssignmentClientBeginDeleteOptions) (*runtime.Poller[AccessPolicyAssignmentClientDeleteResponse], error)`
- New function `*AccessPolicyAssignmentClient.Get(context.Context, string, string, string, string, *AccessPolicyAssignmentClientGetOptions) (AccessPolicyAssignmentClientGetResponse, error)`
- New function `*AccessPolicyAssignmentClient.NewListPager(string, string, string, *AccessPolicyAssignmentClientListOptions) *runtime.Pager[AccessPolicyAssignmentClientListResponse]`
- New function `*Client.ListSKUsForScaling(context.Context, string, string, *ClientListSKUsForScalingOptions) (ClientListSKUsForScalingResponse, error)`
- New function `*ClientFactory.NewAccessPolicyAssignmentClient() *AccessPolicyAssignmentClient`
- New function `*DatabasesClient.BeginForceLinkToReplicationGroup(context.Context, string, string, string, ForceLinkParameters, *DatabasesClientBeginForceLinkToReplicationGroupOptions) (*runtime.Poller[DatabasesClientForceLinkToReplicationGroupResponse], error)`
- New function `*DatabasesClient.BeginUpgradeDBRedisVersion(context.Context, string, string, string, *DatabasesClientBeginUpgradeDBRedisVersionOptions) (*runtime.Poller[DatabasesClientUpgradeDBRedisVersionResponse], error)`
- New struct `AccessPolicyAssignment`
- New struct `AccessPolicyAssignmentList`
- New struct `AccessPolicyAssignmentProperties`
- New struct `AccessPolicyAssignmentPropertiesUser`
- New struct `ClusterCommonProperties`
- New struct `ClusterCommonPropertiesEncryption`
- New struct `ClusterCommonPropertiesEncryptionCustomerManagedKeyEncryption`
- New struct `ClusterCommonPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity`
- New struct `ClusterCreateProperties`
- New struct `ClusterUpdateProperties`
- New struct `DatabaseCommonProperties`
- New struct `DatabaseCommonPropertiesGeoReplication`
- New struct `DatabaseCreateProperties`
- New struct `DatabaseUpdateProperties`
- New struct `ErrorDetailAutoGenerated`
- New struct `ErrorResponseAutoGenerated`
- New struct `ForceLinkParameters`
- New struct `ForceLinkParametersGeoReplication`
- New struct `ProxyResourceAutoGenerated`
- New struct `ResourceAutoGenerated`
- New struct `SKUDetails`
- New struct `SKUDetailsList`
- New struct `SystemData`
- New field `Kind` in struct `Cluster`
- New field `SystemData` in struct `Database`
- New field `Parameters` in struct `DatabasesClientBeginFlushOptions`
- New field `SystemData` in struct `ProxyResource`


## 2.1.0-beta.3 (2025-04-23)
### Breaking Changes

- Function `*DatabasesClient.BeginFlush` parameter(s) have been changed from `(context.Context, string, string, string, FlushParameters, *DatabasesClientBeginFlushOptions)` to `(context.Context, string, string, string, *DatabasesClientBeginFlushOptions)`
- Field `GroupNickname`, `LinkedDatabases` of struct `ForceLinkParameters` has been removed

### Features Added

- New value `ClusteringPolicyNoCluster` added to enum type `ClusteringPolicy`
- New value `ResourceStateMoving` added to enum type `ResourceState`
- New enum type `Kind` with values `KindV1`, `KindV2`
- New function `*Client.ListSKUsForScaling(context.Context, string, string, *ClientListSKUsForScalingOptions) (ClientListSKUsForScalingResponse, error)`
- New struct `ErrorDetailAutoGenerated`
- New struct `ErrorResponseAutoGenerated`
- New struct `ForceLinkParametersGeoReplication`
- New struct `SKUDetails`
- New struct `SKUDetailsList`
- New field `Kind` in struct `Cluster`
- New field `Parameters` in struct `DatabasesClientBeginFlushOptions`
- New field `GeoReplication` in struct `ForceLinkParameters`


## 2.1.0-beta.2 (2024-09-27)
### Features Added

- New value `SKUNameBalancedB0`, `SKUNameBalancedB1`, `SKUNameBalancedB10`, `SKUNameBalancedB100`, `SKUNameBalancedB1000`, `SKUNameBalancedB150`, `SKUNameBalancedB20`, `SKUNameBalancedB250`, `SKUNameBalancedB3`, `SKUNameBalancedB350`, `SKUNameBalancedB5`, `SKUNameBalancedB50`, `SKUNameBalancedB500`, `SKUNameBalancedB700`, `SKUNameComputeOptimizedX10`, `SKUNameComputeOptimizedX100`, `SKUNameComputeOptimizedX150`, `SKUNameComputeOptimizedX20`, `SKUNameComputeOptimizedX250`, `SKUNameComputeOptimizedX3`, `SKUNameComputeOptimizedX350`, `SKUNameComputeOptimizedX5`, `SKUNameComputeOptimizedX50`, `SKUNameComputeOptimizedX500`, `SKUNameComputeOptimizedX700`, `SKUNameEnterpriseE1`, `SKUNameEnterpriseE200`, `SKUNameEnterpriseE400`, `SKUNameFlashOptimizedA1000`, `SKUNameFlashOptimizedA1500`, `SKUNameFlashOptimizedA2000`, `SKUNameFlashOptimizedA250`, `SKUNameFlashOptimizedA4500`, `SKUNameFlashOptimizedA500`, `SKUNameFlashOptimizedA700`, `SKUNameMemoryOptimizedM10`, `SKUNameMemoryOptimizedM100`, `SKUNameMemoryOptimizedM1000`, `SKUNameMemoryOptimizedM150`, `SKUNameMemoryOptimizedM1500`, `SKUNameMemoryOptimizedM20`, `SKUNameMemoryOptimizedM2000`, `SKUNameMemoryOptimizedM250`, `SKUNameMemoryOptimizedM350`, `SKUNameMemoryOptimizedM50`, `SKUNameMemoryOptimizedM500`, `SKUNameMemoryOptimizedM700` added to enum type `SKUName`
- New enum type `AccessKeysAuthentication` with values `AccessKeysAuthenticationDisabled`, `AccessKeysAuthenticationEnabled`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `HighAvailability` with values `HighAvailabilityDisabled`, `HighAvailabilityEnabled`
- New enum type `RedundancyMode` with values `RedundancyModeLR`, `RedundancyModeNone`, `RedundancyModeZR`
- New function `NewAccessPolicyAssignmentClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessPolicyAssignmentClient, error)`
- New function `*AccessPolicyAssignmentClient.BeginCreateUpdate(context.Context, string, string, string, string, AccessPolicyAssignment, *AccessPolicyAssignmentClientBeginCreateUpdateOptions) (*runtime.Poller[AccessPolicyAssignmentClientCreateUpdateResponse], error)`
- New function `*AccessPolicyAssignmentClient.BeginDelete(context.Context, string, string, string, string, *AccessPolicyAssignmentClientBeginDeleteOptions) (*runtime.Poller[AccessPolicyAssignmentClientDeleteResponse], error)`
- New function `*AccessPolicyAssignmentClient.Get(context.Context, string, string, string, string, *AccessPolicyAssignmentClientGetOptions) (AccessPolicyAssignmentClientGetResponse, error)`
- New function `*AccessPolicyAssignmentClient.NewListPager(string, string, string, *AccessPolicyAssignmentClientListOptions) *runtime.Pager[AccessPolicyAssignmentClientListResponse]`
- New function `*ClientFactory.NewAccessPolicyAssignmentClient() *AccessPolicyAssignmentClient`
- New struct `AccessPolicyAssignment`
- New struct `AccessPolicyAssignmentList`
- New struct `AccessPolicyAssignmentProperties`
- New struct `AccessPolicyAssignmentPropertiesUser`
- New struct `ProxyResourceAutoGenerated`
- New struct `ResourceAutoGenerated`
- New struct `SystemData`
- New field `HighAvailability`, `RedundancyMode` in struct `ClusterProperties`
- New field `SystemData` in struct `Database`
- New field `AccessKeysAuthentication` in struct `DatabaseProperties`
- New field `SystemData` in struct `ProxyResource`


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