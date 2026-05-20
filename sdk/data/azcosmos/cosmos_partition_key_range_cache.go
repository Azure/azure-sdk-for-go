// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// partitionKeyRangeCache provides a client-level cache of partition key ranges
// for containers. It is keyed by container ResourceID (RID) and uses event-driven
// invalidation (no TTL). Refreshes are incremental using the change-feed ETag.
// Keying by RID (rather than name-based link) ensures the cache survives
// container renames and matches the service's partition key range addressing.
//
// Concurrency model (single-pending-I/O per container):
//
//   - At most one refresh runs per container at any time. Concurrent callers
//     that arrive while a refresh is in flight share its result.
//   - The cached routing map remains readable while a refresh is in flight;
//     getRoutingMap returns immediately whenever the entry already has a
//     non-nil routingMap and does not wait for the in-flight refresh.
//   - The refresh goroutine runs on context.Background() so a single caller's
//     ctx cancellation does not abort the shared fetch for other waiters.
//     Each waiter still honors its own ctx and returns ctx.Err() if it fires
//     before the shared refresh completes.
//   - forceRefresh accepts the routing map pointer the caller observed when
//     it decided to refresh ("previous"). If the entry already holds a
//     different (fresher) routing map, the caller is served that map
//     immediately without starting a new refresh — i.e. a refresh triggered
//     by a stale-view caller is suppressed (pointer-identity dedup).
type partitionKeyRangeCache struct {
	mu      sync.RWMutex
	entries map[string]*pkRangeCacheEntry // keyed by container ResourceID
}

// refreshOp represents an in-flight partition-key-range refresh for one
// container. Awaiters receive the (rm, err) pair by reading the fields after
// done is closed.
type refreshOp struct {
	done chan struct{}
	rm   *collectionRoutingMap
	err  error
}

type pkRangeCacheEntry struct {
	mu         sync.Mutex // protects routingMap and inFlight
	routingMap *collectionRoutingMap
	inFlight   *refreshOp
}

func newPartitionKeyRangeCache() *partitionKeyRangeCache {
	return &partitionKeyRangeCache{
		entries: make(map[string]*pkRangeCacheEntry),
	}
}

