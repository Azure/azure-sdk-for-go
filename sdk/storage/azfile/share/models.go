//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

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
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	// Specifies the option include to delete the base share and all of its snapshots.
	DeleteSnapshots *DeleteSnapshotsOptionType
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
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
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
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
}

// ---------------------------------------------------------------------------------------------------------------------

// CreateSnapshotOptions contains the optional parameters for the Client.CreateSnapshot method.
type CreateSnapshotOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
}

// ---------------------------------------------------------------------------------------------------------------------

// GetAccessPolicyOptions contains the optional parameters for the Client.GetAccessPolicy method.
type GetAccessPolicyOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

// SignedIdentifier - Signed identifier.
type SignedIdentifier = generated.SignedIdentifier

// AccessPolicy - An Access policy.
type AccessPolicy = generated.AccessPolicy

// ---------------------------------------------------------------------------------------------------------------------

// SetAccessPolicyOptions contains the optional parameters for the Client.SetAccessPolicy method.
type SetAccessPolicyOptions struct {
	// Specifies the ACL for the share.
	ShareACL []*SignedIdentifier
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// CreatePermissionOptions contains the optional parameters for the Client.CreatePermission method.
type CreatePermissionOptions struct {
	// placeholder for future options
}

// Permission - A permission (a security descriptor) at the share level.
type Permission = generated.SharePermission

// ---------------------------------------------------------------------------------------------------------------------

// GetPermissionOptions contains the optional parameters for the Client.GetPermission method.
type GetPermissionOptions struct {
	// placeholder for future options
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// GetStatisticsOptions contains the optional parameters for the Client.GetStatistics method.
type GetStatisticsOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

// Stats - Stats for the share.
type Stats = generated.ShareStats

// ---------------------------------------------------------------------------------------------------------------------

// AcquireLeaseOptions contains the optional parameters for the Client.AcquireLease method.
type AcquireLeaseOptions struct {
	// Specifies the duration of the lease, in seconds, or negative one (-1) for a lease that never expires. A non-infinite lease
	// can be between 15 and 60 seconds. A lease duration cannot be changed using
	// renew or change.
	Duration *int32
	// Proposed lease ID, in a GUID string format.
	// The File service returns 400 (Invalid request) if the proposed lease ID is not in the correct format.
	ProposedLeaseID *string
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
}

// ---------------------------------------------------------------------------------------------------------------------

// BreakLeaseOptions contains the optional parameters for the Client.BreakLease method.
type BreakLeaseOptions struct {
	// For a break operation, proposed duration the lease should continue before it is broken, in seconds, between 0 and 60. This
	// break period is only used if it is shorter than the time remaining on the
	// lease. If longer, the time remaining on the lease is used. A new lease will not be available before the break period has
	// expired, but the lease may be held for longer than the break period. If this
	// header does not appear with a break operation, a fixed-duration lease breaks after the remaining lease period elapses,
	// and an infinite lease breaks immediately.
	BreakPeriod *int32
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// ChangeLeaseOptions contains the optional parameters for the Client.ChangeLease method.
type ChangeLeaseOptions struct {
	// Proposed lease ID, in a GUID string format.
	// The File service returns 400 (Invalid request) if the proposed lease ID is not in the correct format.
	ProposedLeaseID *string
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
}

// ---------------------------------------------------------------------------------------------------------------------

// ReleaseLeaseOptions contains the optional parameters for the Client.ReleaseLease method.
type ReleaseLeaseOptions struct {
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
}

// ---------------------------------------------------------------------------------------------------------------------

// RenewLeaseOptions contains the optional parameters for the Client.RenewLease method.
type RenewLeaseOptions struct {
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
}
