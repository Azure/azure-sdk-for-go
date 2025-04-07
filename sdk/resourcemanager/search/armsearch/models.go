// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsearch

import "time"

// AdminKeyResult - Response containing the primary and secondary admin API keys for a given Azure AI Search service.
type AdminKeyResult struct {
	// READ-ONLY; The primary admin API key of the search service.
	PrimaryKey *string

	// READ-ONLY; The secondary admin API key of the search service.
	SecondaryKey *string
}

// CheckNameAvailabilityInput - Input of check name availability API.
type CheckNameAvailabilityInput struct {
	// REQUIRED; The search service name to validate. Search service names must only contain lowercase letters, digits or dashes,
	// cannot use dash as the first two or last one characters, cannot contain consecutive
	// dashes, and must be between 2 and 60 characters in length.
	Name *string

	// CONSTANT; The type of the resource whose name is to be validated. This value must always be 'searchServices'.
	// Field has constant value "searchServices", any specified value is ignored.
	Type *string
}

// CheckNameAvailabilityOutput - Output of check name availability API.
type CheckNameAvailabilityOutput struct {
	// READ-ONLY; A value indicating whether the name is available.
	IsNameAvailable *bool

	// READ-ONLY; A message that explains why the name is invalid and provides resource naming requirements. Available only if
	// 'Invalid' is returned in the 'reason' property.
	Message *string

	// READ-ONLY; The reason why the name is not available. 'Invalid' indicates the name provided does not match the naming requirements
	// (incorrect length, unsupported characters, etc.). 'AlreadyExists' indicates that
	// the name is already in use and is therefore unavailable.
	Reason *UnavailableNameReason
}

// DataPlaneAADOrAPIKeyAuthOption - Indicates that either the API key or an access token from a Microsoft Entra ID tenant
// can be used for authentication.
type DataPlaneAADOrAPIKeyAuthOption struct {
	// Describes what response the data plane API of a search service would send for requests that failed authentication.
	AADAuthFailureMode *AADAuthFailureMode
}

// DataPlaneAuthOptions - Defines the options for how the search service authenticates a data plane request. This cannot be
// set if 'disableLocalAuth' is set to true.
type DataPlaneAuthOptions struct {
	// Indicates that either the API key or an access token from a Microsoft Entra ID tenant can be used for authentication.
	AADOrAPIKey *DataPlaneAADOrAPIKeyAuthOption

	// Indicates that only the API key can be used for authentication.
	APIKeyOnly any
}

// EncryptionWithCmk - Describes a policy that determines how resources within the search service are to be encrypted with
// customer managed keys.
type EncryptionWithCmk struct {
	// Describes how a search service should enforce compliance if it finds objects that aren't encrypted with the customer-managed
	// key.
	Enforcement *SearchEncryptionWithCmk

	// READ-ONLY; Returns the status of search service compliance with respect to non-CMK-encrypted objects. If a service has
	// more than one unencrypted object, and enforcement is enabled, the service is marked as
	// noncompliant.
	EncryptionComplianceStatus *SearchEncryptionComplianceStatus
}

type FeatureOffering struct {
	// The name of the feature offered in this region.
	Name *FeatureName
}

// IPRule - The IP restriction rule of the Azure AI Search service.
type IPRule struct {
	// Value corresponding to a single IPv4 address (eg., 123.1.2.3) or an IP range in CIDR format (eg., 123.1.2.3/24) to be allowed.
	Value *string
}

// Identity - Details about the search service identity. A null value indicates that the search service has no identity assigned.
type Identity struct {
	// REQUIRED; The type of identity used for the resource. The type 'SystemAssigned, UserAssigned' includes both an identity
	// created by the system and a set of user assigned identities. The type 'None' will remove
	// all identities from the service.
	Type *IdentityType

	// The list of user identities associated with the resource. The user identity dictionary key references will be ARM resource
	// IDs in the form:
	// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'.
	UserAssignedIdentities map[string]*UserAssignedManagedIdentity

	// READ-ONLY; The principal ID of the system-assigned identity of the search service.
	PrincipalID *string

	// READ-ONLY; The tenant ID of the system-assigned identity of the search service.
	TenantID *string
}

