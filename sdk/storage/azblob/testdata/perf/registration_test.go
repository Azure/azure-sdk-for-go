// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"flag"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func withTestFlagSet(t *testing.T) *flag.FlagSet {
	t.Helper()
	saved := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(t.Name(), flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	t.Cleanup(func() { flag.CommandLine = saved })
	return flag.CommandLine
}

func TestUploadTestRegister(t *testing.T) {
	flags := withTestFlagSet(t)
	uploadTestRegister()
	require.NoError(t, flags.Parse([]string{"--size=42"}))
	require.Equal(t, 42, uploadTestOpts.size)
}

func TestDownloadTestRegister(t *testing.T) {
	flags := withTestFlagSet(t)
	downloadTestRegister()
	require.NoError(t, flags.Parse([]string{"--size=43"}))
	require.Equal(t, 43, downloadTestOpts.size)
}

func TestListTestRegister(t *testing.T) {
	flags := withTestFlagSet(t)
	listTestRegister()
	require.NoError(t, flags.Parse([]string{"--num-blobs=44", "--num-blobs-parallelism=3"}))
	require.Equal(t, listTestOptions{count: 44, parallelism: 3}, listTestOpts)
}
