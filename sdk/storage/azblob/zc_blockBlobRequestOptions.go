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
		BlobTagsString:          SerializeBlobTags(o.BlobTagsMap),
		Metadata:                o.Metadata,
		Tier:                    o.Tier,
		TransactionalContentMd5: o.TransactionalContentMd5,
	}

	return &basics, o.BlobHttpHeaders, o.LeaseAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions
}

type StageBlockOptions struct {
	LeaseAccessConditions      *LeaseAccessConditions
	BlockBlobStageBlockOptions *BlockBlobStageBlockOptions
}

func (o *StageBlockOptions) pointers() (*LeaseAccessConditions, *BlockBlobStageBlockOptions) {
	if o == nil {
		return nil, nil
	}

	return o.LeaseAccessConditions, o.BlockBlobStageBlockOptions
}

type StageBlockFromURLOptions struct {
	LeaseAccessConditions          *LeaseAccessConditions
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	BlockBlobStageBlockOptions     *BlockBlobStageBlockFromURLOptions
	CpkInfo                        *CpkInfo
	CpkScopeInfo                   *CpkScopeInfo
}

func (o *StageBlockFromURLOptions) pointers() (*LeaseAccessConditions, *SourceModifiedAccessConditions, *BlockBlobStageBlockFromURLOptions, *CpkInfo, *CpkScopeInfo) {
	if o == nil {
		return nil, nil, nil, nil, nil
	}

	return o.LeaseAccessConditions, o.SourceModifiedAccessConditions, o.BlockBlobStageBlockOptions, o.CpkInfo, o.CpkScopeInfo
}

type CommitBlockListOptions struct {
	BlockBlobCommitBlockListOptions *BlockBlobCommitBlockListOptions
	BlobHTTPHeaders                 *BlobHttpHeaders
	CpkInfo                         *CpkInfo
	CpkScope                        *CpkScopeInfo
	BlobAccessConditions
}

func (o *CommitBlockListOptions) pointers() (*BlockBlobCommitBlockListOptions, *BlobHttpHeaders, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	return o.BlockBlobCommitBlockListOptions, o.BlobHTTPHeaders, o.CpkInfo, o.CpkScope, o.ModifiedAccessConditions, o.LeaseAccessConditions
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
	BlobCopyFromURLOptions         *BlobCopyFromURLOptions
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	BlobAccessConditions
}

func (o *CopyBlockBlobFromURLOptions) pointers() (*BlobCopyFromURLOptions, *SourceModifiedAccessConditions, *ModifiedAccessConditions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil
	}

	return o.BlobCopyFromURLOptions, o.SourceModifiedAccessConditions, o.ModifiedAccessConditions, o.LeaseAccessConditions
}
