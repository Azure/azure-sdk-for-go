# Release History

## 3.0.0-beta.2 (2026-02-23)
### Breaking Changes

- Enum `Architecture` has been removed
- Enum `BaseImageDependencyType` has been removed
- Enum `BaseImageTriggerType` has been removed
- Enum `OS` has been removed
- Enum `RunStatus` has been removed
- Enum `RunType` has been removed
- Enum `SecretObjectType` has been removed
- Enum `SourceControlType` has been removed
- Enum `SourceRegistryLoginMode` has been removed
- Enum `SourceTriggerEvent` has been removed
- Enum `StepType` has been removed
- Enum `TaskStatus` has been removed
- Enum `TokenType` has been removed
- Enum `UpdateTriggerPayloadType` has been removed
- Enum `Variant` has been removed
- Function `NewAgentPoolsClient` has been removed
- Function `*AgentPoolsClient.BeginCreate` has been removed
- Function `*AgentPoolsClient.BeginDelete` has been removed
- Function `*AgentPoolsClient.Get` has been removed
- Function `*AgentPoolsClient.GetQueueStatus` has been removed
- Function `*AgentPoolsClient.NewListPager` has been removed
- Function `*AgentPoolsClient.BeginUpdate` has been removed
- Function `*ClientFactory.NewAgentPoolsClient` has been removed
- Function `*ClientFactory.NewRunsClient` has been removed
- Function `*ClientFactory.NewTaskRunsClient` has been removed
- Function `*ClientFactory.NewTasksClient` has been removed
- Function `*DockerBuildRequest.GetRunRequest` has been removed
- Function `*DockerBuildStep.GetTaskStepProperties` has been removed
- Function `*DockerBuildStepUpdateParameters.GetTaskStepUpdateParameters` has been removed
- Function `*EncodedTaskRunRequest.GetRunRequest` has been removed
- Function `*EncodedTaskStep.GetTaskStepProperties` has been removed
- Function `*EncodedTaskStepUpdateParameters.GetTaskStepUpdateParameters` has been removed
- Function `*FileTaskRunRequest.GetRunRequest` has been removed
- Function `*FileTaskStep.GetTaskStepProperties` has been removed
- Function `*FileTaskStepUpdateParameters.GetTaskStepUpdateParameters` has been removed
- Function `*RegistriesClient.GetBuildSourceUploadURL` has been removed
- Function `*RegistriesClient.ScheduleRun` has been removed
- Function `*RunRequest.GetRunRequest` has been removed
- Function `NewRunsClient` has been removed
- Function `*RunsClient.Cancel` has been removed
- Function `*RunsClient.Get` has been removed
- Function `*RunsClient.GetLogSasURL` has been removed
- Function `*RunsClient.NewListPager` has been removed
- Function `*RunsClient.Update` has been removed
- Function `*TaskRunRequest.GetRunRequest` has been removed
- Function `NewTaskRunsClient` has been removed
- Function `*TaskRunsClient.BeginCreate` has been removed
- Function `*TaskRunsClient.Delete` has been removed
- Function `*TaskRunsClient.Get` has been removed
- Function `*TaskRunsClient.GetDetails` has been removed
- Function `*TaskRunsClient.NewListPager` has been removed
- Function `*TaskRunsClient.BeginUpdate` has been removed
- Function `*TaskStepProperties.GetTaskStepProperties` has been removed
- Function `*TaskStepUpdateParameters.GetTaskStepUpdateParameters` has been removed
- Function `NewTasksClient` has been removed
- Function `*TasksClient.Create` has been removed
- Function `*TasksClient.Delete` has been removed
- Function `*TasksClient.Get` has been removed
- Function `*TasksClient.GetDetails` has been removed
- Function `*TasksClient.NewListPager` has been removed
- Function `*TasksClient.Update` has been removed
- Struct `AgentPool` has been removed
- Struct `AgentPoolListResult` has been removed
- Struct `AgentPoolProperties` has been removed
- Struct `AgentPoolPropertiesUpdateParameters` has been removed
- Struct `AgentPoolQueueStatus` has been removed
- Struct `AgentPoolUpdateParameters` has been removed
- Struct `AgentProperties` has been removed
- Struct `Argument` has been removed
- Struct `AuthInfo` has been removed
- Struct `AuthInfoUpdateParameters` has been removed
- Struct `BaseImageDependency` has been removed
- Struct `BaseImageTrigger` has been removed
- Struct `BaseImageTriggerUpdateParameters` has been removed
- Struct `Credentials` has been removed
- Struct `CustomRegistryCredentials` has been removed
- Struct `DockerBuildRequest` has been removed
- Struct `DockerBuildStep` has been removed
- Struct `DockerBuildStepUpdateParameters` has been removed
- Struct `EncodedTaskRunRequest` has been removed
- Struct `EncodedTaskStep` has been removed
- Struct `EncodedTaskStepUpdateParameters` has been removed
- Struct `FileTaskRunRequest` has been removed
- Struct `FileTaskStep` has been removed
- Struct `FileTaskStepUpdateParameters` has been removed
- Struct `ImageDescriptor` has been removed
- Struct `ImageUpdateTrigger` has been removed
- Struct `OverrideTaskStepProperties` has been removed
- Struct `PlatformProperties` has been removed
- Struct `PlatformUpdateParameters` has been removed
- Struct `Run` has been removed
- Struct `RunGetLogResult` has been removed
- Struct `RunListResult` has been removed
- Struct `RunProperties` has been removed
- Struct `RunUpdateParameters` has been removed
- Struct `SecretObject` has been removed
- Struct `SetValue` has been removed
- Struct `SourceProperties` has been removed
- Struct `SourceRegistryCredentials` has been removed
- Struct `SourceTrigger` has been removed
- Struct `SourceTriggerDescriptor` has been removed
- Struct `SourceTriggerUpdateParameters` has been removed
- Struct `SourceUpdateParameters` has been removed
- Struct `SourceUploadDefinition` has been removed
- Struct `Task` has been removed
- Struct `TaskListResult` has been removed
- Struct `TaskProperties` has been removed
- Struct `TaskPropertiesUpdateParameters` has been removed
- Struct `TaskRun` has been removed
- Struct `TaskRunListResult` has been removed
- Struct `TaskRunProperties` has been removed
- Struct `TaskRunPropertiesUpdateParameters` has been removed
- Struct `TaskRunRequest` has been removed
- Struct `TaskRunUpdateParameters` has been removed
- Struct `TaskUpdateParameters` has been removed
- Struct `TimerTrigger` has been removed
- Struct `TimerTriggerDescriptor` has been removed
- Struct `TimerTriggerUpdateParameters` has been removed
- Struct `TriggerProperties` has been removed
- Struct `TriggerUpdateParameters` has been removed

### Features Added

