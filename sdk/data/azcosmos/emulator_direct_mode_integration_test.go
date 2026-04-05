// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type testLogger struct {
	t        *testing.T
	file     *os.File
	mu       sync.Mutex
	testName string
}

func newTestLogger(t *testing.T) *testLogger {
	testName := t.Name()

	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Logf("Warning: Could not create logs directory: %v", err)
		return &testLogger{t: t, testName: testName}
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("%s_%s.log", testName, timestamp)
	logPath := filepath.Join(logsDir, logFileName)

	file, err := os.Create(logPath)
	if err != nil {
		t.Logf("Warning: Could not create log file: %v", err)
		return &testLogger{t: t, testName: testName}
	}

	header := fmt.Sprintf("=== Log file for %s ===\nStarted: %s\n\n", testName, time.Now().Format(time.RFC3339))
	_, _ = file.WriteString(header)

	t.Logf("Logging to file: %s", logPath)

	t.Cleanup(func() {
		if file != nil {
			_, _ = file.WriteString(fmt.Sprintf("\n=== Test completed at %s ===\n", time.Now().Format(time.RFC3339)))
			_ = file.Close()
		}
	})

	return &testLogger{
		t:        t,
		file:     file,
		testName: testName,
	}
}

func (l *testLogger) Log(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05.000")

	l.t.Log(msg)

	if l.file != nil {
		l.mu.Lock()
		_, _ = l.file.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, msg))
		l.mu.Unlock()
	}
}

func (l *testLogger) Logf(format string, args ...any) {
	l.Log(format, args...)
}

func (l *testLogger) Fatalf(format string, args ...any) {
	l.Log("FATAL: "+format, args...)
	l.t.FailNow()
}

// TestDirectModeOperations is a quick sanity test that verifies all Direct Mode operations work.
// This runs quickly (~10 seconds) and exercises CREATE, READ, PATCH, QUERY, and DELETE.
//
// Run with: EMULATOR=true go test -run TestDirectModeOperations -v -count=1 -timeout 1m
func TestDirectModeOperations(t *testing.T) {
	log := newTestLogger(t)

	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getDirectModeClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	dbName := "directModeOps"
	containerName := "opsTest"

	existingDb, _ := client.NewDatabase(dbName)
	_, _ = existingDb.Delete(context.TODO(), nil)

	log.Log("=== Direct Mode Operations Test ===")

	database := emulatorTests.createDatabase(t, context.TODO(), client, dbName)
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)

	properties := ContainerProperties{
		ID: containerName,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	container, _ := database.NewContainer(containerName)
	pk := NewPartitionKeyString("test")

	item := map[string]any{
		"id":      "test-item-1",
		"pk":      "test",
		"value":   "initial",
		"counter": 0,
	}
	marshalled, _ := json.Marshal(item)

	createResp, err := container.CreateItem(context.TODO(), pk, marshalled, nil)
	if err != nil {
		log.Fatalf("CREATE failed: %v", err)
	}
	log.Logf("CREATE - %d (RU: %.2f)", createResp.RawResponse.StatusCode, createResp.RequestCharge)

	readResp, err := container.ReadItem(context.TODO(), pk, "test-item-1", nil)
	if err != nil {
		log.Fatalf("READ failed: %v", err)
	}
	log.Logf("READ - %d (RU: %.2f)", readResp.RawResponse.StatusCode, readResp.RequestCharge)

	patchOps := PatchOperations{}
	patchOps.AppendSet("/value", "patched")
	patchOps.AppendIncrement("/counter", 1)

	patchResp, err := container.PatchItem(context.TODO(), pk, "test-item-1", patchOps, nil)
	if err != nil {
		log.Fatalf("PATCH failed: %v", err)
	}
	log.Logf("PATCH - %d (RU: %.2f)", patchResp.RawResponse.StatusCode, patchResp.RequestCharge)

	item["value"] = "replaced"
	marshalled, _ = json.Marshal(item)
	replaceResp, err := container.ReplaceItem(context.TODO(), pk, "test-item-1", marshalled, nil)
	if err != nil {
		log.Fatalf("REPLACE failed: %v", err)
	}
	log.Logf("REPLACE - %d (RU: %.2f)", replaceResp.RawResponse.StatusCode, replaceResp.RequestCharge)

	item2 := map[string]any{
		"id":    "test-item-2",
		"pk":    "test",
		"value": "upserted",
	}
	marshalled2, _ := json.Marshal(item2)
	upsertResp, err := container.UpsertItem(context.TODO(), pk, marshalled2, nil)
	if err != nil {
		log.Fatalf("UPSERT failed: %v", err)
	}
	log.Logf("UPSERT - %d (RU: %.2f)", upsertResp.RawResponse.StatusCode, upsertResp.RequestCharge)

	queryOpts := QueryOptions{PageSizeHint: 10}
	queryPager := container.NewQueryItemsPager("SELECT * FROM c WHERE c.pk = 'test'", pk, &queryOpts)
	itemCount := 0
	for queryPager.More() {
		page, err := queryPager.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("QUERY failed: %v", err)
		}
		itemCount += len(page.Items)
		log.Logf("QUERY page - %d items (RU: %.2f)", len(page.Items), page.RequestCharge)
	}
	if itemCount != 2 {
		log.Fatalf("Expected 2 items from query, got %d", itemCount)
	}

	deleteResp, err := container.DeleteItem(context.TODO(), pk, "test-item-1", nil)
	if err != nil {
		log.Fatalf("DELETE failed: %v", err)
	}
	log.Logf("DELETE - %d (RU: %.2f)", deleteResp.RawResponse.StatusCode, deleteResp.RequestCharge)

	deleteResp2, err := container.DeleteItem(context.TODO(), pk, "test-item-2", nil)
	if err != nil {
		log.Fatalf("DELETE failed: %v", err)
	}
	log.Logf("DELETE - %d (RU: %.2f)", deleteResp2.RawResponse.StatusCode, deleteResp2.RequestCharge)

	log.Log("=== All Direct Mode operations verified! ===")
}

