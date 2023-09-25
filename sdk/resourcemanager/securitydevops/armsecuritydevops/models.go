//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecuritydevops

import "time"

type ActionableRemediation struct {
	// Branch onboarding info.
	BranchConfiguration *TargetBranchConfiguration
	Categories []*RuleCategory
	SeverityLevels []*string
	State *ActionableRemediationState
}

type AuthorizationInfo struct {
	// Gets or sets one-time OAuth code to exchange for refresh and access tokens.
	Code *string
}

type AzureDevOpsConnector struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string
	Properties *AzureDevOpsConnectorProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type AzureDevOpsConnectorListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*AzureDevOpsConnector
}

type AzureDevOpsConnectorProperties struct {
	Authorization *AuthorizationInfo

	// Gets or sets org onboarding information.
	Orgs []*AzureDevOpsOrgMetadata
	ProvisioningState *ProvisioningState
}

type AzureDevOpsConnectorStats struct {
	Properties *AzureDevOpsConnectorStatsProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type AzureDevOpsConnectorStatsListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*AzureDevOpsConnectorStats
}

type AzureDevOpsConnectorStatsProperties struct {
	// Gets or sets orgs count.
	OrgsCount *int64

	// Gets or sets projects count.
	ProjectsCount *int64
	ProvisioningState *ProvisioningState

	// Gets or sets repos count.
	ReposCount *int64
}

// AzureDevOpsOrg - Azure DevOps Org Proxy Resource.
type AzureDevOpsOrg struct {
	// AzureDevOps Org properties.
	Properties *AzureDevOpsOrgProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type AzureDevOpsOrgListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*AzureDevOpsOrg
}

// AzureDevOpsOrgMetadata - Org onboarding info.
type AzureDevOpsOrgMetadata struct {
	AutoDiscovery *AutoDiscovery

	// Gets or sets name of the AzureDevOps Org.
	Name *string
	Projects []*AzureDevOpsProjectMetadata
}

// AzureDevOpsOrgProperties - AzureDevOps Org properties.
type AzureDevOpsOrgProperties struct {
	AutoDiscovery *AutoDiscovery
	ProvisioningState *ProvisioningState
}

// AzureDevOpsProject - Azure DevOps Project Proxy Resource.
type AzureDevOpsProject struct {
	// AzureDevOps Project properties.
	Properties *AzureDevOpsProjectProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type AzureDevOpsProjectListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*AzureDevOpsProject
}

// AzureDevOpsProjectMetadata - Project onboarding info.
type AzureDevOpsProjectMetadata struct {
	AutoDiscovery *AutoDiscovery

	// Gets or sets name of the AzureDevOps Project.
	Name *string

	// Gets or sets repositories.
	Repos []*string
}

// AzureDevOpsProjectProperties - AzureDevOps Project properties.
type AzureDevOpsProjectProperties struct {
	AutoDiscovery *AutoDiscovery

	// Gets or sets AzureDevOps Org Name.
	OrgName *string

	// Gets or sets AzureDevOps Project Id.
	ProjectID *string
	ProvisioningState *ProvisioningState
}

// AzureDevOpsRepo - Azure DevOps Repo Proxy Resource.
type AzureDevOpsRepo struct {
	// AzureDevOps Repo properties.
	Properties *AzureDevOpsRepoProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type AzureDevOpsRepoListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*AzureDevOpsRepo
}

// AzureDevOpsRepoProperties - AzureDevOps Repo properties.
type AzureDevOpsRepoProperties struct {
	ActionableRemediation *ActionableRemediation

	// Gets or sets AzureDevOps Org Name.
	OrgName *string

	// Gets or sets AzureDevOps Project Name.
	ProjectName *string
	ProvisioningState *ProvisioningState

	// Gets or sets Azure DevOps repo id.
	RepoID *string

	// Gets or sets AzureDevOps repo url.
	RepoURL *string

	// Gets or sets AzureDevOps repo visibility, whether it is public or private etc.
	Visibility *string
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info any

	// READ-ONLY; The additional info type.
	Type *string
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo

	// READ-ONLY; The error code.
	Code *string

	// READ-ONLY; The error details.
	Details []*ErrorDetail

	// READ-ONLY; The error message.
	Message *string

	// READ-ONLY; The error target.
	Target *string
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations.
// (This also follows the OData error response format.).
type ErrorResponse struct {
	// The error object.
	Error *ErrorDetail
}

// GitHubConnector - Represents an ARM resource for /subscriptions/xxx/resourceGroups/xxx/providers/Microsoft.SecurityDevOps/gitHubConnectors.
type GitHubConnector struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// Properties of the ARM resource for /subscriptions/xxx/resourceGroups/xxx/providers/Microsoft.SecurityDevOps/gitHubConnectors.
	Properties *GitHubConnectorProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type GitHubConnectorListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*GitHubConnector
}

