//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridconnectivity

import "time"

// AADProfileProperties - The AAD Profile
type AADProfileProperties struct {
	// REQUIRED; The arc ingress gateway server app id.
	ServerID *string

	// REQUIRED; The target resource home tenant id.
	TenantID *string
}

// EndpointAccessResource - The endpoint access for the target resource.
type EndpointAccessResource struct {
	// Azure relay hybrid connection access properties
	Relay *RelayNamespaceAccessProperties
}

// EndpointProperties - Endpoint details
type EndpointProperties struct {
	// REQUIRED; The type of endpoint.
	Type *Type

	// The resource Id of the connectivity endpoint (optional).
	ResourceID *string

	// READ-ONLY; The resource provisioning state.
	ProvisioningState *string
}

// EndpointResource - The endpoint for the target resource.
type EndpointResource struct {
	// The endpoint properties.
	Properties *EndpointProperties

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// EndpointsList - The list of endpoints.
type EndpointsList struct {
	// The link used to get the next page of endpoints list.
	NextLink *string

	// The list of endpoint.
	Value []*EndpointResource
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

// IngressGatewayResource - The ingress gateway access credentials
type IngressGatewayResource struct {
	// Ingress gateway profile
	Ingress *IngressProfileProperties

	// Azure relay hybrid connection access properties
	Relay *RelayNamespaceAccessProperties
}

// IngressProfileProperties - Ingress gateway profile
type IngressProfileProperties struct {
	// REQUIRED; The AAD Profile
	AADProfile *AADProfileProperties

	// REQUIRED; The ingress hostname.
	Hostname *string
}

// ListCredentialsRequest - The details of the service for which credentials needs to be returned.
type ListCredentialsRequest struct {
	// The name of the service. If not provided, the request will by pass the generation of service configuration token
	ServiceName *ServiceName
}

// ListIngressGatewayCredentialsRequest - Represent ListIngressGatewayCredentials Request object.
type ListIngressGatewayCredentialsRequest struct {
	// The name of the service.
	ServiceName *ServiceName
}

// ManagedProxyRequest - Represent ManageProxy Request object.
type ManagedProxyRequest struct {
	// REQUIRED; The name of the service.
	Service *string

	// The target host name.
	Hostname *string

	// The name of the service. It is an optional property, if not provided, service configuration tokens issue code would be
	// by passed.
	ServiceName *ServiceName
}

// ManagedProxyResource - Managed Proxy
type ManagedProxyResource struct {
	// REQUIRED; The expiration time of short lived proxy name in unix epoch.
	ExpiresOn *int64

	// REQUIRED; The short lived proxy name.
	Proxy *string
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
	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// RelayNamespaceAccessProperties - Azure relay hybrid connection access properties
type RelayNamespaceAccessProperties struct {
	// REQUIRED; Azure Relay hybrid connection name for the resource.
	HybridConnectionName *string

	// REQUIRED; The namespace name.
	NamespaceName *string

	// REQUIRED; The suffix domain name of relay namespace.
	NamespaceNameSuffix *string

	// The expiration of access key in unix time.
	ExpiresOn *int64

	// The token to access the enabled service.
	ServiceConfigurationToken *string

	// READ-ONLY; Access key for hybrid connection.
	AccessKey *string
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ServiceConfigurationList - The paginated list of serviceConfigurations
type ServiceConfigurationList struct {
	// The link to fetch the next page of connected cluster
	NextLink *string

	// The list of service configuration
	Value []*ServiceConfigurationResource
}

// ServiceConfigurationProperties - Service configuration details
type ServiceConfigurationProperties struct {
	// REQUIRED; Name of the service.
	ServiceName *ServiceName

	// The port on which service is enabled.
	Port *int64

	// The resource Id of the connectivity endpoint (optional).
	ResourceID *string

	// READ-ONLY; The resource provisioning state.
	ProvisioningState *ProvisioningState
}

// ServiceConfigurationPropertiesPatch - Service configuration details
type ServiceConfigurationPropertiesPatch struct {
	// The port on which service is enabled.
	Port *int64
}

// ServiceConfigurationResource - The service configuration details associated with the target resource.
type ServiceConfigurationResource struct {
	// The service configuration properties.
	Properties *ServiceConfigurationProperties

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ServiceConfigurationResourcePatch - The service details under service configuration for the target endpoint resource.
type ServiceConfigurationResourcePatch struct {
	// The service configuration properties.
	Properties *ServiceConfigurationPropertiesPatch
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
