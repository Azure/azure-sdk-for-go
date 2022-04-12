//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azfile

const (
	// FileMaxUploadRangeBytes indicates the maximum number of bytes that can be sent in a call to UploadRange.
	FileMaxUploadRangeBytes = 4 * 1024 * 1024 // 4MB

	ISO8601 = "2006-01-02T15:04:05.0000000Z"
)

var (
	// DefaultFilePermissionString is a constant for all intents and purposes.
	// But you can't take the address of a constant string, so it's a variable.
	// Inherit inherits permissions from the parent folder (default when creating files/folders)
	DefaultFilePermissionString = "inherit"

	// DefaultCurrentTimeString sets creation/last write times to now
	DefaultCurrentTimeString = "now"

	// DefaultPreserveString preserves old permissions on the file/folder (default when updating properties)
	DefaultPreserveString = "preserve"

	// DefaultFileAttributes is defaults for file attributes
	DefaultFileAttributes = "None"

	DefaultPermissionCopyMode = PermissionCopyModeTypeNone
)

//nolint
const (
	ServiceVersion = "2020-02-10"

	// SASVersion indicates the SAS version.
	SASVersion = ServiceVersion

	headerAuthorization           = "Authorization"
	headerXmsDate                 = "x-ms-date"
	headerContentLength           = "Content-Length"
	headerContentEncoding         = "Content-Encoding"
	headerContentLanguage         = "Content-Language"
	headerContentType             = "Content-Type"
	headerContentMD5              = "Content-MD5"
	headerIfModifiedSince         = "If-Modified-Since"
	headerIfMatch                 = "If-Match"
	headerIfNoneMatch             = "If-None-Match"
	headerIfUnmodifiedSince       = "If-Unmodified-Since"
	headerRange                   = "Range"
	headerDate                    = "Date"
	headerXmsVersion              = "x-ms-version"
	headerAcceptCharset           = "Accept-Charset"
	headerDataServiceVersion      = "DataServiceVersion"
	headerMaxDataServiceVersion   = "MaxDataServiceVersion"
	headerContentTransferEncoding = "Content-Transfer-Encoding"

	etagOData = "odata.etag"
	rfc3339   = "2006-01-02T15:04:05.9999999Z"
	timestamp = "Timestamp"
	etag      = "ETag"

	tokenScope = "https://storage.azure.com/.default"
)
