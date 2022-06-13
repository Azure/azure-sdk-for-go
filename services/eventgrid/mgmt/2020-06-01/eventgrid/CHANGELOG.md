# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. TopicsClient.RegenerateKey
	- Returns
		- From: TopicSharedAccessKeys, error
		- To: TopicsRegenerateKeyFuture, error
1. TopicsClient.RegenerateKeySender
	- Returns
		- From: *http.Response, error
		- To: TopicsRegenerateKeyFuture, error

## Additive Changes

### New Constants

1. CreatedByType.Application
1. CreatedByType.Key
1. CreatedByType.ManagedIdentity
1. CreatedByType.User

### New Funcs

1. *TopicsRegenerateKeyFuture.UnmarshalJSON([]byte) error
1. DomainTopicProperties.MarshalJSON() ([]byte, error)
1. PossibleCreatedByTypeValues() []CreatedByType

### Struct Changes

#### New Structs

1. SystemData
1. TopicsRegenerateKeyFuture

#### New Struct Fields

1. Domain.SystemData
1. DomainTopic.SystemData
1. EventSubscription.SystemData
1. Topic.SystemData
