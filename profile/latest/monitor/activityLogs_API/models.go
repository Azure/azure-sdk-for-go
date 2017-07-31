package activitylogsapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/2015-04-01/activityLogs_API"
)

type (
	 ActivityLogsClient = original.ActivityLogsClient
	 ManagementClient = original.ManagementClient
	 EventLevel = original.EventLevel
	 ErrorResponse = original.ErrorResponse
	 EventData = original.EventData
	 EventDataCollection = original.EventDataCollection
	 HTTPRequestInfo = original.HTTPRequestInfo
	 LocalizableString = original.LocalizableString
	 SenderAuthorization = original.SenderAuthorization
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Critical = original.Critical
	 Error = original.Error
	 Informational = original.Informational
	 Verbose = original.Verbose
	 Warning = original.Warning
)
