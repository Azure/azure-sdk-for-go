// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestEmulatorChangeFeed_F1_SubRangeOfPhysicalRange exercises the original
// F1 bug pattern: a customer-supplied FeedRange that is a STRICT SUB-RANGE
// of a physical PK range (i.e., does not exact-match any pkrange boundary).
// Pre-F1 this returned silent no-data via the toHeaders exact-match path;
// post-F1 the overlap-resolution upstream finds the containing physical
// range and routes the request to it.
func TestEmulatorChangeFeed_F1_SubRangeOfPhysicalRange(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "cfF1SubRange")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)

	properties := ContainerProperties{
		ID:                     "cf-f1-sub",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}},
	}
	throughput := NewManualThroughputProperties(400)
	_, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	require.NoError(t, err)

	container, _ := database.NewContainer("cf-f1-sub")

	// Insert a few items
	for i := 0; i < 5; i++ {
		pkv := fmt.Sprintf("pk%d", i)
		doc := map[string]interface{}{"id": fmt.Sprintf("item%d", i), "pk": pkv}
		b, _ := json.Marshal(doc)
		_, err := container.CreateItem(context.TODO(), NewPartitionKeyString(pkv), b, nil)
		require.NoError(t, err)
	}

	// Wait for change feed propagation.
	time.Sleep(2 * time.Second)

	// FeedRange that is a strict sub-range of the typical "" -> "FF"
	// container partition. Pre-F1 this did not exact-match any physical
	// range and the request was silently dropped; post-F1 it resolves via
	// overlap-match and the page is returned.
	subRange := &FeedRange{MinInclusive: "10", MaxExclusive: "80"}
	resp, err := container.GetChangeFeed(context.TODO(), &ChangeFeedOptions{
		FeedRange:    subRange,
		MaxItemCount: 10,
	})
	require.NoError(t, err, "post-F1 sub-range FeedRange must resolve via overlap-match")
	// The request must have been issued and a continuation token persisted,
	// even if the page had zero results in this physical range.
	require.NotEmpty(t, resp.ContinuationToken, "continuation token must be populated")
}

// TestEmulatorChangeFeed_F1_WideRangeMultiDrain inserts items across a
// container and reads the change feed using a wide FeedRange that overlaps
// every physical range. Verifies the composite continuation token is well-
// formed and that pagination across calls eventually drains every range
// (the multi-range queue path).
func TestEmulatorChangeFeed_F1_WideRangeMultiDrain(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "cfF1WideDrain")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)

	properties := ContainerProperties{
		ID:                     "cf-f1-wide",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}},
	}
	// Higher throughput → emulator may allocate multiple physical PK ranges.
	throughput := NewManualThroughputProperties(10000)
	_, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	require.NoError(t, err)

	container, _ := database.NewContainer("cf-f1-wide")

	const totalItems = 50
	expectedIDs := make(map[string]bool, totalItems)
	for i := 0; i < totalItems; i++ {
		pkv := fmt.Sprintf("pk%d", i)
		id := fmt.Sprintf("item%d", i)
		expectedIDs[id] = true
		doc := map[string]interface{}{"id": id, "pk": pkv, "i": i}
		b, _ := json.Marshal(doc)
		_, err := container.CreateItem(context.TODO(), NewPartitionKeyString(pkv), b, nil)
		require.NoError(t, err)
	}

	time.Sleep(3 * time.Second)

	// Wide FeedRange that overlaps every physical range.
	wide := &FeedRange{MinInclusive: "", MaxExclusive: "FF"}

	seen := make(map[string]bool, totalItems)
	var token *string
	for iter := 0; iter < 20; iter++ { // safety cap
		opts := &ChangeFeedOptions{FeedRange: wide, MaxItemCount: 10, Continuation: token}
		resp, err := container.GetChangeFeed(context.TODO(), opts)
		require.NoError(t, err)

		// Verify token shape.
		var ct compositeContinuationToken
		require.NoError(t, json.Unmarshal([]byte(resp.ContinuationToken), &ct))
		require.Equal(t, cosmosCompositeContinuationTokenVersion, ct.Version)
		require.GreaterOrEqual(t, len(ct.Continuation), 1, "composite token must carry at least one range")

		for _, raw := range resp.Documents {
			var d map[string]interface{}
			require.NoError(t, json.Unmarshal(raw, &d))
			if id, ok := d["id"].(string); ok {
				seen[id] = true
			}
		}

		// 304 with no docs and an unchanged token signals drain complete.
		if resp.Count == 0 && len(seen) >= totalItems {
			break
		}
		// Loop forward with the rotated/composite token.
		ctTok := resp.ContinuationToken
		token = &ctTok
	}

	for id := range expectedIDs {
		require.True(t, seen[id], "expected to drain item %s across all ranges", id)
	}
}

// TestEmulatorChangeFeed_F1_CrossContainerTokenRejected verifies the
// cross-container guard: a continuation token whose embedded ResourceID
// belongs to a different container must be rejected loudly rather than
// misrouted against the wrong routing map.
func TestEmulatorChangeFeed_F1_CrossContainerTokenRejected(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "cfF1XContainer")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)

	throughput := NewManualThroughputProperties(400)
	_, err := database.CreateContainer(context.TODO(), ContainerProperties{
		ID:                     "container-a",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}},
	}, &CreateContainerOptions{ThroughputProperties: &throughput})
	require.NoError(t, err)
	_, err = database.CreateContainer(context.TODO(), ContainerProperties{
		ID:                     "container-b",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}},
	}, &CreateContainerOptions{ThroughputProperties: &throughput})
	require.NoError(t, err)

	containerA, _ := database.NewContainer("container-a")
	containerB, _ := database.NewContainer("container-b")

	// Seed container-a so the change feed has a real token to issue.
	doc := map[string]interface{}{"id": "a-item", "pk": "pk-a"}
	b, _ := json.Marshal(doc)
	_, err = containerA.CreateItem(context.TODO(), NewPartitionKeyString("pk-a"), b, nil)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)

	respA, err := containerA.GetChangeFeed(context.TODO(), &ChangeFeedOptions{
		FeedRange: &FeedRange{MinInclusive: "", MaxExclusive: "FF"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, respA.ContinuationToken)

	// Hand container-a's token to container-b → must be rejected.
	tokenA := respA.ContinuationToken
	_, err = containerB.GetChangeFeed(context.TODO(), &ChangeFeedOptions{
		Continuation: &tokenA,
	})
	require.Error(t, err, "cross-container continuation token must be rejected")
}
