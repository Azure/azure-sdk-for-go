# Change History

## Breaking Changes

### Signature Changes

#### Funcs

1. DomainEventSubscriptionsClient.List
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string, *int32
	- Returns
		- From: EventSubscriptionsListResult, error
		- To: EventSubscriptionsListResultPage, error
1. DomainEventSubscriptionsClient.ListPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string, *int32
1. DomainTopicEventSubscriptionsClient.List
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, string, *int32
	- Returns
		- From: EventSubscriptionsListResult, error
		- To: EventSubscriptionsListResultPage, error
1. DomainTopicEventSubscriptionsClient.ListPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, string, *int32
1. TopicEventSubscriptionsClient.List
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string, *int32
	- Returns
		- From: EventSubscriptionsListResult, error
		- To: EventSubscriptionsListResultPage, error
1. TopicEventSubscriptionsClient.ListPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string, *int32

## Additive Changes

### New Funcs

1. DomainEventSubscriptionsClient.ListComplete(context.Context, string, string, string, *int32) (EventSubscriptionsListResultIterator, error)
1. DomainTopicEventSubscriptionsClient.ListComplete(context.Context, string, string, string, string, *int32) (EventSubscriptionsListResultIterator, error)
1. TopicEventSubscriptionsClient.ListComplete(context.Context, string, string, string, *int32) (EventSubscriptionsListResultIterator, error)
