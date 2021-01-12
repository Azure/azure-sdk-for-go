Generated from https://github.com/Azure/azure-rest-api-specs/tree/4c93f28f89435f6d244f4db61bbf810b5d20f09f/specification/monitor/resource-manager/readme.md tag: `package-2019-11`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Constants

1. KnownDataCollectionRuleAssociationProvisioningState.Creating
1. KnownDataCollectionRuleAssociationProvisioningState.Deleting
1. KnownDataCollectionRuleAssociationProvisioningState.Updating
1. KnownDataCollectionRuleProvisioningState.KnownDataCollectionRuleProvisioningStateCreating
1. KnownDataCollectionRuleProvisioningState.KnownDataCollectionRuleProvisioningStateDeleting
1. KnownDataCollectionRuleProvisioningState.KnownDataCollectionRuleProvisioningStateFailed
1. KnownDataCollectionRuleProvisioningState.KnownDataCollectionRuleProvisioningStateSucceeded
1. KnownDataCollectionRuleProvisioningState.KnownDataCollectionRuleProvisioningStateUpdating
1. KnownDataFlowStreams.MicrosoftAntiMalwareStatus
1. KnownDataFlowStreams.MicrosoftAuditd
1. KnownDataFlowStreams.MicrosoftCISCOASA
1. KnownDataFlowStreams.MicrosoftCommonSecurityLog
1. KnownDataFlowStreams.MicrosoftComputerGroup
1. KnownDataFlowStreams.MicrosoftEvent
1. KnownDataFlowStreams.MicrosoftFirewallLog
1. KnownDataFlowStreams.MicrosoftHealthStateChange
1. KnownDataFlowStreams.MicrosoftHeartbeat
1. KnownDataFlowStreams.MicrosoftInsightsMetrics
1. KnownDataFlowStreams.MicrosoftOperationLog
1. KnownDataFlowStreams.MicrosoftPerf
1. KnownDataFlowStreams.MicrosoftProcessInvestigator
1. KnownDataFlowStreams.MicrosoftProtectionStatus
1. KnownDataFlowStreams.MicrosoftRomeDetectionEvent
1. KnownDataFlowStreams.MicrosoftSecurityBaseline
1. KnownDataFlowStreams.MicrosoftSecurityBaselineSummary
1. KnownDataFlowStreams.MicrosoftSecurityEvent
1. KnownDataFlowStreams.MicrosoftSyslog
1. KnownDataFlowStreams.MicrosoftWindowsEvent
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftAntiMalwareStatus
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftAuditd
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftCISCOASA
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftCommonSecurityLog
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftComputerGroup
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftEvent
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftFirewallLog
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftHealthStateChange
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftHeartbeat
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftInsightsMetrics
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftOperationLog
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftPerf
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftProcessInvestigator
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftProtectionStatus
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftRomeDetectionEvent
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftSecurityBaseline
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftSecurityBaselineSummary
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftSecurityEvent
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftSyslog
1. KnownExtensionDataSourceStreams.KnownExtensionDataSourceStreamsMicrosoftWindowsEvent
1. KnownPerfCounterDataSourceScheduledTransferPeriod.PT15M
1. KnownPerfCounterDataSourceScheduledTransferPeriod.PT1M
1. KnownPerfCounterDataSourceScheduledTransferPeriod.PT30M
1. KnownPerfCounterDataSourceScheduledTransferPeriod.PT5M
1. KnownPerfCounterDataSourceScheduledTransferPeriod.PT60M
1. KnownPerfCounterDataSourceStreams.KnownPerfCounterDataSourceStreamsMicrosoftInsightsMetrics
1. KnownPerfCounterDataSourceStreams.KnownPerfCounterDataSourceStreamsMicrosoftPerf
1. KnownSyslogDataSourceFacilityNames.Auth
1. KnownSyslogDataSourceFacilityNames.Authpriv
1. KnownSyslogDataSourceFacilityNames.Cron
1. KnownSyslogDataSourceFacilityNames.Daemon
1. KnownSyslogDataSourceFacilityNames.Kern
1. KnownSyslogDataSourceFacilityNames.Local0
1. KnownSyslogDataSourceFacilityNames.Local1
1. KnownSyslogDataSourceFacilityNames.Local2
1. KnownSyslogDataSourceFacilityNames.Local3
1. KnownSyslogDataSourceFacilityNames.Local4
1. KnownSyslogDataSourceFacilityNames.Local5
1. KnownSyslogDataSourceFacilityNames.Local6
1. KnownSyslogDataSourceFacilityNames.Local7
1. KnownSyslogDataSourceFacilityNames.Lpr
1. KnownSyslogDataSourceFacilityNames.Mail
1. KnownSyslogDataSourceFacilityNames.Mark
1. KnownSyslogDataSourceFacilityNames.News
1. KnownSyslogDataSourceFacilityNames.Syslog
1. KnownSyslogDataSourceFacilityNames.UUCP
1. KnownSyslogDataSourceFacilityNames.User
1. KnownSyslogDataSourceLogLevels.Alert
1. KnownSyslogDataSourceLogLevels.Critical
1. KnownSyslogDataSourceLogLevels.Debug
1. KnownSyslogDataSourceLogLevels.Emergency
1. KnownSyslogDataSourceLogLevels.Error
1. KnownSyslogDataSourceLogLevels.Info
1. KnownSyslogDataSourceLogLevels.Notice
1. KnownSyslogDataSourceLogLevels.Warning
1. KnownSyslogDataSourceStreams.KnownSyslogDataSourceStreamsMicrosoftSyslog
1. KnownWindowsEventLogDataSourceScheduledTransferPeriod.KnownWindowsEventLogDataSourceScheduledTransferPeriodPT15M
1. KnownWindowsEventLogDataSourceScheduledTransferPeriod.KnownWindowsEventLogDataSourceScheduledTransferPeriodPT1M
1. KnownWindowsEventLogDataSourceScheduledTransferPeriod.KnownWindowsEventLogDataSourceScheduledTransferPeriodPT30M
1. KnownWindowsEventLogDataSourceScheduledTransferPeriod.KnownWindowsEventLogDataSourceScheduledTransferPeriodPT5M
1. KnownWindowsEventLogDataSourceScheduledTransferPeriod.KnownWindowsEventLogDataSourceScheduledTransferPeriodPT60M
1. KnownWindowsEventLogDataSourceStreams.KnownWindowsEventLogDataSourceStreamsMicrosoftEvent
1. KnownWindowsEventLogDataSourceStreams.KnownWindowsEventLogDataSourceStreamsMicrosoftWindowsEvent

### Removed Funcs

1. *DataCollectionRuleAssociationProxyOnlyResource.UnmarshalJSON([]byte) error
1. *DataCollectionRuleAssociationProxyOnlyResourceListResultIterator.Next() error
1. *DataCollectionRuleAssociationProxyOnlyResourceListResultIterator.NextWithContext(context.Context) error
1. *DataCollectionRuleAssociationProxyOnlyResourceListResultPage.Next() error
1. *DataCollectionRuleAssociationProxyOnlyResourceListResultPage.NextWithContext(context.Context) error
1. *DataCollectionRuleResource.UnmarshalJSON([]byte) error
1. *DataCollectionRuleResourceListResultIterator.Next() error
1. *DataCollectionRuleResourceListResultIterator.NextWithContext(context.Context) error
1. *DataCollectionRuleResourceListResultPage.Next() error
1. *DataCollectionRuleResourceListResultPage.NextWithContext(context.Context) error
1. DataCollectionRule.MarshalJSON() ([]byte, error)
1. DataCollectionRuleAssociation.MarshalJSON() ([]byte, error)
1. DataCollectionRuleAssociationProxyOnlyResource.MarshalJSON() ([]byte, error)
1. DataCollectionRuleAssociationProxyOnlyResourceListResult.IsEmpty() bool
1. DataCollectionRuleAssociationProxyOnlyResourceListResultIterator.NotDone() bool
1. DataCollectionRuleAssociationProxyOnlyResourceListResultIterator.Response() DataCollectionRuleAssociationProxyOnlyResourceListResult
1. DataCollectionRuleAssociationProxyOnlyResourceListResultIterator.Value() DataCollectionRuleAssociationProxyOnlyResource
1. DataCollectionRuleAssociationProxyOnlyResourceListResultPage.NotDone() bool
1. DataCollectionRuleAssociationProxyOnlyResourceListResultPage.Response() DataCollectionRuleAssociationProxyOnlyResourceListResult
1. DataCollectionRuleAssociationProxyOnlyResourceListResultPage.Values() []DataCollectionRuleAssociationProxyOnlyResource
1. DataCollectionRuleAssociationProxyOnlyResourceProperties.MarshalJSON() ([]byte, error)
1. DataCollectionRuleAssociationsClient.Create(context.Context, string, string, *DataCollectionRuleAssociationProxyOnlyResource) (DataCollectionRuleAssociationProxyOnlyResource, error)
1. DataCollectionRuleAssociationsClient.CreatePreparer(context.Context, string, string, *DataCollectionRuleAssociationProxyOnlyResource) (*http.Request, error)
1. DataCollectionRuleAssociationsClient.CreateResponder(*http.Response) (DataCollectionRuleAssociationProxyOnlyResource, error)
1. DataCollectionRuleAssociationsClient.CreateSender(*http.Request) (*http.Response, error)
1. DataCollectionRuleAssociationsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. DataCollectionRuleAssociationsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. DataCollectionRuleAssociationsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. DataCollectionRuleAssociationsClient.DeleteSender(*http.Request) (*http.Response, error)
1. DataCollectionRuleAssociationsClient.Get(context.Context, string, string) (DataCollectionRuleAssociationProxyOnlyResource, error)
1. DataCollectionRuleAssociationsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. DataCollectionRuleAssociationsClient.GetResponder(*http.Response) (DataCollectionRuleAssociationProxyOnlyResource, error)
1. DataCollectionRuleAssociationsClient.GetSender(*http.Request) (*http.Response, error)
1. DataCollectionRuleAssociationsClient.ListByResource(context.Context, string) (DataCollectionRuleAssociationProxyOnlyResourceListResultPage, error)
1. DataCollectionRuleAssociationsClient.ListByResourceComplete(context.Context, string) (DataCollectionRuleAssociationProxyOnlyResourceListResultIterator, error)
1. DataCollectionRuleAssociationsClient.ListByResourcePreparer(context.Context, string) (*http.Request, error)
1. DataCollectionRuleAssociationsClient.ListByResourceResponder(*http.Response) (DataCollectionRuleAssociationProxyOnlyResourceListResult, error)
1. DataCollectionRuleAssociationsClient.ListByResourceSender(*http.Request) (*http.Response, error)
1. DataCollectionRuleAssociationsClient.ListByRule(context.Context, string, string) (DataCollectionRuleAssociationProxyOnlyResourceListResultPage, error)
1. DataCollectionRuleAssociationsClient.ListByRuleComplete(context.Context, string, string) (DataCollectionRuleAssociationProxyOnlyResourceListResultIterator, error)
1. DataCollectionRuleAssociationsClient.ListByRulePreparer(context.Context, string, string) (*http.Request, error)
1. DataCollectionRuleAssociationsClient.ListByRuleResponder(*http.Response) (DataCollectionRuleAssociationProxyOnlyResourceListResult, error)
1. DataCollectionRuleAssociationsClient.ListByRuleSender(*http.Request) (*http.Response, error)
1. DataCollectionRuleResource.MarshalJSON() ([]byte, error)
1. DataCollectionRuleResourceListResult.IsEmpty() bool
1. DataCollectionRuleResourceListResultIterator.NotDone() bool
1. DataCollectionRuleResourceListResultIterator.Response() DataCollectionRuleResourceListResult
1. DataCollectionRuleResourceListResultIterator.Value() DataCollectionRuleResource
1. DataCollectionRuleResourceListResultPage.NotDone() bool
1. DataCollectionRuleResourceListResultPage.Response() DataCollectionRuleResourceListResult
1. DataCollectionRuleResourceListResultPage.Values() []DataCollectionRuleResource
1. DataCollectionRuleResourceProperties.MarshalJSON() ([]byte, error)
1. DataCollectionRulesClient.Create(context.Context, string, string, *DataCollectionRuleResource) (DataCollectionRuleResource, error)
1. DataCollectionRulesClient.CreatePreparer(context.Context, string, string, *DataCollectionRuleResource) (*http.Request, error)
1. DataCollectionRulesClient.CreateResponder(*http.Response) (DataCollectionRuleResource, error)
1. DataCollectionRulesClient.CreateSender(*http.Request) (*http.Response, error)
1. DataCollectionRulesClient.Delete(context.Context, string, string) (autorest.Response, error)
1. DataCollectionRulesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. DataCollectionRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. DataCollectionRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. DataCollectionRulesClient.Get(context.Context, string, string) (DataCollectionRuleResource, error)
1. DataCollectionRulesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. DataCollectionRulesClient.GetResponder(*http.Response) (DataCollectionRuleResource, error)
1. DataCollectionRulesClient.GetSender(*http.Request) (*http.Response, error)
1. DataCollectionRulesClient.ListByResourceGroup(context.Context, string) (DataCollectionRuleResourceListResultPage, error)
1. DataCollectionRulesClient.ListByResourceGroupComplete(context.Context, string) (DataCollectionRuleResourceListResultIterator, error)
1. DataCollectionRulesClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. DataCollectionRulesClient.ListByResourceGroupResponder(*http.Response) (DataCollectionRuleResourceListResult, error)
1. DataCollectionRulesClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. DataCollectionRulesClient.ListBySubscription(context.Context) (DataCollectionRuleResourceListResultPage, error)
1. DataCollectionRulesClient.ListBySubscriptionComplete(context.Context) (DataCollectionRuleResourceListResultIterator, error)
1. DataCollectionRulesClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. DataCollectionRulesClient.ListBySubscriptionResponder(*http.Response) (DataCollectionRuleResourceListResult, error)
1. DataCollectionRulesClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. DataCollectionRulesClient.Update(context.Context, string, string, *ResourceForUpdate) (DataCollectionRuleResource, error)
1. DataCollectionRulesClient.UpdatePreparer(context.Context, string, string, *ResourceForUpdate) (*http.Request, error)
1. DataCollectionRulesClient.UpdateResponder(*http.Response) (DataCollectionRuleResource, error)
1. DataCollectionRulesClient.UpdateSender(*http.Request) (*http.Response, error)
1. NewDataCollectionRuleAssociationProxyOnlyResourceListResultIterator(DataCollectionRuleAssociationProxyOnlyResourceListResultPage) DataCollectionRuleAssociationProxyOnlyResourceListResultIterator
1. NewDataCollectionRuleAssociationProxyOnlyResourceListResultPage(DataCollectionRuleAssociationProxyOnlyResourceListResult, func(context.Context, DataCollectionRuleAssociationProxyOnlyResourceListResult) (DataCollectionRuleAssociationProxyOnlyResourceListResult, error)) DataCollectionRuleAssociationProxyOnlyResourceListResultPage
1. NewDataCollectionRuleAssociationsClient(string) DataCollectionRuleAssociationsClient
1. NewDataCollectionRuleAssociationsClientWithBaseURI(string, string) DataCollectionRuleAssociationsClient
1. NewDataCollectionRuleResourceListResultIterator(DataCollectionRuleResourceListResultPage) DataCollectionRuleResourceListResultIterator
1. NewDataCollectionRuleResourceListResultPage(DataCollectionRuleResourceListResult, func(context.Context, DataCollectionRuleResourceListResult) (DataCollectionRuleResourceListResult, error)) DataCollectionRuleResourceListResultPage
1. NewDataCollectionRulesClient(string) DataCollectionRulesClient
1. NewDataCollectionRulesClientWithBaseURI(string, string) DataCollectionRulesClient
1. PossibleKnownDataCollectionRuleAssociationProvisioningStateValues() []KnownDataCollectionRuleAssociationProvisioningState
1. PossibleKnownDataCollectionRuleProvisioningStateValues() []KnownDataCollectionRuleProvisioningState
1. PossibleKnownDataFlowStreamsValues() []KnownDataFlowStreams
1. PossibleKnownExtensionDataSourceStreamsValues() []KnownExtensionDataSourceStreams
1. PossibleKnownPerfCounterDataSourceScheduledTransferPeriodValues() []KnownPerfCounterDataSourceScheduledTransferPeriod
1. PossibleKnownPerfCounterDataSourceStreamsValues() []KnownPerfCounterDataSourceStreams
1. PossibleKnownSyslogDataSourceFacilityNamesValues() []KnownSyslogDataSourceFacilityNames
1. PossibleKnownSyslogDataSourceLogLevelsValues() []KnownSyslogDataSourceLogLevels
1. PossibleKnownSyslogDataSourceStreamsValues() []KnownSyslogDataSourceStreams
1. PossibleKnownWindowsEventLogDataSourceScheduledTransferPeriodValues() []KnownWindowsEventLogDataSourceScheduledTransferPeriod
1. PossibleKnownWindowsEventLogDataSourceStreamsValues() []KnownWindowsEventLogDataSourceStreams
1. ResourceForUpdate.MarshalJSON() ([]byte, error)

