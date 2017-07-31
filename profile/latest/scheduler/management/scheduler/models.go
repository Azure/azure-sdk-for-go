package scheduler

import (
	 original "github.com/Azure/azure-sdk-for-go/service/scheduler/management/2016-03-01/scheduler"
)

type (
	 ManagementClient = original.ManagementClient
	 JobCollectionsClient = original.JobCollectionsClient
	 JobsClient = original.JobsClient
	 DayOfWeek = original.DayOfWeek
	 HTTPAuthenticationType = original.HTTPAuthenticationType
	 JobActionType = original.JobActionType
	 JobCollectionState = original.JobCollectionState
	 JobExecutionStatus = original.JobExecutionStatus
	 JobHistoryActionName = original.JobHistoryActionName
	 JobScheduleDay = original.JobScheduleDay
	 JobState = original.JobState
	 RecurrenceFrequency = original.RecurrenceFrequency
	 RetryType = original.RetryType
	 ServiceBusAuthenticationType = original.ServiceBusAuthenticationType
	 ServiceBusTransportType = original.ServiceBusTransportType
	 SkuDefinition = original.SkuDefinition
	 BasicAuthentication = original.BasicAuthentication
	 ClientCertAuthentication = original.ClientCertAuthentication
	 HTTPAuthentication = original.HTTPAuthentication
	 HTTPRequest = original.HTTPRequest
	 JobAction = original.JobAction
	 JobCollectionDefinition = original.JobCollectionDefinition
	 JobCollectionListResult = original.JobCollectionListResult
	 JobCollectionProperties = original.JobCollectionProperties
	 JobCollectionQuota = original.JobCollectionQuota
	 JobDefinition = original.JobDefinition
	 JobErrorAction = original.JobErrorAction
	 JobHistoryDefinition = original.JobHistoryDefinition
	 JobHistoryDefinitionProperties = original.JobHistoryDefinitionProperties
	 JobHistoryFilter = original.JobHistoryFilter
	 JobHistoryListResult = original.JobHistoryListResult
	 JobListResult = original.JobListResult
	 JobMaxRecurrence = original.JobMaxRecurrence
	 JobProperties = original.JobProperties
	 JobRecurrence = original.JobRecurrence
	 JobRecurrenceSchedule = original.JobRecurrenceSchedule
	 JobRecurrenceScheduleMonthlyOccurrence = original.JobRecurrenceScheduleMonthlyOccurrence
	 JobStateFilter = original.JobStateFilter
	 JobStatus = original.JobStatus
	 OAuthAuthentication = original.OAuthAuthentication
	 RetryPolicy = original.RetryPolicy
	 ServiceBusAuthentication = original.ServiceBusAuthentication
	 ServiceBusBrokeredMessageProperties = original.ServiceBusBrokeredMessageProperties
	 ServiceBusMessage = original.ServiceBusMessage
	 ServiceBusQueueMessage = original.ServiceBusQueueMessage
	 ServiceBusTopicMessage = original.ServiceBusTopicMessage
	 Sku = original.Sku
	 StorageQueueMessage = original.StorageQueueMessage
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Friday = original.Friday
	 Monday = original.Monday
	 Saturday = original.Saturday
	 Sunday = original.Sunday
	 Thursday = original.Thursday
	 Tuesday = original.Tuesday
	 Wednesday = original.Wednesday
	 ActiveDirectoryOAuth = original.ActiveDirectoryOAuth
	 Basic = original.Basic
	 ClientCertificate = original.ClientCertificate
	 NotSpecified = original.NotSpecified
	 HTTP = original.HTTP
	 HTTPS = original.HTTPS
	 ServiceBusQueue = original.ServiceBusQueue
	 ServiceBusTopic = original.ServiceBusTopic
	 StorageQueue = original.StorageQueue
	 Deleted = original.Deleted
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 Suspended = original.Suspended
	 Completed = original.Completed
	 Failed = original.Failed
	 Postponed = original.Postponed
	 ErrorAction = original.ErrorAction
	 MainAction = original.MainAction
	 JobScheduleDayFriday = original.JobScheduleDayFriday
	 JobScheduleDayMonday = original.JobScheduleDayMonday
	 JobScheduleDaySaturday = original.JobScheduleDaySaturday
	 JobScheduleDaySunday = original.JobScheduleDaySunday
	 JobScheduleDayThursday = original.JobScheduleDayThursday
	 JobScheduleDayTuesday = original.JobScheduleDayTuesday
	 JobScheduleDayWednesday = original.JobScheduleDayWednesday
	 JobStateCompleted = original.JobStateCompleted
	 JobStateDisabled = original.JobStateDisabled
	 JobStateEnabled = original.JobStateEnabled
	 JobStateFaulted = original.JobStateFaulted
	 Day = original.Day
	 Hour = original.Hour
	 Minute = original.Minute
	 Month = original.Month
	 Week = original.Week
	 Fixed = original.Fixed
	 None = original.None
	 ServiceBusAuthenticationTypeNotSpecified = original.ServiceBusAuthenticationTypeNotSpecified
	 ServiceBusAuthenticationTypeSharedAccessKey = original.ServiceBusAuthenticationTypeSharedAccessKey
	 ServiceBusTransportTypeAMQP = original.ServiceBusTransportTypeAMQP
	 ServiceBusTransportTypeNetMessaging = original.ServiceBusTransportTypeNetMessaging
	 ServiceBusTransportTypeNotSpecified = original.ServiceBusTransportTypeNotSpecified
	 Free = original.Free
	 P10Premium = original.P10Premium
	 P20Premium = original.P20Premium
	 Standard = original.Standard
)
