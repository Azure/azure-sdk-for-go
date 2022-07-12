//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armrecoveryservices

import "time"

// AzureMonitorAlertSettings - Settings for Azure Monitor based alerts
type AzureMonitorAlertSettings struct {
	AlertsForAllJobFailures *AlertsState `json:"alertsForAllJobFailures,omitempty"`
}

// CertificateRequest - Details of the certificate to be uploaded to the vault.
type CertificateRequest struct {
	// Raw certificate data.
	Properties *RawCertificateData `json:"properties,omitempty"`
}

// CheckNameAvailabilityParameters - Resource Name availability input parameters - Resource type and resource name
type CheckNameAvailabilityParameters struct {
	// Resource name for which availability needs to be checked
	Name *string `json:"name,omitempty"`

	// Describes the Resource type: Microsoft.RecoveryServices/Vaults
	Type *string `json:"type,omitempty"`
}

// CheckNameAvailabilityResult - Response for check name availability API. Resource provider will set availability as true
// | false.
type CheckNameAvailabilityResult struct {
	Message       *string `json:"message,omitempty"`
	NameAvailable *bool   `json:"nameAvailable,omitempty"`
	Reason        *string `json:"reason,omitempty"`
}

// ClassicAlertSettings - Settings for classic alerts
type ClassicAlertSettings struct {
	AlertsForCriticalOperations *AlertsState `json:"alertsForCriticalOperations,omitempty"`
}

// ClientCheckNameAvailabilityOptions contains the optional parameters for the Client.CheckNameAvailability method.
type ClientCheckNameAvailabilityOptions struct {
	// placeholder for future optional parameters
}

// ClientDiscoveryDisplay - Localized display information of an operation.
type ClientDiscoveryDisplay struct {
	// Description of the operation having details of what operation is about.
	Description *string `json:"description,omitempty"`

	// Operations Name itself.
	Operation *string `json:"operation,omitempty"`

	// Name of the provider for display purposes
	Provider *string `json:"provider,omitempty"`

	// ResourceType for which this Operation can be performed.
	Resource *string `json:"resource,omitempty"`
}

// ClientDiscoveryForLogSpecification - Class to represent shoebox log specification in json client discovery.
type ClientDiscoveryForLogSpecification struct {
	// Blobs created in customer storage account per hour
	BlobDuration *string `json:"blobDuration,omitempty"`

	// Localized display name
	DisplayName *string `json:"displayName,omitempty"`

	// Name of the log.
	Name *string `json:"name,omitempty"`
}

// ClientDiscoveryForProperties - Class to represent shoebox properties in json client discovery.
type ClientDiscoveryForProperties struct {
	// Operation properties.
	ServiceSpecification *ClientDiscoveryForServiceSpecification `json:"serviceSpecification,omitempty"`
}

// ClientDiscoveryForServiceSpecification - Class to represent shoebox service specification in json client discovery.
type ClientDiscoveryForServiceSpecification struct {
	// List of log specifications of this operation.
	LogSpecifications []*ClientDiscoveryForLogSpecification `json:"logSpecifications,omitempty"`
}

// ClientDiscoveryResponse - Operations List response which contains list of available APIs.
type ClientDiscoveryResponse struct {
	// Link to the next chunk of the response
	NextLink *string `json:"nextLink,omitempty"`

	// List of available operations.
	Value []*ClientDiscoveryValueForSingleAPI `json:"value,omitempty"`
}

// ClientDiscoveryValueForSingleAPI - Available operation details.
type ClientDiscoveryValueForSingleAPI struct {
	// Contains the localized display information for this particular operation
	Display *ClientDiscoveryDisplay `json:"display,omitempty"`

	// Name of the Operation.
	Name *string `json:"name,omitempty"`

	// The intended executor of the operation;governs the display of the operation in the RBAC UX and the audit logs UX
	Origin *string `json:"origin,omitempty"`

	// ShoeBox properties for the given operation.
	Properties *ClientDiscoveryForProperties `json:"properties,omitempty"`
}

// CloudError - An error response from Azure Backup.
type CloudError struct {
	// The resource management error response.
	Error *Error `json:"error,omitempty"`
}

