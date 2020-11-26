
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewZoneListResultPage` signature has been changed from `(func(context.Context, ZoneListResult) (ZoneListResult, error))` to `(ZoneListResult,func(context.Context, ZoneListResult) (ZoneListResult, error))`
- Function `NewRecordSetListResultPage` signature has been changed from `(func(context.Context, RecordSetListResult) (RecordSetListResult, error))` to `(RecordSetListResult,func(context.Context, RecordSetListResult) (RecordSetListResult, error))`

## New Content

- Field `MaxNumberOfRecordsPerRecordSet` is added to struct `ZoneProperties`

