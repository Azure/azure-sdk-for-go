# Release History

## 1.1.0 (2026-03-09)
### Features Added

- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New function `*SchedulersClient.BeginCreateOrUpdatePrivateEndpointConnection(ctx context.Context, resourceGroupName string, schedulerName string, privateEndpointConnectionName string, resource PrivateEndpointConnection, options *SchedulersClientBeginCreateOrUpdatePrivateEndpointConnectionOptions) (*runtime.Poller[SchedulersClientCreateOrUpdatePrivateEndpointConnectionResponse], error)`
- New function `*SchedulersClient.BeginDeletePrivateEndpointConnection(ctx context.Context, resourceGroupName string, schedulerName string, privateEndpointConnectionName string, options *SchedulersClientBeginDeletePrivateEndpointConnectionOptions) (*runtime.Poller[SchedulersClientDeletePrivateEndpointConnectionResponse], error)`
- New function `*SchedulersClient.GetPrivateEndpointConnection(ctx context.Context, resourceGroupName string, schedulerName string, privateEndpointConnectionName string, options *SchedulersClientGetPrivateEndpointConnectionOptions) (SchedulersClientGetPrivateEndpointConnectionResponse, error)`
- New function `*SchedulersClient.GetPrivateLink(ctx context.Context, resourceGroupName string, schedulerName string, privateLinkResourceName string, options *SchedulersClientGetPrivateLinkOptions) (SchedulersClientGetPrivateLinkResponse, error)`
- New function `*SchedulersClient.NewListPrivateEndpointConnectionsPager(resourceGroupName string, schedulerName string, options *SchedulersClientListPrivateEndpointConnectionsOptions) *runtime.Pager[SchedulersClientListPrivateEndpointConnectionsResponse]`
- New function `*SchedulersClient.NewListPrivateLinksPager(resourceGroupName string, schedulerName string, options *SchedulersClientListPrivateLinksOptions) *runtime.Pager[SchedulersClientListPrivateLinksResponse]`
- New function `*SchedulersClient.BeginUpdatePrivateEndpointConnection(ctx context.Context, resourceGroupName string, schedulerName string, privateEndpointConnectionName string, properties PrivateEndpointConnectionUpdate, options *SchedulersClientBeginUpdatePrivateEndpointConnectionOptions) (*runtime.Poller[SchedulersClientUpdatePrivateEndpointConnectionResponse], error)`
- New struct `OptionalPropertiesUpdateableProperties`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateEndpointConnectionUpdate`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `SchedulerPrivateLinkResource`
- New struct `SchedulerPrivateLinkResourceListResult`
- New field `PrivateEndpointConnections`, `PublicNetworkAccess` in struct `SchedulerProperties`
- New field `PublicNetworkAccess` in struct `SchedulerPropertiesUpdate`


## 1.0.0 (2025-09-26)
### Breaking Changes

- Type of `SchedulerSKU.Name` has been changed from `*string` to `*SchedulerSKUName`
- Type of `SchedulerSKUUpdate.Name` has been changed from `*string` to `*SchedulerSKUName`

### Features Added

- New enum type `SchedulerSKUName` with values `SchedulerSKUNameConsumption`, `SchedulerSKUNameDedicated`


## 0.2.0 (2025-04-15)
### Features Added

- New enum type `PurgeableOrchestrationState` with values `PurgeableOrchestrationStateCanceled`, `PurgeableOrchestrationStateCompleted`, `PurgeableOrchestrationStateFailed`, `PurgeableOrchestrationStateTerminated`
- New function `*ClientFactory.NewRetentionPoliciesClient() *RetentionPoliciesClient`
- New function `NewRetentionPoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RetentionPoliciesClient, error)`
- New function `*RetentionPoliciesClient.BeginCreateOrReplace(context.Context, string, string, RetentionPolicy, *RetentionPoliciesClientBeginCreateOrReplaceOptions) (*runtime.Poller[RetentionPoliciesClientCreateOrReplaceResponse], error)`
- New function `*RetentionPoliciesClient.BeginDelete(context.Context, string, string, *RetentionPoliciesClientBeginDeleteOptions) (*runtime.Poller[RetentionPoliciesClientDeleteResponse], error)`
- New function `*RetentionPoliciesClient.Get(context.Context, string, string, *RetentionPoliciesClientGetOptions) (RetentionPoliciesClientGetResponse, error)`
- New function `*RetentionPoliciesClient.NewListBySchedulerPager(string, string, *RetentionPoliciesClientListBySchedulerOptions) *runtime.Pager[RetentionPoliciesClientListBySchedulerResponse]`
- New function `*RetentionPoliciesClient.BeginUpdate(context.Context, string, string, RetentionPolicy, *RetentionPoliciesClientBeginUpdateOptions) (*runtime.Poller[RetentionPoliciesClientUpdateResponse], error)`
- New struct `RetentionPolicy`
- New struct `RetentionPolicyDetails`
- New struct `RetentionPolicyListResult`
- New struct `RetentionPolicyProperties`


## 0.1.0 (2025-03-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/durabletask/armdurabletask` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).