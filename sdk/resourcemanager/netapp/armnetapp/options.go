//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetapp

// AccountsClientBeginChangeKeyVaultOptions contains the optional parameters for the AccountsClient.BeginChangeKeyVault method.
type AccountsClientBeginChangeKeyVaultOptions struct {
	// The required parameters to perform encryption migration.
	Body *ChangeKeyVault

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginCreateOrUpdateOptions contains the optional parameters for the AccountsClient.BeginCreateOrUpdate method.
type AccountsClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginDeleteOptions contains the optional parameters for the AccountsClient.BeginDelete method.
type AccountsClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginGetChangeKeyVaultInformationOptions contains the optional parameters for the AccountsClient.BeginGetChangeKeyVaultInformation
// method.
type AccountsClientBeginGetChangeKeyVaultInformationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginRenewCredentialsOptions contains the optional parameters for the AccountsClient.BeginRenewCredentials
// method.
type AccountsClientBeginRenewCredentialsOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginTransitionToCmkOptions contains the optional parameters for the AccountsClient.BeginTransitionToCmk
// method.
type AccountsClientBeginTransitionToCmkOptions struct {
	// The required parameters to perform encryption transition.
	Body *EncryptionTransitionRequest

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientBeginUpdateOptions contains the optional parameters for the AccountsClient.BeginUpdate method.
type AccountsClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// AccountsClientGetOptions contains the optional parameters for the AccountsClient.Get method.
type AccountsClientGetOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListBySubscriptionOptions contains the optional parameters for the AccountsClient.NewListBySubscriptionPager
// method.
type AccountsClientListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// AccountsClientListOptions contains the optional parameters for the AccountsClient.NewListPager method.
type AccountsClientListOptions struct {
	// placeholder for future optional parameters
}

// BackupPoliciesClientBeginCreateOptions contains the optional parameters for the BackupPoliciesClient.BeginCreate method.
type BackupPoliciesClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupPoliciesClientBeginDeleteOptions contains the optional parameters for the BackupPoliciesClient.BeginDelete method.
type BackupPoliciesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupPoliciesClientBeginUpdateOptions contains the optional parameters for the BackupPoliciesClient.BeginUpdate method.
type BackupPoliciesClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupPoliciesClientGetOptions contains the optional parameters for the BackupPoliciesClient.Get method.
type BackupPoliciesClientGetOptions struct {
	// placeholder for future optional parameters
}

// BackupPoliciesClientListOptions contains the optional parameters for the BackupPoliciesClient.NewListPager method.
type BackupPoliciesClientListOptions struct {
	// placeholder for future optional parameters
}

// BackupVaultsClientBeginCreateOrUpdateOptions contains the optional parameters for the BackupVaultsClient.BeginCreateOrUpdate
// method.
type BackupVaultsClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupVaultsClientBeginDeleteOptions contains the optional parameters for the BackupVaultsClient.BeginDelete method.
type BackupVaultsClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupVaultsClientBeginUpdateOptions contains the optional parameters for the BackupVaultsClient.BeginUpdate method.
type BackupVaultsClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupVaultsClientGetOptions contains the optional parameters for the BackupVaultsClient.Get method.
type BackupVaultsClientGetOptions struct {
	// placeholder for future optional parameters
}

// BackupVaultsClientListByNetAppAccountOptions contains the optional parameters for the BackupVaultsClient.NewListByNetAppAccountPager
// method.
type BackupVaultsClientListByNetAppAccountOptions struct {
	// placeholder for future optional parameters
}

// BackupsClientBeginCreateOptions contains the optional parameters for the BackupsClient.BeginCreate method.
type BackupsClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupsClientBeginDeleteOptions contains the optional parameters for the BackupsClient.BeginDelete method.
type BackupsClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupsClientBeginUpdateOptions contains the optional parameters for the BackupsClient.BeginUpdate method.
type BackupsClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupsClientGetLatestStatusOptions contains the optional parameters for the BackupsClient.GetLatestStatus method.
type BackupsClientGetLatestStatusOptions struct {
	// placeholder for future optional parameters
}

// BackupsClientGetOptions contains the optional parameters for the BackupsClient.Get method.
type BackupsClientGetOptions struct {
	// placeholder for future optional parameters
}

// BackupsClientGetVolumeLatestRestoreStatusOptions contains the optional parameters for the BackupsClient.GetVolumeLatestRestoreStatus
// method.
type BackupsClientGetVolumeLatestRestoreStatusOptions struct {
	// placeholder for future optional parameters
}

// BackupsClientListByVaultOptions contains the optional parameters for the BackupsClient.NewListByVaultPager method.
type BackupsClientListByVaultOptions struct {
	// An option to specify the VolumeResourceId. If present, then only returns the backups under the specified volume
	Filter *string
}

// BackupsUnderAccountClientBeginMigrateBackupsOptions contains the optional parameters for the BackupsUnderAccountClient.BeginMigrateBackups
// method.
type BackupsUnderAccountClientBeginMigrateBackupsOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupsUnderBackupVaultClientBeginRestoreFilesOptions contains the optional parameters for the BackupsUnderBackupVaultClient.BeginRestoreFiles
// method.
type BackupsUnderBackupVaultClientBeginRestoreFilesOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// BackupsUnderVolumeClientBeginMigrateBackupsOptions contains the optional parameters for the BackupsUnderVolumeClient.BeginMigrateBackups
// method.
type BackupsUnderVolumeClientBeginMigrateBackupsOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.NewListPager method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// PoolsClientBeginCreateOrUpdateOptions contains the optional parameters for the PoolsClient.BeginCreateOrUpdate method.
type PoolsClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// PoolsClientBeginDeleteOptions contains the optional parameters for the PoolsClient.BeginDelete method.
type PoolsClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// PoolsClientBeginUpdateOptions contains the optional parameters for the PoolsClient.BeginUpdate method.
type PoolsClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// PoolsClientGetOptions contains the optional parameters for the PoolsClient.Get method.
type PoolsClientGetOptions struct {
	// placeholder for future optional parameters
}

// PoolsClientListOptions contains the optional parameters for the PoolsClient.NewListPager method.
type PoolsClientListOptions struct {
	// placeholder for future optional parameters
}

// ResourceClientBeginUpdateNetworkSiblingSetOptions contains the optional parameters for the ResourceClient.BeginUpdateNetworkSiblingSet
// method.
type ResourceClientBeginUpdateNetworkSiblingSetOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ResourceClientCheckFilePathAvailabilityOptions contains the optional parameters for the ResourceClient.CheckFilePathAvailability
// method.
type ResourceClientCheckFilePathAvailabilityOptions struct {
	// placeholder for future optional parameters
}

// ResourceClientCheckNameAvailabilityOptions contains the optional parameters for the ResourceClient.CheckNameAvailability
// method.
type ResourceClientCheckNameAvailabilityOptions struct {
	// placeholder for future optional parameters
}

// ResourceClientCheckQuotaAvailabilityOptions contains the optional parameters for the ResourceClient.CheckQuotaAvailability
// method.
type ResourceClientCheckQuotaAvailabilityOptions struct {
	// placeholder for future optional parameters
}

// ResourceClientQueryNetworkSiblingSetOptions contains the optional parameters for the ResourceClient.QueryNetworkSiblingSet
// method.
type ResourceClientQueryNetworkSiblingSetOptions struct {
	// placeholder for future optional parameters
}

// ResourceClientQueryRegionInfoOptions contains the optional parameters for the ResourceClient.QueryRegionInfo method.
type ResourceClientQueryRegionInfoOptions struct {
	// placeholder for future optional parameters
}

// ResourceQuotaLimitsClientGetOptions contains the optional parameters for the ResourceQuotaLimitsClient.Get method.
type ResourceQuotaLimitsClientGetOptions struct {
	// placeholder for future optional parameters
}

// ResourceQuotaLimitsClientListOptions contains the optional parameters for the ResourceQuotaLimitsClient.NewListPager method.
type ResourceQuotaLimitsClientListOptions struct {
	// placeholder for future optional parameters
}

// ResourceRegionInfosClientGetOptions contains the optional parameters for the ResourceRegionInfosClient.Get method.
type ResourceRegionInfosClientGetOptions struct {
	// placeholder for future optional parameters
}

// ResourceRegionInfosClientListOptions contains the optional parameters for the ResourceRegionInfosClient.NewListPager method.
type ResourceRegionInfosClientListOptions struct {
	// placeholder for future optional parameters
}

// SnapshotPoliciesClientBeginDeleteOptions contains the optional parameters for the SnapshotPoliciesClient.BeginDelete method.
type SnapshotPoliciesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SnapshotPoliciesClientBeginUpdateOptions contains the optional parameters for the SnapshotPoliciesClient.BeginUpdate method.
type SnapshotPoliciesClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SnapshotPoliciesClientCreateOptions contains the optional parameters for the SnapshotPoliciesClient.Create method.
type SnapshotPoliciesClientCreateOptions struct {
	// placeholder for future optional parameters
}

// SnapshotPoliciesClientGetOptions contains the optional parameters for the SnapshotPoliciesClient.Get method.
type SnapshotPoliciesClientGetOptions struct {
	// placeholder for future optional parameters
}

// SnapshotPoliciesClientListOptions contains the optional parameters for the SnapshotPoliciesClient.NewListPager method.
type SnapshotPoliciesClientListOptions struct {
	// placeholder for future optional parameters
}

// SnapshotPoliciesClientListVolumesOptions contains the optional parameters for the SnapshotPoliciesClient.ListVolumes method.
type SnapshotPoliciesClientListVolumesOptions struct {
	// placeholder for future optional parameters
}

// SnapshotsClientBeginCreateOptions contains the optional parameters for the SnapshotsClient.BeginCreate method.
type SnapshotsClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SnapshotsClientBeginDeleteOptions contains the optional parameters for the SnapshotsClient.BeginDelete method.
type SnapshotsClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SnapshotsClientBeginRestoreFilesOptions contains the optional parameters for the SnapshotsClient.BeginRestoreFiles method.
type SnapshotsClientBeginRestoreFilesOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SnapshotsClientBeginUpdateOptions contains the optional parameters for the SnapshotsClient.BeginUpdate method.
type SnapshotsClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SnapshotsClientGetOptions contains the optional parameters for the SnapshotsClient.Get method.
type SnapshotsClientGetOptions struct {
	// placeholder for future optional parameters
}

// SnapshotsClientListOptions contains the optional parameters for the SnapshotsClient.NewListPager method.
type SnapshotsClientListOptions struct {
	// placeholder for future optional parameters
}

// SubvolumesClientBeginCreateOptions contains the optional parameters for the SubvolumesClient.BeginCreate method.
type SubvolumesClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SubvolumesClientBeginDeleteOptions contains the optional parameters for the SubvolumesClient.BeginDelete method.
type SubvolumesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SubvolumesClientBeginGetMetadataOptions contains the optional parameters for the SubvolumesClient.BeginGetMetadata method.
type SubvolumesClientBeginGetMetadataOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SubvolumesClientBeginUpdateOptions contains the optional parameters for the SubvolumesClient.BeginUpdate method.
type SubvolumesClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// SubvolumesClientGetOptions contains the optional parameters for the SubvolumesClient.Get method.
type SubvolumesClientGetOptions struct {
	// placeholder for future optional parameters
}

// SubvolumesClientListByVolumeOptions contains the optional parameters for the SubvolumesClient.NewListByVolumePager method.
type SubvolumesClientListByVolumeOptions struct {
	// placeholder for future optional parameters
}

// VolumeGroupsClientBeginCreateOptions contains the optional parameters for the VolumeGroupsClient.BeginCreate method.
type VolumeGroupsClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumeGroupsClientBeginDeleteOptions contains the optional parameters for the VolumeGroupsClient.BeginDelete method.
type VolumeGroupsClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumeGroupsClientGetOptions contains the optional parameters for the VolumeGroupsClient.Get method.
type VolumeGroupsClientGetOptions struct {
	// placeholder for future optional parameters
}

// VolumeGroupsClientListByNetAppAccountOptions contains the optional parameters for the VolumeGroupsClient.NewListByNetAppAccountPager
// method.
type VolumeGroupsClientListByNetAppAccountOptions struct {
	// placeholder for future optional parameters
}

// VolumeQuotaRulesClientBeginCreateOptions contains the optional parameters for the VolumeQuotaRulesClient.BeginCreate method.
type VolumeQuotaRulesClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumeQuotaRulesClientBeginDeleteOptions contains the optional parameters for the VolumeQuotaRulesClient.BeginDelete method.
type VolumeQuotaRulesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumeQuotaRulesClientBeginUpdateOptions contains the optional parameters for the VolumeQuotaRulesClient.BeginUpdate method.
type VolumeQuotaRulesClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumeQuotaRulesClientGetOptions contains the optional parameters for the VolumeQuotaRulesClient.Get method.
type VolumeQuotaRulesClientGetOptions struct {
	// placeholder for future optional parameters
}

// VolumeQuotaRulesClientListByVolumeOptions contains the optional parameters for the VolumeQuotaRulesClient.NewListByVolumePager
// method.
type VolumeQuotaRulesClientListByVolumeOptions struct {
	// placeholder for future optional parameters
}

// VolumesClientBeginAuthorizeExternalReplicationOptions contains the optional parameters for the VolumesClient.BeginAuthorizeExternalReplication
// method.
type VolumesClientBeginAuthorizeExternalReplicationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginAuthorizeReplicationOptions contains the optional parameters for the VolumesClient.BeginAuthorizeReplication
// method.
type VolumesClientBeginAuthorizeReplicationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginBreakFileLocksOptions contains the optional parameters for the VolumesClient.BeginBreakFileLocks method.
type VolumesClientBeginBreakFileLocksOptions struct {
	// Optional body to provide the ability to clear file locks with selected options
	Body *BreakFileLocksRequest

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginBreakReplicationOptions contains the optional parameters for the VolumesClient.BeginBreakReplication
// method.
type VolumesClientBeginBreakReplicationOptions struct {
	// Optional body to force break the replication.
	Body *BreakReplicationRequest

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginCreateOrUpdateOptions contains the optional parameters for the VolumesClient.BeginCreateOrUpdate method.
type VolumesClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginDeleteOptions contains the optional parameters for the VolumesClient.BeginDelete method.
type VolumesClientBeginDeleteOptions struct {
	// An option to force delete the volume. Will cleanup resources connected to the particular volume
	ForceDelete *bool

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginDeleteReplicationOptions contains the optional parameters for the VolumesClient.BeginDeleteReplication
// method.
type VolumesClientBeginDeleteReplicationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginFinalizeExternalReplicationOptions contains the optional parameters for the VolumesClient.BeginFinalizeExternalReplication
// method.
type VolumesClientBeginFinalizeExternalReplicationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginFinalizeRelocationOptions contains the optional parameters for the VolumesClient.BeginFinalizeRelocation
// method.
type VolumesClientBeginFinalizeRelocationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginListGetGroupIDListForLdapUserOptions contains the optional parameters for the VolumesClient.BeginListGetGroupIDListForLdapUser
// method.
type VolumesClientBeginListGetGroupIDListForLdapUserOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginPeerExternalClusterOptions contains the optional parameters for the VolumesClient.BeginPeerExternalCluster
// method.
type VolumesClientBeginPeerExternalClusterOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginPerformReplicationTransferOptions contains the optional parameters for the VolumesClient.BeginPerformReplicationTransfer
// method.
type VolumesClientBeginPerformReplicationTransferOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginPoolChangeOptions contains the optional parameters for the VolumesClient.BeginPoolChange method.
type VolumesClientBeginPoolChangeOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginPopulateAvailabilityZoneOptions contains the optional parameters for the VolumesClient.BeginPopulateAvailabilityZone
// method.
type VolumesClientBeginPopulateAvailabilityZoneOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginReInitializeReplicationOptions contains the optional parameters for the VolumesClient.BeginReInitializeReplication
// method.
type VolumesClientBeginReInitializeReplicationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginReestablishReplicationOptions contains the optional parameters for the VolumesClient.BeginReestablishReplication
// method.
type VolumesClientBeginReestablishReplicationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginRelocateOptions contains the optional parameters for the VolumesClient.BeginRelocate method.
type VolumesClientBeginRelocateOptions struct {
	// Relocate volume request
	Body *RelocateVolumeRequest

	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginResetCifsPasswordOptions contains the optional parameters for the VolumesClient.BeginResetCifsPassword
// method.
type VolumesClientBeginResetCifsPasswordOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginResyncReplicationOptions contains the optional parameters for the VolumesClient.BeginResyncReplication
// method.
type VolumesClientBeginResyncReplicationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginRevertOptions contains the optional parameters for the VolumesClient.BeginRevert method.
type VolumesClientBeginRevertOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginRevertRelocationOptions contains the optional parameters for the VolumesClient.BeginRevertRelocation
// method.
type VolumesClientBeginRevertRelocationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientBeginUpdateOptions contains the optional parameters for the VolumesClient.BeginUpdate method.
type VolumesClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// VolumesClientGetOptions contains the optional parameters for the VolumesClient.Get method.
type VolumesClientGetOptions struct {
	// placeholder for future optional parameters
}

// VolumesClientListOptions contains the optional parameters for the VolumesClient.NewListPager method.
type VolumesClientListOptions struct {
	// placeholder for future optional parameters
}

// VolumesClientListReplicationsOptions contains the optional parameters for the VolumesClient.NewListReplicationsPager method.
type VolumesClientListReplicationsOptions struct {
	// placeholder for future optional parameters
}

// VolumesClientReplicationStatusOptions contains the optional parameters for the VolumesClient.ReplicationStatus method.
type VolumesClientReplicationStatusOptions struct {
	// placeholder for future optional parameters
}
