Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/containerregistry/resource-manager/readme.md tag: `package-2019-12-preview`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *AgentPoolsCreateFuture.Result(AgentPoolsClient) (AgentPool, error)
1. *AgentPoolsDeleteFuture.Result(AgentPoolsClient) (autorest.Response, error)
1. *AgentPoolsUpdateFuture.Result(AgentPoolsClient) (AgentPool, error)
1. *ExportPipelinesCreateFuture.Result(ExportPipelinesClient) (ExportPipeline, error)
1. *ExportPipelinesDeleteFuture.Result(ExportPipelinesClient) (autorest.Response, error)
1. *ImportPipelinesCreateFuture.Result(ImportPipelinesClient) (ImportPipeline, error)
1. *ImportPipelinesDeleteFuture.Result(ImportPipelinesClient) (autorest.Response, error)
1. *PipelineRunsCreateFuture.Result(PipelineRunsClient) (PipelineRun, error)
1. *PipelineRunsDeleteFuture.Result(PipelineRunsClient) (autorest.Response, error)
1. *PrivateEndpointConnectionsCreateOrUpdateFuture.Result(PrivateEndpointConnectionsClient) (PrivateEndpointConnection, error)
1. *PrivateEndpointConnectionsDeleteFuture.Result(PrivateEndpointConnectionsClient) (autorest.Response, error)
1. *RegistriesCreateFuture.Result(RegistriesClient) (Registry, error)
1. *RegistriesDeleteFuture.Result(RegistriesClient) (autorest.Response, error)
1. *RegistriesGenerateCredentialsFuture.Result(RegistriesClient) (GenerateCredentialsResult, error)
1. *RegistriesImportImageFuture.Result(RegistriesClient) (autorest.Response, error)
1. *RegistriesScheduleRunFuture.Result(RegistriesClient) (Run, error)
1. *RegistriesUpdateFuture.Result(RegistriesClient) (Registry, error)
1. *ReplicationsCreateFuture.Result(ReplicationsClient) (Replication, error)
1. *ReplicationsDeleteFuture.Result(ReplicationsClient) (autorest.Response, error)
1. *ReplicationsUpdateFuture.Result(ReplicationsClient) (Replication, error)
1. *RunsCancelFuture.Result(RunsClient) (autorest.Response, error)
1. *RunsUpdateFuture.Result(RunsClient) (Run, error)
1. *ScopeMapsCreateFuture.Result(ScopeMapsClient) (ScopeMap, error)
1. *ScopeMapsDeleteFuture.Result(ScopeMapsClient) (autorest.Response, error)
1. *ScopeMapsUpdateFuture.Result(ScopeMapsClient) (ScopeMap, error)
1. *TaskRunsCreateFuture.Result(TaskRunsClient) (TaskRun, error)
1. *TaskRunsDeleteFuture.Result(TaskRunsClient) (autorest.Response, error)
1. *TaskRunsUpdateFuture.Result(TaskRunsClient) (TaskRun, error)
1. *TasksCreateFuture.Result(TasksClient) (Task, error)
1. *TasksDeleteFuture.Result(TasksClient) (autorest.Response, error)
1. *TasksUpdateFuture.Result(TasksClient) (Task, error)
1. *TokensCreateFuture.Result(TokensClient) (Token, error)
1. *TokensDeleteFuture.Result(TokensClient) (autorest.Response, error)
1. *TokensUpdateFuture.Result(TokensClient) (Token, error)
1. *WebhooksCreateFuture.Result(WebhooksClient) (Webhook, error)
1. *WebhooksDeleteFuture.Result(WebhooksClient) (autorest.Response, error)
1. *WebhooksUpdateFuture.Result(WebhooksClient) (Webhook, error)

## Struct Changes

### Removed Struct Fields

1. AgentPoolsCreateFuture.azure.Future
1. AgentPoolsDeleteFuture.azure.Future
1. AgentPoolsUpdateFuture.azure.Future
1. ExportPipelinesCreateFuture.azure.Future
1. ExportPipelinesDeleteFuture.azure.Future
1. ImportPipelinesCreateFuture.azure.Future
1. ImportPipelinesDeleteFuture.azure.Future
1. PipelineRunsCreateFuture.azure.Future
1. PipelineRunsDeleteFuture.azure.Future
1. PrivateEndpointConnectionsCreateOrUpdateFuture.azure.Future
1. PrivateEndpointConnectionsDeleteFuture.azure.Future
1. RegistriesCreateFuture.azure.Future
1. RegistriesDeleteFuture.azure.Future
1. RegistriesGenerateCredentialsFuture.azure.Future
1. RegistriesImportImageFuture.azure.Future
1. RegistriesScheduleRunFuture.azure.Future
1. RegistriesUpdateFuture.azure.Future
1. ReplicationsCreateFuture.azure.Future
1. ReplicationsDeleteFuture.azure.Future
1. ReplicationsUpdateFuture.azure.Future
1. RunsCancelFuture.azure.Future
1. RunsUpdateFuture.azure.Future
1. ScopeMapsCreateFuture.azure.Future
1. ScopeMapsDeleteFuture.azure.Future
1. ScopeMapsUpdateFuture.azure.Future
1. TaskRunsCreateFuture.azure.Future
1. TaskRunsDeleteFuture.azure.Future
1. TaskRunsUpdateFuture.azure.Future
1. TasksCreateFuture.azure.Future
1. TasksDeleteFuture.azure.Future
1. TasksUpdateFuture.azure.Future
1. TokensCreateFuture.azure.Future
1. TokensDeleteFuture.azure.Future
1. TokensUpdateFuture.azure.Future
1. WebhooksCreateFuture.azure.Future
1. WebhooksDeleteFuture.azure.Future
1. WebhooksUpdateFuture.azure.Future

