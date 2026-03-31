# Release History

## 0.11.0 (2026-03-13)
### Breaking Changes

- The following resource clients and their associated types/enums have been split into dedicated packages. Please migrate to the new packages:
  - `AlertProcessingRulesClient` → [`armalertprocessingrules.Client`](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/alertprocessingrules/armalertprocessingrules)
  - `TenantActivityLogAlertsClient` → [`armtenantactivitylogalerts.Client`](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/tenantactivitylogalerts/armtenantactivitylogalerts)
  - `AlertRuleRecommendationsClient` → [`armalertrulerecommendations.Client`](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/alertrulerecommendations/armalertrulerecommendations)
  - `PrometheusRuleGroupsClient` → [`armprometheusrulegroups.Client`](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/prometheusrulegroups/armprometheusrulegroups)
- `SmartGroupsClient` and all related types/enums have been removed
- `AlertModificationEventActionRuleSuppressed`, `AlertModificationEventActionRuleTriggered`, `AlertModificationEventActionsFailed` from enum `AlertModificationEvent` have been removed
- Function `NewAlertsClient` parameter(s) have been changed from `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)` to `(scope string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `NewClientFactory` parameter(s) have been changed from `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)` to `(scope string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Type of `Operation.Origin` has been changed from `*string` to `*Origin`
- Struct `OperationsList` has been renamed to `OperationListResult`; field `OperationsList` of struct `OperationsClientListResponse` has been replaced by `OperationListResult`

### Features Added

- New value `ActionTypeInternal` added to enum type `ActionType`
- New value `MonitorServiceResourceHealth` added to enum type `MonitorService`
- New enum type `AlertModificationType` with values `AlertModificationTypeActionsSuppressed`, `AlertModificationTypeActionsTriggered`, `AlertModificationTypePropertyChange`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ResultStatus` with values `ResultStatusFailed`, `ResultStatusInline`, `ResultStatusNone`, `ResultStatusThrottled`, `ResultStatusThrottledByAlertRule`, `ResultStatusThrottledBySubscription`
- New enum type `RuleType` with values `RuleTypeActionRule`, `RuleTypeAlertRule`
- New enum type `Status` with values `StatusFailed`, `StatusSucceeded`
- New enum type `Type` with values `TypePrometheusEnrichmentItem`, `TypePrometheusInstantQuery`, `TypePrometheusRangeQuery`
- New function `*ActionSuppressedDetails.GetBaseDetails() *BaseDetails`
- New function `*ActionTriggeredDetails.GetBaseDetails() *BaseDetails`
- New function `*AlertEnrichmentItem.GetAlertEnrichmentItem() *AlertEnrichmentItem`
- New function `*AlertsClient.ChangeStateTenant(ctx context.Context, alertID string, newState AlertState, options *AlertsClientChangeStateTenantOptions) (AlertsClientChangeStateTenantResponse, error)`
- New function `*AlertsClient.NewGetAllTenantPager(options *AlertsClientGetAllTenantOptions) *runtime.Pager[AlertsClientGetAllTenantResponse]`
- New function `*AlertsClient.GetByIDTenant(ctx context.Context, alertID string, options *AlertsClientGetByIDTenantOptions) (AlertsClientGetByIDTenantResponse, error)`
- New function `*AlertsClient.NewGetEnrichmentsPager(alertID string, options *AlertsClientGetEnrichmentsOptions) *runtime.Pager[AlertsClientGetEnrichmentsResponse]`
- New function `*AlertsClient.GetHistoryTenant(ctx context.Context, alertID string, options *AlertsClientGetHistoryTenantOptions) (AlertsClientGetHistoryTenantResponse, error)`
- New function `*BaseDetails.GetBaseDetails() *BaseDetails`
- New function `*PrometheusEnrichmentItem.GetAlertEnrichmentItem() *AlertEnrichmentItem`
- New function `*PrometheusEnrichmentItem.GetPrometheusEnrichmentItem() *PrometheusEnrichmentItem`
- New function `*PrometheusInstantQuery.GetAlertEnrichmentItem() *AlertEnrichmentItem`
- New function `*PrometheusInstantQuery.GetPrometheusEnrichmentItem() *PrometheusEnrichmentItem`
- New function `*PrometheusRangeQuery.GetAlertEnrichmentItem() *AlertEnrichmentItem`
- New function `*PrometheusRangeQuery.GetPrometheusEnrichmentItem() *PrometheusEnrichmentItem`
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