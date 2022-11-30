# Release History

## 0.2.0 (2022-11-30)
### Breaking Changes

- Function `*VirtualMachinesClient.NewListPager` parameter(s) have been changed from `(*VirtualMachinesClientListOptions)` to `(string, *VirtualMachinesClientListOptions)`
- Function `*GuestAgentsClient.NewListByVMPager` has been removed
- Function `*HybridIdentityMetadataClient.NewListByVMPager` has been removed
- Function `*VirtualMachinesClient.BeginCreate` has been removed
- Function `*VirtualMachinesClient.NewListByResourceGroupPager` has been removed
- Struct `Condition` has been removed
- Struct `ErrorDefinition` has been removed
- Struct `ErrorResponse` has been removed
- Struct `GuestAgentsClientListByVMOptions` has been removed
- Struct `GuestAgentsClientListByVMResponse` has been removed
- Struct `HybridIdentityMetadataClientListByVMOptions` has been removed
- Struct `HybridIdentityMetadataClientListByVMResponse` has been removed
- Struct `MachineExtensionInstanceView` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `VirtualMachinesClientBeginCreateOptions` has been removed
- Struct `VirtualMachinesClientCreateResponse` has been removed
- Struct `VirtualMachinesClientListByResourceGroupOptions` has been removed
- Struct `VirtualMachinesClientListByResourceGroupResponse` has been removed

### Features Added

- New function `NewAzureArcVMwareManagementServiceAPIClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AzureArcVMwareManagementServiceAPIClient, error)`
- New function `*AzureArcVMwareManagementServiceAPIClient.BeginUpgradeExtensions(context.Context, string, string, MachineExtensionUpgrade, *AzureArcVMwareManagementServiceAPIClientBeginUpgradeExtensionsOptions) (*runtime.Poller[AzureArcVMwareManagementServiceAPIClientUpgradeExtensionsResponse], error)`
- New function `*GuestAgentsClient.NewListPager(string, string, *GuestAgentsClientListOptions) *runtime.Pager[GuestAgentsClientListResponse]`
- New function `*HybridIdentityMetadataClient.NewListPager(string, string, *HybridIdentityMetadataClientListOptions) *runtime.Pager[HybridIdentityMetadataClientListResponse]`
- New function `*VirtualMachinesClient.BeginCreateOrUpdate(context.Context, string, string, VirtualMachine, *VirtualMachinesClientBeginCreateOrUpdateOptions) (*runtime.Poller[VirtualMachinesClientCreateOrUpdateResponse], error)`
- New function `*VirtualMachinesClient.NewListAllPager(*VirtualMachinesClientListAllOptions) *runtime.Pager[VirtualMachinesClientListAllResponse]`
- New struct `AzureArcVMwareManagementServiceAPIClient`
- New struct `AzureArcVMwareManagementServiceAPIClientBeginUpgradeExtensionsOptions`
- New struct `AzureArcVMwareManagementServiceAPIClientUpgradeExtensionsResponse`
- New struct `ErrorAdditionalInfo`
- New struct `ExtensionTargetProperties`
- New struct `GuestAgentProfileUpdate`
- New struct `GuestAgentsClientListOptions`
- New struct `GuestAgentsClientListResponse`
- New struct `HybridIdentityMetadataClientListOptions`
- New struct `HybridIdentityMetadataClientListResponse`
- New struct `MachineExtensionUpgrade`
- New struct `VirtualMachinesClientBeginCreateOrUpdateOptions`
- New struct `VirtualMachinesClientCreateOrUpdateResponse`
- New struct `VirtualMachinesClientListAllOptions`
- New struct `VirtualMachinesClientListAllResponse`
- New field `CapacityGB` in struct `DatastoreProperties`
- New field `FreeSpaceGB` in struct `DatastoreProperties`
- New field `AdditionalInfo` in struct `ErrorDetail`
- New field `ClientPublicKey` in struct `GuestAgentProfile`
- New field `MssqlDiscovered` in struct `GuestAgentProfile`
- New field `DatastoreIDs` in struct `HostProperties`
- New field `NetworkIDs` in struct `HostProperties`
- New field `InventoryType` in struct `InventoryItemDetails`
- New field `DatastoreIDs` in struct `ResourcePoolProperties`
- New field `NetworkIDs` in struct `ResourcePoolProperties`
- New field `Cluster` in struct `VirtualMachineInventoryItem`
- New field `ToolsVersion` in struct `VirtualMachineTemplateInventoryItem`
- New field `ToolsVersionStatus` in struct `VirtualMachineTemplateInventoryItem`
- New field `GuestAgentProfile` in struct `VirtualMachineUpdateProperties`


## 0.1.0 (2022-08-09)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/connectedvmware/armconnectedvmware` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).