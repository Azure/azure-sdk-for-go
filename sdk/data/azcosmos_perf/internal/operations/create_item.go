// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package operations

import (
	"context"
	"math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/seed"
	"github.com/google/uuid"
)

// CreateItemOperation creates a new item and adds it to the shared pool.
type CreateItemOperation struct {
	items           *seed.SharedItems
	excludedRegions []string
}

func NewCreateItemOperation(items *seed.SharedItems, excludedRegions []string) *CreateItemOperation {
	return &CreateItemOperation{items: items, excludedRegions: excludedRegions}
}

func (o *CreateItemOperation) Name() string { return "CreateItem" }

func (o *CreateItemOperation) Execute(ctx context.Context, c *azcosmos.ContainerClient) (*time.Duration, error) {
	ctx, collector := prepareContext(ctx, o.excludedRegions)
	id := uuid.NewString()
	partitionKey := uuid.NewString()
	item := PerfItem{
		ID:           id,
		PartitionKey: partitionKey,
		Value:        rand.Uint64(),
		Payload:      "perf-test-created",
	}
	body, err := marshalItem(item)
	if err != nil {
		return nil, err
	}
	_, err = c.CreateItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), body, nil)
	if err != nil {
		return nil, err
	}
	o.items.Push(seed.SeededItem{ID: id, PartitionKey: partitionKey})
	return collector.duration(), nil
}
