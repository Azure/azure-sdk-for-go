# Release History

## 2.0.0-beta.1 (2025-08-21)
### Breaking Changes

- `ManagedServiceIdentityTypeSystemAndUserAssigned` from enum `ManagedServiceIdentityType` has been removed

### Features Added

- New value `ManagedServiceIdentityTypeSystemAssignedUserAssigned` added to enum type `ManagedServiceIdentityType`
- New enum type `CapacityType` with values `CapacityTypeVCPU`, `CapacityTypeVM`
- New enum type `FleetMode` with values `FleetModeInstance`, `FleetModeManaged`
- New enum type `VMOperationStatus` with values `VMOperationStatusCancelFailedStatusUnknown`, `VMOperationStatusCanceled`, `VMOperationStatusCreating`, `VMOperationStatusFailed`, `VMOperationStatusSucceeded`
- New enum type `ZoneDistributionStrategy` with values `ZoneDistributionStrategyBestEffortSingleZone`, `ZoneDistributionStrategyPrioritized`
- New function `*FleetsClient.BeginCancel(context.Context, string, string, *FleetsClientBeginCancelOptions) (*runtime.Poller[FleetsClientCancelResponse], error)`
- New function `*FleetsClient.NewListVirtualMachinesPager(string, string, *FleetsClientListVirtualMachinesOptions) *runtime.Pager[FleetsClientListVirtualMachinesResponse]`
- New struct `VirtualMachine`
- New struct `VirtualMachineListResult`
- New struct `ZoneAllocationPolicy`
- New struct `ZonePreference`
- New field `CapacityType`, `Mode`, `ZoneAllocationPolicy` in struct `FleetProperties`


## 1.0.0 (2024-10-22)
### Breaking Changes

- `NetworkAPIVersion20201101` from enum `NetworkAPIVersion` has been renamed to `NetworkAPIVersionV20201101`

### Features Added

- New enum type `AcceleratorManufacturer` with values `AcceleratorManufacturerAMD`, `AcceleratorManufacturerNvidia`, `AcceleratorManufacturerXilinx`
- New enum type `AcceleratorType` with values `AcceleratorTypeFPGA`, `AcceleratorTypeGPU`
- New enum type `ArchitectureType` with values `ArchitectureTypeARM64`, `ArchitectureTypeX64`
- New enum type `CPUManufacturer` with values `CPUManufacturerAMD`, `CPUManufacturerAmpere`, `CPUManufacturerIntel`, `CPUManufacturerMicrosoft`
- New enum type `LocalStorageDiskType` with values `LocalStorageDiskTypeHDD`, `LocalStorageDiskTypeSSD`
- New enum type `VMAttributeSupport` with values `VMAttributeSupportExcluded`, `VMAttributeSupportIncluded`, `VMAttributeSupportRequired`
- New enum type `VMCategory` with values `VMCategoryComputeOptimized`, `VMCategoryFpgaAccelerated`, `VMCategoryGeneralPurpose`, `VMCategoryGpuAccelerated`, `VMCategoryHighPerformanceCompute`, `VMCategoryMemoryOptimized`, `VMCategoryStorageOptimized`
- New struct `AdditionalCapabilities`
- New struct `AdditionalLocationsProfile`
- New struct `LocationProfile`
- New struct `VMAttributeMinMaxDouble`
- New struct `VMAttributeMinMaxInteger`
- New struct `VMAttributes`
- New field `AdditionalVirtualMachineCapabilities` in struct `ComputeProfile`
- New field `AdditionalLocationsProfile`, `VMAttributes` in struct `FleetProperties`


## 0.1.0 (2024-07-26)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/computefleet/armcomputefleet` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).