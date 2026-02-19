// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package appendblob

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// Type Declarations ---------------------------------------------------------------------

// AppendPositionAccessConditions contains a group of parameters for the Client.AppendBlock method.
type AppendPositionAccessConditions = generated.AppendPositionAccessConditions

// Request Model Declaration -------------------------------------------------------------------------------------------

// CreateOptions provides set of configurations for Create Append Blob operation
type CreateOptions struct {
	// Specifies the date time when the blobs immutability policy is set to expire.
	ImmutabilityPolicyExpiry *time.Time

	// Specifies the immutability policy mode to set on the blob.
	ImmutabilityPolicyMode *blob.ImmutabilityPolicySetting

	// Specified if a legal hold should be set on the blob.
	LegalHold *bool

	AccessConditions *blob.AccessConditions

	HTTPHeaders *blob.HTTPHeaders

	CPKInfo *blob.CPKInfo

	CPKScopeInfo *blob.CPKScopeInfo

	// Optional. Used to set blob tags in various blob operations.
	Tags map[string]string

	// Optional. Specifies a user-defined name-value pair associated with the blob. If no name-value pairs are specified, the
	// operation will copy the metadata from the source blob or file to the destination blob. If one or more name-value pairs
	// are specified, the destination blob is created with the specified metadata, and metadata is not copied from the source
	// blob or file. Note that beginning with version 2009-09-19, metadata names must adhere to the naming rules for C# identifiers.
	// See Naming and Referencing Containers, Blobs, and Metadata for more information.
	Metadata map[string]*string
}

