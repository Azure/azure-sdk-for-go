// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// GetChangeFeed retrieves a single page of the change feed using the provided options.
// ctx - The context for the request.
// options - Options for the operation
//
// Routes via overlap-match against the current PK-range cache (split-aware). When
// the customer's FeedRange straddles a split, it's expanded into one queue entry
// per child so no events are missed at the boundary. 410/Gone responses with a
// PK-range substatus trigger a cache refresh and bounded retry. The continuation
// token is a multi-range composite; subsequent calls drain remaining ranges.
//
// Returns ErrFeedRangeUnresolved (wrapped) when the customer's FeedRange/token
// doesn't overlap any current physical range even after a forced refresh — a
// signal to re-derive FeedRanges from GetFeedRanges.
//
// Returns an error wrapping *azcore.ResponseError on persistent 410/Gone or any
// non-retryable HTTP error.
func (c *ContainerClient) GetChangeFeed(
	ctx context.Context,
	options *ChangeFeedOptions,
) (ChangeFeedResponse, error) {
	if options == nil {
		options = &ChangeFeedOptions{}
	}

	var err error
	spanName, err := c.getSpanForItems(operationTypeRead)
	if err != nil {
		return ChangeFeedResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()

	// Cross-container token guard. We only need the container's current
	// ResourceID when a continuation token carries one to validate. Building
	// the initial queue takes care of that lazily so the no-token path stays
	// at one extra request (pk-ranges fetch).
	token, partitionKeyRanges, err := c.buildChangeFeedInitialQueue(ctx, options)
	if err != nil {
		return ChangeFeedResponse{}, err
	}

	// Capture into err so the deferred endSpan closure observes drain
	// failures (410 budget exhaustion, send errors, refresh failures,
	// serialization errors, ErrFeedRangeUnresolved) instead of recording
	// success on every drain that actually failed.
	var resp ChangeFeedResponse
	resp, err = c.getChangeFeedForQueue(ctx, options, token, partitionKeyRanges)
	return resp, err
}

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
			if compositeToken.ResourceID != "" {
				currentRID, ridErr := c.getContainerRID(ctx)
				if ridErr != nil {
					return nil, nil, ridErr
				}
				if currentRID != "" && compositeToken.ResourceID != currentRID {
					return nil, nil, fmt.Errorf(
						"continuation token ResourceID %q does not match the current container's ResourceID %q; the token was issued for a different container",
						compositeToken.ResourceID, currentRID,
					)
				}
			}
			queue := append([]changeFeedRange(nil), compositeToken.Continuation...)
			token := compositeToken
			token.Continuation = queue

			// Fetch a PK-range snapshot for the drain loop. Reused so the
			// loop doesn't issue an extra request per iteration.
			pkrResp, err := c.getPartitionKeyRanges(ctx, nil)
			if err != nil {
				return nil, nil, err
			}
			return &token, pkrResp.PartitionKeyRanges, nil
		}
		// Not a composite token — fall through to FeedRange path; the legacy
		// ETag-only continuation is handled by buildRequestHeaders via the
		// queue head's ContinuationToken field.
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
	token := compositeContinuationToken{
		Version:      cosmosCompositeContinuationTokenVersion,
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
//
// Degraded fallback: when the PK-range cache is unavailable AND the direct
// fetch returns no ranges (e.g., test scaffolding), routing information is
// effectively missing. Rather than failing loudly, we drop any continuation
// context the caller may have provided, log a warning so the misroute is
// observable, and return a passthrough entry representing the customer's
// FeedRange (placed both in the children list and in the snapshot). The
// drain loop then issues the request fresh — without a PK-range-id header
// and without an If-None-Match ETag — and the server resolves routing.
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
		// Cache is fresh and still no overlap → unresolvable.
		return nil, nil, &feedRangeUnresolvedError{feedRange: feedRange}
	}

	// No cache wired up. If we got any PK ranges in the direct fetch but none
	// overlap, that's a real customer mistake — same loud-fail semantics as
	// the cache path.
	if len(pkrResp.PartitionKeyRanges) > 0 {
		return nil, nil, &feedRangeUnresolvedError{feedRange: feedRange}
	}

	// Degraded fallback (cache absent AND direct fetch returned no ranges).
	// Log a warning, then issue the read fresh: passthrough range with no
	// continuation token. The passthrough is also placed in the snapshot so
	// the drain loop's overlap check is satisfied without re-entering this
	// branch.
	log.Writef(azlog.EventResponse,
		"azcosmos: GetChangeFeed: routing information unavailable for FeedRange [%s, %s); reading fresh without continuation token or PK-range header",
		feedRange.MinInclusive, feedRange.MaxExclusive,
	)
	passthrough := partitionKeyRange{
		MinInclusive: feedRange.MinInclusive,
		MaxExclusive: feedRange.MaxExclusive,
	}
	return []partitionKeyRange{passthrough}, []partitionKeyRange{passthrough}, nil
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

