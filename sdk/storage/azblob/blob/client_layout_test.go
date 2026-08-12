// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/stretchr/testify/require"
)

// This file covers layout-aware routing at the Client level: what requests DownloadBuffer actually
// puts on the wire. The pure layout parsing/selection logic lives in layout_test.go.

// ======================================================================================== //
// Helper methods for layout mock tests

type fakeLayoutResponder struct {
	l                     layout
	layoutResponses       map[string]*http.Response
	getPropertiesResponse *http.Response

	// mu guards the fields below. Do is invoked concurrently by the chunk download goroutines,
	// so every access must be synchronized.
	mu                  sync.Mutex
	layoutCalls         int
	getPropertiesCalled bool
	localityGets        int
	normalGets          int
	// downloadIfMatch records the If-Match header of every chunk download request.
	downloadIfMatch []string
	// layoutIfMatch records the If-Match header of every get layout request.
	layoutIfMatch []string
	// layoutStatusOverride, when non-nil, is consulted on each layout call to override the
	// canned response. It receives the 1-based call number.
	layoutStatusOverride func(call int) *http.Response
}

func (f *fakeLayoutResponder) reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.layoutCalls = 0
	f.getPropertiesCalled = false
	f.localityGets = 0
	f.normalGets = 0
	f.downloadIfMatch = nil
	f.layoutIfMatch = nil
}

// counts returns a consistent snapshot of the request counters.
func (f *fakeLayoutResponder) counts() (layoutCalls, localityGets, normalGets int, getPropertiesCalled bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.layoutCalls, f.localityGets, f.normalGets, f.getPropertiesCalled
}

// ifMatchHeaders returns a copy of the If-Match header seen on each chunk download.
func (f *fakeLayoutResponder) ifMatchHeaders() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.downloadIfMatch...)
}

// layoutIfMatchHeaders returns a copy of the If-Match header seen on each get layout request.
func (f *fakeLayoutResponder) layoutIfMatchHeaders() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.layoutIfMatch...)
}

func newFakeLayoutResponder(l layout, getPropsResponse *http.Response) *fakeLayoutResponder {
	layoutResponses := make(map[string]*http.Response)
	pages := splitLayoutToPages(l, 3) // Use a small page size to create multiple pages for testing
	for i, page := range pages {
		if i == 0 {
			layoutResponses[""] = newMockLayoutResponse(l.contentLength, string(*l.eTag), page, 0)
		} else {
			layoutResponses[strconv.Itoa(i)] = newMockLayoutResponse(l.contentLength, string(*l.eTag), page, 0)
		}
	}
	return &fakeLayoutResponder{
		l:                     l,
		layoutResponses:       layoutResponses,
		getPropertiesResponse: getPropsResponse,
	}
}

func (f *fakeLayoutResponder) Do(req *http.Request) (*http.Response, error) {
	// Layout
	qp := req.URL.Query()
	if comp := qp.Get("comp"); comp == "layout" {
		f.mu.Lock()
		f.layoutCalls++
		call := f.layoutCalls
		override := f.layoutStatusOverride
		f.layoutIfMatch = append(f.layoutIfMatch, req.Header.Get("If-Match"))
		f.mu.Unlock()
		if override != nil {
			if resp := override(call); resp != nil {
				return resp, nil
			}
		}
		marker := qp.Get("marker")
		return f.layoutResponses[marker], nil
	}

	// Get properties
	if req.Method == http.MethodHead {
		f.mu.Lock()
		f.getPropertiesCalled = true
		f.mu.Unlock()
		return f.getPropertiesResponse, nil
	}

	// Validate that the request range is going to the right layout
	if req.Method == http.MethodGet {
		f.mu.Lock()
		// If the request Host is different from the URL host
		if req.Host != req.URL.Host {
			f.localityGets++
		} else {
			f.normalGets++
		}
		f.downloadIfMatch = append(f.downloadIfMatch, req.Header.Get("If-Match"))
		f.mu.Unlock()
		// Download
		header := http.Header{}
		header.Set("Content-Length", strconv.FormatInt(rangeLength(req.Header.Get("x-ms-range")), 10))
		return &http.Response{
			StatusCode: http.StatusPartialContent,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
			Header:     header,
		}, nil
	}
	return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
}

// rangeLength returns the number of bytes covered by an "bytes=start-end" range header.
func rangeLength(r string) int64 {
	var start, end int64
	if _, err := fmt.Sscanf(r, "bytes=%d-%d", &start, &end); err != nil {
		return 0
	}
	return end - start + 1
}

