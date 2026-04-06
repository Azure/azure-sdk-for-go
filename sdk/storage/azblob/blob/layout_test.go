// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/stretchr/testify/require"
)

func TestGetIdealEndpoint_EmptyLayoutRanges(t *testing.T) {
	l := layout{
		layoutRanges:  []layoutRange{},
		contentLength: 100,
	}
	result := l.getIdealEndpoint(50)
	require.Equal(t, "", result)
}

func TestGetIdealEndpoint_SingleRange(t *testing.T) {
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 100, endpoint: "endpoint1"},
		},
		contentLength: 100,
	}

	// Offset at start
	require.Equal(t, "endpoint1", l.getIdealEndpoint(0))

	// Offset in middle
	require.Equal(t, "endpoint1", l.getIdealEndpoint(50))

	// Offset at end
	require.Equal(t, "endpoint1", l.getIdealEndpoint(100))
}

func TestGetIdealEndpoint_MultipleRanges(t *testing.T) {
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 99, endpoint: "endpoint1"},
			{start: 100, end: 199, endpoint: "endpoint2"},
			{start: 200, end: 299, endpoint: "endpoint3"},
		},
		contentLength: 300,
	}

	// Offset in first range
	require.Equal(t, "endpoint1", l.getIdealEndpoint(0))
	require.Equal(t, "endpoint1", l.getIdealEndpoint(50))
	require.Equal(t, "endpoint1", l.getIdealEndpoint(99))

	// Offset in second range
	require.Equal(t, "endpoint2", l.getIdealEndpoint(100))
	require.Equal(t, "endpoint2", l.getIdealEndpoint(150))
	require.Equal(t, "endpoint2", l.getIdealEndpoint(199))

	// Offset in third range
	require.Equal(t, "endpoint3", l.getIdealEndpoint(200))
	require.Equal(t, "endpoint3", l.getIdealEndpoint(250))
	require.Equal(t, "endpoint3", l.getIdealEndpoint(299))
}

func TestGetIdealEndpoint_BinarySearchBoundary(t *testing.T) {
	// Test with more ranges to exercise binary search properly
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 9, endpoint: "ep0"},
			{start: 10, end: 19, endpoint: "ep1"},
			{start: 20, end: 29, endpoint: "ep2"},
			{start: 30, end: 39, endpoint: "ep3"},
			{start: 40, end: 49, endpoint: "ep4"},
			{start: 50, end: 59, endpoint: "ep5"},
			{start: 60, end: 69, endpoint: "ep6"},
		},
		contentLength: 70,
	}

	// Test boundaries at each range
	require.Equal(t, "ep0", l.getIdealEndpoint(0))
	require.Equal(t, "ep0", l.getIdealEndpoint(9))
	require.Equal(t, "ep1", l.getIdealEndpoint(10))
	require.Equal(t, "ep3", l.getIdealEndpoint(35))
	require.Equal(t, "ep6", l.getIdealEndpoint(65))
	require.Equal(t, "ep6", l.getIdealEndpoint(69))
}

func TestGetIdealEndpoint_SameEndpointDifferentRanges(t *testing.T) {
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 49, endpoint: "endpointA"},
			{start: 50, end: 99, endpoint: "endpointB"},
			{start: 100, end: 149, endpoint: "endpointA"},
		},
		contentLength: 150,
	}

	require.Equal(t, "endpointA", l.getIdealEndpoint(25))
	require.Equal(t, "endpointB", l.getIdealEndpoint(75))
	require.Equal(t, "endpointA", l.getIdealEndpoint(125))
}

func TestGetLayout_SinglePageWithLayout(t *testing.T) {
	ctx := context.Background()
	contentLength := int64(1000)
	etag := azcore.ETag("test-etag")

	responses := []GetLayoutResponse{
		{
			BlobLayout: generated.BlobLayout{
				Endpoints: &generated.BlobLayoutEndpoints{
					Endpoint: []*generated.BlobLayoutEndpointsEndpointItem{
						{Index: to.Ptr(int32(0)), Value: to.Ptr("endpoint1")},
						{Index: to.Ptr(int32(1)), Value: to.Ptr("endpoint2")},
					},
				},
				Ranges: &generated.BlobLayoutRanges{
					Range: []*generated.BlobLayoutRangesRangeItem{
						{Start: to.Ptr(int64(0)), End: to.Ptr(int64(499)), EndpointIndex: to.Ptr(int32(0))},
						{Start: to.Ptr(int64(500)), End: to.Ptr(int64(999)), EndpointIndex: to.Ptr(int32(1))},
					},
				},
			},
			BlobContentLength: &contentLength,
			ETag:              &etag,
		},
	}

	pager := createMockPager(responses, nil)
	state := layoutState{ctx: ctx}

	result, expiry, err := getLayout(state, pager)

	require.NoError(t, err)
	require.Len(t, result.layoutRanges, 2)
	require.Equal(t, int64(1000), result.contentLength)
	require.NotNil(t, result.eTag)
	require.Equal(t, etag, *result.eTag)
	require.False(t, expiry.IsZero())

	// Verify ranges
	require.Equal(t, int64(0), result.layoutRanges[0].start)
	require.Equal(t, int64(499), result.layoutRanges[0].end)
	require.Equal(t, "endpoint1", result.layoutRanges[0].endpoint)
	require.Equal(t, int64(500), result.layoutRanges[1].start)
	require.Equal(t, int64(999), result.layoutRanges[1].end)
	require.Equal(t, "endpoint2", result.layoutRanges[1].endpoint)
}

