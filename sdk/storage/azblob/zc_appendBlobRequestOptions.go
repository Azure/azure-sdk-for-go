package azblob

type CreateAppendBlobOptions struct {
	BlobAccessConditions

	BlobHttpHeaders *BlobHttpHeaders

	CpkInfo *CpkInfo

	CpkScopeInfo *CpkScopeInfo
	// Optional. Used to set blob tags in various blob operations.
	BlobTagsMap *map[string]string
	// Optional. Specifies a user-defined name-value pair associated with the blob. If no name-value pairs are specified, the
	// operation will copy the metadata from the source blob or file to the destination blob. If one or more name-value pairs
	// are specified, the destination blob is created with the specified metadata, and metadata is not copied from the source
	// blob or file. Note that beginning with version 2009-09-19, metadata names must adhere to the naming rules for C# identifiers.
	// See Naming and Referencing Containers, Blobs, and Metadata for more information.
	Metadata *map[string]string
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage analytics logging is enabled.
	RequestId *string

	Timeout *int32
}

func (o *CreateAppendBlobOptions) pointers() (*AppendBlobCreateOptions, *BlobHttpHeaders, *LeaseAccessConditions, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	options := AppendBlobCreateOptions{
		BlobTagsString: SerializeBlobTagsToStrPtr(o.BlobTagsMap),
		Metadata:       o.Metadata,
		RequestId:      o.RequestId,
		Timeout:        o.Timeout,
	}
	return &options, o.BlobHttpHeaders, o.LeaseAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions
}

type AppendBlockOptions struct {
	// Specify the transactional crc64 for the body, to be validated by the service.
	TransactionalContentCrc64 *[]byte
	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMd5 *[]byte

	AppendPositionAccessConditions *AppendPositionAccessConditions
	CpkInfo                        *CpkInfo
	CpkScopeInfo                   *CpkScopeInfo
	BlobAccessConditions
}

func (o *AppendBlockOptions) pointers() (*AppendBlobAppendBlockOptions, *AppendPositionAccessConditions, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	options := &AppendBlobAppendBlockOptions{
		TransactionalContentCrc64: o.TransactionalContentCrc64,
		TransactionalContentMd5:   o.TransactionalContentMd5,
	}

	return options, o.AppendPositionAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions, o.LeaseAccessConditions
}

type AppendBlockURLOptions struct {
	// Specify the md5 calculated for the range of bytes that must be read from the copy source.
	SourceContentMd5 *[]byte
	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCrc64 *[]byte
	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMd5 *[]byte

	AppendPositionAccessConditions *AppendPositionAccessConditions
	CpkInfo                        *CpkInfo
	CpkScopeInfo                   *CpkScopeInfo
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	BlobAccessConditions
	// Optional, you can specify whether a particular range of the blob is read
	Offset *int64
	Count  *int64
}

func (o *AppendBlockURLOptions) pointers() (*AppendBlobAppendBlockFromURLOptions, *AppendPositionAccessConditions, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions, *LeaseAccessConditions, *SourceModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil, nil
	}

	options := &AppendBlobAppendBlockFromURLOptions{
		SourceRange:             httpRange{count: *o.Count, offset: *o.Offset}.pointers(),
		SourceContentMd5:        o.SourceContentMd5,
		SourceContentcrc64:      o.SourceContentCrc64,
		TransactionalContentMd5: o.TransactionalContentMd5,
	}

	return options, o.AppendPositionAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions, o.LeaseAccessConditions, o.SourceModifiedAccessConditions
}
