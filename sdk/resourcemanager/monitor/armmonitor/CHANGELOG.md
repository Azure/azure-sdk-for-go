# Release History

## 0.8.0 (2022-10-18)
### Breaking Changes

- Type alias `AlertSeverity` type has been changed from `string` to `int64`
- Function `*ScheduledQueryRulesClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, LogSearchRuleResource, *ScheduledQueryRulesClientCreateOrUpdateOptions)` to `(context.Context, string, string, ScheduledQueryRuleResource, *ScheduledQueryRulesClientCreateOrUpdateOptions)`
- Function `*ScheduledQueryRulesClient.Update` parameter(s) have been changed from `(context.Context, string, string, LogSearchRuleResourcePatch, *ScheduledQueryRulesClientUpdateOptions)` to `(context.Context, string, string, ScheduledQueryRuleResourcePatch, *ScheduledQueryRulesClientUpdateOptions)`
- Type of `OperationStatus.Error` has been changed from `*ErrorResponseCommon` to `*ErrorDetail`
- Type of `PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState` has been changed from `*PrivateLinkServiceConnectionStateProperty` to `*PrivateLinkServiceConnectionState`
- Type of `PrivateEndpointConnectionProperties.PrivateEndpoint` has been changed from `*PrivateEndpointProperty` to `*PrivateEndpoint`
- Type of `PrivateEndpointConnectionProperties.ProvisioningState` has been changed from `*string` to `*PrivateEndpointConnectionProvisioningState`
- Type of `ErrorContract.Error` has been changed from `*ErrorResponse` to `*ErrorResponseDetails`
- Type of `Dimension.Operator` has been changed from `*Operator` to `*DimensionOperator`
- Type alias `ConditionalOperator`, const `ConditionalOperatorLessThanOrEqual`, `ConditionalOperatorEqual`, `ConditionalOperatorGreaterThanOrEqual`, `ConditionalOperatorLessThan`, `ConditionalOperatorGreaterThan` and function `PossibleConditionalOperatorValues` have been removed
- Type alias `Enabled`, const `EnabledTrue`, `EnabledFalse` and function `PossibleEnabledValues` have been removed
- Type alias `QueryType`, const `QueryTypeResultCount` and function `PossibleQueryTypeValues` have been removed
- Type alias `MetricTriggerType`, const `MetricTriggerTypeConsecutive`, `MetricTriggerTypeTotal` and function `PossibleMetricTriggerTypeValues` have been removed
- Type alias `ProvisioningState`, const `ProvisioningStateSucceeded`, `ProvisioningStateFailed`, `ProvisioningStateDeploying`, `ProvisioningStateCanceled` and function `PossibleProvisioningStateValues` have been removed
- Const `OperatorInclude` has been removed
- Function `*DiagnosticSettingsClient.List` has been changed to `*DiagnosticSettingsClient.NewListPager(string, *DiagnosticSettingsClientListOptions) *runtime.Pager[DiagnosticSettingsClientListResponse]`
- Function `*PrivateEndpointConnectionsClient.NewListByPrivateLinkScopePager` has been changed to `*PrivateEndpointConnectionsClient.ListByPrivateLinkScope(context.Context, string, string, *PrivateEndpointConnectionsClientListByPrivateLinkScopeOptions) (PrivateEndpointConnectionsClientListByPrivateLinkScopeResponse, error)`
- Function `*DiagnosticSettingsCategoryClient.List` has been changed to `*DiagnosticSettingsCategoryClient.NewListPager(string, *DiagnosticSettingsCategoryClientListOptions) *runtime.Pager[DiagnosticSettingsCategoryClientListResponse]`
- Function `*PrivateLinkResourcesClient.NewListByPrivateLinkScopePager` has been changed to `*PrivateLinkResourcesClient.ListByPrivateLinkScope(context.Context, string, string, *PrivateLinkResourcesClientListByPrivateLinkScopeOptions) (PrivateLinkResourcesClientListByPrivateLinkScopeResponse, error)`
- Struct `Action` and function `*Action.GetAction` have been removed
- Struct `AlertingAction` and function `*AlertingAction.GetAction` have been removed
- Struct `LogToMetricAction` and function `*LogToMetricAction.GetAction` have been removed
- Struct `AzNsActionGroup` has been removed
- Struct `Criteria` has been removed
- Struct `LogMetricTrigger` has been removed
- Struct `LogSearchRule` has been removed
- Struct `LogSearchRulePatch` has been removed
- Struct `LogSearchRuleResource` has been removed
- Struct `LogSearchRuleResourceCollection` has been removed
- Struct `LogSearchRuleResourcePatch` has been removed
- Struct `PrivateEndpointProperty` has been removed
- Struct `PrivateLinkScopesResource` has been removed
- Struct `PrivateLinkServiceConnectionStateProperty` has been removed
- Struct `Schedule` has been removed
- Struct `Source` has been removed
- Struct `TriggerCondition` has been removed
- Field `Filter` of struct `ScheduledQueryRulesClientListByResourceGroupOptions` has been removed
- Field `LogSearchRuleResource` of struct `ScheduledQueryRulesClientUpdateResponse` has been removed
- Field `LogSearchRuleResourceCollection` of struct `ScheduledQueryRulesClientListBySubscriptionResponse` has been removed
- Field `LogSearchRuleResourceCollection` of struct `ScheduledQueryRulesClientListByResourceGroupResponse` has been removed
- Field `NextLink` of struct `PrivateLinkResourceListResult` has been removed
- Field `Filter` of struct `ScheduledQueryRulesClientListBySubscriptionOptions` has been removed
- Field `LogSearchRuleResource` of struct `ScheduledQueryRulesClientGetResponse` has been removed
- Field `NextLink` of struct `PrivateEndpointConnectionListResult` has been removed
- Field `Identity` of struct `ActionGroupResource` has been removed
- Field `Kind` of struct `ActionGroupResource` has been removed
- Field `Identity` of struct `AzureResource` has been removed
- Field `Kind` of struct `AzureResource` has been removed
- Field `LogSearchRuleResource` of struct `ScheduledQueryRulesClientCreateOrUpdateResponse` has been removed

