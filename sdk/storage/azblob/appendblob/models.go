//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package appendblob

import (
	"encoding/binary"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/hashing"
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

	CpkInfo *blob.CpkInfo

	CpkScopeInfo *blob.CpkScopeInfo

	// Optional. Used to set blob tags in various blob operations.
	Tags map[string]string

	// Optional. Specifies a user-defined name-value pair associated with the blob. If no name-value pairs are specified, the
	// operation will copy the metadata from the source blob or file to the destination blob. If one or more name-value pairs
	// are specified, the destination blob is created with the specified metadata, and metadata is not copied from the source
	// blob or file. Note that beginning with version 2009-09-19, metadata names must adhere to the naming rules for C# identifiers.
	// See Naming and Referencing Containers, Blobs, and Metadata for more information.
	Metadata map[string]string
}

func (o *CreateOptions) format() (*generated.AppendBlobClientCreateOptions, *generated.BlobHTTPHeaders, *generated.LeaseAccessConditions, *generated.CpkInfo, *generated.CpkScopeInfo, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	options := generated.AppendBlobClientCreateOptions{
		BlobTagsString:           shared.SerializeBlobTagsToStrPtr(o.Tags),
		Metadata:                 o.Metadata,
		ImmutabilityPolicyExpiry: o.ImmutabilityPolicyExpiry,
		ImmutabilityPolicyMode:   o.ImmutabilityPolicyMode,
		LegalHold:                o.LegalHold,
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return &options, o.HTTPHeaders, leaseAccessConditions, o.CpkInfo, o.CpkScopeInfo, modifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// AppendBlockOptions contains the optional parameters for the Client.AppendBlock method.
type AppendBlockOptions struct {
	// Let the SDK hash for you (providing the hash type(s) specified voids their relevance in this option).
	// nil = None default
	TransactionalValidationOption hashing.TransferValidationType

	// Specify the transactional crc64 for the body, to be validated by the service. Should be hashed using hashing.CRC64Table or hashing.CRC64Polynomial
	TransactionalContentCRC64 uint64
	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 []byte

	AppendPositionAccessConditions *AppendPositionAccessConditions

	CpkInfo *blob.CpkInfo

	CpkScopeInfo *blob.CpkScopeInfo

	AccessConditions *blob.AccessConditions
}

func (o *AppendBlockOptions) format() (*generated.AppendBlobClientAppendBlockOptions, *generated.AppendPositionAccessConditions,
	*generated.CpkInfo, *generated.CpkScopeInfo, *generated.ModifiedAccessConditions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	options := &generated.AppendBlobClientAppendBlockOptions{
		TransactionalContentMD5: o.TransactionalContentMD5,
	}

	if o.TransactionalValidationOption == hashing.TransferValidationTypeCRC64 || o.TransactionalContentCRC64 != 0 {
		options.TransactionalContentCRC64 = make([]byte, 8)
		// If the validation option is specified & the CRC is 0, it will be 0, and get overwritten anyway with the new hash.
		// Thus, it's OK to run putUint64 here.
		binary.LittleEndian.PutUint64(options.TransactionalContentCRC64, o.TransactionalContentCRC64)
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return options, o.AppendPositionAccessConditions, o.CpkInfo, o.CpkScopeInfo, modifiedAccessConditions, leaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// AppendBlockFromURLOptions contains the optional parameters for the Client.AppendBlockFromURL method.
type AppendBlockFromURLOptions struct {
	// Specify the md5 calculated for the range of bytes that must be read from the copy source.
	SourceContentMD5 []byte
	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCRC64 []byte
	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 []byte

	AppendPositionAccessConditions *AppendPositionAccessConditions

	CpkInfo *blob.CpkInfo

	CpkScopeInfo *blob.CpkScopeInfo

	SourceModifiedAccessConditions *blob.SourceModifiedAccessConditions

	AccessConditions *blob.AccessConditions

	// Range specifies a range of bytes.  The default value is all bytes.
	Range blob.HTTPRange
}

func (o *AppendBlockFromURLOptions) format() (*generated.AppendBlobClientAppendBlockFromURLOptions, *generated.CpkInfo,
	*generated.CpkScopeInfo, *generated.LeaseAccessConditions, *generated.AppendPositionAccessConditions,
	*generated.ModifiedAccessConditions, *generated.SourceModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil, nil
	}

	options := &generated.AppendBlobClientAppendBlockFromURLOptions{
		SourceRange:             exported.FormatHTTPRange(o.Range),
		SourceContentMD5:        o.SourceContentMD5,
		SourceContentcrc64:      o.SourceContentCRC64,
		TransactionalContentMD5: o.TransactionalContentMD5,
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return options, o.CpkInfo, o.CpkScopeInfo, leaseAccessConditions, o.AppendPositionAccessConditions, modifiedAccessConditions, o.SourceModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// SealOptions provides set of configurations for SealAppendBlob operation
type SealOptions struct {
	AccessConditions               *blob.AccessConditions
	AppendPositionAccessConditions *AppendPositionAccessConditions
}

func (o *SealOptions) format() (*generated.LeaseAccessConditions,
	*generated.ModifiedAccessConditions, *generated.AppendPositionAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.AccessConditions)
	return leaseAccessConditions, modifiedAccessConditions, o.AppendPositionAccessConditions

}

// ---------------------------------------------------------------------------------------------------------------------
