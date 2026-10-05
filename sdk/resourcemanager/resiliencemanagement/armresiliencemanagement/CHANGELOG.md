# Release History

## 0.3.0 (2026-10-05)
### Breaking Changes

- Type of `DrillResource.Properties` has been changed from `*DrillResourceProperties` to `DrillResourcePropertiesClassification`
- `UsagePlanTypeBasic` from enum `UsagePlanType` has been removed
- Enum `GoalAssignmentType` has been removed
- Enum `GoalType` has been removed
- Enum `MembershipType` has been removed
- Enum `RequirementSelected` has been removed
- Enum `ResilienceHealthStatus` has been removed
- Enum `UnifiedResilienceItemRequirementSelected` has been removed
- Function `*ClientFactory.NewGoalTemplatesClient` has been removed
- Function `NewGoalTemplatesClient` has been removed
- Function `*GoalTemplatesClient.BeginCreateOrUpdate` has been removed
- Function `*GoalTemplatesClient.BeginDelete` has been removed
- Function `*GoalTemplatesClient.Get` has been removed
- Function `*GoalTemplatesClient.NewListPager` has been removed
- Function `*GoalTemplatesClient.BeginUpdate` has been removed
- Struct `GoalTemplate` has been removed
- Struct `GoalTemplateListResult` has been removed
- Struct `GoalTemplateProperties` has been removed
- Struct `RecommendationsData` has been removed
- Struct `RecommendationsHighAvailabilityData` has been removed
- Struct `ServiceGroupMembership` has been removed
- Field `ActivePhysicalZones`, `AdvisorHaRecommendationID`, `HaStatus`, `RecoveryPhysicalZones` of struct `DrillResourceProperties` has been removed
- Field `GoalAssignmentType`, `GoalTemplateID` of struct `GoalAssignmentProperties` has been removed
- Field `DisasterRecoveryAttestationStatus`, `DisasterRecoveryGoalParticipation`, `ExclusionReasonForDisasterRecoveryGoals`, `ExclusionReasonForHighAvailabilityGoals`, `HighAvailabilityAttestationStatus`, `HighAvailabilityGoalParticipation`, `ServiceGroupMemberships`, `UserConfirmationForHighAvailability` of struct `GoalResourceProperties` has been removed
- Field `RegionalRecoveryPointEstimatedInMinutes`, `RegionalRecoveryPointObjectiveInMinutes`, `RegionalRecoveryPointObjectiveStatus`, `RegionalRecoveryTimeActualInMinutes`, `RegionalRecoveryTimeObjectiveInMinutes`, `RegionalRecoveryTimeObjectiveStatus`, `RequireDisasterRecovery`, `RequireHighAvailability`, `TemplateID` of struct `GoalsData` has been removed
- Field `DiscoveryRuleID` of struct `HealthModelMonitoringProperties` has been removed
- Field `ServiceLevelObjectiveResourceID` of struct `ServiceLevelResource` has been removed
- Field `Recommendations` of struct `UnifiedResilienceItemProperties` has been removed

### Features Added

