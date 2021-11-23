//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package perf

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func int64Ptr(i int64) *int64 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func TestPerf(t *testing.T) {
	duration = int64Ptr(2)
	warmUp = int64Ptr(1)
	iterations = int64Ptr(3)
	runPerf = boolPtr(true)

	setupRan := false
	teardownRan := false
	GlobalSetup(t, func() error {
		setupRan = true
		return nil
	})

	start := time.Now()
	counter := 0
	RunFunc(t, func() {
		counter += 1
	})

	require.Greater(t, time.Since(start), time.Second*7)
	require.Greater(t, counter, 1)
	require.True(t, setupRan)

	GlobalTeardown(t, func() error {
		teardownRan = true
		return nil
	})
	require.True(t, teardownRan)
}

func TestRandomStream(t *testing.T) {
	r := NewRandomStream(1024)
	buffer := make([]byte, 500)

	n, err := r.Read(buffer)
	require.Equal(t, 500, n)
	require.NotEqual(t, make([]byte, 500), buffer)
	require.NoError(t, err)

	n, err = r.Read(buffer)
	require.Equal(t, 500, n)
	require.NotEqual(t, make([]byte, 500), buffer)
	require.NoError(t, err)

	n, err = r.Read(buffer)
	require.Equal(t, 24, n)
	require.NotEqual(t, make([]byte, 500), buffer)
	require.Error(t, err)

	r.Seek(-100, io.SeekEnd)
	fmt.Println(r.position)
	newBuffer := make([]byte, 101)
	n, err = r.Read(newBuffer)
	require.Error(t, err)
	require.Equal(t, 100, n)
	require.NotEqual(t, make([]byte, 101), newBuffer)
	require.Equal(t, newBuffer[len(newBuffer)-1], make([]byte, 1)[0])

	r.Seek(100, io.SeekStart)
	n, err = r.Read(buffer)
	require.Equal(t, 500, n)
	require.NoError(t, err)
	require.NoError(t, r.Close())
}
