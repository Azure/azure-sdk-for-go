# Unreleased

## Additive Changes

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

### Struct Changes

#### New Structs

1. NetworkRuleSetListResult
1. NetworkRuleSetListResultIterator
1. NetworkRuleSetListResultPage

#### New Struct Fields

1. SBNamespaceProperties.Status
