# Release History

## 1.2.0-beta.1 (2026-10-15)
### Features Added

- New enum type `DetectorRenderingType` with values `DetectorRenderingTypeDataSummary`, `DetectorRenderingTypeDynamicInsight`, `DetectorRenderingTypeEmailSnippet`, `DetectorRenderingTypeInsight`, `DetectorRenderingTypeMarkdown`, `DetectorRenderingTypePieChart`, `DetectorRenderingTypeTable`, `DetectorRenderingTypeTimeSeriesGraph`, `DetectorRenderingTypeTimeSeriesPerInstanceGraph`
- New enum type `DetectorStatusID` with values `DetectorStatusIDCritical`, `DetectorStatusIDInfo`, `DetectorStatusIDNone`, `DetectorStatusIDSuccess`, `DetectorStatusIDWarning`
- New enum type `DetectorType` with values `DetectorTypeAnalysis`, `DetectorTypeCategoryOverview`, `DetectorTypeDetector`
- New enum type `FailoverState` with values `FailoverStateFailed`, `FailoverStateInProgress`, `FailoverStateNone`, `FailoverStateSucceeded`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `ReplicationRole` with values `ReplicationRolePrimary`, `ReplicationRoleSecondary`
- New enum type `TransparentDataEncryptionKeySource` with values `TransparentDataEncryptionKeySourceCustomerManaged`, `TransparentDataEncryptionKeySourceMicrosoftManaged`
- New function `*ClientFactory.NewDetectorsClient() *DetectorsClient`
- New function `*ClientFactory.NewFailoverInstancesClient() *FailoverInstancesClient`
- New function `*ClientFactory.NewTaskQueuesClient() *TaskQueuesClient`
- New function `*ClientFactory.NewTransparentDataEncryptionsClient() *TransparentDataEncryptionsClient`
- New function `NewDetectorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DetectorsClient, error)`
- New function `*DetectorsClient.Get(ctx context.Context, resourceGroupName string, schedulerName string, detectorID string, options *DetectorsClientGetOptions) (DetectorsClientGetResponse, error)`
- New function `*DetectorsClient.NewListBySchedulerPager(resourceGroupName string, schedulerName string, options *DetectorsClientListBySchedulerOptions) *runtime.Pager[DetectorsClientListBySchedulerResponse]`
- New function `NewFailoverInstancesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FailoverInstancesClient, error)`
- New function `*FailoverInstancesClient.BeginCreateOrReplace(ctx context.Context, resourceGroupName string, schedulerName string, resource FailoverInstance, options *FailoverInstancesClientBeginCreateOrReplaceOptions) (*runtime.Poller[FailoverInstancesClientCreateOrReplaceResponse], error)`
- New function `*FailoverInstancesClient.BeginDelete(ctx context.Context, resourceGroupName string, schedulerName string, options *FailoverInstancesClientBeginDeleteOptions) (*runtime.Poller[FailoverInstancesClientDeleteResponse], error)`
- New function `*FailoverInstancesClient.BeginFailover(ctx context.Context, resourceGroupName string, schedulerName string, options *FailoverInstancesClientBeginFailoverOptions) (*runtime.Poller[FailoverInstancesClientFailoverResponse], error)`
- New function `*FailoverInstancesClient.BeginForceFailoverAllowDataLoss(ctx context.Context, resourceGroupName string, schedulerName string, options *FailoverInstancesClientBeginForceFailoverAllowDataLossOptions) (*runtime.Poller[FailoverInstancesClientForceFailoverAllowDataLossResponse], error)`
- New function `*FailoverInstancesClient.Get(ctx context.Context, resourceGroupName string, schedulerName string, options *FailoverInstancesClientGetOptions) (FailoverInstancesClientGetResponse, error)`
- New function `*FailoverInstancesClient.NewListBySchedulerPager(resourceGroupName string, schedulerName string, options *FailoverInstancesClientListBySchedulerOptions) *runtime.Pager[FailoverInstancesClientListBySchedulerResponse]`
- New function `*SchedulersClient.Restart(ctx context.Context, resourceGroupName string, schedulerName string, options *SchedulersClientRestartOptions) (SchedulersClientRestartResponse, error)`
- New function `NewTaskQueuesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TaskQueuesClient, error)`
- New function `*TaskQueuesClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, schedulerName string, taskHubName string, taskQueueName string, resource TaskQueue, options *TaskQueuesClientBeginCreateOrUpdateOptions) (*runtime.Poller[TaskQueuesClientCreateOrUpdateResponse], error)`
- New function `*TaskQueuesClient.BeginDelete(ctx context.Context, resourceGroupName string, schedulerName string, taskHubName string, taskQueueName string, options *TaskQueuesClientBeginDeleteOptions) (*runtime.Poller[TaskQueuesClientDeleteResponse], error)`
- New function `*TaskQueuesClient.Get(ctx context.Context, resourceGroupName string, schedulerName string, taskHubName string, taskQueueName string, options *TaskQueuesClientGetOptions) (TaskQueuesClientGetResponse, error)`
- New function `*TaskQueuesClient.NewListByTaskHubPager(resourceGroupName string, schedulerName string, taskHubName string, options *TaskQueuesClientListByTaskHubOptions) *runtime.Pager[TaskQueuesClientListByTaskHubResponse]`
- New function `*TaskQueuesClient.BeginUpdate(ctx context.Context, resourceGroupName string, schedulerName string, taskHubName string, taskQueueName string, properties TaskQueueUpdate, options *TaskQueuesClientBeginUpdateOptions) (*runtime.Poller[TaskQueuesClientUpdateResponse], error)`
- New function `NewTransparentDataEncryptionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TransparentDataEncryptionsClient, error)`
- New function `*TransparentDataEncryptionsClient.BeginCreateOrReplace(ctx context.Context, resourceGroupName string, schedulerName string, resource TransparentDataEncryption, options *TransparentDataEncryptionsClientBeginCreateOrReplaceOptions) (*runtime.Poller[TransparentDataEncryptionsClientCreateOrReplaceResponse], error)`
- New function `*TransparentDataEncryptionsClient.BeginDelete(ctx context.Context, resourceGroupName string, schedulerName string, options *TransparentDataEncryptionsClientBeginDeleteOptions) (*runtime.Poller[TransparentDataEncryptionsClientDeleteResponse], error)`
- New function `*TransparentDataEncryptionsClient.Get(ctx context.Context, resourceGroupName string, schedulerName string, options *TransparentDataEncryptionsClientGetOptions) (TransparentDataEncryptionsClientGetResponse, error)`
- New function `*TransparentDataEncryptionsClient.NewListBySchedulerPager(resourceGroupName string, schedulerName string, options *TransparentDataEncryptionsClientListBySchedulerOptions) *runtime.Pager[TransparentDataEncryptionsClientListBySchedulerResponse]`
- New struct `Detector`
- New struct `DetectorDataset`
- New struct `DetectorListResult`
- New struct `DetectorMetadata`
- New struct `DetectorProperties`
- New struct `DetectorRenderingProperties`
- New struct `DetectorStatus`
- New struct `DetectorTable`
- New struct `DetectorTableColumn`
- New struct `FailoverInstance`
- New struct `FailoverInstanceListResult`
- New struct `FailoverInstanceProperties`
- New struct `ManagedServiceIdentity`
- New struct `TaskQueue`
- New struct `TaskQueueListResult`
- New struct `TaskQueueProperties`
- New struct `TaskQueueThroughput`
- New struct `TaskQueueUpdate`
- New struct `TransparentDataEncryption`
- New struct `TransparentDataEncryptionListResult`
- New struct `TransparentDataEncryptionProperties`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Scheduler`
- New field `Identity` in struct `SchedulerUpdate`
- New field `Capabilities` in struct `TaskHubProperties`


## 1.1.0 (2026-03-09)
### Features Added

- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New function `*SchedulersClient.BeginCreateOrUpdatePrivateEndpointConnection(ctx context.Context, resourceGroupName string, schedulerName string, privateEndpointConnectionName string, resource PrivateEndpointConnection, options *SchedulersClientBeginCreateOrUpdatePrivateEndpointConnectionOptions) (*runtime.Poller[SchedulersClientCreateOrUpdatePrivateEndpointConnectionResponse], error)`
- New function `*SchedulersClient.BeginDeletePrivateEndpointConnection(ctx context.Context, resourceGroupName string, schedulerName string, privateEndpointConnectionName string, options *SchedulersClientBeginDeletePrivateEndpointConnectionOptions) (*runtime.Poller[SchedulersClientDeletePrivateEndpointConnectionResponse], error)`
- New function `*SchedulersClient.GetPrivateEndpointConnection(ctx context.Context, resourceGroupName string, schedulerName string, privateEndpointConnectionName string, options *SchedulersClientGetPrivateEndpointConnectionOptions) (SchedulersClientGetPrivateEndpointConnectionResponse, error)`
- New function `*SchedulersClient.GetPrivateLink(ctx context.Context, resourceGroupName string, schedulerName string, privateLinkResourceName string, options *SchedulersClientGetPrivateLinkOptions) (SchedulersClientGetPrivateLinkResponse, error)`
- New function `*SchedulersClient.NewListPrivateEndpointConnectionsPager(resourceGroupName string, schedulerName string, options *SchedulersClientListPrivateEndpointConnectionsOptions) *runtime.Pager[SchedulersClientListPrivateEndpointConnectionsResponse]`
- New function `*SchedulersClient.NewListPrivateLinksPager(resourceGroupName string, schedulerName string, options *SchedulersClientListPrivateLinksOptions) *runtime.Pager[SchedulersClientListPrivateLinksResponse]`
- New function `*SchedulersClient.BeginUpdatePrivateEndpointConnection(ctx context.Context, resourceGroupName string, schedulerName string, privateEndpointConnectionName string, properties PrivateEndpointConnectionUpdate, options *SchedulersClientBeginUpdatePrivateEndpointConnectionOptions) (*runtime.Poller[SchedulersClientUpdatePrivateEndpointConnectionResponse], error)`
- New struct `OptionalPropertiesUpdateableProperties`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateEndpointConnectionUpdate`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `SchedulerPrivateLinkResource`
- New struct `SchedulerPrivateLinkResourceListResult`
- New field `PrivateEndpointConnections`, `PublicNetworkAccess` in struct `SchedulerProperties`
- New field `PublicNetworkAccess` in struct `SchedulerPropertiesUpdate`


