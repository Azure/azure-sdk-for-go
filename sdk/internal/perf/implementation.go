// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"strings"
	"sync"
)

var (
	// debug is true if --debug is specified
	debug bool
	// duration is the -d/--duration flag
	duration int
	// testProxyURLs is the -x/--test-proxy flag, a semi-colon separated list
	testProxyURLs string
	// warmUpDuration is the -w/--warmup flag
	warmUpDuration int
	// parallelInstances is the -p/--parallel flag
	parallelInstances int

	// wg is used to keep track of the number of goroutines created
	wg sync.WaitGroup

	// number of processes to use, the --maxprocs flag
	numProcesses int
)

// parseProxyURLs splits the --test-proxy input with the delimiter ';'
func parseProxyURLS() []string {
	if testProxyURLs == "" {
		return nil
	}

	testProxyURLs = strings.TrimSuffix(testProxyURLs, ";")

	return strings.Split(testProxyURLs, ";")
}
