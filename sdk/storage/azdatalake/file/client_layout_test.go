// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/stretchr/testify/require"
)

// This file covers data locality for DataLake: GetLayoutPager, the LayoutEndpoint option on
// DownloadStream, and layout aware routing for the managed downloads. The layout parsing,
// endpoint selection, caching and URL rewriting all live in azblob and are tested there; these
// tests assert that the DataLake client forwards to it, and that the requests DataLake puts on
// the wire are the expected ones.

const (
	fakeFileURL     = "https://fake.dfs.core.windows.net/fs/file.txt"
	fakeBlobHost    = "fake.blob.core.windows.net"
	layoutEndpoint0 = "https://locality0.blob.core.windows.net"
	layoutEndpoint1 = "https://locality1.blob.core.windows.net"
	fakeLayoutETag  = "0x8DD00000000AAAA"
	fakeUserETag    = "0x8DD00000000BBBB"
)

// fakeLayoutRange describes one range of a mocked layout response.
type fakeLayoutRange struct {
	start         int64
	end           int64
	endpointIndex int32
}

// getRequest records what a single data download request looked like on the wire.
type getRequest struct {
	// urlHost is the host the request was actually sent to.
	urlHost string
	// hostHeader is the Host header of the request, which the layout policy sets to the account
	// host when it rewrites the URL.
	hostHeader string
	offset     int64
	count      int64
	ifMatch    string
}

// fakeLayoutTransport serves GetLayout, GetProperties and ranged downloads for a file whose
// content is deterministic, so a download can be verified byte for byte.
type fakeLayoutTransport struct {
	content   []byte
	endpoints []string
	// pages holds the layout pages keyed by the marker that requests them; the empty marker is
	// the first page.
	pages map[string][]fakeLayoutRange
	// nextMarkers maps a page's marker to the marker of the page that follows it.
	nextMarkers map[string]string
	// layoutStatus, when non-zero, is returned for every GetLayout request instead of a page.
	layoutStatus int

	mu                  sync.Mutex
	layoutCalls         int
	layoutMarkers       []string
	layoutIfMatch       []string
	layoutRangeHeaders  []string
	layoutMaxResults    []string
	getPropertiesCalled int
	gets                []getRequest
}

func (f *fakeLayoutTransport) Do(req *http.Request) (*http.Response, error) {
	qp := req.URL.Query()

	if qp.Get("comp") == "layout" {
		f.mu.Lock()
		f.layoutCalls++
		marker := qp.Get("marker")
		f.layoutMarkers = append(f.layoutMarkers, marker)
		f.layoutIfMatch = append(f.layoutIfMatch, rawHeader(req, "If-Match"))
		f.layoutRangeHeaders = append(f.layoutRangeHeaders, rawHeader(req, "x-ms-range"))
		f.layoutMaxResults = append(f.layoutMaxResults, qp.Get("maxresults"))
		status := f.layoutStatus
		f.mu.Unlock()

		if status != 0 {
			return &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(bytes.NewReader(nil)),
				Header:     http.Header{},
				Request:    req,
			}, nil
		}

		ranges, ok := f.pages[marker]
		if !ok {
			return nil, fmt.Errorf("unexpected layout marker %q", marker)
		}
		header := http.Header{}
		header.Set("x-ms-blob-content-length", strconv.FormatInt(int64(len(f.content)), 10))
		header.Set("ETag", fakeLayoutETag)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(layoutXML(ranges, f.endpoints, f.nextMarkers[marker]))),
			Header:     header,
			Request:    req,
		}, nil
	}

	if req.Method == http.MethodHead {
		f.mu.Lock()
		f.getPropertiesCalled++
		f.mu.Unlock()
		header := http.Header{}
		header.Set("Content-Length", strconv.FormatInt(int64(len(f.content)), 10))
		header.Set("ETag", fakeLayoutETag)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Header:     header,
			Request:    req,
		}, nil
	}

	if req.Method == http.MethodGet {
		offset, count := parseRangeHeader(rawHeader(req, "x-ms-range"), int64(len(f.content)))
		f.mu.Lock()
		f.gets = append(f.gets, getRequest{
			urlHost:    req.URL.Host,
			hostHeader: req.Host,
			offset:     offset,
			count:      count,
			ifMatch:    rawHeader(req, "If-Match"),
		})
		f.mu.Unlock()

		header := http.Header{}
		header.Set("Content-Length", strconv.FormatInt(count, 10))
		header.Set("ETag", fakeLayoutETag)
		return &http.Response{
			StatusCode: http.StatusPartialContent,
			Body:       io.NopCloser(bytes.NewReader(f.content[offset : offset+count])),
			Header:     header,
			Request:    req,
		}, nil
	}

	return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.String())
}

