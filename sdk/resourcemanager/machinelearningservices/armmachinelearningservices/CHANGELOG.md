# Release History

## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New function `NewClientFactory(string, azcore.TokenCredential, *arm.ClientOptions) (*ClientFactory, error)`
- New function `*ClientFactory.NewComputeClient() *ComputeClient`
- New function `*ClientFactory.NewOperationsClient() *OperationsClient`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `*ClientFactory.NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient`
- New function `*ClientFactory.NewQuotasClient() *QuotasClient`
- New function `*ClientFactory.NewUsagesClient() *UsagesClient`
- New function `*ClientFactory.NewVirtualMachineSizesClient() *VirtualMachineSizesClient`
- New function `*ClientFactory.NewWorkspaceConnectionsClient() *WorkspaceConnectionsClient`
- New function `*ClientFactory.NewWorkspaceFeaturesClient() *WorkspaceFeaturesClient`
- New function `*ClientFactory.NewWorkspaceSKUsClient() *WorkspaceSKUsClient`
- New function `*ClientFactory.NewWorkspacesClient() *WorkspacesClient`
- New struct `ClientFactory`


## 1.0.1 (2022-05-30)

- Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning instead.

## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearningservices/armmachinelearningservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).