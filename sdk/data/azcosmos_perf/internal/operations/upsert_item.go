// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package operations

import (
	"context"
	"math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/seed"
)

// UpsertItemOperation upserts an item in a random seeded partition.
type UpsertItemOperation struct {
	items           *seed.SharedItems
	excludedRegions []string
}

func NewUpsertItemOperation(items *seed.SharedItems, excludedRegions []string) *UpsertItemOperation {
	return &UpsertItemOperation{items: items, excludedRegions: excludedRegions}
}

func (o *UpsertItemOperation) Name() string { return "UpsertItem" }

func (o *UpsertItemOperation) Execute(ctx context.Context, c *azcosmos.ContainerClient) (*time.Duration, error) {
	ctx, collector := prepareContext(ctx, o.excludedRegions)
	seeded := o.items.Random()
	item := PerfItem{
		ID:           seeded.ID,
		PartitionKey: seeded.PartitionKey,
		Value:        rand.Uint64(),
		Payload:      "perf-test-payload",
	}
	body, err := marshalItem(item)
	if err != nil {
		return nil, err
	}
	_, err = c.UpsertItem(ctx, azcosmos.NewPartitionKeyString(item.PartitionKey), body, nil)
	if err != nil {
		return nil, err
	}
	return collector.duration(), nil
}
