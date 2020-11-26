
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewOperationListPage` signature has been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList,func(context.Context, OperationList) (OperationList, error))`
- Function `NewServiceResourceListPage` signature has been changed from `(func(context.Context, ServiceResourceList) (ServiceResourceList, error))` to `(ServiceResourceList,func(context.Context, ServiceResourceList) (ServiceResourceList, error))`

## New Content

- Function `ServiceClient.CheckNameAvailabilitySender(*http.Request) (*http.Response,error)` is added
- Function `ServiceClient.CheckNameAvailabilityResponder(*http.Response) (NameAvailability,error)` is added
- Function `ServiceClient.CheckNameAvailability(context.Context,*NameAvailabilityParameters) (NameAvailability,error)` is added
- Function `ServiceClient.CheckNameAvailabilityPreparer(context.Context,*NameAvailabilityParameters) (*http.Request,error)` is added
- Struct `NameAvailability` is added
- Struct `NameAvailabilityParameters` is added

