// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmulatorCursorConcurrentCorrelation(t *testing.T) {
	container := emulatorContainer(t)
	run := uniqueItemID(t)
	const workers = 16
	queries := make([]Query, workers)
	for worker := range workers {
		for value := range 3 {
			id := fmt.Sprintf("%s-%d-%d", run, worker, value)
			key := NewPartitionKeyString(id)
			body, err := json.Marshal(map[string]any{"id": id, "pk": id, "run": run, "worker": worker, "value": worker*10 + value})
			require.NoError(t, err)
			_, err = container.CreateItem(t.Context(), key, id, body, nil)
			require.NoError(t, err)
			trackEmulatorItem(t, container, key, id)
		}
		query, err := NewQuery("SELECT VALUE c.value FROM c WHERE c.run = @run AND c.worker = @worker ORDER BY c.value").
			WithParameter("@run", run)
		require.NoError(t, err)
		queries[worker], err = query.WithParameter("@worker", worker)
		require.NoError(t, err)
	}
	start := make(chan struct{})
	results := make(chan error, workers)
	var wg sync.WaitGroup
	for worker := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			pager := container.NewQueryItemsPager(queries[worker], NewFeedScopeForFullContainer(),
				&QueryOptions{Feed: FeedOptions{PageSizeHint: 1}})
			defer pager.Close()
			var values []int
			for pages := 0; pager.More(); pages++ {
				if pages > 12 {
					results <- fmt.Errorf("worker %d did not exhaust", worker)
					return
				}
				page, err := pager.NextPage(t.Context())
				if err != nil {
					results <- fmt.Errorf("worker %d: %w", worker, err)
					return
				}
				for _, item := range page.Items {
					var value int
					if err := json.Unmarshal(item, &value); err != nil {
						results <- err
						return
					}
					values = append(values, value)
				}
				if pager.More() && pages == 0 {
					if _, err := pager.ContinuationToken(t.Context()); err != nil {
						results <- fmt.Errorf("worker %d checkpoint: %w", worker, err)
						return
					}
				}
			}
			if len(values) != 3 || values[0] != worker*10 || values[1] != worker*10+1 || values[2] != worker*10+2 {
				results <- fmt.Errorf("worker %d received another cursor's results: %v", worker, values)
				return
			}
			results <- nil
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	for err := range results {
		require.NoError(t, err)
	}
}
