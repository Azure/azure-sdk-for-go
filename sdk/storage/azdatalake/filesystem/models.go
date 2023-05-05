//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
)

// CPKScopeInfo contains a group of parameters for the FilesystemClient.Create method.
type CPKScopeInfo = container.CPKScopeInfo

// AccessConditions identifies filesystem-specific access conditions which you optionally set.
type AccessConditions = exported.AccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions

// SignedIdentifier - signed identifier.
type SignedIdentifier = container.SignedIdentifier

type CreateOptions struct {
	// Specifies whether data in the filesystem may be accessed publicly and the level of access.
	Access *PublicAccessType

	// Optional. Specifies a user-defined name-value pair.
	Metadata map[string]*string

	// Optional. Specifies the encryption scope settings to set on the filesystem.
	CPKScopeInfo *CPKScopeInfo
}

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	AccessConditions *AccessConditions
}

// SetMetadataOptions contains the optional parameters for the Filesystem.SetMetadata method.
type SetMetadataOptions struct {
	Metadata                 map[string]*string
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

// GetPropertiesOptions contains the optional parameters for the Filesystem.GetProperties method.
type GetPropertiesOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

// SetAccessPolicyOptions provides set of configurations for Filesystem.SetAccessPolicy operation.
type SetAccessPolicyOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access.
	// If this header is not included in the request, container data is private to the account owner.
	Access           *PublicAccessType
	AccessConditions *AccessConditions
	ContainerACL     []*SignedIdentifier
}

// GetAccessPolicyOptions contains the optional parameters for the Filesystem.GetAccessPolicy method.
type GetAccessPolicyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

type ListPathsOptions struct {
}
