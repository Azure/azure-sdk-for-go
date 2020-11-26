
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewExemptionListResultPage` signature has been changed from `(func(context.Context, ExemptionListResult) (ExemptionListResult, error))` to `(ExemptionListResult,func(context.Context, ExemptionListResult) (ExemptionListResult, error))`
- Function `NewDefinitionListResultPage` signature has been changed from `(func(context.Context, DefinitionListResult) (DefinitionListResult, error))` to `(DefinitionListResult,func(context.Context, DefinitionListResult) (DefinitionListResult, error))`
- Function `NewSetDefinitionListResultPage` signature has been changed from `(func(context.Context, SetDefinitionListResult) (SetDefinitionListResult, error))` to `(SetDefinitionListResult,func(context.Context, SetDefinitionListResult) (SetDefinitionListResult, error))`
- Function `NewAssignmentListResultPage` signature has been changed from `(func(context.Context, AssignmentListResult) (AssignmentListResult, error))` to `(AssignmentListResult,func(context.Context, AssignmentListResult) (AssignmentListResult, error))`
- Type of `CloudError.Error` has been changed from `*CloudErrorError` to `*ErrorDetail`
- Struct `CloudErrorError` has been removed

## New Content

- Struct `ErrorDetail` is added

