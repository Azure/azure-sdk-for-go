Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewApplicationListResultPage` parameter(s) have been changed from `(func(context.Context, ApplicationListResult) (ApplicationListResult, error))` to `(ApplicationListResult, func(context.Context, ApplicationListResult) (ApplicationListResult, error))`
- Function `NewScriptActionsListPage` parameter(s) have been changed from `(func(context.Context, ScriptActionsList) (ScriptActionsList, error))` to `(ScriptActionsList, func(context.Context, ScriptActionsList) (ScriptActionsList, error))`
- Function `NewScriptActionExecutionHistoryListPage` parameter(s) have been changed from `(func(context.Context, ScriptActionExecutionHistoryList) (ScriptActionExecutionHistoryList, error))` to `(ScriptActionExecutionHistoryList, func(context.Context, ScriptActionExecutionHistoryList) (ScriptActionExecutionHistoryList, error))`
- Function `NewClusterListResultPage` parameter(s) have been changed from `(func(context.Context, ClusterListResult) (ClusterListResult, error))` to `(ClusterListResult, func(context.Context, ClusterListResult) (ClusterListResult, error))`
- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Const `PublicLoadBalancer` has been removed
- Const `UDR` has been removed
- Const `InboundAndOutbound` has been removed
- Const `OutboundOnly` has been removed
- Function `PossiblePublicNetworkAccessValues` has been removed
- Function `PossibleOutboundOnlyPublicNetworkAccessTypeValues` has been removed
- Struct `NetworkSettings` has been removed
- Field `NetworkSettings` of struct `ClusterGetProperties` has been removed
- Field `NetworkSettings` of struct `ClusterCreateProperties` has been removed

## New Content

- New const `Outbound`
- New const `Inbound`
- New const `Enabled`
- New const `Disabled`
- New function `PossibleResourceProviderConnectionValues() []ResourceProviderConnection`
- New function `PossiblePrivateLinkValues() []PrivateLink`
- New struct `NetworkProperties`
- New field `NetworkProperties` in struct `ClusterCreateProperties`
- New field `ClusterID` in struct `ClusterGetProperties`
- New field `NetworkProperties` in struct `ClusterGetProperties`
