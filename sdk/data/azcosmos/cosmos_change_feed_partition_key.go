// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal/epk"
)

// changeFeedPartitionKeyEPK returns the effective partition key (EPK) of
// options.PartitionKey when it identifies exactly one logical partition. It
// returns "" when the read isn't scoped to a single logical partition: no
// PartitionKey was set, or the key is a prefix of a hierarchical partition key,
// which spans an EPK range rather than a single point.
func (c *ContainerClient) changeFeedPartitionKeyEPK(ctx context.Context, options *ChangeFeedOptions) (string, error) {
	if options == nil || options.PartitionKey == nil {
		return "", nil
	}
	pkDef, err := c.getPartitionKeyDefinition(ctx)
	if err != nil {
		return "", err
	}
	r, err := computeEPKRange(options.PartitionKey, pkDef)
	if err != nil {
		return "", fmt.Errorf("ChangeFeedOptions: computing the effective partition key: %w", err)
	}
	if r.Min != r.Max {
		return "", nil
	}
	return r.Min, nil
}

// getPartitionKeyDefinition returns the container's partition key definition,
// from the container properties cache when one is available.
func (c *ContainerClient) getPartitionKeyDefinition(ctx context.Context) (PartitionKeyDefinition, error) {
	if c.database.client.getContainerCache() != nil {
		props, err := c.database.client.getContainerCache().getProperties(ctx, c)
		if err != nil {
			return PartitionKeyDefinition{}, err
		}
		return props.PartitionKeyDefinition, nil
	}
	resp, err := c.Read(ctx, nil)
	if err != nil {
		return PartitionKeyDefinition{}, err
	}
	return resp.ContainerProperties.PartitionKeyDefinition, nil
}

// scopeChangeFeedQueueToEPK reduces a continuation queue to the single entry
// that covers pkEPK, re-scoped to the point [pkEPK, pkEPK] and carrying that
// entry's ETag. This lets a token issued for a whole FeedRange (one entry per
// physical partition) resume a read of one logical partition, and keeps tokens
// issued for a partition key stable across splits, because the key's EPK never
// moves. Returns an error when no entry covers pkEPK, which means the token
// was issued for a different partition key or FeedRange.
func scopeChangeFeedQueueToEPK(queue []changeFeedRange, pkEPK string) ([]changeFeedRange, error) {
	for _, entry := range queue {
		if !changeFeedRangeContainsEPK(entry, pkEPK) {
			continue
		}
		scoped := changeFeedRange{MinInclusive: pkEPK, MaxExclusive: pkEPK}
		if entry.ContinuationToken != nil {
			etagCopy := *entry.ContinuationToken
			scoped.ContinuationToken = &etagCopy
		}
		return []changeFeedRange{scoped}, nil
	}
	return nil, fmt.Errorf("the continuation token doesn't cover ChangeFeedOptions.PartitionKey; it was issued for a different partition key or FeedRange")
}

// changeFeedRangeContainsEPK reports whether entry covers the EPK value. An
// entry whose bounds are equal is a single partition key and covers only that
// EPK; otherwise the entry is the half-open range [MinInclusive, MaxExclusive).
func changeFeedRangeContainsEPK(entry changeFeedRange, value string) bool {
	if entry.MinInclusive == entry.MaxExclusive {
		return epk.CompareEPK(entry.MinInclusive, value) == 0
	}
	return epk.CompareEPK(entry.MinInclusive, value) <= 0 &&
		epk.CompareEPK(value, normalizeMaxBoundary(entry.MaxExclusive)) < 0
}

// isPartitionKeyChangeFeedQueue reports whether queue was issued for a read
// scoped to one partition key. Such a queue holds an entry whose bounds are
// the same EPK point, which never describes a FeedRange.
func isPartitionKeyChangeFeedQueue(queue []changeFeedRange) bool {
	for _, entry := range queue {
		if entry.MinInclusive == entry.MaxExclusive {
			return true
		}
	}
	return false
}
