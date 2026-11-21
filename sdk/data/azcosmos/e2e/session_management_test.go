// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

func setReplication(t *testing.T, managementEndpoint, action string) {
	t.Helper()
	endpoint, err := url.JoinPath(managementEndpoint, "regions", "West US", "replication", action)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	defer func() { require.NoError(t, response.Body.Close()) }()
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestSessionManagement(t *testing.T) {
	profile := selectedProfile(t, "consistency.session-management")
	endpoint, managementEndpoint, databaseID, containerID := requireProfileEmulator(t, profile)
	require.NoError(t, validateEmulatorBackend(os.Getenv("AZURE_COSMOS_E2E_BACKEND")))
	probeEmulatorCapabilities(t, profile, endpoint, managementEndpoint)
	options := &azcosmos.ClientOptions{
		Routing:   azcosmos.PreferredRegions(azcosmos.RegionWestUS, azcosmos.RegionEastUS),
		Operation: azcosmos.OperationOptions{SessionCapturingDisabled: to(true)},
	}
	fx := newFixtureAt(t, endpoint, databaseID, containerID, "", options)
	setReplication(t, managementEndpoint, "pause")
	t.Cleanup(func() { setReplication(t, managementEndpoint, "resume") })

	ctx := context.Background()
	const count = 6
	pkValue := uniqueID(t)
	pk := azcosmos.NewPartitionKeyString(pkValue)
	var session azcosmos.SessionToken
	for i := range count {
		id := fmt.Sprintf("%s-%d", pkValue, i)
		created, err := fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, item(id, pkValue, int64(i))), nil)
		require.NoError(t, err)
		trackItem(t, fx.Container, pk, id)
		session = created.SessionToken
		require.NotEmpty(t, session)
	}
	id := fmt.Sprintf("%s-%d", pkValue, count-1)
	_, err := fx.Container.ReadItem(ctx, pk, id, &azcosmos.ReadItemOptions{
		Operation: azcosmos.OperationOptions{ConsistencyStrategy: azcosmos.ReadConsistencyStrategySession},
	})
	stale := requireCode(t, err, azcosmos.CodeNotFound)
	require.NotNil(t, stale.Diagnostics)
	require.NotEmpty(t, stale.Diagnostics.RegionsContacted)
	require.Equal(t, string(azcosmos.RegionWestUS), stale.Diagnostics.RegionsContacted[0])

	read, err := fx.Container.ReadItem(ctx, pk, id, &azcosmos.ReadItemOptions{
		SessionToken: session,
		Operation:    azcosmos.OperationOptions{ConsistencyStrategy: azcosmos.ReadConsistencyStrategySession},
	})
	require.NoError(t, err)
	require.Equal(t, item(id, pkValue, count-1), mustUnmarshalItem(t, read.Value))

	query, err := azcosmos.NewQuery("SELECT * FROM c WHERE c.pk = @pk ORDER BY c.value ASC").WithParameter("@pk", pkValue)
	require.NoError(t, err)
	scope := azcosmos.NewFeedScopeForPartitionKey(pk)
	control := fx.Container.NewQueryItemsPager[Item](query, scope, &azcosmos.QueryOptions{
		Operation: azcosmos.OperationOptions{ConsistencyStrategy: azcosmos.ReadConsistencyStrategySession},
		Feed:      azcosmos.FeedOptions{PageSizeHint: 1},
	})
	controlPages := 0
	for control.More() {
		page, err := control.NextPage(ctx)
		require.NoError(t, err)
		controlPages++
		require.Empty(t, page.Items, "the stale replica must not expose the writes without a token")
	}
	require.Positive(t, controlPages)
	require.NoError(t, control.Close())
	queryOptions := azcosmos.QueryOptions{
		SessionToken: session,
		Operation:    azcosmos.OperationOptions{ConsistencyStrategy: azcosmos.ReadConsistencyStrategySession},
		Feed:         azcosmos.FeedOptions{PageSizeHint: 1},
	}
	pager := fx.Container.NewQueryItemsPager[Item](query, scope, &queryOptions)
	first, err := pager.NextPage(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, first.Items)
	require.Less(t, len(first.Items), count)
	continuation, err := pager.ContinuationToken(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, continuation)
	require.NoError(t, pager.Close())
	queryOptions.Feed.ContinuationToken = continuation
	resumed := fx.Container.NewQueryItemsPager[Item](query, scope, &queryOptions)
	defer func() { require.NoError(t, resumed.Close()) }()
	items := append([]Item(nil), first.Items...)
	for resumed.More() {
		page, err := resumed.NextPage(ctx)
		require.NoError(t, err)
		items = append(items, page.Items...)
	}
	require.Len(t, items, count)
	for i, got := range items {
		require.Equal(t, item(fmt.Sprintf("%s-%d", pkValue, i), pkValue, int64(i)), got)
	}
}
