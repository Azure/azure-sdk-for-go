# Release History

## 3.0.0 (2025-12-03)
### Breaking Changes

- Type of `Alert.ETag` has been changed from `*string` to `*azcore.ETag`
- Type of `Export.ETag` has been changed from `*string` to `*azcore.ETag`
- Type of `KpiProperties.Type` has been changed from `*KpiType` to `*KpiTypeType`
- Type of `OperationListResult.Value` has been changed from `[]*OperationForCostManagement` to `[]*Operation`
- Type of `PivotProperties.Type` has been changed from `*PivotType` to `*PivotTypeType`
- Type of `ScheduledAction.ETag` has been changed from `*string` to `*azcore.ETag`
- Type of `View.ETag` has been changed from `*string` to `*azcore.ETag`
- Enum `KpiType` has been removed
- Enum `PivotType` has been removed
- Function `*PriceSheetClient.BeginDownload` has been removed
- Operation `*AlertsClient.List` has supported pagination, use `*AlertsClient.NewListPager` instead.
- Operation `*ExportsClient.List` has supported pagination, use `*ExportsClient.NewListPager` instead.
- Struct `OperationForCostManagement` has been removed
- Field `DownloadURL` of struct `PriceSheetClientDownloadByBillingProfileResponse` has been removed

### Features Added

