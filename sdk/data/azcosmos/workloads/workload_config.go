// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package workloads

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// WorkloadConfig is the configuration loaded from environment (mirrors the Python version)
type WorkloadConfig struct {
	Endpoint              string
	Key                   string
	PreferredLocations    []string
	DatabaseID            string
	ContainerID           string
	PartitionKeyFieldName string
	LogicalPartitions     int
	Throughput            int // optional (unused if not supported)
}

const defaultLogicalPartitions = 10000
const defaultThroughput = 100000
const defaultContainerName = "scale_cont"
const defaultDatabaseName = "scale_db"
const defaultPKField = "pk"

func LoadConfig() (WorkloadConfig, error) {
	get := func(name string) (string, error) {
		v := os.Getenv(name)
		if v == "" {
			return "", fmt.Errorf("missing env var %s", name)
		}
		return v, nil
	}

	var cfg WorkloadConfig
	var err error

	if cfg.Endpoint, err = get("COSMOS_URI"); err != nil {
		return cfg, err
	}
	if key := os.Getenv("COSMOS_KEY"); key != "" {
		cfg.Key = key
	}
	if cosmosDatabase := os.Getenv("COSMOS_DATABASE"); cosmosDatabase != "" {
		cfg.DatabaseID = cosmosDatabase
	} else {
		cfg.DatabaseID = defaultDatabaseName
	}
	if cosmosContainer := os.Getenv("COSMOS_CONTAINER"); cosmosContainer != "" {
		cfg.ContainerID = cosmosContainer
	} else {
		cfg.ContainerID = defaultContainerName
	}
	if pk := os.Getenv("PARTITION_KEY"); pk != "" {
		cfg.PartitionKeyFieldName = pk
	} else {
		cfg.PartitionKeyFieldName = defaultPKField
	}

	if lp := os.Getenv("NUMBER_OF_LOGICAL_PARTITIONS"); lp != "" {
		n, convErr := strconv.Atoi(lp)
		if convErr != nil {
			return cfg, fmt.Errorf("invalid NUMBER_OF_LOGICAL_PARTITIONS: %w", convErr)
		}
		cfg.LogicalPartitions = n
	} else {
		cfg.LogicalPartitions = defaultLogicalPartitions
	}

	if tp := os.Getenv("THROUGHPUT"); tp != "" {
		n, convErr := strconv.Atoi(tp)
		if convErr != nil {
			return cfg, fmt.Errorf("invalid THROUGHPUT: %w", convErr)
		}
		cfg.Throughput = n
	} else {
		cfg.Throughput = defaultThroughput
	}

	// Comma-separated preferred locations (optional)
	if pl := os.Getenv("PREFERRED_LOCATIONS"); pl != "" {
		// Simple split on comma;
		cfg.PreferredLocations = strings.Split(pl, ",")
	}
	return cfg, nil
}
