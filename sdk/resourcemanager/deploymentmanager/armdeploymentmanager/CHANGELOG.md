# Release History

## 0.5.0 (2022-07-06)
### Breaking Changes

- Function `*StepsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, *StepsClientCreateOrUpdateOptions)` to `(context.Context, string, string, StepResource, *StepsClientCreateOrUpdateOptions)`
- Function `*ArtifactSourcesClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, *ArtifactSourcesClientCreateOrUpdateOptions)` to `(context.Context, string, string, ArtifactSource, *ArtifactSourcesClientCreateOrUpdateOptions)`
- Function `*RolloutsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, *RolloutsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, RolloutRequest, *RolloutsClientBeginCreateOrUpdateOptions)`
- Field `RolloutRequest` of struct `RolloutsClientBeginCreateOrUpdateOptions` has been removed
- Field `StepInfo` of struct `StepsClientCreateOrUpdateOptions` has been removed
- Field `ArtifactSourceInfo` of struct `ArtifactSourcesClientCreateOrUpdateOptions` has been removed


## 0.4.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deploymentmanager/armdeploymentmanager` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.4.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).