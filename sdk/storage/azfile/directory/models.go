//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"reflect"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NfsFileType specifies the type of the file or directory.
type NfsFileType = generated.NfsFileType

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// DestinationLeaseAccessConditions contains optional parameters to access the destination directory.
type DestinationLeaseAccessConditions = generated.DestinationLeaseAccessConditions

// ---------------------------------------------------------------------------------------------------------------------

// CreateOptions contains the optional parameters for the Client.Create method.
type CreateOptions struct {
	// The default value is 'Directory' for Attributes and 'now' for CreationTime and LastWriteTime fields in file.SMBProperties.
	FileSMBProperties *file.SMBProperties
	// The default value is 'inherit' for Permission field in file.Permissions.
	FilePermissions *file.Permissions
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
	// Optional. Available for version 2023-06-01 and later. Specifies the format in which the permission is returned. Acceptable
	// values are SDDL or binary. If x-ms-file-permission-format is unspecified or
	// explicitly set to SDDL, the permission is returned in SDDL format. If x-ms-file-permission-format is explicitly set to
	// binary, the permission is returned as a base64 string representing the binary
	// encoding of the permission
	FilePermissionFormat *FilePermissionFormat
	// NFS only. The file mode of the file or directory
	FileMode *string
	// NFS only. The owner of the file or directory.
	Owner *string
	// NFS only. The owning group of the file or directory.
	Group *string
	// NFS only. Type of the file or directory.
	NfsFileType *NfsFileType
}

func (o *CreateOptions) format() *generated.DirectoryClientCreateOptions {
	if o == nil {
		return nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.FileSMBProperties, true)
	permission, permissionKey := exported.FormatPermissions(o.FilePermissions)
	createOptions := &generated.DirectoryClientCreateOptions{
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
	}

	if permissionKey != nil && *permissionKey != shared.DefaultFilePermissionString {
		createOptions.FilePermissionFormat = to.Ptr(FilePermissionFormat(shared.DefaultFilePermissionFormat))
	} else if o.FilePermissionFormat != nil {
		createOptions.FilePermissionFormat = to.Ptr(FilePermissionFormat(*o.FilePermissionFormat))
	}

	return createOptions
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	// placeholder for future options
}

func (o *DeleteOptions) format() *generated.DirectoryClientDeleteOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// RenameOptions contains the optional parameters for the Client.Rename method.
type RenameOptions struct {
	// FileSMBProperties contains the optional parameters regarding the SMB/NTFS properties for a file.
	FileSMBProperties *file.SMBProperties
	// FilePermissions contains the optional parameters for the permissions on the file.
	FilePermissions *file.Permissions
	// FilePermissionFormat contains the file permission format, sddl(Default) or Binary.
	FilePermissionFormat *FilePermissionFormat
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
	// DestinationLeaseAccessConditions contains optional parameters to access the destination directory.
	DestinationLeaseAccessConditions *DestinationLeaseAccessConditions
}

