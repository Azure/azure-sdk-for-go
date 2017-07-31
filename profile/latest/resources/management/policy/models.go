package policy

import (
	 original "github.com/Azure/azure-sdk-for-go/service/resources/management/2016-12-01/policy"
)

type (
	 AssignmentsClient = original.AssignmentsClient
	 ManagementClient = original.ManagementClient
	 DefinitionsClient = original.DefinitionsClient
	 Type = original.Type
	 Assignment = original.Assignment
	 AssignmentListResult = original.AssignmentListResult
	 AssignmentProperties = original.AssignmentProperties
	 Definition = original.Definition
	 DefinitionListResult = original.DefinitionListResult
	 DefinitionProperties = original.DefinitionProperties
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 BuiltIn = original.BuiltIn
	 Custom = original.Custom
	 NotSpecified = original.NotSpecified
)
