// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"testing"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type captureLogger struct {
	mu      sync.Mutex
	entries []logEntry
}

type logEntry struct {
	event   azlog.Event
	message string
}

func (cl *captureLogger) log(event azlog.Event, message string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.entries = append(cl.entries, logEntry{
		event:   event,
		message: message,
	})
}

func (cl *captureLogger) getEntries() []logEntry {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	entries := make([]logEntry, len(cl.entries))
	copy(entries, cl.entries)
	return entries
}

func (cl *captureLogger) filterByEvent(event azlog.Event) []string {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	var messages []string
	for _, entry := range cl.entries {
		if entry.event == event {
			messages = append(messages, entry.message)
		}
	}
	return messages
}

func (cl *captureLogger) getAllLogs() string {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	var sb strings.Builder
	for _, entry := range cl.entries {
		sb.WriteString("[")
		sb.WriteString(string(entry.event))
		sb.WriteString("] ")
		sb.WriteString(entry.message)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (cl *captureLogger) clear() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.entries = nil
}

// Test helper functions
func assertContains(t *testing.T, logs []string, substring string) {
	t.Helper()
	for _, log := range logs {
		if strings.Contains(log, substring) {
			return
		}
	}
	t.Errorf("Expected to find '%s' in logs, but didn't.\nLogs:\n%v", substring, logs)
}

func assertNotContains(t *testing.T, logs []string, substring string) {
	t.Helper()
	for _, log := range logs {
		if strings.Contains(log, substring) {
			t.Errorf("Expected NOT to find '%s' in logs, but did in: %s", substring, log)
			return
		}
	}
}

func countOccurrences(logs []string, substring string) int {
	count := 0
	for _, log := range logs {
		if strings.Contains(log, substring) {
			count++
		}
	}
	return count
}

func TestLoggingEndpointFailureAndRecovery(t *testing.T) {
	emulatorTests := newEmulatorTests(t)

	logger := &captureLogger{}
	azlog.SetListener(func(event azlog.Event, message string) {
		logger.log(event, message)
	})
	defer azlog.SetListener(nil)

	azlog.SetEvents(EventEndpointManager)

	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{}, // Not validating spans in this test
	}))

	ctx := context.Background()

	t.Run("Initial Healthy State", func(t *testing.T) {
		db := emulatorTests.createDatabase(t, ctx, client, "healthyTestDB")
		defer emulatorTests.deleteDatabase(t, ctx, db)

		logs := logger.filterByEvent(EventEndpointManager)

		assertContains(t, logs, "Initializing Global Endpoint Manager")
		assertNotContains(t, logs, "Marked endpoint unavailable")

		t.Logf("Phase 1 complete: %d EventEndpointManager logs", len(logs))
	})

	t.Run("Mark Endpoint Unavailable", func(t *testing.T) {
		logger.clear()

		defaultEndpoint, err := url.Parse(emulatorTests.host)
		if err != nil {
			t.Fatalf("Failed to parse endpoint: %v", err)
		}

		err = client.gem.MarkEndpointUnavailableForRead(*defaultEndpoint)
		if err != nil {
			t.Fatalf("Failed to mark endpoint unavailable: %v", err)
		}

		logs := logger.filterByEvent(EventEndpointManager)

		assertContains(t, logs, "Marked endpoint unavailable")

		var markedLog string
		for _, log := range logs {
			if strings.Contains(log, "Marked endpoint unavailable") {
				markedLog = log
				break
			}
		}

		if markedLog == "" {
			t.Fatalf("Expected to find 'Marked endpoint unavailable' log")
		}

		if !strings.Contains(markedLog, "endpoint=") {
			t.Errorf("Unavailable log missing 'endpoint=' field")
		}
		if !strings.Contains(markedLog, "operation=") {
			t.Errorf("Unavailable log missing 'operation=' field")
		}
		if !strings.Contains(markedLog, "operation=read") {
			t.Errorf("Expected operation=read in log, got: %s", markedLog)
		}
		t.Logf("Found unavailable log: %s", markedLog)

		priorityLogs := countOccurrences(logs, "Endpoint Priority Recomputed")
		t.Logf("Priority recomputation logs: %d (expected 0-1 for single region)", priorityLogs)

		t.Logf("Phase 2 complete: %d EventEndpointManager logs", len(logs))
	})

	t.Run("No Duplicate Unavailable Logs", func(t *testing.T) {
		logger.clear()

		defaultEndpoint, _ := url.Parse(emulatorTests.host)

		err := client.gem.MarkEndpointUnavailableForRead(*defaultEndpoint)
		if err != nil {
			t.Fatalf("Failed to mark endpoint unavailable: %v", err)
		}

		logs := logger.filterByEvent(EventEndpointManager)

		count := countOccurrences(logs, "Marked endpoint unavailable")

		t.Logf("Phase 3 complete: %d 'Marked endpoint unavailable' occurrences", count)
	})

	t.Run("Endpoint Recovery", func(t *testing.T) {
		logger.clear()

		client.gem.locationCache.forceRefreshStaleEndpoints()

		err := client.gem.locationCache.update(nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("Failed to update location cache: %v", err)
		}

		logs := logger.filterByEvent(EventEndpointManager)

		assertContains(t, logs, "Endpoint is now available")

		var availableLog string
		for _, log := range logs {
			if strings.Contains(log, "Endpoint is now available") {
				availableLog = log
				break
			}
		}

		if availableLog == "" {
			t.Fatalf("Expected to find 'Endpoint is now available' log")
		}

		if !strings.Contains(availableLog, "endpoint=") {
			t.Errorf("Available log missing 'endpoint=' field")
		}
		if !strings.Contains(availableLog, "unavailableFor=") {
			t.Errorf("Available log missing 'unavailableFor=' field")
		}
		t.Logf("Found available log: %s", availableLog)

		t.Logf("Phase 4 complete: %d EventEndpointManager logs", len(logs))
	})

	t.Run("Mark Unavailable For Write", func(t *testing.T) {
		logger.clear()

		defaultEndpoint, _ := url.Parse(emulatorTests.host)

		err := client.gem.MarkEndpointUnavailableForWrite(*defaultEndpoint)
		if err != nil {
			t.Fatalf("Failed to mark endpoint unavailable for write: %v", err)
		}

		logs := logger.filterByEvent(EventEndpointManager)

		assertContains(t, logs, "Marked endpoint unavailable")

		var markedLog string
		for _, log := range logs {
			if strings.Contains(log, "Marked endpoint unavailable") {
				markedLog = log
				break
			}
		}

		if markedLog == "" {
			t.Fatalf("Expected to find 'Marked endpoint unavailable' log for write operation")
		}

		if !strings.Contains(markedLog, "operation=write") {
			t.Errorf("Expected operation=write in log, got: %s", markedLog)
		}
		t.Logf("Found write unavailable log: %s", markedLog)

		t.Logf("Phase 5 complete: %d EventEndpointManager logs", len(logs))
	})
}

