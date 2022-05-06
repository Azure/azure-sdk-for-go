# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. StorageAccountProperties

#### Removed Struct Fields

1. RegistryProperties.StorageAccount

### Signature Changes

#### Struct Fields

1. ErrorResponseBody.Details changed type from *InnerErrorDescription to *[]InnerErrorDescription

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
