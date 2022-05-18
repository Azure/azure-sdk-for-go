# Release History

## 1.0.0 (2022-05-18)
### Breaking Changes

- Function `*Client.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[ClientDeleteResponse], error)` to `(*runtime.Poller[ClientDeleteResponse], error)`
- Function `*Client.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[ClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[ClientCreateOrUpdateResponse], error)`
- Function `ManagementGroupDetails.MarshalJSON` has been removed
- Function `EntityListResult.MarshalJSON` has been removed
- Function `ManagementGroupProperties.MarshalJSON` has been removed
- Function `HierarchySettingsList.MarshalJSON` has been removed
- Function `ListSubscriptionUnderManagementGroup.MarshalJSON` has been removed
- Function `ManagementGroupListResult.MarshalJSON` has been removed
- Function `ManagementGroupChildInfo.MarshalJSON` has been removed
- Function `OperationListResult.MarshalJSON` has been removed
- Function `EntityHierarchyItemProperties.MarshalJSON` has been removed
- Function `DescendantListResult.MarshalJSON` has been removed
- Function `EntityInfoProperties.MarshalJSON` has been removed


## 0.6.0 (2022-04-18)
### Breaking Changes

- Function `*Client.List` has been removed
- Function `*Client.GetDescendants` has been removed
- Function `*ManagementGroupSubscriptionsClient.GetSubscriptionsUnderManagementGroup` has been removed
- Function `*EntitiesClient.List` has been removed
- Function `*OperationsClient.List` has been removed

### Features Added

