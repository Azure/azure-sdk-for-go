// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore colls

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// This scripted service exercises the shipped driver's ABI, not service change-log fidelity.
type changeFeedContractService struct {
	mu         sync.Mutex
	url        string
	ranges     int
	idleFirst  bool
	headers    []http.Header
	unexpected []string
	fail       bool
}

func newChangeFeedContractService(t *testing.T, ranges int, idleFirst bool) *changeFeedContractService {
	t.Helper()
	service := &changeFeedContractService{ranges: ranges, idleFirst: idleFirst}
	server := httptest.NewUnstartedServer(http.HandlerFunc(service.serveHTTP))
	server.Config.Protocols = new(http.Protocols)
	server.Config.Protocols.SetHTTP1(true)
	server.Config.Protocols.SetUnencryptedHTTP2(true)
	server.Start()
	service.mu.Lock()
	service.url = server.URL
	service.mu.Unlock()
	t.Cleanup(func() {
		server.Close()
		service.mu.Lock()
		defer service.mu.Unlock()
		require.Empty(t, service.unexpected, "the fixture must handle every observed native request")
	})
	return service
}

func (s *changeFeedContractService) container(t *testing.T, id string) *ContainerClient {
	t.Helper()
	credential, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)
	s.mu.Lock()
	endpoint := s.url
	s.mu.Unlock()
	client, err := NewClientWithKey(endpoint, credential, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer("db", id)
	require.NoError(t, err)
	return container
}

func (s *changeFeedContractService) observed() []http.Header {
	s.mu.Lock()
	defer s.mu.Unlock()
	headers := make([]http.Header, len(s.headers))
	for i, header := range s.headers {
		headers[i] = header.Clone()
	}
	return headers
}

func (s *changeFeedContractService) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.UserAgent() == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-ms-request-charge", "3.5")
	w.Header().Set("x-ms-activity-id", "44444444-4444-4444-8444-444444444444")
	w.Header().Set("x-ms-session-token", "0:2")
	switch {
	case r.URL.Path == "/" || r.URL.Path == "":
		s.mu.Lock()
		endpoint := s.url
		s.mu.Unlock()
		_, _ = fmt.Fprintf(w, `{"id":"scripted-account","_rid":"scripted-account","_self":"",
			"continuousBackupEnabled":true,"writableLocations":[{"name":"East US","databaseAccountEndpoint":%q}],
			"readableLocations":[{"name":"East US","databaseAccountEndpoint":%q}]}`, endpoint+"/", endpoint+"/")
	case strings.HasSuffix(r.URL.Path, "/pkranges"):
		w.Header().Set("ETag", "ranges-position")
		if r.Header.Get("If-None-Match") == "ranges-position" {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		ranges := []map[string]string{{"id": "0", "minInclusive": "", "maxExclusive": "FF"}}
		if s.ranges == 2 {
			ranges[0]["maxExclusive"] = "80"
			ranges = append(ranges, map[string]string{"id": "1", "minInclusive": "80", "maxExclusive": "FF"})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"PartitionKeyRanges": ranges})
	case strings.HasSuffix(r.URL.Path, "/docs"):
		s.mu.Lock()
		s.headers = append(s.headers, r.Header.Clone())
		fail := s.fail
		s.mu.Unlock()
		if fail {
			w.Header().Set("ETag", "failed-position")
			w.Header().Set("x-ms-retry-after-ms", "25")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, `{"code":"BadRequest","message":"scripted change-feed failure"}`)
			return
		}
		rangeID := r.Header.Get("x-ms-documentdb-partitionkeyrangeid")
		idle := s.idleFirst && (rangeID == "0" || strings.HasSuffix(rangeID, ",0"))
		position := r.Header.Get("If-None-Match")
		w.Header().Set("ETag", "position-2")
		if idle || position == "position-2" {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		if position == "*" {
			w.Header().Set("ETag", "position-0")
			w.WriteHeader(http.StatusNotModified)
			return
		}
		if position == "position-1" {
			_, _ = fmt.Fprint(w, `{"Documents":[{"current":{"id":"second","n":9007199254740993},"previous":{"id":"second","n":1},"metadata":{"operationType":"replace"},"future":true}],"_count":1}`)
			return
		}
		w.Header().Set("ETag", "position-1")
		_, _ = fmt.Fprint(w, `{"Documents":[{"current":{"id":"first"},"metadata":{"operationType":"create"}}],"_count":1}`)
	case strings.Contains(r.URL.Path, "/colls/"):
		rid := "AQIDBAUGBwg="
		id := "items"
		if strings.HasSuffix(r.URL.Path, "/other") {
			rid, id = "AQIDBAgHBgU=", "other"
		}
		definition := `{"paths":["/key"],"kind":"Hash","version":2}`
		if strings.HasSuffix(r.URL.Path, "/hierarchical") {
			id, rid = "hierarchical", "AQIDBAkKCww="
			definition = `{"paths":["/tenant","/category","/key"],"kind":"MultiHash","version":2}`
		}
		_, _ = fmt.Fprintf(w, `{"id":%q,"_rid":%q,"_self":"","_etag":"container-position",
			"partitionKey":%s}`, id, rid, definition)
	default:
		s.mu.Lock()
		s.unexpected = append(s.unexpected, r.Method+" "+r.URL.String())
		s.mu.Unlock()
		http.Error(w, "unexpected scripted request", http.StatusInternalServerError)
	}
}

