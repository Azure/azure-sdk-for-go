//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcostmanagement

const (
	moduleName    = "armcostmanagement"
	moduleVersion = "v1.1.0"
)

// AccumulatedType - Show costs accumulated over time.
type AccumulatedType string

const (
	AccumulatedTypeFalse AccumulatedType = "false"
	AccumulatedTypeTrue  AccumulatedType = "true"
)

// PossibleAccumulatedTypeValues returns the possible values for the AccumulatedType const type.
func PossibleAccumulatedTypeValues() []AccumulatedType {
	return []AccumulatedType{
		AccumulatedTypeFalse,
		AccumulatedTypeTrue,
	}
}

// AlertCategory - Alert category
type AlertCategory string

const (
	AlertCategoryBilling AlertCategory = "Billing"
	AlertCategoryCost    AlertCategory = "Cost"
	AlertCategorySystem  AlertCategory = "System"
	AlertCategoryUsage   AlertCategory = "Usage"
)

// PossibleAlertCategoryValues returns the possible values for the AlertCategory const type.
func PossibleAlertCategoryValues() []AlertCategory {
	return []AlertCategory{
		AlertCategoryBilling,
		AlertCategoryCost,
		AlertCategorySystem,
		AlertCategoryUsage,
	}
}

// AlertCriteria - Criteria that triggered alert
type AlertCriteria string

const (
	AlertCriteriaCostThresholdExceeded          AlertCriteria = "CostThresholdExceeded"
	AlertCriteriaCreditThresholdApproaching     AlertCriteria = "CreditThresholdApproaching"
	AlertCriteriaCreditThresholdReached         AlertCriteria = "CreditThresholdReached"
	AlertCriteriaCrossCloudCollectionError      AlertCriteria = "CrossCloudCollectionError"
	AlertCriteriaCrossCloudNewDataAvailable     AlertCriteria = "CrossCloudNewDataAvailable"
	AlertCriteriaForecastCostThresholdExceeded  AlertCriteria = "ForecastCostThresholdExceeded"
	AlertCriteriaForecastUsageThresholdExceeded AlertCriteria = "ForecastUsageThresholdExceeded"
	AlertCriteriaGeneralThresholdError          AlertCriteria = "GeneralThresholdError"
	AlertCriteriaInvoiceDueDateApproaching      AlertCriteria = "InvoiceDueDateApproaching"
	AlertCriteriaInvoiceDueDateReached          AlertCriteria = "InvoiceDueDateReached"
	AlertCriteriaMultiCurrency                  AlertCriteria = "MultiCurrency"
	AlertCriteriaQuotaThresholdApproaching      AlertCriteria = "QuotaThresholdApproaching"
	AlertCriteriaQuotaThresholdReached          AlertCriteria = "QuotaThresholdReached"
	AlertCriteriaUsageThresholdExceeded         AlertCriteria = "UsageThresholdExceeded"
)

// PossibleAlertCriteriaValues returns the possible values for the AlertCriteria const type.
func PossibleAlertCriteriaValues() []AlertCriteria {
	return []AlertCriteria{
		AlertCriteriaCostThresholdExceeded,
		AlertCriteriaCreditThresholdApproaching,
		AlertCriteriaCreditThresholdReached,
		AlertCriteriaCrossCloudCollectionError,
		AlertCriteriaCrossCloudNewDataAvailable,
		AlertCriteriaForecastCostThresholdExceeded,
		AlertCriteriaForecastUsageThresholdExceeded,
		AlertCriteriaGeneralThresholdError,
		AlertCriteriaInvoiceDueDateApproaching,
		AlertCriteriaInvoiceDueDateReached,
		AlertCriteriaMultiCurrency,
		AlertCriteriaQuotaThresholdApproaching,
		AlertCriteriaQuotaThresholdReached,
		AlertCriteriaUsageThresholdExceeded,
	}
}

// AlertOperator - operator used to compare currentSpend with amount
type AlertOperator string