func (o *RenameOptions) format() (*generated.DirectoryClientRenameOptions, *generated.DestinationLeaseAccessConditions, *generated.CopyFileSMBInfo) {
	if o == nil {
		return nil, nil, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.FileSMBProperties, true)
	permission, permissionKey := exported.FormatPermissions(o.FilePermissions)

	renameOpts := &generated.DirectoryClientRenameOptions{
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
		IgnoreReadOnly:    o.IgnoreReadOnly,
		Metadata:          o.Metadata,
		ReplaceIfExists:   o.ReplaceIfExists,
	}

	if permissionKey != nil && *permissionKey != shared.DefaultPreserveString {
		renameOpts.FilePermissionFormat = to.Ptr(FilePermissionFormat(shared.DefaultFilePermissionFormat))
	} else if o.FilePermissionFormat != nil {
		renameOpts.FilePermissionFormat = to.Ptr(FilePermissionFormat(*o.FilePermissionFormat))
	}

	smbInfo := &generated.CopyFileSMBInfo{
		FileAttributes:    fileAttributes,
		FileChangeTime:    fileChangeTime,
		FileCreationTime:  fileCreationTime,
		FileLastWriteTime: fileLastWriteTime,
	}

	return renameOpts, o.DestinationLeaseAccessConditions, smbInfo
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method.
type GetPropertiesOptions struct {
	// ShareSnapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query for the directory properties.
	ShareSnapshot *string
}

func (o *GetPropertiesOptions) format() *generated.DirectoryClientGetPropertiesOptions {
	if o == nil {
		return nil
	}

	return &generated.DirectoryClientGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// SetPropertiesOptions contains the optional parameters for the Client.SetProperties method.
type SetPropertiesOptions struct {
	// The default value is 'preserve' for Attributes, CreationTime and LastWriteTime fields in file.SMBProperties.
	FileSMBProperties *file.SMBProperties
	// The default value is 'preserve' for Permission field in file.Permissions.
	FilePermissions *file.Permissions
	// FilePermissionFormat contains the format of the file permissions, Can be sddl (Default) or Binary.
	FilePermissionFormat *FilePermissionFormat
	// NFS only. The file mode of the file or directory
	FileMode *string
	// NFS only. The owner of the file or directory.
	Owner *string
	// NFS only. The owning group of the file or directory.
	Group *string
	// NFS only. Type of the file or directory.
	NfsFileType *NfsFileType
}

func (o *SetPropertiesOptions) format() *generated.DirectoryClientSetPropertiesOptions {
	if o == nil {
		return nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.FileSMBProperties, true)
	permission, permissionKey := exported.FormatPermissions(o.FilePermissions)

	setPropertiesOptions := &generated.DirectoryClientSetPropertiesOptions{
		FileAttributes:    fileAttributes,
		FileChangeTime:    fileChangeTime,
		FileCreationTime:  fileCreationTime,
		FileLastWriteTime: fileLastWriteTime,
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
		FileMode:          o.FileMode,
		Owner:             o.Owner,
		Group:             o.Group,
	}

	if permissionKey != nil && *permissionKey != shared.DefaultPreserveString {
		setPropertiesOptions.FilePermissionFormat = to.Ptr(FilePermissionFormat(shared.DefaultFilePermissionFormat))
	} else if o.FilePermissionFormat != nil {
		setPropertiesOptions.FilePermissionFormat = to.Ptr(FilePermissionFormat(*o.FilePermissionFormat))
	}
	return setPropertiesOptions
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
}

func (o *SetMetadataOptions) format() *generated.DirectoryClientSetMetadataOptions {
	if o == nil {
		return nil
	}

	return &generated.DirectoryClientSetMetadataOptions{
		Metadata: o.Metadata,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// ListFilesAndDirectoriesOptions contains the optional parameters for the Client.NewListFilesAndDirectoriesPager method.
type ListFilesAndDirectoriesOptions struct {
	// Include this parameter to specify one or more datasets to include in the response.
	Include ListFilesInclude
	// Include extended information.
	IncludeExtendedInfo *bool
	// A string value that identifies the portion of the list to be returned with the next list operation. The operation returns
	// a marker value within the response body if the list returned was not complete.
	// The marker value may then be used in a subsequent call to request the next set of list items. The marker value is opaque
	// to the client.
	Marker *string
	// Specifies the maximum number of entries to return. If the request does not specify maxresults, or specifies a value greater
	// than 5,000, the server will return up to 5,000 items.
	MaxResults *int32
	// Filters the results to return only entries whose name begins with the specified prefix.
	Prefix *string
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query for the list of files and directories.
	ShareSnapshot *string
}

// ListFilesInclude specifies one or more datasets to include in the response.
type ListFilesInclude struct {
	Timestamps, ETag, Attributes, PermissionKey bool
}

func (l ListFilesInclude) format() []generated.ListFilesIncludeType {
	if reflect.ValueOf(l).IsZero() {
		return nil
	}

	var include []generated.ListFilesIncludeType

	if l.Timestamps {
		include = append(include, ListFilesIncludeTypeTimestamps)
	}
	if l.ETag {
		include = append(include, ListFilesIncludeTypeETag)
	}
	if l.Attributes {
		include = append(include, ListFilesIncludeTypeAttributes)
	}
	if l.PermissionKey {
		include = append(include, ListFilesIncludeTypePermissionKey)
	}

	return include
}

// FilesAndDirectoriesListSegment - Abstract for entries that can be listed from directory.
type FilesAndDirectoriesListSegment = generated.FilesAndDirectoriesListSegment

// Directory - A listed directory item.
type Directory = generated.Directory

// File - A listed file item.
type File = generated.File

// FileProperty - File properties.
type FileProperty = generated.FileProperty

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
	// Specifies operation should apply to the directory specified in the URI, its files, its subdirectories and their files.
	Recursive *bool
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ListHandlesOptions) format() *generated.DirectoryClientListHandlesOptions {
	if o == nil {
		return nil
	}

	return &generated.DirectoryClientListHandlesOptions{
		Marker:        o.Marker,
		Maxresults:    o.MaxResults,
		Recursive:     o.Recursive,
		Sharesnapshot: o.ShareSnapshot,
	}
}

// Handle - A listed Azure Storage handle item.
type Handle = generated.Handle

// ---------------------------------------------------------------------------------------------------------------------

// ForceCloseHandlesOptions contains the optional parameters for the Client.ForceCloseHandles method.
type ForceCloseHandlesOptions struct {
	// A string value that identifies the portion of the list to be returned with the next list operation. The operation returns
	// a marker value within the response body if the list returned was not complete.
	// The marker value may then be used in a subsequent call to request the next set of list items. The marker value is opaque
	// to the client.
	Marker *string
	// Specifies operation should apply to the directory specified in the URI, its files, its subdirectories and their files.
	Recursive *bool
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ForceCloseHandlesOptions) format() *generated.DirectoryClientForceCloseHandlesOptions {
	if o == nil {
		return nil
	}

	return &generated.DirectoryClientForceCloseHandlesOptions{
		Marker:        o.Marker,
		Recursive:     o.Recursive,
		Sharesnapshot: o.ShareSnapshot,
	}
}
