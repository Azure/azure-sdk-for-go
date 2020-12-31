Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Funcs

1. *MigrationConfigsCreateAndStartMigrationFuture.Result(MigrationConfigsClient) (MigrationConfigProperties, error)
1. *NamespacesCreateOrUpdateFuture.Result(NamespacesClient) (SBNamespace, error)
1. *NamespacesDeleteFuture.Result(NamespacesClient) (autorest.Response, error)
1. *PrivateEndpointConnectionsDeleteFuture.Result(PrivateEndpointConnectionsClient) (autorest.Response, error)

## Struct Changes

### Removed Struct Fields

1. MigrationConfigsCreateAndStartMigrationFuture.azure.Future
1. NamespacesCreateOrUpdateFuture.azure.Future
1. NamespacesDeleteFuture.azure.Future
1. PrivateEndpointConnectionsDeleteFuture.azure.Future

### New Funcs

1. *NetworkRuleSetListResultIterator.Next() error
1. *NetworkRuleSetListResultIterator.NextWithContext(context.Context) error
1. *NetworkRuleSetListResultPage.Next() error
1. *NetworkRuleSetListResultPage.NextWithContext(context.Context) error
1. NamespacesClient.ListNetworkRuleSets(context.Context, string, string) (NetworkRuleSetListResultPage, error)
1. NamespacesClient.ListNetworkRuleSetsComplete(context.Context, string, string) (NetworkRuleSetListResultIterator, error)
1. NamespacesClient.ListNetworkRuleSetsPreparer(context.Context, string, string) (*http.Request, error)
1. NamespacesClient.ListNetworkRuleSetsResponder(*http.Response) (NetworkRuleSetListResult, error)
1. NamespacesClient.ListNetworkRuleSetsSender(*http.Request) (*http.Response, error)
1. NetworkRuleSetListResult.IsEmpty() bool
1. NetworkRuleSetListResultIterator.NotDone() bool
1. NetworkRuleSetListResultIterator.Response() NetworkRuleSetListResult
1. NetworkRuleSetListResultIterator.Value() NetworkRuleSet
1. NetworkRuleSetListResultPage.NotDone() bool
1. NetworkRuleSetListResultPage.Response() NetworkRuleSetListResult
1. NetworkRuleSetListResultPage.Values() []NetworkRuleSet
1. NewNetworkRuleSetListResultIterator(NetworkRuleSetListResultPage) NetworkRuleSetListResultIterator
1. NewNetworkRuleSetListResultPage(NetworkRuleSetListResult, func(context.Context, NetworkRuleSetListResult) (NetworkRuleSetListResult, error)) NetworkRuleSetListResultPage

## Struct Changes

### New Structs

1. NetworkRuleSetListResult
1. NetworkRuleSetListResultIterator
1. NetworkRuleSetListResultPage

### New Struct Fields

1. MigrationConfigsCreateAndStartMigrationFuture.Result
1. MigrationConfigsCreateAndStartMigrationFuture.azure.FutureAPI
1. NamespacesCreateOrUpdateFuture.Result
1. NamespacesCreateOrUpdateFuture.azure.FutureAPI
1. NamespacesDeleteFuture.Result
1. NamespacesDeleteFuture.azure.FutureAPI
1. PrivateEndpointConnectionsDeleteFuture.Result
1. PrivateEndpointConnectionsDeleteFuture.azure.FutureAPI
