
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewTenantListResultPage` signature has been changed from `(func(context.Context, TenantListResult) (TenantListResult, error))` to `(TenantListResult,func(context.Context, TenantListResult) (TenantListResult, error))`
- Function `NewListResultPage` signature has been changed from `(func(context.Context, ListResult) (ListResult, error))` to `(ListResult,func(context.Context, ListResult) (ListResult, error))`

## New Content

- Const `Failed` is added
- Const `Accepted` is added
- Const `Succeeded` is added
- Const `DevTest` is added
- Const `Production` is added
- Function `Client.ListAlias(context.Context) (PutAliasListResult,error)` is added
- Function `Client.ListAliasResponder(*http.Response) (PutAliasListResult,error)` is added
- Function `Client.CreateAliasPreparer(context.Context,string,PutAliasRequest) (*http.Request,error)` is added
- Function `*CreateAliasFuture.Result(Client) (PutAliasResponse,error)` is added
- Function `Client.DeleteAlias(context.Context,string) (autorest.Response,error)` is added
- Function `Client.CreateAliasResponder(*http.Response) (PutAliasResponse,error)` is added
- Function `PossibleWorkloadValues() []Workload` is added
- Function `Client.CreateAlias(context.Context,string,PutAliasRequest) (CreateAliasFuture,error)` is added
- Function `Client.CreateAliasSender(*http.Request) (CreateAliasFuture,error)` is added
- Function `Client.DeleteAliasResponder(*http.Response) (autorest.Response,error)` is added
- Function `Client.GetAliasSender(*http.Request) (*http.Response,error)` is added
- Function `Client.ListAliasSender(*http.Request) (*http.Response,error)` is added
- Function `Client.GetAliasResponder(*http.Response) (PutAliasResponse,error)` is added
- Function `PutAliasResponse.MarshalJSON() ([]byte,error)` is added
- Function `Client.GetAlias(context.Context,string) (PutAliasResponse,error)` is added
- Function `Client.DeleteAliasPreparer(context.Context,string) (*http.Request,error)` is added
- Function `Client.DeleteAliasSender(*http.Request) (*http.Response,error)` is added
- Function `Client.GetAliasPreparer(context.Context,string) (*http.Request,error)` is added
- Function `PutAliasResponseProperties.MarshalJSON() ([]byte,error)` is added
- Function `PossibleProvisioningStateValues() []ProvisioningState` is added
- Function `Client.ListAliasPreparer(context.Context) (*http.Request,error)` is added
- Struct `CreateAliasFuture` is added
- Struct `ErrorResponseBody` is added
- Struct `PutAliasListResult` is added
- Struct `PutAliasRequest` is added
- Struct `PutAliasRequestProperties` is added
- Struct `PutAliasResponse` is added
- Struct `PutAliasResponseProperties` is added

