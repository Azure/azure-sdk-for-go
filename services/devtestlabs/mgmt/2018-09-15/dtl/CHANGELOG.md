Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/devtestlabs/resource-manager/readme.md tag: `package-2018-09`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *CustomImagesCreateOrUpdateFuture.Result(CustomImagesClient) (CustomImage, error)
1. *CustomImagesDeleteFuture.Result(CustomImagesClient) (autorest.Response, error)
1. *DisksAttachFuture.Result(DisksClient) (autorest.Response, error)
1. *DisksCreateOrUpdateFuture.Result(DisksClient) (Disk, error)
1. *DisksDeleteFuture.Result(DisksClient) (autorest.Response, error)
1. *DisksDetachFuture.Result(DisksClient) (autorest.Response, error)
1. *EnvironmentsCreateOrUpdateFuture.Result(EnvironmentsClient) (Environment, error)
1. *EnvironmentsDeleteFuture.Result(EnvironmentsClient) (autorest.Response, error)
1. *FormulasCreateOrUpdateFuture.Result(FormulasClient) (Formula, error)
1. *GlobalSchedulesExecuteFuture.Result(GlobalSchedulesClient) (autorest.Response, error)
1. *GlobalSchedulesRetargetFuture.Result(GlobalSchedulesClient) (autorest.Response, error)
1. *LabsClaimAnyVMFuture.Result(LabsClient) (autorest.Response, error)
1. *LabsCreateEnvironmentFuture.Result(LabsClient) (autorest.Response, error)
1. *LabsCreateOrUpdateFuture.Result(LabsClient) (Lab, error)
1. *LabsDeleteFuture.Result(LabsClient) (autorest.Response, error)
1. *LabsExportResourceUsageFuture.Result(LabsClient) (autorest.Response, error)
1. *LabsImportVirtualMachineFuture.Result(LabsClient) (autorest.Response, error)
1. *SchedulesExecuteFuture.Result(SchedulesClient) (autorest.Response, error)
1. *SecretsCreateOrUpdateFuture.Result(SecretsClient) (Secret, error)
1. *ServiceFabricSchedulesExecuteFuture.Result(ServiceFabricSchedulesClient) (autorest.Response, error)
1. *ServiceFabricsCreateOrUpdateFuture.Result(ServiceFabricsClient) (ServiceFabric, error)
1. *ServiceFabricsDeleteFuture.Result(ServiceFabricsClient) (autorest.Response, error)
1. *ServiceFabricsStartFuture.Result(ServiceFabricsClient) (autorest.Response, error)
1. *ServiceFabricsStopFuture.Result(ServiceFabricsClient) (autorest.Response, error)
1. *UsersCreateOrUpdateFuture.Result(UsersClient) (User, error)
1. *UsersDeleteFuture.Result(UsersClient) (autorest.Response, error)
1. *VirtualMachineSchedulesExecuteFuture.Result(VirtualMachineSchedulesClient) (autorest.Response, error)
1. *VirtualMachinesAddDataDiskFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesApplyArtifactsFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesClaimFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesCreateOrUpdateFuture.Result(VirtualMachinesClient) (LabVirtualMachine, error)
1. *VirtualMachinesDeleteFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesDetachDataDiskFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesRedeployFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesResizeFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesRestartFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesStartFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesStopFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesTransferDisksFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesUnClaimFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualNetworksCreateOrUpdateFuture.Result(VirtualNetworksClient) (VirtualNetwork, error)
1. *VirtualNetworksDeleteFuture.Result(VirtualNetworksClient) (autorest.Response, error)

## Struct Changes

### Removed Struct Fields

1. CustomImagesCreateOrUpdateFuture.azure.Future
1. CustomImagesDeleteFuture.azure.Future
1. DisksAttachFuture.azure.Future
1. DisksCreateOrUpdateFuture.azure.Future
1. DisksDeleteFuture.azure.Future
1. DisksDetachFuture.azure.Future
1. EnvironmentsCreateOrUpdateFuture.azure.Future
1. EnvironmentsDeleteFuture.azure.Future
1. FormulasCreateOrUpdateFuture.azure.Future
1. GlobalSchedulesExecuteFuture.azure.Future
1. GlobalSchedulesRetargetFuture.azure.Future
1. LabsClaimAnyVMFuture.azure.Future
1. LabsCreateEnvironmentFuture.azure.Future
1. LabsCreateOrUpdateFuture.azure.Future
1. LabsDeleteFuture.azure.Future
1. LabsExportResourceUsageFuture.azure.Future
1. LabsImportVirtualMachineFuture.azure.Future
1. SchedulesExecuteFuture.azure.Future
1. SecretsCreateOrUpdateFuture.azure.Future
1. ServiceFabricSchedulesExecuteFuture.azure.Future
1. ServiceFabricsCreateOrUpdateFuture.azure.Future
1. ServiceFabricsDeleteFuture.azure.Future
1. ServiceFabricsStartFuture.azure.Future
1. ServiceFabricsStopFuture.azure.Future
1. UsersCreateOrUpdateFuture.azure.Future
1. UsersDeleteFuture.azure.Future
1. VirtualMachineSchedulesExecuteFuture.azure.Future
1. VirtualMachinesAddDataDiskFuture.azure.Future
1. VirtualMachinesApplyArtifactsFuture.azure.Future
1. VirtualMachinesClaimFuture.azure.Future
1. VirtualMachinesCreateOrUpdateFuture.azure.Future
1. VirtualMachinesDeleteFuture.azure.Future
1. VirtualMachinesDetachDataDiskFuture.azure.Future
1. VirtualMachinesRedeployFuture.azure.Future
1. VirtualMachinesResizeFuture.azure.Future
1. VirtualMachinesRestartFuture.azure.Future
1. VirtualMachinesStartFuture.azure.Future
1. VirtualMachinesStopFuture.azure.Future
1. VirtualMachinesTransferDisksFuture.azure.Future
1. VirtualMachinesUnClaimFuture.azure.Future
1. VirtualNetworksCreateOrUpdateFuture.azure.Future
1. VirtualNetworksDeleteFuture.azure.Future

