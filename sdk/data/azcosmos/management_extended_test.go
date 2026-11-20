// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestManagementRejectsNilContextsAndClosedClients(t *testing.T) {
	client := newTestClient(t)
	database, err := client.NewDatabase("db")
	require.NoError(t, err)
	container := localTestContainer(client, "db", "items")
	properties := ContainerProperties{ID: "items", PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk"}}}
	operations := []func(context.Context) error{
		func(ctx context.Context) error {
			_, err := client.CreateDatabase(ctx, DatabaseProperties{ID: "db"}, nil)
			return err
		},
		func(ctx context.Context) error { _, err := database.Read(ctx, nil); return err },
		func(ctx context.Context) error { _, err := database.Delete(ctx, nil); return err },
		func(ctx context.Context) error { _, err := database.CreateContainer(ctx, properties, nil); return err },
		func(ctx context.Context) error { _, err := container.Read(ctx, nil); return err },
		func(ctx context.Context) error { _, err := container.Replace(ctx, properties, nil); return err },
		func(ctx context.Context) error { _, err := container.Delete(ctx, nil); return err },
		func(ctx context.Context) error { _, err := database.ReadThroughput(ctx, nil); return err },
		func(ctx context.Context) error { _, err := container.ReadThroughput(ctx, nil); return err },
		func(ctx context.Context) error {
			_, err := database.BeginReplaceThroughput(ctx, ThroughputProperties{ManualThroughput: to(uint64(400))}, nil)
			return err
		},
		func(ctx context.Context) error {
			_, err := container.BeginReplaceThroughput(ctx, ThroughputProperties{ManualThroughput: to(uint64(400))}, nil)
			return err
		},
	}
	for _, operation := range operations {
		require.ErrorContains(t, operation(nil), "context must not be nil")
	}
	require.NoError(t, client.Close())
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, operation := range operations {
		var cosmosErr *Error
		require.ErrorAs(t, operation(ctx), &cosmosErr)
		require.Equal(t, CodeClientClosed, cosmosErr.Code)
	}
}

func TestManagementCreateThroughputHeadersAndOwnership(t *testing.T) {
	operation := OperationOptions{CustomHeaders: map[string]string{"x-test": "preserved"}}
	manual, err := managementCreateOptions(operation, &ThroughputProperties{ManualThroughput: to(uint64(400))})
	require.NoError(t, err)
	require.Equal(t, map[string]string{"x-test": "preserved"}, manual.CustomHeaders)
	require.Equal(t, map[string]string{"x-test": "preserved"}, operation.CustomHeaders)
	autoscale, err := managementCreateOptions(operation, &ThroughputProperties{Autoscale: &AutoscaleSettings{
		MaxThroughput: 4000, AutoUpgradePolicy: &AutoscaleAutoUpgradePolicy{ThroughputPolicy: &AutoscaleThroughputPolicy{IncrementPercent: 10}},
	}})
	require.NoError(t, err)
	require.Equal(t, map[string]string{"x-test": "preserved"}, autoscale.CustomHeaders)
	inherited, err := managementCreateOptions(OperationOptions{}, &ThroughputProperties{ManualThroughput: to(uint64(400))})
	require.NoError(t, err)
	require.Nil(t, inherited.CustomHeaders, "typed throughput must not replace inherited headers")
	_, err = managementCreateOptions(OperationOptions{CustomHeaders: map[string]string{"X-MS-Offer-Throughput": "999"}}, &ThroughputProperties{ManualThroughput: to(uint64(400))})
	require.ErrorContains(t, err, "conflict")
	for _, properties := range []ThroughputProperties{{}, {ManualThroughput: to(uint64(0))}, {Autoscale: &AutoscaleSettings{}}, {ManualThroughput: to(uint64(400)), Autoscale: &AutoscaleSettings{MaxThroughput: 4000}}} {
		_, err := managementCreateOptions(OperationOptions{}, &properties)
		require.Error(t, err)
	}
}

func TestManagementPropertiesRetainSupportedAndUnknownConfiguration(t *testing.T) {
	body := `{
		"id":"items","partitionKey":{"paths":["/pk"],"kind":"Hash","version":2,"systemKey":true},
		"indexingPolicy":{"automatic":true,"includedPaths":[{"path":"/*","indexes":[{"kind":"Range","dataType":"String","precision":-1}]}],"spatialIndexes":[{"path":"/location/*","types":["Point"],"boundingBox":{"xmin":-180,"ymin":-90,"xmax":180,"ymax":90}}],"vectorIndexes":[{"path":"/vector","type":"diskANN","quantizationByteSize":8,"indexingSearchListSize":100,"vectorIndexShardKey":["/tenant"]}],"fullTextIndexes":[{"path":"/text"}]},
		"conflictResolutionPolicy":{"mode":"LastWriterWins","conflictResolutionPath":"/_ts","conflictResolutionProcedure":""},
		"vectorEmbeddingPolicy":{"vectorEmbeddings":[{"path":"/vector","dataType":"float32","dimensions":3,"distanceFunction":"cosine","futureField":7}],"futurePolicy":true},
		"fullTextPolicy":{"defaultLanguage":"en-US","fullTextPaths":[{"path":"/text","language":"en-US","futurePath":true}],"futurePolicy":"keep"},
		"changeFeedPolicy":{"retentionDuration":5},"analyticalStorageTtl":-1,"defaultTtl":10,
		"futureContainer":{"enabled":true},"_ts":1735689600
	}`
	var properties ContainerProperties
	require.NoError(t, json.Unmarshal([]byte(body), &properties))
	properties.DefaultTimeToLive = to(int32(30))
	encoded, err := json.Marshal(properties)
	require.NoError(t, err)
	var actual, expected map[string]any
	require.NoError(t, json.Unmarshal(encoded, &actual))
	require.NoError(t, json.Unmarshal([]byte(body), &expected))
	delete(expected, "_ts")
	expected["defaultTtl"] = float64(30)
	require.Equal(t, expected, actual)
	require.Equal(t, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), properties.LastModified)
}

