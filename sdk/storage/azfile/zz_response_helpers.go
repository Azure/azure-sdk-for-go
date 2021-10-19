package azfile

import (
	"context"
	"io"
	"net/http"
)

//----------------------------------------------------------------------------------------------------------------------

// DownloadResponse wraps AutoRest generated DownloadResponse and helps to provide info for retry.
type DownloadResponse struct {
	FileDownloadResponse

	// Fields need for retry.
	ctx  context.Context
	f    FileClient
	info HTTPGetterInfo
}

// Body constructs a stream to read data from with a resilient reader option.
// A zero-value option means to get a raw stream.
func (dr *DownloadResponse) Body(o RetryReaderOptions) io.ReadCloser {
	if o.MaxRetryRequests == 0 {
		return dr.RawResponse.Body
	}

	return NewRetryReader(dr.ctx, dr.RawResponse, dr.info, o,
		func(ctx context.Context, info HTTPGetterInfo) (*http.Response, error) {
			resp, err := dr.f.Download(ctx, info.Offset, info.Count, nil)
			if err != nil {
				return nil, err
			}
			return resp.RawResponse, err
		})
}
