# Release History

## 0.7.0 (2026-03-19)
### Breaking Changes

- Function `*ConfigurationsClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serverGroupName string, configurationName string, options *ConfigurationsClientGetOptions)` to `(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, options *ConfigurationsClientGetOptions)`
- Function `*ConfigurationsClient.NewListByServerPager` parameter(s) have been changed from `(resourceGroupName string, serverGroupName string, serverName string, options *ConfigurationsClientListByServerOptions)` to `(resourceGroupName string, clusterName string, serverName string, options *ConfigurationsClientListByServerOptions)`
- Function `*FirewallRulesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serverGroupName string, firewallRuleName string, parameters FirewallRule, options *FirewallRulesClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, clusterName string, firewallRuleName string, parameters FirewallRule, options *FirewallRulesClientBeginCreateOrUpdateOptions)`
- Function `*FirewallRulesClient.BeginDelete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serverGroupName string, firewallRuleName string, options *FirewallRulesClientBeginDeleteOptions)` to `(ctx context.Context, resourceGroupName string, clusterName string, firewallRuleName string, options *FirewallRulesClientBeginDeleteOptions)`
- Function `*FirewallRulesClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serverGroupName string, firewallRuleName string, options *FirewallRulesClientGetOptions)` to `(ctx context.Context, resourceGroupName string, clusterName string, firewallRuleName string, options *FirewallRulesClientGetOptions)`
- Function `*RolesClient.BeginCreate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serverGroupName string, roleName string, parameters Role, options *RolesClientBeginCreateOptions)` to `(ctx context.Context, resourceGroupName string, clusterName string, roleName string, parameters Role, options *RolesClientBeginCreateOptions)`
- Function `*RolesClient.BeginDelete` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serverGroupName string, roleName string, options *RolesClientBeginDeleteOptions)` to `(ctx context.Context, resourceGroupName string, clusterName string, roleName string, options *RolesClientBeginDeleteOptions)`
- Function `*ServersClient.Get` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serverGroupName string, serverName string, options *ServersClientGetOptions)` to `(ctx context.Context, resourceGroupName string, clusterName string, serverName string, options *ServersClientGetOptions)`
- Type of `NameAvailabilityRequest.Type` has been changed from `*string` to `*CheckNameAvailabilityResourceType`
- Type of `Operation.Origin` has been changed from `*OperationOrigin` to `*Origin`
- Enum `CitusVersion` has been removed
- Enum `CreateMode` has been removed
- Enum `OperationOrigin` has been removed
- Enum `PostgreSQLVersion` has been removed
- Enum `ResourceProviderType` has been removed
- Enum `ServerEdition` has been removed
- Enum `ServerHaState` has been removed
- Enum `ServerState` has been removed
- Function `*ClientFactory.NewServerGroupsClient` has been removed
- Function `*ConfigurationsClient.NewListByServerGroupPager` has been removed
- Function `*ConfigurationsClient.BeginUpdate` has been removed
- Function `*FirewallRulesClient.NewListByServerGroupPager` has been removed
- Function `*RolesClient.NewListByServerGroupPager` has been removed
- Function `*ServersClient.NewListByServerGroupPager` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
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
- Struct `ServerGroupsClientBeginCreateOrUpdateOptions` has been removed
- Struct `ServerGroupsClientBeginDeleteOptions` has been removed
- Struct `ServerGroupsClientBeginRestartOptions` has been removed
- Struct `ServerGroupsClientBeginStartOptions` has been removed
- Struct `ServerGroupsClientBeginStopOptions` has been removed
- Struct `ServerGroupsClientBeginUpdateOptions` has been removed
- Struct `ServerGroupsClientCheckNameAvailabilityOptions` has been removed
- Struct `ServerGroupsClientGetOptions` has been removed
- Struct `ServerGroupsClientListByResourceGroupOptions` has been removed
- Struct `ServerGroupsClientListOptions` has been removed
- Struct `ServerProperties` has been removed
- Struct `ServerRoleGroup` has been removed
- Struct `TrackedResource` has been removed
- Field `Properties` of struct `Operation` has been removed

