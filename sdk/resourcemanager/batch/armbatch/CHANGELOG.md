# Release History

## 2.0.0 (2022-10-26)
### Breaking Changes

- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed
- Field `ActionRequired` of struct `PrivateLinkServiceConnectionState` has been removed
- Field `Username` of struct `CIFSMountConfiguration` has been removed
- Field `DynamicVNetAssignmentScope` of struct `NetworkConfiguration` has been removed

### Features Added

- New const `NodeCommunicationModeClassic`
- New const `EndpointAccessDefaultActionAllow`
- New const `NodeCommunicationModeSimplified`
- New const `EndpointAccessDefaultActionDeny`
- New const `PrivateEndpointConnectionProvisioningStateCreating`
- New const `PrivateEndpointConnectionProvisioningStateCancelled`
- New const `PrivateEndpointConnectionProvisioningStateDeleting`
- New const `NodeCommunicationModeDefault`
- New type alias `EndpointAccessDefaultAction`
- New type alias `NodeCommunicationMode`
- New function `PossibleNodeCommunicationModeValues() []NodeCommunicationMode`
- New function `*PrivateEndpointConnectionClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionClientDeleteResponse], error)`
- New function `PossibleEndpointAccessDefaultActionValues() []EndpointAccessDefaultAction`
- New struct `EndpointAccessProfile`
- New struct `IPRule`
- New struct `NetworkProfile`
- New struct `PrivateEndpointConnectionClientBeginDeleteOptions`
- New struct `PrivateEndpointConnectionClientDeleteResponse`
- New field `UserName` in struct `CIFSMountConfiguration`
- New field `GroupIDs` in struct `PrivateEndpointConnectionProperties`
- New field `NodeManagementEndpoint` in struct `AccountProperties`
- New field `NetworkProfile` in struct `AccountProperties`
- New field `TargetNodeCommunicationMode` in struct `PoolProperties`
- New field `CurrentNodeCommunicationMode` in struct `PoolProperties`
- New field `NetworkProfile` in struct `AccountCreateProperties`
- New field `NetworkProfile` in struct `AccountUpdateProperties`
- New field `PublicNetworkAccess` in struct `AccountUpdateProperties`
- New field `ActionsRequired` in struct `PrivateLinkServiceConnectionState`
- New field `DynamicVnetAssignmentScope` in struct `NetworkConfiguration`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/batch/armbatch` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).