- New enum type `EndpointProtocol` with values `EndpointProtocolIPv4`, `EndpointProtocolIPv4AndIPv6`
- New enum type `MetadataSearch` with values `MetadataSearchDisabled`, `MetadataSearchEnabled`
- New enum type `PackageSourceType` with values `PackageSourceTypeRemote`
- New enum type `PipelineOptions` with values `PipelineOptionsContinueOnErrors`, `PipelineOptionsDeleteSourceBlobOnSuccess`, `PipelineOptionsOverwriteBlobs`, `PipelineOptionsOverwriteTags`
- New enum type `PipelineRunSourceType` with values `PipelineRunSourceTypeAzureStorageBlob`
- New enum type `PipelineRunTargetType` with values `PipelineRunTargetTypeAzureStorageBlob`
- New enum type `PipelineSourceType` with values `PipelineSourceTypeAzureStorageBlobContainer`
- New enum type `RegionalEndpoints` with values `RegionalEndpointsDisabled`, `RegionalEndpointsEnabled`
- New enum type `StorageAccessMode` with values `StorageAccessModeManagedIdentity`, `StorageAccessModeSasToken`
- New enum type `SyncState` with values `SyncStateFailed`, `SyncStateNotActivated`, `SyncStatePending`, `SyncStateSucceeded`, `SyncStateSyncing`, `SyncStateTimedOut`
- New enum type `SyncTrigger` with values `SyncTriggerInitialSync`, `SyncTriggerManualResync`, `SyncTriggerRecoveryService`, `SyncTriggerSyncTokenUpdate`
- New function `NewArchiveVersionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ArchiveVersionsClient, error)`
- New function `*ArchiveVersionsClient.BeginCreate(ctx context.Context, resourceGroupName string, registryName string, packageType string, archiveName string, archiveVersionName string, options *ArchiveVersionsClientBeginCreateOptions) (*runtime.Poller[ArchiveVersionsClientCreateResponse], error)`
- New function `*ArchiveVersionsClient.BeginDelete(ctx context.Context, resourceGroupName string, registryName string, packageType string, archiveName string, archiveVersionName string, options *ArchiveVersionsClientBeginDeleteOptions) (*runtime.Poller[ArchiveVersionsClientDeleteResponse], error)`
- New function `*ArchiveVersionsClient.Get(ctx context.Context, resourceGroupName string, registryName string, packageType string, archiveName string, archiveVersionName string, options *ArchiveVersionsClientGetOptions) (ArchiveVersionsClientGetResponse, error)`
- New function `*ArchiveVersionsClient.NewListPager(resourceGroupName string, registryName string, packageType string, archiveName string, options *ArchiveVersionsClientListOptions) *runtime.Pager[ArchiveVersionsClientListResponse]`
- New function `NewArchivesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ArchivesClient, error)`
- New function `*ArchivesClient.BeginCreate(ctx context.Context, resourceGroupName string, registryName string, packageType string, archiveName string, archiveCreateParameters Archive, options *ArchivesClientBeginCreateOptions) (*runtime.Poller[ArchivesClientCreateResponse], error)`
- New function `*ArchivesClient.BeginDelete(ctx context.Context, resourceGroupName string, registryName string, packageType string, archiveName string, options *ArchivesClientBeginDeleteOptions) (*runtime.Poller[ArchivesClientDeleteResponse], error)`
- New function `*ArchivesClient.Get(ctx context.Context, resourceGroupName string, registryName string, packageType string, archiveName string, options *ArchivesClientGetOptions) (ArchivesClientGetResponse, error)`
- New function `*ArchivesClient.NewListPager(resourceGroupName string, registryName string, packageType string, options *ArchivesClientListOptions) *runtime.Pager[ArchivesClientListResponse]`
- New function `*ArchivesClient.Update(ctx context.Context, resourceGroupName string, registryName string, packageType string, archiveName string, archiveUpdateParameters ArchiveUpdateParameters, options *ArchivesClientUpdateOptions) (ArchivesClientUpdateResponse, error)`
- New function `*ClientFactory.NewArchiveVersionsClient() *ArchiveVersionsClient`
- New function `*ClientFactory.NewArchivesClient() *ArchivesClient`
- New function `*ClientFactory.NewExportPipelinesClient() *ExportPipelinesClient`
- New function `*ClientFactory.NewImportPipelinesClient() *ImportPipelinesClient`
- New function `*ClientFactory.NewPipelineRunsClient() *PipelineRunsClient`
- New function `*ConnectedRegistriesClient.Resync(ctx context.Context, resourceGroupName string, registryName string, connectedRegistryName string, options *ConnectedRegistriesClientResyncOptions) (ConnectedRegistriesClientResyncResponse, error)`
- New function `NewExportPipelinesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ExportPipelinesClient, error)`
- New function `*ExportPipelinesClient.BeginCreate(ctx context.Context, resourceGroupName string, registryName string, exportPipelineName string, exportPipelineCreateParameters ExportPipeline, options *ExportPipelinesClientBeginCreateOptions) (*runtime.Poller[ExportPipelinesClientCreateResponse], error)`
- New function `*ExportPipelinesClient.BeginDelete(ctx context.Context, resourceGroupName string, registryName string, exportPipelineName string, options *ExportPipelinesClientBeginDeleteOptions) (*runtime.Poller[ExportPipelinesClientDeleteResponse], error)`
- New function `*ExportPipelinesClient.Get(ctx context.Context, resourceGroupName string, registryName string, exportPipelineName string, options *ExportPipelinesClientGetOptions) (ExportPipelinesClientGetResponse, error)`
- New function `*ExportPipelinesClient.NewListPager(resourceGroupName string, registryName string, options *ExportPipelinesClientListOptions) *runtime.Pager[ExportPipelinesClientListResponse]`
- New function `NewImportPipelinesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ImportPipelinesClient, error)`
- New function `*ImportPipelinesClient.BeginCreate(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, importPipelineCreateParameters ImportPipeline, options *ImportPipelinesClientBeginCreateOptions) (*runtime.Poller[ImportPipelinesClientCreateResponse], error)`
- New function `*ImportPipelinesClient.BeginDelete(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *ImportPipelinesClientBeginDeleteOptions) (*runtime.Poller[ImportPipelinesClientDeleteResponse], error)`
- New function `*ImportPipelinesClient.Get(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *ImportPipelinesClientGetOptions) (ImportPipelinesClientGetResponse, error)`
- New function `*ImportPipelinesClient.NewListPager(resourceGroupName string, registryName string, options *ImportPipelinesClientListOptions) *runtime.Pager[ImportPipelinesClientListResponse]`
- New function `NewPipelineRunsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PipelineRunsClient, error)`
- New function `*PipelineRunsClient.BeginCreate(ctx context.Context, resourceGroupName string, registryName string, pipelineRunName string, pipelineRunCreateParameters PipelineRun, options *PipelineRunsClientBeginCreateOptions) (*runtime.Poller[PipelineRunsClientCreateResponse], error)`
- New function `*PipelineRunsClient.BeginDelete(ctx context.Context, resourceGroupName string, registryName string, pipelineRunName string, options *PipelineRunsClientBeginDeleteOptions) (*runtime.Poller[PipelineRunsClientDeleteResponse], error)`
- New function `*PipelineRunsClient.Get(ctx context.Context, resourceGroupName string, registryName string, pipelineRunName string, options *PipelineRunsClientGetOptions) (PipelineRunsClientGetResponse, error)`
- New function `*PipelineRunsClient.NewListPager(resourceGroupName string, registryName string, options *PipelineRunsClientListOptions) *runtime.Pager[PipelineRunsClientListResponse]`
- New struct `Archive`
- New struct `ArchiveListResult`
- New struct `ArchivePackageSourceProperties`
- New struct `ArchiveProperties`
- New struct `ArchiveUpdateParameters`
- New struct `ArchiveUpdateProperties`
- New struct `ArchiveVersion`
- New struct `ArchiveVersionListResult`
- New struct `ArchiveVersionProperties`
- New struct `ExportPipeline`
- New struct `ExportPipelineListResult`
- New struct `ExportPipelineProperties`
- New struct `ExportPipelineTargetProperties`
- New struct `ImportPipeline`
- New struct `ImportPipelineListResult`
- New struct `ImportPipelineProperties`
- New struct `ImportPipelineSourceProperties`
- New struct `PipelineRun`
- New struct `PipelineRunListResult`
- New struct `PipelineRunProperties`
- New struct `PipelineRunRequest`
- New struct `PipelineRunResponse`
- New struct `PipelineRunSourceProperties`
- New struct `PipelineRunTargetProperties`
- New struct `PipelineSourceTriggerDescriptor`
- New struct `PipelineSourceTriggerProperties`
- New struct `PipelineTriggerDescriptor`
- New struct `PipelineTriggerProperties`
- New struct `ProgressProperties`
- New struct `RegistrySyncResult`
- New struct `SoftDeletePolicy`
- New field `Identity` in struct `CacheRule`
- New field `Identity` in struct `CacheRuleUpdateParameters`
- New field `RegistrySyncResult` in struct `ConnectedRegistryProperties`
- New field `SoftDeletePolicy` in struct `Policies`
- New field `EndpointProtocol`, `MetadataSearch`, `RegionalEndpointHostNames`, `RegionalEndpoints` in struct `RegistryProperties`
- New field `EndpointProtocol`, `MetadataSearch`, `RegionalEndpoints` in struct `RegistryPropertiesUpdateParameters`