### Features Added

- New enum type `AADEnabledEnum` with values `AADEnabledEnumDisabled`, `AADEnabledEnumEnabled`
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `ActiveDirectoryAuth` with values `ActiveDirectoryAuthDisabled`, `ActiveDirectoryAuthEnabled`
- New enum type `CheckNameAvailabilityResourceType` with values `CheckNameAvailabilityResourceTypeMICROSOFTDBFORPOSTGRESQLSERVERGROUPSV2`
- New enum type `DataEncryptionType` with values `DataEncryptionTypeAzureKeyVault`, `DataEncryptionTypeSystemAssigned`
- New enum type `IdentityType` with values `IdentityTypeSystemAssigned`, `IdentityTypeUserAssigned`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `PasswordAuth` with values `PasswordAuthDisabled`, `PasswordAuthEnabled`
- New enum type `PasswordEnabledEnum` with values `PasswordEnabledEnumDisabled`, `PasswordEnabledEnumEnabled`
- New enum type `PrincipalType` with values `PrincipalTypeGroup`, `PrincipalTypeServicePrincipal`, `PrincipalTypeUser`
- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateFailed`, `ProvisioningStateInProgress`, `ProvisioningStateSucceeded`
- New enum type `RoleType` with values `RoleTypeAdmin`, `RoleTypeUser`
- New function `*ClientFactory.NewClustersClient() *ClustersClient`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `*ClientFactory.NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient`
- New function `NewClustersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClustersClient, error)`
- New function `*ClustersClient.CheckNameAvailability(ctx context.Context, nameAvailabilityRequest NameAvailabilityRequest, options *ClustersClientCheckNameAvailabilityOptions) (ClustersClientCheckNameAvailabilityResponse, error)`
- New function `*ClustersClient.BeginCreate(ctx context.Context, resourceGroupName string, clusterName string, parameters Cluster, options *ClustersClientBeginCreateOptions) (*runtime.Poller[ClustersClientCreateResponse], error)`
- New function `*ClustersClient.BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginDeleteOptions) (*runtime.Poller[ClustersClientDeleteResponse], error)`
- New function `*ClustersClient.Get(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientGetOptions) (ClustersClientGetResponse, error)`
- New function `*ClustersClient.NewListByResourceGroupPager(resourceGroupName string, options *ClustersClientListByResourceGroupOptions) *runtime.Pager[ClustersClientListByResourceGroupResponse]`
- New function `*ClustersClient.NewListPager(options *ClustersClientListOptions) *runtime.Pager[ClustersClientListResponse]`
- New function `*ClustersClient.BeginPromoteReadReplica(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginPromoteReadReplicaOptions) (*runtime.Poller[ClustersClientPromoteReadReplicaResponse], error)`
- New function `*ClustersClient.BeginRestart(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginRestartOptions) (*runtime.Poller[ClustersClientRestartResponse], error)`
- New function `*ClustersClient.BeginStart(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginStartOptions) (*runtime.Poller[ClustersClientStartResponse], error)`
- New function `*ClustersClient.BeginStop(ctx context.Context, resourceGroupName string, clusterName string, options *ClustersClientBeginStopOptions) (*runtime.Poller[ClustersClientStopResponse], error)`
- New function `*ClustersClient.BeginUpdate(ctx context.Context, resourceGroupName string, clusterName string, parameters ClusterForUpdate, options *ClustersClientBeginUpdateOptions) (*runtime.Poller[ClustersClientUpdateResponse], error)`
- New function `*ConfigurationsClient.GetCoordinator(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, options *ConfigurationsClientGetCoordinatorOptions) (ConfigurationsClientGetCoordinatorResponse, error)`
- New function `*ConfigurationsClient.GetNode(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, options *ConfigurationsClientGetNodeOptions) (ConfigurationsClientGetNodeResponse, error)`
- New function `*ConfigurationsClient.NewListByClusterPager(resourceGroupName string, clusterName string, options *ConfigurationsClientListByClusterOptions) *runtime.Pager[ConfigurationsClientListByClusterResponse]`
- New function `*ConfigurationsClient.BeginUpdateOnCoordinator(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, parameters ServerConfiguration, options *ConfigurationsClientBeginUpdateOnCoordinatorOptions) (*runtime.Poller[ConfigurationsClientUpdateOnCoordinatorResponse], error)`
- New function `*ConfigurationsClient.BeginUpdateOnNode(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, parameters ServerConfiguration, options *ConfigurationsClientBeginUpdateOnNodeOptions) (*runtime.Poller[ConfigurationsClientUpdateOnNodeResponse], error)`
- New function `*FirewallRulesClient.NewListByClusterPager(resourceGroupName string, clusterName string, options *FirewallRulesClientListByClusterOptions) *runtime.Pager[FirewallRulesClientListByClusterResponse]`
- New function `NewPrivateEndpointConnectionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string, parameters PrivateEndpointConnection, options *PrivateEndpointConnectionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[PrivateEndpointConnectionsClientCreateOrUpdateResponse], error)`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string, options *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.Get(ctx context.Context, resourceGroupName string, clusterName string, privateEndpointConnectionName string, options *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.NewListByClusterPager(resourceGroupName string, clusterName string, options *PrivateEndpointConnectionsClientListByClusterOptions) *runtime.Pager[PrivateEndpointConnectionsClientListByClusterResponse]`
- New function `NewPrivateLinkResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.Get(ctx context.Context, resourceGroupName string, clusterName string, privateLinkResourceName string, options *PrivateLinkResourcesClientGetOptions) (PrivateLinkResourcesClientGetResponse, error)`
- New function `*PrivateLinkResourcesClient.NewListByClusterPager(resourceGroupName string, clusterName string, options *PrivateLinkResourcesClientListByClusterOptions) *runtime.Pager[PrivateLinkResourcesClientListByClusterResponse]`
- New function `*RolesClient.Get(ctx context.Context, resourceGroupName string, clusterName string, roleName string, options *RolesClientGetOptions) (RolesClientGetResponse, error)`
- New function `*RolesClient.NewListByClusterPager(resourceGroupName string, clusterName string, options *RolesClientListByClusterOptions) *runtime.Pager[RolesClientListByClusterResponse]`
- New function `*ServersClient.NewListByClusterPager(resourceGroupName string, clusterName string, options *ServersClientListByClusterOptions) *runtime.Pager[ServersClientListByClusterResponse]`
- New struct `AuthConfig`
- New struct `Cluster`
- New struct `ClusterConfigurationListResult`
- New struct `ClusterForUpdate`
- New struct `ClusterListResult`
- New struct `ClusterProperties`
- New struct `ClusterPropertiesForUpdate`
- New struct `ClusterServer`
- New struct `ClusterServerListResult`
- New struct `ClusterServerProperties`
- New struct `Configuration`
- New struct `ConfigurationProperties`
- New struct `DataEncryption`
- New struct `IdentityProperties`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateEndpointConnectionSimpleProperties`
- New struct `PrivateEndpointProperty`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `PromoteRequest`
- New struct `RolePropertiesExternalIdentity`
- New struct `SimplePrivateEndpointConnection`
- New struct `UserAssignedIdentity`
- New field `NextLink` in struct `FirewallRuleListResult`
- New field `ProvisioningState` in struct `FirewallRuleProperties`
- New field `ActionType` in struct `Operation`
- New field `NextLink` in struct `RoleListResult`
- New field `ExternalIdentity`, `ProvisioningState`, `RoleType` in struct `RoleProperties`
- New field `ProvisioningState`, `RequiresRestart` in struct `ServerConfigurationProperties`


## 0.6.1 (2023-06-23)
### Other Changes

- Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmosforpostgresql/armcosmosforpostgresql instead.


## 0.6.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresqlhsc/armpostgresqlhsc` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).