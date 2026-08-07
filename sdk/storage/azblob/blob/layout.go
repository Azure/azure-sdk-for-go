// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
)

var layoutRefresh = 5 * time.Minute

type layoutRange struct {
	start    int64
	end      int64
	endpoint string
}

type layout struct {
	layoutRanges  []layoutRange
	contentLength int64
	eTag          *azcore.ETag
	fallback      bool
	expiry        time.Time
}

func (l *layout) Fallback() bool {
	return l.fallback
}

func (l *layout) Expiry() time.Time {
	return l.expiry
}

func getLayout(ctx context.Context, pager *runtime.Pager[GetLayoutResponse]) (layout, time.Time, error) {
	layoutRanges := make([]layoutRange, 0)

	var contentLength int64
	var eTag *azcore.ETag
	expiry := time.Now().Add(layoutRefresh)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			// A 400 or 5xx means the service can't provide a layout. Return a fallback layout with a
			// nil error so temporal.Resource caches the decision; returning an error would leave the
			// resource unset and cause every subsequent call to hit the service again.
			if sc := bloberror.GetStatusCode(err); sc == http.StatusBadRequest || sc >= 500 {
				return layout{fallback: true, expiry: expiry}, expiry, nil
			}
			return layout{}, time.Time{}, err
		}
		if resp.BlobContentLength != nil {
			contentLength = *resp.BlobContentLength
		}
		if eTag == nil {
			eTag = resp.ETag
		}
		if resp.BlobLayout.Endpoints == nil || resp.BlobLayout.Endpoints.Endpoint == nil || len(resp.BlobLayout.Endpoints.Endpoint) == 0 ||
			resp.BlobLayout.Ranges == nil || resp.BlobLayout.Ranges.Range == nil || len(resp.BlobLayout.Ranges.Range) == 0 {
			// No layout means we can download the whole blob from the primary endpoint.
			return layout{contentLength: contentLength, eTag: eTag, expiry: expiry}, expiry, nil
		}
		endpoints := make([]string, len(resp.BlobLayout.Endpoints.Endpoint))
		for _, ep := range resp.BlobLayout.Endpoints.Endpoint {
			endpoints[*ep.Index] = *ep.Value
		}
		for _, r := range resp.BlobLayout.Ranges.Range {
			lr := layoutRange{
				start:    *r.Start,
				end:      *r.End,
				endpoint: endpoints[*r.EndpointIndex],
			}
			layoutRanges = append(layoutRanges, lr)
		}
	}

	return layout{layoutRanges: layoutRanges, contentLength: contentLength, eTag: eTag, expiry: expiry}, expiry, nil
}

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

func shouldRefreshLayout(resource layout, _ context.Context) bool {
	if resource.Fallback() {
		// A fallback layout is a cached "layout is unavailable" decision. Refreshing it early would
		// contact the service before the decision was due to be reconsidered; let it expire instead.
		return false
	}
	return resource.Expiry().Add(-30 * time.Second).Before(time.Now())
}