// ListQueryKeysResult - Response containing the query API keys for a given Azure AI Search service.
type ListQueryKeysResult struct {
	// READ-ONLY; Request URL that can be used to query next page of query keys. Returned when the total number of requested query
	// keys exceed maximum page size.
	NextLink *string

	// READ-ONLY; The query keys for the Azure AI Search service.
	Value []*QueryKey
}

// NSPConfigAccessRule - An access rule for a network security perimeter configuration.
type NSPConfigAccessRule struct {
	Name *string

	// The properties for the access rules in a network security perimeter configuration.
	Properties *NSPConfigAccessRuleProperties
}

// NSPConfigAccessRuleProperties - The properties for the access rules in a network security perimeter configuration.
type NSPConfigAccessRuleProperties struct {
	AddressPrefixes           []*string
	Direction                 *string
	FullyQualifiedDomainNames []*string
	NetworkSecurityPerimeters []*NSPConfigNetworkSecurityPerimeterRule
	Subscriptions             []*string
}

// NSPConfigAssociation - The resource association for the network security perimeter.
type NSPConfigAssociation struct {
	AccessMode *string
	Name       *string
}

// NSPConfigNetworkSecurityPerimeterRule - The network security perimeter properties present in a configuration rule.
type NSPConfigNetworkSecurityPerimeterRule struct {
	ID            *string
	Location      *string
	PerimeterGUID *string
}

// NSPConfigPerimeter - The perimeter for a network security perimeter configuration.
type NSPConfigPerimeter struct {
	ID            *string
	Location      *string
	PerimeterGUID *string
}

// NSPConfigProfile - The profile for a network security perimeter configuration.
type NSPConfigProfile struct {
	AccessRules        []*NSPConfigAccessRule
	AccessRulesVersion *string
	Name               *string
}

// NSPProvisioningIssue - An object to describe any issues with provisioning network security perimeters to a search service.
type NSPProvisioningIssue struct {
	Name *string

	// The properties to describe any issues with provisioning network security perimeters to a search service.
	Properties *NSPProvisioningIssueProperties
}

// NSPProvisioningIssueProperties - The properties to describe any issues with provisioning network security perimeters to
// a search service.
type NSPProvisioningIssueProperties struct {
	Description          *string
	IssueType            *string
	Severity             *string
	SuggestedAccessRules []*string
	SuggestedResourceIDs []*string
}

// NetworkRuleSet - Network specific rules that determine how the Azure AI Search service may be reached.
type NetworkRuleSet struct {
	// Possible origins of inbound traffic that can bypass the rules defined in the 'ipRules' section.
	Bypass *SearchBypass

	// A list of IP restriction rules that defines the inbound network(s) with allowing access to the search service endpoint.
	// At the meantime, all other public IP networks are blocked by the firewall. These
	// restriction rules are applied only when the 'publicNetworkAccess' of the search service is 'enabled'; otherwise, traffic
	// over public interface is not allowed even with any public IP rules, and private
	// endpoint connections would be the exclusive access method.
	IPRules []*IPRule
}