// getChangeFeedForQueue drains the queue, advancing on every response (200 or
// 304). On 200 with documents, returns immediately so the caller can process
// the page; on 304, rotates and tries the next entry until the original queue
// length is fully consumed (with budget bumps on splits). On 410, refreshes
// the cache, re-resolves the head, and retries — capped at maxPKRangeGoneRetries.
//
// partitionKeyRanges is the snapshot fetched once at the start of the call;
// the loop reuses it instead of re-fetching per iteration. The 410-retry path
// re-fetches and replaces the snapshot.
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

	for rotations < originalQueueLen {
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
			// still no overlap, the head is unresolvable. The cache-nil branch
			// is reachable only from test scaffolding (production always wires
			// up caches via acquireCaches) but we still attempt a direct fetch
			// fallback so handcrafted tests behave the same.
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
				// Degraded fallback: cache absent AND the direct fetch
				// returned nothing — routing information is unavailable for
				// this head. Drop the head's continuation token (an ETag
				// from a prior call could correspond to a now-defunct
				// physical range), log a warning so the misroute is
				// observable, and synthesize a passthrough so the request
				// is issued fresh — without a PK-range-id header and
				// without an If-None-Match ETag.
				if c.database.client.getPKRangeCache() == nil && len(partitionKeyRanges) == 0 {
					log.Writef(azlog.EventResponse,
						"azcosmos: GetChangeFeed: routing information unavailable for head [%s, %s); dropping continuation token and reading fresh",
						head.MinInclusive, head.MaxExclusive,
					)
					token.dropHeadContinuation()
					// head is a pointer into token.Continuation[0]; the
					// in-place mutation by dropHeadContinuation is already
					// visible through it. No reload needed.
					headFeedRange = FeedRange{MinInclusive: head.MinInclusive, MaxExclusive: head.MaxExclusive}
					passthrough := partitionKeyRange{
						MinInclusive: head.MinInclusive,
						MaxExclusive: head.MaxExclusive,
					}
					partitionKeyRanges = []partitionKeyRange{passthrough}
					overlaps = []partitionKeyRange{passthrough}
				} else {
					return ChangeFeedResponse{}, &feedRangeUnresolvedError{feedRange: headFeedRange}
				}
			}
		}

		var resolvedPKRangeID string
		if len(overlaps) > 1 {
			// Split-expansion. Replace the head with N children inheriting
			// the head's ETag, and bump the rotation budget so newly-inserted
			// children get visited in this call. Reset the 410 budget too —
			// each newly-inserted child is a fresh physical head and deserves
			// its own retry allowance.
			children := buildChildQueueEntries(overlaps, head.ContinuationToken)
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

		serialized, serErr := serializeCompositeContinuationToken(token)
		if serErr != nil {
			return response, serErr
		}
		response.ContinuationToken = serialized
		lastResp = response

		// 200 with documents → return immediately so the caller can process.
		if response.RawResponse != nil && response.RawResponse.StatusCode == http.StatusOK && response.Count > 0 {
			return response, nil
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
	return lastResp, nil
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
