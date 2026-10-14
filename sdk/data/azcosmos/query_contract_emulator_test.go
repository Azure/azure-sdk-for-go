// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore colls documentdb

package azcosmos

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEmulatorCursorCheckpointCancellation(t *testing.T) {
	key, query, _ := seedQueryContractItems(t, emulatorContainer(t))
	container := queryContractContainer(t, nil, nil)
	client := container.database.client
	pager := container.NewQueryItemsPager(query, NewFeedScopeForPartitionKey(key),
		&QueryOptions{Feed: FeedOptions{PageSizeHint: 1}})
	defer func() { require.NoError(t, pager.Close()) }()
	_, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	cursor, ok := pager.cursor.(*nativeQueryCursor)
	require.True(t, ok)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	type pendingDelivery struct {
		waiter chan cursorDelivery
		result cursorDelivery
	}
	held := make(chan pendingDelivery, 1)
	restore := func() {
		select {
		case delivery := <-held:
			delivery.waiter <- delivery.result
		default:
		}
	}
	// Always unblock cleanup before the test client's cleanup calls Close.
	defer restore()
	var resultKind uint32
	cursor.beforeCompletionWait = func(waiter chan cursorDelivery) {
		select {
		case completion := <-waiter:
			held <- pendingDelivery{waiter, completion}
			if completion.completion != nil {
				resultKind = uint32(completion.completion.result_kind)
			}
			cancel()
		case <-ctx.Done():
		}
	}
	token, err := pager.ContinuationToken(ctx)
	require.Len(t, held, 1, "the test must withhold an admitted checkpoint completion")
	require.EqualValues(t, 3, resultKind, "the withheld completion must be a real checkpoint")
	require.Empty(t, token)
	require.ErrorIs(t, err, context.Canceled)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
	require.False(t, pager.More())

	cleanupStarted := make(chan struct{})
	client.driver.beforePendingWait = func() { close(cleanupStarted) }
	closed := make(chan error, 1)
	go func() { closed <- client.Close() }()
	select {
	case <-cleanupStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("client teardown did not reach pending checkpoint cleanup")
	}
	select {
	case err := <-closed:
		t.Fatalf("client closed before withheld checkpoint completion was released: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	restore()
	select {
	case err := <-closed:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("client did not finish after checkpoint delivery")
	}
	require.Nil(t, cursor.handle)
}

func TestEmulatorCursorCompletionCancellationRace(t *testing.T) {
	key, query, _ := seedQueryContractItems(t, emulatorContainer(t))
	const activity = "33333333-3333-4333-8333-333333333333"
	container := queryContractContainer(t, nil, func(r *http.Response) error {
		if isQueryDataRequest(r.Request) && r.StatusCode == http.StatusOK {
			r.Header.Set("x-ms-request-charge", "6.25")
			r.Header.Set("x-ms-activity-id", activity)
		}
		return nil
	})
	for _, operation := range []string{"page", "checkpoint"} {
		t.Run(operation, func(t *testing.T) {
			for range 20 {
				pager := container.NewQueryItemsPager(query, NewFeedScopeForPartitionKey(key),
					&QueryOptions{Feed: FeedOptions{PageSizeHint: 1}})
				_, err := pager.NextPage(t.Context())
				require.NoError(t, err)
				cursor, ok := pager.cursor.(*nativeQueryCursor)
				require.True(t, ok)
				ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
				delivered := false
				var resultKind uint32
				cursor.beforeCompletionWait = func(waiter chan cursorDelivery) {
					select {
					case completion := <-waiter:
						waiter <- completion
						delivered = true
						if completion.completion != nil {
							resultKind = uint32(completion.completion.result_kind)
						}
						cancel()
					case <-ctx.Done():
					}
				}
				if operation == "page" {
					page, err := pager.NextPage(ctx)
					require.Zero(t, page)
					require.ErrorIs(t, err, context.Canceled)
					var cosmosErr *Error
					require.ErrorAs(t, err, &cosmosErr)
					require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
					require.Equal(t, 6.25, cosmosErr.RequestCharge)
					require.Equal(t, activity, cosmosErr.ActivityID)
				} else {
					token, err := pager.ContinuationToken(ctx)
					require.Empty(t, token)
					require.ErrorIs(t, err, context.Canceled)
					var cosmosErr *Error
					require.ErrorAs(t, err, &cosmosErr)
					require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
				}
				cancel()
				require.True(t, delivered, "both completion and cancellation must be ready before waiting")
				if operation == "page" {
					require.EqualValues(t, 2, resultKind)
				} else {
					require.EqualValues(t, 3, resultKind)
				}
				require.False(t, pager.More())
				require.NoError(t, pager.Close())
			}
		})
	}
}

func queryContractContainer(t *testing.T, intercept func(http.ResponseWriter, *http.Request) bool, response func(*http.Response) error) *ContainerClient {
	t.Helper()
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	target, err := url.Parse(endpoint)
	require.NoError(t, err)
	proxy := httputil.NewSingleHostReverseProxy(target)
	var proxyURL string
	proxy.ModifyResponse = func(r *http.Response) error {
		if err := rewriteQueryAccountResponse(r, proxyURL); err != nil {
			return err
		}
		if response != nil {
			return response(r)
		}
		return nil
	}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if intercept == nil || !intercept(w, r) {
			proxy.ServeHTTP(w, r)
		}
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
	return container
}

func seedQueryContractItems(t *testing.T, container *ContainerClient) (PartitionKey, Query, SessionToken) {
	t.Helper()
	partition := uniqueItemID(t)
	key := NewPartitionKeyString(partition)
	var session SessionToken
	for value := range 6 {
		id := fmt.Sprintf("%s-%d", partition, value)
		body, err := json.Marshal(map[string]any{"id": id, "pk": partition, "value": value})
		require.NoError(t, err)
		item, err := container.CreateItem(t.Context(), key, id, body, nil)
		require.NoError(t, err)
		session = item.SessionToken
		trackEmulatorItem(t, container, key, id)
	}
	require.NotEmpty(t, session)
	return key, NewQuery("SELECT VALUE c.value FROM c ORDER BY c.value"), session
}

func isQueryDataRequest(r *http.Request) bool {
	return r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/docs") &&
		!strings.EqualFold(r.Header.Get("x-ms-cosmos-is-query-plan-request"), "true")
}

func TestEmulatorCursorTimeoutDoesNotExpireBetweenPages(t *testing.T) {
	endpoint, database, containerID := emulatorConfiguration(t)
	shared, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: to(time.Second)}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shared.Close()) })
	client, err := NewClientWithKey(endpoint, KeyCredential{accountKey: emulatorKey}, &ClientOptions{Runtime: shared})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer(database, containerID)
	require.NoError(t, err)
	key, query, _ := seedQueryContractItems(t, container)
	pager := container.NewQueryItemsPager(query, NewFeedScopeForPartitionKey(key), &QueryOptions{Feed: FeedOptions{PageSizeHint: 1}})
	defer func() { require.NoError(t, pager.Close()) }()
	first, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("0")}, first.Items)
	time.Sleep(1100 * time.Millisecond)
	second, err := pager.NextPage(t.Context())
	require.NoError(t, err, "cursor must not retain its admission snapshot's absolute deadline")
	require.Equal(t, [][]byte{[]byte("1")}, second.Items)
}