func newMockLayoutResponse(contentLength int64, eTag string, layout generated.BlobLayout, statusCode int) *http.Response {
	if statusCode == 0 || statusCode == http.StatusOK {
		data, _ := xml.Marshal(layout)
		// NOTE: use Header.Set, not a map literal. "ETag" is not the canonical MIME key ("Etag"),
		// so a literal would be invisible to Header.Get and the ETag would never be parsed.
		header := http.Header{}
		header.Set("x-ms-blob-content-length", fmt.Sprintf("%d", contentLength))
		header.Set("ETag", eTag)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(data)),
			Header:     header,
		}
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}
}

func newMockGetPropertiesResponse(contentLength int64, eTag string) *http.Response {
	header := http.Header{}
	header.Set("Content-Length", fmt.Sprintf("%d", contentLength))
	header.Set("ETag", eTag)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
	}
}

// splitLayoutToPages splits a layout into multiple BlobLayout pages with sequential ranges.
// Each page contains up to maxRangesPerPage ranges.
func splitLayoutToPages(l layout, maxRangesPerPage int) []generated.BlobLayout {
	if len(l.layoutRanges) == 0 {
		return nil
	}

	if maxRangesPerPage <= 0 {
		maxRangesPerPage = 1
	}

	// Build unique endpoints map
	endpointMap := make(map[string]int32)
	var endpointIndex int32
	for _, lr := range l.layoutRanges {
		if _, exists := endpointMap[lr.endpoint]; !exists {
			endpointMap[lr.endpoint] = endpointIndex
			endpointIndex++
		}
	}

	// Convert map to Endpoint slice
	endpoints := make([]*generated.BlobLayoutEndpointsEndpointItem, len(endpointMap))
	for ep, idx := range endpointMap {
		epCopy := ep
		idxCopy := idx
		endpoints[idx] = &generated.BlobLayoutEndpointsEndpointItem{
			Index: &idxCopy,
			Value: &epCopy,
		}
	}

	var pages []generated.BlobLayout
	for i := 0; i < len(l.layoutRanges); i += maxRangesPerPage {
		end := i + maxRangesPerPage
		if end > len(l.layoutRanges) {
			end = len(l.layoutRanges)
		}

		ranges := make([]*generated.BlobLayoutRangesRangeItem, 0, end-i)
		for j := i; j < end; j++ {
			lr := l.layoutRanges[j]
			start := lr.start
			rangeEnd := lr.end
			epIdx := endpointMap[lr.endpoint]
			ranges = append(ranges, &generated.BlobLayoutRangesRangeItem{
				Start:         &start,
				End:           &rangeEnd,
				EndpointIndex: &epIdx,
			})
		}

		// just pass in the index as a marker for testing purposes, the actual value doesn't matter
		iCopy := strconv.Itoa(len(pages) + 1)
		page := generated.BlobLayout{
			Endpoints: &generated.BlobLayoutEndpoints{
				Endpoint: endpoints,
			},
			Ranges: &generated.BlobLayoutRanges{
				Range: ranges,
			},
			NextMarker: &iCopy,
		}
		pages = append(pages, page)
	}

	// Last page should have empty NextMarker to indicate no more pages
	if len(pages) > 0 {
		pages[len(pages)-1].NextMarker = nil
	}

	return pages
}

// newFakeLayoutClient builds a client wired to the given fake transport.
func newFakeLayoutClient(t *testing.T, f *fakeLayoutResponder) *Client {
	t.Helper()
	client, err := NewClientWithNoCredential("https://fake.blob.core.windows.net/container/blob", &ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: f},
	})
	require.NoError(t, err)
	return client
}

// buildLayout produces a layout of n equally sized ranges spread over endpointCount endpoints.
func buildLayout(n int, rangeSize int64, endpointCount int, etag *azcore.ETag) layout {
	ranges := make([]layoutRange, 0, n)
	for i := int64(0); i < int64(n); i++ {
		ranges = append(ranges, layoutRange{
			start:    i * rangeSize,
			end:      i*rangeSize + rangeSize - 1,
			endpoint: fmt.Sprintf("https://locality%d.blob.core.windows.net", i%int64(endpointCount)),
		})
	}
	return layout{layoutRanges: ranges, contentLength: int64(n) * rangeSize, eTag: etag}
}

