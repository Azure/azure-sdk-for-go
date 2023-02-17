//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
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
	return nil, nil
}

func (f *FileClient) generated() *generated.FileClient {
	return base.InnerClient((*base.Client[generated.FileClient])(f.fileClient))
}

// LeaseID returns leaseID of the client.
func (f *FileClient) LeaseID() *string {
	return f.leaseID
}

// Acquire operation can be used to request a new lease.
// The lease duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *FileClient) Acquire(ctx context.Context, duration int32, options *FileAcquireOptions) (FileAcquireResponse, error) {
	// TODO: update generated code to make duration as required parameter
	return FileAcquireResponse{}, nil
}

// Break operation can be used to break the lease, if the file has an active lease. Once a lease is broken, it cannot be renewed.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *FileClient) Break(ctx context.Context, options *FileBreakOptions) (FileBreakResponse, error) {
	return FileBreakResponse{}, nil
}

// Change operation can be used to change the lease ID of an active lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *FileClient) Change(ctx context.Context, leaseID string, proposedLeaseID string, options *FileChangeOptions) (FileChangeResponse, error) {
	// TODO: update generated code to make proposedLeaseID as required parameter
	return FileChangeResponse{}, nil
}

// Release operation can be used to free the lease if it is no longer needed so that another client may immediately acquire a lease against the file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *FileClient) Release(ctx context.Context, leaseID string, options *FileReleaseOptions) (FileReleaseResponse, error) {
	return FileReleaseResponse{}, nil
}
