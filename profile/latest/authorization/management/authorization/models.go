package authorization

import (
	 original "github.com/Azure/azure-sdk-for-go/service/authorization/management/2015-07-01/authorization"
)

type (
	 RoleDefinitionsClient = original.RoleDefinitionsClient
	 ClassicAdministratorsClient = original.ClassicAdministratorsClient
	 ManagementClient = original.ManagementClient
	 ClassicAdministrator = original.ClassicAdministrator
	 ClassicAdministratorListResult = original.ClassicAdministratorListResult
	 ClassicAdministratorProperties = original.ClassicAdministratorProperties
	 Permission = original.Permission
	 PermissionGetResult = original.PermissionGetResult
	 ProviderOperation = original.ProviderOperation
	 ProviderOperationsMetadata = original.ProviderOperationsMetadata
	 ProviderOperationsMetadataListResult = original.ProviderOperationsMetadataListResult
	 ResourceType = original.ResourceType
	 RoleAssignment = original.RoleAssignment
	 RoleAssignmentCreateParameters = original.RoleAssignmentCreateParameters
	 RoleAssignmentFilter = original.RoleAssignmentFilter
	 RoleAssignmentListResult = original.RoleAssignmentListResult
	 RoleAssignmentProperties = original.RoleAssignmentProperties
	 RoleAssignmentPropertiesWithScope = original.RoleAssignmentPropertiesWithScope
	 RoleDefinition = original.RoleDefinition
	 RoleDefinitionFilter = original.RoleDefinitionFilter
	 RoleDefinitionListResult = original.RoleDefinitionListResult
	 RoleDefinitionProperties = original.RoleDefinitionProperties
	 PermissionsClient = original.PermissionsClient
	 ProviderOperationsMetadataClient = original.ProviderOperationsMetadataClient
	 RoleAssignmentsClient = original.RoleAssignmentsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)
