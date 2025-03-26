// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armregulatedenvironmentmanagement

import "time"

// CreateLZConfigurationCopyRequest - The request for create duplicate landing zone configuration.
type CreateLZConfigurationCopyRequest struct {
	// REQUIRED; The name of the duplicate landing zone configuration resource.
	Name *string
}

// CreateLZConfigurationCopyResult - The response of the create duplicate landing zone configuration.
type CreateLZConfigurationCopyResult struct {
	// REQUIRED; The ID of the duplicate landing zone configuration resource.
	CopiedLandingZoneConfigurationID *string
}

// CustomNamingConvention - The details for the custom naming convention override for a specific resource type.
type CustomNamingConvention struct {
	// REQUIRED; The custom naming formula for the resource type.
	Formula *string

	// REQUIRED; The type of the resource.
	ResourceType *ResourceType
}

// DecommissionedManagementGroupProperties - The 'Decommissioned' management group properties.
type DecommissionedManagementGroupProperties struct {
	// REQUIRED; This parameter determines whether the 'Decommissioned' management group will be created. If set to true, the
	// group will be created; if set to false, it will not be created. The default value is false.
	Create *bool

	// REQUIRED; Array of policy initiatives applied to the management group.
	PolicyInitiativesAssignmentProperties []*PolicyInitiativeAssignmentProperties
}

// GenerateLandingZoneRequest - The request to generate Infrastructure as Code (IaC) for a landing zone.
type GenerateLandingZoneRequest struct {
	// REQUIRED; The Azure region where the landing zone will be deployed. All Azure regions are supported.
	DeploymentLocation *string

	// REQUIRED; The prefix that will be added to all resources created by this deployment. Use between 2 and 5 characters, consisting
	// only of letters, digits, '-', '.', or '_'. No other special characters are supported.
	DeploymentPrefix *string

	// REQUIRED; The export options available for code generation.
	InfrastructureAsCodeOutputOptions *InfrastructureAsCodeOutputOptions

	// REQUIRED; The display name assigned to the top management group of the landing zone deployment hierarchy. It is recommended
	// to use unique names for each landing zone deployment.
	TopLevelMgDisplayName *string

	// The optional suffix that will be appended to all resources created by this deployment, maximum 5 characters.
	DeploymentSuffix *string

	// The environment where the landing zone is being deployed, such as ppe, prod, test, etc.
	Environment *string

	// Existing 'Connectivity' subscription ID to be linked with this deployment when reusing instead of creating a new subscription.
	ExistingConnectivitySubscriptionID *string

	// Existing 'Identity' subscription ID to be linked with this deployment when reusing instead of creating a new subscription.
	ExistingIdentitySubscriptionID *string

	// Existing 'Management' subscription ID to be linked with this deployment when reusing instead of creating a new subscription.
	ExistingManagementSubscriptionID *string

	// Optional parent for the management group hierarchy, serving as an intermediate root management group parent if specified.
	// If left empty, the default will be to deploy under the tenant root management group.
	ExistingTopLevelMgParentID *string

	// The name of the organization or agency for which the landing zone is being deployed. This is optional.
	Organization *string

	// The complete resource ID of the billing scope linked to the EA, MCA, or MPA account where you want to create the subscription.
	SubscriptionBillingScope *string
}

// GenerateLandingZoneResult - The response payload for generating infrastructure-as-code for the landing zone.
type GenerateLandingZoneResult struct {
	// REQUIRED; The storage account blob name to access the generated code.
	BlobName *string

	// REQUIRED; The storage account container to access the generated code.
	ContainerName *string

	// REQUIRED; The url to access the generated code.
	GeneratedCodeURL *string

	// REQUIRED; The name of the Landing zone configuration resource.
	LandingZoneConfigurationName *string

	// REQUIRED; The storage account name to access the generated code.
	StorageAccountName *string

	// REQUIRED; The parent management group name of the landing zone deployment.
	TopLevelMgDisplayName *string

	// The generated code content in JSON string format.
	GeneratedArmTemplate *string
}