// ======================================================================================== //
// Tests

func TestDownloadBufferWithLayoutAwareRoutingError(t *testing.T) {
	f := &fakeLayoutResponder{}
	client, err := NewClientWithNoCredential("https://fake/blob/path", &ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: f,
		},
	})
	require.NoError(t, err)

	buff := make([]byte, 0)
	// 412 should trigger an error
	f.layoutResponses = map[string]*http.Response{"": newMockLayoutResponse(0, "etag", generated.BlobLayout{}, http.StatusPreconditionFailed)}
	_, err = client.DownloadBuffer(context.Background(), buff, &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "412")

	// 400 should trigger a fallback to get  properties
	f.reset()
	f.layoutResponses = map[string]*http.Response{"": newMockLayoutResponse(0, "etag", generated.BlobLayout{}, http.StatusBadRequest)}
	f.getPropertiesResponse = newMockGetPropertiesResponse(0, "etag")
	_, err = client.DownloadBuffer(context.Background(), buff, &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
	})
	require.NoError(t, err)
	layoutCalls, localityGets, _, getPropsCalled := f.counts()
	require.Equal(t, 1, layoutCalls)
	require.True(t, getPropsCalled)
	require.Zero(t, localityGets)
}

func TestDownloadBufferWithLayoutAwareRoutingNoLayout(t *testing.T) {
	f := &fakeLayoutResponder{}
	client, err := NewClientWithNoCredential("https://fake/blob/path", &ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: f,
		},
	})
	require.NoError(t, err)

	buff := make([]byte, 0)
	f.layoutResponses = map[string]*http.Response{"": newMockLayoutResponse(10, "etag", generated.BlobLayout{}, http.StatusOK)}
	_, err = client.DownloadBuffer(context.Background(), buff, &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
	})
	require.NoError(t, err)
	layoutCalls, localityGets, normalGets, getPropsCalled := f.counts()
	require.Equal(t, 1, layoutCalls)
	require.False(t, getPropsCalled)
	require.Equal(t, 1, normalGets)
	require.Zero(t, localityGets)
}

func TestDownloadBufferWithLayoutAwareRoutingWithLayout(t *testing.T) {
	etag := azcore.ETag("test-etag")
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 99, endpoint: "https://locality1.blob.core.windows.net"},
			{start: 100, end: 199, endpoint: "https://locality2.blob.core.windows.net"},
			{start: 200, end: 299, endpoint: "https://locality1.blob.core.windows.net"},
		},
		contentLength: 300,
		eTag:          &etag,
	}

	f := newFakeLayoutResponder(l, nil)
	client := newFakeLayoutClient(t, f)

	buff := make([]byte, 300)
	_, err := client.DownloadBuffer(context.Background(), buff, &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
	})
	require.NoError(t, err)
	layoutCalls, localityGets, normalGets, getPropsCalled := f.counts()
	require.GreaterOrEqual(t, layoutCalls, 1)
	require.False(t, getPropsCalled)
	require.Greater(t, localityGets, 0)
	require.Zero(t, normalGets)
}

func TestDownloadBufferWithLayoutAwareRoutingMultiplePages(t *testing.T) {
	etag := azcore.ETag("multi-page-etag")
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 99, endpoint: "https://locality1.blob.core.windows.net"},
			{start: 100, end: 199, endpoint: "https://locality2.blob.core.windows.net"},
			{start: 200, end: 299, endpoint: "https://locality3.blob.core.windows.net"},
			{start: 300, end: 399, endpoint: "https://locality1.blob.core.windows.net"},
			{start: 400, end: 499, endpoint: "https://locality2.blob.core.windows.net"},
		},
		contentLength: 500,
		eTag:          &etag,
	}

	f := newFakeLayoutResponder(l, nil)
	client := newFakeLayoutClient(t, f)

	buff := make([]byte, 500)
	_, err := client.DownloadBuffer(context.Background(), buff, &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
	})
	require.NoError(t, err)
	// With maxRangesPerPage=3 in splitLayoutToPages, 5 ranges should create 2 pages
	layoutCalls, localityGets, normalGets, getPropsCalled := f.counts()
	require.GreaterOrEqual(t, layoutCalls, 2)
	require.False(t, getPropsCalled)
	require.Greater(t, localityGets, 0)
	require.Zero(t, normalGets)
}

