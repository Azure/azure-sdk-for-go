//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"time"
)

const _1MiB = 1024 * 1024

// UploadOption identifies options used by the UploadBuffer and UploadFile functions.
type UploadOption struct {
	// BlockSize specifies the block size to use; the default (and maximum size) is BlockBlobMaxStageBlockBytes.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are sent to the BlockBlobClient.
	// Note that the progress reporting is not always increasing; it can go down when retrying a request.
	Progress func(bytesTransferred int64)

	// HTTPHeaders indicates the HTTP headers to be associated with the blob.
	HTTPHeaders *BlobHTTPHeaders

	// Metadata indicates the metadata to be associated with the blob when PutBlockList is called.
	Metadata map[string]string

	// BlobAccessConditions indicates the access conditions for the block blob.
	BlobAccessConditions *BlobAccessConditions

	// AccessTier indicates the tier of blob
	AccessTier *AccessTier

	// BlobTags
	BlobTags map[string]string

	// ClientProvidedKeyOptions indicates the client provided key by name and/or by value to encrypt/decrypt data.
	CpkInfo      *CpkInfo
	CpkScopeInfo *CpkScopeInfo

	// Parallelism indicates the maximum number of blocks to upload in parallel (0=default)
	Parallelism uint16
	// Optional header, Specifies the transactional crc64 for the body, to be validated by the service.
	TransactionalContentCRC64 *[]byte
	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMD5 *[]byte
}

type UploadResponse struct {
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// ContentMD5 contains the information returned from the Content-MD5 header response.
	ContentMD5 []byte

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// ETag contains the information returned from the ETag header response.
	ETag *string

	// EncryptionKeySHA256 contains the information returned from the x-ms-encryption-key-sha256 header response.
	EncryptionKeySHA256 *string

	// EncryptionScope contains the information returned from the x-ms-encryption-scope header response.
	EncryptionScope *string

	// IsServerEncrypted contains the information returned from the x-ms-request-server-encrypted header response.
	IsServerEncrypted *bool

	// LastModified contains the information returned from the Last-Modified header response.
	LastModified *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string

	// VersionID contains the information returned from the x-ms-version-id header response.
	VersionID *string

	// XMSContentCRC64 contains the information returned from the x-ms-content-crc64 header response.
	// Will be a part of response only if uploading data >= internal.BlockBlobMaxUploadBlobBytes (= 256 * 1024 * 1024 // 256MB)
	XMSContentCRC64 []byte
}

func toUploadResponseFromBlockBlobUploadResponse(resp BlockBlobUploadResponse) UploadResponse {
	return UploadResponse{
		ClientRequestID:     resp.ClientRequestID,
		ContentMD5:          resp.ContentMD5,
		Date:                resp.Date,
		ETag:                resp.ETag,
		EncryptionKeySHA256: resp.EncryptionKeySHA256,
		EncryptionScope:     resp.EncryptionScope,
		IsServerEncrypted:   resp.IsServerEncrypted,
		LastModified:        resp.LastModified,
		RequestID:           resp.RequestID,
		Version:             resp.Version,
		VersionID:           resp.VersionID,
	}
}

func toUploadResponseFromBlockBlobCommitBlockListResponse(resp BlockBlobCommitBlockListResponse) UploadResponse {
	return UploadResponse{
		ClientRequestID:     resp.ClientRequestID,
		ContentMD5:          resp.ContentMD5,
		Date:                resp.Date,
		ETag:                resp.ETag,
		EncryptionKeySHA256: resp.EncryptionKeySHA256,
		EncryptionScope:     resp.EncryptionScope,
		IsServerEncrypted:   resp.IsServerEncrypted,
		LastModified:        resp.LastModified,
		RequestID:           resp.RequestID,
		Version:             resp.Version,
		VersionID:           resp.VersionID,
		XMSContentCRC64:     resp.XMSContentCRC64,
	}
}

func (o *UploadOption) getStageBlockOptions() *BlockBlobStageBlockOptions {
	leaseAccessConditions, _ := o.BlobAccessConditions.format()
	return &BlockBlobStageBlockOptions{
		CpkInfo:               o.CpkInfo,
		CpkScopeInfo:          o.CpkScopeInfo,
		LeaseAccessConditions: leaseAccessConditions,
	}
}

func (o *UploadOption) getUploadBlockBlobOptions() *BlockBlobUploadOptions {
	return &BlockBlobUploadOptions{
		BlobTags:             o.BlobTags,
		Metadata:             o.Metadata,
		Tier:                 o.AccessTier,
		HTTPHeaders:          o.HTTPHeaders,
		BlobAccessConditions: o.BlobAccessConditions,
		CpkInfo:              o.CpkInfo,
		CpkScopeInfo:         o.CpkScopeInfo,
	}
}

