// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// NewListBatchesPager returns a pager for listing batches.
func (c *Client) NewListBatchesPager(options *ListBatchesOptions) *runtime.Pager[ListBatchesResponse] {
	nextPageFn := func(ctx context.Context, options *ListBatchesOptions) (ListBatchesResponse, error) {
		return c.listBatches(ctx, options)
	}

	return newOpenAIPager(c, nextPageFn, options)
}

func (o *ListBatchesOptions) updateAfter(after *string) { o.After = after }
func (r ListBatchesResponse) lastID() *string           { return r.LastID }
func (r ListBatchesResponse) hasMore() bool             { return *r.HasMore }
