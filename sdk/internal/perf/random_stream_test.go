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
	baseData, err := getRandomBytes(1024)
	require.NoError(t, err)
	r := &randomStream{
		offset:   0,
		baseData: baseData,
	}
	require.NoError(t, err)

	a := make([]byte, 500)
	n, err := r.Read(a)
	require.NoError(t, err)
	require.Equal(t, 500, n)

	b := make([]byte, 500)
	n, err = r.Read(b)
	require.NoError(t, err)
	require.Equal(t, 500, n)

	c := make([]byte, 500)
	n, err = r.Read(c)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 24, n)

	require.NotEqual(t, a, b)
	require.NotEqual(t, a, c)
	require.NotEqual(t, c, b)
}
