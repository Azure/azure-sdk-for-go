# Release History

## 1.1.0-beta.1 (2023-03-24)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New enum type `CmkIdentityType` with values `CmkIdentityTypeSystemAssignedIdentity`, `CmkIdentityTypeUserAssignedIdentity`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New function `*DatabasesClient.BeginFlush(context.Context, string, string, string, FlushParameters, *DatabasesClientBeginFlushOptions) (*runtime.Poller[DatabasesClientFlushResponse], error)`
- New function `NewSKUsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SKUsClient, error)`
- New function `*SKUsClient.NewListPager(string, *SKUsClientListOptions) *runtime.Pager[SKUsClientListResponse]`
- New function `timeRFC3339.MarshalText() ([]byte, error)`
- New function `*timeRFC3339.Parse(string) error`
- New function `*timeRFC3339.UnmarshalText([]byte) error`
- New struct `Capability`
- New struct `ClusterPropertiesEncryption`
- New struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryption`
- New struct `ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyIdentity`
- New struct `FlushParameters`
- New struct `LocationInfo`
- New struct `ManagedServiceIdentity`
- New struct `RegionSKUDetail`
- New struct `RegionSKUDetails`
- New struct `SKUDetail`
- New struct `SystemData`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Cluster`
- New field `SystemData` in struct `Cluster`
- New field `Encryption` in struct `ClusterProperties`
- New field `Identity` in struct `ClusterUpdate`
- New field `SystemData` in struct `Database`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `SystemData` in struct `ProxyResource`
- New field `SystemData` in struct `Resource`
- New field `SystemData` in struct `TrackedResource`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).