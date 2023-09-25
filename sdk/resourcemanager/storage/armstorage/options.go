//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstorage

// AccountsClientBeginAbortHierarchicalNamespaceMigrationOptions contains the optional parameters for the AccountsClient.BeginAbortHierarchicalNamespaceMigration
// method.
type AccountsClientBeginAbortHierarchicalNamespaceMigrationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginCreateOptions contains the optional parameters for the AccountsClient.BeginCreate method.
type AccountsClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginCustomerInitiatedMigrationOptions contains the optional parameters for the AccountsClient.BeginCustomerInitiatedMigration
// method.
type AccountsClientBeginCustomerInitiatedMigrationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginFailoverOptions contains the optional parameters for the AccountsClient.BeginFailover method.
type AccountsClientBeginFailoverOptions struct {
	// The parameter is set to 'Planned' to indicate whether a Planned failover is requested.. Specifying any value will set the
// value to Planned.
	FailoverType *string

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginHierarchicalNamespaceMigrationOptions contains the optional parameters for the AccountsClient.BeginHierarchicalNamespaceMigration
// method.
type AccountsClientBeginHierarchicalNamespaceMigrationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginRestoreBlobRangesOptions contains the optional parameters for the AccountsClient.BeginRestoreBlobRanges
// method.
type AccountsClientBeginRestoreBlobRangesOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientCheckNameAvailabilityOptions contains the optional parameters for the AccountsClient.CheckNameAvailability
// method.
type AccountsClientCheckNameAvailabilityOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientDeleteOptions contains the optional parameters for the AccountsClient.Delete method.
type AccountsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientGetCustomerInitiatedMigrationOptions contains the optional parameters for the AccountsClient.GetCustomerInitiatedMigration
// method.
type AccountsClientGetCustomerInitiatedMigrationOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientGetPropertiesOptions contains the optional parameters for the AccountsClient.GetProperties method.
type AccountsClientGetPropertiesOptions struct {
	// May be used to expand the properties within account's properties. By default, data is not included when fetching properties.
// Currently we only support geoReplicationStats and blobRestoreStatus.
	Expand *StorageAccountExpand
}

// AccountsClientListAccountSASOptions contains the optional parameters for the AccountsClient.ListAccountSAS method.
type AccountsClientListAccountSASOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListByResourceGroupOptions contains the optional parameters for the AccountsClient.NewListByResourceGroupPager
// method.
type AccountsClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListKeysOptions contains the optional parameters for the AccountsClient.ListKeys method.
type AccountsClientListKeysOptions struct {
	// Specifies type of the key to be listed. Possible value is kerb.. Specifying any value will set the value to kerb.
	Expand *string
}

// AccountsClientListOptions contains the optional parameters for the AccountsClient.NewListPager method.
type AccountsClientListOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListServiceSASOptions contains the optional parameters for the AccountsClient.ListServiceSAS method.
type AccountsClientListServiceSASOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientRegenerateKeyOptions contains the optional parameters for the AccountsClient.RegenerateKey method.
type AccountsClientRegenerateKeyOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientRevokeUserDelegationKeysOptions contains the optional parameters for the AccountsClient.RevokeUserDelegationKeys
// method.
type AccountsClientRevokeUserDelegationKeysOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientUpdateOptions contains the optional parameters for the AccountsClient.Update method.
type AccountsClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// BlobContainersClientBeginObjectLevelWormOptions contains the optional parameters for the BlobContainersClient.BeginObjectLevelWorm
// method.
type BlobContainersClientBeginObjectLevelWormOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BlobContainersClientClearLegalHoldOptions contains the optional parameters for the BlobContainersClient.ClearLegalHold
// method.
type BlobContainersClientClearLegalHoldOptions struct {
	// placeholder for future optional parameters
}

// BlobContainersClientCreateOptions contains the optional parameters for the BlobContainersClient.Create method.
type BlobContainersClientCreateOptions struct {
	// placeholder for future optional parameters
}

// BlobContainersClientCreateOrUpdateImmutabilityPolicyOptions contains the optional parameters for the BlobContainersClient.CreateOrUpdateImmutabilityPolicy
// method.
type BlobContainersClientCreateOrUpdateImmutabilityPolicyOptions struct {
	// The entity state (ETag) version of the immutability policy to update. A value of "*" can be used to apply the operation
// only if the immutability policy already exists. If omitted, this operation will
// always be applied.
	IfMatch *string

	// The ImmutabilityPolicy Properties that will be created or updated to a blob container.
	Parameters *ImmutabilityPolicy
}

// BlobContainersClientDeleteImmutabilityPolicyOptions contains the optional parameters for the BlobContainersClient.DeleteImmutabilityPolicy
// method.
type BlobContainersClientDeleteImmutabilityPolicyOptions struct {
	// placeholder for future optional parameters
}

// BlobContainersClientDeleteOptions contains the optional parameters for the BlobContainersClient.Delete method.
type BlobContainersClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// BlobContainersClientExtendImmutabilityPolicyOptions contains the optional parameters for the BlobContainersClient.ExtendImmutabilityPolicy
// method.
type BlobContainersClientExtendImmutabilityPolicyOptions struct {
	// The ImmutabilityPolicy Properties that will be extended for a blob container.
	Parameters *ImmutabilityPolicy
}

// BlobContainersClientGetImmutabilityPolicyOptions contains the optional parameters for the BlobContainersClient.GetImmutabilityPolicy
// method.
type BlobContainersClientGetImmutabilityPolicyOptions struct {
	// The entity state (ETag) version of the immutability policy to update. A value of "*" can be used to apply the operation
// only if the immutability policy already exists. If omitted, this operation will
// always be applied.
	IfMatch *string
}

// BlobContainersClientGetOptions contains the optional parameters for the BlobContainersClient.Get method.
type BlobContainersClientGetOptions struct {
	// placeholder for future optional parameters
}

// BlobContainersClientLeaseOptions contains the optional parameters for the BlobContainersClient.Lease method.
type BlobContainersClientLeaseOptions struct {
	// Lease Container request body.
	Parameters *LeaseContainerRequest
}

// BlobContainersClientListOptions contains the optional parameters for the BlobContainersClient.NewListPager method.
type BlobContainersClientListOptions struct {
	// Optional. When specified, only container names starting with the filter will be listed.
	Filter *string

	// Optional, used to include the properties for soft deleted blob containers.
	Include *ListContainersInclude

	// Optional. Specified maximum number of containers that can be included in the list.
	Maxpagesize *string
}

// BlobContainersClientLockImmutabilityPolicyOptions contains the optional parameters for the BlobContainersClient.LockImmutabilityPolicy
// method.
type BlobContainersClientLockImmutabilityPolicyOptions struct {
	// placeholder for future optional parameters
}

// BlobContainersClientSetLegalHoldOptions contains the optional parameters for the BlobContainersClient.SetLegalHold method.
type BlobContainersClientSetLegalHoldOptions struct {
	// placeholder for future optional parameters
}

// BlobContainersClientUpdateOptions contains the optional parameters for the BlobContainersClient.Update method.
type BlobContainersClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// BlobInventoryPoliciesClientCreateOrUpdateOptions contains the optional parameters for the BlobInventoryPoliciesClient.CreateOrUpdate
// method.
type BlobInventoryPoliciesClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// BlobInventoryPoliciesClientDeleteOptions contains the optional parameters for the BlobInventoryPoliciesClient.Delete method.
type BlobInventoryPoliciesClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// BlobInventoryPoliciesClientGetOptions contains the optional parameters for the BlobInventoryPoliciesClient.Get method.
type BlobInventoryPoliciesClientGetOptions struct {
	// placeholder for future optional parameters
}

// BlobInventoryPoliciesClientListOptions contains the optional parameters for the BlobInventoryPoliciesClient.NewListPager
// method.
type BlobInventoryPoliciesClientListOptions struct {
	// placeholder for future optional parameters
}

// BlobServicesClientGetServicePropertiesOptions contains the optional parameters for the BlobServicesClient.GetServiceProperties
// method.
type BlobServicesClientGetServicePropertiesOptions struct {
	// placeholder for future optional parameters
}

// BlobServicesClientListOptions contains the optional parameters for the BlobServicesClient.NewListPager method.
type BlobServicesClientListOptions struct {
	// placeholder for future optional parameters
}

// BlobServicesClientSetServicePropertiesOptions contains the optional parameters for the BlobServicesClient.SetServiceProperties
// method.
type BlobServicesClientSetServicePropertiesOptions struct {
	// placeholder for future optional parameters
}

// DeletedAccountsClientGetOptions contains the optional parameters for the DeletedAccountsClient.Get method.
type DeletedAccountsClientGetOptions struct {
	// placeholder for future optional parameters
}

// DeletedAccountsClientListOptions contains the optional parameters for the DeletedAccountsClient.NewListPager method.
type DeletedAccountsClientListOptions struct {
	// placeholder for future optional parameters
}

// EncryptionScopesClientGetOptions contains the optional parameters for the EncryptionScopesClient.Get method.
type EncryptionScopesClientGetOptions struct {
	// placeholder for future optional parameters
}

// EncryptionScopesClientListOptions contains the optional parameters for the EncryptionScopesClient.NewListPager method.
type EncryptionScopesClientListOptions struct {
	// Optional. When specified, only encryption scope names starting with the filter will be listed.
	Filter *string

	// Optional, when specified, will list encryption scopes with the specific state. Defaults to All
	Include *ListEncryptionScopesInclude

	// Optional, specifies the maximum number of encryption scopes that will be included in the list response.
	Maxpagesize *int32
}

// EncryptionScopesClientPatchOptions contains the optional parameters for the EncryptionScopesClient.Patch method.
type EncryptionScopesClientPatchOptions struct {
	// placeholder for future optional parameters
}

// EncryptionScopesClientPutOptions contains the optional parameters for the EncryptionScopesClient.Put method.
type EncryptionScopesClientPutOptions struct {
	// placeholder for future optional parameters
}

// FileServicesClientGetServicePropertiesOptions contains the optional parameters for the FileServicesClient.GetServiceProperties
// method.
type FileServicesClientGetServicePropertiesOptions struct {
	// placeholder for future optional parameters
}

// FileServicesClientListOptions contains the optional parameters for the FileServicesClient.List method.
type FileServicesClientListOptions struct {
	// placeholder for future optional parameters
}

// FileServicesClientSetServicePropertiesOptions contains the optional parameters for the FileServicesClient.SetServiceProperties
// method.
type FileServicesClientSetServicePropertiesOptions struct {
	// placeholder for future optional parameters
}

// FileSharesClientCreateOptions contains the optional parameters for the FileSharesClient.Create method.
type FileSharesClientCreateOptions struct {
	// Optional, used to expand the properties within share's properties. Valid values are: snapshots. Should be passed as a string
// with delimiter ','
	Expand *string
}

// FileSharesClientDeleteOptions contains the optional parameters for the FileSharesClient.Delete method.
type FileSharesClientDeleteOptions struct {
	// Optional. Valid values are: snapshots, leased-snapshots, none. The default value is snapshots. For 'snapshots', the file
// share is deleted including all of its file share snapshots. If the file share
// contains leased-snapshots, the deletion fails. For 'leased-snapshots', the file share is deleted included all of its file
// share snapshots (leased/unleased). For 'none', the file share is deleted if it
// has no share snapshots. If the file share contains any snapshots (leased or unleased), the deletion fails.
	Include *string

	// Optional, used to delete a snapshot.
	XMSSnapshot *string
}

// FileSharesClientGetOptions contains the optional parameters for the FileSharesClient.Get method.
type FileSharesClientGetOptions struct {
	// Optional, used to expand the properties within share's properties. Valid values are: stats. Should be passed as a string
// with delimiter ','.
	Expand *string

	// Optional, used to retrieve properties of a snapshot.
	XMSSnapshot *string
}

// FileSharesClientLeaseOptions contains the optional parameters for the FileSharesClient.Lease method.
type FileSharesClientLeaseOptions struct {
	// Lease Share request body.
	Parameters *LeaseShareRequest

	// Optional. Specify the snapshot time to lease a snapshot.
	XMSSnapshot *string
}

// FileSharesClientListOptions contains the optional parameters for the FileSharesClient.NewListPager method.
type FileSharesClientListOptions struct {
	// Optional, used to expand the properties within share's properties. Valid values are: deleted, snapshots. Should be passed
// as a string with delimiter ','
	Expand *string

	// Optional. When specified, only share names starting with the filter will be listed.
	Filter *string

	// Optional. Specified maximum number of shares that can be included in the list.
	Maxpagesize *string
}

// FileSharesClientRestoreOptions contains the optional parameters for the FileSharesClient.Restore method.
type FileSharesClientRestoreOptions struct {
	// placeholder for future optional parameters
}

// FileSharesClientUpdateOptions contains the optional parameters for the FileSharesClient.Update method.
type FileSharesClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// LocalUsersClientCreateOrUpdateOptions contains the optional parameters for the LocalUsersClient.CreateOrUpdate method.
type LocalUsersClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// LocalUsersClientDeleteOptions contains the optional parameters for the LocalUsersClient.Delete method.
type LocalUsersClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// LocalUsersClientGetOptions contains the optional parameters for the LocalUsersClient.Get method.
type LocalUsersClientGetOptions struct {
	// placeholder for future optional parameters
}

// LocalUsersClientListKeysOptions contains the optional parameters for the LocalUsersClient.ListKeys method.
type LocalUsersClientListKeysOptions struct {
	// placeholder for future optional parameters
}

// LocalUsersClientListOptions contains the optional parameters for the LocalUsersClient.NewListPager method.
type LocalUsersClientListOptions struct {
	// placeholder for future optional parameters
}

// LocalUsersClientRegeneratePasswordOptions contains the optional parameters for the LocalUsersClient.RegeneratePassword
// method.
type LocalUsersClientRegeneratePasswordOptions struct {
	// placeholder for future optional parameters
}

// ManagementPoliciesClientCreateOrUpdateOptions contains the optional parameters for the ManagementPoliciesClient.CreateOrUpdate
// method.
type ManagementPoliciesClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// ManagementPoliciesClientDeleteOptions contains the optional parameters for the ManagementPoliciesClient.Delete method.
type ManagementPoliciesClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// ManagementPoliciesClientGetOptions contains the optional parameters for the ManagementPoliciesClient.Get method.
type ManagementPoliciesClientGetOptions struct {
	// placeholder for future optional parameters
}

// ObjectReplicationPoliciesClientCreateOrUpdateOptions contains the optional parameters for the ObjectReplicationPoliciesClient.CreateOrUpdate
// method.
type ObjectReplicationPoliciesClientCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// ObjectReplicationPoliciesClientDeleteOptions contains the optional parameters for the ObjectReplicationPoliciesClient.Delete
// method.
type ObjectReplicationPoliciesClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// ObjectReplicationPoliciesClientGetOptions contains the optional parameters for the ObjectReplicationPoliciesClient.Get
// method.
type ObjectReplicationPoliciesClientGetOptions struct {
	// placeholder for future optional parameters
}

// ObjectReplicationPoliciesClientListOptions contains the optional parameters for the ObjectReplicationPoliciesClient.NewListPager
// method.
type ObjectReplicationPoliciesClientListOptions struct {
	// placeholder for future optional parameters
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.NewListPager method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsClientDeleteOptions contains the optional parameters for the PrivateEndpointConnectionsClient.Delete
// method.
type PrivateEndpointConnectionsClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsClientGetOptions contains the optional parameters for the PrivateEndpointConnectionsClient.Get
// method.
type PrivateEndpointConnectionsClientGetOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsClientListOptions contains the optional parameters for the PrivateEndpointConnectionsClient.NewListPager
// method.
type PrivateEndpointConnectionsClientListOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsClientPutOptions contains the optional parameters for the PrivateEndpointConnectionsClient.Put
// method.
type PrivateEndpointConnectionsClientPutOptions struct {
	// placeholder for future optional parameters
}

// PrivateLinkResourcesClientListByStorageAccountOptions contains the optional parameters for the PrivateLinkResourcesClient.ListByStorageAccount
// method.
type PrivateLinkResourcesClientListByStorageAccountOptions struct {
	// placeholder for future optional parameters
}

// QueueClientCreateOptions contains the optional parameters for the QueueClient.Create method.
type QueueClientCreateOptions struct {
	// placeholder for future optional parameters
}

// QueueClientDeleteOptions contains the optional parameters for the QueueClient.Delete method.
type QueueClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// QueueClientGetOptions contains the optional parameters for the QueueClient.Get method.
type QueueClientGetOptions struct {
	// placeholder for future optional parameters
}

// QueueClientListOptions contains the optional parameters for the QueueClient.NewListPager method.
type QueueClientListOptions struct {
	// Optional, When specified, only the queues with a name starting with the given filter will be listed.
	Filter *string

	// Optional, a maximum number of queues that should be included in a list queue response
	Maxpagesize *string
}

// QueueClientUpdateOptions contains the optional parameters for the QueueClient.Update method.
type QueueClientUpdateOptions struct {
	// placeholder for future optional parameters
}

// QueueServicesClientGetServicePropertiesOptions contains the optional parameters for the QueueServicesClient.GetServiceProperties
// method.
type QueueServicesClientGetServicePropertiesOptions struct {
	// placeholder for future optional parameters
}

// QueueServicesClientListOptions contains the optional parameters for the QueueServicesClient.List method.
type QueueServicesClientListOptions struct {
	// placeholder for future optional parameters
}

// QueueServicesClientSetServicePropertiesOptions contains the optional parameters for the QueueServicesClient.SetServiceProperties
// method.
type QueueServicesClientSetServicePropertiesOptions struct {
	// placeholder for future optional parameters
}

// SKUsClientListOptions contains the optional parameters for the SKUsClient.NewListPager method.
type SKUsClientListOptions struct {
	// placeholder for future optional parameters
}

// TableClientCreateOptions contains the optional parameters for the TableClient.Create method.
type TableClientCreateOptions struct {
	// The parameters to provide to create a table.
	Parameters *Table
}

// TableClientDeleteOptions contains the optional parameters for the TableClient.Delete method.
type TableClientDeleteOptions struct {
	// placeholder for future optional parameters
}

// TableClientGetOptions contains the optional parameters for the TableClient.Get method.
type TableClientGetOptions struct {
	// placeholder for future optional parameters
}

// TableClientListOptions contains the optional parameters for the TableClient.NewListPager method.
type TableClientListOptions struct {
	// placeholder for future optional parameters
}

// TableClientUpdateOptions contains the optional parameters for the TableClient.Update method.
type TableClientUpdateOptions struct {
	// The parameters to provide to create a table.
	Parameters *Table
}

// TableServicesClientGetServicePropertiesOptions contains the optional parameters for the TableServicesClient.GetServiceProperties
// method.
type TableServicesClientGetServicePropertiesOptions struct {
	// placeholder for future optional parameters
}

// TableServicesClientListOptions contains the optional parameters for the TableServicesClient.List method.
type TableServicesClientListOptions struct {
	// placeholder for future optional parameters
}

// TableServicesClientSetServicePropertiesOptions contains the optional parameters for the TableServicesClient.SetServiceProperties
// method.
type TableServicesClientSetServicePropertiesOptions struct {
	// placeholder for future optional parameters
}

// UsagesClientListByLocationOptions contains the optional parameters for the UsagesClient.NewListByLocationPager method.
type UsagesClientListByLocationOptions struct {
	// placeholder for future optional parameters
}