func TestManagementPagerCleanupAndValidation(t *testing.T) {
	client := newTestClient(t)
	for _, options := range []*ManagementQueryOptions{{PageSizeHint: -1}, {ContinuationToken: "bad\x00token"}} {
		pager := client.NewQueryDatabasesPager(NewQuery("SELECT * FROM c"), options)
		_, err := pager.NextPage(t.Context())
		require.Error(t, err)
		require.NoError(t, pager.Close())
	}
	cursor := &testQueryCursor{}
	pager := &ManagementPager[json.RawMessage]{client: client, cursor: cursor}
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Empty(t, page.Items)
	require.True(t, pager.More())
	token, err := pager.ContinuationToken(t.Context())
	require.NoError(t, err)
	require.Equal(t, "snapshot", token)
	page, err = pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Len(t, page.Items, 2)
	_, err = pager.NextPage(t.Context())
	require.NoError(t, err)
	require.False(t, pager.More())
	require.Equal(t, 1, cursor.closed)
	require.NoError(t, pager.Close())
	require.Equal(t, 1, cursor.closed)
	_, err = pager.NextPage(t.Context())
	require.ErrorContains(t, err, "exhausted")

	cursor = &testQueryCursor{}
	invalidModels := &ManagementPager[DatabaseProperties]{client: client, cursor: cursor}
	_, err = invalidModels.NextPage(t.Context())
	require.NoError(t, err)
	result, err := invalidModels.NextPage(t.Context())
	require.Equal(t, ManagementQueryResponse[DatabaseProperties]{}, result)
	require.Error(t, err)
	require.Equal(t, 1, cursor.closed)
	require.False(t, invalidModels.More())
}

func TestManagementPagerPreservesPageSessionToken(t *testing.T) {
	cursor := &modelQueryCursor{page: queryPage{
		SessionToken: "0:42",
		Items:        [][]byte{[]byte(`{"id":"db"}`)},
	}}
	pager := &ManagementPager[DatabaseProperties]{client: newTestClient(t), cursor: cursor}
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Equal(t, SessionToken("0:42"), page.SessionToken)
	require.Equal(t, "db", page.Items[0].ID)
	require.NoError(t, pager.Close())
}

func TestThroughputPollerTerminalFailureAndCancellation(t *testing.T) {
	calls := 0
	poller := &ThroughputPoller{result: ThroughputResponse{Response: Response{RequestCharge: 4}}, poll: func(context.Context) (managementResult, error) {
		calls++
		switch calls {
		case 1:
			return managementResult{}, &Error{Code: CodeOperationCancelled, RequestCharge: 1, cause: context.Canceled}
		case 2:
			return managementResult{response: ItemResponse{Response: Response{RequestCharge: 2}}, body: []byte(`{"content":{"offerThroughput":400}}`), pending: true}, nil
		default:
			return managementResult{}, &Error{Code: CodeForbidden, StatusCode: 403, RequestCharge: 1.5, ActivityID: "failed"}
		}
	}}
	_, err := poller.Result()
	require.ErrorContains(t, err, "pending")
	_, err = poller.Poll(t.Context())
	require.ErrorIs(t, err, context.Canceled)
	require.False(t, poller.Done())
	result, err := poller.Poll(t.Context())
	require.NoError(t, err)
	require.Equal(t, float64(7), result.RequestCharge)
	require.False(t, poller.Done())
	result, err = poller.Poll(t.Context())
	require.Equal(t, ThroughputResponse{}, result)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeForbidden, cosmosErr.Code)
	require.Equal(t, 8.5, cosmosErr.RequestCharge)
	require.True(t, poller.Done())
	_, again := poller.Poll(t.Context())
	require.Same(t, err, again)
	require.Equal(t, 3, calls)
	_, again = poller.Result()
	require.Same(t, err, again)
}

func TestThroughputPollerPollUntilDoneAndDecodeFailure(t *testing.T) {
	poller := &ThroughputPoller{poll: func(context.Context) (managementResult, error) {
		return managementResult{body: []byte(`{"content":{"offerAutopilotSettings":{"maxThroughput":8000}}}`)}, nil
	}}
	result, err := poller.PollUntilDone(t.Context(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	require.NoError(t, err)
	require.Equal(t, uint64(8000), result.ThroughputProperties.Autoscale.MaxThroughput)
	require.True(t, poller.Done())
	poller = &ThroughputPoller{poll: func(context.Context) (managementResult, error) { return managementResult{body: []byte(`{}`)}, nil }}
	_, err = poller.Poll(t.Context())
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
	require.True(t, poller.Done())
	_, err = poller.Result()
	require.True(t, errors.As(err, &cosmosErr))
}
