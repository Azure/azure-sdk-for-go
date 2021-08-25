# Unreleased

## Breaking Changes

### Removed Constants

1. NameAvailability.False
1. NameAvailability.True

### Removed Funcs

1. PossibleNameAvailabilityValues() []NameAvailability
1. RemoteRenderingAccountIdentity.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. RemoteRenderingAccountIdentity

### Signature Changes

#### Struct Fields

1. CheckNameAvailabilityResponse.NameAvailable changed type from NameAvailability to *bool
1. RemoteRenderingAccount.Identity changed type from *RemoteRenderingAccountIdentity to *Identity

## Additive Changes

### New Constants

1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User

### New Funcs

1. PossibleCreatedByTypeValues() []CreatedByType

### Struct Changes

#### New Structs

1. LogSpecification
1. MetricDimension
1. MetricSpecification
1. OperationProperties
1. ServiceSpecification
1. SystemData

#### New Struct Fields

1. AccountProperties.StorageAccountName
1. Operation.IsDataAction
1. Operation.Origin
1. Operation.Properties
1. RemoteRenderingAccount.Kind
1. RemoteRenderingAccount.Plan
1. RemoteRenderingAccount.Sku
1. RemoteRenderingAccount.SystemData
1. SpatialAnchorsAccount.Identity
1. SpatialAnchorsAccount.Kind
1. SpatialAnchorsAccount.Plan
1. SpatialAnchorsAccount.Sku
1. SpatialAnchorsAccount.SystemData
