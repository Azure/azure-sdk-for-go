// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

type UploadBlockBlobOptions struct {
	// Optional. Used to set blob tags in various blob operations.
	BlobTagsMap *map[string]string

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]string

	// Optional. Indicates the tier to be set on the blob.
	Tier *AccessTier

	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 []byte

	BlobHTTPHeaders          *BlobHTTPHeaders
	LeaseAccessConditions    *LeaseAccessConditions
	CpkInfo                  *CpkInfo
	CpkScopeInfo             *CpkScopeInfo
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *UploadBlockBlobOptions) pointers() (*BlockBlobUploadOptions, *BlobHTTPHeaders, *LeaseAccessConditions,
	*CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	basics := BlockBlobUploadOptions{
		BlobTagsString:          serializeBlobTagsToStrPtr(o.BlobTagsMap),
		Metadata:                o.Metadata,
		Tier:                    o.Tier,
		TransactionalContentMD5: o.TransactionalContentMD5,
	}

	return &basics, o.BlobHTTPHeaders, o.LeaseAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions
}

type StageBlockOptions struct {
	CpkInfo                    *CpkInfo
	CpkScopeInfo               *CpkScopeInfo
	LeaseAccessConditions      *LeaseAccessConditions
	BlockBlobStageBlockOptions *BlockBlobStageBlockOptions
}

func (o *StageBlockOptions) pointers() (*LeaseAccessConditions, *BlockBlobStageBlockOptions, *CpkInfo, *CpkScopeInfo) {
	if o == nil {
		return nil, nil, nil, nil
	}

	return o.LeaseAccessConditions, o.BlockBlobStageBlockOptions, o.CpkInfo, o.CpkScopeInfo
}

type StageBlockFromURLOptions struct {
	LeaseAccessConditions          *LeaseAccessConditions
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage analytics logging is enabled.
	RequestID *string
	// Specify the md5 calculated for the range of bytes that must be read from the copy source.
	SourceContentMD5 []byte
	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentcrc64 []byte

	Offset *int64

	Count *int64
	// The timeout parameter is expressed in seconds.
	Timeout *int32

	CpkInfo      *CpkInfo
	CpkScopeInfo *CpkScopeInfo
}

func (o *StageBlockFromURLOptions) pointers() (*LeaseAccessConditions, *SourceModifiedAccessConditions, *BlockBlobStageBlockFromURLOptions, *CpkInfo, *CpkScopeInfo) {
	if o == nil {
		return nil, nil, nil, nil, nil
	}

	options := &BlockBlobStageBlockFromURLOptions{
		RequestID:          o.RequestID,
		SourceContentMD5:   o.SourceContentMD5,
		SourceContentcrc64: o.SourceContentcrc64,
		SourceRange:        getSourceRange(o.Offset, o.Count),
		Timeout:            o.Timeout,
	}

	return o.LeaseAccessConditions, o.SourceModifiedAccessConditions, options, o.CpkInfo, o.CpkScopeInfo
}

type CommitBlockListOptions struct {
	BlobTagsMap               *map[string]string
	Metadata                  map[string]string
	RequestID                 *string
	Tier                      *AccessTier
	Timeout                   *int32
	TransactionalContentCRC64 []byte
	TransactionalContentMD5   []byte
	BlobHTTPHeaders           *BlobHTTPHeaders
	CpkInfo                   *CpkInfo
	CpkScopeInfo              *CpkScopeInfo
	BlobAccessConditions
}

func (o *CommitBlockListOptions) pointers() (*BlockBlobCommitBlockListOptions, *BlobHTTPHeaders, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	options := &BlockBlobCommitBlockListOptions{
		BlobTagsString:            serializeBlobTagsToStrPtr(o.BlobTagsMap),
		Metadata:                  o.Metadata,
		RequestID:                 o.RequestID,
		Tier:                      o.Tier,
		Timeout:                   o.Timeout,
		TransactionalContentCRC64: o.TransactionalContentCRC64,
		TransactionalContentMD5:   o.TransactionalContentMD5,
	}

	return options, o.BlobHTTPHeaders, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions, o.LeaseAccessConditions
}

type GetBlockListOptions struct {
	BlockBlobGetBlockListOptions *BlockBlobGetBlockListOptions
	BlobAccessConditions
}

func (o *GetBlockListOptions) pointers() (*BlockBlobGetBlockListOptions, *ModifiedAccessConditions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	return o.BlockBlobGetBlockListOptions, o.ModifiedAccessConditions, o.LeaseAccessConditions
}

type CopyBlockBlobFromURLOptions struct {
	BlobTagsMap                    *map[string]string
	Metadata                       map[string]string
	RequestID                      *string
	SourceContentMD5               []byte
	Tier                           *AccessTier
	Timeout                        *int32
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	BlobAccessConditions
}

func (o *CopyBlockBlobFromURLOptions) pointers() (*BlobCopyFromURLOptions, *SourceModifiedAccessConditions, *ModifiedAccessConditions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil
	}

	options := &BlobCopyFromURLOptions{
		BlobTagsString:   serializeBlobTagsToStrPtr(o.BlobTagsMap),
		Metadata:         o.Metadata,
		RequestID:        o.RequestID,
		SourceContentMD5: o.SourceContentMD5,
		Tier:             o.Tier,
		Timeout:          o.Timeout,
	}

	return options, o.SourceModifiedAccessConditions, o.ModifiedAccessConditions, o.LeaseAccessConditions
}
