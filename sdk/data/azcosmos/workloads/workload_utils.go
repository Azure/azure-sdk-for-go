// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package workloads

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Azure/azure-cosmos-client-engine/go/azcosmoscx"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

const defaultConcurrency = 32

// GetPKValue returns the partition key value for a given logical index.
// If pkField is "id" the PK equals the item id (test-<index>), otherwise it uses "pk-<index>".
func GetPKValue(pkField string, index int) string {
	if pkField == "id" {
		return fmt.Sprintf("test-%d", index)
	}
	return fmt.Sprintf("pk-%d", index)
}

// rwOperation enumerates random read/write/query operations.
type rwOperation int

const (
	opUpsert rwOperation = iota
	opRead
	opQuery
	opOpCount // sentinel for number of operations
)

func createRandomItem(i int) map[string]interface{} {
	return map[string]interface{}{
		"type":      "testItem",
		"createdAt": time.Now().UTC().Format(time.RFC3339Nano),
		"seq":       i,
		"value":     rand.Int63(), // pseudo-random payload
		"embedding": createRandomEmbedding(),
	}
}

func createRandomEmbedding() []float64 {
	return []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
}

// runConcurrent executes count indexed jobs across at most workers goroutines.
// jf should be idempotent per index; it receives a per-worker RNG (not safe to share across workers).
func runConcurrent(ctx context.Context, count, workers int, jobFunction func(ctx context.Context, index int, rng *rand.Rand) error) error {
	if count <= 0 {
		return errors.New("count must be > 0")
	}
	if workers <= 0 {
		workers = 1
	}
	if count < workers {
		workers = count
	}

	type job struct{ i int }
	jobs := make(chan job, workers)
	errs := make(chan error, count)
	wg := &sync.WaitGroup{}

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// Seed rng per worker (unique-ish seed)
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)<<32))
			for j := range jobs {
				if ctx.Err() != nil {
					return
				}
				if err := jobFunction(ctx, j.i, rng); err != nil {
					select {
					case errs <- err:
					default: // channel full; drop to avoid blocking
					}
				}
			}
		}(w)
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

	var collected []error
	for e := range errs {
		if e != nil {
			collected = append(collected, e)
		}
	}
	return errors.Join(collected...)
}

func upsertItemsConcurrently(ctx context.Context, container *azcosmos.ContainerClient, count int, pkField string) error {
	return runConcurrent(ctx, count, defaultConcurrency, func(ctx context.Context, i int, rng *rand.Rand) error {
		item := createRandomItem(i)
		id := fmt.Sprintf("test-%d", i)
		pkVal := GetPKValue(pkField, i)
		item["id"] = id
		item[pkField] = pkVal
		body, err := json.Marshal(item)
		if err != nil {
			return err
		}
		pk := azcosmos.NewPartitionKeyString(pkVal)
		_, err = container.UpsertItem(ctx, pk, body, nil)
		return err
	})
}

func randomReadWriteQueries(ctx context.Context, container *azcosmos.ContainerClient, count int, pkField string) error {
	return runConcurrent(ctx, count, defaultConcurrency, func(ctx context.Context, i int, rng *rand.Rand) error {
		// pick a random existing (or future) document index to operate on
		num := rng.Intn(count) + 1
		id := fmt.Sprintf("test-%d", num)
		pkVal := GetPKValue(pkField, num)

		pk := azcosmos.NewPartitionKeyString(pkVal)

		op := rwOperation(rng.Intn(int(opOpCount)))
		switch op {
		case opUpsert:
			item := createRandomItem(i)
			item["id"] = id
			item[pkField] = pkVal
			body, err := json.Marshal(item)
			if err != nil {
				log.Printf("randomRW marshal error id=%s pk=%s: %v", id, pkVal, err)
				return err
			}
			if _, err := container.UpsertItem(ctx, pk, body, nil); err != nil {
				log.Printf("upsert error id=%s pk=%s: %v", id, pkVal, err)
				return err
			}
		case opRead:
			if _, err := container.ReadItem(ctx, pk, id, nil); err != nil {
				log.Printf("read error id=%s pk=%s: %v", id, pkVal, err)
				return err
			}
		case opQuery:
			pager := container.NewQueryItemsPager(
				"SELECT * FROM c WHERE c.id = @id",
				pk,
				&azcosmos.QueryOptions{QueryParameters: []azcosmos.QueryParameter{{Name: "@id", Value: id}}},
			)
			for pager.More() {
				if _, err := pager.NextPage(ctx); err != nil {
					log.Printf("query error id=%s pk=%s: %v", id, pkVal, err)
					return err
				}
			}
		}
		return nil
	})
}

func vectorSearchQueries(ctx context.Context, container *azcosmos.ContainerClient, count int, pkField string) error {
	return runConcurrent(ctx, count, defaultConcurrency, func(ctx context.Context, i int, rng *rand.Rand) error {
		embedding := createRandomEmbedding()

		pkVal := GetPKValue(pkField, rng.Intn(count)+1)

		pager := container.NewQueryItemsPager(
			"SELECT TOP 10 c.id FROM c ORDER BY VectorDistance(c.embedding, @embedding)",
			azcosmos.NewPartitionKeyString(pkVal),
			&azcosmos.QueryOptions{
				QueryParameters: []azcosmos.QueryParameter{{Name: "@embedding", Value: embedding}},
				QueryEngine:     azcosmoscx.NewQueryEngine(),
			},
		)
		for pager.More() {
			if _, err := pager.NextPage(ctx); err != nil {
				log.Printf("vs query error: %v", err)
				return err
			}
		}
		return nil
	})
}

// CreateClient creates a Cosmos DB client based on the workload configuration.
// cfg - The workload configuration
func CreateClient(cfg WorkloadConfig) (*azcosmos.Client, error) {
	log.Printf("Creating client for endpoint %s with preferred regions: %v", cfg.Endpoint, cfg.PreferredLocations)
	telemetryOpts := policy.TelemetryOptions{
		// create a random guid for the application ID to distinguish this workload's telemetry
		ApplicationID: uuid.New().String(),
	}
	// log application ID for correlating requests in service
	log.Printf("Using ApplicationID: %s", telemetryOpts.ApplicationID)
	opts := &azcosmos.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Telemetry: telemetryOpts,
		},
		PreferredRegions: cfg.PreferredLocations,
	}
	if cfg.Key != "" {
		cred, err := azcosmos.NewKeyCredential(cfg.Key)
		if err != nil {
			return nil, err
		}

		return azcosmos.NewClientWithKey(cfg.Endpoint, cred, opts)
	} else {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}
		return azcosmos.NewClient(cfg.Endpoint, cred, opts)
	}
}
