# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. PrivateEndpointConnectionsClient.Delete
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string
1. PrivateEndpointConnectionsClient.DeletePreparer
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string
1. PrivateEndpointConnectionsClient.Get
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string
1. PrivateEndpointConnectionsClient.GetPreparer
	- Params
		- From: context.Context, string, string, string, string
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string
1. PrivateEndpointConnectionsClient.ListByResource
	- Params
		- From: context.Context, string, string, string, string, *int32
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string, *int32
1. PrivateEndpointConnectionsClient.ListByResourceComplete
	- Params
		- From: context.Context, string, string, string, string, *int32
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string, *int32
1. PrivateEndpointConnectionsClient.ListByResourcePreparer
	- Params
		- From: context.Context, string, string, string, string, *int32
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string, *int32
1. PrivateEndpointConnectionsClient.Update
	- Params
		- From: context.Context, string, string, string, string, PrivateEndpointConnection
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string, PrivateEndpointConnection
1. PrivateEndpointConnectionsClient.UpdatePreparer
	- Params
		- From: context.Context, string, string, string, string, PrivateEndpointConnection
		- To: context.Context, string, PrivateEndpointConnectionsParentType, string, string, PrivateEndpointConnection

## Additive Changes

### New Constants

1. PrivateEndpointConnectionsParentType.PrivateEndpointConnectionsParentTypeDomains
1. PrivateEndpointConnectionsParentType.PrivateEndpointConnectionsParentTypeTopics

### New Funcs

1. PossiblePrivateEndpointConnectionsParentTypeValues() []PrivateEndpointConnectionsParentType

### Struct Changes

#### New Struct Fields

1. Operation.IsDataAction
