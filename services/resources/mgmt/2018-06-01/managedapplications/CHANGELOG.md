Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `ApplicationsClient.UpdatePreparer` parameter(s) have been changed from `(context.Context, string, string, *Application)` to `(context.Context, string, string, *ApplicationPatchable)`
- Function `ApplicationsClient.Update` parameter(s) have been changed from `(context.Context, string, string, *Application)` to `(context.Context, string, string, *ApplicationPatchable)`
- Function `ApplicationDefinitionsClient.CreateOrUpdateByID` parameter(s) have been changed from `(context.Context, string, ApplicationDefinition)` to `(context.Context, string, string, ApplicationDefinition)`
- Function `ApplicationDefinitionsClient.GetByIDPreparer` parameter(s) have been changed from `(context.Context, string)` to `(context.Context, string, string)`
- Function `NewApplicationListResultPage` parameter(s) have been changed from `(func(context.Context, ApplicationListResult) (ApplicationListResult, error))` to `(ApplicationListResult, func(context.Context, ApplicationListResult) (ApplicationListResult, error))`
- Function `ApplicationDefinitionsClient.DeleteByID` parameter(s) have been changed from `(context.Context, string)` to `(context.Context, string, string)`
- Function `ApplicationDefinitionsClient.GetByID` parameter(s) have been changed from `(context.Context, string)` to `(context.Context, string, string)`
- Function `NewApplicationDefinitionListResultPage` parameter(s) have been changed from `(func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error))` to `(ApplicationDefinitionListResult, func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error))`
- Function `ApplicationDefinitionsClient.CreateOrUpdateByIDPreparer` parameter(s) have been changed from `(context.Context, string, ApplicationDefinition)` to `(context.Context, string, string, ApplicationDefinition)`
- Function `ApplicationDefinitionsClient.DeleteByIDPreparer` parameter(s) have been changed from `(context.Context, string)` to `(context.Context, string, string)`

## New Content

- New function `OperationListResultIterator.Value() Operation`
- New function `BaseClient.ListOperations(context.Context) (OperationListResultPage, error)`
- New function `OperationListResultPage.Values() []Operation`
- New function `*OperationListResultIterator.Next() error`
- New function `BaseClient.ListOperationsComplete(context.Context) (OperationListResultIterator, error)`
- New function `BaseClient.ListOperationsSender(*http.Request) (*http.Response, error)`
- New function `BaseClient.ListOperationsResponder(*http.Response) (OperationListResult, error)`
- New function `*OperationListResultPage.Next() error`
- New function `OperationListResultPage.Response() OperationListResult`
- New function `NewOperationListResultIterator(OperationListResultPage) OperationListResultIterator`
- New function `OperationListResultIterator.NotDone() bool`
- New function `OperationListResult.IsEmpty() bool`
- New function `*OperationListResultIterator.NextWithContext(context.Context) error`
- New function `BaseClient.ListOperationsPreparer(context.Context) (*http.Request, error)`
- New function `OperationListResultPage.NotDone() bool`
- New function `OperationListResultIterator.Response() OperationListResult`
- New function `*OperationListResultPage.NextWithContext(context.Context) error`
- New function `NewOperationListResultPage(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `OperationListResultIterator`
- New struct `OperationListResultPage`