// CmkKekIdentity - The details of the identity used for CMK
type CmkKekIdentity struct {
	// Indicate that system assigned identity should be used. Mutually exclusive with 'userAssignedIdentity' field
	UseSystemAssignedIdentity *bool `json:"useSystemAssignedIdentity,omitempty"`

	// The user assigned identity to be used to grant permissions in case the type of identity used is UserAssigned
	UserAssignedIdentity *string `json:"userAssignedIdentity,omitempty"`
}

// CmkKeyVaultProperties - The properties of the Key Vault which hosts CMK
type CmkKeyVaultProperties struct {
	// The key uri of the Customer Managed Key
	KeyURI *string `json:"keyUri,omitempty"`
}

// Error - The resource management error response.
type Error struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*Error `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// IdentityData - Identity for the resource.
type IdentityData struct {
	// REQUIRED; The type of managed identity used. The type 'SystemAssigned, UserAssigned' includes both an implicitly created
	// identity and a set of user-assigned identities. The type 'None' will remove any
	// identities.
	Type *ResourceIdentityType `json:"type,omitempty"`

	// The list of user-assigned identities associated with the resource. The user-assigned identity dictionary keys will be ARM
	// resource ids in the form:
	// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'.
	UserAssignedIdentities map[string]*UserIdentity `json:"userAssignedIdentities,omitempty"`

	// READ-ONLY; The principal ID of resource identity.
	PrincipalID *string `json:"principalId,omitempty" azure:"ro"`

	// READ-ONLY; The tenant ID of resource.
	TenantID *string `json:"tenantId,omitempty" azure:"ro"`
}

// JobsSummary - Summary of the replication job data for this vault.
type JobsSummary struct {
	// Count of failed jobs.
	FailedJobs *int32 `json:"failedJobs,omitempty"`

	// Count of in-progress jobs.
	InProgressJobs *int32 `json:"inProgressJobs,omitempty"`

	// Count of suspended jobs.
	SuspendedJobs *int32 `json:"suspendedJobs,omitempty"`
}

// MonitoringSettings - Monitoring Settings of the vault
type MonitoringSettings struct {
	// Settings for Azure Monitor based alerts
	AzureMonitorAlertSettings *AzureMonitorAlertSettings `json:"azureMonitorAlertSettings,omitempty"`

	// Settings for classic alerts
	ClassicAlertSettings *ClassicAlertSettings `json:"classicAlertSettings,omitempty"`
}

// MonitoringSummary - Summary of the replication monitoring data for this vault.
type MonitoringSummary struct {
	// Count of all deprecated recovery service providers.
	DeprecatedProviderCount *int32 `json:"deprecatedProviderCount,omitempty"`

	// Count of all critical warnings.
	EventsCount *int32 `json:"eventsCount,omitempty"`

	// Count of all the supported recovery service providers.
	SupportedProviderCount *int32 `json:"supportedProviderCount,omitempty"`

	// Count of unhealthy replication providers.
	UnHealthyProviderCount *int32 `json:"unHealthyProviderCount,omitempty"`

	// Count of unhealthy VMs.
	UnHealthyVMCount *int32 `json:"unHealthyVmCount,omitempty"`

	// Count of all the unsupported recovery service providers.
	UnsupportedProviderCount *int32 `json:"unsupportedProviderCount,omitempty"`
}

// NameInfo - The name of usage.
type NameInfo struct {
	// Localized value of usage.
	LocalizedValue *string `json:"localizedValue,omitempty"`

	// Value of usage.
	Value *string `json:"value,omitempty"`
}

// OperationResource - Operation Resource
type OperationResource struct {
	// End time of the operation
	EndTime *time.Time `json:"endTime,omitempty"`

	// Required if status == failed or status == canceled. This is the OData v4 error format, used by the RPC and will go into
	// the v2.2 Azure REST API guidelines.
	Error *Error `json:"error,omitempty"`

	// It should match what is used to GET the operation result
	ID *string `json:"id,omitempty"`

	// It must match the last segment of the "id" field, and will typically be a GUID / system generated value
	Name *string `json:"name,omitempty"`

	// Start time of the operation
	StartTime *time.Time `json:"startTime,omitempty"`

	// The status of the operation. (InProgress/Success/Failed/Cancelled)
	Status *string `json:"status,omitempty"`
}

// OperationsClientGetOperationResultOptions contains the optional parameters for the OperationsClient.GetOperationResult
// method.
type OperationsClientGetOperationResultOptions struct {
	// placeholder for future optional parameters
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.List method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// OperationsClientOperationStatusGetOptions contains the optional parameters for the OperationsClient.OperationStatusGet
// method.
type OperationsClientOperationStatusGetOptions struct {
	// placeholder for future optional parameters
}

// PatchTrackedResource - Tracked resource with location.
type PatchTrackedResource struct {
	// Optional ETag.
	Etag *string `json:"etag,omitempty"`

	// Resource location.
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Resource Id represents the complete path to the resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name associated with the resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Resource type represents the complete path of the form Namespace/ResourceType/ResourceType/…
	Type *string `json:"type,omitempty" azure:"ro"`
}

// PatchVault - Patch Resource information, as returned by the resource provider.
type PatchVault struct {
	// Optional ETag.
	Etag *string `json:"etag,omitempty"`

	// Identity for the resource.
	Identity *IdentityData `json:"identity,omitempty"`

	// Resource location.
	Location *string `json:"location,omitempty"`

	// Properties of the vault.
	Properties *VaultProperties `json:"properties,omitempty"`

	// Identifies the unique system identifier for each Azure resource.
	SKU *SKU `json:"sku,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Resource Id represents the complete path to the resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name associated with the resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Resource type represents the complete path of the form Namespace/ResourceType/ResourceType/…
	Type *string `json:"type,omitempty" azure:"ro"`
}

