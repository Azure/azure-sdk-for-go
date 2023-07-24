//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
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
	// AccessConditions specifies parameters for accessing the directory
	AccessConditions *AccessConditions
}

func (o *DeleteOptions) format() (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, error) {
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	return leaseAccessConditions, modifiedAccessConditions, nil
}

type RenameOptions struct {
	// SourceModifiedAccessConditions specifies parameters for accessing the source directory
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	// AccessConditions specifies parameters for accessing the destination directory
	AccessConditions *AccessConditions
}

// ===================================== PATH IMPORTS ===========================================

// SetAccessControlRecursiveOptions contains the optional parameters when calling the SetAccessControlRecursive operation. TODO: Design formatter
type SetAccessControlRecursiveOptions struct {
	// ACL is the access control list for the path.
	ACL *string
	// BatchSize is the number of paths to set access control recursively in a single call.
	BatchSize *int32
	// MaxBatches is the maximum number of batches to perform the operation on.
	MaxBatches *int32
	// ContinueOnFailure indicates whether to continue on failure when the operation encounters an error.
	ContinueOnFailure *bool
	// Marker is the continuation token to use when continuing the operation.
	Marker *string
}

func (o *SetAccessControlRecursiveOptions) format() (*generated.PathClientSetAccessControlRecursiveOptions, error) {
	// TODO: design formatter
	return nil, nil
}

// UpdateAccessControlRecursiveOptions contains the optional parameters when calling the UpdateAccessControlRecursive operation. TODO: Design formatter
type UpdateAccessControlRecursiveOptions struct {
	// ACL is the access control list for the path.
	ACL *string
	// BatchSize is the number of paths to set access control recursively in a single call.
	BatchSize *int32
	// MaxBatches is the maximum number of batches to perform the operation on.
	MaxBatches *int32
	// ContinueOnFailure indicates whether to continue on failure when the operation encounters an error.
	ContinueOnFailure *bool
	// Marker is the continuation token to use when continuing the operation.
	Marker *string
}

func (o *UpdateAccessControlRecursiveOptions) format() (*generated.PathClientSetAccessControlRecursiveOptions, error) {
	// TODO: design formatter - similar to SetAccessControlRecursiveOptions
	return nil, nil
}

// RemoveAccessControlRecursiveOptions contains the optional parameters when calling the RemoveAccessControlRecursive operation. TODO: Design formatter
type RemoveAccessControlRecursiveOptions struct {
	// ACL is the access control list for the path.
	ACL *string
	// BatchSize is the number of paths to set access control recursively in a single call.
	BatchSize *int32
	// MaxBatches is the maximum number of batches to perform the operation on.
	MaxBatches *int32
	// ContinueOnFailure indicates whether to continue on failure when the operation encounters an error.
	ContinueOnFailure *bool
	// Marker is the continuation token to use when continuing the operation.
	Marker *string
}

func (o *RemoveAccessControlRecursiveOptions) format() (*generated.PathClientSetAccessControlRecursiveOptions, error) {
	// TODO: design formatter - similar to SetAccessControlRecursiveOptions
	return nil, nil
}

// ================================= path imports ==================================

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method
type GetPropertiesOptions = path.GetPropertiesOptions

// SetAccessControlOptions contains the optional parameters when calling the SetAccessControl operation. dfs endpoint
type SetAccessControlOptions = path.SetAccessControlOptions

// GetAccessControlOptions contains the optional parameters when calling the GetAccessControl operation.
type GetAccessControlOptions = path.GetAccessControlOptions

// CPKInfo contains a group of parameters for the PathClient.Download method.
type CPKInfo = path.CPKInfo

// GetSASURLOptions contains the optional parameters for the Client.GetSASURL method.
type GetSASURLOptions = path.GetSASURLOptions

// SetHTTPHeadersOptions contains the optional parameters for the Client.SetHTTPHeaders method.
type SetHTTPHeadersOptions = path.SetHTTPHeadersOptions

// HTTPHeaders contains the HTTP headers for path operations.
type HTTPHeaders = path.HTTPHeaders

// SetMetadataOptions provides set of configurations for Set Metadata on path operation
type SetMetadataOptions = path.SetMetadataOptions

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = path.SharedKeyCredential

// AccessConditions identifies blob-specific access conditions which you optionally set.
type AccessConditions = path.AccessConditions

// SourceAccessConditions identifies blob-specific access conditions which you optionally set.
type SourceAccessConditions = path.SourceAccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = path.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = path.ModifiedAccessConditions

// SourceModifiedAccessConditions contains a group of parameters for specifying access conditions.
type SourceModifiedAccessConditions = path.SourceModifiedAccessConditions

// CPKScopeInfo contains a group of parameters for the PathClient.SetMetadata method.
type CPKScopeInfo path.CPKScopeInfo
