# Release History

## 0.7.0 (2022-09-15)
### Features Added

- New const `IdentityTypeSystemAssigned`
- New const `IdentityTypeNone`
- New type alias `IdentityType`
- New function `NewAccessConnectorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccessConnectorsClient, error)`
- New function `*AccessConnectorsClient.BeginDelete(context.Context, string, string, *AccessConnectorsClientBeginDeleteOptions) (*runtime.Poller[AccessConnectorsClientDeleteResponse], error)`
- New function `*AccessConnectorsClient.BeginUpdate(context.Context, string, string, AccessConnectorUpdate, *AccessConnectorsClientBeginUpdateOptions) (*runtime.Poller[AccessConnectorsClientUpdateResponse], error)`
- New function `*AccessConnectorsClient.Get(context.Context, string, string, *AccessConnectorsClientGetOptions) (AccessConnectorsClientGetResponse, error)`
- New function `*AccessConnectorsClient.NewListByResourceGroupPager(string, *AccessConnectorsClientListByResourceGroupOptions) *runtime.Pager[AccessConnectorsClientListByResourceGroupResponse]`
- New function `PossibleIdentityTypeValues() []IdentityType`
- New function `*AccessConnectorsClient.NewListBySubscriptionPager(*AccessConnectorsClientListBySubscriptionOptions) *runtime.Pager[AccessConnectorsClientListBySubscriptionResponse]`
- New function `*AccessConnectorsClient.BeginCreateOrUpdate(context.Context, string, string, AccessConnector, *AccessConnectorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AccessConnectorsClientCreateOrUpdateResponse], error)`
- New struct `AccessConnector`
- New struct `AccessConnectorListResult`
- New struct `AccessConnectorProperties`
- New struct `AccessConnectorUpdate`
- New struct `AccessConnectorsClient`
- New struct `AccessConnectorsClientBeginCreateOrUpdateOptions`
- New struct `AccessConnectorsClientBeginDeleteOptions`
- New struct `AccessConnectorsClientBeginUpdateOptions`
- New struct `AccessConnectorsClientCreateOrUpdateResponse`
- New struct `AccessConnectorsClientDeleteResponse`
- New struct `AccessConnectorsClientGetOptions`
- New struct `AccessConnectorsClientGetResponse`
- New struct `AccessConnectorsClientListByResourceGroupOptions`
- New struct `AccessConnectorsClientListByResourceGroupResponse`
- New struct `AccessConnectorsClientListBySubscriptionOptions`
- New struct `AccessConnectorsClientListBySubscriptionResponse`
- New struct `AccessConnectorsClientUpdateResponse`
- New struct `IdentityData`


## 0.6.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databricks/armdatabricks` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).