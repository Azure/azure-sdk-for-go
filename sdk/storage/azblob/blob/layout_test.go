// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/stretchr/testify/require"
)

// This file covers the pure layout logic in layout.go: parsing a GetLayout response into a layout,
// choosing an endpoint for an offset, and deciding when a cached layout should be refreshed.
// Client-level behavior (what DownloadBuffer puts on the wire) lives in client_layout_test.go.

// createMockPager creates a mock pager for testing getLayout function
func createMockPager(responses []GetLayoutResponse, err error) *runtime.Pager[GetLayoutResponse] {
	index := 0
	return runtime.NewPager(runtime.PagingHandler[GetLayoutResponse]{
		More: func(resp GetLayoutResponse) bool {
			return resp.BlobLayout.NextMarker != nil && *resp.BlobLayout.NextMarker != ""
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
// getIdealEndpoint

func TestGetIdealEndpoint_EmptyLayoutRanges(t *testing.T) {
	l := layout{
		layoutRanges:  []layoutRange{},
		contentLength: 100,
	}
	result := getIdealEndpoint(50, l)
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
	require.Equal(t, "endpoint1", getIdealEndpoint(0, l))

	// Offset in middle
	require.Equal(t, "endpoint1", getIdealEndpoint(50, l))

	// Offset at end
	require.Equal(t, "endpoint1", getIdealEndpoint(100, l))
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
	require.Equal(t, "endpoint1", getIdealEndpoint(0, l))
	require.Equal(t, "endpoint1", getIdealEndpoint(50, l))
	require.Equal(t, "endpoint1", getIdealEndpoint(99, l))

	// Offset in second range
	require.Equal(t, "endpoint2", getIdealEndpoint(100, l))
	require.Equal(t, "endpoint2", getIdealEndpoint(150, l))
	require.Equal(t, "endpoint2", getIdealEndpoint(199, l))

	// Offset in third range
	require.Equal(t, "endpoint3", getIdealEndpoint(200, l))
	require.Equal(t, "endpoint3", getIdealEndpoint(250, l))
	require.Equal(t, "endpoint3", getIdealEndpoint(299, l))
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
	require.Equal(t, "ep0", getIdealEndpoint(0, l))
	require.Equal(t, "ep0", getIdealEndpoint(9, l))
	require.Equal(t, "ep1", getIdealEndpoint(10, l))
	require.Equal(t, "ep3", getIdealEndpoint(35, l))
	require.Equal(t, "ep6", getIdealEndpoint(65, l))
	require.Equal(t, "ep6", getIdealEndpoint(69, l))
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

	require.Equal(t, "endpointA", getIdealEndpoint(25, l))
	require.Equal(t, "endpointB", getIdealEndpoint(75, l))
	require.Equal(t, "endpointA", getIdealEndpoint(125, l))
}

// TestGetIdealEndpoint_OffsetOutOfRange documents the behavior for an offset outside the covered
// ranges: the binary search clamps to the first/last range rather than returning "".
func TestGetIdealEndpoint_OffsetOutOfRange(t *testing.T) {
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 99, endpoint: "endpoint1"},
			{start: 100, end: 199, endpoint: "endpoint2"},
		},
		contentLength: 200,
	}
	require.Equal(t, "endpoint2", getIdealEndpoint(500, l))
	require.Equal(t, "endpoint1", getIdealEndpoint(-1, l))
}

// ======================================================================================== //
// getLayout

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
			ContentLength:     &contentLength,
			BlobContentLength: &contentLength,
			ETag:              &etag,
		},
	}

	pager := createMockPager(responses, nil)

	result, expiry, err := getLayout(ctx, pager)

	require.NoError(t, err)
	require.Len(t, result.layoutRanges, 2)
	require.Equal(t, int64(1000), result.contentLength)
	require.NotNil(t, result.eTag)
	require.Equal(t, etag, *result.eTag)
	require.False(t, expiry.IsZero())
	require.False(t, result.Fallback())

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
			ContentLength:     &contentLength,
			BlobContentLength: &contentLength,
			ETag:              &etag,
		},
	}

	pager := createMockPager(responses, nil)

	result, expiry, err := getLayout(ctx, pager)

	require.NoError(t, err)
	require.Len(t, result.layoutRanges, 0)
	require.Equal(t, int64(500), result.contentLength)
	require.NotNil(t, result.eTag)
	require.Equal(t, etag, *result.eTag)
	require.False(t, expiry.IsZero())
	require.False(t, result.Fallback())
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
			ContentLength:     &contentLength,
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
			ContentLength:     &contentLength,
			BlobContentLength: &contentLength,
			ETag:              &etag,
		},
	}

	pager := createMockPager(responses, nil)

	result, expiry, err := getLayout(ctx, pager)

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

	result, expiry, err := getLayout(ctx, pager)

	require.Error(t, err)
	require.Equal(t, testErr, err)
	require.Empty(t, result.layoutRanges)
	require.True(t, expiry.IsZero())
}

