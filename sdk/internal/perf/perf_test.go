// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type globalNoOpPerfTest struct {
	PerfTestOptions
}

func NewNoOpTest(ctx context.Context, options PerfTestOptions) (GlobalPerfTest, error) {
	return &globalNoOpPerfTest{
		PerfTestOptions: options,
	}, nil
}

func (g *globalNoOpPerfTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

type noOpPerTest struct {
	*PerfTestOptions
}

func (g *globalNoOpPerfTest) NewPerfTest(ctx context.Context, options *PerfTestOptions) (PerfTest, error) {
	return &noOpPerTest{options}, nil
}

func (n *noOpPerTest) Run(ctx context.Context) error {
	return nil
}

func (n *noOpPerTest) Cleanup(ctx context.Context) error {
	return nil
}

func TestRun(t *testing.T) {
	duration = 2
	warmUpDuration = 0
	parallelInstances = 1

	err := runPerfTest("Sleep", NewNoOpTest)
	require.NoError(t, err)
}

func TestParseProxyURLs(t *testing.T) {
	testProxyURLs = ""
	result := parseProxyURLS()
	require.Nil(t, result)

	testProxyURLs = "https://localhost:5001"
	result = parseProxyURLS()
	require.Equal(t, []string{"https://localhost:5001"}, result)

	testProxyURLs = "https://localhost:5001;https://abc;"
	result = parseProxyURLS()
	require.Equal(t, []string{"https://localhost:5001", "https://abc"}, result)

	testProxyURLs = "https://localhost:5001;https://abc;https://def"
	result = parseProxyURLS()
	require.Equal(t, []string{"https://localhost:5001", "https://abc", "https://def"}, result)
}
