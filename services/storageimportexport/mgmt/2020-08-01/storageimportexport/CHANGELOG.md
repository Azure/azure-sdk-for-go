# Unreleased

## Breaking Changes

### Removed Funcs

1. PossibleKekTypeValues() []KekType
1. PossibleTypeValues() []Type

### Signature Changes

#### Const Types

1. CustomerManaged changed type from KekType to EncryptionKekType
1. MicrosoftManaged changed type from KekType to EncryptionKekType
1. None changed type from Type to IdentityType
1. SystemAssigned changed type from Type to IdentityType
1. UserAssigned changed type from Type to IdentityType

#### Struct Fields

1. EncryptionKeyDetails.KekType changed type from KekType to EncryptionKekType
1. IdentityDetails.Type changed type from Type to IdentityType

## Additive Changes

### New Constants

1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User

### New Funcs

1. PossibleCreatedByTypeValues() []CreatedByType
1. PossibleEncryptionKekTypeValues() []EncryptionKekType
1. PossibleIdentityTypeValues() []IdentityType

### Struct Changes

#### New Structs

1. SystemData

#### New Struct Fields

1. JobResponse.SystemData
1. LocationProperties.AdditionalShippingInformation
