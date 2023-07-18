//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"net/http"
	"strconv"
	"time"
)

// CreateOptions contains the optional parameters when calling the Create operation. dfs endpoint.
type CreateOptions struct {
	// AccessConditions contains parameters for accessing the file.
	AccessConditions *AccessConditions
	// Metadata is a map of name-value pairs to associate with the file storage object.
	Metadata map[string]*string
	// CPKInfo contains a group of parameters for client provided encryption key.
	CPKInfo *CPKInfo
	// HTTPHeaders contains the HTTP headers for path operations.
	HTTPHeaders *HTTPHeaders
	// Expiry specifies the type and time of expiry for the file.
	Expiry CreationExpiryType
	// LeaseDuration specifies the duration of the lease, in seconds, or negative one
	// (-1) for a lease that never expires. A non-infinite lease can be
	// between 15 and 60 seconds.
	LeaseDuration *int64
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

func (o *CreateOptions) format() (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, *generated.PathHTTPHeaders, *generated.PathClientCreateOptions, *generated.CPKInfo) {
	resource := generated.PathResourceTypeFile
	createOpts := &generated.PathClientCreateOptions{
		Resource: &resource,
	}
	if o == nil {
		return nil, nil, nil, createOpts, nil
	}
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	if o.Expiry == nil {
		createOpts.ExpiryOptions = nil
		createOpts.ExpiresOn = nil
	} else {
		expOpts, expiresOn := o.Expiry.Format()
		createOpts.ExpiryOptions = (*generated.PathExpiryOptions)(&expOpts)
		createOpts.ExpiresOn = expiresOn
	}
	createOpts.ACL = o.ACL
	createOpts.Group = o.Group
	createOpts.Owner = o.Owner
	createOpts.Umask = o.Umask
	createOpts.Permissions = o.Permissions
	createOpts.ProposedLeaseID = o.ProposedLeaseID
	createOpts.LeaseDuration = o.LeaseDuration

	var httpHeaders *generated.PathHTTPHeaders
	var cpkOpts *generated.CPKInfo

	if o.HTTPHeaders != nil {
		httpHeaders = o.HTTPHeaders.formatPathHTTPHeaders()
	}
	if o.CPKInfo != nil {
		cpkOpts = &generated.CPKInfo{
			EncryptionAlgorithm: (*generated.EncryptionAlgorithmType)(o.CPKInfo.EncryptionAlgorithm),
			EncryptionKey:       o.CPKInfo.EncryptionKey,
			EncryptionKeySHA256: o.CPKInfo.EncryptionKeySHA256,
		}
	}
	return leaseAccessConditions, modifiedAccessConditions, httpHeaders, createOpts, cpkOpts
}

// DeleteOptions contains the optional parameters when calling the Delete operation. dfs endpoint
type DeleteOptions struct {
	// AccessConditions contains parameters for accessing the file.
	AccessConditions *AccessConditions
}

func (o *DeleteOptions) format() (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, *generated.PathClientDeleteOptions) {
	recursive := false
	deleteOpts := &generated.PathClientDeleteOptions{
		Recursive: &recursive,
	}
	if o == nil {
		return nil, nil, deleteOpts
	}
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	return leaseAccessConditions, modifiedAccessConditions, deleteOpts
}

// RenameOptions contains the optional parameters when calling the Rename operation.
type RenameOptions struct {
	// SourceAccessConditions identifies the source path access conditions.
	SourceAccessConditions *SourceAccessConditions
	// AccessConditions contains parameters for accessing the file.
	AccessConditions *AccessConditions
}

func (o *RenameOptions) format() (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, *generated.SourceModifiedAccessConditions, *generated.PathClientCreateOptions) {
	// we don't need sourceModAccCond since this is not rename
	mode := generated.PathRenameModeLegacy
	resource := generated.PathResourceTypeFile
	createOpts := &generated.PathClientCreateOptions{
		Mode:     &mode,
		Resource: &resource,
	}
	if o == nil {
		return nil, nil, nil, createOpts
	}
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	if o.SourceAccessConditions != nil {
		if o.SourceAccessConditions.SourceModifiedAccessConditions != nil {
			sourceModifiedAccessConditions := &generated.SourceModifiedAccessConditions{
				SourceIfMatch:           o.SourceAccessConditions.SourceModifiedAccessConditions.SourceIfMatch,
				SourceIfModifiedSince:   o.SourceAccessConditions.SourceModifiedAccessConditions.SourceIfModifiedSince,
				SourceIfNoneMatch:       o.SourceAccessConditions.SourceModifiedAccessConditions.SourceIfNoneMatch,
				SourceIfUnmodifiedSince: o.SourceAccessConditions.SourceModifiedAccessConditions.SourceIfUnmodifiedSince,
			}
			createOpts.SourceLeaseID = o.SourceAccessConditions.SourceLeaseAccessConditions.LeaseID
			return leaseAccessConditions, modifiedAccessConditions, sourceModifiedAccessConditions, createOpts
		}
	}
	return leaseAccessConditions, modifiedAccessConditions, nil, createOpts
}

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method
type GetPropertiesOptions struct {
	AccessConditions *AccessConditions
	CPKInfo          *CPKInfo
}

func (o *GetPropertiesOptions) format() *blob.GetPropertiesOptions {
	if o == nil {
		return nil
	}
	accessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return &blob.GetPropertiesOptions{
		AccessConditions: accessConditions,
		CPKInfo: &blob.CPKInfo{
			EncryptionKey:       o.CPKInfo.EncryptionKey,
			EncryptionAlgorithm: o.CPKInfo.EncryptionAlgorithm,
			EncryptionKeySHA256: o.CPKInfo.EncryptionKeySHA256,
		},
	}
}

// ===================================== PATH IMPORTS ===========================================

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

func (o *SetAccessControlOptions) format() (*generated.PathClientSetAccessControlOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}
	// call path formatter since we're hitting dfs in this operation
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	return &generated.PathClientSetAccessControlOptions{
		Owner:       o.Owner,
		Group:       o.Group,
		ACL:         o.ACL,
		Permissions: o.Permissions,
	}, leaseAccessConditions, modifiedAccessConditions
}

