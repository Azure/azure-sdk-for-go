// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEagerConstructorCancellationDrainsRuntimeAttachment(t *testing.T) {
	shared, err := NewRuntime(nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shared.Close()) })
	credential := &gatedFailingTokenCredential{
		started: make(chan struct{}), release: make(chan struct{}),
	}
	unblock := sync.OnceFunc(func() { close(credential.release) })
	t.Cleanup(unblock)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	type outcome struct {
		client *Client
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		client, err := NewClient(ctx, "https://construction.documents.azure.com", credential, &ClientOptions{Runtime: shared})
		done <- outcome{client, err}
	}()
	select {
	case <-credential.started:
	case <-time.After(5 * time.Second):
		t.Fatal("constructor did not start native initialization")
	}
	require.NoError(t, shared.Close())
	cancel()
	select {
	case result := <-done:
		require.Nil(t, result.client)
		require.ErrorIs(t, result.err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("constructor blocked on abandoned native work")
	}
	shared.mu.Lock()
	attached, native := shared.clients, shared.native
	shared.mu.Unlock()
	require.Equal(t, 1, attached, "native work must retain its runtime attachment")
	require.NotNil(t, native)
	unblock()
	require.Eventually(t, func() bool {
		shared.mu.Lock()
		defer shared.mu.Unlock()
		return shared.clients == 0 && shared.native == nil
	}, 5*time.Second, time.Millisecond, "late completion must release the runtime attachment")
}

