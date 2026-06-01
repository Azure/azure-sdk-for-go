// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

/*
Due to differences in the authoring of swagger to typespec, not all the models are generating.
Adding the missing models manually.
*/

// DestinationLeaseAccessConditions contains a group of parameters for the DirectoryClient.Rename method.
type DestinationLeaseAccessConditions struct {
	// Required if the destination file has an active infinite lease. The lease ID specified for this header must match the lease
	// ID of the destination file. If the request does not include the lease ID or
	// it is not valid, the operation fails with status code 412 (Precondition Failed). If this header is specified and the destination
	// file does not currently have an active lease, the operation will also
	// fail with status code 412 (Precondition Failed).
	DestinationLeaseID *string
}

// LeaseAccessConditions contains a group of parameters for the ShareClient.GetProperties method.
type LeaseAccessConditions struct {
	// If specified, the operation only succeeds if the resource's lease is active and matches this ID.
	LeaseID *string
}

// ShareFileHTTPHeaders contains a group of parameters for the FileClient.Create method.
type ShareFileHTTPHeaders struct {
	// Sets the file's cache control. The File service stores this value but does not use or modify it.
	CacheControl *string

	// Sets the file's Content-Disposition header.
	ContentDisposition *string

	// Specifies which content encodings have been applied to the file.
	ContentEncoding *string

	// Specifies the natural languages used by this resource.
	ContentLanguage *string

	// Sets the file's MD5 hash.
	ContentMD5 []byte

	// Sets the MIME content type of the file. The default type is 'application/octet-stream'.
	ContentType *string
}

// SourceModifiedAccessConditions contains a group of parameters for the FileClient.UploadRangeFromURL method.
type SourceModifiedAccessConditions struct {
	// Specify the crc64 value to operate only on range with a matching crc64 checksum.
	SourceIfMatchCRC64 []byte

	// Specify the crc64 value to operate only on range without a matching crc64 checksum.
	SourceIfNoneMatchCRC64 []byte
}

// SourceLeaseAccessConditions contains a group of parameters for the DirectoryClient.Rename method.
type SourceLeaseAccessConditions struct {
	// Required if the source file has an active infinite lease.
	SourceLeaseID *string
}
