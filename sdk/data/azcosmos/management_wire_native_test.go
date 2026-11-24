// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestManagementNativeDatabasePagingAndThroughput(t *testing.T) {
	var mu sync.Mutex
	var requests []string
	polls := 0
	server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
		if r.URL.Path == "/" {
			return false
		}

		mu.Lock()
		defer mu.Unlock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.Header().Set("x-ms-request-charge", "2")
		w.Header().Set("x-ms-activity-id", "management-activity")
		switch {
		case r.URL.Path == "/dbs" && r.Method == http.MethodPost && r.Header.Get("x-ms-documentdb-isquery") == "":
			require.Equal(t, "400", r.Header.Get("x-ms-offer-throughput"))
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"db","_rid":"db-rid"}`)
		case r.URL.Path == "/dbs":
			require.Equal(t, "1", r.Header.Get("x-ms-max-item-count"))
			if r.Header.Get("x-ms-continuation") == "" {
				w.Header().Set("x-ms-continuation", "second")
				_, _ = io.WriteString(w, `{"Databases":[{"id":"one","_rid":"one-rid"}],"_count":1}`)
			} else {
				require.Equal(t, "second", r.Header.Get("x-ms-continuation"))
				_, _ = io.WriteString(w, `{"Databases":[{"id":"two","_rid":"two-rid"}],"_count":1}`)
			}
		case r.URL.Path == "/dbs/db":
			_, _ = io.WriteString(w, `{"id":"db","_rid":"db-rid"}`)
		case r.URL.Path == "/offers":
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.JSONEq(t, `{"query":"SELECT * FROM c WHERE c.offerResourceId = @rid","parameters":[{"name":"@rid","value":"db-rid"}]}`, string(body))
			_, _ = io.WriteString(w, `{"Offers":[{"id":"offer-id","offerResourceId":"db-rid","offerVersion":"V2","content":{"offerThroughput":400}}],"_count":1}`)
		case r.URL.Path == "/offers/offer-id" && r.Method == http.MethodPut:
			require.Equal(t, "inherited", r.Header.Get("x-test-inherited"))
			require.Equal(t, `"offer-current"`, r.Header.Get("If-Match"))
			var body map[string]json.RawMessage
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.JSONEq(t, `{"offerThroughput":600}`, string(body["content"]))
			w.Header().Set("x-ms-offer-replace-pending", "true")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{"id":"offer-id","content":{"offerThroughput":400}}`)
		case r.URL.Path == "/offers/offer-id":
			polls++
			require.Empty(t, r.Header.Get("If-Match"), "replacement preconditions must not leak into polling")
			w.Header().Set("ETag", `"offer-etag"`)
			_, _ = io.WriteString(w, `{"id":"offer-id","content":{"offerThroughput":600}}`)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, `{"code":"NotFound","message":"unexpected test request"}`)
		}
		return true
	})
	client, err := NewClientWithKey(t.Context(), server.URL, mustKeyCredential(t), &ClientOptions{
		Operation: OperationOptions{CustomHeaders: map[string]string{"x-test-inherited": "inherited"}},
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	_, err = client.CreateDatabase(ctx, DatabaseProperties{ID: "db"}, &CreateDatabaseOptions{
		Operation:  OperationOptions{CustomHeaders: map[string]string{}},
		Throughput: &ThroughputProperties{ManualThroughput: to(uint64(400))},
	})
	require.NoError(t, err)
	pager := client.NewQueryDatabasesPager(NewQuery("SELECT * FROM c"), &ManagementQueryOptions{
		PageSizeHint: 1, Operation: OperationOptions{EndToEndTimeout: to(time.Second)},
	})
	defer pager.Close()
	first, err := pager.NextPage(ctx)
	require.NoError(t, err)
	require.Len(t, first.Items, 1)
	token, err := pager.ContinuationToken(ctx)
	require.NoError(t, err)
	require.Equal(t, "second", token)
	time.Sleep(1100 * time.Millisecond)
	ids := []string{first.Items[0].ID}
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		for _, database := range page.Items {
			ids = append(ids, database.ID)
		}
	}
	require.Equal(t, []string{"one", "two"}, ids)
	resumed := client.NewQueryDatabasesPager(NewQuery("SELECT * FROM c"), &ManagementQueryOptions{PageSizeHint: 1, ContinuationToken: token})
	resumedPage, err := resumed.NextPage(ctx)
	require.NoError(t, err)
	require.Len(t, resumedPage.Items, 1)
	require.Equal(t, "two", resumedPage.Items[0].ID)
	require.NoError(t, resumed.Close())
	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	poller, err := database.BeginReplaceThroughput(ctx, ThroughputProperties{ManualThroughput: to(uint64(600))}, &BeginReplaceThroughputOptions{IfMatchETag: to(azcore.ETag(`"offer-current"`))})
	require.NoError(t, err)
	require.False(t, poller.Done(), "HTTP 200 plus the pending header is not terminal success")
	result, err := poller.Poll(ctx)
	require.NoError(t, err)
	require.True(t, poller.Done())
	require.Equal(t, uint64(600), *result.ThroughputProperties.ManualThroughput)
	require.Equal(t, float64(8), result.RequestCharge)
	require.Equal(t, `"offer-etag"`, string(result.ETag))
	mu.Lock()
	defer mu.Unlock()
	require.Equal(t, 1, polls)
}