- New value `ResourceProtectionSolutionTypeAzureCosmosDB`, `ResourceProtectionSolutionTypeAzureNetAppFiles`, `ResourceProtectionSolutionTypeAzureServiceBus`, `ResourceProtectionSolutionTypeAzureStorageAccount`, `ResourceProtectionSolutionTypeAzureTemplate` added to enum type `ResourceProtectionSolutionType`
- New enum type `FaultEligibility` with values `FaultEligibilityEligible`, `FaultEligibilityIneligible`, `FaultEligibilityUnknown`
- New enum type `FaultIneligibleReason` with values `FaultIneligibleReasonRecoveryPlanNotConfigured`, `FaultIneligibleReasonResourceNotIncludedInRecoveryPlan`
- New enum type `RegionalResiliencyStatus` with values `RegionalResiliencyStatusNotResilient`, `RegionalResiliencyStatusResilient`
- New enum type `ReplicationMode` with values `ReplicationModeActiveActive`, `ReplicationModeActivePassive`, `ReplicationModeNone`
- New enum type `ResourceInclusionDisabledReason` with values `ResourceInclusionDisabledReasonResourceActiveActiveProtection`, `ResourceInclusionDisabledReasonResourceHighlyAvailable`
- New function `*DrillResourceProperties.GetDrillResourceProperties() *DrillResourceProperties`
- New function `*RegionalDrillResourceProperties.GetDrillResourceProperties() *DrillResourceProperties`
- New function `*ResourceAzureTemplateProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceCosmosDBProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceNetAppFilesProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceServiceBusProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceStorageAccountProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ZonalDrillResourceProperties.GetDrillResourceProperties() *DrillResourceProperties`
- New struct `GoalAssignmentPropertiesOfDrill`
- New struct `RegionalDrillResourceProperties`
- New struct `RegionalObjectives`
- New struct `ResourceAzureTemplateProtectionSetting`
- New struct `ResourceCosmosDBProtectionSetting`
- New struct `ResourceNetAppFilesProtectionSetting`
- New struct `ResourceServiceBusProtectionSetting`
- New struct `ResourceStorageAccountProtectionSetting`
- New struct `UnifiedResilienceItemBillingInfo`
- New struct `UnifiedResilienceItemGoalRequirement`
- New struct `UnifiedResilienceItemRegionalResiliencyPosture`
- New struct `UnifiedResilienceItemResiliencyPosture`
- New struct `UnifiedResilienceItemZonalResiliencyPosture`
- New struct `ZonalDrillResourceProperties`
- New field `DrillRbacOnGoalAssignment`, `GoalAssignment`, `HealthModelAssociatedWithServiceGroup`, `RbacNeededForDrillOnGoalAssignment`, `RecoveryPlan` in struct `AttentionReason`
- New field `GoalAssignmentProperties` in struct `DrillProperties`
- New field `RecoveryTimeObjective` in struct `DrillRunProperties`
- New field `GoalAssignmentProperties` in struct `DrillUpdateProperties`
- New field `RegionalObjectives`, `RequireRegionalResiliency` in struct `GoalAssignmentProperties`
- New field `RegionalResiliency` in struct `GoalResourceProperties`
- New field `RegionalResiliency`, `ZonalResiliency` in struct `GoalsData`
- New field `HealthModelID` in struct `HealthModelMonitoringProperties`
- New field `LastRunRecoveryTimeActual` in struct `LastRunProperties`
- New field `InclusionDisabledReasons` in struct `RecoveryResourceProperties`
- New field `GoalAssignmentProperties` in struct `RegionalDrillProperties`
- New field `ReplicationMode` in struct `ResourceProtectionSolutionSettings`
- New field `BillingInfo`, `ResiliencyPosture` in struct `UnifiedResilienceItemProperties`
- New field `GoalAssignmentProperties` in struct `ZonalDrillProperties`


## 0.2.0 (2026-09-23)
### Breaking Changes

- Function `*DrillRunsClient.BeginFailOver` parameter(s) have been changed from `(ctx context.Context, serviceGroupName string, operationID string, drillName string, drillRunName string, body DrillRunFailoverRequest, options *DrillRunsClientBeginFailOverOptions)` to `(ctx context.Context, serviceGroupName string, operationID string, drillName string, drillRunName string, options *DrillRunsClientBeginFailOverOptions)`
- Type of `GoalResourceProperties.UserConfirmationForHighAvailability` has been changed from `[]*UserConfirmationForHighAvailabilityItem` to `[]*UserConfirmationItem`
- Struct `ManagedOnBehalfOfConfiguration` has been removed
- Struct `MoboBrokerResource` has been removed
- Struct `UserConfirmationForHighAvailabilityItem` has been removed
- Field `ManagedOnBehalfOfConfiguration` of struct `DrillProperties` has been removed
- Field `ManagedOnBehalfOfConfiguration` of struct `RegionalDrillProperties` has been removed
- Field `ManagedOnBehalfOfConfiguration` of struct `ZonalDrillProperties` has been removed

### Features Added

