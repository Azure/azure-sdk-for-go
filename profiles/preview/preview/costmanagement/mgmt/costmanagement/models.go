// +build go1.9

// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package costmanagement

import original "github.com/Azure/azure-sdk-for-go/services/preview/costmanagement/mgmt/2018-08-01-preview/costmanagement"

type BillingAccountDimensionsClient = original.BillingAccountDimensionsClient

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type ConnectorClient = original.ConnectorClient
type ConnectorStatus = original.ConnectorStatus

const (
	Active    ConnectorStatus = original.Active
	Error     ConnectorStatus = original.Error
	Suspended ConnectorStatus = original.Suspended
)

type ExecutionStatus = original.ExecutionStatus

const (
	Completed  ExecutionStatus = original.Completed
	Failed     ExecutionStatus = original.Failed
	InProgress ExecutionStatus = original.InProgress
	Queud      ExecutionStatus = original.Queud
	Timeout    ExecutionStatus = original.Timeout
)

type ExecutionType = original.ExecutionType

const (
	OnDemand  ExecutionType = original.OnDemand
	Scheduled ExecutionType = original.Scheduled
)

type FormatType = original.FormatType

const (
	Csv FormatType = original.Csv
)

type GranularityType = original.GranularityType

const (
	Daily GranularityType = original.Daily
)

type RecurrenceType = original.RecurrenceType

const (
	RecurrenceTypeAnnually RecurrenceType = original.RecurrenceTypeAnnually
	RecurrenceTypeDaily    RecurrenceType = original.RecurrenceTypeDaily
	RecurrenceTypeMonthly  RecurrenceType = original.RecurrenceTypeMonthly
	RecurrenceTypeWeekly   RecurrenceType = original.RecurrenceTypeWeekly
)

type ReportColumnType = original.ReportColumnType

const (
	ReportColumnTypeDimension ReportColumnType = original.ReportColumnTypeDimension
	ReportColumnTypeTag       ReportColumnType = original.ReportColumnTypeTag
)

type StatusType = original.StatusType

const (
	StatusTypeActive   StatusType = original.StatusTypeActive
	StatusTypeInactive StatusType = original.StatusTypeInactive
)

type TimeframeType = original.TimeframeType

const (
	Custom      TimeframeType = original.Custom
	MonthToDate TimeframeType = original.MonthToDate
	WeekToDate  TimeframeType = original.WeekToDate
)

type CommonReportProperties = original.CommonReportProperties
type ConnectorCollectionErrorInfo = original.ConnectorCollectionErrorInfo
type ConnectorCollectionInfo = original.ConnectorCollectionInfo
type ConnectorDefinition = original.ConnectorDefinition
type ConnectorDefinitionListResult = original.ConnectorDefinitionListResult
type ConnectorProperties = original.ConnectorProperties
type Dimension = original.Dimension
type DimensionProperties = original.DimensionProperties
type DimensionsListResult = original.DimensionsListResult
type ErrorBase = original.ErrorBase
type ErrorDetails = original.ErrorDetails
type ErrorResponse = original.ErrorResponse
type Query = original.Query
type QueryColumn = original.QueryColumn
type QueryProperties = original.QueryProperties
type QueryResult = original.QueryResult
type Report = original.Report
type ReportAggregation = original.ReportAggregation
type ReportComparisonExpression = original.ReportComparisonExpression
type ReportDataset = original.ReportDataset
type ReportDatasetConfiguration = original.ReportDatasetConfiguration
type ReportDefinition = original.ReportDefinition
type ReportDeliveryDestination = original.ReportDeliveryDestination
type ReportDeliveryInfo = original.ReportDeliveryInfo
type ReportExecution = original.ReportExecution
type ReportExecutionListResult = original.ReportExecutionListResult
type ReportExecutionProperties = original.ReportExecutionProperties
type ReportFilter = original.ReportFilter
type ReportGrouping = original.ReportGrouping
type ReportListResult = original.ReportListResult
type ReportProperties = original.ReportProperties
type ReportRecurrencePeriod = original.ReportRecurrencePeriod
type ReportSchedule = original.ReportSchedule
type ReportTimePeriod = original.ReportTimePeriod
type Resource = original.Resource
type ReportClient = original.ReportClient
type ReportsClient = original.ReportsClient
type ResourceGroupDimensionsClient = original.ResourceGroupDimensionsClient
type SubscriptionDimensionsClient = original.SubscriptionDimensionsClient

func NewBillingAccountDimensionsClient(subscriptionID string) BillingAccountDimensionsClient {
	return original.NewBillingAccountDimensionsClient(subscriptionID)
}
func NewBillingAccountDimensionsClientWithBaseURI(baseURI string, subscriptionID string) BillingAccountDimensionsClient {
	return original.NewBillingAccountDimensionsClientWithBaseURI(baseURI, subscriptionID)
}
func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func NewConnectorClient(subscriptionID string) ConnectorClient {
	return original.NewConnectorClient(subscriptionID)
}
func NewConnectorClientWithBaseURI(baseURI string, subscriptionID string) ConnectorClient {
	return original.NewConnectorClientWithBaseURI(baseURI, subscriptionID)
}
func PossibleConnectorStatusValues() []ConnectorStatus {
	return original.PossibleConnectorStatusValues()
}
func PossibleExecutionStatusValues() []ExecutionStatus {
	return original.PossibleExecutionStatusValues()
}
func PossibleExecutionTypeValues() []ExecutionType {
	return original.PossibleExecutionTypeValues()
}
func PossibleFormatTypeValues() []FormatType {
	return original.PossibleFormatTypeValues()
}
func PossibleGranularityTypeValues() []GranularityType {
	return original.PossibleGranularityTypeValues()
}
func PossibleRecurrenceTypeValues() []RecurrenceType {
	return original.PossibleRecurrenceTypeValues()
}
func PossibleReportColumnTypeValues() []ReportColumnType {
	return original.PossibleReportColumnTypeValues()
}
func PossibleStatusTypeValues() []StatusType {
	return original.PossibleStatusTypeValues()
}
func PossibleTimeframeTypeValues() []TimeframeType {
	return original.PossibleTimeframeTypeValues()
}
func NewReportClient(subscriptionID string) ReportClient {
	return original.NewReportClient(subscriptionID)
}
func NewReportClientWithBaseURI(baseURI string, subscriptionID string) ReportClient {
	return original.NewReportClientWithBaseURI(baseURI, subscriptionID)
}
func NewReportsClient(subscriptionID string) ReportsClient {
	return original.NewReportsClient(subscriptionID)
}
func NewReportsClientWithBaseURI(baseURI string, subscriptionID string) ReportsClient {
	return original.NewReportsClientWithBaseURI(baseURI, subscriptionID)
}
func NewResourceGroupDimensionsClient(subscriptionID string) ResourceGroupDimensionsClient {
	return original.NewResourceGroupDimensionsClient(subscriptionID)
}
func NewResourceGroupDimensionsClientWithBaseURI(baseURI string, subscriptionID string) ResourceGroupDimensionsClient {
	return original.NewResourceGroupDimensionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewSubscriptionDimensionsClient(subscriptionID string) SubscriptionDimensionsClient {
	return original.NewSubscriptionDimensionsClient(subscriptionID)
}
func NewSubscriptionDimensionsClientWithBaseURI(baseURI string, subscriptionID string) SubscriptionDimensionsClient {
	return original.NewSubscriptionDimensionsClientWithBaseURI(baseURI, subscriptionID)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
