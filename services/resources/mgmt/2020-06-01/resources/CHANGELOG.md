
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewProviderListResultPage` signature has been changed from `(func(context.Context, ProviderListResult) (ProviderListResult, error))` to `(ProviderListResult,func(context.Context, ProviderListResult) (ProviderListResult, error))`
- Function `NewGroupListResultPage` signature has been changed from `(func(context.Context, GroupListResult) (GroupListResult, error))` to `(GroupListResult,func(context.Context, GroupListResult) (GroupListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewDeploymentListResultPage` signature has been changed from `(func(context.Context, DeploymentListResult) (DeploymentListResult, error))` to `(DeploymentListResult,func(context.Context, DeploymentListResult) (DeploymentListResult, error))`
- Function `NewListResultPage` signature has been changed from `(func(context.Context, ListResult) (ListResult, error))` to `(ListResult,func(context.Context, ListResult) (ListResult, error))`
- Function `NewTagsListResultPage` signature has been changed from `(func(context.Context, TagsListResult) (TagsListResult, error))` to `(TagsListResult,func(context.Context, TagsListResult) (TagsListResult, error))`
- Function `NewDeploymentOperationsListResultPage` signature has been changed from `(func(context.Context, DeploymentOperationsListResult) (DeploymentOperationsListResult, error))` to `(DeploymentOperationsListResult,func(context.Context, DeploymentOperationsListResult) (DeploymentOperationsListResult, error))`

## New Content

- Const `ExpressionEvaluationOptionsScopeTypeNotSpecified` is added
- Const `ExpressionEvaluationOptionsScopeTypeOuter` is added
- Const `ExpressionEvaluationOptionsScopeTypeInner` is added
- Function `ProvidersClient.RegisterAtManagementGroupScopePreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `ProvidersClient.RegisterAtManagementGroupScopeResponder(*http.Response) (autorest.Response,error)` is added
- Function `ProvidersClient.RegisterAtManagementGroupScope(context.Context,string,string) (autorest.Response,error)` is added
- Function `ProvidersClient.RegisterAtManagementGroupScopeSender(*http.Request) (*http.Response,error)` is added
- Function `PossibleExpressionEvaluationOptionsScopeTypeValues() []ExpressionEvaluationOptionsScopeType` is added
- Struct `ExpressionEvaluationOptions` is added
- Field `ExpressionEvaluationOptions` is added to struct `DeploymentWhatIfProperties`
- Field `ExpressionEvaluationOptions` is added to struct `DeploymentProperties`

