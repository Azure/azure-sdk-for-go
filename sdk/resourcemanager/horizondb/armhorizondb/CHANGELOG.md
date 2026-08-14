# Release History

## 0.2.0 (2026-08-14)
### Breaking Changes

- Function `*PrivateEndpointConnectionsClient.BeginDelete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, privateEndpointConnectionName string, options *PrivateEndpointConnectionsClientBeginDeleteOptions)` to `(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string, options *PrivateEndpointConnectionsClientBeginDeleteOptions)`
- Function `*PrivateEndpointConnectionsClient.BeginUpdate` has been removed
- Struct `OptionalPropertiesUpdateableProperties` has been removed
- Struct `PrivateEndpointConnection` has been removed
- Struct `PrivateEndpointConnectionUpdate` has been removed

### Features Added

- New value `StateSucceeded`, `StateUpgrading` added to enum type `State`
- New enum type `AuthenticationState` with values `AuthenticationStateDisabled`, `AuthenticationStateEnabled`
- New enum type `ComputeModelType` with values `ComputeModelTypeProvisioned`, `ComputeModelTypeServerless`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `PrincipalTypes` with values `PrincipalTypesGroup`, `PrincipalTypesServicePrincipal`, `PrincipalTypesUnknown`, `PrincipalTypesUser`
- New function `NewAdministratorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AdministratorsClient, error)`
- New function `*AdministratorsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, objectID string, resource AdministratorAdd, options *AdministratorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AdministratorsClientCreateOrUpdateResponse], error)`
- New function `*AdministratorsClient.BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, objectID string, options *AdministratorsClientBeginDeleteOptions) (*runtime.Poller[AdministratorsClientDeleteResponse], error)`
- New function `*AdministratorsClient.Get(ctx context.Context, resourceGroupName string, clusterName string, objectID string, options *AdministratorsClientGetOptions) (AdministratorsClientGetResponse, error)`
- New function `*AdministratorsClient.NewListPager(resourceGroupName string, clusterName string, options *AdministratorsClientListOptions) *runtime.Pager[AdministratorsClientListResponse]`
- New function `*ClientFactory.NewAdministratorsClient() *AdministratorsClient`
- New function `*ClustersClient.BeginRestart(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginRestartOptions) (*runtime.Poller[ClustersClientRestartResponse], error)`
- New function `*ClustersClient.BeginStart(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginStartOptions) (*runtime.Poller[ClustersClientStartResponse], error)`
- New function `*ClustersClient.BeginStop(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginStopOptions) (*runtime.Poller[ClustersClientStopResponse], error)`
- New function `*PrivateEndpointConnectionsClient.UpdateStatus(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string, resource PrivateEndpointConnectionResource, options *PrivateEndpointConnectionsClientUpdateStatusOptions) (PrivateEndpointConnectionsClientUpdateStatusResponse, error)`
- New struct `Administrator`
- New struct `AdministratorAdd`
- New struct `AdministratorListResult`
- New struct `AdministratorProperties`
- New struct `AdministratorPropertiesForAdd`
- New struct `ClusterAuthConfig`
- New struct `ClusterMirroring`
- New struct `ComputeModel`
- New struct `ManagedServiceIdentity`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Cluster`
- New field `Identity` in struct `ClusterForPatchUpdate`
- New field `AuthConfig`, `ComputeModel`, `Mirroring` in struct `ClusterProperties`
- New field `AuthConfig`, `ComputeModel`, `Mirroring` in struct `ClusterPropertiesForPatchUpdate`


## 0.1.0 (2026-03-30)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/horizondb/armhorizondb` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).