## 3.0.0 (2026-01-28)
### Breaking Changes
- Registry Task related clients, functions and types has been moved to new package: ‎`sdk/resourcemanager/containerregistry/armcontainerregistrytasks`.
- Field `AutoGeneratedDomainNameLabelScope`, `ResourceGroupName` of struct `RegistryNameCheckRequest` has been removed
- Field `AvailableLoginServerName` of struct `RegistryNameStatus` has been removed
- Field `AutoGeneratedDomainNameLabelScope` of struct `RegistryProperties` has been removed


## 2.0.0 (2025-11-05)
### Breaking Changes

- Type of `SystemData.LastModifiedByType` has been changed from `*LastModifiedByType` to `*CreatedByType`
- Enum `LastModifiedByType` has been removed

### Features Added

- New enum type `ActivationStatus` with values `ActivationStatusActive`, `ActivationStatusInactive`
- New enum type `AuditLogStatus` with values `AuditLogStatusDisabled`, `AuditLogStatusEnabled`
- New enum type `AutoGeneratedDomainNameLabelScope` with values `AutoGeneratedDomainNameLabelScopeNoReuse`, `AutoGeneratedDomainNameLabelScopeResourceGroupReuse`, `AutoGeneratedDomainNameLabelScopeSubscriptionReuse`, `AutoGeneratedDomainNameLabelScopeTenantReuse`, `AutoGeneratedDomainNameLabelScopeUnsecure`
- New enum type `AzureADAuthenticationAsArmPolicyStatus` with values `AzureADAuthenticationAsArmPolicyStatusDisabled`, `AzureADAuthenticationAsArmPolicyStatusEnabled`
- New enum type `CertificateType` with values `CertificateTypeLocalDirectory`
- New enum type `ConnectedRegistryMode` with values `ConnectedRegistryModeMirror`, `ConnectedRegistryModeReadOnly`, `ConnectedRegistryModeReadWrite`, `ConnectedRegistryModeRegistry`
- New enum type `ConnectionState` with values `ConnectionStateOffline`, `ConnectionStateOnline`, `ConnectionStateSyncing`, `ConnectionStateUnhealthy`
- New enum type `LogLevel` with values `LogLevelDebug`, `LogLevelError`, `LogLevelInformation`, `LogLevelNone`, `LogLevelWarning`
- New enum type `RoleAssignmentMode` with values `RoleAssignmentModeAbacRepositoryPermissions`, `RoleAssignmentModeLegacyRegistryPermissions`
- New enum type `TLSStatus` with values `TLSStatusDisabled`, `TLSStatusEnabled`
- New function `*ClientFactory.NewConnectedRegistriesClient() *ConnectedRegistriesClient`
- New function `NewConnectedRegistriesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedRegistriesClient, error)`
- New function `*ConnectedRegistriesClient.BeginCreate(context.Context, string, string, string, ConnectedRegistry, *ConnectedRegistriesClientBeginCreateOptions) (*runtime.Poller[ConnectedRegistriesClientCreateResponse], error)`
- New function `*ConnectedRegistriesClient.BeginDeactivate(context.Context, string, string, string, *ConnectedRegistriesClientBeginDeactivateOptions) (*runtime.Poller[ConnectedRegistriesClientDeactivateResponse], error)`
- New function `*ConnectedRegistriesClient.BeginDelete(context.Context, string, string, string, *ConnectedRegistriesClientBeginDeleteOptions) (*runtime.Poller[ConnectedRegistriesClientDeleteResponse], error)`
- New function `*ConnectedRegistriesClient.Get(context.Context, string, string, string, *ConnectedRegistriesClientGetOptions) (ConnectedRegistriesClientGetResponse, error)`
- New function `*ConnectedRegistriesClient.NewListPager(string, string, *ConnectedRegistriesClientListOptions) *runtime.Pager[ConnectedRegistriesClientListResponse]`
- New function `*ConnectedRegistriesClient.BeginUpdate(context.Context, string, string, string, ConnectedRegistryUpdateParameters, *ConnectedRegistriesClientBeginUpdateOptions) (*runtime.Poller[ConnectedRegistriesClientUpdateResponse], error)`
- New struct `ActivationProperties`
- New struct `AzureADAuthenticationAsArmPolicy`
- New struct `ConnectedRegistry`
- New struct `ConnectedRegistryListResult`
- New struct `ConnectedRegistryProperties`
- New struct `ConnectedRegistryUpdateParameters`
- New struct `ConnectedRegistryUpdateProperties`
- New struct `GarbageCollectionProperties`
- New struct `LoggingProperties`
- New struct `LoginServerProperties`
- New struct `ParentProperties`
- New struct `StatusDetailProperties`
- New struct `SyncProperties`
- New struct `SyncUpdateProperties`
- New struct `TLSCertificateProperties`
- New struct `TLSProperties`
- New field `AzureADAuthenticationAsArmPolicy` in struct `Policies`
- New field `SystemData` in struct `PrivateLinkResource`
- New field `AutoGeneratedDomainNameLabelScope`, `ResourceGroupName` in struct `RegistryNameCheckRequest`
- New field `AvailableLoginServerName` in struct `RegistryNameStatus`
- New field `AnonymousPullEnabled`, `AutoGeneratedDomainNameLabelScope`, `NetworkRuleBypassAllowedForTasks`, `RoleAssignmentMode` in struct `RegistryProperties`
- New field `AnonymousPullEnabled`, `NetworkRuleBypassAllowedForTasks`, `RoleAssignmentMode` in struct `RegistryPropertiesUpdateParameters`


