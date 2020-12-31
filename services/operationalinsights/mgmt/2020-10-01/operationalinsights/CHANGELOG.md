Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Funcs

1. *ClustersCreateOrUpdateFuture.Result(ClustersClient) (Cluster, error)
1. *ClustersDeleteFuture.Result(ClustersClient) (autorest.Response, error)
1. *WorkspacesCreateOrUpdateFuture.Result(WorkspacesClient) (Workspace, error)
1. *WorkspacesDeleteFuture.Result(WorkspacesClient) (autorest.Response, error)

## Struct Changes

### Removed Struct Fields

1. ClustersCreateOrUpdateFuture.azure.Future
1. ClustersDeleteFuture.azure.Future
1. WorkspacesCreateOrUpdateFuture.azure.Future
1. WorkspacesDeleteFuture.azure.Future

### New Funcs

1. *Table.UnmarshalJSON([]byte) error
1. NewTablesClient(string) TablesClient
1. NewTablesClientWithBaseURI(string, string) TablesClient
1. Table.MarshalJSON() ([]byte, error)
1. TableProperties.MarshalJSON() ([]byte, error)
1. TablesClient.Get(context.Context, string, string, string) (Table, error)
1. TablesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. TablesClient.GetResponder(*http.Response) (Table, error)
1. TablesClient.GetSender(*http.Request) (*http.Response, error)
1. TablesClient.ListByWorkspace(context.Context, string, string) (TablesListResult, error)
1. TablesClient.ListByWorkspacePreparer(context.Context, string, string) (*http.Request, error)
1. TablesClient.ListByWorkspaceResponder(*http.Response) (TablesListResult, error)
1. TablesClient.ListByWorkspaceSender(*http.Request) (*http.Response, error)
1. TablesClient.Update(context.Context, string, string, string, Table) (Table, error)
1. TablesClient.UpdatePreparer(context.Context, string, string, string, Table) (*http.Request, error)
1. TablesClient.UpdateResponder(*http.Response) (Table, error)
1. TablesClient.UpdateSender(*http.Request) (*http.Response, error)

## Struct Changes

### New Structs

1. Table
1. TableProperties
1. TablesClient
1. TablesListResult

### New Struct Fields

1. ClustersCreateOrUpdateFuture.Result
1. ClustersCreateOrUpdateFuture.azure.FutureAPI
1. ClustersDeleteFuture.Result
1. ClustersDeleteFuture.azure.FutureAPI
1. WorkspacesCreateOrUpdateFuture.Result
1. WorkspacesCreateOrUpdateFuture.azure.FutureAPI
1. WorkspacesDeleteFuture.Result
1. WorkspacesDeleteFuture.azure.FutureAPI