// PrivateEndpoint - The Private Endpoint network resource that is linked to the Private Endpoint connection.
type PrivateEndpoint struct {
	// READ-ONLY; Gets or sets id.
	ID *string `json:"id,omitempty" azure:"ro"`
}

// PrivateEndpointConnection - Private Endpoint Connection Response Properties.
type PrivateEndpointConnection struct {
	// READ-ONLY; The Private Endpoint network resource that is linked to the Private Endpoint connection.
	PrivateEndpoint *PrivateEndpoint `json:"privateEndpoint,omitempty" azure:"ro"`

	// READ-ONLY; Gets or sets private link service connection state.
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty" azure:"ro"`

	// READ-ONLY; Gets or sets provisioning state of the private endpoint connection.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// PrivateEndpointConnectionVaultProperties - Information to be stored in Vault properties as an element of privateEndpointConnections
// List.
type PrivateEndpointConnectionVaultProperties struct {
	// READ-ONLY; Format of id subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.[Service]/{resource}/{resourceName}/privateEndpointConnections/{connectionName}.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The location of the private Endpoint connection
	Location *string `json:"location,omitempty" azure:"ro"`

	// READ-ONLY; The name of the private Endpoint Connection
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Private Endpoint Connection Response Properties.
	Properties *PrivateEndpointConnection `json:"properties,omitempty" azure:"ro"`

	// READ-ONLY; The type, which will be of the format, Microsoft.RecoveryServices/vaults/privateEndpointConnections
	Type *string `json:"type,omitempty" azure:"ro"`
}

// PrivateLinkResource - Information of the private link resource.
type PrivateLinkResource struct {
	// Resource properties
	Properties *PrivateLinkResourceProperties `json:"properties,omitempty"`

	// READ-ONLY; Fully qualified identifier of the resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Name of the resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; e.g. Microsoft.RecoveryServices/vaults/privateLinkResources
	Type *string `json:"type,omitempty" azure:"ro"`
}

// PrivateLinkResourceProperties - Properties of the private link resource.
type PrivateLinkResourceProperties struct {
	// READ-ONLY; e.g. f9ad6492-33d4-4690-9999-6bfd52a0d081 (Backup) or f9ad6492-33d4-4690-9999-6bfd52a0d082 (SiteRecovery)
	GroupID *string `json:"groupId,omitempty" azure:"ro"`

	// READ-ONLY; [backup-ecs1, backup-prot1, backup-prot1b, backup-prot1c, backup-id1]
	RequiredMembers []*string `json:"requiredMembers,omitempty" azure:"ro"`

	// READ-ONLY; The private link resource Private link DNS zone name.
	RequiredZoneNames []*string `json:"requiredZoneNames,omitempty" azure:"ro"`
}

// PrivateLinkResources - Class which represent the stamps associated with the vault.
type PrivateLinkResources struct {
	// Link to the next chunk of the response
	NextLink *string `json:"nextLink,omitempty"`

	// A collection of private link resources
	Value []*PrivateLinkResource `json:"value,omitempty"`
}

// PrivateLinkResourcesClientGetOptions contains the optional parameters for the PrivateLinkResourcesClient.Get method.
type PrivateLinkResourcesClientGetOptions struct {
	// placeholder for future optional parameters
}

// PrivateLinkResourcesClientListOptions contains the optional parameters for the PrivateLinkResourcesClient.List method.
type PrivateLinkResourcesClientListOptions struct {
	// placeholder for future optional parameters
}

// PrivateLinkServiceConnectionState - Gets or sets private link service connection state.
type PrivateLinkServiceConnectionState struct {
	// READ-ONLY; Gets or sets actions required.
	ActionsRequired *string `json:"actionsRequired,omitempty" azure:"ro"`

	// READ-ONLY; Gets or sets description.
	Description *string `json:"description,omitempty" azure:"ro"`

	// READ-ONLY; Gets or sets the status.
	Status *PrivateEndpointConnectionStatus `json:"status,omitempty" azure:"ro"`
}

// RawCertificateData - Raw certificate data.
type RawCertificateData struct {
	// Specifies the authentication type.
	AuthType *AuthType `json:"authType,omitempty"`

	// The base64 encoded certificate raw data string
	Certificate []byte `json:"certificate,omitempty"`
}

// RegisteredIdentitiesClientDeleteOptions contains the optional parameters for the RegisteredIdentitiesClient.Delete method.
type RegisteredIdentitiesClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// ReplicationUsage - Replication usages of a vault.
type ReplicationUsage struct {
	// Summary of the replication jobs data for this vault.
	JobsSummary *JobsSummary `json:"jobsSummary,omitempty"`

	// Summary of the replication monitoring data for this vault.
	MonitoringSummary *MonitoringSummary `json:"monitoringSummary,omitempty"`

	// Number of replication protected items for this vault.
	ProtectedItemCount *int32 `json:"protectedItemCount,omitempty"`

	// Number of replication recovery plans for this vault.
	RecoveryPlanCount *int32 `json:"recoveryPlanCount,omitempty"`

	// The authentication type of recovery service providers in the vault.
	RecoveryServicesProviderAuthType *int32 `json:"recoveryServicesProviderAuthType,omitempty"`

	// Number of servers registered to this vault.
	RegisteredServersCount *int32 `json:"registeredServersCount,omitempty"`
}

// ReplicationUsageList - Replication usages for vault.
type ReplicationUsageList struct {
	// The list of replication usages for the given vault.
	Value []*ReplicationUsage `json:"value,omitempty"`
}

// ReplicationUsagesClientListOptions contains the optional parameters for the ReplicationUsagesClient.List method.
type ReplicationUsagesClientListOptions struct {
	// placeholder for future optional parameters
}

// Resource - ARM Resource.
type Resource struct {
	// Optional ETag.
	Etag *string `json:"etag,omitempty"`

	// READ-ONLY; Resource Id represents the complete path to the resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name associated with the resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Resource type represents the complete path of the form Namespace/ResourceType/ResourceType/…
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ResourceCertificateAndAADDetails - Certificate details representing the Vault credentials for AAD.
type ResourceCertificateAndAADDetails struct {
	// REQUIRED; AAD tenant authority.
	AADAuthority *string `json:"aadAuthority,omitempty"`

	// REQUIRED; AAD tenant Id.
	AADTenantID *string `json:"aadTenantId,omitempty"`

	// REQUIRED; This property will be used as the discriminator for deciding the specific types in the polymorphic chain of types.
	AuthType *string `json:"authType,omitempty"`

	// REQUIRED; Azure Management Endpoint Audience.
	AzureManagementEndpointAudience *string `json:"azureManagementEndpointAudience,omitempty"`

	// REQUIRED; AAD service principal clientId.
	ServicePrincipalClientID *string `json:"servicePrincipalClientId,omitempty"`

	// REQUIRED; AAD service principal ObjectId.
	ServicePrincipalObjectID *string `json:"servicePrincipalObjectId,omitempty"`

	// The base64 encoded certificate raw data string.
	Certificate []byte `json:"certificate,omitempty"`

	// Certificate friendly name.
	FriendlyName *string `json:"friendlyName,omitempty"`

	// Certificate issuer.
	Issuer *string `json:"issuer,omitempty"`

	// Resource ID of the vault.
	ResourceID *int64 `json:"resourceId,omitempty"`

	// Service Resource Id.
	ServiceResourceID *string `json:"serviceResourceId,omitempty"`

	// Certificate Subject Name.
	Subject *string `json:"subject,omitempty"`

	// Certificate thumbprint.
	Thumbprint *string `json:"thumbprint,omitempty"`

	// Certificate Validity start Date time.
	ValidFrom *time.Time `json:"validFrom,omitempty"`

	// Certificate Validity End Date time.
	ValidTo *time.Time `json:"validTo,omitempty"`
}

// GetResourceCertificateDetails implements the ResourceCertificateDetailsClassification interface for type ResourceCertificateAndAADDetails.
func (r *ResourceCertificateAndAADDetails) GetResourceCertificateDetails() *ResourceCertificateDetails {
	return &ResourceCertificateDetails{
		AuthType:     r.AuthType,
		Certificate:  r.Certificate,
		FriendlyName: r.FriendlyName,
		Issuer:       r.Issuer,
		ResourceID:   r.ResourceID,
		Subject:      r.Subject,
		Thumbprint:   r.Thumbprint,
		ValidFrom:    r.ValidFrom,
		ValidTo:      r.ValidTo,
	}
}

// ResourceCertificateAndAcsDetails - Certificate details representing the Vault credentials for ACS.
type ResourceCertificateAndAcsDetails struct {
	// REQUIRED; This property will be used as the discriminator for deciding the specific types in the polymorphic chain of types.
	AuthType *string `json:"authType,omitempty"`

	// REQUIRED; Acs mgmt host name to connect to.
	GlobalAcsHostName *string `json:"globalAcsHostName,omitempty"`

	// REQUIRED; ACS namespace name - tenant for our service.
	GlobalAcsNamespace *string `json:"globalAcsNamespace,omitempty"`

	// REQUIRED; Global ACS namespace RP realm.
	GlobalAcsRPRealm *string `json:"globalAcsRPRealm,omitempty"`

	// The base64 encoded certificate raw data string.
	Certificate []byte `json:"certificate,omitempty"`

	// Certificate friendly name.
	FriendlyName *string `json:"friendlyName,omitempty"`

	// Certificate issuer.
	Issuer *string `json:"issuer,omitempty"`

	// Resource ID of the vault.
	ResourceID *int64 `json:"resourceId,omitempty"`

	// Certificate Subject Name.
	Subject *string `json:"subject,omitempty"`

	// Certificate thumbprint.
	Thumbprint *string `json:"thumbprint,omitempty"`

	// Certificate Validity start Date time.
	ValidFrom *time.Time `json:"validFrom,omitempty"`

	// Certificate Validity End Date time.
	ValidTo *time.Time `json:"validTo,omitempty"`
}

// GetResourceCertificateDetails implements the ResourceCertificateDetailsClassification interface for type ResourceCertificateAndAcsDetails.
func (r *ResourceCertificateAndAcsDetails) GetResourceCertificateDetails() *ResourceCertificateDetails {
	return &ResourceCertificateDetails{
		AuthType:     r.AuthType,
		Certificate:  r.Certificate,
		FriendlyName: r.FriendlyName,
		Issuer:       r.Issuer,
		ResourceID:   r.ResourceID,
		Subject:      r.Subject,
		Thumbprint:   r.Thumbprint,
		ValidFrom:    r.ValidFrom,
		ValidTo:      r.ValidTo,
	}
}

// ResourceCertificateDetailsClassification provides polymorphic access to related types.
// Call the interface's GetResourceCertificateDetails() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *ResourceCertificateAndAADDetails, *ResourceCertificateAndAcsDetails, *ResourceCertificateDetails
type ResourceCertificateDetailsClassification interface {
	// GetResourceCertificateDetails returns the ResourceCertificateDetails content of the underlying type.
	GetResourceCertificateDetails() *ResourceCertificateDetails
}

// ResourceCertificateDetails - Certificate details representing the Vault credentials.
type ResourceCertificateDetails struct {
	// REQUIRED; This property will be used as the discriminator for deciding the specific types in the polymorphic chain of types.
	AuthType *string `json:"authType,omitempty"`

	// The base64 encoded certificate raw data string.
	Certificate []byte `json:"certificate,omitempty"`

	// Certificate friendly name.
	FriendlyName *string `json:"friendlyName,omitempty"`

	// Certificate issuer.
	Issuer *string `json:"issuer,omitempty"`

	// Resource ID of the vault.
	ResourceID *int64 `json:"resourceId,omitempty"`

	// Certificate Subject Name.
	Subject *string `json:"subject,omitempty"`

	// Certificate thumbprint.
	Thumbprint *string `json:"thumbprint,omitempty"`

	// Certificate Validity start Date time.
	ValidFrom *time.Time `json:"validFrom,omitempty"`

	// Certificate Validity End Date time.
	ValidTo *time.Time `json:"validTo,omitempty"`
}

// GetResourceCertificateDetails implements the ResourceCertificateDetailsClassification interface for type ResourceCertificateDetails.
func (r *ResourceCertificateDetails) GetResourceCertificateDetails() *ResourceCertificateDetails {
	return r
}

// SKU - Identifies the unique system identifier for each Azure resource.
type SKU struct {
	// REQUIRED; The Sku name.
	Name *SKUName `json:"name,omitempty"`

	// The sku capacity
	Capacity *string `json:"capacity,omitempty"`

	// The sku family
	Family *string `json:"family,omitempty"`

	// The sku size
	Size *string `json:"size,omitempty"`

	// The Sku tier.
	Tier *string `json:"tier,omitempty"`
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}

// TrackedResource - Tracked resource with location.
type TrackedResource struct {
	// REQUIRED; Resource location.
	Location *string `json:"location,omitempty"`

	// Optional ETag.
	Etag *string `json:"etag,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Resource Id represents the complete path to the resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name associated with the resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Resource type represents the complete path of the form Namespace/ResourceType/ResourceType/…
	Type *string `json:"type,omitempty" azure:"ro"`
}

// UpgradeDetails - Details for upgrading vault.
type UpgradeDetails struct {
	// READ-ONLY; UTC time at which the upgrade operation has ended.
	EndTimeUTC *time.Time `json:"endTimeUtc,omitempty" azure:"ro"`

	// READ-ONLY; UTC time at which the upgrade operation status was last updated.
	LastUpdatedTimeUTC *time.Time `json:"lastUpdatedTimeUtc,omitempty" azure:"ro"`

	// READ-ONLY; Message to the user containing information about the upgrade operation.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; ID of the vault upgrade operation.
	OperationID *string `json:"operationId,omitempty" azure:"ro"`

	// READ-ONLY; Resource ID of the vault before the upgrade.
	PreviousResourceID *string `json:"previousResourceId,omitempty" azure:"ro"`

	// READ-ONLY; UTC time at which the upgrade operation has started.
	StartTimeUTC *time.Time `json:"startTimeUtc,omitempty" azure:"ro"`

	// READ-ONLY; Status of the vault upgrade operation.
	Status *VaultUpgradeState `json:"status,omitempty" azure:"ro"`

	// READ-ONLY; The way the vault upgrade was triggered.
	TriggerType *TriggerType `json:"triggerType,omitempty" azure:"ro"`

	// READ-ONLY; Resource ID of the upgraded vault.
	UpgradedResourceID *string `json:"upgradedResourceId,omitempty" azure:"ro"`
}

// UsagesClientListByVaultsOptions contains the optional parameters for the UsagesClient.ListByVaults method.
type UsagesClientListByVaultsOptions struct {
	// placeholder for future optional parameters
}

// UserIdentity - A resource identity that is managed by the user of the service.
type UserIdentity struct {
	// READ-ONLY; The client ID of the user-assigned identity.
	ClientID *string `json:"clientId,omitempty" azure:"ro"`

	// READ-ONLY; The principal ID of the user-assigned identity.
	PrincipalID *string `json:"principalId,omitempty" azure:"ro"`
}

// Vault - Resource information, as returned by the resource provider.
type Vault struct {
	// REQUIRED; Resource location.
	Location *string `json:"location,omitempty"`

	// Optional ETag.
	Etag *string `json:"etag,omitempty"`

	// Identity for the resource.
	Identity *IdentityData `json:"identity,omitempty"`

	// Properties of the vault.
	Properties *VaultProperties `json:"properties,omitempty"`

	// Identifies the unique system identifier for each Azure resource.
	SKU *SKU `json:"sku,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Resource Id represents the complete path to the resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name associated with the resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; Resource type represents the complete path of the form Namespace/ResourceType/ResourceType/…
	Type *string `json:"type,omitempty" azure:"ro"`
}

// VaultCertificateResponse - Certificate corresponding to a vault that can be used by clients to register themselves with
// the vault.
type VaultCertificateResponse struct {
	// Certificate details representing the Vault credentials.
	Properties ResourceCertificateDetailsClassification `json:"properties,omitempty"`

	// READ-ONLY; Resource Id represents the complete path to the resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name associated with the resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Resource type represents the complete path of the form Namespace/ResourceType/ResourceType/…
	Type *string `json:"type,omitempty" azure:"ro"`
}

// VaultCertificatesClientCreateOptions contains the optional parameters for the VaultCertificatesClient.Create method.
type VaultCertificatesClientCreateOptions struct {
	// placeholder for future optional parameters
}

// VaultExtendedInfo - Vault extended information.
type VaultExtendedInfo struct {
	// Algorithm for Vault ExtendedInfo
	Algorithm *string `json:"algorithm,omitempty"`

	// Encryption key.
	EncryptionKey *string `json:"encryptionKey,omitempty"`

	// Encryption key thumbprint.
	EncryptionKeyThumbprint *string `json:"encryptionKeyThumbprint,omitempty"`

	// Integrity key.
	IntegrityKey *string `json:"integrityKey,omitempty"`
}

// VaultExtendedInfoClientCreateOrUpdateOptions contains the optional parameters for the VaultExtendedInfoClient.CreateOrUpdate
// method.
type VaultExtendedInfoClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// VaultExtendedInfoClientGetOptions contains the optional parameters for the VaultExtendedInfoClient.Get method.
type VaultExtendedInfoClientGetOptions struct {
	// placeholder for future optional parameters
}

// VaultExtendedInfoClientUpdateOptions contains the optional parameters for the VaultExtendedInfoClient.Update method.
type VaultExtendedInfoClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// VaultExtendedInfoResource - Vault extended information.
type VaultExtendedInfoResource struct {
	// Optional ETag.
	Etag *string `json:"etag,omitempty"`

	// Vault extended information.
	Properties *VaultExtendedInfo `json:"properties,omitempty"`

	// READ-ONLY; Resource Id represents the complete path to the resource.
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; Resource name associated with the resource.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Resource type represents the complete path of the form Namespace/ResourceType/ResourceType/…
	Type *string `json:"type,omitempty" azure:"ro"`
}

// VaultList - The response model for a list of Vaults.
type VaultList struct {
	Value []*Vault `json:"value,omitempty"`

	// READ-ONLY
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`
}

// VaultProperties - Properties of the vault.
type VaultProperties struct {
	// Customer Managed Key details of the resource.
	Encryption *VaultPropertiesEncryption `json:"encryption,omitempty"`

	// Monitoring Settings of the vault
	MonitoringSettings *MonitoringSettings `json:"monitoringSettings,omitempty"`

	// The details of the latest move operation performed on the Azure Resource
	MoveDetails *VaultPropertiesMoveDetails `json:"moveDetails,omitempty"`

	// Details for upgrading vault.
	UpgradeDetails *UpgradeDetails `json:"upgradeDetails,omitempty"`

	// READ-ONLY; Backup storage version
	BackupStorageVersion *BackupStorageVersion `json:"backupStorageVersion,omitempty" azure:"ro"`

	// READ-ONLY; The State of the Resource after the move operation
	MoveState *ResourceMoveState `json:"moveState,omitempty" azure:"ro"`

	// READ-ONLY; List of private endpoint connection.
	PrivateEndpointConnections []*PrivateEndpointConnectionVaultProperties `json:"privateEndpointConnections,omitempty" azure:"ro"`

	// READ-ONLY; Private endpoint state for backup.
	PrivateEndpointStateForBackup *VaultPrivateEndpointState `json:"privateEndpointStateForBackup,omitempty" azure:"ro"`

	// READ-ONLY; Private endpoint state for site recovery.
	PrivateEndpointStateForSiteRecovery *VaultPrivateEndpointState `json:"privateEndpointStateForSiteRecovery,omitempty" azure:"ro"`

	// READ-ONLY; Provisioning State.
	ProvisioningState *string `json:"provisioningState,omitempty" azure:"ro"`
}

// VaultPropertiesEncryption - Customer Managed Key details of the resource.
type VaultPropertiesEncryption struct {
	// Enabling/Disabling the Double Encryption state
	InfrastructureEncryption *InfrastructureEncryptionState `json:"infrastructureEncryption,omitempty"`

	// The details of the identity used for CMK
	KekIdentity *CmkKekIdentity `json:"kekIdentity,omitempty"`

	// The properties of the Key Vault which hosts CMK
	KeyVaultProperties *CmkKeyVaultProperties `json:"keyVaultProperties,omitempty"`
}

// VaultPropertiesMoveDetails - The details of the latest move operation performed on the Azure Resource
type VaultPropertiesMoveDetails struct {
	// READ-ONLY; End Time of the Resource Move Operation
	CompletionTimeUTC *time.Time `json:"completionTimeUtc,omitempty" azure:"ro"`

	// READ-ONLY; OperationId of the Resource Move Operation
	OperationID *string `json:"operationId,omitempty" azure:"ro"`

	// READ-ONLY; Source Resource of the Resource Move Operation
	SourceResourceID *string `json:"sourceResourceId,omitempty" azure:"ro"`

	// READ-ONLY; Start Time of the Resource Move Operation
	StartTimeUTC *time.Time `json:"startTimeUtc,omitempty" azure:"ro"`

	// READ-ONLY; Target Resource of the Resource Move Operation
	TargetResourceID *string `json:"targetResourceId,omitempty" azure:"ro"`
}

// VaultUsage - Usages of a vault.
type VaultUsage struct {
	// Current value of usage.
	CurrentValue *int64 `json:"currentValue,omitempty"`

	// Limit of usage.
	Limit *int64 `json:"limit,omitempty"`

	// Name of usage.
	Name *NameInfo `json:"name,omitempty"`

	// Next reset time of usage.
	NextResetTime *time.Time `json:"nextResetTime,omitempty"`

	// Quota period of usage.
	QuotaPeriod *string `json:"quotaPeriod,omitempty"`

	// Unit of the usage.
	Unit *UsagesUnit `json:"unit,omitempty"`
}

// VaultUsageList - Usage for vault.
type VaultUsageList struct {
	// The list of usages for the given vault.
	Value []*VaultUsage `json:"value,omitempty"`
}

// VaultsClientBeginCreateOrUpdateOptions contains the optional parameters for the VaultsClient.BeginCreateOrUpdate method.
type VaultsClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VaultsClientBeginUpdateOptions contains the optional parameters for the VaultsClient.BeginUpdate method.
type VaultsClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VaultsClientDeleteOptions contains the optional parameters for the VaultsClient.Delete method.
type VaultsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// VaultsClientGetOptions contains the optional parameters for the VaultsClient.Get method.
type VaultsClientGetOptions struct {
	// placeholder for future optional parameters
}

// VaultsClientListByResourceGroupOptions contains the optional parameters for the VaultsClient.ListByResourceGroup method.
type VaultsClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// VaultsClientListBySubscriptionIDOptions contains the optional parameters for the VaultsClient.ListBySubscriptionID method.
type VaultsClientListBySubscriptionIDOptions struct {
	// placeholder for future optional parameters
}
