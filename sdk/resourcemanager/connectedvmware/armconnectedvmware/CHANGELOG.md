# Release History

## 1.1.1 (2023-12-22)
### Other Changes

- Fixed README.md


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.0.0 (2023-10-27)
### Breaking Changes

- Type of `ClusterInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `ClusterProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `DatastoreInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `GuestAgentProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `HostInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `HostProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `InventoryItemProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `ResourcePoolInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `ResourcePoolProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VCenterProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineTemplateInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineTemplateProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualNetworkInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualNetworkProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Enum `IdentityType` has been removed
- Enum `OsTypeUM` has been removed
- Enum `PatchOperationStartedBy` has been removed
- Enum `PatchOperationStatus` has been removed
- Enum `PatchServiceUsed` has been removed
- Enum `StatusLevelTypes` has been removed
- Enum `StatusTypes` has been removed
- Enum `VMGuestPatchClassificationLinux` has been removed
- Enum `VMGuestPatchClassificationWindows` has been removed
- Enum `VMGuestPatchRebootSetting` has been removed
- Enum `VMGuestPatchRebootStatus` has been removed
- Function `*ClientFactory.NewGuestAgentsClient` has been removed
- Function `*ClientFactory.NewHybridIdentityMetadataClient` has been removed
- Function `*ClientFactory.NewMachineExtensionsClient` has been removed
- Function `*ClientFactory.NewVirtualMachinesClient` has been removed
- Function `NewGuestAgentsClient` has been removed
- Function `*GuestAgentsClient.BeginCreate` has been removed
- Function `*GuestAgentsClient.BeginDelete` has been removed
- Function `*GuestAgentsClient.Get` has been removed
- Function `*GuestAgentsClient.NewListByVMPager` has been removed
- Function `NewHybridIdentityMetadataClient` has been removed
- Function `*HybridIdentityMetadataClient.Create` has been removed
- Function `*HybridIdentityMetadataClient.Delete` has been removed
- Function `*HybridIdentityMetadataClient.Get` has been removed
- Function `*HybridIdentityMetadataClient.NewListByVMPager` has been removed
- Function `NewMachineExtensionsClient` has been removed
- Function `*MachineExtensionsClient.BeginCreateOrUpdate` has been removed
- Function `*MachineExtensionsClient.BeginDelete` has been removed
- Function `*MachineExtensionsClient.Get` has been removed
- Function `*MachineExtensionsClient.NewListPager` has been removed
- Function `*MachineExtensionsClient.BeginUpdate` has been removed
- Function `NewVirtualMachinesClient` has been removed
- Function `*VirtualMachinesClient.BeginAssessPatches` has been removed
- Function `*VirtualMachinesClient.BeginCreate` has been removed
- Function `*VirtualMachinesClient.BeginDelete` has been removed
- Function `*VirtualMachinesClient.Get` has been removed
- Function `*VirtualMachinesClient.BeginInstallPatches` has been removed
- Function `*VirtualMachinesClient.NewListByResourceGroupPager` has been removed
- Function `*VirtualMachinesClient.NewListPager` has been removed
- Function `*VirtualMachinesClient.BeginRestart` has been removed
- Function `*VirtualMachinesClient.BeginStart` has been removed
- Function `*VirtualMachinesClient.BeginStop` has been removed
- Function `*VirtualMachinesClient.BeginUpdate` has been removed
- Struct `AvailablePatchCountByClassification` has been removed
- Struct `ErrorDetail` has been removed
- Struct `GuestAgentProfile` has been removed
- Struct `HybridIdentityMetadata` has been removed
- Struct `HybridIdentityMetadataList` has been removed
- Struct `HybridIdentityMetadataProperties` has been removed
- Struct `Identity` has been removed
- Struct `LinuxParameters` has been removed
- Struct `MachineExtension` has been removed
- Struct `MachineExtensionInstanceViewStatus` has been removed
- Struct `MachineExtensionProperties` has been removed
- Struct `MachineExtensionPropertiesInstanceView` has been removed
- Struct `MachineExtensionUpdate` has been removed
- Struct `MachineExtensionUpdateProperties` has been removed
- Struct `MachineExtensionsListResult` has been removed
- Struct `OsProfile` has been removed
- Struct `OsProfileLinuxConfiguration` has been removed
- Struct `OsProfileUpdate` has been removed
- Struct `OsProfileUpdateLinuxConfiguration` has been removed
- Struct `OsProfileUpdateWindowsConfiguration` has been removed
- Struct `OsProfileWindowsConfiguration` has been removed
- Struct `PatchSettings` has been removed
- Struct `VirtualMachine` has been removed
- Struct `VirtualMachineAssessPatchesResult` has been removed
- Struct `VirtualMachineInstallPatchesParameters` has been removed
- Struct `VirtualMachineInstallPatchesResult` has been removed
- Struct `VirtualMachineProperties` has been removed
- Struct `VirtualMachineUpdate` has been removed
- Struct `VirtualMachineUpdateProperties` has been removed
- Struct `VirtualMachinesList` has been removed
- Struct `WindowsParameters` has been removed

### Features Added

- New function `*ClientFactory.NewVMInstanceGuestAgentsClient() *VMInstanceGuestAgentsClient`
- New function `*ClientFactory.NewVMInstanceHybridIdentityMetadataClient() *VMInstanceHybridIdentityMetadataClient`
- New function `*ClientFactory.NewVirtualMachineInstancesClient() *VirtualMachineInstancesClient`
- New function `NewVMInstanceGuestAgentsClient(azcore.TokenCredential, *arm.ClientOptions) (*VMInstanceGuestAgentsClient, error)`
- New function `*VMInstanceGuestAgentsClient.BeginCreate(context.Context, string, GuestAgent, *VMInstanceGuestAgentsClientBeginCreateOptions) (*runtime.Poller[VMInstanceGuestAgentsClientCreateResponse], error)`
- New function `*VMInstanceGuestAgentsClient.BeginDelete(context.Context, string, *VMInstanceGuestAgentsClientBeginDeleteOptions) (*runtime.Poller[VMInstanceGuestAgentsClientDeleteResponse], error)`
- New function `*VMInstanceGuestAgentsClient.Get(context.Context, string, *VMInstanceGuestAgentsClientGetOptions) (VMInstanceGuestAgentsClientGetResponse, error)`
- New function `*VMInstanceGuestAgentsClient.NewListPager(string, *VMInstanceGuestAgentsClientListOptions) *runtime.Pager[VMInstanceGuestAgentsClientListResponse]`
- New function `NewVMInstanceHybridIdentityMetadataClient(azcore.TokenCredential, *arm.ClientOptions) (*VMInstanceHybridIdentityMetadataClient, error)`
- New function `*VMInstanceHybridIdentityMetadataClient.Get(context.Context, string, *VMInstanceHybridIdentityMetadataClientGetOptions) (VMInstanceHybridIdentityMetadataClientGetResponse, error)`
- New function `*VMInstanceHybridIdentityMetadataClient.NewListPager(string, *VMInstanceHybridIdentityMetadataClientListOptions) *runtime.Pager[VMInstanceHybridIdentityMetadataClientListResponse]`
- New function `NewVirtualMachineInstancesClient(azcore.TokenCredential, *arm.ClientOptions) (*VirtualMachineInstancesClient, error)`
- New function `*VirtualMachineInstancesClient.BeginCreateOrUpdate(context.Context, string, VirtualMachineInstance, *VirtualMachineInstancesClientBeginCreateOrUpdateOptions) (*runtime.Poller[VirtualMachineInstancesClientCreateOrUpdateResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginDelete(context.Context, string, *VirtualMachineInstancesClientBeginDeleteOptions) (*runtime.Poller[VirtualMachineInstancesClientDeleteResponse], error)`
- New function `*VirtualMachineInstancesClient.Get(context.Context, string, *VirtualMachineInstancesClientGetOptions) (VirtualMachineInstancesClientGetResponse, error)`
- New function `*VirtualMachineInstancesClient.NewListPager(string, *VirtualMachineInstancesClientListOptions) *runtime.Pager[VirtualMachineInstancesClientListResponse]`
- New function `*VirtualMachineInstancesClient.BeginRestart(context.Context, string, *VirtualMachineInstancesClientBeginRestartOptions) (*runtime.Poller[VirtualMachineInstancesClientRestartResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginStart(context.Context, string, *VirtualMachineInstancesClientBeginStartOptions) (*runtime.Poller[VirtualMachineInstancesClientStartResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginStop(context.Context, string, *VirtualMachineInstancesClientBeginStopOptions) (*runtime.Poller[VirtualMachineInstancesClientStopResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginUpdate(context.Context, string, VirtualMachineInstanceUpdate, *VirtualMachineInstancesClientBeginUpdateOptions) (*runtime.Poller[VirtualMachineInstancesClientUpdateResponse], error)`
- New struct `InfrastructureProfile`
- New struct `OsProfileForVMInstance`
- New struct `VMInstanceHybridIdentityMetadata`
- New struct `VMInstanceHybridIdentityMetadataList`
- New struct `VMInstanceHybridIdentityMetadataProperties`
- New struct `VirtualMachineInstance`
- New struct `VirtualMachineInstanceProperties`
- New struct `VirtualMachineInstanceUpdate`
- New struct `VirtualMachineInstanceUpdateProperties`
- New struct `VirtualMachineInstancesList`
- New field `TotalCPUMHz`, `TotalMemoryGB`, `UsedCPUMHz`, `UsedMemoryGB` in struct `ClusterProperties`
- New field `CapacityGB`, `FreeSpaceGB` in struct `DatastoreProperties`
- New field `PrivateLinkScopeResourceID` in struct `GuestAgentProperties`
- New field `CPUMhz`, `DatastoreIDs`, `MemorySizeGB`, `NetworkIDs`, `OverallCPUUsageMHz`, `OverallMemoryUsageGB` in struct `HostProperties`
- New field `InventoryType` in struct `InventoryItemDetails`
- New field `CPUCapacityMHz`, `CPUOverallUsageMHz`, `DatastoreIDs`, `MemCapacityGB`, `MemOverallUsageGB`, `NetworkIDs` in struct `ResourcePoolProperties`
- New field `Cluster` in struct `VirtualMachineInventoryItem`
- New field `ToolsVersion`, `ToolsVersionStatus` in struct `VirtualMachineTemplateInventoryItem`


## 0.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.2.0 (2023-03-28)
### Breaking Changes

- Struct `Condition` has been removed
- Struct `MachineExtensionInstanceView` has been removed

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.1.0 (2022-08-09)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/connectedvmware/armconnectedvmware` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).