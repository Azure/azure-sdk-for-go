// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package setup creates Cosmos DB databases and containers for perf runs.
package setup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

const (
	maxRetries     = 10
	initialBackoff = time.Second
	maxBackoff     = 30 * time.Second
)

// EnsureDatabase ensures a database exists and is readable.
func EnsureDatabase(ctx context.Context, client *azcosmos.Client, name string) error {
	db, err := client.NewDatabase(name)
	if err != nil {
		return err
	}

	_, err = db.Read(ctx, nil)
	if err == nil {
		fmt.Printf("Database '%s' already exists.\n", name)
		return nil
	}
	if !hasStatus(err, http.StatusNotFound) {
		return err
	}
	fmt.Printf("Database '%s' not found, creating...\n", name)

	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: name}, nil)
	if err == nil {
		fmt.Printf("Database '%s' created.\n", name)
	} else if hasStatus(err, http.StatusConflict) {
		fmt.Printf("Database '%s' was created concurrently.\n", name)
	} else {
		return err
	}

	backoff := initialBackoff
	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err = db.Read(ctx, nil)
		if err == nil {
			fmt.Printf("Database '%s' confirmed readable.\n", name)
			return nil
		}
		if !hasStatus(err, http.StatusNotFound) {
			return err
		}
		fmt.Printf("Database not yet visible (attempt %d/%d), retrying in %s...\n", attempt, maxRetries, backoff)
		if err := sleep(ctx, backoff); err != nil {
			return err
		}
		backoff = minDuration(backoff*2, maxBackoff)
	}
	return fmt.Errorf("database '%s' not readable after %d retries", name, maxRetries)
}

// EnsureContainer ensures a container exists and is readable.
func EnsureContainer(ctx context.Context, db *azcosmos.DatabaseClient, name string, throughput int32, defaultTTL *int32) (*azcosmos.ContainerClient, error) {
	container, err := db.NewContainer(name)
	if err != nil {
		return nil, err
	}

	_, err = container.Read(ctx, nil)
	if err == nil {
		fmt.Printf("Container '%s' already exists.\n", name)
		return container, nil
	}
	if !hasStatus(err, http.StatusNotFound) {
		return nil, err
	}
	fmt.Printf("Container '%s' not found, creating with %d RU/s...\n", name, throughput)

	props := azcosmos.ContainerProperties{
		ID: name,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/partition_key"},
		},
		DefaultTimeToLive: defaultTTL,
	}
	throughputProps := azcosmos.NewManualThroughputProperties(throughput)
	_, err = db.CreateContainer(ctx, props, &azcosmos.CreateContainerOptions{ThroughputProperties: &throughputProps})
	if err == nil {
		fmt.Printf("Container '%s' created.\n", name)
	} else if hasStatus(err, http.StatusConflict) {
		fmt.Printf("Container '%s' was created concurrently.\n", name)
	} else {
		return nil, err
	}

	backoff := initialBackoff
	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err = container.Read(ctx, nil)
		if err == nil {
			fmt.Printf("Container '%s' confirmed readable.\n", name)
			return container, nil
		}
		if !hasStatus(err, http.StatusNotFound) {
			return nil, err
		}
		fmt.Printf("Container not yet visible (attempt %d/%d), retrying in %s...\n", attempt, maxRetries, backoff)
		if err := sleep(ctx, backoff); err != nil {
			return nil, err
		}
		backoff = minDuration(backoff*2, maxBackoff)
	}
	return nil, fmt.Errorf("container '%s' not readable after %d retries", name, maxRetries)
}

func hasStatus(err error, status int) bool {
	var responseErr *azcore.ResponseError
	return errors.As(err, &responseErr) && responseErr.StatusCode == status
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
