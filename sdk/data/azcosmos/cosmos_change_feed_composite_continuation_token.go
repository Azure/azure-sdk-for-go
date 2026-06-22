// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// Version 1 is the initial version of the composite continuation token.
const cosmosCompositeContinuationTokenVersion = 1

type compositeContinuationToken struct {
	// Version is the version of the continuation token format.
	Version int `json:"version,omitempty"`
	// ResourceID is the ID of the resource for which the continuation token is valid.
	ResourceID string `json:"resourceId"`
	// Continuation is the list of Epk Ranges part of the continuation token
	Continuation []changeFeedRange `json:"continuation"`
}

// newCompositeContinuationToken creates a new CompositeContinuationToken with the specified resource ID and continuation ranges.
// This function is used to create a continuation token for the Cosmos DB change feed.
// It is designed for internal use only and should not be used directly by clients.
func newCompositeContinuationToken(resourceID string, continuation []changeFeedRange) compositeContinuationToken {
	return compositeContinuationToken{
		Version:      cosmosCompositeContinuationTokenVersion,
		ResourceID:   resourceID,
		Continuation: continuation,
	}
}

// head returns a pointer to the head queue entry, or nil if the queue is empty.
// Callers MUST NOT mutate the returned entry; use replaceHeadWithChildren or
// advance instead, which preserve queue invariants.
func (t *compositeContinuationToken) head() *changeFeedRange {
	if t == nil || len(t.Continuation) == 0 {
		return nil
	}
	return &t.Continuation[0]
}

// advance rotates the head entry to the tail of the queue, updating its
// ContinuationToken to the freshly-returned ETag from the just-completed
// request. This is the FIFO rotation used after every successful 200 (or
// 304 — both progress the per-range ETag).
//
// No-op if the queue is empty. If newETag is empty, the head's existing
// ContinuationToken is preserved (the service didn't issue a new ETag,
// e.g., on a 304 without one).
func (t *compositeContinuationToken) advance(newETag azcore.ETag) {
	if t == nil || len(t.Continuation) == 0 {
		return
	}
	head := t.Continuation[0]
	if newETag != "" {
		etagCopy := newETag
		head.ContinuationToken = &etagCopy
	}

	// Rotate: drop head, append to tail.
	t.Continuation = append(t.Continuation[1:], head)
}

// replaceHeadWithChildren replaces the head queue entry with the provided
// child entries, preserving the order of the children at the front of the
// queue. Used when a 410/Gone refresh reveals a split (parent → N children)
// or when initial overlap resolution returns multiple children for one
// customer-supplied FeedRange.
//
// No-op if children is empty (degenerate split that produces no overlap).
// No-op if the queue is empty (caller should not have called this).
func (t *compositeContinuationToken) replaceHeadWithChildren(children []changeFeedRange) {
	if t == nil || len(t.Continuation) == 0 || len(children) == 0 {
		return
	}
	rest := t.Continuation[1:]
	t.Continuation = append(append(make([]changeFeedRange, 0, len(children)+len(rest)), children...), rest...)
}
