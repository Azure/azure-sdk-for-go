# Unreleased

## Breaking Changes

### Removed Funcs

1. *AccountUpdateParameters.UnmarshalJSON([]byte) error
1. *DeletedAccount.UnmarshalJSON([]byte) error
1. *PrivateLinkResource.UnmarshalJSON([]byte) error
1. DeletedAccount.MarshalJSON() ([]byte, error)
1. DeletedAccountProperties.MarshalJSON() ([]byte, error)
1. DeletedAccountPropertiesModel.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. DeletedAccount
1. DeletedAccountList
1. DeletedAccountProperties
1. DeletedAccountPropertiesModel

#### Removed Struct Fields

1. AccountUpdateParameters.*AccountProperties
1. PrivateLinkResource.*PrivateLinkResourceProperties

## Additive Changes

### New Constants

1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User
1. LastModifiedByType.LastModifiedByTypeApplication
1. LastModifiedByType.LastModifiedByTypeKey
1. LastModifiedByType.LastModifiedByTypeManagedIdentity
1. LastModifiedByType.LastModifiedByTypeUser

### New Funcs

1. AccountPropertiesSystemData.MarshalJSON() ([]byte, error)
1. PossibleCreatedByTypeValues() []CreatedByType
1. PossibleLastModifiedByTypeValues() []LastModifiedByType
1. SystemData.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. AccountPropertiesSystemData
1. SystemData

#### New Struct Fields

1. AccountProperties.ManagedResourceGroupName
1. AccountProperties.SystemData
1. AccountUpdateParameters.Properties
1. PrivateLinkResource.Properties
