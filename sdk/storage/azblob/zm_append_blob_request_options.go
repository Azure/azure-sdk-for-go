// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

type CreateAppendBlobOptions struct {
	BlobAccessConditions

	BlobHTTPHeaders *BlobHTTPHeaders

	CpkInfo *CpkInfo

	CpkScopeInfo *CpkScopeInfo
	// Optional. Used to set blob tags in various blob operations.
	BlobTagsMap *map[string]string
	// Optional. Specifies a user-defined name-value pair associated with the blob. If no name-value pairs are specified, the
	// operation will copy the metadata from the source blob or file to the destination blob. If one or more name-value pairs
	// are specified, the destination blob is created with the specified metadata, and metadata is not copied from the source
	// blob or file. Note that beginning with version 2009-09-19, metadata names must adhere to the naming rules for C# identifiers.
	// See Naming and Referencing Containers, Blobs, and Metadata for more information.
	Metadata map[string]string
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage analytics logging is enabled.
	RequestID *string

	Timeout *int32
}

func (o *CreateAppendBlobOptions) pointers() (*AppendBlobCreateOptions, *BlobHTTPHeaders, *LeaseAccessConditions, *CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	options := AppendBlobCreateOptions{
		BlobTagsString: serializeBlobTagsToStrPtr(o.BlobTagsMap),
		Metadata:       o.Metadata,
		RequestID:      o.RequestID,
		Timeout:        o.Timeout,
	}
	return &options, o.BlobHTTPHeaders, o.LeaseAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions
}

type AppendBlockOptions struct {
	// Specify the transactional crc64 for the body, to be validated by the service.
	TransactionalContentCRC64 []byte
	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 []byte

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
		TransactionalContentCRC64: o.TransactionalContentCRC64,
		TransactionalContentMD5:   o.TransactionalContentMD5,
	}

	return options, o.AppendPositionAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions, o.LeaseAccessConditions
}

type AppendBlockURLOptions struct {
	// Specify the md5 calculated for the range of bytes that must be read from the copy source.
	SourceContentMD5 []byte
	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCRC64 []byte
	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 []byte

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
		SourceRange:             getSourceRange(o.Offset, o.Count),
		SourceContentMD5:        o.SourceContentMD5,
		SourceContentcrc64:      o.SourceContentCRC64,
		TransactionalContentMD5: o.TransactionalContentMD5,
	}

	return options, o.AppendPositionAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions, o.LeaseAccessConditions, o.SourceModifiedAccessConditions
}

type SealAppendBlobOptions struct {
	BlobAccessConditions           *BlobAccessConditions
	AppendPositionAccessConditions *AppendPositionAccessConditions
}

func (o *SealAppendBlobOptions) pointers() (leaseAccessConditions *LeaseAccessConditions,
	modifiedAccessConditions *ModifiedAccessConditions, appendPositionAccessConditions *AppendPositionAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	return
}