- New function `*Client.NewGetDescendantsPager(string, *ClientGetDescendantsOptions) *runtime.Pager[ClientGetDescendantsResponse]`
- New function `*Client.NewListPager(*ClientListOptions) *runtime.Pager[ClientListResponse]`
- New function `*ManagementGroupSubscriptionsClient.NewGetSubscriptionsUnderManagementGroupPager(string, *ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupOptions) *runtime.Pager[ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResponse]`
- New function `*EntitiesClient.NewListPager(*EntitiesClientListOptions) *runtime.Pager[EntitiesClientListResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.5.0 (2022-04-12)
### Breaking Changes

- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*EntitiesClient.List` return value(s) have been changed from `(*EntitiesClientListPager)` to `(*runtime.Pager[EntitiesClientListResponse])`
- Function `*ManagementGroupSubscriptionsClient.GetSubscriptionsUnderManagementGroup` return value(s) have been changed from `(*ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager)` to `(*runtime.Pager[ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResponse])`
- Function `NewAPIClient` return value(s) have been changed from `(*APIClient)` to `(*APIClient, error)`
- Function `NewClient` return value(s) have been changed from `(*Client)` to `(*Client, error)`
- Function `*Client.BeginCreateOrUpdate` return value(s) have been changed from `(ClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[ClientCreateOrUpdateResponse], error)`
- Function `NewManagementGroupSubscriptionsClient` return value(s) have been changed from `(*ManagementGroupSubscriptionsClient)` to `(*ManagementGroupSubscriptionsClient, error)`
- Function `NewEntitiesClient` return value(s) have been changed from `(*EntitiesClient)` to `(*EntitiesClient, error)`
- Function `*Client.BeginDelete` return value(s) have been changed from `(ClientDeletePollerResponse, error)` to `(*armruntime.Poller[ClientDeleteResponse], error)`
- Function `*Client.GetDescendants` return value(s) have been changed from `(*ClientGetDescendantsPager)` to `(*runtime.Pager[ClientGetDescendantsResponse])`
- Function `NewHierarchySettingsClient` return value(s) have been changed from `(*HierarchySettingsClient)` to `(*HierarchySettingsClient, error)`
- Function `*Client.List` return value(s) have been changed from `(*ClientListPager)` to `(*runtime.Pager[ClientListResponse])`
- Function `Permissions.ToPtr` has been removed
- Function `*ClientGetDescendantsPager.PageResponse` has been removed
- Function `*EntitiesClientListPager.PageResponse` has been removed
- Function `Status.ToPtr` has been removed
- Function `*ClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*EntitiesClientListPager.Err` has been removed
- Function `*ClientListPager.Err` has been removed
- Function `Reason.ToPtr` has been removed
- Function `*ClientListPager.PageResponse` has been removed
- Function `*ClientCreateOrUpdatePoller.Done` has been removed
- Function `*ClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `EntityViewParameterType.ToPtr` has been removed
- Function `ClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*ClientDeletePoller.Poll` has been removed
- Function `*ClientDeletePollerResponse.Resume` has been removed
- Function `EntitySearchType.ToPtr` has been removed
- Function `*ClientCreateOrUpdatePoller.Poll` has been removed
- Function `*ClientListPager.NextPage` has been removed
- Function `*ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager.PageResponse` has been removed
- Function `*ClientGetDescendantsPager.NextPage` has been removed
- Function `ManagementGroupChildType.ToPtr` has been removed
- Function `*ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager.Err` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*ClientDeletePoller.ResumeToken` has been removed
- Function `*EntitiesClientListPager.NextPage` has been removed
- Function `*ClientGetDescendantsPager.Err` has been removed
- Function `ClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*ClientDeletePoller.Done` has been removed
- Function `*ClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `ManagementGroupExpandType.ToPtr` has been removed
- Function `*ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager.NextPage` has been removed
- Function `*ClientDeletePoller.FinalResponse` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Struct `APIClientCheckNameAvailabilityResult` has been removed
- Struct `APIClientStartTenantBackfillResult` has been removed
- Struct `APIClientTenantBackfillStatusResult` has been removed
- Struct `ClientCreateOrUpdatePoller` has been removed
- Struct `ClientCreateOrUpdatePollerResponse` has been removed
- Struct `ClientCreateOrUpdateResult` has been removed
- Struct `ClientDeletePoller` has been removed
- Struct `ClientDeletePollerResponse` has been removed
- Struct `ClientDeleteResult` has been removed
- Struct `ClientGetDescendantsPager` has been removed
- Struct `ClientGetDescendantsResult` has been removed
- Struct `ClientGetResult` has been removed
- Struct `ClientListPager` has been removed
- Struct `ClientListResult` has been removed
- Struct `ClientUpdateResult` has been removed
- Struct `EntitiesClientListPager` has been removed
- Struct `EntitiesClientListResult` has been removed
- Struct `HierarchySettingsClientCreateOrUpdateResult` has been removed
- Struct `HierarchySettingsClientGetResult` has been removed
- Struct `HierarchySettingsClientListResult` has been removed
- Struct `HierarchySettingsClientUpdateResult` has been removed
- Struct `ManagementGroupSubscriptionsClientCreateResult` has been removed
- Struct `ManagementGroupSubscriptionsClientGetSubscriptionResult` has been removed
- Struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager` has been removed
- Struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `HierarchySettingsClientListResult` of struct `HierarchySettingsClientListResponse` has been removed
- Field `RawResponse` of struct `HierarchySettingsClientListResponse` has been removed
- Field `ManagementGroupSubscriptionsClientCreateResult` of struct `ManagementGroupSubscriptionsClientCreateResponse` has been removed
- Field `RawResponse` of struct `ManagementGroupSubscriptionsClientCreateResponse` has been removed
- Field `EntitiesClientListResult` of struct `EntitiesClientListResponse` has been removed
- Field `RawResponse` of struct `EntitiesClientListResponse` has been removed
- Field `HierarchySettingsClientGetResult` of struct `HierarchySettingsClientGetResponse` has been removed
- Field `RawResponse` of struct `HierarchySettingsClientGetResponse` has been removed
- Field `ClientDeleteResult` of struct `ClientDeleteResponse` has been removed
- Field `RawResponse` of struct `ClientDeleteResponse` has been removed
- Field `HierarchySettingsClientCreateOrUpdateResult` of struct `HierarchySettingsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `HierarchySettingsClientCreateOrUpdateResponse` has been removed
- Field `ClientListResult` of struct `ClientListResponse` has been removed
- Field `RawResponse` of struct `ClientListResponse` has been removed
- Field `ClientGetResult` of struct `ClientGetResponse` has been removed
- Field `RawResponse` of struct `ClientGetResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `APIClientStartTenantBackfillResult` of struct `APIClientStartTenantBackfillResponse` has been removed
- Field `RawResponse` of struct `APIClientStartTenantBackfillResponse` has been removed
- Field `RawResponse` of struct `ManagementGroupSubscriptionsClientDeleteResponse` has been removed
- Field `ClientGetDescendantsResult` of struct `ClientGetDescendantsResponse` has been removed
- Field `RawResponse` of struct `ClientGetDescendantsResponse` has been removed
- Field `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResult` of struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResponse` has been removed
- Field `RawResponse` of struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResponse` has been removed
- Field `HierarchySettingsClientUpdateResult` of struct `HierarchySettingsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `HierarchySettingsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `HierarchySettingsClientDeleteResponse` has been removed
- Field `ClientCreateOrUpdateResult` of struct `ClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ClientCreateOrUpdateResponse` has been removed
- Field `ManagementGroupSubscriptionsClientGetSubscriptionResult` of struct `ManagementGroupSubscriptionsClientGetSubscriptionResponse` has been removed
- Field `RawResponse` of struct `ManagementGroupSubscriptionsClientGetSubscriptionResponse` has been removed
- Field `APIClientCheckNameAvailabilityResult` of struct `APIClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `APIClientCheckNameAvailabilityResponse` has been removed
- Field `ClientUpdateResult` of struct `ClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ClientUpdateResponse` has been removed
- Field `APIClientTenantBackfillStatusResult` of struct `APIClientTenantBackfillStatusResponse` has been removed
- Field `RawResponse` of struct `APIClientTenantBackfillStatusResponse` has been removed

### Features Added

- New function `EntityHierarchyItemProperties.MarshalJSON() ([]byte, error)`
- New struct `EntityHierarchyItem`
- New struct `EntityHierarchyItemProperties`
- New struct `ErrorDetails`
- New struct `ErrorResponse`
- New struct `OperationResults`
- New anonymous field `CheckNameAvailabilityResult` in struct `APIClientCheckNameAvailabilityResponse`
- New anonymous field `SubscriptionUnderManagementGroup` in struct `ManagementGroupSubscriptionsClientCreateResponse`
- New anonymous field `TenantBackfillStatusResult` in struct `APIClientStartTenantBackfillResponse`
- New field `ResumeToken` in struct `ClientBeginCreateOrUpdateOptions`
- New anonymous field `HierarchySettingsList` in struct `HierarchySettingsClientListResponse`
- New anonymous field `ManagementGroup` in struct `ClientUpdateResponse`
- New anonymous field `HierarchySettings` in struct `HierarchySettingsClientGetResponse`
- New anonymous field `DescendantListResult` in struct `ClientGetDescendantsResponse`
- New anonymous field `HierarchySettings` in struct `HierarchySettingsClientCreateOrUpdateResponse`
- New field `ResumeToken` in struct `ClientBeginDeleteOptions`
- New anonymous field `ManagementGroup` in struct `ClientCreateOrUpdateResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `ListSubscriptionUnderManagementGroup` in struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResponse`
- New anonymous field `ManagementGroup` in struct `ClientGetResponse`
- New anonymous field `AzureAsyncOperationResults` in struct `ClientDeleteResponse`
- New anonymous field `SubscriptionUnderManagementGroup` in struct `ManagementGroupSubscriptionsClientGetSubscriptionResponse`
- New anonymous field `HierarchySettings` in struct `HierarchySettingsClientUpdateResponse`
- New anonymous field `EntityListResult` in struct `EntitiesClientListResponse`
- New anonymous field `ManagementGroupListResult` in struct `ClientListResponse`
- New anonymous field `TenantBackfillStatusResult` in struct `APIClientTenantBackfillStatusResponse`


## 0.4.0 (2022-02-22)
### Breaking Changes

- Type of `ClientGetOptions.Expand` has been changed from `*Enum0` to `*ManagementGroupExpandType`
- Type of `EntitiesClientListOptions.Search` has been changed from `*Enum2` to `*EntitySearchType`
- Type of `EntitiesClientListOptions.View` has been changed from `*Enum3` to `*EntityViewParameterType`
- Const `Enum2ChildrenOnly` has been removed
- Const `Enum0Path` has been removed
- Const `Enum3FullHierarchy` has been removed
- Const `Enum2ParentAndFirstLevelChildren` has been removed
- Const `Enum2ParentOnly` has been removed
- Const `Enum3Audit` has been removed
- Const `Enum3SubscriptionsOnly` has been removed
- Const `Enum2AllowedChildren` has been removed
- Const `Enum0Ancestors` has been removed
- Const `Enum3GroupsOnly` has been removed
- Const `Enum0Children` has been removed
- Const `Enum2AllowedParents` has been removed
- Function `PossibleEnum3Values` has been removed
- Function `PossibleEnum2Values` has been removed
- Function `Enum3.ToPtr` has been removed
- Function `Enum0.ToPtr` has been removed
- Function `Enum2.ToPtr` has been removed
- Function `PossibleEnum0Values` has been removed
- Function `EntityHierarchyItemProperties.MarshalJSON` has been removed
- Struct `EntityHierarchyItem` has been removed
- Struct `EntityHierarchyItemProperties` has been removed
- Struct `ErrorDetails` has been removed
- Struct `ErrorResponse` has been removed
- Struct `OperationResults` has been removed

### Features Added

- New const `EntitySearchTypeAllowedChildren`
- New const `EntityViewParameterTypeSubscriptionsOnly`
- New const `EntitySearchTypeAllowedParents`
- New const `ManagementGroupExpandTypePath`
- New const `EntitySearchTypeParentOnly`
- New const `EntityViewParameterTypeAudit`
- New const `EntitySearchTypeParentAndFirstLevelChildren`
- New const `ManagementGroupExpandTypeChildren`
- New const `EntityViewParameterTypeFullHierarchy`
- New const `EntitySearchTypeChildrenOnly`
- New const `EntityViewParameterTypeGroupsOnly`
- New const `ManagementGroupExpandTypeAncestors`
- New function `EntityViewParameterType.ToPtr() *EntityViewParameterType`
- New function `PossibleManagementGroupExpandTypeValues() []ManagementGroupExpandType`
- New function `PossibleEntitySearchTypeValues() []EntitySearchType`
- New function `PossibleEntityViewParameterTypeValues() []EntityViewParameterType`
- New function `ManagementGroupExpandType.ToPtr() *ManagementGroupExpandType`
- New function `EntitySearchType.ToPtr() *EntitySearchType`


## 0.3.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.3.0 (2022-01-13)
### Breaking Changes

- Function `*ManagementGroupSubscriptionsClient.Create` parameter(s) have been changed from `(context.Context, string, string, *ManagementGroupSubscriptionsCreateOptions)` to `(context.Context, string, string, *ManagementGroupSubscriptionsClientCreateOptions)`
- Function `*ManagementGroupSubscriptionsClient.Create` return value(s) have been changed from `(ManagementGroupSubscriptionsCreateResponse, error)` to `(ManagementGroupSubscriptionsClientCreateResponse, error)`
- Function `*ManagementGroupSubscriptionsClient.GetSubscriptionsUnderManagementGroup` parameter(s) have been changed from `(string, *ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupOptions)` to `(string, *ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupOptions)`
- Function `*ManagementGroupSubscriptionsClient.GetSubscriptionsUnderManagementGroup` return value(s) have been changed from `(*ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupPager)` to `(*ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager)`
- Function `*ManagementGroupSubscriptionsClient.GetSubscription` parameter(s) have been changed from `(context.Context, string, string, *ManagementGroupSubscriptionsGetSubscriptionOptions)` to `(context.Context, string, string, *ManagementGroupSubscriptionsClientGetSubscriptionOptions)`
- Function `*ManagementGroupSubscriptionsClient.GetSubscription` return value(s) have been changed from `(ManagementGroupSubscriptionsGetSubscriptionResponse, error)` to `(ManagementGroupSubscriptionsClientGetSubscriptionResponse, error)`
- Function `*HierarchySettingsClient.Delete` parameter(s) have been changed from `(context.Context, string, *HierarchySettingsDeleteOptions)` to `(context.Context, string, *HierarchySettingsClientDeleteOptions)`
- Function `*HierarchySettingsClient.Delete` return value(s) have been changed from `(HierarchySettingsDeleteResponse, error)` to `(HierarchySettingsClientDeleteResponse, error)`
- Function `*HierarchySettingsClient.Update` parameter(s) have been changed from `(context.Context, string, CreateOrUpdateSettingsRequest, *HierarchySettingsUpdateOptions)` to `(context.Context, string, CreateOrUpdateSettingsRequest, *HierarchySettingsClientUpdateOptions)`
- Function `*HierarchySettingsClient.Update` return value(s) have been changed from `(HierarchySettingsUpdateResponse, error)` to `(HierarchySettingsClientUpdateResponse, error)`
- Function `*HierarchySettingsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, CreateOrUpdateSettingsRequest, *HierarchySettingsCreateOrUpdateOptions)` to `(context.Context, string, CreateOrUpdateSettingsRequest, *HierarchySettingsClientCreateOrUpdateOptions)`
- Function `*HierarchySettingsClient.CreateOrUpdate` return value(s) have been changed from `(HierarchySettingsCreateOrUpdateResponse, error)` to `(HierarchySettingsClientCreateOrUpdateResponse, error)`
- Function `*EntitiesClient.List` parameter(s) have been changed from `(*EntitiesListOptions)` to `(*EntitiesClientListOptions)`
- Function `*EntitiesClient.List` return value(s) have been changed from `(*EntitiesListPager)` to `(*EntitiesClientListPager)`
- Function `*ManagementGroupSubscriptionsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *ManagementGroupSubscriptionsDeleteOptions)` to `(context.Context, string, string, *ManagementGroupSubscriptionsClientDeleteOptions)`
- Function `*ManagementGroupSubscriptionsClient.Delete` return value(s) have been changed from `(ManagementGroupSubscriptionsDeleteResponse, error)` to `(ManagementGroupSubscriptionsClientDeleteResponse, error)`
- Function `*HierarchySettingsClient.Get` parameter(s) have been changed from `(context.Context, string, *HierarchySettingsGetOptions)` to `(context.Context, string, *HierarchySettingsClientGetOptions)`
- Function `*HierarchySettingsClient.Get` return value(s) have been changed from `(HierarchySettingsGetResponse, error)` to `(HierarchySettingsClientGetResponse, error)`
- Function `*HierarchySettingsClient.List` parameter(s) have been changed from `(context.Context, string, *HierarchySettingsListOptions)` to `(context.Context, string, *HierarchySettingsClientListOptions)`
- Function `*HierarchySettingsClient.List` return value(s) have been changed from `(HierarchySettingsListResponse, error)` to `(HierarchySettingsClientListResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*ManagementGroupsGetDescendantsPager.Err` has been removed
- Function `*EntitiesListPager.NextPage` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*ManagementGroupsDeletePoller.FinalResponse` has been removed
- Function `*ManagementGroupsClient.BeginDelete` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*ManagementGroupsListPager.NextPage` has been removed
- Function `*ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupPager.NextPage` has been removed
- Function `*ManagementGroupsAPIClient.TenantBackfillStatus` has been removed
- Function `*ManagementGroupsListPager.PageResponse` has been removed
- Function `*EntitiesListPager.PageResponse` has been removed
- Function `ManagementGroupsCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*ManagementGroupsCreateOrUpdatePoller.Done` has been removed
- Function `*ManagementGroupsGetDescendantsPager.NextPage` has been removed
- Function `*ManagementGroupsClient.GetDescendants` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*ManagementGroupsCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*ManagementGroupsClient.List` has been removed
- Function `*ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupPager.Err` has been removed
- Function `*ManagementGroupsCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*ManagementGroupsDeletePoller.ResumeToken` has been removed
- Function `*ManagementGroupsCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupPager.PageResponse` has been removed
- Function `*ManagementGroupsListPager.Err` has been removed
- Function `*ManagementGroupsGetDescendantsPager.PageResponse` has been removed
- Function `*ManagementGroupsClient.Get` has been removed
- Function `*ManagementGroupsClient.BeginCreateOrUpdate` has been removed
- Function `*ManagementGroupsDeletePoller.Done` has been removed
- Function `*ManagementGroupsDeletePollerResponse.Resume` has been removed
- Function `*ManagementGroupsCreateOrUpdatePoller.Poll` has been removed
- Function `*ManagementGroupsDeletePoller.Poll` has been removed
- Function `*EntitiesListPager.Err` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*ManagementGroupsAPIClient.CheckNameAvailability` has been removed
- Function `*ManagementGroupsAPIClient.StartTenantBackfill` has been removed
- Function `ManagementGroupsDeletePollerResponse.PollUntilDone` has been removed
- Function `*ManagementGroupsClient.Update` has been removed
- Function `NewManagementGroupsClient` has been removed
- Function `NewManagementGroupsAPIClient` has been removed
- Struct `EntitiesListOptions` has been removed
- Struct `EntitiesListPager` has been removed
- Struct `EntitiesListResponse` has been removed
- Struct `EntitiesListResult` has been removed
- Struct `HierarchySettingsCreateOrUpdateOptions` has been removed
- Struct `HierarchySettingsCreateOrUpdateResponse` has been removed
- Struct `HierarchySettingsCreateOrUpdateResult` has been removed
- Struct `HierarchySettingsDeleteOptions` has been removed
- Struct `HierarchySettingsDeleteResponse` has been removed
- Struct `HierarchySettingsGetOptions` has been removed
- Struct `HierarchySettingsGetResponse` has been removed
- Struct `HierarchySettingsGetResult` has been removed
- Struct `HierarchySettingsListOptions` has been removed
- Struct `HierarchySettingsListResponse` has been removed
- Struct `HierarchySettingsListResult` has been removed
- Struct `HierarchySettingsUpdateOptions` has been removed
- Struct `HierarchySettingsUpdateResponse` has been removed
- Struct `HierarchySettingsUpdateResult` has been removed
- Struct `ManagementGroupSubscriptionsCreateOptions` has been removed
- Struct `ManagementGroupSubscriptionsCreateResponse` has been removed
- Struct `ManagementGroupSubscriptionsCreateResult` has been removed
- Struct `ManagementGroupSubscriptionsDeleteOptions` has been removed
- Struct `ManagementGroupSubscriptionsDeleteResponse` has been removed
- Struct `ManagementGroupSubscriptionsGetSubscriptionOptions` has been removed
- Struct `ManagementGroupSubscriptionsGetSubscriptionResponse` has been removed
- Struct `ManagementGroupSubscriptionsGetSubscriptionResult` has been removed
- Struct `ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupOptions` has been removed
- Struct `ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupPager` has been removed
- Struct `ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupResponse` has been removed
- Struct `ManagementGroupSubscriptionsGetSubscriptionsUnderManagementGroupResult` has been removed
- Struct `ManagementGroupsAPICheckNameAvailabilityOptions` has been removed
- Struct `ManagementGroupsAPICheckNameAvailabilityResponse` has been removed
- Struct `ManagementGroupsAPICheckNameAvailabilityResult` has been removed
- Struct `ManagementGroupsAPIClient` has been removed
- Struct `ManagementGroupsAPIStartTenantBackfillOptions` has been removed
- Struct `ManagementGroupsAPIStartTenantBackfillResponse` has been removed
- Struct `ManagementGroupsAPIStartTenantBackfillResult` has been removed
- Struct `ManagementGroupsAPITenantBackfillStatusOptions` has been removed
- Struct `ManagementGroupsAPITenantBackfillStatusResponse` has been removed
- Struct `ManagementGroupsAPITenantBackfillStatusResult` has been removed
- Struct `ManagementGroupsBeginCreateOrUpdateOptions` has been removed
- Struct `ManagementGroupsBeginDeleteOptions` has been removed
- Struct `ManagementGroupsClient` has been removed
- Struct `ManagementGroupsCreateOrUpdatePoller` has been removed
- Struct `ManagementGroupsCreateOrUpdatePollerResponse` has been removed
- Struct `ManagementGroupsCreateOrUpdateResponse` has been removed
- Struct `ManagementGroupsCreateOrUpdateResult` has been removed
- Struct `ManagementGroupsDeletePoller` has been removed
- Struct `ManagementGroupsDeletePollerResponse` has been removed
- Struct `ManagementGroupsDeleteResponse` has been removed
- Struct `ManagementGroupsDeleteResult` has been removed
- Struct `ManagementGroupsGetDescendantsOptions` has been removed
- Struct `ManagementGroupsGetDescendantsPager` has been removed
- Struct `ManagementGroupsGetDescendantsResponse` has been removed
- Struct `ManagementGroupsGetDescendantsResult` has been removed
- Struct `ManagementGroupsGetOptions` has been removed
- Struct `ManagementGroupsGetResponse` has been removed
- Struct `ManagementGroupsGetResult` has been removed
- Struct `ManagementGroupsListOptions` has been removed
- Struct `ManagementGroupsListPager` has been removed
- Struct `ManagementGroupsListResponse` has been removed
- Struct `ManagementGroupsListResult` has been removed
- Struct `ManagementGroupsUpdateOptions` has been removed
- Struct `ManagementGroupsUpdateResponse` has been removed
- Struct `ManagementGroupsUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*ClientListPager.NextPage(context.Context) bool`
- New function `*ClientDeletePoller.Done() bool`
- New function `*ClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ClientGetDescendantsPager.PageResponse() ClientGetDescendantsResponse`
- New function `*ClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager.NextPage(context.Context) bool`
- New function `*ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager.Err() error`
- New function `*ClientDeletePollerResponse.Resume(context.Context, *Client, string) error`
- New function `NewAPIClient(azcore.TokenCredential, *arm.ClientOptions) *APIClient`
- New function `*EntitiesClientListPager.Err() error`
- New function `*Client.Get(context.Context, string, *ClientGetOptions) (ClientGetResponse, error)`
- New function `*ClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*EntitiesClientListPager.NextPage(context.Context) bool`
- New function `*APIClient.CheckNameAvailability(context.Context, CheckNameAvailabilityRequest, *APIClientCheckNameAvailabilityOptions) (APIClientCheckNameAvailabilityResponse, error)`
- New function `*Client.BeginDelete(context.Context, string, *ClientBeginDeleteOptions) (ClientDeletePollerResponse, error)`
- New function `*APIClient.StartTenantBackfill(context.Context, *APIClientStartTenantBackfillOptions) (APIClientStartTenantBackfillResponse, error)`
- New function `ClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (ClientDeleteResponse, error)`
- New function `*ClientListPager.PageResponse() ClientListResponse`
- New function `*APIClient.TenantBackfillStatus(context.Context, *APIClientTenantBackfillStatusOptions) (APIClientTenantBackfillStatusResponse, error)`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*Client.BeginCreateOrUpdate(context.Context, string, CreateManagementGroupRequest, *ClientBeginCreateOrUpdateOptions) (ClientCreateOrUpdatePollerResponse, error)`
- New function `*Client.Update(context.Context, string, PatchManagementGroupRequest, *ClientUpdateOptions) (ClientUpdateResponse, error)`
- New function `*ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager.PageResponse() ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResponse`
- New function `NewClient(azcore.TokenCredential, *arm.ClientOptions) *Client`
- New function `*ClientListPager.Err() error`
- New function `*ClientDeletePoller.FinalResponse(context.Context) (ClientDeleteResponse, error)`
- New function `*ClientGetDescendantsPager.Err() error`
- New function `ClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (ClientCreateOrUpdateResponse, error)`
- New function `*EntitiesClientListPager.PageResponse() EntitiesClientListResponse`
- New function `*ClientCreateOrUpdatePollerResponse.Resume(context.Context, *Client, string) error`
- New function `*Client.GetDescendants(string, *ClientGetDescendantsOptions) *ClientGetDescendantsPager`
- New function `*ClientCreateOrUpdatePoller.Done() bool`
- New function `*ClientCreateOrUpdatePoller.FinalResponse(context.Context) (ClientCreateOrUpdateResponse, error)`
- New function `*ClientDeletePoller.ResumeToken() (string, error)`
- New function `*ClientGetDescendantsPager.NextPage(context.Context) bool`
- New function `*Client.List(*ClientListOptions) *ClientListPager`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.Err() error`
- New struct `APIClient`
- New struct `APIClientCheckNameAvailabilityOptions`
- New struct `APIClientCheckNameAvailabilityResponse`
- New struct `APIClientCheckNameAvailabilityResult`
- New struct `APIClientStartTenantBackfillOptions`
- New struct `APIClientStartTenantBackfillResponse`
- New struct `APIClientStartTenantBackfillResult`
- New struct `APIClientTenantBackfillStatusOptions`
- New struct `APIClientTenantBackfillStatusResponse`
- New struct `APIClientTenantBackfillStatusResult`
- New struct `Client`
- New struct `ClientBeginCreateOrUpdateOptions`
- New struct `ClientBeginDeleteOptions`
- New struct `ClientCreateOrUpdatePoller`
- New struct `ClientCreateOrUpdatePollerResponse`
- New struct `ClientCreateOrUpdateResponse`
- New struct `ClientCreateOrUpdateResult`
- New struct `ClientDeletePoller`
- New struct `ClientDeletePollerResponse`
- New struct `ClientDeleteResponse`
- New struct `ClientDeleteResult`
- New struct `ClientGetDescendantsOptions`
- New struct `ClientGetDescendantsPager`
- New struct `ClientGetDescendantsResponse`
- New struct `ClientGetDescendantsResult`
- New struct `ClientGetOptions`
- New struct `ClientGetResponse`
- New struct `ClientGetResult`
- New struct `ClientListOptions`
- New struct `ClientListPager`
- New struct `ClientListResponse`
- New struct `ClientListResult`
- New struct `ClientUpdateOptions`
- New struct `ClientUpdateResponse`
- New struct `ClientUpdateResult`
- New struct `EntitiesClientListOptions`
- New struct `EntitiesClientListPager`
- New struct `EntitiesClientListResponse`
- New struct `EntitiesClientListResult`
- New struct `HierarchySettingsClientCreateOrUpdateOptions`
- New struct `HierarchySettingsClientCreateOrUpdateResponse`
- New struct `HierarchySettingsClientCreateOrUpdateResult`
- New struct `HierarchySettingsClientDeleteOptions`
- New struct `HierarchySettingsClientDeleteResponse`
- New struct `HierarchySettingsClientGetOptions`
- New struct `HierarchySettingsClientGetResponse`
- New struct `HierarchySettingsClientGetResult`
- New struct `HierarchySettingsClientListOptions`
- New struct `HierarchySettingsClientListResponse`
- New struct `HierarchySettingsClientListResult`
- New struct `HierarchySettingsClientUpdateOptions`
- New struct `HierarchySettingsClientUpdateResponse`
- New struct `HierarchySettingsClientUpdateResult`
- New struct `ManagementGroupSubscriptionsClientCreateOptions`
- New struct `ManagementGroupSubscriptionsClientCreateResponse`
- New struct `ManagementGroupSubscriptionsClientCreateResult`
- New struct `ManagementGroupSubscriptionsClientDeleteOptions`
- New struct `ManagementGroupSubscriptionsClientDeleteResponse`
- New struct `ManagementGroupSubscriptionsClientGetSubscriptionOptions`
- New struct `ManagementGroupSubscriptionsClientGetSubscriptionResponse`
- New struct `ManagementGroupSubscriptionsClientGetSubscriptionResult`
- New struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupOptions`
- New struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupPager`
- New struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResponse`
- New struct `ManagementGroupSubscriptionsClientGetSubscriptionsUnderManagementGroupResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Error` in struct `ErrorResponse`


## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-10-26)

- Initial preview release.
