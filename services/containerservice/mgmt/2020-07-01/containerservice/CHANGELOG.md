
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewManagedClusterListResultPage` signature has been changed from `(func(context.Context, ManagedClusterListResult) (ManagedClusterListResult, error))` to `(ManagedClusterListResult,func(context.Context, ManagedClusterListResult) (ManagedClusterListResult, error))`
- Function `NewListResultPage` signature has been changed from `(func(context.Context, ListResult) (ListResult, error))` to `(ListResult,func(context.Context, ListResult) (ListResult, error))`
- Function `NewOpenShiftManagedClusterListResultPage` signature has been changed from `(func(context.Context, OpenShiftManagedClusterListResult) (OpenShiftManagedClusterListResult, error))` to `(OpenShiftManagedClusterListResult,func(context.Context, OpenShiftManagedClusterListResult) (OpenShiftManagedClusterListResult, error))`
- Function `NewAgentPoolListResultPage` signature has been changed from `(func(context.Context, AgentPoolListResult) (AgentPoolListResult, error))` to `(AgentPoolListResult,func(context.Context, AgentPoolListResult) (AgentPoolListResult, error))`
- Function `*ManagedClustersUpgradeNodeImageVersionFuture.Result` has been removed
- Function `ManagedClustersClient.UpgradeNodeImageVersion` has been removed
- Function `ManagedClustersClient.UpgradeNodeImageVersionPreparer` has been removed
- Function `ManagedClustersClient.UpgradeNodeImageVersionSender` has been removed
- Function `ManagedClustersClient.UpgradeNodeImageVersionResponder` has been removed
- Struct `ManagedClustersUpgradeNodeImageVersionFuture` has been removed

## New Content

- Function `AgentPoolsClient.UpgradeNodeImageVersionSender(*http.Request) (AgentPoolsUpgradeNodeImageVersionFuture,error)` is added
- Function `AgentPoolsClient.UpgradeNodeImageVersion(context.Context,string,string,string) (AgentPoolsUpgradeNodeImageVersionFuture,error)` is added
- Function `AgentPoolsClient.UpgradeNodeImageVersionResponder(*http.Response) (AgentPool,error)` is added
- Function `*AgentPoolsUpgradeNodeImageVersionFuture.Result(AgentPoolsClient) (AgentPool,error)` is added
- Function `AgentPoolsClient.UpgradeNodeImageVersionPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Struct `AgentPoolsUpgradeNodeImageVersionFuture` is added

