// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
)

type layoutRange struct {
	start    int64
	end      int64
	endpoint string
}

type layout = []layoutRange

// layoutState is the state needed to refresh the layout
type layoutState struct {
	client *Client
	opts   *GetLayoutOptions
	ctx    context.Context
}

// layoutCache manages layout caching for a single download operation
type layoutCache = *temporal.Resource[*layout, *layoutState]

func refresh(state layoutState) (layout layout, expiration time.Time, err error) {
	layout = make([]layoutRange, 0)

	pager := state.client.getLayoutPager(state.opts)

	for pager.More() {
		resp, err := pager.NextPage(state.ctx)
		if err != nil {
			return nil, time.Time{}, err
		}
		if len(resp.BlobLayout.Endpoints.Endpoint) == 0 {
			// No layout means we can download the whole blob from the primary endpoint.
			return layout, time.Time{}, nil
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
			layout = append(layout, lr)
		}
	}

	// Expire the cache after 9 minutes so that we refresh the layout at 4 minutes by default.
	// The default refresh time of temporal resource is 5 minutes.
	return layout, time.Now().Add(9 * time.Minute), nil
}