// TestDownloadBufferLayoutFetchedOncePerDownload is the regression test for layout caching: the
// layout must be fetched once for the whole download, not once per chunk. Before the fallback/
// caching fix, each of the N chunk goroutines re-issued the GetLayout request.
func TestDownloadBufferLayoutFetchedOncePerDownload(t *testing.T) {
	etag := azcore.ETag("etag")
	const numRanges, rangeSize = 20, int64(100)
	l := buildLayout(numRanges, rangeSize, 3, &etag)

	f := newFakeLayoutResponder(l, nil)
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, l.contentLength), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
		BlockSize:          rangeSize, // one chunk per layout range
		Concurrency:        8,
	})
	require.NoError(t, err)

	layoutCalls, localityGets, normalGets, getPropsCalled := f.counts()
	require.Equal(t, numRanges, localityGets, "every chunk should be routed to a locality endpoint")
	require.Zero(t, normalGets)
	require.False(t, getPropsCalled)
	// splitLayoutToPages uses maxRangesPerPage=3, so 20 ranges => 7 pages, each fetched exactly once.
	require.Equal(t, 7, layoutCalls, "layout must be fetched once per download, not once per chunk")
}

// TestDownloadBufferFallbackCachedAcrossChunks verifies the "layout unavailable" decision is
// fetched once and reused by every chunk instead of re-requested per chunk.
func TestDownloadBufferFallbackCachedAcrossChunks(t *testing.T) {
	f := &fakeLayoutResponder{
		layoutResponses:       map[string]*http.Response{"": newMockLayoutResponse(0, "etag", generated.BlobLayout{}, http.StatusBadRequest)},
		getPropertiesResponse: newMockGetPropertiesResponse(2000, "etag"),
	}
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, 2000), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
		BlockSize:          100, // 20 chunks
		Concurrency:        8,
	})
	require.NoError(t, err)

	layoutCalls, localityGets, normalGets, getPropsCalled := f.counts()
	require.Equal(t, 1, layoutCalls, "the fallback decision must be cached, not re-requested per chunk")
	require.True(t, getPropsCalled)
	require.Equal(t, 20, normalGets)
	require.Zero(t, localityGets)
}

// TestDownloadBufferLayoutDisabled verifies the legacy path is untouched: no GetLayout request is
// issued at all when layout-aware routing is explicitly disabled, and no If-Match is added.
func TestDownloadBufferLayoutDisabled(t *testing.T) {
	f := &fakeLayoutResponder{getPropertiesResponse: newMockGetPropertiesResponse(300, "etag")}
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, 300), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingDisabled,
		BlockSize:          100,
	})
	require.NoError(t, err)

	layoutCalls, localityGets, normalGets, getPropsCalled := f.counts()
	require.Zero(t, layoutCalls, "no layout request should be made when routing is disabled")
	require.True(t, getPropsCalled)
	require.Equal(t, 3, normalGets)
	require.Zero(t, localityGets)
	for _, v := range f.ifMatchHeaders() {
		require.Empty(t, v, "the legacy path must not add an If-Match of its own")
	}
}

// TestDownloadBufferFallbackNoETagLock verifies that when the service can't supply a layout and the
// download falls back to GetProperties, no If-Match is added either.
func TestDownloadBufferFallbackNoETagLock(t *testing.T) {
	f := &fakeLayoutResponder{
		layoutResponses:       map[string]*http.Response{"": newMockLayoutResponse(0, "etag", generated.BlobLayout{}, http.StatusBadRequest)},
		getPropertiesResponse: newMockGetPropertiesResponse(300, "etag"),
	}
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, 300), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
		BlockSize:          100,
	})
	require.NoError(t, err)

	_, _, _, getPropsCalled := f.counts()
	require.True(t, getPropsCalled)
	for _, v := range f.ifMatchHeaders() {
		require.Empty(t, v, "the GetProperties fallback must not add an If-Match")
	}
}

// TestDownloadBufferLayoutDefaultEnabled verifies that leaving LayoutAwareRouting unset (the zero
// value, LayoutAwareRoutingAuto) resolves to enabled: the layout is fetched and used for routing.
func TestDownloadBufferLayoutDefaultEnabled(t *testing.T) {
	etag := azcore.ETag("etag")
	l := buildLayout(3, 100, 2, &etag)

	f := newFakeLayoutResponder(l, nil)
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, l.contentLength), &DownloadBufferOptions{
		BlockSize: 100,
	})
	require.NoError(t, err)

	layoutCalls, localityGets, normalGets, getPropsCalled := f.counts()
	require.NotZero(t, layoutCalls, "the default must fetch the layout")
	require.Equal(t, 3, localityGets, "every chunk should be routed to a locality endpoint by default")
	require.Zero(t, normalGets)
	require.False(t, getPropsCalled, "the layout response supplies the length, so GetProperties is unnecessary")
}