## 1.3.0-beta.3 (2025-04-24)
### Breaking Changes

- Operation `*RegistriesClient.BeginScheduleRun` has been changed to non-LRO, use `*RegistriesClient.ScheduleRun` instead.
- Operation `*RunsClient.BeginCancel` has been changed to non-LRO, use `*RunsClient.Cancel` instead.
- Operation `*RunsClient.BeginUpdate` has been changed to non-LRO, use `*RunsClient.Update` instead.
- Operation `*TaskRunsClient.BeginDelete` has been changed to non-LRO, use `*TaskRunsClient.Delete` instead.
- Operation `*TasksClient.BeginCreate` has been changed to non-LRO, use `*TasksClient.Create` instead.
- Operation `*TasksClient.BeginDelete` has been changed to non-LRO, use `*TasksClient.Delete` instead.
- Operation `*TasksClient.BeginUpdate` has been changed to non-LRO, use `*TasksClient.Update` instead.

### Features Added

- New enum type `RoleAssignmentMode` with values `RoleAssignmentModeAbacRepositoryPermissions`, `RoleAssignmentModeLegacyRegistryPermissions`
- New field `RoleAssignmentMode` in struct `RegistryProperties`
- New field `RoleAssignmentMode` in struct `RegistryPropertiesUpdateParameters`
- New field `Identity` in struct `SourceRegistryCredentials`


## 1.3.0-beta.2 (2024-12-27)
### Features Added

- New enum type `AutoGeneratedDomainNameLabelScope` with values `AutoGeneratedDomainNameLabelScopeNoReuse`, `AutoGeneratedDomainNameLabelScopeResourceGroupReuse`, `AutoGeneratedDomainNameLabelScopeSubscriptionReuse`, `AutoGeneratedDomainNameLabelScopeTenantReuse`, `AutoGeneratedDomainNameLabelScopeUnsecure`
- New struct `GarbageCollectionProperties`
- New field `GarbageCollection` in struct `ConnectedRegistryProperties`
- New field `GarbageCollection` in struct `ConnectedRegistryUpdateProperties`
- New field `AutoGeneratedDomainNameLabelScope`, `ResourceGroupName` in struct `RegistryNameCheckRequest`
- New field `AvailableLoginServerName` in struct `RegistryNameStatus`
- New field `AutoGeneratedDomainNameLabelScope` in struct `RegistryProperties`


## 1.3.0-beta.1 (2023-11-24)
### Features Added