func TestNativeChangeFeedCheckpointResumeOnFreshClient(t *testing.T) {
	for _, full := range []bool{false, true} {
		t.Run(fmt.Sprintf("full-%t", full), func(t *testing.T) {
			service := newChangeFeedContractService(t, 1, false)
			container := service.container(t, "items")
			scope := NewFeedScopeForPartitionKey(NewPartitionKeyString("partition"))
			if full {
				scope = NewFeedScopeForFullContainer()
			}
			pager := container.NewChangeFeedPager(scope, NewChangeFeedStartFromBeginning(), &ChangeFeedOptions{
				Feed: FeedOptions{PageSizeHint: 7}, SessionToken: "0:1",
			})
			defer func() { require.NoError(t, pager.Close()) }()
			require.Empty(t, service.observed(), "construction must not fetch a page")
			page, err := pager.NextPage(t.Context())
			require.NoError(t, err)
			require.Equal(t, [][]byte{[]byte(`{"current":{"id":"first"},"metadata":{"operationType":"create"}}`)}, page.Items)
			charge := 3.5
			if !full {
				charge += 3.5
			}
			require.Equal(t, charge, page.RequestCharge, "include the additional scope-metadata read")
			token, err := pager.ContinuationToken(t.Context())
			require.NoError(t, err)
			require.NotEmpty(t, token)
			observed := len(service.observed())
			repeated, err := pager.ContinuationToken(t.Context())
			require.NoError(t, err)
			require.Equal(t, token, repeated)
			require.Len(t, service.observed(), observed, "checkpointing must not read changes")
			require.NoError(t, pager.Close())
			require.NoError(t, container.database.client.Close())

			fresh := service.container(t, "items")
			resumed := fresh.NewChangeFeedPager(scope, NewChangeFeedStartFromNow(), &ChangeFeedOptions{
				Feed: FeedOptions{PageSizeHint: 7, ContinuationToken: token}, SessionToken: "0:1",
			})
			defer func() { require.NoError(t, resumed.Close()) }()
			page, err = resumed.NextPage(t.Context())
			require.NoError(t, err)
			require.Equal(t, [][]byte{[]byte(`{"current":{"id":"second","n":9007199254740993},"previous":{"id":"second","n":1},"metadata":{"operationType":"replace"},"future":true}`)}, page.Items)
			require.EqualValues(t, "position-2", page.ETag)
			require.Equal(t, SessionToken("0:2"), page.SessionToken)
			page, err = resumed.NextPage(t.Context())
			require.NoError(t, err)
			require.Equal(t, http.StatusNotModified, page.StatusCode)
			require.Empty(t, page.Items)
			require.True(t, resumed.More())
			idleToken, err := resumed.ContinuationToken(t.Context())
			require.NoError(t, err)
			require.NotEmpty(t, idleToken)
			for _, header := range service.observed() {
				require.Equal(t, "7", header.Get("x-ms-max-item-count"))
				require.NotEmpty(t, header.Get("A-IM"))
			}
		})
	}
}

