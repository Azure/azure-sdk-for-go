# Release History

## 0.2.0 (2026-09-15)
### Breaking Changes

- Function `*DrillRunsClient.BeginFailOver` parameter(s) have been changed from `(ctx context.Context, serviceGroupName string, operationID string, drillName string, drillRunName string, body DrillRunFailoverRequest, options *DrillRunsClientBeginFailOverOptions)` to `(ctx context.Context, serviceGroupName string, operationID string, drillName string, drillRunName string, options *DrillRunsClientBeginFailOverOptions)`
- `UsagePlanTypeBasic` from enum `UsagePlanType` has been removed
- Enum `GoalAssignmentType` has been removed
- Enum `GoalType` has been removed
- Enum `MembershipType` has been removed
- Enum `RequirementSelected` has been removed
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
- Struct `ManagedOnBehalfOfConfiguration` has been removed
- Struct `MoboBrokerResource` has been removed
- Struct `ServiceGroupMembership` has been removed
- Struct `UserConfirmationForHighAvailabilityItem` has been removed
- Field `ManagedOnBehalfOfConfiguration` of struct `DrillProperties` has been removed
- Field `GoalAssignmentType`, `GoalTemplateID` of struct `GoalAssignmentProperties` has been removed
- Field `DisasterRecoveryAttestationStatus`, `DisasterRecoveryGoalParticipation`, `ExclusionReasonForDisasterRecoveryGoals`, `ExclusionReasonForHighAvailabilityGoals`, `HighAvailabilityAttestationStatus`, `HighAvailabilityGoalParticipation`, `ServiceGroupMemberships`, `UserConfirmationForHighAvailability` of struct `GoalResourceProperties` has been removed
- Field `ManagedOnBehalfOfConfiguration` of struct `RegionalDrillProperties` has been removed
- Field `ServiceLevelObjectiveResourceID` of struct `ServiceLevelResource` has been removed
- Field `ManagedOnBehalfOfConfiguration` of struct `ZonalDrillProperties` has been removed

### Features Added

