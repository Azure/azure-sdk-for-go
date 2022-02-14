// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type globalNoOpPerfTest struct {
	perf.PerfTestOptions
}

func NewNoOpTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	return &globalNoOpPerfTest{
		PerfTestOptions: options,
	}, nil
}

func (g *globalNoOpPerfTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

type noOpPerTest struct {
	*perf.PerfTestOptions
}

func (g *globalNoOpPerfTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	return &noOpPerTest{options}, nil
}

func (n *noOpPerTest) Run(ctx context.Context) error {
	return nil
}

func (n *noOpPerTest) Cleanup(ctx context.Context) error {
	return nil
}
