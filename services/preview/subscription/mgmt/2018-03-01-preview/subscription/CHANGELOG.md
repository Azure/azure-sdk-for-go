Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Funcs

1. *FactoryCreateSubscriptionInEnrollmentAccountFuture.Result(FactoryClient) (CreationResult, error)

## Struct Changes

### Removed Struct Fields

1. FactoryCreateSubscriptionInEnrollmentAccountFuture.azure.Future

### New Funcs

1. NewOperationClient() OperationClient
1. NewOperationClientWithBaseURI(string) OperationClient
1. OperationClient.Get(context.Context, string) (CreationResult, error)
1. OperationClient.GetPreparer(context.Context, string) (*http.Request, error)
1. OperationClient.GetResponder(*http.Response) (CreationResult, error)
1. OperationClient.GetSender(*http.Request) (*http.Response, error)

## Struct Changes

### New Structs

1. ErrorResponseBody
1. OperationClient

### New Struct Fields

1. FactoryCreateSubscriptionInEnrollmentAccountFuture.Result
1. FactoryCreateSubscriptionInEnrollmentAccountFuture.azure.FutureAPI
