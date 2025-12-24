// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
)

// FileClient provides lease functionality for the underlying file client.
type FileClient struct {
	fileClient *file.Client
	leaseID    *string
}

// FileClientOptions contains the optional values when creating a FileClient.
type FileClientOptions struct {
	// LeaseID contains a caller-provided lease ID.
	LeaseID *string
}

// NewFileClient creates a file lease client for the provided file client.
//   - client - an instance of a file client
//   - options - client options; pass nil to accept the default values
func NewFileClient(client *file.Client, options *FileClientOptions) (*FileClient, error) {
	var leaseID *string
	if options != nil {
		leaseID = options.LeaseID
	}

	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
	}

	return &FileClient{
		fileClient: client,
		leaseID:    leaseID,
	}, nil
}

func (f *FileClient) generated() *generated.FileClient {
	return base.InnerClient((*base.Client[generated.FileClient])(f.fileClient))
}

// LeaseID returns leaseID of the client.
func (f *FileClient) LeaseID() *string {
	return f.leaseID
}

// Acquire operation can be used to request a new lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *FileClient) Acquire(ctx context.Context, options *FileAcquireOptions) (FileAcquireResponse, error) {
	opts := options.format(f.LeaseID())
	resp, err := f.generated().AcquireLease(ctx, (int32)(-1), opts)
	return resp, err
}

// Break operation can be used to break the lease, if the file has an active lease. Once a lease is broken, it cannot be renewed.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *FileClient) Break(ctx context.Context, options *FileBreakOptions) (FileBreakResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := f.generated().BreakLease(ctx, opts, leaseAccessConditions)
	return resp, err
}

// Change operation can be used to change the lease ID of an active lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *FileClient) Change(ctx context.Context, proposedLeaseID string, options *FileChangeOptions) (FileChangeResponse, error) {
	if f.LeaseID() == nil {
		return FileChangeResponse{}, errors.New("leaseID cannot be nil")
	}

	opts := options.format(&proposedLeaseID)
	resp, err := f.generated().ChangeLease(ctx, *f.LeaseID(), opts)

	// If lease has been changed successfully, set the leaseID in client
	if err == nil {
		f.leaseID = &proposedLeaseID
	}

	return resp, err
}

// Release operation can be used to free the lease if it is no longer needed so that another client may immediately acquire a lease against the file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *FileClient) Release(ctx context.Context, options *FileReleaseOptions) (FileReleaseResponse, error) {
	if f.LeaseID() == nil {
		return FileReleaseResponse{}, errors.New("leaseID cannot be nil")
	}

	opts := options.format()
	resp, err := f.generated().ReleaseLease(ctx, *f.LeaseID(), opts)
	return resp, err
}
