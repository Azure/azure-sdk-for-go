// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEmulatorQueryPaginationAfterContainerRecreation(t *testing.T) {
	endpoint, databaseID, _ := emulatorConfiguration(t)
	account, err := url.Parse(endpoint)
	require.NoError(t, err)
	if account.Scheme != "http" || (account.Hostname() != "127.0.0.1" && account.Hostname() != "localhost") {
		t.Skip("requires the local Rust emulator's unauthenticated resource-management API")
	}
	containerID := uniqueItemID(t)
	collectionPath := fmt.Sprintf("dbs/%s/colls", url.PathEscape(databaseID))
	resourcePath := collectionPath + "/" + url.PathEscape(containerID)
	metadata := []byte(fmt.Sprintf(`{"id":%q,"partitionKey":{"paths":["/pk"],"kind":"Hash","version":2}}`, containerID))
	manage := func(method, path string, body []byte, want int) []byte {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		request, err := http.NewRequestWithContext(ctx, method, endpoint+path, bytes.NewReader(body))
		require.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("x-ms-version", "2018-12-31")
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer response.Body.Close()
		result, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		require.Equal(t, want, response.StatusCode, string(result))
		return result
	}
	created := manage(http.MethodPost, collectionPath, metadata, http.StatusCreated)
	t.Cleanup(func() { manage(http.MethodDelete, resourcePath, nil, http.StatusNoContent) })
	var original struct {
		RID string `json:"_rid"`
	}
	require.NoError(t, json.Unmarshal(created, &original))
	require.NotEmpty(t, original.RID)

	// The pinned emulator doesn't emit RID-mismatch errors, so inject the service's recovery
	// signal once, while keeping resource recreation and query execution on the real emulator.
	var invalidate atomic.Bool
	var proxyURL string
	proxy := httputil.NewSingleHostReverseProxy(account)
	proxy.ModifyResponse = func(response *http.Response) error {
		if response.Request.Method != http.MethodGet || response.Request.URL.Path != "/" || response.StatusCode != http.StatusOK {
			return nil
		}
		defer response.Body.Close()
		var properties map[string]json.RawMessage
		if err := json.NewDecoder(response.Body).Decode(&properties); err != nil {
			return err
		}
		for _, field := range []string{"readableLocations", "writableLocations"} {
			var locations []map[string]any
			if err := json.Unmarshal(properties[field], &locations); err != nil {
				return err
			}
			for _, location := range locations {
				location["databaseAccountEndpoint"] = proxyURL + "/"
			}
			encoded, err := json.Marshal(locations)
			if err != nil {
				return err
			}
			properties[field] = encoded
		}
		body, err := json.Marshal(properties)
		if err != nil {
			return err
		}
		response.Body = io.NopCloser(bytes.NewReader(body))
		response.ContentLength = int64(len(body))
		response.Header.Set("Content-Length", strconv.Itoa(len(body)))
		return nil
	}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/"+resourcePath+"/docs/item-0" && invalidate.CompareAndSwap(true, false) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-ms-substatus", "1024")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"code":"BadRequest","message":"container was recreated"}`)
			return
		}
		proxy.ServeHTTP(w, r)
	}))
	proxyURL = "http://" + server.Listener.Addr().String()
	server.Start()
	t.Cleanup(server.Close)
	credential, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)
	client, err := NewClientWithKey(server.URL, credential, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	key := NewPartitionKeyString("partition")
	seed := func(container *ContainerClient) {
		t.Helper()
		for i := 0; i < 5; i++ {
			id := fmt.Sprintf("item-%d", i)
			_, err := container.CreateItem(t.Context(), key, id,
				[]byte(fmt.Sprintf(`{"id":%q,"pk":"partition","value":%d}`, id, i)), nil)
			require.NoError(t, err)
		}
	}
	seed(container) // Populate the original client's Go handle cache.
	query := NewQuery("SELECT VALUE c.value FROM c ORDER BY c.value")
	scope := NewFeedScopeForPartitionKey(key)
	options := &QueryOptions{Feed: FeedOptions{PageSizeHint: 2}}
	oldPage, err := container.NewQueryItemsPager(query, scope, options).NextPage(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, oldPage.ContinuationToken)

	manage(http.MethodDelete, resourcePath, nil, http.StatusNoContent)
	recreated := manage(http.MethodPost, collectionPath, metadata, http.StatusCreated)
	var replacement struct {
		RID string `json:"_rid"`
	}
	require.NoError(t, json.Unmarshal(recreated, &replacement))
	require.NotEmpty(t, replacement.RID)
	require.NotEqual(t, original.RID, replacement.RID)
	freshClient, _, _ := emulatorClient(t)
	freshContainer, err := freshClient.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	seed(freshContainer)

	// A point read lets native recovery discover the replacement while Go retains the old handle.
	invalidate.Store(true)
	_, err = container.ReadItem(t.Context(), key, "item-0", nil)
	require.NoError(t, err)
	require.False(t, invalidate.Load(), "the stale-container signal must reach the native driver")
	pager := container.NewQueryItemsPager(query, scope, options)
	first, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("0"), []byte("1")}, first.Items)
	require.NotEmpty(t, first.ContinuationToken)
	var remaining [][]byte
	for count := 0; pager.More(); count++ {
		require.Less(t, count, 5)
		page, err := pager.NextPage(t.Context())
		require.NoError(t, err, "tokens issued for the replacement must work on the same client")
		remaining = append(remaining, page.Items...)
	}
	require.Equal(t, [][]byte{[]byte("2"), []byte("3"), []byte("4")}, remaining)

	resumed := container.NewQueryItemsPager(query, scope, &QueryOptions{
		Feed: FeedOptions{PageSizeHint: 2, ContinuationToken: first.ContinuationToken},
	})
	var resumedItems [][]byte
	for count := 0; resumed.More(); count++ {
		require.Less(t, count, 5)
		page, err := resumed.NextPage(t.Context())
		require.NoError(t, err)
		resumedItems = append(resumedItems, page.Items...)
	}
	require.Equal(t, remaining, resumedItems)

	stale := container.NewQueryItemsPager(query, scope, &QueryOptions{
		Feed: FeedOptions{ContinuationToken: oldPage.ContinuationToken},
	})
	page, err := stale.NextPage(t.Context())
	require.Zero(t, page)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr, "tokens from the deleted container must still be rejected")
	require.Equal(t, http.StatusBadRequest, cosmosErr.StatusCode)
}

func TestEmulatorQueryPagination(t *testing.T) {
	container := emulatorContainer(t)
	partition := uniqueItemID(t)
	key := NewPartitionKeyString(partition)
	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("%s-%d", partition, i)
		body, err := json.Marshal(map[string]any{"id": id, "pk": partition, "value": i})
		require.NoError(t, err)
		_, err = container.CreateItem(t.Context(), key, id, body, nil)
		require.NoError(t, err)
		trackEmulatorItem(t, container, key, id)
	}
	otherKey := NewPartitionKeyString(partition + "-other")
	otherBody, err := json.Marshal(map[string]any{"id": partition, "pk": partition + "-other", "value": 0})
	require.NoError(t, err)
	_, err = container.CreateItem(t.Context(), otherKey, partition, otherBody, nil)
	require.NoError(t, err)
	trackEmulatorItem(t, container, otherKey, partition)

	query, err := NewQuery("SELECT VALUE c.value FROM c WHERE c.value >= @min ORDER BY c.value").WithParameter("@min", 0)
	require.NoError(t, err)
	scope := NewFeedScopeForPartitionKey(key)
	pager := container.NewQueryItemsPager(query, scope, &QueryOptions{Feed: FeedOptions{PageSizeHint: 2}})
	first, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("0"), []byte("1")}, first.Items, "the ABI must return every item, not just the first")
	require.Positive(t, first.RequestCharge)
	require.NotEmpty(t, first.ActivityID)
	require.NotEmpty(t, first.SessionToken)
	require.NotEmpty(t, first.ContinuationToken)

	remaining := make([][]byte, 0)
	for pages := 0; pager.More(); pages++ {
		require.Less(t, pages, 10, "query failed to terminate")
		page, err := pager.NextPage(t.Context())
		require.NoError(t, err)
		remaining = append(remaining, page.Items...)
	}
	require.Equal(t, [][]byte{[]byte("2"), []byte("3"), []byte("4")}, remaining)

	resumed := container.NewQueryItemsPager(query, scope, &QueryOptions{
		Feed: FeedOptions{PageSizeHint: 2, ContinuationToken: first.ContinuationToken},
	})
	var resumedItems [][]byte
	for pages := 0; resumed.More(); pages++ {
		require.Less(t, pages, 10, "resumed query failed to terminate")
		page, err := resumed.NextPage(t.Context())
		require.NoError(t, err)
		resumedItems = append(resumedItems, page.Items...)
	}
	require.Equal(t, remaining, resumedItems)

	empty := container.NewQueryItemsPager(NewQuery("SELECT * FROM c WHERE c.value = -1"), scope, nil)
	for pages := 0; empty.More(); pages++ {
		require.Less(t, pages, 3, "empty query failed to terminate")
		page, err := empty.NextPage(t.Context())
		require.NoError(t, err)
		require.Empty(t, page.Items)
	}

	invalid := container.NewQueryItemsPager(query, scope, &QueryOptions{Feed: FeedOptions{ContinuationToken: "not-a-continuation"}})
	page, err := invalid.NextPage(t.Context())
	require.Error(t, err)
	require.Zero(t, page)
}

func TestEmulatorQueryHierarchicalScope(t *testing.T) {
	client, databaseID, _ := emulatorClient(t)
	container, err := client.NewContainer(databaseID, "query-hierarchical")
	require.NoError(t, err)
	partition := uniqueItemID(t)
	key := NewPartitionKeyString(partition).AppendString("child")
	body, err := json.Marshal(map[string]any{"id": partition, "pk": partition, "child": "child"})
	require.NoError(t, err)
	_, err = container.CreateItem(t.Context(), key, partition, body, nil)
	require.NoError(t, err)
	trackEmulatorItem(t, container, key, partition)

	pager := container.NewQueryItemsPager(NewQuery("SELECT VALUE c.id FROM c"), NewFeedScopeForPartitionKey(key), nil)
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	var id string
	require.NoError(t, json.Unmarshal(page.Items[0], &id))
	require.Equal(t, partition, id)

	prefix := container.NewQueryItemsPager(NewQuery("SELECT * FROM c"),
		NewFeedScopeForPartitionKey(NewPartitionKeyString(partition)), nil)
	page, err = prefix.NextPage(t.Context())
	require.Zero(t, page)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeBadRequest, cosmosErr.Code)
	require.Contains(t, cosmosErr.Message, "prefix queries are not supported")
	require.Positive(t, cosmosErr.RequestCharge)
	require.NotEmpty(t, cosmosErr.ActivityID)
}

func TestEmulatorQuerySpecialPartitionKeys(t *testing.T) {
	container := emulatorContainer(t)
	for _, tt := range []struct {
		name  string
		key   PartitionKey
		value any
	}{
		{"null", NewPartitionKeyNull(), nil},
		{"empty string", NewPartitionKeyString(""), ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			id := uniqueItemID(t)
			item := map[string]any{"id": id, "pk": tt.value}
			body, err := json.Marshal(item)
			require.NoError(t, err)
			_, err = container.CreateItem(t.Context(), tt.key, id, body, nil)
			require.NoError(t, err)
			trackEmulatorItem(t, container, tt.key, id)
			query, err := NewQuery("SELECT * FROM c WHERE c.id = @id").WithParameter("@id", id)
			require.NoError(t, err)
			page, err := container.NewQueryItemsPager(query, NewFeedScopeForPartitionKey(tt.key), nil).NextPage(t.Context())
			require.NoError(t, err)
			require.Len(t, page.Items, 1)
			var result map[string]any
			require.NoError(t, json.Unmarshal(page.Items[0], &result))
			require.Equal(t, id, result["id"])
		})
	}
}

func TestEmulatorQueryCancellationAndClose(t *testing.T) {
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	for _, action := range []string{"cancel", "close"} {
		t.Run(action, func(t *testing.T) {
			credential := &delayedTokenCredential{started: make(chan struct{}), release: make(chan struct{})}
			client, err := NewClient(endpoint, credential, nil)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			container, err := client.NewContainer(databaseID, containerID)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			type result struct {
				page QueryItemsResponse
				err  error
			}
			results := make(chan result, 1)
			go func() {
				page, err := container.NewQueryItemsPager(NewQuery("SELECT * FROM c"),
					NewFeedScopeForPartitionKey(NewPartitionKeyString("cancel-query")), nil).NextPage(ctx)
				results <- result{page, err}
			}()
			select {
			case <-credential.started:
			case <-time.After(10 * time.Second):
				t.Fatal("query did not reach token acquisition")
			}
			if action == "cancel" {
				cancel()
			} else {
				closed := make(chan error, 1)
				go func() { closed <- client.Close() }()
				select {
				case err := <-closed:
					require.NoError(t, err)
				case <-time.After(10 * time.Second):
					t.Fatal("Close did not release the in-flight query")
				}
			}
			select {
			case got := <-results:
				require.Error(t, got.err)
				require.Zero(t, got.page)
				if action == "cancel" {
					require.ErrorIs(t, got.err, context.Canceled)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("query did not finish after cancellation")
			}
		})
	}
}
