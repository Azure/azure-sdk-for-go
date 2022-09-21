# Change History

## Breaking Changes

### Removed Funcs

1. ActionGroup.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. ActionGroup

### Signature Changes

#### Struct Fields

1. ActionList.ActionGroups changed type from *[]ActionGroup to *[]ActionGroupForActivityLogAlerts

## Additive Changes

### New Funcs

1. ActionGroupForActivityLogAlerts.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ActionGroupForActivityLogAlerts
