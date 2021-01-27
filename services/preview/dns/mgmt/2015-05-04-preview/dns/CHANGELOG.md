Generated from https://github.com/Azure/azure-rest-api-specs/tree/../../../../../azure-rest-api-specs/specification/dns/resource-manager/readme.md tag: `package-2015-05-preview`

Code generator 


### Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. RecordSet.Location
1. RecordSet.Properties
1. RecordSet.Tags

#### New Funcs

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
