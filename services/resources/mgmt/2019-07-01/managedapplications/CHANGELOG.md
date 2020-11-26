
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewApplicationDefinitionListResultPage` signature has been changed from `(func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error))` to `(ApplicationDefinitionListResult,func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error))`
- Function `ApplicationDefinitionsClient.CreateOrUpdateByIDPreparer` signature has been changed from `(context.Context,string,ApplicationDefinition)` to `(context.Context,string,string,ApplicationDefinition)`
- Function `NewApplicationListResultPage` signature has been changed from `(func(context.Context, ApplicationListResult) (ApplicationListResult, error))` to `(ApplicationListResult,func(context.Context, ApplicationListResult) (ApplicationListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `ApplicationDefinitionsClient.GetByIDPreparer` signature has been changed from `(context.Context,string)` to `(context.Context,string,string)`
- Function `ApplicationDefinitionsClient.CreateOrUpdateByID` signature has been changed from `(context.Context,string,ApplicationDefinition)` to `(context.Context,string,string,ApplicationDefinition)`
- Function `ApplicationsClient.Update` signature has been changed from `(context.Context,string,string,*Application)` to `(context.Context,string,string,*ApplicationPatchable)`
- Function `ApplicationDefinitionsClient.DeleteByIDPreparer` signature has been changed from `(context.Context,string)` to `(context.Context,string,string)`
- Function `ApplicationDefinitionsClient.GetByID` signature has been changed from `(context.Context,string)` to `(context.Context,string,string)`
- Function `ApplicationsClient.UpdatePreparer` signature has been changed from `(context.Context,string,string,*Application)` to `(context.Context,string,string,*ApplicationPatchable)`
- Function `ApplicationDefinitionsClient.DeleteByID` signature has been changed from `(context.Context,string)` to `(context.Context,string,string)`
- Field `*ApplicationPropertiesPatchable` of struct `ApplicationPatchable` has been removed

## New Content

- Anonymous field `*ApplicationProperties` is added to struct `ApplicationPatchable`

