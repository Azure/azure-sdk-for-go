// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"strings"
	"sync"
)

var (
	debug             bool
	duration          int
	testProxyURLs     string
	warmUpDuration    int
	parallelInstances int
	wg                sync.WaitGroup
	numProcesses      int
)

// parse the TestProxy input into a slice of strings
func parseProxyURLS() []string {
	var ret []string
	if testProxyURLs == "" {
		return ret
	}

	testProxyURLs = strings.TrimSuffix(testProxyURLs, ";")

	ret = strings.Split(testProxyURLs, ";")

	return ret
}
