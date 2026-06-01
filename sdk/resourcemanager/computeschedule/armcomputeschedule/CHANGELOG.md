# Release History

## 2.0.0-beta.1 (2026-06-01)
### Breaking Changes

- Function `*ScheduledActionsClient.VirtualMachinesCancelOperations` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody CancelOperationsContent, options *ScheduledActionsClientVirtualMachinesCancelOperationsOptions)` to `(ctx context.Context, locationparameter string, requestBody CancelOperationsRequest, options *ScheduledActionsClientVirtualMachinesCancelOperationsOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteCreate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteCreateContent, options *ScheduledActionsClientVirtualMachinesExecuteCreateOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteCreateRequest, options *ScheduledActionsClientVirtualMachinesExecuteCreateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteDeallocate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteDeallocateContent, options *ScheduledActionsClientVirtualMachinesExecuteDeallocateOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteDeallocateRequest, options *ScheduledActionsClientVirtualMachinesExecuteDeallocateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteDelete` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteDeleteContent, options *ScheduledActionsClientVirtualMachinesExecuteDeleteOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteDeleteRequest, options *ScheduledActionsClientVirtualMachinesExecuteDeleteOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteHibernate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteHibernateContent, options *ScheduledActionsClientVirtualMachinesExecuteHibernateOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteHibernateRequest, options *ScheduledActionsClientVirtualMachinesExecuteHibernateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesExecuteStart` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody ExecuteStartContent, options *ScheduledActionsClientVirtualMachinesExecuteStartOptions)` to `(ctx context.Context, locationparameter string, requestBody ExecuteStartRequest, options *ScheduledActionsClientVirtualMachinesExecuteStartOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesGetOperationErrors` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody GetOperationErrorsContent, options *ScheduledActionsClientVirtualMachinesGetOperationErrorsOptions)` to `(ctx context.Context, locationparameter string, requestBody GetOperationErrorsRequest, options *ScheduledActionsClientVirtualMachinesGetOperationErrorsOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesGetOperationStatus` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody GetOperationStatusContent, options *ScheduledActionsClientVirtualMachinesGetOperationStatusOptions)` to `(ctx context.Context, locationparameter string, requestBody GetOperationStatusRequest, options *ScheduledActionsClientVirtualMachinesGetOperationStatusOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesSubmitDeallocate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody SubmitDeallocateContent, options *ScheduledActionsClientVirtualMachinesSubmitDeallocateOptions)` to `(ctx context.Context, locationparameter string, requestBody SubmitDeallocateRequest, options *ScheduledActionsClientVirtualMachinesSubmitDeallocateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesSubmitHibernate` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody SubmitHibernateContent, options *ScheduledActionsClientVirtualMachinesSubmitHibernateOptions)` to `(ctx context.Context, locationparameter string, requestBody SubmitHibernateRequest, options *ScheduledActionsClientVirtualMachinesSubmitHibernateOptions)`
- Function `*ScheduledActionsClient.VirtualMachinesSubmitStart` parameter(s) have been changed from `(ctx context.Context, locationparameter string, requestBody SubmitStartContent, options *ScheduledActionsClientVirtualMachinesSubmitStartOptions)` to `(ctx context.Context, locationparameter string, requestBody SubmitStartRequest, options *ScheduledActionsClientVirtualMachinesSubmitStartOptions)`
- Struct `CancelOperationsContent` has been removed
- Struct `ExecuteCreateContent` has been removed
- Struct `ExecuteDeallocateContent` has been removed
- Struct `ExecuteDeleteContent` has been removed
- Struct `ExecuteHibernateContent` has been removed
- Struct `ExecuteStartContent` has been removed
- Struct `GetOperationErrorsContent` has been removed
- Struct `GetOperationStatusContent` has been removed
- Struct `SubmitDeallocateContent` has been removed
- Struct `SubmitHibernateContent` has been removed
- Struct `SubmitStartContent` has been removed

### Features Added

- New value `ResourceOperationTypeCreate`, `ResourceOperationTypeDelete` added to enum type `ResourceOperationType`
- New enum type `AllocationStrategy` with values `AllocationStrategyCapacityOptimized`, `AllocationStrategyLowestPrice`, `AllocationStrategyPrioritized`
- New enum type `DistributionStrategy` with values `DistributionStrategyBestEffortBalanced`, `DistributionStrategyBestEffortSingleZone`, `DistributionStrategyPrioritized`, `DistributionStrategyStrictBalanced`
- New enum type `OsType` with values `OsTypeLinux`, `OsTypeWindows`
- New enum type `PriorityType` with values `PriorityTypeRegular`, `PriorityTypeSpot`
- New function `*ScheduledActionsClient.VirtualMachinesExecuteCreateFlex(ctx context.Context, locationparameter string, body ExecuteCreateFlexRequest, options *ScheduledActionsClientVirtualMachinesExecuteCreateFlexOptions) (ScheduledActionsClientVirtualMachinesExecuteCreateFlexResponse, error)`
- New struct `CancelOperationsRequest`
- New struct `CreateFlexResourceOperationResponse`
- New struct `ExecuteCreateFlexRequest`
- New struct `ExecuteCreateRequest`
- New struct `ExecuteDeallocateRequest`
- New struct `ExecuteDeleteRequest`
- New struct `ExecuteHibernateRequest`
- New struct `ExecuteStartRequest`
- New struct `FallbackOperationInfo`
- New struct `FlexProperties`
- New struct `GetOperationErrorsRequest`
- New struct `GetOperationStatusRequest`
- New struct `PriorityProfile`
- New struct `ResourceProvisionFlexPayload`
- New struct `SubmitDeallocateRequest`
- New struct `SubmitHibernateRequest`
- New struct `SubmitStartRequest`
- New struct `VMSizeProfile`
- New struct `ZoneAllocationPolicy`
- New struct `ZonePreference`
- New field `FallbackOperationInfo` in struct `ResourceOperationDetails`
- New field `OnFailureAction` in struct `RetryPolicy`


## 1.2.0-beta.2 (2026-04-03)
### Features Added

- New value `ResourceOperationTypeCreate`, `ResourceOperationTypeDelete` added to enum type `ResourceOperationType`
- New enum type `AllocationStrategy` with values `AllocationStrategyCapacityOptimized`, `AllocationStrategyLowestPrice`, `AllocationStrategyPrioritized`
- New enum type `DistributionStrategy` with values `DistributionStrategyBestEffortBalanced`, `DistributionStrategyBestEffortSingleZone`, `DistributionStrategyPrioritized`, `DistributionStrategyStrictBalanced`
- New enum type `OsType` with values `OsTypeLinux`, `OsTypeWindows`
- New enum type `PriorityType` with values `PriorityTypeRegular`, `PriorityTypeSpot`
- New function `*ScheduledActionsClient.VirtualMachinesExecuteCreateFlex(ctx context.Context, locationparameter string, body ExecuteCreateFlexContent, options *ScheduledActionsClientVirtualMachinesExecuteCreateFlexOptions) (ScheduledActionsClientVirtualMachinesExecuteCreateFlexResponse, error)`
- New struct `CreateFlexResourceOperationResponse`
- New struct `ExecuteCreateFlexContent`
- New struct `FallbackOperationInfo`
- New struct `FlexProperties`
- New struct `PriorityProfile`
- New struct `ResourceProvisionFlexPayload`
- New struct `VMSizeProfile`
- New struct `ZoneAllocationPolicy`
- New struct `ZonePreference`
- New field `FallbackOperationInfo` in struct `ResourceOperationDetails`
- New field `OnFailureAction` in struct `RetryPolicy`


## 1.2.0-beta.1 (2025-07-24)
### Features Added

- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `Language` with values `LanguageEnUs`
- New enum type `Month` with values `MonthAll`, `MonthApril`, `MonthAugust`, `MonthDecember`, `MonthFebruary`, `MonthJanuary`, `MonthJuly`, `MonthJune`, `MonthMarch`, `MonthMay`, `MonthNovember`, `MonthOctober`, `MonthSeptember`
- New enum type `NotificationType` with values `NotificationTypeEmail`
- New enum type `OccurrenceState` with values `OccurrenceStateCanceled`, `OccurrenceStateCancelling`, `OccurrenceStateCreated`, `OccurrenceStateFailed`, `OccurrenceStateRescheduling`, `OccurrenceStateScheduled`, `OccurrenceStateSucceeded`
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`
- New enum type `ResourceOperationStatus` with values `ResourceOperationStatusFailed`, `ResourceOperationStatusSucceeded`
- New enum type `ResourceProvisioningState` with values `ResourceProvisioningStateCanceled`, `ResourceProvisioningStateFailed`, `ResourceProvisioningStateSucceeded`
- New enum type `ResourceType` with values `ResourceTypeVirtualMachine`, `ResourceTypeVirtualMachineScaleSet`
- New enum type `ScheduledActionType` with values `ScheduledActionTypeDeallocate`, `ScheduledActionTypeHibernate`, `ScheduledActionTypeStart`
- New enum type `WeekDay` with values `WeekDayAll`, `WeekDayFriday`, `WeekDayMonday`, `WeekDaySaturday`, `WeekDaySunday`, `WeekDayThursday`, `WeekDayTuesday`, `WeekDayWednesday`
- New function `*ClientFactory.NewOccurrenceExtensionClient() *OccurrenceExtensionClient`
- New function `*ClientFactory.NewOccurrencesClient() *OccurrencesClient`
- New function `*ClientFactory.NewScheduledActionExtensionClient() *ScheduledActionExtensionClient`
- New function `NewOccurrenceExtensionClient(azcore.TokenCredential, *arm.ClientOptions) (*OccurrenceExtensionClient, error)`
- New function `*OccurrenceExtensionClient.NewListOccurrenceByVMsPager(string, *OccurrenceExtensionClientListOccurrenceByVMsOptions) *runtime.Pager[OccurrenceExtensionClientListOccurrenceByVMsResponse]`
- New function `NewOccurrencesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OccurrencesClient, error)`
- New function `*OccurrencesClient.Cancel(context.Context, string, string, string, CancelOccurrenceRequest, *OccurrencesClientCancelOptions) (OccurrencesClientCancelResponse, error)`
- New function `*OccurrencesClient.BeginDelay(context.Context, string, string, string, DelayRequest, *OccurrencesClientBeginDelayOptions) (*runtime.Poller[OccurrencesClientDelayResponse], error)`
- New function `*OccurrencesClient.Get(context.Context, string, string, string, *OccurrencesClientGetOptions) (OccurrencesClientGetResponse, error)`
- New function `*OccurrencesClient.NewListByScheduledActionPager(string, string, *OccurrencesClientListByScheduledActionOptions) *runtime.Pager[OccurrencesClientListByScheduledActionResponse]`
- New function `*OccurrencesClient.NewListResourcesPager(string, string, string, *OccurrencesClientListResourcesOptions) *runtime.Pager[OccurrencesClientListResourcesResponse]`
- New function `NewScheduledActionExtensionClient(azcore.TokenCredential, *arm.ClientOptions) (*ScheduledActionExtensionClient, error)`
- New function `*ScheduledActionExtensionClient.NewListByVMsPager(string, *ScheduledActionExtensionClientListByVMsOptions) *runtime.Pager[ScheduledActionExtensionClientListByVMsResponse]`
- New function `*ScheduledActionsClient.AttachResources(context.Context, string, string, ResourceAttachRequest, *ScheduledActionsClientAttachResourcesOptions) (ScheduledActionsClientAttachResourcesResponse, error)`
- New function `*ScheduledActionsClient.CancelNextOccurrence(context.Context, string, string, CancelOccurrenceRequest, *ScheduledActionsClientCancelNextOccurrenceOptions) (ScheduledActionsClientCancelNextOccurrenceResponse, error)`
- New function `*ScheduledActionsClient.BeginCreateOrUpdate(context.Context, string, string, ScheduledAction, *ScheduledActionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ScheduledActionsClientCreateOrUpdateResponse], error)`
- New function `*ScheduledActionsClient.BeginDelete(context.Context, string, string, *ScheduledActionsClientBeginDeleteOptions) (*runtime.Poller[ScheduledActionsClientDeleteResponse], error)`
- New function `*ScheduledActionsClient.DetachResources(context.Context, string, string, ResourceDetachRequest, *ScheduledActionsClientDetachResourcesOptions) (ScheduledActionsClientDetachResourcesResponse, error)`
- New function `*ScheduledActionsClient.Disable(context.Context, string, string, *ScheduledActionsClientDisableOptions) (ScheduledActionsClientDisableResponse, error)`
- New function `*ScheduledActionsClient.Enable(context.Context, string, string, *ScheduledActionsClientEnableOptions) (ScheduledActionsClientEnableResponse, error)`
- New function `*ScheduledActionsClient.Get(context.Context, string, string, *ScheduledActionsClientGetOptions) (ScheduledActionsClientGetResponse, error)`
- New function `*ScheduledActionsClient.NewListByResourceGroupPager(string, *ScheduledActionsClientListByResourceGroupOptions) *runtime.Pager[ScheduledActionsClientListByResourceGroupResponse]`
- New function `*ScheduledActionsClient.NewListBySubscriptionPager(*ScheduledActionsClientListBySubscriptionOptions) *runtime.Pager[ScheduledActionsClientListBySubscriptionResponse]`
- New function `*ScheduledActionsClient.NewListResourcesPager(string, string, *ScheduledActionsClientListResourcesOptions) *runtime.Pager[ScheduledActionsClientListResourcesResponse]`
- New function `*ScheduledActionsClient.PatchResources(context.Context, string, string, ResourcePatchRequest, *ScheduledActionsClientPatchResourcesOptions) (ScheduledActionsClientPatchResourcesResponse, error)`
- New function `*ScheduledActionsClient.TriggerManualOccurrence(context.Context, string, string, *ScheduledActionsClientTriggerManualOccurrenceOptions) (ScheduledActionsClientTriggerManualOccurrenceResponse, error)`
- New function `*ScheduledActionsClient.Update(context.Context, string, string, ScheduledActionUpdate, *ScheduledActionsClientUpdateOptions) (ScheduledActionsClientUpdateResponse, error)`
- New struct `CancelOccurrenceRequest`
- New struct `DelayRequest`
- New struct `Error`
- New struct `InnerError`
- New struct `NotificationProperties`
- New struct `Occurrence`
- New struct `OccurrenceExtensionProperties`
- New struct `OccurrenceExtensionResource`
- New struct `OccurrenceExtensionResourceListResult`
- New struct `OccurrenceListResult`
- New struct `OccurrenceProperties`
- New struct `OccurrenceResource`
- New struct `OccurrenceResourceListResponse`
- New struct `OccurrenceResultSummary`
- New struct `RecurringActionsResourceOperationResult`
- New struct `ResourceAttachRequest`
- New struct `ResourceDetachRequest`
- New struct `ResourceListResponse`
- New struct `ResourcePatchRequest`
- New struct `ResourceResultSummary`
- New struct `ResourceStatus`
- New struct `ScheduledAction`
- New struct `ScheduledActionListResult`
- New struct `ScheduledActionProperties`
- New struct `ScheduledActionResource`
- New struct `ScheduledActionResources`
- New struct `ScheduledActionResourcesListResult`
- New struct `ScheduledActionUpdate`
- New struct `ScheduledActionUpdateProperties`
- New struct `ScheduledActionsSchedule`
- New struct `SystemData`


## 1.1.0 (2025-07-15)
### Features Added

- New function `*ScheduledActionsClient.VirtualMachinesExecuteCreate(context.Context, string, ExecuteCreateRequest, *ScheduledActionsClientVirtualMachinesExecuteCreateOptions) (ScheduledActionsClientVirtualMachinesExecuteCreateResponse, error)`
- New function `*ScheduledActionsClient.VirtualMachinesExecuteDelete(context.Context, string, ExecuteDeleteRequest, *ScheduledActionsClientVirtualMachinesExecuteDeleteOptions) (ScheduledActionsClientVirtualMachinesExecuteDeleteResponse, error)`
- New struct `CreateResourceOperationResponse`
- New struct `DeleteResourceOperationResponse`
- New struct `ExecuteCreateRequest`
- New struct `ExecuteDeleteRequest`
- New struct `ResourceProvisionPayload`


## 1.0.0 (2025-01-24)
### Breaking Changes

- Type of `OperationErrorDetails.ErrorDetails` has been changed from `*time.Time` to `*string`

### Features Added

- New field `AzureOperationName`, `Timestamp` in struct `OperationErrorDetails`
- New field `Timezone` in struct `ResourceOperationDetails`
- New field `Deadline`, `Timezone` in struct `Schedule`


## 0.1.0 (2024-09-27)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/computeschedule/armcomputeschedule` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
