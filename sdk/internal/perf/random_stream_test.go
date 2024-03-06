// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRandomBytes(t *testing.T) {
	a, err := getRandomBytes(1024)
	require.NoError(t, err)
	b, err := getRandomBytes(1024)
	require.NoError(t, err)
	require.NotEqual(t, a, b)
}

func TestRandomStream(t *testing.T) {
	r, err := NewRandomStream(1024)
	require.NoError(t, err)

	a := make([]byte, 500)
	n, err := r.Read(a)
	require.NoError(t, err)
	require.Equal(t, 500, n)
	require.NotEqual(t, a, make([]byte, 500))

	b := make([]byte, 500)
	n, err = r.Read(b)
	require.NoError(t, err)
	require.Equal(t, 500, n)

	c := make([]byte, 500)
	n, err = r.Read(c)
	require.NoError(t, err)
	require.Equal(t, 24, n)

	require.NoError(t, r.Close())

	// Read after finishing should return io.EOF
	n, err = r.Read(c)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)

	// SeekStart should set the position to 0
	pos, err := r.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)
	require.Equal(t, r.(*randomStream).position, int64(0))
	require.Equal(t, r.(*randomStream).remaining, 1024)

	a1 := make([]byte, 500)
	n, err = r.Read(a1)
	require.NoError(t, err)
	require.Equal(t, 500, n)
	require.Equal(t, a, a1)
}
