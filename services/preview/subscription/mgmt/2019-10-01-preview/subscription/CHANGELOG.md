Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewTenantListResultPage` parameter(s) have been changed from `(func(context.Context, TenantListResult) (TenantListResult, error))` to `(TenantListResult, func(context.Context, TenantListResult) (TenantListResult, error))`
- Function `NewListResultPage` parameter(s) have been changed from `(func(context.Context, ListResult) (ListResult, error))` to `(ListResult, func(context.Context, ListResult) (ListResult, error))`

## New Content

- New const `DevTest`
- New const `Failed`
- New const `Succeeded`
- New const `Production`
- New const `Accepted`
- New function `Client.GetAliasResponder(*http.Response) (PutAliasResponse, error)`
- New function `Client.CreateAlias(context.Context, string, PutAliasRequest) (CreateAliasFuture, error)`
- New function `PossibleWorkloadValues() []Workload`
- New function `Client.DeleteAlias(context.Context, string) (autorest.Response, error)`
- New function `Client.ListAliasSender(*http.Request) (*http.Response, error)`
- New function `PossibleProvisioningStateValues() []ProvisioningState`
- New function `Client.CreateAliasSender(*http.Request) (CreateAliasFuture, error)`
- New function `*CreateAliasFuture.Result(Client) (PutAliasResponse, error)`
- New function `Client.GetAliasPreparer(context.Context, string) (*http.Request, error)`
- New function `Client.ListAliasPreparer(context.Context) (*http.Request, error)`
- New function `PutAliasResponse.MarshalJSON() ([]byte, error)`
- New function `Client.GetAlias(context.Context, string) (PutAliasResponse, error)`
- New function `PutAliasResponseProperties.MarshalJSON() ([]byte, error)`
- New function `Client.DeleteAliasPreparer(context.Context, string) (*http.Request, error)`
- New function `Client.DeleteAliasResponder(*http.Response) (autorest.Response, error)`
- New function `Client.CreateAliasResponder(*http.Response) (PutAliasResponse, error)`
- New function `Client.ListAlias(context.Context) (PutAliasListResult, error)`
- New function `Client.GetAliasSender(*http.Request) (*http.Response, error)`
- New function `Client.DeleteAliasSender(*http.Request) (*http.Response, error)`
- New function `Client.ListAliasResponder(*http.Response) (PutAliasListResult, error)`
- New function `Client.CreateAliasPreparer(context.Context, string, PutAliasRequest) (*http.Request, error)`
- New struct `CreateAliasFuture`
- New struct `ErrorResponseBody`
- New struct `PutAliasListResult`
- New struct `PutAliasRequest`
- New struct `PutAliasRequestProperties`
- New struct `PutAliasResponse`
- New struct `PutAliasResponseProperties`
