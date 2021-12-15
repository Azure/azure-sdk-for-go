# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. NamespacesClient.ListKeys
	- Returns
		- From: SharedAccessAuthorizationRuleListResult, error
		- To: ResourceListKeys, error
1. NamespacesClient.ListKeysResponder
	- Returns
		- From: SharedAccessAuthorizationRuleListResult, error
		- To: ResourceListKeys, error
