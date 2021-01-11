package azblob

type UploadBlockBlobOptions struct {
	// Optional. Used to set blob tags in various blob operations.
	BlobTagsMap *map[string]string

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata *map[string]string

	// Optional. Indicates the tier to be set on the blob.
	Tier *AccessTier

	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMd5 *[]byte

	BlobHttpHeaders          *BlobHttpHeaders
	LeaseAccessConditions    *LeaseAccessConditions
	CpkInfo                  *CpkInfo
	CpkScopeInfo             *CpkScopeInfo
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *UploadBlockBlobOptions) pointers() (*BlockBlobUploadOptions, *BlobHttpHeaders, *LeaseAccessConditions,
	*CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	basics := BlockBlobUploadOptions{
		BlobTagsString:          SerializeBlobTagsToStrPtr(o.BlobTagsMap),
		Metadata:                o.Metadata,
		Tier:                    o.Tier,
		TransactionalContentMd5: o.TransactionalContentMd5,
	}

	return &basics, o.BlobHttpHeaders, o.LeaseAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions
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
	RequestId *string
	// Specify the md5 calculated for the range of bytes that must be read from the copy source.
	SourceContentMd5 *[]byte
	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentcrc64 *[]byte

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
		RequestId:          o.RequestId,
		SourceContentMd5:   o.SourceContentMd5,
		SourceContentcrc64: o.SourceContentcrc64,
		SourceRange:        getSourceRange(o.Offset, o.Count),
		Timeout:            o.Timeout,
	}

	return o.LeaseAccessConditions, o.SourceModifiedAccessConditions, options, o.CpkInfo, o.CpkScopeInfo
}

type CommitBlockListOptions struct {
	BlobTagsMap               *map[string]string
	Metadata                  *map[string]string
	RequestId                 *string
	Tier                      *AccessTier
	Timeout                   *int32
	TransactionalContentCrc64 *[]byte
	TransactionalContentMd5   *[]byte
	BlobHTTPHeaders           *BlobHttpHeaders
	CpkInfo                   *CpkInfo
	CpkScope                  *CpkScopeInfo
	BlobAccessConditions
}

func (o *CommitBlockListOptions) pointers() (*BlockBlobCommitBlockListOptions, *BlobHttpHeaders, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	options := &BlockBlobCommitBlockListOptions{
		BlobTagsString:            SerializeBlobTagsToStrPtr(o.BlobTagsMap),
		Metadata:                  o.Metadata,
		RequestId:                 o.RequestId,
		Tier:                      o.Tier,
		Timeout:                   o.Timeout,
		TransactionalContentCrc64: o.TransactionalContentCrc64,
		TransactionalContentMd5:   o.TransactionalContentMd5,
	}

	return options, o.BlobHTTPHeaders, o.CpkInfo, o.CpkScope, o.ModifiedAccessConditions, o.LeaseAccessConditions
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
	Metadata                       *map[string]string
	RequestId                      *string
	SourceContentMd5               *[]byte
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
		BlobTagsString:   SerializeBlobTagsToStrPtr(o.BlobTagsMap),
		Metadata:         o.Metadata,
		RequestId:        o.RequestId,
		SourceContentMd5: o.SourceContentMd5,
		Tier:             o.Tier,
		Timeout:          o.Timeout,
	}

	return options, o.SourceModifiedAccessConditions, o.ModifiedAccessConditions, o.LeaseAccessConditions
}