func TestEmulatorQueryMidPaginationFailureMetadata(t *testing.T) {
	key, query, _ := seedQueryContractItems(t, emulatorContainer(t))
	const failureBody = `{"code":"BadRequest","message":"synthetic query page failure"}`
	const failureActivity = "22222222-2222-4222-8222-222222222222"
	const successActivity = "11111111-1111-4111-8111-111111111111"
	var requestCount atomic.Int32
	var failureToken string
	var mu sync.Mutex
	container := queryContractContainer(t, func(w http.ResponseWriter, r *http.Request) bool {
		if !isQueryDataRequest(r) {
			return false
		}
		if requestCount.Add(1) != 3 {
			return false
		}
		mu.Lock()
		token := failureToken
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-ms-request-charge", "4.75")
		w.Header().Set("x-ms-activity-id", failureActivity)
		w.Header().Set("x-ms-session-token", token)
		w.Header().Set("x-ms-substatus", "1234")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, failureBody)
		return true
	}, func(r *http.Response) error {
		if r.Request.Method == http.MethodGet && strings.Contains(r.Request.URL.Path, "/colls/") &&
			!strings.Contains(strings.Split(r.Request.URL.Path, "/colls/")[1], "/") {
			r.Header.Set("x-ms-request-charge", "1.25")
		}
		if isQueryDataRequest(r.Request) && r.StatusCode == http.StatusOK {
			r.Header.Set("x-ms-request-charge", "2.5")
			r.Header.Set("x-ms-activity-id", successActivity)
			mu.Lock()
			failureToken = r.Header.Get("x-ms-session-token")
			mu.Unlock()
		}
		return nil
	})
	pager := container.NewQueryItemsPager(query, NewFeedScopeForPartitionKey(key), &QueryOptions{Feed: FeedOptions{PageSizeHint: 2}})
	defer func() { require.NoError(t, pager.Close()) }()
	first, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("0"), []byte("1")}, first.Items)
	require.Equal(t, 3.75, first.RequestCharge, "first fetch includes the extra metadata-validation charge")
	require.Equal(t, successActivity, first.ActivityID)
	require.NotEmpty(t, first.SessionToken)
	second, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("2"), []byte("3")}, second.Items)
	require.Equal(t, 2.5, second.RequestCharge, "subsequent pages must not repeat setup charges")
	require.Equal(t, successActivity, second.ActivityID)
	require.NotEmpty(t, second.SessionToken)
	failed, err := pager.NextPage(t.Context())
	require.Zero(t, failed)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeBadRequest, cosmosErr.Code)
	require.Equal(t, http.StatusBadRequest, cosmosErr.StatusCode)
	require.Equal(t, 1234, cosmosErr.SubStatus)
	require.True(t, cosmosErr.FromWire)
	require.Equal(t, 4.75, cosmosErr.RequestCharge)
	require.Equal(t, failureActivity, cosmosErr.ActivityID)
	require.Equal(t, second.SessionToken, cosmosErr.SessionToken)
	require.JSONEq(t, failureBody, string(cosmosErr.Body))
	require.Equal(t, 11.0, first.RequestCharge+second.RequestCharge+cosmosErr.RequestCharge)
	require.False(t, pager.More())
	_, err = pager.NextPage(t.Context())
	require.Error(t, err)
	require.EqualValues(t, 3, requestCount.Load(), "terminal failure must not issue another request")
}

