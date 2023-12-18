# Release History

## 1.0.0 (2023-12-22)
### Breaking Changes

- Type of `AvailabilitySetProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `AvailabilitySetsClientBeginDeleteOptions.Force` has been changed from `*bool` to `*Force`
- Type of `CloudInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `CloudProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `CloudsClientBeginDeleteOptions.Force` has been changed from `*bool` to `*Force`
- Type of `HardwareProfile.IsHighlyAvailable` has been changed from `*string` to `*IsHighlyAvailable`
- Type of `InventoryItemProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `NetworkProfile.NetworkInterfaces` has been changed from `[]*NetworkInterfaces` to `[]*NetworkInterface`
- Type of `NetworkProfileUpdate.NetworkInterfaces` has been changed from `[]*NetworkInterfacesUpdate` to `[]*NetworkInterfaceUpdate`
- Type of `StopVirtualMachineOptions.SkipShutdown` has been changed from `*bool` to `*SkipShutdown`
- Type of `VMMServerProperties.Credentials` has been changed from `*VMMServerPropertiesCredentials` to `*VMMCredential`
- Type of `VMMServerProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineTemplateInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineTemplateProperties.IsHighlyAvailable` has been changed from `*string` to `*IsHighlyAvailable`
- Type of `VirtualMachineTemplateProperties.NetworkInterfaces` has been changed from `[]*NetworkInterfaces` to `[]*NetworkInterface`
- Type of `VirtualMachineTemplateProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualMachineTemplatesClientBeginDeleteOptions.Force` has been changed from `*bool` to `*Force`
- Type of `VirtualNetworkInventoryItem.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualNetworkProperties.ProvisioningState` has been changed from `*string` to `*ProvisioningState`
- Type of `VirtualNetworksClientBeginDeleteOptions.Force` has been changed from `*bool` to `*Force`
- Type of `VmmServersClientBeginDeleteOptions.Force` has been changed from `*bool` to `*Force`
- Function `*ClientFactory.NewVirtualMachinesClient` has been removed
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
- Struct `NetworkInterfaces` has been removed
- Struct `NetworkInterfacesUpdate` has been removed
- Struct `OsProfile` has been removed
- Struct `ResourceProviderOperation` has been removed
- Struct `ResourceProviderOperationDisplay` has been removed
- Struct `ResourceProviderOperationList` has been removed
- Struct `VMMServerPropertiesCredentials` has been removed
- Struct `VirtualMachine` has been removed
- Struct `VirtualMachineListResult` has been removed
- Struct `VirtualMachineProperties` has been removed
- Struct `VirtualMachineUpdate` has been removed
- Struct `VirtualMachineUpdateProperties` has been removed
- Field `ResourceProviderOperationList` of struct `OperationsClientListResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `DeleteFromHost` with values `DeleteFromHostFalse`, `DeleteFromHostTrue`
- New enum type `Force` with values `ForceFalse`, `ForceTrue`
- New enum type `IsHighlyAvailable` with values `IsHighlyAvailableFalse`, `IsHighlyAvailableTrue`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ProvisioningAction` with values `ProvisioningActionInstall`, `ProvisioningActionRepair`, `ProvisioningActionUninstall`
- New enum type `ProvisioningState` with values `ProvisioningStateAccepted`, `ProvisioningStateCanceled`, `ProvisioningStateCreated`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateProvisioning`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New enum type `SkipShutdown` with values `SkipShutdownFalse`, `SkipShutdownTrue`
- New function `*ClientFactory.NewVMInstanceGuestAgentsClient() *VMInstanceGuestAgentsClient`
- New function `*ClientFactory.NewVirtualMachineInstanceHybridIdentityMetadataClient() *VirtualMachineInstanceHybridIdentityMetadataClient`
- New function `*ClientFactory.NewVirtualMachineInstancesClient() *VirtualMachineInstancesClient`
- New function `NewVMInstanceGuestAgentsClient(azcore.TokenCredential, *arm.ClientOptions) (*VMInstanceGuestAgentsClient, error)`
- New function `*VMInstanceGuestAgentsClient.BeginCreate(context.Context, string, *VMInstanceGuestAgentsClientBeginCreateOptions) (*runtime.Poller[VMInstanceGuestAgentsClientCreateResponse], error)`
- New function `*VMInstanceGuestAgentsClient.Delete(context.Context, string, *VMInstanceGuestAgentsClientDeleteOptions) (VMInstanceGuestAgentsClientDeleteResponse, error)`
- New function `*VMInstanceGuestAgentsClient.Get(context.Context, string, *VMInstanceGuestAgentsClientGetOptions) (VMInstanceGuestAgentsClientGetResponse, error)`
- New function `*VMInstanceGuestAgentsClient.NewListPager(string, *VMInstanceGuestAgentsClientListOptions) *runtime.Pager[VMInstanceGuestAgentsClientListResponse]`
- New function `NewVirtualMachineInstanceHybridIdentityMetadataClient(azcore.TokenCredential, *arm.ClientOptions) (*VirtualMachineInstanceHybridIdentityMetadataClient, error)`
- New function `*VirtualMachineInstanceHybridIdentityMetadataClient.Get(context.Context, string, *VirtualMachineInstanceHybridIdentityMetadataClientGetOptions) (VirtualMachineInstanceHybridIdentityMetadataClientGetResponse, error)`
- New function `*VirtualMachineInstanceHybridIdentityMetadataClient.NewListPager(string, *VirtualMachineInstanceHybridIdentityMetadataClientListOptions) *runtime.Pager[VirtualMachineInstanceHybridIdentityMetadataClientListResponse]`
- New function `NewVirtualMachineInstancesClient(azcore.TokenCredential, *arm.ClientOptions) (*VirtualMachineInstancesClient, error)`
- New function `*VirtualMachineInstancesClient.BeginCreateCheckpoint(context.Context, string, *VirtualMachineInstancesClientBeginCreateCheckpointOptions) (*runtime.Poller[VirtualMachineInstancesClientCreateCheckpointResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginCreateOrUpdate(context.Context, string, *VirtualMachineInstancesClientBeginCreateOrUpdateOptions) (*runtime.Poller[VirtualMachineInstancesClientCreateOrUpdateResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginDelete(context.Context, string, *VirtualMachineInstancesClientBeginDeleteOptions) (*runtime.Poller[VirtualMachineInstancesClientDeleteResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginDeleteCheckpoint(context.Context, string, *VirtualMachineInstancesClientBeginDeleteCheckpointOptions) (*runtime.Poller[VirtualMachineInstancesClientDeleteCheckpointResponse], error)`
- New function `*VirtualMachineInstancesClient.Get(context.Context, string, *VirtualMachineInstancesClientGetOptions) (VirtualMachineInstancesClientGetResponse, error)`
- New function `*VirtualMachineInstancesClient.NewListPager(string, *VirtualMachineInstancesClientListOptions) *runtime.Pager[VirtualMachineInstancesClientListResponse]`
- New function `*VirtualMachineInstancesClient.BeginRestart(context.Context, string, *VirtualMachineInstancesClientBeginRestartOptions) (*runtime.Poller[VirtualMachineInstancesClientRestartResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginRestoreCheckpoint(context.Context, string, *VirtualMachineInstancesClientBeginRestoreCheckpointOptions) (*runtime.Poller[VirtualMachineInstancesClientRestoreCheckpointResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginStart(context.Context, string, *VirtualMachineInstancesClientBeginStartOptions) (*runtime.Poller[VirtualMachineInstancesClientStartResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginStop(context.Context, string, *VirtualMachineInstancesClientBeginStopOptions) (*runtime.Poller[VirtualMachineInstancesClientStopResponse], error)`
- New function `*VirtualMachineInstancesClient.BeginUpdate(context.Context, string, *VirtualMachineInstancesClientBeginUpdateOptions) (*runtime.Poller[VirtualMachineInstancesClientUpdateResponse], error)`
- New struct `GuestAgent`
- New struct `GuestAgentList`
- New struct `GuestAgentProperties`
- New struct `GuestCredential`
- New struct `HTTPProxyConfiguration`
- New struct `InfrastructureProfile`
- New struct `InfrastructureProfileUpdate`
- New struct `NetworkInterface`
- New struct `NetworkInterfaceUpdate`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `OsProfileForVMInstance`
- New struct `VMInstanceHybridIdentityMetadata`
- New struct `VMInstanceHybridIdentityMetadataList`
- New struct `VMInstanceHybridIdentityMetadataProperties`
- New struct `VMMCredential`
- New struct `VirtualMachineInstance`
- New struct `VirtualMachineInstanceListResult`
- New struct `VirtualMachineInstanceProperties`
- New struct `VirtualMachineInstanceUpdate`
- New struct `VirtualMachineInstanceUpdateProperties`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New field `BiosGUID`, `ManagedMachineResourceID`, `OSVersion` in struct `VirtualMachineInventoryItem`


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
