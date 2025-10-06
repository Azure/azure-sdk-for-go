# Release History

## 1.3.0-beta.2 (2025-10-06)
### Breaking Changes

- Operation `*DpsCertificateClient.List` has supported pagination, use `*DpsCertificateClient.NewListPager` instead.
- Operation `*IotDpsResourceClient.ListPrivateLinkResources` has supported pagination, use `*IotDpsResourceClient.NewListPrivateLinkResourcesPager` instead.
- Field `CertificateName1` of struct `DpsCertificateClientDeleteOptions` has been removed
- Field `CertificateName1` of struct `DpsCertificateClientGenerateVerificationCodeOptions` has been removed
- Field `CertificateName1` of struct `DpsCertificateClientVerifyCertificateOptions` has been removed

### Features Added

- New enum type `DeviceRegistryNamespaceAuthenticationType` with values `DeviceRegistryNamespaceAuthenticationTypeSystemAssigned`, `DeviceRegistryNamespaceAuthenticationTypeUserAssigned`
- New struct `DeviceRegistryNamespaceDescription`
- New field `CertificateName` in struct `DpsCertificateClientDeleteOptions`
- New field `CertificateName` in struct `DpsCertificateClientGenerateVerificationCodeOptions`
- New field `CertificateName` in struct `DpsCertificateClientVerifyCertificateOptions`
- New field `SystemData` in struct `GroupIDInformation`
- New field `DeviceRegistryNamespace` in struct `IotDpsPropertiesDescription`


## 1.3.0-beta.1 (2023-11-30)
### Features Added

- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New struct `ManagedServiceIdentity`
- New struct `UserAssignedIdentity`
- New field `PortalOperationsHostName` in struct `IotDpsPropertiesDescription`
- New field `Identity`, `Resourcegroup`, `Subscriptionid` in struct `ProvisioningServiceDescription`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0-beta.1 (2023-06-23)
### Features Added

- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New struct `ManagedServiceIdentity`
- New struct `UserAssignedIdentity`
- New field `PortalOperationsHostName` in struct `IotDpsPropertiesDescription`
- New field `Identity`, `Resourcegroup`, `Subscriptionid` in struct `ProvisioningServiceDescription`


## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deviceprovisioningservices/armdeviceprovisioningservices` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).