# Release History

## 2.0.0-beta.2 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0-beta.1 (2023-10-27)
### Breaking Changes

- `StatusConnectedRecently`, `StatusDisconnected`, `StatusError`, `StatusNotConnectedRecently`, `StatusNotYetRegistered` from enum `Status` has been removed
- Enum `ArcSettingAggregateState` has been removed
- Enum `DiagnosticLevel` has been removed
- Enum `ExtensionAggregateState` has been removed
- Enum `ImdsAttestation` has been removed
- Enum `NodeArcState` has been removed
- Enum `NodeExtensionState` has been removed
- Enum `ProvisioningState` has been removed
- Enum `WindowsServerSubscription` has been removed
- Function `NewArcSettingsClient` has been removed
- Function `*ArcSettingsClient.Create` has been removed
- Function `*ArcSettingsClient.BeginCreateIdentity` has been removed
- Function `*ArcSettingsClient.BeginDelete` has been removed
- Function `*ArcSettingsClient.GeneratePassword` has been removed
- Function `*ArcSettingsClient.Get` has been removed
- Function `*ArcSettingsClient.NewListByClusterPager` has been removed
- Function `*ArcSettingsClient.Update` has been removed
- Function `*ClientFactory.NewArcSettingsClient` has been removed
- Function `*ClientFactory.NewClustersClient` has been removed
- Function `*ClientFactory.NewExtensionsClient` has been removed
- Function `NewClustersClient` has been removed
- Function `*ClustersClient.Create` has been removed
- Function `*ClustersClient.BeginCreateIdentity` has been removed
- Function `*ClustersClient.BeginDelete` has been removed
- Function `*ClustersClient.Get` has been removed
- Function `*ClustersClient.NewListByResourceGroupPager` has been removed
- Function `*ClustersClient.NewListBySubscriptionPager` has been removed
- Function `*ClustersClient.Update` has been removed
- Function `*ClustersClient.BeginUploadCertificate` has been removed
- Function `NewExtensionsClient` has been removed
- Function `*ExtensionsClient.BeginCreate` has been removed
- Function `*ExtensionsClient.BeginDelete` has been removed
- Function `*ExtensionsClient.Get` has been removed
- Function `*ExtensionsClient.NewListByArcSettingPager` has been removed
- Function `*ExtensionsClient.BeginUpdate` has been removed
- Operation `*OperationsClient.List` has supported pagination, use `*OperationsClient.NewListPager` instead.
- Struct `ArcConnectivityProperties` has been removed
- Struct `ArcIdentityResponse` has been removed
- Struct `ArcIdentityResponseProperties` has been removed
- Struct `ArcSetting` has been removed
- Struct `ArcSettingList` has been removed
- Struct `ArcSettingProperties` has been removed
- Struct `ArcSettingsPatch` has been removed
- Struct `ArcSettingsPatchProperties` has been removed
- Struct `Cluster` has been removed
- Struct `ClusterDesiredProperties` has been removed
- Struct `ClusterIdentityResponse` has been removed
- Struct `ClusterIdentityResponseProperties` has been removed
- Struct `ClusterList` has been removed
- Struct `ClusterNode` has been removed
- Struct `ClusterPatch` has been removed
- Struct `ClusterPatchProperties` has been removed
- Struct `ClusterProperties` has been removed
- Struct `ClusterReportedProperties` has been removed
- Struct `Extension` has been removed
- Struct `ExtensionList` has been removed
- Struct `ExtensionParameters` has been removed
- Struct `ExtensionProperties` has been removed
- Struct `PasswordCredential` has been removed
- Struct `PerNodeExtensionState` has been removed
- Struct `PerNodeState` has been removed
- Struct `RawCertificateData` has been removed
- Struct `UploadCertificateRequest` has been removed

### Features Added

