// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
)

// ShareClient provides lease functionality for the underlying share client.
type ShareClient struct {
	shareClient *share.Client
	leaseID     *string
}

// ShareClientOptions contains the optional values when creating a ShareClient.
type ShareClientOptions struct {
	// LeaseID contains a caller-provided lease ID.
	LeaseID *string
}

// NewShareClient creates a share lease client for the provided share client.
//   - client - an instance of a share client
//   - options - client options; pass nil to accept the default values
func NewShareClient(client *share.Client, options *ShareClientOptions) (*ShareClient, error) {
	var leaseID *string
	if options != nil {
		leaseID = options.LeaseID
	}

	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
	}

	return &ShareClient{
		shareClient: client,
		leaseID:     leaseID,
	}, nil
}

func (s *ShareClient) generated() *generated.ShareClient {
	return base.InnerClient((*base.Client[generated.ShareClient])(s.shareClient))
}

// LeaseID returns leaseID of the client.
func (s *ShareClient) LeaseID() *string {
	return s.leaseID
}

// Acquire operation can be used to request a new lease.
// The lease duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) Acquire(ctx context.Context, duration int32, options *ShareAcquireOptions) (ShareAcquireResponse, error) {
	opts := options.format(s.LeaseID())
	resp, err := s.generated().AcquireLease(ctx, duration, opts)
	return resp, err
}

// Break operation can be used to break the lease, if the file share has an active lease. Once a lease is broken, it cannot be renewed.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) Break(ctx context.Context, options *ShareBreakOptions) (ShareBreakResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := s.generated().BreakLease(ctx, opts, leaseAccessConditions)
	return resp, err
}

// Change operation can be used to change the lease ID of an active lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) Change(ctx context.Context, proposedLeaseID string, options *ShareChangeOptions) (ShareChangeResponse, error) {
	if s.LeaseID() == nil {
		return ShareChangeResponse{}, errors.New("leaseID cannot be nil")
	}

	opts := options.format(&proposedLeaseID)
	resp, err := s.generated().ChangeLease(ctx, *s.LeaseID(), opts)

	// If lease has been changed successfully, set the leaseID in client
	if err == nil {
		s.leaseID = &proposedLeaseID
	}

	return resp, err
}

// Release operation can be used to free the lease if it is no longer needed so that another client may immediately acquire a lease against the file share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) Release(ctx context.Context, options *ShareReleaseOptions) (ShareReleaseResponse, error) {
	if s.LeaseID() == nil {
		return ShareReleaseResponse{}, errors.New("leaseID cannot be nil")
	}

	opts := options.format()
	resp, err := s.generated().ReleaseLease(ctx, *s.LeaseID(), opts)
	return resp, err
}

// Renew operation can be used to renew an existing lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) Renew(ctx context.Context, options *ShareRenewOptions) (ShareRenewResponse, error) {
	if s.LeaseID() == nil {
		return ShareRenewResponse{}, errors.New("leaseID cannot be nil")
	}

	opts := options.format()
	resp, err := s.generated().RenewLease(ctx, *s.LeaseID(), opts)
	return resp, err
}