func TestNativeChangeFeedIdleRangeDoesNotEndContainer(t *testing.T) {
	service := newChangeFeedContractService(t, 2, true)
	pager := service.container(t, "items").NewChangeFeedPager(NewFeedScopeForFullContainer(),
		NewChangeFeedStartFromBeginning(), nil)
	defer func() { require.NoError(t, pager.Close()) }()
	idle, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Equal(t, http.StatusNotModified, idle.StatusCode)
	require.True(t, pager.More())
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Len(t, page.Items, 1, "an idle range must not discard another range's pending changes")
}

func TestNativeChangeFeedCheckpointCompatibility(t *testing.T) {
	service := newChangeFeedContractService(t, 1, false)
	pager := service.container(t, "items").NewChangeFeedPager(NewFeedScopeForFullContainer(),
		NewChangeFeedStartFromBeginning(), nil)
	_, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	token, err := pager.ContinuationToken(t.Context())
	require.NoError(t, err)
	require.NoError(t, pager.Close())
	for _, tt := range []struct {
		name, id string
		mode     ChangeFeedMode
		token    string
	}{
		{"mode", "items", ChangeFeedModeAllVersionsAndDeletes, token},
		{"container", "other", ChangeFeedModeLatestVersion, token},
		{"malformed", "items", ChangeFeedModeLatestVersion, "not a native token"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			incompatible := service.container(t, tt.id).NewChangeFeedPager(NewFeedScopeForFullContainer(),
				NewChangeFeedStartFromNow(), &ChangeFeedOptions{Mode: tt.mode, Feed: FeedOptions{ContinuationToken: tt.token}})
			defer func() { require.NoError(t, incompatible.Close()) }()
			before := len(service.observed())
			page, err := incompatible.NextPage(t.Context())
			require.Zero(t, page)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.False(t, cosmosErr.FromWire, "compatibility must be decided by the native planner")
			require.False(t, incompatible.More())
			require.Len(t, service.observed(), before, "incompatible checkpoints must not fetch a page")
		})
	}
}

func TestNativeChangeFeedModeAndStartOnWire(t *testing.T) {
	for _, tt := range []struct {
		name, aim, position, timestamp string
		mode                           ChangeFeedMode
		start                          ChangeFeedStartFrom
	}{
		{"latest beginning", "Incremental Feed", "", "", ChangeFeedModeLatestVersion, NewChangeFeedStartFromBeginning()},
		{"latest now", "Incremental Feed", "*", "", ChangeFeedModeLatestVersion, NewChangeFeedStartFromNow()},
		{"latest time", "Incremental Feed", "", "Fri, 02 Jan 2026 03:04:05 GMT", ChangeFeedModeLatestVersion,
			NewChangeFeedStartFromPointInTime(time.Date(2026, 1, 2, 3, 4, 5, 123, time.UTC))},
		{"all versions now", "Full-Fidelity Feed", "*", "", ChangeFeedModeAllVersionsAndDeletes, NewChangeFeedStartFromNow()},
	} {
		t.Run(tt.name, func(t *testing.T) {
			service := newChangeFeedContractService(t, 1, false)
			scope := NewFeedScopeForFullContainer()
			pager := service.container(t, "items").NewChangeFeedPager(scope, tt.start, &ChangeFeedOptions{Mode: tt.mode})
			defer func() { require.NoError(t, pager.Close()) }()
			_, err := pager.NextPage(t.Context())
			require.NoError(t, err)
			headers := service.observed()
			require.NotEmpty(t, headers)
			require.Equal(t, tt.aim, headers[len(headers)-1].Get("A-IM"))
			require.Equal(t, tt.position, headers[0].Get("If-None-Match"))
			require.Equal(t, tt.timestamp, headers[0].Get("If-Modified-Since"))
			token, err := pager.ContinuationToken(t.Context())
			require.NoError(t, err)
			require.NotEmpty(t, token)
			require.NoError(t, pager.Close())
			resumed := service.container(t, "items").NewChangeFeedPager(scope, NewChangeFeedStartFromBeginning(),
				&ChangeFeedOptions{Mode: tt.mode, Feed: FeedOptions{ContinuationToken: token}})
			defer func() { require.NoError(t, resumed.Close()) }()
			_, err = resumed.NextPage(t.Context())
			require.NoError(t, err, "saved positions must override a structurally valid supplied start")
		})
	}
}

