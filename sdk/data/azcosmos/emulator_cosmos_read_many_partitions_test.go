// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

func TestReadMany_MultiplePhysicalPartitions(t *testing.T) {
	var client *Client
	if connectionString := os.Getenv("AZURE_COSMOS_CONNECTION_STRING"); connectionString != "" {
		var err error
		client, err = NewClientFromConnectionString(connectionString, nil)
		require.NoError(t, err)
	} else {
		emulator := newEmulatorTests(t)
		client = emulator.getClient(t, tracing.Provider{})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	databaseID := fmt.Sprintf("readmany-partitions-%d", time.Now().UnixNano())
	_, err := client.CreateDatabase(ctx, DatabaseProperties{ID: databaseID}, nil)
	require.NoError(t, err)
	database, err := client.NewDatabase(databaseID)
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()
		_, err := database.Delete(cleanupCtx, nil)
		require.NoError(t, err)
	})

	for _, version := range []int{0, 1, 2} {
		t.Run(fmt.Sprintf("version-%d", version), func(t *testing.T) {
			// More than 10,000 RU/s requires multiple physical partitions on Azure.
			throughput := NewManualThroughputProperties(11000)
			containerID := fmt.Sprintf("version-%d", version)
			_, err := database.CreateContainer(ctx, ContainerProperties{
				ID: containerID,
				PartitionKeyDefinition: PartitionKeyDefinition{
					Paths: []string{"/pk"}, Kind: PartitionKeyKindHash, Version: version,
				},
			}, &CreateContainerOptions{ThroughputProperties: &throughput})
			require.NoError(t, err)
			container, err := database.NewContainer(containerID)
			require.NoError(t, err)
			t.Cleanup(func() {
				cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), time.Minute)
				defer cleanupCancel()
				_, err := container.Delete(cleanupCtx, nil)
				require.NoError(t, err)
			})

			ranges, err := container.getPartitionKeyRanges(ctx, nil)
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(ranges.PartitionKeyRanges), 2, "fixture must have multiple physical partitions")

			keys := []string{"", "asdf", "asdf", strings.Repeat("a", 100), strings.Repeat("a", 101), strings.Repeat("\u00e9", 101)}
			for i := 0; i < 32; i++ {
				keys = append(keys, fmt.Sprintf("tenant-%02d", i))
			}
			identities := make([]ItemIdentity, len(keys))
			expectedIDs := make([]string, len(keys))
			for i, key := range keys {
				id := fmt.Sprintf("item-%02d", i)
				body, err := json.Marshal(map[string]string{"id": id, "pk": key})
				require.NoError(t, err)
				pk := NewPartitionKeyString(key)
				_, err = container.CreateItem(ctx, pk, body, nil)
				require.NoError(t, err)
				_, err = container.ReadItem(ctx, pk, id, nil)
				require.NoError(t, err)
				identities[i] = ItemIdentity{ID: id, PartitionKey: pk}
				expectedIDs[i] = id
			}

			// Query each service-provided range directly: do not use the EPK encoder
			// under test to establish which ranges actually contain the fixtures.
			var placedIDs []string
			occupiedRanges := 0
			for _, physicalRange := range ranges.PartitionKeyRanges {
				ids := readManyPhysicalRangeIDs(t, ctx, container, physicalRange.ID)
				placedIDs = append(placedIDs, ids...)
				if len(ids) > 0 {
					occupiedRanges++
				}
			}
			require.ElementsMatch(t, expectedIDs, placedIDs)
			require.GreaterOrEqual(t, occupiedRanges, 2, "items must occupy multiple physical partitions")
			t.Logf("%d items occupy %d of %d physical ranges", len(placedIDs), occupiedRanges, len(ranges.PartitionKeyRanges))

			check := func(t *testing.T, items []ItemIdentity, expected []string, options *ReadManyItemsOptions) {
				t.Helper()
				response, err := container.ReadManyItems(ctx, items, options)
				require.NoError(t, err)
				var actual []string
				for _, raw := range response.Items {
					var document struct {
						ID string `json:"id"`
					}
					require.NoError(t, json.Unmarshal(raw, &document))
					actual = append(actual, document.ID)
				}
				require.ElementsMatch(t, expected, actual)
			}
			for _, concurrency := range []int32{1, 4} {
				t.Run(fmt.Sprintf("all-items-concurrency-%d", concurrency), func(t *testing.T) {
					check(t, identities, expectedIDs, &ReadManyItemsOptions{MaxConcurrency: &concurrency})
				})
			}
			t.Run("missing-item", func(t *testing.T) {
				items := append([]ItemIdentity{}, identities...)
				items = append(items, ItemIdentity{ID: "missing", PartitionKey: NewPartitionKeyString("asdf")})
				check(t, items, expectedIDs, nil)
			})
		})
	}
}

func readManyPhysicalRangeIDs(t *testing.T, ctx context.Context, container *ContainerClient, rangeID string) []string {
	t.Helper()
	operationContext := pipelineRequestOptions{
		resourceType: resourceTypeDocument, resourceAddress: container.link,
	}
	path, err := generatePathForNameBased(resourceTypeDocument, container.link, true)
	require.NoError(t, err)
	pageSize := int32(3)
	options := &QueryOptions{PageSizeHint: pageSize}
	var ids []string
	for {
		response, err := container.database.client.sendQueryRequest(
			path, ctx, "SELECT c.id FROM c", nil, operationContext, options,
			func(req *policy.Request) {
				req.Raw().Header.Set(cosmosHeaderPartitionKeyRangeId, rangeID)
			})
		require.NoError(t, err)
		page, err := newQueryResponse(response)
		require.NoError(t, err)
		for _, raw := range page.Items {
			var document struct {
				ID string `json:"id"`
			}
			require.NoError(t, json.Unmarshal(raw, &document))
			ids = append(ids, document.ID)
		}
		if page.ContinuationToken == nil {
			return ids
		}
		options.ContinuationToken = page.ContinuationToken
	}
}
