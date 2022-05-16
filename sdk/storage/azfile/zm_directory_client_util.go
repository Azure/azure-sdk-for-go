//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

// SMBPropertyHolder is an interface designed for SMBPropertyAdapter, to identify valid responseBody types for adapting.
type SMBPropertyHolder interface {
	FileCreationTime() string
	FileLastWriteTime() string
	FileAttributes() string
}

// SMBPropertyAdapter is a wrapper struct that automatically converts the string outputs of FileAttributes, FileCreationTime and FileLastWrite time to time.Time.
// It is _not_ error resistant. It is expected that the responseBody you're inserting into this is a valid responseBody.
// File and directory calls that return such properties are: GetProperties, SetProperties, Create
// File Downloads also return such properties. Insert other responseBody types at your peril.
type SMBPropertyAdapter struct {
	PropertySource SMBPropertyHolder
}

func (s *SMBPropertyAdapter) convertISO8601(input string) time.Time {
	t, err := time.Parse(ISO8601, input)

	if err != nil {
		// This should literally never happen if this struct is used correctly.
		panic("SMBPropertyAdapter expects a successful responseBody fitting the SMBPropertyHolder interface. Failed to parse time:\n" + err.Error())
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

// ---------------------------------------------------------------------------------------------------------------------

type DirectoryCreateOptions struct {
	SMBProperties   *SMBProperties
	FilePermissions *FilePermissions
	Metadata        map[string]string
}

func (o *DirectoryCreateOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string, createOptions *directoryClientCreateOptions, err error) {
	if o == nil {
		return DefaultFileAttributes, DefaultCurrentTimeString, DefaultCurrentTimeString, &directoryClientCreateOptions{
			FilePermission: to.Ptr(DefaultFilePermissionString)}, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false, DefaultFileAttributes, DefaultCurrentTimeString)

	filePermission, filePermissionKey, err := o.FilePermissions.format(&DefaultFilePermissionString)
	if err != nil {
		return
	}

	createOptions = &directoryClientCreateOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
		Metadata:          o.Metadata,
	}

	return
}

type DirectoryCreateResponse struct {
	directoryClientCreateResponse
}

func toDirectoryCreateResponse(resp directoryClientCreateResponse) DirectoryCreateResponse {
	return DirectoryCreateResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type DirectoryDeleteOptions struct {
}

func (o *DirectoryDeleteOptions) format() *directoryClientDeleteOptions {
	return nil
}

type DirectoryDeleteResponse struct {
	directoryClientDeleteResponse
}

func toDirectoryDeleteResponse(resp directoryClientDeleteResponse) DirectoryDeleteResponse {
	return DirectoryDeleteResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type DirectoryGetPropertiesOptions struct {
	ShareSnapshot *string
}

func (o *DirectoryGetPropertiesOptions) format() *directoryClientGetPropertiesOptions {
	if o == nil {
		return nil
	}

	return &directoryClientGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}
}

type DirectoryGetPropertiesResponse struct {
	directoryClientGetPropertiesResponse
}

func toDirectoryGetPropertiesResponse(resp directoryClientGetPropertiesResponse) DirectoryGetPropertiesResponse {
	return DirectoryGetPropertiesResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type DirectorySetPropertiesOptions struct {
	SMBProperties *SMBProperties

	FilePermissions *FilePermissions
}

func (o *DirectorySetPropertiesOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string, setPropertiesOptions *directoryClientSetPropertiesOptions) {
	if o == nil {
		return DefaultPreserveString, DefaultPreserveString, DefaultPreserveString, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false, DefaultPreserveString, DefaultPreserveString)

	filePermission, filePermissionKey, err := o.FilePermissions.format(&DefaultPreserveString)
	if err != nil {
		return
	}
	setPropertiesOptions = &directoryClientSetPropertiesOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
	}
	return
}

type DirectorySetPropertiesResponse struct {
	directoryClientSetPropertiesResponse
}

func toDirectorySetPropertiesResponse(resp directoryClientSetPropertiesResponse) DirectorySetPropertiesResponse {
	return DirectorySetPropertiesResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type DirectorySetMetadataOptions struct {
	// A name-value pair to associate with a file storage object.

}

func (o *DirectorySetMetadataOptions) format(metadata map[string]string) (*directoryClientSetMetadataOptions, error) {
	if metadata == nil {
		return nil, errors.New("metadata cannot be nil")
	}

	return &directoryClientSetMetadataOptions{
		Metadata: metadata,
	}, nil
}

type DirectorySetMetadataResponse struct {
	directoryClientSetMetadataResponse
}

func toDirectorySetMetadataResponse(resp directoryClientSetMetadataResponse) DirectorySetMetadataResponse {
	return DirectorySetMetadataResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

// DirectoryListFilesAndDirectoriesOptions defines options available when calling ListFilesAndDirectoriesSegment.
type DirectoryListFilesAndDirectoriesOptions struct {
	Marker *string
	// Specifies the maximum number of entries to return. If the request does not specify maxresults, or specifies a value greater than 5,000, the server will
	// return up to 5,000 items.
	MaxResults *int32
	// Filters the results to return only entries whose name begins with the specified prefix.
	Prefix *string
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *DirectoryListFilesAndDirectoriesOptions) format() *directoryClientListFilesAndDirectoriesSegmentOptions {
	if o == nil {
		return nil
	}
	return &directoryClientListFilesAndDirectoriesSegmentOptions{
		Marker:        o.Marker,
		Maxresults:    o.MaxResults,
		Prefix:        o.Prefix,
		Sharesnapshot: o.ShareSnapshot,
	}
}

type DirectoryListFilesAndDirectoriesResponse struct {
	directoryClientListFilesAndDirectoriesSegmentResponse
}

func toDirectoryListFilesAndDirectoriesResponse(resp directoryClientListFilesAndDirectoriesSegmentResponse) DirectoryListFilesAndDirectoriesResponse {
	return DirectoryListFilesAndDirectoriesResponse{resp}
}