- New enum type `ActivationStatus` with values `ActivationStatusActive`, `ActivationStatusInactive`
- New enum type `AuditLogStatus` with values `AuditLogStatusDisabled`, `AuditLogStatusEnabled`
- New enum type `AzureADAuthenticationAsArmPolicyStatus` with values `AzureADAuthenticationAsArmPolicyStatusDisabled`, `AzureADAuthenticationAsArmPolicyStatusEnabled`
- New enum type `CertificateType` with values `CertificateTypeLocalDirectory`
- New enum type `ConnectedRegistryMode` with values `ConnectedRegistryModeMirror`, `ConnectedRegistryModeReadOnly`, `ConnectedRegistryModeReadWrite`, `ConnectedRegistryModeRegistry`
- New enum type `ConnectionState` with values `ConnectionStateOffline`, `ConnectionStateOnline`, `ConnectionStateSyncing`, `ConnectionStateUnhealthy`
- New enum type `LogLevel` with values `LogLevelDebug`, `LogLevelError`, `LogLevelInformation`, `LogLevelNone`, `LogLevelWarning`
- New enum type `MetadataSearch` with values `MetadataSearchDisabled`, `MetadataSearchEnabled`
- New enum type `PackageSourceType` with values `PackageSourceTypeRemote`
- New enum type `PipelineOptions` with values `PipelineOptionsContinueOnErrors`, `PipelineOptionsDeleteSourceBlobOnSuccess`, `PipelineOptionsOverwriteBlobs`, `PipelineOptionsOverwriteTags`
- New enum type `PipelineRunSourceType` with values `PipelineRunSourceTypeAzureStorageBlob`
- New enum type `PipelineRunTargetType` with values `PipelineRunTargetTypeAzureStorageBlob`
- New enum type `PipelineSourceType` with values `PipelineSourceTypeAzureStorageBlobContainer`
- New enum type `TLSStatus` with values `TLSStatusDisabled`, `TLSStatusEnabled`
- New function `NewArchiveVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ArchiveVersionsClient, error)`
- New function `*ArchiveVersionsClient.BeginCreate(context.Context, string, string, string, string, string, *ArchiveVersionsClientBeginCreateOptions) (*runtime.Poller[ArchiveVersionsClientCreateResponse], error)`
- New function `*ArchiveVersionsClient.BeginDelete(context.Context, string, string, string, string, string, *ArchiveVersionsClientBeginDeleteOptions) (*runtime.Poller[ArchiveVersionsClientDeleteResponse], error)`
- New function `*ArchiveVersionsClient.Get(context.Context, string, string, string, string, string, *ArchiveVersionsClientGetOptions) (ArchiveVersionsClientGetResponse, error)`
- New function `*ArchiveVersionsClient.NewListPager(string, string, string, string, *ArchiveVersionsClientListOptions) *runtime.Pager[ArchiveVersionsClientListResponse]`
- New function `NewArchivesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ArchivesClient, error)`
- New function `*ArchivesClient.BeginCreate(context.Context, string, string, string, string, Archive, *ArchivesClientBeginCreateOptions) (*runtime.Poller[ArchivesClientCreateResponse], error)`
- New function `*ArchivesClient.BeginDelete(context.Context, string, string, string, string, *ArchivesClientBeginDeleteOptions) (*runtime.Poller[ArchivesClientDeleteResponse], error)`
- New function `*ArchivesClient.Get(context.Context, string, string, string, string, *ArchivesClientGetOptions) (ArchivesClientGetResponse, error)`
- New function `*ArchivesClient.NewListPager(string, string, string, *ArchivesClientListOptions) *runtime.Pager[ArchivesClientListResponse]`
- New function `*ArchivesClient.Update(context.Context, string, string, string, string, ArchiveUpdateParameters, *ArchivesClientUpdateOptions) (ArchivesClientUpdateResponse, error)`
- New function `*ClientFactory.NewArchiveVersionsClient() *ArchiveVersionsClient`
- New function `*ClientFactory.NewArchivesClient() *ArchivesClient`
- New function `*ClientFactory.NewConnectedRegistriesClient() *ConnectedRegistriesClient`
- New function `*ClientFactory.NewExportPipelinesClient() *ExportPipelinesClient`
- New function `*ClientFactory.NewImportPipelinesClient() *ImportPipelinesClient`
- New function `*ClientFactory.NewPipelineRunsClient() *PipelineRunsClient`
- New function `NewConnectedRegistriesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedRegistriesClient, error)`
- New function `*ConnectedRegistriesClient.BeginCreate(context.Context, string, string, string, ConnectedRegistry, *ConnectedRegistriesClientBeginCreateOptions) (*runtime.Poller[ConnectedRegistriesClientCreateResponse], error)`
- New function `*ConnectedRegistriesClient.BeginDeactivate(context.Context, string, string, string, *ConnectedRegistriesClientBeginDeactivateOptions) (*runtime.Poller[ConnectedRegistriesClientDeactivateResponse], error)`
- New function `*ConnectedRegistriesClient.BeginDelete(context.Context, string, string, string, *ConnectedRegistriesClientBeginDeleteOptions) (*runtime.Poller[ConnectedRegistriesClientDeleteResponse], error)`
- New function `*ConnectedRegistriesClient.Get(context.Context, string, string, string, *ConnectedRegistriesClientGetOptions) (ConnectedRegistriesClientGetResponse, error)`
- New function `*ConnectedRegistriesClient.NewListPager(string, string, *ConnectedRegistriesClientListOptions) *runtime.Pager[ConnectedRegistriesClientListResponse]`
- New function `*ConnectedRegistriesClient.BeginUpdate(context.Context, string, string, string, ConnectedRegistryUpdateParameters, *ConnectedRegistriesClientBeginUpdateOptions) (*runtime.Poller[ConnectedRegistriesClientUpdateResponse], error)`
- New function `NewExportPipelinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExportPipelinesClient, error)`
- New function `*ExportPipelinesClient.BeginCreate(context.Context, string, string, string, ExportPipeline, *ExportPipelinesClientBeginCreateOptions) (*runtime.Poller[ExportPipelinesClientCreateResponse], error)`
- New function `*ExportPipelinesClient.BeginDelete(context.Context, string, string, string, *ExportPipelinesClientBeginDeleteOptions) (*runtime.Poller[ExportPipelinesClientDeleteResponse], error)`
- New function `*ExportPipelinesClient.Get(context.Context, string, string, string, *ExportPipelinesClientGetOptions) (ExportPipelinesClientGetResponse, error)`
- New function `*ExportPipelinesClient.NewListPager(string, string, *ExportPipelinesClientListOptions) *runtime.Pager[ExportPipelinesClientListResponse]`
- New function `NewImportPipelinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ImportPipelinesClient, error)`
- New function `*ImportPipelinesClient.BeginCreate(context.Context, string, string, string, ImportPipeline, *ImportPipelinesClientBeginCreateOptions) (*runtime.Poller[ImportPipelinesClientCreateResponse], error)`
- New function `*ImportPipelinesClient.BeginDelete(context.Context, string, string, string, *ImportPipelinesClientBeginDeleteOptions) (*runtime.Poller[ImportPipelinesClientDeleteResponse], error)`
- New function `*ImportPipelinesClient.Get(context.Context, string, string, string, *ImportPipelinesClientGetOptions) (ImportPipelinesClientGetResponse, error)`
- New function `*ImportPipelinesClient.NewListPager(string, string, *ImportPipelinesClientListOptions) *runtime.Pager[ImportPipelinesClientListResponse]`
- New function `NewPipelineRunsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PipelineRunsClient, error)`
- New function `*PipelineRunsClient.BeginCreate(context.Context, string, string, string, PipelineRun, *PipelineRunsClientBeginCreateOptions) (*runtime.Poller[PipelineRunsClientCreateResponse], error)`
- New function `*PipelineRunsClient.BeginDelete(context.Context, string, string, string, *PipelineRunsClientBeginDeleteOptions) (*runtime.Poller[PipelineRunsClientDeleteResponse], error)`
- New function `*PipelineRunsClient.Get(context.Context, string, string, string, *PipelineRunsClientGetOptions) (PipelineRunsClientGetResponse, error)`
- New function `*PipelineRunsClient.NewListPager(string, string, *PipelineRunsClientListOptions) *runtime.Pager[PipelineRunsClientListResponse]`
- New struct `ActivationProperties`
- New struct `Archive`
- New struct `ArchiveListResult`
- New struct `ArchivePackageSourceProperties`
- New struct `ArchiveProperties`
- New struct `ArchiveUpdateParameters`
- New struct `ArchiveUpdateProperties`
- New struct `ArchiveVersion`
- New struct `ArchiveVersionListResult`
- New struct `ArchiveVersionProperties`
- New struct `AzureADAuthenticationAsArmPolicy`
- New struct `ConnectedRegistry`
- New struct `ConnectedRegistryListResult`
- New struct `ConnectedRegistryProperties`
- New struct `ConnectedRegistryUpdateParameters`
- New struct `ConnectedRegistryUpdateProperties`
- New struct `ExportPipeline`
- New struct `ExportPipelineListResult`
- New struct `ExportPipelineProperties`
- New struct `ExportPipelineTargetProperties`
- New struct `ImportPipeline`
- New struct `ImportPipelineListResult`
- New struct `ImportPipelineProperties`
- New struct `ImportPipelineSourceProperties`
- New struct `LoggingProperties`
- New struct `LoginServerProperties`
- New struct `ParentProperties`
- New struct `PipelineRun`
- New struct `PipelineRunListResult`
- New struct `PipelineRunProperties`
- New struct `PipelineRunRequest`
- New struct `PipelineRunResponse`
- New struct `PipelineRunSourceProperties`
- New struct `PipelineRunTargetProperties`
- New struct `PipelineSourceTriggerDescriptor`
- New struct `PipelineSourceTriggerProperties`
- New struct `PipelineTriggerDescriptor`
- New struct `PipelineTriggerProperties`
- New struct `ProgressProperties`
- New struct `SoftDeletePolicy`
- New struct `StatusDetailProperties`
- New struct `SyncProperties`
- New struct `SyncUpdateProperties`
- New struct `TLSCertificateProperties`
- New struct `TLSProperties`
- New field `AzureADAuthenticationAsArmPolicy`, `SoftDeletePolicy` in struct `Policies`
- New field `AnonymousPullEnabled`, `MetadataSearch` in struct `RegistryProperties`
- New field `AnonymousPullEnabled`, `MetadataSearch` in struct `RegistryPropertiesUpdateParameters`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0-beta.2 (2023-10-09)
### Other Changes

- Updated to latest `azcore` beta.

