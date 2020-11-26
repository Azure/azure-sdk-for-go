
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Const `Succeeded` type has been changed from `ProvisioningStateDR` to `EndPointProvisioningState`
- Const `Creating` type has been changed from `EntityStatus` to `EndPointProvisioningState`
- Const `Deleting` type has been changed from `EntityStatus` to `EndPointProvisioningState`
- Const `Failed` type has been changed from `ProvisioningStateDR` to `EndPointProvisioningState`
- Function `NewArmDisasterRecoveryListResultPage` signature has been changed from `(func(context.Context, ArmDisasterRecoveryListResult) (ArmDisasterRecoveryListResult, error))` to `(ArmDisasterRecoveryListResult,func(context.Context, ArmDisasterRecoveryListResult) (ArmDisasterRecoveryListResult, error))`
- Function `NewMessagingRegionsListResultPage` signature has been changed from `(func(context.Context, MessagingRegionsListResult) (MessagingRegionsListResult, error))` to `(MessagingRegionsListResult,func(context.Context, MessagingRegionsListResult) (MessagingRegionsListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewIPFilterRuleListResultPage` signature has been changed from `(func(context.Context, IPFilterRuleListResult) (IPFilterRuleListResult, error))` to `(IPFilterRuleListResult,func(context.Context, IPFilterRuleListResult) (IPFilterRuleListResult, error))`
- Function `NewAuthorizationRuleListResultPage` signature has been changed from `(func(context.Context, AuthorizationRuleListResult) (AuthorizationRuleListResult, error))` to `(AuthorizationRuleListResult,func(context.Context, AuthorizationRuleListResult) (AuthorizationRuleListResult, error))`
- Function `NewVirtualNetworkRuleListResultPage` signature has been changed from `(func(context.Context, VirtualNetworkRuleListResult) (VirtualNetworkRuleListResult, error))` to `(VirtualNetworkRuleListResult,func(context.Context, VirtualNetworkRuleListResult) (VirtualNetworkRuleListResult, error))`
- Function `NewClusterListResultPage` signature has been changed from `(func(context.Context, ClusterListResult) (ClusterListResult, error))` to `(ClusterListResult,func(context.Context, ClusterListResult) (ClusterListResult, error))`
- Function `NewListResultPage` signature has been changed from `(func(context.Context, ListResult) (ListResult, error))` to `(ListResult,func(context.Context, ListResult) (ListResult, error))`
- Function `NewEHNamespaceListResultPage` signature has been changed from `(func(context.Context, EHNamespaceListResult) (EHNamespaceListResult, error))` to `(EHNamespaceListResult,func(context.Context, EHNamespaceListResult) (EHNamespaceListResult, error))`
- Function `NewConsumerGroupListResultPage` signature has been changed from `(func(context.Context, ConsumerGroupListResult) (ConsumerGroupListResult, error))` to `(ConsumerGroupListResult,func(context.Context, ConsumerGroupListResult) (ConsumerGroupListResult, error))`
- Const `Renaming` has been removed
- Const `Active` has been removed
- Const `SendDisabled` has been removed
- Const `Restoring` has been removed
- Const `Disabled` has been removed
- Const `Unknown` has been removed
- Const `Accepted` has been removed
- Const `ReceiveDisabled` has been removed

## New Content

