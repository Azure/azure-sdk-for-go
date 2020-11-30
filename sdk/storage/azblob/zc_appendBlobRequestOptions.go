package azblob

type AppendBlockOptions struct {
	// Specify the transactional crc64 for the body, to be validated by the service.
	TransactionalContentCrc64 *[]byte
	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMd5 *[]byte

	AppendPositionAccessConditions *AppendPositionAccessConditions
	CpkInfo *CpkInfo
	CpkScopeInfo *CpkScopeInfo
	BlobAccessConditions
}

func (o *AppendBlockOptions) pointers() (*AppendBlobAppendBlockOptions, *AppendPositionAccessConditions, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions, *LeaseAccessConditions){
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
	CpkInfo *CpkInfo
	CpkScopeInfo *CpkScopeInfo
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	BlobAccessConditions
}

func (o *AppendBlockURLOptions) pointers() (*AppendBlobAppendBlockFromURLOptions, *AppendPositionAccessConditions, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions, *LeaseAccessConditions, *SourceModifiedAccessConditions){
	if o == nil {
		return nil, nil, nil, nil, nil, nil, nil
	}

	options := &AppendBlobAppendBlockFromURLOptions{
		SourceContentMd5: o.SourceContentMd5,
		SourceContentcrc64: o.SourceContentCrc64,
		TransactionalContentMd5:   o.TransactionalContentMd5,
	}

	return options, o.AppendPositionAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions, o.LeaseAccessConditions, o.SourceModifiedAccessConditions
}