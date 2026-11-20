// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"strings"
)

// FeedRange is an immutable effective-partition-key interval [MinInclusive, MaxExclusive).
// Equal bounds denote a single EPK, not an empty interval, as defined by the native driver.
// It owns no native resources and can be shared across workers. Its zero value is invalid.
type FeedRange struct {
	serialized   string
	minInclusive string
	maxExclusive string
	logical      bool
}

// NewFeedRange constructs an effective-partition-key interval using native validation.
// Empty minInclusive denotes the beginning of the key space; "FF" is its exclusive end.
// This requires a native build exposing the feed-range extension.
func NewFeedRange(minInclusive, maxExclusive string) (FeedRange, error) {
	if strings.ContainsAny(minInclusive+maxExclusive, "\x00") {
		return FeedRange{}, errors.New("azcosmos: feed range bounds must not contain a NUL byte")
	}
	return newFeedRange(minInclusive, maxExclusive)
}

// ParseFeedRange parses the native driver's stable, opaque string representation.
// Serialization preserves EPK bounds, not a logical partition key's identity.
func ParseFeedRange(serialized string) (FeedRange, error) {
	if serialized == "" || strings.IndexByte(serialized, 0) >= 0 {
		return FeedRange{}, errors.New("azcosmos: serialized feed range must not be empty or contain a NUL byte")
	}
	return parseFeedRange(serialized)
}

// String returns the persistable native representation, or "" for the zero value.
// Persist it with the container identity and that range's checkpoint, not a physical range ID.
func (r FeedRange) String() string { return r.serialized }

// MinInclusive returns the inclusive effective-partition-key lower bound.
func (r FeedRange) MinInclusive() string { return r.minInclusive }

// MaxExclusive returns the exclusive effective-partition-key upper bound.
func (r FeedRange) MaxExclusive() string { return r.maxExclusive }

// IsLogicalPartition reports the native range discriminator.
// Public FeedRange values are normalized to persistable EPK ranges and return false.
// Use NewFeedScopeForPartitionKey when logical-key identity is required for routing.
func (r FeedRange) IsLogicalPartition() bool { return r.logical }

// IsSubsetOf reports whether the native driver considers r contained in other.
// Invalid zero values return an error.
// Equal-bound point comparisons inherit https://github.com/Azure/azure-sdk-for-rust/issues/5638.
func (r FeedRange) IsSubsetOf(other FeedRange) (bool, error) {
	return r.compare(other, false)
}

// Overlaps reports whether the native driver considers the two intervals intersecting.
// Adjacent half-open intervals do not overlap. Invalid zero values return an error.
// Equal-bound point comparisons inherit https://github.com/Azure/azure-sdk-for-rust/issues/5638.
func (r FeedRange) Overlaps(other FeedRange) (bool, error) {
	return r.compare(other, true)
}

func (r FeedRange) compare(other FeedRange, overlaps bool) (bool, error) {
	if err := r.validate(); err != nil {
		return false, err
	}
	if err := other.validate(); err != nil {
		return false, err
	}
	return compareFeedRanges(r, other, overlaps)
}

func (r FeedRange) validate() error {
	if r.serialized == "" {
		return errors.New("azcosmos: feed range must not be the zero value")
	}
	return feedRangeSupportError()
}

// GetFeedRangesOptions configures discovery. Nil selects native cached discovery.
type GetFeedRangesOptions struct {
	// ForceRefresh requests a native routing-map refresh. Go maintains no topology cache.
	// A failed refresh can return cached success: https://github.com/Azure/azure-sdk-for-rust/issues/5637.
	ForceRefresh bool
}

// GetFeedRangesResponse contains independently owned ranges in native order.
type GetFeedRangesResponse struct {
	Response
	FeedRanges []FeedRange
}

// FeedRangeResponse contains a constructed range and any native setup-read metadata.
type FeedRangeResponse struct {
	Response
	FeedRange FeedRange
}

// GetFeedRanges discovers the container's current effective-partition-key intervals.
// Ranges are logical scopes, not stable worker IDs; native cursors handle topology changes.
// Context cancellation bounds the Go wait; Client.Close drains admitted native work.
func (c *ContainerClient) GetFeedRanges(ctx context.Context, options *GetFeedRangesOptions) (GetFeedRangesResponse, error) {
	if ctx == nil {
		return GetFeedRangesResponse{}, errors.New("azcosmos: context must not be nil")
	}
	release, err := c.database.client.acquire()
	if err != nil {
		return GetFeedRangesResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return GetFeedRangesResponse{}, err
	}
	forceRefresh := options != nil && options.ForceRefresh
	return c.getFeedRanges(ctx, forceRefresh)
}

// GetFeedRangeForPartitionKey constructs a range using the container's native key definition.
// Complete keys and hierarchical prefixes are normalized to persistable EPK bounds.
// Use NewFeedScopeForPartitionKey to preserve logical-key identity for reads.
func (c *ContainerClient) GetFeedRangeForPartitionKey(ctx context.Context, partitionKey PartitionKey) (FeedRangeResponse, error) {
	if err := validateItemArguments(partitionKey); err != nil {
		return FeedRangeResponse{}, err
	}
	if ctx == nil {
		return FeedRangeResponse{}, errors.New("azcosmos: context must not be nil")
	}
	release, err := c.database.client.acquire()
	if err != nil {
		return FeedRangeResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return FeedRangeResponse{}, err
	}
	return c.getFeedRangeForPartitionKey(ctx, partitionKey)
}
