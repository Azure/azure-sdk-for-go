// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

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
type partitionKeyRangeCache struct {
	mu      sync.RWMutex
	entries map[string]*pkRangeCacheEntry // keyed by container ResourceID
}

type pkRangeCacheEntry struct {
	mu         sync.Mutex // single-flights refresh for this container
	routingMap *collectionRoutingMap
}

func newPartitionKeyRangeCache() *partitionKeyRangeCache {
	return &partitionKeyRangeCache{
		entries: make(map[string]*pkRangeCacheEntry),
	}
}

// getRoutingMap returns the cached routing map for the given container RID.
// If the cache is empty for this container, it fetches from the service.
// containerLink is the name-based path used for the HTTP request.
func (c *partitionKeyRangeCache) getRoutingMap(
	ctx context.Context,
	containerRID string,
	containerLink string,
	client *Client,
) (*collectionRoutingMap, error) {
	// Fast path: read lock check
	c.mu.RLock()
	entry, exists := c.entries[containerRID]
	c.mu.RUnlock()

	if exists {
		entry.mu.Lock()
		if entry.routingMap != nil {
			rm := entry.routingMap
			entry.mu.Unlock()
			return rm, nil
		}
		// Cache entry exists but routing map is nil (invalidated) — refresh under lock
		rm, err := c.refreshEntry(ctx, containerLink, entry, client)
		entry.mu.Unlock()
		return rm, err
	}

	// Slow path: create entry
	c.mu.Lock()
	// Double check after acquiring write lock
	entry, exists = c.entries[containerRID]
	if !exists {
		entry = &pkRangeCacheEntry{}
		c.entries[containerRID] = entry
	}
	c.mu.Unlock()

	entry.mu.Lock()
	if entry.routingMap != nil {
		rm := entry.routingMap
		entry.mu.Unlock()
		return rm, nil
	}
	rm, err := c.refreshEntry(ctx, containerLink, entry, client)
	entry.mu.Unlock()
	return rm, err
}

// forceRefresh triggers an incremental refresh of the routing map for the given
// container. If the incremental merge fails (incomplete covering), it falls back
// to a full refresh. containerRID is the cache key; containerLink is used for HTTP requests.
func (c *partitionKeyRangeCache) forceRefresh(
	ctx context.Context,
	containerRID string,
	containerLink string,
	client *Client,
) (*collectionRoutingMap, error) {
	c.mu.RLock()
	entry, exists := c.entries[containerRID]
	c.mu.RUnlock()

	if !exists {
		// No entry yet — just do a normal get which will create and populate
		return c.getRoutingMap(ctx, containerRID, containerLink, client)
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()
	return c.refreshEntry(ctx, containerLink, entry, client)
}

// invalidate removes the cached routing map for the given container RID,
// forcing the next access to fetch fresh data.
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
// to prevent runaway requests during large-scale splits. Aligned with the Rust SDK.
const maxChangeFeedIterations = 1000

// changeFeedResult holds the result of draining all change-feed pages.
type changeFeedResult struct {
	ranges    []partitionKeyRange
	finalETag string
	completed bool // true only if loop terminated via 304 Not Modified
}

// fetchAllChangeFeedPages fetches change-feed pages starting from startETag
// until 304 Not Modified or the iteration cap. Returns the accumulated ranges,
// the final ETag, and whether the loop completed cleanly (304 received).
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
		result, err := fetchPartitionKeyRanges(ctx, containerLink, currentETag, client)
		if err != nil {
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

// refreshEntry fetches PK ranges from the service and populates the entry.
// It attempts an incremental refresh if a previous routing map with an ETag exists,
// accumulating all change-feed pages before merging. Falls back to a full
// change-feed refresh if the incremental merge is incomplete.
// Caller must hold entry.mu.
func (c *partitionKeyRangeCache) refreshEntry(
	ctx context.Context,
	containerLink string,
	entry *pkRangeCacheEntry,
	client *Client,
) (*collectionRoutingMap, error) {
	previousMap := entry.routingMap

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
				// No changes since last refresh — update ETag if changed
				if result.finalETag != previousMap.changeFeedETag {
					entry.routingMap = &collectionRoutingMap{
						orderedRanges:  previousMap.orderedRanges,
						rangeByID:      previousMap.rangeByID,
						goneRanges:     previousMap.goneRanges,
						changeFeedETag: result.finalETag,
					}
				}
				return entry.routingMap, nil
			}

			merged := previousMap.tryCombine(result.ranges, result.finalETag)
			if merged != nil {
				entry.routingMap = merged
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

	entry.routingMap = newMap
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
