//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"io"
	"time"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// SMBProperties contains the optional parameters regarding the SMB/NTFS properties for a file.
type SMBProperties = exported.SMBProperties

// NTFSFileAttributes for Files and Directories.
// The subset of attributes is listed at: https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties#file-system-attributes.
type NTFSFileAttributes = exported.NTFSFileAttributes

// ParseNTFSFileAttributes parses the file attributes from *string to *NTFSFileAttributes.
// It can be used to convert the file attributes to *NTFSFileAttributes where it is returned as *string type in the response.
// It returns an error for any unknown file attribute.
func ParseNTFSFileAttributes(attributes *string) (*NTFSFileAttributes, error) {
	return exported.ParseNTFSFileAttributes(attributes)
}

// Permissions contains the optional parameters for the permissions on the file.
type Permissions = exported.Permissions

// HTTPHeaders contains optional parameters for the Client.Create method.
type HTTPHeaders = generated.ShareFileHTTPHeaders

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = generated.LeaseAccessConditions

// SourceModifiedAccessConditions contains a group of parameters for the FileClient.UploadRangeFromURL method.
type SourceModifiedAccessConditions = generated.SourceModifiedAccessConditions

// HTTPRange defines a range of bytes within an HTTP resource, starting at offset and
// ending at offset+count. A zero-value HTTPRange indicates the entire resource. An HTTPRange
// which has an offset and zero value count indicates from the offset to the resource's end.
type HTTPRange = exported.HTTPRange

// ShareFileRangeList - The list of file ranges.
type ShareFileRangeList = generated.ShareFileRangeList

// ClearRange - Ranges there were cleared.
type ClearRange = generated.ClearRange

// ShareFileRange - An Azure Storage file range.
type ShareFileRange = generated.FileRange

// SourceLeaseAccessConditions contains optional parameters to access the source directory.
type SourceLeaseAccessConditions = generated.SourceLeaseAccessConditions

// DestinationLeaseAccessConditions contains optional parameters to access the destination directory.
type DestinationLeaseAccessConditions = generated.DestinationLeaseAccessConditions

// ---------------------------------------------------------------------------------------------------------------------

// CreateOptions contains the optional parameters for the Client.Create method.
type CreateOptions struct {
	// The default value is 'None' for Attributes and 'now' for CreationTime and LastWriteTime fields in file.SMBProperties.
	SMBProperties *SMBProperties
	// The default value is 'inherit' for Permission field in file.Permissions.
	Permissions           *Permissions
	FilePermissionFormat  *PermissionFormat
	HTTPHeaders           *HTTPHeaders
	LeaseAccessConditions *LeaseAccessConditions
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// NFS only. The file mode of the file or directory
	FileMode *string
	// NFS only. The owner of the file or directory.
	Owner *string
	// NFS only. The owning group of the file or directory.
	Group *string
	// NFS only. Type of the file or directory.
	NfsFileType *NfsFileType
}

func (o *CreateOptions) format() (*generated.FileClientCreateOptions, *generated.ShareFileHTTPHeaders, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.SMBProperties, false)
	permission, permissionKey := exported.FormatPermissions(o.Permissions)

	createOptions := &generated.FileClientCreateOptions{
		FileAttributes:    fileAttributes,
		FileChangeTime:    fileChangeTime,
		FileCreationTime:  fileCreationTime,
		FileLastWriteTime: fileLastWriteTime,
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
		Metadata:          o.Metadata,
		FileMode:          o.FileMode,
		Owner:             o.Owner,
		Group:             o.Group,
		NfsFileType:       o.NfsFileType,
	}

	if permissionKey != nil && *permissionKey != shared.DefaultFilePermissionString {
		createOptions.FilePermissionFormat = to.Ptr(PermissionFormat(shared.DefaultFilePermissionFormat))
	} else if o.FilePermissionFormat != nil {
		createOptions.FilePermissionFormat = to.Ptr(PermissionFormat(*o.FilePermissionFormat))
	}

	return createOptions, o.HTTPHeaders, o.LeaseAccessConditions
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

