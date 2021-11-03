# Change History

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. BaseClient.ProviderNamespace

### Signature Changes

#### Funcs

1. BaseClient.ListOperations
	- Params
		- From: context.Context, string
		- To: context.Context
1. BaseClient.ListOperationsComplete
	- Params
		- From: context.Context, string
		- To: context.Context
1. BaseClient.ListOperationsPreparer
	- Params
		- From: context.Context, string
		- To: context.Context
1. Client.Get
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. Client.GetPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. Client.List
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. Client.ListAll
	- Params
		- From: context.Context, string
		- To: context.Context
1. Client.ListAllComplete
	- Params
		- From: context.Context, string
		- To: context.Context
1. Client.ListAllPreparer
	- Params
		- From: context.Context, string
		- To: context.Context
1. Client.ListComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. Client.ListPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string
1. Client.Register
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. Client.RegisterPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. Client.Unregister
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. Client.UnregisterPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string
1. New
	- Params
		- From: string, string
		- To: string
1. NewClient
	- Params
		- From: string, string
		- To: string
1. NewClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSubscriptionFeatureRegistrationsClient
	- Params
		- From: string, string
		- To: string
1. NewSubscriptionFeatureRegistrationsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. SubscriptionFeatureRegistrationsClient.ListAllBySubscription
	- Params
		- From: context.Context, string
		- To: context.Context
1. SubscriptionFeatureRegistrationsClient.ListAllBySubscriptionComplete
	- Params
		- From: context.Context, string
		- To: context.Context
1. SubscriptionFeatureRegistrationsClient.ListAllBySubscriptionPreparer
	- Params
		- From: context.Context, string
		- To: context.Context
