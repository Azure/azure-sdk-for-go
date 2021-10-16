package azfile

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"time"
)

type CreateDirectoryOptions struct {
	SMBProperties   *SMBProperties
	FilePermissions *FilePermissions
	Metadata        map[string]string
}

func (o *CreateDirectoryOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string, directoryCreateOptions *DirectoryCreateOptions, err error) {
	if o == nil {
		return defaultFileAttributes, defaultCurrentTimeString, defaultCurrentTimeString, &DirectoryCreateOptions{FilePermission: to.StringPtr(defaultFilePermissionStr)}, nil
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false)

	filePermission, filePermissionKey, err := o.FilePermissions.format()
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
	//	return defaultFileAttributes, defaultCurrentTimeString, defaultCurrentTimeString, &DirectorySetPropertiesOptions{
	//		FilePermission:    &defaultFilePermissionStr,
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