func TestManagementNativeContainerReplaceAndSharedThroughput(t *testing.T) {
	var replacements atomic.Int32
	server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
		w.Header().Set("x-ms-request-charge", "1.5")
		w.Header().Set("x-ms-activity-id", "container-management")
		switch {
		case r.URL.Path == "/dbs/db/colls" && r.Header.Get("x-ms-documentdb-isquery") != "":
			_, _ = io.WriteString(w, `{"DocumentCollections":[{"id":"items","partitionKey":{"paths":["/pk"],"kind":"Hash","version":2}}],"_count":1}`)
		case r.URL.Path == "/dbs/db/colls" && r.Method == http.MethodPost:
			require.JSONEq(t, `{"maxThroughput":4000}`, r.Header.Get("x-ms-cosmos-offer-autopilot-settings"))
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"items","_rid":"AAAAAIAAAAA=","partitionKey":{"paths":["/pk"],"kind":"Hash","version":2}}`)
		case r.URL.Path == "/dbs/db/colls/items" && r.Method == http.MethodGet:
			_, _ = io.WriteString(w, `{"id":"items","_rid":"AAAAAIAAAAA=","partitionKey":{"paths":["/pk"],"kind":"Hash","systemKey":true},"indexingPolicy":{"automatic":true,"includedPaths":[{"path":"/*","indexes":[{"kind":"Range","dataType":"String","precision":-1}]}]}}`)
		case r.URL.Path == "/dbs/db/colls/items" && r.Method == http.MethodPut:
			replacements.Add(1)
			require.Equal(t, "inherited", r.Header.Get("x-test-inherited"))
			require.Equal(t, `"current"`, r.Header.Get("If-Match"))
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var properties ContainerProperties
			require.NoError(t, json.Unmarshal(body, &properties))
			require.Equal(t, int32(30), *properties.DefaultTimeToLive)
			require.Equal(t, 1, properties.PartitionKeyDefinition.Version)
			var fields map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(body, &fields))
			require.JSONEq(t, `{"paths":["/pk"],"kind":"Hash","version":1,"systemKey":true}`, string(fields["partitionKey"]))
			require.JSONEq(t, `{"automatic":true,"includedPaths":[{"path":"/*","indexes":[{"kind":"Range","dataType":"String","precision":-1}]}]}`, string(fields["indexingPolicy"]))
			w.Header().Set("ETag", `"updated"`)
			_, _ = w.Write(body)
		case r.URL.Path == "/offers":
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Contains(t, string(body), "AAAAAIAAAAA=")
			_, _ = io.WriteString(w, `{"Offers":[],"_count":0}`)
		default:
			return false
		}
		return true
	})
	client, err := NewClientWithKey(t.Context(), server.URL, mustKeyCredential(t), &ClientOptions{
		Operation: OperationOptions{CustomHeaders: map[string]string{"x-test-inherited": "inherited"}},
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	_, err = database.CreateContainer(t.Context(), ContainerProperties{ID: "items", PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}}}, &CreateContainerOptions{
		Operation:  OperationOptions{CustomHeaders: map[string]string{}},
		Throughput: &ThroughputProperties{Autoscale: &AutoscaleSettings{MaxThroughput: 4000}},
	})
	require.NoError(t, err)
	pager := database.NewQueryContainersPager(NewQuery("SELECT * FROM c"), nil)
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	require.Equal(t, "items", page.Items[0].ID)
	require.NoError(t, pager.Close())
	container, err := database.NewContainer(t.Context(), "items")
	require.NoError(t, err)
	properties, err := container.Read(t.Context(), nil)
	require.NoError(t, err)
	properties.ContainerProperties.DefaultTimeToLive = to(int32(30))
	replaced, err := container.Replace(t.Context(), properties.ContainerProperties, &ReplaceContainerOptions{IfMatchETag: to(azcore.ETag(`"current"`))})
	require.NoError(t, err)
	require.Equal(t, `"updated"`, string(replaced.ETag))
	require.Equal(t, int32(30), *replaced.ContainerProperties.DefaultTimeToLive)
	response, err := container.ReadThroughput(t.Context(), nil)
	require.NoError(t, err)
	require.Nil(t, response.ThroughputProperties, "shared container must not inherit its database offer")
	require.Positive(t, response.RequestCharge)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, "container-management", response.ActivityID)
	poller, err := container.BeginReplaceThroughput(t.Context(), ThroughputProperties{ManualThroughput: to(uint64(600))}, nil)
	require.Nil(t, poller)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, 20301, cosmosErr.SubStatus)
	require.Positive(t, cosmosErr.RequestCharge)
	require.Equal(t, "container-management", cosmosErr.ActivityID)
	require.Equal(t, int32(1), replacements.Load())
}

func TestManagementNativeErrorsPreserveChargeAndStatus(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusConflict, http.StatusBadRequest, http.StatusPreconditionFailed, http.StatusTooManyRequests} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
				if r.URL.Path != "/dbs" {
					return false
				}

				w.Header().Set("x-ms-request-charge", "3.25")
				w.Header().Set("x-ms-activity-id", "failed-management")
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, `{"code":"SyntheticError","message":"synthetic management rejection"}`)
				return true
			})
			client, err := NewClientWithKey(t.Context(), server.URL, mustKeyCredential(t), &ClientOptions{Operation: OperationOptions{
				ThrottlingRetry: ThrottlingRetryOptions{MaxRetryCount: to(uint32(0))},
			}})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			response, err := client.CreateDatabase(t.Context(), DatabaseProperties{ID: "db"}, nil)
			require.Equal(t, DatabaseResponse{}, response)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Equal(t, status, cosmosErr.StatusCode)
			require.Equal(t, 3.25, cosmosErr.RequestCharge)
			require.Equal(t, "failed-management", cosmosErr.ActivityID)
			require.True(t, cosmosErr.FromWire)
		})
	}
}

func TestManagementNativeThroughputHeaderInheritance(t *testing.T) {
	for _, tt := range []struct {
		name                     string
		runtime, client, request map[string]string
		wantHeader               string
		requiresExtension        bool
	}{
		{name: "no inherited headers"},
		{name: "client inheritance", client: map[string]string{"x-test-inherited": "client", "x-ms-offer-throughput": "999"}, wantHeader: "client", requiresExtension: true},
		{name: "runtime inheritance", runtime: map[string]string{"x-test-inherited": "runtime"}, wantHeader: "runtime", requiresExtension: true},
		{name: "client clears runtime", runtime: map[string]string{"x-test-inherited": "runtime"}, client: map[string]string{}},
		{name: "request clears inheritance", runtime: map[string]string{"x-test-inherited": "runtime"}, client: map[string]string{"x-test-inherited": "client"}, request: map[string]string{}},
		{name: "request replaces inheritance", runtime: map[string]string{"x-test-inherited": "runtime"}, client: map[string]string{"x-test-inherited": "client"}, request: map[string]string{"x-test-inherited": "request"}, wantHeader: "request"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var creates atomic.Int32
			server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
				if r.Method != http.MethodPost || (r.URL.Path != "/dbs" && r.URL.Path != "/dbs/db/colls") {
					return false
				}
				creates.Add(1)
				require.Equal(t, "400", r.Header.Get("x-ms-offer-throughput"))
				require.Equal(t, tt.wantHeader, r.Header.Get("x-test-inherited"))
				w.WriteHeader(http.StatusCreated)
				_, _ = io.WriteString(w, `{"id":"items","partitionKey":{"paths":["/pk"],"kind":"Hash","version":2}}`)
				return true
			})
			runtime, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{CustomHeaders: tt.runtime}})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, runtime.Close()) })
			client, err := NewClientWithKey(t.Context(), server.URL, mustKeyCredential(t), &ClientOptions{
				Runtime: runtime, Operation: OperationOptions{CustomHeaders: tt.client},
			})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			throughput := &ThroughputProperties{ManualThroughput: to(uint64(400))}
			database, err := client.NewDatabase("db")
			require.NoError(t, err)
			_, dbErr := client.CreateDatabase(t.Context(), DatabaseProperties{ID: "db"}, &CreateDatabaseOptions{
				Operation: OperationOptions{CustomHeaders: tt.request}, Throughput: throughput,
			})
			_, containerErr := database.CreateContainer(t.Context(), ContainerProperties{
				ID: "items", PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}},
			}, &CreateContainerOptions{Operation: OperationOptions{CustomHeaders: tt.request}, Throughput: throughput})
			if tt.requiresExtension && !nativeManagementThroughputAvailable {
				require.ErrorContains(t, dbErr, "unpublished native extension")
				require.ErrorContains(t, containerErr, "unpublished native extension")
				require.Zero(t, creates.Load(), "unsafe legacy injection must fail before data I/O")
			} else {
				require.NoError(t, dbErr)
				require.NoError(t, containerErr)
				require.Equal(t, int32(2), creates.Load())
			}
		})
	}
}
