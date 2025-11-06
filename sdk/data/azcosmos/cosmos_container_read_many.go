// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// Executes a query using the provided query engine.
func (c *ContainerClient) executeReadManyWithEngine(queryEngine queryengine.QueryEngine, items []ItemIdentity, readManyOptions *ReadManyOptions, operationContext pipelineRequestOptions, ctx context.Context) (ReadManyItemsResponse, error) {
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)

	// get the partition key ranges for the container
	rawPartitionKeyRanges, err := c.getPartitionKeyRangesRaw(ctx, operationContext)
	if err != nil {
		// if we can't get the partition key ranges, return empty map
		return ReadManyItemsResponse{}, err
	}

	// get the container properties
	containerRsp, err := c.Read(ctx, nil)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}

	// create the item identities for the query engine with json string
	newItemIdentities := make([]queryengine.ItemIdentity, len(items))
	for i := range items {
		pkStr, err := items[i].PartitionKey.toJsonString()
		if err != nil {
			return ReadManyItemsResponse{}, err
		}
		newItemIdentities[i] = queryengine.ItemIdentity{
			PartitionKeyValue: pkStr,
			ID:                items[i].ID,
		}
	}
	var pkVersion int32
	if containerRsp.ContainerProperties.PartitionKeyDefinition.Version == 0 {
		pkVersion = int32(1)
	} else {
		pkVersion = int32(containerRsp.ContainerProperties.PartitionKeyDefinition.Version)
	}

	readManyPipeline, err := queryEngine.CreateReadManyPipeline(rawPartitionKeyRanges, newItemIdentities, string(containerRsp.ContainerProperties.PartitionKeyDefinition.Kind), pkVersion)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}
	log.Writef(EventQueryEngine, "Created readMany pipeline")
	// Initial run to get any requests.
	log.Writef(EventQueryEngine, "Fetching more data from readMany pipeline")
	result, err := readManyPipeline.Run()
	if err != nil {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	}

	concurrency := determineConcurrency(nil)
	if readManyOptions != nil {
		concurrency = determineConcurrency(readManyOptions.MaxConcurrency)
	}
	totalRequestCharge, err := runEngineRequests(ctx, c, path, readManyPipeline, operationContext, result.Requests, concurrency, func(req queryengine.QueryRequest) (string, []QueryParameter, bool) {
		// ReadMany pipeline requests carry a Query (optional override). No parameters and we always page until continuation exhausted.
		return req.Query, nil, true /* treat like drain for full pagination */
	})
	if err != nil {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	}

	// Final run to gather merged items.
	result, err = readManyPipeline.Run()
	if err != nil {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	}

	if readManyPipeline.IsComplete() {
		log.Writef(EventQueryEngine, "ReadMany pipeline is complete")
		readManyPipeline.Close()
		return ReadManyItemsResponse{
			Items:         result.Items,
			RequestCharge: totalRequestCharge,
		}, nil
	} else {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, errors.New("illegal state readMany pipeline did not complete")
	}
}

func (c *ContainerClient) executeReadManyWithPointReads(items []ItemIdentity, readManyOptions *ReadManyOptions, operationContext pipelineRequestOptions, ctx context.Context) (ReadManyItemsResponse, error) {

	// Determine concurrency: use provided MaxConcurrency or number of CPU cores
	concurrency := determineConcurrency(nil)
	if readManyOptions != nil {
		concurrency = determineConcurrency(readManyOptions.MaxConcurrency)
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
