//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"io"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------

type FilePermissions struct {
	// If specified the permission (security descriptor) shall be set for the directory/file. This header can be used if Permission size is <= 8KB, else x-ms-file-permission-key
	// header shall be used. Default value: Inherit. If SDDL is specified as input, it must have owner, group and dacl. Note: Only one of the x-ms-file-permission
	// or x-ms-file-permission-key should be specified.
	PermissionStr *string
	// Key of the permission to be set for the directory/file. Note: Only one of the x-ms-file-permission or x-ms-file-permission-key should be specified.
	PermissionKey *string
}

func (p *FilePermissions) format(defaultFilePermissionStr *string) (filePermission *string, filePermissionKey *string, err error) {
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
// When you pass this into another function (Either literally or via ShareFileHTTPHeaders), the responseBody will probably fit inside SMBPropertyAdapter.
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

type FileCreateOptions struct {
	FileContentLength *int64

	Metadata map[string]string

	ShareFileHTTPHeaders *ShareFileHTTPHeaders

	FilePermissions *FilePermissions

	// In Windows, a 32 bit file attributes integer exists. This is that.
	SMBProperties *SMBProperties

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileCreateOptions) format() (fileContentLength int64, fileAttributes string, fileCreationTime string, fileLastWriteTime string, createOptions *fileClientCreateOptions, fileHTTPHeaders *ShareFileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions, err error) {
	if o == nil {
		return int64(0), DefaultFileAttributes, DefaultCurrentTimeString, DefaultCurrentTimeString,
			&fileClientCreateOptions{FilePermission: to.Ptr(DefaultFilePermissionString)}, nil, nil, nil
	}
	fileContentLength = 0

	if o.FileContentLength != nil {
		fileContentLength = *(o.FileContentLength)
	}

	fileAttributes, fileCreationTime, fileLastWriteTime = o.SMBProperties.format(false, DefaultFileAttributes, DefaultCurrentTimeString)

	filePermission, filePermissionKey, err := o.FilePermissions.format(&DefaultFilePermissionString)
	if err != nil {
		return
	}

	createOptions = &fileClientCreateOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
		Metadata:          o.Metadata,
	}

	fileHTTPHeaders = o.ShareFileHTTPHeaders
	leaseAccessConditions = o.LeaseAccessConditions

	return
}

type FileCreateResponse struct {
	fileClientCreateResponse
}

