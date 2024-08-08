// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type respType interface {
	hasMore() bool
	lastID() *string
}

// newOpenAIPager is a pager that handles the OpenAI style of paging, where you pass a "lastID"
// to indicate where in the chain of items you are.
//
// NOTE: the OptionsT/POptionsT is to handle the odd ambiguity in generics where you have a
// pointer-to-something and something can't be converted.
func newOpenAIPager[ResponseT respType, OptionsT any, POptionsT interface {
	*OptionsT
	updateAfter(after *string)
}](
	client *Client,
	nextPageFn func(ctx context.Context, opts POptionsT) (ResponseT, error),
	options POptionsT) *runtime.Pager[ResponseT] {
	var lastID *string

	first := true

	return runtime.NewPager(runtime.PagingHandler[ResponseT]{
		More: func(clmr ResponseT) bool {
			return clmr.hasMore()
		},
		Fetcher: func(ctx context.Context, clmr *ResponseT) (ResponseT, error) {
			newOptions := options

			if newOptions == nil {
				var zero OptionsT
				newOptions = &zero
			}

			if !first {
				// make sure to respect the callers choice on the first time through.
				newOptions.updateAfter(lastID)
			}

			first = false

			resp, err := nextPageFn(ctx, newOptions)

			if err != nil {
				var zero ResponseT
				return zero, err
			}

			lastID = resp.lastID()
			return resp, nil
		},
		Tracer: client.internal.Tracer(),
	})
}
