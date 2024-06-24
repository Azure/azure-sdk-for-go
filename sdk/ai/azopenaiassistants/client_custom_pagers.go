//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// NewListMessagesPager returns a pager for messages associated with a thread.
func (c *Client) NewListMessagesPager(threadID string, options *ListMessagesOptions) *runtime.Pager[ListMessagesResponse] {
	nextPageFn := func(ctx context.Context, opts *ListMessagesOptions) (ListMessagesResponse, error) {
		return c.internalListMessages(ctx, threadID, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

func (o *ListMessagesOptions) updateAfter(after *string) { o.After = after }
func (r ListMessagesResponse) lastID() *string           { return r.LastID }
func (r ListMessagesResponse) hasMore() bool             { return *r.HasMore }

// NewListAssistantsPager returns a pager for assistants.
func (c *Client) NewListAssistantsPager(options *ListAssistantsOptions) *runtime.Pager[ListAssistantsResponse] {
	nextPageFn := func(ctx context.Context, opts *ListAssistantsOptions) (ListAssistantsResponse, error) {
		return c.internalListAssistants(ctx, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

func (o *ListAssistantsOptions) updateAfter(after *string) { o.After = after }
func (r ListAssistantsResponse) lastID() *string           { return r.LastID }
func (r ListAssistantsResponse) hasMore() bool             { return *r.HasMore }

// NewListRunStepsPager returns a pager for a Run's steps.
func (c *Client) NewListRunStepsPager(threadID string, runID string, options *ListRunStepsOptions) *runtime.Pager[ListRunStepsResponse] {
	nextPageFn := func(ctx context.Context, opts *ListRunStepsOptions) (ListRunStepsResponse, error) {
		return c.internalListRunSteps(ctx, threadID, runID, opts)
	}

	return newOpenAIPager(c, nextPageFn, options)
}

func (o *ListRunStepsOptions) updateAfter(after *string) { o.After = after }
func (r ListRunStepsResponse) lastID() *string           { return r.LastID }
func (r ListRunStepsResponse) hasMore() bool             { return *r.HasMore }

// NewListRunsPager returns a pager for a Thread's runs.
func (c *Client) NewListRunsPager(threadID string, options *ListRunsOptions) *runtime.Pager[ListRunsResponse] {
	nextPageFn := func(ctx context.Context, opts *ListRunsOptions) (ListRunsResponse, error) {
		return c.internalListRuns(ctx, threadID, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

func (o *ListRunsOptions) updateAfter(after *string) { o.After = after }
func (r ListRunsResponse) lastID() *string           { return r.LastID }
func (r ListRunsResponse) hasMore() bool             { return *r.HasMore }

// NewListVectorStoresPager returns a pager for a VectorStores.
func (c *Client) NewListVectorStoresPager(options *ListVectorStoresOptions) *runtime.Pager[ListVectorStoresResponse] {
	nextPageFn := func(ctx context.Context, opts *ListVectorStoresOptions) (ListVectorStoresResponse, error) {
		return c.internalListVectorStores(ctx, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

func (o *ListVectorStoresOptions) updateAfter(after *string) { o.After = after }
func (r ListVectorStoresResponse) lastID() *string           { return r.LastID }
func (r ListVectorStoresResponse) hasMore() bool             { return *r.HasMore }

// NewListVectorStoreFilesPager returns a pager for a vector store files.
func (c *Client) NewListVectorStoreFilesPager(vectorStoreID string, options *ListVectorStoreFilesOptions) *runtime.Pager[ListVectorStoreFilesResponse] {
	nextPageFn := func(ctx context.Context, opts *ListVectorStoreFilesOptions) (ListVectorStoreFilesResponse, error) {
		return c.internalListVectorStoreFiles(ctx, vectorStoreID, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

func (o *ListVectorStoreFilesOptions) updateAfter(after *string) { o.After = after }
func (r ListVectorStoreFilesResponse) lastID() *string           { return r.LastID }
func (r ListVectorStoreFilesResponse) hasMore() bool             { return *r.HasMore }

// NewListVectorStoreFileBatchFilesPager returns a pager for vector store files in a batch.
func (c *Client) NewListVectorStoreFileBatchFilesPager(vectorStoreID string, batchID string, options *ListVectorStoreFileBatchFilesOptions) *runtime.Pager[ListVectorStoreFileBatchFilesResponse] {
	nextPageFn := func(ctx context.Context, opts *ListVectorStoreFileBatchFilesOptions) (ListVectorStoreFileBatchFilesResponse, error) {
		return c.internalListVectorStoreFileBatchFiles(ctx, vectorStoreID, batchID, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

func (o *ListVectorStoreFileBatchFilesOptions) updateAfter(after *string) { o.After = after }
func (r ListVectorStoreFileBatchFilesResponse) lastID() *string           { return r.LastID }
func (r ListVectorStoreFileBatchFilesResponse) hasMore() bool             { return *r.HasMore }

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
