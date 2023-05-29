//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// SetAccessControlOptions contains the optional parameters when calling the SetAccessControl operation. dfs endpoint
type SetAccessControlOptions struct {
	// Owner is the owner of the path.
	Owner *string
	// Group is the owning group of the path.
	Group *string
	// ACL is the access control list for the path.
	ACL *string
	// Permissions is the octal representation of the permissions for user, group and mask.
	Permissions *string
	// AccessConditions contains parameters for accessing the path.
	AccessConditions *AccessConditions
}

func (o *SetAccessControlOptions) format() (*generated.PathClientSetAccessControlOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, error) {
	if o == nil {
		return nil, nil, nil, nil
	}
	// call path formatter since we're hitting dfs in this operation
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	return &generated.PathClientSetAccessControlOptions{
		Owner:       o.Owner,
		Group:       o.Group,
		ACL:         o.ACL,
		Permissions: o.Permissions,
	}, leaseAccessConditions, modifiedAccessConditions, nil
}

// GetAccessControlOptions contains the optional parameters when calling the GetAccessControl operation.
type GetAccessControlOptions struct {
	// UPN is the user principal name.
	UPN *bool
	// AccessConditions contains parameters for accessing the path.
	AccessConditions *AccessConditions
}

func (o *GetAccessControlOptions) format() (*generated.PathClientGetPropertiesOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, error) {
	action := generated.PathGetPropertiesActionGetAccessControl
	if o == nil {
		return &generated.PathClientGetPropertiesOptions{
			Action: &action,
		}, nil, nil, nil
	}
	// call path formatter since we're hitting dfs in this operation
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	return &generated.PathClientGetPropertiesOptions{
		Upn:    o.UPN,
		Action: &action,
	}, leaseAccessConditions, modifiedAccessConditions, nil
}

// SetAccessControlRecursiveOptions contains the optional parameters when calling the SetAccessControlRecursive operation.
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

// UpdateAccessControlRecursiveOptions contains the optional parameters when calling the UpdateAccessControlRecursive operation.
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

// RemoveAccessControlRecursiveOptions contains the optional parameters when calling the RemoveAccessControlRecursive operation.
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

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions

// SetHTTPHeadersOptions contains the optional parameters when calling the SetHTTPHeaders operation.
type SetHTTPHeadersOptions = blob.SetHTTPHeadersOptions

// HTTPHeaders contains the HTTP headers for path operations.
type HTTPHeaders struct {
	// Optional. Sets the blob's cache control. If specified, this property is stored with the blob and returned with a read request.
	CacheControl *string
	// Optional. Sets the blob's Content-Disposition header.
	ContentDisposition *string
	// Optional. Sets the blob's content encoding. If specified, this property is stored with the blob and returned with a read
	// request.
	ContentEncoding *string
	// Optional. Set the blob's content language. If specified, this property is stored with the blob and returned with a read
	// request.
	ContentLanguage *string
	// Specify the transactional md5 for the body, to be validated by the service.
	ContentMD5 []byte
	// Optional. Sets the blob's content type. If specified, this property is stored with the blob and returned with a read request.
	ContentType *string
}

func (o *HTTPHeaders) formatBlobHTTPHeaders() (*blob.HTTPHeaders, error) {
	if o == nil {
		return nil, nil
	}
	opts := blob.HTTPHeaders{
		BlobCacheControl:       o.CacheControl,
		BlobContentDisposition: o.ContentDisposition,
		BlobContentEncoding:    o.ContentEncoding,
		BlobContentLanguage:    o.ContentLanguage,
		BlobContentMD5:         o.ContentMD5,
		BlobContentType:        o.ContentType,
	}
	return &opts, nil
}

func (o *HTTPHeaders) formatPathHTTPHeaders() (*generated.PathHTTPHeaders, error) {
	if o == nil {
		return nil, nil
	}
	opts := generated.PathHTTPHeaders{
		CacheControl:             o.CacheControl,
		ContentDisposition:       o.ContentDisposition,
		ContentEncoding:          o.ContentEncoding,
		ContentLanguage:          o.ContentLanguage,
		ContentMD5:               o.ContentMD5,
		ContentType:              o.ContentType,
		TransactionalContentHash: o.ContentMD5,
	}
	return &opts, nil
}

// SetMetadataOptions provides set of configurations for Set Metadata on blob operation
type SetMetadataOptions struct {
	AccessConditions *AccessConditions
	CPKInfo          *CPKInfo
	CPKScopeInfo     *CPKScopeInfo
}

func (o *SetMetadataOptions) format() *blob.SetMetadataOptions {
	if o == nil {
		return nil
	}
	return &blob.SetMetadataOptions{
		AccessConditions: o.AccessConditions,
		CPKInfo:          o.CPKInfo,
		CPKScopeInfo:     o.CPKScopeInfo,
	}
}

// AccessConditions identifies container-specific access conditions which you optionally set.
type AccessConditions = exported.PathAccessConditions

// SourceModifiedAccessConditions identifies the source path access conditions.
type SourceModifiedAccessConditions = generated.SourceModifiedAccessConditions

// CPKInfo contains a group of parameters for client provided encryption key.
type CPKInfo = blob.CPKInfo

// CPKScopeInfo contains a group of parameters for client provided encryption scope.
type CPKScopeInfo = blob.CPKScopeInfo
