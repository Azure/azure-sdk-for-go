# Unreleased

## Additive Changes

### New Constants

1. SchemaType.SchemaTypeJSON
1. SchemaType.SchemaTypeXML

### New Funcs

1. *GlobalSchemaCollectionIterator.Next() error
1. *GlobalSchemaCollectionIterator.NextWithContext(context.Context) error
1. *GlobalSchemaCollectionPage.Next() error
1. *GlobalSchemaCollectionPage.NextWithContext(context.Context) error
1. *GlobalSchemaContract.UnmarshalJSON([]byte) error
1. *GlobalSchemaCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. GlobalSchemaClient.CreateOrUpdate(context.Context, string, string, string, GlobalSchemaContract, string) (GlobalSchemaCreateOrUpdateFuture, error)
1. GlobalSchemaClient.CreateOrUpdatePreparer(context.Context, string, string, string, GlobalSchemaContract, string) (*http.Request, error)
1. GlobalSchemaClient.CreateOrUpdateResponder(*http.Response) (GlobalSchemaContract, error)
1. GlobalSchemaClient.CreateOrUpdateSender(*http.Request) (GlobalSchemaCreateOrUpdateFuture, error)
1. GlobalSchemaClient.Delete(context.Context, string, string, string, string) (autorest.Response, error)
1. GlobalSchemaClient.DeletePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. GlobalSchemaClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. GlobalSchemaClient.DeleteSender(*http.Request) (*http.Response, error)
1. GlobalSchemaClient.Get(context.Context, string, string, string) (GlobalSchemaContract, error)
1. GlobalSchemaClient.GetEntityTag(context.Context, string, string, string) (autorest.Response, error)
1. GlobalSchemaClient.GetEntityTagPreparer(context.Context, string, string, string) (*http.Request, error)
1. GlobalSchemaClient.GetEntityTagResponder(*http.Response) (autorest.Response, error)
1. GlobalSchemaClient.GetEntityTagSender(*http.Request) (*http.Response, error)
1. GlobalSchemaClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. GlobalSchemaClient.GetResponder(*http.Response) (GlobalSchemaContract, error)
1. GlobalSchemaClient.GetSender(*http.Request) (*http.Response, error)
1. GlobalSchemaClient.ListByService(context.Context, string, string, string, *int32, *int32) (GlobalSchemaCollectionPage, error)
1. GlobalSchemaClient.ListByServiceComplete(context.Context, string, string, string, *int32, *int32) (GlobalSchemaCollectionIterator, error)
1. GlobalSchemaClient.ListByServicePreparer(context.Context, string, string, string, *int32, *int32) (*http.Request, error)
1. GlobalSchemaClient.ListByServiceResponder(*http.Response) (GlobalSchemaCollection, error)
1. GlobalSchemaClient.ListByServiceSender(*http.Request) (*http.Response, error)
1. GlobalSchemaCollection.IsEmpty() bool
1. GlobalSchemaCollection.MarshalJSON() ([]byte, error)
1. GlobalSchemaCollectionIterator.NotDone() bool
1. GlobalSchemaCollectionIterator.Response() GlobalSchemaCollection
1. GlobalSchemaCollectionIterator.Value() GlobalSchemaContract
1. GlobalSchemaCollectionPage.NotDone() bool
1. GlobalSchemaCollectionPage.Response() GlobalSchemaCollection
1. GlobalSchemaCollectionPage.Values() []GlobalSchemaContract
1. GlobalSchemaContract.MarshalJSON() ([]byte, error)
1. NewGlobalSchemaClient(string) GlobalSchemaClient
1. NewGlobalSchemaClientWithBaseURI(string, string) GlobalSchemaClient
1. NewGlobalSchemaCollectionIterator(GlobalSchemaCollectionPage) GlobalSchemaCollectionIterator
1. NewGlobalSchemaCollectionPage(GlobalSchemaCollection, func(context.Context, GlobalSchemaCollection) (GlobalSchemaCollection, error)) GlobalSchemaCollectionPage
1. PossibleSchemaTypeValues() []SchemaType

### Struct Changes

#### New Structs

1. GlobalSchemaClient
1. GlobalSchemaCollection
1. GlobalSchemaCollectionIterator
1. GlobalSchemaCollectionPage
1. GlobalSchemaContract
1. GlobalSchemaContractProperties
1. GlobalSchemaCreateOrUpdateFuture