## 1.2.0-beta.1 (2023-09-22)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New enum type `ActivationStatus` with values `ActivationStatusActive`, `ActivationStatusInactive`
- New enum type `AuditLogStatus` with values `AuditLogStatusDisabled`, `AuditLogStatusEnabled`
- New enum type `AzureADAuthenticationAsArmPolicyStatus` with values `AzureADAuthenticationAsArmPolicyStatusDisabled`, `AzureADAuthenticationAsArmPolicyStatusEnabled`
- New enum type `CertificateType` with values `CertificateTypeLocalDirectory`
- New enum type `ConnectedRegistryMode` with values `ConnectedRegistryModeMirror`, `ConnectedRegistryModeReadOnly`, `ConnectedRegistryModeReadWrite`, `ConnectedRegistryModeRegistry`
- New enum type `ConnectionState` with values `ConnectionStateOffline`, `ConnectionStateOnline`, `ConnectionStateSyncing`, `ConnectionStateUnhealthy`
- New enum type `LogLevel` with values `LogLevelDebug`, `LogLevelError`, `LogLevelInformation`, `LogLevelNone`, `LogLevelWarning`
- New enum type `PackageSourceType` with values `PackageSourceTypeRemote`
- New enum type `PipelineOptions` with values `PipelineOptionsContinueOnErrors`, `PipelineOptionsDeleteSourceBlobOnSuccess`, `PipelineOptionsOverwriteBlobs`, `PipelineOptionsOverwriteTags`
- New enum type `PipelineRunSourceType` with values `PipelineRunSourceTypeAzureStorageBlob`
- New enum type `PipelineRunTargetType` with values `PipelineRunTargetTypeAzureStorageBlob`
- New enum type `PipelineSourceType` with values `PipelineSourceTypeAzureStorageBlobContainer`
- New enum type `TLSStatus` with values `TLSStatusDisabled`, `TLSStatusEnabled`
- New function `NewArchiveVersionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ArchiveVersionsClient, error)`
- New function `*ArchiveVersionsClient.BeginCreate(context.Context, string, string, string, string, string, *ArchiveVersionsClientBeginCreateOptions) (*runtime.Poller[ArchiveVersionsClientCreateResponse], error)`
- New function `*ArchiveVersionsClient.BeginDelete(context.Context, string, string, string, string, string, *ArchiveVersionsClientBeginDeleteOptions) (*runtime.Poller[ArchiveVersionsClientDeleteResponse], error)`
- New function `*ArchiveVersionsClient.Get(context.Context, string, string, string, string, string, *ArchiveVersionsClientGetOptions) (ArchiveVersionsClientGetResponse, error)`
- New function `*ArchiveVersionsClient.NewListPager(string, string, string, string, *ArchiveVersionsClientListOptions) *runtime.Pager[ArchiveVersionsClientListResponse]`
- New function `NewArchivesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ArchivesClient, error)`
- New function `*ArchivesClient.BeginCreate(context.Context, string, string, string, string, Archive, *ArchivesClientBeginCreateOptions) (*runtime.Poller[ArchivesClientCreateResponse], error)`
- New function `*ArchivesClient.BeginDelete(context.Context, string, string, string, string, *ArchivesClientBeginDeleteOptions) (*runtime.Poller[ArchivesClientDeleteResponse], error)`
- New function `*ArchivesClient.Get(context.Context, string, string, string, string, *ArchivesClientGetOptions) (ArchivesClientGetResponse, error)`
- New function `*ArchivesClient.NewListPager(string, string, string, *ArchivesClientListOptions) *runtime.Pager[ArchivesClientListResponse]`
- New function `*ArchivesClient.Update(context.Context, string, string, string, string, ArchiveUpdateParameters, *ArchivesClientUpdateOptions) (ArchivesClientUpdateResponse, error)`
- New function `*ClientFactory.NewArchiveVersionsClient() *ArchiveVersionsClient`
- New function `*ClientFactory.NewArchivesClient() *ArchivesClient`
- New function `*ClientFactory.NewConnectedRegistriesClient() *ConnectedRegistriesClient`
- New function `*ClientFactory.NewExportPipelinesClient() *ExportPipelinesClient`
- New function `*ClientFactory.NewImportPipelinesClient() *ImportPipelinesClient`
- New function `*ClientFactory.NewPipelineRunsClient() *PipelineRunsClient`
- New function `NewConnectedRegistriesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ConnectedRegistriesClient, error)`
- New function `*ConnectedRegistriesClient.BeginCreate(context.Context, string, string, string, ConnectedRegistry, *ConnectedRegistriesClientBeginCreateOptions) (*runtime.Poller[ConnectedRegistriesClientCreateResponse], error)`
- New function `*ConnectedRegistriesClient.BeginDeactivate(context.Context, string, string, string, *ConnectedRegistriesClientBeginDeactivateOptions) (*runtime.Poller[ConnectedRegistriesClientDeactivateResponse], error)`
- New function `*ConnectedRegistriesClient.BeginDelete(context.Context, string, string, string, *ConnectedRegistriesClientBeginDeleteOptions) (*runtime.Poller[ConnectedRegistriesClientDeleteResponse], error)`
- New function `*ConnectedRegistriesClient.Get(context.Context, string, string, string, *ConnectedRegistriesClientGetOptions) (ConnectedRegistriesClientGetResponse, error)`
- New function `*ConnectedRegistriesClient.NewListPager(string, string, *ConnectedRegistriesClientListOptions) *runtime.Pager[ConnectedRegistriesClientListResponse]`
- New function `*ConnectedRegistriesClient.BeginUpdate(context.Context, string, string, string, ConnectedRegistryUpdateParameters, *ConnectedRegistriesClientBeginUpdateOptions) (*runtime.Poller[ConnectedRegistriesClientUpdateResponse], error)`
- New function `NewExportPipelinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ExportPipelinesClient, error)`
- New function `*ExportPipelinesClient.BeginCreate(context.Context, string, string, string, ExportPipeline, *ExportPipelinesClientBeginCreateOptions) (*runtime.Poller[ExportPipelinesClientCreateResponse], error)`
- New function `*ExportPipelinesClient.BeginDelete(context.Context, string, string, string, *ExportPipelinesClientBeginDeleteOptions) (*runtime.Poller[ExportPipelinesClientDeleteResponse], error)`
- New function `*ExportPipelinesClient.Get(context.Context, string, string, string, *ExportPipelinesClientGetOptions) (ExportPipelinesClientGetResponse, error)`
- New function `*ExportPipelinesClient.NewListPager(string, string, *ExportPipelinesClientListOptions) *runtime.Pager[ExportPipelinesClientListResponse]`
- New function `NewImportPipelinesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ImportPipelinesClient, error)`
- New function `*ImportPipelinesClient.BeginCreate(context.Context, string, string, string, ImportPipeline, *ImportPipelinesClientBeginCreateOptions) (*runtime.Poller[ImportPipelinesClientCreateResponse], error)`
- New function `*ImportPipelinesClient.BeginDelete(context.Context, string, string, string, *ImportPipelinesClientBeginDeleteOptions) (*runtime.Poller[ImportPipelinesClientDeleteResponse], error)`
- New function `*ImportPipelinesClient.Get(context.Context, string, string, string, *ImportPipelinesClientGetOptions) (ImportPipelinesClientGetResponse, error)`
- New function `*ImportPipelinesClient.NewListPager(string, string, *ImportPipelinesClientListOptions) *runtime.Pager[ImportPipelinesClientListResponse]`
- New function `NewPipelineRunsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PipelineRunsClient, error)`
- New function `*PipelineRunsClient.BeginCreate(context.Context, string, string, string, PipelineRun, *PipelineRunsClientBeginCreateOptions) (*runtime.Poller[PipelineRunsClientCreateResponse], error)`
- New function `*PipelineRunsClient.BeginDelete(context.Context, string, string, string, *PipelineRunsClientBeginDeleteOptions) (*runtime.Poller[PipelineRunsClientDeleteResponse], error)`
- New function `*PipelineRunsClient.Get(context.Context, string, string, string, *PipelineRunsClientGetOptions) (PipelineRunsClientGetResponse, error)`
- New function `*PipelineRunsClient.NewListPager(string, string, *PipelineRunsClientListOptions) *runtime.Pager[PipelineRunsClientListResponse]`
- New struct `ActivationProperties`
- New struct `Archive`
- New struct `ArchiveListResult`
- New struct `ArchivePackageSourceProperties`
- New struct `ArchiveProperties`
- New struct `ArchiveUpdateParameters`
- New struct `ArchiveUpdateProperties`
- New struct `ArchiveVersion`
- New struct `ArchiveVersionListResult`
- New struct `ArchiveVersionProperties`
- New struct `AzureADAuthenticationAsArmPolicy`
- New struct `ConnectedRegistry`
- New struct `ConnectedRegistryListResult`
- New struct `ConnectedRegistryProperties`
- New struct `ConnectedRegistryUpdateParameters`
- New struct `ConnectedRegistryUpdateProperties`
- New struct `ExportPipeline`
- New struct `ExportPipelineListResult`
- New struct `ExportPipelineProperties`
- New struct `ExportPipelineTargetProperties`
- New struct `ImportPipeline`
- New struct `ImportPipelineListResult`
- New struct `ImportPipelineProperties`
- New struct `ImportPipelineSourceProperties`
- New struct `LoggingProperties`
- New struct `LoginServerProperties`
- New struct `ParentProperties`
- New struct `PipelineRun`
- New struct `PipelineRunListResult`
- New struct `PipelineRunProperties`
- New struct `PipelineRunRequest`
- New struct `PipelineRunResponse`
- New struct `PipelineRunSourceProperties`
- New struct `PipelineRunTargetProperties`
- New struct `PipelineSourceTriggerDescriptor`
- New struct `PipelineSourceTriggerProperties`
- New struct `PipelineTriggerDescriptor`
- New struct `PipelineTriggerProperties`
- New struct `ProgressProperties`
- New struct `SoftDeletePolicy`
- New struct `StatusDetailProperties`
- New struct `SyncProperties`
- New struct `SyncUpdateProperties`
- New struct `TLSCertificateProperties`
- New struct `TLSProperties`
- New field `AzureADAuthenticationAsArmPolicy`, `SoftDeletePolicy` in struct `Policies`
- New field `AnonymousPullEnabled` in struct `RegistryProperties`
- New field `AnonymousPullEnabled` in struct `RegistryPropertiesUpdateParameters`


