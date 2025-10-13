// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package workloads

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

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

	// Build container properties
	props := azcosmos.ContainerProperties{
		ID: containerID,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/" + pkField},
			Kind:  azcosmos.PartitionKeyKindHash,
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

func upsertItemsConcurrently(ctx context.Context, container *azcosmos.ContainerClient, count int, pkField string) error {
	// Use a bounded worker pool to avoid oversaturating resources
	workers := 32
	if count < workers {
		workers = count
	}
	type job struct {
		i int
	}
	jobs := make(chan job, workers)
	errs := make(chan error, count)
	wg := &sync.WaitGroup{}

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				item := createRandomItem(j.i)
				item["id"] = fmt.Sprintf("test-%d", j.i)
				item[pkField] = fmt.Sprintf("pk-%d", j.i)

				// Marshal item to bytes; UpsertItem often takes []byte + partition key value
				body, err := json.Marshal(item)
				if err != nil {
					errs <- err
					continue
				}

				pk := azcosmos.NewPartitionKeyString(item[pkField].(string))
				_, err = container.UpsertItem(ctx, pk, body, nil)
				if err != nil {
					errs <- err
					continue
				}
				println("writing an item")
			}
		}()
	}

sendLoop:
	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			break sendLoop
		default:
			jobs <- job{i: i}
		}
	}
	close(jobs)
	wg.Wait()
	close(errs)

	// Aggregate errors if any
	var firstErr error
	for e := range errs {
		if firstErr == nil {
			firstErr = e
		}
	}
	return firstErr
}

// RunSetup creates the database/container and performs the concurrent upserts.
func RunSetup(ctx context.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	client, err := createClient(cfg)
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	println("Creating database...")
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
