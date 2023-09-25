//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpurview

import "time"

// AccessKeys - The Account access keys.
type AccessKeys struct {
	// Gets or sets the primary connection string.
	AtlasKafkaPrimaryEndpoint *string

	// Gets or sets the secondary connection string.
	AtlasKafkaSecondaryEndpoint *string
}

// Account resource
type Account struct {
	// Identity Info on the tracked resource
	Identity *Identity

	// Gets or sets the location.
	Location *string

	// Gets or sets the properties.
	Properties *AccountProperties

	// Tags on the azure resource.
	Tags map[string]*string

	// READ-ONLY; Gets or sets the identifier.
	ID *string

	// READ-ONLY; Gets or sets the name.
	Name *string

	// READ-ONLY; Gets or sets the Sku.
	SKU *AccountSKU

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *TrackedResourceSystemData

	// READ-ONLY; Gets or sets the type.
	Type *string
}

// AccountEndpoints - The account endpoints
type AccountEndpoints struct {
	// READ-ONLY; Gets the catalog endpoint.
	Catalog *string

	// READ-ONLY; Gets the guardian endpoint.
	Guardian *string

	// READ-ONLY; Gets the scan endpoint.
	Scan *string
}

// AccountList - Paged list of account resources
type AccountList struct {
	// REQUIRED; Collection of items of type results.
	Value []*Account

	// Total item count.
	Count *int64

	// The Url of next result page.
	NextLink *string
}

// AccountProperties - The account properties
type AccountProperties struct {
	// Cloud connectors. External cloud identifier used as part of scanning configuration.
	CloudConnectors *CloudConnectors

	// Gets or sets the managed resource group name
	ManagedResourceGroupName *string

	// Gets or sets the public network access.
	PublicNetworkAccess *PublicNetworkAccess

	// READ-ONLY; Gets the time at which the entity was created.
	CreatedAt *time.Time

	// READ-ONLY; Gets the creator of the entity.
	CreatedBy *string

	// READ-ONLY; Gets the creators of the entity's object id.
	CreatedByObjectID *string

	// READ-ONLY; The URIs that are the public endpoints of the account.
	Endpoints *AccountPropertiesEndpoints

	// READ-ONLY; Gets or sets the friendly name.
	FriendlyName *string

	// READ-ONLY; Gets the resource identifiers of the managed resources.
	ManagedResources *AccountPropertiesManagedResources

	// READ-ONLY; Gets the private endpoint connections information.
	PrivateEndpointConnections []*PrivateEndpointConnection

	// READ-ONLY; Gets or sets the state of the provisioning.
	ProvisioningState *ProvisioningState
}

// AccountPropertiesEndpoints - The URIs that are the public endpoints of the account.
type AccountPropertiesEndpoints struct {
	// READ-ONLY; Gets the catalog endpoint.
	Catalog *string

	// READ-ONLY; Gets the guardian endpoint.
	Guardian *string

	// READ-ONLY; Gets the scan endpoint.
	Scan *string
}

// AccountPropertiesManagedResources - Gets the resource identifiers of the managed resources.
type AccountPropertiesManagedResources struct {
	// READ-ONLY; Gets the managed event hub namespace resource identifier.
	EventHubNamespace *string

	// READ-ONLY; Gets the managed resource group resource identifier. This resource group will host resource dependencies for
// the account.
	ResourceGroup *string

	// READ-ONLY; Gets the managed storage account resource identifier.
	StorageAccount *string
}

// AccountSKU - Gets or sets the Sku.
type AccountSKU struct {
	// Gets or sets the sku capacity.
	Capacity *int32

	// Gets or sets the sku name.
	Name *Name
}

// AccountSKUAutoGenerated - The Sku
type AccountSKUAutoGenerated struct {
	// Gets or sets the sku capacity.
	Capacity *int32

	// Gets or sets the sku name.
	Name *Name
}

// AccountUpdateParameters - The account update properties.
type AccountUpdateParameters struct {
	// Identity related info to add/remove userAssignedIdentities.
	Identity *Identity

	// The account properties.
	Properties *AccountProperties

	// Tags on the azure resource.
	Tags map[string]*string
}

// CheckNameAvailabilityRequest - The request payload for CheckNameAvailability API
type CheckNameAvailabilityRequest struct {
	// Resource name to verify for availability
	Name *string

	// Fully qualified resource type which includes provider namespace
	Type *string
}

// CheckNameAvailabilityResult - The response payload for CheckNameAvailability API
type CheckNameAvailabilityResult struct {
	// Error message
	Message *string

	// Indicates if name is valid and available.
	NameAvailable *bool

	// The reason the name is not available.
	Reason *Reason
}

// CloudConnectors - External Cloud Service connectors
type CloudConnectors struct {
	// READ-ONLY; AWS external identifier. Configured in AWS to allow use of the role arn used for scanning
	AwsExternalID *string
}

