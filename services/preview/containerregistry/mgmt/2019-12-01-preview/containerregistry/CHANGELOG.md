# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. StorageAccountProperties

#### Removed Struct Fields

1. RegistryProperties.StorageAccount

## Additive Changes

### New Constants

1. NetworkRuleBypassOptions.NetworkRuleBypassOptionsAzureServices
1. NetworkRuleBypassOptions.NetworkRuleBypassOptionsNone

### New Funcs

1. PossibleNetworkRuleBypassOptionsValues() []NetworkRuleBypassOptions

### Struct Changes

#### New Structs

1. OperationLogSpecificationDefinition

#### New Struct Fields

1. OperationServiceSpecificationDefinition.LogSpecifications
1. RegistryProperties.NetworkRuleBypassOptions
1. RegistryPropertiesUpdateParameters.NetworkRuleBypassOptions
