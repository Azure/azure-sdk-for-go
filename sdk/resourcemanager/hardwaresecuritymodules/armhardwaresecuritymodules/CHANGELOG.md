# Release History

## 2.0.0-beta.2 (2024-11-13)
### Other Changes

- nothing new

## 2.0.0-beta.1 (2024-01-26)
### Breaking Changes

- Type of `DedicatedHsm.SystemData` has been changed from `*SystemData` to `*DedicatedHsmSystemData`
- Type of `DedicatedHsmOperation.IsDataAction` has been changed from `*string` to `*bool`
- Type of `ResourceListResult.Value` has been changed from `[]*Resource` to `[]*DedicatedHsmResource`
- Type of `SystemData.CreatedByType` has been changed from `*IdentityType` to `*CreatedByType`
- Type of `SystemData.LastModifiedByType` has been changed from `*IdentityType` to `*CreatedByType`
- Field `Location`, `SKU`, `Tags`, `Zones` of struct `Resource` has been removed

### Features Added

- New enum type `CloudHsmClusterSKUFamily` with values `CloudHsmClusterSKUFamilyB`
- New enum type `CloudHsmClusterSKUName` with values `CloudHsmClusterSKUNameStandardB1`, `CloudHsmClusterSKUNameStandardB10`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCanceled`, `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateInternalError`, `PrivateEndpointConnectionProvisioningStateSucceeded`, `PrivateEndpointConnectionProvisioningStateUpdating`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateProvisioning`, `ProvisioningStateSucceeded`
- New function `*ClientFactory.NewCloudHsmClusterPrivateEndpointConnectionsClient() *CloudHsmClusterPrivateEndpointConnectionsClient`
- New function `*ClientFactory.NewCloudHsmClusterPrivateLinkResourcesClient() *CloudHsmClusterPrivateLinkResourcesClient`
- New function `*ClientFactory.NewCloudHsmClustersClient() *CloudHsmClustersClient`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `NewCloudHsmClusterPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CloudHsmClusterPrivateEndpointConnectionsClient, error)`
- New function `*CloudHsmClusterPrivateEndpointConnectionsClient.Create(context.Context, string, string, string, PrivateEndpointConnection, *CloudHsmClusterPrivateEndpointConnectionsClientCreateOptions) (CloudHsmClusterPrivateEndpointConnectionsClientCreateResponse, error)`
- New function `*CloudHsmClusterPrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *CloudHsmClusterPrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[CloudHsmClusterPrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*CloudHsmClusterPrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *CloudHsmClusterPrivateEndpointConnectionsClientGetOptions) (CloudHsmClusterPrivateEndpointConnectionsClientGetResponse, error)`
- New function `NewCloudHsmClusterPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CloudHsmClusterPrivateLinkResourcesClient, error)`
- New function `*CloudHsmClusterPrivateLinkResourcesClient.ListByCloudHsmCluster(context.Context, string, string, *CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterOptions) (CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterResponse, error)`
- New function `NewCloudHsmClustersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CloudHsmClustersClient, error)`
- New function `*CloudHsmClustersClient.BeginCreateOrUpdate(context.Context, string, string, CloudHsmCluster, *CloudHsmClustersClientBeginCreateOrUpdateOptions) (*runtime.Poller[CloudHsmClustersClientCreateOrUpdateResponse], error)`
- New function `*CloudHsmClustersClient.BeginDelete(context.Context, string, string, *CloudHsmClustersClientBeginDeleteOptions) (*runtime.Poller[CloudHsmClustersClientDeleteResponse], error)`
- New function `*CloudHsmClustersClient.Get(context.Context, string, string, *CloudHsmClustersClientGetOptions) (CloudHsmClustersClientGetResponse, error)`
- New function `*CloudHsmClustersClient.NewListByResourceGroupPager(string, *CloudHsmClustersClientListByResourceGroupOptions) *runtime.Pager[CloudHsmClustersClientListByResourceGroupResponse]`
- New function `*CloudHsmClustersClient.NewListBySubscriptionPager(*CloudHsmClustersClientListBySubscriptionOptions) *runtime.Pager[CloudHsmClustersClientListBySubscriptionResponse]`
- New function `*CloudHsmClustersClient.BeginUpdate(context.Context, string, string, CloudHsmClusterPatchParameters, *CloudHsmClustersClientBeginUpdateOptions) (*runtime.Poller[CloudHsmClustersClientUpdateResponse], error)`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.NewListByCloudHsmClusterPager(string, string, *PrivateEndpointConnectionsClientListByCloudHsmClusterOptions) *runtime.Pager[PrivateEndpointConnectionsClientListByCloudHsmClusterResponse]`
- New struct `BackupProperties`
- New struct `CHsmError`
- New struct `CloudHsmCluster`
- New struct `CloudHsmClusterError`
- New struct `CloudHsmClusterListResult`
- New struct `CloudHsmClusterPatchParameters`
- New struct `CloudHsmClusterPatchParametersProperties`
- New struct `CloudHsmClusterProperties`
- New struct `CloudHsmClusterResource`
- New struct `CloudHsmClusterSKU`
- New struct `CloudHsmClusterSecurityDomainProperties`
- New struct `CloudHsmProperties`
- New struct `DedicatedHsmResource`
- New struct `DedicatedHsmSystemData`
- New struct `ManagedServiceIdentity`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `ProxyResource`
- New struct `RestoreProperties`
- New struct `TrackedResource`
- New struct `UserAssignedIdentity`
- New field `Origin` in struct `DedicatedHsmOperation`
- New field `SystemData` in struct `Resource`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hardwaresecuritymodules/armhardwaresecuritymodules` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).