## 1.1.0 (2023-09-22)
### Features Added

- New enum type `CredentialHealthStatus` with values `CredentialHealthStatusHealthy`, `CredentialHealthStatusUnhealthy`
- New enum type `CredentialName` with values `CredentialNameCredential1`
- New function `NewCacheRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CacheRulesClient, error)`
- New function `*CacheRulesClient.BeginCreate(context.Context, string, string, string, CacheRule, *CacheRulesClientBeginCreateOptions) (*runtime.Poller[CacheRulesClientCreateResponse], error)`
- New function `*CacheRulesClient.BeginDelete(context.Context, string, string, string, *CacheRulesClientBeginDeleteOptions) (*runtime.Poller[CacheRulesClientDeleteResponse], error)`
- New function `*CacheRulesClient.Get(context.Context, string, string, string, *CacheRulesClientGetOptions) (CacheRulesClientGetResponse, error)`
- New function `*CacheRulesClient.NewListPager(string, string, *CacheRulesClientListOptions) *runtime.Pager[CacheRulesClientListResponse]`
- New function `*CacheRulesClient.BeginUpdate(context.Context, string, string, string, CacheRuleUpdateParameters, *CacheRulesClientBeginUpdateOptions) (*runtime.Poller[CacheRulesClientUpdateResponse], error)`
- New function `*ClientFactory.NewCacheRulesClient() *CacheRulesClient`
- New function `*ClientFactory.NewCredentialSetsClient() *CredentialSetsClient`
- New function `NewCredentialSetsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CredentialSetsClient, error)`
- New function `*CredentialSetsClient.BeginCreate(context.Context, string, string, string, CredentialSet, *CredentialSetsClientBeginCreateOptions) (*runtime.Poller[CredentialSetsClientCreateResponse], error)`
- New function `*CredentialSetsClient.BeginDelete(context.Context, string, string, string, *CredentialSetsClientBeginDeleteOptions) (*runtime.Poller[CredentialSetsClientDeleteResponse], error)`
- New function `*CredentialSetsClient.Get(context.Context, string, string, string, *CredentialSetsClientGetOptions) (CredentialSetsClientGetResponse, error)`
- New function `*CredentialSetsClient.NewListPager(string, string, *CredentialSetsClientListOptions) *runtime.Pager[CredentialSetsClientListResponse]`
- New function `*CredentialSetsClient.BeginUpdate(context.Context, string, string, string, CredentialSetUpdateParameters, *CredentialSetsClientBeginUpdateOptions) (*runtime.Poller[CredentialSetsClientUpdateResponse], error)`
- New struct `AuthCredential`
- New struct `CacheRule`
- New struct `CacheRuleProperties`
- New struct `CacheRuleUpdateParameters`
- New struct `CacheRuleUpdateProperties`
- New struct `CacheRulesListResult`
- New struct `CredentialHealth`
- New struct `CredentialSet`
- New struct `CredentialSetListResult`
- New struct `CredentialSetProperties`
- New struct `CredentialSetUpdateParameters`
- New struct `CredentialSetUpdateProperties`


## 1.0.0 (2023-03-24)
### Breaking Changes