// CollectionAdminUpdate - Collection administrator update.
type CollectionAdminUpdate struct {
	// Gets or sets the object identifier of the admin.
	ObjectID *string
}

// DefaultAccountPayload - Payload to get and set the default account in the given scope
type DefaultAccountPayload struct {
	// The name of the account that is set as the default.
	AccountName *string

	// The resource group name of the account that is set as the default.
	ResourceGroupName *string

	// The scope object ID. For example, sub ID or tenant ID.
	Scope *string

	// The scope tenant in which the default account is set.
	ScopeTenantID *string

	// The scope where the default account is set.
	ScopeType *ScopeType

	// The subscription ID of the account that is set as the default.
	SubscriptionID *string
}

// DimensionProperties - properties for dimension
type DimensionProperties struct {
	// localized display name of the dimension to customer
	DisplayName *string

	// dimension name
	Name *string

	// flag indicating whether this dimension should be included to the customer in Azure Monitor logs (aka Shoebox)
	ToBeExportedForCustomer *bool
}

// ErrorModel - Default error model
type ErrorModel struct {
	// READ-ONLY; Gets or sets the code.
	Code *string

	// READ-ONLY; Gets or sets the details.
	Details []*ErrorModel

	// READ-ONLY; Gets or sets the messages.
	Message *string

	// READ-ONLY; Gets or sets the target.
	Target *string
}

// ErrorResponseModel - Default error response model
type ErrorResponseModel struct {
	// READ-ONLY; Gets or sets the error.
	Error *ErrorResponseModelError
}

// ErrorResponseModelError - Gets or sets the error.
type ErrorResponseModelError struct {
	// READ-ONLY; Gets or sets the code.
	Code *string

	// READ-ONLY; Gets or sets the details.
	Details []*ErrorModel

	// READ-ONLY; Gets or sets the messages.
	Message *string

	// READ-ONLY; Gets or sets the target.
	Target *string
}

// Identity - The Managed Identity of the resource
type Identity struct {
	// Identity Type
	Type *Type

	// User Assigned Identities
	UserAssignedIdentities map[string]*UserAssignedIdentity

	// READ-ONLY; Service principal object Id
	PrincipalID *string

	// READ-ONLY; Tenant Id
	TenantID *string
}

// ManagedResources - The managed resources in customer subscription.
type ManagedResources struct {
	// READ-ONLY; Gets the managed event hub namespace resource identifier.
	EventHubNamespace *string

	// READ-ONLY; Gets the managed resource group resource identifier. This resource group will host resource dependencies for
// the account.
	ResourceGroup *string

	// READ-ONLY; Gets the managed storage account resource identifier.
	StorageAccount *string
}

// Operation resource
type Operation struct {
	// Properties on the operation
	Display *OperationDisplay

	// Whether operation is a data action
	IsDataAction *bool

	// Operation name for display purposes
	Name *string

	// origin of the operation
	Origin *string

	// properties for the operation meta info
	Properties *OperationProperties
}

// OperationDisplay - The response model for get operation properties
type OperationDisplay struct {
	// Description of the operation for display purposes
	Description *string

	// Name of the operation for display purposes
	Operation *string

	// Name of the provider for display purposes
	Provider *string

	// Name of the resource type for display purposes
	Resource *string
}

// OperationList - Paged list of operation resources
type OperationList struct {
	// REQUIRED; Collection of items of type results.
	Value []*Operation

	// Total item count.
	Count *int64

	// The Url of next result page.
	NextLink *string
}

// OperationMetaLogSpecification - log specifications for operation api
type OperationMetaLogSpecification struct {
	// blob duration of the log
	BlobDuration *string

	// localized name of the log category
	DisplayName *string

	// name of the log category
	Name *string
}

// OperationMetaMetricSpecification - metric specifications for the operation
type OperationMetaMetricSpecification struct {
	// aggregation type of metric
	AggregationType *string

	// properties for dimension
	Dimensions []*DimensionProperties

	// description of the metric
	DisplayDescription *string

	// localized name of the metric
	DisplayName *string

	// enable regional mdm account
	EnableRegionalMdmAccount *string

	// internal metric name
	InternalMetricName *string

	// name of the metric
	Name *string

	// dimension name use to replace resource id if specified
	ResourceIDDimensionNameOverride *string

	// Metric namespace. Only set the namespace if different from the default value, leaving it empty makes it use the value from
// the ARM manifest.
	SourceMdmNamespace *string

	// supported aggregation types
	SupportedAggregationTypes []*string

	// supported time grain types
	SupportedTimeGrainTypes []*string

	// units for the metric
	Unit *string
}

// OperationMetaServiceSpecification - The operation meta service specification
type OperationMetaServiceSpecification struct {
	// log specifications for the operation
	LogSpecifications []*OperationMetaLogSpecification

	// metric specifications for the operation
	MetricSpecifications []*OperationMetaMetricSpecification
}