// RenameOptions contains the optional parameters for the Client.Rename method.
type RenameOptions struct {
	// SMBProperties contains the optional parameters regarding the SMB/NTFS properties for a file.
	SMBProperties *SMBProperties
	// Permissions contains the optional parameters for the permissions on the file.
	Permissions *Permissions
	// Optional. Available for version 2023-06-01 and later. Specifies the format in which the permission is returned. Acceptable
	// values are SDDL or binary. If x-ms-file-permission-format is unspecified or
	// explicitly set to SDDL, the permission is returned in SDDL format. If x-ms-file-permission-format is explicitly set to
	// binary, the permission is returned as a base64 string representing the binary
	// encoding of the permission
	FilePermissionFormat *PermissionFormat
	// ContentType sets the content type of the file.
	ContentType *string
	// IgnoreReadOnly specifies whether the ReadOnly attribute on a pre-existing destination file should be respected.
	// If true, rename will succeed, otherwise, a previous file at the destination with the ReadOnly attribute set will cause rename to fail.
	IgnoreReadOnly *bool
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// ReplaceIfExists specifies that if the destination file already exists, whether this request will overwrite the file or not.
	// If true, rename will succeed and will overwrite the destination file. If not provided or if false and the destination file does exist,
	// the request will not overwrite the destination file.
	// If provided and the destination file does not exist, rename will succeed.
	ReplaceIfExists *bool
	// SourceLeaseAccessConditions contains optional parameters to access the source directory.
	SourceLeaseAccessConditions *SourceLeaseAccessConditions
	// DestinationLeaseAccessConditions contains optional parameters to access the destination directory.
	DestinationLeaseAccessConditions *DestinationLeaseAccessConditions
}

func (o *RenameOptions) format() (*generated.FileClientRenameOptions, *generated.SourceLeaseAccessConditions, *generated.DestinationLeaseAccessConditions, *generated.CopyFileSMBInfo, *generated.ShareFileHTTPHeaders) {
	if o == nil {
		return nil, nil, nil, nil, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.SMBProperties, false)
	permission, permissionKey := exported.FormatPermissions(o.Permissions)

	renameOpts := &generated.FileClientRenameOptions{
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
		IgnoreReadOnly:    o.IgnoreReadOnly,
		Metadata:          o.Metadata,
		ReplaceIfExists:   o.ReplaceIfExists,
	}

	if permissionKey != nil && *permissionKey != shared.DefaultPreserveString {
		renameOpts.FilePermissionFormat = to.Ptr(PermissionFormat(shared.DefaultFilePermissionFormat))
	} else if o.FilePermissionFormat != nil {
		renameOpts.FilePermissionFormat = to.Ptr(PermissionFormat(*o.FilePermissionFormat))
	}

	smbInfo := &generated.CopyFileSMBInfo{
		FileAttributes:    fileAttributes,
		FileChangeTime:    fileChangeTime,
		FileCreationTime:  fileCreationTime,
		FileLastWriteTime: fileLastWriteTime,
	}

	fileHTTPHeaders := &generated.ShareFileHTTPHeaders{
		ContentType: o.ContentType,
	}

	return renameOpts, o.SourceLeaseAccessConditions, o.DestinationLeaseAccessConditions, smbInfo, fileHTTPHeaders
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
	// Optional. Available for version 2023-06-01 and later. Specifies the format in which the permission is returned. Acceptable
	// values are SDDL or binary. If x-ms-file-permission-format is unspecified or
	// explicitly set to SDDL, the permission is returned in SDDL format. If x-ms-file-permission-format is explicitly set to
	// binary, the permission is returned as a base64 string representing the binary
	// encoding of the permission
	FilePermissionFormat *PermissionFormat

	HTTPHeaders *HTTPHeaders
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
	// NFS only. The file mode of the file or directory
	FileMode *string
	// NFS only. The owner of the file or directory.
	Owner *string
	// NFS only. The owning group of the file or directory.
	Group *string
}

