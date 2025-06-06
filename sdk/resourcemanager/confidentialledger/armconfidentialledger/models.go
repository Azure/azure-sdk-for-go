// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armconfidentialledger

import "time"

// AADBasedSecurityPrincipal - AAD based security principal with associated Ledger RoleName
type AADBasedSecurityPrincipal struct {
	// LedgerRole associated with the Security Principal of Ledger
	LedgerRoleName *LedgerRoleName

	// UUID/GUID based Principal Id of the Security Principal
	PrincipalID *string

	// UUID/GUID based Tenant Id of the Security Principal
	TenantID *string
}

// Backup - Object representing Backup properties of a Confidential Ledger Resource.
type Backup struct {
	// REQUIRED; SAS URI used to access the backup Fileshare.
	URI *string

	// The region where the backup of the ledger will eventually be restored to.
	RestoreRegion *string
}

// BackupResponse - Object representing the backup response of a Confidential Ledger Resource.
type BackupResponse struct {
	// READ-ONLY; Response body stating if the ledger is being backed up.
	Message *string
}

// CertBasedSecurityPrincipal - Cert based security principal with Ledger RoleName
type CertBasedSecurityPrincipal struct {
	// Public key of the user cert (.pem or .cer)
	Cert *string

	// LedgerRole associated with the Security Principal of Ledger
	LedgerRoleName *LedgerRoleName
}

// CheckNameAvailabilityRequest - The check availability request body.
type CheckNameAvailabilityRequest struct {
	// The name of the resource for which availability needs to be checked.
	Name *string

	// The resource type.
	Type *string
}

// CheckNameAvailabilityResponse - The check availability result.
type CheckNameAvailabilityResponse struct {
	// Detailed reason why the given name is available.
	Message *string

	// Indicates if the resource name is available.
	NameAvailable *bool

	// The reason why the given name is not available.
	Reason *CheckNameAvailabilityReason
}

// ConfidentialLedger - Confidential Ledger. Contains the properties of Confidential Ledger Resource.
type ConfidentialLedger struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// Properties of Confidential Ledger Resource.
	Properties *LedgerProperties

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

// DeploymentType - Object representing DeploymentType for Managed CCF.
type DeploymentType struct {
	// Source Uri containing ManagedCCF code
	AppSourceURI *string

	// Unique name for the Managed CCF.
	LanguageRuntime *LanguageRuntime
}

// LedgerProperties - Additional Confidential Ledger properties.
type LedgerProperties struct {
	// Array of all AAD based Security Principals.
	AADBasedSecurityPrincipals []*AADBasedSecurityPrincipal

	// Application type of the Confidential Ledger. Default: "Standard".
	// Expected values: "Standard", "Premium".
	ApplicationType *ApplicationType

	// Array of all cert based Security Principals.
	CertBasedSecurityPrincipals []*CertBasedSecurityPrincipal

	// Enclave platform of the Confidential Ledger.
	EnclavePlatform *EnclavePlatform

	// CCF Property for the logging level for the untrusted host: Trace, Debug, Info, Fail, Fatal.
	HostLevel *string

	// SKU associated with the ledger
	LedgerSKU *LedgerSKU

	// Type of Confidential Ledger
	LedgerType *LedgerType

	// CCF Property for the maximum size of the http request body: 1MB, 5MB, 10MB.
	MaxBodySizeInMb *int32

	// Number of CCF nodes in the ACC Ledger.
	NodeCount *int32

	// Object representing RunningState for Ledger.
	RunningState *RunningState

	// CCF Property for the subject name to include in the node certificate. Default: CN=CCF Node.
	SubjectName *string

	// Number of additional threads processing incoming client requests in the enclave (modify with care!)
	WorkerThreads *int32

	// Prefix for the write load balancer. Example: write
	WriteLBAddressPrefix *string

	// READ-ONLY; Endpoint for accessing network identity.
	IdentityServiceURI *string

	// READ-ONLY; Internal namespace for the Ledger
	LedgerInternalNamespace *string

	// READ-ONLY; Unique name for the Confidential Ledger.
	LedgerName *string

	// READ-ONLY; Endpoint for calling Ledger Service.
	LedgerURI *string

	// READ-ONLY; Provisioning state of Ledger Resource
	ProvisioningState *ProvisioningState
}

// List - Object that includes an array of Confidential Ledgers and a possible link for next set.
type List struct {
	// The URL the client should use to fetch the next page (per server side paging).
	NextLink *string

	// List of Confidential Ledgers
	Value []*ConfidentialLedger
}

