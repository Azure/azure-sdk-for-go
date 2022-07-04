//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
)

type AccessTier = generated.AccessTier

// ImmutabilityPolicyMode enum
type ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyMode

type UploadOptions struct {
	// Optional. Used to set blob tags in various blob operations.
	Tags map[string]string

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]string

	// Optional. Indicates the tier to be set on the blob.
	Tier *AccessTier

	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 []byte

	HTTPHeaders      *blob.HTTPHeaders
	CpkInfo          *CpkInfo
	CpkScopeInfo     *CpkScopeInfo
	AccessConditions *AccessConditions
}

type StageBlockOptions struct {
	CpkInfo *CpkInfo

	CpkScopeInfo *CpkScopeInfo

	LeaseAccessConditions *lease.AccessConditions

	// Specify the transactional crc64 for the body, to be validated by the service.
	TransactionalContentCRC64 []byte

	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 []byte
}

type SourceModifiedAccessConditions = generated.SourceModifiedAccessConditions

// StageBlockFromURLOptions provides set of configurations for StageBlockFromURL operation
type StageBlockFromURLOptions struct {
	// Only Bearer type is supported. Credentials should be a valid OAuth access token to copy source.
	CopySourceAuthorization *string

	LeaseAccessConditions *lease.AccessConditions

	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	// Specify the md5 calculated for the range of bytes that must be read from the copy source.
	SourceContentMD5 []byte
	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCRC64 []byte

	Offset *int64

	Count *int64

	CpkInfo *CpkInfo

	CpkScopeInfo *CpkScopeInfo
}

func (o *StageBlockFromURLOptions) format() (*generated.BlockBlobClientStageBlockFromURLOptions, *generated.CpkInfo, *generated.CpkScopeInfo, *generated.LeaseAccessConditions, *generated.SourceModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil
	}

	options := &generated.BlockBlobClientStageBlockFromURLOptions{
		CopySourceAuthorization: o.CopySourceAuthorization,
		SourceContentMD5:        o.SourceContentMD5,
		SourceContentcrc64:      o.SourceContentCRC64,
		SourceRange:             shared.GetSourceRange(o.Offset, o.Count),
	}

	return options, o.CpkInfo, o.CpkScopeInfo, o.LeaseAccessConditions, o.SourceModifiedAccessConditions
}

type CommitBlockListOptions struct {
	Tags                      map[string]string
	Metadata                  map[string]string
	RequestID                 *string
	Tier                      *AccessTier
	Timeout                   *int32
	TransactionalContentCRC64 []byte
	TransactionalContentMD5   []byte
	BlobHTTPHeaders           *blob.HTTPHeaders
	CpkInfo                   *CpkInfo
	CpkScopeInfo              *CpkScopeInfo
	BlobAccessConditions      *AccessConditions
}

type CpkInfo = generated.CpkInfo

type CpkScopeInfo = generated.CpkScopeInfo

type AccessConditions = exported.BlobAccessConditions

// GetBlockListOptions provides set of configurations for GetBlockList operation
type GetBlockListOptions struct {
	Snapshot             *string
	BlobAccessConditions *AccessConditions
}

func (o *GetBlockListOptions) format() (*generated.BlockBlobClientGetBlockListOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatBlobAccessConditions(o.BlobAccessConditions)
	return &generated.BlockBlobClientGetBlockListOptions{Snapshot: o.Snapshot}, leaseAccessConditions, modifiedAccessConditions
}