- Const `EntityStatusCreating` is added
- Const `Disconnected` is added
- Const `ProvisioningStateDRAccepted` is added
- Const `EntityStatusRestoring` is added
- Const `Rejected` is added
- Const `EntityStatusReceiveDisabled` is added
- Const `EntityStatusSendDisabled` is added
- Const `EntityStatusRenaming` is added
- Const `Pending` is added
- Const `Approved` is added
- Const `ProvisioningStateDRFailed` is added
- Const `EntityStatusActive` is added
- Const `EntityStatusUnknown` is added
- Const `EntityStatusDisabled` is added
- Const `ProvisioningStateDRSucceeded` is added
- Const `Canceled` is added
- Const `Updating` is added
- Const `EntityStatusDeleting` is added
- Function `PrivateLinkResource.MarshalJSON() ([]byte,error)` is added
- Function `PrivateEndpointConnectionsClient.ListResponder(*http.Response) (PrivateEndpointConnectionListResult,error)` is added
- Function `PossiblePrivateLinkConnectionStatusValues() []PrivateLinkConnectionStatus` is added
- Function `PrivateLinkResourcesClient.Get(context.Context,string,string) (PrivateLinkResourcesListResult,error)` is added
- Function `PrivateEndpointConnectionsClient.CreateOrUpdateResponder(*http.Response) (PrivateEndpointConnection,error)` is added
- Function `NewPrivateEndpointConnectionsClientWithBaseURI(string,string) PrivateEndpointConnectionsClient` is added
- Function `PrivateEndpointConnectionsClient.GetPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `NewPrivateLinkResourcesClient(string) PrivateLinkResourcesClient` is added
- Function `PrivateEndpointConnectionListResultPage.Values() []PrivateEndpointConnection` is added
- Function `PrivateLinkResourcesClient.GetSender(*http.Request) (*http.Response,error)` is added
- Function `PrivateEndpointConnectionsClient.ListComplete(context.Context,string,string) (PrivateEndpointConnectionListResultIterator,error)` is added
- Function `PrivateEndpointConnectionsClient.DeletePreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `PrivateLinkResourcesClient.GetResponder(*http.Response) (PrivateLinkResourcesListResult,error)` is added
- Function `PrivateEndpointConnectionsClient.Delete(context.Context,string,string,string) (PrivateEndpointConnectionsDeleteFuture,error)` is added
- Function `PossibleEndPointProvisioningStateValues() []EndPointProvisioningState` is added
- Function `*PrivateLinkResource.UnmarshalJSON([]byte) error` is added
- Function `PrivateEndpointConnectionListResult.IsEmpty() bool` is added
- Function `PrivateEndpointConnectionListResultIterator.NotDone() bool` is added
- Function `*PrivateEndpointConnectionListResultIterator.Next() error` is added
- Function `PrivateEndpointConnectionsClient.CreateOrUpdatePreparer(context.Context,string,string,string,PrivateEndpointConnection) (*http.Request,error)` is added
- Function `*PrivateEndpointConnectionListResultIterator.NextWithContext(context.Context) error` is added
- Function `PrivateEndpointConnectionsClient.GetResponder(*http.Response) (PrivateEndpointConnection,error)` is added
- Function `PrivateEndpointConnection.MarshalJSON() ([]byte,error)` is added
- Function `PrivateEndpointConnectionsClient.Get(context.Context,string,string,string) (PrivateEndpointConnection,error)` is added
- Function `PrivateLinkResourcesClient.GetPreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `PrivateEndpointConnectionsClient.DeleteResponder(*http.Response) (autorest.Response,error)` is added
- Function `NewPrivateEndpointConnectionListResultIterator(PrivateEndpointConnectionListResultPage) PrivateEndpointConnectionListResultIterator` is added
- Function `NewPrivateLinkResourcesClientWithBaseURI(string,string) PrivateLinkResourcesClient` is added
- Function `PrivateEndpointConnectionListResultPage.NotDone() bool` is added
- Function `PrivateEndpointConnectionsClient.CreateOrUpdate(context.Context,string,string,string,PrivateEndpointConnection) (PrivateEndpointConnection,error)` is added
- Function `*PrivateEndpointConnectionListResultPage.Next() error` is added
- Function `PrivateEndpointConnectionsClient.GetSender(*http.Request) (*http.Response,error)` is added
- Function `PrivateEndpointConnectionListResultIterator.Response() PrivateEndpointConnectionListResult` is added
- Function `PrivateEndpointConnectionsClient.ListPreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `PrivateEndpointConnectionsClient.DeleteSender(*http.Request) (PrivateEndpointConnectionsDeleteFuture,error)` is added
- Function `*PrivateEndpointConnectionsDeleteFuture.Result(PrivateEndpointConnectionsClient) (autorest.Response,error)` is added
- Function `*PrivateEndpointConnectionListResultPage.NextWithContext(context.Context) error` is added
- Function `PrivateEndpointConnectionListResultIterator.Value() PrivateEndpointConnection` is added
- Function `PrivateEndpointConnectionsClient.List(context.Context,string,string) (PrivateEndpointConnectionListResultPage,error)` is added
- Function `PrivateEndpointConnectionsClient.ListSender(*http.Request) (*http.Response,error)` is added
- Function `PrivateEndpointConnectionListResultPage.Response() PrivateEndpointConnectionListResult` is added
- Function `NewPrivateEndpointConnectionsClient(string) PrivateEndpointConnectionsClient` is added
- Function `PrivateEndpointConnectionsClient.CreateOrUpdateSender(*http.Request) (*http.Response,error)` is added
- Function `*PrivateEndpointConnection.UnmarshalJSON([]byte) error` is added
- Function `NewPrivateEndpointConnectionListResultPage(PrivateEndpointConnectionListResult,func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error)) PrivateEndpointConnectionListResultPage` is added
- Struct `ConnectionState` is added
- Struct `PrivateEndpoint` is added
- Struct `PrivateEndpointConnection` is added
- Struct `PrivateEndpointConnectionListResult` is added
- Struct `PrivateEndpointConnectionListResultIterator` is added
- Struct `PrivateEndpointConnectionListResultPage` is added
- Struct `PrivateEndpointConnectionProperties` is added
- Struct `PrivateEndpointConnectionsClient` is added
- Struct `PrivateEndpointConnectionsDeleteFuture` is added
- Struct `PrivateLinkResource` is added
- Struct `PrivateLinkResourceProperties` is added
- Struct `PrivateLinkResourcesClient` is added
- Struct `PrivateLinkResourcesListResult` is added
- Field `TrustedServiceAccessEnabled` is added to struct `NetworkRuleSetProperties`

