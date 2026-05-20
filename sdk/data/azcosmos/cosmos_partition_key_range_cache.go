// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"syscall"
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
//   - The refresh goroutine runs on a detached context (a child of
//     context.Background() with pkRangeBackgroundRefreshTimeout) so a single
//     caller's ctx cancellation does not abort the shared fetch for other
//     waiters. Each waiter still honors its own ctx and returns ctx.Err() if
//     it fires before the shared refresh completes.
//   - The detached refresh timeout guarantees the inFlight slot is eventually
//     cleared even when the transport hangs, so a wedged fetch cannot wedge
//     the cache slot indefinitely.
//   - invalidate() bumps the entry's generation. A refresh whose result is
//     ready after an invalidate sees a generation mismatch, discards its
//     result from the cache, AND surfaces errPKRangeCacheInvalidatedDuringRefresh
//     to its awaiters so they re-enter and start a fresh post-invalidate op
//     rather than routing against the stale view they were waiting on.
type partitionKeyRangeCache struct {
	mu      sync.RWMutex
	entries map[string]*pkRangeCacheEntry // keyed by container ResourceID
}

// pkRangeBackgroundRefreshTimeout caps how long a detached refresh goroutine
// is allowed to run. Beyond this the goroutine errors out, clears the
// in-flight slot, and the next caller starts a fresh op. Defends against
// hung transports / wedged regions.
//
// Exposed as a var (not a const) so tests can override it without sleeping
// for the production default. The size must be large enough to cover a full
// change-feed drain (up to maxChangeFeedIterations pages) on a container
// undergoing a heavy split storm.
var pkRangeBackgroundRefreshTimeout = 60 * time.Second

// errPKRangeCacheInvalidatedDuringRefresh is returned to awaiters of a
// refresh that was implicitly invalidated mid-flight via invalidate(). The
// awaiter must NOT consume the routing map the refresh produced: by the
// invalidate's definition that map is stale. Awaiters typically retry,
// which starts a fresh post-invalidate refresh.
var errPKRangeCacheInvalidatedDuringRefresh = errors.New("partition key range cache: invalidated during refresh")

// errPKRangeCacheRefreshTimeout wraps a context.DeadlineExceeded that
// originated from the detached refresh's own pkRangeBackgroundRefreshTimeout
// (not from a caller's context). Wrapped so upstream code that uses
// errors.Is(err, context.DeadlineExceeded) to decide "my own deadline fired"
// does not misclassify a cache-refresh timeout as a caller-deadline timeout.
var errPKRangeCacheRefreshTimeout = errors.New("partition key range cache: detached refresh timed out")

// refreshOp represents an in-flight partition-key-range refresh for one
// container. Awaiters receive the (rm, err) pair by reading the fields after
// done is closed. The writes to rm/err happen-before close(op.done), which
// happens-before the receive from op.done — so awaiters reading the fields
// after the select observe a fully-published result (Go memory model).
type refreshOp struct {
	done chan struct{}
	rm   *collectionRoutingMap
	err  error
}

type pkRangeCacheEntry struct {
	mu         sync.Mutex // protects routingMap, inFlight, and generation
	routingMap *collectionRoutingMap
	inFlight   *refreshOp
	// generation is bumped by every invalidate() call. A refresh that
	// completes with a stale generation discards its result rather than
	// installing it, so invalidate cannot be silently overwritten by a
	// refresh that was already in flight when the invalidate happened.
	generation uint64
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
	// If the caller is already canceled and nobody else is fetching, don't
	// spawn a detached refresh just to immediately abandon it.
	if entry.inFlight == nil {
		if err := ctx.Err(); err != nil {
			entry.mu.Unlock()
			return nil, err
		}
	}
	op := c.ensureInFlightLocked(entry, containerLink, client)
	entry.mu.Unlock()

	return awaitRefresh(ctx, op)
}

