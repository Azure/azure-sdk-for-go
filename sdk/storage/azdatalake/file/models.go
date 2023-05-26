//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
	"time"
)

// CreateOptions contains the optional parameters when calling the Create operation. dfs endpoint
type CreateOptions struct {
	// AccessConditions contains parameters for accessing the file.
	AccessConditions *AccessConditions
	// Metadata is a map of name-value pairs to associate with the file storage object.
	Metadata map[string]*string
	// CPKInfo contains a group of parameters for client provided encryption key.
	CPKInfo *CPKInfo
	// HTTPHeaders contains the HTTP headers for path operations.
	HTTPHeaders *HTTPHeaders
	//PathExpiryOptions *ExpiryOptions
	// ExpiresOn specifies the time that the file will expire.
	ExpiresOn *time.Time
	// LeaseDuration specifies the duration of the lease.
	LeaseDuration *time.Duration
	// ProposedLeaseID specifies the proposed lease ID for the file.
	ProposedLeaseID *string
	// Permissions is the octal representation of the permissions for user, group and mask.
	Permissions *string
	// Umask is the umask for the file.
	Umask *string
	// Owner is the owner of the file.
	Owner *string
	// Group is the owning group of the file.
	Group *string
	// ACL is the access control list for the file.
	ACL *string
}

func (o *CreateOptions) format() (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, *generated.PathHTTPHeaders, error) {
	// TODO: add all other required options for the create operation, we don't need sourceModAccCond since this is not rename
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	httpHeaders := &generated.PathHTTPHeaders{
		CacheControl:             o.HTTPHeaders.CacheControl,
		ContentDisposition:       o.HTTPHeaders.ContentDisposition,
		ContentEncoding:          o.HTTPHeaders.ContentEncoding,
		ContentLanguage:          o.HTTPHeaders.ContentLanguage,
		ContentMD5:               o.HTTPHeaders.ContentMD5,
		ContentType:              o.HTTPHeaders.ContentType,
		TransactionalContentHash: o.HTTPHeaders.ContentMD5,
	}
	return leaseAccessConditions, modifiedAccessConditions, httpHeaders, nil
}

// DeleteOptions contains the optional parameters when calling the Delete operation. dfs endpoint
type DeleteOptions struct {
	// AccessConditions contains parameters for accessing the file.
	AccessConditions *AccessConditions
}

func (o *DeleteOptions) format() (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, error) {
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	return leaseAccessConditions, modifiedAccessConditions, nil
}

// RenameOptions contains the optional parameters when calling the Rename operation.
type RenameOptions struct {
	// SourceModifiedAccessConditions identifies the source path access conditions.
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	// AccessConditions contains parameters for accessing the file.
	AccessConditions *AccessConditions
}

// GetPropertiesOptions contains the optional parameters when calling the GetProperties operation.
type GetPropertiesOptions blob.GetPropertiesOptions

// SetExpiryOptions contains the optional parameters when calling the SetExpiry operation.
type SetExpiryOptions struct {
	// ExpiresOn specifies the time that the file will expire.
	ExpiresOn *time.Time
}

// ===================================== PATH IMPORTS ===========================================

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

// ModifiedAccessConditions identifies path-specific access conditions which you optionally set.
type ModifiedAccessConditions = path.ModifiedAccessConditions

// LeaseAccessConditions identifies path-specific access conditions associated with a lease.
type LeaseAccessConditions = path.LeaseAccessConditions