- New value `ExportTypeFocusCost`, `ExportTypePriceSheet`, `ExportTypeReservationDetails`, `ExportTypeReservationRecommendations`, `ExportTypeReservationTransactions` added to enum type `ExportType`
- New value `FormatTypeParquet` added to enum type `FormatType`
- New value `GranularityTypeMonthly` added to enum type `GranularityType`
- New value `TimeframeTypeTheCurrentMonth` added to enum type `TimeframeType`
- New enum type `BenefitUtilizationSummaryReportSchema` with values `BenefitUtilizationSummaryReportSchemaAvgUtilizationPercentage`, `BenefitUtilizationSummaryReportSchemaBenefitID`, `BenefitUtilizationSummaryReportSchemaBenefitOrderID`, `BenefitUtilizationSummaryReportSchemaBenefitType`, `BenefitUtilizationSummaryReportSchemaKind`, `BenefitUtilizationSummaryReportSchemaMaxUtilizationPercentage`, `BenefitUtilizationSummaryReportSchemaMinUtilizationPercentage`, `BenefitUtilizationSummaryReportSchemaUsageDate`, `BenefitUtilizationSummaryReportSchemaUtilizedPercentage`
- New enum type `BudgetNotificationOperatorType` with values `BudgetNotificationOperatorTypeEqualTo`, `BudgetNotificationOperatorTypeGreaterThan`, `BudgetNotificationOperatorTypeGreaterThanOrEqualTo`, `BudgetNotificationOperatorTypeLessThan`
- New enum type `BudgetOperatorType` with values `BudgetOperatorTypeIn`
- New enum type `CategoryType` with values `CategoryTypeCost`, `CategoryTypeReservationUtilization`
- New enum type `CompressionModeType` with values `CompressionModeTypeGzip`, `CompressionModeTypeNone`, `CompressionModeTypeSnappy`
- New enum type `CostAllocationPolicyType` with values `CostAllocationPolicyTypeFixedProportion`
- New enum type `CostAllocationResourceType` with values `CostAllocationResourceTypeDimension`, `CostAllocationResourceTypeTag`
- New enum type `CultureCode` with values `CultureCodeCsCz`, `CultureCodeDaDk`, `CultureCodeDeDe`, `CultureCodeEnGb`, `CultureCodeEnUs`, `CultureCodeEsEs`, `CultureCodeFrFr`, `CultureCodeHuHu`, `CultureCodeItIt`, `CultureCodeJaJp`, `CultureCodeKoKr`, `CultureCodeNbNo`, `CultureCodeNlNl`, `CultureCodePlPl`, `CultureCodePtBr`, `CultureCodePtPt`, `CultureCodeRuRu`, `CultureCodeSvSe`, `CultureCodeTrTr`, `CultureCodeZhCn`, `CultureCodeZhTw`
- New enum type `DataOverwriteBehaviorType` with values `DataOverwriteBehaviorTypeCreateNewReport`, `DataOverwriteBehaviorTypeOverwritePreviousReport`
- New enum type `DestinationType` with values `DestinationTypeAzureBlob`
- New enum type `FilterItemNames` with values `FilterItemNamesLookBackPeriod`, `FilterItemNamesReservationScope`, `FilterItemNamesResourceType`
- New enum type `Frequency` with values `FrequencyDaily`, `FrequencyMonthly`, `FrequencyWeekly`
- New enum type `KpiTypeType` with values `KpiTypeTypeBudget`, `KpiTypeTypeForecast`
- New enum type `PivotTypeType` with values `PivotTypeTypeDimension`, `PivotTypeTypeTagKey`
- New enum type `Reason` with values `ReasonAlreadyExists`, `ReasonInvalid`, `ReasonValid`
- New enum type `RuleStatus` with values `RuleStatusActive`, `RuleStatusNotActive`, `RuleStatusProcessing`
- New enum type `SettingType` with values `SettingTypeTaginheritance`
- New enum type `SettingsKind` with values `SettingsKindTaginheritance`
- New enum type `SystemAssignedServiceIdentityType` with values `SystemAssignedServiceIdentityTypeNone`, `SystemAssignedServiceIdentityTypeSystemAssigned`
- New enum type `ThresholdType` with values `ThresholdTypeActual`, `ThresholdTypeForecasted`
- New enum type `TimeGrainType` with values `TimeGrainTypeAnnually`, `TimeGrainTypeBillingAnnual`, `TimeGrainTypeBillingMonth`, `TimeGrainTypeBillingQuarter`, `TimeGrainTypeLast30Days`, `TimeGrainTypeLast7Days`, `TimeGrainTypeMonthly`, `TimeGrainTypeQuarterly`
- New function `NewBudgetsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*BudgetsClient, error)`
- New function `*BudgetsClient.CreateOrUpdate(ctx context.Context, scope string, budgetName string, parameters Budget, options *BudgetsClientCreateOrUpdateOptions) (BudgetsClientCreateOrUpdateResponse, error)`
- New function `*BudgetsClient.Delete(ctx context.Context, scope string, budgetName string, options *BudgetsClientDeleteOptions) (BudgetsClientDeleteResponse, error)`
- New function `*BudgetsClient.Get(ctx context.Context, scope string, budgetName string, options *BudgetsClientGetOptions) (BudgetsClientGetResponse, error)`
- New function `*BudgetsClient.NewListPager(scope string, options *BudgetsClientListOptions) *runtime.Pager[BudgetsClientListResponse]`
- New function `*ClientFactory.NewBudgetsClient() *BudgetsClient`
- New function `*ClientFactory.NewCostAllocationRulesClient() *CostAllocationRulesClient`
- New function `*ClientFactory.NewGenerateBenefitUtilizationSummariesReportClient() *GenerateBenefitUtilizationSummariesReportClient`
- New function `*ClientFactory.NewSettingsClient() *SettingsClient`
- New function `NewCostAllocationRulesClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*CostAllocationRulesClient, error)`
- New function `*CostAllocationRulesClient.CheckNameAvailability(ctx context.Context, billingAccountID string, costAllocationRuleCheckNameAvailabilityRequest CostAllocationRuleCheckNameAvailabilityRequest, options *CostAllocationRulesClientCheckNameAvailabilityOptions) (CostAllocationRulesClientCheckNameAvailabilityResponse, error)`
- New function `*CostAllocationRulesClient.CreateOrUpdate(ctx context.Context, billingAccountID string, ruleName string, costAllocationRule CostAllocationRuleDefinition, options *CostAllocationRulesClientCreateOrUpdateOptions) (CostAllocationRulesClientCreateOrUpdateResponse, error)`
- New function `*CostAllocationRulesClient.Delete(ctx context.Context, billingAccountID string, ruleName string, options *CostAllocationRulesClientDeleteOptions) (CostAllocationRulesClientDeleteResponse, error)`
- New function `*CostAllocationRulesClient.Get(ctx context.Context, billingAccountID string, ruleName string, options *CostAllocationRulesClientGetOptions) (CostAllocationRulesClientGetResponse, error)`
- New function `*CostAllocationRulesClient.NewListPager(billingAccountID string, options *CostAllocationRulesClientListOptions) *runtime.Pager[CostAllocationRulesClientListResponse]`
- New function `NewGenerateBenefitUtilizationSummariesReportClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*GenerateBenefitUtilizationSummariesReportClient, error)`
- New function `*GenerateBenefitUtilizationSummariesReportClient.BeginGenerateByBillingAccount(ctx context.Context, billingAccountID string, benefitUtilizationSummariesRequest BenefitUtilizationSummariesRequest, options *GenerateBenefitUtilizationSummariesReportClientBeginGenerateByBillingAccountOptions) (*runtime.Poller[GenerateBenefitUtilizationSummariesReportClientGenerateByBillingAccountResponse], error)`
- New function `*GenerateBenefitUtilizationSummariesReportClient.BeginGenerateByBillingProfile(ctx context.Context, billingAccountID string, billingProfileID string, benefitUtilizationSummariesRequest BenefitUtilizationSummariesRequest, options *GenerateBenefitUtilizationSummariesReportClientBeginGenerateByBillingProfileOptions) (*runtime.Poller[GenerateBenefitUtilizationSummariesReportClientGenerateByBillingProfileResponse], error)`
- New function `*GenerateBenefitUtilizationSummariesReportClient.BeginGenerateByReservationID(ctx context.Context, reservationOrderID string, reservationID string, benefitUtilizationSummariesRequest BenefitUtilizationSummariesRequest, options *GenerateBenefitUtilizationSummariesReportClientBeginGenerateByReservationIDOptions) (*runtime.Poller[GenerateBenefitUtilizationSummariesReportClientGenerateByReservationIDResponse], error)`
- New function `*GenerateBenefitUtilizationSummariesReportClient.BeginGenerateByReservationOrderID(ctx context.Context, reservationOrderID string, benefitUtilizationSummariesRequest BenefitUtilizationSummariesRequest, options *GenerateBenefitUtilizationSummariesReportClientBeginGenerateByReservationOrderIDOptions) (*runtime.Poller[GenerateBenefitUtilizationSummariesReportClientGenerateByReservationOrderIDResponse], error)`
- New function `*GenerateBenefitUtilizationSummariesReportClient.BeginGenerateBySavingsPlanID(ctx context.Context, savingsPlanOrderID string, savingsPlanID string, benefitUtilizationSummariesRequest BenefitUtilizationSummariesRequest, options *GenerateBenefitUtilizationSummariesReportClientBeginGenerateBySavingsPlanIDOptions) (*runtime.Poller[GenerateBenefitUtilizationSummariesReportClientGenerateBySavingsPlanIDResponse], error)`
- New function `*GenerateBenefitUtilizationSummariesReportClient.BeginGenerateBySavingsPlanOrderID(ctx context.Context, savingsPlanOrderID string, benefitUtilizationSummariesRequest BenefitUtilizationSummariesRequest, options *GenerateBenefitUtilizationSummariesReportClientBeginGenerateBySavingsPlanOrderIDOptions) (*runtime.Poller[GenerateBenefitUtilizationSummariesReportClientGenerateBySavingsPlanOrderIDResponse], error)`
- New function `*PriceSheetClient.BeginDownloadByBillingAccount(ctx context.Context, billingAccountID string, billingPeriodName string, options *PriceSheetClientBeginDownloadByBillingAccountOptions) (*runtime.Poller[PriceSheetClientDownloadByBillingAccountResponse], error)`
- New function `*PriceSheetClient.BeginDownloadByInvoice(ctx context.Context, billingAccountName string, billingProfileName string, invoiceName string, options *PriceSheetClientBeginDownloadByInvoiceOptions) (*runtime.Poller[PriceSheetClientDownloadByInvoiceResponse], error)`
- New function `*Setting.GetSetting() *Setting`
- New function `NewSettingsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*SettingsClient, error)`
- New function `*SettingsClient.CreateOrUpdateByScope(ctx context.Context, scope string, typeParam SettingType, setting SettingClassification, options *SettingsClientCreateOrUpdateByScopeOptions) (SettingsClientCreateOrUpdateByScopeResponse, error)`
- New function `*SettingsClient.DeleteByScope(ctx context.Context, scope string, typeParam SettingType, options *SettingsClientDeleteByScopeOptions) (SettingsClientDeleteByScopeResponse, error)`
- New function `*SettingsClient.GetByScope(ctx context.Context, scope string, typeParam SettingType, options *SettingsClientGetByScopeOptions) (SettingsClientGetByScopeResponse, error)`
- New function `*SettingsClient.NewListPager(scope string, options *SettingsClientListOptions) *runtime.Pager[SettingsClientListResponse]`
- New function `*TagInheritanceSetting.GetSetting() *Setting`
- New struct `AsyncOperationStatusProperties`
- New struct `BenefitUtilizationSummariesOperationStatus`
- New struct `BenefitUtilizationSummariesRequest`
- New struct `Budget`
- New struct `BudgetComparisonExpression`
- New struct `BudgetFilter`
- New struct `BudgetFilterProperties`
- New struct `BudgetProperties`
- New struct `BudgetTimePeriod`
- New struct `BudgetsListResult`
- New struct `CostAllocationProportion`
- New struct `CostAllocationRuleCheckNameAvailabilityRequest`
- New struct `CostAllocationRuleCheckNameAvailabilityResponse`
- New struct `CostAllocationRuleDefinition`
- New struct `CostAllocationRuleDetails`
- New struct `CostAllocationRuleList`
- New struct `CostAllocationRuleProperties`
- New struct `CurrentSpend`
- New struct `ExportRunRequest`
- New struct `ExportSuspensionContext`
- New struct `FilterItems`
- New struct `ForecastSpend`
- New struct `MCAPriceSheetProperties`
- New struct `Notification`
- New struct `Operation`
- New struct `PricesheetDownloadProperties`
- New struct `SettingsListResult`
- New struct `SourceCostAllocationResource`
- New struct `SystemAssignedServiceIdentity`
- New struct `TagInheritanceProperties`
- New struct `TagInheritanceSetting`
- New struct `TargetCostAllocationResource`
- New field `SystemData` in struct `Alert`
- New field `SystemData` in struct `BenefitRecommendationModel`
- New field `SystemData` in struct `BenefitUtilizationSummary`
- New field `CompressionMode`, `DataOverwriteBehavior`, `ExportDescription`, `SystemSuspensionContext` in struct `CommonExportProperties`
- New field `NextLink` in struct `DimensionsListResult`
- New field `Identity`, `Location`, `SystemData` in struct `Export`
- New field `DataVersion`, `Filters` in struct `ExportDatasetConfiguration`
- New field `Type` in struct `ExportDeliveryDestination`
- New field `CompressionMode`, `DataOverwriteBehavior`, `ExportDescription`, `SystemSuspensionContext` in struct `ExportProperties`
- New field `EndDate`, `ManifestFile`, `StartDate` in struct `ExportRunProperties`
- New field `Parameters` in struct `ExportsClientExecuteOptions`
- New field `AzureConsumptionAsyncOperation` in struct `GenerateDetailedCostReportClientCreateOperationResponse`
- New field `SystemData` in struct `GenerateDetailedCostReportOperationResult`
- New field `SystemData` in struct `GenerateDetailedCostReportOperationStatuses`
- New field `SystemData` in struct `IncludedQuantityUtilizationSummary`
- New anonymous field `PricesheetDownloadProperties` in struct `PriceSheetClientDownloadByBillingProfileResponse`
- New field `ODataEntityID` in struct `PriceSheetClientDownloadByBillingProfileResponse`
- New field `SystemData` in struct `SavingsPlanUtilizationSummary`
- New field `SystemData` in struct `View`


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-05-26)
### Breaking Changes