func (o *SetHTTPHeadersOptions) format() (*generated.FileClientSetHTTPHeadersOptions, *generated.ShareFileHTTPHeaders, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.SMBProperties, false)
	permission, permissionKey := exported.FormatPermissions(o.Permissions)

	opts := &generated.FileClientSetHTTPHeadersOptions{
		FileAttributes:    fileAttributes,
		FileChangeTime:    fileChangeTime,
		FileCreationTime:  fileCreationTime,
		FileLastWriteTime: fileLastWriteTime,
		FileContentLength: o.FileContentLength,
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
		FileMode:          o.FileMode,
		Owner:             o.Owner,
		Group:             o.Group,
	}

	if permissionKey != nil && *permissionKey != shared.DefaultPreserveString {
		opts.FilePermissionFormat = to.Ptr(PermissionFormat(shared.DefaultFilePermissionFormat))
	} else if o.FilePermissionFormat != nil {
		opts.FilePermissionFormat = to.Ptr(PermissionFormat(*o.FilePermissionFormat))
	}

	return opts, o.HTTPHeaders, o.LeaseAccessConditions
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
	// NFS only. The file mode of the file or directory
	FileMode *string
	// NFS only. The owner of the file or directory.
	Owner *string
	// NFS only. The owning group of the file or directory.
	Group *string
	// NFS only. Applicable only when the copy source is a File. Determines the copy behavior of the mode bits of the file. source:
	FileModeCopyMode *ModeCopyMode
	// NFS only. Determines the copy behavior of the owner user identifier (UID) and group identifier (GID) of the file. source:
	FileOwnerCopyMode *OwnerCopyMode
}

func (o *StartCopyFromURLOptions) format() (*generated.FileClientStartCopyOptions, *generated.CopyFileSMBInfo, *generated.LeaseAccessConditions) {
	if o == nil {
		return &generated.FileClientStartCopyOptions{
				FileModeCopyMode:  to.Ptr(ModeCopyModeSource),
				FileOwnerCopyMode: to.Ptr(OwnerCopyModeSource),
			},
			nil, nil
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
		FileMode:          o.FileMode,
		Owner:             o.Owner,
		Group:             o.Group,
		FileModeCopyMode:  o.FileModeCopyMode,
		FileOwnerCopyMode: o.FileOwnerCopyMode,
	}
	return opts, o.CopyFileSMBInfo.format(), o.LeaseAccessConditions
}

// CopyFileSMBInfo contains a group of parameters for the FileClient.StartCopy method.
type CopyFileSMBInfo struct {
	// Specifies either the option to copy file attributes from a source file(source) to a target file or a list of attributes to set on a target file.
	// CopyFileAttributes is an interface and its underlying implementation are:
	//   - SourceCopyFileAttributes - specifies to copy file attributes from a source file to a target file.
	//   - DestinationCopyFileAttributes - specifies a list of attributes to set on a target file.
	Attributes CopyFileAttributes
	// Specifies either the option to copy file change time from a source file(source) to a target file or a time value in
	// ISO 8601 format to set as change time on a target file.
	// CopyFileChangeTime is an interface and its underlying implementation are:
	//   - SourceCopyFileChangeTime - specifies to copy file change time from a source file to a target file.
	//   - DestinationCopyFileChangeTime - specifies a time value in ISO 8601 format to set as change time on a target file.
	ChangeTime CopyFileChangeTime
	// Specifies either the option to copy file creation time from a source file(source) to a target file or a time value in ISO
	// 8601 format to set as creation time on a target file.
	// CopyFileCreationTime is an interface and its underlying implementation are:
	//   - SourceCopyFileCreationTime - specifies to copy file creation time from a source file to a target file.
	//   - DestinationCopyFileCreationTime - specifies a time value in ISO 8601 format to set as creation time on a target file.
	CreationTime CopyFileCreationTime
	// Specifies either the option to copy file last write time from a source file(source) to a target file or a time value in
	// ISO 8601 format to set as last write time on a target file.
	// CopyFileLastWriteTime is an interface and its underlying implementation are:
	//   - SourceCopyFileLastWriteTime - specifies to copy file last write time from a source file to a target file.
	//   - DestinationCopyFileLastWriteTime - specifies a time value in ISO 8601 format to set as last write time on a target file.
	LastWriteTime CopyFileLastWriteTime
	// Specifies the option to copy file security descriptor from source file or to set it using the value which is defined by
	// the header value of x-ms-file-permission or x-ms-file-permission-key.
	PermissionCopyMode *PermissionCopyModeType
	// Specifies the option to overwrite the target file if it already exists and has read-only attribute set.
	IgnoreReadOnly *bool
	// Specifies the option to set archive attribute on a target file. True means archive attribute will be set on a target file
	// despite attribute overrides or a source file state.
	SetArchiveAttribute *bool
}

