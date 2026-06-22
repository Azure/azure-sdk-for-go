// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// ---------------------------------------------------------------------------------------------------------------------

// CreateOptions contains the optional parameters for the Client.Create method.
type CreateOptions struct {
	// Specifies the access tier of the share.
	AccessTier *AccessTier
	// Protocols to enable on the share.
	EnabledProtocols *string
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// Specifies the maximum size of the share, in gigabytes.
	Quota *int32
	// Root squash to set on the share. Only valid for NFS shares.
	RootSquash *RootSquash
	// Specifies whether the snapshot virtual directory should be accessible at the root of share mount point
	// when NFS is enabled.
	EnableSnapshotVirtualDirectoryAccess *bool

	// EnableSMBDirectoryLease contains the information returned from the x-ms-enable-smb-directory-lease header response.
	EnableSMBDirectoryLease *bool

	// Optional. Boolean. Default if not specified is false. This property enables paid bursting.
	PaidBurstingEnabled *bool

	// Optional. Integer. Default if not specified is the maximum throughput the file share can support. Current maximum for a
	// file share is 10,340 MiB/sec.
	PaidBurstingMaxBandwidthMibps *int64

	// Optional. Integer. Default if not specified is the maximum IOPS the file share can support. Current maximum for a file
	// share is 102,400 IOPS.
	PaidBurstingMaxIops *int64

	// Specifies the provisioned bandwidth of the share, in mebibytes per second (MiBps). If this is not
	// specified, the provisioned bandwidth is set to value calculated based on recommendation formula.
	ShareProvisionedBandwidthMibps *int64

	// Specifies the provisioned number of input/output operations per second (IOPS) of the share. If this is
	// not specified, the provisioned IOPS is set to value calculated based on recommendation formula.
	ShareProvisionedIops *int64
}

