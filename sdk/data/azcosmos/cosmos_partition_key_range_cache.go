// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
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
//   - At most one refresh runs per container; concurrent callers share its result.
//   - Cached routing maps remain readable while a refresh is in flight.
//   - The refresh runs on context.Background(), so one caller's cancellation
//     doesn't abort the shared fetch. Each waiter honors its own ctx.
//   - 429 page failures retry indefinitely; all other errors fail fast.
//   - invalidate() bumps entry.generation. A refresh that completes against a
//     stale generation discards its result and surfaces
//     errPKRangeCacheInvalidatedDuringRefresh so awaiters retry against a
//     fresh post-invalidate refresh.
type partitionKeyRangeCache struct {
	mu      sync.RWMutex
	entries map[string]*pkRangeCacheEntry // keyed by container ResourceID
}

// errPKRangeCacheInvalidatedDuringRefresh signals that an in-flight refresh
// was invalidated mid-flight via invalidate(); its result is stale and must
// not be consumed. Awaiters re-enter to start a fresh post-invalidate refresh.
var errPKRangeCacheInvalidatedDuringRefresh = errors.New("partition key range cache: invalidated during refresh")

// refreshOp represents an in-flight refresh. Awaiters read rm/err after
// op.done is closed; writes happen-before the close, so the receive observes
// a fully-published result.
type refreshOp struct {
	done chan struct{}
	rm   *collectionRoutingMap
	err  error
}