// ManagedCCF - Managed CCF. Contains the properties of Managed CCF Resource.
type ManagedCCF struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// Properties of Managed CCF Resource.
	Properties *ManagedCCFProperties

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

// ManagedCCFBackup - Object representing Backup properties of a Managed CCF Resource.
type ManagedCCFBackup struct {
	// REQUIRED; SAS URI used to access the backup Fileshare.
	URI *string

	// The region where the backup of the managed CCF resource will eventually be restored to.
	RestoreRegion *string
}

// ManagedCCFBackupResponse - Object representing the backup response of a Managed CCF Resource.
type ManagedCCFBackupResponse struct {
	// READ-ONLY; Response body stating if the managed CCF resource is being backed up.
	Message *string
}

// ManagedCCFList - Object that includes an array of Managed CCF and a possible link for next set.
type ManagedCCFList struct {
	// The URL the client should use to fetch the next page (per server side paging).
	NextLink *string

	// List of Managed CCF
	Value []*ManagedCCF
}

// ManagedCCFProperties - Additional Managed CCF properties.
type ManagedCCFProperties struct {
	// Deployment Type of Managed CCF
	DeploymentType *DeploymentType

	// Enclave platform of Managed CCF.
	EnclavePlatform *EnclavePlatform

	// List of member identity certificates for Managed CCF
	MemberIdentityCertificates []*MemberIdentityCertificate

	// Number of CCF nodes in the Managed CCF.
	NodeCount *int32

	// Object representing RunningState for Managed CCF.
	RunningState *RunningState

	// READ-ONLY; Unique name for the Managed CCF.
	AppName *string

	// READ-ONLY; Endpoint for calling Managed CCF Service.
	AppURI *string

	// READ-ONLY; Endpoint for accessing network identity.
	IdentityServiceURI *string

	// READ-ONLY; Provisioning state of Managed CCF Resource
	ProvisioningState *ProvisioningState
}

// ManagedCCFRestore - Object representing Restore properties of Managed CCF Resource.
type ManagedCCFRestore struct {
	// REQUIRED; Fileshare where the managed CCF resource backup is stored.
	FileShareName *string

	// REQUIRED; The region the managed CCF resource is being restored to.
	RestoreRegion *string

	// REQUIRED; SAS URI used to access the backup Fileshare.
	URI *string
}

// ManagedCCFRestoreResponse - Object representing the restore response of a Managed CCF Resource.
type ManagedCCFRestoreResponse struct {
	// READ-ONLY; Response body stating if the managed CCF resource is being restored.
	Message *string
}

// MemberIdentityCertificate - Object representing MemberIdentityCertificate for Managed CCF.
type MemberIdentityCertificate struct {
	// Member Identity Certificate
	Certificate *string

	// Member Identity Certificate Encryption Key
	Encryptionkey *string

	// Anything
	Tags any
}

// ResourceProviderOperationDefinition - Describes the Resource Provider Operation.
type ResourceProviderOperationDefinition struct {
	// Details about the operations
	Display *ResourceProviderOperationDisplay

	// Indicates whether the operation is data action or not.
	IsDataAction *bool

	// Resource provider operation name.
	Name *string
}

// ResourceProviderOperationDisplay - Describes the properties of the Operation.
type ResourceProviderOperationDisplay struct {
	// Description of the resource provider operation.
	Description *string

	// Name of the resource provider operation.
	Operation *string

	// Name of the resource provider.
	Provider *string

	// Name of the resource type.
	Resource *string
}

// ResourceProviderOperationList - List containing this Resource Provider's available operations.
type ResourceProviderOperationList struct {
	// READ-ONLY; The URI that can be used to request the next page for list of Azure operations.
	NextLink *string

	// READ-ONLY; Resource provider operations list.
	Value []*ResourceProviderOperationDefinition
}

// Restore - Object representing Restore properties of a Confidential Ledger Resource.
type Restore struct {
	// REQUIRED; Fileshare where the ledger backup is stored.
	FileShareName *string

	// REQUIRED; The region the ledger is being restored to.
	RestoreRegion *string

	// REQUIRED; SAS URI used to access the backup fileshare.
	URI *string
}

// RestoreResponse - Object representing the restore response of a Confidential Ledger Resource.
type RestoreResponse struct {
	// READ-ONLY; Response body stating if the ledger is being restored.
	Message *string
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
