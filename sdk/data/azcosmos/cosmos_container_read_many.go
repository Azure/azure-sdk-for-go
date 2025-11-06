// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
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

	readManyPipeline, err := queryEngine.CreateReadManyPipeline(rawPartitionKeyRanges, newItemIdentities, string(containerRsp.ContainerProperties.PartitionKeyDefinition.Kind), int32(containerRsp.ContainerProperties.PartitionKeyDefinition.Version))
	if err != nil {
		return ReadManyItemsResponse{}, err
	}
	log.Writef(EventQueryEngine, "Created readMany pipeline")
	// Fetch more data from the pipeline
	log.Writef(EventQueryEngine, "Fetching more data from readMany pipeline")
	result, err := readManyPipeline.Run()
	if err != nil {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	}

	// we need to make requests for the items in the queue.
	// If there are no requests, the pipeline should return true for IsComplete
	totalRequestCharge := float32(0)

	// Determine concurrency for engine requests (mirror point-reads behavior)
	var concurrency int
	if readManyOptions != nil && readManyOptions.MaxConcurrency != nil && *readManyOptions.MaxConcurrency > 0 {
		concurrency = int(*readManyOptions.MaxConcurrency)
	} else {
		concurrency = runtime.NumCPU()
		if concurrency <= 0 {
			concurrency = 1
		}
	}

	// Channels and synchronization
	jobs := make(chan queryengine.QueryRequest, len(result.Requests))
	provideCh := make(chan []queryengine.QueryResult)
	errCh := make(chan error, 1)
	done := make(chan struct{})
	providerDone := make(chan struct{})
	var wg sync.WaitGroup
	var chargeMu sync.Mutex

	// Provider goroutine: serializes calls into the pipeline
	go func() {
		defer close(providerDone)
		for res := range provideCh {
			if err := readManyPipeline.ProvideData(res); err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
		}
	}()

	// Start workers
	workerCount := concurrency
	if workerCount > len(result.Requests) {
		workerCount = len(result.Requests)
	}
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case req, ok := <-jobs:
					if !ok {
						return
					}

					log.Writef(azlog.EventRequest, "ReadMany pipeline requested data for PKRange: %s", req.PartitionKeyRangeID)
					// paginate
					fetchMorePages := true
					for fetchMorePages {
						// wrap req so it implements the toHeaders method expected by sendQueryRequest
						qr := queryRequest(req)
						azResponse, err := c.database.client.sendQueryRequest(
							path,
							ctx,
							qr.Query,
							nil,
							operationContext,
							&qr,
							nil)
						if err != nil {
							select {
							case errCh <- err:
							default:
							}
							return
						}
						queryResponse, err := newQueryResponse(azResponse)
						if err != nil {
							select {
							case errCh <- err:
							default:
							}
							return
						}

						// update totalRequestCharge
						chargeMu.Lock()
						totalRequestCharge += queryResponse.RequestCharge
						chargeMu.Unlock()

						// Load the data into a buffer to send it to the pipeline
						buf := new(bytes.Buffer)
						_, err = buf.ReadFrom(azResponse.Body)
						if err != nil {
							select {
							case errCh <- err:
							default:
							}
							return
						}
						data := buf.Bytes()
						continuation := azResponse.Header.Get(cosmosHeaderContinuationToken)
						// only end once all pages have been fetched
						fetchMorePages = continuation != ""

						// Provide the data to the pipeline, make sure it's tagged with the partition key range ID so the pipeline can merge it into the correct partition.
						qres := queryengine.QueryResult{
							PartitionKeyRangeID: req.PartitionKeyRangeID,
							NextContinuation:    continuation,
							RequestId:           req.Id,
							Data:                data,
						}
						log.Writef(EventQueryEngine, "Received response for PKRange: %s.", req.PartitionKeyRangeID)
						select {
						case <-done:
							return
						case provideCh <- []queryengine.QueryResult{qres}:
						}
					}
				}
			}
		}()
	}

	// Feed jobs
	go func() {
		for _, r := range result.Requests {
			select {
			case <-done:
				break
			default:
			}
			jobs <- queryengine.QueryRequest(r)
		}
		close(jobs)
	}()

	// Wait for workers to finish then close provideCh (so provider can exit)
	go func() {
		wg.Wait()
		close(provideCh)
	}()

	// Wait for provider to finish or for an error/context cancellation
	select {
	case err := <-errCh:
		// signal cancellation
		select {
		case <-done:
		default:
			close(done)
		}
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	case <-ctx.Done():
		// cancel and return
		select {
		case <-done:
		default:
			close(done)
		}
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, ctx.Err()
	case <-providerDone:
		// provider finished normally
	}

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
