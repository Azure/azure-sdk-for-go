// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
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
					Endpoint: []*generated.BlobLayoutEndpoint{
						{Index: to.Ptr(int32(0)), Value: to.Ptr("endpoint1")},
						{Index: to.Ptr(int32(1)), Value: to.Ptr("endpoint2")},
					},
				},
				Ranges: &generated.BlobLayoutRanges{
					Range: []*generated.BlobLayoutRange{
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

	result, err := getLayout(ctx, pager)

	require.NoError(t, err)
	require.Len(t, result.layoutRanges, 2)
	require.Equal(t, int64(1000), result.contentLength)
	require.NotNil(t, result.eTag)
	require.Equal(t, etag, *result.eTag)
	require.False(t, result.fallback)

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
					Endpoint: []*generated.BlobLayoutEndpoint{},
				},
			},
			ContentLength:     &contentLength,
			BlobContentLength: &contentLength,
			ETag:              &etag,
		},
	}

	pager := createMockPager(responses, nil)

	result, err := getLayout(ctx, pager)

	require.NoError(t, err)
	require.Len(t, result.layoutRanges, 0)
	require.Equal(t, int64(500), result.contentLength)
	require.NotNil(t, result.eTag)
	require.Equal(t, etag, *result.eTag)
	require.False(t, result.fallback)
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
					Endpoint: []*generated.BlobLayoutEndpoint{
						{Index: to.Ptr(int32(0)), Value: to.Ptr("endpoint1")},
					},
				},
				Ranges: &generated.BlobLayoutRanges{
					Range: []*generated.BlobLayoutRange{
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
					Endpoint: []*generated.BlobLayoutEndpoint{
						{Index: to.Ptr(int32(0)), Value: to.Ptr("endpoint2")},
					},
				},
				Ranges: &generated.BlobLayoutRanges{
					Range: []*generated.BlobLayoutRange{
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

	result, err := getLayout(ctx, pager)

	require.NoError(t, err)
	require.Len(t, result.layoutRanges, 3)
	require.Equal(t, int64(3000), result.contentLength)

	// Verify all ranges from both pages
	require.Equal(t, "endpoint1", result.layoutRanges[0].endpoint)
	require.Equal(t, "endpoint2", result.layoutRanges[1].endpoint)
	require.Equal(t, "endpoint2", result.layoutRanges[2].endpoint)
}

func TestGetLayout_Error(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("pager error")

	pager := createMockPager(nil, testErr)

	result, err := getLayout(ctx, pager)

	require.Error(t, err)
	require.Equal(t, testErr, err)
	require.Empty(t, result.layoutRanges)
}

// TestGetLayout_UnsupportedIsCached verifies that when the service says it can't provide a layout
// (400 or 5xx), getLayout returns a fallback layout with a nil error so the layout cache keeps
// the decision instead of contacting the service on every call.
func TestGetLayout_UnsupportedIsCached(t *testing.T) {
	for _, statusCode := range []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprintf("status_%d", statusCode), func(t *testing.T) {
			pager := createMockPager(nil, &azcore.ResponseError{StatusCode: statusCode})

			result, err := getLayout(context.Background(), pager)

			require.NoError(t, err)
			require.True(t, result.fallback)
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
			result, err := getLayout(context.Background(), createMockPager(nil, &azcore.ResponseError{StatusCode: sc}))
			require.NoError(t, err)
			require.True(t, result.fallback)
		})
	}

	for _, sc := range propagated {
		t.Run(fmt.Sprintf("propagated_%d", sc), func(t *testing.T) {
			result, err := getLayout(context.Background(), createMockPager(nil, &azcore.ResponseError{StatusCode: sc}))
			require.Error(t, err)
			require.False(t, result.fallback, "a propagated error must not produce a cacheable fallback")
		})
	}
}

// ======================================================================================== //
// newLayoutCache

// layoutTestClock replaces the layout clock for the duration of a test.
type layoutTestClock struct {
	mu  sync.Mutex
	now time.Time
}

func useLayoutTestClock(t *testing.T) *layoutTestClock {
	c := &layoutTestClock{now: time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)}
	prev := layoutNow
	layoutNow = c.Now
	t.Cleanup(func() { layoutNow = prev })
	return c
}

