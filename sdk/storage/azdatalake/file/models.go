//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
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
		httpHeaders = path.FormatPathHTTPHeaders(o.HTTPHeaders)
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

func (o *RenameOptions) format(path string) (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, *generated.SourceModifiedAccessConditions, *generated.PathClientCreateOptions) {
	// we don't need sourceModAccCond since this is not rename
	mode := generated.PathRenameModeLegacy
	createOpts := &generated.PathClientCreateOptions{
		Mode:         &mode,
		RenameSource: &path,
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

// ===================================== PATH IMPORTS ===========================================

// UpdateAccessControlOptions contains the optional parameters when calling the UpdateAccessControlRecursive operation.
type UpdateAccessControlOptions struct {
	//placeholder
}

func (o *UpdateAccessControlOptions) format(ACL string) (*generated.PathClientSetAccessControlRecursiveOptions, generated.PathSetAccessControlRecursiveMode) {
	mode := generated.PathSetAccessControlRecursiveModeModify
	return &generated.PathClientSetAccessControlRecursiveOptions{
		ACL: &ACL,
	}, mode
}

// RemoveAccessControlOptions contains the optional parameters when calling the RemoveAccessControlRecursive operation.
type RemoveAccessControlOptions struct {
	//placeholder
}

func (o *RemoveAccessControlOptions) format(ACL string) (*generated.PathClientSetAccessControlRecursiveOptions, generated.PathSetAccessControlRecursiveMode) {
	mode := generated.PathSetAccessControlRecursiveModeRemove
	return &generated.PathClientSetAccessControlRecursiveOptions{
		ACL: &ACL,
	}, mode
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

// SetAccessControlRecursiveResponse contains part of the response data returned by the []OP_AccessControl operations.
type SetAccessControlRecursiveResponse = generated.SetAccessControlRecursiveResponse

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
