# Release History

## 2.0.0 (2022-11-04)
### Breaking Changes

- Type of `ScalingPlanProperties.HostPoolType` has been changed from `*HostPoolType` to `*ScalingHostPoolType`
- Type of `ScalingSchedule.OffPeakStartTime` has been changed from `*time.Time` to `*Time`
- Type of `ScalingSchedule.RampDownStartTime` has been changed from `*time.Time` to `*Time`
- Type of `ScalingSchedule.PeakStartTime` has been changed from `*time.Time` to `*Time`
- Type of `ScalingSchedule.RampUpStartTime` has been changed from `*time.Time` to `*Time`
- Const `OperationUnhide` has been removed
- Const `OperationHide` has been removed
- Const `OperationStart` has been removed
- Const `OperationRevoke` has been removed
- Const `OperationComplete` has been removed
- Type alias `Operation` has been removed
- Function `PossibleOperationValues` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*DesktopsClient.List` has been removed
- Struct `CloudError` has been removed
- Struct `MigrationRequestProperties` has been removed
- Field `HostPoolType` of struct `ScalingPlanPatchProperties` has been removed
- Field `MigrationRequest` of struct `ApplicationGroupProperties` has been removed
- Field `MigrationRequest` of struct `HostPoolProperties` has been removed

### Features Added

- New const `CreatedByTypeKey`
- New const `CreatedByTypeApplication`
- New const `DayOfWeekFriday`
- New const `DayOfWeekWednesday`
- New const `CreatedByTypeUser`
- New const `DayOfWeekMonday`
- New const `ScalingHostPoolTypePooled`
- New const `SessionHostComponentUpdateTypeScheduled`
- New const `DayOfWeekThursday`
- New const `DayOfWeekSaturday`
- New const `DayOfWeekSunday`
- New const `SessionHostComponentUpdateTypeDefault`
- New const `DayOfWeekTuesday`
- New const `CreatedByTypeManagedIdentity`
- New type alias `CreatedByType`
- New type alias `SessionHostComponentUpdateType`
- New type alias `DayOfWeek`
- New type alias `ScalingHostPoolType`
- New function `*ScalingPlanPooledSchedulesClient.NewListPager(string, string, *ScalingPlanPooledSchedulesClientListOptions) *runtime.Pager[ScalingPlanPooledSchedulesClientListResponse]`
- New function `NewScalingPlanPooledSchedulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ScalingPlanPooledSchedulesClient, error)`
- New function `*ScalingPlanPooledSchedulesClient.Create(context.Context, string, string, string, ScalingPlanPooledSchedule, *ScalingPlanPooledSchedulesClientCreateOptions) (ScalingPlanPooledSchedulesClientCreateResponse, error)`
- New function `*DesktopsClient.NewListPager(string, string, *DesktopsClientListOptions) *runtime.Pager[DesktopsClientListResponse]`
- New function `PossibleSessionHostComponentUpdateTypeValues() []SessionHostComponentUpdateType`
- New function `*ScalingPlanPooledSchedulesClient.Update(context.Context, string, string, string, *ScalingPlanPooledSchedulesClientUpdateOptions) (ScalingPlanPooledSchedulesClientUpdateResponse, error)`
- New function `*ScalingPlanPooledSchedulesClient.Get(context.Context, string, string, string, *ScalingPlanPooledSchedulesClientGetOptions) (ScalingPlanPooledSchedulesClientGetResponse, error)`
- New function `PossibleDayOfWeekValues() []DayOfWeek`
- New function `PossibleScalingHostPoolTypeValues() []ScalingHostPoolType`
- New function `*ScalingPlanPooledSchedulesClient.Delete(context.Context, string, string, string, *ScalingPlanPooledSchedulesClientDeleteOptions) (ScalingPlanPooledSchedulesClientDeleteResponse, error)`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New struct `AgentUpdatePatchProperties`
- New struct `AgentUpdateProperties`
- New struct `MaintenanceWindowPatchProperties`
- New struct `MaintenanceWindowProperties`
- New struct `ScalingPlanPooledSchedule`
- New struct `ScalingPlanPooledScheduleList`
- New struct `ScalingPlanPooledSchedulePatch`
- New struct `ScalingPlanPooledScheduleProperties`
- New struct `ScalingPlanPooledSchedulesClient`
- New struct `ScalingPlanPooledSchedulesClientCreateOptions`
- New struct `ScalingPlanPooledSchedulesClientCreateResponse`
- New struct `ScalingPlanPooledSchedulesClientDeleteOptions`
- New struct `ScalingPlanPooledSchedulesClientDeleteResponse`
- New struct `ScalingPlanPooledSchedulesClientGetOptions`
- New struct `ScalingPlanPooledSchedulesClientGetResponse`
- New struct `ScalingPlanPooledSchedulesClientListOptions`
- New struct `ScalingPlanPooledSchedulesClientListResponse`
- New struct `ScalingPlanPooledSchedulesClientUpdateOptions`
- New struct `ScalingPlanPooledSchedulesClientUpdateResponse`
- New struct `SystemData`
- New struct `Time`
- New field `SystemData` in struct `ApplicationGroup`
- New field `PageSize` in struct `ScalingPlansClientListBySubscriptionOptions`
- New field `InitialSkip` in struct `ScalingPlansClientListBySubscriptionOptions`
- New field `IsDescending` in struct `ScalingPlansClientListBySubscriptionOptions`
- New field `SystemData` in struct `Desktop`
- New field `PageSize` in struct `HostPoolsClientListByResourceGroupOptions`
- New field `InitialSkip` in struct `HostPoolsClientListByResourceGroupOptions`
- New field `IsDescending` in struct `HostPoolsClientListByResourceGroupOptions`
- New field `InitialSkip` in struct `SessionHostsClientListOptions`
- New field `IsDescending` in struct `SessionHostsClientListOptions`
- New field `PageSize` in struct `SessionHostsClientListOptions`
- New field `PageSize` in struct `ScalingPlansClientListByHostPoolOptions`
- New field `InitialSkip` in struct `ScalingPlansClientListByHostPoolOptions`
- New field `IsDescending` in struct `ScalingPlansClientListByHostPoolOptions`
- New field `FriendlyName` in struct `SessionHostPatchProperties`
- New field `SystemData` in struct `ScalingPlan`
- New field `PageSize` in struct `ScalingPlansClientListByResourceGroupOptions`
- New field `InitialSkip` in struct `ScalingPlansClientListByResourceGroupOptions`
- New field `IsDescending` in struct `ScalingPlansClientListByResourceGroupOptions`
- New field `PageSize` in struct `HostPoolsClientListOptions`
- New field `InitialSkip` in struct `HostPoolsClientListOptions`
- New field `IsDescending` in struct `HostPoolsClientListOptions`
- New field `SystemData` in struct `Application`
- New field `NextLink` in struct `ResourceProviderOperationList`
- New field `IsDescending` in struct `UserSessionsClientListByHostPoolOptions`
- New field `PageSize` in struct `UserSessionsClientListByHostPoolOptions`
- New field `InitialSkip` in struct `UserSessionsClientListByHostPoolOptions`
- New field `Force` in struct `SessionHostsClientUpdateOptions`
- New field `AgentUpdate` in struct `HostPoolProperties`
- New field `IsDescending` in struct `WorkspacesClientListByResourceGroupOptions`
- New field `PageSize` in struct `WorkspacesClientListByResourceGroupOptions`
- New field `InitialSkip` in struct `WorkspacesClientListByResourceGroupOptions`
- New field `SystemData` in struct `SessionHost`
- New field `SystemData` in struct `UserSession`
- New field `SystemData` in struct `MSIXPackage`
- New field `InitialSkip` in struct `DesktopsClientListOptions`
- New field `IsDescending` in struct `DesktopsClientListOptions`
- New field `PageSize` in struct `DesktopsClientListOptions`
- New field `IsDescending` in struct `UserSessionsClientListOptions`
- New field `PageSize` in struct `UserSessionsClientListOptions`
- New field `InitialSkip` in struct `UserSessionsClientListOptions`
- New field `SystemData` in struct `HostPool`
- New field `SystemData` in struct `Workspace`
- New field `InitialSkip` in struct `MSIXPackagesClientListOptions`
- New field `IsDescending` in struct `MSIXPackagesClientListOptions`
- New field `PageSize` in struct `MSIXPackagesClientListOptions`
- New field `IsDescending` in struct `StartMenuItemsClientListOptions`
- New field `PageSize` in struct `StartMenuItemsClientListOptions`
- New field `InitialSkip` in struct `StartMenuItemsClientListOptions`
- New field `AgentUpdate` in struct `HostPoolPatchProperties`
- New field `InitialSkip` in struct `ApplicationsClientListOptions`
- New field `IsDescending` in struct `ApplicationsClientListOptions`
- New field `PageSize` in struct `ApplicationsClientListOptions`
- New field `FriendlyName` in struct `SessionHostProperties`
- New field `InitialSkip` in struct `ApplicationGroupsClientListByResourceGroupOptions`
- New field `IsDescending` in struct `ApplicationGroupsClientListByResourceGroupOptions`
- New field `PageSize` in struct `ApplicationGroupsClientListByResourceGroupOptions`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/desktopvirtualization/armdesktopvirtualization` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).