const (
	AlertOperatorEqualTo              AlertOperator = "EqualTo"
	AlertOperatorGreaterThan          AlertOperator = "GreaterThan"
	AlertOperatorGreaterThanOrEqualTo AlertOperator = "GreaterThanOrEqualTo"
	AlertOperatorLessThan             AlertOperator = "LessThan"
	AlertOperatorLessThanOrEqualTo    AlertOperator = "LessThanOrEqualTo"
	AlertOperatorNone                 AlertOperator = "None"
)

// PossibleAlertOperatorValues returns the possible values for the AlertOperator const type.
func PossibleAlertOperatorValues() []AlertOperator {
	return []AlertOperator{
		AlertOperatorEqualTo,
		AlertOperatorGreaterThan,
		AlertOperatorGreaterThanOrEqualTo,
		AlertOperatorLessThan,
		AlertOperatorLessThanOrEqualTo,
		AlertOperatorNone,
	}
}

// AlertSource - Source of alert
type AlertSource string

const (
	AlertSourcePreset AlertSource = "Preset"
	AlertSourceUser   AlertSource = "User"
)

// PossibleAlertSourceValues returns the possible values for the AlertSource const type.
func PossibleAlertSourceValues() []AlertSource {
	return []AlertSource{
		AlertSourcePreset,
		AlertSourceUser,
	}
}

// AlertStatus - alert status
type AlertStatus string

const (
	AlertStatusActive     AlertStatus = "Active"
	AlertStatusDismissed  AlertStatus = "Dismissed"
	AlertStatusNone       AlertStatus = "None"
	AlertStatusOverridden AlertStatus = "Overridden"
	AlertStatusResolved   AlertStatus = "Resolved"
)

// PossibleAlertStatusValues returns the possible values for the AlertStatus const type.
func PossibleAlertStatusValues() []AlertStatus {
	return []AlertStatus{
		AlertStatusActive,
		AlertStatusDismissed,
		AlertStatusNone,
		AlertStatusOverridden,
		AlertStatusResolved,
	}
}

// AlertTimeGrainType - Type of timegrain cadence
type AlertTimeGrainType string

const (
	AlertTimeGrainTypeAnnually       AlertTimeGrainType = "Annually"
	AlertTimeGrainTypeBillingAnnual  AlertTimeGrainType = "BillingAnnual"
	AlertTimeGrainTypeBillingMonth   AlertTimeGrainType = "BillingMonth"
	AlertTimeGrainTypeBillingQuarter AlertTimeGrainType = "BillingQuarter"
	AlertTimeGrainTypeMonthly        AlertTimeGrainType = "Monthly"
	AlertTimeGrainTypeNone           AlertTimeGrainType = "None"
	AlertTimeGrainTypeQuarterly      AlertTimeGrainType = "Quarterly"
)

// PossibleAlertTimeGrainTypeValues returns the possible values for the AlertTimeGrainType const type.
func PossibleAlertTimeGrainTypeValues() []AlertTimeGrainType {
	return []AlertTimeGrainType{
		AlertTimeGrainTypeAnnually,
		AlertTimeGrainTypeBillingAnnual,
		AlertTimeGrainTypeBillingMonth,
		AlertTimeGrainTypeBillingQuarter,
		AlertTimeGrainTypeMonthly,
		AlertTimeGrainTypeNone,
		AlertTimeGrainTypeQuarterly,
	}
}

// AlertType - type of alert
type AlertType string

const (
	AlertTypeBudget         AlertType = "Budget"
	AlertTypeBudgetForecast AlertType = "BudgetForecast"
	AlertTypeCredit         AlertType = "Credit"
	AlertTypeGeneral        AlertType = "General"
	AlertTypeInvoice        AlertType = "Invoice"
	AlertTypeQuota          AlertType = "Quota"
	AlertTypeXCloud         AlertType = "xCloud"
)

