// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

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
	// NFS only.
	FileNFSProperties *file.NFSProperties
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
	// SMB only. How attributes and permissions should be set on the directory.
	// New: automatically adds the ARCHIVE file attribute flag and uses Windows create file permissions semantics (ex: inherit from parent).
	// Restore: does not modify file attribute flag and uses Windows update file permissions semantics.
	// If Restore is specified, the file permission must also be provided, otherwise PropertySemantics will default to New.
	FilePropertySemantics *PropertySemantics
}

func (o *CreateOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.DirectoryClientCreateOptions {
	opts := &generated.DirectoryClientCreateOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Metadata = o.Metadata

	if o.FileNFSProperties != nil {
		fileCreationTime, fileLastWriteTime := exported.FormatNFSProperties(o.FileNFSProperties, true)

		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.FileMode = o.FileNFSProperties.FileMode
		opts.Group = o.FileNFSProperties.Group
		opts.Owner = o.FileNFSProperties.Owner
	} else {
		fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.FileSMBProperties, true)
		permission, permissionKey := exported.FormatPermissions(o.FilePermissions)

		opts.FileAttributes = fileAttributes
		opts.FileChangeTime = fileChangeTime
		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.FilePermission = permission
		opts.FilePermissionKey = permissionKey

		if permissionKey != nil && *permissionKey != shared.DefaultFilePermissionString {
			opts.FilePermissionFormat = to.Ptr(FilePermissionFormat(shared.DefaultFilePermissionFormat))
		} else if o.FilePermissionFormat != nil {
			opts.FilePermissionFormat = to.Ptr(FilePermissionFormat(*o.FilePermissionFormat))
		}
		if o.FilePropertySemantics != nil {
			opts.FilePropertySemantics = o.FilePropertySemantics
		}
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for the Client.Delete method.
type DeleteOptions struct {
	// placeholder for future options
}

func (o *DeleteOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.DirectoryClientDeleteOptions {
	return &generated.DirectoryClientDeleteOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
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

func (o *RenameOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool, allowSourceTrailingDot *bool) *generated.DirectoryClientRenameOptions {
	opts := &generated.DirectoryClientRenameOptions{
		FileRequestIntent:      fileRequestIntent,
		AllowTrailingDot:       allowTrailingDot,
		AllowSourceTrailingDot: allowSourceTrailingDot,
	}
	if o == nil {
		return opts
	}

	fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.FileSMBProperties, true)
	permission, permissionKey := exported.FormatPermissions(o.FilePermissions)

	opts.FilePermission = permission
	opts.FilePermissionKey = permissionKey
	opts.IgnoreReadOnly = o.IgnoreReadOnly
	opts.Metadata = o.Metadata
	opts.ReplaceIfExists = o.ReplaceIfExists

	if permissionKey != nil && *permissionKey != shared.DefaultPreserveString {
		opts.FilePermissionFormat = to.Ptr(FilePermissionFormat(shared.DefaultFilePermissionFormat))
	} else if o.FilePermissionFormat != nil {
		opts.FilePermissionFormat = to.Ptr(FilePermissionFormat(*o.FilePermissionFormat))
	}

	opts.FileAttributes = fileAttributes
	opts.FileChangeTime = fileChangeTime
	opts.FileCreationTime = fileCreationTime
	opts.FileLastWriteTime = fileLastWriteTime

	if o.DestinationLeaseAccessConditions != nil {
		opts.DestinationLeaseID = o.DestinationLeaseAccessConditions.DestinationLeaseID
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties method.
type GetPropertiesOptions struct {
	// ShareSnapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query for the directory properties.
	ShareSnapshot *string
}

func (o *GetPropertiesOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.DirectoryClientGetPropertiesOptions {
	opts := &generated.DirectoryClientGetPropertiesOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Sharesnapshot = o.ShareSnapshot

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// SetPropertiesOptions contains the optional parameters for the Client.SetProperties method.
type SetPropertiesOptions struct {
	// The default value is 'preserve' for Attributes, CreationTime and LastWriteTime fields in file.SMBProperties.
	FileSMBProperties *file.SMBProperties
	// NFS only.
	FileNFSProperties *file.NFSProperties
	// The default value is 'preserve' for Permission field in file.Permissions.
	FilePermissions *file.Permissions
	// FilePermissionFormat contains the format of the file permissions, Can be sddl (Default) or Binary.
	FilePermissionFormat *FilePermissionFormat
}

func (o *SetPropertiesOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.DirectoryClientSetPropertiesOptions {
	opts := &generated.DirectoryClientSetPropertiesOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	if o.FileNFSProperties != nil {
		fileCreationTime, fileLastWriteTime := exported.FormatNFSProperties(o.FileNFSProperties, true)

		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.FileMode = o.FileNFSProperties.FileMode
		opts.Owner = o.FileNFSProperties.Owner
		opts.Group = o.FileNFSProperties.Group
	} else {
		fileAttributes, fileCreationTime, fileLastWriteTime, fileChangeTime := exported.FormatSMBProperties(o.FileSMBProperties, true)
		permission, permissionKey := exported.FormatPermissions(o.FilePermissions)

		opts.FileAttributes = fileAttributes
		opts.FileChangeTime = fileChangeTime
		opts.FileCreationTime = fileCreationTime
		opts.FileLastWriteTime = fileLastWriteTime
		opts.FilePermission = permission
		opts.FilePermissionKey = permissionKey

		if permissionKey != nil && *permissionKey != shared.DefaultPreserveString {
			opts.FilePermissionFormat = to.Ptr(FilePermissionFormat(shared.DefaultFilePermissionFormat))
		} else if o.FilePermissionFormat != nil {
			opts.FilePermissionFormat = to.Ptr(FilePermissionFormat(*o.FilePermissionFormat))
		}
	}

	return opts
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions contains the optional parameters for the Client.SetMetadata method.
type SetMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
}

func (o *SetMetadataOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.DirectoryClientSetMetadataOptions {
	opts := &generated.DirectoryClientSetMetadataOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Metadata = o.Metadata
	return opts
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

func (o *ListHandlesOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.DirectoryClientListHandlesOptions {
	opts := &generated.DirectoryClientListHandlesOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Marker = o.Marker
	opts.Maxresults = o.MaxResults
	opts.Recursive = o.Recursive
	opts.Sharesnapshot = o.ShareSnapshot

	return opts
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

func (o *ForceCloseHandlesOptions) format(fileRequestIntent *generated.ShareTokenIntent, allowTrailingDot *bool) *generated.DirectoryClientForceCloseHandlesOptions {
	opts := &generated.DirectoryClientForceCloseHandlesOptions{
		FileRequestIntent: fileRequestIntent,
		AllowTrailingDot:  allowTrailingDot,
	}
	if o == nil {
		return opts
	}

	opts.Marker = o.Marker
	opts.Recursive = o.Recursive
	opts.Sharesnapshot = o.ShareSnapshot

	return opts
}
