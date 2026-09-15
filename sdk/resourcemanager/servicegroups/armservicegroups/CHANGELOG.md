# Release History

## 1.0.0 (2026-08-18)
### Breaking Changes

- Function `*Client.ListAncestors` has been removed
- Function `*ManagementClient.NewClient` has been removed
- Struct `ServiceGroupCollectionResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New function `NewOperationsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationsClient, error)`
- New function `*OperationsClient.NewListPager(options *OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `ServiceGroupAttributes`
- New field `Attributes` in struct `ServiceGroupProperties`


## 0.1.0 (2026-03-27)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicegroups/armservicegroups` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).