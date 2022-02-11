// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type perfTestSample struct {
	PerfTestOptions
}

func (p *perfTestSample) GetMetadata() PerfTestOptions {
	return p.PerfTestOptions
}

func (p *perfTestSample) Setup(ctx context.Context) error {
	return nil
}
func (p *perfTestSample) GlobalSetup(ctx context.Context) error {
	return nil
}
func (p *perfTestSample) Run(ctx context.Context) error {
	time.Sleep(time.Second)
	return nil
}
func (p *perfTestSample) Cleanup(ctx context.Context) error {
	return nil
}
func (p *perfTestSample) GlobalCleanup(ctx context.Context) error {
	return nil
}

func NewPerfTestSample(options *PerfTestOptions) PerfTest {
	return &perfTestSample{PerfTestOptions: PerfTestOptions{Name: "sample"}}
}

func TestRun(t *testing.T) {
	duration = 2
	warmUpDuration = 0
	parallelInstances = 1

	err := runPerfTest(NewPerfTestSample)
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