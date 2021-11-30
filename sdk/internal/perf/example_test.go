//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package perf_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type perfTest struct {
	counter  int
	setup    bool
	teardown bool
}

func (p *perfTest) Setup()    { p.setup = true }
func (p *perfTest) Run()      { p.counter++ }
func (p *perfTest) TearDown() { p.teardown = true }

func ExampleGlobalSetup() {
	perf.RunPerfTest(&perfTest{counter: 0})
}
