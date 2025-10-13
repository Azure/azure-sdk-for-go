// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package workloads

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func createRandomItem(i int) map[string]interface{} {
	return map[string]interface{}{
		"type":      "testItem",
		"createdAt": time.Now().UTC().Format(time.RFC3339Nano),
		"seq":       i,
		"value":     rand.Int63(), // pseudo-random payload
	}
}

func randomUpserts(ctx context.Context, container *azcosmos.ContainerClient, count int, pkField string) error {
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

				var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
				// re-upsert a document already written
				var num = rng.Intn(count) + 1
				item := createRandomItem(j.i)
				item["id"] = fmt.Sprintf("test-%d", num)
				item[pkField] = fmt.Sprintf("pk-%d", num)

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

func createClient(cfg workloadConfig) (*azcosmos.Client, error) {
	cred, err := azcosmos.NewKeyCredential(cfg.Key)
	if err != nil {
		return nil, err
	}
	opts := &azcosmos.ClientOptions{
		PreferredRegions: cfg.PreferredLocations,
		// Add EnableContentResponseOnWrite: true if you want full responses
	}
	return azcosmos.NewClientWithKey(cfg.Endpoint, cred, opts)
}
