
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `VirtualMachineScaleSetVMExtensionsClient.ListResponder` return values have been changed from `(VirtualMachineExtensionsListResult,error)` to `(VirtualMachineScaleSetVMExtensionsListResult,error)`
- Function `NewVirtualMachineScaleSetListOSUpgradeHistoryPage` signature has been changed from `(func(context.Context, VirtualMachineScaleSetListOSUpgradeHistory) (VirtualMachineScaleSetListOSUpgradeHistory, error))` to `(VirtualMachineScaleSetListOSUpgradeHistory,func(context.Context, VirtualMachineScaleSetListOSUpgradeHistory) (VirtualMachineScaleSetListOSUpgradeHistory, error))`
- Function `NewDedicatedHostGroupListResultPage` signature has been changed from `(func(context.Context, DedicatedHostGroupListResult) (DedicatedHostGroupListResult, error))` to `(DedicatedHostGroupListResult,func(context.Context, DedicatedHostGroupListResult) (DedicatedHostGroupListResult, error))`
- Function `NewGalleryApplicationListPage` signature has been changed from `(func(context.Context, GalleryApplicationList) (GalleryApplicationList, error))` to `(GalleryApplicationList,func(context.Context, GalleryApplicationList) (GalleryApplicationList, error))`
- Function `NewVirtualMachineScaleSetVMListResultPage` signature has been changed from `(func(context.Context, VirtualMachineScaleSetVMListResult) (VirtualMachineScaleSetVMListResult, error))` to `(VirtualMachineScaleSetVMListResult,func(context.Context, VirtualMachineScaleSetVMListResult) (VirtualMachineScaleSetVMListResult, error))`
- Function `VirtualMachineScaleSetVMExtensionsClient.GetResponder` return values have been changed from `(VirtualMachineExtension,error)` to `(VirtualMachineScaleSetVMExtension,error)`
- Function `NewGalleryImageVersionListPage` signature has been changed from `(func(context.Context, GalleryImageVersionList) (GalleryImageVersionList, error))` to `(GalleryImageVersionList,func(context.Context, GalleryImageVersionList) (GalleryImageVersionList, error))`
- Function `NewListUsagesResultPage` signature has been changed from `(func(context.Context, ListUsagesResult) (ListUsagesResult, error))` to `(ListUsagesResult,func(context.Context, ListUsagesResult) (ListUsagesResult, error))`
- Function `*VirtualMachineScaleSetVMExtensionsUpdateFuture.Result` return values have been changed from `(VirtualMachineExtension,error)` to `(VirtualMachineScaleSetVMExtension,error)`
- Function `NewVirtualMachineScaleSetListSkusResultPage` signature has been changed from `(func(context.Context, VirtualMachineScaleSetListSkusResult) (VirtualMachineScaleSetListSkusResult, error))` to `(VirtualMachineScaleSetListSkusResult,func(context.Context, VirtualMachineScaleSetListSkusResult) (VirtualMachineScaleSetListSkusResult, error))`
- Function `NewVirtualMachineListResultPage` signature has been changed from `(func(context.Context, VirtualMachineListResult) (VirtualMachineListResult, error))` to `(VirtualMachineListResult,func(context.Context, VirtualMachineListResult) (VirtualMachineListResult, error))`
- Function `NewVirtualMachineScaleSetListWithLinkResultPage` signature has been changed from `(func(context.Context, VirtualMachineScaleSetListWithLinkResult) (VirtualMachineScaleSetListWithLinkResult, error))` to `(VirtualMachineScaleSetListWithLinkResult,func(context.Context, VirtualMachineScaleSetListWithLinkResult) (VirtualMachineScaleSetListWithLinkResult, error))`
- Function `NewVirtualMachineScaleSetListResultPage` signature has been changed from `(func(context.Context, VirtualMachineScaleSetListResult) (VirtualMachineScaleSetListResult, error))` to `(VirtualMachineScaleSetListResult,func(context.Context, VirtualMachineScaleSetListResult) (VirtualMachineScaleSetListResult, error))`
- Function `NewProximityPlacementGroupListResultPage` signature has been changed from `(func(context.Context, ProximityPlacementGroupListResult) (ProximityPlacementGroupListResult, error))` to `(ProximityPlacementGroupListResult,func(context.Context, ProximityPlacementGroupListResult) (ProximityPlacementGroupListResult, error))`
- Function `NewImageListResultPage` signature has been changed from `(func(context.Context, ImageListResult) (ImageListResult, error))` to `(ImageListResult,func(context.Context, ImageListResult) (ImageListResult, error))`
- Function `NewDiskAccessListPage` signature has been changed from `(func(context.Context, DiskAccessList) (DiskAccessList, error))` to `(DiskAccessList,func(context.Context, DiskAccessList) (DiskAccessList, error))`
- Function `VirtualMachineScaleSetVMExtensionsClient.Update` signature has been changed from `(context.Context,string,string,string,string,VirtualMachineExtensionUpdate)` to `(context.Context,string,string,string,string,VirtualMachineScaleSetVMExtensionUpdate)`
- Function `VirtualMachineScaleSetVMExtensionsClient.Get` return values have been changed from `(VirtualMachineExtension,error)` to `(VirtualMachineScaleSetVMExtension,error)`
- Function `NewVirtualMachineScaleSetExtensionListResultPage` signature has been changed from `(func(context.Context, VirtualMachineScaleSetExtensionListResult) (VirtualMachineScaleSetExtensionListResult, error))` to `(VirtualMachineScaleSetExtensionListResult,func(context.Context, VirtualMachineScaleSetExtensionListResult) (VirtualMachineScaleSetExtensionListResult, error))`
- Function `NewRunCommandListResultPage` signature has been changed from `(func(context.Context, RunCommandListResult) (RunCommandListResult, error))` to `(RunCommandListResult,func(context.Context, RunCommandListResult) (RunCommandListResult, error))`
- Function `NewSnapshotListPage` signature has been changed from `(func(context.Context, SnapshotList) (SnapshotList, error))` to `(SnapshotList,func(context.Context, SnapshotList) (SnapshotList, error))`
- Function `VirtualMachineScaleSetVMExtensionsClient.CreateOrUpdateResponder` return values have been changed from `(VirtualMachineExtension,error)` to `(VirtualMachineScaleSetVMExtension,error)`
- Function `NewGalleryApplicationVersionListPage` signature has been changed from `(func(context.Context, GalleryApplicationVersionList) (GalleryApplicationVersionList, error))` to `(GalleryApplicationVersionList,func(context.Context, GalleryApplicationVersionList) (GalleryApplicationVersionList, error))`
- Function `NewDiskListPage` signature has been changed from `(func(context.Context, DiskList) (DiskList, error))` to `(DiskList,func(context.Context, DiskList) (DiskList, error))`
- Function `NewAvailabilitySetListResultPage` signature has been changed from `(func(context.Context, AvailabilitySetListResult) (AvailabilitySetListResult, error))` to `(AvailabilitySetListResult,func(context.Context, AvailabilitySetListResult) (AvailabilitySetListResult, error))`
- Function `NewVirtualMachineRunCommandsListResultPage` signature has been changed from `(func(context.Context, VirtualMachineRunCommandsListResult) (VirtualMachineRunCommandsListResult, error))` to `(VirtualMachineRunCommandsListResult,func(context.Context, VirtualMachineRunCommandsListResult) (VirtualMachineRunCommandsListResult, error))`
- Function `NewGalleryImageListPage` signature has been changed from `(func(context.Context, GalleryImageList) (GalleryImageList, error))` to `(GalleryImageList,func(context.Context, GalleryImageList) (GalleryImageList, error))`
- Function `VirtualMachineScaleSetVMExtensionsClient.CreateOrUpdate` signature has been changed from `(context.Context,string,string,string,string,VirtualMachineExtension)` to `(context.Context,string,string,string,string,VirtualMachineScaleSetVMExtension)`
- Function `NewResourceSkusResultPage` signature has been changed from `(func(context.Context, ResourceSkusResult) (ResourceSkusResult, error))` to `(ResourceSkusResult,func(context.Context, ResourceSkusResult) (ResourceSkusResult, error))`
- Function `*VirtualMachineScaleSetVMExtensionsCreateOrUpdateFuture.Result` return values have been changed from `(VirtualMachineExtension,error)` to `(VirtualMachineScaleSetVMExtension,error)`
- Function `VirtualMachineScaleSetVMExtensionsClient.List` return values have been changed from `(VirtualMachineExtensionsListResult,error)` to `(VirtualMachineScaleSetVMExtensionsListResult,error)`
- Function `NewGalleryListPage` signature has been changed from `(func(context.Context, GalleryList) (GalleryList, error))` to `(GalleryList,func(context.Context, GalleryList) (GalleryList, error))`
- Function `NewDiskEncryptionSetListPage` signature has been changed from `(func(context.Context, DiskEncryptionSetList) (DiskEncryptionSetList, error))` to `(DiskEncryptionSetList,func(context.Context, DiskEncryptionSetList) (DiskEncryptionSetList, error))`
- Function `VirtualMachineScaleSetVMExtensionsClient.UpdateResponder` return values have been changed from `(VirtualMachineExtension,error)` to `(VirtualMachineScaleSetVMExtension,error)`
- Function `VirtualMachineScaleSetVMExtensionsClient.UpdatePreparer` signature has been changed from `(context.Context,string,string,string,string,VirtualMachineExtensionUpdate)` to `(context.Context,string,string,string,string,VirtualMachineScaleSetVMExtensionUpdate)`
- Function `VirtualMachineScaleSetVMExtensionsClient.CreateOrUpdatePreparer` signature has been changed from `(context.Context,string,string,string,string,VirtualMachineExtension)` to `(context.Context,string,string,string,string,VirtualMachineScaleSetVMExtension)`
- Function `NewSSHPublicKeysGroupListResultPage` signature has been changed from `(func(context.Context, SSHPublicKeysGroupListResult) (SSHPublicKeysGroupListResult, error))` to `(SSHPublicKeysGroupListResult,func(context.Context, SSHPublicKeysGroupListResult) (SSHPublicKeysGroupListResult, error))`
- Function `NewContainerServiceListResultPage` signature has been changed from `(func(context.Context, ContainerServiceListResult) (ContainerServiceListResult, error))` to `(ContainerServiceListResult,func(context.Context, ContainerServiceListResult) (ContainerServiceListResult, error))`
- Function `NewResourceURIListPage` signature has been changed from `(func(context.Context, ResourceURIList) (ResourceURIList, error))` to `(ResourceURIList,func(context.Context, ResourceURIList) (ResourceURIList, error))`
- Function `NewDedicatedHostListResultPage` signature has been changed from `(func(context.Context, DedicatedHostListResult) (DedicatedHostListResult, error))` to `(DedicatedHostListResult,func(context.Context, DedicatedHostListResult) (DedicatedHostListResult, error))`

## New Content

- Function `VirtualMachineScaleSetVMExtensionUpdate.MarshalJSON() ([]byte,error)` is added
- Function `VirtualMachineScaleSetVMExtension.MarshalJSON() ([]byte,error)` is added
- Function `*VirtualMachineScaleSetVMExtension.UnmarshalJSON([]byte) error` is added
- Function `*VirtualMachineScaleSetVMExtensionUpdate.UnmarshalJSON([]byte) error` is added
- Struct `VirtualMachineScaleSetVMExtension` is added
- Struct `VirtualMachineScaleSetVMExtensionUpdate` is added
- Struct `VirtualMachineScaleSetVMExtensionsListResult` is added

