# Unreleased

## Additive Changes

### New Funcs

1. NewOperationClient() OperationClient
1. NewOperationClientWithBaseURI(string) OperationClient
1. OperationClient.Get(context.Context, string) (CreationResult, error)
1. OperationClient.GetPreparer(context.Context, string) (*http.Request, error)
1. OperationClient.GetResponder(*http.Response) (CreationResult, error)
1. OperationClient.GetSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. ErrorResponseBody
1. OperationClient
