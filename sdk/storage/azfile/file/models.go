// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"errors"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
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

// NFSProperties contains the optional parameters regarding the NFS properties for a file.
type NFSProperties = exported.NFSProperties

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
	// NFS only. The default value is 'now' for CreationTime and LastWriteTime fields in file.NFSProperties.
	NFSProperties *NFSProperties
	// The default value is 'inherit' for Permission field in file.Permissions.
	Permissions           *Permissions
	FilePermissionFormat  *PermissionFormat
	HTTPHeaders           *HTTPHeaders
	LeaseAccessConditions *LeaseAccessConditions
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// SMB only. How attributes and permissions should be set on the file.
	// New: automatically adds the ARCHIVE file attribute flag and uses Windows create file permissions semantics (ex: inherit from parent).
	// Restore: does not modify file attribute flag and uses Windows update file permissions semantics.
	// If Restore is specified, the file permission must also be provided, otherwise PropertySemantics will default to New.
	FilePropertySemantics *PropertySemantics
	OptionalBody          io.ReadSeekCloser
	ContentLength         *int64
	ContentMD5            []byte
}

func (o *CreateOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientCreateOptions {
	opts := &generated.FileClientCreateOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Metadata = o.Metadata
	opts.Optionalbody = o.OptionalBody

	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}

	if o.HTTPHeaders != nil {
		opts.FileCacheControl = o.HTTPHeaders.CacheControl
		opts.FileContentDisposition = o.HTTPHeaders.ContentDisposition
		opts.FileContentEncoding = o.HTTPHeaders.ContentEncoding
		opts.FileContentLanguage = o.HTTPHeaders.ContentLanguage
		opts.FileContentMD5 = o.HTTPHeaders.ContentMD5
		opts.FileContentType = o.HTTPHeaders.ContentType
	}

	if o.NFSProperties != nil {
		fileCreationTime, fileLastWriteTime := exported.FormatNFSProperties(o.NFSProperties, false)

		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.FileMode = o.NFSProperties.FileMode
		opts.Group = o.NFSProperties.Group
		opts.Owner = o.NFSProperties.Owner
	} else {
		fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.SMBProperties, false)
		permission, permissionKey := exported.FormatPermissions(o.Permissions)

		opts.FileAttributes = fileAttributes
		opts.FileChangeTime = fileChangeTime
		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.FilePermission = permission
		opts.FilePermissionKey = permissionKey
		opts.Optionalbody = o.OptionalBody

		// Refer the documentation for details - https://learn.microsoft.com/en-us/rest/api/storageservices/create-file#smb-only-request-headers
		if permissionKey != nil {
			opts.FilePermissionKey = permissionKey
		} else if permission != nil {
			opts.FilePermission = permission
			if o.FilePermissionFormat != nil {
				opts.FilePermissionFormat = to.Ptr(*o.FilePermissionFormat)
			} else {
				opts.FilePermissionFormat = to.Ptr(FilePermissionFormatSddl) // optional, default
			}
		}
		if o.FilePropertySemantics != nil {
			opts.FilePropertySemantics = o.FilePropertySemantics
		}
	}
	if len(o.ContentMD5) > 0 {
		opts.ContentMD5 = o.ContentMD5
	}
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *DeleteOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientDeleteOptions {
	opts := &generated.FileClientDeleteOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
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

func (o *RenameOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool, allowSourceTrailingDot *bool) *generated.FileClientRenameOptions {
	opts := &generated.FileClientRenameOptions{
		FileRequestIntent:      fileRequestIntent,
		AllowTrailingDot:       allowTrailingDot,
		AllowSourceTrailingDot: allowSourceTrailingDot,
	}
	if o == nil {
		return opts
	}

	fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.SMBProperties, false)
	permission, permissionKey := exported.FormatPermissions(o.Permissions)

	opts.FilePermission = permission
	opts.FilePermissionKey = permissionKey
	opts.IgnoreReadOnly = o.IgnoreReadOnly
	opts.Metadata = o.Metadata
	opts.ReplaceIfExists = o.ReplaceIfExists
	opts.FileContentType = o.ContentType

	if permissionKey != nil && *permissionKey != shared.DefaultPreserveString {
		opts.FilePermissionFormat = to.Ptr(PermissionFormat(shared.DefaultFilePermissionFormat))
	} else if o.FilePermissionFormat != nil {
		opts.FilePermissionFormat = to.Ptr(PermissionFormat(*o.FilePermissionFormat))
	}

	opts.FileAttributes = fileAttributes
	opts.FileChangeTime = fileChangeTime
	opts.FileCreationTime = fileCreationTime
	opts.FileLastWriteTime = fileLastWriteTime

	if o.SourceLeaseAccessConditions != nil {
		opts.SourceLeaseID = o.SourceLeaseAccessConditions.SourceLeaseID
	}
	if o.DestinationLeaseAccessConditions != nil {
		opts.DestinationLeaseID = o.DestinationLeaseAccessConditions.DestinationLeaseID
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method.
type GetPropertiesOptions struct {
	// ShareSnapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query for the file properties.
	ShareSnapshot *string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetPropertiesOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientGetPropertiesOptions {
	opts := &generated.FileClientGetPropertiesOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Sharesnapshot = o.ShareSnapshot
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// SetHTTPHeadersOptions contains the optional parameters for the Client.SetHTTPHeaders method.
type SetHTTPHeadersOptions struct {
	// Resizes a file to the specified size. If the specified byte value is less than the current size of the file, then all ranges
	// above the specified byte value are cleared.
	FileContentLength *int64
	// The default value is 'preserve' for Attributes, CreationTime and LastWriteTime fields in file.SMBProperties.
	SMBProperties *SMBProperties
	// NFS only. The default value is 'now' for CreationTime and LastWriteTime fields in file.NFSProperties.
	NFSProperties *NFSProperties
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
}

func (o *SetHTTPHeadersOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientSetHTTPHeadersOptions {
	opts := &generated.FileClientSetHTTPHeadersOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.FileContentLength = o.FileContentLength

	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}

	if o.HTTPHeaders != nil {
		opts.FileCacheControl = o.HTTPHeaders.CacheControl
		opts.FileContentDisposition = o.HTTPHeaders.ContentDisposition
		opts.FileContentEncoding = o.HTTPHeaders.ContentEncoding
		opts.FileContentLanguage = o.HTTPHeaders.ContentLanguage
		opts.FileContentMD5 = o.HTTPHeaders.ContentMD5
		opts.FileContentType = o.HTTPHeaders.ContentType
	}

	if o.NFSProperties != nil {
		fileCreationTime, fileLastWriteTime := exported.FormatNFSProperties(o.NFSProperties, false)

		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.FileMode = o.NFSProperties.FileMode
		opts.Group = o.NFSProperties.Group
		opts.Owner = o.NFSProperties.Owner
	} else {
		fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.SMBProperties, false)
		permission, permissionKey := exported.FormatPermissions(o.Permissions)

		opts.FileAttributes = fileAttributes
		opts.FileChangeTime = fileChangeTime
		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.FilePermission = permission
		opts.FilePermissionKey = permissionKey

		if permissionKey != nil && *permissionKey != shared.DefaultPreserveString {
			opts.FilePermissionFormat = to.Ptr(PermissionFormat(shared.DefaultFilePermissionFormat))
		} else if o.FilePermissionFormat != nil {
			opts.FilePermissionFormat = to.Ptr(PermissionFormat(*o.FilePermissionFormat))
		}
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetMetadataOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientSetMetadataOptions {
	opts := &generated.FileClientSetMetadataOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}
	opts.Metadata = o.Metadata
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

// CopyFileNFSProperties contains the optional parameters regarding the NFS properties for a file.
type CopyFileNFSProperties = exported.CopyFileNFSProperties

// ---------------------------------------------------------------------------------------------------------------------

// StartCopyFromURLOptions contains the optional parameters for the Client.StartCopyFromURL method.
type StartCopyFromURLOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// required if x-ms-file-permission-copy-mode is specified as override
	Permissions     *Permissions
	CopyFileSMBInfo *CopyFileSMBInfo
	// NFS only.
	CopyFileNFSProperties *CopyFileNFSProperties
	// LeaseAccessConditions contains optional parameters to access leased entity.
	// Required if the destination file has an active lease.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *StartCopyFromURLOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool, allowSourceTrailingDot *bool) *generated.FileClientStartCopyOptions {
	opts := &generated.FileClientStartCopyOptions{
		FileRequestIntent:      fileRequestIntent,
		AllowTrailingDot:       allowTrailingDot,
		AllowSourceTrailingDot: allowSourceTrailingDot,
	}
	if o == nil {
		return opts
	}
	opts.Metadata = o.Metadata
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}

	if o.CopyFileNFSProperties != nil {
		opts.FileMode = o.CopyFileNFSProperties.FileMode
		opts.Owner = o.CopyFileNFSProperties.Owner
		opts.Group = o.CopyFileNFSProperties.Group
		opts.FileModeCopyMode = o.CopyFileNFSProperties.FileModeCopyMode
		opts.FileOwnerCopyMode = o.CopyFileNFSProperties.FileOwnerCopyMode
		if o.CopyFileNFSProperties.CreationTime != nil {
			opts.FileCreationTime = o.CopyFileNFSProperties.CreationTime.FormatCreationTime()
		}
		if o.CopyFileNFSProperties.LastWriteTime != nil {
			opts.FileLastWriteTime = o.CopyFileNFSProperties.LastWriteTime.FormatLastWriteTime()
		}
	} else {
		if o.Permissions != nil {
			opts.FilePermission = o.Permissions.Permission
			opts.FilePermissionKey = o.Permissions.PermissionKey
		}
		if o.CopyFileSMBInfo != nil {
			opts.FilePermissionCopyMode = o.CopyFileSMBInfo.PermissionCopyMode
			opts.IgnoreReadOnly = o.CopyFileSMBInfo.IgnoreReadOnly
			opts.SetArchiveAttribute = o.CopyFileSMBInfo.SetArchiveAttribute
			if o.CopyFileSMBInfo.Attributes != nil {
				opts.FileAttributes = o.CopyFileSMBInfo.Attributes.FormatAttributes()
			}
			if o.CopyFileSMBInfo.CreationTime != nil {
				opts.FileCreationTime = o.CopyFileSMBInfo.CreationTime.FormatCreationTime()
			}
			if o.CopyFileSMBInfo.LastWriteTime != nil {
				opts.FileLastWriteTime = o.CopyFileSMBInfo.LastWriteTime.FormatLastWriteTime()
			}
			if o.CopyFileSMBInfo.ChangeTime != nil {
				opts.FileChangeTime = o.CopyFileSMBInfo.ChangeTime.FormatChangeTime()
			}
		}
	}
	return opts
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

func (o *AbortCopyOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientAbortCopyOptions {
	opts := &generated.FileClientAbortCopyOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
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

func (o *DownloadStreamOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientDownloadOptions {
	opts := &generated.FileClientDownloadOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}
	opts.Range = exported.FormatHTTPRange(o.Range)
	opts.RangeGetContentMD5 = o.RangeGetContentMD5
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
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

func (o *ResizeOptions) format(contentLength int64, fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientSetHTTPHeadersOptions {
	opts := &generated.FileClientSetHTTPHeadersOptions{
		FileContentLength: &contentLength,
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o != nil && o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
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

func (o *UploadRangeOptions) format(offset int64, body io.ReadSeekCloser, fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) (string, int64, *generated.FileClientUploadRangeOptions, error) {
	if offset < 0 || body == nil {
		return "", 0, nil, errors.New("invalid argument: offset must be >= 0 and body must not be nil")
	}

	count, err := shared.ValidateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return "", 0, nil, err
	}

	if count == 0 {
		return "", 0, nil, errors.New("invalid argument: body must contain readable data whose size is > 0")
	}

	httpRange := exported.FormatHTTPRange(HTTPRange{
		Offset: offset,
		Count:  count,
	})
	rangeParam := ""
	if httpRange != nil {
		rangeParam = *httpRange
	}

	uploadRangeOptions := &generated.FileClientUploadRangeOptions{
		Optionalbody:      body,
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}

	if o != nil {
		if o.LeaseAccessConditions != nil {
			uploadRangeOptions.LeaseID = o.LeaseAccessConditions.LeaseID
		}
		uploadRangeOptions.FileLastWrittenMode = o.LastWrittenMode
	}
	if o != nil && o.TransactionalValidation != nil {
		_, err = o.TransactionalValidation.Apply(body, uploadRangeOptions)
		if err != nil {
			return "", 0, nil, err
		}
	}

	return rangeParam, count, uploadRangeOptions, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// ClearRangeOptions contains the optional parameters for the Client.ClearRange method.
type ClearRangeOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ClearRangeOptions) format(contentRange HTTPRange, fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) (string, *generated.FileClientUploadRangeOptions, error) {
	httpRange := exported.FormatHTTPRange(contentRange)
	if httpRange == nil || contentRange.Offset < 0 || contentRange.Count <= 0 {
		return "", nil, errors.New("invalid argument: either offset is < 0 or count <= 0")
	}

	opts := &generated.FileClientUploadRangeOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o != nil && o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}

	return *httpRange, opts, nil
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

func (o *UploadRangeFromURLOptions) format(sourceOffset int64, destinationOffset int64, count int64, fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool, allowSourceTrailingDot *bool) (string, *generated.FileClientUploadRangeFromURLOptions, error) {
	if sourceOffset < 0 || destinationOffset < 0 {
		return "", nil, errors.New("invalid argument: source and destination offsets must be >= 0")
	}

	httpRangeSrc := exported.FormatHTTPRange(HTTPRange{Offset: sourceOffset, Count: count})
	httpRangeDest := exported.FormatHTTPRange(HTTPRange{Offset: destinationOffset, Count: count})
	destRange := ""
	if httpRangeDest != nil {
		destRange = *httpRangeDest
	}

	opts := &generated.FileClientUploadRangeFromURLOptions{
		SourceRange:            httpRangeSrc,
		FileRequestIntent:      fileRequestIntent,
		AllowTrailingDot:       allowTrailingDot,
		AllowSourceTrailingDot: allowSourceTrailingDot,
	}

	if o != nil {
		opts.CopySourceAuthorization = o.CopySourceAuthorization
		opts.FileLastWrittenMode = o.LastWrittenMode
		if o.SourceModifiedAccessConditions != nil {
			opts.SourceIfMatchCRC64 = o.SourceModifiedAccessConditions.SourceIfMatchCRC64
			opts.SourceIfNoneMatchCRC64 = o.SourceModifiedAccessConditions.SourceIfNoneMatchCRC64
		}
		if o.LeaseAccessConditions != nil {
			opts.LeaseID = o.LeaseAccessConditions.LeaseID
		}

		if o.SourceContentValidation != nil {
			err := o.SourceContentValidation.apply(opts)
			if err != nil {
				return "", nil, err
			}
		}
	}

	return destRange, opts, nil
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

func (o *GetRangeListOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientGetRangeListOptions {
	opts := &generated.FileClientGetRangeListOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Prevsharesnapshot = o.PrevShareSnapshot
	opts.Range = exported.FormatHTTPRange(o.Range)
	opts.Sharesnapshot = o.ShareSnapshot
	opts.SupportRename = o.SupportRename
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}

	return opts
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

func (o *ForceCloseHandlesOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientForceCloseHandlesOptions {
	opts := &generated.FileClientForceCloseHandlesOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Marker = o.Marker
	opts.Sharesnapshot = o.ShareSnapshot

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// CreateHardLinkOptions contains the optional parameters for the Client.CreateHardLink method.
type CreateHardLinkOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *CreateHardLinkOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.FileClientCreateHardLinkOptions {
	opts := &generated.FileClientCreateHardLinkOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// CreateSymbolicLinkOptions contains the optional parameters for the Client.CreateSymbolicLink method.
type CreateSymbolicLinkOptions struct {
	// LeaseAccessConditions contains optional parameters to access leased entity.
	LeaseAccessConditions *LeaseAccessConditions
	// NFS only.
	FileNFSProperties *NFSProperties
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// The default value is 'Directory' for Attributes and 'now' for CreationTime and LastWriteTime fields in file.SMBProperties.
	FileSMBProperties *SMBProperties
	// Client request id
	ClientRequestID *string
}

func (o *CreateSymbolicLinkOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.FileClientCreateSymbolicLinkOptions {
	opts := &generated.FileClientCreateSymbolicLinkOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}
	opts.Metadata = o.Metadata
	if o.LeaseAccessConditions != nil {
		opts.LeaseID = o.LeaseAccessConditions.LeaseID
	}
	if o.FileNFSProperties != nil {
		fileCreationTime, fileLastWriteTime := exported.FormatNFSProperties(o.FileNFSProperties, false)
		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.Group = o.FileNFSProperties.Group
		opts.Owner = o.FileNFSProperties.Owner
	}
	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// GetSymbolicLinkOptions contains the optional parameters for the Client.GetSymbolicLink method.
type GetSymbolicLinkOptions struct {

	// The snapshot parameter, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	// Client request id
	ClientRequestID *string
}

func (o *GetSymbolicLinkOptions) format(fileRequestIntent *generated.ShareTokenIntent) *generated.FileClientGetSymbolicLinkOptions {
	opts := &generated.FileClientGetSymbolicLinkOptions{
		FileRequestIntent: fileRequestIntent,
	}
	if o == nil {
		return opts
	}
	opts.Sharesnapshot = o.ShareSnapshot
	return opts
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

func (o *ListHandlesOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.FileClientListHandlesOptions {
	opts := &generated.FileClientListHandlesOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Marker = o.Marker
	opts.Maxresults = o.MaxResults
	opts.Sharesnapshot = o.ShareSnapshot

	return opts
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
