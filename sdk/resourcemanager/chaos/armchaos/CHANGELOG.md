# Release History

## 0.6.0 (2022-09-30)
### Breaking Changes

- Function `*ExperimentsClient.BeginCreateOrUpdate` has been removed
- Function `*ExperimentsClient.BeginCancel` has been removed
- Struct `ExperimentsClientBeginCancelOptions` has been removed
- Struct `ExperimentsClientBeginCreateOrUpdateOptions` has been removed

### Features Added

- New const `FilterTypeSimple`
- New type alias `FilterType`
- New function `*ExperimentsClient.Cancel(context.Context, string, string, *ExperimentsClientCancelOptions) (ExperimentsClientCancelResponse, error)`
- New function `*ExperimentsClient.CreateOrUpdate(context.Context, string, string, Experiment, *ExperimentsClientCreateOrUpdateOptions) (ExperimentsClientCreateOrUpdateResponse, error)`
- New function `*Filter.GetFilter() *Filter`
- New function `*SimpleFilter.GetFilter() *Filter`
- New function `PossibleFilterTypeValues() []FilterType`
- New struct `CapabilityTypePropertiesRuntimeProperties`
- New struct `ExperimentsClientCancelOptions`
- New struct `ExperimentsClientCreateOrUpdateOptions`
- New struct `Filter`
- New struct `SimpleFilter`
- New struct `SimpleFilterParameters`
- New field `Filter` in struct `Selector`
- New field `Kind` in struct `CapabilityTypeProperties`
- New field `RuntimeProperties` in struct `CapabilityTypeProperties`


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/chaos/armchaos` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).