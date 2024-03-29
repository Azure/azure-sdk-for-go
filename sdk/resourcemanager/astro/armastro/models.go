//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armastro

import "time"

// LiftrBaseDataOrganizationProperties - Properties specific to Data Organization resource
type LiftrBaseDataOrganizationProperties struct {
	// REQUIRED; Marketplace details of the resource.
	Marketplace *LiftrBaseMarketplaceDetails

	// REQUIRED; Details of the user.
	User *LiftrBaseUserDetails

	// Organization properties
	PartnerOrganizationProperties *LiftrBaseDataPartnerOrganizationProperties

	// READ-ONLY; Provisioning state of the resource.
	ProvisioningState *ResourceProvisioningState
}

// LiftrBaseDataPartnerOrganizationProperties - Properties specific to Partner's organization
type LiftrBaseDataPartnerOrganizationProperties struct {
	// REQUIRED; Organization name in partner's system
	OrganizationName *string

	// Organization Id in partner's system
	OrganizationID *string

	// Single Sign On properties for the organization
	SingleSignOnProperties *LiftrBaseSingleSignOnProperties

	// Workspace Id in partner's system
	WorkspaceID *string

	// Workspace name in partner's system
	WorkspaceName *string
}

// LiftrBaseDataPartnerOrganizationPropertiesUpdate - Properties specific to Partner's organization
type LiftrBaseDataPartnerOrganizationPropertiesUpdate struct {
	// Organization Id in partner's system
	OrganizationID *string

	// Organization name in partner's system
	OrganizationName *string

	// Single Sign On properties for the organization
	SingleSignOnProperties *LiftrBaseSingleSignOnProperties

	// Workspace Id in partner's system
	WorkspaceID *string

	// Workspace name in partner's system
	WorkspaceName *string
}

// LiftrBaseMarketplaceDetails - Marketplace details for an organization
type LiftrBaseMarketplaceDetails struct {
	// REQUIRED; Offer details for the marketplace that is selected by the user
	OfferDetails *LiftrBaseOfferDetails

	// REQUIRED; Azure subscription id for the the marketplace offer is purchased from
	SubscriptionID *string

	// Marketplace subscription status
	SubscriptionStatus *MarketplaceSubscriptionStatus
}

// LiftrBaseOfferDetails - Offer details for the marketplace that is selected by the user
type LiftrBaseOfferDetails struct {
	// REQUIRED; Offer Id for the marketplace offer
	OfferID *string

	// REQUIRED; Plan Id for the marketplace offer
	PlanID *string

	// REQUIRED; Publisher Id for the marketplace offer
	PublisherID *string

	// Plan Name for the marketplace offer
	PlanName *string

	// Plan Display Name for the marketplace offer
	TermID *string

	// Plan Display Name for the marketplace offer
	TermUnit *string
}

// LiftrBaseSingleSignOnProperties - Properties specific to Single Sign On Resource
type LiftrBaseSingleSignOnProperties struct {
	// List of AAD domains fetched from Microsoft Graph for user.
	AADDomains []*string

	// AAD enterprise application Id used to setup SSO
	EnterpriseAppID *string

	// State of the Single Sign On for the organization
	SingleSignOnState *SingleSignOnStates

	// URL for SSO to be used by the partner to redirect the user to their system
	SingleSignOnURL *string

	// READ-ONLY; Provisioning State of the resource
	ProvisioningState *ResourceProvisioningState
}

// LiftrBaseUserDetails - User details for an organization
type LiftrBaseUserDetails struct {
	// REQUIRED; Email address of the user
	EmailAddress *string

	// REQUIRED; First name of the user
	FirstName *string

	// REQUIRED; Last name of the user
	LastName *string

	// User's phone number
	PhoneNumber *string

	// User's principal name
	Upn *string
}

// LiftrBaseUserDetailsUpdate - User details for an organization
type LiftrBaseUserDetailsUpdate struct {
	// Email address of the user
	EmailAddress *string

	// First name of the user
	FirstName *string

	// Last name of the user
	LastName *string

	// User's phone number
	PhoneNumber *string

	// User's principal name
	Upn *string
}

// ManagedServiceIdentity - Managed service identity (system assigned and/or user assigned identities)
type ManagedServiceIdentity struct {
	// REQUIRED; Type of managed service identity (where both SystemAssigned and UserAssigned types are allowed).
	Type *ManagedServiceIdentityType

	// The set of user assigned identities associated with the resource. The userAssignedIdentities dictionary keys will be ARM
	// resource ids in the form:
	// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}.
	// The dictionary values can be empty objects ({}) in
	// requests.
	UserAssignedIdentities map[string]*UserAssignedIdentity

	// READ-ONLY; The service principal ID of the system assigned identity. This property will only be provided for a system assigned
	// identity.
	PrincipalID *string

	// READ-ONLY; The tenant ID of the system assigned identity. This property will only be provided for a system assigned identity.
	TenantID *string
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

// OrganizationResource - Organization Resource by Astronomer
type OrganizationResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// The managed service identities assigned to this resource.
	Identity *ManagedServiceIdentity

	// The resource-specific properties for this resource.
	Properties *LiftrBaseDataOrganizationProperties

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

// OrganizationResourceListResult - The response of a OrganizationResource list operation.
type OrganizationResourceListResult struct {
	// REQUIRED; The OrganizationResource items on this page
	Value []*OrganizationResource

	// The link to the next page of items
	NextLink *string
}

// OrganizationResourceUpdate - The type used for update operations of the OrganizationResource.
type OrganizationResourceUpdate struct {
	// The managed service identities assigned to this resource.
	Identity *ManagedServiceIdentity

	// The updatable properties of the OrganizationResource.
	Properties *OrganizationResourceUpdateProperties

	// Resource tags.
	Tags map[string]*string
}

// OrganizationResourceUpdateProperties - The updatable properties of the OrganizationResource.
type OrganizationResourceUpdateProperties struct {
	// Organization properties
	PartnerOrganizationProperties *LiftrBaseDataPartnerOrganizationPropertiesUpdate

	// Details of the user.
	User *LiftrBaseUserDetailsUpdate
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

// UserAssignedIdentity - User assigned identity properties
type UserAssignedIdentity struct {
	// READ-ONLY; The client ID of the assigned identity.
	ClientID *string

	// READ-ONLY; The principal ID of the assigned identity.
	PrincipalID *string
}
