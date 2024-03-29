//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armautomanage

import "time"

// AssignmentReportProperties - Data related to the report detail.
type AssignmentReportProperties struct {
	// End time of the configuration profile assignment processing.
	EndTime *string

	// Start time of the configuration profile assignment processing.
	StartTime *string

	// READ-ONLY; The configurationProfile linked to the assignment.
	ConfigurationProfile *string

	// READ-ONLY; Duration of the configuration profile assignment processing.
	Duration *string

	// READ-ONLY; Error message, if any, returned by the configuration profile assignment processing.
	Error *ErrorDetail

	// READ-ONLY; Last modified time of the configuration profile assignment processing.
	LastModifiedTime *string

	// READ-ONLY; Version of the report format
	ReportFormatVersion *string

	// READ-ONLY; List of resources processed by the configuration profile assignment.
	Resources []*ReportResource

	// READ-ONLY; The status of the configuration profile assignment.
	Status *string

	// READ-ONLY; Type of the configuration profile assignment processing (Initial/Consistency).
	Type *string
}

// BestPractice - Definition of the Automanage best practice.
type BestPractice struct {
	// Properties of the best practice.
	Properties *ConfigurationProfileProperties

	// READ-ONLY; The fully qualified ID for the best practice. For example, /providers/Microsoft.Automanage/bestPractices/azureBestPracticesProduction
	ID *string

	// READ-ONLY; The name of the best practice. For example, azureBestPracticesProduction
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. For example, Microsoft.Automanage/bestPractices
	Type *string
}

// BestPracticeList - The response of the list best practice operation.
type BestPracticeList struct {
	// Result of the list best practice operation.
	Value []*BestPractice
}

// ConfigurationProfile - Definition of the configuration profile.
type ConfigurationProfile struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// Properties of the configuration profile.
	Properties *ConfigurationProfileProperties

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

// ConfigurationProfileAssignment - Configuration profile assignment is an association between a VM and automanage profile
// configuration.
type ConfigurationProfileAssignment struct {
	// Properties of the configuration profile assignment.
	Properties *ConfigurationProfileAssignmentProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedBy *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ConfigurationProfileAssignmentList - The response of the list configuration profile assignment operation.
type ConfigurationProfileAssignmentList struct {
	// Result of the list configuration profile assignment operation.
	Value []*ConfigurationProfileAssignment
}

// ConfigurationProfileAssignmentProperties - Automanage configuration profile assignment properties.
type ConfigurationProfileAssignmentProperties struct {
	// The Automanage configurationProfile ARM Resource URI.
	ConfigurationProfile *string

	// READ-ONLY; The status of onboarding, which only appears in the response.
	Status *string

	// READ-ONLY; The target VM resource URI
	TargetID *string
}

// ConfigurationProfileList - The response of the list configuration profile operation.
type ConfigurationProfileList struct {
	// Result of the list ConfigurationProfile operation.
	Value []*ConfigurationProfile
}

// ConfigurationProfileProperties - Automanage configuration profile properties.
type ConfigurationProfileProperties struct {
	// configuration dictionary of the configuration profile.
	Configuration any
}

// ConfigurationProfileUpdate - Definition of the configuration profile.
type ConfigurationProfileUpdate struct {
	// Properties of the configuration profile.
	Properties *ConfigurationProfileProperties

	// The tags of the resource.
	Tags map[string]*string
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

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// Report - Definition of the report.
type Report struct {
	// The properties for the report.
	Properties *AssignmentReportProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ReportList - The response of the list report operation.
type ReportList struct {
	// Result of the list report operation.
	Value []*Report
}

// ReportResource - Details about the resource processed by the configuration profile assignment
type ReportResource struct {
	// READ-ONLY; Error message, if any, returned when deploying the resource.
	Error *ErrorDetail

	// READ-ONLY; ARM id of the resource.
	ID *string

	// READ-ONLY; Name of the resource.
	Name *string

	// READ-ONLY; Status of the resource.
	Status *string

	// READ-ONLY; Type of the resource.
	Type *string
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ServicePrincipal - The Service Principal Id for the subscription.
type ServicePrincipal struct {
	// The Service Principal properties for the subscription
	Properties *ServicePrincipalProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ServicePrincipalListResult - The list of ServicePrincipals.
type ServicePrincipalListResult struct {
	// The list of servicePrincipals.
	Value []*ServicePrincipal
}

// ServicePrincipalProperties - The Service Principal properties for the subscription.
type ServicePrincipalProperties struct {
	// READ-ONLY; Returns the contributor RBAC Role exist or not for the Service Principal Id.
	AuthorizationSet *bool

	// READ-ONLY; The Service Principal Id for the subscription.
	ServicePrincipalID *string
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

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// UpdateResource - Represents an update resource
type UpdateResource struct {
	// The tags of the resource.
	Tags map[string]*string
}
