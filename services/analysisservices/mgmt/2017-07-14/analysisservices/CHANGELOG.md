Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Funcs

1. *ServersCreateFuture.Result(ServersClient) (Server, error)
1. *ServersDeleteFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersResumeFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersSuspendFuture.Result(ServersClient) (autorest.Response, error)
1. *ServersUpdateFuture.Result(ServersClient) (Server, error)

## Struct Changes

### Removed Struct Fields

1. ServersCreateFuture.azure.Future
1. ServersDeleteFuture.azure.Future
1. ServersResumeFuture.azure.Future
1. ServersSuspendFuture.azure.Future
1. ServersUpdateFuture.azure.Future

### New Funcs

1. *OperationListResultIterator.Next() error
1. *OperationListResultIterator.NextWithContext(context.Context) error
1. *OperationListResultPage.Next() error
1. *OperationListResultPage.NextWithContext(context.Context) error
1. NewOperationListResultIterator(OperationListResultPage) OperationListResultIterator
1. NewOperationListResultPage(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage
1. NewOperationsClient(string) OperationsClient
1. NewOperationsClientWithBaseURI(string, string) OperationsClient
1. OperationDetail.MarshalJSON() ([]byte, error)
1. OperationDisplay.MarshalJSON() ([]byte, error)
1. OperationListResult.IsEmpty() bool
1. OperationListResultIterator.NotDone() bool
1. OperationListResultIterator.Response() OperationListResult
1. OperationListResultIterator.Value() OperationDetail
1. OperationListResultPage.NotDone() bool
1. OperationListResultPage.Response() OperationListResult
1. OperationListResultPage.Values() []OperationDetail
1. OperationsClient.List(context.Context) (OperationListResultPage, error)
1. OperationsClient.ListComplete(context.Context) (OperationListResultIterator, error)
1. OperationsClient.ListPreparer(context.Context) (*http.Request, error)
1. OperationsClient.ListResponder(*http.Response) (OperationListResult, error)
1. OperationsClient.ListSender(*http.Request) (*http.Response, error)

## Struct Changes

### New Structs

1. OperationDetail
1. OperationDisplay
1. OperationListResult
1. OperationListResultIterator
1. OperationListResultPage
1. OperationsClient
1. OperationsErrorResponse

### New Struct Fields

1. ServersCreateFuture.Result
1. ServersCreateFuture.azure.FutureAPI
1. ServersDeleteFuture.Result
1. ServersDeleteFuture.azure.FutureAPI
1. ServersResumeFuture.Result
1. ServersResumeFuture.azure.FutureAPI
1. ServersSuspendFuture.Result
1. ServersSuspendFuture.azure.FutureAPI
1. ServersUpdateFuture.Result
1. ServersUpdateFuture.azure.FutureAPI
