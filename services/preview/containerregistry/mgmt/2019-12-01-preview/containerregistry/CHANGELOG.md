
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewScopeMapListResultPage` signature has been changed from `(func(context.Context, ScopeMapListResult) (ScopeMapListResult, error))` to `(ScopeMapListResult,func(context.Context, ScopeMapListResult) (ScopeMapListResult, error))`
- Function `NewImportPipelineListResultPage` signature has been changed from `(func(context.Context, ImportPipelineListResult) (ImportPipelineListResult, error))` to `(ImportPipelineListResult,func(context.Context, ImportPipelineListResult) (ImportPipelineListResult, error))`
- Function `NewAgentPoolListResultPage` signature has been changed from `(func(context.Context, AgentPoolListResult) (AgentPoolListResult, error))` to `(AgentPoolListResult,func(context.Context, AgentPoolListResult) (AgentPoolListResult, error))`
- Function `NewTaskRunListResultPage` signature has been changed from `(func(context.Context, TaskRunListResult) (TaskRunListResult, error))` to `(TaskRunListResult,func(context.Context, TaskRunListResult) (TaskRunListResult, error))`
- Function `NewEventListResultPage` signature has been changed from `(func(context.Context, EventListResult) (EventListResult, error))` to `(EventListResult,func(context.Context, EventListResult) (EventListResult, error))`
- Function `NewReplicationListResultPage` signature has been changed from `(func(context.Context, ReplicationListResult) (ReplicationListResult, error))` to `(ReplicationListResult,func(context.Context, ReplicationListResult) (ReplicationListResult, error))`
- Function `NewTokenListResultPage` signature has been changed from `(func(context.Context, TokenListResult) (TokenListResult, error))` to `(TokenListResult,func(context.Context, TokenListResult) (TokenListResult, error))`
- Function `NewRegistryListResultPage` signature has been changed from `(func(context.Context, RegistryListResult) (RegistryListResult, error))` to `(RegistryListResult,func(context.Context, RegistryListResult) (RegistryListResult, error))`
- Function `NewWebhookListResultPage` signature has been changed from `(func(context.Context, WebhookListResult) (WebhookListResult, error))` to `(WebhookListResult,func(context.Context, WebhookListResult) (WebhookListResult, error))`
- Function `NewPrivateEndpointConnectionListResultPage` signature has been changed from `(func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error))` to `(PrivateEndpointConnectionListResult,func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error))`
- Function `NewTaskListResultPage` signature has been changed from `(func(context.Context, TaskListResult) (TaskListResult, error))` to `(TaskListResult,func(context.Context, TaskListResult) (TaskListResult, error))`
- Function `NewPrivateLinkResourceListResultPage` signature has been changed from `(func(context.Context, PrivateLinkResourceListResult) (PrivateLinkResourceListResult, error))` to `(PrivateLinkResourceListResult,func(context.Context, PrivateLinkResourceListResult) (PrivateLinkResourceListResult, error))`
- Function `NewPipelineRunListResultPage` signature has been changed from `(func(context.Context, PipelineRunListResult) (PipelineRunListResult, error))` to `(PipelineRunListResult,func(context.Context, PipelineRunListResult) (PipelineRunListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewRunListResultPage` signature has been changed from `(func(context.Context, RunListResult) (RunListResult, error))` to `(RunListResult,func(context.Context, RunListResult) (RunListResult, error))`
- Function `NewExportPipelineListResultPage` signature has been changed from `(func(context.Context, ExportPipelineListResult) (ExportPipelineListResult, error))` to `(ExportPipelineListResult,func(context.Context, ExportPipelineListResult) (ExportPipelineListResult, error))`

## New Content

- Const `User` is added
- Const `ManagedIdentity` is added
- Const `Application` is added
- Const `LastModifiedByTypeUser` is added
- Const `LastModifiedByTypeManagedIdentity` is added
- Const `LastModifiedByTypeKey` is added
- Const `LastModifiedByTypeApplication` is added
- Const `Key` is added
- Function `PossibleLastModifiedByTypeValues() []LastModifiedByType` is added
- Function `PossibleCreatedByTypeValues() []CreatedByType` is added
- Struct `SystemData` is added
- Field `SystemData` is added to struct `Replication`
- Field `SystemData` is added to struct `Token`
- Field `SystemData` is added to struct `ImportPipeline`
- Field `SystemData` is added to struct `Webhook`
- Field `SystemData` is added to struct `TaskRun`
- Field `SystemData` is added to struct `AgentPool`
- Field `SystemData` is added to struct `Resource`
- Field `LogTemplate` is added to struct `TaskRunRequest`
- Field `SystemData` is added to struct `Run`
- Field `LogTemplate` is added to struct `EncodedTaskRunRequest`
- Field `LogTemplate` is added to struct `RunRequest`
- Field `SystemData` is added to struct `Registry`
- Field `LogTemplate` is added to struct `TaskPropertiesUpdateParameters`
- Field `LogTemplate` is added to struct `DockerBuildRequest`
- Field `SystemData` is added to struct `ExportPipeline`
- Field `IsSystemTask` is added to struct `TaskProperties`
- Field `LogTemplate` is added to struct `TaskProperties`
- Field `LogArtifact` is added to struct `RunProperties`
- Field `SystemData` is added to struct `ScopeMap`
- Field `SystemData` is added to struct `PipelineRun`
- Field `SystemData` is added to struct `PrivateEndpointConnection`
- Field `SystemData` is added to struct `ProxyResource`
- Field `SystemData` is added to struct `Task`
- Field `KeyRotationEnabled` is added to struct `KeyVaultProperties`
- Field `LastKeyRotationTimestamp` is added to struct `KeyVaultProperties`
- Field `LogTemplate` is added to struct `FileTaskRunRequest`