// GetAccessControlOptions contains the optional parameters when calling the GetAccessControl operation.
type GetAccessControlOptions struct {
	// UPN is the user principal name.
	UPN *bool
	// AccessConditions contains parameters for accessing the path.
	AccessConditions *AccessConditions
}

func (o *GetAccessControlOptions) format() (*generated.PathClientGetPropertiesOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	action := generated.PathGetPropertiesActionGetAccessControl
	if o == nil {
		return &generated.PathClientGetPropertiesOptions{
			Action: &action,
		}, nil, nil
	}
	// call path formatter since we're hitting dfs in this operation
	leaseAccessConditions, modifiedAccessConditions := exported.FormatPathAccessConditions(o.AccessConditions)
	return &generated.PathClientGetPropertiesOptions{
		Upn:    o.UPN,
		Action: &action,
	}, leaseAccessConditions, modifiedAccessConditions
}

// UpdateAccessControlOptions contains the optional parameters when calling the UpdateAccessControlRecursive operation.
type UpdateAccessControlOptions struct {
	// ACL is the access control list for the path.
	ACL *string
}

func (o *UpdateAccessControlOptions) format() (*generated.PathClientSetAccessControlRecursiveOptions, generated.PathSetAccessControlRecursiveMode) {
	mode := generated.PathSetAccessControlRecursiveModeModify
	if o == nil {
		return nil, mode
	}
	opts := &generated.PathClientSetAccessControlRecursiveOptions{
		ACL: o.ACL,
	}
	return opts, mode
}

// RemoveAccessControlOptions contains the optional parameters when calling the RemoveAccessControlRecursive operation.
type RemoveAccessControlOptions struct {
	//placeholder
}

func (o *RemoveAccessControlOptions) format() (*generated.PathClientSetAccessControlRecursiveOptions, generated.PathSetAccessControlRecursiveMode) {
	mode := generated.PathSetAccessControlRecursiveModeRemove
	return nil, mode
}

// SetHTTPHeadersOptions contains the optional parameters for the Client.SetHTTPHeaders method.
type SetHTTPHeadersOptions struct {
	AccessConditions *AccessConditions
}

func (o *SetHTTPHeadersOptions) format(httpHeaders HTTPHeaders) (*blob.SetHTTPHeadersOptions, blob.HTTPHeaders) {
	httpHeaderOpts := blob.HTTPHeaders{
		BlobCacheControl:       httpHeaders.CacheControl,
		BlobContentDisposition: httpHeaders.ContentDisposition,
		BlobContentEncoding:    httpHeaders.ContentEncoding,
		BlobContentLanguage:    httpHeaders.ContentLanguage,
		BlobContentMD5:         httpHeaders.ContentMD5,
		BlobContentType:        httpHeaders.ContentType,
	}
	if o == nil {
		return nil, httpHeaderOpts
	}
	accessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return &blob.SetHTTPHeadersOptions{
		AccessConditions: accessConditions,
	}, httpHeaderOpts
}

// HTTPHeaders contains the HTTP headers for path operations.
type HTTPHeaders struct {
	// Optional. Sets the path's cache control. If specified, this property is stored with the path and returned with a read request.
	CacheControl *string
	// Optional. Sets the path's Content-Disposition header.
	ContentDisposition *string
	// Optional. Sets the path's content encoding. If specified, this property is stored with the blobpath and returned with a read
	// request.
	ContentEncoding *string
	// Optional. Set the path's content language. If specified, this property is stored with the path and returned with a read
	// request.
	ContentLanguage *string
	// Specify the transactional md5 for the body, to be validated by the service.
	ContentMD5 []byte
	// Optional. Sets the path's content type. If specified, this property is stored with the path and returned with a read request.
	ContentType *string
}

