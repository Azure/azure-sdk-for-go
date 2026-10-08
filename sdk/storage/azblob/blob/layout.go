// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/internal/autorefresh"
)

func getStatusCode(err error) int {
	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return 0
	}

	return respErr.StatusCode
}

// These are variables so tests can manipulate them.
var (
	// layoutLifetime is how long the service considers a layout current.
	layoutLifetime = 5 * time.Minute

	// layoutRefreshBuffer is how long before a cached layout expires that a background refresh
	// starts.
	layoutRefreshBuffer = 30 * time.Second

	// layoutBackgroundAcquireTimeout bounds a background layout refresh.
	layoutBackgroundAcquireTimeout = 30 * time.Second

	// layoutNow is the clock for layout caches.
	layoutNow = time.Now
)

type layoutRange struct {
	start    int64
	end      int64
	endpoint string
}

// layout is a cached Get Layout result. It routes reads only when layoutRanges is non-empty: a
// blob with no layout (an empty or 204 response) and a layout the service couldn't provide
// (fallback) both send every read to the client's configured endpoint, and both are cached for
// the layout's full lifetime so the rest of the download doesn't ask again.
type layout struct {
	layoutRanges  []layoutRange
	contentLength int64
	eTag          *azcore.ETag
	fallback      bool
}

// newLayoutCache returns a cache of the layout fetch produces, valid for layoutLifetime from when
// the fetch completes and refreshed in the background layoutRefreshBuffer before that.
func newLayoutCache(fetch func(context.Context) (layout, error)) *autorefresh.Cache[layout] {
	return autorefresh.New(func(ctx context.Context) (autorefresh.Entry[layout], error) {
		l, err := fetch(ctx)
		if err != nil {
			return autorefresh.Entry[layout]{}, err
		}
		expires := layoutNow().Add(layoutLifetime)
		return autorefresh.Entry[layout]{Value: l, ExpiresOn: expires, RefreshOn: expires.Add(-layoutRefreshBuffer)}, nil
	}, &autorefresh.Options{
		BackgroundAcquireTimeout: layoutBackgroundAcquireTimeout,
		Now:                      layoutNow,
	})
}

// getLayout pages through the layout. A 400 or 5xx means the service can't provide one, which is
// returned as a fallback layout rather than an error so that the decision is cached instead of
// asked again for every read. Any other error is returned.
func getLayout(ctx context.Context, pager *runtime.Pager[GetLayoutResponse]) (layout, error) {
	layoutRanges := make([]layoutRange, 0)

	var contentLength int64
	var eTag *azcore.ETag
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			if sc := getStatusCode(err); sc == http.StatusBadRequest || sc >= 500 {
				return layout{fallback: true}, nil
			}
			return layout{}, err
		}
		if resp.BlobContentLength != nil {
			contentLength = *resp.BlobContentLength
		}
		if eTag == nil {
			eTag = resp.ETag
		}
		if resp.Endpoints == nil || len(resp.Endpoints.Endpoint) == 0 ||
			resp.Ranges == nil || len(resp.Ranges.Range) == 0 {
			// No layout means we can download the whole blob from the primary endpoint.
			return layout{contentLength: contentLength, eTag: eTag}, nil
		}
		endpoints := make([]string, len(resp.Endpoints.Endpoint))
		for _, ep := range resp.Endpoints.Endpoint {
			endpoints[*ep.Index] = *ep.Value
		}
		for _, r := range resp.Ranges.Range {
			lr := layoutRange{
				start:    *r.Start,
				end:      *r.End,
				endpoint: endpoints[*r.EndpointIndex],
			}
			layoutRanges = append(layoutRanges, lr)
		}
	}
	return layout{layoutRanges: layoutRanges, contentLength: contentLength, eTag: eTag}, nil
}

// getIdealEndpoint returns the endpoint that serves offset, or "" when the layout doesn't route
// reads.
func getIdealEndpoint(offset int64, l layout) string {
	if len(l.layoutRanges) == 0 {
		return ""
	}
	// Binary search to find the first range whose end >= offset
	left, right := 0, len(l.layoutRanges)-1
	for left < right {
		mid := left + (right-left)/2
		if l.layoutRanges[mid].end < offset {
			left = mid + 1
		} else {
			right = mid
		}
	}
	// Range is guaranteed to exist, return its endpoint
	return l.layoutRanges[left].endpoint
}