- New value `ProvisioningStateNeedsAttention` added to enum type `ProvisioningState`
- New enum type `DrillReportFinalizationState` with values `DrillReportFinalizationStateFinalized`, `DrillReportFinalizationStateNotFinalized`
- New enum type `DrillReportFormat` with values `DrillReportFormatHTML`
- New enum type `DrillReportGenerationStatus` with values `DrillReportGenerationStatusFailed`, `DrillReportGenerationStatusInProgress`, `DrillReportGenerationStatusNotStarted`, `DrillReportGenerationStatusSucceeded`
- New enum type `DrillRunTasks` with values `DrillRunTasksFailover`, `DrillRunTasksFailoverReverse`, `DrillRunTasksReprotect`, `DrillRunTasksReprotectReverse`
- New enum type `ResourceFeasibilityReviewStatus` with values `ResourceFeasibilityReviewStatusFlagged`, `ResourceFeasibilityReviewStatusNotApplicable`, `ResourceFeasibilityReviewStatusPassed`, `ResourceFeasibilityReviewStatusUnavailable`
- New enum type `ResourceFeasibilityReviewType` with values `ResourceFeasibilityReviewTypeSKUCapacity`
- New enum type `SliType` with values `SliTypeAvailability`, `SliTypeLatency`
- New enum type `SliTypeMatchState` with values `SliTypeMatchStateMatched`, `SliTypeMatchStateMismatched`
- New function `*DrillRunsClient.BeginGenerateReport(ctx context.Context, serviceGroupName string, operationID string, drillName string, drillRunName string, options *DrillRunsClientBeginGenerateReportOptions) (*runtime.Poller[DrillRunsClientGenerateReportResponse], error)`
- New function `*DrillRunsClient.BeginListReportDownloadURL(ctx context.Context, serviceGroupName string, operationID string, drillName string, drillRunName string, body ListReportDownloadURLRequest, options *DrillRunsClientBeginListReportDownloadURLOptions) (*runtime.Poller[DrillRunsClientListReportDownloadURLResponse], error)`
- New function `*ResourceCrossZoneVMRecoveryProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New struct `DrillReportSummary`
- New struct `DrillRunReprotectRequest`
- New struct `HealthModelMonitoringProperties`
- New struct `ListReportDownloadURLRequest`
- New struct `ListReportDownloadURLResponse`
- New struct `ReportStageStatus`
- New struct `ResiliencyProperties`
- New struct `ResourceCrossZoneVMRecoveryProtectionSetting`
- New struct `ResourceFeasibilityReview`
- New struct `SKUDetails`
- New struct `SliAttentionStatus`
- New struct `SliMonitoringProperties`
- New struct `SliSelection`
- New struct `UserConfirmationItem`
- New field `DiscoveryRuleExists`, `DrillRbacOnHealthModel`, `DrillRbacOnSli`, `HealthModelExists`, `MonitoringSourceNotConfigured`, `RbacNeededForDrillOnHealthModel`, `SliAttentionStatuses` in struct `AttentionReason`
- New field `HealthModelMonitoringProperties`, `SliMonitoringProperties` in struct `DrillProperties`
- New field `Report` in struct `DrillRunProperties`
- New field `Body` in struct `DrillRunsClientBeginFailOverOptions`
- New field `Body` in struct `DrillRunsClientBeginReprotectOptions`
- New field `HealthModelMonitoringProperties`, `SliMonitoringProperties` in struct `DrillUpdateProperties`
- New field `RequireZonalResiliency` in struct `GoalAssignmentProperties`
- New field `ZonalResiliency` in struct `GoalResourceProperties`
- New field `ResourceFeasibilityReviews` in struct `OperationQualificationDetails`
- New field `HealthModelMonitoringProperties`, `SliMonitoringProperties` in struct `RegionalDrillProperties`
- New field `OperationName` in struct `ValidateForExecutionProperties`
- New field `HealthModelMonitoringProperties`, `SliMonitoringProperties` in struct `ZonalDrillProperties`


## 0.1.0 (2026-06-17)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resiliencemanagement/armresiliencemanagement` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).