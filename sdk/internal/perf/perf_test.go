// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	globalCleanup = false
	newnooptest   = false
	newperftest   = false
	run           = false
	cleanup       = false
)

type globalNoOpPerfTest struct {
	PerfTestOptions
}

func NewNoOpTest(ctx context.Context, options PerfTestOptions) (GlobalPerfTest, error) {
	newnooptest = true
	return &globalNoOpPerfTest{
		PerfTestOptions: options,
	}, nil
}

func (g *globalNoOpPerfTest) GlobalCleanup(ctx context.Context) error {
	globalCleanup = true
	return nil
}

type noOpPerTest struct {
	*PerfTestOptions
}

func (g *globalNoOpPerfTest) NewPerfTest(ctx context.Context, options *PerfTestOptions) (PerfTest, error) {
	newperftest = true
	return &noOpPerTest{options}, nil
}

func (n *noOpPerTest) Run(ctx context.Context) error {
	run = true
	return nil
}

func (n *noOpPerTest) Cleanup(ctx context.Context) error {
	cleanup = true
	return nil
}

func TestRun(t *testing.T) {
	duration = 2
	warmUpDuration = 0
	parallelInstances = 1

	runner := newPerfRunner(PerfMethods{New: NewNoOpTest, Register: nil}, "NoOpTest")
	err := runner.Run()
	require.NoError(t, err)
}

func TestParseProxyURLs(t *testing.T) {
	testProxyURLs = ""
	result := parseProxyURLS()
	require.Equal(t, 0, len(result))

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
