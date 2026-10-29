// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEmulatorCompletionHandoffOwnsHandles(t *testing.T) {
	client, databaseID, containerID := emulatorClient(t)
	for _, driver := range []bool{true, false} {
		for _, discard := range []bool{false, true} {
			t.Run(fmt.Sprintf("driver=%t/discard=%t", driver, discard), func(t *testing.T) {
				result, err := client.driver.captureHandleCompletionForTest(driver, databaseID, containerID)
				require.NoError(t, err)
				defer result.release()
				owned := result.native
				require.NotNil(t, owned)
				require.NotNil(t, owned.record.backing, "fixture must own real Rust backing")
				if driver {
					require.NotNil(t, owned.record.driver)
				} else {
					require.NotNil(t, owned.record.container)
				}
				if discard {
					pending := &pendingOperation{result: make(chan completionResult, 1), delivered: make(chan struct{})}
					_, ok := pending.abandon()
					require.False(t, ok)
					pending.deliver(result)
					<-pending.delivered
					require.Empty(t, pending.result)
				} else {
					prepared := result.prepare()
					require.NoError(t, prepared.err)
					require.Nil(t, prepared.native)
					if driver {
						require.NotNil(t, prepared.driver)
					} else {
						require.NotNil(t, prepared.container)
					}
					prepared.release()
				}
				require.Nil(t, owned.record.backing)
				require.Nil(t, owned.record.driver)
				require.Nil(t, owned.record.container)
				// Neither a stale result copy nor the cleared source may free the handles again.
				result.release()
			})
		}
	}
}

func TestEmulatorCompletionHandoffConcurrentResponsesSurviveClose(t *testing.T) {
	for _, verbosity := range []DiagnosticsVerbosity{DiagnosticsVerbositySummary, DiagnosticsVerbosityDetailed} {
		t.Run(verbosity.String(), func(t *testing.T) {
			client, databaseID, containerID := emulatorClientWithOptions(t, &ClientOptions{DiagnosticsVerbosity: verbosity})
			container, err := client.NewContainer(t.Context(), databaseID, containerID)
			require.NoError(t, err)
			type snapshot struct {
				Read    ItemResponse
				Batch   TransactionalBatchResponse
				Missing *Error
				before  []byte
			}
			const callers = 12
			saved := make([][]snapshot, callers)
			t.Run("concurrent", func(t *testing.T) {
				for i := range callers {
					t.Run(fmt.Sprint(i), func(t *testing.T) {
						t.Parallel()
						id := uniqueItemID(t)
						pk := NewPartitionKeyString(id)
						trackEmulatorItem(t, container, pk, id)
						body := json.RawMessage(fmt.Sprintf(`{"id":%q,"pk":%q,"value":1}`, id, id))
						for range 3 {
							created, err := container.CreateItem(t.Context(), pk, id, body, nil)
							require.NoError(t, err)
							require.Equal(t, http.StatusCreated, created.StatusCode)
							read, err := container.ReadItem(t.Context(), pk, id, &ReadItemOptions{SessionToken: created.SessionToken})
							require.NoError(t, err)
							require.NotEmpty(t, read.ETag)
							require.NotEmpty(t, read.SessionToken)
							require.Positive(t, read.RequestCharge)
							require.NotEmpty(t, read.ActivityID)
							requireEmulatorDiagnostics(t, read.Diagnostics, http.StatusOK, 0, false)
							var document map[string]any
							require.NoError(t, json.Unmarshal(read.Value, &document))
							require.Equal(t, id, document["id"])
							require.Equal(t, float64(1), document["value"])
							var patch PatchOperations
							require.NoError(t, patch.AppendSet("/value", 2))
							_, err = container.PatchItem(t.Context(), pk, id, patch, &PatchItemOptions{Strategy: PatchStrategyServerSide})
							require.NoError(t, err)
							_, err = container.ReplaceItem(t.Context(), pk, id, body, nil)
							require.NoError(t, err)
							_, err = container.UpsertItem(t.Context(), pk, id, body, nil)
							require.NoError(t, err)
							batch := NewTransactionalBatch(pk)
							require.NoError(t, batch.ReadItem(id, nil))
							require.NoError(t, batch.DeleteItem(id, nil))
							response, err := container.ExecuteTransactionalBatch(t.Context(), batch, &TransactionalBatchOptions{
								Operation: OperationOptions{EnableContentResponseOnWrite: to(true)},
							})
							require.NoError(t, err)
							require.Equal(t, http.StatusOK, response.StatusCode)
							require.Len(t, response.OperationResults, 2)
							require.Equal(t, http.StatusOK, response.OperationResults[0].StatusCode)
							require.Equal(t, http.StatusNoContent, response.OperationResults[1].StatusCode)
							require.NotEmpty(t, response.OperationResults[0].ETag)
							require.NoError(t, json.Unmarshal(response.OperationResults[0].ResourceBody, &document))
							require.Equal(t, id, document["id"])
							require.Positive(t, response.RequestCharge)
							_, err = container.ReadItem(t.Context(), pk, id, nil)
							missing := requireWireErrorDetails(t, err, CodeNotFound, http.StatusNotFound)
							requireEmulatorDiagnostics(t, missing.Diagnostics, http.StatusNotFound, 0, true)
							value := snapshot{Read: read, Batch: response, Missing: missing}
							value.before, err = json.Marshal(value)
							require.NoError(t, err)
							saved[i] = append(saved[i], value)
						}
					})
				}
			})
			require.NoError(t, client.Close())
			for _, snapshots := range saved {
				require.Len(t, snapshots, 3)
				for _, value := range snapshots {
					after, err := json.Marshal(value)
					require.NoError(t, err)
					require.Equal(t, value.before, after, "Go-owned data must survive subsequent polls and native teardown")
				}
			}
		})
	}
}

func TestEmulatorCompletionHandoffCloseWaitsForPreparation(t *testing.T) {
	client, databaseID, containerID := emulatorClient(t)
	container, err := client.NewContainer(t.Context(), databaseID, containerID)
	require.NoError(t, err)
	id := uniqueItemID(t)
	started, resume := make(chan struct{}), make(chan struct{})
	restore := stubBlockingDiagnosticsCopy(started, resume)
	defer restore()
	var once sync.Once
	defer once.Do(func() { close(resume) })
	// A missing read exercises a real completion and the client guard without leaving test data.
	done := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := container.ReadItem(ctx, NewPartitionKeyString(id), id, nil)
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("caller did not reach completion preparation")
	}
	closed := make(chan error, 1)
	go func() { closed <- client.Close() }()
	select {
	case err := <-closed:
		t.Fatalf("Close returned while the caller was still preparing a completion: %v", err)
	case <-time.After(25 * time.Millisecond):
	}
	once.Do(func() { close(resume) })
	var cosmosErr *Error
	select {
	case err := <-done:
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeNotFound, cosmosErr.Code)
	case <-time.After(10 * time.Second):
		t.Fatal("read did not finish after preparation resumed")
	}
	select {
	case err := <-closed:
		require.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("Close did not finish after preparation resumed")
	}
}
