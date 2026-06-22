# Release History

## 4.0.0 (2026-02-11)
### Breaking Changes

- Function `*ApplicationClient.Create` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, accountName string, applicationName string, options *ApplicationClientCreateOptions)` to `(ctx context.Context, resourceGroupName string, accountName string, applicationName string, parameters Application, options *ApplicationClientCreateOptions)`
- Function `*ApplicationPackageClient.Create` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, accountName string, applicationName string, versionName string, options *ApplicationPackageClientCreateOptions)` to `(ctx context.Context, resourceGroupName string, accountName string, applicationName string, versionName string, parameters ApplicationPackage, options *ApplicationPackageClientCreateOptions)`
- Enum `CertificateFormat` has been removed
- Enum `CertificateProvisioningState` has been removed
- Enum `CertificateStoreLocation` has been removed
- Enum `CertificateVisibility` has been removed
- Enum `NodeCommunicationMode` has been removed
- Function `NewCertificateClient` has been removed
- Function `*CertificateClient.CancelDeletion` has been removed
- Function `*CertificateClient.Create` has been removed
- Function `*CertificateClient.BeginDelete` has been removed
- Function `*CertificateClient.Get` has been removed
- Function `*CertificateClient.NewListByBatchAccountPager` has been removed
- Function `*CertificateClient.Update` has been removed
- Function `*ClientFactory.NewCertificateClient` has been removed
- Struct `AzureProxyResource` has been removed
- Struct `AzureResource` has been removed
- Struct `Certificate` has been removed
- Struct `CertificateBaseProperties` has been removed
- Struct `CertificateCreateOrUpdateParameters` has been removed
- Struct `CertificateCreateOrUpdateProperties` has been removed
- Struct `CertificateProperties` has been removed
- Struct `CertificateReference` has been removed
- Struct `DeleteCertificateError` has been removed
- Struct `ErrorAdditionalInfo` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ListCertificatesResult` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Field `Parameters` of struct `ApplicationClientCreateOptions` has been removed
- Field `Parameters` of struct `ApplicationPackageClientCreateOptions` has been removed
- Field `StorageAccountType` of struct `DataDisk` has been removed
- Field `ApplicationLicenses`, `Certificates`, `CurrentNodeCommunicationMode`, `ResourceTags`, `TargetNodeCommunicationMode` of struct `PoolProperties` has been removed

### Features Added

- New value `SecurityEncryptionTypesDiskWithVMGuestState` added to enum type `SecurityEncryptionTypes`
- New enum type `HostEndpointSettingsModeTypes` with values `HostEndpointSettingsModeTypesAudit`, `HostEndpointSettingsModeTypesEnforce`
- New enum type `IPFamily` with values `IPFamilyIPv4`, `IPFamilyIPv6`
- New enum type `JobDefaultOrder` with values `JobDefaultOrderCreationTime`, `JobDefaultOrderNone`
- New struct `DiskCustomerManagedKey`
- New struct `DiskEncryptionSetParameters`
- New struct `HostEndpointSettings`
- New struct `IPTag`
- New struct `ProxyAgentSettings`
- New field `SystemData` in struct `Account`
- New field `SystemData` in struct `Application`
- New field `SystemData` in struct `ApplicationPackage`
- New field `ManagedDisk` in struct `DataDisk`
- New field `SystemData` in struct `DetectorResponse`
- New field `CustomerManagedKey` in struct `DiskEncryptionConfiguration`
- New field `DiskEncryptionSet` in struct `ManagedDisk`
- New field `SystemData` in struct `Pool`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `IPFamilies`, `IPTags` in struct `PublicIPAddressConfiguration`
- New field `ProxyAgentSettings` in struct `SecurityProfile`
- New field `JobDefaultOrder` in struct `TaskSchedulingPolicy`
- New field `DiskEncryptionSet` in struct `VMDiskSecurityProfile`


## 3.0.1 (2025-12-17)
### Other Changes

- Consolidate package version used in live test.

## 3.0.0 (2024-09-27)
### Breaking Changes

- Type of `SecurityProfile.SecurityType` has been changed from `*string` to `*SecurityTypes`
- Function `*LocationClient.NewListSupportedCloudServiceSKUsPager` has been removed
- Struct `CloudServiceConfiguration` has been removed
- Field `CloudServiceConfiguration` of struct `DeploymentConfiguration` has been removed
- Field `Etag` of struct `ProxyResource` has been removed
- Field `Location`, `Tags` of struct `Resource` has been removed

### Features Added

- New value `PublicNetworkAccessTypeSecuredByPerimeter` added to enum type `PublicNetworkAccessType`
- New enum type `AccessRuleDirection` with values `AccessRuleDirectionInbound`, `AccessRuleDirectionOutbound`
- New enum type `ContainerHostDataPath` with values `ContainerHostDataPathApplications`, `ContainerHostDataPathJobPrep`, `ContainerHostDataPathShared`, `ContainerHostDataPathStartup`, `ContainerHostDataPathTask`, `ContainerHostDataPathVfsMounts`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `IssueType` with values `IssueTypeConfigurationPropagationFailure`, `IssueTypeMissingIdentityConfiguration`, `IssueTypeMissingPerimeterConfiguration`, `IssueTypeUnknown`
- New enum type `NetworkSecurityPerimeterConfigurationProvisioningState` with values `NetworkSecurityPerimeterConfigurationProvisioningStateAccepted`, `NetworkSecurityPerimeterConfigurationProvisioningStateCanceled`, `NetworkSecurityPerimeterConfigurationProvisioningStateCreating`, `NetworkSecurityPerimeterConfigurationProvisioningStateDeleting`, `NetworkSecurityPerimeterConfigurationProvisioningStateFailed`, `NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded`, `NetworkSecurityPerimeterConfigurationProvisioningStateUpdating`
- New enum type `ResourceAssociationAccessMode` with values `ResourceAssociationAccessModeAudit`, `ResourceAssociationAccessModeEnforced`, `ResourceAssociationAccessModeLearning`
- New enum type `SecurityEncryptionTypes` with values `SecurityEncryptionTypesNonPersistedTPM`, `SecurityEncryptionTypesVMGuestStateOnly`
- New enum type `SecurityTypes` with values `SecurityTypesConfidentialVM`, `SecurityTypesTrustedLaunch`
- New enum type `Severity` with values `SeverityError`, `SeverityWarning`
- New function `*ClientFactory.NewNetworkSecurityPerimeterClient() *NetworkSecurityPerimeterClient`
- New function `NewNetworkSecurityPerimeterClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkSecurityPerimeterClient, error)`
- New function `*NetworkSecurityPerimeterClient.GetConfiguration(context.Context, string, string, string, *NetworkSecurityPerimeterClientGetConfigurationOptions) (NetworkSecurityPerimeterClientGetConfigurationResponse, error)`
- New function `*NetworkSecurityPerimeterClient.NewListConfigurationsPager(string, string, *NetworkSecurityPerimeterClientListConfigurationsOptions) *runtime.Pager[NetworkSecurityPerimeterClientListConfigurationsResponse]`
- New function `*NetworkSecurityPerimeterClient.BeginReconcileConfiguration(context.Context, string, string, string, *NetworkSecurityPerimeterClientBeginReconcileConfigurationOptions) (*runtime.Poller[NetworkSecurityPerimeterClientReconcileConfigurationResponse], error)`
- New struct `AccessRule`
- New struct `AccessRuleProperties`
- New struct `AccessRulePropertiesSubscriptionsItem`
- New struct `AzureProxyResource`
- New struct `AzureResource`
- New struct `ContainerHostBatchBindMountEntry`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `NetworkSecurityPerimeter`
- New struct `NetworkSecurityPerimeterConfiguration`
- New struct `NetworkSecurityPerimeterConfigurationListResult`
- New struct `NetworkSecurityPerimeterConfigurationProperties`
- New struct `NetworkSecurityProfile`
- New struct `ProvisioningIssue`
- New struct `ProvisioningIssueProperties`
- New struct `ResourceAssociation`
- New struct `SystemData`
- New struct `VMDiskSecurityProfile`
- New field `Tags` in struct `Application`
- New field `Tags` in struct `ApplicationPackage`
- New field `Tags` in struct `Certificate`
- New field `Tags` in struct `CertificateCreateOrUpdateParameters`
- New field `Tags` in struct `DetectorResponse`
- New field `CommunityGalleryImageID`, `SharedGalleryImageID` in struct `ImageReference`
- New field `SecurityProfile` in struct `ManagedDisk`
- New field `Tags` in struct `Pool`
- New field `Tags` in struct `PrivateEndpointConnection`
- New field `Tags` in struct `PrivateLinkResource`
- New field `SystemData` in struct `ProxyResource`
- New field `SystemData` in struct `Resource`
- New field `ContainerHostBatchBindMounts` in struct `TaskContainerSettings`


## 2.3.0 (2024-03-22)
### Features Added

- New enum type `UpgradeMode` with values `UpgradeModeAutomatic`, `UpgradeModeManual`, `UpgradeModeRolling`
- New struct `AutomaticOSUpgradePolicy`
- New struct `RollingUpgradePolicy`
- New struct `UpgradePolicy`
- New field `UpgradePolicy` in struct `PoolProperties`
- New field `BatchSupportEndOfLife` in struct `SupportedSKU`


## 2.2.0 (2023-12-22)
### Features Added

- New value `StorageAccountTypeStandardSSDLRS` added to enum type `StorageAccountType`
- New struct `ManagedDisk`
- New struct `SecurityProfile`
- New struct `ServiceArtifactReference`
- New struct `UefiSettings`
- New field `Caching`, `DiskSizeGB`, `ManagedDisk`, `WriteAcceleratorEnabled` in struct `OSDisk`
- New field `ResourceTags` in struct `PoolProperties`
- New field `SecurityProfile`, `ServiceArtifactReference` in struct `VirtualMachineConfiguration`


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-07-28)
### Breaking Changes

- Type of `ContainerConfiguration.Type` has been changed from `*string` to `*ContainerType`

### Features Added

- New enum type `ContainerType` with values `ContainerTypeCriCompatible`, `ContainerTypeDockerCompatible`
- New field `EnableAcceleratedNetworking` in struct `NetworkConfiguration`
- New field `EnableAutomaticUpgrade` in struct `VMExtension`


## 1.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.2.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0 (2022-11-10)
### Features Added

- New const `PrivateEndpointConnectionProvisioningStateCreating`
- New const `NodeCommunicationModeDefault`
- New const `EndpointAccessDefaultActionAllow`
- New const `NodeCommunicationModeClassic`
- New const `PrivateEndpointConnectionProvisioningStateCancelled`
- New const `EndpointAccessDefaultActionDeny`
- New const `PrivateEndpointConnectionProvisioningStateDeleting`
- New const `NodeCommunicationModeSimplified`
- New type alias `EndpointAccessDefaultAction`
- New type alias `NodeCommunicationMode`
- New function `*PrivateEndpointConnectionClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionClientDeleteResponse], error)`
- New function `PossibleEndpointAccessDefaultActionValues() []EndpointAccessDefaultAction`
- New function `PossibleNodeCommunicationModeValues() []NodeCommunicationMode`
- New struct `EndpointAccessProfile`
- New struct `IPRule`
- New struct `NetworkProfile`
- New struct `PrivateEndpointConnectionClientBeginDeleteOptions`
- New struct `PrivateEndpointConnectionClientDeleteResponse`
- New field `NetworkProfile` in struct `AccountUpdateProperties`
- New field `PublicNetworkAccess` in struct `AccountUpdateProperties`
- New field `NodeManagementEndpoint` in struct `AccountProperties`
- New field `NetworkProfile` in struct `AccountProperties`
- New field `GroupIDs` in struct `PrivateEndpointConnectionProperties`
- New field `NetworkProfile` in struct `AccountCreateProperties`
- New field `TargetNodeCommunicationMode` in struct `PoolProperties`
- New field `CurrentNodeCommunicationMode` in struct `PoolProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/batch/armbatch` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).