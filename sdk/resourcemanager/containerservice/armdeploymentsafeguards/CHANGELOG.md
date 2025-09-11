# Release History

## 0.2.0 (2025-09-10)
### Breaking Changes

- Function `*ClientFactory.NewDeploymentSafeguardsClient` has been removed
- Function `NewDeploymentSafeguardsClient` has been removed
- Function `*DeploymentSafeguardsClient.BeginCreate` has been removed
- Function `*DeploymentSafeguardsClient.BeginDelete` has been removed
- Function `*DeploymentSafeguardsClient.Get` has been removed
- Function `*DeploymentSafeguardsClient.NewListPager` has been removed

### Features Added

- New function `NewClient(azcore.TokenCredential, *arm.ClientOptions) (*Client, error)`
- New function `*Client.BeginCreate(context.Context, string, DeploymentSafeguard, *ClientBeginCreateOptions) (*runtime.Poller[ClientCreateResponse], error)`
- New function `*Client.BeginDelete(context.Context, string, *ClientBeginDeleteOptions) (*runtime.Poller[ClientDeleteResponse], error)`
- New function `*Client.Get(context.Context, string, *ClientGetOptions) (ClientGetResponse, error)`
- New function `*Client.NewListPager(string, *ClientListOptions) *runtime.Pager[ClientListResponse]`
- New function `*ClientFactory.NewClient() *Client`


## 0.1.0 (2025-06-27)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armdeploymentsafeguards` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).