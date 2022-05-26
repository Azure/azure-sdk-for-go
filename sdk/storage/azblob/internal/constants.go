//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

var SASVersion = "2019-12-12"

const (
	ModuleName    = "azblob"
	ModuleVersion = "v0.4.1"
)

//nolint
const (
	// BlockBlobMaxUploadBlobBytes indicates the maximum number of bytes that can be sent in a call to Upload.
	BlockBlobMaxUploadBlobBytes = 256 * 1024 * 1024 // 256MB

	// BlockBlobMaxStageBlockBytes indicates the maximum number of bytes that can be sent in a call to StageBlock.
	BlockBlobMaxStageBlockBytes = 4000 * 1024 * 1024 // 4GB

	// BlockBlobMaxBlocks indicates the maximum number of blocks allowed in a block blob.
	BlockBlobMaxBlocks = 50000

	// PageBlobPageBytes indicates the number of bytes in a page (512).
	PageBlobPageBytes = 512

	// BlobDefaultDownloadBlockSize is default block size
	BlobDefaultDownloadBlockSize = int64(4 * 1024 * 1024) // 4MB
)

const (
	HeaderAuthorization     = "Authorization"
	HeaderXmsDate           = "x-ms-date"
	HeaderContentLength     = "Content-Length"
	HeaderContentEncoding   = "Content-Encoding"
	HeaderContentLanguage   = "Content-Language"
	HeaderContentType       = "Content-Type"
	HeaderContentMD5        = "Content-MD5"
	HeaderIfModifiedSince   = "If-Modified-Since"
	HeaderIfMatch           = "If-Match"
	HeaderIfNoneMatch       = "If-None-Match"
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	HeaderRange             = "Range"
)

const (
	TokenScope = "https://storage.azure.com/.default"
)
