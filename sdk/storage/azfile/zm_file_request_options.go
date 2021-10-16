package azfile

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"io"
	"net/url"
	"strings"
	"time"
)

const (
	// FileMaxUploadRangeBytes indicates the maximum number of bytes that can be sent in a call to UploadRange.
	FileMaxUploadRangeBytes = 4 * 1024 * 1024 // 4MB

	ISO8601 = "2006-01-02T15:04:05.0000000Z"
)

var (
	// For all intents and purposes, this is a constant.
	// But you can't take the address of a constant string, so it's a variable.
	// Inherit inherits permissions from the parent folder (default when creating files/folders)
	defaultFilePermissionStr = "inherit"

	// Sets creation/last write times to now
	defaultCurrentTimeString = "now"

	// Preserves old permissions on the file/folder (default when updating properties)
	defaultPreserveString = "preserve"

	// Defaults for file attributes
	defaultFileAttributes = "None"
)

type FilePermissions struct {
	// If specified the permission (security descriptor) shall be set for the directory/file. This header can be used if Permission size is <= 8KB, else x-ms-file-permission-key
	// header shall be used. Default value: Inherit. If SDDL is specified as input, it must have owner, group and dacl. Note: Only one of the x-ms-file-permission
	// or x-ms-file-permission-key should be specified.
	FilePermissionStr *string
	// Key of the permission to be set for the directory/file. Note: Only one of the x-ms-file-permission or x-ms-file-permission-key should be specified.
	FilePermissionKey *string
}

func (fp *FilePermissions) format() (filePermission *string, filePermissionKey *string, err error) {
	if fp == nil {
		return &defaultFilePermissionStr, nil, nil
	}
	filePermission = &defaultFilePermissionStr
	if fp.FilePermissionStr != nil {
		filePermission = fp.FilePermissionStr
	}

	if fp.FilePermissionKey != nil {
		if filePermission == &defaultFilePermissionStr {
			filePermission = nil
		} else if filePermission != nil {
			return nil, nil, errors.New("only permission string OR permission key may be used")
		}

		filePermissionKey = fp.FilePermissionKey
	}
	return
}

// SMBProperties defines a struct that takes in optional parameters regarding SMB/NTFS properties.
// When you pass this into another function (Either literally or via FileHTTPHeaders), the response will probably fit inside SMBPropertyAdapter.
// Nil values of the properties are inferred to be preserved (or when creating, use defaults). Clearing a value can be done by supplying an empty item instead of nil.
type SMBProperties struct {
	FileAttributes *FileAttributeFlags
	// A UTC time-date string is specified below. A value of 'now' defaults to now. 'preserve' defaults to preserving the old case.
	FileCreationTime *time.Time

	FileLastWriteTime *time.Time
}

func (sp *SMBProperties) format(isDir bool) (fileAttributes string, creationTime string, lastWriteTime string) {
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

type CreateFileOptions struct {
	FileContentLength *int64

	Metadata map[string]string

	FileHTTPHeaders *FileHTTPHeaders

	FilePermissions *FilePermissions

	// In Windows, a 32 bit file attributes integer exists. This is that.
	SMBProperties *SMBProperties

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *CreateFileOptions) format() (fileContentLength int64, fileAttributes string, fileCreationTime string, fileLastWriteTime string, fileCreateOptions *FileCreateOptions, fileHTTPHeaders *FileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions, err error) {
	if o == nil {
		return int64(0), defaultFileAttributes, defaultCurrentTimeString, defaultCurrentTimeString, &FileCreateOptions{FilePermission: to.StringPtr(defaultFilePermissionStr)}, nil, nil, nil
	}
	fileContentLength = 0

	if o.FileContentLength != nil {
		fileContentLength = *(o.FileContentLength)
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false)

	filePermission, filePermissionKey, err := o.FilePermissions.format()
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
	CopySource string

	FilePermissions *FilePermissions
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string

	CopyFileSmbInfo *CopyFileSmbInfo

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *StartFileCopyOptions) format() (fileStartCopyOptions *FileStartCopyOptions, copyFileSmbInfo *CopyFileSmbInfo, leaseAccessConditions *LeaseAccessConditions, err error) {
	if o == nil {
		return
	}

	filePermission, filePermissionKey, err := o.FilePermissions.format()
	if err != nil {
		return nil, nil, nil, err
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
	Offset *int64

	Count *int64
	// When this header is set to true and specified together with the Range header,
	// the service returns the MD5 hash for the range, as long as the range is less than or equal to 4 MB in size.
	RangeGetContentMD5 *bool

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *DownloadFileOptions) format() (fileDownloadOptions *FileDownloadOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}
	fileDownloadOptions = &FileDownloadOptions{
		Range:              getRangeParam(o.Offset, o.Count),
		RangeGetContentMD5: o.RangeGetContentMD5,
	}
	leaseAccessConditions = o.LeaseAccessConditions
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

	FilePermissions *FilePermissions

	SMBProperties *SMBProperties

	FileHTTPHeaders *FileHTTPHeaders

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetFileHTTPHeadersOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	fileSetHTTPHeadersOptions *FileSetHTTPHeadersOptions, fileHTTPHeaders *FileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions, err error) {

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false)

	filePermission, filePermissionKey, err := o.FilePermissions.format()
	if err != nil {
		return
	}

	fileSetHTTPHeadersOptions = &FileSetHTTPHeadersOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
		FileContentLength: o.FileContentLength,
	}

	fileHTTPHeaders = o.FileHTTPHeaders
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

//----------------------------------------------------------------------------------------------------------------------

type SetFileMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata              map[string]string
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetFileMetadataOptions) format() (fileSetMetadataOptions *FileSetMetadataOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}

	fileSetMetadataOptions = &FileSetMetadataOptions{
		Metadata: o.Metadata,
	}
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

//----------------------------------------------------------------------------------------------------------------------

type ResizeFileOptions struct {
	Length               *int64
	LeaseAccessCondition *LeaseAccessConditions
}

func (o *ResizeFileOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	fileSetHTTPHeadersOptions *FileSetHTTPHeadersOptions, fileHTTPHeaders *FileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions) {
	fileAttributes, fileCreationTime, fileLastWriteTime = "preserve", "preserve", "preserve"

	var contentLength *int64
	if o != nil && o.Length != nil {
		contentLength = o.Length
	}

	fileSetHTTPHeadersOptions = &FileSetHTTPHeadersOptions{
		FileContentLength: contentLength,
		FilePermission:    &defaultFilePermissionStr,
	}

	fileHTTPHeaders = nil
	leaseAccessConditions = o.LeaseAccessCondition
	return
}

//----------------------------------------------------------------------------------------------------------------------

type UploadFileRangeOptions struct {
	Offset *int64

	// An MD5 hash of the content. This hash is used to verify the integrity of the data during transport.
	// When the Content-MD5 header is specified, the File service compares the hash of the content that has arrived with the header value that was sent.
	// If the two hashes do not match, the operation will fail with error code 400 (Bad Request).
	ContentMD5 []byte

	// Initial data.
	Body io.ReadSeekCloser

	transactionalMD5 []byte

	FileRangeWrite *FileRangeWriteType

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *UploadFileRangeOptions) format() (rangeParam string, fileRangeWrite FileRangeWriteType, contentLength int64,
	fileUploadRangeOptions *FileUploadRangeOptions, leaseAccessConditions *LeaseAccessConditions, err error) {
	if o == nil || o.Body == nil {
		err = errors.New("invalid argument, body must not be nil")
		return
	}

	offset := int64(0)
	if o.Offset != nil {
		offset = *o.Offset
	}

	count := int64(CountToEnd)
	count, err = validateSeekableStreamAt0AndGetCount(o.Body)
	if err != nil {
		return
	}
	if count == 0 {
		err = errors.New("invalid argument, body must contain readable data whose size is > 0")
		return
	}

	rangeParamPtr := getRangeParam(to.Int64Ptr(offset), to.Int64Ptr(count))
	if rangeParamPtr != nil {
		rangeParam = *rangeParamPtr
	}
	fileRangeWrite = FileRangeWriteTypeUpdate
	if o.FileRangeWrite != nil {
		fileRangeWrite = *o.FileRangeWrite
	}
	contentLength = count

	fileUploadRangeOptions = &FileUploadRangeOptions{
		ContentMD5:   o.ContentMD5,
		Optionalbody: o.Body,
	}

	leaseAccessConditions = o.LeaseAccessConditions
	return
}

//----------------------------------------------------------------------------------------------------------------------

type UploadFileRangeFromURLOptions struct {

	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCRC64 []byte

	SourceModifiedAccessConditions *SourceModifiedAccessConditions

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *UploadFileRangeFromURLOptions) format(sourceURL url.URL, sourceOffset, destinationOffset, count int64) (rangeParam string, copySource string,
	contentLength int64, fileUploadRangeFromURLOptions *FileUploadRangeFromURLOptions, sourceModifiedAccessConditions *SourceModifiedAccessConditions,
	leaseAccessConditions *LeaseAccessConditions) {

	rangeParam = *getRangeParam(to.Int64Ptr(destinationOffset), to.Int64Ptr(count))
	copySource = sourceURL.String()
	contentLength = 0
	fileUploadRangeFromURLOptions = &FileUploadRangeFromURLOptions{
		SourceContentCRC64: o.SourceContentCRC64,
		SourceRange:        getRangeParam(to.Int64Ptr(sourceOffset), to.Int64Ptr(count)),
	}
	sourceModifiedAccessConditions = o.SourceModifiedAccessConditions
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

//----------------------------------------------------------------------------------------------------------------------

type ClearFileRangeOptions struct {
	Offset                *int64
	Count                 *int64
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ClearFileRangeOptions) format() (rangeParam string, fileRangeWrite FileRangeWriteType, contentLength int64,
	fileUploadRangeOptions *FileUploadRangeOptions, leaseAccessConditions *LeaseAccessConditions, err error) {

	if o == nil || o.Offset == nil || (*o.Offset <= 0) {
		err = errors.New("invalid argument: Offset is not specified (null) or Offset is <= 0")
		return
	}
	if o == nil || o.Count == nil || (*o.Count <= 0) {
		err = errors.New("invalid argument: either Count is not specified (null) or Count is <= 0")
	}

	rangeParam = *getRangeParam(o.Offset, o.Count)
	fileRangeWrite = FileRangeWriteTypeClear
	contentLength = *o.Count
	fileUploadRangeOptions = nil
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

//----------------------------------------------------------------------------------------------------------------------

type GetFileRangeListOptions struct {
	// The previous snapshot parameter is an opaque DateTime value that, when present, specifies the previous snapshot.
	PrevShareSnapshot *string
	Offset            *int64
	Count             *int64
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetFileRangeListOptions) format() (fileGetRangeListOptions *FileGetRangeListOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}

	return &FileGetRangeListOptions{
		Prevsharesnapshot: o.PrevShareSnapshot,
		Range:             getRangeParam(o.Offset, o.Count),
		Sharesnapshot:     o.ShareSnapshot,
	}, o.LeaseAccessConditions
}