- New value `ProvisioningStateNeedsAttention` added to enum type `ProvisioningState`
- New value `ResourceProtectionSolutionTypeAzureCosmosDB`, `ResourceProtectionSolutionTypeAzureNetAppFiles`, `ResourceProtectionSolutionTypeAzureServiceBus`, `ResourceProtectionSolutionTypeAzureStorageAccount`, `ResourceProtectionSolutionTypeAzureTemplate` added to enum type `ResourceProtectionSolutionType`
- New enum type `DrillReportFinalizationState` with values `DrillReportFinalizationStateFinalized`, `DrillReportFinalizationStateNotFinalized`
- New enum type `DrillReportFormat` with values `DrillReportFormatHTML`
- New enum type `DrillReportGenerationStatus` with values `DrillReportGenerationStatusFailed`, `DrillReportGenerationStatusInProgress`, `DrillReportGenerationStatusNotStarted`, `DrillReportGenerationStatusSucceeded`
- New enum type `DrillRunTasks` with values `DrillRunTasksFailover`, `DrillRunTasksFailoverReverse`, `DrillRunTasksReprotect`, `DrillRunTasksReprotectReverse`
- New enum type `ReplicationMode` with values `ReplicationModeActiveActive`, `ReplicationModeActivePassive`, `ReplicationModeNone`
- New enum type `ResourceFeasibilityReviewStatus` with values `ResourceFeasibilityReviewStatusFlagged`, `ResourceFeasibilityReviewStatusNotApplicable`, `ResourceFeasibilityReviewStatusPassed`, `ResourceFeasibilityReviewStatusUnavailable`
- New enum type `ResourceFeasibilityReviewType` with values `ResourceFeasibilityReviewTypeSKUCapacity`
- New enum type `ResourceInclusionDisabledReason` with values `ResourceInclusionDisabledReasonResourceActiveActiveProtection`, `ResourceInclusionDisabledReasonResourceHighlyAvailable`
- New enum type `SliType` with values `SliTypeAvailability`, `SliTypeLatency`
- New enum type `SliTypeMatchState` with values `SliTypeMatchStateMatched`, `SliTypeMatchStateMismatched`
- New function `*DrillRunsClient.BeginGenerateReport(ctx context.Context, serviceGroupName string, operationID string, drillName string, drillRunName string, options *DrillRunsClientBeginGenerateReportOptions) (*runtime.Poller[DrillRunsClientGenerateReportResponse], error)`
- New function `*DrillRunsClient.BeginListReportDownloadURL(ctx context.Context, serviceGroupName string, operationID string, drillName string, drillRunName string, body ListReportDownloadURLRequest, options *DrillRunsClientBeginListReportDownloadURLOptions) (*runtime.Poller[DrillRunsClientListReportDownloadURLResponse], error)`
- New function `PossibleSliTypeValues() []SliType`
- New function `*ResourceAzureTemplateProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceCosmosDBProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceCrossZoneVMRecoveryProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceNetAppFilesProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceServiceBusProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New function `*ResourceStorageAccountProtectionSetting.GetResourceBaseProtectionSolutionSetting() *ResourceBaseProtectionSolutionSetting`
- New struct `DrillReportSummary`
- New struct `DrillRunReprotectRequest`
- New struct `GoalAssignmentPropertiesOfDrill`
- New struct `HealthModelMonitoringProperties`
- New struct `ListReportDownloadURLRequest`
- New struct `ListReportDownloadURLResponse`
- New struct `RegionalObjectives`
- New struct `ReportStageStatus`
- New struct `ResiliencyProperties`
- New struct `ResourceAzureTemplateProtectionSetting`
- New struct `ResourceCosmosDBProtectionSetting`
- New struct `ResourceCrossZoneVMRecoveryProtectionSetting`
- New struct `ResourceFeasibilityReview`
- New struct `ResourceNetAppFilesProtectionSetting`
- New struct `ResourceServiceBusProtectionSetting`
- New struct `ResourceStorageAccountProtectionSetting`
- New struct `SKUDetails`
- New struct `SliAttentionStatus`
- New struct `SliMonitoringProperties`
- New struct `SliSelection`
- New struct `UserConfirmationItem`
- New field `DiscoveryRuleExists`, `DrillRbacOnGoalAssignment`, `DrillRbacOnHealthModel`, `DrillRbacOnSli`, `GoalAssignment`, `HealthModelExists`, `MonitoringSourceNotConfigured`, `RbacNeededForDrillOnGoalAssignment`, `RbacNeededForDrillOnHealthModel`, `RecoveryPlan`, `SliAttentionStatuses` in struct `AttentionReason`
- New field `GoalAssignmentProperties`, `HealthModelMonitoringProperties`, `SliMonitoringProperties` in struct `DrillProperties`
- New field `RecoveryTimeObjective`, `Report` in struct `DrillRunProperties`
- New field `Body` in struct `DrillRunsClientBeginFailOverOptions`
- New field `Body` in struct `DrillRunsClientBeginReprotectOptions`
- New field `GoalAssignmentProperties`, `HealthModelMonitoringProperties`, `SliMonitoringProperties` in struct `DrillUpdateProperties`
- New field `RegionalObjectives`, `RequireRegionalResiliency`, `RequireZonalResiliency` in struct `GoalAssignmentProperties`
- New field `RegionalResiliency`, `ZonalResiliency` in struct `GoalResourceProperties`
- New field `LastRunRecoveryTimeActual` in struct `LastRunProperties`
- New field `ResourceFeasibilityReviews` in struct `OperationQualificationDetails`
- New field `InclusionDisabledReasons` in struct `RecoveryResourceProperties`
- New field `GoalAssignmentProperties`, `HealthModelMonitoringProperties`, `SliMonitoringProperties` in struct `RegionalDrillProperties`
- New field `ReplicationMode` in struct `ResourceProtectionSolutionSettings`
- New field `OperationName` in struct `ValidateForExecutionProperties`
- New field `GoalAssignmentProperties`, `HealthModelMonitoringProperties`, `SliMonitoringProperties` in struct `ZonalDrillProperties`


## 0.1.0 (2026-06-17)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resiliencemanagement/armresiliencemanagement` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).