# Release History

## 2.0.0 (2023-03-24)
### Breaking Changes

- Type of `ScalingPlanProperties.HostPoolType` has been changed from `*HostPoolType` to `*ScalingHostPoolType`
- Type of `ScalingSchedule.OffPeakStartTime` has been changed from `*time.Time` to `*Time`
- Type of `ScalingSchedule.PeakStartTime` has been changed from `*time.Time` to `*Time`
- Type of `ScalingSchedule.RampDownStartTime` has been changed from `*time.Time` to `*Time`
- Type of `ScalingSchedule.RampUpStartTime` has been changed from `*time.Time` to `*Time`
- Type alias `Operation` has been removed
- Operation `*DesktopsClient.List` has supported pagination, use `*DesktopsClient.NewListPager` instead.
- Operation `*OperationsClient.List` has supported pagination, use `*OperationsClient.NewListPager` instead.
- Struct `MigrationRequestProperties` has been removed
- Field `MigrationRequest` of struct `ApplicationGroupProperties` has been removed
- Field `MigrationRequest` of struct `HostPoolProperties` has been removed
- Field `HostPoolType` of struct `ScalingPlanPatchProperties` has been removed

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DayOfWeek` with values `DayOfWeekFriday`, `DayOfWeekMonday`, `DayOfWeekSaturday`, `DayOfWeekSunday`, `DayOfWeekThursday`, `DayOfWeekTuesday`, `DayOfWeekWednesday`
- New enum type `ScalingHostPoolType` with values `ScalingHostPoolTypePooled`
- New enum type `SessionHostComponentUpdateType` with values `SessionHostComponentUpdateTypeDefault`, `SessionHostComponentUpdateTypeScheduled`
- New function `NewScalingPlanPooledSchedulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ScalingPlanPooledSchedulesClient, error)`
- New function `*ScalingPlanPooledSchedulesClient.Create(context.Context, string, string, string, ScalingPlanPooledSchedule, *ScalingPlanPooledSchedulesClientCreateOptions) (ScalingPlanPooledSchedulesClientCreateResponse, error)`
- New function `*ScalingPlanPooledSchedulesClient.Delete(context.Context, string, string, string, *ScalingPlanPooledSchedulesClientDeleteOptions) (ScalingPlanPooledSchedulesClientDeleteResponse, error)`
- New function `*ScalingPlanPooledSchedulesClient.Get(context.Context, string, string, string, *ScalingPlanPooledSchedulesClientGetOptions) (ScalingPlanPooledSchedulesClientGetResponse, error)`
- New function `*ScalingPlanPooledSchedulesClient.NewListPager(string, string, *ScalingPlanPooledSchedulesClientListOptions) *runtime.Pager[ScalingPlanPooledSchedulesClientListResponse]`
- New function `*ScalingPlanPooledSchedulesClient.Update(context.Context, string, string, string, *ScalingPlanPooledSchedulesClientUpdateOptions) (ScalingPlanPooledSchedulesClientUpdateResponse, error)`
- New struct `AgentUpdatePatchProperties`
- New struct `AgentUpdateProperties`
- New struct `MaintenanceWindowPatchProperties`
- New struct `MaintenanceWindowProperties`
- New struct `ScalingPlanPooledSchedule`
- New struct `ScalingPlanPooledScheduleList`
- New struct `ScalingPlanPooledSchedulePatch`
- New struct `ScalingPlanPooledScheduleProperties`
- New struct `SystemData`
- New struct `Time`
- New field `SystemData` in struct `Application`
- New field `SystemData` in struct `ApplicationGroup`
- New field `InitialSkip` in struct `ApplicationGroupsClientListByResourceGroupOptions`
- New field `IsDescending` in struct `ApplicationGroupsClientListByResourceGroupOptions`
- New field `PageSize` in struct `ApplicationGroupsClientListByResourceGroupOptions`
- New field `InitialSkip` in struct `ApplicationsClientListOptions`
- New field `IsDescending` in struct `ApplicationsClientListOptions`
- New field `PageSize` in struct `ApplicationsClientListOptions`
- New field `SystemData` in struct `Desktop`
- New field `SystemData` in struct `HostPool`
- New field `AgentUpdate` in struct `HostPoolPatchProperties`
- New field `AgentUpdate` in struct `HostPoolProperties`
- New field `InitialSkip` in struct `HostPoolsClientListByResourceGroupOptions`
- New field `IsDescending` in struct `HostPoolsClientListByResourceGroupOptions`
- New field `PageSize` in struct `HostPoolsClientListByResourceGroupOptions`
- New field `InitialSkip` in struct `HostPoolsClientListOptions`
- New field `IsDescending` in struct `HostPoolsClientListOptions`
- New field `PageSize` in struct `HostPoolsClientListOptions`
- New field `SystemData` in struct `MSIXPackage`
- New field `InitialSkip` in struct `MSIXPackagesClientListOptions`
- New field `IsDescending` in struct `MSIXPackagesClientListOptions`
- New field `PageSize` in struct `MSIXPackagesClientListOptions`
- New field `NextLink` in struct `ResourceProviderOperationList`
- New field `SystemData` in struct `ScalingPlan`
- New field `InitialSkip` in struct `ScalingPlansClientListByHostPoolOptions`
- New field `IsDescending` in struct `ScalingPlansClientListByHostPoolOptions`
- New field `PageSize` in struct `ScalingPlansClientListByHostPoolOptions`
- New field `InitialSkip` in struct `ScalingPlansClientListByResourceGroupOptions`
- New field `IsDescending` in struct `ScalingPlansClientListByResourceGroupOptions`
- New field `PageSize` in struct `ScalingPlansClientListByResourceGroupOptions`
- New field `InitialSkip` in struct `ScalingPlansClientListBySubscriptionOptions`
- New field `IsDescending` in struct `ScalingPlansClientListBySubscriptionOptions`
- New field `PageSize` in struct `ScalingPlansClientListBySubscriptionOptions`
- New field `SystemData` in struct `SessionHost`
- New field `FriendlyName` in struct `SessionHostPatchProperties`
- New field `FriendlyName` in struct `SessionHostProperties`
- New field `InitialSkip` in struct `SessionHostsClientListOptions`
- New field `IsDescending` in struct `SessionHostsClientListOptions`
- New field `PageSize` in struct `SessionHostsClientListOptions`
- New field `Force` in struct `SessionHostsClientUpdateOptions`
- New field `InitialSkip` in struct `StartMenuItemsClientListOptions`
- New field `IsDescending` in struct `StartMenuItemsClientListOptions`
- New field `PageSize` in struct `StartMenuItemsClientListOptions`
- New field `SystemData` in struct `UserSession`
- New field `InitialSkip` in struct `UserSessionsClientListByHostPoolOptions`
- New field `IsDescending` in struct `UserSessionsClientListByHostPoolOptions`
- New field `PageSize` in struct `UserSessionsClientListByHostPoolOptions`
- New field `InitialSkip` in struct `UserSessionsClientListOptions`
- New field `IsDescending` in struct `UserSessionsClientListOptions`
- New field `PageSize` in struct `UserSessionsClientListOptions`
- New field `SystemData` in struct `Workspace`
- New field `InitialSkip` in struct `WorkspacesClientListByResourceGroupOptions`
- New field `IsDescending` in struct `WorkspacesClientListByResourceGroupOptions`
- New field `PageSize` in struct `WorkspacesClientListByResourceGroupOptions`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/desktopvirtualization/armdesktopvirtualization` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).