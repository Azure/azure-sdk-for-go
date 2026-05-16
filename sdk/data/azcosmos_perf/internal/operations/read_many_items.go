// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package operations

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/seed"
)

// ReadManyItemsOperation batch-reads random items from the shared pool.
type ReadManyItemsOperation struct {
	items           *seed.SharedItems
	batchSize       int
	excludedRegions []string
}

func NewReadManyItemsOperation(items *seed.SharedItems, batchSize int, excludedRegions []string) *ReadManyItemsOperation {
	return &ReadManyItemsOperation{items: items, batchSize: batchSize, excludedRegions: excludedRegions}
}

func (o *ReadManyItemsOperation) Name() string { return "ReadManyItems" }

func (o *ReadManyItemsOperation) Execute(ctx context.Context, c *azcosmos.ContainerClient) (*time.Duration, error) {
	ctx, collector := prepareContext(ctx, o.excludedRegions)
	sample := o.items.Sample(o.batchSize)
	items := make([]azcosmos.ItemIdentity, len(sample))
	for i, item := range sample {
		items[i] = azcosmos.ItemIdentity{ID: item.ID, PartitionKey: azcosmos.NewPartitionKeyString(item.PartitionKey)}
	}
	_, err := c.ReadManyItems(ctx, items, nil)
	if err != nil {
		return nil, err
	}
	return collector.duration(), nil
}
