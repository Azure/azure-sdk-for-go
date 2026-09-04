# Release History

## 2.1.0-beta.1 (2026-08-21)
### Features Added

- New enum type `ClusterProvisioningState` with values `ClusterProvisioningStateCanceled`, `ClusterProvisioningStateCreating`, `ClusterProvisioningStateDeleting`, `ClusterProvisioningStateFailed`, `ClusterProvisioningStateScaling`, `ClusterProvisioningStateSucceeded`
- New enum type `ClusterSKUName` with values `ClusterSKUNameDedicated`
- New enum type `ClusterSKUScaleType` with values `ClusterSKUScaleTypeAutomatic`
- New enum type `ClusterSKUTier` with values `ClusterSKUTierDedicated`
- New enum type `TLSVersion` with values `TLSVersion12`, `TLSVersion13`
- New function `*ClientFactory.NewClustersClient() *ClustersClient`
- New function `NewClustersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClustersClient, error)`
- New function `*ClustersClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, resource Cluster, options *ClustersClientBeginCreateOrUpdateOptions) (*runtime.Poller[ClustersClientCreateOrUpdateResponse], error)`
- New function `*ClustersClient.BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginDeleteOptions) (*runtime.Poller[ClustersClientDeleteResponse], error)`
- New function `*ClustersClient.Get(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientGetOptions) (ClustersClientGetResponse, error)`
- New function `*ClustersClient.ListAvailableClusterRegion(ctx context.Context, options *ClustersClientListAvailableClusterRegionOptions) (ClustersClientListAvailableClusterRegionResponse, error)`
- New function `*ClustersClient.NewListByResourceGroupPager(resourceGroupName string, options *ClustersClientListByResourceGroupOptions) *runtime.Pager[ClustersClientListByResourceGroupResponse]`
- New function `*ClustersClient.NewListBySubscriptionPager(options *ClustersClientListBySubscriptionOptions) *runtime.Pager[ClustersClientListBySubscriptionResponse]`
- New function `*ClustersClient.ListNamespaces(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientListNamespacesOptions) (ClustersClientListNamespacesResponse, error)`
- New function `*ClustersClient.ListSKUs(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientListSKUsOptions) (ClustersClientListSKUsResponse, error)`
- New function `*ClustersClient.Update(ctx context.Context, resourceGroupName string, clusterName string, properties ClusterUpdate, options *ClustersClientUpdateOptions) (ClustersClientUpdateResponse, error)`
- New struct `AvailableRelayClusterRegion`
- New struct `AvailableRelayClustersList`
- New struct `Cluster`
- New struct `ClusterListResult`
- New struct `ClusterProperties`
- New struct `ClusterSKU`
- New struct `ClusterSKUCapacity`
- New struct `ClusterSKUDetails`
- New struct `ClusterSKUInfo`
- New struct `ClusterSKUListResult`
- New struct `ClusterSKUUpdate`
- New struct `ClusterUpdate`
- New struct `NamespaceIDListResult`
- New struct `NamespaceReference`
- New field `MinimumTLSVersion` in struct `NamespaceProperties`


## 2.0.0 (2026-06-24)
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