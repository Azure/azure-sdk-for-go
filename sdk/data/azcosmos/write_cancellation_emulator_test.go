// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"encoding/json"
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

func TestEmulatorTransactionalBatchCancellationDrainsBeforeClose(t *testing.T) {
	for _, rollback := range []bool{false, true} {
		name := "commit"
		if rollback {
			name = "rollback"
		}
		t.Run(name, func(t *testing.T) {
			independent := emulatorContainer(t)
			id := uniqueItemID(t)
			pk := NewPartitionKeyString(id)
			body := []byte(fmt.Sprintf(`{"id":%q,"pk":%q,"value":1}`, id, id))
			trackEmulatorItem(t, independent, pk, id)
			started, unblock := make(chan struct{}), make(chan struct{})
			var startOnce, unblockOnce sync.Once
			gated := queryContractContainer(t, func(_ http.ResponseWriter, r *http.Request) bool {
				if r.Method == http.MethodPost && strings.EqualFold(r.Header.Get("x-ms-cosmos-is-batch-request"), "true") {
					startOnce.Do(func() { close(started) })
					<-unblock
				}
				return false
			}, nil)
			t.Cleanup(func() { unblockOnce.Do(func() { close(unblock) }) })
			batch := NewTransactionalBatch(pk)
			require.NoError(t, batch.CreateItem(body, nil))
			if rollback {
				require.NoError(t, batch.ReadItem(id+"-missing", nil))
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			type outcome struct {
				response TransactionalBatchResponse
				err      error
			}
			results := make(chan outcome, 1)
			go func() {
				response, err := gated.ExecuteTransactionalBatch(ctx, batch, &TransactionalBatchOptions{
					Operation: OperationOptions{EndToEndTimeout: to(30 * time.Second)},
				})
				results <- outcome{response, err}
			}()
			select {
			case <-started:
			case <-time.After(10 * time.Second):
				t.Fatal("batch request did not reach the service boundary")
			}
			cancel()
			select {
			case result := <-results:
				require.Zero(t, result.response)
				require.ErrorIs(t, result.err, context.Canceled)
				var cosmosErr *Error
				require.ErrorAs(t, result.err, &cosmosErr)
				require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
			case <-time.After(time.Second):
				t.Fatal("admitted batch did not stop waiting after cancellation")
			}
			draining := make(chan struct{})
			client := gated.database.client
			client.driver.beforePendingWait = func() { close(draining) }
			closed := make(chan error, 1)
			go func() { closed <- client.Close() }()
			select {
			case <-draining:
			case <-time.After(5 * time.Second):
				t.Fatal("client teardown did not reach pending batch cleanup")
			}
			select {
			case err := <-closed:
				t.Fatalf("client closed before native batch completion: %v", err)
			case <-time.After(50 * time.Millisecond):
			}
			unblockOnce.Do(func() { close(unblock) })
			select {
			case err := <-closed:
				require.NoError(t, err)
			case <-time.After(10 * time.Second):
				t.Fatal("client did not finish draining the cancelled batch")
			}
			response, err := independent.ReadItem(t.Context(), pk, id, nil)
			if rollback {
				require.Zero(t, response)
				var cosmosErr *Error
				require.ErrorAs(t, err, &cosmosErr)
				require.Equal(t, CodeNotFound, cosmosErr.Code)
			} else {
				require.NoError(t, err, "cancelling the Go wait must not cancel native batch execution")
				require.Equal(t, http.StatusOK, response.StatusCode)
				var persisted map[string]any
				require.NoError(t, json.Unmarshal(response.Value, &persisted))
				require.Equal(t, id, persisted["id"])
				require.Equal(t, float64(1), persisted["value"])
			}
		})
	}
}
