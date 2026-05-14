// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal/epk"
)

// buildChangeFeedInitialQueue assembles the queue this call's drain loop will
// operate on, validating any provided continuation token against the current
// container and resolving any provided FeedRange against the current PK
// range cache. Returns the fetched PK-range snapshot alongside the token so
// the drain loop can reuse it without re-fetching on every iteration; the
// 410-retry path re-fetches on its own.
//
// Returns (token, snapshot, nil) on success. Returns ErrFeedRangeUnresolved
// (wrapped) when no overlap exists even after the cache is fresh.
func (c *ContainerClient) buildChangeFeedInitialQueue(
	ctx context.Context,
	options *ChangeFeedOptions,
) (*compositeContinuationToken, []partitionKeyRange, error) {
	// Path A: continuation token drives the queue.
	if options.Continuation != nil && *options.Continuation != "" {
		var compositeToken compositeContinuationToken
		if err := json.Unmarshal([]byte(*options.Continuation), &compositeToken); err == nil && len(compositeToken.Continuation) > 0 {
			// Reject cross-container token reuse loudly. Customers who hit this
			// have either pasted the wrong token, dropped a container and
			// recreated it under the same name, or fanned out a token to a
			// different client. Continuing would route against the wrong map.
			currentRID, ridErr := c.getContainerRID(ctx)
			if ridErr != nil {
				return nil, nil, ridErr
			}
			if compositeToken.ResourceID != "" && currentRID != "" && compositeToken.ResourceID != currentRID {
				return nil, nil, fmt.Errorf(
					"continuation token ResourceID %q does not match the current container's ResourceID %q; the token was issued for a different container",
					compositeToken.ResourceID, currentRID,
				)
			}
			queue := append([]changeFeedRange(nil), compositeToken.Continuation...)
			token := compositeToken
			token.Continuation = queue
			// Populate ResourceID at construction time so the cross-container
			// guard remains meaningful even if the very first response is a
			// 304 (no body parsed → response.ResourceID empty).
			if token.ResourceID == "" {
				token.ResourceID = currentRID
			}

			// Fetch a PK-range snapshot for the drain loop. Reused so the
			// loop doesn't issue an extra request per iteration.
			pkrResp, err := c.getPartitionKeyRanges(ctx, nil)
			if err != nil {
				return nil, nil, err
			}
			return &token, pkrResp.PartitionKeyRanges, nil
		}
		// options.Continuation was supplied but is not a multi-range composite
		// token. Only composite tokens issued by this SDK's GetChangeFeed are
		// honored on resume; legacy raw-ETag continuation strings are NOT
		// supported and are silently dropped here. The FeedRange path below
		// will start fresh.
	}

	// Path B: FeedRange drives the queue.
	if options.FeedRange == nil {
		return nil, nil, fmt.Errorf("GetChangeFeed requires a FeedRange to be set in the options, or a continuation token that contains a composite continuation token")
	}

	children, pkrs, err := c.resolveFeedRangeToChildren(ctx, *options.FeedRange)
	if err != nil {
		return nil, nil, err
	}
	entries := buildChildQueueEntries(children, nil)
	// Populate ResourceID at construction time so a token persisted from a
	// 304-only first call still triggers the cross-container guard on resume.
	currentRID, ridErr := c.getContainerRID(ctx)
	if ridErr != nil {
		return nil, nil, ridErr
	}
	token := compositeContinuationToken{
		Version:      cosmosCompositeContinuationTokenVersion,
		ResourceID:   currentRID,
		Continuation: entries,
	}
	return &token, pkrs, nil
}

