//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSyncPool(t *testing.T) {
	sp, err := NewSyncPool(_1MiB, 4)
	require.NoError(t, err)
	require.NotNil(t, sp)

	buf := sp.Get()
	require.Len(t, buf, _1MiB)

	const _2MiB = _1MiB * 2
	sp.Put(make([]byte, _2MiB))

	buf = sp.Get()

	// buf is either 1 or 2 MB depending on what sync.Pool decides to do
	if l := len(buf); l != _2MiB && l != _1MiB {
		t.Fatalf("unexpected length %d", l)
	}
}