package azfile

import (
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
	FilePermissions *FilePermissions
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
	Prefix     string // No Prefix header is produced if ""
	MaxResults int32  // 0 means unspecified
}

func (o *ListFilesAndDirectoriesOptions) pointers() (prefix *string, maxResults *int32) {
	if o.Prefix != "" {
		prefix = &o.Prefix
	}
	if o.MaxResults != 0 {
		maxResults = &o.MaxResults
	}
	return
}

type SetDirectoryPropertiesOptions struct {
	FileAttributes *string

	FileCreationTime *time.Time

	FileLastWriteTime *time.Time

	FilePermission *string

	FilePermissionKey *string
}

func (o *SetDirectoryPropertiesOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string, options *DirectorySetPropertiesOptions) {
	//if o == nil {
	//	return DefaultFileAttributes, DefaultCurrentTimeString, DefaultCurrentTimeString, &DirectorySetPropertiesOptions{
	//		FilePermission:    &DefaultFilePermissionStr,
	//	}
	//}
	panic("Write")
}

type SetDirectoryMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string
}

func (o *SetDirectoryMetadataOptions) format() *DirectorySetMetadataOptions {
	if o == nil {
		return nil
	}

	return &DirectorySetMetadataOptions{
		Metadata: o.Metadata,
	}
}
