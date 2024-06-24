# Release History

## 1.0.0 (2024-06-28)
### Breaking Changes

- Function `*AvailabilitySetsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, ResourcePatch, *AvailabilitySetsClientBeginUpdateOptions)` to `(context.Context, string, string, AvailabilitySetTagsUpdate, *AvailabilitySetsClientBeginUpdateOptions)`
- Function `*CloudsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, ResourcePatch, *CloudsClientBeginUpdateOptions)` to `(context.Context, string, string, CloudTagsUpdate, *CloudsClientBeginUpdateOptions)`
- Function `*InventoryItemsClient.Create` parameter(s) have been changed from `(context.Context, string, string, string, *InventoryItemsClientCreateOptions)` to `(context.Context, string, string, string, InventoryItem, *InventoryItemsClientCreateOptions)`
- Function `*VirtualMachineTemplatesClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, ResourcePatch, *VirtualMachineTemplatesClientBeginUpdateOptions)` to `(context.Context, string, string, VirtualMachineTemplateTagsUpdate, *VirtualMachineTemplatesClientBeginUpdateOptions)`
- Function `*VirtualNetworksClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, ResourcePatch, *VirtualNetworksClientBeginUpdateOptions)` to `(context.Context, string, string, VirtualNetworkTagsUpdate, *VirtualNetworksClientBeginUpdateOptions)`
- Function `*VmmServersClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, VMMServer, *VmmServersClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, VmmServer, *VmmServersClientBeginCreateOrUpdateOptions)`
- Function `*VmmServersClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, ResourcePatch, *VmmServersClientBeginUpdateOptions)` to `(context.Context, string, string, VmmServerTagsUpdate, *VmmServersClientBeginUpdateOptions)`
- Type of `AvailabilitySetProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `AvailabilitySetsClientBeginDeleteOptions.Force` has been changed from `*bool` to `*ForceDelete`
- Type of `CloudInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `CloudProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `CloudsClientBeginDeleteOptions.Force` has been changed from `*bool` to `*ForceDelete`
- Type of `HardwareProfile.IsHighlyAvailable` has been changed from `*string` to `*IsHighlyAvailable`
- Type of `InventoryItemProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `NetworkProfile.NetworkInterfaces` has been changed from `[]*NetworkInterfaces` to `[]*NetworkInterface`
- Type of `NetworkProfileUpdate.NetworkInterfaces` has been changed from `[]*NetworkInterfacesUpdate` to `[]*NetworkInterfaceUpdate`
- Type of `StopVirtualMachineOptions.SkipShutdown` has been changed from `*bool` to `*SkipShutdown`
- Type of `VirtualMachineInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineTemplateInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineTemplateProperties.IsHighlyAvailable` has been changed from `*string` to `*IsHighlyAvailable`
- Type of `VirtualMachineTemplateProperties.NetworkInterfaces` has been changed from `[]*NetworkInterfaces` to `[]*NetworkInterface`
- Type of `VirtualMachineTemplateProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineTemplatesClientBeginDeleteOptions.Force` has been changed from `*bool` to `*ForceDelete`
- Type of `VirtualNetworkInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualNetworkProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualNetworksClientBeginDeleteOptions.Force` has been changed from `*bool` to `*ForceDelete`
- Type of `VmmServersClientBeginDeleteOptions.Force` has been changed from `*bool` to `*ForceDelete`
- Function `*ClientFactory.NewVirtualMachinesClient` has been removed
- Function `*InventoryItemsClient.NewListByVMMServerPager` has been removed
- Function `NewVirtualMachinesClient` has been removed
- Function `*VirtualMachinesClient.BeginCreateCheckpoint` has been removed
- Function `*VirtualMachinesClient.BeginCreateOrUpdate` has been removed
- Function `*VirtualMachinesClient.BeginDelete` has been removed
- Function `*VirtualMachinesClient.BeginDeleteCheckpoint` has been removed
- Function `*VirtualMachinesClient.Get` has been removed
- Function `*VirtualMachinesClient.NewListByResourceGroupPager` has been removed
- Function `*VirtualMachinesClient.NewListBySubscriptionPager` has been removed
- Function `*VirtualMachinesClient.BeginRestart` has been removed
- Function `*VirtualMachinesClient.BeginRestoreCheckpoint` has been removed
- Function `*VirtualMachinesClient.BeginStart` has been removed
- Function `*VirtualMachinesClient.BeginStop` has been removed
- Function `*VirtualMachinesClient.BeginUpdate` has been removed
- Struct `ErrorDefinition` has been removed
- Struct `InventoryItemsList` has been removed
- Struct `NetworkInterfaces` has been removed
- Struct `NetworkInterfacesUpdate` has been removed
- Struct `OsProfile` has been removed
- Struct `ResourcePatch` has been removed
- Struct `ResourceProviderOperation` has been removed
- Struct `ResourceProviderOperationDisplay` has been removed
- Struct `ResourceProviderOperationList` has been removed
- Struct `StorageQoSPolicy` has been removed
- Struct `StorageQoSPolicyDetails` has been removed
- Struct `VMMServer` has been removed
- Struct `VMMServerListResult` has been removed
- Struct `VMMServerProperties` has been removed
- Struct `VMMServerPropertiesCredentials` has been removed
- Struct `VirtualMachine` has been removed
- Struct `VirtualMachineListResult` has been removed
- Struct `VirtualMachineProperties` has been removed
- Struct `VirtualMachineUpdate` has been removed
- Struct `VirtualMachineUpdateProperties` has been removed
- Field `StorageQoSPolicies` of struct `CloudProperties` has been removed
- Field `Body` of struct `InventoryItemsClientCreateOptions` has been removed
- Field `ResourceProviderOperationList` of struct `OperationsClientListResponse` has been removed
- Field `StorageQoSPolicy` of struct `VirtualDisk` has been removed
- Field `StorageQoSPolicy` of struct `VirtualDiskUpdate` has been removed
- Field `VMMServer` of struct `VmmServersClientCreateOrUpdateResponse` has been removed
- Field `VMMServer` of struct `VmmServersClientGetResponse` has been removed
- Field `VMMServerListResult` of struct `VmmServersClientListByResourceGroupResponse` has been removed
- Field `VMMServerListResult` of struct `VmmServersClientListBySubscriptionResponse` has been removed
- Field `VMMServer` of struct `VmmServersClientUpdateResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `DeleteFromHost` with values `DeleteFromHostFalse`, `DeleteFromHostTrue`
- New enum type `ForceDelete` with values `ForceDeleteFalse`, `ForceDeleteTrue`
- New enum type `IsHighlyAvailable` with values `IsHighlyAvailableFalse`, `IsHighlyAvailableTrue`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ProvisioningAction` with values `ProvisioningActionInstall`, `ProvisioningActionRepair`, `ProvisioningActionUninstall`
- New enum type `ProvisioningState` with values `ProvisioningStateAccepted`, `ProvisioningStateCanceled`, `ProvisioningStateCreated`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateProvisioning`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New enum type `SkipShutdown` with values `SkipShutdownFalse`, `SkipShutdownTrue`
- New function `*ClientFactory.NewGuestAgentsClient() *GuestAgentsClient`
- New function `*ClientFactory.NewVMInstanceHybridIdentityMetadatasClient() *VMInstanceHybridIdentityMetadatasClient`
- New function `*ClientFactory.NewVirtualMachineInstancesClient() *VirtualMachineInstancesClient`
- New function `NewGuestAgentsClient(azcore.TokenCredential, *arm.ClientOptions) (*GuestAgentsClient, error)`
- New function `*GuestAgentsClient.BeginCreate(context.Context, string, GuestAgent, *GuestAgentsClientBeginCreateOptions) (*runtime.Poller[GuestAgentsClientCreateResponse], error)`
- New function `*GuestAgentsClient.Delete(context.Context, string, *GuestAgentsClientDeleteOptions) (GuestAgentsClientDeleteResponse, error)`
- New function `*GuestAgentsClient.Get(context.Context, string, *GuestAgentsClientGetOptions) (GuestAgentsClientGetResponse, error)`
- New function `*GuestAgentsClient.NewListByVirtualMachineInstancePager(string, *GuestAgentsClientListByVirtualMachineInstanceOptions) *runtime.Pager[GuestAgentsClientListByVirtualMachineInstanceResponse]`
- New function `*InventoryItemsClient.NewListByVmmServerPager(string, string, *InventoryItemsClientListByVmmServerOptions) *runtime.Pager[InventoryItemsClientListByVmmServerResponse]`
- New function `NewVMInstanceHybridIdentityMetadatasClient(azcore.TokenCredential, *arm.ClientOptions) (*VMInstanceHybridIdentityMetadatasClient, error)`
- New function `*VMInstanceHybridIdentityMetadatasClient.Get(context.Context, string, *VMInstanceHybridIdentityMetadatasClientGetOptions) (VMInstanceHybridIdentityMetadatasClientGetResponse, error)`
- New function `*VMInstanceHybridIdentityMetadatasClient.NewListByVirtualMachineInstancePager(string, *VMInstanceHybridIdentityMetadatasClientListByVirtualMachineInstanceOptions) *runtime.Pager[VMInstanceHybridIdentityMetadatasClientListByVirtualMachineInstanceResponse]`
- New function `NewVirtualMachineInstancesClient(azcore.TokenCredential, *arm.ClientOptions) (*VirtualMachineInstancesClient, error)`
- New function `*VirtualMachineInstancesClient.BeginCreateCheckpoint(context.Context, string, VirtualMachineCreateCheckpoint, *VirtualMachineInstancesClientBeginCreateCheckpointOptions) (*runtime.Poller[VirtualMachineInstancesClientCreateCheckpointResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginCreateOrUpdate(context.Context, string, VirtualMachineInstance, *VirtualMachineInstancesClientBeginCreateOrUpdateOptions) (*runtime.Poller[VirtualMachineInstancesClientCreateOrUpdateResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginDelete(context.Context, string, *VirtualMachineInstancesClientBeginDeleteOptions) (*runtime.Poller[VirtualMachineInstancesClientDeleteResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginDeleteCheckpoint(context.Context, string, VirtualMachineDeleteCheckpoint, *VirtualMachineInstancesClientBeginDeleteCheckpointOptions) (*runtime.Poller[VirtualMachineInstancesClientDeleteCheckpointResponse], error)`
- New function `*VirtualMachineInstancesClient.Get(context.Context, string, *VirtualMachineInstancesClientGetOptions) (VirtualMachineInstancesClientGetResponse, error)`
- New function `*VirtualMachineInstancesClient.NewListPager(string, *VirtualMachineInstancesClientListOptions) *runtime.Pager[VirtualMachineInstancesClientListResponse]`
- New function `*VirtualMachineInstancesClient.BeginRestart(context.Context, string, *VirtualMachineInstancesClientBeginRestartOptions) (*runtime.Poller[VirtualMachineInstancesClientRestartResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginRestoreCheckpoint(context.Context, string, VirtualMachineRestoreCheckpoint, *VirtualMachineInstancesClientBeginRestoreCheckpointOptions) (*runtime.Poller[VirtualMachineInstancesClientRestoreCheckpointResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginStart(context.Context, string, *VirtualMachineInstancesClientBeginStartOptions) (*runtime.Poller[VirtualMachineInstancesClientStartResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginStop(context.Context, string, StopVirtualMachineOptions, *VirtualMachineInstancesClientBeginStopOptions) (*runtime.Poller[VirtualMachineInstancesClientStopResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginUpdate(context.Context, string, VirtualMachineInstanceUpdate, *VirtualMachineInstancesClientBeginUpdateOptions) (*runtime.Poller[VirtualMachineInstancesClientUpdateResponse], error)`
- New struct `AvailabilitySetTagsUpdate`
- New struct `CloudTagsUpdate`
- New struct `GuestAgent`
- New struct `GuestAgentListResult`
- New struct `GuestAgentProperties`
- New struct `GuestCredential`
- New struct `HTTPProxyConfiguration`
- New struct `InfrastructureProfile`
- New struct `InfrastructureProfileUpdate`
- New struct `InventoryItemListResult`
- New struct `NetworkInterface`
- New struct `NetworkInterfaceUpdate`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `OsProfileForVMInstance`
- New struct `StorageQosPolicy`
- New struct `StorageQosPolicyDetails`
- New struct `VMInstanceHybridIdentityMetadata`
- New struct `VMInstanceHybridIdentityMetadataListResult`
- New struct `VMInstanceHybridIdentityMetadataProperties`
- New struct `VirtualMachineInstance`
- New struct `VirtualMachineInstanceListResult`
- New struct `VirtualMachineInstanceProperties`
- New struct `VirtualMachineInstanceUpdate`
- New struct `VirtualMachineInstanceUpdateProperties`
- New struct `VirtualMachineTemplateTagsUpdate`
- New struct `VirtualNetworkTagsUpdate`
- New struct `VmmCredential`
- New struct `VmmServer`
- New struct `VmmServerListResult`
- New struct `VmmServerProperties`
- New struct `VmmServerTagsUpdate`
- New field `StorageQosPolicies` in struct `CloudProperties`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New field `StorageQosPolicy` in struct `VirtualDisk`
- New field `StorageQosPolicy` in struct `VirtualDiskUpdate`
- New field `BiosGUID`, `ManagedMachineResourceID`, `OSVersion` in struct `VirtualMachineInventoryItem`
- New anonymous field `VmmServer` in struct `VmmServersClientCreateOrUpdateResponse`
- New anonymous field `VmmServer` in struct `VmmServersClientGetResponse`
- New anonymous field `VmmServerListResult` in struct `VmmServersClientListByResourceGroupResponse`
- New anonymous field `VmmServerListResult` in struct `VmmServersClientListBySubscriptionResponse`
- New anonymous field `VmmServer` in struct `VmmServersClientUpdateResponse`


## 0.4.0 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.3.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.2.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/scvmm/armscvmm` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.2.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).