// rawHeader reads a header the generated code sets with a non-canonical map key
// (for example "x-ms-range"), which http.Header.Get would not find.
func rawHeader(req *http.Request, key string) string {
	if v := req.Header[key]; len(v) > 0 {
		return v[0]
	}
	return req.Header.Get(key)
}

// snapshot returns a consistent copy of the recorded data download requests.
func (f *fakeLayoutTransport) snapshot() []getRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]getRequest(nil), f.gets...)
}

// parseRangeHeader parses a "bytes=start-end" range header. An absent header means the whole file.
func parseRangeHeader(r string, size int64) (offset, count int64) {
	if r == "" {
		return 0, size
	}
	var start, end int64
	if _, err := fmt.Sscanf(r, "bytes=%d-%d", &start, &end); err != nil {
		return 0, size
	}
	if end >= size {
		end = size - 1
	}
	return start, end - start + 1
}

// layoutXML renders a layout page the way the service does: endpoints and ranges are attributes.
func layoutXML(ranges []fakeLayoutRange, endpoints []string, nextMarker string) []byte {
	sb := &strings.Builder{}
	sb.WriteString(`<?xml version="1.0" encoding="utf-8"?><BlobLayout><Endpoints>`)
	for i, ep := range endpoints {
		fmt.Fprintf(sb, `<Endpoint Index="%d" Value="%s"/>`, i, ep)
	}
	sb.WriteString(`</Endpoints><Ranges>`)
	for _, r := range ranges {
		fmt.Fprintf(sb, `<Range Start="%d" End="%d" EndpointIndex="%d"/>`, r.start, r.end, r.endpointIndex)
	}
	sb.WriteString(`</Ranges>`)
	if nextMarker != "" {
		fmt.Fprintf(sb, `<NextMarker>%s</NextMarker>`, nextMarker)
	}
	sb.WriteString(`</BlobLayout>`)
	return []byte(sb.String())
}

// newFakeContent builds deterministic file content so a download can be verified byte for byte.
func newFakeContent(size int) []byte {
	content := make([]byte, size)
	for i := range content {
		content[i] = byte(i % 251)
	}
	return content
}

// newSingleChunkTransport serves a file whose layout is a single range on one endpoint.
func newSingleChunkTransport(size int) *fakeLayoutTransport {
	return &fakeLayoutTransport{
		content:   newFakeContent(size),
		endpoints: []string{layoutEndpoint0},
		pages: map[string][]fakeLayoutRange{
			"": {{start: 0, end: int64(size) - 1, endpointIndex: 0}},
		},
		nextMarkers: map[string]string{},
	}
}

// newStripedTransport serves a file striped over two endpoints in chunkSize sized ranges.
func newStripedTransport(size, chunkSize int) *fakeLayoutTransport {
	var ranges []fakeLayoutRange
	for offset := 0; offset < size; offset += chunkSize {
		end := offset + chunkSize - 1
		if end > size-1 {
			end = size - 1
		}
		ranges = append(ranges, fakeLayoutRange{
			start:         int64(offset),
			end:           int64(end),
			endpointIndex: int32((offset / chunkSize) % 2),
		})
	}
	return &fakeLayoutTransport{
		content:     newFakeContent(size),
		endpoints:   []string{layoutEndpoint0, layoutEndpoint1},
		pages:       map[string][]fakeLayoutRange{"": ranges},
		nextMarkers: map[string]string{},
	}
}

func newFakeFileClient(t *testing.T, transport policy.Transporter) *file.Client {
	t.Helper()
	client, err := file.NewClientWithNoCredential(fakeFileURL, &file.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)
	return client
}

// hostOf returns the host of an endpoint URL, which is what the layout policy puts in the URL.
func hostOf(endpoint string) string {
	return strings.TrimPrefix(endpoint, "https://")
}

// ======================================================================================== //
// GetLayoutPager

