// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package workloads

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// RunWorkload starts up a workload that performs several opertations.
// ctx - The context for the requests
// client - The Cosmos DB client
// cfg - The workload configuration
func RunWorkload(ctx context.Context, client *azcosmos.Client, cfg WorkloadConfig) error {

	dbClient, err := client.NewDatabase(cfg.DatabaseID)
	if err != nil {
		return fmt.Errorf("ensure database: %w", err)
	}

	container, err := dbClient.NewContainer(cfg.ContainerID)
	if err != nil {
		return fmt.Errorf("ensure container: %w", err)
	}

	var count = cfg.LogicalPartitions

	log.Printf("Starting %d concurrent read/write/queries ...", count)

	// Use two goroutines each locked to their own OS thread.
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: random read/write/queries
	go func() {
		// Pin this goroutine to its own OS thread
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := randomReadWriteQueries(ctx, container, count, cfg.PartitionKeyFieldName); err != nil {
				log.Printf("read/write/queries failed: %v", err)
			}

		}
	}()

	// Goroutine 2: vector search queries
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := vectorSearchQueries(ctx, container, count, cfg.PartitionKeyFieldName); err != nil {
				log.Printf("vector search queries failed: %v", err)
			}
		}
	}()

	// Wait until context is cancelled, then wait for goroutines to finish
	<-ctx.Done()
	// Give goroutines a moment to observe ctx.Done and exit; they will return promptly
	wg.Wait()
	return nil

}
