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

package job

import original "github.com/Azure/azure-sdk-for-go/services/datalake/analytics/2017-09-01-preview/job"

type PipelineClient = original.PipelineClient
type Client = original.Client
type CompileMode = original.CompileMode

const (
	Full      CompileMode = original.Full
	Semantic  CompileMode = original.Semantic
	SingleBox CompileMode = original.SingleBox
)

type ResourceType = original.ResourceType

const (
	JobManagerResource             ResourceType = original.JobManagerResource
	JobManagerResourceInUserFolder ResourceType = original.JobManagerResourceInUserFolder
	StatisticsResource             ResourceType = original.StatisticsResource
	StatisticsResourceInUserFolder ResourceType = original.StatisticsResourceInUserFolder
	VertexResource                 ResourceType = original.VertexResource
	VertexResourceInUserFolder     ResourceType = original.VertexResourceInUserFolder
)

type Result = original.Result

const (
	Cancelled Result = original.Cancelled
	Failed    Result = original.Failed
	None      Result = original.None
	Succeeded Result = original.Succeeded
)

type SeverityTypes = original.SeverityTypes

const (
	Deprecated    SeverityTypes = original.Deprecated
	Error         SeverityTypes = original.Error
	Info          SeverityTypes = original.Info
	SevereWarning SeverityTypes = original.SevereWarning
	UserWarning   SeverityTypes = original.UserWarning
	Warning       SeverityTypes = original.Warning
)

type State = original.State

const (
	StateAccepted           State = original.StateAccepted
	StateCompiling          State = original.StateCompiling
	StateEnded              State = original.StateEnded
	StateNew                State = original.StateNew
	StatePaused             State = original.StatePaused
	StateQueued             State = original.StateQueued
	StateRunning            State = original.StateRunning
	StateScheduling         State = original.StateScheduling
	StateStarting           State = original.StateStarting
	StateWaitingForCapacity State = original.StateWaitingForCapacity
)

type Type = original.Type

const (
	TypeHive          Type = original.TypeHive
	TypeJobProperties Type = original.TypeJobProperties
	TypeScope         Type = original.TypeScope
	TypeUSQL          Type = original.TypeUSQL
)

type TypeBasicCreateJobProperties = original.TypeBasicCreateJobProperties

const (
	TypeBasicCreateJobPropertiesTypeCreateJobProperties TypeBasicCreateJobProperties = original.TypeBasicCreateJobPropertiesTypeCreateJobProperties
	TypeBasicCreateJobPropertiesTypeScope               TypeBasicCreateJobProperties = original.TypeBasicCreateJobPropertiesTypeScope
	TypeBasicCreateJobPropertiesTypeUSQL                TypeBasicCreateJobProperties = original.TypeBasicCreateJobPropertiesTypeUSQL
)

type TypeEnum = original.TypeEnum

const (
	Hive  TypeEnum = original.Hive
	Scope TypeEnum = original.Scope
	USQL  TypeEnum = original.USQL
)

type BaseJobParameters = original.BaseJobParameters
type BuildJobParameters = original.BuildJobParameters
type CancelFuture = original.CancelFuture
type CreateJobParameters = original.CreateJobParameters
type BasicCreateJobProperties = original.BasicCreateJobProperties
type CreateJobProperties = original.CreateJobProperties
type CreateScopeJobParameters = original.CreateScopeJobParameters
type CreateScopeJobProperties = original.CreateScopeJobProperties
type CreateUSQLJobProperties = original.CreateUSQLJobProperties
type DataPath = original.DataPath
type Diagnostics = original.Diagnostics
type ErrorDetails = original.ErrorDetails
type HiveJobProperties = original.HiveJobProperties
type InfoListResult = original.InfoListResult
type InfoListResultIterator = original.InfoListResultIterator
type InfoListResultPage = original.InfoListResultPage
type Information = original.Information
type InformationBasic = original.InformationBasic
type InnerError = original.InnerError
type PipelineInformation = original.PipelineInformation
type PipelineInformationListResult = original.PipelineInformationListResult
type PipelineInformationListResultIterator = original.PipelineInformationListResultIterator
type PipelineInformationListResultPage = original.PipelineInformationListResultPage
type PipelineRunInformation = original.PipelineRunInformation
type BasicProperties = original.BasicProperties
type Properties = original.Properties
type RecurrenceInformation = original.RecurrenceInformation
type RecurrenceInformationListResult = original.RecurrenceInformationListResult
type RecurrenceInformationListResultIterator = original.RecurrenceInformationListResultIterator
type RecurrenceInformationListResultPage = original.RecurrenceInformationListResultPage
type RelationshipProperties = original.RelationshipProperties
type Resource = original.Resource
type ResourceUsageStatistics = original.ResourceUsageStatistics
type ScopeJobProperties = original.ScopeJobProperties
type ScopeJobResource = original.ScopeJobResource
type StateAuditRecord = original.StateAuditRecord
type Statistics = original.Statistics
type StatisticsVertex = original.StatisticsVertex
type StatisticsVertexStage = original.StatisticsVertexStage
type UpdateFuture = original.UpdateFuture
type UpdateJobParameters = original.UpdateJobParameters
type USQLJobProperties = original.USQLJobProperties
type YieldFuture = original.YieldFuture
type RecurrenceClient = original.RecurrenceClient

const (
	DefaultAdlaJobDNSSuffix = original.DefaultAdlaJobDNSSuffix
)

type BaseClient = original.BaseClient

func PossibleCompileModeValues() []CompileMode {
	return original.PossibleCompileModeValues()
}
func PossibleResourceTypeValues() []ResourceType {
	return original.PossibleResourceTypeValues()
}
func PossibleResultValues() []Result {
	return original.PossibleResultValues()
}
func PossibleSeverityTypesValues() []SeverityTypes {
	return original.PossibleSeverityTypesValues()
}
func PossibleStateValues() []State {
	return original.PossibleStateValues()
}
func PossibleTypeValues() []Type {
	return original.PossibleTypeValues()
}
func PossibleTypeBasicCreateJobPropertiesValues() []TypeBasicCreateJobProperties {
	return original.PossibleTypeBasicCreateJobPropertiesValues()
}
func PossibleTypeEnumValues() []TypeEnum {
	return original.PossibleTypeEnumValues()
}
func NewRecurrenceClient() RecurrenceClient {
	return original.NewRecurrenceClient()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
func New() BaseClient {
	return original.New()
}
func NewWithoutDefaults(adlaJobDNSSuffix string) BaseClient {
	return original.NewWithoutDefaults(adlaJobDNSSuffix)
}
func NewPipelineClient() PipelineClient {
	return original.NewPipelineClient()
}
func NewClient() Client {
	return original.NewClient()
}
