# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. Resource.Location
1. Resource.Tags

## Additive Changes

### New Constants

1. EventSerializationType.Avro
1. EventSerializationType.Csv
1. EventSerializationType.JSON

### New Funcs

1. PossibleEventSerializationTypeValues() []EventSerializationType
1. ProxyResource.MarshalJSON() ([]byte, error)
1. TrackedResource.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. Error
1. ErrorDetails
1. ErrorError
1. ProxyResource
1. TrackedResource
