# Release History

## 0.11.0 (2026-01-29)
### Breaking Changes

- Function `NewAlertsClient` parameter(s) have been changed from `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)` to `(scope string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `NewClientFactory` parameter(s) have been changed from `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)` to `(scope string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Type of `Operation.Origin` has been changed from `*string` to `*Origin`
- `ActionTypeAddActionGroups`, `ActionTypeRemoveAllActionGroups` from enum `ActionType` has been removed
- `AlertModificationEventActionRuleSuppressed`, `AlertModificationEventActionRuleTriggered`, `AlertModificationEventActionsFailed` from enum `AlertModificationEvent` has been removed
- `TimeRangeOneD`, `TimeRangeOneH`, `TimeRangeSevenD`, `TimeRangeThirtyD` from enum `TimeRange` has been removed
- Enum `DaysOfWeek` has been removed
- Enum `Field` has been removed
- Enum `Operator` has been removed
- Enum `RecurrenceType` has been removed
- Enum `SmartGroupModificationEvent` has been removed
- Enum `SmartGroupsSortByFields` has been removed
- Enum `State` has been removed
- Function `*Action.GetAction` has been removed
- Function `*AddActionGroups.GetAction` has been removed
- Function `NewAlertProcessingRulesClient` has been removed
- Function `*AlertProcessingRulesClient.CreateOrUpdate` has been removed
- Function `*AlertProcessingRulesClient.Delete` has been removed
- Function `*AlertProcessingRulesClient.GetByName` has been removed
- Function `*AlertProcessingRulesClient.NewListByResourceGroupPager` has been removed
- Function `*AlertProcessingRulesClient.NewListBySubscriptionPager` has been removed
- Function `*AlertProcessingRulesClient.Update` has been removed
- Function `NewAlertRuleRecommendationsClient` has been removed
- Function `*AlertRuleRecommendationsClient.NewListByResourcePager` has been removed
- Function `*AlertRuleRecommendationsClient.NewListByTargetTypePager` has been removed
- Function `*ClientFactory.NewAlertProcessingRulesClient` has been removed
- Function `*ClientFactory.NewAlertRuleRecommendationsClient` has been removed
- Function `*ClientFactory.NewPrometheusRuleGroupsClient` has been removed
- Function `*ClientFactory.NewSmartGroupsClient` has been removed
- Function `*ClientFactory.NewTenantActivityLogAlertsClient` has been removed
- Function `*DailyRecurrence.GetRecurrence` has been removed
- Function `*MonthlyRecurrence.GetRecurrence` has been removed
- Function `NewPrometheusRuleGroupsClient` has been removed
- Function `*PrometheusRuleGroupsClient.CreateOrUpdate` has been removed
- Function `*PrometheusRuleGroupsClient.Delete` has been removed
- Function `*PrometheusRuleGroupsClient.Get` has been removed
- Function `*PrometheusRuleGroupsClient.NewListByResourceGroupPager` has been removed
- Function `*PrometheusRuleGroupsClient.NewListBySubscriptionPager` has been removed
- Function `*PrometheusRuleGroupsClient.Update` has been removed
- Function `*Recurrence.GetRecurrence` has been removed
- Function `*RemoveAllActionGroups.GetAction` has been removed
- Function `NewSmartGroupsClient` has been removed
- Function `*SmartGroupsClient.ChangeState` has been removed
- Function `*SmartGroupsClient.NewGetAllPager` has been removed
- Function `*SmartGroupsClient.GetByID` has been removed
- Function `*SmartGroupsClient.GetHistory` has been removed
- Function `NewTenantActivityLogAlertsClient` has been removed
- Function `*TenantActivityLogAlertsClient.CreateOrUpdate` has been removed
- Function `*TenantActivityLogAlertsClient.Delete` has been removed
- Function `*TenantActivityLogAlertsClient.Get` has been removed
- Function `*TenantActivityLogAlertsClient.NewListByManagementGroupPager` has been removed
- Function `*TenantActivityLogAlertsClient.NewListByTenantPager` has been removed
- Function `*TenantActivityLogAlertsClient.Update` has been removed
- Function `*WeeklyRecurrence.GetRecurrence` has been removed
- Struct `ActionGroup` has been removed
- Struct `ActionList` has been removed
- Struct `AddActionGroups` has been removed
- Struct `AlertProcessingRule` has been removed
- Struct `AlertProcessingRuleProperties` has been removed
- Struct `AlertProcessingRulesList` has been removed
- Struct `AlertRuleAllOfCondition` has been removed
- Struct `AlertRuleAnyOfOrLeafCondition` has been removed
- Struct `AlertRuleLeafCondition` has been removed
- Struct `AlertRuleProperties` has been removed
- Struct `AlertRuleRecommendationProperties` has been removed
- Struct `AlertRuleRecommendationResource` has been removed
- Struct `AlertRuleRecommendationsListResponse` has been removed
- Struct `Condition` has been removed
- Struct `DailyRecurrence` has been removed
- Struct `MonthlyRecurrence` has been removed
- Struct `OperationsList` has been removed
- Struct `PatchObject` has been removed
- Struct `PatchProperties` has been removed
- Struct `PrometheusRule` has been removed
- Struct `PrometheusRuleGroupAction` has been removed
- Struct `PrometheusRuleGroupProperties` has been removed
- Struct `PrometheusRuleGroupResource` has been removed
- Struct `PrometheusRuleGroupResourceCollection` has been removed
- Struct `PrometheusRuleGroupResourcePatch` has been removed
- Struct `PrometheusRuleGroupResourcePatchProperties` has been removed
- Struct `PrometheusRuleResolveConfiguration` has been removed
- Struct `RemoveAllActionGroups` has been removed
- Struct `RuleArmTemplate` has been removed
- Struct `Schedule` has been removed
- Struct `SmartGroup` has been removed
- Struct `SmartGroupAggregatedProperty` has been removed
- Struct `SmartGroupModification` has been removed
- Struct `SmartGroupModificationItem` has been removed
- Struct `SmartGroupModificationProperties` has been removed
- Struct `SmartGroupProperties` has been removed
- Struct `SmartGroupsList` has been removed
- Struct `TenantActivityLogAlertResource` has been removed
- Struct `TenantAlertRuleList` has been removed
- Struct `TenantAlertRulePatchObject` has been removed
- Struct `TenantAlertRulePatchProperties` has been removed
- Struct `WeeklyRecurrence` has been removed
- Field `OperationsList` of struct `OperationsClientListResponse` has been removed

### Features Added

- New value `ActionTypeInternal` added to enum type `ActionType`
- New value `MonitorServiceResourceHealth` added to enum type `MonitorService`
- New value `TimeRange1D`, `TimeRange1H`, `TimeRange30D`, `TimeRange7D` added to enum type `TimeRange`
- New enum type `AlertModificationType` with values `AlertModificationTypeActionsSuppressed`, `AlertModificationTypeActionsTriggered`, `AlertModificationTypePropertyChange`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ResultStatus` with values `ResultStatusFailed`, `ResultStatusInline`, `ResultStatusNone`, `ResultStatusThrottled`, `ResultStatusThrottledByAlertRule`, `ResultStatusThrottledBySubscription`
- New enum type `RuleType` with values `RuleTypeActionRule`, `RuleTypeAlertRule`
- New enum type `Status` with values `StatusFailed`, `StatusSucceeded`
- New enum type `Type` with values `TypePrometheusInstantQuery`, `TypePrometheusRangeQuery`
- New function `*ActionSuppressedDetails.GetBaseDetails() *BaseDetails`
- New function `*ActionTriggeredDetails.GetBaseDetails() *BaseDetails`
- New function `*AlertEnrichmentItem.GetAlertEnrichmentItem() *AlertEnrichmentItem`
- New function `*AlertsClient.ChangeStateTenant(ctx context.Context, alertID string, newState AlertState, options *AlertsClientChangeStateTenantOptions) (AlertsClientChangeStateTenantResponse, error)`
- New function `*AlertsClient.NewGetAllTenantPager(options *AlertsClientGetAllTenantOptions) *runtime.Pager[AlertsClientGetAllTenantResponse]`
- New function `*AlertsClient.GetByIDTenant(ctx context.Context, alertID string, options *AlertsClientGetByIDTenantOptions) (AlertsClientGetByIDTenantResponse, error)`
- New function `*AlertsClient.NewGetEnrichmentsPager(alertID string, options *AlertsClientGetEnrichmentsOptions) *runtime.Pager[AlertsClientGetEnrichmentsResponse]`
- New function `*AlertsClient.GetHistoryTenant(ctx context.Context, alertID string, options *AlertsClientGetHistoryTenantOptions) (AlertsClientGetHistoryTenantResponse, error)`
- New function `*BaseDetails.GetBaseDetails() *BaseDetails`
- New function `*PrometheusInstantQuery.GetAlertEnrichmentItem() *AlertEnrichmentItem`
- New function `*PrometheusRangeQuery.GetAlertEnrichmentItem() *AlertEnrichmentItem`
- New function `*PropertyChangeDetails.GetBaseDetails() *BaseDetails`
- New struct `ActionSuppressedDetails`
- New struct `ActionTriggeredDetails`
- New struct `AlertEnrichmentProperties`
- New struct `AlertEnrichmentResponse`
- New struct `AlertEnrichmentsList`
- New struct `NotificationResult`
- New struct `OperationListResult`
- New struct `PrometheusInstantQuery`
- New struct `PrometheusRangeQuery`
- New struct `PropertyChangeDetails`
- New struct `TriggeredRule`
- New field `SystemData` in struct `Alert`
- New field `Details` in struct `AlertModificationItem`
- New field `CustomProperties` in struct `AlertProperties`
- New field `ActionType`, `IsDataAction` in struct `Operation`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`


## 0.10.0 (2024-03-01)
### Breaking Changes

- Type of `AlertsClientChangeStateOptions.Comment` has been changed from `*string` to `*Comments`

### Features Added

- New function `NewAlertRuleRecommendationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AlertRuleRecommendationsClient, error)`
- New function `*AlertRuleRecommendationsClient.NewListByResourcePager(string, *AlertRuleRecommendationsClientListByResourceOptions) *runtime.Pager[AlertRuleRecommendationsClientListByResourceResponse]`
- New function `*AlertRuleRecommendationsClient.NewListByTargetTypePager(string, *AlertRuleRecommendationsClientListByTargetTypeOptions) *runtime.Pager[AlertRuleRecommendationsClientListByTargetTypeResponse]`
- New function `*ClientFactory.NewAlertRuleRecommendationsClient() *AlertRuleRecommendationsClient`
- New function `*ClientFactory.NewPrometheusRuleGroupsClient() *PrometheusRuleGroupsClient`
- New function `*ClientFactory.NewTenantActivityLogAlertsClient() *TenantActivityLogAlertsClient`
- New function `NewPrometheusRuleGroupsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrometheusRuleGroupsClient, error)`
- New function `*PrometheusRuleGroupsClient.CreateOrUpdate(context.Context, string, string, PrometheusRuleGroupResource, *PrometheusRuleGroupsClientCreateOrUpdateOptions) (PrometheusRuleGroupsClientCreateOrUpdateResponse, error)`
- New function `*PrometheusRuleGroupsClient.Delete(context.Context, string, string, *PrometheusRuleGroupsClientDeleteOptions) (PrometheusRuleGroupsClientDeleteResponse, error)`
- New function `*PrometheusRuleGroupsClient.Get(context.Context, string, string, *PrometheusRuleGroupsClientGetOptions) (PrometheusRuleGroupsClientGetResponse, error)`
- New function `*PrometheusRuleGroupsClient.NewListByResourceGroupPager(string, *PrometheusRuleGroupsClientListByResourceGroupOptions) *runtime.Pager[PrometheusRuleGroupsClientListByResourceGroupResponse]`
- New function `*PrometheusRuleGroupsClient.NewListBySubscriptionPager(*PrometheusRuleGroupsClientListBySubscriptionOptions) *runtime.Pager[PrometheusRuleGroupsClientListBySubscriptionResponse]`
- New function `*PrometheusRuleGroupsClient.Update(context.Context, string, string, PrometheusRuleGroupResourcePatch, *PrometheusRuleGroupsClientUpdateOptions) (PrometheusRuleGroupsClientUpdateResponse, error)`
- New function `NewTenantActivityLogAlertsClient(azcore.TokenCredential, *arm.ClientOptions) (*TenantActivityLogAlertsClient, error)`
- New function `*TenantActivityLogAlertsClient.CreateOrUpdate(context.Context, string, string, TenantActivityLogAlertResource, *TenantActivityLogAlertsClientCreateOrUpdateOptions) (TenantActivityLogAlertsClientCreateOrUpdateResponse, error)`
- New function `*TenantActivityLogAlertsClient.Delete(context.Context, string, string, *TenantActivityLogAlertsClientDeleteOptions) (TenantActivityLogAlertsClientDeleteResponse, error)`
- New function `*TenantActivityLogAlertsClient.Get(context.Context, string, string, *TenantActivityLogAlertsClientGetOptions) (TenantActivityLogAlertsClientGetResponse, error)`
- New function `*TenantActivityLogAlertsClient.NewListByManagementGroupPager(string, *TenantActivityLogAlertsClientListByManagementGroupOptions) *runtime.Pager[TenantActivityLogAlertsClientListByManagementGroupResponse]`
- New function `*TenantActivityLogAlertsClient.NewListByTenantPager(*TenantActivityLogAlertsClientListByTenantOptions) *runtime.Pager[TenantActivityLogAlertsClientListByTenantResponse]`
- New function `*TenantActivityLogAlertsClient.Update(context.Context, string, string, TenantAlertRulePatchObject, *TenantActivityLogAlertsClientUpdateOptions) (TenantActivityLogAlertsClientUpdateResponse, error)`
- New struct `ActionGroup`
- New struct `ActionList`
- New struct `AlertRuleAllOfCondition`
- New struct `AlertRuleAnyOfOrLeafCondition`
- New struct `AlertRuleLeafCondition`
- New struct `AlertRuleProperties`
- New struct `AlertRuleRecommendationProperties`
- New struct `AlertRuleRecommendationResource`
- New struct `AlertRuleRecommendationsListResponse`
- New struct `Comments`
- New struct `PrometheusRule`
- New struct `PrometheusRuleGroupAction`
- New struct `PrometheusRuleGroupProperties`
- New struct `PrometheusRuleGroupResource`
- New struct `PrometheusRuleGroupResourceCollection`
- New struct `PrometheusRuleGroupResourcePatch`
- New struct `PrometheusRuleGroupResourcePatchProperties`
- New struct `PrometheusRuleResolveConfiguration`
- New struct `RuleArmTemplate`
- New struct `TenantActivityLogAlertResource`
- New struct `TenantAlertRuleList`
- New struct `TenantAlertRulePatchObject`
- New struct `TenantAlertRulePatchProperties`


## 0.9.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.8.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 0.8.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.7.0 (2022-08-19)
### Features Added

- New field `Origin` in struct `Operation`
- New field `Comment` in struct `AlertsClientChangeStateOptions`


## 0.6.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/alertsmanagement/armalertsmanagement` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).