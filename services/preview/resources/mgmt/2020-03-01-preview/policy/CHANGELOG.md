Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewAssignmentListResultPage` parameter(s) have been changed from `(func(context.Context, AssignmentListResult) (AssignmentListResult, error))` to `(AssignmentListResult, func(context.Context, AssignmentListResult) (AssignmentListResult, error))`
- Function `NewSetDefinitionListResultPage` parameter(s) have been changed from `(func(context.Context, SetDefinitionListResult) (SetDefinitionListResult, error))` to `(SetDefinitionListResult, func(context.Context, SetDefinitionListResult) (SetDefinitionListResult, error))`
- Function `NewExemptionListResultPage` parameter(s) have been changed from `(func(context.Context, ExemptionListResult) (ExemptionListResult, error))` to `(ExemptionListResult, func(context.Context, ExemptionListResult) (ExemptionListResult, error))`
- Function `NewDefinitionListResultPage` parameter(s) have been changed from `(func(context.Context, DefinitionListResult) (DefinitionListResult, error))` to `(DefinitionListResult, func(context.Context, DefinitionListResult) (DefinitionListResult, error))`
- Type of `CloudError.Error` has been changed from `*CloudErrorError` to `*ErrorDetail`
- Struct `CloudErrorError` has been removed

## New Content

- New struct `ErrorDetail`
