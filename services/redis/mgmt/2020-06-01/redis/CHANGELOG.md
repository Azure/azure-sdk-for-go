# Unreleased

## Breaking Changes

### Signature Changes

#### Struct Fields

1. CommonProperties.RedisConfiguration changed type from map[string]*string to *CommonPropertiesRedisConfiguration
1. CreateProperties.RedisConfiguration changed type from map[string]*string to *CommonPropertiesRedisConfiguration
1. Properties.RedisConfiguration changed type from map[string]*string to *CommonPropertiesRedisConfiguration
1. UpdateProperties.RedisConfiguration changed type from map[string]*string to *CommonPropertiesRedisConfiguration

## Additive Changes

### New Funcs

1. *CommonPropertiesRedisConfiguration.UnmarshalJSON([]byte) error
1. CommonPropertiesRedisConfiguration.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. CommonPropertiesRedisConfiguration
