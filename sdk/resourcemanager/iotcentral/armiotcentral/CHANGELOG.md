# Release History

## 2.0.0-beta.1 (2022-06-22)
### Breaking Changes

- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed
- Field `App` of struct `AppsClientUpdateResponse` has been removed
- Field `Location` of struct `Resource` has been removed
- Field `Tags` of struct `Resource` has been removed

### Features Added

- New const `CreatedByTypeApplication`
- New const `PrivateEndpointConnectionProvisioningStateSucceeded`
- New const `PrivateEndpointConnectionProvisioningStateCreating`
- New const `PrivateEndpointConnectionProvisioningStateDeleting`
- New const `PublicNetworkAccessEnabled`
- New const `PrivateEndpointServiceConnectionStatusPending`
- New const `NetworkActionDeny`
- New const `ProvisioningStateUpdating`
- New const `CreatedByTypeKey`
- New const `PrivateEndpointConnectionProvisioningStateFailed`
- New const `CreatedByTypeUser`
- New const `PrivateEndpointServiceConnectionStatusApproved`
- New const `IPRuleActionAllow`
- New const `NetworkActionAllow`
- New const `ProvisioningStateCanceled`
- New const `PublicNetworkAccessDisabled`
- New const `ProvisioningStateDeleting`
- New const `ProvisioningStateSucceeded`
- New const `ProvisioningStateCreating`
- New const `CreatedByTypeManagedIdentity`
- New const `ProvisioningStateFailed`
- New const `PrivateEndpointServiceConnectionStatusRejected`
- New function `PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `timeRFC3339.MarshalText() ([]byte, error)`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.BeginCreate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientBeginCreateOptions) (*runtime.Poller[PrivateEndpointConnectionsClientCreateResponse], error)`
- New function `NewPrivateLinksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinksClient, error)`
- New function `*PrivateEndpointConnectionsClient.NewListPager(string, string, *PrivateEndpointConnectionsClientListOptions) *runtime.Pager[PrivateEndpointConnectionsClientListResponse]`
- New function `*timeRFC3339.Parse(string) error`
- New function `*PrivateLinksClient.Get(context.Context, string, string, string, *PrivateLinksClientGetOptions) (PrivateLinksClientGetResponse, error)`
- New function `PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState`
- New function `*PrivateLinksClient.NewListPager(string, string, *PrivateLinksClientListOptions) *runtime.Pager[PrivateLinksClientListResponse]`
- New function `PossiblePublicNetworkAccessValues() []PublicNetworkAccess`
- New function `PossibleProvisioningStateValues() []ProvisioningState`
- New function `PossibleIPRuleActionValues() []IPRuleAction`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*timeRFC3339.UnmarshalText([]byte) error`
- New function `PossibleNetworkActionValues() []NetworkAction`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `NetworkRuleSetIPRule`
- New struct `NetworkRuleSets`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateEndpointConnectionsClient`
- New struct `PrivateEndpointConnectionsClientBeginCreateOptions`
- New struct `PrivateEndpointConnectionsClientBeginDeleteOptions`
- New struct `PrivateEndpointConnectionsClientCreateResponse`
- New struct `PrivateEndpointConnectionsClientDeleteResponse`
- New struct `PrivateEndpointConnectionsClientGetOptions`
- New struct `PrivateEndpointConnectionsClientGetResponse`
- New struct `PrivateEndpointConnectionsClientListOptions`
- New struct `PrivateEndpointConnectionsClientListResponse`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `PrivateLinksClient`
- New struct `PrivateLinksClientGetOptions`
- New struct `PrivateLinksClientGetResponse`
- New struct `PrivateLinksClientListOptions`
- New struct `PrivateLinksClientListResponse`
- New struct `SystemData`
- New struct `TrackedResource`
- New field `NetworkRuleSets` in struct `AppProperties`
- New field `PrivateEndpointConnections` in struct `AppProperties`
- New field `PublicNetworkAccess` in struct `AppProperties`
- New field `ProvisioningState` in struct `AppProperties`
- New field `SystemData` in struct `App`
- New field `SystemData` in struct `Resource`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotcentral/armiotcentral` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).