//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
	"time"
)

// CreateOptions contains the optional parameters when calling the Create operation. dfs endpoint. TODO: Design formatter
type CreateOptions struct {
	// AccessConditions contains parameters for accessing the file.
	AccessConditions *path.AccessConditions
	// Metadata is a map of name-value pairs to associate with the file storage object.
	Metadata map[string]*string
	// CPKInfo contains a group of parameters for client provided encryption key.
	CPKInfo *path.CPKInfo
	// HTTPHeaders contains the HTTP headers for path operations.
	HTTPHeaders *path.HTTPHeaders
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
	leaseAccessConditions, modifiedAccessConditions := azdatalake.FormatPathAccessConditions(o.AccessConditions)
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
	AccessConditions *path.AccessConditions
}

func (o *DeleteOptions) format() (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions, error) {
	leaseAccessConditions, modifiedAccessConditions := azdatalake.FormatPathAccessConditions(o.AccessConditions)
	return leaseAccessConditions, modifiedAccessConditions, nil
}

// RenameOptions contains the optional parameters when calling the Rename operation. TODO: Design formatter
type RenameOptions struct {
	// SourceModifiedAccessConditions identifies the source path access conditions.
	SourceModifiedAccessConditions *path.SourceModifiedAccessConditions
	// AccessConditions contains parameters for accessing the file.
	AccessConditions *path.AccessConditions
}

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method
type GetPropertiesOptions struct {
	AccessConditions *path.AccessConditions
	CPKInfo          *path.CPKInfo
}

func (o *GetPropertiesOptions) format() *blob.GetPropertiesOptions {
	if o == nil {
		return nil
	}
	accessConditions := azdatalake.FormatBlobAccessConditions(o.AccessConditions)
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

// ExpiryType defines values for ExpiryType.
type ExpiryType = exported.ExpiryType

// ExpiryTypeAbsolute defines the absolute time for the expiry.
type ExpiryTypeAbsolute = exported.ExpiryTypeAbsolute

// ExpiryTypeRelativeToNow defines the duration relative to now for the expiry.
type ExpiryTypeRelativeToNow = exported.ExpiryTypeRelativeToNow

// ExpiryTypeRelativeToCreation defines the duration relative to creation for the expiry.
type ExpiryTypeRelativeToCreation = exported.ExpiryTypeRelativeToCreation

// ExpiryTypeNever defines that will be set to never expire.
type ExpiryTypeNever = exported.ExpiryTypeNever

// SetExpiryOptions contains the optional parameters for the Client.SetExpiry method.
type SetExpiryOptions = exported.SetExpiryOptions