func TestFileGetLayoutPagerReturnsLayout(t *testing.T) {
	transport := newStripedTransport(4096, 1024)
	client := newFakeFileClient(t, transport)

	pager := client.GetLayoutPager(nil)
	require.True(t, pager.More())
	resp, err := pager.NextPage(context.Background())
	require.NoError(t, err)

	require.NotNil(t, resp.BlobContentLength)
	require.Equal(t, int64(4096), *resp.BlobContentLength)
	require.NotNil(t, resp.ETag)
	require.Equal(t, azcore.ETag(fakeLayoutETag), *resp.ETag)

	require.NotNil(t, resp.Endpoints)
	require.Len(t, resp.Endpoints.Endpoint, 2)
	require.Equal(t, layoutEndpoint0, *resp.Endpoints.Endpoint[0].Value)
	require.Equal(t, layoutEndpoint1, *resp.Endpoints.Endpoint[1].Value)

	require.NotNil(t, resp.Ranges)
	require.Len(t, resp.Ranges.Range, 4)
	require.Equal(t, int64(0), *resp.Ranges.Range[0].Start)
	require.Equal(t, int64(1023), *resp.Ranges.Range[0].End)
	require.Equal(t, int32(0), *resp.Ranges.Range[0].EndpointIndex)
	require.Equal(t, int64(1024), *resp.Ranges.Range[1].Start)
	require.Equal(t, int32(1), *resp.Ranges.Range[1].EndpointIndex)

	require.False(t, pager.More())
	require.Equal(t, 1, transport.layoutCalls)
}

func TestFileGetLayoutPagerPagination(t *testing.T) {
	transport := &fakeLayoutTransport{
		content:   newFakeContent(3072),
		endpoints: []string{layoutEndpoint0, layoutEndpoint1},
		pages: map[string][]fakeLayoutRange{
			"":  {{start: 0, end: 1023, endpointIndex: 0}},
			"1": {{start: 1024, end: 2047, endpointIndex: 1}},
			"2": {{start: 2048, end: 3071, endpointIndex: 0}},
		},
		nextMarkers: map[string]string{"": "1", "1": "2"},
	}
	client := newFakeFileClient(t, transport)

	var pages, ranges int
	pager := client.GetLayoutPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		pages++
		ranges += len(resp.Ranges.Range)
	}

	require.Equal(t, 3, pages)
	require.Equal(t, 3, ranges)
	require.Equal(t, 3, transport.layoutCalls)
	// The pager walks the continuation markers the service handed back.
	require.Equal(t, []string{"", "1", "2"}, transport.layoutMarkers)
}

func TestFileGetLayoutPagerLocksETagAcrossPages(t *testing.T) {
	transport := &fakeLayoutTransport{
		content:   newFakeContent(2048),
		endpoints: []string{layoutEndpoint0},
		pages: map[string][]fakeLayoutRange{
			"":  {{start: 0, end: 1023, endpointIndex: 0}},
			"1": {{start: 1024, end: 2047, endpointIndex: 0}},
		},
		nextMarkers: map[string]string{"": "1"},
	}
	client := newFakeFileClient(t, transport)

	pager := client.GetLayoutPager(nil)
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		require.NoError(t, err)
	}

	// The first request carries no If-Match; every subsequent page is locked to the ETag the
	// first response returned, so the layout cannot change underneath the enumeration.
	require.Equal(t, []string{"", fakeLayoutETag}, transport.layoutIfMatch)
}

func TestFileGetLayoutPagerOptionsAreForwarded(t *testing.T) {
	transport := &fakeLayoutTransport{
		content:   newFakeContent(2048),
		endpoints: []string{layoutEndpoint0},
		pages: map[string][]fakeLayoutRange{
			"":  {{start: 0, end: 1023, endpointIndex: 0}},
			"1": {{start: 1024, end: 2047, endpointIndex: 0}},
		},
		nextMarkers: map[string]string{"": "1"},
	}
	client := newFakeFileClient(t, transport)

	userETag := azcore.ETag(fakeUserETag)
	maxResults := int32(1)
	pager := client.GetLayoutPager(&file.GetLayoutOptions{
		MaxResults: &maxResults,
		Range:      &file.HTTPRange{Offset: 512, Count: 1024},
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{IfMatch: &userETag},
		},
	})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		require.NoError(t, err)
	}

	require.Equal(t, 2, transport.layoutCalls)
	// A caller supplied If-Match wins over the ETag captured from the first page.
	require.Equal(t, []string{fakeUserETag, fakeUserETag}, transport.layoutIfMatch)
	require.Equal(t, []string{"bytes=512-1535", "bytes=512-1535"}, transport.layoutRangeHeaders)
	require.Equal(t, []string{"1", "1"}, transport.layoutMaxResults)
}

