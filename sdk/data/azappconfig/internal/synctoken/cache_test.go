//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package synctoken

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	stk := NewCache()
	require.Zero(t, stk.Get())

	require.Error(t, stk.Set(""))
	require.Error(t, stk.Set("  \t"))
	require.Error(t, stk.Set("id=val"))
	require.Zero(t, stk.Get())
	require.Error(t, stk.Set("id=val;"))
	require.Zero(t, stk.Get())
	require.Error(t, stk.Set(";sn=1"))
	require.Zero(t, stk.Get())

	require.NoError(t, stk.Set("id=val1;sn=1"))
	f := stk.Get()
	require.EqualValues(t, "id=val1", f)

	require.NoError(t, stk.Set("id=val2;sn=2"))
	f = stk.Get()
	require.EqualValues(t, "id=val2", f)

	require.NoError(t, stk.Set("id2=some;sn=2"))
	f = stk.Get()
	require.EqualValues(t, "id=val2,id2=some", f)
}