func TestNativeChangeFeedAllVersionsPrimesEveryRange(t *testing.T) {
	service := newChangeFeedContractService(t, 2, false)
	pager := service.container(t, "items").NewChangeFeedPager(NewFeedScopeForFullContainer(),
		NewChangeFeedStartFromNow(), &ChangeFeedOptions{Mode: ChangeFeedModeAllVersionsAndDeletes})
	defer func() { require.NoError(t, pager.Close()) }()
	_, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	primed := make(map[string]bool)
	for _, header := range service.observed() {
		if header.Get("If-None-Match") == "*" {
			primed[header.Get("x-ms-documentdb-partitionkeyrangeid")] = true
		}
	}
	require.Len(t, primed, 2, "all ranges must have concrete positions before the first delivered page")
	token, err := pager.ContinuationToken(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestNativeChangeFeedAllVersionsStartRejectionComesFromService(t *testing.T) {
	for _, start := range []ChangeFeedStartFrom{
		NewChangeFeedStartFromBeginning(),
		NewChangeFeedStartFromPointInTime(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)),
	} {
		service := newChangeFeedContractService(t, 1, false)
		service.mu.Lock()
		service.fail = true
		service.mu.Unlock()
		pager := service.container(t, "items").NewChangeFeedPager(NewFeedScopeForFullContainer(), start,
			&ChangeFeedOptions{Mode: ChangeFeedModeAllVersionsAndDeletes})
		page, err := pager.NextPage(t.Context())
		require.Zero(t, page)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.True(t, cosmosErr.FromWire, "Go must not impose its own mode/start policy")
		require.Equal(t, CodeBadRequest, cosmosErr.Code)
		require.NotEmpty(t, service.observed())
		require.False(t, pager.More())
		require.NoError(t, pager.Close())
	}
}

func TestNativeChangeFeedForwardsHierarchicalKeys(t *testing.T) {
	for _, tt := range []struct {
		name, wire string
		key        PartitionKey
	}{
		{"prefix", `["tenant"]`, NewPartitionKeyString("tenant")},
		{"complete", `["tenant","category","item"]`, NewPartitionKeyString("tenant").AppendString("category").AppendString("item")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			service := newChangeFeedContractService(t, 2, false)
			pager := service.container(t, "hierarchical").NewChangeFeedPager(NewFeedScopeForPartitionKey(tt.key),
				NewChangeFeedStartFromBeginning(), nil)
			defer func() { require.NoError(t, pager.Close()) }()
			_, err := pager.NextPage(t.Context())
			require.NoError(t, err)
			headers := service.observed()
			require.NotEmpty(t, headers)
			require.JSONEq(t, tt.wire, headers[0].Get("x-ms-documentdb-partitionkey"),
				"forward the logical key; this fixture does not certify prefix fan-out")
		})
	}
}

func TestNativeChangeFeedInitialFanOutLimit(t *testing.T) {
	service := newChangeFeedContractService(t, 2, false)
	pager := service.container(t, "items").NewChangeFeedPager(NewFeedScopeForFullContainer(),
		NewChangeFeedStartFromBeginning(), &ChangeFeedOptions{Feed: FeedOptions{MaxFanOut: 1}})
	defer func() { require.NoError(t, pager.Close()) }()
	page, err := pager.NextPage(t.Context())
	require.Zero(t, page)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.False(t, cosmosErr.FromWire)
	require.False(t, pager.More())
	require.Empty(t, service.observed(), "the native planner must apply its fan-out limit")
}