func (o *CreateOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientCreateOptions {
	opts := &generated.ShareClientCreateOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	opts.AccessTier = o.AccessTier
	opts.EnabledProtocols = o.EnabledProtocols
	opts.Metadata = o.Metadata
	opts.Quota = o.Quota
	opts.RootSquash = o.RootSquash
	opts.EnableSnapshotVirtualDirectoryAccess = o.EnableSnapshotVirtualDirectoryAccess
	opts.EnableSMBDirectoryLease = o.EnableSMBDirectoryLease
	opts.PaidBurstingEnabled = o.PaidBurstingEnabled
	opts.PaidBurstingMaxBandwidthMibps = o.PaidBurstingMaxBandwidthMibps
	opts.PaidBurstingMaxIops = o.PaidBurstingMaxIops
	opts.ShareProvisionedBandwidthMibps = o.ShareProvisionedBandwidthMibps
	opts.ShareProvisionedIops = o.ShareProvisionedIops
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	// Specifies the option include to delete the base share and all of its snapshots.
	DeleteSnapshots *DeleteSnapshotsOptionType
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *DeleteOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientDeleteOptions {
	opts := &generated.ShareClientDeleteOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	opts.DeleteSnapshots = o.DeleteSnapshots
	opts.Sharesnapshot = o.ShareSnapshot
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = generated.LeaseAccessConditions

// ---------------------------------------------------------------------------------------------------------------------

// RestoreOptions contains the optional parameters for the Client.Restore method.
type RestoreOptions struct {
	// placeholder for future options
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method.
type GetPropertiesOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetPropertiesOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientGetPropertiesOptions {
	opts := &generated.ShareClientGetPropertiesOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	opts.Sharesnapshot = o.ShareSnapshot
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// SetPropertiesOptions contains the optional parameters for the Client.SetProperties method.
type SetPropertiesOptions struct {
	// Specifies the access tier of the share.
	AccessTier *AccessTier
	// Specifies the maximum size of the share, in gigabytes.
	Quota *int32
	// Root squash to set on the share. Only valid for NFS shares.
	RootSquash *RootSquash
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
	// Specifies whether the snapshot virtual directory should be accessible at the root of share mount point
	// when NFS is enabled.
	EnableSnapshotVirtualDirectoryAccess *bool
	// EnableSMBDirectoryLease contains the information returned from the x-ms-enable-smb-directory-lease header response.
	EnableSMBDirectoryLease *bool
	// Optional. Boolean. Default if not specified is false. This property enables paid bursting.
	PaidBurstingEnabled *bool
	// Optional. Integer. Default if not specified is the maximum throughput the file share can support. Current maximum for a
	// file share is 10,340 MiB/sec.
	PaidBurstingMaxBandwidthMibps *int64
	// Optional. Integer. Default if not specified is the maximum IOPS the file share can support. Current maximum for a file
	// share is 102,400 IOPS.
	PaidBurstingMaxIops *int64
	// Specifies the provisioned bandwidth of the share, in mebibytes per second (MiBps). If this is not
	// specified, the provisioned bandwidth is set to value calculated based on recommendation formula.
	ShareProvisionedBandwidthMibps *int64

	// Specifies the provisioned number of input/output operations per second (IOPS) of the share. If this is
	// not specified, the provisioned IOPS is set to value calculated based on recommendation formula.
	ShareProvisionedIops *int64
}

func (o *SetPropertiesOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientSetPropertiesOptions {
	opts := &generated.ShareClientSetPropertiesOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	opts.AccessTier = o.AccessTier
	opts.Quota = o.Quota
	opts.RootSquash = o.RootSquash
	opts.EnableSnapshotVirtualDirectoryAccess = o.EnableSnapshotVirtualDirectoryAccess
	opts.EnableSMBDirectoryLease = o.EnableSMBDirectoryLease
	opts.PaidBurstingEnabled = o.PaidBurstingEnabled
	opts.PaidBurstingMaxBandwidthMibps = o.PaidBurstingMaxBandwidthMibps
	opts.PaidBurstingMaxIops = o.PaidBurstingMaxIops
	opts.ShareProvisionedIops = o.ShareProvisionedIops
	opts.ShareProvisionedBandwidthMibps = o.ShareProvisionedBandwidthMibps
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// CreateSnapshotOptions contains the optional parameters for the Client.CreateSnapshot method.
type CreateSnapshotOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
}

func (o *CreateSnapshotOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientCreateSnapshotOptions {
	opts := &generated.ShareClientCreateSnapshotOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	opts.Metadata = o.Metadata
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// GetAccessPolicyOptions contains the optional parameters for the Client.GetAccessPolicy method.
type GetAccessPolicyOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetAccessPolicyOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientGetAccessPolicyOptions {
	opts := &generated.ShareClientGetAccessPolicyOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

// SignedIdentifier - Signed identifier.
type SignedIdentifier = generated.SignedIdentifier

// AccessPolicy - An Access policy.
type AccessPolicy = generated.AccessPolicy

// AccessPolicyPermission type simplifies creating the permissions string for a share's access policy.
// Initialize an instance of this type and then call its String method to set AccessPolicy's permission field.
type AccessPolicyPermission = exported.AccessPolicyPermission

// ---------------------------------------------------------------------------------------------------------------------

// SetAccessPolicyOptions contains the optional parameters for the Client.SetAccessPolicy method.
type SetAccessPolicyOptions struct {
	// Specifies the ACL for the share.
	ShareACL []*SignedIdentifier
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetAccessPolicyOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientSetAccessPolicyOptions {
	opts := &generated.ShareClientSetAccessPolicyOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

func formatTime(si *SignedIdentifier) error {
	if si.AccessPolicy == nil {
		return nil
	}

	if si.AccessPolicy.Start != nil {
		st, err := time.Parse(time.RFC3339, si.AccessPolicy.Start.UTC().Format(time.RFC3339))
		if err != nil {
			return err
		}
		si.AccessPolicy.Start = &st
	}
	if si.AccessPolicy.Expiry != nil {
		et, err := time.Parse(time.RFC3339, si.AccessPolicy.Expiry.UTC().Format(time.RFC3339))
		if err != nil {
			return err
		}
		si.AccessPolicy.Expiry = &et
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// CreatePermissionOptions contains the optional parameters for the Client.CreatePermission method.
type CreatePermissionOptions struct {
	// placeholder for future options
}

func (o *CreatePermissionOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientCreatePermissionOptions {
	return &generated.ShareClientCreatePermissionOptions{
		FileRequestIntent: fileRequestIntent,
	}
}

// Permission - A permission (a security descriptor) at the share level.
type Permission = generated.SharePermission

// ---------------------------------------------------------------------------------------------------------------------

// GetPermissionOptions contains the optional parameters for the Client.GetPermission method.
type GetPermissionOptions struct {
	FilePermissionFormat *PermissionFormat
}

func (o *GetPermissionOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientGetPermissionOptions {
	opts := &generated.ShareClientGetPermissionOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	opts.FilePermissionFormat = o.FilePermissionFormat
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetMetadataOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientSetMetadataOptions {
	opts := &generated.ShareClientSetMetadataOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	opts.Metadata = o.Metadata
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// GetStatisticsOptions contains the optional parameters for the Client.GetStatistics method.
type GetStatisticsOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetStatisticsOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientGetStatisticsOptions {
	opts := &generated.ShareClientGetStatisticsOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}

	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

// Stats - Stats for the share.
type Stats = generated.ShareStats

// ---------------------------------------------------------------------------------------------------------------------

// GetSASURLOptions contains the optional parameters for the Client.GetSASURL method.
type GetSASURLOptions struct {
	StartTime *time.Time
}

func (o *GetSASURLOptions) format() time.Time {
	if o == nil {
		return time.Time{}
	}

	var st time.Time
	if o.StartTime != nil {
		st = o.StartTime.UTC()
	} else {
		st = time.Time{}
	}
	return st
}