func (c *layoutTestClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *layoutTestClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// layoutFetcher hands out layouts whose single range is served by "endpoint-<n>" for the nth
// fetch, or the result of fail when it is set.
type layoutFetcher struct {
	calls atomic.Int32
	gate  chan struct{}
	fail  func(n int) (layout, error)
}

func (f *layoutFetcher) fetch(context.Context) (layout, error) {
	n := int(f.calls.Add(1))
	if f.gate != nil {
		<-f.gate
	}
	if f.fail != nil {
		if l, err := f.fail(n); err != nil || l.fallback {
			return l, err
		}
	}
	return layout{layoutRanges: []layoutRange{{start: 0, end: 99, endpoint: fmt.Sprintf("endpoint-%d", n)}}}, nil
}

func cachedEndpoint(t *testing.T, c interface {
	Get(context.Context) (layout, error)
}) string {
	t.Helper()
	l, err := c.Get(context.Background())
	require.NoError(t, err)
	return getIdealEndpoint(0, l)
}

func TestLayoutCacheLifetimeAndProactiveRefresh(t *testing.T) {
	clock := useLayoutTestClock(t)
	f := &layoutFetcher{}
	cache := newLayoutCache(f.fetch)

	require.Equal(t, "endpoint-1", cachedEndpoint(t, cache))

	// Valid and not yet due for refresh until 30 seconds before the 5 minute lifetime ends.
	clock.Advance(4*time.Minute + 29*time.Second)
	require.Equal(t, "endpoint-1", cachedEndpoint(t, cache))
	require.Equal(t, int32(1), f.calls.Load())

	// In the refresh window, consumers keep the current layout while the refresh runs in the
	// background.
	f.gate = make(chan struct{})
	clock.Advance(2 * time.Second)
	require.Equal(t, "endpoint-1", cachedEndpoint(t, cache))
	require.Eventually(t, func() bool { return f.calls.Load() == 2 }, 5*time.Second, time.Millisecond)
	require.Equal(t, "endpoint-1", cachedEndpoint(t, cache), "the old layout is used while the refresh runs")
	close(f.gate)
	require.Eventually(t, func() bool { return cachedEndpoint(t, cache) == "endpoint-2" }, 5*time.Second, time.Millisecond)
}

func TestLayoutCacheLifetimeStartsWhenFetchCompletes(t *testing.T) {
	clock := useLayoutTestClock(t)
	var calls atomic.Int32
	cache := newLayoutCache(func(context.Context) (layout, error) {
		calls.Add(1)
		// a slow fetch: the layout's lifetime must not include the time spent fetching it
		clock.Advance(time.Minute)
		return layout{layoutRanges: []layoutRange{{start: 0, end: 99, endpoint: "e"}}}, nil
	})
	require.Equal(t, "e", cachedEndpoint(t, cache))

	// 4m29s after the fetch completed (5m29s after it started) the layout is still current and
	// not yet due for refresh.
	clock.Advance(4*time.Minute + 29*time.Second)
	require.Equal(t, "e", cachedEndpoint(t, cache))
	require.Equal(t, int32(1), calls.Load())
}

func TestLayoutCacheFailedBackgroundRefreshKeepsLayout(t *testing.T) {
	clock := useLayoutTestClock(t)
	f := &layoutFetcher{fail: func(n int) (layout, error) {
		if n > 1 {
			return layout{}, errors.New("refresh failed")
		}
		return layout{}, nil
	}}
	cache := newLayoutCache(f.fetch)
	require.Equal(t, "endpoint-1", cachedEndpoint(t, cache))

	clock.Advance(4*time.Minute + 31*time.Second)
	require.Equal(t, "endpoint-1", cachedEndpoint(t, cache))
	require.Eventually(t, func() bool { return f.calls.Load() == 2 }, 5*time.Second, time.Millisecond)
	for range 5 {
		require.Equal(t, "endpoint-1", cachedEndpoint(t, cache), "a failed background refresh doesn't fail consumers")
	}
}

func TestLayoutCacheCachesUnavailableAndEmptyLayouts(t *testing.T) {
	for _, tt := range []struct {
		name   string
		result layout
	}{
		{"Unavailable", layout{fallback: true}},
		{"NoLayout", layout{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			clock := useLayoutTestClock(t)
			var calls atomic.Int32
			cache := newLayoutCache(func(context.Context) (layout, error) {
				calls.Add(1)
				return tt.result, nil
			})
			for range 10 {
				require.Equal(t, "", cachedEndpoint(t, cache), "reads go to the configured endpoint")
			}
			clock.Advance(4 * time.Minute)
			require.Equal(t, "", cachedEndpoint(t, cache))
			require.Equal(t, int32(1), calls.Load(), "the result is cached for the layout's lifetime")
		})
	}
}

func TestLayoutCacheConcurrentConsumersFetchOnce(t *testing.T) {
	useLayoutTestClock(t)
	f := &layoutFetcher{gate: make(chan struct{})}
	cache := newLayoutCache(f.fetch)

	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l, err := cache.Get(context.Background())
			require.NoError(t, err)
			require.Equal(t, "endpoint-1", getIdealEndpoint(0, l))
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(f.gate)
	wg.Wait()
	require.Equal(t, int32(1), f.calls.Load(), "concurrent chunks don't stampede Get Layout")
}

func TestLayoutCacheExpiredLayoutIsRefetched(t *testing.T) {
	clock := useLayoutTestClock(t)
	f := &layoutFetcher{}
	cache := newLayoutCache(f.fetch)
	require.Equal(t, "endpoint-1", cachedEndpoint(t, cache))

	clock.Advance(5 * time.Minute)
	require.Equal(t, "endpoint-2", cachedEndpoint(t, cache), "an expired layout is never used")
}
