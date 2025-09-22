# Release History

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
