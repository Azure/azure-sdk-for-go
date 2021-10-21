package azfile

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"io"
	"strings"
	"time"
)

const (
	// FileMaxUploadRangeBytes indicates the maximum number of bytes that can be sent in a call to UploadRange.
	FileMaxUploadRangeBytes = 4 * 1024 * 1024 // 4MB

	ISO8601 = "2006-01-02T15:04:05.0000000Z"
)

var (
	// DefaultFilePermissionStr is a constant for all intents and purposes.
	// But you can't take the address of a constant string, so it's a variable.
	// Inherit inherits permissions from the parent folder (default when creating files/folders)
	DefaultFilePermissionStr = "inherit"

	// DefaultCurrentTimeString sets creation/last write times to now
	DefaultCurrentTimeString = "now"

	// DefaultPreserveString preserves old permissions on the file/folder (default when updating properties)
	DefaultPreserveString = "preserve"

	// DefaultFileAttributes is defaults for file attributes
	DefaultFileAttributes = "None"

	DefaultPermissionCopyMode = PermissionCopyModeType("")
)

//----------------------------------------------------------------------------------------------------------------------

type Permissions struct {
	// If specified the permission (security descriptor) shall be set for the directory/file. This header can be used if Permission size is <= 8KB, else x-ms-file-permission-key
	// header shall be used. Default value: Inherit. If SDDL is specified as input, it must have owner, group and dacl. Note: Only one of the x-ms-file-permission
	// or x-ms-file-permission-key should be specified.
	PermissionStr *string
	// Key of the permission to be set for the directory/file. Note: Only one of the x-ms-file-permission or x-ms-file-permission-key should be specified.
	PermissionKey *string
}