## Struct Changes

### Removed Structs

1. AzureMonitorMetricsDestination
1. DataCollectionRule
1. DataCollectionRuleAssociation
1. DataCollectionRuleAssociationProxyOnlyResource
1. DataCollectionRuleAssociationProxyOnlyResourceListResult
1. DataCollectionRuleAssociationProxyOnlyResourceListResultIterator
1. DataCollectionRuleAssociationProxyOnlyResourceListResultPage
1. DataCollectionRuleAssociationProxyOnlyResourceProperties
1. DataCollectionRuleAssociationsClient
1. DataCollectionRuleDataSources
1. DataCollectionRuleDestinations
1. DataCollectionRuleResource
1. DataCollectionRuleResourceListResult
1. DataCollectionRuleResourceListResultIterator
1. DataCollectionRuleResourceListResultPage
1. DataCollectionRuleResourceProperties
1. DataCollectionRulesClient
1. DataFlow
1. DataSourcesSpec
1. DestinationsSpec
1. DestinationsSpecAzureMonitorMetrics
1. ErrorDetails
1. ErrorResponseError
1. ExtensionDataSource
1. LogAnalyticsDestination
1. PerfCounterDataSource
1. ResourceForUpdate
1. SyslogDataSource
1. WindowsEventLogDataSource

### Removed Struct Fields

1. ErrorResponse.Error

## Signature Changes

### Const Types

1. Failed changed type from KnownDataCollectionRuleAssociationProvisioningState to ProvisioningState
1. Succeeded changed type from KnownDataCollectionRuleAssociationProvisioningState to ProvisioningState

### New Constants

1. AggregationType.Average
1. AggregationType.Count
1. AggregationType.Maximum
1. AggregationType.Minimum
1. AggregationType.None
1. AggregationType.Total
1. AlertSeverity.Four
1. AlertSeverity.One
1. AlertSeverity.Three
1. AlertSeverity.Two
1. AlertSeverity.Zero
1. BaselineSensitivity.High
1. BaselineSensitivity.Low
1. BaselineSensitivity.Medium
1. CategoryType.Logs
1. CategoryType.Metrics
1. ComparisonOperationType.Equals
1. ComparisonOperationType.GreaterThan
1. ComparisonOperationType.GreaterThanOrEqual
1. ComparisonOperationType.LessThan
1. ComparisonOperationType.LessThanOrEqual
1. ComparisonOperationType.NotEquals
1. ConditionOperator.ConditionOperatorGreaterThan
1. ConditionOperator.ConditionOperatorGreaterThanOrEqual
1. ConditionOperator.ConditionOperatorLessThan
1. ConditionOperator.ConditionOperatorLessThanOrEqual
1. ConditionalOperator.ConditionalOperatorEqual
1. ConditionalOperator.ConditionalOperatorGreaterThan
1. ConditionalOperator.ConditionalOperatorLessThan
1. CriterionType.CriterionTypeDynamicThresholdCriterion
1. CriterionType.CriterionTypeMultiMetricCriteria
1. CriterionType.CriterionTypeStaticThresholdCriterion
1. DataStatus.NotPresent
1. DataStatus.Present
1. DynamicThresholdOperator.DynamicThresholdOperatorGreaterOrLessThan
1. DynamicThresholdOperator.DynamicThresholdOperatorGreaterThan
1. DynamicThresholdOperator.DynamicThresholdOperatorLessThan
1. DynamicThresholdSensitivity.DynamicThresholdSensitivityHigh
1. DynamicThresholdSensitivity.DynamicThresholdSensitivityLow
1. DynamicThresholdSensitivity.DynamicThresholdSensitivityMedium
1. Enabled.False
1. Enabled.True
1. EventLevel.EventLevelCritical
1. EventLevel.EventLevelError
1. EventLevel.EventLevelInformational
1. EventLevel.EventLevelVerbose
1. EventLevel.EventLevelWarning
1. MetricStatisticType.MetricStatisticTypeAverage
1. MetricStatisticType.MetricStatisticTypeMax
1. MetricStatisticType.MetricStatisticTypeMin
1. MetricStatisticType.MetricStatisticTypeSum
1. MetricTriggerType.MetricTriggerTypeConsecutive
1. MetricTriggerType.MetricTriggerTypeTotal
1. OdataType.OdataTypeMicrosoftAzureManagementInsightsModelsRuleManagementEventDataSource
1. OdataType.OdataTypeMicrosoftAzureManagementInsightsModelsRuleMetricDataSource
1. OdataType.OdataTypeRuleDataSource
1. OdataTypeBasicAction.OdataTypeAction
1. OdataTypeBasicAction.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesAlertingAction
1. OdataTypeBasicAction.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesLogToMetricAction
1. OdataTypeBasicMetricAlertCriteria.OdataTypeMetricAlertCriteria
1. OdataTypeBasicMetricAlertCriteria.OdataTypeMicrosoftAzureMonitorMultipleResourceMultipleMetricCriteria
1. OdataTypeBasicMetricAlertCriteria.OdataTypeMicrosoftAzureMonitorSingleResourceMultipleMetricCriteria
1. OdataTypeBasicMetricAlertCriteria.OdataTypeMicrosoftAzureMonitorWebtestLocationAvailabilityCriteria
1. OdataTypeBasicRuleAction.OdataTypeMicrosoftAzureManagementInsightsModelsRuleEmailAction
1. OdataTypeBasicRuleAction.OdataTypeMicrosoftAzureManagementInsightsModelsRuleWebhookAction
1. OdataTypeBasicRuleAction.OdataTypeRuleAction
1. OdataTypeBasicRuleCondition.OdataTypeMicrosoftAzureManagementInsightsModelsLocationThresholdRuleCondition
1. OdataTypeBasicRuleCondition.OdataTypeMicrosoftAzureManagementInsightsModelsManagementEventRuleCondition
1. OdataTypeBasicRuleCondition.OdataTypeMicrosoftAzureManagementInsightsModelsThresholdRuleCondition
1. OdataTypeBasicRuleCondition.OdataTypeRuleCondition
1. OnboardingStatus.NotOnboarded
1. OnboardingStatus.Onboarded
1. OnboardingStatus.Unknown
1. Operator.OperatorEquals
1. Operator.OperatorGreaterThan
1. Operator.OperatorGreaterThanOrEqual
1. Operator.OperatorLessThan
1. Operator.OperatorLessThanOrEqual
1. Operator.OperatorNotEquals
1. ProvisioningState.Canceled
1. ProvisioningState.Deploying
1. QueryType.ResultCount
1. ReceiverStatus.ReceiverStatusDisabled
1. ReceiverStatus.ReceiverStatusEnabled
1. ReceiverStatus.ReceiverStatusNotSpecified
1. RecurrenceFrequency.RecurrenceFrequencyDay
1. RecurrenceFrequency.RecurrenceFrequencyHour
1. RecurrenceFrequency.RecurrenceFrequencyMinute
1. RecurrenceFrequency.RecurrenceFrequencyMonth
1. RecurrenceFrequency.RecurrenceFrequencyNone
1. RecurrenceFrequency.RecurrenceFrequencySecond
1. RecurrenceFrequency.RecurrenceFrequencyWeek
1. RecurrenceFrequency.RecurrenceFrequencyYear
1. ResultType.Data
1. ResultType.Metadata
1. ScaleDirection.ScaleDirectionDecrease
1. ScaleDirection.ScaleDirectionIncrease
1. ScaleDirection.ScaleDirectionNone
1. ScaleRuleMetricDimensionOperationType.ScaleRuleMetricDimensionOperationTypeEquals
1. ScaleRuleMetricDimensionOperationType.ScaleRuleMetricDimensionOperationTypeNotEquals
1. ScaleType.ChangeCount
1. ScaleType.ExactCount
1. ScaleType.PercentChangeCount
1. Sensitivity.SensitivityHigh
1. Sensitivity.SensitivityLow
1. Sensitivity.SensitivityMedium
1. TimeAggregationOperator.TimeAggregationOperatorAverage
1. TimeAggregationOperator.TimeAggregationOperatorLast
1. TimeAggregationOperator.TimeAggregationOperatorMaximum
1. TimeAggregationOperator.TimeAggregationOperatorMinimum
1. TimeAggregationOperator.TimeAggregationOperatorTotal
1. TimeAggregationType.TimeAggregationTypeAverage
1. TimeAggregationType.TimeAggregationTypeCount
1. TimeAggregationType.TimeAggregationTypeLast
1. TimeAggregationType.TimeAggregationTypeMaximum
1. TimeAggregationType.TimeAggregationTypeMinimum
1. TimeAggregationType.TimeAggregationTypeTotal
1. Unit.UnitBitsPerSecond
1. Unit.UnitByteSeconds
1. Unit.UnitBytes
1. Unit.UnitBytesPerSecond
1. Unit.UnitCores
1. Unit.UnitCount
1. Unit.UnitCountPerSecond
1. Unit.UnitMilliCores
1. Unit.UnitMilliSeconds
1. Unit.UnitNanoCores
1. Unit.UnitPercent
1. Unit.UnitSeconds
1. Unit.UnitUnspecified