// LZAccount - The Landing zone account resource type. A Landing zone account is the container for configuring, deploying
// and managing multiple landing zones.
type LZAccount struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// READ-ONLY; The landing zone account.
	Name *string

	// The managed service identities assigned to this resource.
	Identity *ManagedServiceIdentity

	// The resource-specific properties for this resource.
	Properties *LZAccountProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// LZAccountProperties - The properties of landing zone account resource type.
type LZAccountProperties struct {
	// REQUIRED; The storage account that will host the generated infrastructure as code (IaC) for a landing zone deployment.
	StorageAccount *string

	// READ-ONLY; The state that reflects the current stage in the creation, updating, or deletion process of the landing zone
	// account.
	ProvisioningState *ProvisioningState
}

// LZConfiguration - Concrete proxy resource types can be created by aliasing this type using a specific property type.
type LZConfiguration struct {
	// The resource-specific properties for this resource.
	Properties *LZConfigurationProperties

	// READ-ONLY; The landing zone configuration name
	Name *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// LZConfigurationProperties - The properties of landing zone configuration resource type.
type LZConfigurationProperties struct {
	// REQUIRED; Parameter used to deploy a Bastion: Select 'Yes' to enable deployment, 'No' to skip it, or 'Existing' to reuse
	// an existing Bastion.
	AzureBastionCreationOption *ResourceCreationOptions

	// REQUIRED; Parameter used to deploy a DDoS protection plan: Select 'Yes' to enable deployment, 'No' to skip it, or 'Existing'
	// to reuse an existing DDoS protection plan.
	DdosProtectionCreationOption *ResourceCreationOptions

	// REQUIRED; Parameter used for deploying a Firewall: Select 'No' to skip deployment, 'Standard' to deploy the Standard SKU,
	// or 'Premium' to deploy the Premium SKU.
	FirewallCreationOption *FirewallCreationOptions

	// REQUIRED; The gateway subnet address used for deploying a virtual network. Specify the subnet using IPv4 CIDR notation.
	GatewaySubnetCidrBlock *string

	// REQUIRED; The Virtual Network address. Specify the address using IPv4 CIDR notation.
	HubNetworkCidrBlock *string

	// REQUIRED; Parameter used to deploy a log analytics workspace: Select 'Yes' to enable deployment, 'No' to skip it, or 'Existing'
	// to reuse an existing log analytics workspace.
	LogAnalyticsWorkspaceCreationOption *ResourceCreationOptions

	// REQUIRED; Parameter to define the retention period for logs, in days. The minimum duration is 30 days and the maximum is
	// 730 days.
	LogRetentionInDays *int64

	// REQUIRED; The managed identity to be assigned to this landing zone configuration.
	ManagedIdentity *ManagedIdentityProperties

	// The Bastion subnet address. Specify the address using IPv4 CIDR notation.
	AzureBastionSubnetCidrBlock *string

	// The custom naming convention applied to specific resource types for this landing zone configuration, which overrides the
	// default naming convention for those resource types. Example - 'customNamingConvention': [{'resourceType': 'azureFirewalls',
	// 'formula': '{DeploymentPrefix}-afwl-{DeploymentSuffix}'}]
	CustomNamingConvention []*CustomNamingConvention

	// The assigned policies of the 'Decommissioned' management group and indicator to create it or not.
	DecommissionedMgMetadata *DecommissionedManagementGroupProperties

	// The resource ID of the Bastion when reusing an existing one.
	ExistingAzureBastionID *string

	// The resource ID of the DDoS protection plan when reusing an existing one.
	ExistingDdosProtectionID *string

	// The resource ID of the log analytics workspace when reusing an existing one.
	ExistingLogAnalyticsWorkspaceID *string

	// The Firewall subnet address used for deploying a firewall. Specify the Firewall subnet using IPv4 CIDR notation.
	FirewallSubnetCidrBlock *string

	// The child management groups of 'Landing Zones' management group and their assigned policies.
	LandingZonesMgChildren []*LZManagementGroupProperties

	// The assigned policies of the 'Landing Zones' management group.
	LandingZonesMgMetadata *ManagementGroupProperties

	// The default naming convention applied to all resources for this landing zone configuration. Example - {DeploymentPrefix}-Contoso-{ResourceTypeAbbreviation}{DeploymentSuffix}-{Environment}-testing
	NamingConventionFormula *string

	// The assigned policies of the 'Connectivity' management group under 'Platform' management group.
	PlatformConnectivityMgMetadata *ManagementGroupProperties

	// The assigned policies of the 'Identity' management group under 'Platform' management group.
	PlatformIdentityMgMetadata *ManagementGroupProperties

	// The assigned policies of the 'Management' management group under 'Platform' management group.
	PlatformManagementMgMetadata *ManagementGroupProperties

	// The names of the 'Platform' child management groups and their assigned policies, excluding the default ones: 'Connectivity',
	// 'Identity', and 'Management'
	PlatformMgChildren []*PlatformManagementGroupProperties

	// The assigned policies of the 'Platform' management group.
	PlatformMgMetadata *ManagementGroupProperties

	// The assigned policies of the 'Sandbox' management group and indicator to create it or not.
	SandboxMgMetadata *SandboxManagementGroupProperties

	// Tags are key-value pairs that can be assigned to a resource to organize and manage it more effectively. Example: {'name':
	// 'a tag name', 'value': 'a tag value'}
	Tags []*Tags

	// The assigned policies of the parent management group.
	TopLevelMgMetadata *ManagementGroupProperties

	// READ-ONLY; The status that indicates the current phase of the configuration process for a deployment.
	AuthoringStatus *AuthoringStatus

	// READ-ONLY; The state that reflects the current stage in the creation, updating, or deletion process of the landing zone
	// configuration.
	ProvisioningState *ProvisioningState
}

// LZManagementGroupProperties - The 'Landing Zones' management group properties..
type LZManagementGroupProperties struct {
	// REQUIRED; Management group name.
	Name *string

	// REQUIRED; Array of policy initiatives applied to the management group.
	PolicyInitiativesAssignmentProperties []*PolicyInitiativeAssignmentProperties
}

// LZRegistration - The Landing zone registration resource type.
type LZRegistration struct {
	// The resource-specific properties for this resource.
	Properties *LZRegistrationProperties

	// READ-ONLY; The name of the landing zone registration resource.
	Name *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// LZRegistrationProperties - The properties of landing zone registration resource type.
type LZRegistrationProperties struct {
	// REQUIRED; The resource id of the associated landing zone configuration.
	ExistingLandingZoneConfigurationID *string

	// REQUIRED; The resource id of the top level management group
	ExistingTopLevelMgID *string

	// The managed identity to be assigned to this landing zone registration.
	ManagedIdentity *ManagedIdentityProperties

	// READ-ONLY; The state that reflects the current stage in the creation, updating, or deletion process of the landing zone
	// registration resource type.
	ProvisioningState *ProvisioningState
}

// LandingZoneAccountResourceListResult - The response of a LandingZoneAccountResource list operation.
type LandingZoneAccountResourceListResult struct {
	// REQUIRED; The LandingZoneAccountResource items on this page
	Value []*LZAccount

	// The link to the next page of items
	NextLink *string
}

// LandingZoneConfigurationResourceListResult - The response of a LandingZoneConfigurationResource list operation.
type LandingZoneConfigurationResourceListResult struct {
	// REQUIRED; The LandingZoneConfigurationResource items on this page
	Value []*LZConfiguration

	// The link to the next page of items
	NextLink *string
}

// LandingZoneRegistrationResourceListResult - The response of a LandingZoneRegistrationResource list operation.
type LandingZoneRegistrationResourceListResult struct {
	// REQUIRED; The LandingZoneRegistrationResource items on this page
	Value []*LZRegistration

	// The link to the next page of items
	NextLink *string
}

// ManagedIdentityProperties - The properties of managed identity, specifically including type and resource ID.
type ManagedIdentityProperties struct {
	// REQUIRED; The type of managed identity.
	Type *ManagedIdentityResourceType

	// The resource id of the managed identity.
	UserAssignedIdentityResourceID *string
}

// ManagedServiceIdentity - Managed service identity (system assigned and/or user assigned identities)
type ManagedServiceIdentity struct {
	// REQUIRED; The type of managed identity assigned to this resource.
	Type *ManagedServiceIdentityType

	// The identities assigned to this resource by the user.
	UserAssignedIdentities map[string]*UserAssignedIdentity

	// READ-ONLY; The service principal ID of the system assigned identity. This property will only be provided for a system assigned
	// identity.
	PrincipalID *string

	// READ-ONLY; The tenant ID of the system assigned identity. This property will only be provided for a system assigned identity.
	TenantID *string
}

// ManagementGroupProperties - The properties of policy initiatives applied to the management group.
type ManagementGroupProperties struct {
	// REQUIRED; Array of policy initiatives applied to the management group.
	PolicyInitiativesAssignmentProperties []*PolicyInitiativeAssignmentProperties
}

// Operation - Details of a REST API operation, returned from the Resource Provider Operations API
type Operation struct {
	// Localized display information for this particular operation.
	Display *OperationDisplay

	// READ-ONLY; Extensible enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
	ActionType *ActionType

	// READ-ONLY; Whether the operation applies to data-plane. This is "true" for data-plane operations and "false" for Azure
	// Resource Manager/control-plane operations.
	IsDataAction *bool

	// READ-ONLY; The name of the operation, as per Resource-Based Access Control (RBAC). Examples: "Microsoft.Compute/virtualMachines/write",
	// "Microsoft.Compute/virtualMachines/capture/action"
	Name *string

	// READ-ONLY; The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
	// value is "user,system"
	Origin *Origin
}

// OperationDisplay - Localized display information for and operation.
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
	// REQUIRED; The Operation items on this page
	Value []*Operation

	// The link to the next page of items
	NextLink *string
}

// PlatformManagementGroupProperties - The 'Platform' management group properties.
type PlatformManagementGroupProperties struct {
	// REQUIRED; Management group name.
	Name *string

	// REQUIRED; Array of policy initiatives applied to the management group.
	PolicyInitiativesAssignmentProperties []*PolicyInitiativeAssignmentProperties
}

// PolicyInitiativeAssignmentProperties - The properties of assigned policy initiatives.
type PolicyInitiativeAssignmentProperties struct {
	// REQUIRED; The parameters of the assigned policy initiative.
	AssignmentParameters map[string]any

	// REQUIRED; The fully qualified id of the policy initiative.
	PolicyInitiativeID *string
}

// SandboxManagementGroupProperties - The 'Sandbox' management group properties.
type SandboxManagementGroupProperties struct {
	// REQUIRED; This parameter determines whether the 'Sandbox' management group will be created. If set to true, the group will
	// be created; if set to false, it will not be created. The default value is false.
	Create *bool

	// REQUIRED; Array of policy initiatives applied to the management group.
	PolicyInitiativesAssignmentProperties []*PolicyInitiativeAssignmentProperties
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

// Tags - Key-value pairs that can be assigned to this resource.
type Tags struct {
	// REQUIRED; A tag name.
	Name *string

	// A tag value.
	Value *string
}

// UpdateAuthoringStatusRequest - The request to update the authoring status of a configuration.
type UpdateAuthoringStatusRequest struct {
	// REQUIRED; The authoring status value to be updated. Possible values include: 'Authoring', 'ReadyForUse' and 'Disabled'.
	AuthoringStatus *AuthoringStatus
}

// UpdateAuthoringStatusResult - The response for authoring status update request.
type UpdateAuthoringStatusResult struct {
	// REQUIRED; The authoring status value to be updated.
	AuthoringStatus *AuthoringStatus

	// REQUIRED; The name of the landing zone configuration resource.
	LandingZoneConfigurationName *string
}

// UserAssignedIdentity - User assigned identity properties
type UserAssignedIdentity struct {
	// READ-ONLY; The client ID of the assigned identity.
	ClientID *string

	// READ-ONLY; The principal ID of the assigned identity.
	PrincipalID *string
}