func TestEmulatorQueryMetricsAndSessionOptions(t *testing.T) {
	key, query, session := seedQueryContractItems(t, emulatorContainer(t))
	const indexMetrics = `{"UtilizedIndexes":[]}`
	const queryMetrics = "retrievedDocumentCount=2;outputDocumentCount=2"
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("metrics=%t", enabled), func(t *testing.T) {
			var mu sync.Mutex
			var requests []http.Header
			var sessions []SessionToken
			container := queryContractContainer(t, nil, func(r *http.Response) error {
				if !isQueryDataRequest(r.Request) || r.StatusCode != http.StatusOK {
					return nil
				}
				mu.Lock()
				requests = append(requests, r.Request.Header.Clone())
				sessions = append(sessions, SessionToken(r.Header.Get("x-ms-session-token")))
				mu.Unlock()
				// The pinned emulator lacks metric generation; emulate service headers
				// only when requested, exercising native decoding and the Go response path.
				r.Header.Del("x-ms-cosmos-index-utilization")
				r.Header.Del("x-ms-documentdb-query-metrics")
				if strings.EqualFold(r.Request.Header.Get("x-ms-cosmos-populateindexmetrics"), "true") {
					r.Header.Set("x-ms-cosmos-index-utilization", base64.StdEncoding.EncodeToString([]byte(indexMetrics)))
				}
				if strings.EqualFold(r.Request.Header.Get("x-ms-documentdb-populatequerymetrics"), "true") {
					r.Header.Set("x-ms-documentdb-query-metrics", queryMetrics)
				}
				return nil
			})
			pager := container.NewQueryItemsPager(query, NewFeedScopeForPartitionKey(key), &QueryOptions{
				Feed: FeedOptions{PageSizeHint: 2}, SessionToken: session,
				PopulateIndexMetrics: &enabled, PopulateQueryMetrics: &enabled,
			})
			defer func() { require.NoError(t, pager.Close()) }()
			for pageNumber := range 2 {
				page, err := pager.NextPage(t.Context())
				require.NoError(t, err)
				require.Equal(t, [][]byte{[]byte(fmt.Sprint(pageNumber * 2)), []byte(fmt.Sprint(pageNumber*2 + 1))}, page.Items)
				mu.Lock()
				observedRequests := append([]http.Header(nil), requests...)
				observedSessions := append([]SessionToken(nil), sessions...)
				mu.Unlock()
				require.Len(t, observedRequests, pageNumber+1)
				header := observedRequests[pageNumber]
				require.Equal(t, string(session), header.Get("x-ms-session-token"), "explicit token must survive retained paging")
				require.Equal(t, enabled, strings.EqualFold(header.Get("x-ms-cosmos-populateindexmetrics"), "true"))
				require.Equal(t, enabled, strings.EqualFold(header.Get("x-ms-documentdb-populatequerymetrics"), "true"))
				require.NotEmpty(t, observedSessions[pageNumber])
				require.Equal(t, observedSessions[pageNumber], page.SessionToken)
				if enabled {
					require.JSONEq(t, indexMetrics, page.IndexMetrics)
					require.Equal(t, queryMetrics, page.QueryMetrics)
				} else {
					require.Empty(t, page.IndexMetrics)
					require.Empty(t, page.QueryMetrics)
				}
			}
		})
	}
}