// getOrCreateEntry returns the entry for the given container RID, creating it
// under the write lock if necessary. The returned entry is safe to use after
// the cache-level lock is released; each entry has its own per-entry mutex
// guarding its routingMap/inFlight fields.
func (c *partitionKeyRangeCache) getOrCreateEntry(containerRID string) *pkRangeCacheEntry {
	c.mu.RLock()
	entry, exists := c.entries[containerRID]
	c.mu.RUnlock()
	if exists {
		return entry
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, exists = c.entries[containerRID]; exists {
		return entry
	}
	entry = &pkRangeCacheEntry{}
	c.entries[containerRID] = entry
	return entry
}

// getRoutingMap returns the cached routing map for the given container RID.
// If a routing map is already cached it is returned immediately, even when a
// refresh is in flight. Otherwise the caller joins or starts a shared
// in-flight refresh and waits on its own context.
func (c *partitionKeyRangeCache) getRoutingMap(
	ctx context.Context,
	containerRID string,
	containerLink string,
	client *Client,
) (*collectionRoutingMap, error) {
	entry := c.getOrCreateEntry(containerRID)

	entry.mu.Lock()
	if entry.routingMap != nil {
		rm := entry.routingMap
		entry.mu.Unlock()
		return rm, nil
	}
	op := c.ensureInFlightLocked(entry, containerLink, client)
	entry.mu.Unlock()

	return awaitRefresh(ctx, op)
}

// forceRefresh starts (or joins) a refresh for the given container and
// returns the resulting routing map. The previous parameter is the routing
// map pointer the caller observed when it decided to refresh: when non-nil
// and the entry already holds a different routing map, the caller is served
// that fresher map immediately without starting a new refresh. Pass nil to
// always start or join a refresh (e.g. when the caller has no prior view).
func (c *partitionKeyRangeCache) forceRefresh(
	ctx context.Context,
	containerRID string,
	containerLink string,
	client *Client,
	previous *collectionRoutingMap,
) (*collectionRoutingMap, error) {
	entry := c.getOrCreateEntry(containerRID)

	entry.mu.Lock()
	// Suppress refresh if another caller already installed a fresher map past
	// our view. Dedup of stale-view triggers via pointer identity.
	if previous != nil && entry.routingMap != nil && entry.routingMap != previous {
		rm := entry.routingMap
		entry.mu.Unlock()
		return rm, nil
	}
	op := c.ensureInFlightLocked(entry, containerLink, client)
	entry.mu.Unlock()

	return awaitRefresh(ctx, op)
}

// ensureInFlightLocked returns the entry's in-flight refresh op, creating
// (and spawning) one if none is already running. Caller MUST hold entry.mu.
func (c *partitionKeyRangeCache) ensureInFlightLocked(
	entry *pkRangeCacheEntry,
	containerLink string,
	client *Client,
) *refreshOp {
	if entry.inFlight != nil {
		return entry.inFlight
	}
	op := &refreshOp{done: make(chan struct{})}
	entry.inFlight = op
	go c.runRefresh(entry, containerLink, client, op)
	return op
}

// runRefresh executes the change-feed refresh on a detached context.Background()
// so caller cancellations do not abort the shared fetch. On completion it
// updates the entry under entry.mu, clears the in-flight slot, and signals
// awaiters by closing op.done.
func (c *partitionKeyRangeCache) runRefresh(
	entry *pkRangeCacheEntry,
	containerLink string,
	client *Client,
	op *refreshOp,
) {
	rm, err := c.refreshEntryDetached(containerLink, entry, client)

	entry.mu.Lock()
	if err == nil && rm != nil {
		entry.routingMap = rm
	}
	op.rm = rm
	op.err = err
	entry.inFlight = nil
	entry.mu.Unlock()

	close(op.done)
}

// awaitRefresh blocks the caller until either the refresh completes or the
// caller's context is cancelled. The refresh continues running in the
// background even when individual awaiters return early via ctx.Err().
func awaitRefresh(ctx context.Context, op *refreshOp) (*collectionRoutingMap, error) {
	select {
	case <-op.done:
		return op.rm, op.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// invalidate clears the cached routing map for the given container RID,
// forcing the next access to fetch fresh data. An in-flight refresh continues
// and its result will replace the cleared map.
func (c *partitionKeyRangeCache) invalidate(containerRID string) {
	c.mu.RLock()
	entry, exists := c.entries[containerRID]
	c.mu.RUnlock()

	if exists {
		entry.mu.Lock()
		entry.routingMap = nil
		entry.mu.Unlock()
	}
}

// maxChangeFeedIterations caps the number of change-feed fetch loops
// to prevent runaway requests during large-scale splits.
const maxChangeFeedIterations = 1000

// changeFeedPageMaxAttempts bounds how many times we retry a single
// change-feed page after a transient failure. The pipeline already does
// per-request retry; this is an additional safety net that preserves
// already-accumulated pages instead of restarting the entire change-feed
// drain from page 1 on a single bad page.
const changeFeedPageMaxAttempts = 3

// changeFeedPageRetryBaseDelay is the base sleep before retrying a failed
// change-feed page. Backoff is linear: base, 2*base, 3*base, ...
const changeFeedPageRetryBaseDelay = 100 * time.Millisecond

// changeFeedResult holds the result of draining all change-feed pages.
type changeFeedResult struct {
	ranges    []partitionKeyRange
	finalETag string
	completed bool // true only if loop terminated via 304 Not Modified
}

// fetchAllChangeFeedPages fetches change-feed pages starting from startETag
// until 304 Not Modified or the iteration cap. Returns the accumulated ranges,
// the final ETag, and whether the loop completed cleanly (304 received).
//
// Each individual page is retried up to changeFeedPageMaxAttempts times on
// transient errors (5xx, 408, 429, network errors) with linear backoff so
// a single bad page doesn't discard the pages already accumulated. Non-
// transient errors (4xx other than 408/429, context errors) fail fast.
func fetchAllChangeFeedPages(
	ctx context.Context,
	containerLink string,
	startETag string,
	client *Client,
) (changeFeedResult, error) {
	var allRanges []partitionKeyRange
	currentETag := startETag
	for i := 0; i < maxChangeFeedIterations; i++ {
		if err := ctx.Err(); err != nil {
			return changeFeedResult{}, err
		}
		result, err := fetchOneChangeFeedPageWithRetry(ctx, containerLink, currentETag, client)
		if err != nil {
			// Even though we're surfacing the error to the caller, log how
			// far we got so operators can correlate partial-drain failures
			// with the next refresh re-starting from scratch.
			log.Writef(azlog.EventResponse, "partition key range change-feed page failed for container %s after %d successful pages (%d ranges accumulated): %v", containerLink, i, len(allRanges), err)
			return changeFeedResult{}, err
		}

		if result.notModified {
			// 304 Not Modified — change feed fully drained
			if result.etag != "" {
				currentETag = result.etag
			}
			return changeFeedResult{ranges: allRanges, finalETag: currentETag, completed: true}, nil
		}
		allRanges = append(allRanges, result.ranges...)
		if result.etag != "" {
			currentETag = result.etag
		}
	}
	// Loop cap reached without 304
	log.Writef(azlog.EventResponse, "partition key range change-feed loop exited without reaching 304 Not Modified after %d iterations for container %s (accumulated %d ranges)", maxChangeFeedIterations, containerLink, len(allRanges))
	return changeFeedResult{ranges: allRanges, finalETag: currentETag, completed: false}, nil
}

// fetchOneChangeFeedPageWithRetry fetches a single change-feed page,
// retrying on transient errors so a transient hiccup mid-pagination
// doesn't discard the already-accumulated pages in the caller. Returns
// the last error if all attempts fail or the caller's context fires.
func fetchOneChangeFeedPageWithRetry(
	ctx context.Context,
	containerLink string,
	currentETag string,
	client *Client,
) (fetchPartitionKeyRangesResult, error) {
	var lastErr error
	for attempt := 0; attempt < changeFeedPageMaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fetchPartitionKeyRangesResult{}, err
		}
		result, err := fetchPartitionKeyRanges(ctx, containerLink, currentETag, client)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !isTransientPKRangeFetchError(err) {
			return fetchPartitionKeyRangesResult{}, err
		}
		// Last attempt — don't sleep, just return.
		if attempt == changeFeedPageMaxAttempts-1 {
			break
		}
		// Linear backoff: 1×base, 2×base, ...
		delay := time.Duration(attempt+1) * changeFeedPageRetryBaseDelay
		log.Writef(azlog.EventResponse, "partition key range change-feed page transient failure for container %s (attempt %d/%d, retrying in %s): %v", containerLink, attempt+1, changeFeedPageMaxAttempts, delay, err)
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return fetchPartitionKeyRangesResult{}, ctx.Err()
		}
	}
	return fetchPartitionKeyRangesResult{}, lastErr
}

