# Unreleased

## Additive Changes

### New Constants

1. PrivateEndpointConnectionProvisioningState.Creating
1. PrivateEndpointConnectionProvisioningState.Deleting
1. PrivateEndpointConnectionProvisioningState.Failed
1. PrivateEndpointConnectionProvisioningState.Succeeded
1. PrivateEndpointServiceConnectionStatus.Approved
1. PrivateEndpointServiceConnectionStatus.Pending
1. PrivateEndpointServiceConnectionStatus.Rejected

### New Funcs

1. *PrivateEndpointConnection.UnmarshalJSON([]byte) error
1. NewPrivateEndpointConnectionsClient(string) PrivateEndpointConnectionsClient
1. NewPrivateEndpointConnectionsClientWithBaseURI(string, string) PrivateEndpointConnectionsClient
1. PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState
1. PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus
1. PrivateEndpoint.MarshalJSON() ([]byte, error)
1. PrivateEndpointConnection.MarshalJSON() ([]byte, error)
1. PrivateEndpointConnectionsClient.Create(context.Context, string, string, string, PrivateEndpointConnection) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.CreatePreparer(context.Context, string, string, string, PrivateEndpointConnection) (*http.Request, error)
1. PrivateEndpointConnectionsClient.CreateResponder(*http.Response) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.CreateSender(*http.Request) (*http.Response, error)
1. PrivateEndpointConnectionsClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. PrivateEndpointConnectionsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. PrivateEndpointConnectionsClient.DeleteSender(*http.Request) (*http.Response, error)
1. PrivateEndpointConnectionsClient.Get(context.Context, string, string, string) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.GetResponder(*http.Response) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateEndpointConnectionsClient.List(context.Context, string, string) (PrivateEndpointConnectionListResult, error)
1. PrivateEndpointConnectionsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.ListResponder(*http.Response) (PrivateEndpointConnectionListResult, error)
1. PrivateEndpointConnectionsClient.ListSender(*http.Request) (*http.Response, error)
1. StatusResult.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. PrivateEndpoint
1. PrivateEndpointConnection
1. PrivateEndpointConnectionListResult
1. PrivateEndpointConnectionProperties
1. PrivateEndpointConnectionsClient
1. PrivateLinkServiceConnectionState

#### New Struct Fields

1. StatusResult.PrivateEndpointConnections
