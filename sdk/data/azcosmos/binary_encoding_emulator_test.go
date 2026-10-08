// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore documentdb isquery

package azcosmos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmulatorQueryBinaryWirePreferenceAndTextResults(t *testing.T) {
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	target, err := url.Parse(endpoint)
	require.NoError(t, err)
	for _, tt := range []struct {
		name    string
		client  *BinaryEncodingOptions
		request *BinaryEncodingOptions
		want    bool
	}{
		{"default", nil, nil, true},
		{"client disabled", &BinaryEncodingOptions{Enabled: to(false)}, nil, false},
		{"request enables", &BinaryEncodingOptions{Enabled: to(false)}, &BinaryEncodingOptions{}, true},
		{"request disables", &BinaryEncodingOptions{}, &BinaryEncodingOptions{Enabled: to(false)}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var mu sync.Mutex
			var negotiated []bool
			var binaryResponses int
			var proxyURL string
			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ModifyResponse = func(r *http.Response) error {
				if err := rewriteQueryAccountResponse(r, proxyURL); err != nil {
					return err
				}
				if isQueryDataRequest(r.Request) && strings.EqualFold(r.Request.Header.Get("x-ms-documentdb-isquery"), "true") && r.StatusCode == http.StatusOK {
					body, err := io.ReadAll(r.Body)
					closeErr := r.Body.Close()
					if err != nil {
						return err
					}
					if closeErr != nil {
						return closeErr
					}
					r.Body = io.NopCloser(bytes.NewReader(body))
					if len(body) > 0 && body[0] == 0x80 {
						mu.Lock()
						binaryResponses++
						mu.Unlock()
					}
				}
				return nil
			}
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if isQueryDataRequest(r) && strings.EqualFold(r.Header.Get("x-ms-documentdb-isquery"), "true") {
					mu.Lock()
					negotiated = append(negotiated, strings.Contains(r.Header.Get("x-ms-cosmos-supported-serialization-formats"), "CosmosBinary"))
					mu.Unlock()
				}
				proxy.ServeHTTP(w, r)
			}))
			proxyURL = "http://" + server.Listener.Addr().String()
			server.Start()
			t.Cleanup(server.Close)
			client, err := NewClientWithKey(server.URL, KeyCredential{accountKey: emulatorKey},
				&ClientOptions{BinaryEncoding: tt.client})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			container, err := client.NewContainer(databaseID, containerID)
			require.NoError(t, err)
			id := uniqueItemID(t)
			key := NewPartitionKeyString(id)
			_, err = container.CreateItem(t.Context(), key, id, []byte(fmt.Sprintf(`{"id":%q,"pk":%q,"value":42}`, id, id)), nil)
			require.NoError(t, err)
			trackEmulatorItem(t, container, key, id)
			query, err := NewQuery("SELECT VALUE c.value FROM c WHERE c.id = @id").WithParameter("@id", id)
			require.NoError(t, err)
			for _, scope := range []FeedScope{NewFeedScopeForPartitionKey(key), NewFeedScopeForFullContainer()} {
				pager := container.NewQueryItemsPager(query, scope, &QueryOptions{
					Operation: OperationOptions{BinaryEncoding: tt.request}, Feed: FeedOptions{PageSizeHint: 1},
				})
				var items [][]byte
				for pageCount := 0; pager.More(); pageCount++ {
					require.Less(t, pageCount, 20)
					page, err := pager.NextPage(t.Context())
					require.NoError(t, err)
					for _, item := range page.Items {
						require.True(t, json.Valid(item))
					}
					items = append(items, page.Items...)
				}
				require.NoError(t, pager.Close())
				require.Equal(t, [][]byte{[]byte("42")}, items)
			}
			mu.Lock()
			defer mu.Unlock()
			require.NotEmpty(t, negotiated)
			for _, binary := range negotiated {
				require.Equal(t, tt.want, binary, "wire negotiation must follow client/request preference")
			}
			if tt.want {
				require.Positive(t, binaryResponses, "must observe binary from the service, not merely request it")
			} else {
				require.Zero(t, binaryResponses)
			}
		})
	}
}