- Type of `ExportExecutionListResult.Value` has been changed from `[]*ExportExecution` to `[]*ExportRun`
- Type of `ForecastDataset.Aggregation` has been changed from `map[string]*QueryAggregation` to `map[string]*ForecastAggregation`
- Type of `ForecastDataset.Configuration` has been changed from `*QueryDatasetConfiguration` to `*ForecastDatasetConfiguration`
- Type of `ForecastDataset.Filter` has been changed from `*QueryFilter` to `*ForecastFilter`
- Type of `ForecastDefinition.TimePeriod` has been changed from `*QueryTimePeriod` to `*ForecastTimePeriod`
- Type of `ForecastDefinition.Timeframe` has been changed from `*ForecastTimeframeType` to `*ForecastTimeframe`
- Type of `OperationListResult.Value` has been changed from `[]*Operation` to `[]*OperationForCostManagement`
- Type of `ReportConfigGrouping.Type` has been changed from `*ReportConfigColumnType` to `*QueryColumnType`
- `QueryColumnTypeTag` from enum `QueryColumnType` has been removed
- Enum `ForecastTimeframeType` has been removed
- Enum `ReportConfigColumnType` has been removed
- Operation `*GenerateDetailedCostReportOperationResultsClient.Get` has been changed to LRO, use `*GenerateDetailedCostReportOperationResultsClient.BeginGet` instead.
- Field `QueryResult` of struct `ForecastClientExternalCloudProviderUsageResponse` has been removed
- Field `QueryResult` of struct `ForecastClientUsageResponse` has been removed

