// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armstoragediscovery

import "time"

// Operation - REST API Operation
//
// Details of a REST API operation, returned from the Resource Provider Operations API
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

// Scope - Storage Discovery Scope. This had added validations
type Scope struct {
	// REQUIRED; Display name of the collection
	DisplayName *string

	// REQUIRED; Resource types for the collection
	ResourceTypes []*ResourceType

	// The storage account tags keys to filter
	TagKeysOnly []*string

	// Resource tags.
	Tags map[string]*string
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

// Workspace - A Storage Discovery Workspace resource. This resource configures the collection of storage account metrics.
type Workspace struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// The resource-specific properties for this resource.
	Properties *WorkspaceProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; The name of the StorageDiscoveryWorkspace
	Name *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// WorkspaceListResult - The response of a StorageDiscoveryWorkspace list operation.
type WorkspaceListResult struct {
	// REQUIRED; The StorageDiscoveryWorkspace items on this page
	Value []*Workspace

	// The link to the next page of items
	NextLink *string
}

// WorkspaceProperties - Storage Discovery Workspace Properties
type WorkspaceProperties struct {
	// REQUIRED; The scopes of the storage discovery workspace.
	Scopes []*Scope

	// REQUIRED; The view level storage discovery data estate
	WorkspaceRoots []*string

	// The description of the storage discovery workspace
	Description *string

	// The storage discovery sku
	SKU *SKU

	// READ-ONLY; The status of the last operation.
	ProvisioningState *ResourceProvisioningState
}

// WorkspacePropertiesUpdate - The template for adding updateable properties.
type WorkspacePropertiesUpdate struct {
	// The description of the storage discovery workspace
	Description *string

	// The storage discovery sku
	SKU *SKU

	// The scopes of the storage discovery workspace.
	Scopes []*Scope

	// The view level storage discovery data estate
	WorkspaceRoots []*string
}

// WorkspaceUpdate - The template for adding updateable properties.
type WorkspaceUpdate struct {
	// The resource-specific properties for this resource.
	Properties *WorkspacePropertiesUpdate

	// Resource tags.
	Tags map[string]*string
}
