package actiongroupsapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/management/2017-04-01/actionGroups_API"
)

type (
	 ActionGroupsClient = original.ActionGroupsClient
	 ManagementClient = original.ManagementClient
	 ReceiverStatus = original.ReceiverStatus
	 ActionGroup = original.ActionGroup
	 ActionGroupList = original.ActionGroupList
	 ActionGroupResource = original.ActionGroupResource
	 EmailReceiver = original.EmailReceiver
	 EnableRequest = original.EnableRequest
	 Resource = original.Resource
	 SmsReceiver = original.SmsReceiver
	 WebhookReceiver = original.WebhookReceiver
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 NotSpecified = original.NotSpecified
)
