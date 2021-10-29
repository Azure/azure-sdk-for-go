# Unreleased

## Additive Changes

### New Funcs

1. Client.ListSubscriptionOperations(context.Context) (OperationsPage, error)
1. Client.ListSubscriptionOperationsComplete(context.Context) (OperationsIterator, error)
1. Client.ListSubscriptionOperationsPreparer(context.Context) (*http.Request, error)
1. Client.ListSubscriptionOperationsResponder(*http.Response) (Operations, error)
1. Client.ListSubscriptionOperationsSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Struct Fields

1. Dimension.InternalMetricName
1. Dimension.InternalName
1. Dimension.SourceMdmNamespace
1. Dimension.ToBeExportedToShoebox
