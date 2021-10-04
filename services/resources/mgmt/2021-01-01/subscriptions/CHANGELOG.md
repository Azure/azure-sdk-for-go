# Unreleased

## Breaking Changes

### Removed Funcs

1. *OperationListResultIterator.Next() error
1. *OperationListResultIterator.NextWithContext(context.Context) error
1. *OperationListResultPage.Next() error
1. *OperationListResultPage.NextWithContext(context.Context) error
1. NewOperationListResultIterator(OperationListResultPage) OperationListResultIterator
1. NewOperationListResultPage(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage
1. NewOperationsClient() OperationsClient
1. NewOperationsClientWithBaseURI(string) OperationsClient
1. OperationListResult.IsEmpty() bool
1. OperationListResultIterator.NotDone() bool
1. OperationListResultIterator.Response() OperationListResult
1. OperationListResultIterator.Value() Operation
1. OperationListResultPage.NotDone() bool
1. OperationListResultPage.Response() OperationListResult
1. OperationListResultPage.Values() []Operation
1. OperationsClient.List(context.Context) (OperationListResultPage, error)
1. OperationsClient.ListComplete(context.Context) (OperationListResultIterator, error)
1. OperationsClient.ListPreparer(context.Context) (*http.Request, error)
1. OperationsClient.ListResponder(*http.Response) (OperationListResult, error)
1. OperationsClient.ListSender(*http.Request) (*http.Response, error)

### Struct Changes

#### Removed Structs

1. OperationListResultIterator
1. OperationListResultPage
1. OperationsClient

#### Removed Struct Fields

1. OperationListResult.autorest.Response
