# Release History

## 2.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.1.0 (2023-09-22)
### Features Added

- New struct `PostBuildDefinition`
- New struct `SubstituteFromDefinition`
- New field `ReconciliationWaitDuration`, `WaitForReconciliation` in struct `FluxConfigurationProperties`
- New field `PostBuild`, `Wait` in struct `KustomizationDefinition`
- New field `PostBuild`, `Wait` in struct `KustomizationPatchDefinition`


## 2.0.0 (2023-05-26)
### Breaking Changes

- Field `InstalledVersion` of struct `ExtensionProperties` has been removed

### Features Added

- New value `SourceKindTypeAzureBlob` added to enum type `SourceKindType`
- New struct `AzureBlobDefinition`
- New struct `AzureBlobPatchDefinition`
- New struct `ManagedIdentityDefinition`
- New struct `ManagedIdentityPatchDefinition`
- New struct `Plan`
- New struct `ServicePrincipalDefinition`
- New struct `ServicePrincipalPatchDefinition`
- New field `Plan` in struct `Extension`
- New field `CurrentVersion`, `IsSystemExtension` in struct `ExtensionProperties`
- New field `AzureBlob` in struct `FluxConfigurationPatchProperties`
- New field `AzureBlob` in struct `FluxConfigurationProperties`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/kubernetesconfiguration/armkubernetesconfiguration` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).