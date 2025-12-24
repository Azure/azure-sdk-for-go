// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"

// FileSystemAcquireResponse contains the response from method FileSystemClient.AcquireLease.
type FileSystemAcquireResponse = lease.ContainerAcquireResponse

// FileSystemBreakResponse contains the response from method FileSystemClient.BreakLease.
type FileSystemBreakResponse = lease.ContainerBreakResponse

// FileSystemChangeResponse contains the response from method FileSystemClient.ChangeLease.
type FileSystemChangeResponse = lease.ContainerChangeResponse

// FileSystemReleaseResponse contains the response from method FileSystemClient.ReleaseLease.
type FileSystemReleaseResponse = lease.ContainerReleaseResponse

// FileSystemRenewResponse contains the response from method FileSystemClient.RenewLease.
type FileSystemRenewResponse = lease.ContainerRenewResponse

// PathAcquireResponse contains the response from method PathClient.AcquireLease.
type PathAcquireResponse = lease.BlobAcquireResponse

// PathBreakResponse contains the response from method PathClient.BreakLease.
type PathBreakResponse = lease.BlobBreakResponse

// PathChangeResponse contains the response from method PathClient.ChangeLease.
type PathChangeResponse = lease.BlobChangeResponse

// PathReleaseResponse contains the response from method PathClient.ReleaseLease.
type PathReleaseResponse = lease.BlobReleaseResponse

// PathRenewResponse contains the response from method PathClient.RenewLease.
type PathRenewResponse = lease.BlobRenewResponse