## Struct Changes

### New Struct Fields

1. AgentPoolsCreateFuture.Result
1. AgentPoolsCreateFuture.azure.FutureAPI
1. AgentPoolsDeleteFuture.Result
1. AgentPoolsDeleteFuture.azure.FutureAPI
1. AgentPoolsUpdateFuture.Result
1. AgentPoolsUpdateFuture.azure.FutureAPI
1. ExportPipelinesCreateFuture.Result
1. ExportPipelinesCreateFuture.azure.FutureAPI
1. ExportPipelinesDeleteFuture.Result
1. ExportPipelinesDeleteFuture.azure.FutureAPI
1. ImportPipelinesCreateFuture.Result
1. ImportPipelinesCreateFuture.azure.FutureAPI
1. ImportPipelinesDeleteFuture.Result
1. ImportPipelinesDeleteFuture.azure.FutureAPI
1. PipelineRunsCreateFuture.Result
1. PipelineRunsCreateFuture.azure.FutureAPI
1. PipelineRunsDeleteFuture.Result
1. PipelineRunsDeleteFuture.azure.FutureAPI
1. PrivateEndpointConnectionsCreateOrUpdateFuture.Result
1. PrivateEndpointConnectionsCreateOrUpdateFuture.azure.FutureAPI
1. PrivateEndpointConnectionsDeleteFuture.Result
1. PrivateEndpointConnectionsDeleteFuture.azure.FutureAPI
1. RegistriesCreateFuture.Result
1. RegistriesCreateFuture.azure.FutureAPI
1. RegistriesDeleteFuture.Result
1. RegistriesDeleteFuture.azure.FutureAPI
1. RegistriesGenerateCredentialsFuture.Result
1. RegistriesGenerateCredentialsFuture.azure.FutureAPI
1. RegistriesImportImageFuture.Result
1. RegistriesImportImageFuture.azure.FutureAPI
1. RegistriesScheduleRunFuture.Result
1. RegistriesScheduleRunFuture.azure.FutureAPI
1. RegistriesUpdateFuture.Result
1. RegistriesUpdateFuture.azure.FutureAPI
1. ReplicationsCreateFuture.Result
1. ReplicationsCreateFuture.azure.FutureAPI
1. ReplicationsDeleteFuture.Result
1. ReplicationsDeleteFuture.azure.FutureAPI
1. ReplicationsUpdateFuture.Result
1. ReplicationsUpdateFuture.azure.FutureAPI
1. RunsCancelFuture.Result
1. RunsCancelFuture.azure.FutureAPI
1. RunsUpdateFuture.Result
1. RunsUpdateFuture.azure.FutureAPI
1. ScopeMapsCreateFuture.Result
1. ScopeMapsCreateFuture.azure.FutureAPI
1. ScopeMapsDeleteFuture.Result
1. ScopeMapsDeleteFuture.azure.FutureAPI
1. ScopeMapsUpdateFuture.Result
1. ScopeMapsUpdateFuture.azure.FutureAPI
1. TaskRunsCreateFuture.Result
1. TaskRunsCreateFuture.azure.FutureAPI
1. TaskRunsDeleteFuture.Result
1. TaskRunsDeleteFuture.azure.FutureAPI
1. TaskRunsUpdateFuture.Result
1. TaskRunsUpdateFuture.azure.FutureAPI
1. TasksCreateFuture.Result
1. TasksCreateFuture.azure.FutureAPI
1. TasksDeleteFuture.Result
1. TasksDeleteFuture.azure.FutureAPI
1. TasksUpdateFuture.Result
1. TasksUpdateFuture.azure.FutureAPI
1. TokensCreateFuture.Result
1. TokensCreateFuture.azure.FutureAPI
1. TokensDeleteFuture.Result
1. TokensDeleteFuture.azure.FutureAPI
1. TokensUpdateFuture.Result
1. TokensUpdateFuture.azure.FutureAPI
1. WebhooksCreateFuture.Result
1. WebhooksCreateFuture.azure.FutureAPI
1. WebhooksDeleteFuture.Result
1. WebhooksDeleteFuture.azure.FutureAPI
1. WebhooksUpdateFuture.Result
1. WebhooksUpdateFuture.azure.FutureAPI
