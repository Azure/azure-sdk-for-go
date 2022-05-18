# Release History

## 0.2.0 (2022-05-18)
### Breaking Changes

- Function `*VirtualMachinesClient.BeginRestart` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientRestartResponse], error)` to `(*runtime.Poller[VirtualMachinesClientRestartResponse], error)`
- Function `*VirtualMachinesClient.BeginRestoreCheckpoint` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientRestoreCheckpointResponse], error)` to `(*runtime.Poller[VirtualMachinesClientRestoreCheckpointResponse], error)`
- Function `*CloudsClient.BeginUpdate` return value(s) have been changed from `(*armruntime.Poller[CloudsClientUpdateResponse], error)` to `(*runtime.Poller[CloudsClientUpdateResponse], error)`
- Function `*VirtualMachineTemplatesClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[VirtualMachineTemplatesClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[VirtualMachineTemplatesClientCreateOrUpdateResponse], error)`
- Function `*VirtualNetworksClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[VirtualNetworksClientDeleteResponse], error)` to `(*runtime.Poller[VirtualNetworksClientDeleteResponse], error)`
- Function `*VmmServersClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[VmmServersClientDeleteResponse], error)` to `(*runtime.Poller[VmmServersClientDeleteResponse], error)`
- Function `*AvailabilitySetsClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[AvailabilitySetsClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[AvailabilitySetsClientCreateOrUpdateResponse], error)`
- Function `*VirtualMachineTemplatesClient.BeginUpdate` return value(s) have been changed from `(*armruntime.Poller[VirtualMachineTemplatesClientUpdateResponse], error)` to `(*runtime.Poller[VirtualMachineTemplatesClientUpdateResponse], error)`
- Function `*VirtualMachinesClient.BeginCreateCheckpoint` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientCreateCheckpointResponse], error)` to `(*runtime.Poller[VirtualMachinesClientCreateCheckpointResponse], error)`
- Function `*VirtualMachinesClient.BeginDeleteCheckpoint` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientDeleteCheckpointResponse], error)` to `(*runtime.Poller[VirtualMachinesClientDeleteCheckpointResponse], error)`
- Function `*VirtualMachinesClient.BeginUpdate` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientUpdateResponse], error)` to `(*runtime.Poller[VirtualMachinesClientUpdateResponse], error)`
- Function `*VirtualMachinesClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[VirtualMachinesClientCreateOrUpdateResponse], error)`
- Function `*VmmServersClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[VmmServersClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[VmmServersClientCreateOrUpdateResponse], error)`
- Function `*VirtualMachinesClient.BeginStop` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientStopResponse], error)` to `(*runtime.Poller[VirtualMachinesClientStopResponse], error)`
- Function `*VirtualNetworksClient.BeginUpdate` return value(s) have been changed from `(*armruntime.Poller[VirtualNetworksClientUpdateResponse], error)` to `(*runtime.Poller[VirtualNetworksClientUpdateResponse], error)`
- Function `*VirtualMachineTemplatesClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[VirtualMachineTemplatesClientDeleteResponse], error)` to `(*runtime.Poller[VirtualMachineTemplatesClientDeleteResponse], error)`
- Function `*CloudsClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[CloudsClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[CloudsClientCreateOrUpdateResponse], error)`
- Function `*AvailabilitySetsClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[AvailabilitySetsClientDeleteResponse], error)` to `(*runtime.Poller[AvailabilitySetsClientDeleteResponse], error)`
- Function `*VirtualNetworksClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[VirtualNetworksClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[VirtualNetworksClientCreateOrUpdateResponse], error)`
- Function `*VirtualMachinesClient.BeginStart` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientStartResponse], error)` to `(*runtime.Poller[VirtualMachinesClientStartResponse], error)`
- Function `*VirtualMachinesClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[VirtualMachinesClientDeleteResponse], error)` to `(*runtime.Poller[VirtualMachinesClientDeleteResponse], error)`
- Function `*VmmServersClient.BeginUpdate` return value(s) have been changed from `(*armruntime.Poller[VmmServersClientUpdateResponse], error)` to `(*runtime.Poller[VmmServersClientUpdateResponse], error)`
- Function `*AvailabilitySetsClient.BeginUpdate` return value(s) have been changed from `(*armruntime.Poller[AvailabilitySetsClientUpdateResponse], error)` to `(*runtime.Poller[AvailabilitySetsClientUpdateResponse], error)`
- Function `*CloudsClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[CloudsClientDeleteResponse], error)` to `(*runtime.Poller[CloudsClientDeleteResponse], error)`
- Function `CloudListResult.MarshalJSON` has been removed
- Function `InventoryItemsList.MarshalJSON` has been removed
- Function `ErrorDefinition.MarshalJSON` has been removed
- Function `VirtualMachineListResult.MarshalJSON` has been removed
- Function `ResourceProviderOperationList.MarshalJSON` has been removed
- Function `VirtualMachineTemplateListResult.MarshalJSON` has been removed
- Function `VMMServerListResult.MarshalJSON` has been removed
- Function `VirtualNetworkListResult.MarshalJSON` has been removed
- Function `AvailabilitySetListResult.MarshalJSON` has been removed


## 0.1.0 (2022-04-21)

- Init release.