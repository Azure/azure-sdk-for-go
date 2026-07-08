# Release History

## 0.3.0 (2026-07-07)
### Breaking Changes

- Field `Dimension` of struct `AzureResourceSignal` has been removed
- Field `Dimension` of struct `ResourceMetricSignalDefinitionProperties` has been removed

### Features Added

- New value `RefreshIntervalPT15M` added to enum type `RefreshInterval`
- New value `SignalOperatorDynamic` added to enum type `SignalOperator`
- New enum type `DynamicThresholdSensitivity` with values `DynamicThresholdSensitivityHigh`, `DynamicThresholdSensitivityLow`, `DynamicThresholdSensitivityMedium`
- New enum type `LookBackWindow` with values `LookBackWindowPT15M`, `LookBackWindowPT1H`, `LookBackWindowPT30M`, `LookBackWindowPT5M`
- New enum type `ResourceHealthAvailabilityState` with values `ResourceHealthAvailabilityStateAvailable`, `ResourceHealthAvailabilityStateDegraded`, `ResourceHealthAvailabilityStateUnavailable`, `ResourceHealthAvailabilityStateUnknown`
- New enum type `ResourceHealthAvailabilityStateSignalBehavior` with values `ResourceHealthAvailabilityStateSignalBehaviorDisabled`, `ResourceHealthAvailabilityStateSignalBehaviorEnabled`
- New enum type `ResourceHealthCategory` with values `ResourceHealthCategoryPlanned`, `ResourceHealthCategoryUnplanned`
- New enum type `ResourceHealthReasonChronicity` with values `ResourceHealthReasonChronicityPersistent`, `ResourceHealthReasonChronicityTransient`
- New enum type `ResourceHealthReasonType` with values `ResourceHealthReasonTypePlanned`, `ResourceHealthReasonTypeUnplanned`, `ResourceHealthReasonTypeUserInitiated`
- New function `*EntitiesClient.AddDataAnnotation(ctx context.Context, resourceGroupName string, healthModelName string, entityName string, body AddDataAnnotationRequest, options *EntitiesClientAddDataAnnotationOptions) (EntitiesClientAddDataAnnotationResponse, error)`
- New function `*EntitiesClient.GetDataAnnotations(ctx context.Context, resourceGroupName string, healthModelName string, entityName string, body GetDataAnnotationsRequest, options *EntitiesClientGetDataAnnotationsOptions) (EntitiesClientGetDataAnnotationsResponse, error)`
- New function `*EntitiesClient.GetSignalRecommendations(ctx context.Context, resourceGroupName string, healthModelName string, entityName string, options *EntitiesClientGetSignalRecommendationsOptions) (EntitiesClientGetSignalRecommendationsResponse, error)`
- New function `PossibleResourceHealthAvailabilityStateValues() []ResourceHealthAvailabilityState`
- New struct `AddDataAnnotationRequest`
- New struct `AzureResourceHealthSignal`
- New struct `AzureResourceHealthSignalStatus`
- New struct `DataAnnotation`
- New struct `GetDataAnnotationsRequest`
- New struct `GetDataAnnotationsResponse`
- New struct `GetSignalRecommendationsResponse`
- New struct `SignalConfiguration`
- New field `ResourceHealth` in struct `AzureResourceSignals`
- New field `AddResourceHealthSignal` in struct `DiscoveryRuleProperties`
- New field `NextMarker`, `Top` in struct `EntityHistoryRequest`
- New field `NextMarker` in struct `EntityHistoryResponse`
- New field `NextMarker`, `Top` in struct `SignalHistoryRequest`
- New field `NextMarker` in struct `SignalHistoryResponse`
- New field `AdditionalContext` in struct `SignalStatus`
- New field `LookBackWindow`, `Sensitivity` in struct `ThresholdRuleV2`


## 0.2.0 (2026-06-01)
### Breaking Changes

