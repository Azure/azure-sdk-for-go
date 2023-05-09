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

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// ---------------------------------------------------------------------------------------------------------------------

// CreateOptions contains the optional parameters for the Client.Create method.
type CreateOptions struct {
	// The default value is 'Directory' for Attributes and 'now' for CreationTime and LastWriteTime fields in file.SMBProperties.
	FileSMBProperties *file.SMBProperties
	// The default value is 'inherit' for Permission field in file.Permissions.
	FilePermissions *file.Permissions
	// A name-value pair to associate with a file storage object.
	Metadata map[string]*string
}

func (o *CreateOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string, createOptions *generated.DirectoryClientCreateOptions) {
	if o == nil {
		return shared.FileAttributesDirectory, shared.DefaultCurrentTimeString, shared.DefaultCurrentTimeString, &generated.DirectoryClientCreateOptions{
			FilePermission: to.Ptr(shared.DefaultFilePermissionString),
		}
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.FileSMBProperties.Format(true, shared.FileAttributesDirectory, shared.DefaultCurrentTimeString)

	permission, permissionKey := o.FilePermissions.Format(shared.DefaultFilePermissionString)

	createOptions = &generated.DirectoryClientCreateOptions{
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
		Metadata:          o.Metadata,
	}

	return
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
}

func (o *SetPropertiesOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string, setPropertiesOptions *generated.DirectoryClientSetPropertiesOptions) {
	if o == nil {
		return shared.DefaultPreserveString, shared.DefaultPreserveString, shared.DefaultPreserveString, &generated.DirectoryClientSetPropertiesOptions{
			FilePermission: to.Ptr(shared.DefaultPreserveString),
		}
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.FileSMBProperties.Format(true, shared.DefaultPreserveString, shared.DefaultPreserveString)

	permission, permissionKey := o.FilePermissions.Format(shared.DefaultPreserveString)

	setPropertiesOptions = &generated.DirectoryClientSetPropertiesOptions{
		FilePermission:    permission,
		FilePermissionKey: permissionKey,
	}
	return
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
