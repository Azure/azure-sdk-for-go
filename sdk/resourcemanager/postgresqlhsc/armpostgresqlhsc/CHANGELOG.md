# Release History

## 1.0.0 (2023-02-28)
### Breaking Changes

- Type of `ServerProperties.ServerEdition` has been changed from `*ServerEdition` to `*string`
- Type alias `CitusVersion` has been removed
- Type alias `CreateMode` has been removed
- Type alias `PostgreSQLVersion` has been removed
- Type alias `ResourceProviderType` has been removed
- Type alias `ServerEdition` has been removed
- Type alias `ServerHaState` has been removed
- Type alias `ServerState` has been removed
- Function `*ConfigurationsClient.NewListByServerGroupPager` has been removed
- Function `*ConfigurationsClient.BeginUpdate` has been removed
- Function `*FirewallRulesClient.NewListByServerGroupPager` has been removed
- Function `*RolesClient.NewListByServerGroupPager` has been removed
- Function `NewServerGroupsClient` has been removed
- Function `*ServerGroupsClient.CheckNameAvailability` has been removed
- Function `*ServerGroupsClient.BeginCreateOrUpdate` has been removed
- Function `*ServerGroupsClient.BeginDelete` has been removed
- Function `*ServerGroupsClient.Get` has been removed
- Function `*ServerGroupsClient.NewListByResourceGroupPager` has been removed
- Function `*ServerGroupsClient.NewListPager` has been removed
- Function `*ServerGroupsClient.BeginRestart` has been removed
- Function `*ServerGroupsClient.BeginStart` has been removed
- Function `*ServerGroupsClient.BeginStop` has been removed
- Function `*ServerGroupsClient.BeginUpdate` has been removed
- Function `*ServersClient.NewListByServerGroupPager` has been removed
- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed
- Struct `ServerGroup` has been removed
- Struct `ServerGroupConfiguration` has been removed
- Struct `ServerGroupConfigurationListResult` has been removed
- Struct `ServerGroupConfigurationProperties` has been removed
- Struct `ServerGroupForUpdate` has been removed
- Struct `ServerGroupListResult` has been removed
- Struct `ServerGroupProperties` has been removed
- Struct `ServerGroupPropertiesDelegatedSubnetArguments` has been removed
- Struct `ServerGroupPropertiesForUpdate` has been removed
- Struct `ServerGroupPropertiesPrivateDNSZoneArguments` has been removed
- Struct `ServerGroupServer` has been removed
- Struct `ServerGroupServerListResult` has been removed
- Struct `ServerGroupServerProperties` has been removed
- Struct `ServerGroupsClient` has been removed
- Struct `ServerRoleGroup` has been removed
- Field `ServerGroupConfiguration` of struct `ConfigurationsClientGetResponse` has been removed
- Field `EnablePublicIP` of struct `ServerProperties` has been removed
- Field `ServerGroupServer` of struct `ServersClientGetResponse` has been removed

### Features Added

