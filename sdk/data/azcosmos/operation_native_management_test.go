// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Pins the request's unset max-item-count sentinel for every management kind, the same hazard
// TestOperationRequestUsesTheDriversUnsetSentinels pins for item kinds: the driver rejects a zero
// hint outright, so a composite literal that leaves this zero would fail every operation.
func TestManagementRequestUsesTheDriversUnsetSentinel(t *testing.T) {
	for _, kind := range []operationKind{
		operationKindCreateDatabase,
		operationKindReadDatabase,
		operationKindDeleteDatabase,
		operationKindCreateContainer,
		operationKindReadContainer,
		operationKindDeleteContainer,
	} {
		request, release := inspectNativeManagementRequest(kind, nil, nil, nil, nil)
		t.Cleanup(release)
		require.Negative(t, request.maxItemCount,
			"the driver reads < 0 as unset and rejects 0, so a zero here fails the operation")
	}
}

func TestManagementRequestCarriesBody(t *testing.T) {
	body := []byte(`{"id":"db"}`)
	request, release := inspectNativeManagementRequest(operationKindCreateDatabase, nil, nil, nil, body)
	t.Cleanup(release)
	require.Equal(t, int32(operationKindCreateDatabase), request.kind)
	require.Equal(t, body, request.body)
}

func TestManagementRequestOmitsBodyWhenEmpty(t *testing.T) {
	request, release := inspectNativeManagementRequest(operationKindReadDatabase, nil, nil, nil, nil)
	t.Cleanup(release)
	require.Nil(t, request.body)
}

// Nil scope handles reach the native struct as nil rather than some non-zero default: the
// account/database/container fields are plain pass-through, and executeManagement is the one
// deciding which the kind needs, not this builder.
func TestManagementRequestLeavesUnsetScopeHandlesNil(t *testing.T) {
	request, release := inspectNativeManagementRequest(operationKindCreateContainer, nil, nil, nil, nil)
	t.Cleanup(release)
	require.False(t, request.hasAccount)
	require.False(t, request.hasDatabase)
	require.False(t, request.hasContainer)
}
