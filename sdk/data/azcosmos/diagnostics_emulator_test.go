// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireEmulatorDiagnostics(t *testing.T, diagnostics *Diagnostics, status, subStatus int, failed bool) {
	t.Helper()
	require.NotNil(t, diagnostics)
	require.Equal(t, status, diagnostics.StatusCode)
	require.Equal(t, subStatus, diagnostics.SubStatus)
	require.True(t, diagnostics.Completed)
	require.Equal(t, failed, diagnostics.Failed)
	require.Equal(t, uint32(1), diagnostics.AttemptCount)
	require.False(t, diagnostics.Compacted)
	require.Len(t, diagnostics.Attempts, 1)
	require.Equal(t, status, diagnostics.Attempts[0].StatusCode)
	require.Equal(t, subStatus, diagnostics.Attempts[0].SubStatus)
	require.NotEmpty(t, diagnostics.Attempts[0].Endpoint)
	require.NotEmpty(t, diagnostics.RegionsContacted)
	require.True(t, json.Valid([]byte(diagnostics.JSON)))
	require.Positive(t, diagnostics.Elapsed)
}

func TestEmulatorDiagnosticsPointOperations(t *testing.T) {
	container := emulatorContainer(t)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	trackEmulatorItem(t, container, pk, id)
	body, err := json.Marshal(map[string]string{"id": id, "pk": id})
	require.NoError(t, err)

	created, err := container.CreateItem(t.Context(), pk, id, body, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, created.StatusCode)
	require.Zero(t, created.SubStatus)
	requireEmulatorDiagnostics(t, created.Diagnostics, http.StatusCreated, 0, false)
	require.Equal(t, created.AttemptCount, created.Diagnostics.AttemptCount)
	require.InDelta(t, created.RequestCharge, created.Diagnostics.TotalRequestCharge, 0.001)
	require.Positive(t, created.Diagnostics.Attempts[0].RequestCharge)

	read, err := container.ReadItem(t.Context(), pk, id, &ReadItemOptions{SessionToken: created.SessionToken})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, read.StatusCode)
	requireEmulatorDiagnostics(t, read.Diagnostics, http.StatusOK, 0, false)
	require.InDelta(t, read.RequestCharge, read.Diagnostics.TotalRequestCharge, 0.001)
	require.NotSame(t, created.Diagnostics, read.Diagnostics)

	_, err = container.DeleteItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	requireEmulatorDiagnostics(t, created.Diagnostics, http.StatusCreated, 0, false)
	require.True(t, json.Valid([]byte(read.Diagnostics.JSON)), "snapshots survive later native completions")
}

func TestEmulatorDiagnosticsWireFailures(t *testing.T) {
	container := emulatorContainer(t)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	body, err := json.Marshal(map[string]string{"id": id, "pk": id})
	require.NoError(t, err)

	_, err = container.ReadItem(t.Context(), pk, id, nil)
	missing := requireWireErrorDetails(t, err, CodeNotFound, http.StatusNotFound)
	requireEmulatorDiagnostics(t, missing.Diagnostics, http.StatusNotFound, 0, true)
	require.Equal(t, missing.AttemptCount, missing.Diagnostics.AttemptCount)
	require.InDelta(t, missing.RequestCharge, missing.Diagnostics.TotalRequestCharge, 0.001)

	trackEmulatorItem(t, container, pk, id)
	_, err = container.CreateItem(t.Context(), pk, id, body, nil)
	require.NoError(t, err)
	response, err := container.CreateItem(t.Context(), pk, id, body, nil)
	require.Equal(t, ItemResponse{}, response)
	conflict := requireWireErrorDetails(t, err, CodeConflict, http.StatusConflict)
	requireEmulatorDiagnostics(t, conflict.Diagnostics, http.StatusConflict, 0, true)
	require.Equal(t, conflict.AttemptCount, conflict.Diagnostics.AttemptCount)
	require.InDelta(t, conflict.RequestCharge, conflict.Diagnostics.TotalRequestCharge, 0.001)
	requireEmulatorDiagnostics(t, missing.Diagnostics, http.StatusNotFound, 0, true)
}

func TestEmulatorDiagnosticsQueryPage(t *testing.T) {
	container := emulatorContainer(t)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	trackEmulatorItem(t, container, pk, id)
	body, err := json.Marshal(map[string]string{"id": id, "pk": id})
	require.NoError(t, err)
	_, err = container.CreateItem(t.Context(), pk, id, body, nil)
	require.NoError(t, err)

	pager := container.NewQueryItemsPager(NewQuery("SELECT * FROM c"), NewFeedScopeForPartitionKey(pk), nil)
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	require.Equal(t, http.StatusOK, page.StatusCode)
	requireEmulatorDiagnostics(t, page.Diagnostics, http.StatusOK, 0, false)
	require.Equal(t, page.AttemptCount, page.Diagnostics.AttemptCount)
	// The first page's Response includes a separate metadata-validation read.
	require.GreaterOrEqual(t, page.RequestCharge, page.Diagnostics.TotalRequestCharge)
}

func TestEmulatorDiagnosticsConcurrentSnapshots(t *testing.T) {
	container := emulatorContainer(t)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	trackEmulatorItem(t, container, pk, id)
	body, err := json.Marshal(map[string]string{"id": id, "pk": id})
	require.NoError(t, err)
	_, err = container.CreateItem(t.Context(), pk, id, body, nil)
	require.NoError(t, err)

	const callers = 8
	responses := make([]ItemResponse, callers)
	errs := make([]error, callers)
	var wg sync.WaitGroup
	for i := range responses {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			responses[i], errs[i] = container.ReadItem(context.Background(), pk, id, nil)
		}(i)
	}
	wg.Wait()
	for i, response := range responses {
		require.NoError(t, errs[i])
		requireEmulatorDiagnostics(t, response.Diagnostics, http.StatusOK, 0, false)
		require.Equal(t, response.AttemptCount, response.Diagnostics.AttemptCount)
		if i > 0 {
			require.NotSame(t, responses[0].Diagnostics, response.Diagnostics)
		}
	}
}

func TestEmulatorDiagnosticsSurviveClientClose(t *testing.T) {
	client, databaseID, containerID := emulatorClient(t)
	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	body, err := json.Marshal(map[string]string{"id": id, "pk": id})
	require.NoError(t, err)
	_, err = container.CreateItem(t.Context(), pk, id, body, nil)
	require.NoError(t, err)
	read, err := container.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	missingID := uniqueItemID(t) + "-missing"
	_, err = container.ReadItem(t.Context(), NewPartitionKeyString(missingID), missingID, nil)
	missing := requireWireErrorDetails(t, err, CodeNotFound, http.StatusNotFound)
	_, err = container.DeleteItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	require.NoError(t, client.Close())

	requireEmulatorDiagnostics(t, read.Diagnostics, http.StatusOK, 0, false)
	requireEmulatorDiagnostics(t, missing.Diagnostics, http.StatusNotFound, 0, true)
}