func TestGetLayout_SinglePageNoLayout(t *testing.T) {
	ctx := context.Background()
	contentLength := int64(500)
	etag := azcore.ETag("no-layout-etag")

	responses := []GetLayoutResponse{
		{
			BlobLayout: generated.BlobLayout{
				Endpoints: &generated.BlobLayoutEndpoints{
					Endpoint: []*generated.BlobLayoutEndpointsEndpointItem{},
				},
			},
			BlobContentLength: &contentLength,
			ETag:              &etag,
		},
	}

	pager := createMockPager(responses, nil)
	state := layoutState{ctx: ctx}

	result, expiry, err := getLayout(state, pager)

	require.NoError(t, err)
	require.Len(t, result.layoutRanges, 0)
	require.Equal(t, int64(500), result.contentLength)
	require.NotNil(t, result.eTag)
	require.Equal(t, etag, *result.eTag)
	require.True(t, expiry.IsZero())
}

func TestGetLayout_MultiplePages(t *testing.T) {
	ctx := context.Background()
	contentLength := int64(3000)
	etag := azcore.ETag("multi-page-etag")

	responses := []GetLayoutResponse{
		{
			BlobLayout: generated.BlobLayout{
				NextMarker: to.Ptr("marker1"),
				Endpoints: &generated.BlobLayoutEndpoints{
					Endpoint: []*generated.BlobLayoutEndpointsEndpointItem{
						{Index: to.Ptr(int32(0)), Value: to.Ptr("endpoint1")},
					},
				},
				Ranges: &generated.BlobLayoutRanges{
					Range: []*generated.BlobLayoutRangesRangeItem{
						{Start: to.Ptr(int64(0)), End: to.Ptr(int64(999)), EndpointIndex: to.Ptr(int32(0))},
					},
				},
			},
			BlobContentLength: &contentLength,
			ETag:              &etag,
		},
		{
			BlobLayout: generated.BlobLayout{
				Endpoints: &generated.BlobLayoutEndpoints{
					Endpoint: []*generated.BlobLayoutEndpointsEndpointItem{
						{Index: to.Ptr(int32(0)), Value: to.Ptr("endpoint2")},
					},
				},
				Ranges: &generated.BlobLayoutRanges{
					Range: []*generated.BlobLayoutRangesRangeItem{
						{Start: to.Ptr(int64(1000)), End: to.Ptr(int64(1999)), EndpointIndex: to.Ptr(int32(0))},
						{Start: to.Ptr(int64(2000)), End: to.Ptr(int64(2999)), EndpointIndex: to.Ptr(int32(0))},
					},
				},
			},
			BlobContentLength: &contentLength,
			ETag:              &etag,
		},
	}

	pager := createMockPager(responses, nil)
	state := layoutState{ctx: ctx}

	result, expiry, err := getLayout(state, pager)

	require.NoError(t, err)
	require.Len(t, result.layoutRanges, 3)
	require.Equal(t, int64(3000), result.contentLength)
	require.False(t, expiry.IsZero())

	// Verify all ranges from both pages
	require.Equal(t, "endpoint1", result.layoutRanges[0].endpoint)
	require.Equal(t, "endpoint2", result.layoutRanges[1].endpoint)
	require.Equal(t, "endpoint2", result.layoutRanges[2].endpoint)
}

func TestGetLayout_Error(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("pager error")

	pager := createMockPager(nil, testErr)
	state := layoutState{ctx: ctx}

	result, expiry, err := getLayout(state, pager)

	require.Error(t, err)
	require.Equal(t, testErr, err)
	require.Empty(t, result.layoutRanges)
	require.True(t, expiry.IsZero())
}