### New Funcs

1. *ActionGroupPatchBody.UnmarshalJSON([]byte) error
1. *ActionGroupResource.UnmarshalJSON([]byte) error
1. *ActivityLogAlertPatchBody.UnmarshalJSON([]byte) error
1. *ActivityLogAlertResource.UnmarshalJSON([]byte) error
1. *AlertRule.UnmarshalJSON([]byte) error
1. *AlertRuleResource.UnmarshalJSON([]byte) error
1. *AlertRuleResourcePatch.UnmarshalJSON([]byte) error
1. *AutoscaleSettingResource.UnmarshalJSON([]byte) error
1. *AutoscaleSettingResourceCollectionIterator.Next() error
1. *AutoscaleSettingResourceCollectionIterator.NextWithContext(context.Context) error
1. *AutoscaleSettingResourceCollectionPage.Next() error
1. *AutoscaleSettingResourceCollectionPage.NextWithContext(context.Context) error
1. *AutoscaleSettingResourcePatch.UnmarshalJSON([]byte) error
1. *AzureMonitorPrivateLinkScope.UnmarshalJSON([]byte) error
1. *AzureMonitorPrivateLinkScopeListResultIterator.Next() error
1. *AzureMonitorPrivateLinkScopeListResultIterator.NextWithContext(context.Context) error
1. *AzureMonitorPrivateLinkScopeListResultPage.Next() error
1. *AzureMonitorPrivateLinkScopeListResultPage.NextWithContext(context.Context) error
1. *BaselineResponse.UnmarshalJSON([]byte) error
1. *DiagnosticSettingsCategoryResource.UnmarshalJSON([]byte) error
1. *DiagnosticSettingsResource.UnmarshalJSON([]byte) error
1. *DynamicMetricCriteria.UnmarshalJSON([]byte) error
1. *EventDataCollectionIterator.Next() error
1. *EventDataCollectionIterator.NextWithContext(context.Context) error
1. *EventDataCollectionPage.Next() error
1. *EventDataCollectionPage.NextWithContext(context.Context) error
1. *LocationThresholdRuleCondition.UnmarshalJSON([]byte) error
1. *LogProfileResource.UnmarshalJSON([]byte) error
1. *LogProfileResourcePatch.UnmarshalJSON([]byte) error
1. *LogSearchRule.UnmarshalJSON([]byte) error
1. *LogSearchRuleResource.UnmarshalJSON([]byte) error
1. *LogSearchRuleResourcePatch.UnmarshalJSON([]byte) error
1. *ManagementEventRuleCondition.UnmarshalJSON([]byte) error
1. *MetricAlertCriteria.UnmarshalJSON([]byte) error
1. *MetricAlertMultipleResourceMultipleMetricCriteria.UnmarshalJSON([]byte) error
1. *MetricAlertProperties.UnmarshalJSON([]byte) error
1. *MetricAlertPropertiesPatch.UnmarshalJSON([]byte) error
1. *MetricAlertResource.UnmarshalJSON([]byte) error
1. *MetricAlertResourcePatch.UnmarshalJSON([]byte) error
1. *MetricAlertSingleResourceMultipleMetricCriteria.UnmarshalJSON([]byte) error
1. *MetricCriteria.UnmarshalJSON([]byte) error
1. *MultiMetricCriteria.UnmarshalJSON([]byte) error
1. *PrivateEndpointConnection.UnmarshalJSON([]byte) error
1. *PrivateEndpointConnectionListResultIterator.Next() error
1. *PrivateEndpointConnectionListResultIterator.NextWithContext(context.Context) error
1. *PrivateEndpointConnectionListResultPage.Next() error
1. *PrivateEndpointConnectionListResultPage.NextWithContext(context.Context) error
1. *PrivateLinkResource.UnmarshalJSON([]byte) error
1. *PrivateLinkResourceListResultIterator.Next() error
1. *PrivateLinkResourceListResultIterator.NextWithContext(context.Context) error
1. *PrivateLinkResourceListResultPage.Next() error
1. *PrivateLinkResourceListResultPage.NextWithContext(context.Context) error
1. *RuleCondition.UnmarshalJSON([]byte) error
1. *ScopedResource.UnmarshalJSON([]byte) error
1. *ScopedResourceListResultIterator.Next() error
1. *ScopedResourceListResultIterator.NextWithContext(context.Context) error
1. *ScopedResourceListResultPage.Next() error
1. *ScopedResourceListResultPage.NextWithContext(context.Context) error
1. *SingleMetricBaseline.UnmarshalJSON([]byte) error
1. *SubscriptionDiagnosticSettingsResource.UnmarshalJSON([]byte) error
1. *ThresholdRuleCondition.UnmarshalJSON([]byte) error
1. *VMInsightsOnboardingStatus.UnmarshalJSON([]byte) error
1. *WebtestLocationAvailabilityCriteria.UnmarshalJSON([]byte) error
1. *WorkspaceInfo.UnmarshalJSON([]byte) error
1. Action.AsAction() (*Action, bool)
1. Action.AsAlertingAction() (*AlertingAction, bool)
1. Action.AsBasicAction() (BasicAction, bool)
1. Action.AsLogToMetricAction() (*LogToMetricAction, bool)
1. Action.MarshalJSON() ([]byte, error)
1. ActionGroupPatchBody.MarshalJSON() ([]byte, error)
1. ActionGroupResource.MarshalJSON() ([]byte, error)
1. ActionGroupsClient.CreateOrUpdate(context.Context, string, string, ActionGroupResource) (ActionGroupResource, error)
1. ActionGroupsClient.CreateOrUpdatePreparer(context.Context, string, string, ActionGroupResource) (*http.Request, error)
1. ActionGroupsClient.CreateOrUpdateResponder(*http.Response) (ActionGroupResource, error)
1. ActionGroupsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ActionGroupsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. ActionGroupsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. ActionGroupsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ActionGroupsClient.DeleteSender(*http.Request) (*http.Response, error)
1. ActionGroupsClient.EnableReceiver(context.Context, string, string, EnableRequest) (autorest.Response, error)
1. ActionGroupsClient.EnableReceiverPreparer(context.Context, string, string, EnableRequest) (*http.Request, error)
1. ActionGroupsClient.EnableReceiverResponder(*http.Response) (autorest.Response, error)
1. ActionGroupsClient.EnableReceiverSender(*http.Request) (*http.Response, error)
1. ActionGroupsClient.Get(context.Context, string, string) (ActionGroupResource, error)
1. ActionGroupsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ActionGroupsClient.GetResponder(*http.Response) (ActionGroupResource, error)
1. ActionGroupsClient.GetSender(*http.Request) (*http.Response, error)
1. ActionGroupsClient.ListByResourceGroup(context.Context, string) (ActionGroupList, error)
1. ActionGroupsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. ActionGroupsClient.ListByResourceGroupResponder(*http.Response) (ActionGroupList, error)
1. ActionGroupsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. ActionGroupsClient.ListBySubscriptionID(context.Context) (ActionGroupList, error)
1. ActionGroupsClient.ListBySubscriptionIDPreparer(context.Context) (*http.Request, error)
1. ActionGroupsClient.ListBySubscriptionIDResponder(*http.Response) (ActionGroupList, error)
1. ActionGroupsClient.ListBySubscriptionIDSender(*http.Request) (*http.Response, error)
1. ActionGroupsClient.Update(context.Context, string, string, ActionGroupPatchBody) (ActionGroupResource, error)
1. ActionGroupsClient.UpdatePreparer(context.Context, string, string, ActionGroupPatchBody) (*http.Request, error)
1. ActionGroupsClient.UpdateResponder(*http.Response) (ActionGroupResource, error)
1. ActionGroupsClient.UpdateSender(*http.Request) (*http.Response, error)
1. ActivityLogAlertActionGroup.MarshalJSON() ([]byte, error)
1. ActivityLogAlertPatchBody.MarshalJSON() ([]byte, error)
1. ActivityLogAlertResource.MarshalJSON() ([]byte, error)
1. ActivityLogAlertsClient.CreateOrUpdate(context.Context, string, string, ActivityLogAlertResource) (ActivityLogAlertResource, error)
1. ActivityLogAlertsClient.CreateOrUpdatePreparer(context.Context, string, string, ActivityLogAlertResource) (*http.Request, error)
1. ActivityLogAlertsClient.CreateOrUpdateResponder(*http.Response) (ActivityLogAlertResource, error)
1. ActivityLogAlertsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ActivityLogAlertsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. ActivityLogAlertsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. ActivityLogAlertsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ActivityLogAlertsClient.DeleteSender(*http.Request) (*http.Response, error)
1. ActivityLogAlertsClient.Get(context.Context, string, string) (ActivityLogAlertResource, error)
1. ActivityLogAlertsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ActivityLogAlertsClient.GetResponder(*http.Response) (ActivityLogAlertResource, error)
1. ActivityLogAlertsClient.GetSender(*http.Request) (*http.Response, error)
1. ActivityLogAlertsClient.ListByResourceGroup(context.Context, string) (ActivityLogAlertList, error)
1. ActivityLogAlertsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. ActivityLogAlertsClient.ListByResourceGroupResponder(*http.Response) (ActivityLogAlertList, error)
1. ActivityLogAlertsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. ActivityLogAlertsClient.ListBySubscriptionID(context.Context) (ActivityLogAlertList, error)
1. ActivityLogAlertsClient.ListBySubscriptionIDPreparer(context.Context) (*http.Request, error)
1. ActivityLogAlertsClient.ListBySubscriptionIDResponder(*http.Response) (ActivityLogAlertList, error)
1. ActivityLogAlertsClient.ListBySubscriptionIDSender(*http.Request) (*http.Response, error)
1. ActivityLogAlertsClient.Update(context.Context, string, string, ActivityLogAlertPatchBody) (ActivityLogAlertResource, error)
1. ActivityLogAlertsClient.UpdatePreparer(context.Context, string, string, ActivityLogAlertPatchBody) (*http.Request, error)
1. ActivityLogAlertsClient.UpdateResponder(*http.Response) (ActivityLogAlertResource, error)
1. ActivityLogAlertsClient.UpdateSender(*http.Request) (*http.Response, error)
1. ActivityLogsClient.List(context.Context, string, string) (EventDataCollectionPage, error)
1. ActivityLogsClient.ListComplete(context.Context, string, string) (EventDataCollectionIterator, error)
1. ActivityLogsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. ActivityLogsClient.ListResponder(*http.Response) (EventDataCollection, error)
1. ActivityLogsClient.ListSender(*http.Request) (*http.Response, error)
1. AlertRule.MarshalJSON() ([]byte, error)
1. AlertRuleIncidentsClient.Get(context.Context, string, string, string) (Incident, error)
1. AlertRuleIncidentsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. AlertRuleIncidentsClient.GetResponder(*http.Response) (Incident, error)
1. AlertRuleIncidentsClient.GetSender(*http.Request) (*http.Response, error)
1. AlertRuleIncidentsClient.ListByAlertRule(context.Context, string, string) (IncidentListResult, error)
1. AlertRuleIncidentsClient.ListByAlertRulePreparer(context.Context, string, string) (*http.Request, error)
1. AlertRuleIncidentsClient.ListByAlertRuleResponder(*http.Response) (IncidentListResult, error)
1. AlertRuleIncidentsClient.ListByAlertRuleSender(*http.Request) (*http.Response, error)
1. AlertRuleResource.MarshalJSON() ([]byte, error)
1. AlertRuleResourcePatch.MarshalJSON() ([]byte, error)
1. AlertRulesClient.CreateOrUpdate(context.Context, string, string, AlertRuleResource) (AlertRuleResource, error)
1. AlertRulesClient.CreateOrUpdatePreparer(context.Context, string, string, AlertRuleResource) (*http.Request, error)
1. AlertRulesClient.CreateOrUpdateResponder(*http.Response) (AlertRuleResource, error)
1. AlertRulesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. AlertRulesClient.Delete(context.Context, string, string) (autorest.Response, error)
1. AlertRulesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. AlertRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. AlertRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. AlertRulesClient.Get(context.Context, string, string) (AlertRuleResource, error)
1. AlertRulesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. AlertRulesClient.GetResponder(*http.Response) (AlertRuleResource, error)
1. AlertRulesClient.GetSender(*http.Request) (*http.Response, error)
1. AlertRulesClient.ListByResourceGroup(context.Context, string) (AlertRuleResourceCollection, error)
1. AlertRulesClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. AlertRulesClient.ListByResourceGroupResponder(*http.Response) (AlertRuleResourceCollection, error)
1. AlertRulesClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. AlertRulesClient.ListBySubscription(context.Context) (AlertRuleResourceCollection, error)
1. AlertRulesClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. AlertRulesClient.ListBySubscriptionResponder(*http.Response) (AlertRuleResourceCollection, error)
1. AlertRulesClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. AlertRulesClient.Update(context.Context, string, string, AlertRuleResourcePatch) (AlertRuleResource, error)
1. AlertRulesClient.UpdatePreparer(context.Context, string, string, AlertRuleResourcePatch) (*http.Request, error)
1. AlertRulesClient.UpdateResponder(*http.Response) (AlertRuleResource, error)
1. AlertRulesClient.UpdateSender(*http.Request) (*http.Response, error)
1. AlertingAction.AsAction() (*Action, bool)
1. AlertingAction.AsAlertingAction() (*AlertingAction, bool)
1. AlertingAction.AsBasicAction() (BasicAction, bool)
1. AlertingAction.AsLogToMetricAction() (*LogToMetricAction, bool)
1. AlertingAction.MarshalJSON() ([]byte, error)
1. AutoscaleSettingResource.MarshalJSON() ([]byte, error)
1. AutoscaleSettingResourceCollection.IsEmpty() bool
1. AutoscaleSettingResourceCollectionIterator.NotDone() bool
1. AutoscaleSettingResourceCollectionIterator.Response() AutoscaleSettingResourceCollection
1. AutoscaleSettingResourceCollectionIterator.Value() AutoscaleSettingResource
1. AutoscaleSettingResourceCollectionPage.NotDone() bool
1. AutoscaleSettingResourceCollectionPage.Response() AutoscaleSettingResourceCollection
1. AutoscaleSettingResourceCollectionPage.Values() []AutoscaleSettingResource
1. AutoscaleSettingResourcePatch.MarshalJSON() ([]byte, error)
1. AutoscaleSettingsClient.CreateOrUpdate(context.Context, string, string, AutoscaleSettingResource) (AutoscaleSettingResource, error)
1. AutoscaleSettingsClient.CreateOrUpdatePreparer(context.Context, string, string, AutoscaleSettingResource) (*http.Request, error)
1. AutoscaleSettingsClient.CreateOrUpdateResponder(*http.Response) (AutoscaleSettingResource, error)
1. AutoscaleSettingsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. AutoscaleSettingsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. AutoscaleSettingsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. AutoscaleSettingsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. AutoscaleSettingsClient.DeleteSender(*http.Request) (*http.Response, error)
1. AutoscaleSettingsClient.Get(context.Context, string, string) (AutoscaleSettingResource, error)
1. AutoscaleSettingsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. AutoscaleSettingsClient.GetResponder(*http.Response) (AutoscaleSettingResource, error)
1. AutoscaleSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. AutoscaleSettingsClient.ListByResourceGroup(context.Context, string) (AutoscaleSettingResourceCollectionPage, error)
1. AutoscaleSettingsClient.ListByResourceGroupComplete(context.Context, string) (AutoscaleSettingResourceCollectionIterator, error)
1. AutoscaleSettingsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. AutoscaleSettingsClient.ListByResourceGroupResponder(*http.Response) (AutoscaleSettingResourceCollection, error)
1. AutoscaleSettingsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. AutoscaleSettingsClient.ListBySubscription(context.Context) (AutoscaleSettingResourceCollectionPage, error)
1. AutoscaleSettingsClient.ListBySubscriptionComplete(context.Context) (AutoscaleSettingResourceCollectionIterator, error)
1. AutoscaleSettingsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. AutoscaleSettingsClient.ListBySubscriptionResponder(*http.Response) (AutoscaleSettingResourceCollection, error)
1. AutoscaleSettingsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. AutoscaleSettingsClient.Update(context.Context, string, string, AutoscaleSettingResourcePatch) (AutoscaleSettingResource, error)
1. AutoscaleSettingsClient.UpdatePreparer(context.Context, string, string, AutoscaleSettingResourcePatch) (*http.Request, error)
1. AutoscaleSettingsClient.UpdateResponder(*http.Response) (AutoscaleSettingResource, error)
1. AutoscaleSettingsClient.UpdateSender(*http.Request) (*http.Response, error)
1. AzureMonitorPrivateLinkScope.MarshalJSON() ([]byte, error)
1. AzureMonitorPrivateLinkScopeListResult.IsEmpty() bool
1. AzureMonitorPrivateLinkScopeListResultIterator.NotDone() bool
1. AzureMonitorPrivateLinkScopeListResultIterator.Response() AzureMonitorPrivateLinkScopeListResult
1. AzureMonitorPrivateLinkScopeListResultIterator.Value() AzureMonitorPrivateLinkScope
1. AzureMonitorPrivateLinkScopeListResultPage.NotDone() bool
1. AzureMonitorPrivateLinkScopeListResultPage.Response() AzureMonitorPrivateLinkScopeListResult
1. AzureMonitorPrivateLinkScopeListResultPage.Values() []AzureMonitorPrivateLinkScope
1. BaselineResponse.MarshalJSON() ([]byte, error)
1. BaselinesClient.List(context.Context, string, string, string, string, *string, string, string, string, ResultType) (MetricBaselinesResponse, error)
1. BaselinesClient.ListPreparer(context.Context, string, string, string, string, *string, string, string, string, ResultType) (*http.Request, error)
1. BaselinesClient.ListResponder(*http.Response) (MetricBaselinesResponse, error)
1. BaselinesClient.ListSender(*http.Request) (*http.Response, error)
1. DiagnosticSettingsCategoryClient.Get(context.Context, string, string) (DiagnosticSettingsCategoryResource, error)
1. DiagnosticSettingsCategoryClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. DiagnosticSettingsCategoryClient.GetResponder(*http.Response) (DiagnosticSettingsCategoryResource, error)
1. DiagnosticSettingsCategoryClient.GetSender(*http.Request) (*http.Response, error)
1. DiagnosticSettingsCategoryClient.List(context.Context, string) (DiagnosticSettingsCategoryResourceCollection, error)
1. DiagnosticSettingsCategoryClient.ListPreparer(context.Context, string) (*http.Request, error)
1. DiagnosticSettingsCategoryClient.ListResponder(*http.Response) (DiagnosticSettingsCategoryResourceCollection, error)
1. DiagnosticSettingsCategoryClient.ListSender(*http.Request) (*http.Response, error)
1. DiagnosticSettingsCategoryResource.MarshalJSON() ([]byte, error)
1. DiagnosticSettingsClient.CreateOrUpdate(context.Context, string, DiagnosticSettingsResource, string) (DiagnosticSettingsResource, error)
1. DiagnosticSettingsClient.CreateOrUpdatePreparer(context.Context, string, DiagnosticSettingsResource, string) (*http.Request, error)
1. DiagnosticSettingsClient.CreateOrUpdateResponder(*http.Response) (DiagnosticSettingsResource, error)
1. DiagnosticSettingsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. DiagnosticSettingsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. DiagnosticSettingsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. DiagnosticSettingsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. DiagnosticSettingsClient.DeleteSender(*http.Request) (*http.Response, error)
1. DiagnosticSettingsClient.Get(context.Context, string, string) (DiagnosticSettingsResource, error)
1. DiagnosticSettingsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. DiagnosticSettingsClient.GetResponder(*http.Response) (DiagnosticSettingsResource, error)
1. DiagnosticSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. DiagnosticSettingsClient.List(context.Context, string) (DiagnosticSettingsResourceCollection, error)
1. DiagnosticSettingsClient.ListPreparer(context.Context, string) (*http.Request, error)
1. DiagnosticSettingsClient.ListResponder(*http.Response) (DiagnosticSettingsResourceCollection, error)
1. DiagnosticSettingsClient.ListSender(*http.Request) (*http.Response, error)
1. DiagnosticSettingsResource.MarshalJSON() ([]byte, error)
1. DynamicMetricCriteria.AsBasicMultiMetricCriteria() (BasicMultiMetricCriteria, bool)
1. DynamicMetricCriteria.AsDynamicMetricCriteria() (*DynamicMetricCriteria, bool)
1. DynamicMetricCriteria.AsMetricCriteria() (*MetricCriteria, bool)
1. DynamicMetricCriteria.AsMultiMetricCriteria() (*MultiMetricCriteria, bool)
1. DynamicMetricCriteria.MarshalJSON() ([]byte, error)
1. EmailReceiver.MarshalJSON() ([]byte, error)
1. ErrorResponseCommon.MarshalJSON() ([]byte, error)
1. EventCategoriesClient.List(context.Context) (EventCategoryCollection, error)
1. EventCategoriesClient.ListPreparer(context.Context) (*http.Request, error)
1. EventCategoriesClient.ListResponder(*http.Response) (EventCategoryCollection, error)
1. EventCategoriesClient.ListSender(*http.Request) (*http.Response, error)
1. EventData.MarshalJSON() ([]byte, error)
1. EventDataCollection.IsEmpty() bool
1. EventDataCollectionIterator.NotDone() bool
1. EventDataCollectionIterator.Response() EventDataCollection
1. EventDataCollectionIterator.Value() EventData
1. EventDataCollectionPage.NotDone() bool
1. EventDataCollectionPage.Response() EventDataCollection
1. EventDataCollectionPage.Values() []EventData
1. LocationThresholdRuleCondition.AsBasicRuleCondition() (BasicRuleCondition, bool)
1. LocationThresholdRuleCondition.AsLocationThresholdRuleCondition() (*LocationThresholdRuleCondition, bool)
1. LocationThresholdRuleCondition.AsManagementEventRuleCondition() (*ManagementEventRuleCondition, bool)
1. LocationThresholdRuleCondition.AsRuleCondition() (*RuleCondition, bool)
1. LocationThresholdRuleCondition.AsThresholdRuleCondition() (*ThresholdRuleCondition, bool)
1. LocationThresholdRuleCondition.MarshalJSON() ([]byte, error)
1. LogProfileResource.MarshalJSON() ([]byte, error)
1. LogProfileResourcePatch.MarshalJSON() ([]byte, error)
1. LogProfilesClient.CreateOrUpdate(context.Context, string, LogProfileResource) (LogProfileResource, error)
1. LogProfilesClient.CreateOrUpdatePreparer(context.Context, string, LogProfileResource) (*http.Request, error)
1. LogProfilesClient.CreateOrUpdateResponder(*http.Response) (LogProfileResource, error)
1. LogProfilesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. LogProfilesClient.Delete(context.Context, string) (autorest.Response, error)
1. LogProfilesClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. LogProfilesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. LogProfilesClient.DeleteSender(*http.Request) (*http.Response, error)
1. LogProfilesClient.Get(context.Context, string) (LogProfileResource, error)
1. LogProfilesClient.GetPreparer(context.Context, string) (*http.Request, error)
1. LogProfilesClient.GetResponder(*http.Response) (LogProfileResource, error)
1. LogProfilesClient.GetSender(*http.Request) (*http.Response, error)
1. LogProfilesClient.List(context.Context) (LogProfileCollection, error)
1. LogProfilesClient.ListPreparer(context.Context) (*http.Request, error)
1. LogProfilesClient.ListResponder(*http.Response) (LogProfileCollection, error)
1. LogProfilesClient.ListSender(*http.Request) (*http.Response, error)
1. LogProfilesClient.Update(context.Context, string, LogProfileResourcePatch) (LogProfileResource, error)
1. LogProfilesClient.UpdatePreparer(context.Context, string, LogProfileResourcePatch) (*http.Request, error)
1. LogProfilesClient.UpdateResponder(*http.Response) (LogProfileResource, error)
1. LogProfilesClient.UpdateSender(*http.Request) (*http.Response, error)
1. LogSearchRule.MarshalJSON() ([]byte, error)
1. LogSearchRuleResource.MarshalJSON() ([]byte, error)
1. LogSearchRuleResourcePatch.MarshalJSON() ([]byte, error)
1. LogToMetricAction.AsAction() (*Action, bool)
1. LogToMetricAction.AsAlertingAction() (*AlertingAction, bool)
1. LogToMetricAction.AsBasicAction() (BasicAction, bool)
1. LogToMetricAction.AsLogToMetricAction() (*LogToMetricAction, bool)
1. LogToMetricAction.MarshalJSON() ([]byte, error)
1. ManagementEventRuleCondition.AsBasicRuleCondition() (BasicRuleCondition, bool)
1. ManagementEventRuleCondition.AsLocationThresholdRuleCondition() (*LocationThresholdRuleCondition, bool)
1. ManagementEventRuleCondition.AsManagementEventRuleCondition() (*ManagementEventRuleCondition, bool)
1. ManagementEventRuleCondition.AsRuleCondition() (*RuleCondition, bool)
1. ManagementEventRuleCondition.AsThresholdRuleCondition() (*ThresholdRuleCondition, bool)
1. ManagementEventRuleCondition.MarshalJSON() ([]byte, error)
1. MetricAlertAction.MarshalJSON() ([]byte, error)
1. MetricAlertCriteria.AsBasicMetricAlertCriteria() (BasicMetricAlertCriteria, bool)
1. MetricAlertCriteria.AsMetricAlertCriteria() (*MetricAlertCriteria, bool)
1. MetricAlertCriteria.AsMetricAlertMultipleResourceMultipleMetricCriteria() (*MetricAlertMultipleResourceMultipleMetricCriteria, bool)
1. MetricAlertCriteria.AsMetricAlertSingleResourceMultipleMetricCriteria() (*MetricAlertSingleResourceMultipleMetricCriteria, bool)
1. MetricAlertCriteria.AsWebtestLocationAvailabilityCriteria() (*WebtestLocationAvailabilityCriteria, bool)
1. MetricAlertCriteria.MarshalJSON() ([]byte, error)
1. MetricAlertMultipleResourceMultipleMetricCriteria.AsBasicMetricAlertCriteria() (BasicMetricAlertCriteria, bool)
1. MetricAlertMultipleResourceMultipleMetricCriteria.AsMetricAlertCriteria() (*MetricAlertCriteria, bool)
1. MetricAlertMultipleResourceMultipleMetricCriteria.AsMetricAlertMultipleResourceMultipleMetricCriteria() (*MetricAlertMultipleResourceMultipleMetricCriteria, bool)
1. MetricAlertMultipleResourceMultipleMetricCriteria.AsMetricAlertSingleResourceMultipleMetricCriteria() (*MetricAlertSingleResourceMultipleMetricCriteria, bool)
1. MetricAlertMultipleResourceMultipleMetricCriteria.AsWebtestLocationAvailabilityCriteria() (*WebtestLocationAvailabilityCriteria, bool)
1. MetricAlertMultipleResourceMultipleMetricCriteria.MarshalJSON() ([]byte, error)
1. MetricAlertProperties.MarshalJSON() ([]byte, error)
1. MetricAlertPropertiesPatch.MarshalJSON() ([]byte, error)
1. MetricAlertResource.MarshalJSON() ([]byte, error)
1. MetricAlertResourcePatch.MarshalJSON() ([]byte, error)
1. MetricAlertSingleResourceMultipleMetricCriteria.AsBasicMetricAlertCriteria() (BasicMetricAlertCriteria, bool)
1. MetricAlertSingleResourceMultipleMetricCriteria.AsMetricAlertCriteria() (*MetricAlertCriteria, bool)
1. MetricAlertSingleResourceMultipleMetricCriteria.AsMetricAlertMultipleResourceMultipleMetricCriteria() (*MetricAlertMultipleResourceMultipleMetricCriteria, bool)
1. MetricAlertSingleResourceMultipleMetricCriteria.AsMetricAlertSingleResourceMultipleMetricCriteria() (*MetricAlertSingleResourceMultipleMetricCriteria, bool)
1. MetricAlertSingleResourceMultipleMetricCriteria.AsWebtestLocationAvailabilityCriteria() (*WebtestLocationAvailabilityCriteria, bool)
1. MetricAlertSingleResourceMultipleMetricCriteria.MarshalJSON() ([]byte, error)
1. MetricAlertStatusProperties.MarshalJSON() ([]byte, error)
1. MetricAlertsClient.CreateOrUpdate(context.Context, string, string, MetricAlertResource) (MetricAlertResource, error)
1. MetricAlertsClient.CreateOrUpdatePreparer(context.Context, string, string, MetricAlertResource) (*http.Request, error)
1. MetricAlertsClient.CreateOrUpdateResponder(*http.Response) (MetricAlertResource, error)
1. MetricAlertsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. MetricAlertsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. MetricAlertsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. MetricAlertsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. MetricAlertsClient.DeleteSender(*http.Request) (*http.Response, error)
1. MetricAlertsClient.Get(context.Context, string, string) (MetricAlertResource, error)
1. MetricAlertsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. MetricAlertsClient.GetResponder(*http.Response) (MetricAlertResource, error)
1. MetricAlertsClient.GetSender(*http.Request) (*http.Response, error)
1. MetricAlertsClient.ListByResourceGroup(context.Context, string) (MetricAlertResourceCollection, error)
1. MetricAlertsClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. MetricAlertsClient.ListByResourceGroupResponder(*http.Response) (MetricAlertResourceCollection, error)
1. MetricAlertsClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. MetricAlertsClient.ListBySubscription(context.Context) (MetricAlertResourceCollection, error)
1. MetricAlertsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. MetricAlertsClient.ListBySubscriptionResponder(*http.Response) (MetricAlertResourceCollection, error)
1. MetricAlertsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. MetricAlertsClient.Update(context.Context, string, string, MetricAlertResourcePatch) (MetricAlertResource, error)
1. MetricAlertsClient.UpdatePreparer(context.Context, string, string, MetricAlertResourcePatch) (*http.Request, error)
1. MetricAlertsClient.UpdateResponder(*http.Response) (MetricAlertResource, error)
1. MetricAlertsClient.UpdateSender(*http.Request) (*http.Response, error)
1. MetricAlertsStatusClient.List(context.Context, string, string) (MetricAlertStatusCollection, error)
1. MetricAlertsStatusClient.ListByName(context.Context, string, string, string) (MetricAlertStatusCollection, error)
1. MetricAlertsStatusClient.ListByNamePreparer(context.Context, string, string, string) (*http.Request, error)
1. MetricAlertsStatusClient.ListByNameResponder(*http.Response) (MetricAlertStatusCollection, error)
1. MetricAlertsStatusClient.ListByNameSender(*http.Request) (*http.Response, error)
1. MetricAlertsStatusClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. MetricAlertsStatusClient.ListResponder(*http.Response) (MetricAlertStatusCollection, error)
1. MetricAlertsStatusClient.ListSender(*http.Request) (*http.Response, error)
1. MetricBaselineClient.CalculateBaseline(context.Context, string, TimeSeriesInformation) (CalculateBaselineResponse, error)
1. MetricBaselineClient.CalculateBaselinePreparer(context.Context, string, TimeSeriesInformation) (*http.Request, error)
1. MetricBaselineClient.CalculateBaselineResponder(*http.Response) (CalculateBaselineResponse, error)
1. MetricBaselineClient.CalculateBaselineSender(*http.Request) (*http.Response, error)
1. MetricBaselineClient.Get(context.Context, string, string, string, *string, string, string, ResultType) (BaselineResponse, error)
1. MetricBaselineClient.GetPreparer(context.Context, string, string, string, *string, string, string, ResultType) (*http.Request, error)
1. MetricBaselineClient.GetResponder(*http.Response) (BaselineResponse, error)
1. MetricBaselineClient.GetSender(*http.Request) (*http.Response, error)
1. MetricCriteria.AsBasicMultiMetricCriteria() (BasicMultiMetricCriteria, bool)
1. MetricCriteria.AsDynamicMetricCriteria() (*DynamicMetricCriteria, bool)
1. MetricCriteria.AsMetricCriteria() (*MetricCriteria, bool)
1. MetricCriteria.AsMultiMetricCriteria() (*MultiMetricCriteria, bool)
1. MetricCriteria.MarshalJSON() ([]byte, error)
1. MetricDefinitionsClient.List(context.Context, string, string) (MetricDefinitionCollection, error)
1. MetricDefinitionsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. MetricDefinitionsClient.ListResponder(*http.Response) (MetricDefinitionCollection, error)
1. MetricDefinitionsClient.ListSender(*http.Request) (*http.Response, error)
1. MetricNamespacesClient.List(context.Context, string, string) (MetricNamespaceCollection, error)
1. MetricNamespacesClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. MetricNamespacesClient.ListResponder(*http.Response) (MetricNamespaceCollection, error)
1. MetricNamespacesClient.ListSender(*http.Request) (*http.Response, error)
1. MetricsClient.List(context.Context, string, string, *string, string, string, *int32, string, string, ResultType, string) (Response, error)
1. MetricsClient.ListPreparer(context.Context, string, string, *string, string, string, *int32, string, string, ResultType, string) (*http.Request, error)
1. MetricsClient.ListResponder(*http.Response) (Response, error)
1. MetricsClient.ListSender(*http.Request) (*http.Response, error)
1. MultiMetricCriteria.AsBasicMultiMetricCriteria() (BasicMultiMetricCriteria, bool)
1. MultiMetricCriteria.AsDynamicMetricCriteria() (*DynamicMetricCriteria, bool)
1. MultiMetricCriteria.AsMetricCriteria() (*MetricCriteria, bool)
1. MultiMetricCriteria.AsMultiMetricCriteria() (*MultiMetricCriteria, bool)
1. MultiMetricCriteria.MarshalJSON() ([]byte, error)
1. NewActionGroupsClient(string) ActionGroupsClient
1. NewActionGroupsClientWithBaseURI(string, string) ActionGroupsClient
1. NewActivityLogAlertsClient(string) ActivityLogAlertsClient
1. NewActivityLogAlertsClientWithBaseURI(string, string) ActivityLogAlertsClient
1. NewActivityLogsClient(string) ActivityLogsClient
1. NewActivityLogsClientWithBaseURI(string, string) ActivityLogsClient
1. NewAlertRuleIncidentsClient(string) AlertRuleIncidentsClient
1. NewAlertRuleIncidentsClientWithBaseURI(string, string) AlertRuleIncidentsClient
1. NewAlertRulesClient(string) AlertRulesClient
1. NewAlertRulesClientWithBaseURI(string, string) AlertRulesClient
1. NewAutoscaleSettingResourceCollectionIterator(AutoscaleSettingResourceCollectionPage) AutoscaleSettingResourceCollectionIterator
1. NewAutoscaleSettingResourceCollectionPage(AutoscaleSettingResourceCollection, func(context.Context, AutoscaleSettingResourceCollection) (AutoscaleSettingResourceCollection, error)) AutoscaleSettingResourceCollectionPage
1. NewAutoscaleSettingsClient(string) AutoscaleSettingsClient
1. NewAutoscaleSettingsClientWithBaseURI(string, string) AutoscaleSettingsClient
1. NewAzureMonitorPrivateLinkScopeListResultIterator(AzureMonitorPrivateLinkScopeListResultPage) AzureMonitorPrivateLinkScopeListResultIterator
1. NewAzureMonitorPrivateLinkScopeListResultPage(AzureMonitorPrivateLinkScopeListResult, func(context.Context, AzureMonitorPrivateLinkScopeListResult) (AzureMonitorPrivateLinkScopeListResult, error)) AzureMonitorPrivateLinkScopeListResultPage
1. NewBaselinesClient(string) BaselinesClient
1. NewBaselinesClientWithBaseURI(string, string) BaselinesClient
1. NewDiagnosticSettingsCategoryClient(string) DiagnosticSettingsCategoryClient
1. NewDiagnosticSettingsCategoryClientWithBaseURI(string, string) DiagnosticSettingsCategoryClient
1. NewDiagnosticSettingsClient(string) DiagnosticSettingsClient
1. NewDiagnosticSettingsClientWithBaseURI(string, string) DiagnosticSettingsClient
1. NewEventCategoriesClient(string) EventCategoriesClient
1. NewEventCategoriesClientWithBaseURI(string, string) EventCategoriesClient
1. NewEventDataCollectionIterator(EventDataCollectionPage) EventDataCollectionIterator
1. NewEventDataCollectionPage(EventDataCollection, func(context.Context, EventDataCollection) (EventDataCollection, error)) EventDataCollectionPage
1. NewLogProfilesClient(string) LogProfilesClient
1. NewLogProfilesClientWithBaseURI(string, string) LogProfilesClient
1. NewMetricAlertsClient(string) MetricAlertsClient
1. NewMetricAlertsClientWithBaseURI(string, string) MetricAlertsClient
1. NewMetricAlertsStatusClient(string) MetricAlertsStatusClient
1. NewMetricAlertsStatusClientWithBaseURI(string, string) MetricAlertsStatusClient
1. NewMetricBaselineClient(string) MetricBaselineClient
1. NewMetricBaselineClientWithBaseURI(string, string) MetricBaselineClient
1. NewMetricDefinitionsClient(string) MetricDefinitionsClient
1. NewMetricDefinitionsClientWithBaseURI(string, string) MetricDefinitionsClient
1. NewMetricNamespacesClient(string) MetricNamespacesClient
1. NewMetricNamespacesClientWithBaseURI(string, string) MetricNamespacesClient
1. NewMetricsClient(string) MetricsClient
1. NewMetricsClientWithBaseURI(string, string) MetricsClient
1. NewOperationsClient(string) OperationsClient
1. NewOperationsClientWithBaseURI(string, string) OperationsClient
1. NewPrivateEndpointConnectionListResultIterator(PrivateEndpointConnectionListResultPage) PrivateEndpointConnectionListResultIterator
1. NewPrivateEndpointConnectionListResultPage(PrivateEndpointConnectionListResult, func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error)) PrivateEndpointConnectionListResultPage
1. NewPrivateEndpointConnectionsClient(string) PrivateEndpointConnectionsClient
1. NewPrivateEndpointConnectionsClientWithBaseURI(string, string) PrivateEndpointConnectionsClient
1. NewPrivateLinkResourceListResultIterator(PrivateLinkResourceListResultPage) PrivateLinkResourceListResultIterator
1. NewPrivateLinkResourceListResultPage(PrivateLinkResourceListResult, func(context.Context, PrivateLinkResourceListResult) (PrivateLinkResourceListResult, error)) PrivateLinkResourceListResultPage
1. NewPrivateLinkResourcesClient(string) PrivateLinkResourcesClient
1. NewPrivateLinkResourcesClientWithBaseURI(string, string) PrivateLinkResourcesClient
1. NewPrivateLinkScopeOperationStatusClient(string) PrivateLinkScopeOperationStatusClient
1. NewPrivateLinkScopeOperationStatusClientWithBaseURI(string, string) PrivateLinkScopeOperationStatusClient
1. NewPrivateLinkScopedResourcesClient(string) PrivateLinkScopedResourcesClient
1. NewPrivateLinkScopedResourcesClientWithBaseURI(string, string) PrivateLinkScopedResourcesClient
1. NewPrivateLinkScopesClient(string) PrivateLinkScopesClient
1. NewPrivateLinkScopesClientWithBaseURI(string, string) PrivateLinkScopesClient
1. NewScheduledQueryRulesClient(string) ScheduledQueryRulesClient
1. NewScheduledQueryRulesClientWithBaseURI(string, string) ScheduledQueryRulesClient
1. NewScopedResourceListResultIterator(ScopedResourceListResultPage) ScopedResourceListResultIterator
1. NewScopedResourceListResultPage(ScopedResourceListResult, func(context.Context, ScopedResourceListResult) (ScopedResourceListResult, error)) ScopedResourceListResultPage
1. NewSubscriptionDiagnosticSettingsClient(string) SubscriptionDiagnosticSettingsClient
1. NewSubscriptionDiagnosticSettingsClientWithBaseURI(string, string) SubscriptionDiagnosticSettingsClient
1. NewTenantActivityLogsClient(string) TenantActivityLogsClient
1. NewTenantActivityLogsClientWithBaseURI(string, string) TenantActivityLogsClient
1. NewVMInsightsClient(string) VMInsightsClient
1. NewVMInsightsClientWithBaseURI(string, string) VMInsightsClient
1. OperationsClient.List(context.Context) (OperationListResult, error)
1. OperationsClient.ListPreparer(context.Context) (*http.Request, error)
1. OperationsClient.ListResponder(*http.Response) (OperationListResult, error)
1. OperationsClient.ListSender(*http.Request) (*http.Response, error)
1. PossibleAggregationTypeValues() []AggregationType
1. PossibleAlertSeverityValues() []AlertSeverity
1. PossibleBaselineSensitivityValues() []BaselineSensitivity
1. PossibleCategoryTypeValues() []CategoryType
1. PossibleComparisonOperationTypeValues() []ComparisonOperationType
1. PossibleConditionOperatorValues() []ConditionOperator
1. PossibleConditionalOperatorValues() []ConditionalOperator
1. PossibleCriterionTypeValues() []CriterionType
1. PossibleDataStatusValues() []DataStatus
1. PossibleDynamicThresholdOperatorValues() []DynamicThresholdOperator
1. PossibleDynamicThresholdSensitivityValues() []DynamicThresholdSensitivity
1. PossibleEnabledValues() []Enabled
1. PossibleEventLevelValues() []EventLevel
1. PossibleMetricStatisticTypeValues() []MetricStatisticType
1. PossibleMetricTriggerTypeValues() []MetricTriggerType
1. PossibleOdataTypeBasicActionValues() []OdataTypeBasicAction
1. PossibleOdataTypeBasicMetricAlertCriteriaValues() []OdataTypeBasicMetricAlertCriteria
1. PossibleOdataTypeBasicRuleActionValues() []OdataTypeBasicRuleAction
1. PossibleOdataTypeBasicRuleConditionValues() []OdataTypeBasicRuleCondition
1. PossibleOdataTypeValues() []OdataType
1. PossibleOnboardingStatusValues() []OnboardingStatus
1. PossibleOperatorValues() []Operator
1. PossibleProvisioningStateValues() []ProvisioningState
1. PossibleQueryTypeValues() []QueryType
1. PossibleReceiverStatusValues() []ReceiverStatus
1. PossibleRecurrenceFrequencyValues() []RecurrenceFrequency
1. PossibleResultTypeValues() []ResultType
1. PossibleScaleDirectionValues() []ScaleDirection
1. PossibleScaleRuleMetricDimensionOperationTypeValues() []ScaleRuleMetricDimensionOperationType
1. PossibleScaleTypeValues() []ScaleType
1. PossibleSensitivityValues() []Sensitivity
1. PossibleTimeAggregationOperatorValues() []TimeAggregationOperator
1. PossibleTimeAggregationTypeValues() []TimeAggregationType
1. PossibleUnitValues() []Unit
1. PrivateEndpointConnection.MarshalJSON() ([]byte, error)
1. PrivateEndpointConnectionListResult.IsEmpty() bool
1. PrivateEndpointConnectionListResultIterator.NotDone() bool
1. PrivateEndpointConnectionListResultIterator.Response() PrivateEndpointConnectionListResult
1. PrivateEndpointConnectionListResultIterator.Value() PrivateEndpointConnection
1. PrivateEndpointConnectionListResultPage.NotDone() bool
1. PrivateEndpointConnectionListResultPage.Response() PrivateEndpointConnectionListResult
1. PrivateEndpointConnectionListResultPage.Values() []PrivateEndpointConnection
1. PrivateEndpointConnectionProperties.MarshalJSON() ([]byte, error)
1. PrivateEndpointConnectionsClient.CreateOrUpdate(context.Context, string, string, string, PrivateEndpointConnection) (PrivateEndpointConnectionsCreateOrUpdateFuture, error)
1. PrivateEndpointConnectionsClient.CreateOrUpdatePreparer(context.Context, string, string, string, PrivateEndpointConnection) (*http.Request, error)
1. PrivateEndpointConnectionsClient.CreateOrUpdateResponder(*http.Response) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.CreateOrUpdateSender(*http.Request) (PrivateEndpointConnectionsCreateOrUpdateFuture, error)
1. PrivateEndpointConnectionsClient.Delete(context.Context, string, string, string) (PrivateEndpointConnectionsDeleteFuture, error)
1. PrivateEndpointConnectionsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. PrivateEndpointConnectionsClient.DeleteSender(*http.Request) (PrivateEndpointConnectionsDeleteFuture, error)
1. PrivateEndpointConnectionsClient.Get(context.Context, string, string, string) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.GetResponder(*http.Response) (PrivateEndpointConnection, error)
1. PrivateEndpointConnectionsClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateEndpointConnectionsClient.ListByPrivateLinkScope(context.Context, string, string) (PrivateEndpointConnectionListResultPage, error)
1. PrivateEndpointConnectionsClient.ListByPrivateLinkScopeComplete(context.Context, string, string) (PrivateEndpointConnectionListResultIterator, error)
1. PrivateEndpointConnectionsClient.ListByPrivateLinkScopePreparer(context.Context, string, string) (*http.Request, error)
1. PrivateEndpointConnectionsClient.ListByPrivateLinkScopeResponder(*http.Response) (PrivateEndpointConnectionListResult, error)
1. PrivateEndpointConnectionsClient.ListByPrivateLinkScopeSender(*http.Request) (*http.Response, error)
1. PrivateLinkResource.MarshalJSON() ([]byte, error)
1. PrivateLinkResourceListResult.IsEmpty() bool
1. PrivateLinkResourceListResultIterator.NotDone() bool
1. PrivateLinkResourceListResultIterator.Response() PrivateLinkResourceListResult
1. PrivateLinkResourceListResultIterator.Value() PrivateLinkResource
1. PrivateLinkResourceListResultPage.NotDone() bool
1. PrivateLinkResourceListResultPage.Response() PrivateLinkResourceListResult
1. PrivateLinkResourceListResultPage.Values() []PrivateLinkResource
1. PrivateLinkResourcesClient.Get(context.Context, string, string, string) (PrivateLinkResource, error)
1. PrivateLinkResourcesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateLinkResourcesClient.GetResponder(*http.Response) (PrivateLinkResource, error)
1. PrivateLinkResourcesClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateLinkResourcesClient.ListByPrivateLinkScope(context.Context, string, string) (PrivateLinkResourceListResultPage, error)
1. PrivateLinkResourcesClient.ListByPrivateLinkScopeComplete(context.Context, string, string) (PrivateLinkResourceListResultIterator, error)
1. PrivateLinkResourcesClient.ListByPrivateLinkScopePreparer(context.Context, string, string) (*http.Request, error)
1. PrivateLinkResourcesClient.ListByPrivateLinkScopeResponder(*http.Response) (PrivateLinkResourceListResult, error)
1. PrivateLinkResourcesClient.ListByPrivateLinkScopeSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopeOperationStatusClient.Get(context.Context, string, string) (OperationStatus, error)
1. PrivateLinkScopeOperationStatusClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. PrivateLinkScopeOperationStatusClient.GetResponder(*http.Response) (OperationStatus, error)
1. PrivateLinkScopeOperationStatusClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopedResourcesClient.CreateOrUpdate(context.Context, string, string, string, ScopedResource) (PrivateLinkScopedResourcesCreateOrUpdateFuture, error)
1. PrivateLinkScopedResourcesClient.CreateOrUpdatePreparer(context.Context, string, string, string, ScopedResource) (*http.Request, error)
1. PrivateLinkScopedResourcesClient.CreateOrUpdateResponder(*http.Response) (ScopedResource, error)
1. PrivateLinkScopedResourcesClient.CreateOrUpdateSender(*http.Request) (PrivateLinkScopedResourcesCreateOrUpdateFuture, error)
1. PrivateLinkScopedResourcesClient.Delete(context.Context, string, string, string) (PrivateLinkScopedResourcesDeleteFuture, error)
1. PrivateLinkScopedResourcesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateLinkScopedResourcesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. PrivateLinkScopedResourcesClient.DeleteSender(*http.Request) (PrivateLinkScopedResourcesDeleteFuture, error)
1. PrivateLinkScopedResourcesClient.Get(context.Context, string, string, string) (ScopedResource, error)
1. PrivateLinkScopedResourcesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. PrivateLinkScopedResourcesClient.GetResponder(*http.Response) (ScopedResource, error)
1. PrivateLinkScopedResourcesClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopedResourcesClient.ListByPrivateLinkScope(context.Context, string, string) (ScopedResourceListResultPage, error)
1. PrivateLinkScopedResourcesClient.ListByPrivateLinkScopeComplete(context.Context, string, string) (ScopedResourceListResultIterator, error)
1. PrivateLinkScopedResourcesClient.ListByPrivateLinkScopePreparer(context.Context, string, string) (*http.Request, error)
1. PrivateLinkScopedResourcesClient.ListByPrivateLinkScopeResponder(*http.Response) (ScopedResourceListResult, error)
1. PrivateLinkScopedResourcesClient.ListByPrivateLinkScopeSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopesClient.CreateOrUpdate(context.Context, string, string, AzureMonitorPrivateLinkScope) (AzureMonitorPrivateLinkScope, error)
1. PrivateLinkScopesClient.CreateOrUpdatePreparer(context.Context, string, string, AzureMonitorPrivateLinkScope) (*http.Request, error)
1. PrivateLinkScopesClient.CreateOrUpdateResponder(*http.Response) (AzureMonitorPrivateLinkScope, error)
1. PrivateLinkScopesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopesClient.Delete(context.Context, string, string) (PrivateLinkScopesDeleteFuture, error)
1. PrivateLinkScopesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. PrivateLinkScopesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. PrivateLinkScopesClient.DeleteSender(*http.Request) (PrivateLinkScopesDeleteFuture, error)
1. PrivateLinkScopesClient.Get(context.Context, string, string) (AzureMonitorPrivateLinkScope, error)
1. PrivateLinkScopesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. PrivateLinkScopesClient.GetResponder(*http.Response) (AzureMonitorPrivateLinkScope, error)
1. PrivateLinkScopesClient.GetSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopesClient.List(context.Context) (AzureMonitorPrivateLinkScopeListResultPage, error)
1. PrivateLinkScopesClient.ListByResourceGroup(context.Context, string) (AzureMonitorPrivateLinkScopeListResultPage, error)
1. PrivateLinkScopesClient.ListByResourceGroupComplete(context.Context, string) (AzureMonitorPrivateLinkScopeListResultIterator, error)
1. PrivateLinkScopesClient.ListByResourceGroupPreparer(context.Context, string) (*http.Request, error)
1. PrivateLinkScopesClient.ListByResourceGroupResponder(*http.Response) (AzureMonitorPrivateLinkScopeListResult, error)
1. PrivateLinkScopesClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopesClient.ListComplete(context.Context) (AzureMonitorPrivateLinkScopeListResultIterator, error)
1. PrivateLinkScopesClient.ListPreparer(context.Context) (*http.Request, error)
1. PrivateLinkScopesClient.ListResponder(*http.Response) (AzureMonitorPrivateLinkScopeListResult, error)
1. PrivateLinkScopesClient.ListSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopesClient.UpdateTags(context.Context, string, string, TagsResource) (AzureMonitorPrivateLinkScope, error)
1. PrivateLinkScopesClient.UpdateTagsPreparer(context.Context, string, string, TagsResource) (*http.Request, error)
1. PrivateLinkScopesClient.UpdateTagsResponder(*http.Response) (AzureMonitorPrivateLinkScope, error)
1. PrivateLinkScopesClient.UpdateTagsSender(*http.Request) (*http.Response, error)
1. PrivateLinkScopesResource.MarshalJSON() ([]byte, error)
1. PrivateLinkServiceConnectionStateProperty.MarshalJSON() ([]byte, error)
1. Resource.MarshalJSON() ([]byte, error)
1. RuleAction.AsBasicRuleAction() (BasicRuleAction, bool)
1. RuleAction.AsRuleAction() (*RuleAction, bool)
1. RuleAction.AsRuleEmailAction() (*RuleEmailAction, bool)
1. RuleAction.AsRuleWebhookAction() (*RuleWebhookAction, bool)
1. RuleAction.MarshalJSON() ([]byte, error)
1. RuleCondition.AsBasicRuleCondition() (BasicRuleCondition, bool)
1. RuleCondition.AsLocationThresholdRuleCondition() (*LocationThresholdRuleCondition, bool)
1. RuleCondition.AsManagementEventRuleCondition() (*ManagementEventRuleCondition, bool)
1. RuleCondition.AsRuleCondition() (*RuleCondition, bool)
1. RuleCondition.AsThresholdRuleCondition() (*ThresholdRuleCondition, bool)
1. RuleCondition.MarshalJSON() ([]byte, error)
1. RuleDataSource.AsBasicRuleDataSource() (BasicRuleDataSource, bool)
1. RuleDataSource.AsRuleDataSource() (*RuleDataSource, bool)
1. RuleDataSource.AsRuleManagementEventDataSource() (*RuleManagementEventDataSource, bool)
1. RuleDataSource.AsRuleMetricDataSource() (*RuleMetricDataSource, bool)
1. RuleDataSource.MarshalJSON() ([]byte, error)
1. RuleEmailAction.AsBasicRuleAction() (BasicRuleAction, bool)
1. RuleEmailAction.AsRuleAction() (*RuleAction, bool)
1. RuleEmailAction.AsRuleEmailAction() (*RuleEmailAction, bool)
1. RuleEmailAction.AsRuleWebhookAction() (*RuleWebhookAction, bool)
1. RuleEmailAction.MarshalJSON() ([]byte, error)
1. RuleManagementEventDataSource.AsBasicRuleDataSource() (BasicRuleDataSource, bool)
1. RuleManagementEventDataSource.AsRuleDataSource() (*RuleDataSource, bool)
1. RuleManagementEventDataSource.AsRuleManagementEventDataSource() (*RuleManagementEventDataSource, bool)
1. RuleManagementEventDataSource.AsRuleMetricDataSource() (*RuleMetricDataSource, bool)
1. RuleManagementEventDataSource.MarshalJSON() ([]byte, error)
1. RuleMetricDataSource.AsBasicRuleDataSource() (BasicRuleDataSource, bool)
1. RuleMetricDataSource.AsRuleDataSource() (*RuleDataSource, bool)
1. RuleMetricDataSource.AsRuleManagementEventDataSource() (*RuleManagementEventDataSource, bool)
1. RuleMetricDataSource.AsRuleMetricDataSource() (*RuleMetricDataSource, bool)
1. RuleMetricDataSource.MarshalJSON() ([]byte, error)
1. RuleWebhookAction.AsBasicRuleAction() (BasicRuleAction, bool)
1. RuleWebhookAction.AsRuleAction() (*RuleAction, bool)
1. RuleWebhookAction.AsRuleEmailAction() (*RuleEmailAction, bool)
1. RuleWebhookAction.AsRuleWebhookAction() (*RuleWebhookAction, bool)
1. RuleWebhookAction.MarshalJSON() ([]byte, error)
1. ScheduledQueryRulesClient.CreateOrUpdate(context.Context, string, string, LogSearchRuleResource) (LogSearchRuleResource, error)
1. ScheduledQueryRulesClient.CreateOrUpdatePreparer(context.Context, string, string, LogSearchRuleResource) (*http.Request, error)
1. ScheduledQueryRulesClient.CreateOrUpdateResponder(*http.Response) (LogSearchRuleResource, error)
1. ScheduledQueryRulesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ScheduledQueryRulesClient.Delete(context.Context, string, string) (autorest.Response, error)
1. ScheduledQueryRulesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. ScheduledQueryRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ScheduledQueryRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. ScheduledQueryRulesClient.Get(context.Context, string, string) (LogSearchRuleResource, error)
1. ScheduledQueryRulesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ScheduledQueryRulesClient.GetResponder(*http.Response) (LogSearchRuleResource, error)
1. ScheduledQueryRulesClient.GetSender(*http.Request) (*http.Response, error)
1. ScheduledQueryRulesClient.ListByResourceGroup(context.Context, string, string) (LogSearchRuleResourceCollection, error)
1. ScheduledQueryRulesClient.ListByResourceGroupPreparer(context.Context, string, string) (*http.Request, error)
1. ScheduledQueryRulesClient.ListByResourceGroupResponder(*http.Response) (LogSearchRuleResourceCollection, error)
1. ScheduledQueryRulesClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)
1. ScheduledQueryRulesClient.ListBySubscription(context.Context, string) (LogSearchRuleResourceCollection, error)
1. ScheduledQueryRulesClient.ListBySubscriptionPreparer(context.Context, string) (*http.Request, error)
1. ScheduledQueryRulesClient.ListBySubscriptionResponder(*http.Response) (LogSearchRuleResourceCollection, error)
1. ScheduledQueryRulesClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. ScheduledQueryRulesClient.Update(context.Context, string, string, LogSearchRuleResourcePatch) (LogSearchRuleResource, error)
1. ScheduledQueryRulesClient.UpdatePreparer(context.Context, string, string, LogSearchRuleResourcePatch) (*http.Request, error)
1. ScheduledQueryRulesClient.UpdateResponder(*http.Response) (LogSearchRuleResource, error)
1. ScheduledQueryRulesClient.UpdateSender(*http.Request) (*http.Response, error)
1. ScopedResource.MarshalJSON() ([]byte, error)
1. ScopedResourceListResult.IsEmpty() bool
1. ScopedResourceListResultIterator.NotDone() bool
1. ScopedResourceListResultIterator.Response() ScopedResourceListResult
1. ScopedResourceListResultIterator.Value() ScopedResource
1. ScopedResourceListResultPage.NotDone() bool
1. ScopedResourceListResultPage.Response() ScopedResourceListResult
1. ScopedResourceListResultPage.Values() []ScopedResource
1. ScopedResourceProperties.MarshalJSON() ([]byte, error)
1. SingleMetricBaseline.MarshalJSON() ([]byte, error)
1. SmsReceiver.MarshalJSON() ([]byte, error)
1. SubscriptionDiagnosticSettingsClient.CreateOrUpdate(context.Context, SubscriptionDiagnosticSettingsResource, string) (SubscriptionDiagnosticSettingsResource, error)
1. SubscriptionDiagnosticSettingsClient.CreateOrUpdatePreparer(context.Context, SubscriptionDiagnosticSettingsResource, string) (*http.Request, error)
1. SubscriptionDiagnosticSettingsClient.CreateOrUpdateResponder(*http.Response) (SubscriptionDiagnosticSettingsResource, error)
1. SubscriptionDiagnosticSettingsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. SubscriptionDiagnosticSettingsClient.Delete(context.Context, string) (autorest.Response, error)
1. SubscriptionDiagnosticSettingsClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. SubscriptionDiagnosticSettingsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. SubscriptionDiagnosticSettingsClient.DeleteSender(*http.Request) (*http.Response, error)
1. SubscriptionDiagnosticSettingsClient.Get(context.Context, string) (SubscriptionDiagnosticSettingsResource, error)
1. SubscriptionDiagnosticSettingsClient.GetPreparer(context.Context, string) (*http.Request, error)
1. SubscriptionDiagnosticSettingsClient.GetResponder(*http.Response) (SubscriptionDiagnosticSettingsResource, error)
1. SubscriptionDiagnosticSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. SubscriptionDiagnosticSettingsClient.List(context.Context) (SubscriptionDiagnosticSettingsResourceCollection, error)
1. SubscriptionDiagnosticSettingsClient.ListPreparer(context.Context) (*http.Request, error)
1. SubscriptionDiagnosticSettingsClient.ListResponder(*http.Response) (SubscriptionDiagnosticSettingsResourceCollection, error)
1. SubscriptionDiagnosticSettingsClient.ListSender(*http.Request) (*http.Response, error)
1. SubscriptionDiagnosticSettingsResource.MarshalJSON() ([]byte, error)
1. SubscriptionProxyOnlyResource.MarshalJSON() ([]byte, error)
1. TagsResource.MarshalJSON() ([]byte, error)
1. TenantActivityLogsClient.List(context.Context, string, string) (EventDataCollectionPage, error)
1. TenantActivityLogsClient.ListComplete(context.Context, string, string) (EventDataCollectionIterator, error)
1. TenantActivityLogsClient.ListPreparer(context.Context, string, string) (*http.Request, error)
1. TenantActivityLogsClient.ListResponder(*http.Response) (EventDataCollection, error)
1. TenantActivityLogsClient.ListSender(*http.Request) (*http.Response, error)
1. ThresholdRuleCondition.AsBasicRuleCondition() (BasicRuleCondition, bool)
1. ThresholdRuleCondition.AsLocationThresholdRuleCondition() (*LocationThresholdRuleCondition, bool)
1. ThresholdRuleCondition.AsManagementEventRuleCondition() (*ManagementEventRuleCondition, bool)
1. ThresholdRuleCondition.AsRuleCondition() (*RuleCondition, bool)
1. ThresholdRuleCondition.AsThresholdRuleCondition() (*ThresholdRuleCondition, bool)
1. ThresholdRuleCondition.MarshalJSON() ([]byte, error)
1. VMInsightsClient.GetOnboardingStatus(context.Context, string) (VMInsightsOnboardingStatus, error)
1. VMInsightsClient.GetOnboardingStatusPreparer(context.Context, string) (*http.Request, error)
1. VMInsightsClient.GetOnboardingStatusResponder(*http.Response) (VMInsightsOnboardingStatus, error)
1. VMInsightsClient.GetOnboardingStatusSender(*http.Request) (*http.Response, error)
1. VMInsightsOnboardingStatus.MarshalJSON() ([]byte, error)
1. WebhookNotification.MarshalJSON() ([]byte, error)
1. WebtestLocationAvailabilityCriteria.AsBasicMetricAlertCriteria() (BasicMetricAlertCriteria, bool)
1. WebtestLocationAvailabilityCriteria.AsMetricAlertCriteria() (*MetricAlertCriteria, bool)
1. WebtestLocationAvailabilityCriteria.AsMetricAlertMultipleResourceMultipleMetricCriteria() (*MetricAlertMultipleResourceMultipleMetricCriteria, bool)
1. WebtestLocationAvailabilityCriteria.AsMetricAlertSingleResourceMultipleMetricCriteria() (*MetricAlertSingleResourceMultipleMetricCriteria, bool)
1. WebtestLocationAvailabilityCriteria.AsWebtestLocationAvailabilityCriteria() (*WebtestLocationAvailabilityCriteria, bool)
1. WebtestLocationAvailabilityCriteria.MarshalJSON() ([]byte, error)
1. WorkspaceInfo.MarshalJSON() ([]byte, error)

