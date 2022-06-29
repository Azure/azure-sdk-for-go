# Release History

## 2.0.0 (2022-06-29)
### Breaking Changes

- Function `*PrivateLinkForAzureAdClient.Update` parameter(s) have been changed from `(context.Context, string, string, *PrivateLinkForAzureAdClientUpdateOptions)` to `(context.Context, string, string, PrivateLinkPolicyUpdateParameter, *PrivateLinkForAzureAdClientUpdateOptions)`
- Field `PrivateLinkPolicy` of struct `PrivateLinkForAzureAdClientUpdateOptions` has been removed


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/aad/armaad` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).