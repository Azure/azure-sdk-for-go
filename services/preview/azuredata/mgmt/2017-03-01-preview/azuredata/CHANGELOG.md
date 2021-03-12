Generated from https://github.com/Azure/azure-rest-api-specs/tree/f0093fcd77ee0c8c82ae5c071a76ccb055aa9927/specification/azuredata/resource-manager/readme.md tag: `package-2017-03-01-preview`

Code generator 


## Breaking Changes

### Removed Funcs

1. *DataControllerResource.UnmarshalJSON([]byte) error
1. *PageOfDataControllerResourceIterator.Next() error
1. *PageOfDataControllerResourceIterator.NextWithContext(context.Context) error
1. *PageOfDataControllerResourcePage.Next() error
1. *PageOfDataControllerResourcePage.NextWithContext(context.Context) error
1. *PostgresInstance.UnmarshalJSON([]byte) error
1. *PostgresInstanceListResultIterator.Next() error
1. *PostgresInstanceListResultIterator.NextWithContext(context.Context) error
1. *PostgresInstanceListResultPage.Next() error
1. *PostgresInstanceListResultPage.NextWithContext(context.Context) error
1. *SQLManagedInstance.UnmarshalJSON([]byte) error
1. *SQLManagedInstanceListResultIterator.Next() error
1. *SQLManagedInstanceListResultIterator.NextWithContext(context.Context) error
1. *SQLManagedInstanceListResultPage.Next() error
1. *SQLManagedInstanceListResultPage.NextWithContext(context.Context) error
1. *SQLServerInstance.UnmarshalJSON([]byte) error
1. *SQLServerInstanceListResultIterator.Next() error
1. *SQLServerInstanceListResultIterator.NextWithContext(context.Context) error
1. *SQLServerInstanceListResultPage.Next() error
1. *SQLServerInstanceListResultPage.NextWithContext(context.Context) error
1. DataControllerResource.MarshalJSON() ([]byte, error)
1. DataControllerUpdate.MarshalJSON() ([]byte, error)
1. DataControllersClient.DeleteDataController(context.Context, string, string) (autorest.Response, error)
1. DataControllersClient.DeleteDataControllerPreparer(context.Context, string, string) (*http.Request, error)
1. DataControllersClient.DeleteDataControllerResponder(*http.Response) (autorest.Response, error)
1. DataControllersClient.DeleteDataControllerSender(*http.Request) (*http.Response, error)
1. DataControllersClient.GetDataController(context.Context, string, string) (DataControllerResource, error)
1. DataControllersClient.GetDataControllerPreparer(context.Context, string, string) (*http.Request, error)
1. DataControllersClient.GetDataControllerResponder(*http.Response) (DataControllerResource, error)
1. DataControllersClient.GetDataControllerSender(*http.Request) (*http.Response, error)
1. DataControllersClient.ListInGroup(context.Context, string) (PageOfDataControllerResourcePage, error)
1. DataControllersClient.ListInGroupComplete(context.Context, string) (PageOfDataControllerResourceIterator, error)
1. DataControllersClient.ListInGroupPreparer(context.Context, string) (*http.Request, error)
1. DataControllersClient.ListInGroupResponder(*http.Response) (PageOfDataControllerResource, error)
1. DataControllersClient.ListInGroupSender(*http.Request) (*http.Response, error)
1. DataControllersClient.ListInSubscription(context.Context) (PageOfDataControllerResourcePage, error)
1. DataControllersClient.ListInSubscriptionComplete(context.Context) (PageOfDataControllerResourceIterator, error)
1. DataControllersClient.ListInSubscriptionPreparer(context.Context) (*http.Request, error)
1. DataControllersClient.ListInSubscriptionResponder(*http.Response) (PageOfDataControllerResource, error)
1. DataControllersClient.ListInSubscriptionSender(*http.Request) (*http.Response, error)
1. DataControllersClient.PatchDataController(context.Context, string, string, DataControllerUpdate) (DataControllerResource, error)
1. DataControllersClient.PatchDataControllerPreparer(context.Context, string, string, DataControllerUpdate) (*http.Request, error)
1. DataControllersClient.PatchDataControllerResponder(*http.Response) (DataControllerResource, error)
1. DataControllersClient.PatchDataControllerSender(*http.Request) (*http.Response, error)
1. DataControllersClient.PutDataController(context.Context, string, DataControllerResource, string) (DataControllerResource, error)
1. DataControllersClient.PutDataControllerPreparer(context.Context, string, DataControllerResource, string) (*http.Request, error)
1. DataControllersClient.PutDataControllerResponder(*http.Response) (DataControllerResource, error)
1. DataControllersClient.PutDataControllerSender(*http.Request) (*http.Response, error)
1. NewDataControllersClient(string, string) DataControllersClient
1. NewDataControllersClientWithBaseURI(string, string, string) DataControllersClient
1. NewPageOfDataControllerResourceIterator(PageOfDataControllerResourcePage) PageOfDataControllerResourceIterator
1. NewPageOfDataControllerResourcePage(PageOfDataControllerResource, func(context.Context, PageOfDataControllerResource) (PageOfDataControllerResource, error)) PageOfDataControllerResourcePage
1. NewPostgresInstanceListResultIterator(PostgresInstanceListResultPage) PostgresInstanceListResultIterator
1. NewPostgresInstanceListResultPage(PostgresInstanceListResult, func(context.Context, PostgresInstanceListResult) (PostgresInstanceListResult, error)) PostgresInstanceListResultPage
1. NewPostgresInstancesClient(string, string) PostgresInstancesClient
1. NewPostgresInstancesClientWithBaseURI(string, string, string) PostgresInstancesClient
1. NewSQLManagedInstanceListResultIterator(SQLManagedInstanceListResultPage) SQLManagedInstanceListResultIterator
1. NewSQLManagedInstanceListResultPage(SQLManagedInstanceListResult, func(context.Context, SQLManagedInstanceListResult) (SQLManagedInstanceListResult, error)) SQLManagedInstanceListResultPage
1. NewSQLManagedInstancesClient(string, string) SQLManagedInstancesClient
1. NewSQLManagedInstancesClientWithBaseURI(string, string, string) SQLManagedInstancesClient
1. NewSQLServerInstanceListResultIterator(SQLServerInstanceListResultPage) SQLServerInstanceListResultIterator
1. NewSQLServerInstanceListResultPage(SQLServerInstanceListResult, func(context.Context, SQLServerInstanceListResult) (SQLServerInstanceListResult, error)) SQLServerInstanceListResultPage
1. NewSQLServerInstancesClient(string, string) SQLServerInstancesClient
1. NewSQLServerInstancesClientWithBaseURI(string, string, string) SQLServerInstancesClient
1. PageOfDataControllerResource.IsEmpty() bool
1. PageOfDataControllerResourceIterator.NotDone() bool
1. PageOfDataControllerResourceIterator.Response() PageOfDataControllerResource
1. PageOfDataControllerResourceIterator.Value() DataControllerResource
1. PageOfDataControllerResourcePage.NotDone() bool
1. PageOfDataControllerResourcePage.Response() PageOfDataControllerResource
1. PageOfDataControllerResourcePage.Values() []DataControllerResource
1. PostgresInstance.MarshalJSON() ([]byte, error)
1. PostgresInstanceListResult.IsEmpty() bool
1. PostgresInstanceListResultIterator.NotDone() bool
1. PostgresInstanceListResultIterator.Response() PostgresInstanceListResult
1. PostgresInstanceListResultIterator.Value() PostgresInstance
1. PostgresInstanceListResultPage.NotDone() bool
1. PostgresInstanceListResultPage.Response() PostgresInstanceListResult
1. PostgresInstanceListResultPage.Values() []PostgresInstance
1. PostgresInstanceUpdate.MarshalJSON() ([]byte, error)
1. PostgresInstancesClient.Create(context.Context, string, string) (PostgresInstance, error)
1. PostgresInstancesClient.CreatePreparer(context.Context, string, string) (*http.Request, error)
1. PostgresInstancesClient.CreateResponder(*http.Response) (PostgresInstance, error)
1. PostgresInstancesClient.CreateSender(*http.Request) (*http.Response, error)
1. PostgresInstancesClient.Delete(context.Context, string, string) (autorest.Response, error)
1. PostgresInstancesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. PostgresInstancesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. PostgresInstancesClient.DeleteSender(*http.Request) (*http.Response, error)
1. PostgresInstancesClient.Get(context.Context, string, string) (PostgresInstance, error)
1. PostgresInstancesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. PostgresInstancesClient.GetResponder(*http.Response) (PostgresInstance, error)
1. PostgresInstancesClient.GetSender(*http.Request) (*http.Response, error)
1. PostgresInstancesClient.List(context.Context) (PostgresInstanceListResultPage, error)
1. PostgresInstancesClient.ListByResourceGroup(context.Context, string) (PostgresInstanceListResultPage, error)
1. PostgresInstancesClient.ListByResourceGroupComplete(context.Context, string) (PostgresInstanceListResultIterator, error)
1. PostgresInstancesClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. PostgresInstancesClient.ListByResourceGroupResponder(*http.Response) (PostgresInstanceListResult, error)
1. PostgresInstancesClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. PostgresInstancesClient.ListComplete(context.Context) (PostgresInstanceListResultIterator, error)
1. PostgresInstancesClient.ListPreparer(context.Context) (*http.Request, error)
1. PostgresInstancesClient.ListResponder(*http.Response) (PostgresInstanceListResult, error)
1. PostgresInstancesClient.ListSender(*http.Request) (*http.Response, error)
1. PostgresInstancesClient.Update(context.Context, string, string, PostgresInstanceUpdate) (PostgresInstance, error)
1. PostgresInstancesClient.UpdatePreparer(context.Context, string, string, PostgresInstanceUpdate) (*http.Request, error)
1. PostgresInstancesClient.UpdateResponder(*http.Response) (PostgresInstance, error)
1. PostgresInstancesClient.UpdateSender(*http.Request) (*http.Response, error)
1. SQLManagedInstance.MarshalJSON() ([]byte, error)
1. SQLManagedInstanceListResult.IsEmpty() bool
1. SQLManagedInstanceListResultIterator.NotDone() bool
1. SQLManagedInstanceListResultIterator.Response() SQLManagedInstanceListResult
1. SQLManagedInstanceListResultIterator.Value() SQLManagedInstance
1. SQLManagedInstanceListResultPage.NotDone() bool
1. SQLManagedInstanceListResultPage.Response() SQLManagedInstanceListResult
1. SQLManagedInstanceListResultPage.Values() []SQLManagedInstance
1. SQLManagedInstanceUpdate.MarshalJSON() ([]byte, error)
1. SQLManagedInstancesClient.Create(context.Context, string, string, SQLManagedInstance) (SQLManagedInstance, error)
1. SQLManagedInstancesClient.CreatePreparer(context.Context, string, string, SQLManagedInstance) (*http.Request, error)
1. SQLManagedInstancesClient.CreateResponder(*http.Response) (SQLManagedInstance, error)
1. SQLManagedInstancesClient.CreateSender(*http.Request) (*http.Response, error)
1. SQLManagedInstancesClient.Delete(context.Context, string, string) (autorest.Response, error)
1. SQLManagedInstancesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. SQLManagedInstancesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. SQLManagedInstancesClient.DeleteSender(*http.Request) (*http.Response, error)
1. SQLManagedInstancesClient.Get(context.Context, string, string) (SQLManagedInstance, error)
1. SQLManagedInstancesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. SQLManagedInstancesClient.GetResponder(*http.Response) (SQLManagedInstance, error)
1. SQLManagedInstancesClient.GetSender(*http.Request) (*http.Response, error)
1. SQLManagedInstancesClient.List(context.Context) (SQLManagedInstanceListResultPage, error)
1. SQLManagedInstancesClient.ListByResourceGroup(context.Context, string) (SQLManagedInstanceListResultPage, error)
1. SQLManagedInstancesClient.ListByResourceGroupComplete(context.Context, string) (SQLManagedInstanceListResultIterator, error)
1. SQLManagedInstancesClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. SQLManagedInstancesClient.ListByResourceGroupResponder(*http.Response) (SQLManagedInstanceListResult, error)
1. SQLManagedInstancesClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. SQLManagedInstancesClient.ListComplete(context.Context) (SQLManagedInstanceListResultIterator, error)
1. SQLManagedInstancesClient.ListPreparer(context.Context) (*http.Request, error)
1. SQLManagedInstancesClient.ListResponder(*http.Response) (SQLManagedInstanceListResult, error)
1. SQLManagedInstancesClient.ListSender(*http.Request) (*http.Response, error)
1. SQLManagedInstancesClient.Update(context.Context, string, string, SQLManagedInstanceUpdate) (SQLManagedInstance, error)
1. SQLManagedInstancesClient.UpdatePreparer(context.Context, string, string, SQLManagedInstanceUpdate) (*http.Request, error)
1. SQLManagedInstancesClient.UpdateResponder(*http.Response) (SQLManagedInstance, error)
1. SQLManagedInstancesClient.UpdateSender(*http.Request) (*http.Response, error)
1. SQLServerInstance.MarshalJSON() ([]byte, error)
1. SQLServerInstanceListResult.IsEmpty() bool
1. SQLServerInstanceListResultIterator.NotDone() bool
1. SQLServerInstanceListResultIterator.Response() SQLServerInstanceListResult
1. SQLServerInstanceListResultIterator.Value() SQLServerInstance
1. SQLServerInstanceListResultPage.NotDone() bool
1. SQLServerInstanceListResultPage.Response() SQLServerInstanceListResult
1. SQLServerInstanceListResultPage.Values() []SQLServerInstance
1. SQLServerInstanceProperties.MarshalJSON() ([]byte, error)
1. SQLServerInstanceUpdate.MarshalJSON() ([]byte, error)
1. SQLServerInstancesClient.Create(context.Context, string, string, SQLServerInstance) (SQLServerInstance, error)
1. SQLServerInstancesClient.CreatePreparer(context.Context, string, string, SQLServerInstance) (*http.Request, error)
1. SQLServerInstancesClient.CreateResponder(*http.Response) (SQLServerInstance, error)
1. SQLServerInstancesClient.CreateSender(*http.Request) (*http.Response, error)
1. SQLServerInstancesClient.Delete(context.Context, string, string) (autorest.Response, error)
1. SQLServerInstancesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. SQLServerInstancesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. SQLServerInstancesClient.DeleteSender(*http.Request) (*http.Response, error)
1. SQLServerInstancesClient.Get(context.Context, string, string) (SQLServerInstance, error)
1. SQLServerInstancesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. SQLServerInstancesClient.GetResponder(*http.Response) (SQLServerInstance, error)
1. SQLServerInstancesClient.GetSender(*http.Request) (*http.Response, error)
1. SQLServerInstancesClient.List(context.Context) (SQLServerInstanceListResultPage, error)
1. SQLServerInstancesClient.ListByResourceGroup(context.Context, string) (SQLServerInstanceListResultPage, error)
1. SQLServerInstancesClient.ListByResourceGroupComplete(context.Context, string) (SQLServerInstanceListResultIterator, error)
1. SQLServerInstancesClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. SQLServerInstancesClient.ListByResourceGroupResponder(*http.Response) (SQLServerInstanceListResult, error)
1. SQLServerInstancesClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. SQLServerInstancesClient.ListComplete(context.Context) (SQLServerInstanceListResultIterator, error)
1. SQLServerInstancesClient.ListPreparer(context.Context) (*http.Request, error)
1. SQLServerInstancesClient.ListResponder(*http.Response) (SQLServerInstanceListResult, error)
1. SQLServerInstancesClient.ListSender(*http.Request) (*http.Response, error)
1. SQLServerInstancesClient.Update(context.Context, string, string, SQLServerInstanceUpdate) (SQLServerInstance, error)
1. SQLServerInstancesClient.UpdatePreparer(context.Context, string, string, SQLServerInstanceUpdate) (*http.Request, error)
1. SQLServerInstancesClient.UpdateResponder(*http.Response) (SQLServerInstance, error)
1. SQLServerInstancesClient.UpdateSender(*http.Request) (*http.Response, error)

## Struct Changes

### Removed Structs

1. DataControllerProperties
1. DataControllerResource
1. DataControllerUpdate
1. DataControllersClient
1. OnPremiseProperty
1. PageOfDataControllerResource
1. PageOfDataControllerResourceIterator
1. PageOfDataControllerResourcePage
1. PostgresInstance
1. PostgresInstanceListResult
1. PostgresInstanceListResultIterator
1. PostgresInstanceListResultPage
1. PostgresInstanceProperties
1. PostgresInstanceUpdate
1. PostgresInstancesClient
1. SQLManagedInstance
1. SQLManagedInstanceListResult
1. SQLManagedInstanceListResultIterator
1. SQLManagedInstanceListResultPage
1. SQLManagedInstanceProperties
1. SQLManagedInstanceUpdate
1. SQLManagedInstancesClient
1. SQLServerInstance
1. SQLServerInstanceListResult
1. SQLServerInstanceListResultIterator
1. SQLServerInstanceListResultPage
1. SQLServerInstanceProperties
1. SQLServerInstanceUpdate
1. SQLServerInstancesClient
