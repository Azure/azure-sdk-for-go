//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// CPKScopeInfo contains a group of parameters for the FilesystemClient.Create method.
type CPKScopeInfo = container.CPKScopeInfo

// AccessConditions identifies filesystem-specific access conditions which you optionally set.
type AccessConditions = container.AccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = container.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = container.ModifiedAccessConditions

// SignedIdentifier - signed identifier.
type SignedIdentifier = container.SignedIdentifier

// CreateOptions contains the optional parameters for the Filesystem.Create method.
type CreateOptions = container.CreateOptions

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions = container.DeleteOptions

// SetMetadataOptions contains the optional parameters for the Filesystem.SetMetadata method.
type SetMetadataOptions = container.SetMetadataOptions

// GetPropertiesOptions contains the optional parameters for the Filesystem.GetProperties method.
type GetPropertiesOptions = container.GetPropertiesOptions

// SetAccessPolicyOptions provides set of configurations for Filesystem.SetAccessPolicy operation.
type SetAccessPolicyOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access.
	// If this header is not included in the request, container data is private to the account owner.
	Access           *PublicAccessType
	AccessConditions *AccessConditions
	FilesystemACL    []*SignedIdentifier
}

// GetAccessPolicyOptions contains the optional parameters for the Filesystem.GetAccessPolicy method.
type GetAccessPolicyOptions = container.GetAccessPolicyOptions

type ListPathsOptions struct {
}

type ListDeletedPathsOptions struct {
}

type UndeletePathOptions struct {
}
