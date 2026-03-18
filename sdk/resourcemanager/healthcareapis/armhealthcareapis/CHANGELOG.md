# Release History

## 3.0.0-beta.1 (2026-03-18)
### Breaking Changes

- Type of `DicomService.Identity` has been changed from `*ServiceManagedIdentityIdentity` to `*ManagedServiceIdentity`
- Type of `FhirService.Identity` has been changed from `*ServiceManagedIdentityIdentity` to `*ManagedServiceIdentity`
- Type of `IotConnector.Identity` has been changed from `*ServiceManagedIdentityIdentity` to `*ManagedServiceIdentity`
- Struct `Error` has been removed
- Struct `ErrorDetails` has been removed
- Struct `ErrorDetailsInternal` has been removed
- Struct `IotDestinationProperties` has been removed
- Struct `LocationBasedResource` has been removed
- Struct `PrivateEndpointConnectionListResult` has been removed
- Struct `PrivateLinkResource` has been removed
- Struct `Resource` has been removed
- Struct `ResourceCore` has been removed
- Struct `ResourceTags` has been removed
- Struct `ServiceManagedIdentity` has been removed
- Struct `ServicesResource` has been removed
- Struct `TaggedResource` has been removed

### Features Added

- New enum type `ArmManagedServiceIdentityType` with values `ArmManagedServiceIdentityTypeNone`, `ArmManagedServiceIdentityTypeSystemAssigned`, `ArmManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ArmManagedServiceIdentityTypeUserAssigned`
- New struct `ManagedServiceIdentity`
- New struct `StorageIndexingConfiguration`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `NextLink` in struct `PrivateEndpointConnectionListResultDescription`
- New field `NextLink` in struct `PrivateLinkResourceListResultDescription`
- New field `StorageIndexingConfiguration` in struct `StorageConfiguration`


## 2.1.0 (2024-04-26)
### Features Added

- New enum type `SmartDataActions` with values `SmartDataActionsRead`
- New struct `SmartIdentityProviderApplication`
- New struct `SmartIdentityProviderConfiguration`
- New struct `StorageConfiguration`
- New field `EnableDataPartitions`, `StorageConfiguration` in struct `DicomServiceProperties`
- New field `SmartIdentityProviders` in struct `FhirServiceAuthenticationConfiguration`


## 2.0.0 (2023-12-22)
### Breaking Changes

- Struct `FhirServiceAccessPolicyEntry` has been removed
- Field `AccessPolicies` of struct `FhirServiceProperties` has been removed

### Features Added

- New struct `CorsConfiguration`
- New struct `Encryption`
- New struct `EncryptionCustomerManagedKeyEncryption`
- New struct `FhirServiceImportConfiguration`
- New struct `ImplementationGuidesConfiguration`
- New struct `ServiceImportConfigurationInfo`
- New field `CorsConfiguration`, `Encryption`, `EventState` in struct `DicomServiceProperties`
- New field `Encryption`, `ImplementationGuidesConfiguration`, `ImportConfiguration` in struct `FhirServiceProperties`
- New field `EnableRegionalMdmAccount`, `IsInternal`, `MetricFilterPattern`, `ResourceIDDimensionNameOverride`, `SourceMdmAccount` in struct `MetricSpecification`
- New field `CrossTenantCmkApplicationID` in struct `ServiceCosmosDbConfigurationInfo`
- New field `ImportConfiguration` in struct `ServicesProperties`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-04-06)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/healthcareapis/armhealthcareapis` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).