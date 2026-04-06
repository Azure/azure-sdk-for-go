// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type layoutRange struct {
	start    int64
	end      int64
	endpoint string
}

type layout struct {
	layoutRanges  []layoutRange
	contentLength int64
	eTag          *azcore.ETag
}

// layoutState is the state needed to refresh the layout
type layoutState struct {
	ctx context.Context
}

// getLayout gets the layout of the blob.
func getLayout(state layoutState, pager *runtime.Pager[GetLayoutResponse]) (layout, time.Time, error) {
	layoutRanges := make([]layoutRange, 0)

	var contentLength int64
	var eTag *azcore.ETag
	for pager.More() {
		resp, err := pager.NextPage(state.ctx)
		if err != nil {
			return layout{}, time.Time{}, err
		}
		contentLength = *resp.BlobContentLength
		if eTag == nil {
			eTag = resp.ETag
		}
		if resp.Endpoints == nil || resp.Endpoints.Endpoint == nil || len(resp.Endpoints.Endpoint) == 0 ||
			resp.Ranges == nil || resp.Ranges.Range == nil || len(resp.Ranges.Range) == 0 {
			// No layout means we can download the whole blob from the primary endpoint.
			return layout{contentLength: contentLength, eTag: eTag}, time.Time{}, nil
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

	// Expire the cache after 9 minutes so that we refresh the layout at 4 minutes by default.
	// The default refresh time of temporal resource is 5 minutes.
	return layout{layoutRanges: layoutRanges, contentLength: contentLength, eTag: eTag}, time.Now().Add(9 * time.Minute), nil
}

// getIdealEndpoint returns the endpoint to download the blob at the given offset based on the layout.
// It considers the first range containing the offset. If the chunk isn't fully available, the service still returns
// correct data.
func (l layout) getIdealEndpoint(offset int64) string {
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