## Struct Changes

### New Struct Fields

1. CustomImagesCreateOrUpdateFuture.Result
1. CustomImagesCreateOrUpdateFuture.azure.FutureAPI
1. CustomImagesDeleteFuture.Result
1. CustomImagesDeleteFuture.azure.FutureAPI
1. DisksAttachFuture.Result
1. DisksAttachFuture.azure.FutureAPI
1. DisksCreateOrUpdateFuture.Result
1. DisksCreateOrUpdateFuture.azure.FutureAPI
1. DisksDeleteFuture.Result
1. DisksDeleteFuture.azure.FutureAPI
1. DisksDetachFuture.Result
1. DisksDetachFuture.azure.FutureAPI
1. EnvironmentsCreateOrUpdateFuture.Result
1. EnvironmentsCreateOrUpdateFuture.azure.FutureAPI
1. EnvironmentsDeleteFuture.Result
1. EnvironmentsDeleteFuture.azure.FutureAPI
1. FormulasCreateOrUpdateFuture.Result
1. FormulasCreateOrUpdateFuture.azure.FutureAPI
1. GlobalSchedulesExecuteFuture.Result
1. GlobalSchedulesExecuteFuture.azure.FutureAPI
1. GlobalSchedulesRetargetFuture.Result
1. GlobalSchedulesRetargetFuture.azure.FutureAPI
1. LabsClaimAnyVMFuture.Result
1. LabsClaimAnyVMFuture.azure.FutureAPI
1. LabsCreateEnvironmentFuture.Result
1. LabsCreateEnvironmentFuture.azure.FutureAPI
1. LabsCreateOrUpdateFuture.Result
1. LabsCreateOrUpdateFuture.azure.FutureAPI
1. LabsDeleteFuture.Result
1. LabsDeleteFuture.azure.FutureAPI
1. LabsExportResourceUsageFuture.Result
1. LabsExportResourceUsageFuture.azure.FutureAPI
1. LabsImportVirtualMachineFuture.Result
1. LabsImportVirtualMachineFuture.azure.FutureAPI
1. SchedulesExecuteFuture.Result
1. SchedulesExecuteFuture.azure.FutureAPI
1. SecretsCreateOrUpdateFuture.Result
1. SecretsCreateOrUpdateFuture.azure.FutureAPI
1. ServiceFabricSchedulesExecuteFuture.Result
1. ServiceFabricSchedulesExecuteFuture.azure.FutureAPI
1. ServiceFabricsCreateOrUpdateFuture.Result
1. ServiceFabricsCreateOrUpdateFuture.azure.FutureAPI
1. ServiceFabricsDeleteFuture.Result
1. ServiceFabricsDeleteFuture.azure.FutureAPI
1. ServiceFabricsStartFuture.Result
1. ServiceFabricsStartFuture.azure.FutureAPI
1. ServiceFabricsStopFuture.Result
1. ServiceFabricsStopFuture.azure.FutureAPI
1. UsersCreateOrUpdateFuture.Result
1. UsersCreateOrUpdateFuture.azure.FutureAPI
1. UsersDeleteFuture.Result
1. UsersDeleteFuture.azure.FutureAPI
1. VirtualMachineSchedulesExecuteFuture.Result
1. VirtualMachineSchedulesExecuteFuture.azure.FutureAPI
1. VirtualMachinesAddDataDiskFuture.Result
1. VirtualMachinesAddDataDiskFuture.azure.FutureAPI
1. VirtualMachinesApplyArtifactsFuture.Result
1. VirtualMachinesApplyArtifactsFuture.azure.FutureAPI
1. VirtualMachinesClaimFuture.Result
1. VirtualMachinesClaimFuture.azure.FutureAPI
1. VirtualMachinesCreateOrUpdateFuture.Result
1. VirtualMachinesCreateOrUpdateFuture.azure.FutureAPI
1. VirtualMachinesDeleteFuture.Result
1. VirtualMachinesDeleteFuture.azure.FutureAPI
1. VirtualMachinesDetachDataDiskFuture.Result
1. VirtualMachinesDetachDataDiskFuture.azure.FutureAPI
1. VirtualMachinesRedeployFuture.Result
1. VirtualMachinesRedeployFuture.azure.FutureAPI
1. VirtualMachinesResizeFuture.Result
1. VirtualMachinesResizeFuture.azure.FutureAPI
1. VirtualMachinesRestartFuture.Result
1. VirtualMachinesRestartFuture.azure.FutureAPI
1. VirtualMachinesStartFuture.Result
1. VirtualMachinesStartFuture.azure.FutureAPI
1. VirtualMachinesStopFuture.Result
1. VirtualMachinesStopFuture.azure.FutureAPI
1. VirtualMachinesTransferDisksFuture.Result
1. VirtualMachinesTransferDisksFuture.azure.FutureAPI
1. VirtualMachinesUnClaimFuture.Result
1. VirtualMachinesUnClaimFuture.azure.FutureAPI
1. VirtualNetworksCreateOrUpdateFuture.Result
1. VirtualNetworksCreateOrUpdateFuture.azure.FutureAPI
1. VirtualNetworksDeleteFuture.Result
1. VirtualNetworksDeleteFuture.azure.FutureAPI
