// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore AVAD

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// These tests cover binding and scope isolation, not historical/AVAD emulator fidelity.
func TestEmulatorChangeFeedPartitionIsolation(t *testing.T) {
	for _, hierarchical := range []bool{false, true} {
		t.Run(fmt.Sprintf("hierarchical-%t", hierarchical), func(t *testing.T) {
			client, databaseID, containerID := emulatorClient(t)
			if hierarchical {
				containerID = "query-hierarchical"
			}
			container, err := client.NewContainer(t.Context(), databaseID, containerID)
			require.NoError(t, err)
			prefix := fmt.Sprintf("change-feed-%d", time.Now().UnixNano())
			key := NewPartitionKeyString(prefix)
			foreignKey := NewPartitionKeyString(prefix + "-other")
			if hierarchical {
				key = key.AppendString("child")
				foreignKey = NewPartitionKeyString(prefix).AppendString("other-child")
			}
			wanted := map[string]bool{prefix + "-first": false, prefix + "-second": false}
			foreignID := prefix + "-foreign"
			for id := range wanted {
				seedChangeFeedEmulatorItem(t, container, key, id, prefix, "child")
			}
			foreignPartition, foreignChild := prefix+"-other", "child"
			if hierarchical {
				foreignPartition, foreignChild = prefix, "other-child"
			}
			seedChangeFeedEmulatorItem(t, container, foreignKey, foreignID, foreignPartition, foreignChild)
			pager := container.NewChangeFeedPager(NewFeedScopeForPartitionKey(key), NewChangeFeedStartFromBeginning(),
				&ChangeFeedOptions{Feed: FeedOptions{PageSizeHint: 100}, Operation: OperationOptions{EndToEndTimeout: to(5 * time.Second)}})
			defer func() { require.NoError(t, pager.Close()) }()
			for range 8 {
				page, err := pager.NextPage(t.Context())
				require.NoError(t, err)
				for _, item := range page.Items {
					id := changeFeedEmulatorItemID(t, item)
					require.NotEqual(t, foreignID, id, "a different complete partition key must remain isolated")
					_, expected := wanted[id]
					require.True(t, expected, "unexpected item %q from another partition", id)
					wanted[id] = true
				}
				require.True(t, pager.More(), "scope polling must not permanently exhaust")
			}
			for id, seen := range wanted {
				require.True(t, seen, "the target partition must include %q", id)
			}
		})
	}
}

func TestEmulatorChangeFeedFullContainerScope(t *testing.T) {
	container := emulatorContainer(t)
	prefix := fmt.Sprintf("change-feed-container-%d", time.Now().UnixNano())
	wanted := map[string]bool{prefix + "-first": false, prefix + "-second": false}
	for id := range wanted {
		seedChangeFeedEmulatorItem(t, container, NewPartitionKeyString(id), id, id, "child")
	}
	pager := container.NewChangeFeedPager(NewFeedScopeForFullContainer(), NewChangeFeedStartFromBeginning(),
		&ChangeFeedOptions{Feed: FeedOptions{PageSizeHint: 1000}, Operation: OperationOptions{EndToEndTimeout: to(5 * time.Second)}})
	defer func() { require.NoError(t, pager.Close()) }()
	for range 32 {
		page, err := pager.NextPage(t.Context())
		require.NoError(t, err)
		for _, item := range page.Items {
			id := changeFeedEmulatorItemID(t, item)
			if _, expected := wanted[id]; expected {
				wanted[id] = true
			}
		}
		require.True(t, pager.More())
		if wanted[prefix+"-first"] && wanted[prefix+"-second"] {
			break
		}
	}
	for id, seen := range wanted {
		require.True(t, seen, "full-container polling must include %q", id)
	}
}

func TestEmulatorChangeFeedAllVersionsStartValidation(t *testing.T) {
	client, databaseID, _ := emulatorClient(t)
	containerID := uniqueItemID(t)
	// ChangeFeedPolicy is intentionally not exposed by the v2 public container model. Create this
	// test-only container through the normal native management path so the backend can validate
	// AVAD starts without changing the public API or the shared emulator seed configuration.
	body, err := json.Marshal(struct {
		ID                     string                 `json:"id"`
		PartitionKeyDefinition PartitionKeyDefinition `json:"partitionKey"`
		ChangeFeedPolicy       struct {
			RetentionDuration string `json:"retentionDuration"`
		} `json:"changeFeedPolicy"`
	}{
		ID:                     containerID,
		PartitionKeyDefinition: PartitionKeyDefinition{Kind: PartitionKeyKindHash, Paths: []string{"/pk"}, Version: 2},
		ChangeFeedPolicy: struct {
			RetentionDuration string `json:"retentionDuration"`
		}{RetentionDuration: "PT5M"},
	})
	require.NoError(t, err)

	release, err := client.acquire()
	require.NoError(t, err)
	_, _, err = client.executeManagement(t.Context(), operationKindCreateContainer, databaseID, containerID, body, OperationOptions{})
	release()
	require.NoError(t, err)

	container, err := client.NewContainer(t.Context(), databaseID, containerID)
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := container.Delete(ctx, nil)
		require.NoError(t, err)
	})

	scope := NewFeedScopeForPartitionKey(NewPartitionKeyString("partition"))
	supported := container.NewChangeFeedPager(scope, NewChangeFeedStartFromNow(),
		&ChangeFeedOptions{Mode: ChangeFeedModeAllVersionsAndDeletes})
	page, err := supported.NextPage(t.Context())
	require.NoError(t, err, "StartFromNow is the supported AVAD positive control")
	require.Contains(t, []int{http.StatusOK, http.StatusNotModified}, page.StatusCode)
	require.True(t, supported.More())
	require.NoError(t, supported.Close())

	for _, tt := range []struct {
		name  string
		start ChangeFeedStartFrom
	}{
		{name: "beginning", start: NewChangeFeedStartFromBeginning()},
		{name: "point-in-time", start: NewChangeFeedStartFromPointInTime(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pager := container.NewChangeFeedPager(scope, tt.start,
				&ChangeFeedOptions{Mode: ChangeFeedModeAllVersionsAndDeletes})
			page, err := pager.NextPage(t.Context())
			require.Zero(t, page)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Equal(t, http.StatusBadRequest, cosmosErr.StatusCode)
			require.Equal(t, CodeBadRequest, cosmosErr.Code)
			require.True(t, cosmosErr.FromWire, "the configured emulator must reject the unsupported start")
			require.False(t, pager.More())
			require.NoError(t, pager.Close())
		})
	}
}

func seedChangeFeedEmulatorItem(t *testing.T, container *ContainerClient, key PartitionKey, id, partition, child string) {
	t.Helper()
	body, err := json.Marshal(map[string]string{"id": id, "pk": partition, "child": child})
	require.NoError(t, err)
	_, err = container.CreateItem(t.Context(), key, id, json.RawMessage(body), nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := container.DeleteItem(ctx, key, id, nil)
		require.NoError(t, err)
	})
}

func changeFeedEmulatorItemID(t *testing.T, body []byte) string {
	t.Helper()
	var item struct {
		ID      string `json:"id"`
		Current *struct {
			ID string `json:"id"`
		} `json:"current"`
	}
	require.NoError(t, json.Unmarshal(body, &item))
	if item.Current != nil {
		return item.Current.ID
	}
	// The pinned snapshot-only emulator can return plain documents instead of v2 envelopes.
	return item.ID
}
