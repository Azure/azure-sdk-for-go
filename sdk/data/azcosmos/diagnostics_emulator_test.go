// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmulatorDiagnosticsSnapshotSurvivesClientClose(t *testing.T) {
	endpoint, database, containerID := emulatorConfiguration(t)
	runtime, err := NewRuntime(nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	observed := make(chan OperationDiagnostic, 2)
	client, err := NewClientWithKey(endpoint, KeyCredential{accountKey: emulatorKey}, &ClientOptions{
		Runtime: runtime,
		Operation: OperationOptions{
			ThrottlingRetry:       ThrottlingRetryOptions{MaxRetryCount: to(uint32(0))},
			HedgingEnabled:        to(false),
			MaxFailoverRetryCount: to(uint32(0)),
		},
		DiagnosticsHandler: func(_ context.Context, diagnostic OperationDiagnostic) { observed <- diagnostic },
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer(database, containerID)
	require.NoError(t, err)
	_, err = container.ReadItem(t.Context(), NewPartitionKeyString("missing"), uniqueItemID(t), nil)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, 404, cosmosErr.StatusCode)
	diagnostic := <-observed
	require.Equal(t, "ReadItem", diagnostic.Operation)
	require.Error(t, diagnostic.Error)
	require.True(t, diagnostic.Diagnostics.Available())
	require.NoError(t, diagnostic.Diagnostics.Err())
	require.True(t, json.Valid([]byte(diagnostic.Diagnostics.JSON())))
	require.Positive(t, diagnostic.Diagnostics.RequestCount())
	jsonSnapshot := diagnostic.Diagnostics.JSON()
	_, err = container.ReadItem(t.Context(), NewPartitionKeyString("missing"), uniqueItemID(t), nil)
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, 404, cosmosErr.StatusCode)
	<-observed
	require.NoError(t, client.Close())
	require.Equal(t, jsonSnapshot, diagnostic.Diagnostics.JSON(), "snapshot must survive native completion and client release")
}

func TestEmulatorDiagnosticsHandlerCanCloseClient(t *testing.T) {
	endpoint, database, containerID := emulatorConfiguration(t)
	var client *Client
	var err error
	client, err = NewClientWithKey(endpoint, KeyCredential{accountKey: emulatorKey}, &ClientOptions{
		DiagnosticsHandler: func(_ context.Context, _ OperationDiagnostic) {
			require.NoError(t, client.Close())
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer(database, containerID)
	require.NoError(t, err)
	_, err = container.ReadItem(t.Context(), NewPartitionKeyString("missing"), uniqueItemID(t), nil)
	require.Error(t, err)
	_, err = client.acquire()
	var closed *Error
	require.ErrorAs(t, err, &closed)
	require.Equal(t, CodeClientClosed, closed.Code)
}
