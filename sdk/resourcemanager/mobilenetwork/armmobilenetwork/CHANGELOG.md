# Release History

## 0.2.0 (2022-03-30)
### Breaking Changes

- Const `ConfigurationStateIncomplete` has been removed
- Const `ConfigurationStateComplete` has been removed
- Function `PossibleConfigurationStateValues` has been removed
- Function `ConfigurationState.ToPtr` has been removed
- Field `ConfigurationState` of struct `SimPropertiesFormat` has been removed

### Features Added

- New const `SimStateEnabled`
- New const `SimStateInvalid`
- New const `SimStateDisabled`
- New function `PossibleSimStateValues() []SimState`
- New function `SimState.ToPtr() *SimState`
- New field `SimState` in struct `SimPropertiesFormat`
- New field `SystemData` in struct `Resource`
- New field `IPv4Address` in struct `InterfaceProperties`
- New field `IPv4Gateway` in struct `InterfaceProperties`
- New field `IPv4Subnet` in struct `InterfaceProperties`
- New field `SystemData` in struct `TrackedResource`


## 0.1.0 (2022-02-28)

- Init release.