type pkRangeCacheEntry struct {
	mu         sync.Mutex
	routingMap *collectionRoutingMap
	inFlight   *refreshOp
	// generation is bumped by invalidate(). A refresh that completes with a
	// stale generation discards its result rather than installing it.
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

// pkRangeCacheSentinelMaxAttempts caps how many times getRoutingMap /
// forceRefresh re-enter on errPKRangeCacheInvalidatedDuringRefresh. The cap
// only guards against a pathological invalidate-on-every-refresh storm.
const pkRangeCacheSentinelMaxAttempts = 3

// getRoutingMap returns the cached routing map for the given container RID.
// Cached maps are returned immediately even while a refresh is in flight;
// otherwise the caller joins or starts a shared refresh and waits on its
// own context. errPKRangeCacheInvalidatedDuringRefresh is handled internally.
func (c *partitionKeyRangeCache) getRoutingMap(
	ctx context.Context,
	containerRID string,
	containerLink string,
	client *Client,
) (*collectionRoutingMap, error) {
	entry := c.getOrCreateEntry(containerRID)
	var lastErr error
	for attempt := 0; attempt < pkRangeCacheSentinelMaxAttempts; attempt++ {
		entry.mu.Lock()
		if entry.routingMap != nil {
			rm := entry.routingMap
			entry.mu.Unlock()
			return rm, nil
		}
		// Don't spawn a detached refresh just to abandon it if the caller
		// is already canceled and nobody else is fetching.
		if entry.inFlight == nil {
			if err := ctx.Err(); err != nil {
				entry.mu.Unlock()
				return nil, err
			}
		}
		op := c.ensureInFlightLocked(entry, containerLink, client)
		entry.mu.Unlock()

		rm, err := awaitRefresh(ctx, op)
		if errors.Is(err, errPKRangeCacheInvalidatedDuringRefresh) {
			lastErr = err
			continue
		}
		return rm, err
	}
	return nil, lastErr
}

// forceRefresh starts or joins a refresh for the given container.
//
// forceRefresh joins any in-flight refresh rather than displacing it: a
// running /pkranges drain already pulls a fresh snapshot from the service,
// and a later 410 will retrigger another forceRefresh if needed (eventually
// consistent). Callers needing strict post-call freshness MUST call
// invalidate() first; the generation counter then guarantees the in-flight
// op's result is discarded and a fresh fetch runs.
//
// errPKRangeCacheInvalidatedDuringRefresh is handled internally — see
// getRoutingMap.
func (c *partitionKeyRangeCache) forceRefresh(
	ctx context.Context,
	containerRID string,
	containerLink string,
	client *Client,
) (*collectionRoutingMap, error) {
	entry := c.getOrCreateEntry(containerRID)
	var lastErr error
	for attempt := 0; attempt < pkRangeCacheSentinelMaxAttempts; attempt++ {
		entry.mu.Lock()
		// Don't spawn a detached refresh just to abandon it if the caller
		// is already canceled and nobody else is fetching.
		if entry.inFlight == nil {
			if err := ctx.Err(); err != nil {
				entry.mu.Unlock()
				return nil, err
			}
		}
		op := c.ensureInFlightLocked(entry, containerLink, client)
		entry.mu.Unlock()

		rm, err := awaitRefresh(ctx, op)
		if errors.Is(err, errPKRangeCacheInvalidatedDuringRefresh) {
			lastErr = err
			continue
		}
		return rm, err
	}
	return nil, lastErr
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

// runRefresh executes the change-feed refresh on context.Background() so
// caller cancellations don't abort the shared fetch. On completion the
// entry is updated and op.done is closed under entry.mu so a racing caller
// either joins the in-flight slot or observes the cleared slot AND a
// closed op.done — never a state that lets it start a duplicate refresh.
//
// If entry.generation changed during the refresh, the result is discarded
// and awaiters receive errPKRangeCacheInvalidatedDuringRefresh so they
// re-enter against a fresh post-invalidate refresh.
//
// 429 page failures retry indefinitely; there is no detached timeout. A
// wedged background fetch only affects cold-start awaiters until the next
// invalidate() bumps the generation.
func (c *partitionKeyRangeCache) runRefresh(
	entry *pkRangeCacheEntry,
	containerLink string,
	client *Client,
	op *refreshOp,
	startGeneration uint64,
) {
	ctx := context.Background()

	start := time.Now()
	rm, err := c.refreshEntryDetached(ctx, containerLink, entry, client)
	elapsed := time.Since(start)

	entry.mu.Lock()
	invalidated := entry.generation != startGeneration
	if err == nil && rm != nil && !invalidated {
		entry.routingMap = rm
	}
	if invalidated {
		// Don't hand awaiters a routing map we refused to install.
		op.rm = nil
		op.err = errPKRangeCacheInvalidatedDuringRefresh
	} else {
		op.rm = rm
		op.err = err
	}
	entry.inFlight = nil
	close(op.done)
	entry.mu.Unlock()

	log.Writef(azlog.EventResponse,
		"partition key range cache refresh for container %s finished in %s (err=%v, invalidated=%t)",
		containerLink, elapsed, op.err, invalidated)
}

// awaitRefresh blocks until the refresh completes or the caller's context
// is cancelled. The refresh continues in the background even when individual
// awaiters return early via ctx.Err().
func awaitRefresh(ctx context.Context, op *refreshOp) (*collectionRoutingMap, error) {
	select {
	case <-op.done:
		return op.rm, op.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// invalidate clears the cached routing map and bumps the entry's generation.
// An in-flight refresh continues but its result is discarded (via the
// generation check in runRefresh) so a pre-invalidate fetch cannot silently
// overwrite the cache.
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

// Linear backoff with jitter for 429 retries: base, 2*base, ... capped at
// changeFeedPageRetryMaxDelay. Exposed as vars so tests can shrink them.
var (
	changeFeedPageRetryBaseDelay = 100 * time.Millisecond
	changeFeedPageRetryJitter    = 50 * time.Millisecond
	changeFeedPageRetryMaxDelay  = 5 * time.Second
)

// changeFeedResult holds the result of draining all change-feed pages.
type changeFeedResult struct {
	ranges    []partitionKeyRange
	finalETag string
	completed bool // true only if loop terminated via 304 Not Modified
}

// fetchAllChangeFeedPages drains change-feed pages from startETag until 304
// Not Modified or the iteration cap. 429 page failures retry indefinitely
// (with capped linear backoff + jitter) so accumulated pages aren't
// discarded under throttling; all other errors fail fast and the azcore
// pipeline owns retry for those classes.
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
			// Log how far we got so operators can correlate partial-drain
			// failures with the next refresh restarting from scratch.
			log.Writef(azlog.EventResponse, "partition key range change-feed page failed for container %s after %d successful pages (%d ranges accumulated): %v", containerLink, i, len(allRanges), err)
			return changeFeedResult{}, err
		}

		if result.notModified {
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
	log.Writef(azlog.EventResponse, "partition key range change-feed loop exited without reaching 304 Not Modified after %d iterations for container %s (accumulated %d ranges)", maxChangeFeedIterations, containerLink, len(allRanges))
	return changeFeedResult{ranges: allRanges, finalETag: currentETag, completed: false}, nil
}

// fetchOneChangeFeedPageWithRetry fetches one change-feed page, retrying
// indefinitely on 429 so a mid-drain throttle doesn't discard accumulated
// pages. All other errors are surfaced immediately: the azcore pipeline
// already retries 5xx / 408 / transport with backoff and Retry-After
// honoring, and re-retrying here would double-retry.
func fetchOneChangeFeedPageWithRetry(
	ctx context.Context,
	containerLink string,
	currentETag string,
	client *Client,
) (fetchPartitionKeyRangesResult, error) {
	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return fetchPartitionKeyRangesResult{}, err
		}
		result, err := fetchPartitionKeyRanges(ctx, containerLink, currentETag, client)
		if err == nil {
			return result, nil
		}
		if !isThrottleError(err) {
			return fetchPartitionKeyRangesResult{}, err
		}
		// Linear backoff (capped) with jitter to avoid synchronized retry
		// waves across containers.
		delay := time.Duration(attempt+1) * changeFeedPageRetryBaseDelay
		if delay > changeFeedPageRetryMaxDelay {
			delay = changeFeedPageRetryMaxDelay
		}
		delay += jitter(changeFeedPageRetryJitter)
		log.Writef(azlog.EventResponse, "partition key range change-feed page throttled for container %s (attempt %d, retrying in %s): %v", containerLink, attempt+1, delay, err)
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return fetchPartitionKeyRangesResult{}, ctx.Err()
		}
	}
}

// isThrottleError reports whether err is a 429 Too Many Requests response.
func isThrottleError(err error) bool {
	var respErr *azcore.ResponseError
	return errors.As(err, &respErr) && respErr.StatusCode == http.StatusTooManyRequests
}

// jitter returns a uniform random duration in [0, max).
func jitter(maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return 0
	}
	return time.Duration(rand.Int63n(int64(maxJitter)))
}

// refreshEntryDetached fetches PK ranges and returns a fresh routing map.
// All network I/O runs without holding any lock; the caller installs the
// returned map on the entry. ctx is the detached refresh context created
// by runRefresh — never a caller-scoped context.
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
