# Release History

## 0.3.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.2.0 (2022-10-13)
### Breaking Changes

- Function `NewWorkflowClient` parameter(s) have been changed from `(string, *string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewDeveloperHubServiceClient` parameter(s) have been changed from `(string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*DeveloperHubServiceClient.GitHubOAuthCallback` parameter(s) have been changed from `(context.Context, string, *DeveloperHubServiceClientGitHubOAuthCallbackOptions)` to `(context.Context, string, string, string, *DeveloperHubServiceClientGitHubOAuthCallbackOptions)`

### Features Added

- New field `ManagedClusterResource` in struct `WorkflowClientListByResourceGroupOptions`


## 0.1.1 (2022-10-12)
### Other Changes
- Loosen Go version requirement.

## 0.1.0 (2022-09-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devhub/armdevhub` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).