Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewPipelineRunListResultPage` parameter(s) have been changed from `(func(context.Context, PipelineRunListResult) (PipelineRunListResult, error))` to `(PipelineRunListResult, func(context.Context, PipelineRunListResult) (PipelineRunListResult, error))`
- Function `NewImportPipelineListResultPage` parameter(s) have been changed from `(func(context.Context, ImportPipelineListResult) (ImportPipelineListResult, error))` to `(ImportPipelineListResult, func(context.Context, ImportPipelineListResult) (ImportPipelineListResult, error))`
- Function `NewTaskRunListResultPage` parameter(s) have been changed from `(func(context.Context, TaskRunListResult) (TaskRunListResult, error))` to `(TaskRunListResult, func(context.Context, TaskRunListResult) (TaskRunListResult, error))`
- Function `NewPrivateLinkResourceListResultPage` parameter(s) have been changed from `(func(context.Context, PrivateLinkResourceListResult) (PrivateLinkResourceListResult, error))` to `(PrivateLinkResourceListResult, func(context.Context, PrivateLinkResourceListResult) (PrivateLinkResourceListResult, error))`
- Function `NewAgentPoolListResultPage` parameter(s) have been changed from `(func(context.Context, AgentPoolListResult) (AgentPoolListResult, error))` to `(AgentPoolListResult, func(context.Context, AgentPoolListResult) (AgentPoolListResult, error))`
- Function `NewReplicationListResultPage` parameter(s) have been changed from `(func(context.Context, ReplicationListResult) (ReplicationListResult, error))` to `(ReplicationListResult, func(context.Context, ReplicationListResult) (ReplicationListResult, error))`
- Function `NewRunListResultPage` parameter(s) have been changed from `(func(context.Context, RunListResult) (RunListResult, error))` to `(RunListResult, func(context.Context, RunListResult) (RunListResult, error))`
- Function `NewPrivateEndpointConnectionListResultPage` parameter(s) have been changed from `(func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error))` to `(PrivateEndpointConnectionListResult, func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error))`
- Function `NewWebhookListResultPage` parameter(s) have been changed from `(func(context.Context, WebhookListResult) (WebhookListResult, error))` to `(WebhookListResult, func(context.Context, WebhookListResult) (WebhookListResult, error))`
- Function `NewTaskListResultPage` parameter(s) have been changed from `(func(context.Context, TaskListResult) (TaskListResult, error))` to `(TaskListResult, func(context.Context, TaskListResult) (TaskListResult, error))`
- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewExportPipelineListResultPage` parameter(s) have been changed from `(func(context.Context, ExportPipelineListResult) (ExportPipelineListResult, error))` to `(ExportPipelineListResult, func(context.Context, ExportPipelineListResult) (ExportPipelineListResult, error))`
- Function `NewEventListResultPage` parameter(s) have been changed from `(func(context.Context, EventListResult) (EventListResult, error))` to `(EventListResult, func(context.Context, EventListResult) (EventListResult, error))`
- Function `NewTokenListResultPage` parameter(s) have been changed from `(func(context.Context, TokenListResult) (TokenListResult, error))` to `(TokenListResult, func(context.Context, TokenListResult) (TokenListResult, error))`
- Function `NewScopeMapListResultPage` parameter(s) have been changed from `(func(context.Context, ScopeMapListResult) (ScopeMapListResult, error))` to `(ScopeMapListResult, func(context.Context, ScopeMapListResult) (ScopeMapListResult, error))`
- Function `NewRegistryListResultPage` parameter(s) have been changed from `(func(context.Context, RegistryListResult) (RegistryListResult, error))` to `(RegistryListResult, func(context.Context, RegistryListResult) (RegistryListResult, error))`

## New Content

- New const `LastModifiedByTypeUser`
- New const `LastModifiedByTypeKey`
- New const `ManagedIdentity`
- New const `Application`
- New const `LastModifiedByTypeApplication`
- New const `Key`
- New const `LastModifiedByTypeManagedIdentity`
- New const `User`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New function `PossibleLastModifiedByTypeValues() []LastModifiedByType`
- New struct `SystemData`
- New field `LogTemplate` in struct `FileTaskRunRequest`
- New field `SystemData` in struct `ImportPipeline`
- New field `LogTemplate` in struct `TaskRunRequest`
- New field `SystemData` in struct `ProxyResource`
- New field `SystemData` in struct `Resource`
- New field `SystemData` in struct `TaskRun`
- New field `LogTemplate` in struct `EncodedTaskRunRequest`
- New field `SystemData` in struct `Replication`
- New field `SystemData` in struct `Run`
- New field `SystemData` in struct `Task`
- New field `SystemData` in struct `Webhook`
- New field `SystemData` in struct `PrivateEndpointConnection`
- New field `SystemData` in struct `AgentPool`
- New field `LogTemplate` in struct `TaskProperties`
- New field `IsSystemTask` in struct `TaskProperties`
- New field `SystemData` in struct `Token`
- New field `LogTemplate` in struct `RunRequest`
- New field `SystemData` in struct `PipelineRun`
- New field `LastKeyRotationTimestamp` in struct `KeyVaultProperties`
- New field `KeyRotationEnabled` in struct `KeyVaultProperties`
- New field `SystemData` in struct `ExportPipeline`
- New field `LogTemplate` in struct `TaskPropertiesUpdateParameters`
- New field `SystemData` in struct `ScopeMap`
- New field `SystemData` in struct `Registry`
- New field `LogTemplate` in struct `DockerBuildRequest`
- New field `LogArtifact` in struct `RunProperties`