// PossibleAlertTypeValues returns the possible values for the AlertType const type.
func PossibleAlertTypeValues() []AlertType {
	return []AlertType{
		AlertTypeBudget,
		AlertTypeBudgetForecast,
		AlertTypeCredit,
		AlertTypeGeneral,
		AlertTypeInvoice,
		AlertTypeQuota,
		AlertTypeXCloud,
	}
}

// ChartType - Chart type of the main view in Cost Analysis. Required.
type ChartType string

const (
	ChartTypeArea          ChartType = "Area"
	ChartTypeGroupedColumn ChartType = "GroupedColumn"
	ChartTypeLine          ChartType = "Line"
	ChartTypeStackedColumn ChartType = "StackedColumn"
	ChartTypeTable         ChartType = "Table"
)

// PossibleChartTypeValues returns the possible values for the ChartType const type.
func PossibleChartTypeValues() []ChartType {
	return []ChartType{
		ChartTypeArea,
		ChartTypeGroupedColumn,
		ChartTypeLine,
		ChartTypeStackedColumn,
		ChartTypeTable,
	}
}

// ExecutionStatus - The last known status of the export execution.
type ExecutionStatus string

const (
	ExecutionStatusCompleted           ExecutionStatus = "Completed"
	ExecutionStatusDataNotAvailable    ExecutionStatus = "DataNotAvailable"
	ExecutionStatusFailed              ExecutionStatus = "Failed"
	ExecutionStatusInProgress          ExecutionStatus = "InProgress"
	ExecutionStatusNewDataNotAvailable ExecutionStatus = "NewDataNotAvailable"
	ExecutionStatusQueued              ExecutionStatus = "Queued"
	ExecutionStatusTimeout             ExecutionStatus = "Timeout"
)

// PossibleExecutionStatusValues returns the possible values for the ExecutionStatus const type.
func PossibleExecutionStatusValues() []ExecutionStatus {
	return []ExecutionStatus{
		ExecutionStatusCompleted,
		ExecutionStatusDataNotAvailable,
		ExecutionStatusFailed,
		ExecutionStatusInProgress,
		ExecutionStatusNewDataNotAvailable,
		ExecutionStatusQueued,
		ExecutionStatusTimeout,
	}
}

// ExecutionType - The type of the export execution.
type ExecutionType string

const (
	ExecutionTypeOnDemand  ExecutionType = "OnDemand"
	ExecutionTypeScheduled ExecutionType = "Scheduled"
)

// PossibleExecutionTypeValues returns the possible values for the ExecutionType const type.
func PossibleExecutionTypeValues() []ExecutionType {
	return []ExecutionType{
		ExecutionTypeOnDemand,
		ExecutionTypeScheduled,
	}
}

// ExportType - The type of the export. Note that 'Usage' is equivalent to 'ActualCost' and is applicable to exports that
// do not yet provide data for charges or amortization for service reservations.
type ExportType string

const (
	ExportTypeActualCost    ExportType = "ActualCost"
	ExportTypeAmortizedCost ExportType = "AmortizedCost"
	ExportTypeUsage         ExportType = "Usage"
)

// PossibleExportTypeValues returns the possible values for the ExportType const type.
func PossibleExportTypeValues() []ExportType {
	return []ExportType{
		ExportTypeActualCost,
		ExportTypeAmortizedCost,
		ExportTypeUsage,
	}
}

type ExternalCloudProviderType string

const (
	ExternalCloudProviderTypeExternalBillingAccounts ExternalCloudProviderType = "externalBillingAccounts"
	ExternalCloudProviderTypeExternalSubscriptions   ExternalCloudProviderType = "externalSubscriptions"
)

// PossibleExternalCloudProviderTypeValues returns the possible values for the ExternalCloudProviderType const type.
func PossibleExternalCloudProviderTypeValues() []ExternalCloudProviderType {
	return []ExternalCloudProviderType{
		ExternalCloudProviderTypeExternalBillingAccounts,
		ExternalCloudProviderTypeExternalSubscriptions,
	}
}

