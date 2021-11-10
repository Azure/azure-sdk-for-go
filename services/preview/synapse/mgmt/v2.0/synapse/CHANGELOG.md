# Change History

## Breaking Changes

### Removed Funcs

1. KustoPoolClient.ListSkus(context.Context) (SkuDescriptionList, error)
1. KustoPoolClient.ListSkusPreparer(context.Context) (*http.Request, error)
1. KustoPoolClient.ListSkusResponder(*http.Response) (SkuDescriptionList, error)
1. KustoPoolClient.ListSkusSender(*http.Request) (*http.Response, error)
1. NewKustoPoolClient(string) KustoPoolClient
1. NewKustoPoolClientWithBaseURI(string, string) KustoPoolClient

### Struct Changes

#### Removed Structs

1. KustoPoolClient

## Additive Changes

### New Funcs

1. KustoPoolsClient.ListSkus(context.Context) (SkuDescriptionList, error)
1. KustoPoolsClient.ListSkusPreparer(context.Context) (*http.Request, error)
1. KustoPoolsClient.ListSkusResponder(*http.Response) (SkuDescriptionList, error)
1. KustoPoolsClient.ListSkusSender(*http.Request) (*http.Response, error)
