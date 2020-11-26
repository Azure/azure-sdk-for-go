
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewClusterListResultPage` signature has been changed from `(func(context.Context, ClusterListResult) (ClusterListResult, error))` to `(ClusterListResult,func(context.Context, ClusterListResult) (ClusterListResult, error))`
- Function `NewScriptActionsListPage` signature has been changed from `(func(context.Context, ScriptActionsList) (ScriptActionsList, error))` to `(ScriptActionsList,func(context.Context, ScriptActionsList) (ScriptActionsList, error))`
- Function `NewScriptActionExecutionHistoryListPage` signature has been changed from `(func(context.Context, ScriptActionExecutionHistoryList) (ScriptActionExecutionHistoryList, error))` to `(ScriptActionExecutionHistoryList,func(context.Context, ScriptActionExecutionHistoryList) (ScriptActionExecutionHistoryList, error))`
- Function `NewApplicationListResultPage` signature has been changed from `(func(context.Context, ApplicationListResult) (ApplicationListResult, error))` to `(ApplicationListResult,func(context.Context, ApplicationListResult) (ApplicationListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Const `InboundAndOutbound` has been removed
- Const `PublicLoadBalancer` has been removed
- Const `OutboundOnly` has been removed
- Const `UDR` has been removed
- Function `PossibleOutboundOnlyPublicNetworkAccessTypeValues` has been removed
- Function `PossiblePublicNetworkAccessValues` has been removed
- Struct `NetworkSettings` has been removed
- Field `NetworkSettings` of struct `ClusterCreateProperties` has been removed
- Field `NetworkSettings` of struct `ClusterGetProperties` has been removed

## New Content

- Const `Disabled` is added
- Const `Outbound` is added
- Const `Inbound` is added
- Const `Enabled` is added
- Function `PossibleResourceProviderConnectionValues() []ResourceProviderConnection` is added
- Function `PossiblePrivateLinkValues() []PrivateLink` is added
- Struct `NetworkProperties` is added
- Field `ClusterID` is added to struct `ClusterGetProperties`
- Field `NetworkProperties` is added to struct `ClusterGetProperties`
- Field `NetworkProperties` is added to struct `ClusterCreateProperties`