func toFileCreateResponse(resp fileClientCreateResponse) FileCreateResponse {
	return FileCreateResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileStartCopyOptions struct {
	FilePermissions *FilePermissions
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string

	CopyFileSmbInfo *CopyFileSmbInfo

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileStartCopyOptions) format() (startCopyOptions *fileClientStartCopyOptions, copyFileSmbInfo *CopyFileSmbInfo, leaseAccessConditions *LeaseAccessConditions, err error) {
	if o == nil {
		return
	}

	filePermission, filePermissionKey, err := o.FilePermissions.format(nil)
	if err != nil {
		return
	}

	startCopyOptions = &fileClientStartCopyOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
		Metadata:          o.Metadata,
	}

	copyFileSmbInfo = o.CopyFileSmbInfo
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

type FileStartCopyResponse struct {
	fileClientStartCopyResponse
}

func toFileStartCopyResponse(resp fileClientStartCopyResponse) FileStartCopyResponse {
	return FileStartCopyResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileAbortCopyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileAbortCopyOptions) format() (*fileClientAbortCopyOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}

type FileAbortCopyResponse struct {
	fileClientAbortCopyResponse
}

func toFileAbortCopyResponse(resp fileClientAbortCopyResponse) FileAbortCopyResponse {
	return FileAbortCopyResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileDownloadOptions struct {
	// When this header is set to true and specified together with the Range header,
	// the service returns the MD5 hash for the range, as long as the range is less than or equal to 4 MB in size.
	RangeGetContentMD5 *bool

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileDownloadOptions) format(offset, count int64) (downloadOptions *fileClientDownloadOptions, leaseAccessConditions *LeaseAccessConditions) {
	downloadOptions = &fileClientDownloadOptions{Range: NewHttpRange(offset, count).format()}

	if o != nil {
		downloadOptions.RangeGetContentMD5 = o.RangeGetContentMD5
		leaseAccessConditions = o.LeaseAccessConditions
	}

	return
}

// FileDownloadResponse wraps AutoRest generated DownloadResponse and helps to provide info for retry.
type FileDownloadResponse struct {
	fileClientDownloadResponse
	// Fields need for retry.
	ctx  context.Context
	f    *FileClient
	info HTTPGetterInfo
}

func toFileDownloadResponse(ctx context.Context, f *FileClient, downloadResponse fileClientDownloadResponse, offset int64, count int64) FileDownloadResponse {
	return FileDownloadResponse{
		fileClientDownloadResponse: downloadResponse,
		f:                          f,
		ctx:                        ctx,
		info:                       HTTPGetterInfo{Offset: offset, Count: count, ETag: downloadResponse.ETag},
	}
}

// FileBody constructs a stream to read data from with a resilient reader option.
// A zero-value option means to get a raw stream.
func (dr *FileDownloadResponse) FileBody(o RetryReaderOptions) io.ReadCloser {
	if o.MaxRetryRequests == 0 {
		return dr.Body
	}

	return NewRetryReader(dr.ctx, dr.Body, dr.info, o,
		func(ctx context.Context, info HTTPGetterInfo) (io.ReadCloser, error) {
			resp, err := dr.f.Download(ctx, info.Offset, info.Count, nil)
			if err != nil {
				return nil, err
			}
			return resp.Body, err
		})
}

//----------------------------------------------------------------------------------------------------------------------

type FileDeleteOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileDeleteOptions) format() (deleteOptions *fileClientDeleteOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

type FileDeleteResponse struct {
	fileClientDeleteResponse
}

func toFileDeleteResponse(resp fileClientDeleteResponse) FileDeleteResponse {
	return FileDeleteResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileGetPropertiesOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileGetPropertiesOptions) format() (getPropertiesOptions *fileClientGetPropertiesOptions, leaseAccessConditions *LeaseAccessConditions) {
	if o == nil {
		return
	}
	getPropertiesOptions = &fileClientGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}
	leaseAccessConditions = o.LeaseAccessConditions
	return
}

type FileGetPropertiesResponse struct {
	fileClientGetPropertiesResponse
}

func toFileGetPropertiesResponse(resp fileClientGetPropertiesResponse) FileGetPropertiesResponse {
	return FileGetPropertiesResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type SetFileHTTPHeadersOptions struct {
	// Resizes a file to the specified size.
	// If the specified byte value is less than the current size of the file,
	// then all ranges above the specified byte value are cleared.
	FileContentLength *int64

	FilePermissions *FilePermissions

	SMBProperties *SMBProperties

	ShareFileHTTPHeaders *ShareFileHTTPHeaders

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetFileHTTPHeadersOptions) format() (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	setHTTPHeadersOptions *fileClientSetHTTPHeadersOptions, fileHTTPHeaders *ShareFileHTTPHeaders,
	leaseAccessConditions *LeaseAccessConditions, err error) {

	fileAttributes, fileCreationTime, fileLastWriteTime = DefaultPreserveString, DefaultPreserveString, DefaultPreserveString

	filePermission, filePermissionKey, err := o.FilePermissions.format(&DefaultPreserveString)
	if err != nil {
		return
	}

	setHTTPHeadersOptions = &fileClientSetHTTPHeadersOptions{
		FilePermission:    filePermission,
		FilePermissionKey: filePermissionKey,
	}

	if o != nil {
		setHTTPHeadersOptions.FileContentLength = o.FileContentLength
		fileHTTPHeaders = o.ShareFileHTTPHeaders
		leaseAccessConditions = o.LeaseAccessConditions
	}
	return
}

type FileSetHTTPHeadersResponse struct {
	fileClientSetHTTPHeadersResponse
}

func toFileSetHTTPHeadersResponse(resp fileClientSetHTTPHeadersResponse) FileSetHTTPHeadersResponse {
	return FileSetHTTPHeadersResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileSetMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileSetMetadataOptions) format(metadata map[string]string) (setMetadataOptions *fileClientSetMetadataOptions, leaseAccessConditions *LeaseAccessConditions) {
	setMetadataOptions = &fileClientSetMetadataOptions{
		Metadata: metadata,
	}

	if o != nil {
		leaseAccessConditions = o.LeaseAccessConditions
	}
	return
}

type FileSetMetadataResponse struct {
	fileClientSetMetadataResponse
}

func toFileSetMetadataResponse(resp fileClientSetMetadataResponse) FileSetMetadataResponse {
	return FileSetMetadataResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileResizeOptions struct {
	LeaseAccessCondition *LeaseAccessConditions
}

func (o *FileResizeOptions) format(contentLength int64) (fileAttributes string, fileCreationTime string, fileLastWriteTime string,
	setHTTPHeadersOptions *fileClientSetHTTPHeadersOptions, fileHTTPHeaders *ShareFileHTTPHeaders, leaseAccessConditions *LeaseAccessConditions) {

	fileAttributes, fileCreationTime, fileLastWriteTime = DefaultPreserveString, DefaultPreserveString, DefaultPreserveString

	setHTTPHeadersOptions = &fileClientSetHTTPHeadersOptions{
		FileContentLength: to.Ptr(contentLength),
		FilePermission:    &DefaultPreserveString,
	}

	fileHTTPHeaders = nil

	if o != nil {
		leaseAccessConditions = o.LeaseAccessCondition
	}
	return
}

type FileResizeResponse struct {
	fileClientSetHTTPHeadersResponse
}

func toFileResizeResponse(resp fileClientSetHTTPHeadersResponse) FileResizeResponse {
	return FileResizeResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileUploadRangeOptions struct {
	// An MD5 hash of the content. This hash is used to verify the integrity of the data during transport. When the Content-MD5
	// header is specified, the File service compares the hash of the content that has
	// arrived with the header value that was sent. If the two hashes do not match, the operation will fail with error code 400
	// (Bad Request).
	ContentMD5 []byte

	TransactionalMD5 []byte

	FileRangeWrite *FileRangeWriteType

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileUploadRangeOptions) format(offset int64, body io.ReadSeekCloser) (string, FileRangeWriteType, int64,
	*fileClientUploadRangeOptions, *LeaseAccessConditions, error) {

	if offset < 0 || body == nil {
		return "", "", 0, nil, nil, errors.New("invalid argument: offset must be >= 0 and body must not be nil")
	}

	count := int64(CountToEnd)
	count, err := validateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return "", "", 0, nil, nil, err
	}

	if count == 0 {
		return "", "", 0, nil, nil, errors.New("invalid argument: body must contain readable data whose size is > 0")
	}

	httpRange := getSourceRange(to.Ptr(offset), to.Ptr(count))
	fileRangeWrite := FileRangeWriteTypeUpdate
	var contentMD5 []byte
	var leaseAccessConditions *LeaseAccessConditions

	if o != nil {
		if o.FileRangeWrite != nil {
			fileRangeWrite = *o.FileRangeWrite
		}
		if o.ContentMD5 != nil {
			contentMD5 = o.ContentMD5
		}
		leaseAccessConditions = o.LeaseAccessConditions
	}

	uploadRangeOptions := &fileClientUploadRangeOptions{ContentMD5: contentMD5, Optionalbody: body}
	return httpRange, fileRangeWrite, count, uploadRangeOptions, leaseAccessConditions, nil
}

type FileUploadRangeResponse struct {
	fileClientUploadRangeResponse
}

func toFileUploadRangeResponse(resp fileClientUploadRangeResponse) FileUploadRangeResponse {
	return FileUploadRangeResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileUploadRangeFromURLOptions struct {
	// Specify the crc64 calculated for the range of bytes that must be read from the copy source.
	SourceContentCRC64 []byte

	SourceModifiedAccessConditions *SourceModifiedAccessConditions

	LeaseAccessConditions *LeaseAccessConditions

	// Only Bearer type is supported. Credentials should be a valid OAuth access token to copy source.
	CopySourceAuthorization *string
	// Bytes of source data in the specified range.
	SourceRange *string
}

func (o *FileUploadRangeFromURLOptions) format(sourceURL string, sourceOffset, destinationOffset, count int64) (string,
	string, int64, *fileClientUploadRangeFromURLOptions, *SourceModifiedAccessConditions, *LeaseAccessConditions) {

	destinationDataRange := getSourceRange(to.Ptr(destinationOffset), to.Ptr(count))
	uploadRangeFromURLOptions := &fileClientUploadRangeFromURLOptions{
		SourceRange: NewHttpRange(sourceOffset, count).format(),
	}

	if o != nil {
		uploadRangeFromURLOptions.SourceContentCRC64 = o.SourceContentCRC64
		uploadRangeFromURLOptions.CopySourceAuthorization = o.CopySourceAuthorization
		return destinationDataRange, sourceURL, count, uploadRangeFromURLOptions, o.SourceModifiedAccessConditions, o.LeaseAccessConditions
	}

	return destinationDataRange, sourceURL, count, uploadRangeFromURLOptions, nil, nil
}

type FileUploadRangeFromURLResponse struct {
	fileClientUploadRangeFromURLResponse
}

func toFileUploadRangeFromURLResponse(resp fileClientUploadRangeFromURLResponse) FileUploadRangeFromURLResponse {
	return FileUploadRangeFromURLResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileClearRangeOptions struct {
	// An MD5 hash of the content. This hash is used to verify the integrity of the data during transport. When the Content-MD5
	// header is specified, the File service compares the hash of the content that has
	// arrived with the header value that was sent. If the two hashes do not match, the operation will fail with error code 400
	// (Bad Request).
	ContentMD5 []byte
	// Initial data.
	OptionalBody io.ReadSeekCloser

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileClearRangeOptions) format(offset, count int64) (string, FileRangeWriteType, int64, *fileClientUploadRangeOptions, *LeaseAccessConditions, error) {
	httpRange := getSourceRange(to.Ptr(offset), to.Ptr(count))
	fileRangeWrite := FileRangeWriteTypeClear
	contentLength := int64(0)

	if offset < 0 || count <= 0 {
		return httpRange, fileRangeWrite, contentLength, &fileClientUploadRangeOptions{}, nil, errors.New("invalid argument: either offset is < 0 or count <= 0")
	}

	if o == nil {
		return httpRange, fileRangeWrite, contentLength, &fileClientUploadRangeOptions{}, nil, nil
	}

	uploadRangeOptions := fileClientUploadRangeOptions{
		ContentMD5:   o.ContentMD5,
		Optionalbody: o.OptionalBody,
	}
	return httpRange, fileRangeWrite, contentLength, &uploadRangeOptions, o.LeaseAccessConditions, nil
}

type FileClearRangeResponse struct {
	fileClientUploadRangeResponse
}

func toFileClearRangeResponse(resp fileClientUploadRangeResponse) FileClearRangeResponse {
	return FileClearRangeResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------

type FileGetRangeListOptions struct {
	// The previous snapshot parameter is an opaque DateTime value that, when present, specifies the previous snapshot.
	PrevShareSnapshot *string
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *FileGetRangeListOptions) format(offset, count int64) (*fileClientGetRangeListOptions, *LeaseAccessConditions) {
	getRangeListOptions := &fileClientGetRangeListOptions{
		Range: NewHttpRange(offset, count).format(),
	}

	if o != nil {
		getRangeListOptions.Prevsharesnapshot = o.PrevShareSnapshot
		getRangeListOptions.Sharesnapshot = o.ShareSnapshot
		return getRangeListOptions, o.LeaseAccessConditions
	}

	return getRangeListOptions, nil
}

type FileGetRangeListResponse struct {
	fileClientGetRangeListResponse
}

func toFileGetRangeListResponse(resp fileClientGetRangeListResponse) FileGetRangeListResponse {
	return FileGetRangeListResponse{resp}
}

//----------------------------------------------------------------------------------------------------------------------
