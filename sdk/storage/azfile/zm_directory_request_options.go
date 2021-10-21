// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"time"
)

// SMBPropertyHolder is an interface designed for SMBPropertyAdapter, to identify valid response types for adapting.
type SMBPropertyHolder interface {
	FileCreationTime() string
	FileLastWriteTime() string
	FileAttributes() string
}

// SMBPropertyAdapter is a wrapper struct that automatically converts the string outputs of FileAttributes, FileCreationTime and FileLastWrite time to time.Time.
// It is _not_ error resistant. It is expected that the response you're inserting into this is a valid response.
// File and directory calls that return such properties are: GetProperties, SetProperties, Create
// File Downloads also return such properties. Insert other response types at your peril.
type SMBPropertyAdapter struct {
	PropertySource SMBPropertyHolder
}

func (s *SMBPropertyAdapter) convertISO8601(input string) time.Time {
	t, err := time.Parse(ISO8601, input)

	if err != nil {
		// This should literally never happen if this struct is used correctly.
		panic("SMBPropertyAdapter expects a successful response fitting the SMBPropertyHolder interface. Failed to parse time:\n" + err.Error())
	}

	return t
}

func (s *SMBPropertyAdapter) FileCreationTime() time.Time {
	return s.convertISO8601(s.PropertySource.FileCreationTime()).UTC()
}

func (s *SMBPropertyAdapter) FileLastWriteTime() time.Time {
	return s.convertISO8601(s.PropertySource.FileLastWriteTime()).UTC()
}

func (s *SMBPropertyAdapter) FileAttributes() FileAttributeFlags {
	return ParseFileAttributeFlagsString(s.PropertySource.FileAttributes())
}

type CreateDirectoryOptions struct {
	SMBProperties   *SMBProperties
	FilePermissions *Permissions
	Metadata        map[string]string
}

func (o *CreateDirectoryOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string, directoryCreateOptions *DirectoryCreateOptions, err error) {
	if o == nil {
		return DefaultFileAttributes, DefaultCurrentTimeString, DefaultCurrentTimeString, &DirectoryCreateOptions{FilePermission: to.StringPtr(DefaultFilePermissionStr)}, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false, DefaultFileAttributes, DefaultCurrentTimeString)

	filePermission, filePermissionKey, err := o.FilePermissions.format(&DefaultFilePermissionStr)
	if err != nil {
		return
	}

	directoryCreateOptions = &DirectoryCreateOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
		Metadata:          o.Metadata,
	}

	return
}

type DeleteDirectoryOptions struct {
}

//----------------------------------------------------------------------------------------------------------------------

type GetDirectoryPropertiesOptions struct {
	ShareSnapshot *string
}

func (o *GetDirectoryPropertiesOptions) format() *DirectoryGetPropertiesOptions {
	if o == nil {
		return nil
	}
	directoryGetPropertiesOptions := &DirectoryGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}
	return directoryGetPropertiesOptions
}

// ListFilesAndDirectoriesOptions defines options available when calling ListFilesAndDirectoriesSegment.
type ListFilesAndDirectoriesOptions struct {
	// Specifies the maximum number of entries to return. If the request does not specify maxresults, or specifies a value greater than 5,000, the server will
	// return up to 5,000 items.
	MaxResults *int32
	// Filters the results to return only entries whose name begins with the specified prefix.
	Prefix *string
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *ListFilesAndDirectoriesOptions) format(marker *string) *DirectoryListFilesAndDirectoriesSegmentOptions {
	if o == nil {
		return &DirectoryListFilesAndDirectoriesSegmentOptions{Marker: marker}
	}
	return &DirectoryListFilesAndDirectoriesSegmentOptions{
		Marker:        marker,
		Maxresults:    o.MaxResults,
		Prefix:        o.Prefix,
		Sharesnapshot: o.ShareSnapshot,
	}
}

type SetDirectoryPropertiesOptions struct {
	SMBProperties *SMBProperties

	FilePermissions *Permissions
}

func (o *SetDirectoryPropertiesOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string, directorySetPropertiesOptions *DirectorySetPropertiesOptions) {
	if o == nil {
		return DefaultPreserveString, DefaultPreserveString, DefaultPreserveString, nil
	}
	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false, DefaultPreserveString, DefaultPreserveString)

	filePermission, filePermissionKey, err := o.FilePermissions.format(&DefaultPreserveString)
	if err != nil {
		return
	}
	directorySetPropertiesOptions = &DirectorySetPropertiesOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
	}
	return
}

type SetDirectoryMetadataOptions struct {
	// A name-value pair to associate with a file storage object.

}

func (o *SetDirectoryMetadataOptions) format(metadata map[string]string) (*DirectorySetMetadataOptions, error) {
	if metadata == nil {
		return nil, errors.New("metadata cannot be nil")

	}

	return &DirectorySetMetadataOptions{
		Metadata: metadata,
	}, nil
}
