//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"

// FilesystemAcquireResponse contains the response from method BlobClient.AcquireLease.
type FilesystemAcquireResponse = lease.ContainerAcquireResponse

// FilesystemBreakResponse contains the response from method BlobClient.BreakLease.
type FilesystemBreakResponse = lease.ContainerBreakResponse

// FilesystemChangeResponse contains the response from method BlobClient.ChangeLease.
type FilesystemChangeResponse = lease.ContainerChangeResponse

// FilesystemReleaseResponse contains the response from method BlobClient.ReleaseLease.
type FilesystemReleaseResponse = lease.ContainerReleaseResponse

// FilesystemRenewResponse contains the response from method BlobClient.RenewLease.
type FilesystemRenewResponse = lease.ContainerRenewResponse

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
