//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
	"time"
)

// CreateOptions contains the optional parameters when calling the Create operation. dfs endpoint
type CreateOptions struct {
	// AccessConditions specifies parameters for accessing the directory
	AccessConditions *AccessConditions
	// Metadata defines user passed key-value pairs that are associated with the directory.
	Metadata map[string]*string
	// CPKInfo contains a group of parameters for client provided encryption key.
	CPKInfo *CPKInfo
	// HTTPHeaders contains the HTTP headers for path operations.
	HTTPHeaders *HTTPHeaders
	// LeaseDuration defines the duration of the lease.
	LeaseDuration *time.Duration
	// ProposedLeaseID defines the proposed lease ID for the directory.
	ProposedLeaseID *string
	// Permissions defines the file permission for the directory.
	Permissions *string
	// Umask defines the umask permission of the directory.
	Umask *string
	// Owner defines the owner of the directory.
	Owner *string
	// Group defines the group of the directory.
	Group *string
	// ACL defines the ACL of the directory.
	ACL *string
}

// DeleteOptions contains the optional parameters when calling the Delete operation. dfs endpoint
type DeleteOptions struct {
	// AccessConditions specifies parameters for accessing the directory
	AccessConditions *AccessConditions
}

type RenameOptions struct {
	// SourceModifiedAccessConditions specifies parameters for accessing the source directory
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	// AccessConditions specifies parameters for accessing the destination directory
	AccessConditions *AccessConditions
}

type GetPropertiesOptions struct {
	// AccessConditions specifies parameters for accessing the directory
	AccessConditions *AccessConditions
	// CPKInfo contains a group of parameters for client provided encryption key.
	CPKInfo *CPKInfo
}

// ===================================== PATH IMPORTS ===========================================

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = path.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = path.ModifiedAccessConditions

// CPKInfo contains a group of parameters for client provided encryption key.
type CPKInfo = path.CPKInfo

// CPKScopeInfo contains a group of parameters for client provided encryption scope.
type CPKScopeInfo = path.CPKScopeInfo

// AccessConditions identifies container-specific access conditions which you optionally set.
type AccessConditions = path.AccessConditions

// HTTPHeaders contains the HTTP headers for path operations.
type HTTPHeaders = path.HTTPHeaders

// SourceModifiedAccessConditions identifies the source path access conditions.
type SourceModifiedAccessConditions = path.SourceModifiedAccessConditions

// SetAccessControlRecursiveOptions contains the optional parameters when calling the SetAccessControlRecursive operation.
type SetAccessControlRecursiveOptions = path.SetAccessControlRecursiveOptions

// SetMetadataOptions contains the optional parameters when calling the SetMetadata operation.
type SetMetadataOptions = path.SetMetadataOptions

// SetHTTPHeadersOptions contains the optional parameters when calling the SetHTTPHeaders operation.
type SetHTTPHeadersOptions = path.SetHTTPHeadersOptions

// RemoveAccessControlRecursiveOptions contains the optional parameters when calling the RemoveAccessControlRecursive operation.
type RemoveAccessControlRecursiveOptions = path.RemoveAccessControlRecursiveOptions

// UpdateAccessControlRecursiveOptions contains the optional parameters when calling the UpdateAccessControlRecursive operation.
type UpdateAccessControlRecursiveOptions = path.UpdateAccessControlRecursiveOptions

// SetAccessControlOptions contains the optional parameters when calling the SetAccessControl operation.
type SetAccessControlOptions = path.SetAccessControlOptions