// OperationProperties - properties on meta info
type OperationProperties struct {
	// meta service specification
	ServiceSpecification *OperationMetaServiceSpecification
}

// PrivateEndpoint - A private endpoint class.
type PrivateEndpoint struct {
	// The private endpoint identifier.
	ID *string
}

// PrivateEndpointConnection - A private endpoint connection class.
type PrivateEndpointConnection struct {
	// The connection identifier.
	Properties *PrivateEndpointConnectionProperties

	// READ-ONLY; Gets or sets the identifier.
	ID *string

	// READ-ONLY; Gets or sets the name.
	Name *string

	// READ-ONLY; Gets or sets the type.
	Type *string
}

// PrivateEndpointConnectionList - Paged list of private endpoint connections
type PrivateEndpointConnectionList struct {
	// REQUIRED; Collection of items of type results.
	Value []*PrivateEndpointConnection

	// Total item count.
	Count *int64

	// The Url of next result page.
	NextLink *string
}

// PrivateEndpointConnectionProperties - A private endpoint connection properties class.
type PrivateEndpointConnectionProperties struct {
	// The private endpoint information.
	PrivateEndpoint *PrivateEndpoint

	// The private link service connection state.
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState

	// READ-ONLY; The provisioning state.
	ProvisioningState *string
}

// PrivateLinkResource - A privately linkable resource.
type PrivateLinkResource struct {
	// READ-ONLY; The private link resource identifier.
	ID *string

	// READ-ONLY; The private link resource name.
	Name *string

	// READ-ONLY; The private link resource properties.
	Properties *PrivateLinkResourceProperties

	// READ-ONLY; The private link resource type.
	Type *string
}

// PrivateLinkResourceList - Paged list of private link resources
type PrivateLinkResourceList struct {
	// REQUIRED; Collection of items of type results.
	Value []*PrivateLinkResource

	// Total item count.
	Count *int64

	// The Url of next result page.
	NextLink *string
}

// PrivateLinkResourceProperties - A privately linkable resource properties.
type PrivateLinkResourceProperties struct {
	// READ-ONLY; The private link resource group identifier.
	GroupID *string

	// READ-ONLY; This translates to how many Private IPs should be created for each privately linkable resource.
	RequiredMembers []*string

	// READ-ONLY; The required zone names for private link resource.
	RequiredZoneNames []*string
}

// PrivateLinkServiceConnectionState - The private link service connection state.
type PrivateLinkServiceConnectionState struct {
	// The required actions.
	ActionsRequired *string

	// The description.
	Description *string

	// The status.
	Status *Status
}

// ProxyResource - Proxy Azure Resource
type ProxyResource struct {
	// READ-ONLY; Gets or sets the identifier.
	ID *string

	// READ-ONLY; Gets or sets the name.
	Name *string

	// READ-ONLY; Gets or sets the type.
	Type *string
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// READ-ONLY; The timestamp of resource creation (UTC).
	CreatedAt *time.Time

	// READ-ONLY; The identity that created the resource.
	CreatedBy *string

	// READ-ONLY; The type of identity that created the resource.
	CreatedByType *CreatedByType

	// READ-ONLY; The timestamp of the last modification the resource (UTC).
	LastModifiedAt *time.Time

	// READ-ONLY; The identity that last modified the resource.
	LastModifiedBy *string

	// READ-ONLY; The type of identity that last modified the resource.
	LastModifiedByType *LastModifiedByType
}

// TrackedResource - Azure ARM Tracked Resource
type TrackedResource struct {
	// Identity Info on the tracked resource
	Identity *Identity

	// Gets or sets the location.
	Location *string

	// Tags on the azure resource.
	Tags map[string]*string

	// READ-ONLY; Gets or sets the identifier.
	ID *string

	// READ-ONLY; Gets or sets the name.
	Name *string

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *TrackedResourceSystemData

	// READ-ONLY; Gets or sets the type.
	Type *string
}

// TrackedResourceSystemData - Metadata pertaining to creation and last modification of the resource.
type TrackedResourceSystemData struct {
	// READ-ONLY; The timestamp of resource creation (UTC).
	CreatedAt *time.Time

	// READ-ONLY; The identity that created the resource.
	CreatedBy *string

	// READ-ONLY; The type of identity that created the resource.
	CreatedByType *CreatedByType

	// READ-ONLY; The timestamp of the last modification the resource (UTC).
	LastModifiedAt *time.Time

	// READ-ONLY; The identity that last modified the resource.
	LastModifiedBy *string

	// READ-ONLY; The type of identity that last modified the resource.
	LastModifiedByType *LastModifiedByType
}

// UserAssignedIdentity - Uses client ID and Principal ID
type UserAssignedIdentity struct {
	// READ-ONLY; Gets or Sets Client ID
	ClientID *string

	// READ-ONLY; Gets or Sets Principal ID
	PrincipalID *string
}