- Type alias `ActivationStatus` has been removed
- Type alias `AuditLogStatus` has been removed
- Type alias `AzureADAuthenticationAsArmPolicyStatus` has been removed
- Type alias `CertificateType` has been removed
- Type alias `ConnectedRegistryMode` has been removed
- Type alias `ConnectionState` has been removed
- Type alias `CredentialHealthStatus` has been removed
- Type alias `CredentialName` has been removed
- Type alias `LogLevel` has been removed
- Type alias `PipelineOptions` has been removed
- Type alias `PipelineRunSourceType` has been removed
- Type alias `PipelineRunTargetType` has been removed
- Type alias `PipelineSourceType` has been removed
- Type alias `TLSStatus` has been removed
- Function `NewCacheRulesClient` has been removed
- Function `*CacheRulesClient.BeginCreate` has been removed
- Function `*CacheRulesClient.BeginDelete` has been removed
- Function `*CacheRulesClient.Get` has been removed
- Function `*CacheRulesClient.NewListPager` has been removed
- Function `*CacheRulesClient.BeginUpdate` has been removed
- Function `NewConnectedRegistriesClient` has been removed
- Function `*ConnectedRegistriesClient.BeginCreate` has been removed
- Function `*ConnectedRegistriesClient.BeginDeactivate` has been removed
- Function `*ConnectedRegistriesClient.BeginDelete` has been removed
- Function `*ConnectedRegistriesClient.Get` has been removed
- Function `*ConnectedRegistriesClient.NewListPager` has been removed
- Function `*ConnectedRegistriesClient.BeginUpdate` has been removed
- Function `NewCredentialSetsClient` has been removed
- Function `*CredentialSetsClient.BeginCreate` has been removed
- Function `*CredentialSetsClient.BeginDelete` has been removed
- Function `*CredentialSetsClient.Get` has been removed
- Function `*CredentialSetsClient.NewListPager` has been removed
- Function `*CredentialSetsClient.BeginUpdate` has been removed
- Function `NewExportPipelinesClient` has been removed
- Function `*ExportPipelinesClient.BeginCreate` has been removed
- Function `*ExportPipelinesClient.BeginDelete` has been removed
- Function `*ExportPipelinesClient.Get` has been removed
- Function `*ExportPipelinesClient.NewListPager` has been removed
- Function `NewImportPipelinesClient` has been removed
- Function `*ImportPipelinesClient.BeginCreate` has been removed
- Function `*ImportPipelinesClient.BeginDelete` has been removed
- Function `*ImportPipelinesClient.Get` has been removed
- Function `*ImportPipelinesClient.NewListPager` has been removed
- Function `NewPipelineRunsClient` has been removed
- Function `*PipelineRunsClient.BeginCreate` has been removed
- Function `*PipelineRunsClient.BeginDelete` has been removed
- Function `*PipelineRunsClient.Get` has been removed
- Function `*PipelineRunsClient.NewListPager` has been removed
- Struct `ActivationProperties` has been removed
- Struct `AuthCredential` has been removed
- Struct `AzureADAuthenticationAsArmPolicy` has been removed
- Struct `CacheRule` has been removed
- Struct `CacheRuleProperties` has been removed
- Struct `CacheRuleUpdateParameters` has been removed
- Struct `CacheRuleUpdateProperties` has been removed
- Struct `CacheRulesClient` has been removed
- Struct `CacheRulesListResult` has been removed
- Struct `ConnectedRegistriesClient` has been removed
- Struct `ConnectedRegistry` has been removed
- Struct `ConnectedRegistryListResult` has been removed
- Struct `ConnectedRegistryProperties` has been removed
- Struct `ConnectedRegistryUpdateParameters` has been removed
- Struct `ConnectedRegistryUpdateProperties` has been removed
- Struct `CredentialHealth` has been removed
- Struct `CredentialSet` has been removed
- Struct `CredentialSetListResult` has been removed
- Struct `CredentialSetProperties` has been removed
- Struct `CredentialSetUpdateParameters` has been removed
- Struct `CredentialSetUpdateProperties` has been removed
- Struct `CredentialSetsClient` has been removed
- Struct `ExportPipeline` has been removed
- Struct `ExportPipelineListResult` has been removed
- Struct `ExportPipelineProperties` has been removed
- Struct `ExportPipelineTargetProperties` has been removed
- Struct `ExportPipelinesClient` has been removed
- Struct `ImportPipeline` has been removed
- Struct `ImportPipelineListResult` has been removed
- Struct `ImportPipelineProperties` has been removed
- Struct `ImportPipelineSourceProperties` has been removed
- Struct `ImportPipelinesClient` has been removed
- Struct `LoggingProperties` has been removed
- Struct `LoginServerProperties` has been removed
- Struct `ParentProperties` has been removed
- Struct `PipelineRun` has been removed
- Struct `PipelineRunListResult` has been removed
- Struct `PipelineRunProperties` has been removed
- Struct `PipelineRunRequest` has been removed
- Struct `PipelineRunResponse` has been removed
- Struct `PipelineRunSourceProperties` has been removed
- Struct `PipelineRunTargetProperties` has been removed
- Struct `PipelineRunsClient` has been removed
- Struct `PipelineSourceTriggerDescriptor` has been removed
- Struct `PipelineSourceTriggerProperties` has been removed
- Struct `PipelineTriggerDescriptor` has been removed
- Struct `PipelineTriggerProperties` has been removed
- Struct `ProgressProperties` has been removed
- Struct `SoftDeletePolicy` has been removed
- Struct `StatusDetailProperties` has been removed
- Struct `SyncProperties` has been removed
- Struct `SyncUpdateProperties` has been removed
- Struct `TLSCertificateProperties` has been removed
- Struct `TLSProperties` has been removed
- Field `AzureADAuthenticationAsArmPolicy` of struct `Policies` has been removed
- Field `SoftDeletePolicy` of struct `Policies` has been removed
- Field `AnonymousPullEnabled` of struct `RegistryProperties` has been removed
- Field `AnonymousPullEnabled` of struct `RegistryPropertiesUpdateParameters` has been removed

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.7.0 (2023-02-24)
### Features Added

- New type alias `CredentialHealthStatus` with values `CredentialHealthStatusHealthy`, `CredentialHealthStatusUnhealthy`
- New type alias `CredentialName` with values `CredentialNameCredential1`
- New function `NewCacheRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CacheRulesClient, error)`
- New function `*CacheRulesClient.BeginCreate(context.Context, string, string, string, CacheRule, *CacheRulesClientBeginCreateOptions) (*runtime.Poller[CacheRulesClientCreateResponse], error)`
- New function `*CacheRulesClient.BeginDelete(context.Context, string, string, string, *CacheRulesClientBeginDeleteOptions) (*runtime.Poller[CacheRulesClientDeleteResponse], error)`
- New function `*CacheRulesClient.Get(context.Context, string, string, string, *CacheRulesClientGetOptions) (CacheRulesClientGetResponse, error)`
- New function `*CacheRulesClient.NewListPager(string, string, *CacheRulesClientListOptions) *runtime.Pager[CacheRulesClientListResponse]`
- New function `*CacheRulesClient.BeginUpdate(context.Context, string, string, string, CacheRuleUpdateParameters, *CacheRulesClientBeginUpdateOptions) (*runtime.Poller[CacheRulesClientUpdateResponse], error)`
- New function `NewCredentialSetsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CredentialSetsClient, error)`
- New function `*CredentialSetsClient.BeginCreate(context.Context, string, string, string, CredentialSet, *CredentialSetsClientBeginCreateOptions) (*runtime.Poller[CredentialSetsClientCreateResponse], error)`
- New function `*CredentialSetsClient.BeginDelete(context.Context, string, string, string, *CredentialSetsClientBeginDeleteOptions) (*runtime.Poller[CredentialSetsClientDeleteResponse], error)`
- New function `*CredentialSetsClient.Get(context.Context, string, string, string, *CredentialSetsClientGetOptions) (CredentialSetsClientGetResponse, error)`
- New function `*CredentialSetsClient.NewListPager(string, string, *CredentialSetsClientListOptions) *runtime.Pager[CredentialSetsClientListResponse]`
- New function `*CredentialSetsClient.BeginUpdate(context.Context, string, string, string, CredentialSetUpdateParameters, *CredentialSetsClientBeginUpdateOptions) (*runtime.Poller[CredentialSetsClientUpdateResponse], error)`
- New struct `AuthCredential`
- New struct `CacheRule`
- New struct `CacheRuleProperties`
- New struct `CacheRuleUpdateParameters`
- New struct `CacheRuleUpdateProperties`
- New struct `CacheRulesClient`
- New struct `CacheRulesClientCreateResponse`
- New struct `CacheRulesClientDeleteResponse`
- New struct `CacheRulesClientListResponse`
- New struct `CacheRulesClientUpdateResponse`
- New struct `CacheRulesListResult`
- New struct `CredentialHealth`
- New struct `CredentialSet`
- New struct `CredentialSetListResult`
- New struct `CredentialSetProperties`
- New struct `CredentialSetUpdateParameters`
- New struct `CredentialSetUpdateProperties`
- New struct `CredentialSetsClient`
- New struct `CredentialSetsClientCreateResponse`
- New struct `CredentialSetsClientDeleteResponse`
- New struct `CredentialSetsClientListResponse`
- New struct `CredentialSetsClientUpdateResponse`


## 0.6.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
