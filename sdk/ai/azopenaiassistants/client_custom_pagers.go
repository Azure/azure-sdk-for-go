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
	nextPageFn := func(client *Client, ctx context.Context, opts *ListMessagesOptions) (ListMessagesResponse, error) {
		return c.internalListMessages(ctx, threadID, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

// NewListAssistantsPager returns a pager for assistants.
func (c *Client) NewListAssistantsPager(options *ListAssistantsOptions) *runtime.Pager[ListAssistantsResponse] {
	nextPageFn := func(client *Client, ctx context.Context, opts *ListAssistantsOptions) (ListAssistantsResponse, error) {
		return c.internalListAssistants(ctx, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

// NewListAssistantFilesPager returns a pager for an assistant's files.
func (c *Client) NewListAssistantFilesPager(assistantID string, options *ListAssistantFilesOptions) *runtime.Pager[ListAssistantFilesResponse] {
	nextPageFn := func(client *Client, ctx context.Context, opts *ListAssistantFilesOptions) (ListAssistantFilesResponse, error) {
		return c.internalListAssistantFiles(ctx, assistantID, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

// NewListMessageFilesPager returns a pager for a message's files.
func (c *Client) NewListMessageFilesPager(threadID string, messageID string, options *ListMessageFilesOptions) *runtime.Pager[ListMessageFilesResponse] {
	nextPageFn := func(client *Client, ctx context.Context, opts *ListMessageFilesOptions) (ListMessageFilesResponse, error) {
		return c.internalListMessageFiles(ctx, threadID, messageID, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

// NewListRunStepsPager returns a pager for a Run's steps.
func (c *Client) NewListRunStepsPager(threadID string, runID string, options *ListRunStepsOptions) *runtime.Pager[ListRunStepsResponse] {
	nextPageFn := func(client *Client, ctx context.Context, opts *ListRunStepsOptions) (ListRunStepsResponse, error) {
		return c.internalListRunSteps(ctx, threadID, runID, opts)
	}

	return newOpenAIPager(c, nextPageFn, options)
}

// NewListRunsPager returns a pager for a Thread's runs.
func (c *Client) NewListRunsPager(threadID string, options *ListRunsOptions) *runtime.Pager[ListRunsResponse] {
	nextPageFn := func(client *Client, ctx context.Context, opts *ListRunsOptions) (ListRunsResponse, error) {
		return c.internalListRuns(ctx, threadID, opts)
	}
	return newOpenAIPager(c, nextPageFn, options)
}

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
	nextPageFn func(client *Client, ctx context.Context, opts POptionsT) (ResponseT, error),
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

			resp, err := nextPageFn(client, ctx, newOptions)

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

func (r ListAssistantsResponse) lastID() *string     { return r.LastID }
func (r ListMessagesResponse) lastID() *string       { return r.LastID }
func (r ListAssistantFilesResponse) lastID() *string { return r.LastID }
func (r ListMessageFilesResponse) lastID() *string   { return r.LastID }
func (r ListRunStepsResponse) lastID() *string       { return r.LastID }
func (r ListRunsResponse) lastID() *string           { return r.LastID }

func (r ListAssistantsResponse) hasMore() bool     { return *r.HasMore }
func (r ListMessagesResponse) hasMore() bool       { return *r.HasMore }
func (r ListAssistantFilesResponse) hasMore() bool { return *r.HasMore }
func (r ListMessageFilesResponse) hasMore() bool   { return *r.HasMore }
func (r ListRunStepsResponse) hasMore() bool       { return *r.HasMore }
func (r ListRunsResponse) hasMore() bool           { return *r.HasMore }

func (o *ListAssistantsOptions) updateAfter(after *string)     { o.After = after }
func (o *ListMessagesOptions) updateAfter(after *string)       { o.After = after }
func (o *ListAssistantFilesOptions) updateAfter(after *string) { o.After = after }
func (o *ListMessageFilesOptions) updateAfter(after *string)   { o.After = after }
func (o *ListRunStepsOptions) updateAfter(after *string)       { o.After = after }
func (o *ListRunsOptions) updateAfter(after *string)           { o.After = after }