func (o *CreateOptions) format() *generated.AppendBlobClientCreateOptions {
	if o == nil {
		return nil
	}
	return &generated.AppendBlobClientCreateOptions{
		ImmutabilityPolicyExpiry: o.ImmutabilityPolicyExpiry,
		ImmutabilityPolicyMode:   o.ImmutabilityPolicyMode,
		LegalHold:                o.LegalHold,
		LeaseID:                  o.AccessConditions.LeaseAccessConditions.LeaseID,
		IfMatch:                  o.AccessConditions.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:          o.AccessConditions.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:              o.AccessConditions.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince:        o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince,
		BlobCacheControl:         o.HTTPHeaders.BlobCacheControl,
		BlobContentDisposition:   o.HTTPHeaders.BlobContentDisposition,
		BlobContentEncoding:      o.HTTPHeaders.BlobContentEncoding,
		BlobContentLanguage:      o.HTTPHeaders.BlobContentLanguage,
		BlobContentMD5:           o.HTTPHeaders.BlobContentMD5,
		BlobContentType:          o.HTTPHeaders.BlobContentType,
		EncryptionAlgorithm:      o.CPKInfo.EncryptionAlgorithm,
		EncryptionKey:            o.CPKInfo.EncryptionKey,
		EncryptionKeySHA256:      o.CPKInfo.EncryptionKeySHA256,
		EncryptionScope:          o.CPKScopeInfo.EncryptionScope,
		BlobTagsString:           shared.SerializeBlobTagsToStrPtr(o.Tags),
		Metadata:                 o.Metadata,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// AppendBlockOptions contains the optional parameters for the Client.AppendBlock method.
type AppendBlockOptions struct {
	// TransactionalValidation specifies the transfer validation type to use.
	// The default is nil (no transfer validation).
	TransactionalValidation blob.TransferValidationType

	AppendPositionAccessConditions *AppendPositionAccessConditions

	CPKInfo *blob.CPKInfo

	CPKScopeInfo *blob.CPKScopeInfo

	AccessConditions *blob.AccessConditions
}

// TODO transfer validation fix
func (o *AppendBlockOptions) format() *generated.AppendBlobClientAppendBlockOptions {
	if o == nil {
		return nil
	}
	return &generated.AppendBlobClientAppendBlockOptions{
		AppendPosition:      o.AppendPositionAccessConditions.AppendPosition,
		MaxSize:             o.AppendPositionAccessConditions.MaxSize,
		LeaseID:             o.AccessConditions.LeaseAccessConditions.LeaseID,
		IfMatch:             o.AccessConditions.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:     o.AccessConditions.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:         o.AccessConditions.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince:   o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince,
		EncryptionAlgorithm: o.CPKInfo.EncryptionAlgorithm,
		EncryptionKey:       o.CPKInfo.EncryptionKey,
		EncryptionKeySHA256: o.CPKInfo.EncryptionKeySHA256,
		EncryptionScope:     o.CPKScopeInfo.EncryptionScope,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// AppendBlockFromURLOptions contains the optional parameters for the Client.AppendBlockFromURL method.
type AppendBlockFromURLOptions struct {
	// Only Bearer type is supported. Credentials should be a valid OAuth access token to copy source.
	CopySourceAuthorization *string

	// SourceContentValidation contains the validation mechanism used on the range of bytes read from the source.
	SourceContentValidation blob.SourceContentValidationType

	AppendPositionAccessConditions *AppendPositionAccessConditions

	CPKInfo *blob.CPKInfo

	CPKScopeInfo *blob.CPKScopeInfo

	FileRequestIntent *blob.FileRequestIntentType

	SourceModifiedAccessConditions *blob.SourceModifiedAccessConditions

	AccessConditions *blob.AccessConditions

	// Range specifies a range of bytes.  The default value is all bytes.
	Range blob.HTTPRange
}

func (o *AppendBlockFromURLOptions) format() *generated.AppendBlobClientAppendBlockFromURLOptions {
	if o == nil {
		return nil
	}

	// Notes: no mapping for o.SourceModifiedAccessConditions.SourceIfTags
	options := &generated.AppendBlobClientAppendBlockFromURLOptions{
		CopySourceAuthorization: o.CopySourceAuthorization,
		AppendPosition:          o.AppendPositionAccessConditions.AppendPosition,
		MaxSize:                 o.AppendPositionAccessConditions.MaxSize,
		LeaseID:                 o.AccessConditions.LeaseAccessConditions.LeaseID,
		IfMatch:                 o.AccessConditions.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:         o.AccessConditions.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:             o.AccessConditions.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince:       o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince,
		EncryptionAlgorithm:     o.CPKInfo.EncryptionAlgorithm,
		EncryptionKey:           o.CPKInfo.EncryptionKey,
		EncryptionKeySHA256:     o.CPKInfo.EncryptionKeySHA256,
		EncryptionScope:         o.CPKScopeInfo.EncryptionScope,
		FileRequestIntent:       o.FileRequestIntent,
		SourceIfMatch:           o.SourceModifiedAccessConditions.SourceIfMatch,
		SourceIfModifiedSince:   o.SourceModifiedAccessConditions.SourceIfModifiedSince,
		SourceIfNoneMatch:       o.SourceModifiedAccessConditions.SourceIfNoneMatch,
		SourceIfUnmodifiedSince: o.SourceModifiedAccessConditions.SourceIfUnmodifiedSince,
		SourceRange:             exported.FormatHTTPRange(o.Range),
	}

	if o.SourceContentValidation != nil {
		o.SourceContentValidation.Apply(options)
	}

	return options
}

// ---------------------------------------------------------------------------------------------------------------------

// SealOptions provides set of configurations for SealAppendBlob operation
type SealOptions struct {
	AccessConditions               *blob.AccessConditions
	AppendPositionAccessConditions *AppendPositionAccessConditions
}

func (o *SealOptions) format() *generated.AppendBlobClientSealOptions {
	if o == nil {
		return nil
	}
	// NOTE: no mapping for o.AppendPositionAccessConditions.MaxSize
	return &generated.AppendBlobClientSealOptions{
		AppendPosition:    o.AppendPositionAccessConditions.AppendPosition,
		LeaseID:           o.AccessConditions.LeaseAccessConditions.LeaseID,
		IfMatch:           o.AccessConditions.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:   o.AccessConditions.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:       o.AccessConditions.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince: o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// ExpiryType defines values for ExpiryType
type ExpiryType = exported.ExpiryType

// ExpiryTypeAbsolute defines the absolute time for the blob expiry
type ExpiryTypeAbsolute = exported.ExpiryTypeAbsolute

// ExpiryTypeRelativeToNow defines the duration relative to now for the blob expiry
type ExpiryTypeRelativeToNow = exported.ExpiryTypeRelativeToNow

// ExpiryTypeRelativeToCreation defines the duration relative to creation for the blob expiry
type ExpiryTypeRelativeToCreation = exported.ExpiryTypeRelativeToCreation

// ExpiryTypeNever defines that the blob will be set to never expire
type ExpiryTypeNever = exported.ExpiryTypeNever

// SetExpiryOptions contains the optional parameters for the Client.SetExpiry method.
type SetExpiryOptions = exported.SetExpiryOptions
