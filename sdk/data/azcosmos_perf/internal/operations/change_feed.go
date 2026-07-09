// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package operations

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// ChangeFeedItemsOperation reads one incremental change feed page per execution.
type ChangeFeedItemsOperation struct {
	mu              sync.Mutex
	ranges          []azcosmos.FeedRange
	continuations   map[string]*string
	maxItems        int32
	excludedRegions []string
}

func NewChangeFeedItemsOperation(maxItems int32, excludedRegions []string) *ChangeFeedItemsOperation {
	return &ChangeFeedItemsOperation{
		continuations:   make(map[string]*string),
		maxItems:        maxItems,
		excludedRegions: excludedRegions,
	}
}

func (o *ChangeFeedItemsOperation) Name() string { return "ChangeFeedItems" }

func (o *ChangeFeedItemsOperation) Execute(ctx context.Context, c *azcosmos.ContainerClient) (*time.Duration, error) {
	ctx, collector := prepareContext(ctx, o.excludedRegions)
	if err := o.ensureRanges(ctx, c); err != nil {
		return nil, err
	}

	feedRange, key, continuation, err := o.pickRange()
	if err != nil {
		return nil, err
	}

	// TODO: Add AVAD/FFCF mode when the Go SDK exposes a public ChangeFeedMode option.
	response, err := c.GetChangeFeed(ctx, &azcosmos.ChangeFeedOptions{
		FeedRange:    &feedRange,
		Continuation: continuation,
		MaxItemCount: o.maxItems,
	})
	if err != nil {
		return nil, err
	}

	var next *string
	if response.RawResponse != nil && response.RawResponse.StatusCode == http.StatusNotModified {
		next = nil
	} else if response.ContinuationToken != "" {
		ct := response.ContinuationToken
		next = &ct
	}
	o.mu.Lock()
	o.continuations[key] = next
	o.mu.Unlock()

	return collector.duration(), nil
}

func (o *ChangeFeedItemsOperation) ensureRanges(ctx context.Context, c *azcosmos.ContainerClient) error {
	o.mu.Lock()
	initialized := len(o.ranges) > 0
	o.mu.Unlock()
	if initialized {
		return nil
	}

	ranges, err := c.GetFeedRanges(ctx)
	if err != nil {
		return err
	}
	if len(ranges) == 0 {
		return fmt.Errorf("container returned no feed ranges")
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.ranges) == 0 {
		o.ranges = append([]azcosmos.FeedRange(nil), ranges...)
		for _, feedRange := range ranges {
			o.continuations[feedRangeKey(feedRange)] = nil
		}
	}
	return nil
}

func (o *ChangeFeedItemsOperation) pickRange() (azcosmos.FeedRange, string, *string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.ranges) == 0 {
		return azcosmos.FeedRange{}, "", nil, fmt.Errorf("change feed ranges are not initialized")
	}
	feedRange := o.ranges[rand.Intn(len(o.ranges))]
	key := feedRangeKey(feedRange)
	var continuationCopy *string
	if continuation := o.continuations[key]; continuation != nil {
		ct := *continuation
		continuationCopy = &ct
	}
	return feedRange, key, continuationCopy, nil
}

func feedRangeKey(feedRange azcosmos.FeedRange) string {
	return feedRange.MinInclusive + "|" + feedRange.MaxExclusive
}
