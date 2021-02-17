package azblob

import (
	"context"
	"io"
	"net/http"
)

// NewHTTPHeaders returns the user-modifiable properties for this blob.
func (bgpr BlobGetPropertiesResponse) NewHTTPHeaders() BlobHttpHeaders {
	return BlobHttpHeaders{
		BlobContentType:        bgpr.ContentType,
		BlobContentEncoding:    bgpr.ContentEncoding,
		BlobContentLanguage:    bgpr.ContentLanguage,
		BlobContentDisposition: bgpr.ContentDisposition,
		BlobCacheControl:       bgpr.CacheControl,
		BlobContentMd5:         bgpr.ContentMD5,
	}
}

///////////////////////////////////////////////////////////////////////////////

// NewHTTPHeaders returns the user-modifiable properties for this blob.
func (dr BlobDownloadResponse) NewHTTPHeaders() BlobHttpHeaders {
	return BlobHttpHeaders{
		BlobContentType:        dr.ContentType,
		BlobContentEncoding:    dr.ContentEncoding,
		BlobContentLanguage:    dr.ContentLanguage,
		BlobContentDisposition: dr.ContentDisposition,
		BlobCacheControl:       dr.CacheControl,
		BlobContentMd5:         dr.ContentMD5,
	}
}

///////////////////////////////////////////////////////////////////////////////

// DownloadResponse wraps AutoRest generated DownloadResponse and helps to provide info for retry.
type DownloadResponse struct {
	BlobDownloadResponse
	ctx     context.Context
	b       BlobClient
	getInfo HTTPGetterInfo
}

// Body constructs new RetryReader stream for reading data. If a connection fails
// while reading, it will make additional requests to reestablish a connection and
// continue reading. Specifying a RetryReaderOption's with MaxRetryRequests set to 0
// (the default), returns the original response body and no retries will be performed.
func (r *DownloadResponse) Body(o RetryReaderOptions) io.ReadCloser {
	if o.MaxRetryRequests == 0 { // No additional retries
		return r.RawResponse.Body
	}
	return NewRetryReader(r.ctx, r.RawResponse, r.getInfo, o,
		func(ctx context.Context, getInfo HTTPGetterInfo) (*http.Response, error) {
			accessConditions := ModifiedAccessConditions{IfMatch: &getInfo.ETag}
			options := DownloadBlobOptions{
				Offset:                   &getInfo.Offset,
				Count:                    &getInfo.Count,
				ModifiedAccessConditions: &accessConditions,
				CpkInfo:                  o.CpkInfo,
				//CpkScopeInfo: 			  o.CpkScopeInfo,
			}
			resp, err := r.b.Download(ctx, &options)
			if err != nil {
				return nil, err
			}
			return resp.RawResponse, err
		},
	)
}
