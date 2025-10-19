// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/workloads"
)

func main() {
	// max run the workload for 12 hours
	ctx, cancel := context.WithTimeout(context.Background(), 12*60*time.Minute)
	defer cancel()

	cfg, err := workloads.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client, err := workloads.CreateClient(cfg)
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}

	if err := workloads.RunSetup(ctx, client, cfg); err != nil {
		log.Fatalf("setup failed: %v", err)
	}
	log.Println("setup completed")
	if err := workloads.RunWorkload(ctx, client, cfg); err != nil {
		log.Fatalf("workload failed: %v", err)
	}
}