func (c *CopyFileSMBInfo) format() *generated.CopyFileSMBInfo {
	if c == nil {
		return nil
	}

	opts := &generated.CopyFileSMBInfo{
		FilePermissionCopyMode: c.PermissionCopyMode,
		IgnoreReadOnly:         c.IgnoreReadOnly,
		SetArchiveAttribute:    c.SetArchiveAttribute,
	}

	if c.Attributes != nil {
		opts.FileAttributes = c.Attributes.FormatAttributes()
	}
	if c.CreationTime != nil {
		opts.FileCreationTime = c.CreationTime.FormatCreationTime()
	}
	if c.LastWriteTime != nil {
		opts.FileLastWriteTime = c.LastWriteTime.FormatLastWriteTime()
	}
	if c.ChangeTime != nil {
		opts.FileChangeTime = c.ChangeTime.FormatChangeTime()
	}

	return opts
}

// CopyFileAttributes specifies either the option to copy file attributes from a source file(source) to a target file or
// a list of attributes to set on a target file.
type CopyFileAttributes = exported.CopyFileAttributes

// SourceCopyFileAttributes specifies to copy file attributes from a source file(source) to a target file
type SourceCopyFileAttributes = exported.SourceCopyFileAttributes

// DestinationCopyFileAttributes specifies a list of attributes to set on a target file.
type DestinationCopyFileAttributes = exported.DestinationCopyFileAttributes

// CopyFileChangeTime specifies either the option to copy file change time from a source file(source) to a target file or
// a time value in ISO 8601 format to set as change time on a target file.
type CopyFileChangeTime = exported.CopyFileChangeTime

// SourceCopyFileChangeTime specifies to copy file change time from a source file(source) to a target file.
type SourceCopyFileChangeTime = exported.SourceCopyFileChangeTime

// DestinationCopyFileChangeTime specifies a time value in ISO 8601 format to set as change time on a target file.
type DestinationCopyFileChangeTime = exported.DestinationCopyFileChangeTime

// CopyFileCreationTime specifies either the option to copy file creation time from a source file(source) to a target file or
// a time value in ISO 8601 format to set as creation time on a target file.
type CopyFileCreationTime = exported.CopyFileCreationTime

// SourceCopyFileCreationTime specifies to copy file creation time from a source file(source) to a target file.
type SourceCopyFileCreationTime = exported.SourceCopyFileCreationTime

// DestinationCopyFileCreationTime specifies a time value in ISO 8601 format to set as creation time on a target file.
type DestinationCopyFileCreationTime = exported.DestinationCopyFileCreationTime

// CopyFileLastWriteTime specifies either the option to copy file last write time from a source file(source) to a target file or
// a time value in ISO 8601 format to set as last write time on a target file.
type CopyFileLastWriteTime = exported.CopyFileLastWriteTime

// SourceCopyFileLastWriteTime specifies to copy file last write time from a source file(source) to a target file.
type SourceCopyFileLastWriteTime = exported.SourceCopyFileLastWriteTime

