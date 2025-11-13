// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package workloads

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// createDatabaseIfNotExists attempts to create the database, tolerating conflicts.
func createDatabaseIfNotExists(ctx context.Context, client *azcosmos.Client, dbID string) (*azcosmos.DatabaseClient, error) {
	dbClient, err := client.NewDatabase(dbID)
	if err != nil {
		return nil, err
	}
	props := azcosmos.DatabaseProperties{ID: dbID}
	_, err = client.CreateDatabase(ctx, props, nil)
	if err != nil {
		var azErr *azcore.ResponseError
		if errors.As(err, &azErr) {
			if azErr.StatusCode == 409 {
				return dbClient, nil // already exists
			}
		}
		return nil, err
	}
	return dbClient, nil
}

func createContainerIfNotExists(ctx context.Context, db *azcosmos.DatabaseClient, containerID, pkField string, desiredThroughput int32) (*azcosmos.ContainerClient, error) {
	containerClient, err := db.NewContainer(containerID)
	if err != nil {
		return nil, err
	}

	// Build container properties with vector indexing policy
	props := azcosmos.ContainerProperties{
		ID: containerID,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/" + pkField},
			Kind:  azcosmos.PartitionKeyKindHash,
		},
		VectorEmbeddingPolicy: &azcosmos.VectorEmbeddingPolicy{
			VectorEmbeddings: []azcosmos.VectorEmbedding{
				{
					Path:             "/embedding",
					DataType:         azcosmos.VectorDataTypeFloat32,
					DistanceFunction: azcosmos.VectorDistanceFunctionCosine,
					Dimensions:       10,
				},
			},
		},
		IndexingPolicy: &azcosmos.IndexingPolicy{
			Automatic:    true,
			IndexingMode: azcosmos.IndexingModeConsistent,
			IncludedPaths: []azcosmos.IncludedPath{
				{Path: "/*"},
			},
			ExcludedPaths: []azcosmos.ExcludedPath{
				{Path: "/\"_etag\"/?"},
				{Path: "/embedding/*"}, // Exclude vector path from standard indexing
			},
			VectorIndexes: []azcosmos.VectorIndex{
				{
					Path: "/embedding",
					Type: azcosmos.VectorIndexTypeDiskANN,
				},
			},
		},
	}

	throughput := azcosmos.NewManualThroughputProperties(desiredThroughput)
	// Try create
	_, err = db.CreateContainer(ctx, props, &azcosmos.CreateContainerOptions{
		ThroughputProperties: &throughput,
	})
	if err != nil {
		var azErr *azcore.ResponseError
		if errors.As(err, &azErr) {
			if azErr.StatusCode == 409 {
				return containerClient, nil // already exists
			}
		}
		return nil, err
	}

	return containerClient, nil
}

// RunSetup creates the database/container and performs the concurrent upserts.
func RunSetup(ctx context.Context, client *azcosmos.Client, cfg workloadConfig) error {

	dbClient, err := createDatabaseIfNotExists(ctx, client, cfg.DatabaseID)
	if err != nil {
		return fmt.Errorf("ensure database: %w", err)
	}

	container, err := createContainerIfNotExists(ctx, dbClient, cfg.ContainerID, cfg.PartitionKeyFieldName, int32(cfg.Throughput))
	if err != nil {
		return fmt.Errorf("ensure container: %w", err)
	}

	var count = cfg.LogicalPartitions

	log.Printf("Starting %d concurrent upserts...", count)

	if err := upsertItemsConcurrently(ctx, container, count, cfg.PartitionKeyFieldName); err != nil {
		return fmt.Errorf("upserts failed: %w", err)
	}

	log.Printf("Completed %d upserts.", count)
	return nil
}