// isTransientPKRangeFetchError reports whether a /pkranges fetch error is
// worth retrying mid-pagination. Returns true for 5xx, 408, 429, and
// network-class errors (any error without an HTTP response). Returns false
// for 4xx (other than 408/429) and for context cancellation.
func isTransientPKRangeFetchError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		// Network / transport error — treat as transient.
		return true
	}
	switch {
	case respErr.StatusCode >= 500:
		return true
	case respErr.StatusCode == http.StatusRequestTimeout: // 408
		return true
	case respErr.StatusCode == http.StatusTooManyRequests: // 429
		return true
	}
	return false
}

// refreshEntryDetached fetches PK ranges from the service and returns a fresh
// routing map. It snapshots the entry's previous routing map under entry.mu
// (briefly), then performs all network I/O without holding any lock. The
// caller is responsible for installing the returned map onto the entry.
//
// The function uses context.Background() internally via the runRefresh
// goroutine; callers must not pass a caller-scoped context.
func (c *partitionKeyRangeCache) refreshEntryDetached(
	containerLink string,
	entry *pkRangeCacheEntry,
	client *Client,
) (*collectionRoutingMap, error) {
	entry.mu.Lock()
	previousMap := entry.routingMap
	entry.mu.Unlock()

	ctx := context.Background()

	if previousMap != nil && previousMap.changeFeedETag != "" {
		// Incremental refresh: accumulate ALL change-feed pages first, then
		// call tryCombine once with the entire batch. This matches .NET behavior
		// and handles cascading splits (A→B+C, B→D+E) that span multiple pages.
		// Note: unlike Python, we do not retry the incremental path on failure —
		// we fall through to full refresh immediately, matching .NET behavior.
		result, err := fetchAllChangeFeedPages(ctx, containerLink, previousMap.changeFeedETag, client)
		if err != nil {
			return nil, err
		}

		if result.completed {
			if len(result.ranges) == 0 {
				// No changes since last refresh — surface the ETag bump if any,
				// otherwise return the previous map unchanged.
				if result.finalETag != previousMap.changeFeedETag {
					return &collectionRoutingMap{
						orderedRanges:  previousMap.orderedRanges,
						rangeByID:      previousMap.rangeByID,
						goneRanges:     previousMap.goneRanges,
						changeFeedETag: result.finalETag,
					}, nil
				}
				return previousMap, nil
			}

			merged := previousMap.tryCombine(result.ranges, result.finalETag)
			if merged != nil {
				return merged, nil
			}
		}

		// Incremental merge failed or loop cap exhausted — fall through to full refresh.
	}

	// Full change-feed refresh: fetch all ranges from the beginning using A-IM
	// without an ETag, looping until 304 Not Modified.
	result, err := fetchAllChangeFeedPages(ctx, containerLink, "", client)
	if err != nil {
		return nil, err
	}

	if !result.completed {
		return nil, fmt.Errorf("partition key range cache refresh failed: change-feed pagination did not terminate after %d iterations for container %s (accumulated %d ranges)", maxChangeFeedIterations, containerLink, len(result.ranges))
	}

	newMap := newCollectionRoutingMap(result.ranges, result.finalETag)
	if !isCompleteSetOfRanges(newMap.orderedRanges) {
		issue := describeRangeDiscontinuity(newMap.orderedRanges)
		return nil, fmt.Errorf("partition key range cache refresh failed: service returned an incomplete set of ranges for container %s (raw ranges=%d, final ranges=%d, issue: %s). This may indicate a transient issue during a partition split", containerLink, len(result.ranges), len(newMap.orderedRanges), issue)
	}

	return newMap, nil
}