- Type of `EvaluationRule.DegradedRule` has been changed from `*ThresholdRule` to `*ThresholdRuleV2`
- Type of `EvaluationRule.UnhealthyRule` has been changed from `*ThresholdRule` to `*ThresholdRuleV2`
- `DependenciesAggregationTypeThresholds` from enum `DependenciesAggregationType` has been removed
- `HealthStateError` from enum `HealthState` has been removed
- `SignalOperatorEquals`, `SignalOperatorGreaterOrEquals`, `SignalOperatorLowerOrEquals`, `SignalOperatorLowerThan` from enum `SignalOperator` has been removed
- Enum `DynamicThresholdDirection` has been removed
- Enum `DynamicThresholdModel` has been removed
- Operation `*AuthenticationSettingsClient.CreateOrUpdate` has been changed to LRO, use `*AuthenticationSettingsClient.BeginCreateOrUpdate` instead.
- Operation `*AuthenticationSettingsClient.Delete` has been changed to LRO, use `*AuthenticationSettingsClient.BeginDelete` instead.
- Operation `*DiscoveryRulesClient.CreateOrUpdate` has been changed to LRO, use `*DiscoveryRulesClient.BeginCreateOrUpdate` instead.
- Operation `*DiscoveryRulesClient.Delete` has been changed to LRO, use `*DiscoveryRulesClient.BeginDelete` instead.
- Operation `*EntitiesClient.CreateOrUpdate` has been changed to LRO, use `*EntitiesClient.BeginCreateOrUpdate` instead.
- Operation `*EntitiesClient.Delete` has been changed to LRO, use `*EntitiesClient.BeginDelete` instead.
- Operation `*RelationshipsClient.CreateOrUpdate` has been changed to LRO, use `*RelationshipsClient.BeginCreateOrUpdate` instead.
- Operation `*RelationshipsClient.Delete` has been changed to LRO, use `*RelationshipsClient.BeginDelete` instead.
- Operation `*SignalDefinitionsClient.CreateOrUpdate` has been changed to LRO, use `*SignalDefinitionsClient.BeginCreateOrUpdate` instead.
- Operation `*SignalDefinitionsClient.Delete` has been changed to LRO, use `*SignalDefinitionsClient.BeginDelete` instead.
- Struct `AzureMonitorWorkspaceSignalGroup` has been removed
- Struct `AzureResourceSignalGroup` has been removed
- Struct `DependenciesSignalGroup` has been removed
- Struct `DynamicDetectionRule` has been removed
- Struct `HealthModelUpdateProperties` has been removed
- Struct `LogAnalyticsSignalGroup` has been removed
- Struct `ModelDiscoverySettings` has been removed
- Struct `SignalAssignment` has been removed
- Struct `SignalGroup` has been removed
- Struct `ThresholdRule` has been removed
- Field `DeletionDate`, `ErrorMessage`, `NumberOfDiscoveredEntities`, `ResourceGraphQuery` of struct `DiscoveryRuleProperties` has been removed
- Field `DeletionDate`, `Kind`, `Labels`, `Signals` of struct `EntityProperties` has been removed
- Field `DynamicDetectionRule` of struct `EvaluationRule` has been removed
- Field `DataplaneEndpoint`, `Discovery` of struct `HealthModelProperties` has been removed
- Field `Properties` of struct `HealthModelUpdate` has been removed
- Field `DeletionDate`, `Labels` of struct `LogAnalyticsQuerySignalDefinitionProperties` has been removed
- Field `DeletionDate`, `Labels` of struct `PrometheusMetricsSignalDefinitionProperties` has been removed
- Field `DeletionDate`, `Labels` of struct `RelationshipProperties` has been removed
- Field `DeletionDate`, `Labels` of struct `ResourceMetricSignalDefinitionProperties` has been removed
- Field `DeletionDate`, `Labels` of struct `SignalDefinitionProperties` has been removed

### Features Added

