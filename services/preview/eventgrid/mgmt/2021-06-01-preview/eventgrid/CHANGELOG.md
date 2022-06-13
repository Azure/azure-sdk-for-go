# Unreleased

## Breaking Changes

### Signature Changes

#### Struct Fields

1. TopicTypeProperties.SupportedScopesForSource changed type from *[]string to *[]TopicTypeSourceScope

## Additive Changes

### New Constants

1. TopicTypeSourceScope.TopicTypeSourceScopeAzureSubscription
1. TopicTypeSourceScope.TopicTypeSourceScopeResource
1. TopicTypeSourceScope.TopicTypeSourceScopeResourceGroup

### New Funcs

1. PossibleTopicTypeSourceScopeValues() []TopicTypeSourceScope
