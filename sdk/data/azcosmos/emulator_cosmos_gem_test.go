// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azcosmos

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGlobalEndpointManager(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	_, err := url.Parse(client.Endpoint())
	if err != nil {
		fmt.Println("Unable to parse default endpoint URI")
		os.Exit(1)
	}

	preferredRegions := []string{"West US", "Central US"}
	gem, err := newGlobalEndpointManager(client, preferredRegions, 5)
	if err != nil {
		t.Fatalf("failed to create globalEndpointManager: %v", err)
	}
	fmt.Printf("gem: %v", gem)

	t.Run("TestGetAccountProperties", func(t *testing.T) {
		accountProperties, err := gem.GetAccountProperties()
		assert.NoError(t, err, "GetAccountProperties should not return an error")
		assert.NotNil(t, accountProperties, "accountProperties should not be nil")

		ReadRegions := accountProperties.ReadRegions
		assert.NotNil(t, ReadRegions, "ReadRegions should not be nil")
		assert.Len(t, ReadRegions, 2, "ReadRegions should contain 2 regions")

		WriteRegion := accountProperties.WriteRegions
		assert.NotNil(t, WriteRegion, "WriteRegion should not be nil")
		assert.Len(t, WriteRegion, 1, "WriteRegion should contain 1 region")
	})

	t.Run("TestGetWriteEndpoints", func(t *testing.T) {
		writeEndpoints, err := gem.GetWriteEndpoints()
		assert.NoError(t, err, "GetWriteEndpoints should not return an error")
		assert.NotNil(t, writeEndpoints, "writeEndpoints should not be nil")
		assert.Len(t, writeEndpoints, 2, "writeEndpoints should contain 2 endpoints")
		assert.Equal(t, loc1Endpoint, writeEndpoints[0], "writeEndpoints[0] should match loc1Endpoint")
		assert.Equal(t, loc2Endpoint, writeEndpoints[1], "writeEndpoints[1] should match loc2Endpoint")
	})

	t.Run("TestGetReadEndpoints", func(t *testing.T) {
		readEndpoints, err := gem.GetReadEndpoints()
		assert.NoError(t, err, "GetReadEndpoints should not return an error")
		assert.NotNil(t, readEndpoints, "readEndpoints should not be nil")
		assert.Len(t, readEndpoints, 2, "readEndpoints should contain 2 endpoints")
		assert.Equal(t, loc1Endpoint, readEndpoints[0], "readEndpoints[0] should match loc1Endpoint")
		assert.Equal(t, loc2Endpoint, readEndpoints[1], "readEndpoints[1] should match loc2Endpoint")
	})

	t.Run("TestGetLocation", func(t *testing.T) {
		location := gem.GetLocation(*loc1Endpoint)
		assert.Equal(t, "location1", location, "location should match expected location")
	})

	t.Run("TestMarkEndpointUnavailableForRead", func(t *testing.T) {
		err := gem.MarkEndpointUnavailableForRead(*loc1Endpoint)
		assert.NoError(t, err, "MarkEndpointUnavailableForRead should not return an error")
	})

	t.Run("TestMarkEndpointUnavailableForWrite", func(t *testing.T) {
		err := gem.MarkEndpointUnavailableForWrite(*loc2Endpoint)
		assert.NoError(t, err, "MarkEndpointUnavailableForWrite should not return an error")
	})

	t.Run("TestUpdate", func(t *testing.T) {
		err = gem.Update()
		assert.NoError(t, err, "Update should not return an error")
	})

	t.Run("TestRefreshStaleEndpoints", func(t *testing.T) {
		// Perform refresh of stale endpoints
		gem.RefreshStaleEndpoints()
	})

	t.Run("TestIsEndpointUnavailable", func(t *testing.T) {

		isUnavailable := gem.IsEndpointUnavailable(*loc1Endpoint, read)
		assert.True(t, isUnavailable, "IsEndpointUnavailable should return true for marked Read endpoint")

		isUnavailable = gem.IsEndpointUnavailable(*loc2Endpoint, write)
		assert.True(t, isUnavailable, "IsEndpointUnavailable should return true for marked Write endpoint")
	})

	t.Run("TestCanUseMultipleWriteLocations", func(t *testing.T) {
		// canUseMultipleWrite := gem.CanUseMultipleWriteLocations()
	})

	t.Run("TestBackgroundRefresh", func(t *testing.T) {

		// Start background refresh
		gem.startBackgroundRefresh(5 * time.Minute)

		// Wait for background refresh to occur (assuming it takes less than 1 second)
		time.Sleep(1 * time.Second)

		// Stop background refresh
		gem.stopBackgroundRefresh()
	})

	t.Run("TestConcurrentAccess", func(t *testing.T) {

		// Concurrently access GetWriteEndpoints from multiple goroutines
		numGoroutines := 10
		done := make(chan struct{})
		errors := make(chan error, numGoroutines)
		wg := sync.WaitGroup{}

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				endpoints, err := gem.GetWriteEndpoints()
				if err != nil {
					errors <- err
					return
				}
				assert.Equal(t, []*url.URL{loc1Endpoint, loc2Endpoint}, endpoints, "concurrent access should return the same write endpoints")
			}()
		}

		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// All goroutines completed successfully
		case err := <-errors:
			t.Errorf("concurrent access error: %v", err)
		case <-time.After(5 * time.Second):
			t.Errorf("concurrent access timed out")
		}
	})

	t.Run("TestErrorConditions", func(t *testing.T) {
		t.Run("TestGetWriteEndpointsError", func(t *testing.T) {

			_, err := gem.GetWriteEndpoints()
			assert.Error(t, err, "GetWriteEndpoints should return an error")
		})

		t.Run("TestGetReadEndpointsError", func(t *testing.T) {

			_, err := gem.GetReadEndpoints()
			assert.Error(t, err, "GetReadEndpoints should return an error")
		})
	})

	t.Run("TestIntegration", func(t *testing.T) {
		client := emulatorTests.getClient(t)

		preferredRegions := []string{"location1", "location2"}
		gem, err := newGlobalEndpointManager(client, preferredRegions, (5 * time.Minute))
		if err != nil {
			t.Fatalf("failed to create globalEndpointManager: %v", err)
		}

		t.Run("TestGetWriteEndpointsIntegration", func(t *testing.T) {
			writeEndpoints, err := gem.GetWriteEndpoints()
			if err != nil {
				t.Errorf("GetWriteEndpoints failed: %v", err)
			}
			assert.Equal(t, 2, len(writeEndpoints), "GetWriteEndpoints should return 2 endpoints")
		})

		t.Run("TestGetReadEndpointsIntegration", func(t *testing.T) {
			readEndpoints, err := gem.GetReadEndpoints()
			if err != nil {
				t.Errorf("GetReadEndpoints failed: %v", err)
			}
			assert.Equal(t, 2, len(readEndpoints), "GetReadEndpoints should return 2 endpoints")
		})
	})
}
