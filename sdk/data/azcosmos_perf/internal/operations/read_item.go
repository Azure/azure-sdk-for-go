// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package operations

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/seed"
)

// ReadItemOperation reads a random seeded item by ID and partition key.
type ReadItemOperation struct {
	items           *seed.SharedItems
	excludedRegions []string
}

func NewReadItemOperation(items *seed.SharedItems, excludedRegions []string) *ReadItemOperation {
	return &ReadItemOperation{items: items, excludedRegions: excludedRegions}
}

func (o *ReadItemOperation) Name() string { return "ReadItem" }

func (o *ReadItemOperation) Execute(ctx context.Context, c *azcosmos.ContainerClient) (*time.Duration, error) {
	ctx, collector := prepareContext(ctx, o.excludedRegions)
	item := o.items.Random()
	_, err := c.ReadItem(ctx, azcosmos.NewPartitionKeyString(item.PartitionKey), item.ID, nil)
	if err != nil {
		return nil, err
	}
	return collector.duration(), nil
}
