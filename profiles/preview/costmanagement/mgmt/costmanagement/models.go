// +build go1.9

// Copyright 2020 Microsoft Corporation
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

import (
	"context"

	original "github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2019-10-01/costmanagement"
)

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type ExecutionStatus = original.ExecutionStatus

const (
	Completed           ExecutionStatus = original.Completed
	DataNotAvailable    ExecutionStatus = original.DataNotAvailable
	Failed              ExecutionStatus = original.Failed
	InProgress          ExecutionStatus = original.InProgress
	NewDataNotAvailable ExecutionStatus = original.NewDataNotAvailable
	Queued              ExecutionStatus = original.Queued
	Timeout             ExecutionStatus = original.Timeout
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
	Daily  GranularityType = original.Daily
	Hourly GranularityType = original.Hourly
)

type QueryColumnType = original.QueryColumnType

const (
	QueryColumnTypeDimension QueryColumnType = original.QueryColumnTypeDimension
	QueryColumnTypeTag       QueryColumnType = original.QueryColumnTypeTag
)

type RecurrenceType = original.RecurrenceType

const (
	RecurrenceTypeAnnually RecurrenceType = original.RecurrenceTypeAnnually
	RecurrenceTypeDaily    RecurrenceType = original.RecurrenceTypeDaily
	RecurrenceTypeMonthly  RecurrenceType = original.RecurrenceTypeMonthly
	RecurrenceTypeWeekly   RecurrenceType = original.RecurrenceTypeWeekly
)

type SortDirection = original.SortDirection

const (
	Ascending  SortDirection = original.Ascending
	Descending SortDirection = original.Descending
)

type StatusType = original.StatusType

const (
	Active   StatusType = original.Active
	Inactive StatusType = original.Inactive
)

type TimeframeType = original.TimeframeType

const (
	Custom       TimeframeType = original.Custom
	MonthToDate  TimeframeType = original.MonthToDate
	TheLastMonth TimeframeType = original.TheLastMonth
	TheLastWeek  TimeframeType = original.TheLastWeek
	TheLastYear  TimeframeType = original.TheLastYear
	WeekToDate   TimeframeType = original.WeekToDate
	YearToDate   TimeframeType = original.YearToDate
)

type BaseClient = original.BaseClient
type CommonExportProperties = original.CommonExportProperties
type Dimension = original.Dimension
type DimensionProperties = original.DimensionProperties
type DimensionsClient = original.DimensionsClient
type DimensionsListResult = original.DimensionsListResult
type ErrorDetails = original.ErrorDetails
type ErrorResponse = original.ErrorResponse
type Export = original.Export
type ExportDeliveryDestination = original.ExportDeliveryDestination
type ExportDeliveryInfo = original.ExportDeliveryInfo
type ExportExecution = original.ExportExecution
type ExportExecutionListResult = original.ExportExecutionListResult
type ExportExecutionProperties = original.ExportExecutionProperties
type ExportListResult = original.ExportListResult
type ExportProperties = original.ExportProperties
type ExportRecurrencePeriod = original.ExportRecurrencePeriod
type ExportSchedule = original.ExportSchedule
type ExportsClient = original.ExportsClient
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationListResult = original.OperationListResult
type OperationListResultIterator = original.OperationListResultIterator
type OperationListResultPage = original.OperationListResultPage
type OperationsClient = original.OperationsClient
type QueryAggregation = original.QueryAggregation
type QueryClient = original.QueryClient
type QueryColumn = original.QueryColumn
type QueryComparisonExpression = original.QueryComparisonExpression
type QueryDataset = original.QueryDataset
type QueryDatasetConfiguration = original.QueryDatasetConfiguration
type QueryDefinition = original.QueryDefinition
type QueryFilter = original.QueryFilter
type QueryGrouping = original.QueryGrouping
type QueryProperties = original.QueryProperties
type QueryResult = original.QueryResult
type QuerySortingConfiguration = original.QuerySortingConfiguration
type QueryTimePeriod = original.QueryTimePeriod
type Resource = original.Resource

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewDimensionsClient(subscriptionID string) DimensionsClient {
	return original.NewDimensionsClient(subscriptionID)
}
func NewDimensionsClientWithBaseURI(baseURI string, subscriptionID string) DimensionsClient {
	return original.NewDimensionsClientWithBaseURI(baseURI, subscriptionID)
}
func NewExportsClient(subscriptionID string) ExportsClient {
	return original.NewExportsClient(subscriptionID)
}
func NewExportsClientWithBaseURI(baseURI string, subscriptionID string) ExportsClient {
	return original.NewExportsClientWithBaseURI(baseURI, subscriptionID)
}
func NewOperationListResultIterator(page OperationListResultPage) OperationListResultIterator {
	return original.NewOperationListResultIterator(page)
}
func NewOperationListResultPage(cur OperationListResult, getNextPage func(context.Context, OperationListResult) (OperationListResult, error)) OperationListResultPage {
	return original.NewOperationListResultPage(cur, getNextPage)
}
func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewQueryClient(subscriptionID string) QueryClient {
	return original.NewQueryClient(subscriptionID)
}
func NewQueryClientWithBaseURI(baseURI string, subscriptionID string) QueryClient {
	return original.NewQueryClientWithBaseURI(baseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
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
func PossibleQueryColumnTypeValues() []QueryColumnType {
	return original.PossibleQueryColumnTypeValues()
}
func PossibleRecurrenceTypeValues() []RecurrenceType {
	return original.PossibleRecurrenceTypeValues()
}
func PossibleSortDirectionValues() []SortDirection {
	return original.PossibleSortDirectionValues()
}
func PossibleStatusTypeValues() []StatusType {
	return original.PossibleStatusTypeValues()
}
func PossibleTimeframeTypeValues() []TimeframeType {
	return original.PossibleTimeframeTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
