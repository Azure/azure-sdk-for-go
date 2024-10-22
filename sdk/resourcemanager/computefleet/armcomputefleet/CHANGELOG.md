# Release History

## 1.0.0 (2024-10-22)
### Breaking Changes

- `NetworkAPIVersion20201101` from enum `NetworkAPIVersion` has been removed

### Features Added

- New value `NetworkAPIVersionV20201101` added to enum type `NetworkAPIVersion`
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