// NetworkSecurityPerimeterConfiguration - Network security perimeter configuration for a server.
type NetworkSecurityPerimeterConfiguration struct {
	// Resource properties.
	Properties *NetworkSecurityPerimeterConfigurationProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// NetworkSecurityPerimeterConfigurationListResult - A list of network security perimeter configurations for a server.
type NetworkSecurityPerimeterConfigurationListResult struct {
	// READ-ONLY; Link to retrieve next page of results.
	NextLink *string

	// READ-ONLY; Array of results.
	Value []*NetworkSecurityPerimeterConfiguration
}

// NetworkSecurityPerimeterConfigurationProperties - The properties of a network security perimeter configuration.
type NetworkSecurityPerimeterConfigurationProperties struct {
	// The perimeter for a network security perimeter configuration.
	NetworkSecurityPerimeter *NSPConfigPerimeter

	// The profile for a network security perimeter configuration.
	Profile            *NSPConfigProfile
	ProvisioningIssues []*NSPProvisioningIssue

	// The resource association for the network security perimeter.
	ResourceAssociation *NSPConfigAssociation

	// READ-ONLY
	ProvisioningState *string
}

type OfferingsByRegion struct {
	// The list of features offered in this region.
	Features []*FeatureOffering

	// The name of the region.
	RegionName *string

	// The list of SKUs offered in this region.
	SKUs []*SKUOffering
}

// OfferingsListResult - The response containing a list of features and SKUs offered in various regions.
type OfferingsListResult struct {
	// The list of regions with their respective features and SKUs offered.
	Value []*OfferingsByRegion

	// READ-ONLY; The URL to get the next set of offerings, if any.
	NextLink *string
}

// Operation - Describes a REST API operation.
type Operation struct {
	// READ-ONLY; The object that describes the operation.
	Display *OperationDisplay

	// READ-ONLY; Describes if the specified operation is a data plane API operation. Operations where this value is not true
	// are supported directly by the resource provider.
	IsDataAction *bool

	// READ-ONLY; The name of the operation. This name is of the form {provider}/{resource}/{operation}.
	Name *string

	// READ-ONLY; Describes which originating entities are allowed to invoke this operation.
	Origin *string

	// READ-ONLY; Describes additional properties for this operation.
	Properties *OperationProperties
}

// OperationAvailability - Describes a particular availability for the metric specification.
type OperationAvailability struct {
	// READ-ONLY; The blob duration for the dimension.
	BlobDuration *string

	// READ-ONLY; The time grain for the dimension.
	TimeGrain *string
}

// OperationDisplay - The object that describes the operation.
type OperationDisplay struct {
	// READ-ONLY; The friendly name of the operation.
	Description *string

	// READ-ONLY; The operation type: read, write, delete, listKeys/action, etc.
	Operation *string

	// READ-ONLY; The friendly name of the resource provider.
	Provider *string

	// READ-ONLY; The resource type on which the operation is performed.
	Resource *string
}

// OperationListResult - The result of the request to list REST API operations. It contains a list of operations and a URL
// to get the next set of results.
type OperationListResult struct {
	// READ-ONLY; The URL to get the next set of operation list results, if any.
	NextLink *string

	// READ-ONLY; The list of operations by Azure AI Search, some supported by the resource provider and others by data plane
	// APIs.
	Value []*Operation
}

// OperationLogsSpecification - Specifications of one type of log for this operation.
type OperationLogsSpecification struct {
	// READ-ONLY; The blob duration for the log specification.
	BlobDuration *string

	// READ-ONLY; The display name of the log specification.
	DisplayName *string

	// READ-ONLY; The name of the log specification.
	Name *string
}

// OperationMetricDimension - Describes a particular dimension for the metric specification.
type OperationMetricDimension struct {
	// READ-ONLY; The display name of the dimension.
	DisplayName *string

	// READ-ONLY; The name of the dimension.
	Name *string
}

// OperationMetricsSpecification - Specifications of one type of metric for this operation.
type OperationMetricsSpecification struct {
	// READ-ONLY; The type of aggregation for the metric specification.
	AggregationType *string

	// READ-ONLY; Availabilities for the metric specification.
	Availabilities []*OperationAvailability

	// READ-ONLY; Dimensions for the metric specification.
	Dimensions []*OperationMetricDimension

	// READ-ONLY; The display description of the metric specification.
	DisplayDescription *string

	// READ-ONLY; The display name of the metric specification.
	DisplayName *string

	// READ-ONLY; The name of the metric specification.
	Name *string

	// READ-ONLY; The unit for the metric specification.
	Unit *string
}

// OperationProperties - Describes additional properties for this operation.
type OperationProperties struct {
	// READ-ONLY; Specifications of the service for this operation.
	ServiceSpecification *OperationServiceSpecification
}

// OperationServiceSpecification - Specifications of the service for this operation.
type OperationServiceSpecification struct {
	// READ-ONLY; Specifications of logs for this operation.
	LogSpecifications []*OperationLogsSpecification

	// READ-ONLY; Specifications of metrics for this operation.
	MetricSpecifications []*OperationMetricsSpecification
}

// PrivateEndpointConnection - Describes an existing private endpoint connection to the Azure AI Search service.
type PrivateEndpointConnection struct {
	// Describes the properties of an existing private endpoint connection to the Azure AI Search service.
	Properties *PrivateEndpointConnectionProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// PrivateEndpointConnectionListResult - Response containing a list of private endpoint connections.
type PrivateEndpointConnectionListResult struct {
	// READ-ONLY; Request URL that can be used to query next page of private endpoint connections. Returned when the total number
	// of requested private endpoint connections exceed maximum page size.
	NextLink *string

	// READ-ONLY; The list of private endpoint connections.
	Value []*PrivateEndpointConnection
}

// PrivateEndpointConnectionProperties - Describes the properties of an existing private endpoint connection to the search
// service.
type PrivateEndpointConnectionProperties struct {
	// The group ID of the Azure resource for which the private link service is for.
	GroupID *string

	// The private endpoint resource from Microsoft.Network provider.
	PrivateEndpoint *PrivateEndpointConnectionPropertiesPrivateEndpoint

	// Describes the current state of an existing Azure Private Link service connection to the private endpoint.
	PrivateLinkServiceConnectionState *PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState

	// The provisioning state of the private link service connection. Valid values are Updating, Deleting, Failed, Succeeded,
	// Incomplete, or Canceled.
	ProvisioningState *PrivateLinkServiceConnectionProvisioningState
}

// PrivateEndpointConnectionPropertiesPrivateEndpoint - The private endpoint resource from Microsoft.Network provider.
type PrivateEndpointConnectionPropertiesPrivateEndpoint struct {
	// The resource ID of the private endpoint resource from Microsoft.Network provider.
	ID *string
}

// PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState - Describes the current state of an existing Azure
// Private Link service connection to the private endpoint.
type PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState struct {
	// A description of any extra actions that may be required.
	ActionsRequired *string

	// The description for the private link service connection state.
	Description *string

	// Status of the the private link service connection. Valid values are Pending, Approved, Rejected, or Disconnected.
	Status *PrivateLinkServiceConnectionStatus
}

// PrivateLinkResource - Describes a supported private link resource for the Azure AI Search service.
type PrivateLinkResource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Describes the properties of a supported private link resource for the Azure AI Search service.
	Properties *PrivateLinkResourceProperties

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// PrivateLinkResourceProperties - Describes the properties of a supported private link resource for the Azure AI Search service.
// For a given API version, this represents the 'supported' groupIds when creating a shared private link
// resource.
type PrivateLinkResourceProperties struct {
	// READ-ONLY; The group ID of the private link resource.
	GroupID *string

	// READ-ONLY; The list of required members of the private link resource.
	RequiredMembers []*string

	// READ-ONLY; The list of required DNS zone names of the private link resource.
	RequiredZoneNames []*string

	// READ-ONLY; The list of resources that are onboarded to private link service, that are supported by Azure AI Search.
	ShareablePrivateLinkResourceTypes []*ShareablePrivateLinkResourceType
}

// PrivateLinkResourcesResult - Response containing a list of supported Private Link Resources.
type PrivateLinkResourcesResult struct {
	// READ-ONLY; The list of supported Private Link Resources.
	Value []*PrivateLinkResource
}

// QueryKey - Describes an API key for a given Azure AI Search service that conveys read-only permissions on the docs collection
// of an index.
type QueryKey struct {
	// READ-ONLY; The value of the query API key.
	Key *string

	// READ-ONLY; The name of the query API key. Query names are optional, but assigning a name can help you remember how it's
	// used.
	Name *string
}

// QuotaUsageResult - Describes the quota usage for a particular SKU.
type QuotaUsageResult struct {
	// The currently used up value for the particular search SKU.
	CurrentValue *int32

	// The resource ID of the quota usage SKU endpoint for Microsoft.Search provider.
	ID *string

	// The quota limit for the particular search SKU.
	Limit *int32

	// The unit of measurement for the search SKU.
	Unit *string

	// READ-ONLY; The name of the SKU supported by Azure AI Search.
	Name *QuotaUsageResultName
}

// QuotaUsageResultName - The name of the SKU supported by Azure AI Search.
type QuotaUsageResultName struct {
	// The localized string value for the SKU name.
	LocalizedValue *string

	// The SKU name supported by Azure AI Search.
	Value *string
}

// QuotaUsagesListResult - Response containing the quota usage information for all the supported SKUs of Azure AI Search.
type QuotaUsagesListResult struct {
	// READ-ONLY; Request URL that can be used to query next page of quota usages. Returned when the total number of requested
	// quota usages exceed maximum page size.
	NextLink *string

	// READ-ONLY; The quota usages for the SKUs supported by Azure AI Search.
	Value []*QuotaUsageResult
}

// SKU - Defines the SKU of a search service, which determines billing rate and capacity limits.
type SKU struct {
	// The SKU of the search service. Valid values include: 'free': Shared service. 'basic': Dedicated service with up to 3 replicas.
	// 'standard': Dedicated service with up to 12 partitions and 12 replicas.
	// 'standard2': Similar to standard, but with more capacity per search unit. 'standard3': The largest Standard offering with
	// up to 12 partitions and 12 replicas (or up to 3 partitions with more indexes
	// if you also set the hostingMode property to 'highDensity'). 'storageoptimizedl1': Supports 1TB per partition, up to 12
	// partitions. 'storageoptimizedl2': Supports 2TB per partition, up to 12
	// partitions.'
	Name *SKUName
}

type SKUOffering struct {
	// The limits associated with this SKU offered in this region.
	Limits *SKUOfferingLimits

	// Defines the SKU of a search service, which determines billing rate and capacity limits.
	SKU *SKU
}

// SKUOfferingLimits - The limits associated with this SKU offered in this region.
type SKUOfferingLimits struct {
	// The maximum number of indexers available for this SKU.
	Indexers *int32

	// The maximum number of indexes available for this SKU.
	Indexes *int32

	// The maximum storage size in Gigabytes available for this SKU per partition.
	PartitionStorageInGigabytes *float32

	// The maximum vector storage size in Gigabytes available for this SKU per partition.
	PartitionVectorStorageInGigabytes *float32

	// The maximum number of partitions available for this SKU.
	Partitions *int32

	// The maximum number of replicas available for this SKU.
	Replicas *int32

	// The maximum number of search units available for this SKU.
	SearchUnits *int32
}

// Service - Describes an Azure AI Search service and its current state.
type Service struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// The identity of the resource.
	Identity *Identity

	// Properties of the search service.
	Properties *ServiceProperties

	// The SKU of the search service, which determines price tier and capacity limits. This property is required when creating
	// a new search service.
	SKU *SKU

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata of the search service containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ServiceListResult - Response containing a list of Azure AI Search services.
type ServiceListResult struct {
	// READ-ONLY; Request URL that can be used to query next page of search services. Returned when the total number of requested
	// search services exceed maximum page size.
	NextLink *string

	// READ-ONLY; The list of search services.
	Value []*Service
}

// ServiceProperties - Properties of the search service.
type ServiceProperties struct {
	// Defines the options for how the data plane API of a search service authenticates requests. This cannot be set if 'disableLocalAuth'
	// is set to true.
	AuthOptions *DataPlaneAuthOptions

	// Configure this property to support the search service using either the default compute or Azure Confidential Compute.
	ComputeType *ComputeType

	// When set to true, calls to the search service will not be permitted to utilize API keys for authentication. This cannot
	// be set to true if 'dataPlaneAuthOptions' are defined.
	DisableLocalAuth *bool

	// A list of data exfiltration scenarios that are explicitly disallowed for the search service. Currently, the only supported
	// value is 'All' to disable all possible data export scenarios with more fine
	// grained controls planned for the future.
	DisabledDataExfiltrationOptions []*SearchDisabledDataExfiltrationOption

	// Specifies any policy regarding encryption of resources (such as indexes) using customer manager keys within a search service.
	EncryptionWithCmk *EncryptionWithCmk

	// The endpoint of the Azure AI Search service.
	Endpoint *string

	// Applicable only for the standard3 SKU. You can set this property to enable up to 3 high density partitions that allow up
	// to 1000 indexes, which is much higher than the maximum indexes allowed for any
	// other SKU. For the standard3 SKU, the value is either 'default' or 'highDensity'. For all other SKUs, this value must be
	// 'default'.
	HostingMode *HostingMode

	// Network specific rules that determine how the Azure AI Search service may be reached.
	NetworkRuleSet *NetworkRuleSet

	// The number of partitions in the search service; if specified, it can be 1, 2, 3, 4, 6, or 12. Values greater than 1 are
	// only valid for standard SKUs. For 'standard3' services with hostingMode set to
	// 'highDensity', the allowed values are between 1 and 3.
	PartitionCount *int32

	// This value can be set to 'enabled' to avoid breaking changes on existing customer resources and templates. If set to 'disabled',
	// traffic over public interface is not allowed, and private endpoint
	// connections would be the exclusive access method.
	PublicNetworkAccess *PublicNetworkAccess

	// The number of replicas in the search service. If specified, it must be a value between 1 and 12 inclusive for standard
	// SKUs or between 1 and 3 inclusive for basic SKU.
	ReplicaCount *int32

	// Sets options that control the availability of semantic search. This configuration is only possible for certain Azure AI
	// Search SKUs in certain locations.
	SemanticSearch *SearchSemanticSearch

	// READ-ONLY; A system generated property representing the service's etag that can be for optimistic concurrency control during
	// updates.
	ETag *string

	// READ-ONLY; The list of private endpoint connections to the Azure AI Search service.
	PrivateEndpointConnections []*PrivateEndpointConnection

	// READ-ONLY; The state of the last provisioning operation performed on the search service. Provisioning is an intermediate
	// state that occurs while service capacity is being established. After capacity is set up,
	// provisioningState changes to either 'Succeeded' or 'Failed'. Client applications can poll provisioning status (the recommended
	// polling interval is from 30 seconds to one minute) by using the Get
	// Search Service operation to see when an operation is completed. If you are using the free service, this value tends to
	// come back as 'Succeeded' directly in the call to Create search service. This is
	// because the free service uses capacity that is already set up.
	ProvisioningState *ProvisioningState

	// READ-ONLY; The date and time the search service was last upgraded. This field will be null until the service gets upgraded
	// for the first time.
	ServiceUpgradeDate *time.Time

	// READ-ONLY; The list of shared private link resources managed by the Azure AI Search service.
	SharedPrivateLinkResources []*SharedPrivateLinkResource

	// READ-ONLY; The status of the search service. Possible values include: 'running': The search service is running and no provisioning
	// operations are underway. 'provisioning': The search service is being provisioned
	// or scaled up or down. 'deleting': The search service is being deleted. 'degraded': The search service is degraded. This
	// can occur when the underlying search units are not healthy. The search service
	// is most likely operational, but performance might be slow and some requests might be dropped. 'disabled': The search service
	// is disabled. In this state, the service will reject all API requests.
	// 'error': The search service is in an error state. 'stopped': The search service is in a subscription that's disabled. If
	// your service is in the degraded, disabled, or error states, it means the Azure
	// AI Search team is actively investigating the underlying issue. Dedicated services in these states are still chargeable
	// based on the number of search units provisioned.
	Status *SearchServiceStatus

	// READ-ONLY; The details of the search service status.
	StatusDetails *string

	// READ-ONLY; Indicates whether or not the search service has an upgrade available.
	UpgradeAvailable *bool
}

// ServiceUpdate - The parameters used to update an Azure AI Search service.
type ServiceUpdate struct {
	// Details about the search service identity. A null value indicates that the search service has no identity assigned.
	Identity *Identity

	// The geographic location of the resource. This must be one of the supported and registered Azure geo regions (for example,
	// West US, East US, Southeast Asia, and so forth). This property is required
	// when creating a new resource.
	Location *string

	// Properties of the search service.
	Properties *ServiceProperties

	// The SKU of the search service, which determines price tier and capacity limits. This property is required when creating
	// a new search service.
	SKU *SKU

	// Tags to help categorize the resource in the Azure portal.
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

// ShareablePrivateLinkResourceProperties - Describes the properties of a resource type that has been onboarded to private
// link service, supported by Azure AI Search.
type ShareablePrivateLinkResourceProperties struct {
	// READ-ONLY; The description of the resource type that has been onboarded to private link service, supported by Azure AI
	// Search.
	Description *string

	// READ-ONLY; The resource provider group id for the resource that has been onboarded to private link service, supported by
	// Azure AI Search.
	GroupID *string

	// READ-ONLY; The resource provider type for the resource that has been onboarded to private link service, supported by Azure
	// AI Search.
	Type *string
}

// ShareablePrivateLinkResourceType - Describes an resource type that has been onboarded to private link service, supported
// by Azure AI Search.
type ShareablePrivateLinkResourceType struct {
	// READ-ONLY; The name of the resource type that has been onboarded to private link service, supported by Azure AI Search.
	Name *string

	// READ-ONLY; Describes the properties of a resource type that has been onboarded to private link service, supported by Azure
	// AI Search.
	Properties *ShareablePrivateLinkResourceProperties
}

// SharedPrivateLinkResource - Describes a shared private link resource managed by the Azure AI Search service.
type SharedPrivateLinkResource struct {
	// Describes the properties of a shared private link resource managed by the Azure AI Search service.
	Properties *SharedPrivateLinkResourceProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// SharedPrivateLinkResourceListResult - Response containing a list of shared private link resources.
type SharedPrivateLinkResourceListResult struct {
	// The URL to get the next set of shared private link resources, if there are any.
	NextLink *string

	// READ-ONLY; The list of shared private link resources.
	Value []*SharedPrivateLinkResource
}

// SharedPrivateLinkResourceProperties - Describes the properties of an existing shared private link resource managed by the
// Azure AI Search service.
type SharedPrivateLinkResourceProperties struct {
	// The group ID from the provider of resource the shared private link resource is for.
	GroupID *string

	// The resource ID of the resource the shared private link resource is for.
	PrivateLinkResourceID *string

	// The provisioning state of the shared private link resource. Valid values are Updating, Deleting, Failed, Succeeded or Incomplete.
	ProvisioningState *SharedPrivateLinkResourceProvisioningState

	// The message for requesting approval of the shared private link resource.
	RequestMessage *string

	// Optional. Can be used to specify the Azure Resource Manager location of the resource for which a shared private link is
	// being created. This is only required for those resources whose DNS configuration
	// are regional (such as Azure Kubernetes Service).
	ResourceRegion *string

	// Status of the shared private link resource. Valid values are Pending, Approved, Rejected or Disconnected.
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

// UserAssignedManagedIdentity - The details of the user assigned managed identity assigned to the search service.
type UserAssignedManagedIdentity struct {
	// READ-ONLY; The client ID of user assigned identity.
	ClientID *string

	// READ-ONLY; The principal ID of user assigned identity.
	PrincipalID *string
}
