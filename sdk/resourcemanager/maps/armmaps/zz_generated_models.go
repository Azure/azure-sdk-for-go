//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmaps

import "time"

// Account - An Azure resource which represents access to a suite of Maps REST APIs.
type Account struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// REQUIRED; The SKU of this account.
	SKU *SKU `json:"sku,omitempty"`

	// Sets the identity property for maps account.
	Identity *ManagedServiceIdentity `json:"identity,omitempty"`

	// Get or Set Kind property.
	Kind *Kind `json:"kind,omitempty"`

	// The map account properties.
	Properties *AccountProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The system meta data relating to this resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// AccountKeys - The set of keys which can be used to access the Maps REST APIs. Two keys are provided for key rotation without
// interruption.
type AccountKeys struct {
	// READ-ONLY; The primary key for accessing the Maps REST APIs.
	PrimaryKey *string `json:"primaryKey,omitempty" azure:"ro"`

	// READ-ONLY; The last updated date and time of the primary key.
	PrimaryKeyLastUpdated *string `json:"primaryKeyLastUpdated,omitempty" azure:"ro"`

	// READ-ONLY; The secondary key for accessing the Maps REST APIs.
	SecondaryKey *string `json:"secondaryKey,omitempty" azure:"ro"`

	// READ-ONLY; The last updated date and time of the secondary key.
	SecondaryKeyLastUpdated *string `json:"secondaryKeyLastUpdated,omitempty" azure:"ro"`
}

// AccountProperties - Additional Map account properties
type AccountProperties struct {
	// Specifies CORS rules for the Blob service. You can include up to five CorsRule elements in the request. If no CorsRule
	// elements are included in the request body, all CORS rules will be deleted, and
	// CORS will be disabled for the Blob service.
	Cors *CorsRules `json:"cors,omitempty"`

	// Allows toggle functionality on Azure Policy to disable Azure Maps local authentication support. This will disable Shared
	// Keys authentication from any usage.
	DisableLocalAuth *bool `json:"disableLocalAuth,omitempty"`

	// Sets the resources to be used for Managed Identities based operations for the Map account resource.
	LinkedResources []*LinkedResource `json:"linkedResources,omitempty"`

	// READ-ONLY; The provisioning state of the Map account resource.
	ProvisioningState *string `json:"provisioningState,omitempty" azure:"ro"`

	// READ-ONLY; A unique identifier for the maps account
	UniqueID *string `json:"uniqueId,omitempty" azure:"ro"`
}

// AccountSasParameters - Parameters used to create an account Shared Access Signature (SAS) token. The REST API access control
// is provided by Azure Maps Role Based Access (RBAC) identity and access.
type AccountSasParameters struct {
	// REQUIRED; The date time offset of when the token validity expires. For example "2017-05-24T10:42:03.1567373Z"
	Expiry *string `json:"expiry,omitempty"`

	// REQUIRED; Required parameter which represents the desired maximum request per second to allowed for the given SAS token.
	// This does not guarantee perfect accuracy in measurements but provides application safe
	// guards of abuse with eventual enforcement.
	MaxRatePerSecond *int32 `json:"maxRatePerSecond,omitempty"`

	// REQUIRED; The principal Id also known as the object Id of a User Assigned Managed Identity currently assigned to the Map
	// Account. To assign a Managed Identity of the account, use operation Create or Update an
	// assign a User Assigned Identity resource Id.
	PrincipalID *string `json:"principalId,omitempty"`

	// REQUIRED; The Map account key to use for signing.
	SigningKey *SigningKey `json:"signingKey,omitempty"`

	// REQUIRED; The date time offset of when the token validity begins. For example "2017-05-24T10:42:03.1567373Z".
	Start *string `json:"start,omitempty"`

	// Optional, allows control of which region locations are permitted access to Azure Maps REST APIs with the SAS token. Example:
	// "eastus", "westus2". Omitting this parameter will allow all region
	// locations to be accessible.
	Regions []*string `json:"regions,omitempty"`
}

// AccountSasToken - A new Sas token which can be used to access the Maps REST APIs and is controlled by the specified Managed
// identity permissions on Azure (IAM) Role Based Access Control.
type AccountSasToken struct {
	// READ-ONLY; The shared access signature access token.
	AccountSasToken *string `json:"accountSasToken,omitempty" azure:"ro"`
}

// AccountUpdateParameters - Parameters used to update an existing Maps Account.
type AccountUpdateParameters struct {
	// Sets the identity property for maps account.
	Identity *ManagedServiceIdentity `json:"identity,omitempty"`

	// Get or Set Kind property.
	Kind *Kind `json:"kind,omitempty"`

	// The map account properties.
	Properties *AccountProperties `json:"properties,omitempty"`

	// The SKU of this account.
	SKU *SKU `json:"sku,omitempty"`

	// Gets or sets a list of key value pairs that describe the resource. These tags can be used in viewing and grouping this
	// resource (across resource groups). A maximum of 15 tags can be provided for a
	// resource. Each tag must have a key no greater than 128 characters and value no greater than 256 characters.
	Tags map[string]*string `json:"tags,omitempty"`
}

