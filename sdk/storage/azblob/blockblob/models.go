// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// Type Declarations ---------------------------------------------------------------------

// Block - Represents a single block in a block blob. It describes the block's ID and size.
type Block = generated.Block

// BlockList - can be uncommitted or committed blocks (committed/uncommitted)
type BlockList = generated.BlockList

// Request Model Declaration -------------------------------------------------------------------------------------------

// UploadOptions contains the optional parameters for the Client.Upload method.
type UploadOptions struct {
	// Optional. Used to set blob tags in various blob operations.
	Tags map[string]string

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]*string

	// Optional. Indicates the tier to be set on the blob.
	Tier *blob.AccessTier

	// TransactionalValidation specifies the transfer validation type to use.
	// The default is nil (no transfer validation).
	TransactionalValidation blob.TransferValidationType

	HTTPHeaders                  *blob.HTTPHeaders
	CPKInfo                      *blob.CPKInfo
	CPKScopeInfo                 *blob.CPKScopeInfo
	AccessConditions             *blob.AccessConditions
	LegalHold                    *bool
	ImmutabilityPolicyMode       *blob.ImmutabilityPolicySetting
	ImmutabilityPolicyExpiryTime *time.Time

	// Deprecated: TransactionalContentMD5 can be set by using TransactionalValidation instead
	TransactionalContentMD5 []byte
}