// ForecastTimeframeType - The time frame for pulling data for the forecast. If custom, then a specific time period must be
// provided.
type ForecastTimeframeType string

const (
	ForecastTimeframeTypeBillingMonthToDate  ForecastTimeframeType = "BillingMonthToDate"
	ForecastTimeframeTypeCustom              ForecastTimeframeType = "Custom"
	ForecastTimeframeTypeMonthToDate         ForecastTimeframeType = "MonthToDate"
	ForecastTimeframeTypeTheLastBillingMonth ForecastTimeframeType = "TheLastBillingMonth"
	ForecastTimeframeTypeTheLastMonth        ForecastTimeframeType = "TheLastMonth"
	ForecastTimeframeTypeWeekToDate          ForecastTimeframeType = "WeekToDate"
)

// PossibleForecastTimeframeTypeValues returns the possible values for the ForecastTimeframeType const type.
func PossibleForecastTimeframeTypeValues() []ForecastTimeframeType {
	return []ForecastTimeframeType{
		ForecastTimeframeTypeBillingMonthToDate,
		ForecastTimeframeTypeCustom,
		ForecastTimeframeTypeMonthToDate,
		ForecastTimeframeTypeTheLastBillingMonth,
		ForecastTimeframeTypeTheLastMonth,
		ForecastTimeframeTypeWeekToDate,
	}
}

// ForecastType - The type of the forecast.
type ForecastType string

const (
	ForecastTypeActualCost    ForecastType = "ActualCost"
	ForecastTypeAmortizedCost ForecastType = "AmortizedCost"
	ForecastTypeUsage         ForecastType = "Usage"
)

// PossibleForecastTypeValues returns the possible values for the ForecastType const type.
func PossibleForecastTypeValues() []ForecastType {
	return []ForecastType{
		ForecastTypeActualCost,
		ForecastTypeAmortizedCost,
		ForecastTypeUsage,
	}
}

// FormatType - The format of the export being delivered. Currently only 'Csv' is supported.
type FormatType string

const (
	FormatTypeCSV FormatType = "Csv"
)

// PossibleFormatTypeValues returns the possible values for the FormatType const type.
func PossibleFormatTypeValues() []FormatType {
	return []FormatType{
		FormatTypeCSV,
	}
}

// FunctionType - The name of the aggregation function to use.
type FunctionType string

const (
	FunctionTypeSum FunctionType = "Sum"
)

// PossibleFunctionTypeValues returns the possible values for the FunctionType const type.
func PossibleFunctionTypeValues() []FunctionType {
	return []FunctionType{
		FunctionTypeSum,
	}
}

// GenerateDetailedCostReportMetricType - The type of the detailed report. By default ActualCost is provided
type GenerateDetailedCostReportMetricType string

const (
	GenerateDetailedCostReportMetricTypeActualCost    GenerateDetailedCostReportMetricType = "ActualCost"
	GenerateDetailedCostReportMetricTypeAmortizedCost GenerateDetailedCostReportMetricType = "AmortizedCost"
)

// PossibleGenerateDetailedCostReportMetricTypeValues returns the possible values for the GenerateDetailedCostReportMetricType const type.
func PossibleGenerateDetailedCostReportMetricTypeValues() []GenerateDetailedCostReportMetricType {
	return []GenerateDetailedCostReportMetricType{
		GenerateDetailedCostReportMetricTypeActualCost,
		GenerateDetailedCostReportMetricTypeAmortizedCost,
	}
}

// GranularityType - The granularity of rows in the export. Currently only 'Daily' is supported.
type GranularityType string

const (
	GranularityTypeDaily GranularityType = "Daily"
)

// PossibleGranularityTypeValues returns the possible values for the GranularityType const type.
func PossibleGranularityTypeValues() []GranularityType {
	return []GranularityType{
		GranularityTypeDaily,
	}
}

// KpiType - KPI type (Forecast, Budget).
type KpiType string

const (
	KpiTypeBudget   KpiType = "Budget"
	KpiTypeForecast KpiType = "Forecast"
)