## Struct Changes

### New Structs

1. Action
1. ActionGroup
1. ActionGroupList
1. ActionGroupPatch
1. ActionGroupPatchBody
1. ActionGroupResource
1. ActionGroupsClient
1. ActivityLogAlert
1. ActivityLogAlertActionGroup
1. ActivityLogAlertActionList
1. ActivityLogAlertAllOfCondition
1. ActivityLogAlertLeafCondition
1. ActivityLogAlertList
1. ActivityLogAlertPatch
1. ActivityLogAlertPatchBody
1. ActivityLogAlertResource
1. ActivityLogAlertsClient
1. ActivityLogsClient
1. AlertRule
1. AlertRuleIncidentsClient
1. AlertRuleResource
1. AlertRuleResourceCollection
1. AlertRuleResourcePatch
1. AlertRulesClient
1. AlertingAction
1. ArmRoleReceiver
1. AutomationRunbookReceiver
1. AutoscaleNotification
1. AutoscaleProfile
1. AutoscaleSetting
1. AutoscaleSettingResource
1. AutoscaleSettingResourceCollection
1. AutoscaleSettingResourceCollectionIterator
1. AutoscaleSettingResourceCollectionPage
1. AutoscaleSettingResourcePatch
1. AutoscaleSettingsClient
1. AzNsActionGroup
1. AzureAppPushReceiver
1. AzureFunctionReceiver
1. AzureMonitorPrivateLinkScope
1. AzureMonitorPrivateLinkScopeListResult
1. AzureMonitorPrivateLinkScopeListResultIterator
1. AzureMonitorPrivateLinkScopeListResultPage
1. AzureMonitorPrivateLinkScopeProperties
1. Baseline
1. BaselineMetadata
1. BaselineMetadataValue
1. BaselineProperties
1. BaselineResponse
1. BaselinesClient
1. CalculateBaselineResponse
1. Criteria
1. DataContainer
1. DiagnosticSettings
1. DiagnosticSettingsCategory
1. DiagnosticSettingsCategoryClient
1. DiagnosticSettingsCategoryResource
1. DiagnosticSettingsCategoryResourceCollection
1. DiagnosticSettingsClient
1. DiagnosticSettingsResource
1. DiagnosticSettingsResourceCollection
1. Dimension
1. DynamicMetricCriteria
1. DynamicThresholdFailingPeriods
1. EmailNotification
1. EmailReceiver
1. EnableRequest
1. Error
1. ErrorResponseCommon
1. EventCategoriesClient
1. EventCategoryCollection
1. EventData
1. EventDataCollection
1. EventDataCollectionIterator
1. EventDataCollectionPage
1. HTTPRequestInfo
1. Incident
1. IncidentListResult
1. ItsmReceiver
1. LocalizableString
1. LocationThresholdRuleCondition
1. LogMetricTrigger
1. LogProfileCollection
1. LogProfileProperties
1. LogProfileResource
1. LogProfileResourcePatch
1. LogProfilesClient
1. LogSearchRule
1. LogSearchRulePatch
1. LogSearchRuleResource
1. LogSearchRuleResourceCollection
1. LogSearchRuleResourcePatch
1. LogSettings
1. LogToMetricAction
1. LogicAppReceiver
1. ManagementEventAggregationCondition
1. ManagementEventRuleCondition
1. MetadataValue
1. Metric
1. MetricAlertAction
1. MetricAlertCriteria
1. MetricAlertMultipleResourceMultipleMetricCriteria
1. MetricAlertProperties
1. MetricAlertPropertiesPatch
1. MetricAlertResource
1. MetricAlertResourceCollection
1. MetricAlertResourcePatch
1. MetricAlertSingleResourceMultipleMetricCriteria
1. MetricAlertStatus
1. MetricAlertStatusCollection
1. MetricAlertStatusProperties
1. MetricAlertsClient
1. MetricAlertsStatusClient
1. MetricAvailability
1. MetricBaselineClient
1. MetricBaselinesProperties
1. MetricBaselinesResponse
1. MetricCriteria
1. MetricDefinition
1. MetricDefinitionCollection
1. MetricDefinitionsClient
1. MetricDimension
1. MetricNamespace
1. MetricNamespaceCollection
1. MetricNamespaceName
1. MetricNamespacesClient
1. MetricSettings
1. MetricSingleDimension
1. MetricTrigger
1. MetricValue
1. MetricsClient
1. MultiMetricCriteria
1. Operation
1. OperationDisplay
1. OperationListResult
1. OperationStatus
1. OperationsClient
1. PrivateEndpointConnection
1. PrivateEndpointConnectionListResult
1. PrivateEndpointConnectionListResultIterator
1. PrivateEndpointConnectionListResultPage
1. PrivateEndpointConnectionProperties
1. PrivateEndpointConnectionsClient
1. PrivateEndpointConnectionsCreateOrUpdateFuture
1. PrivateEndpointConnectionsDeleteFuture
1. PrivateEndpointProperty
1. PrivateLinkResource
1. PrivateLinkResourceListResult
1. PrivateLinkResourceListResultIterator
1. PrivateLinkResourceListResultPage
1. PrivateLinkResourceProperties
1. PrivateLinkResourcesClient
1. PrivateLinkScopeOperationStatusClient
1. PrivateLinkScopedResourcesClient
1. PrivateLinkScopedResourcesCreateOrUpdateFuture
1. PrivateLinkScopedResourcesDeleteFuture
1. PrivateLinkScopesClient
1. PrivateLinkScopesDeleteFuture
1. PrivateLinkScopesResource
1. PrivateLinkServiceConnectionStateProperty
1. ProxyOnlyResource
1. ProxyResource
1. Recurrence
1. RecurrentSchedule
1. Resource
1. Response
1. ResponseWithError
1. RetentionPolicy
1. RuleAction
1. RuleCondition
1. RuleDataSource
1. RuleEmailAction
1. RuleManagementEventClaimsDataSource
1. RuleManagementEventDataSource
1. RuleMetricDataSource
1. RuleWebhookAction
1. ScaleAction
1. ScaleCapacity
1. ScaleRule
1. ScaleRuleMetricDimension
1. Schedule
1. ScheduledQueryRulesClient
1. ScopedResource
1. ScopedResourceListResult
1. ScopedResourceListResultIterator
1. ScopedResourceListResultPage
1. ScopedResourceProperties
1. SenderAuthorization
1. SingleBaseline
1. SingleMetricBaseline
1. SmsReceiver
1. Source
1. SubscriptionDiagnosticSettings
1. SubscriptionDiagnosticSettingsClient
1. SubscriptionDiagnosticSettingsResource
1. SubscriptionDiagnosticSettingsResourceCollection
1. SubscriptionLogSettings
1. SubscriptionProxyOnlyResource
1. TagsResource
1. TenantActivityLogsClient
1. ThresholdRuleCondition
1. TimeSeriesBaseline
1. TimeSeriesElement
1. TimeSeriesInformation
1. TimeWindow
1. TriggerCondition
1. VMInsightsClient
1. VMInsightsOnboardingStatus
1. VMInsightsOnboardingStatusProperties
1. VoiceReceiver
1. WebhookNotification
1. WebhookReceiver
1. WebtestLocationAvailabilityCriteria
1. WorkspaceInfo
1. WorkspaceInfoProperties

### New Struct Fields

1. ErrorResponse.Code
1. ErrorResponse.Message
