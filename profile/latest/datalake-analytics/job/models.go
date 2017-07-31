package job

import (
	 original "github.com/Azure/azure-sdk-for-go/service/datalake-analytics/2016-11-01/job"
)

type (
	 PipelineClient = original.PipelineClient
	 RecurrenceClient = original.RecurrenceClient
	 ManagementClient = original.ManagementClient
	 GroupClient = original.GroupClient
	 CompileMode = original.CompileMode
	 ResourceType = original.ResourceType
	 Result = original.Result
	 SeverityTypes = original.SeverityTypes
	 State = original.State
	 Type = original.Type
	 BaseJobParameters = original.BaseJobParameters
	 BuildJobParameters = original.BuildJobParameters
	 CreateJobParameters = original.CreateJobParameters
	 CreateJobProperties = original.CreateJobProperties
	 CreateUSQLJobProperties = original.CreateUSQLJobProperties
	 DataPath = original.DataPath
	 Diagnostics = original.Diagnostics
	 ErrorDetails = original.ErrorDetails
	 HiveJobProperties = original.HiveJobProperties
	 InfoListResult = original.InfoListResult
	 Information = original.Information
	 InformationBasic = original.InformationBasic
	 InnerError = original.InnerError
	 PipelineInformation = original.PipelineInformation
	 PipelineInformationListResult = original.PipelineInformationListResult
	 PipelineRunInformation = original.PipelineRunInformation
	 Properties = original.Properties
	 RecurrenceInformation = original.RecurrenceInformation
	 RecurrenceInformationListResult = original.RecurrenceInformationListResult
	 RelationshipProperties = original.RelationshipProperties
	 Resource = original.Resource
	 StateAuditRecord = original.StateAuditRecord
	 Statistics = original.Statistics
	 StatisticsVertexStage = original.StatisticsVertexStage
	 USQLJobProperties = original.USQLJobProperties
)

const (
	 DefaultAdlaJobDNSSuffix = original.DefaultAdlaJobDNSSuffix
	 Full = original.Full
	 Semantic = original.Semantic
	 SingleBox = original.SingleBox
	 JobManagerResource = original.JobManagerResource
	 JobManagerResourceInUserFolder = original.JobManagerResourceInUserFolder
	 StatisticsResource = original.StatisticsResource
	 StatisticsResourceInUserFolder = original.StatisticsResourceInUserFolder
	 VertexResource = original.VertexResource
	 VertexResourceInUserFolder = original.VertexResourceInUserFolder
	 Cancelled = original.Cancelled
	 Failed = original.Failed
	 None = original.None
	 Succeeded = original.Succeeded
	 Deprecated = original.Deprecated
	 Error = original.Error
	 Info = original.Info
	 SevereWarning = original.SevereWarning
	 UserWarning = original.UserWarning
	 Warning = original.Warning
	 StateAccepted = original.StateAccepted
	 StateCompiling = original.StateCompiling
	 StateEnded = original.StateEnded
	 StateNew = original.StateNew
	 StatePaused = original.StatePaused
	 StateQueued = original.StateQueued
	 StateRunning = original.StateRunning
	 StateScheduling = original.StateScheduling
	 StateStarting = original.StateStarting
	 StateWaitingForCapacity = original.StateWaitingForCapacity
	 Hive = original.Hive
	 USQL = original.USQL
)
