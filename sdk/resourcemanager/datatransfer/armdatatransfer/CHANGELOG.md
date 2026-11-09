# Release History

## 0.2.0 (2026-11-09)
### Features Added

- New value `TargetTypeMessagingTransferFlow`, `TargetTypeStorageTransferFlow` added to enum type `TargetType`
- New enum type `CommonTransferFlowStatus` with values `CommonTransferFlowStatusApprovedWaitingToBeLinked`, `CommonTransferFlowStatusAwaitingApproval`, `CommonTransferFlowStatusDisabled`, `CommonTransferFlowStatusEnabled`, `CommonTransferFlowStatusRejected`
- New enum type `TransferRole` with values `TransferRoleDestination`, `TransferRoleSource`
- New function `*ClientFactory.NewMessagingTransferFlowsClient() *MessagingTransferFlowsClient`
- New function `*ClientFactory.NewStorageTransferFlowsClient() *StorageTransferFlowsClient`
- New function `NewMessagingTransferFlowsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MessagingTransferFlowsClient, error)`
- New function `*MessagingTransferFlowsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, messagingTransferFlowName string, messagingTransferFlow MessagingTransferFlow, options *MessagingTransferFlowsClientBeginCreateOrUpdateOptions) (*runtime.Poller[MessagingTransferFlowsClientCreateOrUpdateResponse], error)`
- New function `*MessagingTransferFlowsClient.BeginDelete(ctx context.Context, resourceGroupName string, messagingTransferFlowName string, options *MessagingTransferFlowsClientBeginDeleteOptions) (*runtime.Poller[MessagingTransferFlowsClientDeleteResponse], error)`
- New function `*MessagingTransferFlowsClient.BeginDisableTransfer(ctx context.Context, resourceGroupName string, messagingTransferFlowName string, options *MessagingTransferFlowsClientBeginDisableTransferOptions) (*runtime.Poller[MessagingTransferFlowsClientDisableTransferResponse], error)`
- New function `*MessagingTransferFlowsClient.BeginEnableTransfer(ctx context.Context, resourceGroupName string, messagingTransferFlowName string, options *MessagingTransferFlowsClientBeginEnableTransferOptions) (*runtime.Poller[MessagingTransferFlowsClientEnableTransferResponse], error)`
- New function `*MessagingTransferFlowsClient.Get(ctx context.Context, resourceGroupName string, messagingTransferFlowName string, options *MessagingTransferFlowsClientGetOptions) (MessagingTransferFlowsClientGetResponse, error)`
- New function `*MessagingTransferFlowsClient.BeginLinkRemoteFlow(ctx context.Context, resourceGroupName string, messagingTransferFlowName string, options *MessagingTransferFlowsClientBeginLinkRemoteFlowOptions) (*runtime.Poller[MessagingTransferFlowsClientLinkRemoteFlowResponse], error)`
- New function `*MessagingTransferFlowsClient.NewListByResourceGroupPager(resourceGroupName string, options *MessagingTransferFlowsClientListByResourceGroupOptions) *runtime.Pager[MessagingTransferFlowsClientListByResourceGroupResponse]`
- New function `*MessagingTransferFlowsClient.NewListBySubscriptionPager(options *MessagingTransferFlowsClientListBySubscriptionOptions) *runtime.Pager[MessagingTransferFlowsClientListBySubscriptionResponse]`
- New function `*MessagingTransferFlowsClient.BeginUpdate(ctx context.Context, resourceGroupName string, messagingTransferFlowName string, messagingTransferFlow MessagingTransferFlowsPatch, options *MessagingTransferFlowsClientBeginUpdateOptions) (*runtime.Poller[MessagingTransferFlowsClientUpdateResponse], error)`
- New function `*PipelinesClient.BeginApproveTransferFlow(ctx context.Context, resourceGroupName string, pipelineName string, body TransferFlowRequestBody, options *PipelinesClientBeginApproveTransferFlowOptions) (*runtime.Poller[PipelinesClientApproveTransferFlowResponse], error)`
- New function `*PipelinesClient.NewListTransferFlowsPager(resourceGroupName string, pipelineName string, options *PipelinesClientListTransferFlowsOptions) *runtime.Pager[PipelinesClientListTransferFlowsResponse]`
- New function `*PipelinesClient.BeginRejectTransferFlow(ctx context.Context, resourceGroupName string, pipelineName string, body TransferFlowRequestBody, options *PipelinesClientBeginRejectTransferFlowOptions) (*runtime.Poller[PipelinesClientRejectTransferFlowResponse], error)`
- New function `NewStorageTransferFlowsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*StorageTransferFlowsClient, error)`
- New function `*StorageTransferFlowsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, storageTransferFlowName string, storageTransferFlow StorageTransferFlow, options *StorageTransferFlowsClientBeginCreateOrUpdateOptions) (*runtime.Poller[StorageTransferFlowsClientCreateOrUpdateResponse], error)`
- New function `*StorageTransferFlowsClient.BeginDelete(ctx context.Context, resourceGroupName string, storageTransferFlowName string, options *StorageTransferFlowsClientBeginDeleteOptions) (*runtime.Poller[StorageTransferFlowsClientDeleteResponse], error)`
- New function `*StorageTransferFlowsClient.BeginDisableTransfer(ctx context.Context, resourceGroupName string, storageTransferFlowName string, options *StorageTransferFlowsClientBeginDisableTransferOptions) (*runtime.Poller[StorageTransferFlowsClientDisableTransferResponse], error)`
- New function `*StorageTransferFlowsClient.BeginEnableTransfer(ctx context.Context, resourceGroupName string, storageTransferFlowName string, options *StorageTransferFlowsClientBeginEnableTransferOptions) (*runtime.Poller[StorageTransferFlowsClientEnableTransferResponse], error)`
- New function `*StorageTransferFlowsClient.Get(ctx context.Context, resourceGroupName string, storageTransferFlowName string, options *StorageTransferFlowsClientGetOptions) (StorageTransferFlowsClientGetResponse, error)`
- New function `*StorageTransferFlowsClient.BeginLinkRemoteFlow(ctx context.Context, resourceGroupName string, storageTransferFlowName string, options *StorageTransferFlowsClientBeginLinkRemoteFlowOptions) (*runtime.Poller[StorageTransferFlowsClientLinkRemoteFlowResponse], error)`
- New function `*StorageTransferFlowsClient.NewListByResourceGroupPager(resourceGroupName string, options *StorageTransferFlowsClientListByResourceGroupOptions) *runtime.Pager[StorageTransferFlowsClientListByResourceGroupResponse]`
- New function `*StorageTransferFlowsClient.NewListBySubscriptionPager(options *StorageTransferFlowsClientListBySubscriptionOptions) *runtime.Pager[StorageTransferFlowsClientListBySubscriptionResponse]`
- New function `*StorageTransferFlowsClient.BeginUpdate(ctx context.Context, resourceGroupName string, storageTransferFlowName string, storageTransferFlow StorageTransferFlowsPatch, options *StorageTransferFlowsClientBeginUpdateOptions) (*runtime.Poller[StorageTransferFlowsClientUpdateResponse], error)`
- New struct `ListTransferFlowsByPipelineRequestBody`
- New struct `ListTransferFlowsByPipelineResult`
- New struct `ListTransferFlowsParameterBody`
- New struct `ListTransferFlowsRequestBody`
- New struct `LocationMapping`
- New struct `MessagingTransferFlow`
- New struct `MessagingTransferFlowListResult`
- New struct `MessagingTransferFlowProperties`
- New struct `MessagingTransferFlowsPatch`
- New struct `StorageTransferFlow`
- New struct `StorageTransferFlowListResult`
- New struct `StorageTransferFlowProperties`
- New struct `StorageTransferFlowsPatch`
- New struct `TransferFlowIDBody`
- New struct `TransferFlowRequestBody`
- New struct `TransferFlowResponseBody`
- New struct `TransferFlowSummary`
- New field `MessagingLocationMappings`, `StorageLocationMappings` in struct `DirectionSettings`


## 0.1.0 (2026-10-15)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datatransfer/armdatatransfer` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).