- New value `DependenciesAggregationTypeMaxNotHealthy`, `DependenciesAggregationTypeMinHealthy` added to enum type `DependenciesAggregationType`
- New value `HealthStateUnhealthy` added to enum type `HealthState`
- New value `SignalKindExternalSignal` added to enum type `SignalKind`
- New value `SignalOperatorEqual`, `SignalOperatorGreaterThanOrEqual`, `SignalOperatorLessThan`, `SignalOperatorLessThanOrEqual`, `SignalOperatorNotEqual` added to enum type `SignalOperator`
- New enum type `DependenciesAggregationUnit` with values `DependenciesAggregationUnitAbsolute`, `DependenciesAggregationUnitPercentage`
- New enum type `DiscoveryRuleKind` with values `DiscoveryRuleKindApplicationInsightsTopology`, `DiscoveryRuleKindResourceGraphQuery`
- New function `*ApplicationInsightsTopologySpecification.GetDiscoveryRuleSpecification() *DiscoveryRuleSpecification`
- New function `*AzureResourceSignal.GetSignalInstanceProperties() *SignalInstanceProperties`
- New function `*DiscoveryRuleSpecification.GetDiscoveryRuleSpecification() *DiscoveryRuleSpecification`
- New function `*EntitiesClient.GetHistory(ctx context.Context, resourceGroupName string, healthModelName string, entityName string, body EntityHistoryRequest, options *EntitiesClientGetHistoryOptions) (EntitiesClientGetHistoryResponse, error)`
- New function `*EntitiesClient.GetSignalHistory(ctx context.Context, resourceGroupName string, healthModelName string, entityName string, body SignalHistoryRequest, options *EntitiesClientGetSignalHistoryOptions) (EntitiesClientGetSignalHistoryResponse, error)`
- New function `*EntitiesClient.IngestHealthReport(ctx context.Context, resourceGroupName string, healthModelName string, entityName string, body HealthReportRequest, options *EntitiesClientIngestHealthReportOptions) (EntitiesClientIngestHealthReportResponse, error)`
- New function `*ExternalSignal.GetSignalInstanceProperties() *SignalInstanceProperties`
- New function `*LogAnalyticsSignal.GetSignalInstanceProperties() *SignalInstanceProperties`
- New function `*PrometheusMetricsSignal.GetSignalInstanceProperties() *SignalInstanceProperties`
- New function `*ResourceGraphQuerySpecification.GetDiscoveryRuleSpecification() *DiscoveryRuleSpecification`
- New function `*SignalInstanceProperties.GetSignalInstanceProperties() *SignalInstanceProperties`
- New struct `ApplicationInsightsTopologySpecification`
- New struct `AzureMonitorWorkspaceSignals`
- New struct `AzureResourceSignal`
- New struct `AzureResourceSignals`
- New struct `DependenciesSignalGroupV2`
- New struct `DiscoveryError`
- New struct `EntityHistoryRequest`
- New struct `EntityHistoryResponse`
- New struct `ExternalSignal`
- New struct `ExternalSignalGroup`
- New struct `HealthReportEvaluationRule`
- New struct `HealthReportRequest`
- New struct `HealthStateTransition`
- New struct `LogAnalyticsSignal`
- New struct `LogAnalyticsSignals`
- New struct `PrometheusMetricsSignal`
- New struct `ResourceGraphQuerySpecification`
- New struct `SignalGroups`
- New struct `SignalHistoryDataPoint`
- New struct `SignalHistoryRequest`
- New struct `SignalHistoryResponse`
- New struct `SignalStatus`
- New struct `ThresholdRuleV2`
- New field `Error`, `Specification` in struct `DiscoveryRuleProperties`
- New field `SignalGroups`, `Tags` in struct `EntityProperties`
- New field `Tags` in struct `LogAnalyticsQuerySignalDefinitionProperties`
- New field `Tags` in struct `PrometheusMetricsSignalDefinitionProperties`
- New field `Tags` in struct `RelationshipProperties`
- New field `Tags` in struct `ResourceMetricSignalDefinitionProperties`
- New field `Tags` in struct `SignalDefinitionProperties`


## 0.1.0 (2025-06-27)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cloudhealth/armcloudhealth` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).