// Accounts - A list of Maps Accounts.
type Accounts struct {
	// URL client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string `json:"nextLink,omitempty"`

	// READ-ONLY; a Maps Account.
	Value []*Account `json:"value,omitempty" azure:"ro"`
}

// AccountsClientCreateOrUpdateOptions contains the optional parameters for the AccountsClient.CreateOrUpdate method.
type AccountsClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientDeleteOptions contains the optional parameters for the AccountsClient.Delete method.
type AccountsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientGetOptions contains the optional parameters for the AccountsClient.Get method.
type AccountsClientGetOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListByResourceGroupOptions contains the optional parameters for the AccountsClient.ListByResourceGroup method.
type AccountsClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListBySubscriptionOptions contains the optional parameters for the AccountsClient.ListBySubscription method.
type AccountsClientListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListKeysOptions contains the optional parameters for the AccountsClient.ListKeys method.
type AccountsClientListKeysOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListSasOptions contains the optional parameters for the AccountsClient.ListSas method.
type AccountsClientListSasOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientRegenerateKeysOptions contains the optional parameters for the AccountsClient.RegenerateKeys method.
type AccountsClientRegenerateKeysOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientUpdateOptions contains the optional parameters for the AccountsClient.Update method.
type AccountsClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// ClientListOperationsOptions contains the optional parameters for the Client.ListOperations method.
type ClientListOperationsOptions struct {
	// placeholder for future optional parameters
}

// ClientListSubscriptionOperationsOptions contains the optional parameters for the Client.ListSubscriptionOperations method.
type ClientListSubscriptionOperationsOptions struct {
	// placeholder for future optional parameters
}

type Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties struct {
	// READ-ONLY; The client id of user assigned identity.
	ClientID *string `json:"clientId,omitempty" azure:"ro"`

	// READ-ONLY; The principal id of user assigned identity.
	PrincipalID *string `json:"principalId,omitempty" azure:"ro"`
}

// CorsRule - Specifies a CORS rule for the Map Account.
type CorsRule struct {
	// REQUIRED; Required if CorsRule element is present. A list of origin domains that will be allowed via CORS, or "*" to allow
	// all domains
	AllowedOrigins []*string `json:"allowedOrigins,omitempty"`
}

// CorsRules - Sets the CORS rules. You can include up to five CorsRule elements in the request.
type CorsRules struct {
	// The list of CORS rules. You can include up to five CorsRule elements in the request.
	CorsRules []*CorsRule `json:"corsRules,omitempty"`
}

// Creator - An Azure resource which represents Maps Creator product and provides ability to manage private location data.
type Creator struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// REQUIRED; The Creator resource properties.
	Properties *CreatorProperties `json:"properties,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The system meta data relating to this resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// CreatorList - A list of Creator resources.
type CreatorList struct {
	// URL client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string `json:"nextLink,omitempty"`

	// READ-ONLY; a Creator account.
	Value []*Creator `json:"value,omitempty" azure:"ro"`
}

// CreatorProperties - Creator resource properties
type CreatorProperties struct {
	// REQUIRED; The storage units to be allocated. Integer values from 1 to 100, inclusive.
	StorageUnits *int32 `json:"storageUnits,omitempty"`

	// READ-ONLY; The state of the resource provisioning, terminal states: Succeeded, Failed, Canceled
	ProvisioningState *string `json:"provisioningState,omitempty" azure:"ro"`
}

// CreatorUpdateParameters - Parameters used to update an existing Creator resource.
type CreatorUpdateParameters struct {
	// Creator resource properties.
	Properties *CreatorProperties `json:"properties,omitempty"`

	// Gets or sets a list of key value pairs that describe the resource. These tags can be used in viewing and grouping this
	// resource (across resource groups). A maximum of 15 tags can be provided for a
	// resource. Each tag must have a key no greater than 128 characters and value no greater than 256 characters.
	Tags map[string]*string `json:"tags,omitempty"`
}

// CreatorsClientCreateOrUpdateOptions contains the optional parameters for the CreatorsClient.CreateOrUpdate method.
type CreatorsClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// CreatorsClientDeleteOptions contains the optional parameters for the CreatorsClient.Delete method.
type CreatorsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// CreatorsClientGetOptions contains the optional parameters for the CreatorsClient.Get method.
type CreatorsClientGetOptions struct {
	// placeholder for future optional parameters
}

// CreatorsClientListByAccountOptions contains the optional parameters for the CreatorsClient.ListByAccount method.
type CreatorsClientListByAccountOptions struct {
	// placeholder for future optional parameters
}

// CreatorsClientUpdateOptions contains the optional parameters for the CreatorsClient.Update method.
type CreatorsClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// Dimension of map account, for example API Category, Api Name, Result Type, and Response Code.
type Dimension struct {
	// Display name of dimension.
	DisplayName *string `json:"displayName,omitempty"`

	// Internal metric name of the dimension.
	InternalMetricName *string `json:"internalMetricName,omitempty"`

	// Internal name of the dimension.
	InternalName *string `json:"internalName,omitempty"`

	// Display name of dimension.
	Name *string `json:"name,omitempty"`

	// Source Mdm Namespace of the dimension.
	SourceMdmNamespace *string `json:"sourceMdmNamespace,omitempty"`

	// Flag to indicate exporting to Azure Monitor.
	ToBeExportedToShoebox *bool `json:"toBeExportedToShoebox,omitempty"`
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*ErrorDetail `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations.
// (This also follows the OData error response format.).
type ErrorResponse struct {
	// The error object.
	Error *ErrorDetail `json:"error,omitempty"`
}

