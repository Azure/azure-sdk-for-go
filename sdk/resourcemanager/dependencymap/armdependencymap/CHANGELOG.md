# Release History

## 0.2.0 (2025-09-25)
### Features Added

- New enum type `ExportDependenciesStatusCode` with values `ExportDependenciesStatusCodeCompleteMatch`, `ExportDependenciesStatusCodeNoMatch`, `ExportDependenciesStatusCodePartialMatch`
- New function `*MapsClient.BeginGetDependencyViewForAllMachines(context.Context, string, string, GetDependencyViewForAllMachinesRequest, *MapsClientBeginGetDependencyViewForAllMachinesOptions) (*runtime.Poller[MapsClientGetDependencyViewForAllMachinesResponse], error)`
- New struct `DependencyProcessFilter`
- New struct `ErrorAdditionalInfo`
- New struct `ErrorDetail`
- New struct `ExportDependenciesAdditionalInfo`
- New struct `ExportDependenciesOperationResult`
- New struct `ExportDependenciesResultProperties`
- New struct `GetDependencyViewForAllMachinesOperationResult`
- New struct `GetDependencyViewForAllMachinesRequest`
- New struct `GetDependencyViewForAllMachinesResultProperties`
- New field `ApplianceNameList` in struct `ExportDependenciesRequest`
- New anonymous field `ExportDependenciesOperationResult` in struct `MapsClientExportDependenciesResponse`


## 0.1.0 (2025-04-15)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dependencymap/armdependencymap` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).