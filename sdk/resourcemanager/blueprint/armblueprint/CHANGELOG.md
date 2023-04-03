# Release History

## 0.6.0 (2023-04-03)
### Breaking Changes

- Struct `ArtifactPropertiesBase` has been removed
- Struct `CloudError` has been removed
- Struct `ResourcePropertiesBase` has been removed
- Struct `SharedBlueprintProperties` has been removed

### Features Added

- New function `NewClientFactory(azcore.TokenCredential, *arm.ClientOptions) (*ClientFactory, error)`
- New function `*ClientFactory.NewArtifactsClient() *ArtifactsClient`
- New function `*ClientFactory.NewAssignmentOperationsClient() *AssignmentOperationsClient`
- New function `*ClientFactory.NewAssignmentsClient() *AssignmentsClient`
- New function `*ClientFactory.NewBlueprintsClient() *BlueprintsClient`
- New function `*ClientFactory.NewPublishedArtifactsClient() *PublishedArtifactsClient`
- New function `*ClientFactory.NewPublishedBlueprintsClient() *PublishedBlueprintsClient`
- New struct `ClientFactory`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/blueprint/armblueprint` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).