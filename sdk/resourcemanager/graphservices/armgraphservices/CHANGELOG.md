# Release History

## 1.0.0 (2023-06-23)
### Breaking Changes

- Function `NewAccountClient` has been removed
- Function `*AccountClient.BeginCreateAndUpdate` has been removed
- Function `*AccountClient.Delete` has been removed
- Function `*AccountClient.Get` has been removed
- Function `*AccountClient.Update` has been removed
- Function `*ClientFactory.NewAccountClient` has been removed
- Function `*ClientFactory.NewOperationClient` has been removed
- Function `NewOperationClient` has been removed
- Function `*OperationClient.NewListPager` has been removed

### Features Added

- New function `*AccountsClient.BeginCreateAndUpdate(context.Context, string, string, AccountResource, *AccountsClientBeginCreateAndUpdateOptions) (*runtime.Poller[AccountsClientCreateAndUpdateResponse], error)`
- New function `*AccountsClient.Delete(context.Context, string, string, *AccountsClientDeleteOptions) (AccountsClientDeleteResponse, error)`
- New function `*AccountsClient.Get(context.Context, string, string, *AccountsClientGetOptions) (AccountsClientGetResponse, error)`
- New function `*AccountsClient.Update(context.Context, string, string, AccountPatchResource, *AccountsClientUpdateOptions) (AccountsClientUpdateResponse, error)`
- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New function `NewOperationsClient(azcore.TokenCredential, *arm.ClientOptions) (*OperationsClient, error)`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.1.0 (2023-03-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/graphservices/armgraphservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).