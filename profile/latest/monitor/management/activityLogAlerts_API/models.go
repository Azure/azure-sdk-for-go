package activitylogalertsapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/management/2017-04-01/activityLogAlerts_API"
)

type (
	 ActivityLogAlertsClient = original.ActivityLogAlertsClient
	 ManagementClient = original.ManagementClient
	 ActivityLogAlert = original.ActivityLogAlert
	 ActivityLogAlertActionGroup = original.ActivityLogAlertActionGroup
	 ActivityLogAlertActionList = original.ActivityLogAlertActionList
	 ActivityLogAlertAllOfCondition = original.ActivityLogAlertAllOfCondition
	 ActivityLogAlertLeafCondition = original.ActivityLogAlertLeafCondition
	 ActivityLogAlertList = original.ActivityLogAlertList
	 ActivityLogAlertPatch = original.ActivityLogAlertPatch
	 ActivityLogAlertPatchBody = original.ActivityLogAlertPatchBody
	 ActivityLogAlertResource = original.ActivityLogAlertResource
	 ErrorResponse = original.ErrorResponse
	 Resource = original.Resource
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)
