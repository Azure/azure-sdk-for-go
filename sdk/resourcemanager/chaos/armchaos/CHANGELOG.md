# Release History

## 0.7.0 (2023-08-25)
### Breaking Changes

- Type of `ExperimentProperties.Selectors` has been changed from `[]*Selector` to `[]SelectorClassification`
- Type of `TargetReference.Type` has been changed from `*string` to `*TargetReferenceType`
- `SelectorTypePercent`, `SelectorTypeRandom`, `SelectorTypeTag` from enum `SelectorType` has been removed
- Operation `*ExperimentsClient.BeginCancel` has been changed to non-LRO, use `*ExperimentsClient.Cancel` instead.
- Operation `*ExperimentsClient.BeginCreateOrUpdate` has been changed to non-LRO, use `*ExperimentsClient.CreateOrUpdate` instead.
- Field `Targets` of struct `Selector` has been removed

### Features Added

- New value `ResourceIdentityTypeUserAssigned` added to enum type `ResourceIdentityType`
- New value `SelectorTypeQuery` added to enum type `SelectorType`
- New enum type `FilterType` with values `FilterTypeSimple`
- New enum type `TargetReferenceType` with values `TargetReferenceTypeChaosTarget`
- New function `*ExperimentsClient.Update(context.Context, string, string, ExperimentUpdate, *ExperimentsClientUpdateOptions) (ExperimentsClientUpdateResponse, error)`
- New function `*Filter.GetFilter() *Filter`
- New function `*ListSelector.GetSelector() *Selector`
- New function `*QuerySelector.GetSelector() *Selector`
- New function `*Selector.GetSelector() *Selector`
- New function `*SimpleFilter.GetFilter() *Filter`
- New struct `CapabilityTypePropertiesRuntimeProperties`
- New struct `ExperimentUpdate`
- New struct `ListSelector`
- New struct `QuerySelector`
- New struct `SimpleFilter`
- New struct `SimpleFilterParameters`
- New struct `UserAssignedIdentity`
- New field `AzureRbacActions`, `AzureRbacDataActions`, `Kind`, `RuntimeProperties` in struct `CapabilityTypeProperties`
- New field `UserAssignedIdentities` in struct `ResourceIdentity`


## 0.6.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 0.6.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/chaos/armchaos` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).