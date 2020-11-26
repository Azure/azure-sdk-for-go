
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewResourceProviderOperationListPage` signature has been changed from `(func(context.Context, ResourceProviderOperationList) (ResourceProviderOperationList, error))` to `(ResourceProviderOperationList,func(context.Context, ResourceProviderOperationList) (ResourceProviderOperationList, error))`
- Function `NewDashboardListResultPage` signature has been changed from `(func(context.Context, DashboardListResult) (DashboardListResult, error))` to `(DashboardListResult,func(context.Context, DashboardListResult) (DashboardListResult, error))`
- Type of `ErrorDefinition.Code` has been changed from `*string` to `*int32`

## New Content

- Function `TenantConfigurationsClient.CreateSender(*http.Request) (*http.Response,error)` is added
- Function `TenantConfigurationsClient.List(context.Context) (ConfigurationList,error)` is added
- Function `TenantConfigurationsClient.CreateResponder(*http.Response) (Configuration,error)` is added
- Function `TenantConfigurationsClient.ListResponder(*http.Response) (ConfigurationList,error)` is added
- Function `TenantConfigurationsClient.Delete(context.Context) (autorest.Response,error)` is added
- Function `TenantConfigurationsClient.ListPreparer(context.Context) (*http.Request,error)` is added
- Function `TenantConfigurationsClient.Create(context.Context,Configuration) (Configuration,error)` is added
- Function `NewTenantConfigurationsClient(string) TenantConfigurationsClient` is added
- Function `TenantConfigurationsClient.DeleteSender(*http.Request) (*http.Response,error)` is added
- Function `TenantConfigurationsClient.DeleteResponder(*http.Response) (autorest.Response,error)` is added
- Function `TenantConfigurationsClient.GetResponder(*http.Response) (Configuration,error)` is added
- Function `TenantConfigurationsClient.CreatePreparer(context.Context,Configuration) (*http.Request,error)` is added
- Function `TenantConfigurationsClient.GetPreparer(context.Context) (*http.Request,error)` is added
- Function `TrackedResource.MarshalJSON() ([]byte,error)` is added
- Function `TenantConfigurationsClient.ListSender(*http.Request) (*http.Response,error)` is added
- Function `NewTenantConfigurationsClientWithBaseURI(string,string) TenantConfigurationsClient` is added
- Function `*Configuration.UnmarshalJSON([]byte) error` is added
- Function `TenantConfigurationsClient.Get(context.Context) (Configuration,error)` is added
- Function `TenantConfigurationsClient.GetSender(*http.Request) (*http.Response,error)` is added
- Function `Configuration.MarshalJSON() ([]byte,error)` is added
- Function `TenantConfigurationsClient.DeletePreparer(context.Context) (*http.Request,error)` is added
- Struct `AzureEntityResource` is added
- Struct `Configuration` is added
- Struct `ConfigurationList` is added
- Struct `ConfigurationProperties` is added
- Struct `ProxyResource` is added
- Struct `Resource` is added
- Struct `TenantConfigurationsClient` is added
- Struct `TrackedResource` is added

