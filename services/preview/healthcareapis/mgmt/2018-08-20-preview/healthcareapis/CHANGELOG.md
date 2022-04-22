# Unreleased

## Additive Changes

### New Constants

1. ManagedServiceIdentityType.None
1. ManagedServiceIdentityType.SystemAssigned

### New Funcs

1. PossibleManagedServiceIdentityTypeValues() []ManagedServiceIdentityType
1. ResourceIdentity.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ResourceIdentity

#### New Struct Fields

1. Operation.IsDataAction
1. Operation.Properties
1. Resource.Identity
1. ServicesDescription.Identity
