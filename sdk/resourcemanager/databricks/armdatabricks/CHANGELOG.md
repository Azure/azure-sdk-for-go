# Release History

## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.0.0 (2023-10-27)
### Breaking Changes

- Field `ActionRequired` of struct `PrivateLinkServiceConnectionState` has been removed

### Features Added

- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New function `NewAccessConnectorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessConnectorsClient, error)`
- New function `*AccessConnectorsClient.BeginCreateOrUpdate(context.Context, string, string, AccessConnector, *AccessConnectorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AccessConnectorsClientCreateOrUpdateResponse], error)`
- New function `*AccessConnectorsClient.BeginDelete(context.Context, string, string, *AccessConnectorsClientBeginDeleteOptions) (*runtime.Poller[AccessConnectorsClientDeleteResponse], error)`
- New function `*AccessConnectorsClient.Get(context.Context, string, string, *AccessConnectorsClientGetOptions) (AccessConnectorsClientGetResponse, error)`
- New function `*AccessConnectorsClient.NewListByResourceGroupPager(string, *AccessConnectorsClientListByResourceGroupOptions) *runtime.Pager[AccessConnectorsClientListByResourceGroupResponse]`
- New function `*AccessConnectorsClient.NewListBySubscriptionPager(*AccessConnectorsClientListBySubscriptionOptions) *runtime.Pager[AccessConnectorsClientListBySubscriptionResponse]`
- New function `*AccessConnectorsClient.BeginUpdate(context.Context, string, string, AccessConnectorUpdate, *AccessConnectorsClientBeginUpdateOptions) (*runtime.Poller[AccessConnectorsClientUpdateResponse], error)`
- New function `*ClientFactory.NewAccessConnectorsClient() *AccessConnectorsClient`
- New struct `AccessConnector`
- New struct `AccessConnectorListResult`
- New struct `AccessConnectorProperties`
- New struct `AccessConnectorUpdate`
- New struct `ManagedDiskEncryption`
- New struct `ManagedDiskEncryptionKeyVaultProperties`
- New struct `ManagedServiceIdentity`
- New struct `UserAssignedIdentity`
- New field `ManagedDisk` in struct `EncryptionEntitiesDefinition`
- New field `Description` in struct `OperationDisplay`
- New field `GroupIDs` in struct `PrivateEndpointConnectionProperties`
- New field `ActionsRequired` in struct `PrivateLinkServiceConnectionState`
- New field `DiskEncryptionSetID`, `ManagedDiskIdentity` in struct `WorkspaceProperties`


## 0.7.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.7.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.6.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databricks/armdatabricks` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).