// createMockPager creates a mock pager for testing getLayout function
func createMockPager(responses []GetLayoutResponse, err error) *runtime.Pager[GetLayoutResponse] {
	index := 0
	return runtime.NewPager(runtime.PagingHandler[GetLayoutResponse]{
		More: func(resp GetLayoutResponse) bool {
			return resp.NextMarker != nil && *resp.NextMarker != ""
		},
		Fetcher: func(ctx context.Context, current *GetLayoutResponse) (GetLayoutResponse, error) {
			if err != nil {
				return GetLayoutResponse{}, err
			}
			if index >= len(responses) {
				return GetLayoutResponse{}, errors.New("no more pages")
			}
			resp := responses[index]
			index++
			return resp, nil
		},
	})
}

// ======================================================================================== //
// Helper methods for layout mock tests
type fakeLayoutResponder struct {
	l                     layout
	layoutResponses       map[string]*http.Response
	getPropertiesResponse *http.Response

	// populated by the pipeline policy
	layoutCalls         int
	getPropertiesCalled bool
	localityGets        int
	normalGets          int
}

func (f *fakeLayoutResponder) reset() {
	f.layoutCalls = 0
	f.getPropertiesCalled = false
	f.localityGets = 0
	f.normalGets = 0
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
		f.layoutCalls++
		marker := qp.Get("marker")
		return f.layoutResponses[marker], nil
	}

	// Get properties
	if req.Method == http.MethodHead {
		f.getPropertiesCalled = true
		return f.getPropertiesResponse, nil
	}

	// Validate that the request range is going to the right layout
	if req.Method == http.MethodGet {
		// If the request Host is different from the URL host
		if req.Host != req.URL.Host {
			f.localityGets++
		} else {
			f.normalGets++
		}
		// Download
		return &http.Response{
			StatusCode: http.StatusPartialContent,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}, nil
	}
	return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
}

func newMockLayoutResponse(contentLength int64, eTag string, layout generated.BlobLayout, statusCode int) *http.Response {
	if statusCode == 0 || statusCode == http.StatusOK {
		data, _ := xml.Marshal(layout)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(data)),
			Header: http.Header{
				"X-Ms-Blob-Content-Length": []string{fmt.Sprintf("%d", contentLength)},
				"ETag":                     []string{eTag},
			},
		}
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}
}

func newMockGetPropertiesResponse(contentLength int64, eTag string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Length": []string{fmt.Sprintf("%d", contentLength)},
			"ETag":           []string{eTag},
		},
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
		EnableLayoutAwareRouting: true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "412")

	// 400 should trigger a fallback to get  properties
	f.reset()
	f.layoutResponses = map[string]*http.Response{"": newMockLayoutResponse(0, "etag", generated.BlobLayout{}, http.StatusBadRequest)}
	f.getPropertiesResponse = newMockGetPropertiesResponse(0, "etag")
	_, err = client.DownloadBuffer(context.Background(), buff, &DownloadBufferOptions{
		EnableLayoutAwareRouting: true,
	})
	require.NoError(t, err)
	require.Equal(t, 1, f.layoutCalls)
	require.True(t, f.getPropertiesCalled)
	require.Zero(t, f.localityGets)
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
		EnableLayoutAwareRouting: true,
	})
	require.NoError(t, err)
	require.Equal(t, 1, f.layoutCalls)
	require.False(t, f.getPropertiesCalled)
	require.Equal(t, 1, f.normalGets)
	require.Zero(t, f.localityGets)
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
	client, err := NewClientWithNoCredential("https://fake.blob.core.windows.net/container/blob", &ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: f,
		},
	})
	require.NoError(t, err)

	buff := make([]byte, 300)
	_, err = client.DownloadBuffer(context.Background(), buff, &DownloadBufferOptions{
		EnableLayoutAwareRouting: true,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, f.layoutCalls, 1)
	require.False(t, f.getPropertiesCalled)
	require.Greater(t, f.localityGets, 0)
	require.Zero(t, f.normalGets)
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
	client, err := NewClientWithNoCredential("https://fake.blob.core.windows.net/container/blob", &ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: f,
		},
	})
	require.NoError(t, err)

	buff := make([]byte, 500)
	_, err = client.DownloadBuffer(context.Background(), buff, &DownloadBufferOptions{
		EnableLayoutAwareRouting: true,
	})
	require.NoError(t, err)
	// With maxRangesPerPage=3 in splitLayoutToPages, 5 ranges should create 2 pages
	require.GreaterOrEqual(t, f.layoutCalls, 2)
	require.False(t, f.getPropertiesCalled)
	require.Greater(t, f.localityGets, 0)
	require.Zero(t, f.normalGets)
}
