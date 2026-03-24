# Release History

## 1.0.0 (2026-03-24)
### Features Added

- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New function `NewOperationsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationsClient, error)`
- New function `*OperationsClient.NewListPager(options *OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.1.0 (2024-10-21)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/databoundaries/armdataboundaries` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).