// PossibleKpiTypeValues returns the possible values for the KpiType const type.
func PossibleKpiTypeValues() []KpiType {
	return []KpiType{
		KpiTypeBudget,
		KpiTypeForecast,
	}
}

// MetricType - Metric to use when displaying costs.
type MetricType string

const (
	MetricTypeAHUB          MetricType = "AHUB"
	MetricTypeActualCost    MetricType = "ActualCost"
	MetricTypeAmortizedCost MetricType = "AmortizedCost"
)

// PossibleMetricTypeValues returns the possible values for the MetricType const type.
func PossibleMetricTypeValues() []MetricType {
	return []MetricType{
		MetricTypeAHUB,
		MetricTypeActualCost,
		MetricTypeAmortizedCost,
	}
}

// OperationStatusType - The status of the long running operation.
type OperationStatusType string

const (
	OperationStatusTypeCompleted OperationStatusType = "Completed"
	OperationStatusTypeFailed    OperationStatusType = "Failed"
	OperationStatusTypeRunning   OperationStatusType = "Running"
)

// PossibleOperationStatusTypeValues returns the possible values for the OperationStatusType const type.
func PossibleOperationStatusTypeValues() []OperationStatusType {
	return []OperationStatusType{
		OperationStatusTypeCompleted,
		OperationStatusTypeFailed,
		OperationStatusTypeRunning,
	}
}

// OperatorType - The operator to use for comparison.
type OperatorType string

const (
	OperatorTypeContains OperatorType = "Contains"
	OperatorTypeIn       OperatorType = "In"
)

// PossibleOperatorTypeValues returns the possible values for the OperatorType const type.
func PossibleOperatorTypeValues() []OperatorType {
	return []OperatorType{
		OperatorTypeContains,
		OperatorTypeIn,
	}
}

// PivotType - Data type to show in view.
type PivotType string

const (
	PivotTypeDimension PivotType = "Dimension"
	PivotTypeTagKey    PivotType = "TagKey"
)

// PossiblePivotTypeValues returns the possible values for the PivotType const type.
func PossiblePivotTypeValues() []PivotType {
	return []PivotType{
		PivotTypeDimension,
		PivotTypeTagKey,
	}
}

// QueryColumnType - The type of the column in the export.
type QueryColumnType string

const (
	QueryColumnTypeDimension QueryColumnType = "Dimension"
	QueryColumnTypeTag       QueryColumnType = "Tag"
)

// PossibleQueryColumnTypeValues returns the possible values for the QueryColumnType const type.
func PossibleQueryColumnTypeValues() []QueryColumnType {
	return []QueryColumnType{
		QueryColumnTypeDimension,
		QueryColumnTypeTag,
	}
}

// QueryOperatorType - The operator to use for comparison.
type QueryOperatorType string

const (
	QueryOperatorTypeIn QueryOperatorType = "In"
)

// PossibleQueryOperatorTypeValues returns the possible values for the QueryOperatorType const type.
func PossibleQueryOperatorTypeValues() []QueryOperatorType {
	return []QueryOperatorType{
		QueryOperatorTypeIn,
	}
}

// RecurrenceType - The schedule recurrence.
type RecurrenceType string

const (
	RecurrenceTypeAnnually RecurrenceType = "Annually"
	RecurrenceTypeDaily    RecurrenceType = "Daily"
	RecurrenceTypeMonthly  RecurrenceType = "Monthly"
	RecurrenceTypeWeekly   RecurrenceType = "Weekly"
)

// PossibleRecurrenceTypeValues returns the possible values for the RecurrenceType const type.
func PossibleRecurrenceTypeValues() []RecurrenceType {
	return []RecurrenceType{
		RecurrenceTypeAnnually,
		RecurrenceTypeDaily,
		RecurrenceTypeMonthly,
		RecurrenceTypeWeekly,
	}
}

// ReportConfigColumnType - The type of the column in the report.
type ReportConfigColumnType string

