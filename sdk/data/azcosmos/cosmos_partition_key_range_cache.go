// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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

// maxIncrementalRefreshIterations caps the number of incremental fetch loops
// to prevent runaway requests during large-scale splits. Aligned with the Rust SDK.
const maxIncrementalRefreshIterations = 1000

// refreshEntry fetches PK ranges from the service and populates the entry.
// It attempts an incremental refresh if a previous routing map with an ETag exists,
// looping until 304 Not Modified (capped at maxIncrementalRefreshIterations).
// Falls back to a full change-feed refresh if the incremental merge is incomplete.
// The full refresh also uses the change-feed mechanism (A-IM: Incremental feed)
// and loops until 304, matching the behavior of the .NET and Python SDKs.
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
		var allIncrementalRanges []partitionKeyRange
		currentETag := previousMap.changeFeedETag
		incrementalComplete := false
		for i := 0; i < maxIncrementalRefreshIterations; i++ {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			ranges, newETag, err := fetchPartitionKeyRanges(ctx, containerLink, currentETag, client)
			if err != nil {
				return nil, err
			}

			if len(ranges) == 0 {
				// 304 Not Modified — all incremental changes collected
				if newETag != "" {
					currentETag = newETag
				}
				incrementalComplete = true
				break
			}
			allIncrementalRanges = append(allIncrementalRanges, ranges...)
			currentETag = newETag
		}

		if incrementalComplete {
			if len(allIncrementalRanges) == 0 {
				// No changes since last refresh — update ETag if changed
				if currentETag != previousMap.changeFeedETag {
					entry.routingMap = &collectionRoutingMap{
						orderedRanges:  previousMap.orderedRanges,
						rangeByID:      previousMap.rangeByID,
						goneRanges:     previousMap.goneRanges,
						changeFeedETag: currentETag,
					}
				}
				return entry.routingMap, nil
			}

			merged := previousMap.tryCombine(allIncrementalRanges, currentETag)
			if merged != nil {
				entry.routingMap = merged
				return merged, nil
			}
		}

		// Incremental merge failed or loop exhausted — fall through to full refresh.
	}

	// Full change-feed refresh: fetch all ranges using A-IM without an ETag,
	// looping until 304 Not Modified to accumulate all pages.
	var allRanges []partitionKeyRange
	var currentETag string
	for i := 0; i < maxIncrementalRefreshIterations; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		ranges, newETag, err := fetchPartitionKeyRanges(ctx, containerLink, currentETag, client)
		if err != nil {
			return nil, err
		}
		if len(ranges) == 0 {
			// 304 Not Modified — all ranges collected
			if newETag != "" {
				currentETag = newETag
			}
			break
		}
		allRanges = append(allRanges, ranges...)
		currentETag = newETag
	}

	newMap := newCollectionRoutingMap(allRanges, currentETag)
	if !isCompleteSetOfRanges(newMap.orderedRanges) {
		gap := describeRangeGap(newMap.orderedRanges)
		return nil, fmt.Errorf("partition key range cache refresh failed: service returned an incomplete set of ranges for container %s (got %d ranges, issue: %s). This may indicate a transient issue during a partition split", containerLink, len(newMap.orderedRanges), gap)
	}

	entry.routingMap = newMap
	return newMap, nil
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
) ([]partitionKeyRange, string, error) {
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypePartitionKeyRange,
		resourceAddress: containerLink,
	}

	path, err := generatePathForNameBased(resourceTypePartitionKeyRange, operationContext.resourceAddress, true)
	if err != nil {
		return nil, "", err
	}

	o := &partitionKeyRangeOptions{}

	azResponse, err := client.sendGetRequest(
		path,
		ctx,
		operationContext,
		o,
		func(req *policy.Request) {
			req.Raw().Header.Set(cosmosHeaderChangeFeed, cosmosHeaderValuesChangeFeed)
			req.Raw().Header.Set(cosmosHeaderMaxItemCount, "-1")
			if changeFeedETag != "" {
				req.Raw().Header.Set(headerIfNoneMatch, changeFeedETag)
			}
		})
	if err != nil {
		return nil, "", err
	}

	newETag := azResponse.Header.Get(cosmosHeaderEtag)

	// 304 Not Modified means no changes
	if azResponse.StatusCode == http.StatusNotModified {
		_ = azResponse.Body.Close()
		return nil, newETag, nil
	}

	body, err := azruntime.Payload(azResponse)
	if err != nil {
		return nil, "", err
	}
	_ = azResponse.Body.Close()

	var response partitionKeyRangeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, "", err
	}

	return response.PartitionKeyRanges, newETag, nil
}