## 1.0.0 (2025-09-26)
### Breaking Changes

- Type of `SchedulerSKU.Name` has been changed from `*string` to `*SchedulerSKUName`
- Type of `SchedulerSKUUpdate.Name` has been changed from `*string` to `*SchedulerSKUName`

### Features Added

- New enum type `SchedulerSKUName` with values `SchedulerSKUNameConsumption`, `SchedulerSKUNameDedicated`


## 0.2.0 (2025-04-15)
### Features Added

- New enum type `PurgeableOrchestrationState` with values `PurgeableOrchestrationStateCanceled`, `PurgeableOrchestrationStateCompleted`, `PurgeableOrchestrationStateFailed`, `PurgeableOrchestrationStateTerminated`
- New function `*ClientFactory.NewRetentionPoliciesClient() *RetentionPoliciesClient`
- New function `NewRetentionPoliciesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*RetentionPoliciesClient, error)`
- New function `*RetentionPoliciesClient.BeginCreateOrReplace(context.Context, string, string, RetentionPolicy, *RetentionPoliciesClientBeginCreateOrReplaceOptions) (*runtime.Poller[RetentionPoliciesClientCreateOrReplaceResponse], error)`
- New function `*RetentionPoliciesClient.BeginDelete(context.Context, string, string, *RetentionPoliciesClientBeginDeleteOptions) (*runtime.Poller[RetentionPoliciesClientDeleteResponse], error)`
- New function `*RetentionPoliciesClient.Get(context.Context, string, string, *RetentionPoliciesClientGetOptions) (RetentionPoliciesClientGetResponse, error)`
- New function `*RetentionPoliciesClient.NewListBySchedulerPager(string, string, *RetentionPoliciesClientListBySchedulerOptions) *runtime.Pager[RetentionPoliciesClientListBySchedulerResponse]`
- New function `*RetentionPoliciesClient.BeginUpdate(context.Context, string, string, RetentionPolicy, *RetentionPoliciesClientBeginUpdateOptions) (*runtime.Poller[RetentionPoliciesClientUpdateResponse], error)`
- New struct `RetentionPolicy`
- New struct `RetentionPolicyDetails`
- New struct `RetentionPolicyListResult`
- New struct `RetentionPolicyProperties`


## 0.1.0 (2025-03-20)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/durabletask/armdurabletask` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).