//
//func (o HTTPHeaders) formatBlobHTTPHeaders() blob.HTTPHeaders {
//
//	opts := blob.HTTPHeaders{
//		BlobCacheControl:       o.CacheControl,
//		BlobContentDisposition: o.ContentDisposition,
//		BlobContentEncoding:    o.ContentEncoding,
//		BlobContentLanguage:    o.ContentLanguage,
//		BlobContentMD5:         o.ContentMD5,
//		BlobContentType:        o.ContentType,
//	}
//	return opts
//}

func (o *HTTPHeaders) formatPathHTTPHeaders() *generated.PathHTTPHeaders {
	// TODO: will be used for file related ops, like append
	if o == nil {
		return nil
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
	return &opts
}

// SetMetadataOptions provides set of configurations for Set Metadata on path operation
type SetMetadataOptions struct {
	AccessConditions *AccessConditions
	CPKInfo          *CPKInfo
	CPKScopeInfo     *CPKScopeInfo
}

func (o *SetMetadataOptions) format() *blob.SetMetadataOptions {
	if o == nil {
		return nil
	}
	accessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return &blob.SetMetadataOptions{
		AccessConditions: accessConditions,
		CPKInfo: &blob.CPKInfo{
			EncryptionKey:       o.CPKInfo.EncryptionKey,
			EncryptionAlgorithm: o.CPKInfo.EncryptionAlgorithm,
			EncryptionKeySHA256: o.CPKInfo.EncryptionKeySHA256,
		},
		CPKScopeInfo: (*blob.CPKScopeInfo)(o.CPKScopeInfo),
	}
}

// CPKInfo contains a group of parameters for the PathClient.Download method.
type CPKInfo struct {
	EncryptionAlgorithm *EncryptionAlgorithmType
	EncryptionKey       *string
	EncryptionKeySHA256 *string
}

// CreationExpiryType defines values for Create() ExpiryType
type CreationExpiryType interface {
	Format() (generated.ExpiryOptions, *string)
	notPubliclyImplementable()
}

// CreationExpiryTypeAbsolute defines the absolute time for the blob expiry
type CreationExpiryTypeAbsolute time.Time

// CreationExpiryTypeRelativeToNow defines the duration relative to now for the blob expiry
type CreationExpiryTypeRelativeToNow time.Duration

// CreationExpiryTypeNever defines that the blob will be set to never expire
type CreationExpiryTypeNever struct {
	// empty struct since NeverExpire expiry type does not require expiry time
}

func (e CreationExpiryTypeAbsolute) Format() (generated.ExpiryOptions, *string) {
	return generated.ExpiryOptionsAbsolute, to.Ptr(time.Time(e).UTC().Format(http.TimeFormat))

}

func (e CreationExpiryTypeAbsolute) notPubliclyImplementable() {}

func (e CreationExpiryTypeRelativeToNow) Format() (generated.ExpiryOptions, *string) {
	return generated.ExpiryOptionsRelativeToNow, to.Ptr(strconv.FormatInt(time.Duration(e).Milliseconds(), 10))
}

func (e CreationExpiryTypeRelativeToNow) notPubliclyImplementable() {}

func (e CreationExpiryTypeNever) Format() (generated.ExpiryOptions, *string) {
	return generated.ExpiryOptionsNeverExpire, nil
}

func (e CreationExpiryTypeNever) notPubliclyImplementable() {}

// ACLFailedEntry contains the failed ACL entry (response model).
type ACLFailedEntry = generated.ACLFailedEntry

// CPKScopeInfo contains a group of parameters for the PathClient.SetMetadata method.
type CPKScopeInfo blob.CPKScopeInfo

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// SetExpiryType defines values for ExpiryType.
type SetExpiryType = exported.SetExpiryType

// SetExpiryTypeAbsolute defines the absolute time for the expiry.
type SetExpiryTypeAbsolute = exported.SetExpiryTypeAbsolute

// SetExpiryTypeRelativeToNow defines the duration relative to now for the expiry.
type SetExpiryTypeRelativeToNow = exported.SetExpiryTypeRelativeToNow

// SetExpiryTypeRelativeToCreation defines the duration relative to creation for the expiry.
type SetExpiryTypeRelativeToCreation = exported.SetExpiryTypeRelativeToCreation

// SetExpiryTypeNever defines that will be set to never expire.
type SetExpiryTypeNever = exported.SetExpiryTypeNever

// SetExpiryOptions contains the optional parameters for the Client.SetExpiry method.
type SetExpiryOptions = exported.SetExpiryOptions

// AccessConditions identifies blob-specific access conditions which you optionally set.
type AccessConditions = exported.AccessConditions

// SourceAccessConditions identifies blob-specific access conditions which you optionally set.
type SourceAccessConditions = exported.SourceAccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions

// SourceModifiedAccessConditions contains a group of parameters for specifying access conditions.
type SourceModifiedAccessConditions = exported.SourceModifiedAccessConditions
