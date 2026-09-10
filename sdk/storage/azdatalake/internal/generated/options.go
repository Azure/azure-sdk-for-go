// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// CPKInfo contains a group of parameters for the PathClient.Create method.
type CPKInfo struct {
	// The algorithm used to produce the encryption key hash. Currently, the only accepted value is "AES256". Must be provided
	// if the x-ms-encryption-key header is provided.
	EncryptionAlgorithm *EncryptionAlgorithmType

	// Optional. Specifies the encryption key to use to encrypt the data provided in the request. If not specified, encryption
	// is performed with the root account encryption key. For more information, see
	// Encryption at Rest for Azure Storage Services.
	EncryptionKey *string

	// The SHA-256 hash of the provided encryption key. Must be provided if the x-ms-encryption-key header is provided.
	EncryptionKeySHA256 *string
}

// LeaseAccessConditions contains a group of parameters for the PathClient.Create method.
type LeaseAccessConditions struct {
	// If specified, the operation only succeeds if the resource's lease is active and matches this ID.
	LeaseID *string
}

// ModifiedAccessConditions contains a group of parameters for the FileSystemClient.SetProperties method.
type ModifiedAccessConditions struct {
	// Specify an ETag value to operate only on blobs with a matching value.
	IfMatch *azcore.ETag

	// Specify this header value to operate only on a blob if it has been modified since the specified date/time.
	IfModifiedSince *time.Time

	// Specify an ETag value to operate only on blobs without a matching value.
	IfNoneMatch *azcore.ETag

	// Specify this header value to operate only on a blob if it has not been modified since the specified date/time.
	IfUnmodifiedSince *time.Time
}

// PathHTTPHeaders contains a group of parameters for the PathClient.Create method.
type PathHTTPHeaders struct {
	// Optional. Sets the blob's cache control. If specified, this property is stored with the blob and returned with a read request.
	CacheControl *string

	// Optional. Sets the blob's Content-Disposition header.
	ContentDisposition *string

	// Optional. Sets the blob's content encoding. If specified, this property is stored with the blob and returned with a read
	// request.
	ContentEncoding *string

	// Optional. Set the blob's content language. If specified, this property is stored with the blob and returned with a read
	// request.
	ContentLanguage *string

	// Specify the transactional md5 for the body, to be validated by the service.
	ContentMD5 []byte

	// Optional. Sets the blob's content type. If specified, this property is stored with the blob and returned with a read request.
	ContentType *string

	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentHash []byte
}

// SourceModifiedAccessConditions contains a group of parameters for the PathClient.Create method.
type SourceModifiedAccessConditions struct {
	// Specify an ETag value to operate only on blobs with a matching value.
	SourceIfMatch *azcore.ETag

	// Specify this header value to operate only on a blob if it has been modified since the specified date/time.
	SourceIfModifiedSince *time.Time

	// Specify an ETag value to operate only on blobs without a matching value.
	SourceIfNoneMatch *azcore.ETag

	// Specify this header value to operate only on a blob if it has not been modified since the specified date/time.
	SourceIfUnmodifiedSince *time.Time
}