### Features Added

- New const `ConditionOperatorEquals`
- New const `KindLogAlert`
- New const `PrivateEndpointServiceConnectionStatusPending`
- New const `AccessModePrivateOnly`
- New const `PredictiveAutoscalePolicyScaleModeEnabled`
- New const `TimeAggregationTotal`
- New const `TimeAggregationAverage`
- New const `AccessModeOpen`
- New const `PredictiveAutoscalePolicyScaleModeForecastOnly`
- New const `PrivateEndpointConnectionProvisioningStateFailed`
- New const `PredictiveAutoscalePolicyScaleModeDisabled`
- New const `PrivateEndpointConnectionProvisioningStateCreating`
- New const `TimeAggregationMinimum`
- New const `DimensionOperatorExclude`
- New const `DimensionOperatorInclude`
- New const `PrivateEndpointConnectionProvisioningStateDeleting`
- New const `PrivateEndpointServiceConnectionStatusRejected`
- New const `PrivateEndpointConnectionProvisioningStateSucceeded`
- New const `PrivateEndpointServiceConnectionStatusApproved`
- New const `KindLogToMetric`
- New const `TimeAggregationCount`
- New const `TimeAggregationMaximum`
- New type alias `PrivateEndpointServiceConnectionStatus`
- New type alias `DimensionOperator`
- New type alias `PredictiveAutoscalePolicyScaleMode`
- New type alias `AccessMode`
- New type alias `Kind`
- New type alias `TimeAggregation`
- New type alias `PrivateEndpointConnectionProvisioningState`
- New function `PossiblePrivateEndpointServiceConnectionStatusValues() []PrivateEndpointServiceConnectionStatus`
- New function `PossiblePredictiveAutoscalePolicyScaleModeValues() []PredictiveAutoscalePolicyScaleMode`
- New function `PossiblePrivateEndpointConnectionProvisioningStateValues() []PrivateEndpointConnectionProvisioningState`
- New function `NewPredictiveMetricClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PredictiveMetricClient, error)`
- New function `PossibleTimeAggregationValues() []TimeAggregation`
- New function `PossibleAccessModeValues() []AccessMode`
- New function `*PredictiveMetricClient.Get(context.Context, string, string, string, string, string, string, string, *PredictiveMetricClientGetOptions) (PredictiveMetricClientGetResponse, error)`
- New function `PossibleKindValues() []Kind`
- New function `PossibleDimensionOperatorValues() []DimensionOperator`
- New struct `AccessModeSettings`
- New struct `AccessModeSettingsExclusion`
- New struct `Actions`
- New struct `AutoscaleErrorResponse`
- New struct `AutoscaleErrorResponseError`
- New struct `Condition`
- New struct `ConditionFailingPeriods`
- New struct `DefaultErrorResponse`
- New struct `ErrorResponseAdditionalInfo`
- New struct `ErrorResponseDetails`
- New struct `PredictiveAutoscalePolicy`
- New struct `PredictiveMetricClient`
- New struct `PredictiveMetricClientGetOptions`
- New struct `PredictiveMetricClientGetResponse`
- New struct `PredictiveResponse`
- New struct `PredictiveValue`
- New struct `PrivateEndpoint`
- New struct `PrivateLinkServiceConnectionState`
- New struct `ProxyResourceAutoGenerated`
- New struct `ResourceAutoGenerated`
- New struct `ResourceAutoGenerated2`
- New struct `ResourceAutoGenerated3`
- New struct `ResourceAutoGenerated4`
- New struct `ScheduledQueryRuleCriteria`
- New struct `ScheduledQueryRuleProperties`
- New struct `ScheduledQueryRuleResource`
- New struct `ScheduledQueryRuleResourceCollection`
- New struct `ScheduledQueryRuleResourcePatch`
- New struct `TrackedResource`
- New anonymous field `TestNotificationDetailsResponse` in struct `ActionGroupsClientCreateNotificationsAtActionGroupResourceLevelResponse`
- New field `SystemData` in struct `DiagnosticSettingsCategoryResource`
- New anonymous field `ScheduledQueryRuleResource` in struct `ScheduledQueryRulesClientCreateOrUpdateResponse`
- New anonymous field `TestNotificationDetailsResponse` in struct `ActionGroupsClientCreateNotificationsAtResourceGroupLevelResponse`
- New field `PredictiveAutoscalePolicy` in struct `AutoscaleSetting`
- New field `SystemData` in struct `Resource`
- New field `CategoryGroups` in struct `DiagnosticSettingsCategory`
- New field `SystemData` in struct `AzureMonitorPrivateLinkScope`
- New field `RequiredZoneNames` in struct `PrivateLinkResourceProperties`
- New anonymous field `ScheduledQueryRuleResource` in struct `ScheduledQueryRulesClientGetResponse`
- New field `MarketplacePartnerID` in struct `DiagnosticSettings`
- New anonymous field `ScheduledQueryRuleResourceCollection` in struct `ScheduledQueryRulesClientListByResourceGroupResponse`
- New field `SystemData` in struct `ScopedResource`
- New field `SystemData` in struct `DiagnosticSettingsResource`
- New field `CategoryGroup` in struct `LogSettings`
- New field `AccessModeSettings` in struct `AzureMonitorPrivateLinkScopeProperties`
- New anonymous field `ScheduledQueryRuleResource` in struct `ScheduledQueryRulesClientUpdateResponse`
- New field `SystemData` in struct `AutoscaleSettingResource`
- New anonymous field `TestNotificationDetailsResponse` in struct `ActionGroupsClientPostTestNotificationsResponse`
- New anonymous field `ScheduledQueryRuleResourceCollection` in struct `ScheduledQueryRulesClientListBySubscriptionResponse`


## 0.7.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.7.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