func TestFileGetLayoutPagerError(t *testing.T) {
	transport := newSingleChunkTransport(1024)
	transport.layoutStatus = http.StatusNotFound
	client := newFakeFileClient(t, transport)

	pager := client.GetLayoutPager(nil)
	_, err := pager.NextPage(context.Background())
	require.Error(t, err)
}

// ======================================================================================== //
// One-shot download: LayoutEndpoint

func TestFileDownloadStreamForwardsLayoutEndpoint(t *testing.T) {
	transport := newSingleChunkTransport(1024)
	client := newFakeFileClient(t, transport)

	resp, err := client.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range:          &file.HTTPRange{Offset: 0, Count: 512},
		LayoutEndpoint: layoutEndpoint0,
	})
	require.NoError(t, err)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, transport.content[:512], data)

	gets := transport.snapshot()
	require.Len(t, gets, 1)
	// The request is sent to the layout endpoint, but still addressed to the account host.
	require.Equal(t, hostOf(layoutEndpoint0), gets[0].urlHost)
	require.Equal(t, fakeBlobHost, gets[0].hostHeader)
	require.Equal(t, 0, transport.layoutCalls)
}

func TestFileDownloadStreamWithoutLayoutEndpoint(t *testing.T) {
	transport := newSingleChunkTransport(1024)
	client := newFakeFileClient(t, transport)

	resp, err := client.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: &file.HTTPRange{Offset: 0, Count: 512},
	})
	require.NoError(t, err)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, transport.content[:512], data)

	gets := transport.snapshot()
	require.Len(t, gets, 1)
	// No layout endpoint: the URL host and the Host header still agree, which is what an
	// unrewritten request looks like -- exactly as before data locality.
	require.Equal(t, fakeBlobHost, gets[0].urlHost)
	require.Equal(t, fakeBlobHost, gets[0].hostHeader)
	require.Equal(t, 0, transport.layoutCalls)
}

// ======================================================================================== //
// Managed downloads: LayoutAwareRouting

// requireStripedRouting asserts every chunk was fetched from the endpoint that owns its offset.
func requireStripedRouting(t *testing.T, gets []getRequest, chunkSize int64) {
	t.Helper()
	for _, g := range gets {
		want := hostOf(layoutEndpoint0)
		if (g.offset/chunkSize)%2 == 1 {
			want = hostOf(layoutEndpoint1)
		}
		require.Equal(t, want, g.urlHost, "chunk at offset %d went to the wrong endpoint", g.offset)
		require.Equal(t, fakeBlobHost, g.hostHeader)
		// Managed downloads ETag lock on the layout response so every chunk sees one file.
		require.Equal(t, fakeLayoutETag, g.ifMatch)
	}
}

func TestFileDownloadBufferWithLayoutAwareRouting(t *testing.T) {
	const size, chunkSize = 8192, 1024
	transport := newStripedTransport(size, chunkSize)
	client := newFakeFileClient(t, transport)

	buffer := make([]byte, size)
	n, err := client.DownloadBuffer(context.Background(), buffer, &file.DownloadBufferOptions{
		ChunkSize:          chunkSize,
		LayoutAwareRouting: file.LayoutAwareRoutingEnabled,
	})
	require.NoError(t, err)
	require.Equal(t, int64(size), n)
	require.Equal(t, transport.content, buffer)

	gets := transport.snapshot()
	require.Len(t, gets, size/chunkSize)
	requireStripedRouting(t, gets, chunkSize)
	// The layout is fetched once and reused for every chunk, and it carries the length, so
	// GetProperties is not needed.
	require.Equal(t, 1, transport.layoutCalls)
	require.Equal(t, 0, transport.getPropertiesCalled)
}