- New type alias `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New type alias `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New type alias `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateFailed`, `ProvisioningStateInProgress`, `ProvisioningStateSucceeded`
- New function `NewClustersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ClustersClient, error)`
- New function `*ClustersClient.CheckNameAvailability(context.Context, NameAvailabilityRequest, *ClustersClientCheckNameAvailabilityOptions) (ClustersClientCheckNameAvailabilityResponse, error)`
- New function `*ClustersClient.BeginCreateOrUpdate(context.Context, string, string, Cluster, *ClustersClientBeginCreateOrUpdateOptions) (*runtime.Poller[ClustersClientCreateOrUpdateResponse], error)`
- New function `*ClustersClient.BeginDelete(context.Context, string, string, *ClustersClientBeginDeleteOptions) (*runtime.Poller[ClustersClientDeleteResponse], error)`
- New function `*ClustersClient.Get(context.Context, string, string, *ClustersClientGetOptions) (ClustersClientGetResponse, error)`
- New function `*ClustersClient.NewListByResourceGroupPager(string, *ClustersClientListByResourceGroupOptions) *runtime.Pager[ClustersClientListByResourceGroupResponse]`
- New function `*ClustersClient.NewListPager(*ClustersClientListOptions) *runtime.Pager[ClustersClientListResponse]`
- New function `*ClustersClient.BeginPromoteReadReplica(context.Context, string, string, *ClustersClientBeginPromoteReadReplicaOptions) (*runtime.Poller[ClustersClientPromoteReadReplicaResponse], error)`
- New function `*ClustersClient.BeginRestart(context.Context, string, string, *ClustersClientBeginRestartOptions) (*runtime.Poller[ClustersClientRestartResponse], error)`
- New function `*ClustersClient.BeginStart(context.Context, string, string, *ClustersClientBeginStartOptions) (*runtime.Poller[ClustersClientStartResponse], error)`
- New function `*ClustersClient.BeginStop(context.Context, string, string, *ClustersClientBeginStopOptions) (*runtime.Poller[ClustersClientStopResponse], error)`
- New function `*ClustersClient.BeginUpdate(context.Context, string, string, ClusterForUpdate, *ClustersClientBeginUpdateOptions) (*runtime.Poller[ClustersClientUpdateResponse], error)`
- New function `*ConfigurationsClient.BeginCreateOrUpdateCoordinator(context.Context, string, string, string, ServerConfiguration, *ConfigurationsClientBeginCreateOrUpdateCoordinatorOptions) (*runtime.Poller[ConfigurationsClientCreateOrUpdateCoordinatorResponse], error)`
- New function `*ConfigurationsClient.BeginCreateOrUpdateNode(context.Context, string, string, string, ServerConfiguration, *ConfigurationsClientBeginCreateOrUpdateNodeOptions) (*runtime.Poller[ConfigurationsClientCreateOrUpdateNodeResponse], error)`
- New function `*ConfigurationsClient.GetCoordinator(context.Context, string, string, string, *ConfigurationsClientGetCoordinatorOptions) (ConfigurationsClientGetCoordinatorResponse, error)`
- New function `*ConfigurationsClient.GetNode(context.Context, string, string, string, *ConfigurationsClientGetNodeOptions) (ConfigurationsClientGetNodeResponse, error)`
- New function `*ConfigurationsClient.NewListByClusterPager(string, string, *ConfigurationsClientListByClusterOptions) *runtime.Pager[ConfigurationsClientListByClusterResponse]`
- New function `*FirewallRulesClient.NewListByClusterPager(string, string, *FirewallRulesClientListByClusterOptions) *runtime.Pager[FirewallRulesClientListByClusterResponse]`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.BeginCreateOrUpdate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[PrivateEndpointConnectionsClientCreateOrUpdateResponse], error)`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.NewListByClusterPager(string, string, *PrivateEndpointConnectionsClientListByClusterOptions) *runtime.Pager[PrivateEndpointConnectionsClientListByClusterResponse]`
- New function `NewPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.Get(context.Context, string, string, string, *PrivateLinkResourcesClientGetOptions) (PrivateLinkResourcesClientGetResponse, error)`
- New function `*PrivateLinkResourcesClient.NewListByClusterPager(string, string, *PrivateLinkResourcesClientListByClusterOptions) *runtime.Pager[PrivateLinkResourcesClientListByClusterResponse]`
- New function `*RolesClient.Get(context.Context, string, string, string, *RolesClientGetOptions) (RolesClientGetResponse, error)`
- New function `*RolesClient.NewListByClusterPager(string, string, *RolesClientListByClusterOptions) *runtime.Pager[RolesClientListByClusterResponse]`
- New function `*ServersClient.NewListByClusterPager(string, string, *ServersClientListByClusterOptions) *runtime.Pager[ServersClientListByClusterResponse]`
- New struct `Cluster`
- New struct `ClusterConfigurationListResult`
- New struct `ClusterForUpdate`
- New struct `ClusterListResult`
- New struct `ClusterProperties`
- New struct `ClusterPropertiesForUpdate`
- New struct `ClusterServer`
- New struct `ClusterServerListResult`
- New struct `ClusterServerProperties`
- New struct `ClustersClient`
- New struct `Configuration`
- New struct `ConfigurationProperties`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateEndpointConnectionSimpleProperties`
- New struct `PrivateEndpointConnectionsClient`
- New struct `PrivateEndpointProperty`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkResourcesClient`
- New struct `PrivateLinkServiceConnectionState`
- New struct `SimplePrivateEndpointConnection`
- New anonymous field `Configuration` in struct `ConfigurationsClientGetResponse`
- New field `ProvisioningState` in struct `FirewallRuleProperties`
- New field `SystemData` in struct `ProxyResource`
- New field `SystemData` in struct `Resource`
- New field `ProvisioningState` in struct `RoleProperties`
- New field `ProvisioningState` in struct `ServerConfigurationProperties`
- New field `RequiresRestart` in struct `ServerConfigurationProperties`
- New field `AdministratorLogin` in struct `ServerProperties`
- New field `EnablePublicIPAccess` in struct `ServerProperties`
- New field `IsReadOnly` in struct `ServerProperties`
- New anonymous field `ClusterServer` in struct `ServersClientGetResponse`
- New field `SystemData` in struct `TrackedResource`


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresqlhsc/armpostgresqlhsc` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).