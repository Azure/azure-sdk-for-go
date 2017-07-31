package webhook

import (
	 original "github.com/Azure/azure-sdk-for-go/service/automation/management/2015-10-31/webhook"
)

type (
	 ManagementClient = original.ManagementClient
	 CreateOrUpdateParameters = original.CreateOrUpdateParameters
	 CreateOrUpdateProperties = original.CreateOrUpdateProperties
	 ErrorResponse = original.ErrorResponse
	 ListResult = original.ListResult
	 Model = original.Model
	 Properties = original.Properties
	 RunbookAssociationProperty = original.RunbookAssociationProperty
	 String = original.String
	 UpdateParameters = original.UpdateParameters
	 UpdateProperties = original.UpdateProperties
	 GroupClient = original.GroupClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)
