# Release History

## 0.2.0 (2022-10-13)
### Breaking Changes

- Function `*WorkflowClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Workflow, *WorkflowClientCreateOrUpdateOptions)` to `(context.Context, string, Workflow, *WorkflowClientCreateOrUpdateOptions)`
- Function `*WorkflowClient.UpdateTags` parameter(s) have been changed from `(context.Context, string, string, TagsObject, *WorkflowClientUpdateTagsOptions)` to `(context.Context, string, TagsObject, *WorkflowClientUpdateTagsOptions)`
- Function `*WorkflowClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *WorkflowClientDeleteOptions)` to `(context.Context, string, *WorkflowClientDeleteOptions)`
- Function `NewWorkflowClient` parameter(s) have been changed from `(string, *string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, string, *string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*WorkflowClient.Get` parameter(s) have been changed from `(context.Context, string, string, *WorkflowClientGetOptions)` to `(context.Context, string, *WorkflowClientGetOptions)`


## 0.1.1 (2022-10-12)
### Other Changes
- Loosen Go version requirement.

## 0.1.0 (2022-09-24)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devhub/armdevhub` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).