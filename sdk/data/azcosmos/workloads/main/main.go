// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/workloads"
)

func main() {
	// setup logger
	f, err := os.OpenFile("workloads.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()

	// Send the default logger output to the file
	log.SetOutput(f)
	// Optional: set flags to include date/time/file info
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	ctx := context.Background()

	cfg, err := workloads.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client, err := workloads.CreateClient(cfg)
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}

	if cfg.Key != "" {
		if err := workloads.RunSetup(ctx, client, cfg); err != nil {
			log.Fatalf("setup failed: %v", err)
		}
	} else {
		log.Printf("Setup skipped as AAD is enabled.")
	}
	log.Println("setup completed")
	if err := workloads.RunWorkload(ctx, client, cfg); err != nil {
		log.Fatalf("workload failed: %v", err)
	}
}
