// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package workloads

import (
	"context"
	"fmt"
	"log"
)

func RunWorkload(ctx context.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	client, err := createClient(cfg)
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

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

	for {
		if err := randomReadWriteQueries(ctx, container, count, cfg.PartitionKeyFieldName); err != nil {
			log.Printf("read/write/queries failed: %v", err)
		}
		if err := vectorSearchQueries(ctx, container, count, cfg.PartitionKeyFieldName); err != nil {
			log.Printf("vector search queries failed: %v", err)
		}
	}

}
