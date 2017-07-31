package resources

import (
	 original "github.com/Azure/azure-sdk-for-go/service/resources/management/2017-05-10/resources"
)

type (
	 DeploymentsClient = original.DeploymentsClient
	 GroupsClient = original.GroupsClient
	 ManagementClient = original.ManagementClient
	 DeploymentOperationsClient = original.DeploymentOperationsClient
	 DeploymentMode = original.DeploymentMode
	 ResourceIdentityType = original.ResourceIdentityType
	 AliasPathType = original.AliasPathType
	 AliasType = original.AliasType
	 BasicDependency = original.BasicDependency
	 DebugSetting = original.DebugSetting
	 Dependency = original.Dependency
	 Deployment = original.Deployment
	 DeploymentExportResult = original.DeploymentExportResult
	 DeploymentExtended = original.DeploymentExtended
	 DeploymentExtendedFilter = original.DeploymentExtendedFilter
	 DeploymentListResult = original.DeploymentListResult
	 DeploymentOperation = original.DeploymentOperation
	 DeploymentOperationProperties = original.DeploymentOperationProperties
	 DeploymentOperationsListResult = original.DeploymentOperationsListResult
	 DeploymentProperties = original.DeploymentProperties
	 DeploymentPropertiesExtended = original.DeploymentPropertiesExtended
	 DeploymentValidateResult = original.DeploymentValidateResult
	 ExportTemplateRequest = original.ExportTemplateRequest
	 GenericResource = original.GenericResource
	 GenericResourceFilter = original.GenericResourceFilter
	 Group = original.Group
	 GroupExportResult = original.GroupExportResult
	 GroupFilter = original.GroupFilter
	 GroupListResult = original.GroupListResult
	 GroupPatchable = original.GroupPatchable
	 GroupProperties = original.GroupProperties
	 HTTPMessage = original.HTTPMessage
	 Identity = original.Identity
	 ListResult = original.ListResult
	 ManagementErrorWithDetails = original.ManagementErrorWithDetails
	 MoveInfo = original.MoveInfo
	 ParametersLink = original.ParametersLink
	 Plan = original.Plan
	 Provider = original.Provider
	 ProviderListResult = original.ProviderListResult
	 ProviderOperationDisplayProperties = original.ProviderOperationDisplayProperties
	 ProviderResourceType = original.ProviderResourceType
	 Resource = original.Resource
	 Sku = original.Sku
	 SubResource = original.SubResource
	 TagCount = original.TagCount
	 TagDetails = original.TagDetails
	 TagsListResult = original.TagsListResult
	 TagValue = original.TagValue
	 TargetResource = original.TargetResource
	 TemplateLink = original.TemplateLink
	 ProvidersClient = original.ProvidersClient
	 GroupClient = original.GroupClient
	 TagsClient = original.TagsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Complete = original.Complete
	 Incremental = original.Incremental
	 SystemAssigned = original.SystemAssigned
)