// DestinationCopyFileLastWriteTime specifies a time value in ISO 8601 format to set as last write time on a target file.
type DestinationCopyFileLastWriteTime = exported.DestinationCopyFileLastWriteTime

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

func (o *DownloadStreamOptions) format() (*generated.FileClientDownloadOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}
	return &generated.FileClientDownloadOptions{
		Range:              exported.FormatHTTPRange(o.Range),
		RangeGetContentMD5: o.RangeGetContentMD5,
	}, o.LeaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// downloadOptions contains common options used by the Client.DownloadBuffer and Client.DownloadFile methods.
type downloadOptions struct {
	// Range specifies a range of bytes. The default value is all bytes.
	Range HTTPRange

	// ChunkSize specifies the chunk size to use for each parallel download; the default size is 4MB.
	ChunkSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	Progress func(bytesTransferred int64)

	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions

	// Concurrency indicates the maximum number of chunks to download in parallel (0=default).
	Concurrency uint16

	// RetryReaderOptionsPerChunk is used when downloading each chunk.
	RetryReaderOptionsPerChunk RetryReaderOptions
}

func (o *downloadOptions) getFilePropertiesOptions() *GetPropertiesOptions {
	if o == nil {
		return nil
	}
	return &GetPropertiesOptions{
		LeaseAccessConditions: o.LeaseAccessConditions,
	}
}

func (o *downloadOptions) getDownloadFileOptions(rng HTTPRange) *DownloadStreamOptions {
	downloadFileOptions := &DownloadStreamOptions{
		Range: rng,
	}
	if o != nil {
		downloadFileOptions.LeaseAccessConditions = o.LeaseAccessConditions
	}
	return downloadFileOptions
}

// DownloadBufferOptions contains the optional parameters for the Client.DownloadBuffer method.
type DownloadBufferOptions struct {
	// Range specifies a range of bytes. The default value is all bytes.
	Range HTTPRange

	// ChunkSize specifies the chunk size to use for each parallel download; the default size is 4MB.
	ChunkSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	Progress func(bytesTransferred int64)

	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions

	// Concurrency indicates the maximum number of chunks to download in parallel (0=default).
	Concurrency uint16

	// RetryReaderOptionsPerChunk is used when downloading each chunk.
	RetryReaderOptionsPerChunk RetryReaderOptions
}

// ---------------------------------------------------------------------------------------------------------------------

// DownloadFileOptions contains the optional parameters for the Client.DownloadFile method.
type DownloadFileOptions struct {
	// Range specifies a range of bytes. The default value is all bytes.
	Range HTTPRange

	// ChunkSize specifies the chunk size to use for each parallel download; the default size is 4MB.
	ChunkSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	Progress func(bytesTransferred int64)

	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions

	// Concurrency indicates the maximum number of chunks to download in parallel (0=default).
	Concurrency uint16

	// RetryReaderOptionsPerChunk is used when downloading each chunk.
	RetryReaderOptionsPerChunk RetryReaderOptions
}

// ---------------------------------------------------------------------------------------------------------------------