// resolveFeedRangeToChildren returns the routing-map ranges that overlap the
// given customer-supplied FeedRange. On no-overlap, performs a single forced
// refresh and retries; on still no overlap, returns ErrFeedRangeUnresolved.
//
// Also returns the PK-range snapshot it fetched, so the drain loop can reuse
// it for the rest of this GetChangeFeed call.
func (c *ContainerClient) resolveFeedRangeToChildren(
	ctx context.Context,
	feedRange FeedRange,
) ([]partitionKeyRange, []partitionKeyRange, error) {
	pkrResp, err := c.getPartitionKeyRanges(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	overlaps := overlappingPartitionKeyRanges(feedRange, pkrResp.PartitionKeyRanges)
	if len(overlaps) > 0 {
		return overlaps, pkrResp.PartitionKeyRanges, nil
	}

	// No overlap on the cached map. Try a forced refresh once if a cache
	// exists; if still no overlap, the customer's FeedRange genuinely doesn't
	// apply to this container.
	if c.database.client.getPKRangeCache() != nil {
		if refreshErr := c.refreshPKRangeCache(ctx); refreshErr != nil {
			return nil, nil, refreshErr
		}
		pkrResp, err = c.getPartitionKeyRanges(ctx, nil)
		if err != nil {
			return nil, nil, err
		}
		overlaps = overlappingPartitionKeyRanges(feedRange, pkrResp.PartitionKeyRanges)
		if len(overlaps) > 0 {
			return overlaps, pkrResp.PartitionKeyRanges, nil
		}
	}

	return nil, nil, &feedRangeUnresolvedError{feedRange: feedRange}
}

// buildChildQueueEntries materializes [].changeFeedRange entries for each
// child range, copying the inheritETag pointer onto every child so no events
// are skipped at the split boundary. inheritETag may be nil for fresh ranges
// that have never been read.
func buildChildQueueEntries(children []partitionKeyRange, inheritETag *azcore.ETag) []changeFeedRange {
	out := make([]changeFeedRange, 0, len(children))
	for _, ch := range children {
		entry := changeFeedRange{
			MinInclusive: ch.MinInclusive,
			MaxExclusive: ch.MaxExclusive,
		}
		if inheritETag != nil {
			etagCopy := *inheritETag
			entry.ContinuationToken = &etagCopy
		}
		out = append(out, entry)
	}
	return out
}

// clampChildrenToParent narrows each child's [Min, Max) to the intersection
// with the parent head's [Min, Max), preserving the customer's original sub-
// range intent across split-expansion. Children that fall entirely outside
// the parent are dropped. Mirrors Java's
// FeedRangeCompositeContinuationImpl.createChildRanges clamping behavior.
func clampChildrenToParent(children []partitionKeyRange, parent changeFeedRange) []partitionKeyRange {
	parentMin := parent.MinInclusive
	parentMax := normalizeMaxBoundary(parent.MaxExclusive)
	out := make([]partitionKeyRange, 0, len(children))
	for _, ch := range children {
		childMin := ch.MinInclusive
		childMax := normalizeMaxBoundary(ch.MaxExclusive)
		// Intersection: max of mins, min of maxes.
		newMin := childMin
		if epk.CompareEPK(parentMin, childMin) > 0 {
			newMin = parentMin
		}
		newMax := childMax
		if epk.CompareEPK(parentMax, childMax) < 0 {
			newMax = parentMax
		}
		// Discard children whose intersection with the parent is empty.
		if epk.CompareEPK(newMin, newMax) >= 0 {
			continue
		}
		clamped := ch
		clamped.MinInclusive = newMin
		clamped.MaxExclusive = newMax
		out = append(out, clamped)
	}
	return out
}

// getChangeFeedForQueue drains the queue, advancing on every response (200 or
// 304). On 200 with documents, returns immediately so the caller can process
// the page; on 304, rotates and tries the next entry until the original queue
// length is fully consumed (with budget bumps on splits). On 410, refreshes
// the cache, re-resolves the head, and retries — capped at maxPKRangeGoneRetries.
//
// partitionKeyRanges is the snapshot fetched once at the start of the call;
// the loop reuses it instead of re-fetching per iteration. The 410-retry path
// re-fetches and replaces the snapshot.
//
// RequestCharge is aggregated across all iterations so the returned response
// reports total RU consumed by the call (matching cosmos_container_read_many).
//
// On a 410-budget exhaustion mid-drain, the partial response (with the queue
// state rotated past every successfully-drained head) is returned alongside
// the error so callers can resume from where the drain failed instead of
// re-querying already-drained heads.
func (c *ContainerClient) getChangeFeedForQueue(
	ctx context.Context,
	options *ChangeFeedOptions,
	token *compositeContinuationToken,
	partitionKeyRanges []partitionKeyRange,
) (ChangeFeedResponse, error) {
	if token == nil || len(token.Continuation) == 0 {
		return ChangeFeedResponse{}, fmt.Errorf("GetChangeFeed has nothing to drain: no FeedRange and no continuation token entries")
	}

	// Drain budget: how many rotations we'll perform before we give up and
	// return an empty page so the caller can poll again. Starts at the queue
	// length and grows whenever a split-expansion inserts children.
	originalQueueLen := len(token.Continuation)
	rotations := 0
	pkRangeGoneAttempts := 0

	var lastResp ChangeFeedResponse
	var totalRequestCharge float32

	// finalize attaches the rotated continuation token and the aggregated
	// request charge to the given response. Used both on success and on
	// the 410-budget-exhausted partial-state return path.
	finalize := func(r ChangeFeedResponse) (ChangeFeedResponse, error) {
		serialized, serErr := serializeCompositeContinuationToken(token)
		if serErr != nil {
			return r, serErr
		}
		r.ContinuationToken = serialized
		r.RequestCharge = totalRequestCharge
		return r, nil
	}

	for rotations < originalQueueLen {
		// Honor caller cancellation between sub-requests so a long drain
		// doesn't keep working past a deadline / explicit cancel.
		select {
		case <-ctx.Done():
			return ChangeFeedResponse{}, ctx.Err()
		default:
		}

		head := token.head()
		if head == nil {
			break
		}

		// Resolve the head's EPK range to a single PK-range ID against the
		// current routing-map snapshot.
		headFeedRange := FeedRange{MinInclusive: head.MinInclusive, MaxExclusive: head.MaxExclusive}
		overlaps := overlappingPartitionKeyRanges(headFeedRange, partitionKeyRanges)
		if len(overlaps) == 0 {
			// No overlap on the cached map. Force a refresh and re-fetch; if
			// still no overlap, the head is unresolvable.
			if c.database.client.getPKRangeCache() != nil {
				if refreshErr := c.refreshPKRangeCache(ctx); refreshErr != nil {
					return ChangeFeedResponse{}, refreshErr
				}
			}
			pkrResp, err := c.getPartitionKeyRanges(ctx, nil)
			if err != nil {
				return ChangeFeedResponse{}, err
			}
			partitionKeyRanges = pkrResp.PartitionKeyRanges
			overlaps = overlappingPartitionKeyRanges(headFeedRange, partitionKeyRanges)
			if len(overlaps) == 0 {
				return ChangeFeedResponse{}, &feedRangeUnresolvedError{feedRange: headFeedRange}
			}
		}

		var resolvedPKRangeID string
		if len(overlaps) > 1 {
			// Split-expansion. Replace the head with N children inheriting
			// the head's ETag, and bump the rotation budget so newly-inserted
			// children get visited in this call. Reset the 410 budget too —
			// each newly-inserted child is a fresh physical head and deserves
			// its own retry allowance. Children are clamped to the parent
			// head's bounds to preserve the customer's original sub-range
			// intent (matches Java's createChildRanges behavior).
			clamped := clampChildrenToParent(overlaps, *head)
			if len(clamped) == 0 {
				// Defensive: every overlap fell outside the parent. Treat as
				// unresolvable rather than spin in an empty loop.
				return ChangeFeedResponse{}, &feedRangeUnresolvedError{feedRange: headFeedRange}
			}
			children := buildChildQueueEntries(clamped, head.ContinuationToken)
			token.replaceHeadWithChildren(children)
			originalQueueLen += len(children) - 1
			pkRangeGoneAttempts = 0
			continue
		}
		resolvedPKRangeID = overlaps[0].ID

		headers, headerErr := options.buildRequestHeaders(*head, resolvedPKRangeID)
		if headerErr != nil {
			return ChangeFeedResponse{}, headerErr
		}

		addHeaders := func(r *policy.Request) {
			for k, v := range headers {
				r.Raw().Header.Set(k, v)
			}
		}

		operationContext := pipelineRequestOptions{
			resourceType:    resourceTypeDocument,
			resourceAddress: c.link,
			headerOptionsOverride: &headerOptionsOverride{
				priorityLevel:    options.PriorityLevel,
				throughputBucket: options.ThroughputBucket,
			},
		}
		path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
		if err != nil {
			return ChangeFeedResponse{}, err
		}

		azResponse, sendErr := c.database.client.sendGetRequest(
			path, ctx, operationContext, nil, addHeaders,
		)
		if sendErr != nil {
			// 410/Gone with a PK-range substatus → refresh + retry.
			if isPKRangeGoneResponseError(sendErr) {
				if pkRangeGoneAttempts >= maxPKRangeGoneRetries {
					// Surface partial drain progress alongside the error so
					// the caller can resume from the failed head instead of
					// re-querying already-drained heads. lastResp may be
					// zero-value if no head succeeded yet — in that case the
					// caller still sees the original token state.
					if lastResp.RawResponse != nil {
						partial, _ := finalize(lastResp)
						return partial, sendErr
					}
					return ChangeFeedResponse{}, sendErr
				}
				pkRangeGoneAttempts++
				if refreshErr := c.refreshPKRangeCache(ctx); refreshErr != nil {
					return ChangeFeedResponse{}, refreshErr
				}
				// Re-fetch the routing map after the cache was invalidated.
				pkrResp, fetchErr := c.getPartitionKeyRanges(ctx, nil)
				if fetchErr != nil {
					return ChangeFeedResponse{}, fetchErr
				}
				partitionKeyRanges = pkrResp.PartitionKeyRanges
				// Retry the same head against the refreshed snapshot.
				continue
			}
			return ChangeFeedResponse{}, sendErr
		}

		response, err := newChangeFeedResponse(azResponse)
		if err != nil {
			return response, err
		}

		// Aggregate RU charge across every sub-request so the returned
		// response reports the true total cost of the drain (single-iter
		// callers are unaffected; multi-range drains stop under-reporting).
		totalRequestCharge += response.RequestCharge

		// Capture the response body's _rid into the token's ResourceID on first
		// successful response. This keeps the cross-container guard meaningful
		// across resume — token-issued-by-this-container always carries the
		// container's RID — and matches pre-F1 PopulateCompositeContinuationToken
		// semantics that downstream tests rely on.
		if token.ResourceID == "" && response.ResourceID != "" {
			token.ResourceID = response.ResourceID
		}

		// Always rotate the head with the freshly-issued ETag, regardless of
		// status. This preserves drain progress even across 304s.
		newETag := response.ETag
		feedRangeForResp := &FeedRange{MinInclusive: head.MinInclusive, MaxExclusive: head.MaxExclusive}
		token.advance(newETag)
		rotations++
		// Head advanced to a new physical range; reset the 410 budget so the
		// next head gets its own allowance instead of inheriting prior 410s.
		pkRangeGoneAttempts = 0

		response.FeedRange = feedRangeForResp
		lastResp = response

		// 200 with a non-empty page → return immediately so the caller can
		// process. The Documents-length belt covers the (shouldn't-happen)
		// case where the server omits _count but still ships docs.
		if response.RawResponse != nil &&
			response.RawResponse.StatusCode == http.StatusOK &&
			(response.Count > 0 || len(response.Documents) > 0) {
			return finalize(response)
		}

		// 304 (or 200 with zero documents) → keep draining the rest of the queue.
	}

	// Whole queue drained without finding documents. Return the last (empty)
	// response with the rotated continuation token so the caller knows the
	// drain progressed and can poll again later.
	if lastResp.RawResponse == nil {
		// Nothing was issued (queue was empty). Synthesize an empty response.
		return ChangeFeedResponse{}, nil
	}
	return finalize(lastResp)
}

// serializeCompositeContinuationToken marshals the token as JSON for emission
// to the customer. Returns "" if the token is nil or has an empty queue.
func serializeCompositeContinuationToken(token *compositeContinuationToken) (string, error) {
	if token == nil || len(token.Continuation) == 0 {
		return "", nil
	}
	if token.Version == 0 {
		token.Version = cosmosCompositeContinuationTokenVersion
	}
	b, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
