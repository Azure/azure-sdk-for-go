//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// AccessConditions contains optional parameters to access leased entity.
type AccessConditions = generated.LeaseAccessConditions

// FileAcquireOptions contains the optional parameters for the FileClient.Acquire method.
type FileAcquireOptions struct {
	// Proposed lease ID, in a GUID string format.
	// The File service returns 400 (Invalid request) if the proposed lease ID is not in the correct format.
	ProposedLeaseID *string
}

// FileBreakOptions contains the optional parameters for the FileClient.Break method.
type FileBreakOptions struct {
	// AccessConditions contains optional parameters to access leased entity.
	AccessConditions *AccessConditions
}

// FileChangeOptions contains the optional parameters for the FileClient.Change method.
type FileChangeOptions struct {
	// placeholder for future options
}

// FileReleaseOptions contains the optional parameters for the FileClient.Release method.
type FileReleaseOptions struct {
	// placeholder for future options
}

// ---------------------------------------------------------------------------------------------------------------------

// ShareAcquireOptions contains the optional parameters for the ShareClient.Acquire method.
type ShareAcquireOptions struct {
	// Proposed lease ID, in a GUID string format.
	// The File service returns 400 (Invalid request) if the proposed lease ID is not in the correct format.
	ProposedLeaseID *string
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
}

// ShareBreakOptions contains the optional parameters for the ShareClient.Break method.
type ShareBreakOptions struct {
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
	// AccessConditions contains optional parameters to access leased entity.
	AccessConditions *AccessConditions
}

// ShareChangeOptions contains the optional parameters for the ShareClient.Change method.
type ShareChangeOptions struct {
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
}

// ShareReleaseOptions contains the optional parameters for the ShareClient.Release method.
type ShareReleaseOptions struct {
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
}

// ShareRenewOptions contains the optional parameters for the ShareClient.Renew method.
type ShareRenewOptions struct {
	// TODO: Should snapshot be removed from the option bag
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	Snapshot *string
}