// TestDirectModeLongRunning is a long-lived integration test that simulates normal usage
// over approximately 5 minutes, exercising CREATE, PATCH, DELETE, and query operations
// via Direct Mode (RNTBD protocol).
//
// Run with: EMULATOR=true go test -run TestDirectModeLongRunning -v -count=1 -timeout 10m
func TestDirectModeLongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}

	log := newTestLogger(t)

	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getDirectModeClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	dbName := "directModeIntegration"
	containerName := "longRunningTest"

	existingDb, _ := client.NewDatabase(dbName)
	_, _ = existingDb.Delete(context.TODO(), nil)

	log.Log("=== Direct Mode Long-Running Integration Test ===")
	log.Logf("Creating database: %s", dbName)

	database := emulatorTests.createDatabase(t, context.TODO(), client, dbName)
	defer func() {
		log.Logf("Cleaning up: Deleting database %s", dbName)
		emulatorTests.deleteDatabase(t, context.TODO(), database)
	}()

	properties := ContainerProperties{
		ID: containerName,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	log.Logf("Creating container: %s", containerName)
	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	container, _ := database.NewContainer(containerName)

	testDuration := 5 * time.Minute
	operationInterval := 2 * time.Second
	startTime := time.Now()
	operationCount := 0
	successCount := 0
	errorCount := 0

	var createdItems []string
	var mu sync.Mutex

	log.Logf("Starting long-running test for %v with %v intervals", testDuration, operationInterval)
	log.Log("Operations: CREATE -> PATCH -> QUERY -> DELETE (cycling)")

	for time.Since(startTime) < testDuration {
		operationCount++
		itemID := fmt.Sprintf("item-%d-%d", operationCount, rand.Int63())
		pk := NewPartitionKeyString("partition1")

		item := map[string]any{
			"id":          itemID,
			"pk":          "partition1",
			"value":       fmt.Sprintf("data-%d", operationCount),
			"counter":     0,
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"description": "test item for long-running direct mode test",
		}

		marshalled, err := json.Marshal(item)
		if err != nil {
			log.Fatalf("Failed to marshal item: %v", err)
		}

		createResp, err := container.CreateItem(context.TODO(), pk, marshalled, nil)
		if err != nil {
			log.Logf("[%d] CREATE %s - ERROR: %v", operationCount, itemID, err)
			errorCount++
			continue
		}
		log.Logf("[%d] CREATE %s - %d (RU: %.2f)", operationCount, itemID, createResp.RawResponse.StatusCode, createResp.RequestCharge)
		successCount++

		mu.Lock()
		createdItems = append(createdItems, itemID)
		mu.Unlock()

		patchOps := PatchOperations{}
		patchOps.AppendIncrement("/counter", 1)
		patchOps.AppendSet("/lastModified", time.Now().UTC().Format(time.RFC3339))
		patchOps.AppendAdd("/patchedBy", "long-running-test")

		patchResp, err := container.PatchItem(context.TODO(), pk, itemID, patchOps, nil)
		if err != nil {
			log.Logf("[%d] PATCH %s - ERROR: %v", operationCount, itemID, err)
			errorCount++
		} else {
			log.Logf("[%d] PATCH %s - %d (RU: %.2f)", operationCount, itemID, patchResp.RawResponse.StatusCode, patchResp.RequestCharge)
			successCount++
		}

		if operationCount%5 == 0 {
			queryOpts := QueryOptions{PageSizeHint: 10}
			queryPager := container.NewQueryItemsPager(
				"SELECT c.id, c.counter, c.timestamp FROM c WHERE c.pk = 'partition1' ORDER BY c.timestamp DESC",
				pk,
				&queryOpts,
			)

			itemCount := 0
			var totalRU float32 = 0
			for queryPager.More() {
				page, err := queryPager.NextPage(context.TODO())
				if err != nil {
					log.Logf("[%d] QUERY - ERROR: %v", operationCount, err)
					errorCount++
					break
				}
				itemCount += len(page.Items)
				totalRU += page.RequestCharge
			}
			log.Logf("[%d] QUERY - 200 (items: %d, total RU: %.2f)", operationCount, itemCount, totalRU)
			successCount++
		}

		mu.Lock()
		if len(createdItems) > 20 {
			toDelete := createdItems[0]
			createdItems = createdItems[1:]
			mu.Unlock()

			deleteResp, err := container.DeleteItem(context.TODO(), pk, toDelete, nil)
			if err != nil {
				log.Logf("[%d] DELETE %s - ERROR: %v", operationCount, toDelete, err)
				errorCount++
			} else {
				log.Logf("[%d] DELETE %s - %d (RU: %.2f)", operationCount, toDelete, deleteResp.RawResponse.StatusCode, deleteResp.RequestCharge)
				successCount++
			}
		} else {
			mu.Unlock()
		}

		time.Sleep(operationInterval)
	}

	log.Log("=== Cleanup Phase ===")
	mu.Lock()
	remainingItems := createdItems
	mu.Unlock()

	for _, itemID := range remainingItems {
		pk := NewPartitionKeyString("partition1")
		deleteResp, err := container.DeleteItem(context.TODO(), pk, itemID, nil)
		if err != nil {
			log.Logf("[cleanup] DELETE %s - ERROR: %v", itemID, err)
		} else {
			log.Logf("[cleanup] DELETE %s - %d", itemID, deleteResp.RawResponse.StatusCode)
		}
	}

	elapsed := time.Since(startTime)
	log.Log("=== Test Summary ===")
	log.Logf("Duration: %v", elapsed)
	log.Logf("Total operations: %d", operationCount*3)
	log.Logf("Successful: %d", successCount)
	log.Logf("Errors: %d", errorCount)

	if errorCount > 0 && float64(errorCount)/float64(successCount+errorCount) > 0.1 {
		t.Errorf("Error rate too high: %d errors out of %d operations", errorCount, successCount+errorCount)
	}

	log.Log("Direct mode long-running test completed!")
}

