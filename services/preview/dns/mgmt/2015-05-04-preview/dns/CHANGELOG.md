# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. RecordSet.Location
1. RecordSet.Properties
1. RecordSet.Tags

## Additive Changes

### New Funcs

1. *RecordSet.UnmarshalJSON([]byte) error
1. RecordSetProperties.MarshalJSON() ([]byte, error)
1. ZoneProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. CloudError
1. CloudErrorBody

#### New Struct Fields

1. RecordSet.*RecordSetProperties
1. RecordSetProperties.Fqdn
1. ZoneProperties.MaxNumberOfRecordsPerRecordSet
