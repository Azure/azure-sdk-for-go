// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"testing"

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
	Duration = 1
	WarmUp = 0
	Parallel = 1

	p := NewPerfTestSample(nil)

	err := runGlobalCleanup(p)
	require.NoError(t, err)

	err = runSetup(p)
	require.NoError(t, err)

	err = runCleanup(p)
	require.NoError(t, err)
	err = runGlobalCleanup(p)
	require.NoError(t, err)
}
