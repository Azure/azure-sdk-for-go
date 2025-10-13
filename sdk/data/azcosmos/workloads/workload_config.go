package workloads

import (
	"fmt"
	"os"
	"strconv"
)

// Configuration loaded from environment (mirrors the Python version)
type workloadConfig struct {
	Endpoint              string
	Key                   string
	PreferredLocations    []string
	DatabaseID            string
	ContainerID           string
	PartitionKeyFieldName string
	LogicalPartitions     int
	Throughput            int // optional (unused if not supported)
}

func loadConfig() (workloadConfig, error) {
	get := func(name string) (string, error) {
		v := os.Getenv(name)
		if v == "" {
			return "", fmt.Errorf("missing env var %s", name)
		}
		return v, nil
	}

	var cfg workloadConfig
	var err error

	if cfg.Endpoint, err = get("COSMOS_URI"); err != nil {
		return cfg, err
	}
	if cfg.Key, err = get("COSMOS_KEY"); err != nil {
		return cfg, err
	}
	if cfg.DatabaseID, err = get("COSMOS_DATABASE"); err != nil {
		return cfg, err
	}
	if cfg.ContainerID, err = get("COSMOS_CONTAINER"); err != nil {
		return cfg, err
	}
	if pk := os.Getenv("PARTITION_KEY"); pk != "" {
		cfg.PartitionKeyFieldName = pk
	} else {
		cfg.PartitionKeyFieldName = "pk"
	}

	if lp := os.Getenv("NUMBER_OF_LOGICAL_PARTITIONS"); lp != "" {
		n, convErr := strconv.Atoi(lp)
		if convErr != nil {
			return cfg, fmt.Errorf("invalid NUMBER_OF_LOGICAL_PARTITIONS: %w", convErr)
		}
		cfg.LogicalPartitions = n
	} else {
		cfg.LogicalPartitions = 10000
	}

	if tp := os.Getenv("THROUGHPUT"); tp != "" {
		n, convErr := strconv.Atoi(tp)
		if convErr != nil {
			return cfg, fmt.Errorf("invalid THROUGHPUT: %w", convErr)
		}
		cfg.Throughput = n
	} else {
		cfg.Throughput = 10000
	}

	// Comma-separated preferred locations (optional)
	if pl := os.Getenv("PREFERRED_LOCATIONS"); pl != "" {
		// Simple split on comma; whitespace trimming omitted for brevity
		cfg.PreferredLocations = splitAndTrim(pl, ',')
	}
	return cfg, nil
}

func splitAndTrim(s string, sep rune) []string {
	if s == "" {
		return nil
	}
	out := []string{}
	cur := ""
	for _, r := range s {
		if r == sep {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
