// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

// layoutState is the state needed to getLayout the layout
type layoutState struct {
	client *Client
	opts   *GetLayoutOptions
	ctx    context.Context
}

func getLayout(state layoutState) (layout, time.Time, error) {
	layoutRanges := make([]layoutRange, 0)

	pager := state.client.GetLayoutPager(state.opts)
	var contentLength int64
	var eTag *azcore.ETag
	for pager.More() {
		resp, err := pager.NextPage(state.ctx)
		if err != nil {
			return layout{}, time.Time{}, err
		}
		contentLength = *resp.ContentLength
		eTag = resp.ETag
		if len(resp.BlobLayout.Endpoints.Endpoint) == 0 {
			// No layout means we can download the whole blob from the primary endpoint.
			return layout{contentLength: contentLength, eTag: eTag}, time.Time{}, nil
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

	// Expire the cache after 9 minutes so that we getLayout the layout at 4 minutes by default.
	// The default getLayout time of temporal resource is 5 minutes.
	return layout{layoutRanges: layoutRanges, contentLength: contentLength, eTag: eTag}, time.Now().Add(9 * time.Minute), nil
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
