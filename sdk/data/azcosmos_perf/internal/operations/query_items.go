// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package operations

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/seed"
)

// QueryItemsOperation runs a single-partition query against a random seeded key.
type QueryItemsOperation struct {
	items           *seed.SharedItems
	excludedRegions []string
}

func NewQueryItemsOperation(items *seed.SharedItems, excludedRegions []string) *QueryItemsOperation {
	return &QueryItemsOperation{items: items, excludedRegions: excludedRegions}
}

func (o *QueryItemsOperation) Name() string { return "QueryItems" }

func (o *QueryItemsOperation) Execute(ctx context.Context, c *azcosmos.ContainerClient) (*time.Duration, error) {
	ctx, collector := prepareContext(ctx, o.excludedRegions)
	item := o.items.Random()
	pk := azcosmos.NewPartitionKeyString(item.PartitionKey)
	query := "SELECT * FROM c WHERE c.partition_key = @pk"
	pager := c.NewQueryItemsPager(query, pk, &azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{{Name: "@pk", Value: item.PartitionKey}},
	})
	for pager.More() {
		if _, err := pager.NextPage(ctx); err != nil {
			return nil, err
		}
	}
	return collector.duration(), nil
}