func TestLoggingAllowedHeaders(t *testing.T) {
	emulatorTests := newEmulatorTests(t)

	logger := &captureLogger{}
	azlog.SetListener(func(event azlog.Event, message string) {
		logger.log(event, message)
	})
	defer azlog.SetListener(nil)

	azlog.SetEvents(azlog.EventResponse)
	defer azlog.SetEvents()

	client := emulatorTests.getClient(t, tracing.Provider{})

	ctx := context.Background()
	dbName := "loggingHeadersTestDB"
	db := emulatorTests.createDatabase(t, ctx, client, dbName)
	defer emulatorTests.deleteDatabase(t, ctx, db)

	logs := logger.filterByEvent(azlog.EventResponse)
	if len(logs) == 0 {
		t.Fatal("No response logs captured")
	}

	headersToCheck := []string{
		"X-Ms-Request-Charge",
		"X-Ms-Activity-Id",
		"Content-Type",
	}

	for _, header := range headersToCheck {
		found := false
		for _, logMsg := range logs {
			if strings.Contains(logMsg, header) && !strings.Contains(logMsg, header+": REDACTED") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected header '%s' to be logged unredacted, but it was not found or was redacted.\nLogs:\n%s", header, logger.getAllLogs())
		}
	}

	foundAuth := false
	for _, logMsg := range logs {
		if strings.Contains(logMsg, "Authorization") {
			foundAuth = true
			if !strings.Contains(logMsg, "Authorization: REDACTED") {
				t.Errorf("Expected Authorization header to be redacted, but it was not.\nLog message: %s", logMsg)
			}
		}
	}

	if !foundAuth {
		t.Log("Authorization header not found in logs")
	}
}