// TestDownloadBufferLayoutETagLock verifies the ETag returned by GetLayout is used to lock every
// chunk download, which is what guarantees a consistent blob across the parallel reads.
func TestDownloadBufferLayoutETagLock(t *testing.T) {
	etag := azcore.ETag("layout-etag")
	l := buildLayout(6, 100, 2, &etag)

	f := newFakeLayoutResponder(l, nil)
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, l.contentLength), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
		BlockSize:          100,
		Concurrency:        4,
	})
	require.NoError(t, err)

	headers := f.ifMatchHeaders()
	require.Len(t, headers, 6)
	for _, v := range headers {
		require.Equal(t, string(etag), v, "every chunk must be ETag-locked to the layout response")
	}
}

// TestDownloadBufferUserETagWins verifies a caller-supplied If-Match is not overwritten by the
// ETag from the layout response.
func TestDownloadBufferUserETagWins(t *testing.T) {
	layoutETag := azcore.ETag("layout-etag")
	userETag := azcore.ETag("user-etag")
	l := buildLayout(3, 100, 2, &layoutETag)

	f := newFakeLayoutResponder(l, nil)
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, l.contentLength), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
		BlockSize:          100,
		AccessConditions: &AccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &userETag},
		},
	})
	require.NoError(t, err)

	headers := f.ifMatchHeaders()
	require.NotEmpty(t, headers)
	for _, v := range headers {
		require.Equal(t, string(userETag), v, "a caller-supplied If-Match must take precedence")
	}
}

// TestDownloadBufferRefreshFailureIsNonFatal verifies that once a layout has been cached, a failed
// eager refresh doesn't fail the download: temporal.Resource serves the stale-but-valid layout.
func TestDownloadBufferRefreshFailureIsNonFatal(t *testing.T) {
	// Make the cached layout go stale immediately so the chunk goroutines attempt a refresh.
	defer func(refresh time.Duration) { layoutRefresh = refresh }(layoutRefresh)
	layoutRefresh = time.Millisecond

	etag := azcore.ETag("etag")
	l := buildLayout(3, 100, 2, &etag)
	f := newFakeLayoutResponder(l, nil)

	// splitLayoutToPages(_, 3) yields a single page, so the initial fetch is call #1.
	// Fail every call after that to simulate the refresh failing.
	f.layoutStatusOverride = func(call int) *http.Response {
		if call == 1 {
			return nil // serve the canned successful page
		}
		return newMockLayoutResponse(0, "", generated.BlobLayout{}, http.StatusInternalServerError)
	}
	client := newFakeLayoutClient(t, f)

	time.Sleep(5 * time.Millisecond) // ensure the first fetch is already stale when chunks run

	_, err := client.DownloadBuffer(context.Background(), make([]byte, l.contentLength), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
		BlockSize:          100,
		Concurrency:        1,
	})
	require.NoError(t, err, "a failed layout refresh must not fail the download")

	_, localityGets, normalGets, _ := f.counts()
	require.Equal(t, 3, localityGets+normalGets, "all chunks should still be downloaded")
}

// TestClientLayoutFallbackCachedSingleRequest verifies that, when driven through the client's
// GetLayoutPager, a cached fallback decision is only fetched from the service once.
func TestClientLayoutFallbackCachedSingleRequest(t *testing.T) {
	f := &fakeLayoutResponder{
		layoutResponses:       map[string]*http.Response{"": newMockLayoutResponse(0, "etag", generated.BlobLayout{}, http.StatusBadRequest)},
		getPropertiesResponse: newMockGetPropertiesResponse(1024, "etag"),
	}
	client := newFakeLayoutClient(t, f)

	temporalLayout := temporal.NewResourceWithOptions(
		func(ctx context.Context) (layout, time.Time, error) {
			return getLayout(ctx, client.GetLayoutPager(nil))
		}, temporal.ResourceOptions[layout, context.Context]{
			ShouldRefresh: shouldRefreshLayout,
		})

	for i := 0; i < 3; i++ {
		l, err := temporalLayout.Get(context.Background())
		require.NoError(t, err)
		require.True(t, l.fallback)
	}
	layoutCalls, _, _, _ := f.counts()
	require.Equal(t, 1, layoutCalls, "the failure should be cached, not re-requested")
}

