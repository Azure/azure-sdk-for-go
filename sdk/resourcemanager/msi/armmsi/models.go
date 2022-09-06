//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armmsi

// AssociatedResourcesListResult - Azure resources returned by the resource action to get a list of assigned resources.
type AssociatedResourcesListResult struct {
	// READ-ONLY; The url to get the next page of results, if any.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; Total number of Azure resources assigned to the identity.
	TotalCount *float32 `json:"totalCount,omitempty" azure:"ro"`

	// READ-ONLY; The collection of Azure resources returned by the resource action to get a list of assigned resources.
	Value []*AzureResource `json:"value,omitempty" azure:"ro"`
}

// AzureResource - Describes an Azure resource that is attached to an identity.
type AzureResource struct {
	// READ-ONLY; The ID of this resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of this resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource group this resource belongs to.
	ResourceGroup *string `json:"resourceGroup,omitempty" azure:"ro"`

	// READ-ONLY; The name of the subscription this resource belongs to.
	SubscriptionDisplayName *string `json:"subscriptionDisplayName,omitempty" azure:"ro"`

	// READ-ONLY; The ID of the subscription this resource belongs to.
	SubscriptionID *string `json:"subscriptionId,omitempty" azure:"ro"`

	// READ-ONLY; The type of this resource.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// CloudError - An error response from the ManagedServiceIdentity service.
type CloudError struct {
	// A list of additional details about the error.
	Error *CloudErrorBody `json:"error,omitempty"`
}

// CloudErrorBody - An error response from the ManagedServiceIdentity service.
type CloudErrorBody struct {
	// An identifier for the error.
	Code *string `json:"code,omitempty"`

	// A list of additional details about the error.
	Details []*CloudErrorBody `json:"details,omitempty"`

	// A message describing the error, intended to be suitable for display in a user interface.
	Message *string `json:"message,omitempty"`

	// The target of the particular error. For example, the name of the property in error.
	Target *string `json:"target,omitempty"`
}

// FederatedIdentityCredential - Describes a federated identity credential.
type FederatedIdentityCredential struct {
	// The properties associated with the federated identity credential.
	Properties *FederatedIdentityCredentialProperties `json:"properties,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// FederatedIdentityCredentialProperties - The properties associated with a federated identity credential.
type FederatedIdentityCredentialProperties struct {
	// REQUIRED; The list of audiences that can appear in the issued token.
	Audiences []*string `json:"audiences,omitempty"`

	// REQUIRED; The URL of the issuer to be trusted.
	Issuer *string `json:"issuer,omitempty"`

	// REQUIRED; The identifier of the external identity.
	Subject *string `json:"subject,omitempty"`
}

// FederatedIdentityCredentialsClientCreateOrUpdateOptions contains the optional parameters for the FederatedIdentityCredentialsClient.CreateOrUpdate
// method.
type FederatedIdentityCredentialsClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// FederatedIdentityCredentialsClientDeleteOptions contains the optional parameters for the FederatedIdentityCredentialsClient.Delete
// method.
type FederatedIdentityCredentialsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// FederatedIdentityCredentialsClientGetOptions contains the optional parameters for the FederatedIdentityCredentialsClient.Get
// method.
type FederatedIdentityCredentialsClientGetOptions struct {
	// placeholder for future optional parameters
}

// FederatedIdentityCredentialsClientListOptions contains the optional parameters for the FederatedIdentityCredentialsClient.List
// method.
type FederatedIdentityCredentialsClientListOptions struct {
	// A skip token is used to continue retrieving items after an operation returns a partial result. If a previous response contains
	// a nextLink element, the value of the nextLink element will include a
	// skipToken parameter that specifies a starting point to use for subsequent calls.
	Skiptoken *string
	// Number of records to return.
	Top *int32
}

// FederatedIdentityCredentialsListResult - Values returned by the List operation for federated identity credentials.
type FederatedIdentityCredentialsListResult struct {
	// The url to get the next page of results, if any.
	NextLink *string `json:"nextLink,omitempty"`

	// The collection of federated identity credentials returned by the listing operation.
	Value []*FederatedIdentityCredential `json:"value,omitempty"`
}

// Identity - Describes an identity resource.
type Identity struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The properties associated with the identity.
	Properties *UserAssignedIdentityProperties `json:"properties,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// IdentityUpdate - Describes an identity resource.
type IdentityUpdate struct {
	// The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The properties associated with the identity.
	Properties *UserAssignedIdentityProperties `json:"properties,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// Operation supported by the Microsoft.ManagedIdentity REST API.
type Operation struct {
	// The object that describes the operation.
	Display *OperationDisplay `json:"display,omitempty"`

	// The name of the REST Operation. This is of the format {provider}/{resource}/{operation}.
	Name *string `json:"name,omitempty"`
}

// OperationDisplay - The object that describes the operation.
type OperationDisplay struct {
	// A description of the operation.
	Description *string `json:"description,omitempty"`

	// The type of operation. For example: read, write, delete.
	Operation *string `json:"operation,omitempty"`

	// Friendly name of the resource provider.
	Provider *string `json:"provider,omitempty"`

	// The resource type on which the operation is performed.
	Resource *string `json:"resource,omitempty"`
}

// OperationListResult - A list of operations supported by Microsoft.ManagedIdentity Resource Provider.
type OperationListResult struct {
	// The url to get the next page of results, if any.
	NextLink *string `json:"nextLink,omitempty"`

	// A list of operations supported by Microsoft.ManagedIdentity Resource Provider.
	Value []*Operation `json:"value,omitempty"`
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.List method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// ProxyResource - The resource model definition for a Azure Resource Manager proxy resource. It will not have tags and a
// location
type ProxyResource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// SystemAssignedIdentitiesClientGetByScopeOptions contains the optional parameters for the SystemAssignedIdentitiesClient.GetByScope
// method.
type SystemAssignedIdentitiesClientGetByScopeOptions struct {
	// placeholder for future optional parameters
}

// SystemAssignedIdentity - Describes a system assigned identity resource.
type SystemAssignedIdentity struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The properties associated with the identity.
	Properties *SystemAssignedIdentityProperties `json:"properties,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// SystemAssignedIdentityProperties - The properties associated with the system assigned identity.
type SystemAssignedIdentityProperties struct {
	// READ-ONLY; The id of the app associated with the identity. This is a random generated UUID by MSI.
	ClientID *string `json:"clientId,omitempty" azure:"ro"`

	// READ-ONLY; The ManagedServiceIdentity DataPlane URL that can be queried to obtain the identity credentials.
	ClientSecretURL *string `json:"clientSecretUrl,omitempty" azure:"ro"`

	// READ-ONLY; The id of the service principal object associated with the created identity.
	PrincipalID *string `json:"principalId,omitempty" azure:"ro"`

	// READ-ONLY; The id of the tenant which the identity belongs to.
	TenantID *string `json:"tenantId,omitempty" azure:"ro"`
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// UserAssignedIdentitiesClientCreateOrUpdateOptions contains the optional parameters for the UserAssignedIdentitiesClient.CreateOrUpdate
// method.
type UserAssignedIdentitiesClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// UserAssignedIdentitiesClientDeleteOptions contains the optional parameters for the UserAssignedIdentitiesClient.Delete
// method.
type UserAssignedIdentitiesClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// UserAssignedIdentitiesClientGetOptions contains the optional parameters for the UserAssignedIdentitiesClient.Get method.
type UserAssignedIdentitiesClientGetOptions struct {
	// placeholder for future optional parameters
}

// UserAssignedIdentitiesClientListAssociatedResourcesOptions contains the optional parameters for the UserAssignedIdentitiesClient.ListAssociatedResources
// method.
type UserAssignedIdentitiesClientListAssociatedResourcesOptions struct {
	// OData filter expression to apply to the query.
	Filter *string
	// OData orderBy expression to apply to the query.
	Orderby *string
	// Number of records to skip.
	Skip *int32
	// A skip token is used to continue retrieving items after an operation returns a partial result. If a previous response contains
	// a nextLink element, the value of the nextLink element will include a
	// skipToken parameter that specifies a starting point to use for subsequent calls.
	Skiptoken *string
	// Number of records to return.
	Top *int32
}

// UserAssignedIdentitiesClientListByResourceGroupOptions contains the optional parameters for the UserAssignedIdentitiesClient.ListByResourceGroup
// method.
type UserAssignedIdentitiesClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// UserAssignedIdentitiesClientListBySubscriptionOptions contains the optional parameters for the UserAssignedIdentitiesClient.ListBySubscription
// method.
type UserAssignedIdentitiesClientListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// UserAssignedIdentitiesClientUpdateOptions contains the optional parameters for the UserAssignedIdentitiesClient.Update
// method.
type UserAssignedIdentitiesClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// UserAssignedIdentitiesListResult - Values returned by the List operation.
type UserAssignedIdentitiesListResult struct {
	// The url to get the next page of results, if any.
	NextLink *string `json:"nextLink,omitempty"`

	// The collection of userAssignedIdentities returned by the listing operation.
	Value []*Identity `json:"value,omitempty"`
}

// UserAssignedIdentityProperties - The properties associated with the user assigned identity.
type UserAssignedIdentityProperties struct {
	// READ-ONLY; The id of the app associated with the identity. This is a random generated UUID by MSI.
	ClientID *string `json:"clientId,omitempty" azure:"ro"`

	// READ-ONLY; The id of the service principal object associated with the created identity.
	PrincipalID *string `json:"principalId,omitempty" azure:"ro"`

	// READ-ONLY; The id of the tenant which the identity belongs to.
	TenantID *string `json:"tenantId,omitempty" azure:"ro"`
}