// KeySpecification - Whether the operation refers to the primary or secondary key.
type KeySpecification struct {
	// REQUIRED; Whether the operation refers to the primary or secondary key.
	KeyType *KeyType `json:"keyType,omitempty"`
}

// LinkedResource - Linked resource is reference to a resource deployed in an Azure subscription, add the linked resource
// uniqueName value as an optional parameter for operations on Azure Maps Geospatial REST APIs.
type LinkedResource struct {
	// REQUIRED; ARM resource id in the form: '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/accounts/{storageName}'.
	ID *string `json:"id,omitempty"`

	// REQUIRED; A provided name which uniquely identifies the linked resource.
	UniqueName *string `json:"uniqueName,omitempty"`
}

// ManagedServiceIdentity - Identity for the resource.
type ManagedServiceIdentity struct {
	// The identity type.
	Type *ResourceIdentityType `json:"type,omitempty"`

	// The list of user identities associated with the resource. The user identity dictionary key references will be ARM resource
	// ids in the form:
	// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'.
	UserAssignedIdentities map[string]*Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties `json:"userAssignedIdentities,omitempty"`

	// READ-ONLY; The principal ID of resource identity.
	PrincipalID *string `json:"principalId,omitempty" azure:"ro"`

	// READ-ONLY; The tenant ID of resource.
	TenantID *string `json:"tenantId,omitempty" azure:"ro"`
}

// MetricSpecification - Metric specification of operation.
type MetricSpecification struct {
	// Aggregation type could be Average.
	AggregationType *string `json:"aggregationType,omitempty"`

	// The category this metric specification belong to, could be Capacity.
	Category *string `json:"category,omitempty"`

	// Dimensions of map account.
	Dimensions []*Dimension `json:"dimensions,omitempty"`

	// Display description of metric specification.
	DisplayDescription *string `json:"displayDescription,omitempty"`

	// Display name of metric specification.
	DisplayName *string `json:"displayName,omitempty"`

	// The property to decide fill gap with zero or not.
	FillGapWithZero *bool `json:"fillGapWithZero,omitempty"`

	// Internal metric name.
	InternalMetricName *string `json:"internalMetricName,omitempty"`

	// Name of metric specification.
	Name *string `json:"name,omitempty"`

	// Account Resource Id.
	ResourceIDDimensionNameOverride *string `json:"resourceIdDimensionNameOverride,omitempty"`

	// Source metrics account.
	SourceMdmAccount *string `json:"sourceMdmAccount,omitempty"`

	// Unit could be Count.
	Unit *string `json:"unit,omitempty"`
}

// OperationDetail - Operation detail payload
type OperationDetail struct {
	// Display of the operation
	Display *OperationDisplay `json:"display,omitempty"`

	// Indicates whether the operation is a data action
	IsDataAction *bool `json:"isDataAction,omitempty"`

	// Name of the operation
	Name *string `json:"name,omitempty"`

	// Origin of the operation
	Origin *string `json:"origin,omitempty"`

	// Properties of the operation
	Properties *OperationProperties `json:"properties,omitempty"`
}

// OperationDisplay - Operation display payload
type OperationDisplay struct {
	// Localized friendly description for the operation
	Description *string `json:"description,omitempty"`

	// Localized friendly name for the operation
	Operation *string `json:"operation,omitempty"`

	// Resource provider of the operation
	Provider *string `json:"provider,omitempty"`

	// Resource of the operation
	Resource *string `json:"resource,omitempty"`
}

// OperationProperties - Properties of operation, include metric specifications.
type OperationProperties struct {
	// One property of operation, include metric specifications.
	ServiceSpecification *ServiceSpecification `json:"serviceSpecification,omitempty"`
}

// Operations - The set of operations available for Maps.
type Operations struct {
	// URL client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string `json:"nextLink,omitempty"`

	// READ-ONLY; An operation available for Maps.
	Value []*OperationDetail `json:"value,omitempty" azure:"ro"`
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

// SKU - The SKU of the Maps Account.
type SKU struct {
	// REQUIRED; The name of the SKU, in standard format (such as S0).
	Name *Name `json:"name,omitempty"`

	// READ-ONLY; Gets the sku tier. This is based on the SKU name.
	Tier *string `json:"tier,omitempty" azure:"ro"`
}

// ServiceSpecification - One property of operation, include metric specifications.
type ServiceSpecification struct {
	// Metric specifications of operation.
	MetricSpecifications []*MetricSpecification `json:"metricSpecifications,omitempty"`
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
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