func newMetadataServer(t *testing.T, intercept func(http.ResponseWriter, *http.Request) bool) *httptest.Server {
	t.Helper()
	var endpoint string
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() == "" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if intercept != nil && intercept(w, r) {
			return
		}
		if r.URL.Path == "/" {
			_, _ = fmt.Fprintf(w, `{
				"id":"synthetic","_rid":"synthetic.documents.azure.com","_self":"",
				"media":"//media/","addresses":"//addresses/","_dbs":"//dbs/",
				"readableLocations":[{"name":"West US","databaseAccountEndpoint":%q}],
				"writableLocations":[{"name":"West US","databaseAccountEndpoint":%q}],
				"userConsistencyPolicy":{"defaultConsistencyLevel":"Session"},
				"userReplicationPolicy":{"asyncReplication":false,"minReplicaSetSize":3,"maxReplicasetSize":4},
				"systemReplicationPolicy":{"asyncReplication":false,"minReplicaSetSize":3,"maxReplicasetSize":4},
				"readPolicy":{"primaryReadCoefficient":1,"secondaryReadCoefficient":1},
				"queryEngineConfiguration":"{}","enableMultipleWriteLocations":false
			}`, endpoint, endpoint)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/pkranges") {
			w.Header().Set("ETag", `"topology"`)
			if r.Header.Get("If-None-Match") == `"topology"` {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			_, _ = fmt.Fprint(w, `{"PartitionKeyRanges":[{"id":"0","minInclusive":"","maxExclusive":"FF","parents":[]}],"_count":1}`)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/colls/items") {
			_, _ = fmt.Fprint(w, `{"id":"items","_rid":"AAAAAIAAAAA=","_etag":"synthetic","partitionKey":{"paths":["/pk"],"kind":"Hash","version":2}}`)
			return
		}
		w.Header().Set("x-ms-request-charge", "2.5")
		w.Header().Set("x-ms-activity-id", "synthetic-metadata")
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"code":"NotFound","message":"synthetic missing container"}`)
	}))
	server.Config.Protocols = new(http.Protocols)
	server.Config.Protocols.SetHTTP1(true)
	server.Config.Protocols.SetUnencryptedHTTP2(true)
	endpoint = "http://" + server.Listener.Addr().String() + "/"
	server.Start()
	t.Cleanup(server.Close)
	return server
}

func TestEagerConstructorsAndContainerResolution(t *testing.T) {
	var accountReads, containerReads atomic.Int32
	server := newMetadataServer(t, func(_ http.ResponseWriter, r *http.Request) bool {
		if r.URL.Path == "/" {
			accountReads.Add(1)
		} else {
			containerReads.Add(1)
		}
		return false
	})
	endpoint := server.URL
	for _, token := range []bool{false, true} {
		t.Run(fmt.Sprint(token), func(t *testing.T) {
			shared, err := NewRuntime(nil)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, shared.Close()) })
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			var client *Client
			if token {
				client, err = NewClient(ctx, endpoint, fakeTokenCredential{}, &ClientOptions{Runtime: shared})
			} else {
				client, err = NewClientWithKey(ctx, endpoint, mustKeyCredential(t), &ClientOptions{Runtime: shared})
			}

			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			require.True(t, client.driver.created)
			require.Positive(t, accountReads.Load())
			before := containerReads.Load()
			database, err := client.NewDatabase("db")
			require.NoError(t, err)
			require.Equal(t, before, containerReads.Load(), "NewDatabase must remain local")
			for _, factory := range []func() (*ContainerClient, error){
				func() (*ContainerClient, error) { return client.NewContainer(ctx, "db", "missing") },
				func() (*ContainerClient, error) { return database.NewContainer(ctx, "missing") },
			} {
				container, err := factory()
				require.Nil(t, container)
				var failure *Error
				require.ErrorAs(t, err, &failure)
				require.Equal(t, CodeNotFound, failure.Code)
				require.Equal(t, 404, failure.StatusCode)
			}
			require.Greater(t, containerReads.Load(), before)
			container, err := client.NewContainer(ctx, "db", "items")
			require.NoError(t, err)
			require.Equal(t, "items", container.ID())
			before = containerReads.Load()
			again, err := database.NewContainer(ctx, "items")
			require.NoError(t, err)
			require.Equal(t, container.ID(), again.ID())
			require.Equal(t, before, containerReads.Load(), "warm metadata and topology must be reused")
		})
	}
}

func TestEagerCancellationReleasesLateSuccessfulHandles(t *testing.T) {
	for _, path := range []string{"/", "/dbs/db/colls/items"} {
		t.Run(path, func(t *testing.T) {
			entered, release := make(chan struct{}), make(chan struct{})
			var once sync.Once
			unblock := sync.OnceFunc(func() { close(release) })
			server := newMetadataServer(t, func(_ http.ResponseWriter, r *http.Request) bool {
				if r.URL.Path == path {
					once.Do(func() { close(entered) })
					<-release
				}
				return false
			})
			shared, err := NewRuntime(nil)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, shared.Close()) })
			t.Cleanup(unblock)
			var client *Client
			if path != "/" {
				client, err = NewClientWithKey(t.Context(), server.URL, mustKeyCredential(t), &ClientOptions{Runtime: shared})
				require.NoError(t, err)
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			result := make(chan error, 1)
			go func() {
				if path == "/" {
					created, err := NewClientWithKey(ctx, server.URL, KeyCredential{accountKey: testAccountKey}, &ClientOptions{Runtime: shared})
					if created != nil {
						_ = created.Close()
						result <- fmt.Errorf("cancelled constructor returned a client")
						return
					}
					result <- err
				} else {
					container, err := client.NewContainer(ctx, "db", "items")
					if container != nil {
						result <- fmt.Errorf("cancelled factory returned a container")
						return
					}
					result <- err
				}
			}()
			select {
			case <-entered:
			case <-time.After(5 * time.Second):
				t.Fatal("native metadata request did not arrive")
			}
			cancel()
			select {
			case err := <-result:
				require.ErrorIs(t, err, context.Canceled)
			case <-time.After(time.Second):
				t.Fatal("metadata wait ignored cancellation")
			}
			require.NoError(t, shared.Close())
			closed := make(chan error, 1)
			if client != nil {
				go func() { closed <- client.Close() }()
			}
			unblock()
			if client != nil {
				select {
				case err := <-closed:
					require.NoError(t, err)
				case <-time.After(5 * time.Second):
					t.Fatal("client did not drain late container resolution")
				}
			}
			require.Eventually(t, func() bool {
				shared.mu.Lock()
				defer shared.mu.Unlock()
				return shared.clients == 0 && shared.native == nil
			}, 5*time.Second, time.Millisecond)
		})
	}
}

func TestNativeBinaryEncodingLayerPrecedence(t *testing.T) {
	for _, tt := range []struct {
		name                     string
		env                      string
		runtime, client, request *BinaryEncodingOptions
		want                     bool
	}{
		{name: "environment disables", env: "false"},
		{name: "environment enables", env: "true", want: true},
		{name: "runtime overrides environment", env: "false", runtime: &BinaryEncodingOptions{}, want: true},
		{name: "client overrides runtime", env: "true", runtime: &BinaryEncodingOptions{}, client: &BinaryEncodingOptions{Disabled: true}},
		{name: "request overrides client", env: "false", client: &BinaryEncodingOptions{Disabled: true}, request: &BinaryEncodingOptions{}, want: true},
		{name: "request disables", env: "true", client: &BinaryEncodingOptions{}, request: &BinaryEncodingOptions{Disabled: true}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED", tt.env)
			var observed atomic.Bool
			var requests atomic.Int32
			server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
				if strings.Contains(r.URL.Path, "/docs/") {
					observed.Store(strings.Contains(r.Header.Get("x-ms-cosmos-supported-serialization-formats"), "CosmosBinary"))
					requests.Add(1)
					_, _ = fmt.Fprint(w, `{"id":"one","pk":"one"}`)
					return true
				}
				return false
			})
			shared, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{BinaryEncoding: tt.runtime}})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, shared.Close()) })
			t.Setenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED", fmt.Sprint(tt.env != "true"))
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			client, err := NewClientWithKey(ctx, server.URL, mustKeyCredential(t), &ClientOptions{
				Runtime: shared, Operation: OperationOptions{BinaryEncoding: tt.client},
			})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			container, err := client.NewContainer(ctx, "db", "items")
			require.NoError(t, err)
			response, err := container.ReadItem(ctx, NewPartitionKeyString("one"), "one", &ReadItemOptions{Operation: OperationOptions{BinaryEncoding: tt.request}})
			require.NoError(t, err)
			require.JSONEq(t, `{"id":"one","pk":"one"}`, string(response.Value))
			require.Positive(t, requests.Load())
			require.Equal(t, tt.want, observed.Load(), "must assert actual native wire negotiation")
		})
	}
}
