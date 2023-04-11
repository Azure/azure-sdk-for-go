//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"time"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// SMBProperties contains the optional parameters regarding the SMB/NTFS properties for a file.
type SMBProperties = exported.SMBProperties

// NTFSFileAttributes for Files and Directories.
// The subset of attributes is listed at: https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties#file-system-attributes.
type NTFSFileAttributes = exported.NTFSFileAttributes

// Permissions contains the optional parameters for the permissions on the file.
type Permissions = exported.Permissions

// HTTPHeaders contains optional parameters for the Client.Create method.
type HTTPHeaders = generated.ShareFileHTTPHeaders

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = generated.LeaseAccessConditions

// SourceModifiedAccessConditions contains a group of parameters for the FileClient.UploadRangeFromURL method.
type SourceModifiedAccessConditions = generated.SourceModifiedAccessConditions

// CopyFileSMBInfo contains a group of parameters for the FileClient.StartCopy method.
type CopyFileSMBInfo = generated.CopyFileSMBInfo

// HTTPRange defines a range of bytes within an HTTP resource, starting at offset and
// ending at offset+count. A zero-value HTTPRange indicates the entire resource. An HTTPRange
// which has an offset but no zero value count indicates from the offset to the resource's end.
type HTTPRange = exported.HTTPRange

// ShareFileRangeList - The list of file ranges.
type ShareFileRangeList = generated.ShareFileRangeList

// ClearRange - Ranges there were cleared.
type ClearRange = generated.ClearRange

// ShareFileRange - An Azure Storage file range.
type ShareFileRange = generated.FileRange

// ---------------------------------------------------------------------------------------------------------------------

// CreateOptions contains the optional parameters for the Client.Create method.
type CreateOptions struct {
	// The default value is 'None' for Attributes and 'now' for CreationTime and LastWriteTime fields in file.SMBProperties.
	SMBProperties *SMBProperties
	// The default value is 'inherit' for Permission field in file.Permissions.
	Permissions           *Permissions
	HTTPHeaders           *HTTPHeaders
	LeaseAccessConditions *LeaseAccessConditions
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
}

func (o *CreateOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	createOptions *generated.FileClientCreateOptions, fileHTTPHeaders *generated.ShareFileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return shared.FileAttributesNone, shared.DefaultCurrentTimeString, shared.DefaultCurrentTimeString, &generated.FileClientCreateOptions{
			FilePermission: to.Ptr(shared.DefaultFilePermissionString),
		}, nil, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.Format(false, shared.FileAttributesNone, shared.DefaultCurrentTimeString)

	permission, permissionKey := o.Permissions.Format(shared.DefaultFilePermissionString)

	createOptions = &generated.FileClientCreateOptions{
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
		Metadata:          o.Metadata,
	}

	fileHTTPHeaders = o.HTTPHeaders
	leaseAccessConditions = o.LeaseAccessConditions

	return
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *DeleteOptions) format() (*generated.FileClientDeleteOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}
	return nil, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method.
type GetPropertiesOptions struct {
	// ShareSnapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query for the file properties.
	ShareSnapshot *string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetPropertiesOptions) format() (*generated.FileClientGetPropertiesOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &generated.FileClientGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// SetHTTPHeadersOptions contains the optional parameters for the Client.SetHTTPHeaders method.
type SetHTTPHeadersOptions struct {
	// Resizes a file to the specified size. If the specified byte value is less than the current size of the file, then all ranges
	// above the specified byte value are cleared.
	FileContentLength *int64
	// The default value is 'preserve' for Attributes, CreationTime and LastWriteTime fields in file.SMBProperties.
	SMBProperties *SMBProperties
	// The default value is 'preserve' for Permission field in file.Permissions.
	Permissions *Permissions
	HTTPHeaders *HTTPHeaders
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetHTTPHeadersOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	opts *generated.FileClientSetHTTPHeadersOptions, fileHTTPHeaders *generated.ShareFileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return shared.DefaultPreserveString, shared.DefaultPreserveString, shared.DefaultPreserveString, &generated.FileClientSetHTTPHeadersOptions{
			FilePermission: to.Ptr(shared.DefaultPreserveString),
		}, nil, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.Format(false, shared.DefaultPreserveString, shared.DefaultPreserveString)

	permission, permissionKey := o.Permissions.Format(shared.DefaultPreserveString)

	opts = &generated.FileClientSetHTTPHeadersOptions{
		FileContentLength: o.FileContentLength,
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
	}

	fileHTTPHeaders = o.HTTPHeaders
	leaseAccessConditions = o.LeaseAccessConditions

	return
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetMetadataOptions) format() (*generated.FileClientSetMetadataOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}
	return &generated.FileClientSetMetadataOptions{
		Metadata: o.Metadata,
	}, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// StartCopyFromURLOptions contains the optional parameters for the Client.StartCopyFromURL method.
type StartCopyFromURLOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// required if x-ms-file-permission-copy-mode is specified as override
	Permissions     *Permissions
	CopyFileSMBInfo *CopyFileSMBInfo
	// LeaseAccessConditions contains optional parameters to access leased entity.
	// Required if the destination file has an active lease.
	LeaseAccessConditions *LeaseAccessConditions
}