// fetchPartitionKeyRangesResult holds the result of a single change-feed fetch.
type fetchPartitionKeyRangesResult struct {
	ranges      []partitionKeyRange
	etag        string
	notModified bool // true only when the service returns 304 Not Modified
}

// fetchPartitionKeyRanges fetches partition key ranges from the service using
// the change-feed mechanism. It always sets A-IM: Incremental feed and
// x-ms-max-item-count: -1 headers, matching the behavior of the .NET and
// Python SDKs. If changeFeedETag is non-empty, it sets If-None-Match for
// incremental updates; otherwise it fetches all ranges from the beginning.
func fetchPartitionKeyRanges(
	ctx context.Context,
	containerLink string,
	changeFeedETag string,
	client *Client,
) (fetchPartitionKeyRangesResult, error) {
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypePartitionKeyRange,
		resourceAddress: containerLink,
	}

	path, err := generatePathForNameBased(resourceTypePartitionKeyRange, operationContext.resourceAddress, true)
	if err != nil {
		return fetchPartitionKeyRangesResult{}, err
	}

	o := &partitionKeyRangeOptions{}

	azResponse, err := client.sendGetRequest(
		path,
		ctx,
		operationContext,
		o,
		func(req *policy.Request) {
			req.Raw().Header.Set(cosmosHeaderChangeFeed, cosmosHeaderValuesChangeFeed)
			req.Raw().Header.Set(cosmosHeaderMaxItemCount, cosmosHeaderValuesMaxItemAll)
			if changeFeedETag != "" {
				req.Raw().Header.Set(headerIfNoneMatch, changeFeedETag)
			}
		})
	if err != nil {
		return fetchPartitionKeyRangesResult{}, err
	}

	newETag := azResponse.Header.Get(cosmosHeaderEtag)

	// 304 Not Modified means no changes
	if azResponse.StatusCode == http.StatusNotModified {
		_ = azResponse.Body.Close()
		return fetchPartitionKeyRangesResult{etag: newETag, notModified: true}, nil
	}

	body, err := azruntime.Payload(azResponse)
	if err != nil {
		return fetchPartitionKeyRangesResult{}, err
	}
	_ = azResponse.Body.Close()

	var response partitionKeyRangeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fetchPartitionKeyRangesResult{}, err
	}

	return fetchPartitionKeyRangesResult{ranges: response.PartitionKeyRanges, etag: newETag}, nil
}
