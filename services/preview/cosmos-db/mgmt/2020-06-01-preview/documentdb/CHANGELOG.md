
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Struct `ErrorResponseUpdatedFormat` has been removed

## New Content

- Const `BuiltInRole` is added
- Const `CustomRole` is added
- Function `SQLRoleDefinitionCreateUpdateParameters.MarshalJSON() ([]byte,error)` is added
- Function `SQLResourcesClient.CreateUpdateSQLRoleAssignmentPreparer(context.Context,string,string,string,SQLRoleAssignmentCreateUpdateParameters) (*http.Request,error)` is added
- Function `SQLResourcesClient.ListSQLRoleDefinitions(context.Context,string,string) (SQLRoleDefinitionListResult,error)` is added
- Function `SQLResourcesClient.CreateUpdateSQLRoleAssignmentResponder(*http.Response) (SQLRoleAssignmentGetResults,error)` is added
- Function `PossibleRoleDefinitionTypeValues() []RoleDefinitionType` is added
- Function `SQLResourcesClient.GetSQLRoleAssignmentResponder(*http.Response) (SQLRoleAssignmentGetResults,error)` is added
- Function `SQLResourcesClient.GetSQLRoleDefinitionSender(*http.Request) (*http.Response,error)` is added
- Function `SQLResourcesClient.DeleteSQLRoleDefinitionPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `SQLResourcesClient.CreateUpdateSQLRoleAssignment(context.Context,string,string,string,SQLRoleAssignmentCreateUpdateParameters) (SQLResourcesCreateUpdateSQLRoleAssignmentFuture,error)` is added
- Function `SQLResourcesClient.CreateUpdateSQLRoleDefinitionSender(*http.Request) (SQLResourcesCreateUpdateSQLRoleDefinitionFuture,error)` is added
- Function `SQLResourcesClient.DeleteSQLRoleAssignmentPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `SQLRoleAssignmentCreateUpdateParameters.MarshalJSON() ([]byte,error)` is added
- Function `SQLResourcesClient.GetSQLRoleDefinitionResponder(*http.Response) (SQLRoleDefinitionGetResults,error)` is added
- Function `SQLResourcesClient.DeleteSQLRoleDefinitionSender(*http.Request) (SQLResourcesDeleteSQLRoleDefinitionFuture,error)` is added
- Function `SQLResourcesClient.DeleteSQLRoleAssignment(context.Context,string,string,string) (SQLResourcesDeleteSQLRoleAssignmentFuture,error)` is added
- Function `*SQLRoleAssignmentCreateUpdateParameters.UnmarshalJSON([]byte) error` is added
- Function `SQLResourcesClient.DeleteSQLRoleAssignmentSender(*http.Request) (SQLResourcesDeleteSQLRoleAssignmentFuture,error)` is added
- Function `SQLResourcesClient.DeleteSQLRoleDefinitionResponder(*http.Response) (autorest.Response,error)` is added
- Function `*SQLRoleAssignmentGetResults.UnmarshalJSON([]byte) error` is added
- Function `*SQLRoleDefinitionCreateUpdateParameters.UnmarshalJSON([]byte) error` is added
- Function `SQLRoleDefinitionGetResults.MarshalJSON() ([]byte,error)` is added
- Function `SQLResourcesClient.CreateUpdateSQLRoleDefinition(context.Context,string,string,string,SQLRoleDefinitionCreateUpdateParameters) (SQLResourcesCreateUpdateSQLRoleDefinitionFuture,error)` is added
- Function `SQLResourcesClient.DeleteSQLRoleDefinition(context.Context,string,string,string) (SQLResourcesDeleteSQLRoleDefinitionFuture,error)` is added
- Function `*SQLResourcesCreateUpdateSQLRoleAssignmentFuture.Result(SQLResourcesClient) (SQLRoleAssignmentGetResults,error)` is added
- Function `SQLResourcesClient.ListSQLRoleAssignmentsPreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `SQLResourcesClient.ListSQLRoleAssignments(context.Context,string,string) (SQLRoleAssignmentListResult,error)` is added
- Function `SQLResourcesClient.CreateUpdateSQLRoleAssignmentSender(*http.Request) (SQLResourcesCreateUpdateSQLRoleAssignmentFuture,error)` is added
- Function `SQLResourcesClient.CreateUpdateSQLRoleDefinitionResponder(*http.Response) (SQLRoleDefinitionGetResults,error)` is added
- Function `SQLResourcesClient.ListSQLRoleDefinitionsPreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `SQLResourcesClient.ListSQLRoleDefinitionsResponder(*http.Response) (SQLRoleDefinitionListResult,error)` is added
- Function `SQLResourcesClient.GetSQLRoleAssignment(context.Context,string,string,string) (SQLRoleAssignmentGetResults,error)` is added
- Function `SQLResourcesClient.GetSQLRoleAssignmentSender(*http.Request) (*http.Response,error)` is added
- Function `SQLResourcesClient.DeleteSQLRoleAssignmentResponder(*http.Response) (autorest.Response,error)` is added
- Function `SQLResourcesClient.GetSQLRoleAssignmentPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `SQLResourcesClient.CreateUpdateSQLRoleDefinitionPreparer(context.Context,string,string,string,SQLRoleDefinitionCreateUpdateParameters) (*http.Request,error)` is added
- Function `SQLResourcesClient.GetSQLRoleDefinitionPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `SQLResourcesClient.ListSQLRoleAssignmentsResponder(*http.Response) (SQLRoleAssignmentListResult,error)` is added
- Function `SQLRoleAssignmentGetResults.MarshalJSON() ([]byte,error)` is added
- Function `*SQLResourcesDeleteSQLRoleDefinitionFuture.Result(SQLResourcesClient) (autorest.Response,error)` is added
- Function `*SQLResourcesDeleteSQLRoleAssignmentFuture.Result(SQLResourcesClient) (autorest.Response,error)` is added
- Function `SQLResourcesClient.GetSQLRoleDefinition(context.Context,string,string,string) (SQLRoleDefinitionGetResults,error)` is added
- Function `SQLResourcesClient.ListSQLRoleDefinitionsSender(*http.Request) (*http.Response,error)` is added
- Function `SQLResourcesClient.ListSQLRoleAssignmentsSender(*http.Request) (*http.Response,error)` is added
- Function `*SQLRoleDefinitionGetResults.UnmarshalJSON([]byte) error` is added
- Function `*SQLResourcesCreateUpdateSQLRoleDefinitionFuture.Result(SQLResourcesClient) (SQLRoleDefinitionGetResults,error)` is added
- Struct `CorsPolicy` is added
- Struct `DefaultErrorResponse` is added
- Struct `ManagedServiceIdentityUserAssignedIdentitiesValue` is added
- Struct `Permission` is added
- Struct `SQLResourcesCreateUpdateSQLRoleAssignmentFuture` is added
- Struct `SQLResourcesCreateUpdateSQLRoleDefinitionFuture` is added
- Struct `SQLResourcesDeleteSQLRoleAssignmentFuture` is added
- Struct `SQLResourcesDeleteSQLRoleDefinitionFuture` is added
- Struct `SQLRoleAssignmentCreateUpdateParameters` is added
- Struct `SQLRoleAssignmentGetResults` is added
- Struct `SQLRoleAssignmentListResult` is added
- Struct `SQLRoleAssignmentResource` is added
- Struct `SQLRoleDefinitionCreateUpdateParameters` is added
- Struct `SQLRoleDefinitionGetResults` is added
- Struct `SQLRoleDefinitionListResult` is added
- Struct `SQLRoleDefinitionResource` is added
- Field `Cors` is added to struct `RestoreReqeustDatabaseAccountCreateUpdateProperties`
- Field `Cors` is added to struct `DatabaseAccountGetProperties`
- Field `UserAssignedIdentities` is added to struct `ManagedServiceIdentity`
- Field `Identity` is added to struct `DatabaseAccountUpdateParameters`
- Field `Cors` is added to struct `DatabaseAccountUpdateProperties`
- Field `Cors` is added to struct `DatabaseAccountCreateUpdateProperties`
- Field `Cors` is added to struct `DefaultRequestDatabaseAccountCreateUpdateProperties`

