# Release History

## 1.0.0 (2023-03-22)
### Features Added

- New function `NewClientFactory(azcore.TokenCredential, *arm.ClientOptions) (*ClientFactory, error)`
- New function `*ClientFactory.NewClient() *Client`
- New function `*ClientFactory.NewOperationClient() *OperationClient`
- New function `*ClientFactory.NewRequestStatusClient() *RequestStatusClient`
- New function `*ClientFactory.NewUsagesClient() *UsagesClient`
- New struct `ClientFactory`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).