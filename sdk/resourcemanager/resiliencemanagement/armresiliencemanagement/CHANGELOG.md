# Release History

## 0.2.0 (2026-08-22)
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
- New function `*DrillRunsClient.ListReportDownloadURL(ctx context.Context, serviceGroupName string, drillName string, drillRunName string, options *DrillRunsClientListReportDownloadURLOptions) (DrillRunsClientListReportDownloadURLResponse, error)`
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