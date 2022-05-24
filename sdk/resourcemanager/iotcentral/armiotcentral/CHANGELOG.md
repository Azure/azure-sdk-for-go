# Release History

## 2.0.0-beta.1 (2022-05-24)
### Breaking Changes

- Function `Resource.MarshalJSON` has been removed
- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed
- Field `Tags` of struct `Resource` has been removed
- Field `Location` of struct `Resource` has been removed
- Field `App` of struct `AppsClientUpdateResponse` has been removed

### Features Added

- New const `NetworkActionAllow`
- New const `PrivateEndpointServiceConnectionStatusRejected`
- New const `ProvisioningStateUpdating`
- New const `CreatedByTypeApplication`
- New const `ProvisioningStateFailed`
- New const `ProvisioningStateCreating`
- New const `CreatedByTypeManagedIdentity`
- New const `ProvisioningStateCanceled`
- New const `PrivateEndpointConnectionProvisioningStateFailed`
- New const `PublicNetworkAccessEnabled`
- New const `IPRuleActionAllow`
- New const `PrivateEndpointConnectionProvisioningStateSucceeded`
- New const `CreatedByTypeUser`
- New const `PrivateEndpointServiceConnectionStatusPending`
- New const `CreatedByTypeKey`
- New const `PrivateEndpointServiceConnectionStatusApproved`
- New const `PrivateEndpointConnectionProvisioningStateDeleting`
- New const `NetworkActionDeny`
- New const `PublicNetworkAccessDisabled`
- New const `ProvisioningStateDeleting`
- New const `ProvisioningStateSucceeded`
- New const `PrivateEndpointConnectionProvisioningStateCreating`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New function `PossibleNetworkActionValues() []NetworkAction`
- New function `PrivateLinkResourceProperties.MarshalJSON() ([]byte, error)`
- New function `*SystemData.UnmarshalJSON([]byte) error`
- New function `TrackedResource.MarshalJSON() ([]byte, error)`
- New function `AppProperties.MarshalJSON() ([]byte, error)`
- New function `PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState`
- New function `PossiblePublicNetworkAccessValues() []PublicNetworkAccess`
- New function `PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus`
- New function `NetworkRuleSets.MarshalJSON() ([]byte, error)`
- New function `PossibleProvisioningStateValues() []ProvisioningState`
- New function `PrivateEndpointConnectionProperties.MarshalJSON() ([]byte, error)`
- New function `SystemData.MarshalJSON() ([]byte, error)`
- New function `PossibleIPRuleActionValues() []IPRuleAction`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `NetworkRuleSetIPRule`
- New struct `NetworkRuleSets`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
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
- New struct `PrivateLinksClientGetOptions`
- New struct `PrivateLinksClientGetResponse`
- New struct `PrivateLinksClientListOptions`
- New struct `PrivateLinksClientListResponse`
- New struct `SystemData`
- New struct `TrackedResource`
- New field `PrivateEndpointConnections` in struct `AppProperties`
- New field `NetworkRuleSets` in struct `AppProperties`
- New field `PublicNetworkAccess` in struct `AppProperties`
- New field `ProvisioningState` in struct `AppProperties`
- New field `SystemData` in struct `App`
- New field `SystemData` in struct `Resource`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotcentral/armiotcentral` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).