// GitHubConnectorProperties - Properties of the ARM resource for /subscriptions/xxx/resourceGroups/xxx/providers/Microsoft.SecurityDevOps/gitHubConnectors.
type GitHubConnectorProperties struct {
	// Gets or sets one-time OAuth code to exchange for refresh and access tokens.
	Code *string
	ProvisioningState *ProvisioningState
}

type GitHubConnectorStats struct {
	Properties *GitHubConnectorStatsProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type GitHubConnectorStatsListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*GitHubConnectorStats
}

type GitHubConnectorStatsProperties struct {
	// Gets or sets owners count.
	OwnersCount *int64
	ProvisioningState *ProvisioningState

	// Gets or sets repos count.
	ReposCount *int64
}

// GitHubOwner - GitHub repo owner Proxy Resource.
type GitHubOwner struct {
	// GitHub Repo Owner properties.
	Properties *GitHubOwnerProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type GitHubOwnerListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*GitHubOwner
}

// GitHubOwnerProperties - GitHub Repo Owner properties.
type GitHubOwnerProperties struct {
	// Gets or sets gitHub owner url.
	OwnerURL *string
	ProvisioningState *ProvisioningState
}

// GitHubRepo - GitHub repo Proxy Resource.
type GitHubRepo struct {
	// GitHub Repo properties.
	Properties *GitHubRepoProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type GitHubRepoListResponse struct {
	// Gets or sets next link to scroll over the results.
	NextLink *string

	// Gets or sets list of resources.
	Value []*GitHubRepo
}

// GitHubRepoProperties - GitHub Repo properties.
type GitHubRepoProperties struct {
	// Gets or sets gitHub repo account id.
	AccountID *int64

	// Gets or sets GitHub Owner Name.
	OwnerName *string
	ProvisioningState *ProvisioningState

	// Gets or sets gitHub repo url.
	RepoURL *string
}

type GitHubReposProperties struct {
	// Gets or sets gitHub repo account id.
	AccountID *int64
	ProvisioningState *ProvisioningState

	// Gets or sets gitHub repo name.
	RepoName *string

	// Gets or sets gitHub repo url.
	RepoURL *string
}

// Operation - Details of a REST API operation, returned from the Resource Provider Operations API
type Operation struct {
	// Localized display information for this particular operation.
	Display *OperationDisplay

	// READ-ONLY; Enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
	ActionType *ActionType

	// READ-ONLY; Whether the operation applies to data-plane. This is "true" for data-plane operations and "false" for ARM/control-plane
// operations.
	IsDataAction *bool

	// READ-ONLY; The name of the operation, as per Resource-Based Access Control (RBAC). Examples: "Microsoft.Compute/virtualMachines/write",
// "Microsoft.Compute/virtualMachines/capture/action"
	Name *string

	// READ-ONLY; The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
// value is "user,system"
	Origin *Origin
}

// OperationDisplay - Localized display information for this particular operation.
type OperationDisplay struct {
	// READ-ONLY; The short, localized friendly description of the operation; suitable for tool tips and detailed views.
	Description *string

	// READ-ONLY; The concise, localized friendly name for the operation; suitable for dropdowns. E.g. "Create or Update Virtual
// Machine", "Restart Virtual Machine".
	Operation *string

	// READ-ONLY; The localized friendly form of the resource provider name, e.g. "Microsoft Monitoring Insights" or "Microsoft
// Compute".
	Provider *string

	// READ-ONLY; The localized friendly name of the resource type related to this operation. E.g. "Virtual Machines" or "Job
// Schedule Collections".
	Resource *string
}

// OperationListResult - A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to
// get the next set of results.
type OperationListResult struct {
	// READ-ONLY; URL to get the next set of operation list results (if there are any).
	NextLink *string

	// READ-ONLY; List of operations supported by the resource provider
	Value []*Operation
}

// ProxyResource - The resource model definition for a Azure Resource Manager proxy resource. It will not have tags and a
// location
type ProxyResource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time

	// The identity that created the resource.
	CreatedBy *string

	// The type of identity that created the resource.
	CreatedByType *CreatedByType

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time

	// The identity that last modified the resource.
	LastModifiedBy *string

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType
}

// TargetBranchConfiguration - Branch onboarding info.
type TargetBranchConfiguration struct {
	// Gets or sets branches that should have annotations.
// For Ignite, we will be supporting a single default branch configuration in the UX.
	Names []*string
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

