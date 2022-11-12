# Release History

## 2.0.0 (2022-11-12)
### Breaking Changes

- Field `InstalledVersion` of struct `ExtensionProperties` has been removed

### Features Added

- New const `SourceKindTypeAzureBlob`
- New struct `AzureBlobDefinition`
- New struct `AzureBlobPatchDefinition`
- New struct `ManagedIdentityDefinition`
- New struct `ManagedIdentityPatchDefinition`
- New struct `Plan`
- New struct `ServicePrincipalDefinition`
- New struct `ServicePrincipalPatchDefinition`
- New field `CurrentVersion` in struct `ExtensionProperties`
- New field `IsSystemExtension` in struct `ExtensionProperties`
- New field `AzureBlob` in struct `FluxConfigurationProperties`
- New field `Plan` in struct `Extension`
- New field `AzureBlob` in struct `FluxConfigurationPatchProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kubernetesconfiguration/armkubernetesconfiguration` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).