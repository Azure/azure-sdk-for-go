//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsignalr

import "time"

// CorsSettings - Cross-Origin Resource Sharing (CORS) settings.
type CorsSettings struct {
	// Gets or sets the list of origins that should be allowed to make cross-origin calls (for example: http://example.com:12345).
// Use "*" to allow all. If omitted, allow all by default.
	AllowedOrigins []*string
}

// CustomCertificate - A custom certificate.
type CustomCertificate struct {
	// REQUIRED; Custom certificate properties.
	Properties *CustomCertificateProperties

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// CustomCertificateList - Custom certificates list.
type CustomCertificateList struct {
	// The URL the client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string

	// List of custom certificates of this resource.
	Value []*CustomCertificate
}

// CustomCertificateProperties - Custom certificate properties.
type CustomCertificateProperties struct {
	// REQUIRED; Base uri of the KeyVault that stores certificate.
	KeyVaultBaseURI *string

	// REQUIRED; Certificate secret name.
	KeyVaultSecretName *string

	// Certificate secret version.
	KeyVaultSecretVersion *string

	// READ-ONLY; Provisioning state of the resource.
	ProvisioningState *ProvisioningState
}

// CustomDomain - A custom domain
type CustomDomain struct {
	// REQUIRED; Properties of a custom domain.
	Properties *CustomDomainProperties

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// CustomDomainList - Custom domains list
type CustomDomainList struct {
	// The URL the client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string

	// List of custom domains that bind to this resource.
	Value []*CustomDomain
}

// CustomDomainProperties - Properties of a custom domain.
type CustomDomainProperties struct {
	// REQUIRED; Reference to a resource.
	CustomCertificate *ResourceReference

	// REQUIRED; The custom domain name.
	DomainName *string

	// READ-ONLY; Provisioning state of the resource.
	ProvisioningState *ProvisioningState
}

// Dimension - Specifications of the Dimension of metrics.
type Dimension struct {
	// Localized friendly display name of the dimension.
	DisplayName *string

	// Name of the dimension as it appears in MDM.
	InternalName *string

	// The public facing name of the dimension.
	Name *string

	// A Boolean flag indicating whether this dimension should be included for the shoebox export scenario.
	ToBeExportedForShoebox *bool
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

// Feature of a resource, which controls the runtime behavior.
type Feature struct {
	// REQUIRED; FeatureFlags is the supported features of Azure SignalR service.
// * ServiceMode: Flag for backend server for SignalR service. Values allowed: "Default": have your own backend server; "Serverless":
// your application doesn't have a backend server; "Classic": for
// backward compatibility. Support both Default and Serverless mode but not recommended; "PredefinedOnly": for future use.
// * EnableConnectivityLogs: "true"/"false", to enable/disable the connectivity log category respectively.
// * EnableMessagingLogs: "true"/"false", to enable/disable the connectivity log category respectively.
// * EnableLiveTrace: Live Trace allows you to know what's happening inside Azure SignalR service, it will give you live traces
// in real time, it will be helpful when you developing your own Azure
// SignalR based web application or self-troubleshooting some issues. Please note that live traces are counted as outbound
// messages that will be charged. Values allowed: "true"/"false", to
// enable/disable live trace feature.
	Flag *FeatureFlags

	// REQUIRED; Value of the feature flag. See Azure SignalR service document https://docs.microsoft.com/azure/azure-signalr/
// for allowed values.
	Value *string

	// Optional properties related to this feature.
	Properties map[string]*string
}

// Keys - A class represents the access keys of the resource.
type Keys struct {
	// Connection string constructed via the primaryKey
	PrimaryConnectionString *string

	// The primary access key.
	PrimaryKey *string

	// Connection string constructed via the secondaryKey
	SecondaryConnectionString *string

	// The secondary access key.
	SecondaryKey *string
}

// LiveTraceCategory - Live trace category configuration of a Microsoft.SignalRService resource.
type LiveTraceCategory struct {
	// Indicates whether or the live trace category is enabled. Available values: true, false. Case insensitive.
	Enabled *string

	// Gets or sets the live trace category's name. Available values: ConnectivityLogs, MessagingLogs. Case insensitive.
	Name *string
}

// LiveTraceConfiguration - Live trace configuration of a Microsoft.SignalRService resource.
type LiveTraceConfiguration struct {
	// Gets or sets the list of category configurations.
	Categories []*LiveTraceCategory

	// Indicates whether or not enable live trace. When it's set to true, live trace client can connect to the service. Otherwise,
// live trace client can't connect to the service, so that you are unable to
// receive any log, no matter what you configure in "categories". Available values: true, false. Case insensitive.
	Enabled *string
}

// LogSpecification - Specifications of the Logs for Azure Monitoring.
type LogSpecification struct {
	// Localized friendly display name of the log.
	DisplayName *string

	// Name of the log.
	Name *string
}

// ManagedIdentity - A class represent managed identities used for request and response
type ManagedIdentity struct {
	// Represents the identity type: systemAssigned, userAssigned, None
	Type *ManagedIdentityType

	// Get or set the user assigned identities
	UserAssignedIdentities map[string]*UserAssignedIdentityProperty

	// READ-ONLY; Get the principal id for the system assigned identity. Only be used in response.
	PrincipalID *string

	// READ-ONLY; Get the tenant id for the system assigned identity. Only be used in response
	TenantID *string
}

// ManagedIdentitySettings - Managed identity settings for upstream.
type ManagedIdentitySettings struct {
	// The Resource indicating the App ID URI of the target resource. It also appears in the aud (audience) claim of the issued
// token.
	Resource *string
}

// MetricSpecification - Specifications of the Metrics for Azure Monitoring.
type MetricSpecification struct {
	// Only provide one value for this field. Valid values: Average, Minimum, Maximum, Total, Count.
	AggregationType *string

	// The name of the metric category that the metric belongs to. A metric can only belong to a single category.
	Category *string

	// The dimensions of the metrics.
	Dimensions []*Dimension

	// Localized friendly description of the metric.
	DisplayDescription *string

	// Localized friendly display name of the metric.
	DisplayName *string

	// Optional. If set to true, then zero will be returned for time duration where no metric is emitted/published. Ex. a metric
// that returns the number of times a particular error code was emitted. The
// error code may not appear often, instead of the RP publishing 0, Shoebox can auto fill in 0s for time periods where nothing
// was emitted.
	FillGapWithZero *string

	// Name of the metric.
	Name *string

	// The unit that makes sense for the metric.
	Unit *string
}

// NameAvailability - Result of the request to check name availability. It contains a flag and possible reason of failure.
type NameAvailability struct {
	// The message of the operation.
	Message *string

	// Indicates whether the name is available or not.
	NameAvailable *bool

	// The reason of the availability. Required if name is not available.
	Reason *string
}

// NameAvailabilityParameters - Data POST-ed to the nameAvailability action
type NameAvailabilityParameters struct {
	// REQUIRED; The resource name to validate. e.g."my-resource-name"
	Name *string

	// REQUIRED; The resource type. Can be "Microsoft.SignalRService/SignalR", "Microsoft.SignalRService/WebPubSub", "Microsoft.SignalRService/SignalR/replicas"
// or "Microsoft.SignalRService/WebPubSub/replicas"
	Type *string
}

// NetworkACL - Network ACL
type NetworkACL struct {
	// Allowed request types. The value can be one or more of: ClientConnection, ServerConnection, RESTAPI.
	Allow []*SignalRRequestType

	// Denied request types. The value can be one or more of: ClientConnection, ServerConnection, RESTAPI.
	Deny []*SignalRRequestType
}

// NetworkACLs - Network ACLs for the resource
type NetworkACLs struct {
	// Azure Networking ACL Action.
	DefaultAction *ACLAction

	// ACLs for requests from private endpoints
	PrivateEndpoints []*PrivateEndpointACL

	// Network ACL
	PublicNetwork *NetworkACL
}

// Operation - REST API operation supported by resource provider.
type Operation struct {
	// The object that describes a operation.
	Display *OperationDisplay

	// If the operation is a data action. (for data plane rbac)
	IsDataAction *bool

	// Name of the operation with format: {provider}/{resource}/{operation}
	Name *string

	// Optional. The intended executor of the operation; governs the display of the operation in the RBAC UX and the audit logs
// UX.
	Origin *string

	// Extra Operation properties.
	Properties *OperationProperties
}

// OperationDisplay - The object that describes a operation.
type OperationDisplay struct {
	// The localized friendly description for the operation
	Description *string

	// The localized friendly name for the operation.
	Operation *string

	// Friendly name of the resource provider
	Provider *string

	// Resource type on which the operation is performed.
	Resource *string
}

// OperationList - Result of the request to list REST API operations. It contains a list of operations.
type OperationList struct {
	// The URL the client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string

	// List of operations supported by the resource provider.
	Value []*Operation
}

// OperationProperties - Extra Operation properties.
type OperationProperties struct {
	// An object that describes a specification.
	ServiceSpecification *ServiceSpecification
}

// PrivateEndpoint - Private endpoint
type PrivateEndpoint struct {
	// Full qualified Id of the private endpoint
	ID *string
}

// PrivateEndpointACL - ACL for a private endpoint
type PrivateEndpointACL struct {
	// REQUIRED; Name of the private endpoint connection
	Name *string

	// Allowed request types. The value can be one or more of: ClientConnection, ServerConnection, RESTAPI.
	Allow []*SignalRRequestType

	// Denied request types. The value can be one or more of: ClientConnection, ServerConnection, RESTAPI.
	Deny []*SignalRRequestType
}

// PrivateEndpointConnection - A private endpoint connection to an azure resource
type PrivateEndpointConnection struct {
	// Private endpoint connection properties
	Properties *PrivateEndpointConnectionProperties

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// PrivateEndpointConnectionList - A list of private endpoint connections
type PrivateEndpointConnectionList struct {
	// Request URL that can be used to query next page of private endpoint connections. Returned when the total number of requested
// private endpoint connections exceed maximum page size.
	NextLink *string

	// The list of the private endpoint connections
	Value []*PrivateEndpointConnection
}

// PrivateEndpointConnectionProperties - Private endpoint connection properties
type PrivateEndpointConnectionProperties struct {
	// Private endpoint
	PrivateEndpoint *PrivateEndpoint

	// Connection state of the private endpoint connection
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState

	// READ-ONLY; Group IDs
	GroupIDs []*string

	// READ-ONLY; Provisioning state of the resource.
	ProvisioningState *ProvisioningState
}

// PrivateLinkResource - Private link resource
type PrivateLinkResource struct {
	// Private link resource properties
	Properties *PrivateLinkResourceProperties

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// PrivateLinkResourceList - Contains a list of PrivateLinkResource and a possible link to query more results
type PrivateLinkResourceList struct {
	// The URL the client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string

	// List of PrivateLinkResource
	Value []*PrivateLinkResource
}

// PrivateLinkResourceProperties - Private link resource properties
type PrivateLinkResourceProperties struct {
	// Group Id of the private link resource
	GroupID *string

	// Required members of the private link resource
	RequiredMembers []*string

	// Required private DNS zone names
	RequiredZoneNames []*string

	// The list of resources that are onboarded to private link service
	ShareablePrivateLinkResourceTypes []*ShareablePrivateLinkResourceType
}

// PrivateLinkServiceConnectionState - Connection state of the private endpoint connection
type PrivateLinkServiceConnectionState struct {
	// A message indicating if changes on the service provider require any updates on the consumer.
	ActionsRequired *string

	// The reason for approval/rejection of the connection.
	Description *string

	// Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.
	Status *PrivateLinkServiceConnectionStatus
}

// Properties - A class that describes the properties of the resource
type Properties struct {
	// Cross-Origin Resource Sharing (CORS) settings.
	Cors *CorsSettings

	// DisableLocalAuth Enable or disable aad auth When set as true, connection with AuthType=aad won't work.
	DisableAADAuth *bool

	// DisableLocalAuth Enable or disable local auth with AccessKey When set as true, connection with AccessKey=xxx won't work.
	DisableLocalAuth *bool

	// List of the featureFlags.
// FeatureFlags that are not included in the parameters for the update operation will not be modified. And the response will
// only include featureFlags that are explicitly set. When a featureFlag is not
// explicitly set, its globally default value will be used But keep in mind, the default value doesn't mean "false". It varies
// in terms of different FeatureFlags.
	Features []*Feature

	// Live trace configuration of a Microsoft.SignalRService resource.
	LiveTraceConfiguration *LiveTraceConfiguration

	// Network ACLs for the resource
	NetworkACLs *NetworkACLs

	// Enable or disable public network access. Default to "Enabled". When it's Enabled, network ACLs still apply. When it's Disabled,
// public network access is always disabled no matter what you set in
// network ACLs.
	PublicNetworkAccess *string

	// Resource log configuration of a Microsoft.SignalRService resource.
	ResourceLogConfiguration *ResourceLogConfiguration

	// Serverless settings.
	Serverless *ServerlessSettings

	// TLS settings for the resource
	TLS *TLSSettings

	// The settings for the Upstream when the service is in server-less mode.
	Upstream *ServerlessUpstreamSettings

	// READ-ONLY; The publicly accessible IP of the resource.
	ExternalIP *string

	// READ-ONLY; FQDN of the service instance.
	HostName *string

	// READ-ONLY; Deprecated.
	HostNamePrefix *string

	// READ-ONLY; Private endpoint connections to the resource.
	PrivateEndpointConnections []*PrivateEndpointConnection

	// READ-ONLY; Provisioning state of the resource.
	ProvisioningState *ProvisioningState

	// READ-ONLY; The publicly accessible port of the resource which is designed for browser/client side usage.
	PublicPort *int32

	// READ-ONLY; The publicly accessible port of the resource which is designed for customer server side usage.
	ServerPort *int32

	// READ-ONLY; The list of shared private link resources.
	SharedPrivateLinkResources []*SharedPrivateLinkResource

	// READ-ONLY; Version of the resource. Probably you need the same or higher version of client SDKs.
	Version *string
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

// RegenerateKeyParameters - Parameters describes the request to regenerate access keys
type RegenerateKeyParameters struct {
	// The type of access key.
	KeyType *KeyType
}

// Replica - A class represent a replica resource.
type Replica struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string
	Properties *ReplicaProperties

	// The billing information of the resource.
	SKU *ResourceSKU

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

type ReplicaList struct {
	// The URL the client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string

	// List of the replica
	Value []*Replica
}

type ReplicaProperties struct {
	// READ-ONLY; Provisioning state of the resource.
	ProvisioningState *ProvisioningState
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

// ResourceInfo - A class represent a resource.
type ResourceInfo struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// A class represent managed identities used for request and response
	Identity *ManagedIdentity

	// The kind of the service
	Kind *ServiceKind

	// A class that describes the properties of the resource
	Properties *Properties

	// The billing information of the resource.
	SKU *ResourceSKU

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ResourceInfoList - Object that includes an array of resources and a possible link for next set.
type ResourceInfoList struct {
	// The URL the client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string

	// List of the resources
	Value []*ResourceInfo
}

// ResourceLogCategory - Resource log category configuration of a Microsoft.SignalRService resource.
type ResourceLogCategory struct {
	// Indicates whether or the resource log category is enabled. Available values: true, false. Case insensitive.
	Enabled *string

	// Gets or sets the resource log category's name. Available values: ConnectivityLogs, MessagingLogs. Case insensitive.
	Name *string
}

// ResourceLogConfiguration - Resource log configuration of a Microsoft.SignalRService resource.
type ResourceLogConfiguration struct {
	// Gets or sets the list of category configurations.
	Categories []*ResourceLogCategory
}

// ResourceReference - Reference to a resource.
type ResourceReference struct {
	// Resource ID.
	ID *string
}

// ResourceSKU - The billing information of the resource.
type ResourceSKU struct {
	// REQUIRED; The name of the SKU. Required.
// Allowed values: StandardS1, FreeF1, Premium_P1
	Name *string

	// Optional, integer. The unit count of the resource. 1 by default.
// If present, following values are allowed: Free: 1; Standard: 1,2,3,4,5,6,7,8,9,10,20,30,40,50,60,70,80,90,100; Premium:
// 1,2,3,4,5,6,7,8,9,10,20,30,40,50,60,70,80,90,100;
	Capacity *int32

	// Optional tier of this particular SKU. 'Standard' or 'Free'.
// Basic is deprecated, use Standard instead.
	Tier *SignalRSKUTier

	// READ-ONLY; Not used. Retained for future use.
	Family *string

	// READ-ONLY; Not used. Retained for future use.
	Size *string
}

// SKU - Describes an available sku."
type SKU struct {
	// READ-ONLY; Describes scaling information of a sku.
	Capacity *SKUCapacity

	// READ-ONLY; The resource type that this object applies to
	ResourceType *string

	// READ-ONLY; The billing information of the resource.
	SKU *ResourceSKU
}

// SKUCapacity - Describes scaling information of a sku.
type SKUCapacity struct {
	// READ-ONLY; Allows capacity value list.
	AllowedValues []*int32

	// READ-ONLY; The default capacity.
	Default *int32

	// READ-ONLY; The highest permitted capacity for this resource
	Maximum *int32

	// READ-ONLY; The lowest permitted capacity for this resource
	Minimum *int32

	// READ-ONLY; The scale type applicable to the sku.
	ScaleType *ScaleType
}

// SKUList - The list skus operation response
type SKUList struct {
	// READ-ONLY; The URL the client should use to fetch the next page (per server side paging). It's null for now, added for
// future use.
	NextLink *string

	// READ-ONLY; The list of skus available for the resource.
	Value []*SKU
}

// ServerlessSettings - Serverless settings.
type ServerlessSettings struct {
	// Gets or sets Client Connection Timeout. Optional to be set. Value in seconds. Default value is 30 seconds. Customer should
// set the timeout to a shorter period if messages are expected to be sent in
// shorter intervals, and want the client to disconnect more quickly after the last message is sent. You can set the timeout
// to a longer period if messages are expected to be sent in longer intervals,
// and they want to keep the same client connection alive during this session. The service considers the client disconnected
// if it hasn't received a message (including keep-alive) in this interval.
	ConnectionTimeoutInSeconds *int32
}

// ServerlessUpstreamSettings - The settings for the Upstream when the service is in server-less mode.
type ServerlessUpstreamSettings struct {
	// Gets or sets the list of Upstream URL templates. Order matters, and the first matching template takes effects.
	Templates []*UpstreamTemplate
}

// ServiceSpecification - An object that describes a specification.
type ServiceSpecification struct {
	// Specifications of the Logs for Azure Monitoring.
	LogSpecifications []*LogSpecification

	// Specifications of the Metrics for Azure Monitoring.
	MetricSpecifications []*MetricSpecification
}

// ShareablePrivateLinkResourceProperties - Describes the properties of a resource type that has been onboarded to private
// link service
type ShareablePrivateLinkResourceProperties struct {
	// The description of the resource type that has been onboarded to private link service
	Description *string

	// The resource provider group id for the resource that has been onboarded to private link service
	GroupID *string

	// The resource provider type for the resource that has been onboarded to private link service
	Type *string
}

// ShareablePrivateLinkResourceType - Describes a resource type that has been onboarded to private link service
type ShareablePrivateLinkResourceType struct {
	// The name of the resource type that has been onboarded to private link service
	Name *string

	// Describes the properties of a resource type that has been onboarded to private link service
	Properties *ShareablePrivateLinkResourceProperties
}

// SharedPrivateLinkResource - Describes a Shared Private Link Resource
type SharedPrivateLinkResource struct {
	// Describes the properties of an existing Shared Private Link Resource
	Properties *SharedPrivateLinkResourceProperties

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// SharedPrivateLinkResourceList - A list of shared private link resources
type SharedPrivateLinkResourceList struct {
	// Request URL that can be used to query next page of private endpoint connections. Returned when the total number of requested
// private endpoint connections exceed maximum page size.
	NextLink *string

	// The list of the shared private link resources
	Value []*SharedPrivateLinkResource
}

// SharedPrivateLinkResourceProperties - Describes the properties of an existing Shared Private Link Resource
type SharedPrivateLinkResourceProperties struct {
	// REQUIRED; The group id from the provider of resource the shared private link resource is for
	GroupID *string

	// REQUIRED; The resource id of the resource the shared private link resource is for
	PrivateLinkResourceID *string

	// The request message for requesting approval of the shared private link resource
	RequestMessage *string

	// READ-ONLY; Provisioning state of the resource.
	ProvisioningState *ProvisioningState

	// READ-ONLY; Status of the shared private link resource
	Status *SharedPrivateLinkResourceStatus
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

// TLSSettings - TLS settings for the resource
type TLSSettings struct {
	// Request client certificate during TLS handshake if enabled. Not supported for free tier. Any input will be ignored for
// free tier.
	ClientCertEnabled *bool
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. E.g. "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}"
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// UpstreamAuthSettings - Upstream auth settings. If not set, no auth is used for upstream messages.
type UpstreamAuthSettings struct {
	// Managed identity settings for upstream.
	ManagedIdentity *ManagedIdentitySettings

	// Upstream auth type enum.
	Type *UpstreamAuthType
}

// UpstreamTemplate - Upstream template item settings. It defines the Upstream URL of the incoming requests. The template
// defines the pattern of the event, the hub or the category of the incoming request that matches
// current URL template.
type UpstreamTemplate struct {
	// REQUIRED; Gets or sets the Upstream URL template. You can use 3 predefined parameters {hub}, {category} {event} inside
// the template, the value of the Upstream URL is dynamically calculated when the client
// request comes in. For example, if the urlTemplate is http://example.com/{hub}/api/{event}, with a client request from hub
// chat connects, it will first POST to this URL:
// http://example.com/chat/api/connect.
	URLTemplate *string

	// Upstream auth settings. If not set, no auth is used for upstream messages.
	Auth *UpstreamAuthSettings

	// Gets or sets the matching pattern for category names. If not set, it matches any category. There are 3 kind of patterns
// supported: 1. "*", it to matches any category name. 2. Combine multiple
// categories with ",", for example "connections,messages", it matches category "connections" and "messages". 3. The single
// category name, for example, "connections", it matches the category
// "connections".
	CategoryPattern *string

	// Gets or sets the matching pattern for event names. If not set, it matches any event. There are 3 kind of patterns supported:
// 1. "*", it to matches any event name. 2. Combine multiple events with ",",
// for example "connect,disconnect", it matches event "connect" and "disconnect". 3. The single event name, for example, "connect",
// it matches "connect".
	EventPattern *string

	// Gets or sets the matching pattern for hub names. If not set, it matches any hub. There are 3 kind of patterns supported:
// 1. "*", it to matches any hub name. 2. Combine multiple hubs with ",", for
// example "hub1,hub2", it matches "hub1" and "hub2". 3. The single hub name, for example, "hub1", it matches "hub1".
	HubPattern *string
}

// Usage - Object that describes a specific usage of the resources.
type Usage struct {
	// Current value for the usage quota.
	CurrentValue *int64

	// Fully qualified ARM resource id
	ID *string

	// The maximum permitted value for the usage quota. If there is no limit, this value will be -1.
	Limit *int64

	// Localizable String object containing the name and a localized value.
	Name *UsageName

	// Representing the units of the usage quota. Possible values are: Count, Bytes, Seconds, Percent, CountPerSecond, BytesPerSecond.
	Unit *string
}

// UsageList - Object that includes an array of the resource usages and a possible link for next set.
type UsageList struct {
	// The URL the client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string

	// List of the resource usages
	Value []*Usage
}

// UsageName - Localizable String object containing the name and a localized value.
type UsageName struct {
	// Localized name of the usage.
	LocalizedValue *string

	// The identifier of the usage.
	Value *string
}

// UserAssignedIdentityProperty - Properties of user assigned identity.
type UserAssignedIdentityProperty struct {
	// READ-ONLY; Get the client id for the user assigned identity
	ClientID *string

	// READ-ONLY; Get the principal id for the user assigned identity
	PrincipalID *string
}

