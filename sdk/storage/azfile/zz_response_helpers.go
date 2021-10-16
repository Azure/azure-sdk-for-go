package azfile

import (
	"context"
	"net/http"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------

// RetryableDownloadResponse wraps AutoRest generated DownloadResponse and helps to provide info for retry.
type RetryableDownloadResponse struct {
	dr FileDownloadResponse

	// Fields need for retry.
	ctx  context.Context
	f    FileClient
	info HTTPGetterInfo
}

// Response returns the raw HTTP response object.
func (dr RetryableDownloadResponse) Response() *http.Response {
	return dr.dr.RawResponse
}

// StatusCode returns the HTTP status code of the response, e.g. 200.
func (dr RetryableDownloadResponse) StatusCode() int {
	return dr.dr.RawResponse.StatusCode
}

// Status returns the HTTP status message of the response, e.g. "200 OK".
func (dr RetryableDownloadResponse) Status() string {
	return dr.dr.RawResponse.Status
}

// AcceptRanges returns the value for header Accept-Ranges.
func (dr RetryableDownloadResponse) AcceptRanges() *string {
	return dr.dr.AcceptRanges
}

// CacheControl returns the value for header Cache-Control.
func (dr RetryableDownloadResponse) CacheControl() *string {
	return dr.dr.CacheControl
}

// ContentDisposition returns the value for header Content-Disposition.
func (dr RetryableDownloadResponse) ContentDisposition() *string {
	return dr.dr.ContentDisposition
}

// ContentEncoding returns the value for header Content-Encoding.
func (dr RetryableDownloadResponse) ContentEncoding() *string {
	return dr.dr.ContentEncoding
}

// ContentLanguage returns the value for header Content-Language.
func (dr RetryableDownloadResponse) ContentLanguage() *string {
	return dr.dr.ContentLanguage
}

// ContentLength returns the value for header Content-Length.
func (dr RetryableDownloadResponse) ContentLength() *int64 {
	return dr.dr.ContentLength
}

// ContentRange returns the value for header Content-Range.
func (dr RetryableDownloadResponse) ContentRange() *string {
	return dr.dr.ContentRange
}

// ContentType returns the value for header Content-Type.
func (dr RetryableDownloadResponse) ContentType() *string {
	return dr.dr.ContentType
}

// CopyCompletionTime returns the value for header x-ms-copy-completion-time.
func (dr RetryableDownloadResponse) CopyCompletionTime() *time.Time {
	return dr.dr.CopyCompletionTime
}

// CopyID returns the value for header x-ms-copy-id.
func (dr RetryableDownloadResponse) CopyID() *string {
	return dr.dr.CopyID
}

// CopyProgress returns the value for header x-ms-copy-progress.
func (dr RetryableDownloadResponse) CopyProgress() *string {
	return dr.dr.CopyProgress
}

// CopySource returns the value for header x-ms-copy-source.
func (dr RetryableDownloadResponse) CopySource() *string {
	return dr.dr.CopySource
}

// CopyStatus returns the value for header x-ms-copy-status.
func (dr RetryableDownloadResponse) CopyStatus() *CopyStatusType {
	return dr.dr.CopyStatus
}

// CopyStatusDescription returns the value for header x-ms-copy-status-description.
func (dr RetryableDownloadResponse) CopyStatusDescription() *string {
	return dr.dr.CopyStatusDescription
}

// Date returns the value for header Date.
func (dr RetryableDownloadResponse) Date() *time.Time {
	return dr.dr.Date
}

// ETag returns the value for header ETag.
func (dr RetryableDownloadResponse) ETag() *string {
	return dr.dr.ETag
}

// IsServerEncrypted returns the value for header x-ms-server-encrypted.
func (dr RetryableDownloadResponse) IsServerEncrypted() *bool {
	return dr.dr.IsServerEncrypted
}

// LastModified returns the value for header Last-Modified.
func (dr RetryableDownloadResponse) LastModified() *time.Time {
	return dr.dr.LastModified
}

// RequestID returns the value for header x-ms-request-id.
func (dr RetryableDownloadResponse) RequestID() *string {
	return dr.dr.RequestID
}

// Version returns the value for header x-ms-version.
func (dr RetryableDownloadResponse) Version() *string {
	return dr.dr.Version
}

// NewMetadata returns user-defined key/value pairs.
func (dr RetryableDownloadResponse) NewMetadata() map[string]string {
	return dr.dr.Metadata
}

// FileContentMD5 returns the value for header x-ms-content-md5.
func (dr RetryableDownloadResponse) FileContentMD5() []byte {
	return dr.dr.FileContentMD5
}

// ContentMD5 returns the value for header Content-MD5.
func (dr RetryableDownloadResponse) ContentMD5() []byte {
	return dr.dr.ContentMD5
}
