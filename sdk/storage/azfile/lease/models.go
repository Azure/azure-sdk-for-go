// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// AccessConditions contains optional parameters to access leased entity.
type AccessConditions = generated.LeaseAccessConditions

// FileAcquireOptions contains the optional parameters for the FileClient.Acquire method.
type FileAcquireOptions struct {
	// placeholder for future options
}

func (o *FileAcquireOptions) format(proposedLeaseID *string, fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientAcquireLeaseOptions {
	leaseDuration := int32(-1)
	opts := &generated.FileClientAcquireLeaseOptions{
		ProposedLeaseID:   proposedLeaseID,
		FileRequestIntent: fileRequestIntent,
		LeaseDuration:     &leaseDuration,
		AllowTrailingDot:  allowTrailingDot,
	}
	return opts
}

// FileBreakOptions contains the optional parameters for the FileClient.Break method.
type FileBreakOptions struct {
	// AccessConditions contains optional parameters to access leased entity.
	AccessConditions *AccessConditions
}

func (o *FileBreakOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientBreakLeaseOptions {
	opts := &generated.FileClientBreakLeaseOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	if o.AccessConditions != nil {
		opts.LeaseID = o.AccessConditions.LeaseID
	}
	return opts
}

// FileChangeOptions contains the optional parameters for the FileClient.Change method.
type FileChangeOptions struct {
	// placeholder for future options
}

func (o *FileChangeOptions) format(proposedLeaseID *string, fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientChangeLeaseOptions {
	opts := &generated.FileClientChangeLeaseOptions{
		ProposedLeaseID:   proposedLeaseID,
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	return opts
}

// FileReleaseOptions contains the optional parameters for the FileClient.Release method.
type FileReleaseOptions struct {
	// placeholder for future options
}

func (o *FileReleaseOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientReleaseLeaseOptions {
	opts := &generated.FileClientReleaseLeaseOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// ShareAcquireOptions contains the optional parameters for the ShareClient.Acquire method.
type ShareAcquireOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ShareAcquireOptions) format(proposedLeaseID *string, duration int32, fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientAcquireLeaseOptions {
	opts := &generated.ShareClientAcquireLeaseOptions{
		ProposedLeaseID:   proposedLeaseID,
		LeaseDuration:     &duration,
		FileRequestIntent: fileRequestIntent,
	}
	if o != nil {
		opts.Sharesnapshot = o.ShareSnapshot
	}
	return opts
}

// ShareBreakOptions contains the optional parameters for the ShareClient.Break method.
type ShareBreakOptions struct {
	// For a break operation, this is the proposed duration the lease should continue before it is broken, in seconds, between 0 and 60. This
	// break period is only used if it is shorter than the time remaining on the
	// lease. If longer, the time remaining on the lease is used. A new lease will not be available before the break period has
	// expired, but the lease may be held for longer than the break period. If this
	// header does not appear with a break operation, a fixed-duration lease breaks after the remaining lease period elapses,
	// and an infinite lease breaks immediately.
	BreakPeriod *int32
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
	// AccessConditions contains optional parameters to access leased entity.
	AccessConditions *AccessConditions
}

func (o *ShareBreakOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientBreakLeaseOptions {
	opts := &generated.ShareClientBreakLeaseOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o != nil {
		opts.BreakPeriod = o.BreakPeriod
		opts.Sharesnapshot = o.ShareSnapshot
		if o.AccessConditions != nil {
			opts.LeaseID = o.AccessConditions.LeaseID
		}
	}
	return opts
}

// ShareChangeOptions contains the optional parameters for the ShareClient.Change method.
type ShareChangeOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ShareChangeOptions) format(proposedLeaseID *string, fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientChangeLeaseOptions {
	opts := &generated.ShareClientChangeLeaseOptions{
		ProposedLeaseID:   proposedLeaseID,
		FileRequestIntent: fileRequestIntent,
	}
	if o != nil {
		opts.Sharesnapshot = o.ShareSnapshot
	}
	return opts
}

// ShareReleaseOptions contains the optional parameters for the ShareClient.Release method.
type ShareReleaseOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ShareReleaseOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientReleaseLeaseOptions {
	opts := &generated.ShareClientReleaseLeaseOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o != nil {
		opts.Sharesnapshot = o.ShareSnapshot
	}
	return opts
}

// ShareRenewOptions contains the optional parameters for the ShareClient.Renew method.
type ShareRenewOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ShareRenewOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.ShareClientRenewLeaseOptions {
	opts := &generated.ShareClientRenewLeaseOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o != nil {
		opts.Sharesnapshot = o.ShareSnapshot
	}
	return opts
}