func (o *UploadOptions) format() *generated.BlockBlobClientUploadOptions {
	if o == nil {
		return nil
	}

	return &generated.BlockBlobClientUploadOptions{
		BlobTagsString:           shared.SerializeBlobTagsToStrPtr(o.Tags),
		Metadata:                 o.Metadata,
		Tier:                     o.Tier,
		TransactionalContentMD5:  o.TransactionalContentMD5,
		LegalHold:                o.LegalHold,
		ImmutabilityPolicyMode:   o.ImmutabilityPolicyMode,
		ImmutabilityPolicyExpiry: o.ImmutabilityPolicyExpiryTime,
		LeaseID:                  o.AccessConditions.LeaseAccessConditions.LeaseID,
		IfMatch:                  o.AccessConditions.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:          o.AccessConditions.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:              o.AccessConditions.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince:        o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince,
		EncryptionAlgorithm:      o.CPKInfo.EncryptionAlgorithm,
		EncryptionKey:            o.CPKInfo.EncryptionKey,
		EncryptionKeySHA256:      o.CPKInfo.EncryptionKeySHA256,
		EncryptionScope:          o.CPKScopeInfo.EncryptionScope,
		BlobCacheControl:         o.HTTPHeaders.BlobCacheControl,
		BlobContentDisposition:   o.HTTPHeaders.BlobContentDisposition,
		BlobContentEncoding:      o.HTTPHeaders.BlobContentEncoding,
		BlobContentLanguage:      o.HTTPHeaders.BlobContentLanguage,
		BlobContentMD5:           o.HTTPHeaders.BlobContentMD5,
		BlobContentType:          o.HTTPHeaders.BlobContentType,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// UploadBlobFromURLOptions contains the optional parameters for the Client.UploadBlobFromURL method.
type UploadBlobFromURLOptions struct {
	// Optional. Used to set blob tags in various blob operations.
	Tags map[string]string

	// Only Bearer type is supported. Credentials should be a valid OAuth access token to copy source.
	CopySourceAuthorization *string

	// Valid value is backup
	FileRequestIntent *blob.FileRequestIntentType

	// Optional, default is true. Indicates if properties from the source blob should be copied.
	CopySourceBlobProperties *bool

	// Optional, default 'replace'. Indicates if source tags should be copied or replaced with the tags specified by x-ms-tags.
	CopySourceTags *BlobCopySourceTags

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]*string

	// Optional. Specifies the md5 calculated for the range of bytes that must be read from the copy source.
	SourceContentMD5 []byte

	// Optional. Indicates the tier to be set on the blob.
	Tier *blob.AccessTier

	// Additional optional headers
	HTTPHeaders                    *blob.HTTPHeaders
	AccessConditions               *blob.AccessConditions
	CPKInfo                        *blob.CPKInfo
	CPKScopeInfo                   *blob.CPKScopeInfo
	SourceModifiedAccessConditions *blob.SourceModifiedAccessConditions
}

func (o *UploadBlobFromURLOptions) format() *generated.BlockBlobClientUploadBlobFromURLOptions {
	if o == nil {
		return nil
	}
	return &generated.BlockBlobClientUploadBlobFromURLOptions{
		BlobTagsString:           shared.SerializeBlobTagsToStrPtr(o.Tags),
		CopySourceAuthorization:  o.CopySourceAuthorization,
		FileRequestIntent:        o.FileRequestIntent,
		CopySourceBlobProperties: o.CopySourceBlobProperties,
		CopySourceTags:           o.CopySourceTags,
		Metadata:                 o.Metadata,
		SourceContentMD5:         o.SourceContentMD5,
		Tier:                     o.Tier,
		LeaseID:                  o.AccessConditions.LeaseAccessConditions.LeaseID,
		IfMatch:                  o.AccessConditions.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:          o.AccessConditions.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:              o.AccessConditions.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince:        o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince,
		EncryptionAlgorithm:      o.CPKInfo.EncryptionAlgorithm,
		EncryptionKey:            o.CPKInfo.EncryptionKey,
		EncryptionKeySHA256:      o.CPKInfo.EncryptionKeySHA256,
		EncryptionScope:          o.CPKScopeInfo.EncryptionScope,
		BlobCacheControl:         o.HTTPHeaders.BlobCacheControl,
		BlobContentDisposition:   o.HTTPHeaders.BlobContentDisposition,
		BlobContentEncoding:      o.HTTPHeaders.BlobContentEncoding,
		BlobContentLanguage:      o.HTTPHeaders.BlobContentLanguage,
		BlobContentMD5:           o.HTTPHeaders.BlobContentMD5,
		BlobContentType:          o.HTTPHeaders.BlobContentType,
		SourceIfMatch:            o.SourceModifiedAccessConditions.SourceIfMatch,
		SourceIfModifiedSince:    o.SourceModifiedAccessConditions.SourceIfModifiedSince,
		SourceIfNoneMatch:        o.SourceModifiedAccessConditions.SourceIfNoneMatch,
		SourceIfUnmodifiedSince:  o.SourceModifiedAccessConditions.SourceIfUnmodifiedSince,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// StageBlockOptions contains the optional parameters for the Client.StageBlock method.
type StageBlockOptions struct {
	CPKInfo *blob.CPKInfo

	CPKScopeInfo *blob.CPKScopeInfo

	LeaseAccessConditions *blob.LeaseAccessConditions

	// TransactionalValidation specifies the transfer validation type to use.
	// The default is nil (no transfer validation).
	TransactionalValidation blob.TransferValidationType
}

// StageBlockOptions contains the optional parameters for the Client.StageBlock method.
func (o *StageBlockOptions) format() *generated.BlockBlobClientStageBlockOptions {
	if o == nil {
		return nil
	}

	return &generated.BlockBlobClientStageBlockOptions{
		LeaseID:             o.LeaseAccessConditions.LeaseID,
		EncryptionAlgorithm: o.CPKInfo.EncryptionAlgorithm,
		EncryptionKey:       o.CPKInfo.EncryptionKey,
		EncryptionKeySHA256: o.CPKInfo.EncryptionKeySHA256,
		EncryptionScope:     o.CPKScopeInfo.EncryptionScope,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// StageBlockFromURLOptions contains the optional parameters for the Client.StageBlockFromURL method.
type StageBlockFromURLOptions struct {
	// Only Bearer type is supported. Credentials should be a valid OAuth access token to copy source.
	CopySourceAuthorization *string

	LeaseAccessConditions *blob.LeaseAccessConditions

	SourceModifiedAccessConditions *blob.SourceModifiedAccessConditions

	// SourceContentValidation contains the validation mechanism used on the range of bytes read from the source.
	SourceContentValidation blob.SourceContentValidationType

	// File request Intent. Valid value is backup.
	FileRequestIntent *blob.FileRequestIntentType

	// Range specifies a range of bytes.  The default value is all bytes.
	Range blob.HTTPRange

	CPKInfo *blob.CPKInfo

	CPKScopeInfo *blob.CPKScopeInfo
}

func (o *StageBlockFromURLOptions) format() *generated.BlockBlobClientStageBlockFromURLOptions {
	if o == nil {
		return nil
	}

	options := &generated.BlockBlobClientStageBlockFromURLOptions{
		CopySourceAuthorization: o.CopySourceAuthorization,
		SourceRange:             exported.FormatHTTPRange(o.Range),
		FileRequestIntent:       o.FileRequestIntent,
		SourceIfMatch:           o.SourceModifiedAccessConditions.SourceIfMatch,
		SourceIfModifiedSince:   o.SourceModifiedAccessConditions.SourceIfModifiedSince,
		SourceIfNoneMatch:       o.SourceModifiedAccessConditions.SourceIfNoneMatch,
		SourceIfUnmodifiedSince: o.SourceModifiedAccessConditions.SourceIfUnmodifiedSince,
		EncryptionAlgorithm:     o.CPKInfo.EncryptionAlgorithm,
		EncryptionKey:           o.CPKInfo.EncryptionKey,
		EncryptionKeySHA256:     o.CPKInfo.EncryptionKeySHA256,
		EncryptionScope:         o.CPKScopeInfo.EncryptionScope,
	}

	if o.SourceContentValidation != nil {
		o.SourceContentValidation.Apply(options)
	}

	return options
}

// ---------------------------------------------------------------------------------------------------------------------

// CommitBlockListOptions contains the optional parameters for Client.CommitBlockList method.
type CommitBlockListOptions struct {
	Tags                         map[string]string
	Metadata                     map[string]*string
	RequestID                    *string
	Tier                         *blob.AccessTier
	Timeout                      *int32
	HTTPHeaders                  *blob.HTTPHeaders
	CPKInfo                      *blob.CPKInfo
	CPKScopeInfo                 *blob.CPKScopeInfo
	AccessConditions             *blob.AccessConditions
	LegalHold                    *bool
	ImmutabilityPolicyMode       *blob.ImmutabilityPolicySetting
	ImmutabilityPolicyExpiryTime *time.Time

	// Deprecated: TransactionalContentCRC64 cannot be generated
	TransactionalContentCRC64 []byte

	// Deprecated: TransactionalContentMD5 cannot be generated
	TransactionalContentMD5 []byte
}

func (o *CommitBlockListOptions) format() *generated.BlockBlobClientCommitBlockListOptions {
	if o == nil {
		return nil
	}
	// TODO figure out request id
	return &generated.BlockBlobClientCommitBlockListOptions{
		BlobTagsString:            shared.SerializeBlobTagsToStrPtr(o.Tags),
		Metadata:                  o.Metadata,
		RequestID:                 o.RequestID,
		Tier:                      o.Tier,
		Timeout:                   o.Timeout,
		TransactionalContentCRC64: o.TransactionalContentCRC64,
		TransactionalContentMD5:   o.TransactionalContentMD5,
		LegalHold:                 o.LegalHold,
		ImmutabilityPolicyMode:    o.ImmutabilityPolicyMode,
		ImmutabilityPolicyExpiry:  o.ImmutabilityPolicyExpiryTime,
		BlobCacheControl:          o.HTTPHeaders.BlobCacheControl,
		BlobContentDisposition:    o.HTTPHeaders.BlobContentDisposition,
		BlobContentEncoding:       o.HTTPHeaders.BlobContentEncoding,
		BlobContentLanguage:       o.HTTPHeaders.BlobContentLanguage,
		BlobContentMD5:            o.HTTPHeaders.BlobContentMD5,
		BlobContentType:           o.HTTPHeaders.BlobContentType,
		LeaseID:                   o.AccessConditions.LeaseAccessConditions.LeaseID,
		IfMatch:                   o.AccessConditions.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:           o.AccessConditions.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:               o.AccessConditions.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince:         o.AccessConditions.ModifiedAccessConditions.IfUnmodifiedSince,
		EncryptionAlgorithm:       o.CPKInfo.EncryptionAlgorithm,
		EncryptionKey:             o.CPKInfo.EncryptionKey,
		EncryptionKeySHA256:       o.CPKInfo.EncryptionKeySHA256,
		EncryptionScope:           o.CPKScopeInfo.EncryptionScope,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// GetBlockListOptions contains the optional parameters for the Client.GetBlockList method.
type GetBlockListOptions struct {
	Snapshot         *string
	AccessConditions *blob.AccessConditions
}

func (o *GetBlockListOptions) format() *generated.BlockBlobClientGetBlockListOptions {
	if o == nil {
		return nil
	}
	// Note: missing mapping for o.AccessConditions.ModifiedAccessConditions
	return &generated.BlockBlobClientGetBlockListOptions{
		Snapshot: o.Snapshot,
		LeaseID:  o.AccessConditions.LeaseAccessConditions.LeaseID,
	}
}

// ------------------------------------------------------------

// uploadFromReaderOptions identifies options used by the UploadBuffer and UploadFile functions.
type uploadFromReaderOptions struct {
	// BlockSize specifies the block size to use; the default (and maximum size) is MaxStageBlockBytes.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are sent to the BlockBlobClient.
	// Note that the progress reporting is not always increasing; it can go down when retrying a request.
	Progress func(bytesTransferred int64)

	// HTTPHeaders indicates the HTTP headers to be associated with the blob.
	HTTPHeaders *blob.HTTPHeaders

	// Metadata indicates the metadata to be associated with the blob when PutBlockList is called.
	Metadata map[string]*string

	// AccessConditions indicates the access conditions for the block blob.
	AccessConditions *blob.AccessConditions

	// AccessTier indicates the tier of blob
	AccessTier *blob.AccessTier

	// BlobTags
	Tags map[string]string

	// ClientProvidedKeyOptions indicates the client provided key by name and/or by value to encrypt/decrypt data.
	CPKInfo      *blob.CPKInfo
	CPKScopeInfo *blob.CPKScopeInfo

	// Concurrency indicates the maximum number of blocks to upload in parallel (0=default)
	Concurrency uint16

	TransactionalValidation blob.TransferValidationType

	// Deprecated: TransactionalContentCRC64 cannot be generated at block level
	TransactionalContentCRC64 uint64

	// Deprecated: TransactionalContentMD5 cannot be generated at block level
	TransactionalContentMD5 []byte
}

// UploadBufferOptions provides set of configurations for UploadBuffer operation.
type UploadBufferOptions = uploadFromReaderOptions

// UploadFileOptions provides set of configurations for UploadFile operation.
type UploadFileOptions = uploadFromReaderOptions

func (o *uploadFromReaderOptions) getStageBlockOptions() *StageBlockOptions {
	leaseAccessConditions, _ := exported.FormatBlobAccessConditions(o.AccessConditions)
	return &StageBlockOptions{
		CPKInfo:               o.CPKInfo,
		CPKScopeInfo:          o.CPKScopeInfo,
		LeaseAccessConditions: leaseAccessConditions,

		TransactionalValidation: o.TransactionalValidation,
	}
}

func (o *uploadFromReaderOptions) getUploadBlockBlobOptions() *UploadOptions {
	return &UploadOptions{
		Tags:             o.Tags,
		Metadata:         o.Metadata,
		Tier:             o.AccessTier,
		HTTPHeaders:      o.HTTPHeaders,
		AccessConditions: o.AccessConditions,
		CPKInfo:          o.CPKInfo,
		CPKScopeInfo:     o.CPKScopeInfo,
	}
}

func (o *uploadFromReaderOptions) getCommitBlockListOptions() *CommitBlockListOptions {
	return &CommitBlockListOptions{
		Tags:         o.Tags,
		Metadata:     o.Metadata,
		Tier:         o.AccessTier,
		HTTPHeaders:  o.HTTPHeaders,
		CPKInfo:      o.CPKInfo,
		CPKScopeInfo: o.CPKScopeInfo,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// UploadStreamOptions provides set of configurations for UploadStream operation.
type UploadStreamOptions struct {
	// BlockSize defines the size of the buffer used during upload. The default and minimum value is 1 MiB.
	BlockSize int64

	// Concurrency defines the max number of concurrent uploads to be performed to upload the file.
	// Each concurrent upload will create a buffer of size BlockSize.  The default value is one.
	Concurrency int

	TransactionalValidation blob.TransferValidationType

	HTTPHeaders      *blob.HTTPHeaders
	Metadata         map[string]*string
	AccessConditions *blob.AccessConditions
	AccessTier       *blob.AccessTier
	Tags             map[string]string
	CPKInfo          *blob.CPKInfo
	CPKScopeInfo     *blob.CPKScopeInfo
}

func (u *UploadStreamOptions) setDefaults() {
	if u.Concurrency == 0 {
		u.Concurrency = 1
	}

	if u.BlockSize < _1MiB {
		u.BlockSize = _1MiB
	}
}

func (u *UploadStreamOptions) getStageBlockOptions() *StageBlockOptions {
	if u == nil {
		return nil
	}

	leaseAccessConditions, _ := exported.FormatBlobAccessConditions(u.AccessConditions)
	return &StageBlockOptions{
		TransactionalValidation: u.TransactionalValidation,
		CPKInfo:                 u.CPKInfo,
		CPKScopeInfo:            u.CPKScopeInfo,
		LeaseAccessConditions:   leaseAccessConditions,
	}
}

func (u *UploadStreamOptions) getCommitBlockListOptions() *CommitBlockListOptions {
	if u == nil {
		return nil
	}

	return &CommitBlockListOptions{
		Tags:             u.Tags,
		Metadata:         u.Metadata,
		Tier:             u.AccessTier,
		HTTPHeaders:      u.HTTPHeaders,
		CPKInfo:          u.CPKInfo,
		CPKScopeInfo:     u.CPKScopeInfo,
		AccessConditions: u.AccessConditions,
	}
}

func (u *UploadStreamOptions) getUploadOptions() *UploadOptions {
	if u == nil {
		return nil
	}

	return &UploadOptions{
		Tags:             u.Tags,
		Metadata:         u.Metadata,
		Tier:             u.AccessTier,
		HTTPHeaders:      u.HTTPHeaders,
		CPKInfo:          u.CPKInfo,
		CPKScopeInfo:     u.CPKScopeInfo,
		AccessConditions: u.AccessConditions,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// ExpiryType defines values for ExpiryType.
type ExpiryType = exported.ExpiryType

// ExpiryTypeAbsolute defines the absolute time for the blob expiry.
type ExpiryTypeAbsolute = exported.ExpiryTypeAbsolute

// ExpiryTypeRelativeToNow defines the duration relative to now for the blob expiry.
type ExpiryTypeRelativeToNow = exported.ExpiryTypeRelativeToNow

// ExpiryTypeRelativeToCreation defines the duration relative to creation for the blob expiry.
type ExpiryTypeRelativeToCreation = exported.ExpiryTypeRelativeToCreation

// ExpiryTypeNever defines that the blob will be set to never expire.
type ExpiryTypeNever = exported.ExpiryTypeNever

// SetExpiryOptions contains the optional parameters for the Client.SetExpiry method.
type SetExpiryOptions = exported.SetExpiryOptions