### Features Added

- New value `QueryColumnTypeTagKey` added to enum type `QueryColumnType`
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `BenefitKind` with values `BenefitKindIncludedQuantity`, `BenefitKindReservation`, `BenefitKindSavingsPlan`
- New enum type `CheckNameAvailabilityReason` with values `CheckNameAvailabilityReasonAlreadyExists`, `CheckNameAvailabilityReasonInvalid`
- New enum type `CostDetailsDataFormat` with values `CostDetailsDataFormatCSVCostDetailsDataFormat`
- New enum type `CostDetailsMetricType` with values `CostDetailsMetricTypeActualCostCostDetailsMetricType`, `CostDetailsMetricTypeAmortizedCostCostDetailsMetricType`
- New enum type `CostDetailsStatusType` with values `CostDetailsStatusTypeCompletedCostDetailsStatusType`, `CostDetailsStatusTypeFailedCostDetailsStatusType`, `CostDetailsStatusTypeNoDataFoundCostDetailsStatusType`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DaysOfWeek` with values `DaysOfWeekFriday`, `DaysOfWeekMonday`, `DaysOfWeekSaturday`, `DaysOfWeekSunday`, `DaysOfWeekThursday`, `DaysOfWeekTuesday`, `DaysOfWeekWednesday`
- New enum type `FileFormat` with values `FileFormatCSV`
- New enum type `ForecastOperatorType` with values `ForecastOperatorTypeIn`
- New enum type `ForecastTimeframe` with values `ForecastTimeframeCustom`
- New enum type `FunctionName` with values `FunctionNameCost`, `FunctionNameCostUSD`, `FunctionNamePreTaxCost`, `FunctionNamePreTaxCostUSD`
- New enum type `Grain` with values `GrainDaily`, `GrainHourly`, `GrainMonthly`
- New enum type `GrainParameter` with values `GrainParameterDaily`, `GrainParameterHourly`, `GrainParameterMonthly`
- New enum type `LookBackPeriod` with values `LookBackPeriodLast30Days`, `LookBackPeriodLast60Days`, `LookBackPeriodLast7Days`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ScheduleFrequency` with values `ScheduleFrequencyDaily`, `ScheduleFrequencyMonthly`, `ScheduleFrequencyWeekly`
- New enum type `ScheduledActionKind` with values `ScheduledActionKindEmail`, `ScheduledActionKindInsightAlert`
- New enum type `ScheduledActionStatus` with values `ScheduledActionStatusDisabled`, `ScheduledActionStatusEnabled`, `ScheduledActionStatusExpired`
- New enum type `Scope` with values `ScopeShared`, `ScopeSingle`
- New enum type `Term` with values `TermP1Y`, `TermP3Y`
- New enum type `WeeksOfMonth` with values `WeeksOfMonthFirst`, `WeeksOfMonthFourth`, `WeeksOfMonthLast`, `WeeksOfMonthSecond`, `WeeksOfMonthThird`
- New function `*BenefitRecommendationProperties.GetBenefitRecommendationProperties() *BenefitRecommendationProperties`
- New function `NewBenefitRecommendationsClient(azcore.TokenCredential, *arm.ClientOptions) (*BenefitRecommendationsClient, error)`
- New function `*BenefitRecommendationsClient.NewListPager(string, *BenefitRecommendationsClientListOptions) *runtime.Pager[BenefitRecommendationsClientListResponse]`
- New function `NewBenefitUtilizationSummariesClient(azcore.TokenCredential, *arm.ClientOptions) (*BenefitUtilizationSummariesClient, error)`
- New function `*BenefitUtilizationSummariesClient.NewListByBillingAccountIDPager(string, *BenefitUtilizationSummariesClientListByBillingAccountIDOptions) *runtime.Pager[BenefitUtilizationSummariesClientListByBillingAccountIDResponse]`
- New function `*BenefitUtilizationSummariesClient.NewListByBillingProfileIDPager(string, string, *BenefitUtilizationSummariesClientListByBillingProfileIDOptions) *runtime.Pager[BenefitUtilizationSummariesClientListByBillingProfileIDResponse]`
- New function `*BenefitUtilizationSummariesClient.NewListBySavingsPlanIDPager(string, string, *BenefitUtilizationSummariesClientListBySavingsPlanIDOptions) *runtime.Pager[BenefitUtilizationSummariesClientListBySavingsPlanIDResponse]`
- New function `*BenefitUtilizationSummariesClient.NewListBySavingsPlanOrderPager(string, *BenefitUtilizationSummariesClientListBySavingsPlanOrderOptions) *runtime.Pager[BenefitUtilizationSummariesClientListBySavingsPlanOrderResponse]`
- New function `*BenefitUtilizationSummary.GetBenefitUtilizationSummary() *BenefitUtilizationSummary`
- New function `*ClientFactory.NewBenefitRecommendationsClient() *BenefitRecommendationsClient`
- New function `*ClientFactory.NewBenefitUtilizationSummariesClient() *BenefitUtilizationSummariesClient`
- New function `*ClientFactory.NewGenerateCostDetailsReportClient() *GenerateCostDetailsReportClient`
- New function `*ClientFactory.NewPriceSheetClient() *PriceSheetClient`
- New function `*ClientFactory.NewScheduledActionsClient() *ScheduledActionsClient`
- New function `NewGenerateCostDetailsReportClient(azcore.TokenCredential, *arm.ClientOptions) (*GenerateCostDetailsReportClient, error)`
- New function `*GenerateCostDetailsReportClient.BeginCreateOperation(context.Context, string, GenerateCostDetailsReportRequestDefinition, *GenerateCostDetailsReportClientBeginCreateOperationOptions) (*runtime.Poller[GenerateCostDetailsReportClientCreateOperationResponse], error)`
- New function `*GenerateCostDetailsReportClient.BeginGetOperationResults(context.Context, string, string, *GenerateCostDetailsReportClientBeginGetOperationResultsOptions) (*runtime.Poller[GenerateCostDetailsReportClientGetOperationResultsResponse], error)`
- New function `*IncludedQuantityUtilizationSummary.GetBenefitUtilizationSummary() *BenefitUtilizationSummary`
- New function `PossibleGrainValues() []Grain`
- New function `NewPriceSheetClient(azcore.TokenCredential, *arm.ClientOptions) (*PriceSheetClient, error)`
- New function `*PriceSheetClient.BeginDownload(context.Context, string, string, string, *PriceSheetClientBeginDownloadOptions) (*runtime.Poller[PriceSheetClientDownloadResponse], error)`
- New function `*PriceSheetClient.BeginDownloadByBillingProfile(context.Context, string, string, *PriceSheetClientBeginDownloadByBillingProfileOptions) (*runtime.Poller[PriceSheetClientDownloadByBillingProfileResponse], error)`
- New function `*SavingsPlanUtilizationSummary.GetBenefitUtilizationSummary() *BenefitUtilizationSummary`
- New function `NewScheduledActionsClient(azcore.TokenCredential, *arm.ClientOptions) (*ScheduledActionsClient, error)`
- New function `*ScheduledActionsClient.CheckNameAvailability(context.Context, CheckNameAvailabilityRequest, *ScheduledActionsClientCheckNameAvailabilityOptions) (ScheduledActionsClientCheckNameAvailabilityResponse, error)`
- New function `*ScheduledActionsClient.CheckNameAvailabilityByScope(context.Context, string, CheckNameAvailabilityRequest, *ScheduledActionsClientCheckNameAvailabilityByScopeOptions) (ScheduledActionsClientCheckNameAvailabilityByScopeResponse, error)`
- New function `*ScheduledActionsClient.CreateOrUpdate(context.Context, string, ScheduledAction, *ScheduledActionsClientCreateOrUpdateOptions) (ScheduledActionsClientCreateOrUpdateResponse, error)`
- New function `*ScheduledActionsClient.CreateOrUpdateByScope(context.Context, string, string, ScheduledAction, *ScheduledActionsClientCreateOrUpdateByScopeOptions) (ScheduledActionsClientCreateOrUpdateByScopeResponse, error)`
- New function `*ScheduledActionsClient.Delete(context.Context, string, *ScheduledActionsClientDeleteOptions) (ScheduledActionsClientDeleteResponse, error)`
- New function `*ScheduledActionsClient.DeleteByScope(context.Context, string, string, *ScheduledActionsClientDeleteByScopeOptions) (ScheduledActionsClientDeleteByScopeResponse, error)`
- New function `*ScheduledActionsClient.Get(context.Context, string, *ScheduledActionsClientGetOptions) (ScheduledActionsClientGetResponse, error)`
- New function `*ScheduledActionsClient.GetByScope(context.Context, string, string, *ScheduledActionsClientGetByScopeOptions) (ScheduledActionsClientGetByScopeResponse, error)`
- New function `*ScheduledActionsClient.NewListByScopePager(string, *ScheduledActionsClientListByScopeOptions) *runtime.Pager[ScheduledActionsClientListByScopeResponse]`
- New function `*ScheduledActionsClient.NewListPager(*ScheduledActionsClientListOptions) *runtime.Pager[ScheduledActionsClientListResponse]`
- New function `*ScheduledActionsClient.Run(context.Context, string, *ScheduledActionsClientRunOptions) (ScheduledActionsClientRunResponse, error)`
- New function `*ScheduledActionsClient.RunByScope(context.Context, string, string, *ScheduledActionsClientRunByScopeOptions) (ScheduledActionsClientRunByScopeResponse, error)`
- New function `*SharedScopeBenefitRecommendationProperties.GetBenefitRecommendationProperties() *BenefitRecommendationProperties`
- New function `*SingleScopeBenefitRecommendationProperties.GetBenefitRecommendationProperties() *BenefitRecommendationProperties`
- New struct `AllSavingsBenefitDetails`
- New struct `AllSavingsList`
- New struct `BenefitRecommendationModel`
- New struct `BenefitRecommendationsListResult`
- New struct `BenefitUtilizationSummariesListResult`
- New struct `BlobInfo`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityResponse`
- New struct `CostDetailsOperationResults`
- New struct `CostDetailsTimePeriod`
- New struct `ExportRun`
- New struct `ExportRunProperties`
- New struct `FileDestination`
- New struct `ForecastAggregation`
- New struct `ForecastColumn`
- New struct `ForecastComparisonExpression`
- New struct `ForecastDatasetConfiguration`
- New struct `ForecastFilter`
- New struct `ForecastProperties`
- New struct `ForecastResult`
- New struct `ForecastTimePeriod`
- New struct `GenerateCostDetailsReportRequestDefinition`
- New struct `IncludedQuantityUtilizationSummary`
- New struct `IncludedQuantityUtilizationSummaryProperties`
- New struct `NotificationProperties`
- New struct `OperationForCostManagement`
- New struct `RecommendationUsageDetails`
- New struct `ReportManifest`
- New struct `RequestContext`
- New struct `SavingsPlanUtilizationSummary`
- New struct `SavingsPlanUtilizationSummaryProperties`
- New struct `ScheduleProperties`
- New struct `ScheduledAction`
- New struct `ScheduledActionListResult`
- New struct `ScheduledActionProperties`
- New struct `SharedScopeBenefitRecommendationProperties`
- New struct `SingleScopeBenefitRecommendationProperties`
- New struct `SystemData`
- New field `ExpiryTime` in struct `DownloadURL`
- New anonymous field `ForecastResult` in struct `ForecastClientExternalCloudProviderUsageResponse`
- New anonymous field `ForecastResult` in struct `ForecastClientUsageResponse`
- New field `EndTime`, `StartTime` in struct `GenerateDetailedCostReportOperationStatuses`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.



## 1.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).