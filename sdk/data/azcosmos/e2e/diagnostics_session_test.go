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
// an attempt count, and a rendered JSON summary; the 404 additionally carries no sub-status,
// matching a plain not-found rather than a session-unavailable 404/1002.
//
// Rust also asserts RegionsContacted names the account's East US region. The pinned emulator's
// RegionsContacted population for a single-region deployment has not been separately verified
// here, so this port checks only that RegionsContacted is non-empty rather than asserting an
// exact region name, to avoid over-claiming an untested contract.
func TestDiagnosticsCoverSuccessAndError(t *testing.T) {
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

	_, err = fx.Container.ReadItem(ctx, pk, uniqueID(t), nil)
	cosmosErr := requireCode(t, err, azcosmos.CodeNotFound)
	require.Zero(t, cosmosErr.SubStatus)
	requireCriticalDiagnostics(t, cosmosErr.Diagnostics, 404)
	require.NotEmpty(t, cosmosErr.Diagnostics.RegionsContacted)
}

// TestSessionTokenExplicitManagement ports Rust scenario consistency.session-management: a
// write's own response token, passed explicitly on a subsequent read's SessionToken option,
// is honored by that read.
//
// Rust's version additionally disables the client's automatic session-token capture, to prove
// the explicit token is a genuine substitute rather than something still captured behind the
// scenes. azcosmos/v2 has no public option to disable automatic session capture (see
// [OperationOptions] and [ClientOptions]), so this port verifies only that an explicitly
// supplied token is honored, not that it substitutes for capture being disabled; that narrower
// claim is the honest one the current public API supports.
func TestSessionTokenExplicitManagement(t *testing.T) {
	fx := newFixture(t, "consistency.session-management", "")
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
