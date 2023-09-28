# Release History

## 1.2.0-beta.1 (2023-09-28)
### Features Added

- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `*ClientFactory.NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.CreateOrUpdate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientCreateOrUpdateOptions) (PrivateEndpointConnectionsClientCreateOrUpdateResponse, error)`
- New function `*PrivateEndpointConnectionsClient.Delete(context.Context, string, string, string, *PrivateEndpointConnectionsClientDeleteOptions) (PrivateEndpointConnectionsClientDeleteResponse, error)`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.ListByEnvironment(context.Context, string, string, *PrivateEndpointConnectionsClientListByEnvironmentOptions) (PrivateEndpointConnectionsClientListByEnvironmentResponse, error)`
- New function `NewPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.NewListSupportedPager(string, string, *PrivateLinkResourcesClientListSupportedOptions) *runtime.Pager[PrivateLinkResourcesClientListSupportedResponse]`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New field `PrivateEndpointConnections`, `PublicNetworkAccess` in struct `Gen2EnvironmentCreationProperties`
- New field `PrivateEndpointConnections`, `PublicNetworkAccess` in struct `Gen2EnvironmentResourceProperties`


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-06-10)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/timeseriesinsights/armtimeseriesinsights` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).