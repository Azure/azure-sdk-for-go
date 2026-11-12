// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"testing"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

// TestDiagnosticsCoverSuccessAndError ports Rust scenario diagnostics.success-and-error: a
// successful read and a 404 read both carry diagnostics identifying the attempted status code,
// an attempt count, a rendered JSON summary, and an activity ID, on both outcomes; the 404
// additionally carries no sub-status, matching a plain not-found rather than a
// session-unavailable 404/1002.
//
// The catalog also requires identifying "the logical operation" on both outcomes. azcosmos/v2
// exposes no operation-type/name field on Response, Error, or Diagnostics, so that part of the
// contract cannot be certified through the current public API; this port covers status,
// activity, attempts, and region only.
//
// Rust also asserts RegionsContacted names the account's East US region. The pinned emulator's
// RegionsContacted population for a single-region deployment has not been separately verified
// here, so this port checks only that RegionsContacted is non-empty rather than asserting an
// exact region name, to avoid over-claiming an untested contract.
func TestDiagnosticsCoverSuccessAndError(t *testing.T) {
	t.Run("diagnostics.success-and-error", func(t *testing.T) {
		fx := newFixture(t, "diagnostics.success-and-error", "")
		ctx := context.Background()
		id := uniqueID(t)
		pk := azcosmos.NewPartitionKeyString("A")

		_, err := fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 1)), nil)
		require.NoError(t, err)
		trackItem(t, fx.Container, pk, id)

		success, err := fx.Container.ReadItem(ctx, pk, id, nil)
		require.NoError(t, err)
		require.Equal(t, 200, success.StatusCode)
		requireCriticalDiagnostics(t, success.Diagnostics, 200)
		require.NotEmpty(t, success.Diagnostics.RegionsContacted)
		require.NotEmptyf(t, success.ActivityID, "successful operation must carry an activity id")

		_, err = fx.Container.ReadItem(ctx, pk, uniqueID(t), nil)
		cosmosErr := requireCode(t, err, azcosmos.CodeNotFound)
		require.Zero(t, cosmosErr.SubStatus)
		requireCriticalDiagnostics(t, cosmosErr.Diagnostics, 404)
		require.NotEmpty(t, cosmosErr.Diagnostics.RegionsContacted)
		require.NotEmptyf(t, cosmosErr.ActivityID, "failed operation must carry an activity id")
	})
}

// TestSessionTokenExplicitManagement is a supplementary regression test for consistency.session-
// management: a write's own response token, passed explicitly on a subsequent read's SessionToken
// option, is honored by that read. It is not registered as certification evidence for that
// scenario (see internal/e2e.Implementations) because the catalog also requires disabling
// automatic session-token capture and exercising continuation-based queries alongside the
// explicit token, and azcosmos/v2 has no public option to disable automatic session capture (see
// [OperationOptions] and [ClientOptions]); only the narrower "explicit token honored" claim is
// tested here.
//
// This also uses newRawFixture, not newFixture: the scenario's only declared profile,
// configurationResilience, needs a two-region, replication-delayed account this fixture's single
// pre-provisioned emulator account cannot provide, so newFixture's profile validation would
// correctly reject it.
func TestSessionTokenExplicitManagement(t *testing.T) {
	fx := newRawFixture(t, "")
	ctx := context.Background()
	id := uniqueID(t)
	pk := azcosmos.NewPartitionKeyString("A")

	created, err := fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 1)), &azcosmos.CreateItemOptions{})
	require.NoError(t, err)
	trackItem(t, fx.Container, pk, id)
	require.NotEmptyf(t, string(created.SessionToken), "write response must expose a session token")

	explicit, err := fx.Container.ReadItem(ctx, pk, id, &azcosmos.ReadItemOptions{
		SessionToken: created.SessionToken,
		Operation:    azcosmos.OperationOptions{ConsistencyStrategy: azcosmos.ReadConsistencyStrategySession},
	})
	require.NoError(t, err)
	require.Equal(t, item(id, "A", 1), mustUnmarshalItem(t, explicit.Value))
}