// ResizeOptions contains the optional parameters for the Client.Resize method.
type ResizeOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ResizeOptions) format(contentLength int64) (*generated.FileClientSetHTTPHeadersOptions, *LeaseAccessConditions) {
	opts := &generated.FileClientSetHTTPHeadersOptions{
		FileContentLength: &contentLength,
	}

	var leaseAccessConditions *LeaseAccessConditions = nil
	if o != nil {
		leaseAccessConditions = o.LeaseAccessConditions
	}

	return opts, leaseAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// UploadRangeOptions contains the optional parameters for the Client.UploadRange method.
type UploadRangeOptions struct {
	// TransactionalValidation specifies the transfer validation type to use.
	// The default is nil (no transfer validation).
	TransactionalValidation TransferValidationType
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
	// LastWrittenMode specifies if the file last write time should be preserved or overwritten.
	LastWrittenMode *LastWrittenMode
}

func (o *UploadRangeOptions) format(offset int64, body io.ReadSeekCloser) (string, int64, *generated.FileClientUploadRangeOptions, *generated.LeaseAccessConditions, error) {
	if offset < 0 || body == nil {
		return "", 0, nil, nil, errors.New("invalid argument: offset must be >= 0 and body must not be nil")
	}

	count, err := shared.ValidateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return "", 0, nil, nil, err
	}

	if count == 0 {
		return "", 0, nil, nil, errors.New("invalid argument: body must contain readable data whose size is > 0")
	}

	httpRange := exported.FormatHTTPRange(HTTPRange{
		Offset: offset,
		Count:  count,
	})
	rangeParam := ""
	if httpRange != nil {
		rangeParam = *httpRange
	}

	var leaseAccessConditions *LeaseAccessConditions
	uploadRangeOptions := &generated.FileClientUploadRangeOptions{}

	if o != nil {
		leaseAccessConditions = o.LeaseAccessConditions
		uploadRangeOptions.FileLastWrittenMode = o.LastWrittenMode
	}
	if o != nil && o.TransactionalValidation != nil {
		_, err = o.TransactionalValidation.Apply(body, uploadRangeOptions)
		if err != nil {
			return "", 0, nil, nil, err
		}
	}

	return rangeParam, count, uploadRangeOptions, leaseAccessConditions, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// ClearRangeOptions contains the optional parameters for the Client.ClearRange method.
type ClearRangeOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ClearRangeOptions) format(contentRange HTTPRange) (string, *generated.LeaseAccessConditions, error) {
	httpRange := exported.FormatHTTPRange(contentRange)
	if httpRange == nil || contentRange.Offset < 0 || contentRange.Count <= 0 {
		return "", nil, errors.New("invalid argument: either offset is < 0 or count <= 0")
	}

	if o == nil {
		return *httpRange, nil, nil
	}

	return *httpRange, o.LeaseAccessConditions, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// UploadRangeFromURLOptions contains the optional parameters for the Client.UploadRangeFromURL method.
type UploadRangeFromURLOptions struct {
	// Only Bearer type is supported. Credentials should be a valid OAuth access token to copy source.
	CopySourceAuthorization *string
	// SourceContentValidation contains the validation mechanism used on the range of bytes read from the source.
	SourceContentValidation        SourceContentValidationType
	SourceModifiedAccessConditions *SourceModifiedAccessConditions
	LeaseAccessConditions          *LeaseAccessConditions
	// LastWrittenMode specifies if the file last write time should be preserved or overwritten.
	LastWrittenMode *LastWrittenMode

	// Deprecated: Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCRC64 uint64
}

func (o *UploadRangeFromURLOptions) format(sourceOffset int64, destinationOffset int64, count int64) (string, *generated.FileClientUploadRangeFromURLOptions, *generated.SourceModifiedAccessConditions, *generated.LeaseAccessConditions, error) {
	if sourceOffset < 0 || destinationOffset < 0 {
		return "", nil, nil, nil, errors.New("invalid argument: source and destination offsets must be >= 0")
	}

	httpRangeSrc := exported.FormatHTTPRange(HTTPRange{Offset: sourceOffset, Count: count})
	httpRangeDest := exported.FormatHTTPRange(HTTPRange{Offset: destinationOffset, Count: count})
	destRange := ""
	if httpRangeDest != nil {
		destRange = *httpRangeDest
	}

	opts := &generated.FileClientUploadRangeFromURLOptions{
		SourceRange: httpRangeSrc,
	}

	var sourceModifiedAccessConditions *SourceModifiedAccessConditions
	var leaseAccessConditions *LeaseAccessConditions

	if o != nil {
		opts.CopySourceAuthorization = o.CopySourceAuthorization
		opts.FileLastWrittenMode = o.LastWrittenMode
		sourceModifiedAccessConditions = o.SourceModifiedAccessConditions
		leaseAccessConditions = o.LeaseAccessConditions

		if o.SourceContentValidation != nil {
			err := o.SourceContentValidation.apply(opts)
			if err != nil {
				return "", nil, nil, nil, err
			}
		}
	}

	return destRange, opts, sourceModifiedAccessConditions, leaseAccessConditions, nil
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
	// SupportRename determines whether the changed ranges for a file should be listed when the file's location in the
	// previous snapshot is different from the location in the Request URI, as a result of rename or move operations.
	SupportRename *bool
}

func (o *GetRangeListOptions) format() (*generated.FileClientGetRangeListOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &generated.FileClientGetRangeListOptions{
		Prevsharesnapshot: o.PrevShareSnapshot,
		Range:             exported.FormatHTTPRange(o.Range),
		Sharesnapshot:     o.ShareSnapshot,
		SupportRename:     o.SupportRename,
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

// CreateHardLinkOptions contains the optional parameters for the Client.CreateHardLink method.
type CreateHardLinkOptions struct {
	// NFS only. Required. Specifies the path of the target file to which the link will be created, up to 2 KiB in length.
	TargetFile string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *CreateHardLinkOptions) format() (*generated.FileClientCreateHardLinkOptions, string, *generated.LeaseAccessConditions) {
	return &generated.FileClientCreateHardLinkOptions{}, o.TargetFile, o.LeaseAccessConditions
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

// ---------------------------------------------------------------------------------------------------------------------

// uploadFromReaderOptions identifies options used by the UploadBuffer and UploadFile functions.
type uploadFromReaderOptions struct {
	// ChunkSize specifies the chunk size to use in bytes; the default (and maximum size) is MaxUpdateRangeBytes.
	ChunkSize int64

	// Progress is a function that is invoked periodically as bytes are sent to the FileClient.
	// Note that the progress reporting is not always increasing; it can go down when retrying a request.
	Progress func(bytesTransferred int64)

	// Concurrency indicates the maximum number of chunks to upload in parallel (default is 5)
	Concurrency uint16

	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

// UploadBufferOptions provides set of configurations for Client.UploadBuffer operation.
type UploadBufferOptions = uploadFromReaderOptions

// UploadFileOptions provides set of configurations for Client.UploadFile operation.
type UploadFileOptions = uploadFromReaderOptions

func (o *uploadFromReaderOptions) getUploadRangeOptions() *UploadRangeOptions {
	return &UploadRangeOptions{
		LeaseAccessConditions: o.LeaseAccessConditions,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// UploadStreamOptions provides set of configurations for Client.UploadStream operation.
type UploadStreamOptions struct {
	// ChunkSize defines the size of the buffer used during upload. The default and minimum value is 1 MiB.
	// Maximum size of a chunk is MaxUpdateRangeBytes.
	ChunkSize int64

	// Concurrency defines the max number of concurrent uploads to be performed to upload the file.
	// Each concurrent upload will create a buffer of size ChunkSize.  The default value is one.
	Concurrency int

	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (u *UploadStreamOptions) setDefaults() {
	if u.Concurrency == 0 {
		u.Concurrency = 1
	}

	if u.ChunkSize < _1MiB {
		u.ChunkSize = _1MiB
	}
}

func (u *UploadStreamOptions) getUploadRangeOptions() *UploadRangeOptions {
	return &UploadRangeOptions{
		LeaseAccessConditions: u.LeaseAccessConditions,
	}
}

// URLParts object represents the components that make up an Azure Storage Share/Directory/File URL.
// NOTE: Changing any SAS-related field requires computing a new SAS signature.
type URLParts = sas.URLParts

// ParseURL parses a URL initializing URLParts' fields including any SAS-related & snapshot query parameters. Any other
// query parameters remain in the UnparsedParams field. This method overwrites all fields in the URLParts object.
func ParseURL(u string) (URLParts, error) {
	return sas.ParseURL(u)
}