func (o *UploadOption) getCommitBlockListOptions() *BlockBlobCommitBlockListOptions {
	return &BlockBlobCommitBlockListOptions{
		BlobTags:        o.BlobTags,
		Metadata:        o.Metadata,
		Tier:            o.AccessTier,
		BlobHTTPHeaders: o.HTTPHeaders,
		CpkInfo:         o.CpkInfo,
		CpkScopeInfo:    o.CpkScopeInfo,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// UploadStreamOptions provides set of configurations for UploadStream operation
type UploadStreamOptions struct {
	// TransferManager provides a TransferManager that controls buffer allocation/reuse and
	// concurrency. This overrides BufferSize and MaxBuffers if set.
	TransferManager      internal.TransferManager
	transferMangerNotSet bool
	// BufferSize sizes the buffer used to read data from source. If < 1 MiB, format to 1 MiB.
	BufferSize int
	// MaxBuffers defines the number of simultaneous uploads will be performed to upload the file.
	MaxBuffers           int
	HTTPHeaders          *BlobHTTPHeaders
	Metadata             map[string]string
	BlobAccessConditions *BlobAccessConditions
	AccessTier           *AccessTier
	BlobTags             map[string]string
	CpkInfo              *CpkInfo
	CpkScopeInfo         *CpkScopeInfo
}

func (u *UploadStreamOptions) format() error {
	if u == nil || u.TransferManager != nil {
		return nil
	}

	if u.MaxBuffers == 0 {
		u.MaxBuffers = 1
	}

	if u.BufferSize < _1MiB {
		u.BufferSize = _1MiB
	}

	var err error
	u.TransferManager, err = internal.NewStaticBuffer(u.BufferSize, u.MaxBuffers)
	if err != nil {
		return fmt.Errorf("bug: default transfer manager could not be created: %s", err)
	}
	u.transferMangerNotSet = true
	return nil
}

func (u *UploadStreamOptions) getStageBlockOptions() *BlockBlobStageBlockOptions {
	leaseAccessConditions, _ := u.BlobAccessConditions.format()
	return &BlockBlobStageBlockOptions{
		CpkInfo:               u.CpkInfo,
		CpkScopeInfo:          u.CpkScopeInfo,
		LeaseAccessConditions: leaseAccessConditions,
	}
}

func (u *UploadStreamOptions) getCommitBlockListOptions() *BlockBlobCommitBlockListOptions {
	options := &BlockBlobCommitBlockListOptions{
		BlobTags:             u.BlobTags,
		Metadata:             u.Metadata,
		Tier:                 u.AccessTier,
		BlobHTTPHeaders:      u.HTTPHeaders,
		CpkInfo:              u.CpkInfo,
		CpkScopeInfo:         u.CpkScopeInfo,
		BlobAccessConditions: u.BlobAccessConditions,
	}

	return options
}

// ---------------------------------------------------------------------------------------------------------------------

// DownloadOptions identifies options used by the DownloadToBuffer and DownloadToFile functions.
type DownloadOptions struct {
	// BlockSize specifies the block size to use for each parallel download; the default size is BlobDefaultDownloadBlockSize.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	Progress func(bytesTransferred int64)

	// BlobAccessConditions indicates the access conditions used when making HTTP GET requests against the blob.
	BlobAccessConditions *BlobAccessConditions

	// ClientProvidedKeyOptions indicates the client provided key by name and/or by value to encrypt/decrypt data.
	CpkInfo      *CpkInfo
	CpkScopeInfo *CpkScopeInfo

	// Parallelism indicates the maximum number of blocks to download in parallel (0=default)
	Parallelism uint16

	// RetryReaderOptionsPerBlock is used when downloading each block.
	RetryReaderOptionsPerBlock RetryReaderOptions
}

func (o *DownloadOptions) getBlobPropertiesOptions() *BlobGetPropertiesOptions {
	return &BlobGetPropertiesOptions{
		BlobAccessConditions: o.BlobAccessConditions,
		CpkInfo:              o.CpkInfo,
	}
}

func (o *DownloadOptions) getDownloadBlobOptions(offSet, count int64, rangeGetContentMD5 *bool) *BlobDownloadOptions {
	return &BlobDownloadOptions{
		BlobAccessConditions: o.BlobAccessConditions,
		CpkInfo:              o.CpkInfo,
		CpkScopeInfo:         o.CpkScopeInfo,
		Offset:               &offSet,
		Count:                &count,
		RangeGetContentMD5:   rangeGetContentMD5,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// batchTransferOptions identifies options used by doBatchTransfer.
type batchTransferOptions struct {
	TransferSize  int64
	ChunkSize     int64
	Parallelism   uint16
	Operation     func(offset int64, chunkSize int64, ctx context.Context) error
	OperationName string
}
