// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type noOpPerfTest struct {
	perf.PerfTestOptions
}

func (n *noOpPerfTest) GlobalSetup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) Run(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) Cleanup(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) RegisterArguments(ctx context.Context) error {
	return nil
}

func (n *noOpPerfTest) GetMetadata() perf.PerfTestOptions {
	return n.PerfTestOptions
}

func NewNoOpTest(options *perf.PerfTestOptions) perf.PerfTest {
	if options == nil {
		options = &perf.PerfTestOptions{}
	}
	options.Name = "NoOpTest"
	return &noOpPerfTest{
		PerfTestOptions: *options,
	}
}
