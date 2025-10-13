// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package workloads

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func randomReadWriteQueries(ctx context.Context, container *azcosmos.ContainerClient, count int, pkField string) error {
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
				rng := rand.New(rand.NewSource(time.Now().UnixNano()))
				// re-upsert/read/query a document (some may not exist yet which can surface 404s)
				num := rng.Intn(count) + 1
				item := createRandomItem(j.i)
				id := fmt.Sprintf("test-%d", num)
				pkVal := fmt.Sprintf("pk-%d", num)
				item["id"] = id
				item[pkField] = pkVal

				body, err := json.Marshal(item)
				if err != nil {
					log.Printf("randomRW marshal error id=%s pk=%s: %v", id, pkVal, err)
					errs <- err
					continue
				}

				pk := azcosmos.NewPartitionKeyString(pkVal)
				// Include query op (0=upsert,1=read,2=query)
				operationNum := rng.Intn(3)
				switch operationNum {
				case 0: // Upsert
					if _, err = container.UpsertItem(ctx, pk, body, nil); err != nil {
						log.Printf("upsert error id=%s pk=%s: %v", id, pkVal, err)
						errs <- err
						continue
					}
				case 1: // Read
					if _, err = container.ReadItem(ctx, pk, id, nil); err != nil {
						log.Printf("read error id=%s pk=%s: %v", id, pkVal, err)
						errs <- err
						continue
					}
				case 2: // Query by id
					pager := container.NewQueryItemsPager(
						"SELECT * FROM c WHERE c.id = @id",
						azcosmos.NewPartitionKeyString(pkVal),
						&azcosmos.QueryOptions{
							QueryParameters: []azcosmos.QueryParameter{{Name: "@id", Value: id}},
						},
					)
					for pager.More() {
						if _, err = pager.NextPage(ctx); err != nil {
							log.Printf("query error id=%s pk=%s: %v", id, pkVal, err)
							errs <- err
							break
						}
					}
				}
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