func (p *Permissions) format(defaultFilePermissionStr *string) (filePermission *string, filePermissionKey *string, err error) {
	if p == nil {
		return defaultFilePermissionStr, nil, nil
	}
	filePermission = defaultFilePermissionStr
	if p.PermissionStr != nil {
		filePermission = p.PermissionStr
	}

	if p.PermissionKey != nil {
		if filePermission == defaultFilePermissionStr {
			filePermission = nil
		} else if filePermission != nil {
			return nil, nil, errors.New("only permission string OR permission key may be used")
		}

		filePermissionKey = p.PermissionKey
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

// SMBProperties defines a struct that takes in optional parameters regarding SMB/NTFS properties.
// When you pass this into another function (Either literally or via FileHTTPHeaders), the response will probably fit inside SMBPropertyAdapter.
// Nil values of the properties are inferred to be preserved (or when creating, use defaults). Clearing a value can be done by supplying an empty item instead of nil.
type SMBProperties struct {
	FileAttributes *FileAttributeFlags
	// A UTC time-date string is specified below. A value of 'now' defaults to now. 'preserve' defaults to preserving the old case.
	FileCreationTime *time.Time

	FileLastWriteTime *time.Time
}

func (sp *SMBProperties) format(isDir bool, defaultFileAttributes, defaultCurrentTimeString string) (fileAttributes string, creationTime string, lastWriteTime string) {
	if sp == nil {
		return defaultFileAttributes, defaultCurrentTimeString, defaultCurrentTimeString
	}

	fileAttributes = defaultFileAttributes
	if sp.FileAttributes != nil {
		fileAttributes = sp.FileAttributes.String()
		if isDir && strings.ToLower(fileAttributes) != "none" { // must test string, not sp.FileAttributes, since it may contain set bits that we don't convert
			// Directories need to have this attribute included, if setting any attributes.
			// We don't expose it in FileAttributes because it doesn't do anything useful to consumers of
			// this SDK. And because it always needs to be set for directories and not for non-directories,
			// so it makes sense to automate that here.
			fileAttributes += "|Directory"
		}
	}

	creationTime = defaultCurrentTimeString
	if sp.FileCreationTime != nil {
		creationTime = sp.FileCreationTime.UTC().Format(ISO8601)
	}

	lastWriteTime = defaultCurrentTimeString
	if sp.FileLastWriteTime != nil {
		lastWriteTime = sp.FileLastWriteTime.UTC().Format(ISO8601)
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------

type CreateFileOptions struct {
	FileContentLength *int64

	Metadata map[string]string

	FileHTTPHeaders *FileHTTPHeaders

	FilePermissions *Permissions

	// In Windows, a 32 bit file attributes integer exists. This is that.
	SMBProperties *SMBProperties

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *CreateFileOptions) format() (fileContentLength int64, fileAttributes string, fileCreationTime string, fileLastWriteTime string, fileCreateOptions *FileCreateOptions, fileHTTPHeaders *FileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions, err error) {
	if o == nil {
		return int64(0), DefaultFileAttributes, DefaultCurrentTimeString, DefaultCurrentTimeString, &FileCreateOptions{FilePermission: to.StringPtr(DefaultFilePermissionStr)}, nil, nil, nil
	}
	fileContentLength = 0

	if o.FileContentLength != nil {
		fileContentLength = *(o.FileContentLength)
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false, DefaultFileAttributes, DefaultCurrentTimeString)

	filePermission, filePermissionKey, err := o.FilePermissions.format(&DefaultFilePermissionStr)
	if err != nil {
		return
	}

	fileCreateOptions = &FileCreateOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
		Metadata:          o.Metadata,
	}

	fileHTTPHeaders = o.FileHTTPHeaders
	leaseAccessConditions = o.LeaseAccessConditions

	return
}

//----------------------------------------------------------------------------------------------------------------------

type StartFileCopyOptions struct {
	FilePermissions *Permissions
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string

	CopyFileSmbInfo *CopyFileSmbInfo

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *StartFileCopyOptions) format() (fileStartCopyOptions *FileStartCopyOptions, copyFileSmbInfo *CopyFileSmbInfo, leaseAccessConditions *LeaseAccessConditions, err error) {
	if o == nil {
		return
	}

	filePermission, filePermissionKey, err := o.FilePermissions.format(nil)
	if err != nil {
		return
	}

	fileStartCopyOptions = &FileStartCopyOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
		Metadata:          o.Metadata,
	}

	copyFileSmbInfo = o.CopyFileSmbInfo
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

//----------------------------------------------------------------------------------------------------------------------

type AbortFileCopyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *AbortFileCopyOptions) format() (fileAbortCopyOptions *FileAbortCopyOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}

	return nil, o.LeaseAccessConditions
}

//----------------------------------------------------------------------------------------------------------------------

type DownloadFileOptions struct {
	// When this header is set to true and specified together with the Range header,
	// the service returns the MD5 hash for the range, as long as the range is less than or equal to 4 MB in size.
	RangeGetContentMD5 *bool

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *DownloadFileOptions) format(offset, count *int64) (fileDownloadOptions *FileDownloadOptions, leaseAccessConditions *LeaseAccessConditions) {
	fileDownloadOptions = &FileDownloadOptions{
		Range: toRange(offset, count),
	}

	if o != nil {
		fileDownloadOptions.RangeGetContentMD5 = o.RangeGetContentMD5
		leaseAccessConditions = o.LeaseAccessConditions
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------

type DeleteFileOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *DeleteFileOptions) format() (fileDownloadOptions *FileDeleteOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

//----------------------------------------------------------------------------------------------------------------------

type GetFilePropertiesOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetFilePropertiesOptions) format() (fileGetPropertiesOptions *FileGetPropertiesOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}
	fileGetPropertiesOptions = &FileGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

//----------------------------------------------------------------------------------------------------------------------

type SetFileHTTPHeadersOptions struct {
	// Resizes a file to the specified size.
	// If the specified byte value is less than the current size of the file,
	// then all ranges above the specified byte value are cleared.
	FileContentLength *int64

	FilePermissions *Permissions

	SMBProperties *SMBProperties

	FileHTTPHeaders *FileHTTPHeaders

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetFileHTTPHeadersOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	fileSetHTTPHeadersOptions *FileSetHTTPHeadersOptions, fileHTTPHeaders *FileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions, err error) {

	fileAttributes, fileCreationTime, fileLastWriteTime = "preserve", "preserve", "preserve"

	filePermission, filePermissionKey, err := o.FilePermissions.format(&DefaultPreserveString)
	if err != nil {
		return
	}

	fileSetHTTPHeadersOptions = &FileSetHTTPHeadersOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
	}

	if o != nil {
		fileSetHTTPHeadersOptions.FileContentLength = o.FileContentLength
		fileHTTPHeaders = o.FileHTTPHeaders
		leaseAccessConditions = o.LeaseAccessConditions
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

type SetFileMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetFileMetadataOptions) format(metadata map[string]string) (fileSetMetadataOptions *FileSetMetadataOptions, leaseAccessConditions *LeaseAccessConditions) {
	fileSetMetadataOptions = &FileSetMetadataOptions{
		Metadata: metadata,
	}

	if o != nil {
		leaseAccessConditions = o.LeaseAccessConditions
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

type ResizeFileOptions struct {
	LeaseAccessCondition *LeaseAccessConditions
}

func (o *ResizeFileOptions) format(contentLength int64) (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	fileSetHTTPHeadersOptions *FileSetHTTPHeadersOptions, fileHTTPHeaders *FileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions) {

	fileAttributes, fileCreationTime, fileLastWriteTime = DefaultPreserveString, DefaultPreserveString, DefaultPreserveString

	fileSetHTTPHeadersOptions = &FileSetHTTPHeadersOptions{
		FileContentLength: to.Int64Ptr(contentLength),
		FilePermission:    &DefaultPreserveString,
	}

	if o != nil {
		leaseAccessConditions = o.LeaseAccessCondition
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

type UploadFileRangeOptions struct {
	// An MD5 hash of the content. This hash is used to verify the integrity of the data during transport.
	// When the Content-MD5 header is specified, the File service compares the hash of the content that has arrived with the header value that was sent.
	// If the two hashes do not match, the operation will fail with error code 400 (Bad Request).
	ContentMD5 []byte

	TransactionalMD5 []byte

	FileRangeWrite *FileRangeWriteType

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *UploadFileRangeOptions) format(offset int64, body io.ReadSeekCloser) (rangeParam string, fileRangeWrite FileRangeWriteType, contentLength int64,
	fileUploadRangeOptions *FileUploadRangeOptions, leaseAccessConditions *LeaseAccessConditions, err error) {
	if offset < 0 {
		err = errors.New("invalid argument, offset must be >= 0")
		return
	}
	if body == nil {
		err = errors.New("invalid argument, body must not be nil")
		return
	}

	count := int64(CountToEnd)
	count, err = validateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return
	}

	if count == 0 {
		err = errors.New("invalid argument, body must contain readable data whose size is > 0")
		return
	}

	rangeParam = *toRange(to.Int64Ptr(offset), to.Int64Ptr(count))
	fileRangeWrite = FileRangeWriteTypeUpdate
	contentLength = count

	var contentMD5 []byte

	if o != nil {
		if o.FileRangeWrite != nil {
			fileRangeWrite = *o.FileRangeWrite
		}
		if o.ContentMD5 != nil {
			contentMD5 = o.ContentMD5
		}
		leaseAccessConditions = o.LeaseAccessConditions
	}

	fileUploadRangeOptions = &FileUploadRangeOptions{
		ContentMD5:   contentMD5,
		Optionalbody: body,
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------

type UploadFileRangeFromURLOptions struct {

	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCRC64 []byte

	SourceModifiedAccessConditions *SourceModifiedAccessConditions

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *UploadFileRangeFromURLOptions) format(sourceURL string, sourceOffset, destinationOffset, count *int64) (rangeParam string, copySource string,
	contentLength int64, fileUploadRangeFromURLOptions *FileUploadRangeFromURLOptions, sourceModifiedAccessConditions *SourceModifiedAccessConditions,
	leaseAccessConditions *LeaseAccessConditions) {

	rangeParam = *toRange(destinationOffset, count)
	copySource = sourceURL
	contentLength = *count
	fileUploadRangeFromURLOptions = &FileUploadRangeFromURLOptions{SourceRange: toRange(sourceOffset, count)}
	if o != nil {
		fileUploadRangeFromURLOptions.SourceContentCRC64 = o.SourceContentCRC64
		sourceModifiedAccessConditions = o.SourceModifiedAccessConditions
		leaseAccessConditions = o.LeaseAccessConditions
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

type ClearFileRangeOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ClearFileRangeOptions) format(offset, count *int64) (rangeParam string, fileRangeWrite FileRangeWriteType, contentLength int64,
	fileUploadRangeOptions *FileUploadRangeOptions, leaseAccessConditions *LeaseAccessConditions, err error) {

	if *offset < 0 || *count <= 0 {
		err = errors.New("invalid argument: either offset is < 0 or count <= 0")
		return
	}
	rangeParam = *toRange(offset, count)
	fileRangeWrite = FileRangeWriteTypeClear
	contentLength = 0

	if o != nil {
		leaseAccessConditions = o.LeaseAccessConditions
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------

type GetFileRangeListOptions struct {
	// The previous snapshot parameter is an opaque DateTime value that, when present, specifies the previous snapshot.
	PrevShareSnapshot *string
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetFileRangeListOptions) format(offset, count *int64) (fileGetRangeListOptions *FileGetRangeListOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}

	return &FileGetRangeListOptions{
		Prevsharesnapshot: o.PrevShareSnapshot,
		Range:             toRange(offset, count),
		Sharesnapshot:     o.ShareSnapshot,
	}, o.LeaseAccessConditions
}