const (
	ReportConfigColumnTypeDimension ReportConfigColumnType = "Dimension"
	ReportConfigColumnTypeTag       ReportConfigColumnType = "Tag"
)

// PossibleReportConfigColumnTypeValues returns the possible values for the ReportConfigColumnType const type.
func PossibleReportConfigColumnTypeValues() []ReportConfigColumnType {
	return []ReportConfigColumnType{
		ReportConfigColumnTypeDimension,
		ReportConfigColumnTypeTag,
	}
}

// ReportConfigSortingType - Direction of sort.
type ReportConfigSortingType string

const (
	ReportConfigSortingTypeAscending  ReportConfigSortingType = "Ascending"
	ReportConfigSortingTypeDescending ReportConfigSortingType = "Descending"
)

// PossibleReportConfigSortingTypeValues returns the possible values for the ReportConfigSortingType const type.
func PossibleReportConfigSortingTypeValues() []ReportConfigSortingType {
	return []ReportConfigSortingType{
		ReportConfigSortingTypeAscending,
		ReportConfigSortingTypeDescending,
	}
}

// ReportGranularityType - The granularity of rows in the report.
type ReportGranularityType string

const (
	ReportGranularityTypeDaily   ReportGranularityType = "Daily"
	ReportGranularityTypeMonthly ReportGranularityType = "Monthly"
)

// PossibleReportGranularityTypeValues returns the possible values for the ReportGranularityType const type.
func PossibleReportGranularityTypeValues() []ReportGranularityType {
	return []ReportGranularityType{
		ReportGranularityTypeDaily,
		ReportGranularityTypeMonthly,
	}
}

// ReportOperationStatusType - The status of the long running operation.
type ReportOperationStatusType string

const (
	ReportOperationStatusTypeCompleted       ReportOperationStatusType = "Completed"
	ReportOperationStatusTypeFailed          ReportOperationStatusType = "Failed"
	ReportOperationStatusTypeInProgress      ReportOperationStatusType = "InProgress"
	ReportOperationStatusTypeNoDataFound     ReportOperationStatusType = "NoDataFound"
	ReportOperationStatusTypeQueued          ReportOperationStatusType = "Queued"
	ReportOperationStatusTypeReadyToDownload ReportOperationStatusType = "ReadyToDownload"
	ReportOperationStatusTypeTimedOut        ReportOperationStatusType = "TimedOut"
)

// PossibleReportOperationStatusTypeValues returns the possible values for the ReportOperationStatusType const type.
func PossibleReportOperationStatusTypeValues() []ReportOperationStatusType {
	return []ReportOperationStatusType{
		ReportOperationStatusTypeCompleted,
		ReportOperationStatusTypeFailed,
		ReportOperationStatusTypeInProgress,
		ReportOperationStatusTypeNoDataFound,
		ReportOperationStatusTypeQueued,
		ReportOperationStatusTypeReadyToDownload,
		ReportOperationStatusTypeTimedOut,
	}
}

// ReportTimeframeType - The time frame for pulling data for the report. If custom, then a specific time period must be provided.
type ReportTimeframeType string

const (
	ReportTimeframeTypeCustom      ReportTimeframeType = "Custom"
	ReportTimeframeTypeMonthToDate ReportTimeframeType = "MonthToDate"
	ReportTimeframeTypeWeekToDate  ReportTimeframeType = "WeekToDate"
	ReportTimeframeTypeYearToDate  ReportTimeframeType = "YearToDate"
)

// PossibleReportTimeframeTypeValues returns the possible values for the ReportTimeframeType const type.
func PossibleReportTimeframeTypeValues() []ReportTimeframeType {
	return []ReportTimeframeType{
		ReportTimeframeTypeCustom,
		ReportTimeframeTypeMonthToDate,
		ReportTimeframeTypeWeekToDate,
		ReportTimeframeTypeYearToDate,
	}
}