// TODO: discuss on the types of FileAttributes, FileCreationTime and FileLastWriteTime in CopyFileSMBInfo.
func (o *StartCopyFromURLOptions) format() (*generated.FileClientStartCopyOptions, *generated.CopyFileSMBInfo, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	var permission, permissionKey *string
	if o.Permissions != nil {
		permission = o.Permissions.Permission
		permissionKey = o.Permissions.PermissionKey
	}

	opts := &generated.FileClientStartCopyOptions{
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
		Metadata:          o.Metadata,
	}
	return opts, o.CopyFileSMBInfo, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// AbortCopyOptions contains the optional parameters for the Client.AbortCopy method.
type AbortCopyOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	// Required if the destination file has an active lease.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *AbortCopyOptions) format() (*generated.FileClientAbortCopyOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// DownloadStreamOptions contains the optional parameters for the Client.DownloadStream method.
type DownloadStreamOptions struct {
	// Range specifies a range of bytes. The default value is all bytes.
	Range HTTPRange
	// When this header is set to true and specified together with the Range header, the service returns the MD5 hash for the
	// range, as long as the range is less than or equal to 4 MB in size.
	RangeGetContentMD5 *bool
	// LeaseAccessConditions contains optional parameters to access leased entity.
	// If specified, the operation is performed only if the file's lease is currently active and
	// the lease ID that's specified in the request matches the lease ID of the file.
	// Otherwise, the operation fails with status code 412 (Precondition Failed).
	LeaseAccessConditions *LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// DownloadBufferOptions contains the optional parameters for the Client.DownloadBuffer method.
type DownloadBufferOptions struct {
	// Range specifies a range of bytes. The default value is all bytes.
	Range HTTPRange

	// When this header is set to true and specified together with the Range header, the service returns the MD5 hash for the
	// range, as long as the range is less than or equal to 4 MB in size.
	RangeGetContentMD5 *bool

	// ChunkSize specifies the block size to use for each parallel download; the default size is 4MB.
	ChunkSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	Progress func(bytesTransferred int64)

	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions

	// Concurrency indicates the maximum number of blocks to download in parallel (0=default).
	Concurrency uint16

	// RetryReaderOptionsPerBlock is used when downloading each block.
	RetryReaderOptionsPerBlock RetryReaderOptions
}

// ---------------------------------------------------------------------------------------------------------------------

// DownloadFileOptions contains the optional parameters for the Client.DownloadFile method.
type DownloadFileOptions struct {
	// Range specifies a range of bytes. The default value is all bytes.
	Range HTTPRange

	// When this header is set to true and specified together with the Range header, the service returns the MD5 hash for the
	// range, as long as the range is less than or equal to 4 MB in size.
	RangeGetContentMD5 *bool

	// ChunkSize specifies the block size to use for each parallel download; the default size is 4MB.
	ChunkSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	Progress func(bytesTransferred int64)

	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions

	// Concurrency indicates the maximum number of blocks to download in parallel (0=default).
	Concurrency uint16

	// RetryReaderOptionsPerBlock is used when downloading each block.
	RetryReaderOptionsPerBlock RetryReaderOptions
}

// ---------------------------------------------------------------------------------------------------------------------

// ResizeOptions contains the optional parameters for the Client.Resize method.
type ResizeOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ResizeOptions) format(contentLength int64) (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	opts *generated.FileClientSetHTTPHeadersOptions, leaseAccessConditions *LeaseAccessConditions) {
	fileAttributes, fileCreationTime, fileLastWriteTime = shared.DefaultPreserveString, shared.DefaultPreserveString, shared.DefaultPreserveString

	opts = &generated.FileClientSetHTTPHeadersOptions{
		FileContentLength: &contentLength,
		FilePermission:    to.Ptr(shared.DefaultPreserveString),
	}

	if o != nil {
		leaseAccessConditions = o.LeaseAccessConditions
	}

	return
}

// ---------------------------------------------------------------------------------------------------------------------

// UploadRangeOptions contains the optional parameters for the Client.UploadRange method.
type UploadRangeOptions struct {
	// An MD5 hash of the content. This hash is used to verify the integrity of the data during transport. When the Content-MD5
	// header is specified, the File service compares the hash of the content that has
	// arrived with the header value that was sent. If the two hashes do not match, the operation will fail with error code 400 (Bad Request).
	ContentMD5 []byte
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// ClearRangeOptions contains the optional parameters for the Client.ClearRange method.
type ClearRangeOptions struct {
	// An MD5 hash of the content. This hash is used to verify the integrity of the data during transport. When the Content-MD5
	// header is specified, the File service compares the hash of the content that has
	// arrived with the header value that was sent. If the two hashes do not match, the operation will fail with error code 400 (Bad Request).
	ContentMD5 []byte
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// UploadRangeFromURLOptions contains the optional parameters for the Client.UploadRangeFromURL method.
type UploadRangeFromURLOptions struct {
	// Only Bearer type is supported. Credentials should be a valid OAuth access token to copy source.
	CopySourceAuthorization *string
	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCRC64             []byte
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	LeaseAccessConditions          *LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// GetRangeListOptions contains the optional parameters for the Client.GetRangeList method.
type GetRangeListOptions struct {
	// The previous snapshot parameter is an opaque DateTime value that, when present, specifies the previous snapshot.
	PrevShareSnapshot *string
	// Specifies the range of bytes over which to list ranges, inclusively.
	Range HTTPRange
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetRangeListOptions) format() (*generated.FileClientGetRangeListOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &generated.FileClientGetRangeListOptions{
		Prevsharesnapshot: o.PrevShareSnapshot,
		Range:             exported.FormatHTTPRange(o.Range),
		Sharesnapshot:     o.ShareSnapshot,
	}, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// GetSASURLOptions contains the optional parameters for the Client.GetSASURL method.
type GetSASURLOptions struct {
	StartTime *time.Time
}

func (o *GetSASURLOptions) format() time.Time {
	if o == nil {
		return time.Time{}
	}

	var st time.Time
	if o.StartTime != nil {
		st = o.StartTime.UTC()
	} else {
		st = time.Time{}
	}
	return st
}

// ---------------------------------------------------------------------------------------------------------------------

// ForceCloseHandlesOptions contains the optional parameters for the Client.ForceCloseHandles method.
type ForceCloseHandlesOptions struct {
	// A string value that identifies the portion of the list to be returned with the next list operation. The operation returns
	// a marker value within the response body if the list returned was not complete.
	// The marker value may then be used in a subsequent call to request the next set of list items. The marker value is opaque
	// to the client.
	Marker *string
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ForceCloseHandlesOptions) format() *generated.FileClientForceCloseHandlesOptions {
	if o == nil {
		return nil
	}

	return &generated.FileClientForceCloseHandlesOptions{
		Marker:        o.Marker,
		Sharesnapshot: o.ShareSnapshot,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// ListHandlesOptions contains the optional parameters for the Client.ListHandles method.
type ListHandlesOptions struct {
	// A string value that identifies the portion of the list to be returned with the next list operation. The operation returns
	// a marker value within the response body if the list returned was not complete.
	// The marker value may then be used in a subsequent call to request the next set of list items. The marker value is opaque
	// to the client.
	Marker *string
	// Specifies the maximum number of entries to return. If the request does not specify maxresults, or specifies a value greater
	// than 5,000, the server will return up to 5,000 items.
	MaxResults *int32
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ListHandlesOptions) format() *generated.FileClientListHandlesOptions {
	if o == nil {
		return nil
	}

	return &generated.FileClientListHandlesOptions{
		Marker:        o.Marker,
		Maxresults:    o.MaxResults,
		Sharesnapshot: o.ShareSnapshot,
	}
}

// Handle - A listed Azure Storage handle item.
type Handle = generated.Handle