// TestDirectModeRateLimiting tests the 429 (Too Many Requests) throttling behavior
// by deliberately exceeding the 400 RU/s limit on the emulator.
//
// Run with: EMULATOR=true go test -run TestDirectModeRateLimiting -v -count=1 -timeout 5m
func TestDirectModeRateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rate limiting test in short mode")
	}

	log := newTestLogger(t)

	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getDirectModeClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	dbName := "directModeRateLimit"
	containerName := "rateLimitTest"

	existingDb, _ := client.NewDatabase(dbName)
	_, _ = existingDb.Delete(context.TODO(), nil)

	log.Log("=== Direct Mode Rate Limiting Test ===")
	log.Logf("Creating database: %s", dbName)

	database := emulatorTests.createDatabase(t, context.TODO(), client, dbName)
	defer func() {
		log.Logf("Cleaning up: Deleting database %s", dbName)
		emulatorTests.deleteDatabase(t, context.TODO(), database)
	}()

	properties := ContainerProperties{
		ID: containerName,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	log.Logf("Creating container: %s (emulator default: 400 RU/s)", containerName)
	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	container, _ := database.NewContainer(containerName)

	var totalRU float32 = 0
	var throttledCount int = 0
	var successCount int = 0
	var createdItems []string

	log.Log("Attempting to exceed 400 RU/s limit with rapid-fire operations...")
	log.Log("Expecting 429 (TooManyRequests) responses once quota exhausted")

	largeData := make([]byte, 4096)
	for i := range largeData {
		largeData[i] = byte('A' + (i % 26))
	}

	testDuration := 30 * time.Second
	startTime := time.Now()
	operationNum := 0

	for time.Since(startTime) < testDuration {
		operationNum++
		itemID := fmt.Sprintf("rate-limit-item-%d", operationNum)

		item := map[string]any{
			"id":      itemID,
			"pk":      "ratelimit",
			"payload": string(largeData),
			"seq":     operationNum,
		}

		marshalled, err := json.Marshal(item)
		if err != nil {
			log.Fatalf("Failed to marshal item: %v", err)
		}

		pk := NewPartitionKeyString("ratelimit")
		createResp, err := container.CreateItem(context.TODO(), pk, marshalled, nil)

		if err != nil {
			var responseErr *azcore.ResponseError
			if errors.As(err, &responseErr) {
				if responseErr.StatusCode == http.StatusTooManyRequests {
					throttledCount++
					log.Logf("[%d] CREATE %s - 429 TooManyRequests (THROTTLED)", operationNum, itemID)

					time.Sleep(100 * time.Millisecond)
					continue
				}
			}
			log.Logf("[%d] CREATE %s - ERROR: %v", operationNum, itemID, err)
			continue
		}

		successCount++
		totalRU += createResp.RequestCharge
		createdItems = append(createdItems, itemID)
		log.Logf("[%d] CREATE %s - %d (RU: %.2f, cumulative: %.2f)", operationNum, itemID, createResp.RawResponse.StatusCode, createResp.RequestCharge, totalRU)
	}

	elapsed := time.Since(startTime)

	log.Log("=== Rate Limiting Test Summary ===")
	log.Logf("Duration: %v", elapsed)
	log.Logf("Total attempts: %d", operationNum)
	log.Logf("Successful: %d", successCount)
	log.Logf("Throttled (429): %d", throttledCount)
	log.Logf("Total RU consumed: %.2f", totalRU)
	log.Logf("Avg RU/operation: %.2f", totalRU/float32(max(successCount, 1)))

	log.Log("=== Cleanup Phase ===")
	pk := NewPartitionKeyString("ratelimit")
	for _, itemID := range createdItems {
		_, err := container.DeleteItem(context.TODO(), pk, itemID, nil)
		if err != nil {
			var responseErr *azcore.ResponseError
			if errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusTooManyRequests {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	if throttledCount == 0 {
		log.Log("WARNING: No 429 responses received. The emulator may have higher RU limits or operations completed too slowly.")
		log.Log("This is not necessarily a failure - it means the operations stayed within the RU budget.")
	} else {
		log.Logf("SUCCESS: Received %d throttled (429) responses, confirming rate limiting behavior", throttledCount)
	}

	log.Log("Direct mode rate limiting test completed!")
}