// ReportType - The type of the report. Usage represents actual usage, forecast represents forecasted data and UsageAndForecast
// represents both usage and forecasted data. Actual usage and forecasted data can be
// differentiated based on dates.
type ReportType string

const (
	ReportTypeUsage ReportType = "Usage"
)

// PossibleReportTypeValues returns the possible values for the ReportType const type.
func PossibleReportTypeValues() []ReportType {
	return []ReportType{
		ReportTypeUsage,
	}
}

// ReservationReportSchema - The CSV file from the reportUrl blob link consists of reservation usage data with the following
// schema at daily granularity
type ReservationReportSchema string

const (
	ReservationReportSchemaInstanceFlexibilityGroup ReservationReportSchema = "InstanceFlexibilityGroup"
	ReservationReportSchemaInstanceFlexibilityRatio ReservationReportSchema = "InstanceFlexibilityRatio"
	ReservationReportSchemaInstanceID               ReservationReportSchema = "InstanceId"
	ReservationReportSchemaKind                     ReservationReportSchema = "Kind"
	ReservationReportSchemaReservationID            ReservationReportSchema = "ReservationId"
	ReservationReportSchemaReservationOrderID       ReservationReportSchema = "ReservationOrderId"
	ReservationReportSchemaReservedHours            ReservationReportSchema = "ReservedHours"
	ReservationReportSchemaSKUName                  ReservationReportSchema = "SkuName"
	ReservationReportSchemaTotalReservedQuantity    ReservationReportSchema = "TotalReservedQuantity"
	ReservationReportSchemaUsageDate                ReservationReportSchema = "UsageDate"
	ReservationReportSchemaUsedHours                ReservationReportSchema = "UsedHours"
)

// PossibleReservationReportSchemaValues returns the possible values for the ReservationReportSchema const type.
func PossibleReservationReportSchemaValues() []ReservationReportSchema {
	return []ReservationReportSchema{
		ReservationReportSchemaInstanceFlexibilityGroup,
		ReservationReportSchemaInstanceFlexibilityRatio,
		ReservationReportSchemaInstanceID,
		ReservationReportSchemaKind,
		ReservationReportSchemaReservationID,
		ReservationReportSchemaReservationOrderID,
		ReservationReportSchemaReservedHours,
		ReservationReportSchemaSKUName,
		ReservationReportSchemaTotalReservedQuantity,
		ReservationReportSchemaUsageDate,
		ReservationReportSchemaUsedHours,
	}
}

// StatusType - The status of the export's schedule. If 'Inactive', the export's schedule is paused.
type StatusType string

const (
	StatusTypeActive   StatusType = "Active"
	StatusTypeInactive StatusType = "Inactive"
)

// PossibleStatusTypeValues returns the possible values for the StatusType const type.
func PossibleStatusTypeValues() []StatusType {
	return []StatusType{
		StatusTypeActive,
		StatusTypeInactive,
	}
}

// TimeframeType - The time frame for pulling data for the export. If custom, then a specific time period must be provided.
type TimeframeType string

const (
	TimeframeTypeBillingMonthToDate  TimeframeType = "BillingMonthToDate"
	TimeframeTypeCustom              TimeframeType = "Custom"
	TimeframeTypeMonthToDate         TimeframeType = "MonthToDate"
	TimeframeTypeTheLastBillingMonth TimeframeType = "TheLastBillingMonth"
	TimeframeTypeTheLastMonth        TimeframeType = "TheLastMonth"
	TimeframeTypeWeekToDate          TimeframeType = "WeekToDate"
)

// PossibleTimeframeTypeValues returns the possible values for the TimeframeType const type.
func PossibleTimeframeTypeValues() []TimeframeType {
	return []TimeframeType{
		TimeframeTypeBillingMonthToDate,
		TimeframeTypeCustom,
		TimeframeTypeMonthToDate,
		TimeframeTypeTheLastBillingMonth,
		TimeframeTypeTheLastMonth,
		TimeframeTypeWeekToDate,
	}
}