- New value `StatusFailed`, `StatusInProgress`, `StatusSucceeded` added to enum type `Status`
- New enum type `CloudInitDataSource` with values `CloudInitDataSourceAzure`, `CloudInitDataSourceNoCloud`
- New enum type `DiskFileFormat` with values `DiskFileFormatVhd`, `DiskFileFormatVhdx`
- New enum type `ExtendedLocationTypes` with values `ExtendedLocationTypesCustomLocation`
- New enum type `HyperVGeneration` with values `HyperVGenerationV1`, `HyperVGenerationV2`
- New enum type `IPAllocationMethodEnum` with values `IPAllocationMethodEnumDynamic`, `IPAllocationMethodEnumStatic`
- New enum type `IPPoolTypeEnum` with values `IPPoolTypeEnumVM`, `IPPoolTypeEnumVippool`
- New enum type `OperatingSystemTypes` with values `OperatingSystemTypesLinux`, `OperatingSystemTypesWindows`
- New enum type `PowerStateEnum` with values `PowerStateEnumDeallocated`, `PowerStateEnumDeallocating`, `PowerStateEnumRunning`, `PowerStateEnumStarting`, `PowerStateEnumStopped`, `PowerStateEnumStopping`, `PowerStateEnumUnknown`
- New enum type `ProvisioningAction` with values `ProvisioningActionInstall`, `ProvisioningActionRepair`, `ProvisioningActionUninstall`
- New enum type `ProvisioningStateEnum` with values `ProvisioningStateEnumAccepted`, `ProvisioningStateEnumCanceled`, `ProvisioningStateEnumDeleting`, `ProvisioningStateEnumFailed`, `ProvisioningStateEnumInProgress`, `ProvisioningStateEnumSucceeded`
- New enum type `SecurityTypes` with values `SecurityTypesConfidentialVM`, `SecurityTypesTrustedLaunch`
- New enum type `StatusLevelTypes` with values `StatusLevelTypesError`, `StatusLevelTypesInfo`, `StatusLevelTypesWarning`
- New enum type `StatusTypes` with values `StatusTypesFailed`, `StatusTypesInProgress`, `StatusTypesSucceeded`
- New enum type `VMSizeEnum` with values `VMSizeEnumCustom`, `VMSizeEnumDefault`, `VMSizeEnumStandardA2V2`, `VMSizeEnumStandardA4V2`, `VMSizeEnumStandardD16SV3`, `VMSizeEnumStandardD2SV3`, `VMSizeEnumStandardD32SV3`, `VMSizeEnumStandardD4SV3`, `VMSizeEnumStandardD8SV3`, `VMSizeEnumStandardDS13V2`, `VMSizeEnumStandardDS2V2`, `VMSizeEnumStandardDS3V2`, `VMSizeEnumStandardDS4V2`, `VMSizeEnumStandardDS5V2`, `VMSizeEnumStandardK8S2V1`, `VMSizeEnumStandardK8S3V1`, `VMSizeEnumStandardK8S4V1`, `VMSizeEnumStandardK8S5V1`, `VMSizeEnumStandardK8SV1`, `VMSizeEnumStandardNK12`, `VMSizeEnumStandardNK6`, `VMSizeEnumStandardNV12`, `VMSizeEnumStandardNV6`
- New function `*ClientFactory.NewGalleryImagesClient() *GalleryImagesClient`
- New function `*ClientFactory.NewGuestAgentClient() *GuestAgentClient`
- New function `*ClientFactory.NewGuestAgentsClient() *GuestAgentsClient`
- New function `*ClientFactory.NewHybridIdentityMetadataClient() *HybridIdentityMetadataClient`
- New function `*ClientFactory.NewLogicalNetworksClient() *LogicalNetworksClient`
- New function `*ClientFactory.NewMarketplaceGalleryImagesClient() *MarketplaceGalleryImagesClient`
- New function `*ClientFactory.NewNetworkInterfacesClient() *NetworkInterfacesClient`
- New function `*ClientFactory.NewStorageContainersClient() *StorageContainersClient`
- New function `*ClientFactory.NewVirtualHardDisksClient() *VirtualHardDisksClient`
- New function `*ClientFactory.NewVirtualMachineInstancesClient() *VirtualMachineInstancesClient`
- New function `NewGalleryImagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GalleryImagesClient, error)`
- New function `*GalleryImagesClient.BeginCreateOrUpdate(context.Context, string, string, GalleryImages, *GalleryImagesClientBeginCreateOrUpdateOptions) (*runtime.Poller[GalleryImagesClientCreateOrUpdateResponse], error)`
- New function `*GalleryImagesClient.BeginDelete(context.Context, string, string, *GalleryImagesClientBeginDeleteOptions) (*runtime.Poller[GalleryImagesClientDeleteResponse], error)`
- New function `*GalleryImagesClient.Get(context.Context, string, string, *GalleryImagesClientGetOptions) (GalleryImagesClientGetResponse, error)`
- New function `*GalleryImagesClient.NewListAllPager(*GalleryImagesClientListAllOptions) *runtime.Pager[GalleryImagesClientListAllResponse]`
- New function `*GalleryImagesClient.NewListPager(string, *GalleryImagesClientListOptions) *runtime.Pager[GalleryImagesClientListResponse]`
- New function `*GalleryImagesClient.BeginUpdate(context.Context, string, string, GalleryImagesUpdateRequest, *GalleryImagesClientBeginUpdateOptions) (*runtime.Poller[GalleryImagesClientUpdateResponse], error)`
- New function `NewGuestAgentClient(azcore.TokenCredential, *arm.ClientOptions) (*GuestAgentClient, error)`
- New function `*GuestAgentClient.BeginCreate(context.Context, string, *GuestAgentClientBeginCreateOptions) (*runtime.Poller[GuestAgentClientCreateResponse], error)`
- New function `*GuestAgentClient.BeginDelete(context.Context, string, *GuestAgentClientBeginDeleteOptions) (*runtime.Poller[GuestAgentClientDeleteResponse], error)`
- New function `*GuestAgentClient.Get(context.Context, string, *GuestAgentClientGetOptions) (GuestAgentClientGetResponse, error)`
- New function `NewGuestAgentsClient(azcore.TokenCredential, *arm.ClientOptions) (*GuestAgentsClient, error)`
- New function `*GuestAgentsClient.NewListPager(string, *GuestAgentsClientListOptions) *runtime.Pager[GuestAgentsClientListResponse]`
- New function `NewHybridIdentityMetadataClient(azcore.TokenCredential, *arm.ClientOptions) (*HybridIdentityMetadataClient, error)`
- New function `*HybridIdentityMetadataClient.Get(context.Context, string, *HybridIdentityMetadataClientGetOptions) (HybridIdentityMetadataClientGetResponse, error)`
- New function `*HybridIdentityMetadataClient.NewListPager(string, *HybridIdentityMetadataClientListOptions) *runtime.Pager[HybridIdentityMetadataClientListResponse]`
- New function `NewLogicalNetworksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LogicalNetworksClient, error)`
- New function `*LogicalNetworksClient.BeginCreateOrUpdate(context.Context, string, string, LogicalNetworks, *LogicalNetworksClientBeginCreateOrUpdateOptions) (*runtime.Poller[LogicalNetworksClientCreateOrUpdateResponse], error)`
- New function `*LogicalNetworksClient.BeginDelete(context.Context, string, string, *LogicalNetworksClientBeginDeleteOptions) (*runtime.Poller[LogicalNetworksClientDeleteResponse], error)`
- New function `*LogicalNetworksClient.Get(context.Context, string, string, *LogicalNetworksClientGetOptions) (LogicalNetworksClientGetResponse, error)`
- New function `*LogicalNetworksClient.NewListAllPager(*LogicalNetworksClientListAllOptions) *runtime.Pager[LogicalNetworksClientListAllResponse]`
- New function `*LogicalNetworksClient.NewListPager(string, *LogicalNetworksClientListOptions) *runtime.Pager[LogicalNetworksClientListResponse]`
- New function `*LogicalNetworksClient.BeginUpdate(context.Context, string, string, LogicalNetworksUpdateRequest, *LogicalNetworksClientBeginUpdateOptions) (*runtime.Poller[LogicalNetworksClientUpdateResponse], error)`
- New function `NewMarketplaceGalleryImagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MarketplaceGalleryImagesClient, error)`
- New function `*MarketplaceGalleryImagesClient.BeginCreateOrUpdate(context.Context, string, string, MarketplaceGalleryImages, *MarketplaceGalleryImagesClientBeginCreateOrUpdateOptions) (*runtime.Poller[MarketplaceGalleryImagesClientCreateOrUpdateResponse], error)`
- New function `*MarketplaceGalleryImagesClient.BeginDelete(context.Context, string, string, *MarketplaceGalleryImagesClientBeginDeleteOptions) (*runtime.Poller[MarketplaceGalleryImagesClientDeleteResponse], error)`
- New function `*MarketplaceGalleryImagesClient.Get(context.Context, string, string, *MarketplaceGalleryImagesClientGetOptions) (MarketplaceGalleryImagesClientGetResponse, error)`
- New function `*MarketplaceGalleryImagesClient.NewListAllPager(*MarketplaceGalleryImagesClientListAllOptions) *runtime.Pager[MarketplaceGalleryImagesClientListAllResponse]`
- New function `*MarketplaceGalleryImagesClient.NewListPager(string, *MarketplaceGalleryImagesClientListOptions) *runtime.Pager[MarketplaceGalleryImagesClientListResponse]`
- New function `*MarketplaceGalleryImagesClient.BeginUpdate(context.Context, string, string, MarketplaceGalleryImagesUpdateRequest, *MarketplaceGalleryImagesClientBeginUpdateOptions) (*runtime.Poller[MarketplaceGalleryImagesClientUpdateResponse], error)`
- New function `NewStorageContainersClient(string, azcore.TokenCredential, *arm.ClientOptions) (*StorageContainersClient, error)`
- New function `*StorageContainersClient.BeginCreateOrUpdate(context.Context, string, string, StorageContainers, *StorageContainersClientBeginCreateOrUpdateOptions) (*runtime.Poller[StorageContainersClientCreateOrUpdateResponse], error)`
- New function `*StorageContainersClient.BeginDelete(context.Context, string, string, *StorageContainersClientBeginDeleteOptions) (*runtime.Poller[StorageContainersClientDeleteResponse], error)`
- New function `*StorageContainersClient.Get(context.Context, string, string, *StorageContainersClientGetOptions) (StorageContainersClientGetResponse, error)`
- New function `*StorageContainersClient.NewListAllPager(*StorageContainersClientListAllOptions) *runtime.Pager[StorageContainersClientListAllResponse]`
- New function `*StorageContainersClient.NewListPager(string, *StorageContainersClientListOptions) *runtime.Pager[StorageContainersClientListResponse]`
- New function `*StorageContainersClient.BeginUpdate(context.Context, string, string, StorageContainersUpdateRequest, *StorageContainersClientBeginUpdateOptions) (*runtime.Poller[StorageContainersClientUpdateResponse], error)`
- New function `NewVirtualHardDisksClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VirtualHardDisksClient, error)`
- New function `*VirtualHardDisksClient.BeginCreateOrUpdate(context.Context, string, string, VirtualHardDisks, *VirtualHardDisksClientBeginCreateOrUpdateOptions) (*runtime.Poller[VirtualHardDisksClientCreateOrUpdateResponse], error)`
- New function `*VirtualHardDisksClient.BeginDelete(context.Context, string, string, *VirtualHardDisksClientBeginDeleteOptions) (*runtime.Poller[VirtualHardDisksClientDeleteResponse], error)`
- New function `*VirtualHardDisksClient.Get(context.Context, string, string, *VirtualHardDisksClientGetOptions) (VirtualHardDisksClientGetResponse, error)`
- New function `*VirtualHardDisksClient.NewListAllPager(*VirtualHardDisksClientListAllOptions) *runtime.Pager[VirtualHardDisksClientListAllResponse]`
- New function `*VirtualHardDisksClient.NewListPager(string, *VirtualHardDisksClientListOptions) *runtime.Pager[VirtualHardDisksClientListResponse]`
- New function `*VirtualHardDisksClient.BeginUpdate(context.Context, string, string, VirtualHardDisksUpdateRequest, *VirtualHardDisksClientBeginUpdateOptions) (*runtime.Poller[VirtualHardDisksClientUpdateResponse], error)`
- New function `NewVirtualMachineInstancesClient(azcore.TokenCredential, *arm.ClientOptions) (*VirtualMachineInstancesClient, error)`
- New function `*VirtualMachineInstancesClient.BeginCreateOrUpdate(context.Context, string, VirtualMachineInstance, *VirtualMachineInstancesClientBeginCreateOrUpdateOptions) (*runtime.Poller[VirtualMachineInstancesClientCreateOrUpdateResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginDelete(context.Context, string, *VirtualMachineInstancesClientBeginDeleteOptions) (*runtime.Poller[VirtualMachineInstancesClientDeleteResponse], error)`
- New function `*VirtualMachineInstancesClient.Get(context.Context, string, *VirtualMachineInstancesClientGetOptions) (VirtualMachineInstancesClientGetResponse, error)`
- New function `*VirtualMachineInstancesClient.NewListPager(string, *VirtualMachineInstancesClientListOptions) *runtime.Pager[VirtualMachineInstancesClientListResponse]`
- New function `*VirtualMachineInstancesClient.BeginRestart(context.Context, string, *VirtualMachineInstancesClientBeginRestartOptions) (*runtime.Poller[VirtualMachineInstancesClientRestartResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginStart(context.Context, string, *VirtualMachineInstancesClientBeginStartOptions) (*runtime.Poller[VirtualMachineInstancesClientStartResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginStop(context.Context, string, *VirtualMachineInstancesClientBeginStopOptions) (*runtime.Poller[VirtualMachineInstancesClientStopResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginUpdate(context.Context, string, VirtualMachineInstanceUpdateRequest, *VirtualMachineInstancesClientBeginUpdateOptions) (*runtime.Poller[VirtualMachineInstancesClientUpdateResponse], error)`
- New function `NewNetworkInterfacesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NetworkInterfacesClient, error)`
- New function `*NetworkInterfacesClient.BeginCreateOrUpdate(context.Context, string, string, NetworkInterfaces, *NetworkInterfacesClientBeginCreateOrUpdateOptions) (*runtime.Poller[NetworkInterfacesClientCreateOrUpdateResponse], error)`
- New function `*NetworkInterfacesClient.BeginDelete(context.Context, string, string, *NetworkInterfacesClientBeginDeleteOptions) (*runtime.Poller[NetworkInterfacesClientDeleteResponse], error)`
- New function `*NetworkInterfacesClient.Get(context.Context, string, string, *NetworkInterfacesClientGetOptions) (NetworkInterfacesClientGetResponse, error)`
- New function `*NetworkInterfacesClient.NewListAllPager(*NetworkInterfacesClientListAllOptions) *runtime.Pager[NetworkInterfacesClientListAllResponse]`
- New function `*NetworkInterfacesClient.NewListPager(string, *NetworkInterfacesClientListOptions) *runtime.Pager[NetworkInterfacesClientListResponse]`
- New function `*NetworkInterfacesClient.BeginUpdate(context.Context, string, string, NetworkInterfacesUpdateRequest, *NetworkInterfacesClientBeginUpdateOptions) (*runtime.Poller[NetworkInterfacesClientUpdateResponse], error)`
- New struct `ExtendedLocation`
- New struct `GalleryImageIdentifier`
- New struct `GalleryImageProperties`
- New struct `GalleryImageStatus`
- New struct `GalleryImageStatusDownloadStatus`
- New struct `GalleryImageStatusProvisioningStatus`
- New struct `GalleryImageVersion`
- New struct `GalleryImageVersionProperties`
- New struct `GalleryImageVersionStorageProfile`
- New struct `GalleryImages`
- New struct `GalleryImagesListResult`
- New struct `GalleryImagesUpdateRequest`
- New struct `GalleryOSDiskImage`
- New struct `GuestAgent`
- New struct `GuestAgentInstallStatus`
- New struct `GuestAgentList`
- New struct `GuestAgentProperties`
- New struct `GuestCredential`
- New struct `HTTPProxyConfiguration`
- New struct `HardwareProfileUpdate`
- New struct `HybridIdentityMetadata`
- New struct `HybridIdentityMetadataList`
- New struct `HybridIdentityMetadataProperties`
- New struct `IPConfiguration`
- New struct `IPConfigurationProperties`
- New struct `IPConfigurationPropertiesSubnet`
- New struct `IPPool`
- New struct `IPPoolInfo`
- New struct `Identity`
- New struct `InstanceViewStatus`
- New struct `InterfaceDNSSettings`
- New struct `LogicalNetworkProperties`
- New struct `LogicalNetworkPropertiesDhcpOptions`
- New struct `LogicalNetworkStatus`
- New struct `LogicalNetworkStatusProvisioningStatus`
- New struct `LogicalNetworks`
- New struct `LogicalNetworksListResult`
- New struct `LogicalNetworksUpdateRequest`
- New struct `MarketplaceGalleryImageProperties`
- New struct `MarketplaceGalleryImageStatus`
- New struct `MarketplaceGalleryImageStatusDownloadStatus`
- New struct `MarketplaceGalleryImageStatusProvisioningStatus`
- New struct `MarketplaceGalleryImages`
- New struct `MarketplaceGalleryImagesListResult`
- New struct `MarketplaceGalleryImagesUpdateRequest`
- New struct `NetworkInterfaceProperties`
- New struct `NetworkInterfaceStatus`
- New struct `NetworkInterfaceStatusProvisioningStatus`
- New struct `NetworkInterfaces`
- New struct `NetworkInterfacesListResult`
- New struct `NetworkInterfacesUpdateRequest`
- New struct `NetworkProfileUpdate`
- New struct `NetworkProfileUpdateNetworkInterfacesItem`
- New struct `OsProfileUpdate`
- New struct `OsProfileUpdateLinuxConfiguration`
- New struct `OsProfileUpdateWindowsConfiguration`
- New struct `Route`
- New struct `RoutePropertiesFormat`
- New struct `RouteTable`
- New struct `RouteTablePropertiesFormat`
- New struct `SSHConfiguration`
- New struct `SSHPublicKey`
- New struct `StorageContainerProperties`
- New struct `StorageContainerStatus`
- New struct `StorageContainerStatusProvisioningStatus`
- New struct `StorageContainers`
- New struct `StorageContainersListResult`
- New struct `StorageContainersUpdateRequest`
- New struct `StorageProfileUpdate`
- New struct `StorageProfileUpdateDataDisksItem`
- New struct `Subnet`
- New struct `SubnetPropertiesFormat`
- New struct `SubnetPropertiesFormatIPConfigurationReferencesItem`
- New struct `VirtualHardDiskProperties`
- New struct `VirtualHardDiskStatus`
- New struct `VirtualHardDiskStatusProvisioningStatus`
- New struct `VirtualHardDisks`
- New struct `VirtualHardDisksListResult`
- New struct `VirtualHardDisksUpdateRequest`
- New struct `VirtualMachineConfigAgentInstanceView`
- New struct `VirtualMachineInstance`
- New struct `VirtualMachineInstanceListResult`
- New struct `VirtualMachineInstanceProperties`
- New struct `VirtualMachineInstancePropertiesHardwareProfile`
- New struct `VirtualMachineInstancePropertiesHardwareProfileDynamicMemoryConfig`
- New struct `VirtualMachineInstancePropertiesNetworkProfile`
- New struct `VirtualMachineInstancePropertiesNetworkProfileNetworkInterfacesItem`
- New struct `VirtualMachineInstancePropertiesOsProfile`
- New struct `VirtualMachineInstancePropertiesOsProfileLinuxConfiguration`
- New struct `VirtualMachineInstancePropertiesOsProfileWindowsConfiguration`
- New struct `VirtualMachineInstancePropertiesSecurityProfile`
- New struct `VirtualMachineInstancePropertiesSecurityProfileUefiSettings`
- New struct `VirtualMachineInstancePropertiesStorageProfile`
- New struct `VirtualMachineInstancePropertiesStorageProfileDataDisksItem`
- New struct `VirtualMachineInstancePropertiesStorageProfileImageReference`
- New struct `VirtualMachineInstancePropertiesStorageProfileOsDisk`
- New struct `VirtualMachineInstanceStatus`
- New struct `VirtualMachineInstanceStatusProvisioningStatus`
- New struct `VirtualMachineInstanceUpdateProperties`
- New struct `VirtualMachineInstanceUpdateRequest`
- New struct `VirtualMachineInstanceView`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/azurestackhci/armazurestackhci` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).