func TestFileDownloadFileWithLayoutAwareRouting(t *testing.T) {
	const size, chunkSize = 8192, 1024
	transport := newStripedTransport(size, chunkSize)
	client := newFakeFileClient(t, transport)

	destPath := filepath.Join(t.TempDir(), "downloaded.txt")
	destFile, err := os.Create(destPath)
	require.NoError(t, err)
	defer func() { require.NoError(t, destFile.Close()) }()

	n, err := client.DownloadFile(context.Background(), destFile, &file.DownloadFileOptions{
		ChunkSize:          chunkSize,
		LayoutAwareRouting: file.LayoutAwareRoutingEnabled,
	})
	require.NoError(t, err)
	require.Equal(t, int64(size), n)

	written, err := os.ReadFile(destPath)
	require.NoError(t, err)
	require.Equal(t, transport.content, written)

	gets := transport.snapshot()
	require.Len(t, gets, size/chunkSize)
	requireStripedRouting(t, gets, chunkSize)
	require.Equal(t, 1, transport.layoutCalls)
}

func TestFileDownloadBufferLayoutAwareRoutingDisabled(t *testing.T) {
	const size, chunkSize = 8192, 1024
	transport := newStripedTransport(size, chunkSize)
	client := newFakeFileClient(t, transport)

	buffer := make([]byte, size)
	n, err := client.DownloadBuffer(context.Background(), buffer, &file.DownloadBufferOptions{
		ChunkSize:          chunkSize,
		LayoutAwareRouting: file.LayoutAwareRoutingDisabled,
	})
	require.NoError(t, err)
	require.Equal(t, int64(size), n)
	require.Equal(t, transport.content, buffer)

	// Disabled means the layout is never requested and every chunk goes to the configured
	// host, with the URL host and Host header still in agreement (no rewrite).
	require.Equal(t, 0, transport.layoutCalls)
	require.Equal(t, 1, transport.getPropertiesCalled)
	gets := transport.snapshot()
	require.Len(t, gets, size/chunkSize)
	for _, g := range gets {
		require.Equal(t, fakeBlobHost, g.urlHost)
		require.Equal(t, fakeBlobHost, g.hostHeader)
		require.Empty(t, g.ifMatch)
	}
}

func TestFileDownloadBufferDefaultIsLayoutAware(t *testing.T) {
	const size, chunkSize = 4096, 1024
	transport := newStripedTransport(size, chunkSize)
	client := newFakeFileClient(t, transport)

	// No data locality options at all: the default (Auto) resolves to enabled, matching Blob.
	buffer := make([]byte, size)
	n, err := client.DownloadBuffer(context.Background(), buffer, &file.DownloadBufferOptions{
		ChunkSize: chunkSize,
	})
	require.NoError(t, err)
	require.Equal(t, int64(size), n)
	require.Equal(t, transport.content, buffer)

	require.Equal(t, 1, transport.layoutCalls)
	requireStripedRouting(t, transport.snapshot(), chunkSize)
}

func TestFileDownloadBufferFallsBackWhenLayoutUnavailable(t *testing.T) {
	const size, chunkSize = 8192, 1024
	transport := newStripedTransport(size, chunkSize)
	// A service that cannot provide a layout answers GetLayout with 400.
	transport.layoutStatus = http.StatusBadRequest
	client := newFakeFileClient(t, transport)

	buffer := make([]byte, size)
	n, err := client.DownloadBuffer(context.Background(), buffer, &file.DownloadBufferOptions{
		ChunkSize: chunkSize,
	})
	require.NoError(t, err)
	require.Equal(t, int64(size), n)
	require.Equal(t, transport.content, buffer)

	// The download completes exactly as it did before data locality: properties, then plain
	// ranged reads against the configured host, with the caller's access conditions untouched.
	require.Equal(t, 1, transport.getPropertiesCalled)
	gets := transport.snapshot()
	require.Len(t, gets, size/chunkSize)
	for _, g := range gets {
		require.Equal(t, fakeBlobHost, g.urlHost)
		require.Equal(t, fakeBlobHost, g.hostHeader)
		require.Empty(t, g.ifMatch)
	}
	// The "layout unavailable" answer is cached, so the chunks do not each retry it.
	require.Equal(t, 1, transport.layoutCalls)
}

func TestFileDownloadBufferNilOptions(t *testing.T) {
	const size = 2048
	transport := newSingleChunkTransport(size)
	client := newFakeFileClient(t, transport)

	buffer := make([]byte, size)
	n, err := client.DownloadBuffer(context.Background(), buffer, nil)
	require.NoError(t, err)
	require.Equal(t, int64(size), n)
	require.Equal(t, transport.content, buffer)
}