// TestDownloadBufferCountSkipsGetProperties verifies that no GetProperties request is issued when
// the caller specifies a count: the blob's length is already known, so the extra round trip would
// be wasted.
func TestDownloadBufferCountSkipsGetProperties(t *testing.T) {
	f := &fakeLayoutResponder{getPropertiesResponse: newMockGetPropertiesResponse(300, "etag")}
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, 200), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingDisabled,
		Range:              HTTPRange{Offset: 0, Count: 200},
		BlockSize:          100,
	})
	require.NoError(t, err)

	layoutCalls, localityGets, normalGets, getPropsCalled := f.counts()
	require.Zero(t, layoutCalls, "no layout request should be made when routing is disabled")
	require.False(t, getPropsCalled, "the caller supplied a count, so GetProperties is unnecessary")
	require.Equal(t, 2, normalGets)
	require.Zero(t, localityGets)

	// Without an initial call there's no ETag to lock on, so no If-Match is sent.
	for _, v := range f.ifMatchHeaders() {
		require.Empty(t, v)
	}
}

// TestDownloadBufferCountWithLayoutStillETagLocks verifies that specifying a count doesn't disable
// the ETag lock when layout aware routing is on: the layout response supplies the ETag without an
// extra GetProperties call.
func TestDownloadBufferCountWithLayoutStillETagLocks(t *testing.T) {
	etag := azcore.ETag("layout-etag")
	l := buildLayout(3, 100, 2, &etag)

	f := newFakeLayoutResponder(l, nil)
	client := newFakeLayoutClient(t, f)

	_, err := client.DownloadBuffer(context.Background(), make([]byte, 200), &DownloadBufferOptions{
		LayoutAwareRouting: LayoutAwareRoutingEnabled,
		Range:              HTTPRange{Offset: 0, Count: 200},
		BlockSize:          100,
	})
	require.NoError(t, err)

	_, _, _, getPropsCalled := f.counts()
	require.False(t, getPropsCalled)

	headers := f.ifMatchHeaders()
	require.Len(t, headers, 2)
	for _, v := range headers {
		require.Equal(t, string(etag), v)
	}
}

// TestGetLayoutPagerLocksETagAcrossPages verifies that, when the caller doesn't supply an If-Match
// condition, the ETag returned by the first layout page is used to lock every subsequent page to
// the same version of the blob.
func TestGetLayoutPagerLocksETagAcrossPages(t *testing.T) {
	etag := azcore.ETag("layout-etag")
	f := newFakeLayoutResponder(buildLayout(5, 100, 2, &etag), nil)
	client := newFakeLayoutClient(t, f)

	// nil options must be supported: the pager has to synthesize its own options to carry the marker.
	pager := client.GetLayoutPager(nil)
	pages := 0
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		pages++
	}
	require.Equal(t, 2, pages, "the fake responder splits the layout into two pages")

	headers := f.layoutIfMatchHeaders()
	require.Len(t, headers, 2)
	require.Empty(t, headers[0], "the initial page has no ETag to lock on yet")
	require.Equal(t, string(etag), headers[1], "subsequent pages must be locked to the initial ETag")
}

// TestGetLayoutPagerUserIfMatchWins verifies that a caller supplied If-Match is honored on every
// page and that the caller's access conditions are not mutated by the pager.
func TestGetLayoutPagerUserIfMatchWins(t *testing.T) {
	etag := azcore.ETag("layout-etag")
	f := newFakeLayoutResponder(buildLayout(5, 100, 2, &etag), nil)
	client := newFakeLayoutClient(t, f)

	userETag := azcore.ETag("user-etag")
	mac := &ModifiedAccessConditions{IfMatch: &userETag}
	options := &GetLayoutOptions{AccessConditions: &AccessConditions{ModifiedAccessConditions: mac}}

	pager := client.GetLayoutPager(options)
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		require.NoError(t, err)
	}

	headers := f.layoutIfMatchHeaders()
	require.Len(t, headers, 2)
	for _, v := range headers {
		require.Equal(t, string(userETag), v, "the caller's If-Match must be used on every page")
	}
	require.Same(t, &userETag, mac.IfMatch, "the caller's access conditions must not be mutated")
}
