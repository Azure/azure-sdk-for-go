//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
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
	return nil, nil
}

func (s *ShareClient) generated() *generated.ShareClient {
	return base.InnerClient((*base.Client[generated.ShareClient])(s.shareClient))
}

// LeaseID returns leaseID of the client.
func (s *ShareClient) LeaseID() *string {
	return s.leaseID
}

// AcquireLease operation can be used to request a new lease.
// The lease duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) AcquireLease(ctx context.Context, duration int32, options *ShareAcquireOptions) (ShareAcquireResponse, error) {
	// TODO: update generated code to make duration as required parameter
	return ShareAcquireResponse{}, nil
}

// BreakLease operation can be used to break the lease, if the file share has an active lease. Once a lease is broken, it cannot be renewed.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) BreakLease(ctx context.Context, options *ShareBreakOptions) (ShareBreakResponse, error) {
	return ShareBreakResponse{}, nil
}

// ChangeLease operation can be used to change the lease ID of an active lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) ChangeLease(ctx context.Context, leaseID string, proposedLeaseID string, options *ShareChangeOptions) (ShareChangeResponse, error) {
	// TODO: update generated code to make proposedLeaseID as required parameter
	return ShareChangeResponse{}, nil
}

// ReleaseLease operation can be used to free the lease if it is no longer needed so that another client may immediately acquire a lease against the file share.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) ReleaseLease(ctx context.Context, leaseID string, options *ShareReleaseOptions) (ShareReleaseResponse, error) {
	return ShareReleaseResponse{}, nil
}

// RenewLease operation can be used to renew an existing lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-share.
func (s *ShareClient) RenewLease(ctx context.Context, leaseID string, options *ShareRenewOptions) (ShareRenewResponse, error) {
	return ShareRenewResponse{}, nil
}
