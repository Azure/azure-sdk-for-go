// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
)

// executeReadManyWithEngine executes a query using the provided query engine.
func (c *ContainerClient) executeReadManyWithEngine(queryEngine queryengine.QueryEngine, items []ItemIdentity, readManyOptions *ReadManyOptions, operationContext pipelineRequestOptions, ctx context.Context) (ReadManyItemsResponse, error) {
	// throw error that this is unsupported
	return ReadManyItemsResponse{}, errors.New("ReadMany with query engine is not supported yet")
}

func (c *ContainerClient) executeReadManyWithPointReads(items []ItemIdentity, readManyOptions *ReadManyOptions, operationContext pipelineRequestOptions, ctx context.Context) (ReadManyItemsResponse, error) {

	// if empty list of items, return empty list
	if len(items) == 0 {
		return ReadManyItemsResponse{}, nil
	}
	// Determine concurrency: use provided MaxConcurrency or number of CPU cores
	var concurrency int
	if readManyOptions != nil && readManyOptions.MaxConcurrency != nil && *readManyOptions.MaxConcurrency > 0 {
		concurrency = int(*readManyOptions.MaxConcurrency)
	} else {
		concurrency = runtime.NumCPU()
		if concurrency <= 0 {
			concurrency = 1
		}
	}

	// Prepare result slots to preserve input order
	type slot struct {
		value         []byte
		requestCharge float32
		err           error
	}

	results := make([]slot, len(items))

	// Worker pool
	var wg sync.WaitGroup
	jobs := make(chan int)

	// cancellation channel to short-circuit on first error
	done := make(chan struct{})

	// Start workers
	workerCount := concurrency
	if workerCount > len(items) {
		workerCount = len(items)
	}
	itemOptions := ItemOptions{}
	if readManyOptions != nil {
		itemOptions.ConsistencyLevel = readManyOptions.ConsistencyLevel
		itemOptions.SessionToken = readManyOptions.SessionToken
	}
	for worker := 0; worker < workerCount; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				select {
				case <-done:
					return
				default:
				}
				item := items[idx]

				itemResponse, err := c.ReadItem(ctx, item.PartitionKey, item.ID, &itemOptions)
				if err != nil {
					var azErr *azcore.ResponseError
					// for 404, just continue without error
					if errors.As(err, &azErr) {
						if azErr.StatusCode == 404 {
							continue
						}
					}
					results[idx].err = err
					// signal cancellation
					select {
					case <-done:
					default:
						close(done)
					}
					// store error and continue to allow workers to exit
					return
				}
				results[idx].value = itemResponse.Value
				results[idx].requestCharge = itemResponse.RequestCharge
			}
		}()
	}

	// Start a goroutine to distribute item indices to the worker pool via the jobs channel.
	go func() {
		for i := range items {
			select {
			case <-done:
				return
			default:
			}
			jobs <- i
		}
		close(jobs)
	}()

	wg.Wait()

	// Check for errors and build response in original order
	var readManyResponse ReadManyItemsResponse
	for i := range results {
		if results[i].err != nil {
			return ReadManyItemsResponse{}, results[i].err
		}
		if results[i].value != nil {
			readManyResponse.Items = append(readManyResponse.Items, results[i].value)
			readManyResponse.RequestCharge += results[i].requestCharge
		}
	}

	return readManyResponse, nil
}
