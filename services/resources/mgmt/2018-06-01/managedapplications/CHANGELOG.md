
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewApplicationListResultPage` signature has been changed from `(func(context.Context, ApplicationListResult) (ApplicationListResult, error))` to `(ApplicationListResult,func(context.Context, ApplicationListResult) (ApplicationListResult, error))`
- Function `ApplicationDefinitionsClient.CreateOrUpdateByIDPreparer` signature has been changed from `(context.Context,string,ApplicationDefinition)` to `(context.Context,string,string,ApplicationDefinition)`
- Function `NewApplicationDefinitionListResultPage` signature has been changed from `(func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error))` to `(ApplicationDefinitionListResult,func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error))`
- Function `ApplicationsClient.UpdatePreparer` signature has been changed from `(context.Context,string,string,*Application)` to `(context.Context,string,string,*ApplicationPatchable)`
- Function `ApplicationDefinitionsClient.DeleteByIDPreparer` signature has been changed from `(context.Context,string)` to `(context.Context,string,string)`
- Function `ApplicationDefinitionsClient.GetByIDPreparer` signature has been changed from `(context.Context,string)` to `(context.Context,string,string)`
- Function `ApplicationDefinitionsClient.CreateOrUpdateByID` signature has been changed from `(context.Context,string,ApplicationDefinition)` to `(context.Context,string,string,ApplicationDefinition)`
- Function `ApplicationDefinitionsClient.DeleteByID` signature has been changed from `(context.Context,string)` to `(context.Context,string,string)`
- Function `ApplicationsClient.Update` signature has been changed from `(context.Context,string,string,*Application)` to `(context.Context,string,string,*ApplicationPatchable)`
- Function `ApplicationDefinitionsClient.GetByID` signature has been changed from `(context.Context,string)` to `(context.Context,string,string)`

## New Content

- Function `BaseClient.ListOperations(context.Context) (OperationListResultPage,error)` is added
- Function `BaseClient.ListOperationsSender(*http.Request) (*http.Response,error)` is added
- Function `*OperationListResultIterator.Next() error` is added
- Function `OperationListResultIterator.Response() OperationListResult` is added
- Function `NewOperationListResultIterator(OperationListResultPage) OperationListResultIterator` is added
- Function `OperationListResultPage.Response() OperationListResult` is added
- Function `OperationListResultIterator.Value() Operation` is added
- Function `OperationListResultPage.NotDone() bool` is added
- Function `*OperationListResultPage.NextWithContext(context.Context) error` is added
- Function `BaseClient.ListOperationsResponder(*http.Response) (OperationListResult,error)` is added
- Function `BaseClient.ListOperationsPreparer(context.Context) (*http.Request,error)` is added
- Function `OperationListResult.IsEmpty() bool` is added
- Function `BaseClient.ListOperationsComplete(context.Context) (OperationListResultIterator,error)` is added
- Function `OperationListResultIterator.NotDone() bool` is added
- Function `*OperationListResultPage.Next() error` is added
- Function `NewOperationListResultPage(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage` is added
- Function `*OperationListResultIterator.NextWithContext(context.Context) error` is added
- Function `OperationListResultPage.Values() []Operation` is added
- Struct `Operation` is added
- Struct `OperationDisplay` is added
- Struct `OperationListResult` is added
- Struct `OperationListResultIterator` is added
- Struct `OperationListResultPage` is added