func TestNativeChangeFeedWireErrorAndSetupCharge(t *testing.T) {
	for _, full := range []bool{false, true} {
		service := newChangeFeedContractService(t, 1, false)
		service.mu.Lock()
		service.fail = true
		service.mu.Unlock()
		scope := NewFeedScopeForFullContainer()
		if !full {
			scope = NewFeedScopeForPartitionKey(NewPartitionKeyString("partition"))
		}
		pager := service.container(t, "items").NewChangeFeedPager(scope, NewChangeFeedStartFromBeginning(), nil)
		page, err := pager.NextPage(t.Context())
		require.Zero(t, page)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeBadRequest, cosmosErr.Code)
		require.Equal(t, http.StatusBadRequest, cosmosErr.StatusCode)
		require.True(t, cosmosErr.FromWire)
		charge := 3.5
		if !full {
			charge += 3.5
		}
		require.Equal(t, charge, cosmosErr.RequestCharge)
		require.Equal(t, "44444444-4444-4444-8444-444444444444", cosmosErr.ActivityID)
		require.Equal(t, SessionToken("0:2"), cosmosErr.SessionToken)
		require.EqualValues(t, "failed-position", cosmosErr.ETag)
		require.Equal(t, 25*time.Millisecond, cosmosErr.RetryAfter)
		require.Equal(t, `{"code":"BadRequest","message":"scripted change-feed failure"}`, string(cosmosErr.Body))
		require.NotNil(t, cosmosErr.Diagnostics)
		require.GreaterOrEqual(t, cosmosErr.AttemptCount, uint32(1))
		require.False(t, pager.More())
		require.NoError(t, pager.Close())
	}
}

func TestNativeChangeFeedCompletionCancellationRace(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNotModified} {
		t.Run(fmt.Sprintf("status-%d", status), func(t *testing.T) {
			service := newChangeFeedContractService(t, 1, false)
			pager := service.container(t, "items").NewChangeFeedPager(NewFeedScopeForFullContainer(),
				NewChangeFeedStartFromBeginning(), nil)
			defer func() { require.NoError(t, pager.Close()) }()
			_, err := pager.NextPage(t.Context())
			require.NoError(t, err)
			if status == http.StatusNotModified {
				_, err = pager.NextPage(t.Context())
				require.NoError(t, err)
			}
			cursor, ok := pager.cursor.(*nativeChangeFeedCursor)
			require.True(t, ok)
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			delivered := false
			cursor.beforeCompletionWait = func(waiter chan cursorDelivery) {
				select {
				case completion := <-waiter:
					waiter <- completion
					delivered = true
					cancel()
				case <-ctx.Done():
				}
			}
			page, err := pager.NextPage(ctx)
			require.True(t, delivered, "completion and cancellation must both be ready")
			require.Zero(t, page)
			require.ErrorIs(t, err, context.Canceled)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
			require.Equal(t, status, cosmosErr.StatusCode)
			require.Equal(t, 3.5, cosmosErr.RequestCharge)
			require.Equal(t, "44444444-4444-4444-8444-444444444444", cosmosErr.ActivityID)
			require.EqualValues(t, "position-2", cosmosErr.ETag)
			require.Equal(t, SessionToken("0:2"), cosmosErr.SessionToken)
			require.NotNil(t, cosmosErr.Diagnostics)
			require.GreaterOrEqual(t, cosmosErr.AttemptCount, uint32(1))
			require.False(t, pager.More())
			require.Nil(t, cursor.handle)
		})
	}
}

func TestNativeChangeFeedCancelledCheckpointWaitsForCleanup(t *testing.T) {
	service := newChangeFeedContractService(t, 1, false)
	clientContainer := service.container(t, "items")
	client := clientContainer.database.client
	pager := clientContainer.NewChangeFeedPager(NewFeedScopeForFullContainer(), NewChangeFeedStartFromBeginning(), nil)
	defer func() { require.NoError(t, pager.Close()) }()
	_, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	cursor, ok := pager.cursor.(*nativeChangeFeedCursor)
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
	defer restore()
	cursor.beforeCompletionWait = func(waiter chan cursorDelivery) {
		select {
		case completion := <-waiter:
			held <- pendingDelivery{waiter, completion}
			cancel()
		case <-ctx.Done():
		}
	}
	token, err := pager.ContinuationToken(ctx)
	require.Len(t, held, 1, "a real admitted checkpoint must be withheld")
	require.Empty(t, token)
	require.ErrorIs(t, err, context.Canceled)
	require.False(t, pager.More())
	cleanupStarted := make(chan struct{})
	client.driver.beforePendingWait = func() { close(cleanupStarted) }
	closed := make(chan error, 1)
	go func() { closed <- client.Close() }()
	select {
	case <-cleanupStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("client did not reach pending change-feed cleanup")
	}
	select {
	case err := <-closed:
		t.Fatalf("client closed before checkpoint delivery: %v", err)
	case <-time.After(20 * time.Millisecond):
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
