# Change History

## Breaking Changes

### Struct Changes

#### Removed Structs

1. LimitValue

### Signature Changes

#### Struct Fields

1. LimitObject.LimitObjectType changed type from LimitType to LimitObjectType

## Additive Changes

### New Constants

1. LimitObjectType.LimitObjectTypeLimitValue

### New Funcs

1. LimitJSONObject.AsLimitObject() (*LimitObject, bool)
1. LimitObject.AsBasicLimitJSONObject() (BasicLimitJSONObject, bool)
1. LimitObject.AsLimitJSONObject() (*LimitJSONObject, bool)
1. LimitObject.AsLimitObject() (*LimitObject, bool)
1. LimitObject.MarshalJSON() ([]byte, error)