// forceRefresh starts (or joins) a refresh for the given container and
// returns the resulting routing map.
func (c *partitionKeyRangeCache) forceRefresh(
	ctx context.Context,
	containerRID string,
	containerLink string,
	client *Client,
) (*collectionRoutingMap, error) {
	entry := c.getOrCreateEntry(containerRID)

	entry.mu.Lock()
	// If the caller is already canceled and nobody else is fetching, don't
	// spawn a detached refresh just to immediately abandon it.
	if entry.inFlight == nil {
		if err := ctx.Err(); err != nil {
			entry.mu.Unlock()
			return nil, err
		}
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
	startGeneration := entry.generation
	go c.runRefresh(entry, containerLink, client, op, startGeneration)
	return op
}

// runRefresh executes the change-feed refresh on a detached context bounded
// by pkRangeBackgroundRefreshTimeout. Caller cancellations do not abort the
// shared fetch. On completion the entry is updated under entry.mu, the
// in-flight slot is cleared, and awaiters are signaled by closing op.done.
//
// IMPORTANT: close(op.done) is performed while holding entry.mu so that any
// caller racing the completion edge observes either (a) the in-flight slot
// still set and joins, or (b) the slot cleared AND op.done already closed —
// never a state that would let it start a duplicate refresh in between.
//
// If entry.generation changed during the refresh (i.e., invalidate() was
// called), the result is discarded from the cache AND the awaiters of this
// op receive errPKRangeCacheInvalidatedDuringRefresh so they re-enter and
// start a fresh post-invalidate refresh instead of routing against the stale
// view they were waiting on.
//
// A context.DeadlineExceeded that originated from the detached refresh's own
// timeout is wrapped with errPKRangeCacheRefreshTimeout so awaiters can
// distinguish "the cache's own deadline fired" from "my own deadline fired."
func (c *partitionKeyRangeCache) runRefresh(
	entry *pkRangeCacheEntry,
	containerLink string,
	client *Client,
	op *refreshOp,
	startGeneration uint64,
) {
	ctx, cancel := newDetachedRefreshContext()
	defer cancel()

	start := time.Now()
	rm, err := c.refreshEntryDetached(ctx, containerLink, entry, client)
	elapsed := time.Since(start)

	// Wrap detached-timeout DeadlineExceeded so upstream callers using
	// errors.Is(err, context.DeadlineExceeded) don't misclassify it as
	// their own deadline firing.
	if err != nil && errors.Is(err, context.DeadlineExceeded) && ctx.Err() == context.DeadlineExceeded {
		err = fmt.Errorf("%w (container %s, elapsed %s, timeout %s): %w",
			errPKRangeCacheRefreshTimeout, containerLink, elapsed, pkRangeBackgroundRefreshTimeout, err)
	}

	entry.mu.Lock()
	invalidated := entry.generation != startGeneration
	if err == nil && rm != nil && !invalidated {
		entry.routingMap = rm
	}
	if invalidated {
		// Don't hand awaiters a routing map we just refused to install.
		// They must re-enter and start a post-invalidate refresh.
		op.rm = nil
		op.err = errPKRangeCacheInvalidatedDuringRefresh
	} else {
		op.rm = rm
		op.err = err
	}
	entry.inFlight = nil
	close(op.done)
	entry.mu.Unlock()

	// Single log line per refresh on the result side. Operators correlate the
	// detached HTTP spans by container link + elapsed; logging on start as
	// well would double the log volume on busy clients with frequent 410s.
	log.Writef(azlog.EventResponse,
		"partition key range cache refresh for container %s finished in %s (err=%v, invalidated=%t)",
		containerLink, elapsed, op.err, invalidated)
}

// newDetachedRefreshContext returns a context whose deadline is bounded by
// pkRangeBackgroundRefreshTimeout. It does NOT inherit cancellation from any
// caller — the caller's cancellation must not abort the shared refresh that
// other waiters are observing.
func newDetachedRefreshContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), pkRangeBackgroundRefreshTimeout)
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
// running but its result will be discarded (via the generation check in
// runRefresh) so the cached value cannot be silently replaced by data that
// pre-dates the invalidation.
func (c *partitionKeyRangeCache) invalidate(containerRID string) {
	c.mu.RLock()
	entry, exists := c.entries[containerRID]
	c.mu.RUnlock()

	if exists {
		entry.mu.Lock()
		entry.routingMap = nil
		entry.generation++
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
// change-feed page. Backoff is linear: base, 2*base, 3*base, ... with
// per-attempt jitter to avoid synchronized retries across containers.
const changeFeedPageRetryBaseDelay = 100 * time.Millisecond

// changeFeedPageRetryJitter caps the random jitter added to each retry
// delay. Keeps the worst-case extra latency small while breaking up
// synchronized retry storms across containers.
const changeFeedPageRetryJitter = 50 * time.Millisecond

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
//
// The retry layer here is intentionally narrow: the azcore pipeline already
// retries 5xx / 408 / 429 with backoff and Retry-After honoring. This loop
// only handles the residual cases where the pipeline gave up (transport
// errors, body-read timeouts, etc.) and where retrying the next page would
// preserve the pages already accumulated.
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
		// Linear backoff with jitter so a fleet of clients retrying across
		// many containers doesn't produce synchronized retry waves.
		delay := time.Duration(attempt+1)*changeFeedPageRetryBaseDelay + jitter(changeFeedPageRetryJitter)
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

// jitter returns a uniform random duration in [0, max).
func jitter(maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return 0
	}
	return time.Duration(rand.Int63n(int64(maxJitter)))
}

// isTransientPKRangeFetchError reports whether a /pkranges fetch error is
// worth retrying mid-pagination. Scope is intentionally narrow: the azcore
// pipeline already retries 5xx / 408 / 429 (with Retry-After honoring) on
// its own. This predicate handles only the residual transport-level cases
// the pipeline does not retry — actual network errors and bare 408s that
// slipped through. Returns false for:
//   - context cancellation / deadline
//   - any HTTP response error (including 5xx, 429, and 4xx)
//   - body-read / JSON-decode / programming errors
func isTransientPKRangeFetchError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	// HTTP responses: retry on server errors, throttling, and request timeout.
	// The pipeline may have already retried these, but the per-page retry here
	// preserves accumulated pages instead of restarting the entire drain.
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		switch {
		case respErr.StatusCode >= 500:
			return true
		case respErr.StatusCode == http.StatusTooManyRequests: // 429
			return true
		case respErr.StatusCode == http.StatusRequestTimeout: // 408
			return true
		}
		return false
	}
	// Transport-class errors only: connection reset, unexpected EOF,
	// net.Error timeouts. Body-read / JSON-decode errors are NOT
	// transient and should not be retried.
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ETIMEDOUT) {
		return true
	}
	return false
}

// refreshEntryDetached fetches PK ranges from the service and returns a fresh
// routing map. It snapshots the entry's previous routing map under entry.mu
// (briefly), then performs all network I/O without holding any lock. The
// caller is responsible for installing the returned map onto the entry.
//
// ctx is the detached, timeout-bounded refresh context created by
// runRefresh — it must NOT be a caller-scoped context.
func (c *partitionKeyRangeCache) refreshEntryDetached(
	ctx context.Context,
	containerLink string,
	entry *pkRangeCacheEntry,
	client *Client,
) (*collectionRoutingMap, error) {
	entry.mu.Lock()
	previousMap := entry.routingMap
	entry.mu.Unlock()

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