// TestGetLayout_UnsupportedIsCached verifies that when the service says it can't provide a layout
// (400 or 5xx), getLayout returns a fallback layout with a nil error so temporal.Resource caches
// the decision instead of contacting the service on every call.
func TestGetLayout_UnsupportedIsCached(t *testing.T) {
	for _, statusCode := range []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprintf("status_%d", statusCode), func(t *testing.T) {
			pager := createMockPager(nil, &azcore.ResponseError{StatusCode: statusCode})

			result, expiry, err := getLayout(context.Background(), pager)

			require.NoError(t, err)
			require.True(t, result.Fallback())
			require.False(t, expiry.IsZero())
			require.False(t, shouldRefreshLayout(result, context.Background()), "a cached fallback must not be eagerly refreshed")
		})
	}
}

// TestGetLayout_ErrorStatusClassification pins down which status codes are cached as a fallback and
// which are surfaced to the caller. 499 in particular guards the `>= 500` boundary.
func TestGetLayout_ErrorStatusClassification(t *testing.T) {
	cached := []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusServiceUnavailable, 599}
	propagated := []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound,
		http.StatusConflict, http.StatusPreconditionFailed, 499}

	for _, sc := range cached {
		t.Run(fmt.Sprintf("cached_%d", sc), func(t *testing.T) {
			result, expiry, err := getLayout(context.Background(), createMockPager(nil, &azcore.ResponseError{StatusCode: sc}))
			require.NoError(t, err)
			require.True(t, result.Fallback())
			require.False(t, expiry.IsZero())
		})
	}

	for _, sc := range propagated {
		t.Run(fmt.Sprintf("propagated_%d", sc), func(t *testing.T) {
			result, expiry, err := getLayout(context.Background(), createMockPager(nil, &azcore.ResponseError{StatusCode: sc}))
			require.Error(t, err)
			require.False(t, result.Fallback(), "a propagated error must not produce a cacheable fallback")
			require.True(t, expiry.IsZero())
		})
	}
}

// TestGetLayout_MissingBlobContentLength verifies a response without the x-ms-blob-content-length
// header doesn't panic.
func TestGetLayout_MissingBlobContentLength(t *testing.T) {
	etag := azcore.ETag("etag")
	responses := []GetLayoutResponse{
		{
			BlobLayout: generated.BlobLayout{
				Endpoints: &generated.BlobLayoutEndpoints{
					Endpoint: []*generated.BlobLayoutEndpointsEndpointItem{
						{Index: to.Ptr(int32(0)), Value: to.Ptr("endpoint1")},
					},
				},
				Ranges: &generated.BlobLayoutRanges{
					Range: []*generated.BlobLayoutRangesRangeItem{
						{Start: to.Ptr(int64(0)), End: to.Ptr(int64(99)), EndpointIndex: to.Ptr(int32(0))},
					},
				},
			},
			// BlobContentLength intentionally omitted
			ETag: &etag,
		},
	}

	require.NotPanics(t, func() {
		result, _, err := getLayout(context.Background(), createMockPager(responses, nil))
		require.NoError(t, err)
		require.Zero(t, result.contentLength)
		require.Len(t, result.layoutRanges, 1)
	})
}

// ======================================================================================== //
// shouldRefreshLayout

func TestShouldRefreshLayout(t *testing.T) {
	ctx := context.Background()

	t.Run("fallback is never eagerly refreshed", func(t *testing.T) {
		// Even an already-expired fallback returns false; temporal.Resource refreshes it on expiry.
		require.False(t, shouldRefreshLayout(layout{fallback: true, expiry: time.Now().Add(-time.Hour)}, ctx))
		require.False(t, shouldRefreshLayout(layout{fallback: true, expiry: time.Now().Add(time.Hour)}, ctx))
	})

	t.Run("fresh layout is not refreshed", func(t *testing.T) {
		require.False(t, shouldRefreshLayout(layout{expiry: time.Now().Add(time.Hour)}, ctx))
	})

	t.Run("layout expiring within 30s is refreshed", func(t *testing.T) {
		require.True(t, shouldRefreshLayout(layout{expiry: time.Now().Add(10 * time.Second)}, ctx))
		require.True(t, shouldRefreshLayout(layout{expiry: time.Now().Add(-time.Second)}, ctx))
	})
}
