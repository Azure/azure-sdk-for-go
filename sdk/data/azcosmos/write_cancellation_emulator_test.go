// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEmulatorSubmittedWritesStopWaitingAndDrain(t *testing.T) {
	for _, operation := range []string{"create", "replace", "upsert", "delete", "patch"} {
		t.Run(operation, func(t *testing.T) {
			started, unblock := make(chan struct{}), make(chan struct{})
			var once, startOnce sync.Once
			id := uniqueItemID(t)
			pk := NewPartitionKeyString(id)
			body := []byte(fmt.Sprintf(`{"id":%q,"pk":%q,"value":1}`, id, id))
			seed := emulatorContainer(t)
			_, err := seed.CreateItem(t.Context(), pk, id, body, nil)
			require.NoError(t, err)
			trackEmulatorItem(t, seed, pk, id)
			gated := queryContractContainer(t, func(w http.ResponseWriter, r *http.Request) bool {
				isWrite := r.Method == http.MethodPost || r.Method == http.MethodPut ||
					r.Method == http.MethodDelete || r.Method == http.MethodPatch
				if isWrite && strings.Contains(r.URL.Path, "/docs") {
					startOnce.Do(func() { close(started) })
					<-unblock
				}
				return false
			}, nil)
			t.Cleanup(func() { once.Do(func() { close(unblock) }) })
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			type result struct {
				response ItemResponse
				err      error
			}
			results := make(chan result, 1)
			var patch PatchOperations
			require.NoError(t, patch.AppendSet("/value", 2))
			go func() {
				var response ItemResponse
				var err error
				switch operation {
				case "create":
					response, err = gated.CreateItem(ctx, pk, id, body, nil)
				case "replace":
					response, err = gated.ReplaceItem(ctx, pk, id, body, nil)
				case "upsert":
					response, err = gated.UpsertItem(ctx, pk, id, body, nil)
				case "delete":
					response, err = gated.DeleteItem(ctx, pk, id, nil)
				case "patch":
					response, err = gated.PatchItem(ctx, pk, id, patch, &PatchItemOptions{Strategy: PatchStrategyServerSide})
				}
				results <- result{response, err}
			}()
			select {
			case <-started:
			case <-time.After(10 * time.Second):
				t.Fatal("write did not reach service boundary")
			}
			cancel()
			select {
			case result := <-results:
				require.Zero(t, result.response)
				require.ErrorIs(t, result.err, context.Canceled)
			case <-time.After(time.Second):
				t.Fatal("submitted write did not stop waiting")
			}
			closing := make(chan error, 1)
			go func() { closing <- gated.database.client.Close() }()
			select {
			case <-closing:
				t.Fatal("client freed resources before native write completed")
			case <-time.After(50 * time.Millisecond):
			}
			once.Do(func() { close(unblock) })
			select {
			case err := <-closing:
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				t.Fatal("client did not finish draining write")
			}
		})
	}
}
