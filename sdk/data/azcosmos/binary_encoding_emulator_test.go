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
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmulatorBinaryNumericValues(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		t.Run(fmt.Sprintf("binary=%t", enabled), func(t *testing.T) {
			var binaryResponses atomic.Int32
			container := queryContractContainer(t, nil, func(response *http.Response) error {
				if response.Request.Method != http.MethodGet || !strings.Contains(response.Request.URL.Path, "/docs/") ||
					response.StatusCode != http.StatusOK {
					return nil
				}
				body, err := io.ReadAll(response.Body)
				closeErr := response.Body.Close()
				if err != nil {
					return err
				}
				if closeErr != nil {
					return closeErr
				}
				response.Body = io.NopCloser(bytes.NewReader(body))
				if len(body) > 0 && body[0] == 0x80 {
					binaryResponses.Add(1)
				}
				return nil
			})
			options := OperationOptions{BinaryEncoding: &BinaryEncodingOptions{Enabled: to(enabled)}}
			for _, number := range []string{
				"9007199254740993", "9223372036854775807", "-9223372036854775808",
				"0.1", "0.125", "-1234.5678", "1e-6",
			} {
				t.Run(number, func(t *testing.T) {
					id := uniqueItemID(t)
					key := NewPartitionKeyString(id)
					body := json.RawMessage(fmt.Sprintf(`{"id":%q,"pk":%q,"number":%s}`, id, id, number))
					_, err := container.CreateItem(t.Context(), key, id, body, &CreateItemOptions{Operation: options})
					require.NoError(t, err)
					trackEmulatorItem(t, container, key, id)
					response, err := container.ReadItem(t.Context(), key, id, &ReadItemOptions{Operation: options})
					require.NoError(t, err)
					var item struct {
						Number json.Number `json:"number"`
					}
					require.NoError(t, json.Unmarshal(response.Value, &item))
					want, ok := new(big.Rat).SetString(number)
					require.True(t, ok)
					got, ok := new(big.Rat).SetString(item.Number.String())
					require.True(t, ok, "response must contain a JSON number")
					require.Zero(t, want.Cmp(got), "numeric value changed: input=%s output=%s", number, item.Number)
				})
			}
			if enabled {
				require.Positive(t, binaryResponses.Load(), "must exercise actual binary-to-text conversion")
			} else {
				require.Zero(t, binaryResponses.Load())
			}
		})
	}
}

func TestEmulatorQueryBinaryWirePreferenceAndTextResults(t *testing.T) {
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	t.Setenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED", "")
	require.NoError(t, os.Unsetenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED"))
	target, err := url.Parse(endpoint)
	require.NoError(t, err)
	for _, tt := range []struct {
		name    string
		client  *BinaryEncodingOptions
		request *BinaryEncodingOptions
		want    bool
		runtime *BinaryEncodingOptions
		common  *BinaryEncodingOptions
	}{
		{name: "default", want: true},
		{name: "client disabled", client: &BinaryEncodingOptions{Enabled: to(false)}},
		{name: "request enables", client: &BinaryEncodingOptions{Enabled: to(false)}, request: &BinaryEncodingOptions{}, want: true},
		{name: "request disables", client: &BinaryEncodingOptions{}, request: &BinaryEncodingOptions{Enabled: to(false)}},
		{name: "runtime disables", runtime: &BinaryEncodingOptions{Enabled: to(false)}},
		{name: "common client disables", common: &BinaryEncodingOptions{Enabled: to(false)}},
		{name: "common client overrides runtime", runtime: &BinaryEncodingOptions{Enabled: to(false)}, common: &BinaryEncodingOptions{}, want: true},
		{name: "dedicated overrides common", common: &BinaryEncodingOptions{Enabled: to(false)}, client: &BinaryEncodingOptions{}, want: true},
		{name: "request overrides runtime", runtime: &BinaryEncodingOptions{Enabled: to(false)}, request: &BinaryEncodingOptions{}, want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var mu sync.Mutex
			var negotiated []bool
			var itemNegotiated []bool
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
				if strings.Contains(r.URL.Path, "/docs/") && r.Method == http.MethodGet {
					mu.Lock()
					itemNegotiated = append(itemNegotiated, strings.Contains(r.Header.Get("x-ms-cosmos-supported-serialization-formats"), "CosmosBinary"))
					mu.Unlock()
				}
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
			runtime, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{BinaryEncoding: tt.runtime}})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, runtime.Close()) })
			client, err := NewClientWithKey(server.URL, KeyCredential{accountKey: emulatorKey},
				&ClientOptions{Runtime: runtime, BinaryEncoding: tt.client, Operation: OperationOptions{BinaryEncoding: tt.common}})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			container, err := client.NewContainer(databaseID, containerID)
			require.NoError(t, err)
			id := uniqueItemID(t)
			key := NewPartitionKeyString(id)
			_, err = container.CreateItem(t.Context(), key, id, json.RawMessage(fmt.Sprintf(`{"id":%q,"pk":%q,"value":42}`, id, id)), nil)
			require.NoError(t, err)
			trackEmulatorItem(t, container, key, id)
			item, err := container.ReadItem(t.Context(), key, id, &ReadItemOptions{Operation: OperationOptions{BinaryEncoding: tt.request}})
			require.NoError(t, err)
			require.True(t, json.Valid(item.Value))
			query, err := NewQuery("SELECT VALUE c.value FROM c WHERE c.id = @id").WithParameter("@id", id)
			require.NoError(t, err)
			for _, scope := range []FeedScope{NewFeedScopeForPartitionKey(key), NewFeedScopeForFullContainer()} {
				pager := container.NewQueryItemsPager[json.RawMessage](query, scope, &QueryOptions{
					Operation: OperationOptions{BinaryEncoding: tt.request}, Feed: FeedOptions{PageSizeHint: 1},
				})
				var items []json.RawMessage
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
				require.Equal(t, []json.RawMessage{[]byte("42")}, items)
			}
			mu.Lock()
			defer mu.Unlock()
			require.NotEmpty(t, negotiated)
			require.NotEmpty(t, itemNegotiated)
			for _, binary := range itemNegotiated {
				require.Equal(t, tt.want, binary, "item wire negotiation must follow inherited preferences")
			}
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
