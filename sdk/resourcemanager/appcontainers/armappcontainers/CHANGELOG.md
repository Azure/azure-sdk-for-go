# Release History

## 1.1.0 (2023-04-07)
### Features Added

- New function `NewClientFactory(string, azcore.TokenCredential, *arm.ClientOptions) (*ClientFactory, error)`
- New function `*ClientFactory.NewCertificatesClient() *CertificatesClient`
- New function `*ClientFactory.NewContainerAppsAuthConfigsClient() *ContainerAppsAuthConfigsClient`
- New function `*ClientFactory.NewContainerAppsClient() *ContainerAppsClient`
- New function `*ClientFactory.NewContainerAppsRevisionReplicasClient() *ContainerAppsRevisionReplicasClient`
- New function `*ClientFactory.NewContainerAppsRevisionsClient() *ContainerAppsRevisionsClient`
- New function `*ClientFactory.NewContainerAppsSourceControlsClient() *ContainerAppsSourceControlsClient`
- New function `*ClientFactory.NewDaprComponentsClient() *DaprComponentsClient`
- New function `*ClientFactory.NewManagedEnvironmentsClient() *ManagedEnvironmentsClient`
- New function `*ClientFactory.NewManagedEnvironmentsStoragesClient() *ManagedEnvironmentsStoragesClient`
- New function `*ClientFactory.NewNamespacesClient() *NamespacesClient`
- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New struct `ClientFactory`


## 1.0.0 (2022-05-25)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).