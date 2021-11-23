//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package perf_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

func ExampleGlobalSetup() {
	func(t *testing.T) {
		// 1. Global Set Up
		perf.GlobalSetup(t, func() error {
			// Some setup method here for creating resources
			return nil
		})

		defer perf.GlobalTeardown(t, func() error {
			// Some resource clean up methods here
			return nil
		})

		perf.RunFunc(t, func() {
			// Methods you are trying to judge for performance here
		})

	}(&testing.T{})
}
